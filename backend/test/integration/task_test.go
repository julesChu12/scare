//go:build integration

package integration

import (
	"net/http"
	"testing"

	"community-elderly-care-platform/test/integration/testutil"

	"github.com/stretchr/testify/assert"
)

func TestTask(t *testing.T) {
	env := testutil.Setup(t)
	adminToken := testutil.AdminToken()
	staffToken := testutil.StaffToken()

	t.Run("B端_任务列表_Admin", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet,
			"/api/v1/b/tasks", adminToken)
		testutil.AssertOK(t, w)
	})

	t.Run("B端_任务池_Admin", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet,
			"/api/v1/b/tasks/pool", adminToken)
		testutil.AssertOK(t, w)
	})

	t.Run("B端_我的任务_Staff", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet,
			"/api/v1/b/tasks/my", staffToken)
		testutil.AssertOK(t, w)
	})

	t.Run("B端_任务详情_Admin", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet,
			"/api/v1/b/tasks/1", adminToken)
		// 种子数据中无任务记录，404 是预期行为
		assert.Contains(t,
			[]int{http.StatusOK, http.StatusNotFound}, w.Code)
	})

	t.Run("Staff_可访问任务池_公共权限", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet,
			"/api/v1/b/tasks/pool", staffToken)
		// 任务池是公共权限，Staff 应当可以访问
		assert.Equal(t, http.StatusOK, w.Code)
		testutil.AssertOK(t, w)
	})
}
