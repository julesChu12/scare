package handler

import (
	"net/http"

	"community-elderly-care-platform/internal/service"

	"github.com/gin-gonic/gin"
)

type BannerHandler struct {
	service *service.BannerService
}

func NewBannerHandler(service *service.BannerService) *BannerHandler {
	return &BannerHandler{service: service}
}

// ListForC 获取 C 端 Banner（公开接口）
// @Summary      获取轮播图列表
// @Description  获取当前站点或全局的轮播图
// @Tags         public
// @Accept       json
// @Produce      json
// @Param        station_id query int false "站点ID（0或不传表示全局）"
// @Success      200  {object} APIResponse "获取成功"
// @Router       /c/banners [get]
func (h *BannerHandler) ListForC(c *gin.Context) {
	stationID := int64(0)
	if sid, err := parseInt64Param(c.Query("station_id")); err == nil && sid > 0 {
		stationID = sid
	}

	var banners interface{}
	var err error

	if stationID > 0 {
		banners, err = h.service.ListForStation(stationID)
	} else {
		banners, err = h.service.ListGlobal()
	}

	if err != nil {
		RespondError(c, http.StatusInternalServerError, "failed to get banners")
		return
	}

	Respond(c, http.StatusOK, "ok", banners)
}

type bannerRequest struct {
	StationID int64  `json:"station_id"`
	Title     string `json:"title"`
	ImageURL  string `json:"image_url" binding:"required"`
	LinkType  string `json:"link_type"`
	LinkValue string `json:"link_value"`
	Sort      int32  `json:"sort"`
	Status    string `json:"status"`
}

// List 获取 Banner 列表（管理端）
// @Summary      获取轮播图列表（管理端）
// @Description  分页获取轮播图列表
// @Tags         b_banner
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        page query int false "页码" default(1)
// @Param        page_size query int false "每页数量" default(10)
// @Param        station_id query int false "站点ID筛选"
// @Success      200  {object} APIResponse "获取成功"
// @Router       /b/banners [get]
func (h *BannerHandler) List(c *gin.Context) {
	page, pageSize := GetPagination(c)

	var stationID *int64
	if sid, err := parseInt64Param(c.Query("station_id")); err == nil {
		stationID = &sid
	}

	banners, total, err := h.service.List(page, pageSize, stationID)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "failed to list banners")
		return
	}

	RespondPage(c, http.StatusOK, "ok", banners, page, pageSize, total)
}

// Create 创建 Banner
// @Summary      创建轮播图
// @Description  创建新的轮播图
// @Tags         b_banner
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        request body bannerRequest true "轮播图信息"
// @Success      200  {object} APIResponse "创建成功"
// @Failure      400  {object} APIResponse "请求参数错误"
// @Router       /b/banners [post]
func (h *BannerHandler) Create(c *gin.Context) {
	var req bannerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "invalid request")
		return
	}

	banner, err := h.service.Create(service.BannerInput{
		StationID: req.StationID,
		Title:     req.Title,
		ImageURL:  req.ImageURL,
		LinkType:  req.LinkType,
		LinkValue: req.LinkValue,
		Sort:      req.Sort,
		Status:    req.Status,
	})
	if err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	Respond(c, http.StatusOK, "ok", banner)
}

// Update 更新 Banner
// @Summary      更新轮播图
// @Description  更新指定轮播图信息
// @Tags         b_banner
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id path int true "轮播图ID"
// @Param        request body bannerRequest true "轮播图信息"
// @Success      200  {object} APIResponse "更新成功"
// @Failure      400  {object} APIResponse "请求参数错误"
// @Router       /b/banners/{id} [put]
func (h *BannerHandler) Update(c *gin.Context) {
	id, err := parseInt64Param(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	var req bannerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "invalid request")
		return
	}

	banner, err := h.service.Update(service.BannerInput{
		ID:        id,
		StationID: req.StationID,
		Title:     req.Title,
		ImageURL:  req.ImageURL,
		LinkType:  req.LinkType,
		LinkValue: req.LinkValue,
		Sort:      req.Sort,
		Status:    req.Status,
	})
	if err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	Respond(c, http.StatusOK, "ok", banner)
}

// Delete 删除 Banner
// @Summary      删除轮播图
// @Description  删除指定轮播图
// @Tags         b_banner
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id path int true "轮播图ID"
// @Success      200  {object} APIResponse "删除成功"
// @Router       /b/banners/{id} [delete]
func (h *BannerHandler) Delete(c *gin.Context) {
	id, err := parseInt64Param(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.service.Delete(id); err != nil {
		RespondError(c, http.StatusInternalServerError, "failed to delete banner")
		return
	}

	Respond(c, http.StatusOK, "ok", nil)
}
