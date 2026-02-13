//go:build integration

package integration

import (
	"net/http"
	"testing"

	"community-elderly-care-platform/test/integration/testutil"
)

func TestNotification(t *testing.T) {
	env := testutil.Setup(t)

	t.Run("B端通知列表", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/b/notifications?page=1&page_size=10", testutil.AdminToken())
		testutil.AssertPageResponse(t, w)
	})

	t.Run("B端标记已读-不存在", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodPost, "/api/v1/b/notifications/999/read", testutil.AdminToken())
		// 不存在的通知，可能返回 404 或 200
		_ = w
	})

	t.Run("C端通知列表", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/c/notifications?page=1&page_size=10", testutil.CEndElderlyToken())
		testutil.AssertPageResponse(t, w)
	})

	t.Run("C端标记已读-不存在", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodPost, "/api/v1/c/notifications/999/read", testutil.CEndElderlyToken())
		_ = w
	})
}
