package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"community-elderly-care-platform/internal/service"
	"community-elderly-care-platform/pkg/crypto"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service        *service.UserService
	idCardTokenKey []byte
	idCardTokenTTL time.Duration
	encryptKey     []byte // AES-256 加密密钥（32 字节）
}

// NewUserHandler 创建 UserHandler
func NewUserHandler(service *service.UserService, tokenSecret string, encryptKeyBase64 string) *UserHandler {
	if tokenSecret == "" {
		tokenSecret = "scare-default-id-card-token-secret"
	}
	var encryptKey []byte
	if encryptKeyBase64 != "" {
		if k, err := base64.StdEncoding.DecodeString(encryptKeyBase64); err == nil && len(k) == 32 {
			encryptKey = k
		}
	}
	return &UserHandler{
		service:        service,
		idCardTokenKey: []byte(tokenSecret + ":id_card_token"),
		idCardTokenTTL: 24 * time.Hour,
		encryptKey:     encryptKey,
	}
}

type userCreateRequest struct {
	Phone        string `json:"phone" binding:"required"`
	Password     string `json:"password" binding:"required"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	IdentityType string `json:"identity_type" binding:"required"` // 身份类型
	StationID    int64  `json:"station_id"`
	Status       string `json:"status"`
}

type userUpdateRequest struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Avatar      string `json:"avatar"`
	Gender      string `json:"gender"`
	BirthDate   string `json:"birth_date"` // 格式：2006-01-02
	IDCard      string `json:"id_card"`
	IDCardHash  string `json:"id_card_hash"`
	IDCardToken string `json:"id_card_token"`
	StationID   int64  `json:"station_id"`
	Status      string `json:"status"`
	Password    string `json:"password"`
}

// Create 创建用户
// @Summary      创建B端用户
// @Description  创建新的B端工作人员账号
// @Tags         b_user
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        request body userCreateRequest true "用户信息"
// @Success      200  {object} APIResponse{data=UserResponse} "创建成功"
// @Failure      400  {object} APIResponse "请求参数错误"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      403  {object} APIResponse "无权限"
// @Router       /b/users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req userCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "invalid payload")
		return
	}

	user, err := h.service.Create(service.UserInput{
		Phone:        req.Phone,
		Password:     req.Password,
		Name:         req.Name,
		Email:        req.Email,
		IdentityType: req.IdentityType,
		StationID:    req.StationID,
		Status:       req.Status,
	})
	if err != nil {
		RespondError(c, http.StatusBadRequest, "create user failed: "+err.Error())
		return
	}

	Respond(c, http.StatusOK, "ok", h.toUserResponse(user))
}

// Update 更新用户
// @Summary      更新B端用户
// @Description  更新指定B端用户的信息
// @Tags         b_user
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id path int true "用户ID"
// @Param        request body userUpdateRequest true "用户信息"
// @Success      200  {object} APIResponse{data=UserResponse} "更新成功"
// @Failure      400  {object} APIResponse "请求参数错误"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      403  {object} APIResponse "无权限"
// @Router       /b/users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	id, err := parseInt64Param(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	var req userUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "invalid payload")
		return
	}
	current, err := h.service.GetByID(id)
	if err != nil {
		RespondError(c, http.StatusNotFound, "user not found")
		return
	}
	currentIDCardHash := current.IDCardHmac
	if currentIDCardHash == "" && current.IDCard != "" {
		currentIDCardHash = h.idCardDigest(current.IDCard)
	}
	if req.IDCard != "" && req.IDCardHash != "" {
		RespondError(c, http.StatusBadRequest, "id_card_hash should not be sent when id_card is provided")
		return
	}
	if req.IDCard != "" && req.IDCardToken != "" {
		RespondError(c, http.StatusBadRequest, "id_card_token should not be sent when id_card is provided")
		return
	}
	if req.IDCard == "" && req.IDCardToken != "" && !h.verifyIDCardToken(id, currentIDCardHash, req.IDCardToken) {
		RespondError(c, http.StatusBadRequest, "invalid id_card_token")
		return
	}
	if req.IDCard == "" && req.IDCardHash != "" && !secureStringEqual(req.IDCardHash, currentIDCardHash) {
		RespondError(c, http.StatusBadRequest, "invalid id_card_hash")
		return
	}

	// 解析出生日期
	var birthDate time.Time
	if req.BirthDate != "" {
		birthDate, err = time.Parse("2006-01-02", req.BirthDate)
		if err != nil {
			RespondError(c, http.StatusBadRequest, "invalid birth_date format, expected YYYY-MM-DD")
			return
		}
	}

	var newIDCardHash string
	var newIDCardEncrypted string
	var newIDCardMasked string
	if req.IDCard != "" {
		newIDCardHash = h.idCardDigest(req.IDCard)
		newIDCardMasked = maskIDCard(req.IDCard)
		if h.encryptKey != nil {
			encrypted, err := crypto.Encrypt(req.IDCard, h.encryptKey)
			if err != nil {
				RespondError(c, http.StatusInternalServerError, "encrypt id_card failed")
				return
			}
			newIDCardEncrypted = encrypted
		} else {
			newIDCardEncrypted = req.IDCard // 未配置密钥时保持明文（兼容）
		}
	}

	user, err := h.service.Update(service.UserInput{
		ID:           id,
		Name:         req.Name,
		Email:        req.Email,
		Avatar:       req.Avatar,
		Gender:       req.Gender,
		BirthDate:    birthDate,
		IDCard:       newIDCardEncrypted,
		IDCardHMAC:   newIDCardHash,
		IDCardMasked: newIDCardMasked,
		StationID:    req.StationID,
		Status:       req.Status,
		Password:     req.Password,
	})
	if err != nil {
		RespondError(c, http.StatusBadRequest, "update user failed")
		return
	}

	Respond(c, http.StatusOK, "ok", h.toUserResponse(user))
}

// List 获取用户列表
// @Summary      获取B端用户列表
// @Description  分页获取B端用户列表，支持筛选
// @Tags         b_user
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        page query int false "页码" default(1)
// @Param        page_size query int false "每页数量" default(10)
// @Param        role query string false "身份筛选(elderly/family/staff/station_manager/admin)"
// @Param        status query string false "状态筛选(active/inactive/suspended)"
// @Param        station_id query int false "站点ID筛选"
// @Param        keyword query string false "关键词搜索(姓名/手机号)"
// @Success      200  {object} APIResponse{data=UserListResponse} "获取成功"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /b/users [get]
func (h *UserHandler) List(c *gin.Context) {
	page, pageSize := GetPagination(c)

	// 解析筛选参数
	filter := service.UserFilter{
		Role:      c.Query("role"),
		Status:    c.Query("status"),
		StationID: parseInt64Query(c, "station_id"),
		Keyword:   c.Query("keyword"),
	}

	users, total, err := h.service.ListWithFilter(page, pageSize, filter)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "list users failed")
		return
	}

	items := make([]gin.H, 0, len(users))
	for _, u := range users {
		items = append(items, h.toUserResponse(u))
	}
	Respond(c, http.StatusOK, "ok", gin.H{
		"items": items,
		"total": total,
	})
}

// GetByID 获取用户详情
// @Summary      获取B端用户详情
// @Description  根据ID获取B端用户详细信息
// @Tags         b_user
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id path int true "用户ID"
// @Success      200  {object} APIResponse{data=UserResponse} "获取成功"
// @Failure      400  {object} APIResponse "请求参数错误"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      404  {object} APIResponse "用户不存在"
// @Router       /b/users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := parseInt64Param(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	user, err := h.service.GetByID(id)
	if err != nil {
		RespondError(c, http.StatusNotFound, "user not found")
		return
	}
	Respond(c, http.StatusOK, "ok", h.toUserResponse(user))
}

// toUserResponse 将用户对象转换为 API 响应格式
func (h *UserHandler) toUserResponse(user *service.UserWithIdentities) gin.H {
	var primaryIdentity string
	if user.PrimaryIdentity != nil {
		primaryIdentity = user.PrimaryIdentity.IdentityType
	}
	var birthDate string
	var age *int
	if !user.BirthDate.IsZero() {
		birthDate = user.BirthDate.Format("2006-01-02")
		calculatedAge := calculateAge(user.BirthDate, time.Now())
		age = &calculatedAge
	}
	idCardMasked := user.IDCardMasked
	if idCardMasked == "" && user.IDCard != "" {
		// 兼容：旧数据未存储脱敏值时，尝试从明文计算（仅未加密场景）
		idCardMasked = maskIDCard(user.IDCard)
	}
	idCardHash := user.IDCardHmac
	if idCardHash == "" && user.IDCard != "" {
		idCardHash = h.idCardDigest(user.IDCard)
	}
	idCardToken := ""
	if idCardHash != "" {
		token, err := h.generateIDCardToken(user.ID, idCardHash)
		if err == nil {
			idCardToken = token
		}
	}

	return gin.H{
		"id":               user.ID,
		"phone":            user.Phone,
		"name":             user.Name,
		"email":            user.Email,
		"avatar":           user.Avatar,
		"gender":           user.Gender,
		"birth_date":       birthDate,
		"age":              age,
		"id_card_masked":   idCardMasked,
		"id_card_hash":     idCardHash,
		"id_card_token":    idCardToken,
		"b_end_identities": user.BEndIdentities,
		"c_end_identities": user.CEndIdentities,
		"primary_identity": primaryIdentity,
		"station_id":       user.StationID,
		"status":           user.Status,
		"created_at":       user.CreatedAt,
		"updated_at":       user.UpdatedAt,
	}
}

// parseInt64Query 解析 query 参数为 int64
func parseInt64Query(c *gin.Context, key string) int64 {
	val := c.Query(key)
	if val == "" {
		return 0
	}
	id, _ := parseInt64Param(val)
	return id
}

// calculateAge 根据出生日期计算年龄
func calculateAge(birthDate time.Time, now time.Time) int {
	age := now.Year() - birthDate.Year()
	if now.Month() < birthDate.Month() || (now.Month() == birthDate.Month() && now.Day() < birthDate.Day()) {
		age--
	}
	if age < 0 {
		return 0
	}
	return age
}

type idCardTokenPayload struct {
	UserID int64  `json:"uid"`
	Digest string `json:"dig"`
	Exp    int64  `json:"exp"`
}

// generateIDCardToken 生成身份证 Token
func (h *UserHandler) generateIDCardToken(userID int64, idCardHash string) (string, error) {
	payload := idCardTokenPayload{
		UserID: userID,
		Digest: idCardHash,
		Exp:    time.Now().Add(h.idCardTokenTTL).Unix(),
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	payloadEncoded := base64.RawURLEncoding.EncodeToString(payloadBytes)
	signature := h.sign(payloadEncoded)
	signatureEncoded := base64.RawURLEncoding.EncodeToString(signature)
	return payloadEncoded + "." + signatureEncoded, nil
}

// verifyIDCardToken 验证身份证 Token
func (h *UserHandler) verifyIDCardToken(userID int64, idCardHash string, token string) bool {
	if idCardHash == "" {
		return false
	}
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return false
	}
	payloadEncoded := parts[0]
	signatureEncoded := parts[1]

	expectedSig := h.sign(payloadEncoded)
	actualSig, err := base64.RawURLEncoding.DecodeString(signatureEncoded)
	if err != nil {
		return false
	}
	if len(actualSig) != len(expectedSig) {
		return false
	}
	if subtle.ConstantTimeCompare(actualSig, expectedSig) != 1 {
		return false
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(payloadEncoded)
	if err != nil {
		return false
	}
	var payload idCardTokenPayload
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return false
	}
	if payload.UserID != userID {
		return false
	}
	if payload.Exp <= time.Now().Unix() {
		return false
	}
	if !secureStringEqual(payload.Digest, idCardHash) {
		return false
	}
	return true
}

// idCardDigest 计算身份证 HMAC 摘要
func (h *UserHandler) idCardDigest(idCard string) string {
	mac := hmac.New(sha256.New, h.idCardTokenKey)
	mac.Write([]byte("id_card:"))
	mac.Write([]byte(idCard))
	return hex.EncodeToString(mac.Sum(nil))
}

// sign 生成 HMAC 签名
func (h *UserHandler) sign(payloadEncoded string) []byte {
	mac := hmac.New(sha256.New, h.idCardTokenKey)
	mac.Write([]byte(payloadEncoded))
	return mac.Sum(nil)
}

// secureStringEqual 安全字符串比较（防时序攻击）
func secureStringEqual(a string, b string) bool {
	aBytes := []byte(a)
	bBytes := []byte(b)
	if len(aBytes) != len(bBytes) {
		return false
	}
	return subtle.ConstantTimeCompare(aBytes, bBytes) == 1
}

// maskIDCard 脱敏身份证号
func maskIDCard(idCard string) string {
	if idCard == "" {
		return ""
	}
	runes := []rune(idCard)
	n := len(runes)
	if n <= 8 {
		if n <= 2 {
			return strings.Repeat("*", n)
		}
		return string(runes[:1]) + strings.Repeat("*", n-2) + string(runes[n-1:])
	}
	return string(runes[:4]) + strings.Repeat("*", n-8) + string(runes[n-4:])
}
