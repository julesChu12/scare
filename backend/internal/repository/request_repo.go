package repository

import (
	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/dao/query"
	"time"

	"gorm.io/gorm"
)

type RequestRepository struct {
	q *query.Query
}

func NewRequestRepository(db *gorm.DB) *RequestRepository {
	return &RequestRepository{
		q: query.Use(db),
	}
}

// Create 创建需求
func (r *RequestRepository) Create(req *model.ServiceRequest) error {
	// 使用 Omit 排除零值的 appointment_time 字段
	return r.q.ServiceRequest.Omit(r.q.ServiceRequest.AppointmentTime).Create(req)
}

// GetByID 根据ID获取需求
func (r *RequestRepository) GetByID(id int64) (*model.ServiceRequest, error) {
	var req model.ServiceRequest
	err := r.q.ServiceRequest.UnderlyingDB().
		Table("service_requests").
		Select(`
			service_requests.*,
			assigned_station.name as station_name,
			source_station.name as source_station_name
		`).
		Joins("LEFT JOIN service_stations AS assigned_station ON service_requests.station_id = assigned_station.id").
		Joins("LEFT JOIN service_stations AS source_station ON service_requests.source_station_id = source_station.id").
		Where("service_requests.id = ?", id).
		Where("service_requests.deleted_at IS NULL").
		First(&req).Error
	if err != nil {
		return nil, err
	}
	return &req, nil
}

// GetByRequestNo 根据需求编号获取需求
func (r *RequestRepository) GetByRequestNo(requestNo string) (*model.ServiceRequest, error) {
	s := r.q.ServiceRequest
	return s.Where(s.RequestNo.Eq(requestNo)).First()
}

// ListByUser 根据用户ID查询需求列表（分页）
func (r *RequestRepository) ListByUser(userID int64, status string, offset, limit int) ([]*model.ServiceRequest, int64, error) {
	s := r.q.ServiceRequest

	// 使用原生 DB 处理条件查询
	db := s.UnderlyingDB().
		Table("service_requests").
		Select(`
			service_requests.*,
			assigned_station.name as station_name,
			source_station.name as source_station_name
		`).
		Joins("LEFT JOIN service_stations AS assigned_station ON service_requests.station_id = assigned_station.id").
		Joins("LEFT JOIN service_stations AS source_station ON service_requests.source_station_id = source_station.id").
		Where("service_requests.user_id = ?", userID).
		Where("service_requests.deleted_at IS NULL")

	if status != "" {
		db = db.Where("service_requests.status = ?", status)
	}

	// 获取总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	var reqs []*model.ServiceRequest
	if err := db.Order("service_requests.id desc").Offset(offset).Limit(limit).Find(&reqs).Error; err != nil {
		return nil, 0, err
	}

	return reqs, total, nil
}

// UpdateRating 更新评价
func (r *RequestRepository) UpdateRating(id int64, rating int, feedback string) error {
	s := r.q.ServiceRequest
	_, err := s.Where(s.ID.Eq(id)).Updates(map[string]interface{}{
		"rating":   rating,
		"feedback": feedback,
	})
	return err
}

// WithTx 事务支持
func (r *RequestRepository) WithTx(tx *gorm.DB) *RequestRepository {
	return &RequestRepository{q: query.Use(tx)}
}

// RequestWithStation 带站点信息的服务请求
type RequestWithStation struct {
	model.ServiceRequest
}

// ListAll B端查询所有需求（支持站点筛选，返回站点信息）
func (r *RequestRepository) ListAll(stationID int64, status string, offset, limit int) ([]*RequestWithStation, int64, error) {
	s := r.q.ServiceRequest
	db := s.UnderlyingDB().Table("service_requests").
		Select(`
			service_requests.*,
			assigned_station.name as station_name,
			source_station.name as source_station_name
		`).
		Joins("LEFT JOIN service_stations AS assigned_station ON service_requests.station_id = assigned_station.id").
		Joins("LEFT JOIN service_stations AS source_station ON service_requests.source_station_id = source_station.id")

	if stationID > 0 {
		db = db.Where("service_requests.station_id = ?", stationID)
	}
	if status != "" {
		db = db.Where("service_requests.status = ?", status)
	}
	db = db.Where("service_requests.deleted_at IS NULL")

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var reqs []*RequestWithStation
	if err := db.Order("service_requests.id desc").Offset(offset).Limit(limit).Find(&reqs).Error; err != nil {
		return nil, 0, err
	}

	return reqs, total, nil
}

// UpdateStatusByAdmin B端更新状态
func (r *RequestRepository) UpdateStatusByAdmin(id int64, status string, rejectReason string) error {
	s := r.q.ServiceRequest
	updates := map[string]any{"status": status}
	if rejectReason != "" {
		updates["reject_reason"] = rejectReason
	}
	_, err := s.Where(s.ID.Eq(id)).Updates(updates)
	return err
}

// CountByStatus 按状态统计需求数量
func (r *RequestRepository) CountByStatus(stationID int64, status string, isAdmin bool) (int64, error) {
	s := r.q.ServiceRequest
	db := s.UnderlyingDB().Model(&model.ServiceRequest{}).Where("status = ?", status)

	if !isAdmin && stationID > 0 {
		db = db.Where("station_id = ?", stationID)
	}

	var count int64
	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountTodayNew 统计今日新增需求数量
func (r *RequestRepository) CountTodayNew(stationID int64, isAdmin bool) (int64, error) {
	s := r.q.ServiceRequest
	start, end := todayRange()
	db := s.UnderlyingDB().Model(&model.ServiceRequest{}).
		Where("created_at >= ? AND created_at < ?", start, end)

	if !isAdmin && stationID > 0 {
		db = db.Where("station_id = ?", stationID)
	}

	var count int64
	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountSince 统计指定时间之后的需求总数
func (r *RequestRepository) CountSince(stationID int64, isAdmin bool, since time.Time) (int64, error) {
	db := r.requestBaseQuery().Where("created_at >= ?", since)
	db = applyRequestScope(db, stationID, isAdmin)

	var count int64
	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountByStatusSince 统计指定时间之后某状态的需求数量
func (r *RequestRepository) CountByStatusSince(stationID int64, status string, isAdmin bool, since time.Time) (int64, error) {
	db := r.requestBaseQuery().Where("status = ? AND created_at >= ?", status, since)
	db = applyRequestScope(db, stationID, isAdmin)

	var count int64
	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountByServiceType 按服务类型统计需求数量
func (r *RequestRepository) CountByServiceType(stationID int64, isAdmin bool, since time.Time) (map[string]int64, error) {
	s := r.q.ServiceRequest
	db := s.UnderlyingDB().Model(&model.ServiceRequest{}).
		Select("service_type, COUNT(*) as count").
		Where("created_at >= ?", since).
		Group("service_type")
	db = applyRequestScope(db, stationID, isAdmin)

	var results []struct {
		ServiceType string
		Count       int64
	}
	if err := db.Find(&results).Error; err != nil {
		return nil, err
	}

	typeCounts := make(map[string]int64)
	for _, r := range results {
		typeCounts[r.ServiceType] = r.Count
	}
	return typeCounts, nil
}

// CountBetween 统计指定时间范围内的需求总数
func (r *RequestRepository) CountBetween(stationID int64, isAdmin bool, startDate, endDate time.Time) (int64, error) {
	db := r.requestBaseQuery().
		Where("created_at >= ? AND created_at < ?", startDate, endExclusive(endDate))
	db = applyRequestScope(db, stationID, isAdmin)

	var count int64
	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountByStatusBetween 统计指定时间范围内某状态的需求数量
func (r *RequestRepository) CountByStatusBetween(stationID int64, status string, isAdmin bool, startDate, endDate time.Time) (int64, error) {
	db := r.requestBaseQuery().
		Where("status = ? AND created_at >= ? AND created_at < ?", status, startDate, endExclusive(endDate))
	db = applyRequestScope(db, stationID, isAdmin)

	var count int64
	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountByServiceTypeBetween 按服务类型统计指定时间范围内的需求数量
func (r *RequestRepository) CountByServiceTypeBetween(stationID int64, isAdmin bool, startDate, endDate time.Time) (map[string]int64, error) {
	s := r.q.ServiceRequest
	db := s.UnderlyingDB().Model(&model.ServiceRequest{}).
		Select("service_type, COUNT(*) as count").
		Where("created_at >= ? AND created_at < ?", startDate, endExclusive(endDate)).
		Group("service_type")
	db = applyRequestScope(db, stationID, isAdmin)

	var results []struct {
		ServiceType string
		Count       int64
	}
	if err := db.Find(&results).Error; err != nil {
		return nil, err
	}

	typeCounts := make(map[string]int64)
	for _, r := range results {
		typeCounts[r.ServiceType] = r.Count
	}
	return typeCounts, nil
}

func (r *RequestRepository) requestBaseQuery() *gorm.DB {
	return r.q.ServiceRequest.UnderlyingDB().Model(&model.ServiceRequest{})
}

func applyRequestScope(db *gorm.DB, stationID int64, isAdmin bool) *gorm.DB {
	if !isAdmin && stationID > 0 {
		return db.Where("station_id = ?", stationID)
	}
	return db
}

func endExclusive(endDate time.Time) time.Time {
	return endDate.AddDate(0, 0, 1)
}

func todayRange() (time.Time, time.Time) {
	start := startOfDay(time.Now())
	return start, start.AddDate(0, 0, 1)
}

func startOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

// DailyTrendItem 每日趋势数据项
type DailyTrendItem struct {
	Date  string
	Count int64
}

// GetDailyTrend 获取每日需求趋势
func (r *RequestRepository) GetDailyTrend(stationID int64, isAdmin bool, days int) ([]DailyTrendItem, error) {
	s := r.q.ServiceRequest
	startDate := startOfDay(time.Now().AddDate(0, 0, -days))
	db := s.UnderlyingDB().Model(&model.ServiceRequest{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("created_at >= ?", startDate).
		Group("DATE(created_at)").
		Order("date DESC")

	if !isAdmin && stationID > 0 {
		db = db.Where("station_id = ?", stationID)
	}

	var results []DailyTrendItem
	if err := db.Find(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

// GetDailyTrendBetween 获取指定时间范围内的每日需求趋势
func (r *RequestRepository) GetDailyTrendBetween(stationID int64, isAdmin bool, startDate, endDate time.Time) ([]DailyTrendItem, error) {
	s := r.q.ServiceRequest
	db := s.UnderlyingDB().Model(&model.ServiceRequest{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("created_at >= ? AND created_at < ?", startDate, endExclusive(endDate)).
		Group("DATE(created_at)").
		Order("date ASC")
	db = applyRequestScope(db, stationID, isAdmin)

	var results []DailyTrendItem
	if err := db.Find(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

// GetAvgRating 获取平均评分
func (r *RequestRepository) GetAvgRating(stationID int64, isAdmin bool, since time.Time) (float64, error) {
	s := r.q.ServiceRequest
	db := s.UnderlyingDB().Model(&model.ServiceRequest{}).
		Select("AVG(rating) as avg_rating").
		Where("rating > 0 AND created_at >= ?", since)

	if !isAdmin && stationID > 0 {
		db = db.Where("station_id = ?", stationID)
	}

	var result struct {
		AvgRating float64
	}
	if err := db.First(&result).Error; err != nil {
		return 0, err
	}

	return result.AvgRating, nil
}
