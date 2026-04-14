package handler

import (
	"net/http"

	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/service"

	"github.com/gin-gonic/gin"
)

type StationHandler struct {
	service *service.StationService
}

func NewStationHandler(service *service.StationService) *StationHandler {
	return &StationHandler{service: service}
}

type stationRequest struct {
	Address   string  `json:"address"`
	Name      string  `json:"name" binding:"required"`
	Code      string  `json:"code"`
	Phone     string  `json:"phone"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Status    string  `json:"status"`
}

type matchStationRequest struct {
	Latitude  *float64 `json:"latitude" form:"lat"`
	Longitude *float64 `json:"longitude" form:"lng"`
}

type matchStationResponse struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Code        string  `json:"code"`
	Address     string  `json:"address"`
	Phone       string  `json:"phone"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	ServiceArea string  `json:"service_area"`
	WorkHours   string  `json:"work_hours"`
}

// Create 创建服务站点
// @Summary      创建服务站点
// @Description  创建新的服务站点
// @Tags         b_station
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        request body stationRequest true "站点信息"
// @Success      200  {object} APIResponse "创建成功"
// @Failure      400  {object} APIResponse "请求参数错误"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      403  {object} APIResponse "无权限"
// @Router       /b/stations [post]
func (h *StationHandler) Create(c *gin.Context) {
	var req stationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "invalid payload")
		return
	}

	station := &model.ServiceStation{
		Address:   req.Address,
		Name:      req.Name,
		Code:      req.Code,
		Phone:     req.Phone,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Status:    req.Status,
	}
	if station.Status == "" {
		station.Status = "active"
	}

	if err := h.service.Create(station); err != nil {
		RespondError(c, http.StatusBadRequest, "create station failed")
		return
	}

	Respond(c, http.StatusOK, "ok", station)
}

// Update 更新服务站点
// @Summary      更新服务站点
// @Description  更新指定服务站点的信息
// @Tags         b_station
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id path int true "站点ID"
// @Param        request body stationRequest true "站点信息"
// @Success      200  {object} APIResponse "更新成功"
// @Failure      400  {object} APIResponse "请求参数错误"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      403  {object} APIResponse "无权限"
// @Router       /b/stations/{id} [put]
func (h *StationHandler) Update(c *gin.Context) {
	var req stationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "invalid payload")
		return
	}
	id, err := parseInt64Param(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	station := &model.ServiceStation{
		ID:        id,
		Address:   req.Address,
		Name:      req.Name,
		Code:      req.Code,
		Phone:     req.Phone,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Status:    req.Status,
	}

	if err := h.service.Update(station); err != nil {
		RespondError(c, http.StatusBadRequest, "update station failed")
		return
	}

	Respond(c, http.StatusOK, "ok", station)
}

// Delete 删除服务站点
// @Summary      删除服务站点
// @Description  删除指定的服务站点
// @Tags         b_station
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id path int true "站点ID"
// @Success      200  {object} APIResponse "删除成功"
// @Failure      400  {object} APIResponse "请求参数错误"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      403  {object} APIResponse "无权限"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /b/stations/{id} [delete]
func (h *StationHandler) Delete(c *gin.Context) {
	id, err := parseInt64Param(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.service.Delete(id); err != nil {
		RespondError(c, http.StatusInternalServerError, "delete station failed")
		return
	}
	Respond(c, http.StatusOK, "ok", nil)
}

// Get 获取服务站点详情
// @Summary      获取服务站点详情
// @Description  获取指定服务站点的详细信息
// @Tags         b_station
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id path int true "站点ID"
// @Success      200  {object} APIResponse "获取成功"
// @Failure      400  {object} APIResponse "请求参数错误"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      404  {object} APIResponse "站点不存在"
// @Router       /b/stations/{id} [get]
func (h *StationHandler) Get(c *gin.Context) {
	id, err := parseInt64Param(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	if !containsRole(GetUserIdentities(c), "admin") {
		stationID, ok := GetStationID(c)
		if !ok || stationID == 0 || stationID != id {
			RespondError(c, http.StatusForbidden, "forbidden")
			return
		}
	}

	station, err := h.service.GetByID(id)
	if err != nil {
		RespondError(c, http.StatusNotFound, "station not found")
		return
	}
	Respond(c, http.StatusOK, "ok", station)
}

// GetPublic 获取公开服务站点详情
// @Summary      获取公开服务站点详情
// @Description  获取指定服务站点的公开信息，仅返回启用中的站点
// @Tags         public
// @Accept       json
// @Produce      json
// @Param        id path int true "站点ID"
// @Success      200  {object} APIResponse "获取成功"
// @Failure      400  {object} APIResponse "请求参数错误"
// @Failure      404  {object} APIResponse "站点不存在"
// @Router       /c/stations/{id} [get]
func (h *StationHandler) GetPublic(c *gin.Context) {
	id, err := parseInt64Param(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	station, err := h.service.GetByID(id)
	if err != nil || station.Status != "active" {
		RespondError(c, http.StatusNotFound, "station not found")
		return
	}

	resp := matchStationResponse{
		ID:          station.ID,
		Name:        station.Name,
		Code:        station.Code,
		Address:     station.Address,
		Phone:       station.Phone,
		Latitude:    station.Latitude,
		Longitude:   station.Longitude,
		ServiceArea: station.ServiceArea,
		WorkHours:   station.WorkHours,
	}

	Respond(c, http.StatusOK, "ok", resp)
}

// List 获取服务站点列表
// @Summary      获取服务站点列表
// @Description  分页获取服务站点列表
// @Tags         b_station
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        page query int false "页码" default(1)
// @Param        page_size query int false "每页数量" default(10)
// @Param        keyword query string false "关键词搜索(站点名称/编号/地址/电话)"
// @Param        name query string false "站点名称(兼容旧参数，等同于 keyword)"
// @Param        status query string false "状态筛选(active/inactive)"
// @Success      200  {object} APIResponse "获取成功"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /b/stations [get]
func (h *StationHandler) List(c *gin.Context) {
	page, pageSize := GetPagination(c)
	keyword := c.Query("keyword")
	if keyword == "" {
		keyword = c.Query("name")
	}

	filter := service.StationListFilter{
		Keyword: keyword,
		Status:  c.Query("status"),
	}

	if !containsRole(GetUserIdentities(c), "admin") {
		stationID, ok := GetStationID(c)
		if !ok || stationID == 0 {
			RespondError(c, http.StatusForbidden, "forbidden")
			return
		}
		filter.StationID = stationID
	}

	stations, total, err := h.service.List(page, pageSize, filter)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "list stations failed")
		return
	}
	RespondPage(c, http.StatusOK, "ok", stations, page, pageSize, total)
}

// MatchStation 匹配服务站点
// @Summary      匹配服务站点
// @Description  根据经纬度匹配最近的服务站点
// @Tags         public
// @Accept       json
// @Produce      json
// @Param        lat query float64 false "纬度"
// @Param        lng query float64 false "经度"
// @Success      200  {object} APIResponse "匹配成功"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /c/stations/match [get]
func (h *StationHandler) MatchStation(c *gin.Context) {
	var req matchStationRequest

	// 根据请求方法选择绑定方式，避免 GET/POST 参数混淆
	if c.Request.Method == "GET" {
		c.ShouldBindQuery(&req)
	} else {
		c.ShouldBindJSON(&req)
	}

	// 从 JWT 中获取用户ID（如果已登录）
	var userID int64
	if uid, exists := c.Get("user_id"); exists {
		if id, ok := uid.(int64); ok {
			userID = id
		}
	}

	input := service.MatchStationInput{
		UserID:    userID,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	}

	station, err := h.service.MatchStation(input)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "match station failed: "+err.Error())
		return
	}

	resp := matchStationResponse{
		ID:          station.ID,
		Name:        station.Name,
		Code:        station.Code,
		Address:     station.Address,
		Phone:       station.Phone,
		Latitude:    station.Latitude,
		Longitude:   station.Longitude,
		ServiceArea: station.ServiceArea,
		WorkHours:   station.WorkHours,
	}

	Respond(c, http.StatusOK, "ok", resp)
}
