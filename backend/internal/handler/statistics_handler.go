package handler

import (
	"net/http"

	"community-elderly-care-platform/internal/service"

	"github.com/gin-gonic/gin"
)

type StatisticsHandler struct {
	service *service.StatisticsService
}

func NewStatisticsHandler(service *service.StatisticsService) *StatisticsHandler {
	return &StatisticsHandler{service: service}
}

// GetDashboardStats 获取工作台统计数据
// @Summary      获取工作台统计数据
// @Description  获取当前用户/站点的工作台统计数据，包括任务和需求统计
// @Tags         b_statistics
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200  {object} APIResponse "获取成功"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /b/statistics/dashboard [get]
func (h *StatisticsHandler) GetDashboardStats(c *gin.Context) {
	userID, ok := GetUserID(c)
	if !ok {
		RespondError(c, http.StatusUnauthorized, "missing user")
		return
	}

	stationID, _ := GetStationID(c)
	roles := GetUserRoles(c)
	isAdmin := containsRole(roles, "admin")

	stats, err := h.service.GetDashboardStats(userID, stationID, isAdmin)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "get statistics failed")
		return
	}

	Respond(c, http.StatusOK, "ok", stats)
}

// GetTaskStats 获取任务统计
// @Summary      获取任务统计
// @Description  获取任务各状态数量统计
// @Tags         b_statistics
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        station_id query int false "站点ID（管理员可指定）"
// @Success      200  {object} APIResponse "获取成功"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /b/statistics/tasks [get]
func (h *StatisticsHandler) GetTaskStats(c *gin.Context) {
	stationID, _ := GetStationID(c)
	roles := GetUserRoles(c)
	isAdmin := containsRole(roles, "admin")

	// 管理员可以指定站点
	if isAdmin {
		if sid := c.Query("station_id"); sid != "" {
			if parsed, err := parseInt64Param(sid); err == nil {
				stationID = parsed
			}
		}
	}

	stats, err := h.service.GetTaskStats(stationID, isAdmin)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "get task statistics failed")
		return
	}

	Respond(c, http.StatusOK, "ok", stats)
}

// GetRequestStats 获取需求统计
// @Summary      获取需求统计
// @Description  获取需求各状态数量统计
// @Tags         b_statistics
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        station_id query int false "站点ID（管理员可指定）"
// @Success      200  {object} APIResponse "获取成功"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /b/statistics/requests [get]
func (h *StatisticsHandler) GetRequestStats(c *gin.Context) {
	stationID, _ := GetStationID(c)
	roles := GetUserRoles(c)
	isAdmin := containsRole(roles, "admin")

	// 管理员可以指定站点
	if isAdmin {
		if sid := c.Query("station_id"); sid != "" {
			if parsed, err := parseInt64Param(sid); err == nil {
				stationID = parsed
			}
		}
	}

	stats, err := h.service.GetRequestStats(stationID, isAdmin)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "get request statistics failed")
		return
	}

	Respond(c, http.StatusOK, "ok", stats)
}

// GetTodayStats 获取今日统计
// @Summary      获取今日统计
// @Description  获取今日新增需求、完成任务等统计
// @Tags         b_statistics
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200  {object} APIResponse "获取成功"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /b/statistics/today [get]
func (h *StatisticsHandler) GetTodayStats(c *gin.Context) {
	stationID, _ := GetStationID(c)
	roles := GetUserRoles(c)
	isAdmin := containsRole(roles, "admin")

	stats, err := h.service.GetTodayStats(stationID, isAdmin)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "get today statistics failed")
		return
	}

	Respond(c, http.StatusOK, "ok", stats)
}
