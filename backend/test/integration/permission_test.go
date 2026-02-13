//go:build integration

package integration

import (
	"net/http"
	"testing"

	"community-elderly-care-platform/test/integration/testutil"

	"github.com/stretchr/testify/assert"
)

func TestPermission(t *testing.T) {
	env := testutil.Setup(t)

	t.Run("权限树", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/b/permissions/tree", testutil.AdminToken())
		testutil.AssertOK(t, w)
	})

	t.Run("角色权限查询", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/b/roles/admin/permissions", testutil.AdminToken())
		testutil.AssertOK(t, w)
	})

	t.Run("角色权限查询-station_manager", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/b/roles/station_manager/permissions", testutil.AdminToken())
		testutil.AssertOK(t, w)
	})

	t.Run("更新角色权限", func(t *testing.T) {
		body := `{"permissions":["service:request:list","service:task:pool"]}`
		w := testutil.DoRequest(env.Engine, http.MethodPut, "/api/v1/b/roles/staff/permissions", testutil.AdminToken(), body)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Staff无权访问权限树", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/b/permissions/tree", testutil.StaffToken())
		testutil.AssertError(t, w, http.StatusForbidden, "forbidden")
	})

	t.Run("更新用户身份", func(t *testing.T) {
		body := `{"identities":[{"identity_type":"staff","station_id":1,"is_primary":true}]}`
		w := testutil.DoRequest(env.Engine, http.MethodPut, "/api/v1/b/users/4/identities", testutil.AdminToken(), body)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
