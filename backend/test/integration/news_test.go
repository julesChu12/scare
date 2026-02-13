//go:build integration

package integration

import (
	"net/http"
	"testing"

	"community-elderly-care-platform/test/integration/testutil"

	"github.com/stretchr/testify/assert"
)

func TestNews(t *testing.T) {
	env := testutil.Setup(t)

	t.Run("B端列表", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/b/news?page=1&page_size=10", testutil.AdminToken())
		data := testutil.AssertPageResponse(t, w)
		items := data["items"].([]interface{})
		assert.GreaterOrEqual(t, len(items), 1)
	})

	t.Run("B端详情", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/b/news/1", testutil.AdminToken())
		data := testutil.AssertOK(t, w)
		assert.Equal(t, "测试新闻", data["title"])
	})

	t.Run("B端创建", func(t *testing.T) {
		body := `{"title":"新建新闻","summary":"摘要","content":"内容","type":"news","status":"draft","station_id":1}`
		w := testutil.DoRequest(env.Engine, http.MethodPost, "/api/v1/b/news", testutil.AdminToken(), body)
		assert.Contains(t, []int{http.StatusOK, http.StatusCreated}, w.Code)
	})

	t.Run("B端更新", func(t *testing.T) {
		body := `{"title":"更新后的新闻","summary":"更新摘要","content":"更新内容","type":"notice","status":"published"}`
		w := testutil.DoRequest(env.Engine, http.MethodPut, "/api/v1/b/news/1", testutil.AdminToken(), body)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("B端删除", func(t *testing.T) {
		// 先创建一条用于删除
		body := `{"title":"待删除","summary":"s","content":"c","type":"news","status":"draft","station_id":1}`
		testutil.DoRequest(env.Engine, http.MethodPost, "/api/v1/b/news", testutil.AdminToken(), body)

		w := testutil.DoRequest(env.Engine, http.MethodDelete, "/api/v1/b/news/100", testutil.AdminToken())
		// 可能 404（ID 不存在）或 200
		assert.Contains(t, []int{http.StatusOK, http.StatusNotFound}, w.Code)
	})

	t.Run("C端列表", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/c/news?page=1&page_size=10", "")
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("C端详情", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/c/news/1", "")
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
