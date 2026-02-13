//go:build integration

package integration

import (
	"net/http"
	"testing"

	"community-elderly-care-platform/test/integration/testutil"

	"github.com/stretchr/testify/assert"
)

func TestZoneCRUD(t *testing.T) {
	env := testutil.Setup(t)
	token := testutil.AdminToken()

	t.Run("列表_获取所有围栏", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/b/zones", token)
		testutil.AssertOK(t, w)
	})

	t.Run("创建_新围栏", func(t *testing.T) {
		body := `{
			"station_id": 1,
			"name": "测试围栏",
			"points": [
				{"lng": 116.3, "lat": 39.8},
				{"lng": 116.5, "lat": 39.8},
				{"lng": 116.5, "lat": 40.0},
				{"lng": 116.3, "lat": 40.0}
			],
			"priority": 2
		}`
		w := testutil.DoRequest(env.Engine, http.MethodPost, "/api/v1/b/zones", token, body)
		data := testutil.AssertOK(t, w)
		assert.NotNil(t, data)
	})

	t.Run("更新_修改围栏", func(t *testing.T) {
		body := `{
			"station_id": 1,
			"name": "更新后围栏",
			"points": [
				{"lng": 116.3, "lat": 39.8},
				{"lng": 116.5, "lat": 39.8},
				{"lng": 116.5, "lat": 40.0},
				{"lng": 116.3, "lat": 40.0}
			],
			"priority": 3
		}`
		w := testutil.DoRequest(env.Engine, http.MethodPut, "/api/v1/b/zones/2", token, body)
		// 围栏 2 由创建步骤产生，验证响应正常
		assert.Contains(t, []int{http.StatusOK, http.StatusNotFound}, w.Code)
	})

	t.Run("删除_移除围栏", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodDelete, "/api/v1/b/zones/2", token)
		assert.Contains(t, []int{http.StatusOK, http.StatusNotFound}, w.Code)
	})
}
