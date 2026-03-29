package repository

import (
	"community-elderly-care-platform/internal/consts"
	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/dao/query"
	"time"

	"gorm.io/gorm"
)

type TaskRepository struct {
	q *query.Query
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{
		q: query.Use(db),
	}
}

// Create 创建任务
func (r *TaskRepository) Create(task *model.TaskAssignment) error {
	// 使用 Omit 排除零值的时间字段和空的 JSON 字段
	return r.q.TaskAssignment.Omit(r.q.TaskAssignment.ClaimedAt, r.q.TaskAssignment.CompletedAt, r.q.TaskAssignment.Images).Create(task)
}

// ListByStaff 根据工作人员ID获取任务列表（分页，包含关联的服务请求）
func (r *TaskRepository) ListByStaff(staffID int64, offset, limit int) ([]*TaskWithRequest, int64, error) {
	db := r.q.TaskAssignment.UnderlyingDB()

	query := db.Table("task_assignments").
		Joins("LEFT JOIN service_requests ON task_assignments.request_id = service_requests.id").
		Where("task_assignments.deleted_at IS NULL AND task_assignments.staff_id = ?", staffID)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var tasks []*model.TaskAssignment
	if err := query.Select("task_assignments.*").Order("task_assignments.id DESC").Offset(offset).Limit(limit).Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	if len(tasks) == 0 {
		return []*TaskWithRequest{}, total, nil
	}

	requestIDs := make([]int64, 0, len(tasks))
	for _, t := range tasks {
		requestIDs = append(requestIDs, t.RequestID)
	}

	var requests []*model.ServiceRequest
	if err := db.Session(&gorm.Session{}).Table("service_requests").Where("id IN ?", requestIDs).Find(&requests).Error; err != nil {
		return nil, 0, err
	}
	requestMap := make(map[int64]*model.ServiceRequest)
	for _, req := range requests {
		requestMap[req.ID] = req
	}

	var staffName string
	db.Session(&gorm.Session{}).Table("users").Select("name").Where("id = ?", staffID).Scan(&staffName)

	result := make([]*TaskWithRequest, len(tasks))
	for i, t := range tasks {
		result[i] = &TaskWithRequest{
			TaskAssignment: *t,
			Request:        requestMap[t.RequestID],
			StaffName:      staffName,
		}
	}

	return result, total, nil
}

// TaskPoolFilter 任务池筛选条件
type TaskPoolFilter struct {
	StationID int64  // 站点ID，0 表示所有站点（仅 admin）
	Status    string // 任务状态，空表示默认 dispatched
}

// ListPool 根据筛选条件查询任务池（包含关联的服务请求）
func (r *TaskRepository) ListPool(filter TaskPoolFilter, offset, limit int) ([]*TaskWithRequest, int64, error) {
	db := r.q.TaskAssignment.UnderlyingDB()

	// 构建基础查询
	query := db.Table("task_assignments").
		Joins("LEFT JOIN service_requests ON task_assignments.request_id = service_requests.id").
		Where("task_assignments.deleted_at IS NULL")

	// 站点筛选：0 表示所有站点
	if filter.StationID > 0 {
		query = query.Where("task_assignments.station_id = ?", filter.StationID)
	}

	// 状态筛选：空表示所有状态
	if status := filter.Status; status != "" {
		query = query.Where("task_assignments.status = ?", status)
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询任务
	var tasks []*model.TaskAssignment
	if err := query.Select("task_assignments.*").Order("task_assignments.id DESC").Offset(offset).Limit(limit).Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	// 批量查询关联的服务请求
	if len(tasks) == 0 {
		return []*TaskWithRequest{}, total, nil
	}

	requestIDs := make([]int64, 0, len(tasks))
	staffIDs := make([]int64, 0, len(tasks))
	for _, t := range tasks {
		requestIDs = append(requestIDs, t.RequestID)
		if t.StaffID > 0 {
			staffIDs = append(staffIDs, t.StaffID)
		}
	}

	// 查询服务请求
	var requests []*model.ServiceRequest
	if err := db.Session(&gorm.Session{}).Table("service_requests").Where("id IN ?", requestIDs).Find(&requests).Error; err != nil {
		return nil, 0, err
	}
	requestMap := make(map[int64]*model.ServiceRequest)
	for _, req := range requests {
		requestMap[req.ID] = req
	}

	// 查询工作人员姓名
	staffMap := make(map[int64]string)
	if len(staffIDs) > 0 {
		var staffUsers []struct {
			ID   int64
			Name string
		}
		if err := db.Table("users").Select("id, name").Where("id IN ?", staffIDs).Find(&staffUsers).Error; err == nil {
			for _, s := range staffUsers {
				staffMap[s.ID] = s.Name
			}
		}
	}

	// 组装结果
	result := make([]*TaskWithRequest, len(tasks))
	for i, t := range tasks {
		result[i] = &TaskWithRequest{
			TaskAssignment: *t,
			Request:        requestMap[t.RequestID],
			StaffName:      staffMap[t.StaffID],
		}
	}

	return result, total, nil
}

// WithTx 事务支持
func (r *TaskRepository) WithTx(tx *gorm.DB) *TaskRepository {
	return &TaskRepository{q: query.Use(tx)}
}

// CountByStatus 按状态统计任务数量
func (r *TaskRepository) CountByStatus(stationID int64, status string, isAdmin bool) (int64, error) {
	t := r.q.TaskAssignment
	db := t.UnderlyingDB().Model(&model.TaskAssignment{}).Where("status = ?", status)

	if !isAdmin && stationID > 0 {
		db = db.Where("station_id = ?", stationID)
	}

	var count int64
	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountTodayCompleted 统计今日完成的任务数量
func (r *TaskRepository) CountTodayCompleted(stationID int64, isAdmin bool) (int64, error) {
	t := r.q.TaskAssignment
	start, end := todayRange()
	db := t.UnderlyingDB().Model(&model.TaskAssignment{}).
		Where("status = ?", consts.TaskStatusCompleted).
		Where("completed_at >= ? AND completed_at < ?", start, end)

	if !isAdmin && stationID > 0 {
		db = db.Where("station_id = ?", stationID)
	}

	var count int64
	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountByStaffAndStatus 按工作人员和状态统计任务数量
func (r *TaskRepository) CountByStaffAndStatus(staffID int64, status string) (int64, error) {
	t := r.q.TaskAssignment
	return t.Where(t.StaffID.Eq(staffID), t.Status.Eq(status)).Count()
}

// TaskWithRequest 任务及关联的服务请求
type TaskWithRequest struct {
	model.TaskAssignment
	Request   *model.ServiceRequest `json:"request,omitempty"`
	StaffName string                `json:"staff_name,omitempty"`
}

// TaskListFilter 任务列表筛选条件
type TaskListFilter struct {
	StationID   int64  // 站点ID，0 表示所有站点
	Status      string // 任务状态
	ServiceType string // 服务类型
	RequestNo   string // 需求编号
	StaffID     int64  // 工作人员ID
}

// ListWithRequest 查询任务列表并关联服务请求
func (r *TaskRepository) ListWithRequest(filter TaskListFilter, offset, limit int) ([]*TaskWithRequest, int64, error) {
	db := r.q.TaskAssignment.UnderlyingDB()

	// 构建基础查询 - 不使用别名
	query := db.Table("task_assignments").
		Joins("LEFT JOIN service_requests ON task_assignments.request_id = service_requests.id").
		Where("task_assignments.deleted_at IS NULL")

	// 站点筛选
	if filter.StationID > 0 {
		query = query.Where("task_assignments.station_id = ?", filter.StationID)
	}

	// 状态筛选
	if filter.Status != "" {
		query = query.Where("task_assignments.status = ?", filter.Status)
	}

	// 服务类型筛选
	if filter.ServiceType != "" {
		query = query.Where("service_requests.service_type = ?", filter.ServiceType)
	}

	// 需求编号筛选
	if filter.RequestNo != "" {
		query = query.Where("service_requests.request_no LIKE ?", "%"+filter.RequestNo+"%")
	}

	// 工作人员筛选
	if filter.StaffID > 0 {
		query = query.Where("task_assignments.staff_id = ?", filter.StaffID)
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询任务
	var tasks []*model.TaskAssignment
	if err := query.Select("task_assignments.*").Order("task_assignments.id DESC").Offset(offset).Limit(limit).Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	// 批量查询关联的服务请求
	if len(tasks) == 0 {
		return []*TaskWithRequest{}, total, nil
	}

	requestIDs := make([]int64, 0, len(tasks))
	staffIDs := make([]int64, 0, len(tasks))
	for _, t := range tasks {
		requestIDs = append(requestIDs, t.RequestID)
		if t.StaffID > 0 {
			staffIDs = append(staffIDs, t.StaffID)
		}
	}

	// 查询服务请求 - 使用新的 DB 会话
	var requests []*model.ServiceRequest
	if err := db.Session(&gorm.Session{}).Table("service_requests").Where("id IN ?", requestIDs).Find(&requests).Error; err != nil {
		return nil, 0, err
	}
	requestMap := make(map[int64]*model.ServiceRequest)
	for _, req := range requests {
		requestMap[req.ID] = req
	}

	// 查询工作人员姓名
	staffMap := make(map[int64]string)
	if len(staffIDs) > 0 {
		var staffUsers []struct {
			ID   int64
			Name string
		}
		if err := db.Table("users").Select("id, name").Where("id IN ?", staffIDs).Find(&staffUsers).Error; err == nil {
			for _, s := range staffUsers {
				staffMap[s.ID] = s.Name
			}
		}
	}

	// 组装结果
	result := make([]*TaskWithRequest, len(tasks))
	for i, t := range tasks {
		result[i] = &TaskWithRequest{
			TaskAssignment: *t,
			Request:        requestMap[t.RequestID],
			StaffName:      staffMap[t.StaffID],
		}
	}

	return result, total, nil
}

// GetByIDWithRequest 根据ID获取任务及关联的服务请求
func (r *TaskRepository) GetByIDWithRequest(id int64) (*TaskWithRequest, error) {
	t := r.q.TaskAssignment
	task, err := t.Where(t.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}

	db := t.UnderlyingDB()

	// 查询关联的服务请求 - 使用新的 DB 会话
	var request model.ServiceRequest
	if err := db.Session(&gorm.Session{}).Table("service_requests").Where("id = ?", task.RequestID).First(&request).Error; err == nil {
		result := &TaskWithRequest{
			TaskAssignment: *task,
			Request:        &request,
		}

		// 查询工作人员姓名
		if task.StaffID > 0 {
			var staffName string
			db.Session(&gorm.Session{}).Table("users").Select("name").Where("id = ?", task.StaffID).Scan(&staffName)
			result.StaffName = staffName
		}

		return result, nil
	}

	return &TaskWithRequest{TaskAssignment: *task}, nil
}

func (r *TaskRepository) GetAvgResponseTime(stationID int64, isAdmin bool, since time.Time) (int64, error) {
	t := r.q.TaskAssignment
	db := t.UnderlyingDB().Model(&model.TaskAssignment{}).
		Select("AVG(TIMESTAMPDIFF(MINUTE, created_at, claimed_at)) as avg_time").
		Where("claimed_at IS NOT NULL AND created_at >= ?", since)

	if !isAdmin && stationID > 0 {
		db = db.Where("station_id = ?", stationID)
	}

	var result struct {
		AvgTime *float64
	}
	if err := db.First(&result).Error; err != nil {
		return 0, err
	}

	if result.AvgTime == nil {
		return 0, nil
	}
	return int64(*result.AvgTime), nil
}

func (r *TaskRepository) GetAvgProcessTime(stationID int64, isAdmin bool, since time.Time) (int64, error) {
	t := r.q.TaskAssignment
	db := t.UnderlyingDB().Model(&model.TaskAssignment{}).
		Select("AVG(TIMESTAMPDIFF(MINUTE, claimed_at, completed_at)) as avg_time").
		Where("completed_at IS NOT NULL AND claimed_at IS NOT NULL AND created_at >= ?", since)

	if !isAdmin && stationID > 0 {
		db = db.Where("station_id = ?", stationID)
	}

	var result struct {
		AvgTime *float64
	}
	if err := db.First(&result).Error; err != nil {
		return 0, err
	}

	if result.AvgTime == nil {
		return 0, nil
	}
	return int64(*result.AvgTime), nil
}

func (r *TaskRepository) GetStaffRanking(stationID int64, isAdmin bool, since time.Time, limit int) ([]StaffRankingItem, error) {
	t := r.q.TaskAssignment
	db := t.UnderlyingDB().
		Table("task_assignments").
		Select(`
				task_assignments.staff_id as id,
				users.name,
				COUNT(*) as completed_count,
				COALESCE(AVG(CASE WHEN service_requests.rating > 0 THEN service_requests.rating END), 0) as avg_rating
			`).
		Joins("LEFT JOIN users ON task_assignments.staff_id = users.id").
		Joins("LEFT JOIN service_requests ON task_assignments.request_id = service_requests.id").
		Where("task_assignments.status = ? AND task_assignments.completed_at >= ?", consts.TaskStatusCompleted, since).
		Group("task_assignments.staff_id, users.name").
		Order("completed_count DESC").
		Limit(limit)

	if !isAdmin && stationID > 0 {
		db = db.Where("task_assignments.station_id = ?", stationID)
	}

	var results []StaffRankingItem
	if err := db.Find(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (r *TaskRepository) GetStaffRankingBetween(stationID int64, isAdmin bool, startDate, endDate time.Time, limit int) ([]StaffRankingItem, error) {
	t := r.q.TaskAssignment
	db := t.UnderlyingDB().
		Table("task_assignments").
		Select(`
				task_assignments.staff_id as id,
				users.name,
				COUNT(*) as completed_count,
				COALESCE(AVG(CASE WHEN service_requests.rating > 0 THEN service_requests.rating END), 0) as avg_rating
			`).
		Joins("LEFT JOIN users ON task_assignments.staff_id = users.id").
		Joins("LEFT JOIN service_requests ON task_assignments.request_id = service_requests.id").
		Where("task_assignments.status = ? AND task_assignments.completed_at >= ? AND task_assignments.completed_at < ?", consts.TaskStatusCompleted, startDate, taskEndExclusive(endDate)).
		Group("task_assignments.staff_id, users.name").
		Order("completed_count DESC").
		Limit(limit)

	if !isAdmin && stationID > 0 {
		db = db.Where("task_assignments.station_id = ?", stationID)
	}

	var results []StaffRankingItem
	if err := db.Find(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func taskEndExclusive(endDate time.Time) time.Time {
	return endDate.AddDate(0, 0, 1)
}

type StaffRankingItem struct {
	ID             int64   `json:"id"`
	Name           string  `json:"name"`
	CompletedCount int64   `json:"completed_count"`
	AvgRating      float64 `json:"avg_rating"`
	IsOnline       bool    `json:"is_online"`
}
