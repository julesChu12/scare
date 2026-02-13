//go:build integration

package integration

import (
	"net/http"
	"testing"

	"community-elderly-care-platform/test/integration/testutil"

	"github.com/stretchr/testify/assert"
)

func TestStation(t *testing.T) {
	env := testutil.Setup(t)

	t.Run("Admin获取站点列表", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/b/stations?page=1&page_size=10", testutil.AdminToken())
		data := testutil.AssertPageResponse(t, w)
		items := data["items"].([]interface{})
		assert.GreaterOrEqual(t, len(items), 1)
	})

	t.Run("Admin获取站点详情", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/b/stations/1", testutil.AdminToken())
		data := testutil.AssertOK(t, w)
		assert.Equal(t, "霍营街道养老服务站", data["name"])
	})

	t.Run("Admin获取不存在的站点", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/b/stations/9999", testutil.AdminToken())
		assert.Contains(t, []int{http.StatusNotFound, http.StatusBadRequest}, w.Code)
	})

	t.Run("Admin创建站点", func(t *testing.T) {
		body := `{
			"name": "回龙观养老服务站",
			"code": "HLG001",
			"latitude": 40.07,
			"longitude": 116.33,
			"status": "active"
		}`
		w := testutil.DoRequest(env.Engine, http.MethodPost, "/api/v1/b/stations", testutil.AdminToken(), body)
		assert.Contains(t, []int{http.StatusOK, http.StatusCreated}, w.Code)
	})

	t.Run("Admin创建站点_缺少名称", func(t *testing.T) {
		body := `{"code": "NONAME", "status": "active"}`
		w := testutil.DoRequest(env.Engine, http.MethodPost, "/api/v1/b/stations", testutil.AdminToken(), body)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Admin更新站点", func(t *testing.T) {
		body := `{"name": "霍营街道养老服务站(更新)", "code": "HY001", "status": "active"}`
		w := testutil.DoRequest(env.Engine, http.MethodPut, "/api/v1/b/stations/1", testutil.AdminToken(), body)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Admin删除站点", func(t *testing.T) {
		// 先创建一个用于删除的站点
		body := `{"name": "待删除站点", "code": "DEL001", "status": "active"}`
		testutil.DoRequest(env.Engine, http.MethodPost, "/api/v1/b/stations", testutil.AdminToken(), body)

		// 删除 ID=2 的站点（刚创建的）
		w := testutil.DoRequest(env.Engine, http.MethodDelete, "/api/v1/b/stations/2", testutil.AdminToken())
		assert.Contains(t, []int{http.StatusOK, http.StatusNoContent}, w.Code)
	})

	t.Run("Admin删除不存在的站点", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodDelete, "/api/v1/b/stations/9999", testutil.AdminToken())
		assert.Contains(t, []int{http.StatusOK, http.StatusNotFound}, w.Code)
	})

	t.Run("Staff无权限创建站点", func(t *testing.T) {
		body := `{"name": "无权限站点", "code": "NOPERM", "status": "active"}`
		w := testutil.DoRequest(env.Engine, http.MethodPost, "/api/v1/b/stations", testutil.StaffToken(), body)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}
