package repository

import (
	"testing"

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
