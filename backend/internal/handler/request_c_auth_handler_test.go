package handler

import (
	"net/http"
	"testing"

	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/repository"
	"community-elderly-care-platform/internal/service"
	"community-elderly-care-platform/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func TestRequestHandler_List_BEndStationManagerUsesOwnStationScope(t *testing.T) {
	db := openHandlerTestDB(t, "request_handler_list_scope.db")
	createHandlerTables(t, db)
	stationA := seedHandlerStation(t, db, "站点A")
	stationB := seedHandlerStation(t, db, "站点B")
	reqA := seedHandlerRequest(t, db, 10, stationA.ID, "REQ-A")
	seedHandlerRequest(t, db, 11, stationB.ID, "REQ-B")

	reqRepo := repository.NewRequestRepository(db)
	handler := NewRequestHandler(service.NewRequestService(
		db,
		reqRepo,
		repository.NewTaskRepository(db),
		repository.NewStationRepository(db),
		nil,
		nil,
		nil,
		nil,
	))

	c, w := newJSONTestContext(t, http.MethodGet, "/b/requests?station_id=2&page=1&page_size=10", nil)
	setBEndClaims(c, 2, stationA.ID, "station_manager")

	handler.List(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	resp := decodeResponseMap(t, w)
	data := resp["data"].(map[string]any)
	if int(data["total"].(float64)) != 1 {
		t.Fatalf("expected total 1, got %v", data["total"])
	}
	items := data["items"].([]any)
	item := items[0].(map[string]any)
	if int64(item["id"].(float64)) != reqA.ID {
		t.Fatalf("expected request %d, got %v", reqA.ID, item["id"])
	}
}

func TestRequestHandler_Get_ForbiddenForForeignCEndUser(t *testing.T) {
	db := openHandlerTestDB(t, "request_handler_get_owner.db")
	createHandlerTables(t, db)
	station := seedHandlerStation(t, db, "站点A")
	req := seedHandlerRequest(t, db, 10, station.ID, "REQ-OWNER")

	handler := NewRequestHandler(service.NewRequestService(
		db,
		repository.NewRequestRepository(db),
		repository.NewTaskRepository(db),
		repository.NewStationRepository(db),
		nil,
		nil,
		nil,
		nil,
	))

	c, w := newJSONTestContext(t, http.MethodGet, "/c/requests/1", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	setCEndClaims(c, 11)

	handler.Get(c)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
	_ = req
}

func TestRequestHandler_Create_RequiresAuthenticatedUser(t *testing.T) {
	handler := NewRequestHandler(nil)
	c, w := newJSONTestContext(t, http.MethodPost, "/c/requests", map[string]any{
		"service_type": "meal",
	})

	handler.Create(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestTaskHandler_ListPool_NonAdminUsesOwnStationScope(t *testing.T) {
	db := openHandlerTestDB(t, "task_handler_pool_scope.db")
	createHandlerTables(t, db)
	stationA := seedHandlerStation(t, db, "站点A")
	stationB := seedHandlerStation(t, db, "站点B")
	reqA := seedHandlerRequest(t, db, 10, stationA.ID, "REQ-TASK-A")
	reqB := seedHandlerRequest(t, db, 11, stationB.ID, "REQ-TASK-B")

	if err := db.Create(&model.TaskAssignment{RequestID: reqA.ID, StationID: stationA.ID, Status: "dispatched"}).Error; err != nil {
		t.Fatalf("failed to create task A: %v", err)
	}
	if err := db.Create(&model.TaskAssignment{RequestID: reqB.ID, StationID: stationB.ID, Status: "dispatched"}).Error; err != nil {
		t.Fatalf("failed to create task B: %v", err)
	}

	handler := NewTaskHandler(service.NewTaskService(
		db,
		repository.NewTaskRepository(db),
		repository.NewRequestRepository(db),
		nil,
	))
	c, w := newJSONTestContext(t, http.MethodGet, "/b/tasks/pool?station_id=2&page=1&page_size=10", nil)
	setBEndClaims(c, 3, stationA.ID, "station_manager")

	handler.ListPool(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	resp := decodeResponseMap(t, w)
	data := resp["data"].(map[string]any)
	if int(data["total"].(float64)) != 1 {
		t.Fatalf("expected total 1, got %v", data["total"])
	}
	items := data["items"].([]any)
	item := items[0].(map[string]any)
	if int64(item["station_id"].(float64)) != stationA.ID {
		t.Fatalf("expected station %d, got %v", stationA.ID, item["station_id"])
	}
}

func TestCAuthHandler_QuickStart_RequiresAddressOrLocation(t *testing.T) {
	db := openHandlerTestDB(t, "c_auth_quickstart_addr.db")
	createHandlerTables(t, db)

	smsSvc := service.NewSMSService(nil, "development")
	authSvc := service.NewAuthService(
		repository.NewUserRepository(db),
		repository.NewUserIdentityRepository(db),
		repository.NewCustomerRepository(db),
		jwt.NewManager("test-secret", 1, 2),
		smsSvc,
		db,
	)
	authSvc.SetStationRepo(repository.NewStationRepository(db))
	authSvc.SetGeofenceService(&service.GeofenceService{})

	handler := NewCAuthHandler(authSvc, repository.NewUserRepository(db), repository.NewCustomerRepository(db), smsSvc)
	c, w := newJSONTestContext(t, http.MethodPost, "/c/auth/quick-start", map[string]any{
		"phone":        "13800138000",
		"code":         "000000",
		"name":         "测试用户",
		"service_type": "meal",
	})

	handler.QuickStart(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	resp := decodeResponseMap(t, w)
	if resp["msg"] != "服务地址不能为空" {
		t.Fatalf("unexpected message: %v", resp["msg"])
	}
}

func TestCAuthHandler_CheckToken_ReturnsUserAndProfile(t *testing.T) {
	db := openHandlerTestDB(t, "c_auth_check_token.db")
	createHandlerTables(t, db)
	seedHandlerUserAndProfile(t, db, 10, "13900000001", "张大爷", "elderly", "测试地址")

	userRepo := repository.NewUserRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	authService := service.NewAuthService(
		userRepo,
		repository.NewUserIdentityRepository(db),
		customerRepo,
		jwt.NewManager("test-secret", 1, 2),
		service.NewSMSService(nil, "development"),
		db,
	)
	handler := NewCAuthHandler(authService, userRepo, customerRepo, service.NewSMSService(nil, "development"))

	c, w := newJSONTestContext(t, http.MethodGet, "/c/auth/check", nil)
	setCEndClaims(c, 10)

	handler.CheckToken(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	resp := decodeResponseMap(t, w)
	data := resp["data"].(map[string]any)
	user := data["user"].(map[string]any)
	if user["phone"] != "13900000001" {
		t.Fatalf("expected phone 13900000001, got %v", user["phone"])
	}
	profile := data["profile"].(map[string]any)
	if profile["address"] != "测试地址" || profile["user_type"] != "elderly" {
		t.Fatalf("unexpected profile payload: %+v", profile)
	}
}
