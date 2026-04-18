package handler

import (
	"net/http"
	"strconv"

	"community-elderly-care-platform/internal/service"

	"github.com/gin-gonic/gin"
)

type StatisticsHandler struct {
	service *service.StatisticsService
}

// NewStatisticsHandler 创建 StatisticsHandler
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

// GetOverviewStats 获取统计概览数据
// @Summary      获取统计概览数据
// @Description  获取指定时间范围内的统计概览数据
// @Tags         b_statistics
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        days query int false "时间范围（天数，默认7）"
// @Param        station_id query int false "站点ID（管理员可指定）"
// @Success      200  {object} APIResponse "获取成功"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /b/statistics/overview [get]
func (h *StatisticsHandler) GetOverviewStats(c *gin.Context) {
	stationID, _ := GetStationID(c)
	roles := GetUserRoles(c)
	isAdmin := containsRole(roles, "admin")

	if isAdmin {
		if sid := c.Query("station_id"); sid != "" {
			if parsed, err := parseInt64Param(sid); err == nil {
				stationID = parsed
			}
		}
	}

	days := 7
	if d := c.Query("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 {
			days = parsed
		}
	}

	stats, err := h.service.GetOverviewStats(stationID, isAdmin, days)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "get overview statistics failed")
		return
	}

	Respond(c, http.StatusOK, "ok", stats)
}

// GetServiceTypeStats 获取服务类型分布统计
// @Summary      获取服务类型分布统计
// @Description  获取指定时间范围内各服务类型的需求数量分布
// @Tags         b_statistics
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        days query int false "时间范围（天数，默认7）"
// @Param        station_id query int false "站点ID（管理员可指定）"
// @Success      200  {object} APIResponse "获取成功"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /b/statistics/service-types [get]
func (h *StatisticsHandler) GetServiceTypeStats(c *gin.Context) {
	stationID, _ := GetStationID(c)
	roles := GetUserRoles(c)
	isAdmin := containsRole(roles, "admin")

	if isAdmin {
		if sid := c.Query("station_id"); sid != "" {
			if parsed, err := parseInt64Param(sid); err == nil {
				stationID = parsed
			}
		}
	}

	days := 7
	if d := c.Query("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 {
			days = parsed
		}
	}

	stats, err := h.service.GetServiceTypeStats(stationID, isAdmin, days)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "get service type statistics failed")
		return
	}

	Respond(c, http.StatusOK, "ok", stats)
}

// GetRequestTrend 获取需求趋势数据
// @Summary      获取需求趋势数据
// @Description  获取指定时间范围内每日需求数量趋势
// @Tags         b_statistics
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        days query int false "时间范围（天数，默认7）"
// @Param        station_id query int false "站点ID（管理员可指定）"
// @Success      200  {object} APIResponse "获取成功"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /b/statistics/trend [get]
func (h *StatisticsHandler) GetRequestTrend(c *gin.Context) {
	stationID, _ := GetStationID(c)
	roles := GetUserRoles(c)
	isAdmin := containsRole(roles, "admin")

	if isAdmin {
		if sid := c.Query("station_id"); sid != "" {
			if parsed, err := parseInt64Param(sid); err == nil {
				stationID = parsed
			}
		}
	}

	days := 7
	if d := c.Query("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 {
			days = parsed
		}
	}

	trend, err := h.service.GetRequestTrend(stationID, isAdmin, days)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "get request trend failed")
		return
	}

	Respond(c, http.StatusOK, "ok", trend)
}

// GetEfficiencyStats 获取处理效率统计
// @Summary      获取处理效率统计
// @Description  获取平均响应时间、处理时间和满意度
// @Tags         b_statistics
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        days query int false "时间范围（天数，默认7）"
// @Param        station_id query int false "站点ID（管理员可指定）"
// @Success      200  {object} APIResponse "获取成功"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /b/statistics/efficiency [get]
func (h *StatisticsHandler) GetEfficiencyStats(c *gin.Context) {
	stationID, _ := GetStationID(c)
	roles := GetUserRoles(c)
	isAdmin := containsRole(roles, "admin")

	if isAdmin {
		if sid := c.Query("station_id"); sid != "" {
			if parsed, err := parseInt64Param(sid); err == nil {
				stationID = parsed
			}
		}
	}

	days := 7
	if d := c.Query("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 {
			days = parsed
		}
	}

	stats, err := h.service.GetEfficiencyStats(stationID, isAdmin, days)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "get efficiency statistics failed")
		return
	}

	Respond(c, http.StatusOK, "ok", stats)
}

// GetStaffRanking 获取服务人员排行
// @Summary      获取服务人员排行
// @Description  获取指定时间范围内服务人员完成任务排行
// @Tags         b_statistics
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        days query int false "时间范围（天数，默认7）"
// @Param        limit query int false "返回数量（默认10）"
// @Param        station_id query int false "站点ID（管理员可指定）"
// @Success      200  {object} APIResponse "获取成功"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /b/statistics/staff-ranking [get]
func (h *StatisticsHandler) GetStaffRanking(c *gin.Context) {
	stationID, _ := GetStationID(c)
	roles := GetUserRoles(c)
	isAdmin := containsRole(roles, "admin")

	if isAdmin {
		if sid := c.Query("station_id"); sid != "" {
			if parsed, err := parseInt64Param(sid); err == nil {
				stationID = parsed
			}
		}
	}

	days := 7
	if d := c.Query("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 {
			days = parsed
		}
	}

	limit := 10
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	ranking, err := h.service.GetStaffRanking(stationID, isAdmin, days, limit)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "get staff ranking failed")
		return
	}

	Respond(c, http.StatusOK, "ok", ranking)
}
