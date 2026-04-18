package service

import (
	"community-elderly-care-platform/internal/consts"
	"community-elderly-care-platform/internal/repository"
	"time"
)

// DashboardStats 工作台统计数据
type DashboardStats struct {
	// 任务统计
	TaskStats TaskStats `json:"task_stats"`
	// 需求统计
	RequestStats RequestStats `json:"request_stats"`
	// 今日统计
	TodayStats TodayStats `json:"today_stats"`
	// 我的任务统计（仅工作人员）
	MyTaskStats *MyTaskStats `json:"my_task_stats,omitempty"`
}

// TaskStats 任务状态统计
type TaskStats struct {
	Total      int64 `json:"total"`
	Dispatched int64 `json:"dispatched"`
	Claimed    int64 `json:"claimed"`
	Completed  int64 `json:"completed"`
	Cancelled  int64 `json:"cancelled"`
}

// RequestStats 需求状态统计
type RequestStats struct {
	Total      int64 `json:"total"`
	Pending    int64 `json:"pending"`
	Dispatched int64 `json:"dispatched"`
	Processing int64 `json:"processing"`
	Completed  int64 `json:"completed"`
	Cancelled  int64 `json:"cancelled"`
}

// TodayStats 今日统计
type TodayStats struct {
	NewRequests     int64 `json:"new_requests"`
	CompletedTasks  int64 `json:"completed_tasks"`
	NewUsers        int64 `json:"new_users"`
	AvgResponseTime int64 `json:"avg_response_time"` // 平均响应时间（分钟）
}

// MyTaskStats 我的任务统计
type MyTaskStats struct {
	Claimed   int64 `json:"claimed"`
	Completed int64 `json:"completed"`
	Total     int64 `json:"total"`
}

type StatisticsService struct {
	taskRepo    *repository.TaskRepository
	requestRepo *repository.RequestRepository
	userRepo    *repository.UserRepository
}

// NewStatisticsService 创建统计服务
func NewStatisticsService(
	taskRepo *repository.TaskRepository,
	requestRepo *repository.RequestRepository,
	userRepo *repository.UserRepository,
) *StatisticsService {
	return &StatisticsService{
		taskRepo:    taskRepo,
		requestRepo: requestRepo,
		userRepo:    userRepo,
	}
}

// GetDashboardStats 获取工作台统计数据。
//
// 说明：
// - stationID=0 且 isAdmin=true 时返回全局视角
// - 非管理员视角下，调用方应提前传入已收口后的站点范围
// - 工作人员额外返回“我的任务”统计，管理员不返回该字段
func (s *StatisticsService) GetDashboardStats(userID, stationID int64, isAdmin bool) (*DashboardStats, error) {
	taskStats, err := s.GetTaskStats(stationID, isAdmin)
	if err != nil {
		return nil, err
	}

	requestStats, err := s.GetRequestStats(stationID, isAdmin)
	if err != nil {
		return nil, err
	}

	todayStats, err := s.GetTodayStats(stationID, isAdmin)
	if err != nil {
		return nil, err
	}

	stats := &DashboardStats{
		TaskStats:    *taskStats,
		RequestStats: *requestStats,
		TodayStats:   *todayStats,
	}

	// 如果是工作人员，获取个人任务统计
	if !isAdmin && userID > 0 {
		myStats, err := s.GetMyTaskStats(userID)
		if err == nil {
			stats.MyTaskStats = myStats
		}
	}

	return stats, nil
}

// GetTaskStats 获取任务统计。
func (s *StatisticsService) GetTaskStats(stationID int64, isAdmin bool) (*TaskStats, error) {
	stats := &TaskStats{}

	// 获取各状态数量
	dispatched, err := s.taskRepo.CountByStatus(stationID, consts.TaskStatusDispatched, isAdmin)
	if err != nil {
		return nil, err
	}
	stats.Dispatched = dispatched

	claimed, err := s.taskRepo.CountByStatus(stationID, consts.TaskStatusClaimed, isAdmin)
	if err != nil {
		return nil, err
	}
	stats.Claimed = claimed

	completed, err := s.taskRepo.CountByStatus(stationID, consts.TaskStatusCompleted, isAdmin)
	if err != nil {
		return nil, err
	}
	stats.Completed = completed

	cancelled, err := s.taskRepo.CountByStatus(stationID, consts.TaskStatusCancelled, isAdmin)
	if err != nil {
		return nil, err
	}
	stats.Cancelled = cancelled

	stats.Total = dispatched + claimed + completed + cancelled

	return stats, nil
}

// GetRequestStats 获取需求统计。
func (s *StatisticsService) GetRequestStats(stationID int64, isAdmin bool) (*RequestStats, error) {
	stats := &RequestStats{}

	pending, err := s.requestRepo.CountByStatus(stationID, consts.RequestStatusPending, isAdmin)
	if err != nil {
		return nil, err
	}
	stats.Pending = pending

	dispatched, err := s.requestRepo.CountByStatus(stationID, consts.RequestStatusDispatched, isAdmin)
	if err != nil {
		return nil, err
	}
	stats.Dispatched = dispatched

	processing, err := s.requestRepo.CountByStatus(stationID, consts.RequestStatusProcessing, isAdmin)
	if err != nil {
		return nil, err
	}
	stats.Processing = processing

	completed, err := s.requestRepo.CountByStatus(stationID, consts.RequestStatusCompleted, isAdmin)
	if err != nil {
		return nil, err
	}
	stats.Completed = completed

	cancelled, err := s.requestRepo.CountByStatus(stationID, consts.RequestStatusCancelled, isAdmin)
	if err != nil {
		return nil, err
	}
	stats.Cancelled = cancelled

	stats.Total = pending + dispatched + processing + completed + cancelled

	return stats, nil
}

// GetTodayStats 获取今日统计。
//
// 说明：
// - NewUsers 目前仅在管理员全局视角下返回
// - AvgResponseTime 字段预留在返回结构中，当前尚未计算写入
func (s *StatisticsService) GetTodayStats(stationID int64, isAdmin bool) (*TodayStats, error) {
	stats := &TodayStats{}

	newRequests, err := s.requestRepo.CountTodayNew(stationID, isAdmin)
	if err != nil {
		return nil, err
	}
	stats.NewRequests = newRequests

	completedTasks, err := s.taskRepo.CountTodayCompleted(stationID, isAdmin)
	if err != nil {
		return nil, err
	}
	stats.CompletedTasks = completedTasks

	// 今日新用户（仅管理员可见全局数据）
	if isAdmin {
		newUsers, err := s.userRepo.CountTodayNew()
		if err == nil {
			stats.NewUsers = newUsers
		}
	}

	return stats, nil
}

// GetMyTaskStats 获取我的任务统计。
func (s *StatisticsService) GetMyTaskStats(userID int64) (*MyTaskStats, error) {
	stats := &MyTaskStats{}

	claimed, err := s.taskRepo.CountByStaffAndStatus(userID, consts.TaskStatusClaimed)
	if err != nil {
		return nil, err
	}
	stats.Claimed = claimed

	completed, err := s.taskRepo.CountByStaffAndStatus(userID, consts.TaskStatusCompleted)
	if err != nil {
		return nil, err
	}
	stats.Completed = completed

	stats.Total = claimed + completed

	return stats, nil
}

// ========== 数据中心 - 统计概览 ==========

// OverviewStats 统计概览数据
type OverviewStats struct {
	TotalRequests int64 `json:"total_requests"`
	Pending       int64 `json:"pending"`
	Completed     int64 `json:"completed"`
	InProgress    int64 `json:"in_progress"`
}

// ServiceTypeStats 服务类型统计
type ServiceTypeStats struct {
	Type       string  `json:"type"`
	Name       string  `json:"name"`
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
}

// TrendItem 趋势数据项
type TrendItem struct {
	Date       string  `json:"date"`
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
}

// EfficiencyStats 处理效率统计
type EfficiencyStats struct {
	AvgResponseTime  int64   `json:"avg_response_time"` // 平均响应时间（分钟）
	AvgProcessTime   int64   `json:"avg_process_time"`  // 平均处理时间（分钟）
	SatisfactionRate float64 `json:"satisfaction_rate"` // 满意度（百分比）
}

// StaffRankingItem 服务人员排行项
type StaffRankingItem = repository.StaffRankingItem

// GetOverviewStats 获取统计概览数据（总需求数、各状态数量等）
func (s *StatisticsService) GetOverviewStats(stationID int64, isAdmin bool, days int) (*OverviewStats, error) {
	stats := &OverviewStats{}

	// 获取指定时间范围内的需求统计
	startDate := time.Now().AddDate(0, 0, -days)

	total, err := s.requestRepo.CountSince(stationID, isAdmin, startDate)
	if err != nil {
		return nil, err
	}
	stats.TotalRequests = total

	pending, err := s.requestRepo.CountByStatusSince(stationID, consts.RequestStatusPending, isAdmin, startDate)
	if err != nil {
		return nil, err
	}
	stats.Pending = pending

	completed, err := s.requestRepo.CountByStatusSince(stationID, consts.RequestStatusCompleted, isAdmin, startDate)
	if err != nil {
		return nil, err
	}
	stats.Completed = completed

	// 进行中 = 已派发 + 处理中
	dispatched, _ := s.requestRepo.CountByStatusSince(stationID, consts.RequestStatusDispatched, isAdmin, startDate)
	processing, _ := s.requestRepo.CountByStatusSince(stationID, consts.RequestStatusProcessing, isAdmin, startDate)
	stats.InProgress = dispatched + processing

	return stats, nil
}

// GetServiceTypeStats 获取服务类型分布统计
func (s *StatisticsService) GetServiceTypeStats(stationID int64, isAdmin bool, days int) ([]ServiceTypeStats, error) {
	startDate := time.Now().AddDate(0, 0, -days)

	// 从数据库获取各服务类型的数量
	typeCounts, err := s.requestRepo.CountByServiceType(stationID, isAdmin, startDate)
	if err != nil {
		return nil, err
	}

	// 计算总数
	var total int64
	for _, count := range typeCounts {
		total += count
	}

	// 构建结果
	var result []ServiceTypeStats
	for _, serviceType := range consts.ServiceTypes {
		count := typeCounts[serviceType]
		if count > 0 {
			percentage := float64(0)
			if total > 0 {
				percentage = float64(count) * 100 / float64(total)
			}
			result = append(result, ServiceTypeStats{
				Type:       serviceType,
				Name:       consts.GetServiceTypeName(serviceType),
				Count:      count,
				Percentage: percentage,
			})
		}
	}

	return result, nil
}

// GetRequestTrend 获取需求趋势数据（每日需求量）
func (s *StatisticsService) GetRequestTrend(stationID int64, isAdmin bool, days int) ([]TrendItem, error) {
	trend, err := s.requestRepo.GetDailyTrend(stationID, isAdmin, days)
	if err != nil {
		return nil, err
	}

	// 找出最大值用于计算百分比
	var maxCount int64
	for _, item := range trend {
		if item.Count > maxCount {
			maxCount = item.Count
		}
	}

	// 计算百分比
	result := make([]TrendItem, len(trend))
	for i, item := range trend {
		percentage := float64(0)
		if maxCount > 0 {
			percentage = float64(item.Count) * 100 / float64(maxCount)
		}
		result[i] = TrendItem{
			Date:       item.Date,
			Count:      item.Count,
			Percentage: percentage,
		}
	}

	return result, nil
}

// GetEfficiencyStats 获取处理效率统计（平均响应时间、满意度等）
func (s *StatisticsService) GetEfficiencyStats(stationID int64, isAdmin bool, days int) (*EfficiencyStats, error) {
	startDate := time.Now().AddDate(0, 0, -days)

	stats := &EfficiencyStats{}

	// 获取平均响应时间（从创建到派发的时间）
	avgResponse, err := s.taskRepo.GetAvgResponseTime(stationID, isAdmin, startDate)
	if err == nil {
		stats.AvgResponseTime = avgResponse
	}

	// 获取平均处理时间（从认领到完成的时间）
	avgProcess, err := s.taskRepo.GetAvgProcessTime(stationID, isAdmin, startDate)
	if err == nil {
		stats.AvgProcessTime = avgProcess
	}

	// 获取满意度（有评分的需求的平均分）
	satisfaction, err := s.requestRepo.GetAvgRating(stationID, isAdmin, startDate)
	if err == nil {
		// 将5分制转换为百分比
		stats.SatisfactionRate = satisfaction * 20
	}

	return stats, nil
}

// GetStaffRanking 获取服务人员绩效排行
func (s *StatisticsService) GetStaffRanking(stationID int64, isAdmin bool, days int, limit int) ([]StaffRankingItem, error) {
	startDate := time.Now().AddDate(0, 0, -days)

	ranking, err := s.taskRepo.GetStaffRanking(stationID, isAdmin, startDate, limit)
	if err != nil {
		return nil, err
	}

	return ranking, nil
}
