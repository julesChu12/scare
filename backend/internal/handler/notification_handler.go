package handler

import (
	"net/http"

	"community-elderly-care-platform/internal/service"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	service *service.NotificationService
}

func NewNotificationHandler(service *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

// List 获取通知列表
// @Summary      获取通知列表
// @Description  获取当前用户的通知列表，支持分页
// @Tags         c_notification,b_notification
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        page query int false "页码" default(1)
// @Param        page_size query int false "每页数量"
// @Success      200  {object} APIResponse "获取成功"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /c/notifications [get]
// @Router       /b/notifications [get]
func (h *NotificationHandler) List(c *gin.Context) {
	userID, ok := GetUserID(c)
	if !ok {
		RespondError(c, http.StatusUnauthorized, "missing user")
		return
	}
	page, pageSize := GetPagination(c)
	notifications, total, err := h.service.List(userID, page, pageSize)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "list notifications failed")
		return
	}
	RespondPage(c, http.StatusOK, "ok", notifications, page, pageSize, total)
}

// MarkRead 标记通知为已读
// @Summary      标记通知为已读
// @Description  将指定的通知标记为已读状态
// @Tags         c_notification,b_notification
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id path int true "通知ID"
// @Success      200  {object} APIResponse "标记成功"
// @Failure      400  {object} APIResponse "ID无效"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /c/notifications/{id}/read [post]
// @Router       /b/notifications/{id}/read [post]
func (h *NotificationHandler) MarkRead(c *gin.Context) {
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
	if err := h.service.MarkRead(userID, id); err != nil {
		RespondError(c, http.StatusInternalServerError, "mark read failed")
		return
	}
	Respond(c, http.StatusOK, "ok", nil)
}
