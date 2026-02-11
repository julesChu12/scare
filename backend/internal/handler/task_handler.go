package handler

import (
	"net/http"

	"community-elderly-care-platform/internal/repository"
	"community-elderly-care-platform/internal/service"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(service *service.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

type taskCompleteRequest struct {
	Images []string `json:"images"`
}

type taskTransferRequest struct {
	StaffID int64 `json:"staff_id" binding:"required"`
}

// GetByID 获取任务详情
// @Summary      获取任务详情
// @Description  根据任务ID获取任务详情，包含关联的服务请求信息
// @Tags         b_task
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id path int true "任务ID"
// @Success      200  {object} APIResponse "获取成功"
// @Failure      400  {object} APIResponse "参数错误"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      404  {object} APIResponse "任务不存在"
// @Router       /b/tasks/{id} [get]
func (h *TaskHandler) GetByID(c *gin.Context) {
	id, err := parseInt64Param(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	task, err := h.service.GetByID(id)
	if err != nil {
		RespondError(c, http.StatusNotFound, "task not found")
		return
	}
	Respond(c, http.StatusOK, "ok", task)
}

// List 获取任务列表（支持多条件筛选）
// @Summary      获取任务列表
// @Description  获取任务列表，支持按状态、服务类型、需求编号等筛选
// @Tags         b_task
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        page query int false "页码" default(1)
// @Param        page_size query int false "每页数量" default(10)
// @Param        station_id query int false "站点ID（管理员可指定）"
// @Param        status query string false "任务状态"
// @Param        service_type query string false "服务类型"
// @Param        request_no query string false "需求编号（模糊匹配）"
// @Success      200  {object} APIResponse "获取成功"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /b/tasks [get]
func (h *TaskHandler) List(c *gin.Context) {
	// 获取用户站点和角色
	userStationID, _ := GetStationID(c)
	roles := GetUserRoles(c)
	isAdmin := containsRole(roles, "admin")

	// 解析筛选参数
	queryStationID := int64(0)
	if sid := c.Query("station_id"); sid != "" {
		if parsed, err := parseInt64Param(sid); err == nil {
			queryStationID = parsed
		}
	}

	// 确定最终使用的站点ID
	var stationID int64
	if isAdmin {
		stationID = queryStationID
	} else {
		if userStationID == 0 {
			RespondError(c, http.StatusBadRequest, "missing station")
			return
		}
		stationID = userStationID
	}

	page, pageSize := GetPagination(c)
	filter := repository.TaskListFilter{
		StationID:   stationID,
		Status:      c.Query("status"),
		ServiceType: c.Query("service_type"),
		RequestNo:   c.Query("request_no"),
	}

	tasks, total, err := h.service.ListWithFilter(filter, page, pageSize)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "list tasks failed")
		return
	}
	RespondPage(c, http.StatusOK, "ok", tasks, page, pageSize, total)
}

// ListPool 获取任务池列表
// @Summary      获取任务池列表
// @Description  获取当前站点待认领的任务列表，管理员可查看所有站点
// @Tags         b_task
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        page query int false "页码" default(1)
// @Param        page_size query int false "每页数量" default(10)
// @Param        station_id query int false "站点ID（管理员可指定，0表示所有站点）"
// @Param        status query string false "任务状态" default(dispatched)
// @Success      200  {object} APIResponse "获取成功"
// @Failure      400  {object} APIResponse "缺少站点信息"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /b/tasks/pool [get]
func (h *TaskHandler) ListPool(c *gin.Context) {
	// 获取用户站点和角色
	userStationID, _ := GetStationID(c)
	roles := GetUserRoles(c)
	isAdmin := containsRole(roles, "admin")

	// 解析筛选参数
	queryStationID := int64(0)
	if sid := c.Query("station_id"); sid != "" {
		if parsed, err := parseInt64Param(sid); err == nil {
			queryStationID = parsed
		}
	}
	status := c.Query("status")

	// 确定最终使用的站点ID
	var stationID int64
	if isAdmin {
		// 管理员：可以指定站点或查看所有（0）
		stationID = queryStationID
	} else {
		// 非管理员：只能查看自己站点
		if userStationID == 0 {
			RespondError(c, http.StatusBadRequest, "missing station")
			return
		}
		stationID = userStationID
	}

	page, pageSize := GetPagination(c)
	filter := service.TaskPoolFilter{
		StationID: stationID,
		Status:    status,
	}
	tasks, total, err := h.service.ListPoolWithFilter(filter, page, pageSize)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "list tasks failed")
		return
	}
	RespondPage(c, http.StatusOK, "ok", tasks, page, pageSize, total)
}

// containsRole 检查角色列表中是否包含指定角色
func containsRole(roles []string, role string) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

// ListMy 获取我的任务列表
// @Summary      获取我的任务列表
// @Description  获取当前用户已认领的任务列表
// @Tags         b_task
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        page query int false "页码" default(1)
// @Param        page_size query int false "每页数量" default(10)
// @Success      200  {object} APIResponse "获取成功"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /b/tasks/my [get]
func (h *TaskHandler) ListMy(c *gin.Context) {
	userID, ok := GetUserID(c)
	if !ok {
		RespondError(c, http.StatusUnauthorized, "missing user")
		return
	}
	page, pageSize := GetPagination(c)
	tasks, total, err := h.service.ListByStaff(userID, page, pageSize)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "list tasks failed")
		return
	}
	RespondPage(c, http.StatusOK, "ok", tasks, page, pageSize, total)
}

// Claim 认领任务
// @Summary      认领任务
// @Description  工作人员认领指定任务
// @Tags         b_task
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id path int true "任务ID"
// @Success      200  {object} APIResponse "认领成功"
// @Failure      400  {object} APIResponse "任务无效"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      409  {object} APIResponse "任务已被认领"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /b/tasks/{id}/claim [post]
func (h *TaskHandler) Claim(c *gin.Context) {
	userID, ok := GetUserID(c)
	if !ok {
		RespondError(c, http.StatusUnauthorized, "missing user")
		return
	}
	id, err := parseInt64Param(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	task, _, err := h.service.Claim(id, userID)
	if err != nil {
		switch err {
		case service.ErrTaskInvalid:
			RespondError(c, http.StatusBadRequest, "invalid task")
		case service.ErrTaskConflict:
			RespondError(c, http.StatusConflict, "task conflict")
		default:
			RespondError(c, http.StatusInternalServerError, "claim failed")
		}
		return
	}
	Respond(c, http.StatusOK, "ok", task)
}

// Complete 完成任务
// @Summary      完成任务
// @Description  标记任务为已完成，可上传服务图片
// @Tags         b_task
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id path int true "任务ID"
// @Param        request body taskCompleteRequest false "完成信息"
// @Success      200  {object} APIResponse "完成成功"
// @Failure      400  {object} APIResponse "任务无效"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      409  {object} APIResponse "任务状态冲突"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /b/tasks/{id}/complete [post]
func (h *TaskHandler) Complete(c *gin.Context) {
	userID, ok := GetUserID(c)
	if !ok {
		RespondError(c, http.StatusUnauthorized, "missing user")
		return
	}
	id, err := parseInt64Param(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	var req taskCompleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "invalid payload")
		return
	}

	task, _, err := h.service.Complete(id, userID, req.Images)
	if err != nil {
		switch err {
		case service.ErrTaskInvalid:
			RespondError(c, http.StatusBadRequest, "invalid task")
		case service.ErrTaskConflict:
			RespondError(c, http.StatusConflict, "task conflict")
		default:
			RespondError(c, http.StatusInternalServerError, "complete failed")
		}
		return
	}
	Respond(c, http.StatusOK, "ok", task)
}

// Transfer 转派任务
// @Summary      转派任务
// @Description  管理员或站长将任务转派给其他工作人员
// @Tags         b_task
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id path int true "任务ID"
// @Param        request body taskTransferRequest true "转派信息"
// @Success      200  {object} APIResponse "转派成功"
// @Failure      400  {object} APIResponse "参数错误"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      409  {object} APIResponse "任务状态冲突"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /b/tasks/{id}/transfer [post]
func (h *TaskHandler) Transfer(c *gin.Context) {
	id, err := parseInt64Param(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	var req taskTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "invalid payload: "+err.Error())
		return
	}

	task, err := h.service.Transfer(id, req.StaffID)
	if err != nil {
		switch err {
		case service.ErrTaskInvalid:
			RespondError(c, http.StatusBadRequest, "invalid task")
		case service.ErrTaskConflict:
			RespondError(c, http.StatusConflict, "task conflict")
		default:
			RespondError(c, http.StatusInternalServerError, "transfer failed")
		}
		return
	}
	Respond(c, http.StatusOK, "ok", task)
}
