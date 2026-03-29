//go:build integration

package integration

import (
	"net/http"
	"testing"

	"community-elderly-care-platform/test/integration/testutil"

	"github.com/stretchr/testify/assert"
)

func TestMenu(t *testing.T) {
	env := testutil.Setup(t)

	t.Run("菜单树", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/b/menus", testutil.AdminToken())
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("菜单详情", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/b/menus/1", testutil.AdminToken())
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("用户菜单", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/b/menus/user", testutil.AdminToken())
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("创建菜单", func(t *testing.T) {
		body := `{"parent_id":0,"name":"测试菜单","path":"/test","icon":"test","sort":99,"hidden":false,"status":"active"}`
		w := testutil.DoRequest(env.Engine, http.MethodPost, "/api/v1/b/menus", testutil.AdminToken(), body)
		assert.Contains(t, []int{http.StatusOK, http.StatusCreated}, w.Code)
	})

	t.Run("更新菜单", func(t *testing.T) {
		body := `{"name":"更新菜单","path":"/updated","sort":1,"hidden":false,"status":"active"}`
		w := testutil.DoRequest(env.Engine, http.MethodPut, "/api/v1/b/menus/1", testutil.AdminToken(), body)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("批量排序", func(t *testing.T) {
		body := `{"updates":[{"id":1,"sort":2},{"id":2,"sort":1}]}`
		w := testutil.DoRequest(env.Engine, http.MethodPut, "/api/v1/b/menus/sort", testutil.AdminToken(), body)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("删除菜单", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodDelete, "/api/v1/b/menus/100", testutil.AdminToken())
		assert.Contains(t, []int{http.StatusOK, http.StatusNotFound}, w.Code)
	})
}
