package handler

import (
	"net/http"

	"community-elderly-care-platform/internal/service"
	"community-elderly-care-platform/pkg/geo"

	"github.com/gin-gonic/gin"
)

type ZoneHandler struct {
	service *service.ZoneService
}

func NewZoneHandler(service *service.ZoneService) *ZoneHandler {
	return &ZoneHandler{service: service}
}

type zoneRequest struct {
	StationID int64            `json:"station_id" binding:"required"`
	Name      string           `json:"name" binding:"required"`
	Points    []geo.Point `json:"points" binding:"required"`
	Priority  int              `json:"priority"`
	Status    string           `json:"status"`
}

// Create 创建服务围栏
// @Summary      创建服务围栏
// @Description  创建新的地理围栏区域
// @Tags         b_zone
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        request body zoneRequest true "围栏信息"
// @Success      200  {object} APIResponse "创建成功"
// @Failure      400  {object} APIResponse "请求参数错误"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      403  {object} APIResponse "无权限"
// @Router       /b/zones [post]
func (h *ZoneHandler) Create(c *gin.Context) {
	var req zoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "invalid payload")
		return
	}
	if len(req.Points) < 3 {
		RespondError(c, http.StatusBadRequest, "invalid points")
		return
	}

	zone, err := h.service.Create(service.ZoneInput{
		StationID: req.StationID,
		Name:      req.Name,
		Points:    req.Points,
		Priority:  req.Priority,
		Status:    req.Status,
	})
	if err != nil {
		RespondError(c, http.StatusBadRequest, "create zone failed")
		return
	}

	Respond(c, http.StatusOK, "ok", zone)
}

// Update 更新服务围栏
// @Summary      更新服务围栏
// @Description  更新指定地理围栏的信息
// @Tags         b_zone
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id path int true "围栏ID"
// @Param        request body zoneRequest true "围栏信息"
// @Success      200  {object} APIResponse "更新成功"
// @Failure      400  {object} APIResponse "请求参数错误"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      403  {object} APIResponse "无权限"
// @Router       /b/zones/{id} [put]
func (h *ZoneHandler) Update(c *gin.Context) {
	var req zoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "invalid payload")
		return
	}
	id, err := parseInt64Param(c.Param("id"))
	if err != nil || len(req.Points) < 3 {
		RespondError(c, http.StatusBadRequest, "invalid payload")
		return
	}

	zone, err := h.service.Update(service.ZoneInput{
		ID:        id,
		StationID: req.StationID,
		Name:      req.Name,
		Points:    req.Points,
		Priority:  req.Priority,
		Status:    req.Status,
	})
	if err != nil {
		RespondError(c, http.StatusBadRequest, "update zone failed")
		return
	}

	Respond(c, http.StatusOK, "ok", zone)
}

// Delete 删除服务围栏
// @Summary      删除服务围栏
// @Description  删除指定的地理围栏
// @Tags         b_zone
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id path int true "围栏ID"
// @Success      200  {object} APIResponse "删除成功"
// @Failure      400  {object} APIResponse "请求参数错误"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      403  {object} APIResponse "无权限"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /b/zones/{id} [delete]
func (h *ZoneHandler) Delete(c *gin.Context) {
	id, err := parseInt64Param(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.service.Delete(id); err != nil {
		RespondError(c, http.StatusInternalServerError, "delete zone failed")
		return
	}
	Respond(c, http.StatusOK, "ok", nil)
}

// List 获取服务围栏列表
// @Summary      获取服务围栏列表
// @Description  分页获取地理围栏列表
// @Tags         b_zone
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        page query int false "页码" default(1)
// @Param        page_size query int false "每页数量" default(10)
// @Success      200  {object} APIResponse "获取成功"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /b/zones [get]
func (h *ZoneHandler) List(c *gin.Context) {
	page, pageSize := GetPagination(c)
	zones, total, err := h.service.List(page, pageSize)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "list zones failed")
		return
	}
	RespondPage(c, http.StatusOK, "ok", zones, page, pageSize, total)
}
