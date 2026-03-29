package service

import (
	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/repository"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

const (
	ReportTypeService     = "service"
	ReportTypePerformance = "performance"
	ReportTypeRequest     = "request"
	ReportTypeStation     = "station"

	ReportFormatXLSX = "xlsx"
	ReportFormatCSV  = "csv"
)

type ReportService struct {
	reportRepo  *repository.ReportRepository
	requestRepo *repository.RequestRepository
	taskRepo    *repository.TaskRepository
	stationRepo *repository.StationRepository
	userRepo    *repository.UserRepository
	storagePath string
}

func NewReportService(
	reportRepo *repository.ReportRepository,
	requestRepo *repository.RequestRepository,
	taskRepo *repository.TaskRepository,
	stationRepo *repository.StationRepository,
	userRepo *repository.UserRepository,
	storagePath string,
) *ReportService {
	return &ReportService{
		reportRepo:  reportRepo,
		requestRepo: requestRepo,
		taskRepo:    taskRepo,
		stationRepo: stationRepo,
		userRepo:    userRepo,
		storagePath: storagePath,
	}
}

type GenerateReportInput struct {
	Type      string
	Format    string
	StationID int64
	StartDate time.Time
	EndDate   time.Time
	UserID    int64
}

type GenerateReportOutput struct {
	Report   *model.Report
	FilePath string
}

type ReportPreviewData struct {
	RequestCount          int64   `json:"request_count"`
	CompletedRequestCount int64   `json:"completed_request_count"`
	CompletionRate        float64 `json:"completion_rate"`
	ServiceTypeCount      int     `json:"service_type_count"`
	RankedStaffCount      int     `json:"ranked_staff_count"`
	CompletedTaskCount    int64   `json:"completed_task_count"`
	AvgRating             float64 `json:"avg_rating"`
	TrendDays             int     `json:"trend_days"`
	StationCount          int64   `json:"station_count"`
}

func (s *ReportService) GenerateReport(input GenerateReportInput) (*GenerateReportOutput, error) {
	reportName := s.generateReportName(input.Type, input.StartDate, input.EndDate)
	filePath := s.generateFilePath(input.Type, input.Format)

	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return nil, fmt.Errorf("create directory failed: %w", err)
	}

	var err error
	switch input.Format {
	case ReportFormatXLSX:
		err = s.generateExcel(input, filePath)
	case ReportFormatCSV:
		err = s.generateCSV(input, filePath)
	default:
		return nil, fmt.Errorf("unsupported format: %s", input.Format)
	}

	if err != nil {
		return nil, fmt.Errorf("generate report failed: %w", err)
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("get file info failed: %w", err)
	}

	report := &model.Report{
		Name:      reportName,
		Type:      input.Type,
		Format:    input.Format,
		FilePath:  filePath,
		FileSize:  fileInfo.Size(),
		StationID: input.StationID,
		StartDate: input.StartDate,
		EndDate:   input.EndDate,
		CreatedBy: input.UserID,
	}

	if err := s.reportRepo.Create(report); err != nil {
		os.Remove(filePath)
		return nil, fmt.Errorf("save report record failed: %w", err)
	}

	go s.cleanupOldReports()

	return &GenerateReportOutput{
		Report:   report,
		FilePath: filePath,
	}, nil
}

func (s *ReportService) GetReport(id int64) (*model.Report, error) {
	return s.reportRepo.GetByID(id)
}

func (s *ReportService) GetPreview(input GenerateReportInput) (*ReportPreviewData, error) {
	isAdmin := input.StationID == 0
	totalRequests, err := s.requestRepo.CountBetween(input.StationID, isAdmin, input.StartDate, input.EndDate)
	if err != nil {
		return nil, err
	}
	completedRequests, err := s.requestRepo.CountByStatusBetween(input.StationID, "completed", isAdmin, input.StartDate, input.EndDate)
	if err != nil {
		return nil, err
	}

	preview := &ReportPreviewData{
		RequestCount:          totalRequests,
		CompletedRequestCount: completedRequests,
		CompletionRate:        calculateCompletionRate(totalRequests, completedRequests),
		TrendDays:             int(input.EndDate.Sub(input.StartDate).Hours()/24) + 1,
	}

	switch input.Type {
	case ReportTypeService:
		typeCounts, err := s.requestRepo.CountByServiceTypeBetween(input.StationID, isAdmin, input.StartDate, input.EndDate)
		if err != nil {
			return nil, err
		}
		preview.ServiceTypeCount = len(typeCounts)
	case ReportTypePerformance:
		ranking, err := s.taskRepo.GetStaffRankingBetween(input.StationID, isAdmin, input.StartDate, input.EndDate, 50)
		if err != nil {
			return nil, err
		}
		preview.RankedStaffCount = len(ranking)
		var totalCompleted int64
		var weightedRating float64
		for _, item := range ranking {
			totalCompleted += item.CompletedCount
			weightedRating += item.AvgRating * float64(item.CompletedCount)
		}
		preview.CompletedTaskCount = totalCompleted
		if totalCompleted > 0 {
			preview.AvgRating = weightedRating / float64(totalCompleted)
		}
	case ReportTypeStation:
		stations, err := s.listReportStations(input)
		if err != nil {
			return nil, err
		}
		preview.StationCount = int64(len(stations))
	}

	return preview, nil
}

func (s *ReportService) ListReports(filter repository.ReportFilter, page, pageSize int) ([]*model.Report, int64, error) {
	offset := (page - 1) * pageSize
	return s.reportRepo.List(filter, offset, pageSize)
}

func (s *ReportService) DeleteReport(id int64) error {
	report, err := s.reportRepo.GetByID(id)
	if err != nil {
		return err
	}

	if err := s.reportRepo.Delete(id); err != nil {
		return err
	}

	os.Remove(report.FilePath)
	return nil
}

func (s *ReportService) cleanupOldReports() {
	filePaths, err := s.reportRepo.DeleteOlderThan(30)
	if err != nil {
		return
	}

	for _, fp := range filePaths {
		os.Remove(fp)
	}
}

func (s *ReportService) generateReportName(reportType string, startDate, endDate time.Time) string {
	typeNames := map[string]string{
		ReportTypeService:     "服务统计报表",
		ReportTypePerformance: "人员绩效报表",
		ReportTypeRequest:     "需求分析报表",
		ReportTypeStation:     "站点运营报表",
	}

	typeName := typeNames[reportType]
	if typeName == "" {
		typeName = "统计报表"
	}

	return fmt.Sprintf("%s_%s至%s",
		typeName,
		startDate.Format("20060102"),
		endDate.Format("20060102"),
	)
}

func (s *ReportService) generateFilePath(reportType, format string) string {
	now := time.Now()
	dir := filepath.Join(s.storagePath, "reports", now.Format("2006"), now.Format("01"))
	filename := fmt.Sprintf("%s-%s-%s.%s", uuid.New().String()[:8], reportType, now.Format("20060102"), format)
	return filepath.Join(dir, filename)
}

func (s *ReportService) generateExcel(input GenerateReportInput, filePath string) error {
	f := excelize.NewFile()
	defer f.Close()

	switch input.Type {
	case ReportTypeService:
		if err := s.generateServiceReport(f, input); err != nil {
			return err
		}
	case ReportTypePerformance:
		if err := s.generatePerformanceReport(f, input); err != nil {
			return err
		}
	case ReportTypeRequest:
		if err := s.generateRequestReport(f, input); err != nil {
			return err
		}
	case ReportTypeStation:
		if err := s.generateStationReport(f, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported report type: %s", input.Type)
	}

	return f.SaveAs(filePath)
}

func (s *ReportService) generateCSV(input GenerateReportInput, filePath string) error {
	f := excelize.NewFile()
	defer f.Close()

	switch input.Type {
	case ReportTypeService:
		if err := s.generateServiceReport(f, input); err != nil {
			return err
		}
	case ReportTypePerformance:
		if err := s.generatePerformanceReport(f, input); err != nil {
			return err
		}
	case ReportTypeRequest:
		if err := s.generateRequestReport(f, input); err != nil {
			return err
		}
	case ReportTypeStation:
		if err := s.generateStationReport(f, input); err != nil {
			return err
		}
	}

	tmpXlsx := filePath + ".xlsx"
	if err := f.SaveAs(tmpXlsx); err != nil {
		return err
	}
	defer os.Remove(tmpXlsx)

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return fmt.Errorf("no worksheet generated")
	}

	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	file.WriteString("\xEF\xBB\xBF")

	for _, row := range rows {
		for i, cell := range row {
			if i > 0 {
				file.WriteString(",")
			}
			file.WriteString(fmt.Sprintf("\"%s\"", cell))
		}
		file.WriteString("\n")
	}

	return nil
}

func (s *ReportService) generateServiceReport(f *excelize.File, input GenerateReportInput) error {
	sheetName := "服务统计"
	f.SetSheetName("Sheet1", sheetName)

	headers := []string{"指标", "数值"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, h)
	}

	isAdmin := input.StationID == 0
	totalRequests, _ := s.requestRepo.CountBetween(input.StationID, isAdmin, input.StartDate, input.EndDate)
	completedRequests, _ := s.requestRepo.CountByStatusBetween(input.StationID, "completed", isAdmin, input.StartDate, input.EndDate)
	pendingRequests, _ := s.requestRepo.CountByStatusBetween(input.StationID, "pending", isAdmin, input.StartDate, input.EndDate)
	completionRate := calculateCompletionRate(totalRequests, completedRequests)

	data := [][]interface{}{
		{"总需求数", totalRequests},
		{"已完成", completedRequests},
		{"待处理", pendingRequests},
		{"完成率", fmt.Sprintf("%.1f%%", completionRate)},
	}

	for i, row := range data {
		for j, val := range row {
			cell, _ := excelize.CoordinatesToCellName(j+1, i+2)
			f.SetCellValue(sheetName, cell, val)
		}
	}

	f.NewSheet("服务类型分布")
	typeHeaders := []string{"服务类型", "数量", "占比"}
	for i, h := range typeHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue("服务类型分布", cell, h)
	}

	typeCounts, _ := s.requestRepo.CountByServiceTypeBetween(input.StationID, isAdmin, input.StartDate, input.EndDate)
	row := 2
	for serviceType, count := range typeCounts {
		percentage := calculateCompletionRate(totalRequests, count)
		f.SetCellValue("服务类型分布", fmt.Sprintf("A%d", row), serviceType)
		f.SetCellValue("服务类型分布", fmt.Sprintf("B%d", row), count)
		f.SetCellValue("服务类型分布", fmt.Sprintf("C%d", row), fmt.Sprintf("%.1f%%", percentage))
		row++
	}

	return nil
}

func (s *ReportService) generatePerformanceReport(f *excelize.File, input GenerateReportInput) error {
	sheetName := "人员绩效"
	f.SetSheetName("Sheet1", sheetName)

	headers := []string{"排名", "姓名", "完成数", "平均评分"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, h)
	}

	isAdmin := input.StationID == 0
	ranking, _ := s.taskRepo.GetStaffRankingBetween(input.StationID, isAdmin, input.StartDate, input.EndDate, 50)

	for i, item := range ranking {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), i+1)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), item.Name)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), item.CompletedCount)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), fmt.Sprintf("%.1f", item.AvgRating))
	}

	return nil
}

func (s *ReportService) generateRequestReport(f *excelize.File, input GenerateReportInput) error {
	sheetName := "需求分析"
	f.SetSheetName("Sheet1", sheetName)

	headers := []string{"日期", "新增需求数"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, h)
	}

	isAdmin := input.StationID == 0
	trend, _ := s.requestRepo.GetDailyTrendBetween(input.StationID, isAdmin, input.StartDate, input.EndDate)

	for i, item := range trend {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), item.Date)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), item.Count)
	}

	f.NewSheet("状态分布")
	statusHeaders := []string{"状态", "数量"}
	for i, h := range statusHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue("状态分布", cell, h)
	}

	statuses := []string{"pending", "dispatched", "processing", "completed", "cancelled"}
	statusNames := map[string]string{
		"pending":    "待处理",
		"dispatched": "已派发",
		"processing": "处理中",
		"completed":  "已完成",
		"cancelled":  "已取消",
	}

	row := 2
	for _, status := range statuses {
		count, _ := s.requestRepo.CountByStatusBetween(input.StationID, status, isAdmin, input.StartDate, input.EndDate)
		f.SetCellValue("状态分布", fmt.Sprintf("A%d", row), statusNames[status])
		f.SetCellValue("状态分布", fmt.Sprintf("B%d", row), count)
		row++
	}

	return nil
}

func (s *ReportService) generateStationReport(f *excelize.File, input GenerateReportInput) error {
	sheetName := "站点运营"
	f.SetSheetName("Sheet1", sheetName)

	headers := []string{"站点名称", "需求数", "完成数", "完成率"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, h)
	}

	stations, err := s.listReportStations(input)
	if err != nil {
		return err
	}

	row := 2
	for _, station := range stations {
		total, _ := s.requestRepo.CountBetween(station.ID, false, input.StartDate, input.EndDate)
		completed, _ := s.requestRepo.CountByStatusBetween(station.ID, "completed", false, input.StartDate, input.EndDate)
		completionRate := calculateCompletionRate(total, completed)

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), station.Name)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), total)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), completed)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), fmt.Sprintf("%.1f%%", completionRate))
		row++
	}

	return nil
}

func (s *ReportService) listReportStations(input GenerateReportInput) ([]*model.ServiceStation, error) {
	if input.StationID > 0 {
		station, err := s.stationRepo.GetByID(input.StationID)
		if err != nil {
			return nil, err
		}
		return []*model.ServiceStation{station}, nil
	}

	stations, _, err := s.stationRepo.List(0, 100, repository.StationListFilter{})
	if err != nil {
		return nil, err
	}
	return stations, nil
}

func calculateCompletionRate(total, completed int64) float64 {
	if total == 0 {
		return 0
	}
	return float64(completed) * 100 / float64(total)
}
