package handler

import (
	"net/http"

	"community-elderly-care-platform/internal/service"

	"github.com/gin-gonic/gin"
)

// CProfileHandler C端个人资料处理器
type CProfileHandler struct {
	profileService *service.CEndProfileService
	geocodeService *service.GeocodeService
}

// NewCProfileHandler 创建 CProfileHandler
func NewCProfileHandler(profileService *service.CEndProfileService, geocodeService *service.GeocodeService) *CProfileHandler {
	return &CProfileHandler{
		profileService: profileService,
		geocodeService: geocodeService,
	}
}

// UpdateProfile 更新个人资料
// @Summary      更新个人资料
// @Description  更新C端用户的个人资料信息
// @Tags         c_profile
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        request body object{name=string,id_number=string,address=string,user_type=string} false "个人资料更新"
// @Success      200  {object} APIResponse "更新成功"
// @Failure      400  {object} APIResponse "请求参数错误"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      404  {object} APIResponse "档案不存在"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /c/profile [put]
func (h *CProfileHandler) UpdateProfile(c *gin.Context) {
	userID, ok := GetUserID(c)
	if !ok {
		RespondError(c, http.StatusUnauthorized, "missing user")
		return
	}

	var req struct {
		Name     *string `json:"name"`
		IDNumber *string `json:"id_number"`
		Address  *string `json:"address"`
		UserType *string `json:"user_type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "invalid payload")
		return
	}

	profile, err := h.profileService.Update(userID, service.CEndProfileUpdateInput{
		Name:     req.Name,
		IDNumber: req.IDNumber,
		Address:  req.Address,
		UserType: req.UserType,
	})
	if err != nil {
		if err == service.ErrNoCustomerProfile {
			RespondError(c, http.StatusNotFound, "profile not found")
			return
		}
		if err == service.ErrUserNotFound {
			RespondError(c, http.StatusNotFound, "user not found")
			return
		}
		RespondError(c, http.StatusInternalServerError, "update failed")
		return
	}

	Respond(c, http.StatusOK, "ok", profile)
}

// Geocode 地址解析（地址字符串 → 经纬度）
// @Summary      地址解析
// @Description  将地址字符串转换为经纬度坐标
// @Tags         c_profile
// @Accept       json
// @Produce      json
// @Param        request body object{address=string} true "地址"
// @Success      200  {object} APIResponse "解析成功"
// @Failure      400  {object} APIResponse "请求参数错误"
// @Failure      404  {object} APIResponse "地址未找到"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /c/geocode [post]
func (h *CProfileHandler) Geocode(c *gin.Context) {
	var req struct {
		Address string `json:"address" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "invalid payload")
		return
	}

	// 调用地址解析服务
	result, err := h.geocodeService.Geocode(req.Address)
	if err != nil {
		if err == service.ErrGeocodeNotFound {
			RespondError(c, http.StatusNotFound, "address not found")
			return
		}
		RespondError(c, http.StatusInternalServerError, "geocode failed")
		return
	}

	Respond(c, http.StatusOK, "ok", gin.H{
		"latitude":          result.Latitude,
		"longitude":         result.Longitude,
		"formatted_address": result.FormattedAddress,
	})
}

// ReverseGeocode 逆地理编码（经纬度 → 地址信息）
// @Summary      逆地理编码
// @Description  将经纬度坐标转换为地址信息
// @Tags         c_profile
// @Accept       json
// @Produce      json
// @Param        lat query float64 true "纬度"
// @Param        lng query float64 true "经度"
// @Success      200  {object} APIResponse "转换成功"
// @Failure      400  {object} APIResponse "请求参数错误"
// @Failure      404  {object} APIResponse "位置未找到"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /c/geocode/reverse [get]
func (h *CProfileHandler) ReverseGeocode(c *gin.Context) {
	var req struct {
		Latitude  float64 `form:"lat" binding:"required"`
		Longitude float64 `form:"lng" binding:"required"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "invalid parameters")
		return
	}

	// 调用逆地理编码服务
	result, err := h.geocodeService.ReverseGeocode(req.Latitude, req.Longitude)
	if err != nil {
		if err == service.ErrReverseGeocodeNotFound {
			RespondError(c, http.StatusNotFound, "location not found")
			return
		}
		RespondError(c, http.StatusInternalServerError, "reverse geocode failed")
		return
	}

	Respond(c, http.StatusOK, "ok", gin.H{
		"province": result.Province,
		"city":     result.City,
		"district": result.District,
		"address":  result.FormattedAddress,
	})
}
