//go:build integration

package integration

import (
	"net/http"
	"testing"

	"community-elderly-care-platform/test/integration/testutil"

	"github.com/stretchr/testify/assert"
)

func TestReports(t *testing.T) {
	env := testutil.Setup(t)

	t.Run("生成报表-服务统计Excel", func(t *testing.T) {
		body := `{
			"type": "service",
			"format": "xlsx",
			"start_date": "2026-01-01",
			"end_date": "2026-01-31"
		}`
		w := testutil.DoRequest(env.Engine, http.MethodPost, "/api/v1/b/reports/generate", testutil.AdminToken(), body)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/octet-stream", w.Header().Get("Content-Type"))
		assert.Contains(t, w.Header().Get("Content-Disposition"), "attachment;")
	})

	t.Run("生成报表-人员绩效Excel", func(t *testing.T) {
		body := `{
			"type": "performance",
			"format": "xlsx",
			"start_date": "2026-01-01",
			"end_date": "2026-01-31"
		}`
		w := testutil.DoRequest(env.Engine, http.MethodPost, "/api/v1/b/reports/generate", testutil.AdminToken(), body)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/octet-stream", w.Header().Get("Content-Type"))
		assert.Contains(t, w.Header().Get("Content-Disposition"), "attachment;")
	})

	t.Run("生成报表-需求分析Excel", func(t *testing.T) {
		body := `{
			"type": "request",
			"format": "xlsx",
			"start_date": "2026-01-01",
			"end_date": "2026-01-31"
		}`
		w := testutil.DoRequest(env.Engine, http.MethodPost, "/api/v1/b/reports/generate", testutil.AdminToken(), body)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/octet-stream", w.Header().Get("Content-Type"))
		assert.Contains(t, w.Header().Get("Content-Disposition"), "attachment;")
	})

	t.Run("生成报表-站点运营Excel", func(t *testing.T) {
		body := `{
			"type": "station",
			"format": "xlsx",
			"start_date": "2026-01-01",
			"end_date": "2026-01-31"
		}`
		w := testutil.DoRequest(env.Engine, http.MethodPost, "/api/v1/b/reports/generate", testutil.AdminToken(), body)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/octet-stream", w.Header().Get("Content-Type"))
		assert.Contains(t, w.Header().Get("Content-Disposition"), "attachment;")
	})

	t.Run("生成报表-CSV格式", func(t *testing.T) {
		body := `{
				"type": "service",
			"format": "csv",
			"start_date": "2026-01-01",
			"end_date": "2026-01-31"
		}`
		w := testutil.DoRequest(env.Engine, http.MethodPost, "/api/v1/b/reports/generate", testutil.AdminToken(), body)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/octet-stream", w.Header().Get("Content-Type"))
		assert.Contains(t, w.Header().Get("Content-Disposition"), "attachment;")
	})

	t.Run("预览报表-服务统计", func(t *testing.T) {
		body := `{
				"type": "service",
				"format": "xlsx",
				"start_date": "2026-01-01",
				"end_date": "2026-01-31",
				"preview": true
			}`
		w := testutil.DoRequest(env.Engine, http.MethodPost, "/api/v1/b/reports/generate", testutil.AdminToken(), body)
		assert.Equal(t, http.StatusOK, w.Code)
		testutil.AssertOK(t, w)
		assert.Contains(t, w.Body.String(), "request_count")
	})

	t.Run("生成报表-参数错误-结束日期早于开始日期", func(t *testing.T) {
		body := `{
			"type": "service",
			"format": "xlsx",
			"start_date": "2026-02-01",
			"end_date": "2026-01-31"
		}`
		w := testutil.DoRequest(env.Engine, http.MethodPost, "/api/v1/b/reports/generate", testutil.AdminToken(), body)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("生成报表-参数错误-日期格式错误", func(t *testing.T) {
		body := `{
			"type": "service",
			"format": "xlsx",
			"start_date": "2026/13/01",
			"end_date": "2026-01-31"
		}`
		w := testutil.DoRequest(env.Engine, http.MethodPost, "/api/v1/b/reports/generate", testutil.AdminToken(), body)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("生成报表-参数错误-无效类型", func(t *testing.T) {
		body := `{
			"type": "invalid_type",
			"format": "xlsx",
			"start_date": "2026-01-01",
			"end_date": "2026-01-31"
		}`
		w := testutil.DoRequest(env.Engine, http.MethodPost, "/api/v1/b/reports/generate", testutil.AdminToken(), body)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("生成报表-参数错误-无效格式", func(t *testing.T) {
		body := `{
			"type": "service",
			"format": "invalid_format",
			"start_date": "2026-01-01",
			"end_date": "2026-01-31"
		}`
		w := testutil.DoRequest(env.Engine, http.MethodPost, "/api/v1/b/reports/generate", testutil.AdminToken(), body)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("获取历史报表列表", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/b/reports?page=1&page_size=10", testutil.AdminToken())
		assert.Equal(t, http.StatusOK, w.Code)
		testutil.AssertOK(t, w)
	})

	t.Run("获取历史报表列表-筛选类型", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/b/reports?type=service", testutil.AdminToken())
		assert.Equal(t, http.StatusOK, w.Code)
		testutil.AssertOK(t, w)
	})

	t.Run("获取历史报表列表-分页", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/b/reports?page=2&page_size=5", testutil.AdminToken())
		assert.Equal(t, http.StatusOK, w.Code)
		testutil.AssertOK(t, w)
	})

	t.Run("下载历史报表", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/b/reports/1/download", testutil.AdminToken())
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Header().Get("Content-Disposition"), "attachment")
	})

	t.Run("下载历史报表-报表不存在", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/b/reports/99999/download", testutil.AdminToken())
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("删除报表", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodDelete, "/api/v1/b/reports/1", testutil.AdminToken())
		assert.Equal(t, http.StatusOK, w.Code)
		testutil.AssertOK(t, w)
	})

	t.Run("删除报表-无权限", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodDelete, "/api/v1/b/reports/2", testutil.StationManagerToken())
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("StationManager访问统计概览", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/b/statistics/overview?days=7", testutil.StationManagerToken())
		assert.Equal(t, http.StatusOK, w.Code)
		testutil.AssertOK(t, w)
	})

	t.Run("获取服务类型分布", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/b/statistics/service-types?days=7", testutil.AdminToken())
		assert.Equal(t, http.StatusOK, w.Code)
		testutil.AssertOK(t, w)
	})

	t.Run("获取需求趋势", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/b/statistics/trend?days=7", testutil.AdminToken())
		assert.Equal(t, http.StatusOK, w.Code)
		testutil.AssertOK(t, w)
	})

	t.Run("获取处理效率", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/b/statistics/efficiency?days=7", testutil.AdminToken())
		assert.Equal(t, http.StatusOK, w.Code)
		testutil.AssertOK(t, w)
	})

	t.Run("获取服务人员排行", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/b/statistics/staff-ranking?days=7&limit=10", testutil.AdminToken())
		assert.Equal(t, http.StatusOK, w.Code)
		testutil.AssertOK(t, w)
	})
}
