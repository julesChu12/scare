package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"community-elderly-care-platform/internal/repository"
	"community-elderly-care-platform/internal/service"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	service *service.ReportService
}

// NewReportHandler 创建 ReportHandler
func NewReportHandler(service *service.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

type generateReportRequest struct {
	Type      string `json:"type" binding:"required,oneof=service performance request station"`
	Format    string `json:"format" binding:"required,oneof=xlsx csv"`
	StationID *int64 `json:"station_id"`
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
	Preview   bool   `json:"preview"`
}

// GenerateReport 生成并下载报表
// @Summary      生成并下载报表
// @Description  生成指定类型的统计报表并直接下载
// @Tags         b_reports
// @Accept       json
// @Produce      application/octet-stream
// @Security     Bearer
// @Param        request body generateReportRequest true "报表参数"
// @Success      200  {file} file "报表文件"
// @Failure      400  {object} APIResponse "参数错误"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      500  {object} APIResponse "生成失败"
// @Router       /b/reports/generate [post]
func (h *ReportHandler) GenerateReport(c *gin.Context) {
	var req generateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid start_date format")
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid end_date format")
		return
	}

	if endDate.Before(startDate) {
		RespondError(c, http.StatusBadRequest, "end_date must be after start_date")
		return
	}

	userID, _ := GetUserID(c)
	stationID, _ := GetStationID(c)
	roles := GetUserRoles(c)
	isAdmin := containsRole(roles, "admin")

	var reportStationID int64
	if isAdmin && req.StationID != nil {
		reportStationID = *req.StationID
	} else if !isAdmin {
		reportStationID = stationID
	}

	input := service.GenerateReportInput{
		Type:      req.Type,
		Format:    req.Format,
		StationID: reportStationID,
		StartDate: startDate,
		EndDate:   endDate,
		UserID:    userID,
	}

	if req.Preview {
		preview, err := h.service.GetPreview(input)
		if err != nil {
			RespondError(c, http.StatusInternalServerError, "preview report failed: "+err.Error())
			return
		}
		Respond(c, http.StatusOK, "ok", preview)
		return
	}

	output, err := h.service.GenerateReport(input)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "generate report failed: "+err.Error())
		return
	}

	filename := output.Report.Name + "." + output.Report.Format
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
	c.Header("Content-Type", "application/octet-stream")
	c.File(output.FilePath)
}

// ListReports 获取历史报表列表
// @Summary      获取历史报表列表
// @Description  获取当前用户可见的历史报表列表
// @Tags         b_reports
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        page query int false "页码"
// @Param        page_size query int false "每页数量"
// @Param        type query string false "报表类型筛选"
// @Success      200  {object} APIResponse "获取成功"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /b/reports [get]
func (h *ReportHandler) ListReports(c *gin.Context) {
	page, pageSize := GetPagination(c)

	userID, _ := GetUserID(c)
	stationID, _ := GetStationID(c)
	roles := GetUserRoles(c)
	isAdmin := containsRole(roles, "admin")

	filter := repository.ReportFilter{
		Type:      c.Query("type"),
		StationID: stationID,
		CreatedBy: userID,
		IsAdmin:   isAdmin,
	}

	reports, total, err := h.service.ListReports(filter, page, pageSize)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "list reports failed")
		return
	}

	RespondPage(c, http.StatusOK, "ok", reports, page, pageSize, total)
}

// DownloadReport 下载历史报表
// @Summary      下载历史报表
// @Description  根据ID下载历史报表文件
// @Tags         b_reports
// @Accept       json
// @Produce      application/octet-stream
// @Security     Bearer
// @Param        id path int true "报表ID"
// @Success      200  {file} file "报表文件"
// @Failure      400  {object} APIResponse "参数错误"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      403  {object} APIResponse "无权限"
// @Failure      404  {object} APIResponse "报表不存在"
// @Router       /b/reports/{id}/download [get]
func (h *ReportHandler) DownloadReport(c *gin.Context) {
	id, err := parseInt64Param(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid report id")
		return
	}

	report, err := h.service.GetReport(id)
	if err != nil {
		RespondError(c, http.StatusNotFound, "report not found")
		return
	}

	userID, _ := GetUserID(c)
	stationID, _ := GetStationID(c)
	roles := GetUserRoles(c)
	isAdmin := containsRole(roles, "admin")

	if !isAdmin && report.CreatedBy != userID && report.StationID != stationID {
		RespondError(c, http.StatusForbidden, "no permission")
		return
	}

	if _, err := os.Stat(report.FilePath); os.IsNotExist(err) {
		RespondError(c, http.StatusNotFound, "file not found")
		return
	}

	filename := filepath.Base(report.FilePath)
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
	c.Header("Content-Type", "application/octet-stream")
	c.File(report.FilePath)
}

// DeleteReport 删除报表
// @Summary      删除报表
// @Description  删除指定的历史报表
// @Tags         b_reports
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id path int true "报表ID"
// @Success      200  {object} APIResponse "删除成功"
// @Failure      400  {object} APIResponse "参数错误"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      403  {object} APIResponse "无权限"
// @Failure      404  {object} APIResponse "报表不存在"
// @Router       /b/reports/{id} [delete]
func (h *ReportHandler) DeleteReport(c *gin.Context) {
	id, err := parseInt64Param(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid report id")
		return
	}

	report, err := h.service.GetReport(id)
	if err != nil {
		RespondError(c, http.StatusNotFound, "report not found")
		return
	}

	userID, _ := GetUserID(c)
	roles := GetUserRoles(c)
	isAdmin := containsRole(roles, "admin")

	if !isAdmin && report.CreatedBy != userID {
		RespondError(c, http.StatusForbidden, "no permission")
		return
	}

	if err := h.service.DeleteReport(id); err != nil {
		RespondError(c, http.StatusInternalServerError, "delete report failed")
		return
	}

	Respond(c, http.StatusOK, "ok", nil)
}
