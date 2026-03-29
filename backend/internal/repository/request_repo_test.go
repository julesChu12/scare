package repository

import (
	"fmt"
	"testing"
	"time"

	"community-elderly-care-platform/internal/consts"
	"community-elderly-care-platform/internal/dao/model"
)

func TestRequestRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRequestRepository(db)

	req := &model.ServiceRequest{
		RequestNo:         fmt.Sprintf("REQ-%d", time.Now().UnixNano()),
		UserID:            1,
		ServiceType:       consts.ServiceTypeMeal,
		Status:            consts.RequestStatusPending,
		SubmitLocationLat: 30.0,
		SubmitLocationLng: 120.0,
		Address:           "测试地址",
	}

	err := repo.Create(req)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	if req.ID == 0 {
		t.Error("expected request ID to be set after creation")
	}
}

func TestRequestRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRequestRepository(db)

	// 创建测试需求
	req := &model.ServiceRequest{
		RequestNo:   fmt.Sprintf("REQ-%d", time.Now().UnixNano()),
		UserID:      1,
		ServiceType: consts.ServiceTypeMeal,
		Status:      consts.RequestStatusPending,
		Address:     "测试地址",
	}
	repo.Create(req)

	// 查询需求
	result, err := repo.GetByID(req.ID)
	if err != nil {
		t.Fatalf("failed to get request: %v", err)
	}

	if result.UserID != 1 {
		t.Errorf("expected user_id 1, got %d", result.UserID)
	}
}

func TestRequestRepository_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRequestRepository(db)

	_, err := repo.GetByID(999)
	if err == nil {
		t.Error("expected error for non-existent request, got nil")
	}
}

func TestRequestRepository_GetByRequestNo(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRequestRepository(db)

	requestNo := fmt.Sprintf("REQ-%d", time.Now().UnixNano())
	req := &model.ServiceRequest{
		RequestNo:   requestNo,
		UserID:      1,
		ServiceType: consts.ServiceTypeMeal,
		Status:      consts.RequestStatusPending,
		Address:     "测试地址",
	}
	repo.Create(req)

	// 根据编号查询
	result, err := repo.GetByRequestNo(requestNo)
	if err != nil {
		t.Fatalf("failed to get request by request_no: %v", err)
	}

	if result.ID != req.ID {
		t.Errorf("expected request ID %d, got %d", req.ID, result.ID)
	}
}

func TestRequestRepository_ListByUser(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRequestRepository(db)

	// 创建用户1的需求
	for i := 0; i < 5; i++ {
		req := &model.ServiceRequest{
			RequestNo:   fmt.Sprintf("REQ-U1-%d-%d", time.Now().UnixNano(), i),
			UserID:      1,
			ServiceType: consts.ServiceTypeMeal,
			Status:      consts.RequestStatusPending,
			Address:     "测试地址",
		}
		repo.Create(req)
	}

	// 创建用户2的需求
	req := &model.ServiceRequest{
		RequestNo:   fmt.Sprintf("REQ-U2-%d", time.Now().UnixNano()),
		UserID:      2,
		ServiceType: consts.ServiceTypeMeal,
		Status:      consts.RequestStatusPending,
		Address:     "测试地址",
	}
	repo.Create(req)

	// 测试分页
	reqs, total, err := repo.ListByUser(1, "", 0, 3)
	if err != nil {
		t.Fatalf("failed to list requests by user: %v", err)
	}

	if total != 5 {
		t.Errorf("expected total 5, got %d", total)
	}

	if len(reqs) != 3 {
		t.Errorf("expected 3 requests in page, got %d", len(reqs))
	}
}

func TestRequestRepository_ListByUser_WithStatus(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRequestRepository(db)

	// 创建不同状态的需求
	statuses := []string{
		consts.RequestStatusPending,
		consts.RequestStatusPending,
		consts.RequestStatusDispatched,
		consts.RequestStatusCompleted,
	}

	for i, status := range statuses {
		req := &model.ServiceRequest{
			RequestNo:   fmt.Sprintf("REQ-S-%d-%d", time.Now().UnixNano(), i),
			UserID:      1,
			ServiceType: consts.ServiceTypeMeal,
			Status:      status,
			Address:     "测试地址",
		}
		repo.Create(req)
	}

	// 按状态筛选
	reqs, total, err := repo.ListByUser(1, consts.RequestStatusPending, 0, 10)
	if err != nil {
		t.Fatalf("failed to list requests with status filter: %v", err)
	}

	if total != 2 {
		t.Errorf("expected 2 pending requests, got %d", total)
	}

	if len(reqs) != 2 {
		t.Errorf("expected 2 requests, got %d", len(reqs))
	}
}

func TestRequestRepository_UpdateRating(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRequestRepository(db)

	// 创建测试需求
	req := &model.ServiceRequest{
		RequestNo:   fmt.Sprintf("REQ-%d", time.Now().UnixNano()),
		UserID:      1,
		ServiceType: consts.ServiceTypeMeal,
		Status:      consts.RequestStatusCompleted,
		Address:     "测试地址",
	}
	repo.Create(req)

	// 更新评价
	err := repo.UpdateRating(req.ID, 5, "服务很好")
	if err != nil {
		t.Fatalf("failed to update rating: %v", err)
	}

	// 验证更新
	updated, _ := repo.GetByID(req.ID)
	if updated.Rating != 5 {
		t.Errorf("expected rating 5, got %d", updated.Rating)
	}
	if updated.Feedback != "服务很好" {
		t.Errorf("expected feedback '服务很好', got %s", updated.Feedback)
	}
}

func TestRequestRepository_CountByStatus(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRequestRepository(db)

	// 创建不同状态的需求
	statuses := []string{
		consts.RequestStatusPending,
		consts.RequestStatusPending,
		consts.RequestStatusDispatched,
		consts.RequestStatusCompleted,
	}

	for i, status := range statuses {
		req := &model.ServiceRequest{
			RequestNo:   fmt.Sprintf("REQ-C-%d-%d", time.Now().UnixNano(), i),
			UserID:      1,
			ServiceType: consts.ServiceTypeMeal,
			Status:      status,
			StationID:   1,
			Address:     "测试地址",
		}
		repo.Create(req)
	}

	// 统计待处理的需求
	count, err := repo.CountByStatus(1, consts.RequestStatusPending, false)
	if err != nil {
		t.Fatalf("failed to count requests: %v", err)
	}

	if count != 2 {
		t.Errorf("expected 2 pending requests, got %d", count)
	}

	// admin 统计所有站点
	count, err = repo.CountByStatus(0, consts.RequestStatusPending, true)
	if err != nil {
		t.Fatalf("failed to count requests as admin: %v", err)
	}

	if count != 2 {
		t.Errorf("expected 2 pending requests for admin, got %d", count)
	}
}

func TestRequestRepository_CountBetween(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRequestRepository(db)

	requests := []*model.ServiceRequest{
		{
			RequestNo:   fmt.Sprintf("REQ-RANGE-%d-1", time.Now().UnixNano()),
			UserID:      1,
			ServiceType: consts.ServiceTypeMeal,
			Status:      consts.RequestStatusPending,
			StationID:   1,
			Address:     "测试地址1",
			CreatedAt:   time.Date(2026, 1, 3, 10, 0, 0, 0, time.Local),
		},
		{
			RequestNo:   fmt.Sprintf("REQ-RANGE-%d-2", time.Now().UnixNano()),
			UserID:      1,
			ServiceType: consts.ServiceTypeMedical,
			Status:      consts.RequestStatusCompleted,
			StationID:   1,
			Address:     "测试地址2",
			CreatedAt:   time.Date(2026, 1, 20, 15, 0, 0, 0, time.Local),
		},
		{
			RequestNo:   fmt.Sprintf("REQ-RANGE-%d-3", time.Now().UnixNano()),
			UserID:      1,
			ServiceType: consts.ServiceTypeMeal,
			Status:      consts.RequestStatusCompleted,
			StationID:   1,
			Address:     "测试地址3",
			CreatedAt:   time.Date(2026, 2, 2, 9, 0, 0, 0, time.Local),
		},
	}

	for _, req := range requests {
		if err := repo.Create(req); err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
	}

	startDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.Local)
	endDate := time.Date(2026, 1, 31, 0, 0, 0, 0, time.Local)

	total, err := repo.CountBetween(1, false, startDate, endDate)
	if err != nil {
		t.Fatalf("failed to count requests between dates: %v", err)
	}
	if total != 2 {
		t.Errorf("expected 2 requests in range, got %d", total)
	}

	completed, err := repo.CountByStatusBetween(1, consts.RequestStatusCompleted, false, startDate, endDate)
	if err != nil {
		t.Fatalf("failed to count completed requests between dates: %v", err)
	}
	if completed != 1 {
		t.Errorf("expected 1 completed request in range, got %d", completed)
	}

	typeCounts, err := repo.CountByServiceTypeBetween(1, false, startDate, endDate)
	if err != nil {
		t.Fatalf("failed to count request types between dates: %v", err)
	}
	if typeCounts[consts.ServiceTypeMeal] != 1 || typeCounts[consts.ServiceTypeMedical] != 1 {
		t.Errorf("unexpected type counts: %#v", typeCounts)
	}
}

func TestRequestRepository_GetDailyTrendBetween(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRequestRepository(db)

	requests := []*model.ServiceRequest{
		{
			RequestNo:   fmt.Sprintf("REQ-TREND-%d-1", time.Now().UnixNano()),
			UserID:      1,
			ServiceType: consts.ServiceTypeMeal,
			Status:      consts.RequestStatusPending,
			StationID:   1,
			Address:     "测试地址1",
			CreatedAt:   time.Date(2026, 1, 2, 8, 0, 0, 0, time.Local),
		},
		{
			RequestNo:   fmt.Sprintf("REQ-TREND-%d-2", time.Now().UnixNano()),
			UserID:      1,
			ServiceType: consts.ServiceTypeMeal,
			Status:      consts.RequestStatusPending,
			StationID:   1,
			Address:     "测试地址2",
			CreatedAt:   time.Date(2026, 1, 2, 18, 0, 0, 0, time.Local),
		},
		{
			RequestNo:   fmt.Sprintf("REQ-TREND-%d-3", time.Now().UnixNano()),
			UserID:      1,
			ServiceType: consts.ServiceTypeMedical,
			Status:      consts.RequestStatusCompleted,
			StationID:   1,
			Address:     "测试地址3",
			CreatedAt:   time.Date(2026, 1, 4, 9, 0, 0, 0, time.Local),
		},
	}

	for _, req := range requests {
		if err := repo.Create(req); err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
	}

	trend, err := repo.GetDailyTrendBetween(
		1,
		false,
		time.Date(2026, 1, 1, 0, 0, 0, 0, time.Local),
		time.Date(2026, 1, 3, 0, 0, 0, 0, time.Local),
	)
	if err != nil {
		t.Fatalf("failed to get daily trend between dates: %v", err)
	}

	if len(trend) != 1 {
		t.Fatalf("expected 1 trend point, got %d", len(trend))
	}
	if trend[0].Date != "2026-01-02" || trend[0].Count != 2 {
		t.Errorf("unexpected trend item: %#v", trend[0])
	}
}

func TestRequestRepository_WithTx(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRequestRepository(db)

	// 测试事务
	tx := db.Begin()
	txRepo := repo.WithTx(tx)

	req := &model.ServiceRequest{
		RequestNo:   fmt.Sprintf("REQ-%d", time.Now().UnixNano()),
		UserID:      1,
		ServiceType: consts.ServiceTypeMeal,
		Status:      consts.RequestStatusPending,
		Address:     "测试地址",
	}
	txRepo.Create(req)

	// 回滚事务
	tx.Rollback()

	// 验证需求未创建
	_, err := repo.GetByID(req.ID)
	if err == nil {
		t.Error("expected request to not exist after rollback")
	}
}
