//go:build integration

package integration

import (
	"fmt"
	"net/http"
	"testing"

	"community-elderly-care-platform/test/integration/testutil"

	"github.com/stretchr/testify/assert"
)

func TestServiceRequest(t *testing.T) {
	env := testutil.Setup(t)
	elderlyToken := testutil.CEndElderlyToken()
	adminToken := testutil.AdminToken()
	var requestID int64

	t.Run("C端_创建服务请求", func(t *testing.T) {
		body := `{
			"service_type": "care",
			"description": "需要日常照护",
			"contact_name": "张大爷",
			"contact_phone": "13900000001",
			"lng": 116.4,
			"lat": 39.9
		}`
		w := testutil.DoRequest(env.Engine, http.MethodPost,
			"/api/v1/c/requests", elderlyToken, body)
		data := testutil.AssertOK(t, w)
		assert.NotNil(t, data)
		requestID = int64(data["id"].(float64))
	})

	t.Run("C端_请求列表", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet,
			"/api/v1/c/requests", elderlyToken)
		testutil.AssertOK(t, w)
	})

	t.Run("C端_请求详情", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet,
			"/api/v1/c/requests/"+itoa(requestID), elderlyToken)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("B端_请求列表", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet,
			"/api/v1/b/requests", adminToken)
		testutil.AssertOK(t, w)
	})

	t.Run("B端_请求详情", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet,
			"/api/v1/b/requests/"+itoa(requestID), adminToken)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("C端_取消请求", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodPost,
			"/api/v1/c/requests/"+itoa(requestID)+"/cancel", elderlyToken)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func itoa(id int64) string {
	return fmt.Sprintf("%d", id)
}
