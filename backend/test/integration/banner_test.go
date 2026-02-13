//go:build integration

package integration

import (
	"net/http"
	"testing"

	"community-elderly-care-platform/test/integration/testutil"

	"github.com/stretchr/testify/assert"
)

func TestBanner(t *testing.T) {
	env := testutil.Setup(t)

	t.Run("B端列表", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/b/banners?page=1&page_size=10", testutil.AdminToken())
		testutil.AssertPageResponse(t, w)
	})

	t.Run("B端创建", func(t *testing.T) {
		body := `{"station_id":0,"title":"新Banner","image_url":"/new.jpg","link_type":"none","sort":2,"status":"active"}`
		w := testutil.DoRequest(env.Engine, http.MethodPost, "/api/v1/b/banners", testutil.AdminToken(), body)
		assert.Contains(t, []int{http.StatusOK, http.StatusCreated}, w.Code)
	})

	t.Run("B端更新", func(t *testing.T) {
		body := `{"title":"更新Banner","image_url":"/updated.jpg","link_type":"url","link_value":"https://example.com","sort":1,"status":"active"}`
		w := testutil.DoRequest(env.Engine, http.MethodPut, "/api/v1/b/banners/1", testutil.AdminToken(), body)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("B端删除", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodDelete, "/api/v1/b/banners/100", testutil.AdminToken())
		assert.Contains(t, []int{http.StatusOK, http.StatusNotFound}, w.Code)
	})

	t.Run("C端列表", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/c/banners", "")
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
