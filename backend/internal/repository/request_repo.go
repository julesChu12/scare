package repository

import (
	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/dao/query"
	"community-elderly-care-platform/internal/consts"

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
	s := r.q.ServiceRequest
	return s.Where(s.ID.Eq(id)).First()
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
	db := s.UnderlyingDB().Model(&model.ServiceRequest{}).Where("user_id = ?", userID)
	
	if status != "" {
		db = db.Where("status = ?", status)
	}
	
	// 获取总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 分页查询
	var reqs []*model.ServiceRequest
	if err := db.Order("id desc").Offset(offset).Limit(limit).Find(&reqs).Error; err != nil {
		return nil, 0, err
	}
	
	return reqs, total, nil
}

// ListDispatchedByStation 根据站点ID查询已分派的需求
func (r *RequestRepository) ListDispatchedByStation(stationID int64, offset, limit int) ([]*model.ServiceRequest, int64, error) {
	s := r.q.ServiceRequest
	
	// 获取总数
	total, err := s.Where(s.StationID.Eq(stationID), s.Status.Eq(consts.RequestStatusDispatched)).Count()
	if err != nil {
		return nil, 0, err
	}
	
	// 分页查询
	reqs, err := s.Where(s.StationID.Eq(stationID), s.Status.Eq(consts.RequestStatusDispatched)).
		Order(s.ID.Desc()).
		Offset(offset).
		Limit(limit).
		Find()
	
	if err != nil {
		return nil, 0, err
	}
	
	return reqs, total, nil
}

// UpdateStatus 更新需求状态
func (r *RequestRepository) UpdateStatus(id int64, status string) error {
	s := r.q.ServiceRequest
	_, err := s.Where(s.ID.Eq(id)).Update(s.Status, status)
	return err
}

// UpdateStationAndStatus 更新站点和状态
func (r *RequestRepository) UpdateStationAndStatus(id int64, stationID int64, status string) error {
	s := r.q.ServiceRequest
	_, err := s.Where(s.ID.Eq(id)).Updates(map[string]interface{}{
		"station_id": stationID,
		"status":     status,
	})
	return err
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
	StationName string `json:"station_name"`
}

// ListAll B端查询所有需求（支持站点筛选，返回站点信息）
func (r *RequestRepository) ListAll(stationID int64, status string, offset, limit int) ([]*RequestWithStation, int64, error) {
	s := r.q.ServiceRequest
	db := s.UnderlyingDB().Table("service_requests").
		Select("service_requests.*, service_stations.name as station_name").
		Joins("LEFT JOIN service_stations ON service_requests.station_id = service_stations.id")

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
	db := s.UnderlyingDB().Model(&model.ServiceRequest{}).
		Where("DATE(created_at) = CURDATE()")

	if !isAdmin && stationID > 0 {
		db = db.Where("station_id = ?", stationID)
	}

	var count int64
	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
