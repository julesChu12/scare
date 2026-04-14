package handler

import (
	"net/http"

	"community-elderly-care-platform/internal/service"

	"github.com/gin-gonic/gin"
)

type RequestHandler struct {
	service *service.RequestService
}

func NewRequestHandler(service *service.RequestService) *RequestHandler {
	return &RequestHandler{service: service}
}

type requestCreate struct {
	RequestNo       string   `json:"request_no"`
	ServiceType     string   `json:"service_type" binding:"required"`
	Lat             *float64 `json:"lat"` // 兼容旧字段，等同 submit_lat
	Lng             *float64 `json:"lng"` // 兼容旧字段，等同 submit_lng
	SubmitLat       *float64 `json:"submit_lat"`
	SubmitLng       *float64 `json:"submit_lng"`
	ServiceLat      *float64 `json:"service_lat"`
	ServiceLng      *float64 `json:"service_lng"`
	SourceStationID *int64   `json:"source_station_id"`
	ContactName     string   `json:"contact_name"`
	ContactPhone    string   `json:"contact_phone"`
	Address         string   `json:"address"`
	Images          []string `json:"images"`
}

// Create 创建服务需求
// @Summary      创建服务需求
// @Description  C端用户创建新的服务需求，系统自动匹配站点并派发任务
// @Tags         c_request
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        request body requestCreate true "服务需求信息"
// @Success      200  {object} APIResponse{data=RequestResponse} "创建成功"
// @Failure      400  {object} APIResponse "请求参数错误"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      409  {object} APIResponse "请求冲突"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /c/requests [post]
func (h *RequestHandler) Create(c *gin.Context) {
	userID, ok := GetUserID(c)
	if !ok {
		RespondError(c, http.StatusUnauthorized, "missing user")
		return
	}

	var req requestCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "invalid payload")
		return
	}

	submitLat := req.SubmitLat
	submitLng := req.SubmitLng
	if submitLat == nil && submitLng == nil {
		submitLat = req.Lat
		submitLng = req.Lng
	}
	serviceLat := req.ServiceLat
	serviceLng := req.ServiceLng
	if serviceLat == nil && serviceLng == nil && req.Lat != nil && req.Lng != nil {
		serviceLat = req.Lat
		serviceLng = req.Lng
	}

	request, created, err := h.service.Create(service.RequestInput{
		UserID:          userID,
		RequestNo:       req.RequestNo,
		ServiceType:     req.ServiceType,
		SubmitLat:       submitLat,
		SubmitLng:       submitLng,
		ServiceLat:      serviceLat,
		ServiceLng:      serviceLng,
		SourceStationID: req.SourceStationID,
		ContactName:     req.ContactName,
		ContactPhone:    req.ContactPhone,
		Address:         req.Address,
		Images:          req.Images,
	})
	if err != nil {
		switch err {
		case service.ErrAddressRequired:
			RespondError(c, http.StatusBadRequest, "address required")
		case service.ErrInvalidRequest:
			RespondError(c, http.StatusBadRequest, "invalid request")
		case service.ErrServiceLocationRequired:
			RespondError(c, http.StatusBadRequest, "service location required")
		case service.ErrRequestConflict:
			RespondError(c, http.StatusConflict, "request conflict")
		case service.ErrNoStation:
			RespondError(c, http.StatusConflict, "no station available")
		default:
			RespondError(c, http.StatusInternalServerError, "create request failed")
		}
		return
	}

	if !created {
		Respond(c, http.StatusOK, "ok", request)
		return
	}
	Respond(c, http.StatusOK, "ok", request)
}

// List 获取服务需求列表
// @Summary      获取服务需求列表
// @Description  获取当前用户的服务需求列表，支持按状态筛选和分页
// @Tags         c_request,b_request
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        status query string false "状态筛选"
// @Param        page query int false "页码" default(1)
// @Param        page_size query int false "每页数量"
// @Success      200  {object} APIResponse{data=RequestResponse} "获取成功"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /c/requests [get]
// @Router       /b/requests [get]
func (h *RequestHandler) List(c *gin.Context) {
	userType, _ := GetUserType(c)
	status := c.Query("status")
	page, pageSize := GetPagination(c)

	if userType == "b_end" {
		// B端：查询所有请求，支持站点筛选
		// 安全修复：非管理员只能查看自己站点的数据
		userStationID, _ := GetStationID(c)
		roles := GetUserRoles(c)
		isAdmin := containsRole(roles, "admin")

		queryStationID, _ := parseInt64Param(c.Query("station_id"))

		var stationID int64
		if isAdmin {
			// 管理员可以指定站点或查看所有（0）
			stationID = queryStationID
		} else {
			// 非管理员强制使用自己的站点
			if userStationID == 0 {
				RespondError(c, http.StatusBadRequest, "missing station")
				return
			}
			stationID = userStationID
		}

		requests, total, err := h.service.ListAll(stationID, status, page, pageSize)
		if err != nil {
			RespondError(c, http.StatusInternalServerError, "list requests failed")
			return
		}
		RespondPage(c, http.StatusOK, "ok", requests, page, pageSize, total)
		return
	}

	// C端：只查询当前用户的请求
	userID, ok := GetUserID(c)
	if !ok {
		RespondError(c, http.StatusUnauthorized, "missing user")
		return
	}
	requests, total, err := h.service.ListByUser(userID, status, page, pageSize)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "list requests failed")
		return
	}
	RespondPage(c, http.StatusOK, "ok", requests, page, pageSize, total)
}

// Get 获取服务需求详情
// @Summary      获取服务需求详情
// @Description  根据ID获取服务需求的详细信息
// @Tags         c_request,b_request
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id path int true "需求ID"
// @Success      200  {object} APIResponse{data=RequestResponse} "获取成功"
// @Failure      400  {object} APIResponse "ID无效"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      403  {object} APIResponse "无权限"
// @Failure      404  {object} APIResponse "需求不存在"
// @Router       /c/requests/{id} [get]
// @Router       /b/requests/{id} [get]
func (h *RequestHandler) Get(c *gin.Context) {
	id, err := parseInt64Param(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	request, err := h.service.GetByID(id)
	if err != nil {
		RespondError(c, http.StatusNotFound, "request not found")
		return
	}
	userID, ok := GetUserID(c)
	if !ok {
		RespondError(c, http.StatusUnauthorized, "missing user")
		return
	}
	userType, _ := GetUserType(c)
	if userType == "c_end" && request.UserID != userID {
		RespondError(c, http.StatusForbidden, "forbidden")
		return
	}
	// B端：非管理员需验证站点归属
	if userType == "b_end" {
		roles := GetUserRoles(c)
		if !containsRole(roles, "admin") {
			userStationID, _ := GetStationID(c)
			// StationID > 0 表示已分配站点，需验证归属
			if request.StationID > 0 && request.StationID != userStationID {
				RespondError(c, http.StatusForbidden, "forbidden")
				return
			}
		}
	}
	Respond(c, http.StatusOK, "ok", request)
}

// Cancel 取消服务需求
// @Summary      取消服务需求
// @Description  取消待处理的服务需求
// @Tags         c_request
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id path int true "需求ID"
// @Success      200  {object} APIResponse{data=RequestResponse} "取消成功"
// @Failure      400  {object} APIResponse "ID无效"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      409  {object} APIResponse "状态不允许取消"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /c/requests/{id}/cancel [post]
func (h *RequestHandler) Cancel(c *gin.Context) {
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
	request, changed, err := h.service.Cancel(id, userID)
	if err != nil {
		switch err {
		case service.ErrRequestConflict:
			RespondError(c, http.StatusConflict, "request conflict")
		default:
			RespondError(c, http.StatusInternalServerError, "cancel failed")
		}
		return
	}
	if !changed {
		Respond(c, http.StatusOK, "ok", request)
		return
	}
	Respond(c, http.StatusOK, "ok", request)
}

type updateRequestBody struct {
	ServiceType  string `json:"service_type"`
	ContactName  string `json:"contact_name"`
	ContactPhone string `json:"contact_phone"`
	Address      string `json:"address"`
	Description  string `json:"description"`
	Urgency      string `json:"urgency"`
	StationID   *int64 `json:"station_id"`
}

// Update B端编辑服务请求
// @Summary      编辑服务请求
// @Description  B端编辑服务请求基本信息（仅 pending/dispatched 状态可编辑）
// @Tags         b_request
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id path int true "需求ID"
// @Param        request body updateRequestBody true "更新信息"
// @Success      200  {object} APIResponse{data=RequestResponse} "更新成功"
// @Failure      400  {object} APIResponse "请求参数错误"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      409  {object} APIResponse "状态不允许编辑"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /b/requests/{id} [put]
func (h *RequestHandler) Update(c *gin.Context) {
	id, err := parseInt64Param(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	var req updateRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "invalid payload")
		return
	}

	result, err := h.service.Update(id, service.UpdateInput{
		ServiceType:  req.ServiceType,
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		Address:      req.Address,
		Description:  req.Description,
		Urgency:      req.Urgency,
		StationID:    req.StationID,
	})
	if err != nil {
		switch err {
		case service.ErrInvalidRequest:
			RespondError(c, http.StatusBadRequest, "invalid request")
		case service.ErrRequestConflict:
			RespondError(c, http.StatusConflict, "当前状态不允许编辑")
		default:
			RespondError(c, http.StatusInternalServerError, "update failed")
		}
		return
	}

	Respond(c, http.StatusOK, "ok", result)
}

// CancelByAdmin B端取消服务请求
// @Summary      B端取消服务请求
// @Description  B端管理人员取消服务请求，同步取消关联任务
// @Tags         b_request
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id path int true "需求ID"
// @Success      200  {object} APIResponse{data=RequestResponse} "取消成功"
// @Failure      400  {object} APIResponse "ID无效"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      409  {object} APIResponse "状态不允许取消"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /b/requests/{id}/cancel [post]
func (h *RequestHandler) CancelByAdmin(c *gin.Context) {
	id, err := parseInt64Param(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	result, err := h.service.CancelByAdmin(id)
	if err != nil {
		switch err {
		case service.ErrRequestConflict:
			RespondError(c, http.StatusConflict, "当前状态不允许取消")
		default:
			RespondError(c, http.StatusInternalServerError, "cancel failed")
		}
		return
	}

	Respond(c, http.StatusOK, "ok", result)
}

type updateStatusRequest struct {
	Status       string `json:"status" binding:"required"`
	RejectReason string `json:"reject_reason"`
}

// UpdateStatus B端更新服务请求状态
// @Summary      更新服务请求状态
// @Description  B端更新服务请求状态（如拒绝、完成等）
// @Tags         b_request
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id path int true "需求ID"
// @Param        request body updateStatusRequest true "状态信息"
// @Success      200  {object} APIResponse{data=RequestResponse} "更新成功"
// @Failure      400  {object} APIResponse "请求参数错误"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /b/requests/{id}/status [put]
func (h *RequestHandler) UpdateStatus(c *gin.Context) {
	id, err := parseInt64Param(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	var req updateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "invalid payload")
		return
	}

	if err := h.service.UpdateStatus(id, req.Status, req.RejectReason); err != nil {
		if err == service.ErrInvalidRequest {
			RespondError(c, http.StatusBadRequest, "invalid status")
			return
		}
		RespondError(c, http.StatusInternalServerError, "update failed")
		return
	}

	Respond(c, http.StatusOK, "ok", nil)
}

type rateRequest struct {
	Rating   int    `json:"rating" binding:"required,min=1,max=5"`
	Feedback string `json:"feedback"`
}

// Rate 评价服务
// @Summary      评价服务
// @Description  对已完成的服务进行评价
// @Tags         c_request
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id path int true "需求ID"
// @Param        request body rateRequest true "评价信息"
// @Success      200  {object} APIResponse{data=RequestResponse} "评价成功"
// @Failure      400  {object} APIResponse "请求参数错误"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      403  {object} APIResponse "无权限"
// @Failure      409  {object} APIResponse "状态不允许评价"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /c/requests/{id}/rate [post]
func (h *RequestHandler) Rate(c *gin.Context) {
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

	var req rateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "invalid payload")
		return
	}

	request, err := h.service.Rate(id, userID, req.Rating, req.Feedback)
	if err != nil {
		switch err {
		case service.ErrInvalidRequest:
			RespondError(c, http.StatusBadRequest, "invalid rating")
		case service.ErrRequestConflict:
			RespondError(c, http.StatusForbidden, "forbidden")
		case service.ErrNotCompleted:
			RespondError(c, http.StatusConflict, "request not completed")
		case service.ErrAlreadyRated:
			RespondError(c, http.StatusConflict, "already rated")
		default:
			RespondError(c, http.StatusInternalServerError, "rate failed")
		}
		return
	}

	Respond(c, http.StatusOK, "ok", request)
}
