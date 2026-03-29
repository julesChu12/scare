package repository

import (
	"fmt"
	"testing"
	"time"

	"community-elderly-care-platform/internal/consts"
	"community-elderly-care-platform/internal/dao/model"
)

func TestTaskRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTaskRepository(db)

	task := &model.TaskAssignment{
		RequestID: 1,
		StationID: 1,
		Status:    consts.TaskStatusDispatched,
	}

	err := repo.Create(task)
	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}

	if task.ID == 0 {
		t.Error("expected task ID to be set after creation")
	}
}

func TestTaskRepository_ListByStaff(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTaskRepository(db)

	// 创建多个测试任务
	for i := 0; i < 5; i++ {
		task := &model.TaskAssignment{
			RequestID: int64(i + 1),
			StationID: 1,
			StaffID:   10,
			Status:    consts.TaskStatusClaimed,
		}
		repo.Create(task)
	}

	// 创建其他工作人员的任务
	task := &model.TaskAssignment{
		RequestID: 100,
		StationID: 1,
		StaffID:   20,
		Status:    consts.TaskStatusClaimed,
	}
	repo.Create(task)

	// 测试分页
	tasks, total, err := repo.ListByStaff(10, 0, 3)
	if err != nil {
		t.Fatalf("failed to list tasks by staff: %v", err)
	}

	if total != 5 {
		t.Errorf("expected total 5, got %d", total)
	}

	if len(tasks) != 3 {
		t.Errorf("expected 3 tasks in page, got %d", len(tasks))
	}
}

func TestTaskRepository_ListPool(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTaskRepository(db)

	// 创建站点1的任务
	for i := 0; i < 3; i++ {
		task := &model.TaskAssignment{
			RequestID: int64(i + 1),
			StationID: 1,
			Status:    consts.TaskStatusDispatched,
		}
		repo.Create(task)
	}

	// 创建站点2的任务
	for i := 0; i < 2; i++ {
		task := &model.TaskAssignment{
			RequestID: int64(i + 10),
			StationID: 2,
			Status:    consts.TaskStatusDispatched,
		}
		repo.Create(task)
	}

	// 查询站点1的任务池
	tasks, total, err := repo.ListPool(TaskPoolFilter{StationID: 1}, 0, 10)
	if err != nil {
		t.Fatalf("failed to list task pool: %v", err)
	}

	if total != 3 {
		t.Errorf("expected 3 tasks in station 1, got %d", total)
	}

	// 查询所有站点的任务池（admin）
	tasks, total, err = repo.ListPool(TaskPoolFilter{StationID: 0}, 0, 10)
	if err != nil {
		t.Fatalf("failed to list all task pool: %v", err)
	}

	if total != 5 {
		t.Errorf("expected 5 total tasks, got %d", total)
	}

	if len(tasks) != 5 {
		t.Errorf("expected 5 tasks, got %d", len(tasks))
	}
}

func TestTaskRepository_CountByStatus(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTaskRepository(db)

	// 创建不同状态的任务
	statuses := []string{
		consts.TaskStatusDispatched,
		consts.TaskStatusDispatched,
		consts.TaskStatusClaimed,
		consts.TaskStatusCompleted,
	}

	for i, status := range statuses {
		task := &model.TaskAssignment{
			RequestID: int64(i + 1),
			StationID: 1,
			Status:    status,
		}
		repo.Create(task)
	}

	// 统计已分派的任务
	count, err := repo.CountByStatus(1, consts.TaskStatusDispatched, false)
	if err != nil {
		t.Fatalf("failed to count tasks: %v", err)
	}

	if count != 2 {
		t.Errorf("expected 2 dispatched tasks, got %d", count)
	}

	// admin 统计所有站点
	count, err = repo.CountByStatus(0, consts.TaskStatusDispatched, true)
	if err != nil {
		t.Fatalf("failed to count tasks as admin: %v", err)
	}

	if count != 2 {
		t.Errorf("expected 2 dispatched tasks for admin, got %d", count)
	}
}

func TestTaskRepository_CountByStaffAndStatus(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTaskRepository(db)

	// 创建工作人员的任务
	for i := 0; i < 3; i++ {
		task := &model.TaskAssignment{
			RequestID: int64(i + 1),
			StationID: 1,
			StaffID:   10,
			Status:    consts.TaskStatusClaimed,
		}
		repo.Create(task)
	}

	// 创建已完成的任务
	task := &model.TaskAssignment{
		RequestID: 100,
		StationID: 1,
		StaffID:   10,
		Status:    consts.TaskStatusCompleted,
	}
	repo.Create(task)

	// 统计已认领的任务
	count, err := repo.CountByStaffAndStatus(10, consts.TaskStatusClaimed)
	if err != nil {
		t.Fatalf("failed to count tasks by staff: %v", err)
	}

	if count != 3 {
		t.Errorf("expected 3 claimed tasks, got %d", count)
	}
}

func TestTaskRepository_GetStaffRankingBetween(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTaskRepository(db)

	users := []model.User{
		{ID: 10, Name: "张三", Phone: "13800000010"},
		{ID: 20, Name: "李四", Phone: "13800000020"},
	}
	for _, user := range users {
		if err := db.Create(&user).Error; err != nil {
			t.Fatalf("failed to create user: %v", err)
		}
	}

	requests := []model.ServiceRequest{
		{
			ID:          1,
			RequestNo:   fmt.Sprintf("REQ-RANK-%d-1", time.Now().UnixNano()),
			UserID:      1,
			ServiceType: consts.ServiceTypeMeal,
			Status:      consts.RequestStatusCompleted,
			StationID:   1,
			Address:     "测试地址1",
			Rating:      5,
		},
		{
			ID:          2,
			RequestNo:   fmt.Sprintf("REQ-RANK-%d-2", time.Now().UnixNano()),
			UserID:      1,
			ServiceType: consts.ServiceTypeMeal,
			Status:      consts.RequestStatusCompleted,
			StationID:   1,
			Address:     "测试地址2",
			Rating:      4,
		},
		{
			ID:          3,
			RequestNo:   fmt.Sprintf("REQ-RANK-%d-3", time.Now().UnixNano()),
			UserID:      1,
			ServiceType: consts.ServiceTypeMeal,
			Status:      consts.RequestStatusCompleted,
			StationID:   1,
			Address:     "测试地址3",
			Rating:      3,
		},
		{
			ID:          4,
			RequestNo:   fmt.Sprintf("REQ-RANK-%d-4", time.Now().UnixNano()),
			UserID:      1,
			ServiceType: consts.ServiceTypeMeal,
			Status:      consts.RequestStatusCompleted,
			StationID:   1,
			Address:     "测试地址4",
			Rating:      0,
		},
	}
	for _, req := range requests {
		if err := db.Create(&req).Error; err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
	}

	tasks := []*model.TaskAssignment{
		{
			RequestID:   1,
			StationID:   1,
			StaffID:     10,
			Status:      consts.TaskStatusCompleted,
			CompletedAt: time.Date(2026, 1, 5, 10, 0, 0, 0, time.Local),
			CreatedAt:   time.Date(2026, 1, 5, 8, 0, 0, 0, time.Local),
			ClaimedAt:   time.Date(2026, 1, 5, 9, 0, 0, 0, time.Local),
		},
		{
			RequestID:   2,
			StationID:   1,
			StaffID:     10,
			Status:      consts.TaskStatusCompleted,
			CompletedAt: time.Date(2026, 2, 6, 10, 0, 0, 0, time.Local),
			CreatedAt:   time.Date(2026, 2, 6, 8, 0, 0, 0, time.Local),
			ClaimedAt:   time.Date(2026, 2, 6, 9, 0, 0, 0, time.Local),
		},
		{
			RequestID:   3,
			StationID:   1,
			StaffID:     20,
			Status:      consts.TaskStatusCompleted,
			CompletedAt: time.Date(2026, 1, 10, 10, 0, 0, 0, time.Local),
			CreatedAt:   time.Date(2026, 1, 10, 8, 0, 0, 0, time.Local),
			ClaimedAt:   time.Date(2026, 1, 10, 9, 0, 0, 0, time.Local),
		},
		{
			RequestID:   4,
			StationID:   1,
			StaffID:     20,
			Status:      consts.TaskStatusCompleted,
			CompletedAt: time.Date(2026, 1, 12, 10, 0, 0, 0, time.Local),
			CreatedAt:   time.Date(2026, 1, 12, 8, 0, 0, 0, time.Local),
			ClaimedAt:   time.Date(2026, 1, 12, 9, 0, 0, 0, time.Local),
		},
	}
	for _, task := range tasks {
		if err := db.Create(task).Error; err != nil {
			t.Fatalf("failed to create task: %v", err)
		}
	}

	ranking, err := repo.GetStaffRankingBetween(
		1,
		false,
		time.Date(2026, 1, 1, 0, 0, 0, 0, time.Local),
		time.Date(2026, 1, 31, 0, 0, 0, 0, time.Local),
		10,
	)
	if err != nil {
		t.Fatalf("failed to get staff ranking between dates: %v", err)
	}

	if len(ranking) != 2 {
		t.Fatalf("expected 2 staff in ranking, got %d", len(ranking))
	}

	if ranking[0].ID != 20 || ranking[0].CompletedCount != 2 {
		t.Errorf("unexpected first ranking item: %#v", ranking[0])
	}
	if ranking[0].AvgRating != 3 {
		t.Errorf("expected first staff avg rating 3.0, got %#v", ranking[0].AvgRating)
	}

	if ranking[1].ID != 10 || ranking[1].CompletedCount != 1 {
		t.Errorf("unexpected second ranking item: %#v", ranking[1])
	}
	if ranking[1].AvgRating != 5 {
		t.Errorf("expected second staff avg rating 5.0, got %#v", ranking[1].AvgRating)
	}
}

func TestTaskRepository_WithTx(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTaskRepository(db)

	// 测试事务
	tx := db.Begin()
	txRepo := repo.WithTx(tx)

	task := &model.TaskAssignment{
		RequestID: 1,
		StationID: 1,
		Status:    consts.TaskStatusDispatched,
	}
	txRepo.Create(task)

	// 回滚事务
	tx.Rollback()

	// 验证任务未创建（通过列表查询确认无记录）
	tasks, total, err := repo.ListByStaff(0, 0, 10)
	if err != nil {
		t.Fatalf("failed to list tasks: %v", err)
	}
	if total != 0 || len(tasks) != 0 {
		t.Error("expected task to not exist after rollback")
	}
}
