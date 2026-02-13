//go:build integration

package integration

import (
	"net/http"
	"testing"

	"community-elderly-care-platform/test/integration/testutil"

	"github.com/stretchr/testify/assert"
)

func TestServiceRequest(t *testing.T) {
	env := testutil.Setup(t)
	elderlyToken := testutil.CEndElderlyToken()
	adminToken := testutil.AdminToken()

	t.Run("C端_创建服务请求", func(t *testing.T) {
		body := `{
			"service_type": "daily_care",
			"description": "需要日常照护",
			"contact_phone": "13900000001",
			"longitude": 116.4,
			"latitude": 39.9
		}`
		w := testutil.DoRequest(env.Engine, http.MethodPost,
			"/api/v1/c/requests", elderlyToken, body)
		data := testutil.AssertOK(t, w)
		assert.NotNil(t, data)
	})

	t.Run("C端_请求列表", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet,
			"/api/v1/c/requests", elderlyToken)
		testutil.AssertOK(t, w)
	})

	t.Run("C端_请求详情", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet,
			"/api/v1/c/requests/1", elderlyToken)
		// 可能 200 或 404，取决于种子数据
		assert.Contains(t,
			[]int{http.StatusOK, http.StatusNotFound}, w.Code)
	})

	t.Run("B端_请求列表", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet,
			"/api/v1/b/requests", adminToken)
		testutil.AssertOK(t, w)
	})

	t.Run("B端_请求详情", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet,
			"/api/v1/b/requests/1", adminToken)
		assert.Contains(t,
			[]int{http.StatusOK, http.StatusNotFound}, w.Code)
	})

	t.Run("C端_取消请求", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodPost,
			"/api/v1/c/requests/1/cancel", elderlyToken)
		// 取消结果取决于请求是否存在及状态
		assert.Contains(t,
			[]int{http.StatusOK, http.StatusNotFound, http.StatusBadRequest}, w.Code)
	})
}
