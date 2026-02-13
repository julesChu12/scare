//go:build integration

package integration

import (
	"net/http"
	"testing"

	"community-elderly-care-platform/test/integration/testutil"

	"github.com/stretchr/testify/assert"
)

func TestStatistics(t *testing.T) {
	env := testutil.Setup(t)

	t.Run("仪表盘统计", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/b/statistics/dashboard", testutil.AdminToken())
		assert.Equal(t, http.StatusOK, w.Code)
		testutil.AssertOK(t, w)
	})

	t.Run("任务统计", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/b/statistics/tasks", testutil.AdminToken())
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("请求统计", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/b/statistics/requests", testutil.AdminToken())
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("今日统计", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/b/statistics/today", testutil.AdminToken())
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("StationManager可访问统计", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/b/statistics/dashboard", testutil.StationManagerToken())
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
