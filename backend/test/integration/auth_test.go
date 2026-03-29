//go:build integration

package integration

import (
	"net/http"
	"testing"

	"community-elderly-care-platform/test/integration/testutil"

	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	env := testutil.Setup(t)

	t.Run("B端登录_密码错误", func(t *testing.T) {
		body := `{"phone":"13800000001","password":"WrongPassword"}`
		w := testutil.DoRequest(env.Engine, http.MethodPost, "/api/v1/b/auth/login", "", body)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("B端登录_手机号不存在", func(t *testing.T) {
		body := `{"phone":"19999999999","password":"Test@123"}`
		w := testutil.DoRequest(env.Engine, http.MethodPost, "/api/v1/b/auth/login", "", body)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("B端登录_参数缺失", func(t *testing.T) {
		body := `{"phone":"13800000001"}`
		w := testutil.DoRequest(env.Engine, http.MethodPost, "/api/v1/b/auth/login", "", body)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("B端获取当前用户信息", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/b/auth/me", testutil.AdminToken())
		data := testutil.AssertOK(t, w)
		user, ok := data["user"].(map[string]interface{})
		assert.True(t, ok)
		assert.NotNil(t, user["id"])
		assert.Equal(t, "13800000001", user["phone"])
	})

	t.Run("B端登出", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodPost, "/api/v1/b/auth/logout", testutil.AdminToken())
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("B端刷新Token_缺少refresh_token", func(t *testing.T) {
		body := `{}`
		w := testutil.DoRequest(env.Engine, http.MethodPost, "/api/v1/b/auth/refresh", "", body)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("B端刷新Token_无效refresh_token", func(t *testing.T) {
		body := `{"refresh_token":"invalid-token"}`
		w := testutil.DoRequest(env.Engine, http.MethodPost, "/api/v1/b/auth/refresh", "", body)
		assert.Contains(t, []int{http.StatusUnauthorized, http.StatusBadRequest}, w.Code)
	})

	t.Run("端类型隔离_B端Token访问C端", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/c/auth/me", testutil.AdminToken())
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("端类型隔离_C端Token访问B端", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/b/auth/me", testutil.CEndElderlyToken())
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("未认证访问B端受保护接口", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/b/auth/me", "")
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("未认证访问C端受保护接口", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/c/auth/me", "")
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("无效Token", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/b/auth/me", "invalid.jwt.token")
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
