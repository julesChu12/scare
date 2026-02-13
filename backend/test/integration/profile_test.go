//go:build integration

package integration

import (
	"net/http"
	"testing"

	"community-elderly-care-platform/test/integration/testutil"

	"github.com/stretchr/testify/assert"
)

func TestProfile(t *testing.T) {
	env := testutil.Setup(t)

	t.Run("C端获取个人信息", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/c/auth/me", testutil.CEndElderlyToken())
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("C端更新个人资料", func(t *testing.T) {
		body := `{"name":"张大爷更新","gender":"male","address":"北京市昌平区"}`
		w := testutil.DoRequest(env.Engine, http.MethodPut, "/api/v1/c/profile", testutil.CEndElderlyToken(), body)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("C端逆地理编码", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/c/geocode/reverse?lng=116.4&lat=39.9", "")
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("C端站点匹配", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/c/stations/match?lng=116.4&lat=39.9", "")
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("C端Token检查", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/c/auth/check", testutil.CEndElderlyToken())
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Family角色访问", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/c/auth/me", testutil.CEndFamilyToken())
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
