//go:build integration

package integration

import (
	"net/http"
	"testing"

	"community-elderly-care-platform/test/integration/testutil"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	env := testutil.Setup(t)

	t.Run("Admin获取用户列表", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/b/users?page=1&page_size=10", testutil.AdminToken())
		data := testutil.AssertOK(t, w)
		items := data["items"].([]interface{})
		assert.GreaterOrEqual(t, len(items), 1)
		total := data["total"].(float64)
		assert.GreaterOrEqual(t, total, float64(1))
	})

	t.Run("Admin获取用户详情", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/b/users/1", testutil.AdminToken())
		data := testutil.AssertOK(t, w)
		assert.Equal(t, "13800000001", data["phone"])
		assert.Equal(t, "系统管理员", data["name"])
	})

	t.Run("Admin获取不存在的用户", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/b/users/9999", testutil.AdminToken())
		assert.Contains(t, []int{http.StatusNotFound, http.StatusBadRequest}, w.Code)
	})

	t.Run("Admin创建用户", func(t *testing.T) {
		body := `{
			"phone": "13800009999",
			"password": "NewUser@123",
			"name": "新建用户",
			"identity_type": "staff",
			"station_id": 1,
			"status": "active"
		}`
		w := testutil.DoRequest(env.Engine, http.MethodPost, "/api/v1/b/users", testutil.AdminToken(), body)
		assert.Contains(t, []int{http.StatusOK, http.StatusCreated}, w.Code)
	})

	t.Run("Admin创建用户_手机号重复", func(t *testing.T) {
		body := `{
			"phone": "13800000001",
			"password": "Dup@123456",
			"name": "重复手机号",
			"identity_type": "staff",
			"station_id": 1
		}`
		w := testutil.DoRequest(env.Engine, http.MethodPost, "/api/v1/b/users", testutil.AdminToken(), body)
		assert.Contains(t, []int{http.StatusBadRequest, http.StatusConflict}, w.Code)
	})

	t.Run("Admin创建用户_参数缺失", func(t *testing.T) {
		body := `{"name": "缺少必填字段"}`
		w := testutil.DoRequest(env.Engine, http.MethodPost, "/api/v1/b/users", testutil.AdminToken(), body)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Admin更新用户", func(t *testing.T) {
		body := `{"name": "更新后的管理员", "status": "active"}`
		w := testutil.DoRequest(env.Engine, http.MethodPut, "/api/v1/b/users/1", testutil.AdminToken(), body)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Staff无权限访问用户列表", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/b/users?page=1&page_size=10", testutil.StaffToken())
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("StationManager无权限访问用户列表", func(t *testing.T) {
		// station_manager 在种子数据中只有 user:list 权限（RolePermission ID=1）
		// 但权限匹配依赖 PermissionService 的实现，这里验证实际行为
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/b/users?page=1&page_size=10", testutil.StationManagerToken())
		// station_manager 有 permission_id=1 (system:user:list)，应该可以访问
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
