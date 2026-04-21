package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/service"

	"github.com/gin-gonic/gin"
)

type NewsHandler struct {
	newsService *service.NewsService
}

// NewNewsHandler 创建 NewsHandler
func NewNewsHandler(newsService *service.NewsService) *NewsHandler {
	return &NewsHandler{newsService: newsService}
}

// newsRequest 新闻请求结构
type newsRequest struct {
	StationID int64  `json:"station_id"`
	Title     string `json:"title" binding:"required"`
	Summary   string `json:"summary"`
	Content   string `json:"content"`
	CoverURL  string `json:"cover_url"`
	Type      string `json:"type"`
	Status    string `json:"status"`
}

// List godoc
// @Summary 获取新闻列表
// @Description 获取已发布的新闻列表（分页）
// @Tags C端-新闻
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param type query string false "新闻类型(news/notice/activity)"
// @Success 200 {object} APIResponse
// @Router /api/v1/c/news [get]
func (h *NewsHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	newsType := c.Query("type")
	status := c.Query("status") // B端筛选状态

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 10
	}

	// 获取站点ID，优先从 Query 参数获取，其次从 Context 获取（JWT）
	var stationID *int64
	if sid, err := strconv.ParseInt(c.Query("station_id"), 10, 64); err == nil && sid >= 0 {
		stationID = &sid
	} else if id, exists := GetStationID(c); exists {
		stationID = &id
	}

	var news []*model.News
	var total int64
	var err error

	// 根据路由判断是 B端还是 C端
	isBEnd := strings.HasPrefix(c.Request.URL.Path, "/api/v1/b")

	if isBEnd {
		// B端管理列表：包含草稿，有站点名称
		news, total, err = h.newsService.List(page, pageSize, newsType, status, stationID)
	} else {
		// C端公开列表：仅发布，按站点+时间排序
		news, total, err = h.newsService.ListPublished(page, pageSize, newsType, stationID)
	}

	if err != nil {
		Respond(c, http.StatusInternalServerError, "获取新闻列表失败", nil)
		return
	}

	Respond(c, http.StatusOK, "ok", gin.H{
		"items":     news,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// Get godoc
// @Summary 获取新闻详情
// @Description 获取单条新闻的详细内容
// @Tags C端-新闻
// @Accept json
// @Produce json
// @Param id path int true "新闻ID"
// @Success 200 {object} Response
// @Router /api/v1/c/news/{id} [get]
func (h *NewsHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		Respond(c, http.StatusBadRequest, "无效的新闻ID", nil)
		return
	}

	news, err := h.newsService.GetByID(id)
	if err != nil {
		Respond(c, http.StatusNotFound, "新闻不存在", nil)
		return
	}

	Respond(c, http.StatusOK, "ok", news)
}

// Create godoc
// @Summary 创建新闻
// @Description 创建新的新闻/公告
// @Tags B端-新闻管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body newsRequest true "新闻信息"
// @Success 200 {object} Response
// @Router /api/v1/b/news [post]
func (h *NewsHandler) Create(c *gin.Context) {
	var req newsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	userID, _ := GetUserID(c)

	news := &model.News{
		StationID: req.StationID,
		Title:     req.Title,
		Summary:   req.Summary,
		Content:   req.Content,
		CoverURL:  req.CoverURL,
		Type:      req.Type,
		Status:    req.Status,
		AuthorID:  userID,
	}

	if news.Type == "" {
		news.Type = "news"
	}
	if news.Status == "" {
		news.Status = "draft"
	}
	if news.Status == "published" {
		now := time.Now()
		news.PublishAt = now
	}

	if err := h.newsService.Create(news); err != nil {
		RespondError(c, http.StatusInternalServerError, "创建新闻失败")
		return
	}

	Respond(c, http.StatusOK, "创建成功", news)
}

// Update godoc
// @Summary 更新新闻
// @Description 更新新闻信息
// @Tags B端-新闻管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "新闻ID"
// @Param request body newsRequest true "新闻信息"
// @Success 200 {object} Response
// @Router /api/v1/b/news/{id} [put]
func (h *NewsHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "无效的新闻ID")
		return
	}

	var req newsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	// 获取现有新闻
	existing, err := h.newsService.GetByID(id)
	if err != nil {
		RespondError(c, http.StatusNotFound, "新闻不存在")
		return
	}

	// 更新字段
	existing.StationID = req.StationID
	existing.Title = req.Title
	existing.Summary = req.Summary
	existing.Content = req.Content
	existing.CoverURL = req.CoverURL
	if req.Type != "" {
		existing.Type = req.Type
	}
	if req.Status != "" {
		// 如果从非发布状态变为发布状态，设置发布时间
		if existing.Status != "published" && req.Status == "published" {
			now := time.Now()
			existing.PublishAt = now
		}
		existing.Status = req.Status
	}

	if err := h.newsService.Update(existing); err != nil {
		RespondError(c, http.StatusInternalServerError, "更新新闻失败")
		return
	}

	Respond(c, http.StatusOK, "更新成功", existing)
}

// Delete godoc
// @Summary 删除新闻
// @Description 删除新闻（软删除）
// @Tags B端-新闻管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "新闻ID"
// @Success 200 {object} Response
// @Router /api/v1/b/news/{id} [delete]
func (h *NewsHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "无效的新闻ID")
		return
	}

	if err := h.newsService.Delete(id); err != nil {
		RespondError(c, http.StatusInternalServerError, "删除新闻失败")
		return
	}

	Respond(c, http.StatusOK, "删除成功", nil)
}
