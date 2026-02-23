package handler

import (
	"net/http"
	"time"

	"community-elderly-care-platform/internal/service"

	"github.com/gin-gonic/gin"
)

// CustomerHandler B端老年人档案管理
type CustomerHandler struct {
	service *service.ElderlyService
}

func NewCustomerHandler(service *service.ElderlyService) *CustomerHandler {
	return &CustomerHandler{service: service}
}

type customerCreateRequest struct {
	Name            string `json:"name" binding:"required"`
	Phone           string `json:"phone" binding:"required"`
	Gender          string `json:"gender"`
	BirthDate       string `json:"birth_date"`
	IDCard          string `json:"id_card"`
	Address         string `json:"address"`
	StationID       int64  `json:"station_id"`
	HealthStatus    string `json:"health_status"`
	DisabilityLevel string `json:"disability_level"`
	MedicalHistory  string `json:"medical_history"`
	SpecialNeeds    string `json:"special_needs"`
}

type customerUpdateRequest struct {
	Name            string `json:"name"`
	Gender          string `json:"gender"`
	BirthDate       string `json:"birth_date"`
	IDCard          string `json:"id_card"`
	Address         string `json:"address"`
	StationID       int64  `json:"station_id"`
	HealthStatus    string `json:"health_status"`
	DisabilityLevel string `json:"disability_level"`
	MedicalHistory  string `json:"medical_history"`
	SpecialNeeds    string `json:"special_needs"`
}

// List 获取老人档案列表
// @Summary      获取老人档案列表
// @Description  分页查询老年人档案，支持关键词/站点/健康状况筛选
// @Tags         b_customer
// @Produce      json
// @Security     Bearer
// @Param        page query int false "页码" default(1)
// @Param        page_size query int false "每页数量" default(10)
// @Param        keyword query string false "搜索关键词(姓名/手机号)"
// @Param        station_id query int false "站点ID"
// @Param        health_status query string false "健康状况(good/normal/poor)"
// @Success      200  {object} Response "获取成功"
// @Router       /b/customers [get]
func (h *CustomerHandler) List(c *gin.Context) {
	page, pageSize := GetPagination(c)
	stationID, _ := parseInt64Param(c.Query("station_id"))

	filter := service.ElderlyFilter{
		Keyword:      c.Query("keyword"),
		StationID:    stationID,
		HealthStatus: c.Query("health_status"),
		Page:         page,
		PageSize:     pageSize,
	}

	items, total, err := h.service.List(filter)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "查询失败: "+err.Error())
		return
	}

	RespondPage(c, http.StatusOK, "ok", items, page, pageSize, total)
}

// Get 获取老人档案详情
// @Summary      获取老人档案详情
// @Description  根据用户ID获取老年人档案详细信息
// @Tags         b_customer
// @Produce      json
// @Security     Bearer
// @Param        id path int true "用户ID"
// @Success      200  {object} Response "获取成功"
// @Router       /b/customers/{id} [get]
func (h *CustomerHandler) Get(c *gin.Context) {
	id, err := parseInt64Param(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	info, err := h.service.GetByID(id)
	if err != nil {
		if err == service.ErrElderlyNotFound {
			RespondError(c, http.StatusNotFound, "老人档案不存在")
			return
		}
		RespondError(c, http.StatusInternalServerError, "查询失败: "+err.Error())
		return
	}

	Respond(c, http.StatusOK, "ok", info)
}

// Create 创建老人档案
// @Summary      创建老人档案
// @Description  B端为老人创建档案（同时创建用户账号和身份）
// @Tags         b_customer
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        request body customerCreateRequest true "老人信息"
// @Success      200  {object} Response "创建成功"
// @Router       /b/customers [post]
func (h *CustomerHandler) Create(c *gin.Context) {
	var req customerCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	var birthDate time.Time
	if req.BirthDate != "" {
		parsed, err := time.Parse("2006-01-02", req.BirthDate)
		if err != nil {
			RespondError(c, http.StatusBadRequest, "出生日期格式错误，应为 YYYY-MM-DD")
			return
		}
		birthDate = parsed
	}

	input := service.ElderlyInput{
		Name:            req.Name,
		Phone:           req.Phone,
		Gender:          req.Gender,
		BirthDate:       birthDate,
		IDCard:          req.IDCard,
		Address:         req.Address,
		StationID:       req.StationID,
		HealthStatus:    req.HealthStatus,
		DisabilityLevel: req.DisabilityLevel,
		MedicalHistory:  req.MedicalHistory,
		SpecialNeeds:    req.SpecialNeeds,
	}

	info, err := h.service.Create(input)
	if err != nil {
		if err == service.ErrPhoneDuplicate {
			RespondError(c, http.StatusBadRequest, "手机号已注册")
			return
		}
		RespondError(c, http.StatusInternalServerError, "创建失败: "+err.Error())
		return
	}

	Respond(c, http.StatusOK, "ok", info)
}

// Update 更新老人档案
// @Summary      更新老人档案
// @Description  更新老年人档案信息
// @Tags         b_customer
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id path int true "用户ID"
// @Param        request body customerUpdateRequest true "更新信息"
// @Success      200  {object} Response "更新成功"
// @Router       /b/customers/{id} [put]
func (h *CustomerHandler) Update(c *gin.Context) {
	id, err := parseInt64Param(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	var req customerUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	var birthDate time.Time
	if req.BirthDate != "" {
		parsed, err := time.Parse("2006-01-02", req.BirthDate)
		if err != nil {
			RespondError(c, http.StatusBadRequest, "出生日期格式错误，应为 YYYY-MM-DD")
			return
		}
		birthDate = parsed
	}

	input := service.ElderlyInput{
		Name:            req.Name,
		Gender:          req.Gender,
		BirthDate:       birthDate,
		IDCard:          req.IDCard,
		Address:         req.Address,
		StationID:       req.StationID,
		HealthStatus:    req.HealthStatus,
		DisabilityLevel: req.DisabilityLevel,
		MedicalHistory:  req.MedicalHistory,
		SpecialNeeds:    req.SpecialNeeds,
	}

	info, err := h.service.Update(id, input)
	if err != nil {
		if err == service.ErrElderlyNotFound {
			RespondError(c, http.StatusNotFound, "老人档案不存在")
			return
		}
		RespondError(c, http.StatusInternalServerError, "更新失败: "+err.Error())
		return
	}

	Respond(c, http.StatusOK, "ok", info)
}

// GetServiceRecords 获取老人服务记录
// @Summary      获取老人服务记录
// @Description  获取指定老人的历史服务记录
// @Tags         b_customer
// @Produce      json
// @Security     Bearer
// @Param        id path int true "用户ID"
// @Param        page query int false "页码" default(1)
// @Param        page_size query int false "每页数量" default(10)
// @Success      200  {object} Response "获取成功"
// @Router       /b/customers/{id}/service-records [get]
func (h *CustomerHandler) GetServiceRecords(c *gin.Context) {
	id, err := parseInt64Param(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	page, pageSize := GetPagination(c)

	items, total, err := h.service.GetServiceRecords(id, page, pageSize)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "查询失败: "+err.Error())
		return
	}

	RespondPage(c, http.StatusOK, "ok", items, page, pageSize, total)
}
