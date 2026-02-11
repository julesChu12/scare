package service

import (
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"community-elderly-care-platform/internal/consts"
	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/repository"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var taskServiceTestSeq int64

func setupTaskServiceForTest(t *testing.T) (*TaskService, *gorm.DB) {
	t.Helper()

	tmpFile := t.TempDir() + "/task_service_test.db"
	dsn := tmpFile + "?_loc=auto&parseTime=true"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		NowFunc: func() time.Time {
			return time.Now().Local().Truncate(time.Second)
		},
	})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	if err := db.Exec(`
		CREATE TABLE service_requests (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			request_no TEXT,
			user_id INTEGER,
			service_type TEXT,
			status TEXT,
			description TEXT,
			submit_location_lat REAL,
			submit_location_lng REAL,
			contact_name TEXT,
			contact_phone TEXT,
			address TEXT,
			appointment_time DATETIME,
			urgency TEXT,
			station_id INTEGER,
			reject_reason TEXT,
			images TEXT,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME,
			rating INTEGER,
			feedback TEXT
		);
	`).Error; err != nil {
		t.Fatalf("failed to create service_requests table: %v", err)
	}

	if err := db.Exec(`
		CREATE TABLE task_assignments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			request_id INTEGER,
			station_id INTEGER,
			staff_id INTEGER,
			status TEXT,
			claimed_at DATETIME,
			completed_at DATETIME,
			rating INTEGER,
			feedback TEXT,
			staff_notes TEXT,
			images TEXT,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		);
	`).Error; err != nil {
		t.Fatalf("failed to create task_assignments table: %v", err)
	}

	taskRepo := repository.NewTaskRepository(db)
	requestRepo := repository.NewRequestRepository(db)
	svc := NewTaskService(db, taskRepo, requestRepo, nil)
	return svc, db
}

func seedTaskScenario(t *testing.T, db *gorm.DB, reqStatus, taskStatus string, staffID int64) (*model.ServiceRequest, *model.TaskAssignment) {
	t.Helper()
	seq := atomic.AddInt64(&taskServiceTestSeq, 1)

	req := &model.ServiceRequest{
		RequestNo:         fmt.Sprintf("REQTASK-%d-%d", time.Now().UnixNano(), seq),
		UserID:            1001,
		ServiceType:       consts.ServiceTypeMeal,
		Status:            reqStatus,
		SubmitLocationLat: 30.100001,
		SubmitLocationLng: 120.100001,
		Address:           "task test",
		StationID:         1,
	}
	if err := db.Create(req).Error; err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	task := &model.TaskAssignment{
		RequestID: req.ID,
		StationID: 1,
		StaffID:   staffID,
		Status:    taskStatus,
	}
	if err := db.Create(task).Error; err != nil {
		t.Fatalf("failed to create task: %v", err)
	}

	return req, task
}

func TestTaskService_Claim(t *testing.T) {
	svc, db := setupTaskServiceForTest(t)

	if _, _, err := svc.Claim(0, 1); !errors.Is(err, ErrTaskInvalid) {
		t.Fatalf("expected ErrTaskInvalid, got %v", err)
	}

	_, completedTask := seedTaskScenario(t, db, consts.RequestStatusCompleted, consts.TaskStatusCompleted, 11)
	if _, _, err := svc.Claim(completedTask.ID, 11); !errors.Is(err, ErrTaskConflict) {
		t.Fatalf("expected ErrTaskConflict for completed task, got %v", err)
	}

	_, claimedSameStaff := seedTaskScenario(t, db, consts.RequestStatusClaimed, consts.TaskStatusClaimed, 21)
	gotTask, changed, err := svc.Claim(claimedSameStaff.ID, 21)
	if err != nil {
		t.Fatalf("idempotent claim failed: %v", err)
	}
	if changed {
		t.Fatalf("expected changed=false for idempotent claim")
	}
	if gotTask.Status != consts.TaskStatusClaimed || gotTask.StaffID != 21 {
		t.Fatalf("unexpected task from idempotent claim: status=%s staff=%d", gotTask.Status, gotTask.StaffID)
	}

	_, claimedOtherStaff := seedTaskScenario(t, db, consts.RequestStatusClaimed, consts.TaskStatusClaimed, 31)
	if _, _, err := svc.Claim(claimedOtherStaff.ID, 32); !errors.Is(err, ErrTaskConflict) {
		t.Fatalf("expected ErrTaskConflict for claiming task already held by another staff, got %v", err)
	}

	reqSuccess, taskSuccess := seedTaskScenario(t, db, consts.RequestStatusDispatched, consts.TaskStatusDispatched, 0)
	updatedTask, changed, err := svc.Claim(taskSuccess.ID, 41)
	if err != nil {
		t.Fatalf("claim failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected changed=true on successful claim")
	}
	if updatedTask.Status != consts.TaskStatusClaimed || updatedTask.StaffID != 41 {
		t.Fatalf("unexpected claimed task values: status=%s staff=%d", updatedTask.Status, updatedTask.StaffID)
	}

	var updatedReq model.ServiceRequest
	if err := db.First(&updatedReq, reqSuccess.ID).Error; err != nil {
		t.Fatalf("failed to reload request: %v", err)
	}
	if updatedReq.Status != consts.RequestStatusClaimed {
		t.Fatalf("expected request status claimed, got %s", updatedReq.Status)
	}

	reqRollback, taskRollback := seedTaskScenario(t, db, consts.RequestStatusClaimed, consts.TaskStatusDispatched, 0)
	if _, _, err := svc.Claim(taskRollback.ID, 51); !errors.Is(err, ErrTaskConflict) {
		t.Fatalf("expected ErrTaskConflict when request status is not dispatched, got %v", err)
	}
	var rollbackTask model.TaskAssignment
	if err := db.First(&rollbackTask, taskRollback.ID).Error; err != nil {
		t.Fatalf("failed to reload rollback task: %v", err)
	}
	if rollbackTask.Status != consts.TaskStatusDispatched || rollbackTask.StaffID != 0 {
		t.Fatalf("expected task update rollback, got status=%s staff=%d", rollbackTask.Status, rollbackTask.StaffID)
	}
	var rollbackReq model.ServiceRequest
	if err := db.First(&rollbackReq, reqRollback.ID).Error; err != nil {
		t.Fatalf("failed to reload rollback request: %v", err)
	}
	if rollbackReq.Status != consts.RequestStatusClaimed {
		t.Fatalf("expected request status to remain claimed, got %s", rollbackReq.Status)
	}
}

func TestTaskService_Complete(t *testing.T) {
	svc, db := setupTaskServiceForTest(t)

	if _, _, err := svc.Complete(0, 1, nil); !errors.Is(err, ErrTaskInvalid) {
		t.Fatalf("expected ErrTaskInvalid, got %v", err)
	}

	_, taskWrongStaff := seedTaskScenario(t, db, consts.RequestStatusClaimed, consts.TaskStatusClaimed, 61)
	if _, _, err := svc.Complete(taskWrongStaff.ID, 62, []string{"a.jpg"}); !errors.Is(err, ErrTaskConflict) {
		t.Fatalf("expected ErrTaskConflict for wrong staff completion, got %v", err)
	}

	_, taskCancelled := seedTaskScenario(t, db, consts.RequestStatusCancelled, consts.TaskStatusCancelled, 71)
	if _, _, err := svc.Complete(taskCancelled.ID, 71, []string{"a.jpg"}); !errors.Is(err, ErrTaskConflict) {
		t.Fatalf("expected ErrTaskConflict for cancelled task completion, got %v", err)
	}

	_, taskCompleted := seedTaskScenario(t, db, consts.RequestStatusCompleted, consts.TaskStatusCompleted, 81)
	gotTask, changed, err := svc.Complete(taskCompleted.ID, 81, []string{"a.jpg"})
	if err != nil {
		t.Fatalf("idempotent complete failed: %v", err)
	}
	if changed {
		t.Fatalf("expected changed=false for idempotent complete")
	}
	if gotTask.Status != consts.TaskStatusCompleted {
		t.Fatalf("unexpected task status from idempotent complete: %s", gotTask.Status)
	}

	reqSuccess, taskSuccess := seedTaskScenario(t, db, consts.RequestStatusClaimed, consts.TaskStatusClaimed, 91)
	updatedTask, changed, err := svc.Complete(taskSuccess.ID, 91, []string{"proof1.jpg", "proof2.jpg"})
	if err != nil {
		t.Fatalf("complete failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected changed=true on successful complete")
	}
	if updatedTask.Status != consts.TaskStatusCompleted {
		t.Fatalf("expected completed task status, got %s", updatedTask.Status)
	}
	if updatedTask.Images != "[\"proof1.jpg\",\"proof2.jpg\"]" {
		t.Fatalf("unexpected task images payload: %s", updatedTask.Images)
	}

	var updatedReq model.ServiceRequest
	if err := db.First(&updatedReq, reqSuccess.ID).Error; err != nil {
		t.Fatalf("failed to reload request: %v", err)
	}
	if updatedReq.Status != consts.RequestStatusCompleted {
		t.Fatalf("expected request status completed, got %s", updatedReq.Status)
	}
}

func TestTaskService_Transfer(t *testing.T) {
	svc, db := setupTaskServiceForTest(t)

	if _, err := svc.Transfer(0, 1); !errors.Is(err, ErrTaskInvalid) {
		t.Fatalf("expected ErrTaskInvalid, got %v", err)
	}

	_, completedTask := seedTaskScenario(t, db, consts.RequestStatusCompleted, consts.TaskStatusCompleted, 101)
	if _, err := svc.Transfer(completedTask.ID, 102); !errors.Is(err, ErrTaskConflict) {
		t.Fatalf("expected ErrTaskConflict for completed task transfer, got %v", err)
	}

	_, sameStaffTask := seedTaskScenario(t, db, consts.RequestStatusClaimed, consts.TaskStatusClaimed, 111)
	task, err := svc.Transfer(sameStaffTask.ID, 111)
	if err != nil {
		t.Fatalf("same staff transfer should be no-op, got err: %v", err)
	}
	if task.StaffID != 111 || task.Status != consts.TaskStatusClaimed {
		t.Fatalf("unexpected task for no-op transfer: staff=%d status=%s", task.StaffID, task.Status)
	}

	reqSuccess, taskSuccess := seedTaskScenario(t, db, consts.RequestStatusDispatched, consts.TaskStatusDispatched, 0)
	updatedTask, err := svc.Transfer(taskSuccess.ID, 121)
	if err != nil {
		t.Fatalf("transfer failed: %v", err)
	}
	if updatedTask.StaffID != 121 || updatedTask.Status != consts.TaskStatusClaimed {
		t.Fatalf("unexpected transferred task values: staff=%d status=%s", updatedTask.StaffID, updatedTask.Status)
	}

	var dbTask model.TaskAssignment
	if err := db.First(&dbTask, taskSuccess.ID).Error; err != nil {
		t.Fatalf("failed to reload transferred task: %v", err)
	}
	if dbTask.StaffID != 121 || dbTask.Status != consts.TaskStatusClaimed {
		t.Fatalf("transfer not persisted: staff=%d status=%s", dbTask.StaffID, dbTask.Status)
	}

	var updatedReq model.ServiceRequest
	if err := db.First(&updatedReq, reqSuccess.ID).Error; err != nil {
		t.Fatalf("failed to reload request: %v", err)
	}
	if updatedReq.Status != consts.RequestStatusClaimed {
		t.Fatalf("expected request status claimed after transfer, got %s", updatedReq.Status)
	}
}
