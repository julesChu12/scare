package service

import (
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"community-elderly-care-platform/internal/consts"
	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/repository"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var requestServiceTestSeq int64

func setupRequestServiceForTest(t *testing.T) (*RequestService, *gorm.DB) {
	t.Helper()

	tmpFile := t.TempDir() + "/request_service_test.db"
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

	reqRepo := repository.NewRequestRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	stationRepo := repository.NewStationRepository(db)
	svc := NewRequestService(db, reqRepo, taskRepo, stationRepo, &GeofenceService{}, nil, nil, nil)
	return svc, db
}

func createRequestWithTask(t *testing.T, db *gorm.DB, reqStatus string, userID int64) (*model.ServiceRequest, *model.TaskAssignment) {
	t.Helper()
	seq := atomic.AddInt64(&requestServiceTestSeq, 1)

	req := &model.ServiceRequest{
		RequestNo:         fmt.Sprintf("REQTEST-%d-%d", time.Now().UnixNano(), seq),
		UserID:            userID,
		ServiceType:       consts.ServiceTypeMeal,
		Status:            reqStatus,
		SubmitLocationLat: 30.000001,
		SubmitLocationLng: 120.000001,
		Address:           "test",
		StationID:         1,
	}
	if err := db.Create(req).Error; err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	task := &model.TaskAssignment{
		RequestID: req.ID,
		StationID: 1,
		Status:    consts.TaskStatusDispatched,
	}
	if err := db.Create(task).Error; err != nil {
		t.Fatalf("failed to create task: %v", err)
	}
	return req, task
}

func TestRequestService_ResolveCoordinates(t *testing.T) {
	svc, _ := setupRequestServiceForTest(t)

	lat := 31.22
	lng := 121.48
	_, _, _, err := svc.resolveCoordinates(RequestInput{
		Lat: &lat,
	})
	if !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("expected ErrInvalidRequest for partial coordinates, got %v", err)
	}

	badLat := 91.0
	_, _, _, err = svc.resolveCoordinates(RequestInput{
		Lat: &badLat,
		Lng: &lng,
	})
	if !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("expected ErrInvalidRequest for out-of-range coordinates, got %v", err)
	}

	_, _, _, err = svc.resolveCoordinates(RequestInput{
		Address: "test address",
	})
	if err == nil || !strings.Contains(err.Error(), "geocode service not configured") {
		t.Fatalf("expected geocode service not configured error, got %v", err)
	}

	gotLat, gotLng, addr, err := svc.resolveCoordinates(RequestInput{
		Lat:     &lat,
		Lng:     &lng,
		Address: "input address",
	})
	if err != nil {
		t.Fatalf("resolveCoordinates failed: %v", err)
	}
	if gotLat != lat || gotLng != lng {
		t.Fatalf("unexpected coordinates: got (%v, %v), want (%v, %v)", gotLat, gotLng, lat, lng)
	}
	if addr != "input address" {
		t.Fatalf("unexpected address: got %q", addr)
	}
}

func TestRequestService_UpdateStatus_InvalidStatus(t *testing.T) {
	svc, _ := setupRequestServiceForTest(t)

	err := svc.UpdateStatus(1, "unknown_status", "")
	if !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("expected ErrInvalidRequest, got %v", err)
	}
}

func TestRequestService_Cancel_SuccessAndIdempotent(t *testing.T) {
	svc, db := setupRequestServiceForTest(t)
	req, _ := createRequestWithTask(t, db, consts.RequestStatusDispatched, 101)

	cancelled, changed, err := svc.Cancel(req.ID, 101)
	if err != nil {
		t.Fatalf("cancel failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected changed=true on first cancel")
	}
	if cancelled.Status != consts.RequestStatusCancelled {
		t.Fatalf("unexpected request status: %s", cancelled.Status)
	}

	var dbTask model.TaskAssignment
	if err := db.Where("request_id = ?", req.ID).First(&dbTask).Error; err != nil {
		t.Fatalf("failed to reload task: %v", err)
	}
	if dbTask.Status != consts.TaskStatusCancelled {
		t.Fatalf("expected task status cancelled, got %s", dbTask.Status)
	}

	cancelledAgain, changedAgain, err := svc.Cancel(req.ID, 101)
	if err != nil {
		t.Fatalf("second cancel failed: %v", err)
	}
	if changedAgain {
		t.Fatalf("expected changed=false for idempotent cancel")
	}
	if cancelledAgain.Status != consts.RequestStatusCancelled {
		t.Fatalf("unexpected status on second cancel: %s", cancelledAgain.Status)
	}
}

func TestRequestService_Cancel_Conflict(t *testing.T) {
	svc, db := setupRequestServiceForTest(t)

	req1, _ := createRequestWithTask(t, db, consts.RequestStatusDispatched, 201)
	_, _, err := svc.Cancel(req1.ID, 202)
	if !errors.Is(err, ErrRequestConflict) {
		t.Fatalf("expected ErrRequestConflict for non-owner cancel, got %v", err)
	}

	req2, _ := createRequestWithTask(t, db, consts.RequestStatusCompleted, 203)
	_, _, err = svc.Cancel(req2.ID, 203)
	if !errors.Is(err, ErrRequestConflict) {
		t.Fatalf("expected ErrRequestConflict for completed request cancel, got %v", err)
	}
}

func TestRequestService_Rate(t *testing.T) {
	svc, db := setupRequestServiceForTest(t)

	if _, err := svc.Rate(1, 1, 0, "bad"); !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("expected ErrInvalidRequest for invalid rating, got %v", err)
	}

	conflictReq, _ := createRequestWithTask(t, db, consts.RequestStatusCompleted, 301)
	if _, err := svc.Rate(conflictReq.ID, 302, 5, "ok"); !errors.Is(err, ErrRequestConflict) {
		t.Fatalf("expected ErrRequestConflict for ownership mismatch, got %v", err)
	}

	notCompletedReq, _ := createRequestWithTask(t, db, consts.RequestStatusClaimed, 303)
	if _, err := svc.Rate(notCompletedReq.ID, 303, 5, "ok"); !errors.Is(err, ErrNotCompleted) {
		t.Fatalf("expected ErrNotCompleted, got %v", err)
	}

	alreadyRatedReq, _ := createRequestWithTask(t, db, consts.RequestStatusCompleted, 304)
	if err := db.Model(&model.ServiceRequest{}).Where("id = ?", alreadyRatedReq.ID).Update("rating", 4).Error; err != nil {
		t.Fatalf("failed to seed rated request: %v", err)
	}
	if _, err := svc.Rate(alreadyRatedReq.ID, 304, 5, "again"); !errors.Is(err, ErrAlreadyRated) {
		t.Fatalf("expected ErrAlreadyRated, got %v", err)
	}

	successReq, _ := createRequestWithTask(t, db, consts.RequestStatusCompleted, 305)
	updated, err := svc.Rate(successReq.ID, 305, 5, "good service")
	if err != nil {
		t.Fatalf("rate failed: %v", err)
	}
	if updated.Rating != 5 || updated.Feedback != "good service" {
		t.Fatalf("unexpected returned rating payload: rating=%d feedback=%q", updated.Rating, updated.Feedback)
	}

	var dbReq model.ServiceRequest
	if err := db.First(&dbReq, successReq.ID).Error; err != nil {
		t.Fatalf("failed to reload rated request: %v", err)
	}
	if dbReq.Rating != 5 || dbReq.Feedback != "good service" {
		t.Fatalf("rating not persisted, got rating=%d feedback=%q", dbReq.Rating, dbReq.Feedback)
	}
}
