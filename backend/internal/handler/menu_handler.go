package handler

import (
	"community-elderly-care-platform/internal/consts"
	"community-elderly-care-platform/internal/repository"
	"community-elderly-care-platform/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MenuHandler struct {
	menuService       *service.MenuService
	userRepo          *repository.UserRepository
	userIdentityRepo  *repository.UserIdentityRepository
	permissionService *service.PermissionService
}

// NewMenuHandler 创建 MenuHandler
func NewMenuHandler(menuService *service.MenuService, userRepo *repository.UserRepository, userIdentityRepo *repository.UserIdentityRepository, permissionService *service.PermissionService) *MenuHandler {
	return &MenuHandler{
		menuService:       menuService,
		userRepo:          userRepo,
		userIdentityRepo:  userIdentityRepo,
		permissionService: permissionService,
	}
}

// CreateMenuRequest 创建菜单请求
type CreateMenuRequest struct {
	ParentID       int64  `json:"parent_id"`
	Name           string `json:"name" binding:"required"`
	Path           string `json:"path"`
	Component      string `json:"component"`
	Icon           string `json:"icon"`
	PermissionCode string `json:"permission_code"`
	Sort           int    `json:"sort"`
	Hidden         bool   `json:"hidden"`
	Status         string `json:"status"`
}

// UpdateMenuRequest 更新菜单请求
type UpdateMenuRequest struct {
	ParentID       int64  `json:"parent_id"`
	Name           string `json:"name" binding:"required"`
	Path           string `json:"path"`
	Component      string `json:"component"`
	Icon           string `json:"icon"`
	PermissionCode string `json:"permission_code"`
	Sort           int    `json:"sort"`
	Hidden         bool   `json:"hidden"`
	Status         string `json:"status"`
}

// BatchUpdateSortRequest 批量更新排序请求
type BatchUpdateSortRequest struct {
	Updates []struct {
		ID   int64 `json:"id" binding:"required"`
		Sort int   `json:"sort"`
	} `json:"updates" binding:"required"`
}

// List 获取菜单树（管理用）
// @Summary 获取菜单树
// @Description 获取完整菜单树，用于菜单管理
// @Tags B端-菜单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} Response
// @Failure 500 {object} Response
// @Router /b/menus [get]
func (h *MenuHandler) List(c *gin.Context) {
	menus, err := h.menuService.GetMenuTree()
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "获取菜单失败")
		return
	}
	Respond(c, http.StatusOK, "success", menus)
}

// GetUserMenus 获取当前用户可见菜单
// @Summary 获取用户菜单
// @Description 根据当前用户权限获取可见菜单树
// @Tags B端-菜单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} Response
// @Failure 401 {object} Response
// @Failure 500 {object} Response
// @Router /b/menus/user [get]
func (h *MenuHandler) GetUserMenus(c *gin.Context) {
	userID, ok := GetUserID(c)
	if !ok {
		RespondError(c, http.StatusUnauthorized, "未登录")
		return
	}

	// 获取用户所有激活的 B 端身份
	identities, err := h.userIdentityRepo.GetBEndIdentities(userID)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "获取身份失败")
		return
	}

	var userIdentities []string
	for _, identity := range identities {
		userIdentities = append(userIdentities, identity.IdentityType)
	}

	// 如果没有身份，返回空菜单
	if len(userIdentities) == 0 {
		Respond(c, http.StatusOK, "success", []interface{}{})
		return
	}

	// admin 身份拥有所有菜单
	for _, identity := range userIdentities {
		if identity == consts.IdentityAdmin {
			menus, err := h.menuService.GetMenuTree()
			if err != nil {
				RespondError(c, http.StatusInternalServerError, "获取菜单失败")
				return
			}
			Respond(c, http.StatusOK, "success", menus)
			return
		}
	}

	// 获取权限码列表
	permissionCodes, err := h.permissionService.GetUserPermissionCodes(userIdentities)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "获取权限失败")
		return
	}

	// 根据权限码获取菜单
	menus, err := h.menuService.GetUserMenus(permissionCodes)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "获取菜单失败")
		return
	}

	Respond(c, http.StatusOK, "success", menus)
}

// GetByID 获取单个菜单
// @Summary 获取菜单详情
// @Description 根据ID获取菜单详情
// @Tags B端-菜单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "菜单ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Router /b/menus/{id} [get]
func (h *MenuHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "无效的菜单ID")
		return
	}

	menu, err := h.menuService.GetByID(id)
	if err != nil {
		RespondError(c, http.StatusNotFound, "菜单不存在")
		return
	}

	Respond(c, http.StatusOK, "success", menu)
}

// Create 创建菜单
// @Summary 创建菜单
// @Description 创建新菜单
// @Tags B端-菜单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateMenuRequest true "菜单信息"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /b/menus [post]
func (h *MenuHandler) Create(c *gin.Context) {
	var req CreateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	menu := &consts.Menu{
		ParentID:       req.ParentID,
		Name:           req.Name,
		Path:           req.Path,
		Component:      req.Component,
		Icon:           req.Icon,
		PermissionCode: req.PermissionCode,
		Sort:           req.Sort,
		Hidden:         req.Hidden,
		Status:         req.Status,
	}

	if menu.Status == "" {
		menu.Status = "active"
	}

	if err := h.menuService.Create(menu); err != nil {
		RespondError(c, http.StatusInternalServerError, "创建菜单失败")
		return
	}

	Respond(c, http.StatusOK, "创建成功", menu)
}

// Update 更新菜单
// @Summary 更新菜单
// @Description 更新菜单信息
// @Tags B端-菜单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "菜单ID"
// @Param request body UpdateMenuRequest true "菜单信息"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /b/menus/{id} [put]
func (h *MenuHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "无效的菜单ID")
		return
	}

	var req UpdateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	menu := &consts.Menu{
		ParentID:       req.ParentID,
		Name:           req.Name,
		Path:           req.Path,
		Component:      req.Component,
		Icon:           req.Icon,
		PermissionCode: req.PermissionCode,
		Sort:           req.Sort,
		Hidden:         req.Hidden,
		Status:         req.Status,
	}

	if err := h.menuService.Update(id, menu); err != nil {
		RespondError(c, http.StatusInternalServerError, "更新菜单失败: "+err.Error())
		return
	}

	Respond(c, http.StatusOK, "更新成功", nil)
}

// Delete 删除菜单
// @Summary 删除菜单
// @Description 删除菜单（软删除）
// @Tags B端-菜单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "菜单ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /b/menus/{id} [delete]
func (h *MenuHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "无效的菜单ID")
		return
	}

	if err := h.menuService.Delete(id); err != nil {
		RespondError(c, http.StatusInternalServerError, "删除菜单失败")
		return
	}

	Respond(c, http.StatusOK, "删除成功", nil)
}

// BatchUpdateSort 批量更新排序
// @Summary 批量更新菜单排序
// @Description 批量更新多个菜单的排序值
// @Tags B端-菜单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body BatchUpdateSortRequest true "排序信息"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /b/menus/sort [put]
func (h *MenuHandler) BatchUpdateSort(c *gin.Context) {
	var req BatchUpdateSortRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	updates := make([]struct {
		ID   int64
		Sort int
	}, len(req.Updates))

	for i, u := range req.Updates {
		updates[i].ID = u.ID
		updates[i].Sort = u.Sort
	}

	if err := h.menuService.BatchUpdateSort(updates); err != nil {
		RespondError(c, http.StatusInternalServerError, "更新排序失败")
		return
	}

	Respond(c, http.StatusOK, "更新成功", nil)
}
