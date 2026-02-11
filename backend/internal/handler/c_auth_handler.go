package handler

import (
	"net/http"

	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/repository"
	"community-elderly-care-platform/internal/service"

	"github.com/gin-gonic/gin"
)

// CAuthHandler C端认证处理器
type CAuthHandler struct {
	authService  *service.AuthService
	userRepo     *repository.UserRepository
	customerRepo *repository.CustomerRepository
	smsService   *service.SMSService
}

func NewCAuthHandler(authService *service.AuthService, userRepo *repository.UserRepository, customerRepo *repository.CustomerRepository, smsService *service.SMSService) *CAuthHandler {
	return &CAuthHandler{
		authService:  authService,
		userRepo:     userRepo,
		customerRepo: customerRepo,
		smsService:   smsService,
	}
}

type cLoginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password"`
	Code     string `json:"code"`
	Type     string `json:"type"` // "password" or "code", default "code"
}

type cRefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Login C端登录接口（支持密码和验证码两种方式）
// @Summary      C端用户登录
// @Description  支持密码登录和验证码登录两种方式
// @Tags         c_auth
// @Accept       json
// @Produce      json
// @Param        request body cLoginRequest true "登录请求"
// @Success      200  {object} APIResponse "登录成功"
// @Failure      400  {object} APIResponse "请求参数错误"
// @Failure      401  {object} APIResponse "认证失败"
// @Failure      404  {object} APIResponse "用户不存在"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /c/auth/login [post]
func (h *CAuthHandler) Login(c *gin.Context) {
	var req cLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "invalid payload")
		return
	}

	// 默认使用验证码登录
	loginType := req.Type
	if loginType == "" {
		loginType = "code"
	}

	var tokens *service.Tokens
	var user *model.User
	var err error

	if loginType == "password" {
		// 密码登录
		if req.Password == "" {
			RespondError(c, http.StatusBadRequest, "password required")
			return
		}
		tokens, user, err = h.authService.LoginCEnd(req.Phone, req.Password)
	} else {
		// 验证码登录
		if req.Code == "" {
			RespondError(c, http.StatusBadRequest, "code required")
			return
		}
		tokens, user, err = h.authService.LoginCEndByCode(req.Phone, req.Code)
	}

	if err != nil {
		if err == service.ErrInvalidCredentials {
			RespondError(c, http.StatusUnauthorized, "invalid credentials")
			return
		}
		if err == service.ErrCodeInvalid {
			RespondError(c, http.StatusUnauthorized, "验证码错误或已过期")
			return
		}
		if err == service.ErrUserInactive {
			RespondError(c, http.StatusForbidden, "user inactive")
			return
		}
		if err == service.ErrUserNotFound {
			RespondError(c, http.StatusNotFound, "用户不存在")
			return
		}
		if err == service.ErrNoCustomerProfile {
			RespondError(c, http.StatusForbidden, "user has no customer profile")
			return
		}
		RespondError(c, http.StatusInternalServerError, "login failed")
		return
	}

	// 获取客户档案信息
	profile, err := h.customerRepo.GetByUserID(user.ID)
	var customerType string
	if err == nil && profile.CustomerType != "" {
		customerType = profile.CustomerType
	}

	data := gin.H{
		"token":         tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
		"user_id":       user.ID,
		"type":          "c_end",
		"customer_type": customerType,
		"name":          user.Name,
		"phone":         user.Phone,
		"status":        user.Status,
	}
	Respond(c, http.StatusOK, "ok", data)
}

// Refresh C端刷新Token
// @Summary      刷新访问令牌
// @Description  使用刷新令牌获取新的访问令牌
// @Tags         c_auth
// @Accept       json
// @Produce      json
// @Param        request body cRefreshRequest true "刷新令牌请求"
// @Success      200  {object} APIResponse "刷新成功"
// @Failure      401  {object} APIResponse "令牌无效"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /c/auth/refresh [post]
func (h *CAuthHandler) Refresh(c *gin.Context) {
	var req cRefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "invalid payload")
		return
	}

	tokens, err := h.authService.Refresh(req.RefreshToken)
	if err != nil {
		RespondError(c, http.StatusUnauthorized, "invalid token")
		return
	}

	Respond(c, http.StatusOK, "ok", gin.H{
		"token":         tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
	})
}

// Me 获取当前登录用户信息（C端）
// @Summary      获取当前登录用户信息
// @Description  获取当前C端用户的详细信息
// @Tags         c_auth
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200  {object} APIResponse "获取成功"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      404  {object} APIResponse "用户不存在"
// @Router       /c/auth/me [get]
func (h *CAuthHandler) Me(c *gin.Context) {
	userID, ok := GetUserID(c)
	if !ok {
		RespondError(c, http.StatusUnauthorized, "missing user")
		return
	}

	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		RespondError(c, http.StatusNotFound, "user not found")
		return
	}

	// 获取客户档案信息
	profile, err := h.customerRepo.GetByUserID(userID)
	if err != nil {
		RespondError(c, http.StatusNotFound, "customer profile not found")
		return
	}

	var customerType string
	if profile.CustomerType != "" {
		customerType = profile.CustomerType
	}

	Respond(c, http.StatusOK, "ok", gin.H{
		"user_id":       user.ID,
		"type":          "c_end",
		"customer_type": customerType,
		"name":          user.Name,
		"phone":         user.Phone,
		"status":        user.Status,
	})
}

// SendCode 发送短信验证码
// @Summary      发送短信验证码
// @Description  向手机号发送6位数字验证码，有效期5分钟
// @Tags         c_auth
// @Accept       json
// @Produce      json
// @Param        request body object{phone=string} true "手机号"
// @Success      200  {object} APIResponse "发送成功"
// @Failure      400  {object} APIResponse "请求参数错误"
// @Failure      429  {object} APIResponse "发送频率限制"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /c/auth/send-code [post]
func (h *CAuthHandler) SendCode(c *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "invalid payload")
		return
	}

	// 发送验证码
	if err := h.smsService.SendCode(req.Phone); err != nil {
		if err == service.ErrRateLimitMinute {
			RespondError(c, http.StatusTooManyRequests, "发送过于频繁，请1分钟后再试")
			return
		}
		if err == service.ErrRateLimitDaily {
			RespondError(c, http.StatusTooManyRequests, "今日发送次数已达上限")
			return
		}
		RespondError(c, http.StatusInternalServerError, "send code failed")
		return
	}

	Respond(c, http.StatusOK, "验证码已发送", nil)
}

// QuickStart 快速开通（注册+登录+创建服务请求）
// @Summary      快速开通（注册+登录+创建服务请求）
// @Description  一次性完成注册、登录和创建服务需求的全流程
// @Tags         c_auth
// @Accept       json
// @Produce      json
// @Param        request body object{phone=string,code=string,name=string,address=string,latitude=float64,longitude=float64,service_type=string,description=string} true "快速开通请求"
// @Success      200  {object} APIResponse "快速开通成功"
// @Failure      400  {object} APIResponse "请求参数错误"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /c/auth/quick-start [post]
func (h *CAuthHandler) QuickStart(c *gin.Context) {
	var req struct {
		Phone       string   `json:"phone" binding:"required"`
		Code        string   `json:"code" binding:"required"`
		Name        string   `json:"name" binding:"required"`
		Address     string   `json:"address"`
		Latitude    *float64 `json:"latitude"`
		Longitude   *float64 `json:"longitude"`
		ServiceType string   `json:"service_type" binding:"required"`
		Description *string  `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "invalid payload")
		return
	}

	// 地址或经纬度至少需要一个
	if req.Address == "" && (req.Latitude == nil || req.Longitude == nil) {
		RespondError(c, http.StatusBadRequest, "address or location required")
		return
	}

	// 调用快速开通服务
	result, err := h.authService.QuickStart(service.QuickStartInput{
		Phone:       req.Phone,
		Code:        req.Code,
		Name:        req.Name,
		Address:     req.Address,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		ServiceType: req.ServiceType,
		Description: req.Description,
	})

	if err != nil {
		if err == service.ErrCodeInvalid {
			RespondError(c, http.StatusBadRequest, "验证码错误或已过期")
			return
		}
		if err == service.ErrNoStation {
			RespondError(c, http.StatusBadRequest, "无法找到服务站点，请检查地址")
			return
		}
		RespondError(c, http.StatusInternalServerError, "quick start failed")
		return
	}

	// 返回结果
	data := gin.H{
		"token":         result.Token,
		"refresh_token": result.RefreshToken,
		"user": gin.H{
			"id":    result.User.ID,
			"phone": result.User.Phone,
			"role":  "c_end",
		},
		"profile": gin.H{
			"name":    result.User.Name,
			"address": result.Profile.Address,
		},
		"request": gin.H{
			"id":           result.Request.ID,
			"request_no":   result.Request.RequestNo,
			"service_type": result.Request.ServiceType,
			"status":       result.Request.Status,
			"created_at":   result.Request.CreatedAt,
		},
	}

	Respond(c, http.StatusOK, "ok", data)
}

// CheckToken 检查Token状态并返回用户信息（用于预填充）
// @Summary      检查Token状态并返回用户信息
// @Description  验证Token有效性并返回用户信息，用于表单预填充
// @Tags         c_auth
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200  {object} APIResponse "Token有效"
// @Failure      401  {object} APIResponse "Token无效或过期"
// @Failure      404  {object} APIResponse "用户不存在"
// @Router       /c/auth/check [get]
func (h *CAuthHandler) CheckToken(c *gin.Context) {
	userID, ok := GetUserID(c)
	if !ok {
		RespondError(c, http.StatusUnauthorized, "missing user")
		return
	}

	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		RespondError(c, http.StatusNotFound, "user not found")
		return
	}

	// 获取客户档案信息（可能不存在）
	profile, err := h.customerRepo.GetByUserID(userID)
	var profileData *gin.H
	if err == nil {
		profileData = &gin.H{
			"name":      user.Name,
			"id_number": profile.IDCard,
			"address":   profile.Address,
			"user_type": profile.CustomerType,
		}
	}

	Respond(c, http.StatusOK, "ok", gin.H{
		"user": gin.H{
			"id":    user.ID,
			"phone": user.Phone,
			"role":  "c_end",
		},
		"profile": profileData,
	})
}
