package service

import (
	"community-elderly-care-platform/internal/consts"
	"community-elderly-care-platform/internal/repository"
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
	NewRequests      int64 `json:"new_requests"`
	CompletedTasks   int64 `json:"completed_tasks"`
	NewUsers         int64 `json:"new_users"`
	AvgResponseTime  int64 `json:"avg_response_time"` // 平均响应时间（分钟）
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

// GetDashboardStats 获取工作台统计数据
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

// GetTaskStats 获取任务统计
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

// GetRequestStats 获取需求统计
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

// GetTodayStats 获取今日统计
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

// GetMyTaskStats 获取我的任务统计
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
