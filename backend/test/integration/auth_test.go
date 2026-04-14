//go:build integration

package integration

import (
	"net/http"
	"testing"

	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/pkg/jwt"
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

	t.Run("C端设置密码后可使用密码登录", func(t *testing.T) {
		user := &model.User{
			Phone:  "13900000999",
			Name:   "测试用户",
			Status: "active",
		}
		assert.NoError(t, env.DB.DB.Omit("PasswordHash").Create(user).Error)
		assert.NoError(t, env.DB.DB.Create(&model.CustomerProfile{
			UserID:           user.ID,
			CustomerType:     "elderly",
			EmergencyContact: `{}`,
		}).Error)
		assert.NoError(t, env.DB.DB.Create(&model.UserIdentity{
			UserID:       user.ID,
			IdentityType: "elderly",
			IsPrimary:    true,
			Status:       "active",
		}).Error)

		token, _ := jwt.NewManager("test-jwt-secret-key-for-integration-tests", 24, 168).
			GenerateToken(user.ID, "c_end", 0, []string{"elderly"}, "elderly")

		setBody := `{"current_password":"InitPass@123","new_password":"NewPass@123"}`
		setResp := testutil.DoRequest(env.Engine, http.MethodPost, "/api/v1/c/auth/password", token, setBody)
		assert.Equal(t, http.StatusOK, setResp.Code)

		loginBody := `{"phone":"13900000999","password":"NewPass@123","type":"password"}`
		loginResp := testutil.DoRequest(env.Engine, http.MethodPost, "/api/v1/c/auth/login", "", loginBody)
		assert.Equal(t, http.StatusOK, loginResp.Code)

		data := testutil.AssertOK(t, loginResp)
		assert.Equal(t, "13900000999", data["phone"])
		assert.NotEmpty(t, data["token"])
	})

	t.Run("C端验证码重置密码后可使用新密码登录", func(t *testing.T) {
		resetBody := `{"phone":"13900000001","code":"000000","new_password":"ResetPass@123"}`
		resetResp := testutil.DoRequest(env.Engine, http.MethodPost, "/api/v1/c/auth/reset-password", "", resetBody)
		assert.Equal(t, http.StatusOK, resetResp.Code)

		loginBody := `{"phone":"13900000001","password":"ResetPass@123","type":"password"}`
		loginResp := testutil.DoRequest(env.Engine, http.MethodPost, "/api/v1/c/auth/login", "", loginBody)
		assert.Equal(t, http.StatusOK, loginResp.Code)

		data := testutil.AssertOK(t, loginResp)
		assert.Equal(t, "13900000001", data["phone"])
		assert.NotEmpty(t, data["token"])
	})

	t.Run("无效Token", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/b/auth/me", "invalid.jwt.token")
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
