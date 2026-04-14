package service

import (
	"errors"
	"fmt"
	"testing"

	"community-elderly-care-platform/internal/consts"
	"community-elderly-care-platform/internal/dao/model"
)

func TestRequestService_Create_InvalidAndIdempotentBranches(t *testing.T) {
	svc, db := setupRequestServiceForTest(t)
	svc.geocodeSvc = NewGeocodeService("")
	createTestStation(t, db, 39.908823, 116.397470)

	if _, _, err := svc.Create(RequestInput{UserID: 0, ServiceType: consts.ServiceTypeMeal}); !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("expected ErrInvalidRequest for missing user, got %v", err)
	}
	if _, _, err := svc.Create(RequestInput{UserID: 1, ServiceType: "unknown"}); !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("expected ErrInvalidRequest for invalid service type, got %v", err)
	}
	if _, _, err := svc.Create(RequestInput{UserID: 1, ServiceType: consts.ServiceTypeMeal, Address: " "}); !errors.Is(err, ErrAddressRequired) {
		t.Fatalf("expected ErrAddressRequired for blank address, got %v", err)
	}

	createdReq, created, err := svc.Create(RequestInput{
		UserID:       1,
		RequestNo:    "REQ-IDEMPOTENT-1",
		ServiceType:  consts.ServiceTypeMeal,
		ContactName:  "测试人",
		ContactPhone: "13800138002",
		Address:      "测试地址",
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if !created {
		t.Fatalf("expected first create to return created=true")
	}

	gotReq, created, err := svc.Create(RequestInput{
		UserID:       1,
		RequestNo:    "REQ-IDEMPOTENT-1",
		ServiceType:  consts.ServiceTypeMeal,
		ContactName:  "测试人",
		ContactPhone: "13800138002",
		Address:      "测试地址",
	})
	if err != nil {
		t.Fatalf("second Create returned error: %v", err)
	}
	if created {
		t.Fatalf("expected idempotent create to return created=false")
	}
	if gotReq.ID != createdReq.ID {
		t.Fatalf("expected same request ID %d, got %d", createdReq.ID, gotReq.ID)
	}

	_, _, err = svc.Create(RequestInput{
		UserID:       2,
		RequestNo:    "REQ-IDEMPOTENT-1",
		ServiceType:  consts.ServiceTypeMeal,
		ContactName:  "冲突用户",
		ContactPhone: "13800138003",
		Address:      "测试地址",
	})
	if !errors.Is(err, ErrRequestConflict) {
		t.Fatalf("expected ErrRequestConflict, got %v", err)
	}
}

func TestRequestService_Create_UsesGeocodeAddressWhenCoordinatesMissing(t *testing.T) {
	svc, db := setupRequestServiceForTest(t)
	svc.geocodeSvc = NewGeocodeService("")
	expected := createTestStation(t, db, 39.908823, 116.397470)

	req, created, err := svc.Create(RequestInput{
		UserID:       11,
		ServiceType:  consts.ServiceTypeMedical,
		Address:      "北京市东城区",
		ContactName:  "地理编码用户",
		ContactPhone: "13800138004",
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if !created {
		t.Fatalf("expected created=true")
	}
	if req.StationID != expected.ID {
		t.Fatalf("expected station %d, got %d", expected.ID, req.StationID)
	}
	if req.ServiceLocationLat != 39.908823 || req.ServiceLocationLng != 116.397470 {
		t.Fatalf("unexpected service coordinates: (%f, %f)", req.ServiceLocationLat, req.ServiceLocationLng)
	}
}

func TestRequestService_Create_ManualReviewWhenAddressCannotBeResolved(t *testing.T) {
	svc, db := setupRequestServiceForTest(t)
	sourceStation := createTestStation(t, db, 39.9, 116.3)
	sourceStationID := sourceStation.ID

	req, created, err := svc.Create(RequestInput{
		UserID:          12,
		ServiceType:     consts.ServiceTypeMedical,
		Address:         "无法解析的地址",
		SourceStationID: &sourceStationID,
		ContactName:     "人工复核用户",
		ContactPhone:    "13800138006",
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if !created {
		t.Fatalf("expected created=true")
	}
	if req.Status != consts.RequestStatusPending {
		t.Fatalf("expected pending status, got %s", req.Status)
	}
	if req.StationID != 0 || !req.NeedsManualVerify || req.DispatchBasis != DispatchBasisAddressManualReview {
		t.Fatalf("unexpected manual review request: %+v", req)
	}
	if req.SourceStationID != sourceStationID {
		t.Fatalf("expected source station %d to be recorded, got %d", sourceStationID, req.SourceStationID)
	}

	var taskCount int64
	if err := db.Model(&model.TaskAssignment{}).Where("request_id = ?", req.ID).Count(&taskCount).Error; err != nil {
		t.Fatalf("failed to count tasks: %v", err)
	}
	if taskCount != 0 {
		t.Fatalf("expected no task to be created for manual review request, got %d", taskCount)
	}
}

func TestRequestService_ListByUser_UpdateAndGetByID(t *testing.T) {
	svc, db := setupRequestServiceForTest(t)

	if _, _, err := svc.ListByUser(0, "", 1, 10); !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("expected ErrInvalidRequest for zero user ID, got %v", err)
	}
	if _, err := svc.GetByID(0); !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("expected ErrInvalidRequest for zero request ID, got %v", err)
	}

	req1, _ := createRequestWithTask(t, db, consts.RequestStatusDispatched, 21)
	req2, _ := createRequestWithTask(t, db, consts.RequestStatusCompleted, 21)
	createRequestWithTask(t, db, consts.RequestStatusDispatched, 22)

	list, total, err := svc.ListByUser(21, consts.RequestStatusDispatched, 1, 10)
	if err != nil {
		t.Fatalf("ListByUser returned error: %v", err)
	}
	if total != 1 || len(list) != 1 || list[0].ID != req1.ID {
		t.Fatalf("unexpected ListByUser result: total=%d len=%d ids=%v", total, len(list), list)
	}

	unchanged, err := svc.Update(req1.ID, UpdateInput{})
	if err != nil {
		t.Fatalf("Update with empty payload returned error: %v", err)
	}
	if unchanged.ID != req1.ID {
		t.Fatalf("expected same request ID after noop update, got %d", unchanged.ID)
	}

	if _, err := svc.Update(req2.ID, UpdateInput{ContactName: "不可编辑"}); !errors.Is(err, ErrRequestConflict) {
		t.Fatalf("expected ErrRequestConflict for completed request update, got %v", err)
	}

	if _, err := svc.Update(req1.ID, UpdateInput{ServiceType: "unknown"}); !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("expected ErrInvalidRequest for invalid service type, got %v", err)
	}

	updated, err := svc.Update(req1.ID, UpdateInput{
		ServiceType:  consts.ServiceTypeMedical,
		ContactName:  "新联系人",
		ContactPhone: "13800138005",
		Address:      "新地址",
		Description:  "新描述",
		Urgency:      "high",
	})
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if updated.ServiceType != consts.ServiceTypeMedical || updated.ContactName != "新联系人" || updated.Address != "新地址" {
		t.Fatalf("unexpected updated request: %+v", updated)
	}
}

func TestRequestService_CancelByAdmin_ListAllAndUpdateStatus(t *testing.T) {
	svc, db := setupRequestServiceForTest(t)
	stationA := createTestStation(t, db, 30, 120)
	stationB := createTestStation(t, db, 31, 121)

	req1, _ := createRequestWithTask(t, db, consts.RequestStatusDispatched, 31)
	req2, _ := createRequestWithTask(t, db, consts.RequestStatusCompleted, 32)

	if err := db.Model(&model.ServiceRequest{}).Where("id = ?", req1.ID).Update("station_id", stationA.ID).Error; err != nil {
		t.Fatalf("failed to seed req1 station_id: %v", err)
	}
	if err := db.Model(&model.TaskAssignment{}).Where("request_id = ?", req1.ID).Update("station_id", stationA.ID).Error; err != nil {
		t.Fatalf("failed to seed task1 station_id: %v", err)
	}
	if err := db.Model(&model.ServiceRequest{}).Where("id = ?", req2.ID).Update("station_id", stationB.ID).Error; err != nil {
		t.Fatalf("failed to seed req2 station_id: %v", err)
	}

	cancelled, err := svc.CancelByAdmin(req1.ID)
	if err != nil {
		t.Fatalf("CancelByAdmin returned error: %v", err)
	}
	if cancelled.Status != consts.RequestStatusCancelled {
		t.Fatalf("expected request cancelled, got %s", cancelled.Status)
	}

	cancelledAgain, err := svc.CancelByAdmin(req1.ID)
	if err != nil {
		t.Fatalf("second CancelByAdmin returned error: %v", err)
	}
	if cancelledAgain.Status != consts.RequestStatusCancelled {
		t.Fatalf("expected request to remain cancelled, got %s", cancelledAgain.Status)
	}

	if _, err := svc.CancelByAdmin(req2.ID); !errors.Is(err, ErrRequestConflict) {
		t.Fatalf("expected ErrRequestConflict for completed request, got %v", err)
	}

	req3 := &model.ServiceRequest{
		RequestNo:   fmt.Sprintf("REQ-LISTALL-%d", req1.ID),
		UserID:      33,
		ServiceType: consts.ServiceTypeMeal,
		Status:      consts.RequestStatusDispatched,
		Address:     "列表地址",
		StationID:   stationA.ID,
	}
	if err := db.Create(req3).Error; err != nil {
		t.Fatalf("failed to create req3: %v", err)
	}

	if err := svc.UpdateStatus(req3.ID, consts.RequestStatusRejected, "资源不足"); err != nil {
		t.Fatalf("UpdateStatus returned error: %v", err)
	}

	list, total, err := svc.ListAll(stationA.ID, consts.RequestStatusRejected, 1, 10)
	if err != nil {
		t.Fatalf("ListAll returned error: %v", err)
	}
	if total != 1 || len(list) != 1 {
		t.Fatalf("unexpected ListAll result: total=%d len=%d", total, len(list))
	}
	if list[0].StationName != stationA.Name || list[0].RejectReason != "资源不足" {
		t.Fatalf("unexpected ListAll payload: %+v", list[0])
	}
}

func TestRequestService_UpdateStationID(t *testing.T) {
	svc, db := setupRequestServiceForTest(t)
	stationA := createTestStation(t, db, 30, 120)
	stationB := createTestStation(t, db, 31, 121)

	req, _ := createRequestWithTask(t, db, consts.RequestStatusPending, 41)

	// 初始 station_id 为 1（来自 createRequestWithTask）
	var initialStationID int64
	if err := db.Model(&model.ServiceRequest{}).Where("id = ?", req.ID).Pluck("station_id", &initialStationID).Error; err != nil {
		t.Fatalf("failed to get initial station_id: %v", err)
	}
	if initialStationID != 1 {
		t.Fatalf("expected initial station_id=1, got %d", initialStationID)
	}

	// 更新 station_id 到 stationA
	newStationID := stationA.ID
	updated, err := svc.Update(req.ID, UpdateInput{
		StationID: &newStationID,
	})
	if err != nil {
		t.Fatalf("Update station_id failed: %v", err)
	}
	if updated.StationID != stationA.ID {
		t.Fatalf("expected station_id=%d after update, got %d", stationA.ID, updated.StationID)
	}

	// 再次更新 station_id 到 stationB
	newStationID = stationB.ID
	updated2, err := svc.Update(req.ID, UpdateInput{
		StationID: &newStationID,
	})
	if err != nil {
		t.Fatalf("Update station_id to stationB failed: %v", err)
	}
	if updated2.StationID != stationB.ID {
		t.Fatalf("expected station_id=%d after second update, got %d", stationB.ID, updated2.StationID)
	}

	// 已完成状态的任务不能更新 station_id
	completedReq, _ := createRequestWithTask(t, db, consts.RequestStatusCompleted, 42)
	_, err = svc.Update(completedReq.ID, UpdateInput{StationID: &stationA.ID})
	if !errors.Is(err, ErrRequestConflict) {
		t.Fatalf("expected ErrRequestConflict for completed request, got %v", err)
	}
}
