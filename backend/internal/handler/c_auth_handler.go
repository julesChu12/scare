package handler

import (
	"net/http"

	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/service"

	"github.com/gin-gonic/gin"
)

// CAuthHandler C端认证处理器
type CAuthHandler struct {
	authService    *service.AuthService
	accountService *service.CEndAccountService
	smsService     *service.SMSService
}

func NewCAuthHandler(authService *service.AuthService, accountService *service.CEndAccountService, smsService *service.SMSService) *CAuthHandler {
	return &CAuthHandler{
		authService:    authService,
		accountService: accountService,
		smsService:     smsService,
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

func buildCEndUserPayload(user *model.User) gin.H {
	return gin.H{
		"id":           user.ID,
		"phone":        user.Phone,
		"role":         "c_end",
		"has_password": service.HasPasswordHash(user.PasswordHash),
	}
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
		if err == service.ErrPasswordNotSet {
			RespondError(c, http.StatusBadRequest, "用户未设置密码，请使用验证码登录")
			return
		}
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

	accountInfo, err := h.accountService.GetAccountInfo(user.ID)
	if err != nil {
		if err == service.ErrNoCustomerProfile {
			RespondError(c, http.StatusForbidden, "user has no customer profile")
			return
		}
		if err == service.ErrUserNotFound {
			RespondError(c, http.StatusNotFound, "用户不存在")
			return
		}
		RespondError(c, http.StatusInternalServerError, "load current user failed")
		return
	}

	data := gin.H{
		"token":         tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
		"user_id":       accountInfo.UserID,
		"type":          accountInfo.Type,
		"customer_type": accountInfo.CustomerType,
		"name":          accountInfo.Name,
		"phone":         accountInfo.Phone,
		"status":        accountInfo.Status,
		"has_password":  accountInfo.HasPassword,
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

// Register 用户注册
// @Summary      用户注册
// @Description  创建新账号，完成手机号验证和密码设置
// @Tags         c_auth
// @Accept       json
// @Produce      json
// @Param        request body object{phone=string,code=string,password=string,name=string} true "注册信息"
// @Success      200  {object} APIResponse "注册成功"
// @Failure      400  {object} APIResponse "请求参数错误"
// @Failure      409  {object} APIResponse "手机号已注册"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /c/auth/register [post]
func (h *CAuthHandler) Register(c *gin.Context) {
	var req struct {
		Phone    string `json:"phone" binding:"required"`
		Code     string `json:"code" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
		Name     string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "invalid payload")
		return
	}

	result, err := h.authService.Register(service.RegisterInput{
		Phone:    req.Phone,
		Code:     req.Code,
		Password: req.Password,
		Name:     req.Name,
	})

	if err != nil {
		if err.Error() == "该手机号已注册，请直接登录" {
			RespondError(c, http.StatusConflict, err.Error())
			return
		}
		if err == service.ErrCodeInvalid {
			RespondError(c, http.StatusBadRequest, "验证码错误或已过期")
			return
		}
		RespondError(c, http.StatusInternalServerError, "register failed: "+err.Error())
		return
	}

	Respond(c, http.StatusOK, "ok", gin.H{
		"token":         result.Token,
		"refresh_token": result.RefreshToken,
		"user":          buildCEndUserPayload(result.User),
		"profile": gin.H{
			"name": result.User.Name,
		},
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

	accountInfo, err := h.accountService.GetAccountInfo(userID)
	if err != nil {
		if err == service.ErrUserNotFound {
			RespondError(c, http.StatusNotFound, "user not found")
			return
		}
		if err == service.ErrNoCustomerProfile {
			RespondError(c, http.StatusNotFound, "customer profile not found")
			return
		}
		RespondError(c, http.StatusInternalServerError, "load current user failed")
		return
	}

	Respond(c, http.StatusOK, "ok", gin.H{
		"user_id":       accountInfo.UserID,
		"type":          accountInfo.Type,
		"customer_type": accountInfo.CustomerType,
		"name":          accountInfo.Name,
		"phone":         accountInfo.Phone,
		"status":        accountInfo.Status,
		"has_password":  accountInfo.HasPassword,
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

// ResetPassword 通过手机号验证码重置密码
// @Summary      重置登录密码
// @Description  未登录用户通过手机号验证码重置 C 端登录密码
// @Tags         c_auth
// @Accept       json
// @Produce      json
// @Param        request body object{phone=string,code=string,new_password=string} true "重置密码信息"
// @Success      200  {object} APIResponse "重置成功"
// @Failure      400  {object} APIResponse "请求参数错误"
// @Failure      404  {object} APIResponse "用户不存在"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /c/auth/reset-password [post]
func (h *CAuthHandler) ResetPassword(c *gin.Context) {
	var req struct {
		Phone       string `json:"phone" binding:"required"`
		Code        string `json:"code" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "invalid payload")
		return
	}

	err := h.authService.ResetCEndPassword(req.Phone, req.Code, req.NewPassword)
	if err != nil {
		switch err {
		case service.ErrCodeInvalid:
			RespondError(c, http.StatusBadRequest, "验证码错误或已过期")
			return
		case service.ErrUserInactive:
			RespondError(c, http.StatusForbidden, "user inactive")
			return
		case service.ErrUserNotFound, service.ErrNoCustomerProfile:
			RespondError(c, http.StatusNotFound, "用户不存在")
			return
		default:
			RespondError(c, http.StatusInternalServerError, "reset password failed")
			return
		}
	}

	Respond(c, http.StatusOK, "密码重置成功", nil)
}

// QuickStart 快速开通（注册+登录+创建服务请求）
// @Summary      快速开通（注册+登录+创建服务请求）
// @Description  一次性完成注册、登录和创建服务需求的全流程
// @Tags         c_auth
// @Accept       json
// @Produce      json
// @Param        request body object{phone=string,code=string,name=string,address=string,latitude=float64,longitude=float64,submit_lat=float64,submit_lng=float64,service_lat=float64,service_lng=float64,source_station_id=int,service_type=string,description=string,images=[]string,contact_name=string,contact_phone=string} true "快速开通请求"
// @Success      200  {object} APIResponse "快速开通成功"
// @Failure      400  {object} APIResponse "请求参数错误"
// @Failure      403  {object} APIResponse "用户已停用"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /c/auth/quick-start [post]
func (h *CAuthHandler) QuickStart(c *gin.Context) {
	var req struct {
		Phone           string   `json:"phone" binding:"required"`
		Code            string   `json:"code" binding:"required"`
		Name            string   `json:"name" binding:"required"`
		Address         string   `json:"address"`
		Latitude        *float64 `json:"latitude"`  // 兼容旧字段，等同 submit_lat
		Longitude       *float64 `json:"longitude"` // 兼容旧字段，等同 submit_lng
		SubmitLat       *float64 `json:"submit_lat"`
		SubmitLng       *float64 `json:"submit_lng"`
		ServiceLat      *float64 `json:"service_lat"`
		ServiceLng      *float64 `json:"service_lng"`
		SourceStationID *int64   `json:"source_station_id"`
		ServiceType     string   `json:"service_type" binding:"required"`
		Description     *string  `json:"description"`
		Images          []string `json:"images"`
		ContactName     string   `json:"contact_name"`  // 联系人姓名
		ContactPhone    string   `json:"contact_phone"` // 联系人电话
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "invalid payload")
		return
	}

	submitLat := req.SubmitLat
	submitLng := req.SubmitLng
	if submitLat == nil && submitLng == nil {
		submitLat = req.Latitude
		submitLng = req.Longitude
	}
	serviceLat := req.ServiceLat
	serviceLng := req.ServiceLng
	if serviceLat == nil && serviceLng == nil && req.Latitude != nil && req.Longitude != nil {
		serviceLat = req.Latitude
		serviceLng = req.Longitude
	}

	// 调用快速开通服务
	result, err := h.authService.QuickStart(service.QuickStartInput{
		Phone:            req.Phone,
		Code:             req.Code,
		Name:             req.Name,
		Address:          req.Address,
		SubmitLatitude:   submitLat,
		SubmitLongitude:  submitLng,
		ServiceLatitude:  serviceLat,
		ServiceLongitude: serviceLng,
		SourceStationID:  req.SourceStationID,
		ServiceType:      req.ServiceType,
		Description:      req.Description,
		Images:           req.Images,
		ContactName:      req.ContactName,
		ContactPhone:     req.ContactPhone,
	})

	if err != nil {
		if err == service.ErrCodeInvalid {
			RespondError(c, http.StatusBadRequest, "验证码错误或已过期")
			return
		}
		if err == service.ErrAddressRequired {
			RespondError(c, http.StatusBadRequest, "服务地址不能为空")
			return
		}
		if err == service.ErrServiceLocationRequired {
			RespondError(c, http.StatusBadRequest, "无法确定服务地点，请完善服务地址或确认当前位置")
			return
		}
		if err == service.ErrNoStation {
			RespondError(c, http.StatusBadRequest, "无法找到服务站点，请检查地址")
			return
		}
		if err == service.ErrUserInactive {
			RespondError(c, http.StatusForbidden, "user inactive")
			return
		}
		RespondError(c, http.StatusInternalServerError, "quick start failed")
		return
	}

	// 返回结果
	data := gin.H{
		"token":         result.Token,
		"refresh_token": result.RefreshToken,
		"user":          buildCEndUserPayload(result.User),
		"profile": gin.H{
			"name":    result.User.Name,
			"address": result.Profile.Address,
		},
		"request": gin.H{
			"id":            result.Request.ID,
			"request_no":    result.Request.RequestNo,
			"service_type":  result.Request.ServiceType,
			"status":        result.Request.Status,
			"contact_name":  result.Request.ContactName,
			"contact_phone": result.Request.ContactPhone,
			"images":        result.Request.Images,
			"created_at":    result.Request.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
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

	payload, err := h.accountService.GetCheckPayload(userID)
	if err != nil {
		if err == service.ErrUserNotFound {
			RespondError(c, http.StatusNotFound, "user not found")
			return
		}
		RespondError(c, http.StatusInternalServerError, "load current user failed")
		return
	}

	response := gin.H{
		"user": gin.H{
			"id":           payload.User.ID,
			"phone":        payload.User.Phone,
			"role":         payload.User.Role,
			"has_password": payload.User.HasPassword,
		},
		"profile": nil,
	}
	if payload.Profile != nil {
		response["profile"] = gin.H{
			"name":      payload.Profile.Name,
			"id_number": payload.Profile.IDNumber,
			"address":   payload.Profile.Address,
			"user_type": payload.Profile.UserType,
		}
	}

	Respond(c, http.StatusOK, "ok", response)
}

// SetPassword 设置或更新当前登录用户的密码
// @Summary      设置登录密码
// @Description  已登录的 C 端用户设置登录密码；首次设置无需提供旧密码，修改已有密码时必须提供旧密码
// @Tags         c_auth
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        request body object{current_password=string,new_password=string} true "密码信息"
// @Success      200  {object} APIResponse "设置成功"
// @Failure      400  {object} APIResponse "请求参数错误"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      404  {object} APIResponse "用户不存在"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /c/auth/password [post]
func (h *CAuthHandler) SetPassword(c *gin.Context) {
	userID, ok := GetUserID(c)
	if !ok {
		RespondError(c, http.StatusUnauthorized, "missing user")
		return
	}

	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "invalid payload")
		return
	}

	err := h.authService.SetCEndPassword(userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		switch err {
		case service.ErrCurrentPasswordRequired:
			RespondError(c, http.StatusBadRequest, "请输入当前密码")
			return
		case service.ErrCurrentPasswordInvalid:
			RespondError(c, http.StatusBadRequest, "当前密码错误")
			return
		case service.ErrUserInactive:
			RespondError(c, http.StatusForbidden, "user inactive")
			return
		case service.ErrUserNotFound:
			RespondError(c, http.StatusNotFound, "用户不存在")
			return
		default:
			RespondError(c, http.StatusInternalServerError, "set password failed")
			return
		}
	}

	Respond(c, http.StatusOK, "密码设置成功", nil)
}
