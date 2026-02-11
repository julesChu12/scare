//go:build integration

package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/repository"
	"community-elderly-care-platform/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// NOTE: 这些测试需要兼容的数据库环境
// 运行方式: go test -tags=integration ./...

func setupRoleHandlerTestDB(t *testing.T) *gorm.DB {
	tmpFile := t.TempDir() + "/test.db"
	dsn := tmpFile + "?_loc=Local&_parseTime=true"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	err = db.AutoMigrate(&model.User{}, &model.UserRole{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	return db
}

func createRoleHandlerTestUser(db *gorm.DB, phone string, roles []string, stationID int64) *model.User {
	hash, _ := service.HashPassword("Test@123")
	user := &model.User{
		Phone:        phone,
		PasswordHash: hash,
		Name:         "Test User",
		Status:       "active",
		StationID:    stationID,
	}
	db.Create(user)

	for i, role := range roles {
		userRole := &model.UserRole{
			UserID:    user.ID,
			Role:      role,
			IsPrimary: i == 0,
			Status:    "active",
		}
		db.Create(userRole)
	}

	return user
}

func TestRoleHandler_UpdateUserRoles(t *testing.T) {
	db := setupRoleHandlerTestDB(t)
	userRepo := repository.NewUserRepository(db)
	userRoleService := service.NewUserRoleService(db, userRepo, nil)
	permissionService := &service.PermissionService{}

	handler := NewRoleHandler(permissionService, userRoleService)

	stationID := int64(1)
	user := createRoleHandlerTestUser(db, "13800000001", []string{"staff"}, stationID)

	reqBody := map[string]interface{}{
		"roles": []string{"admin", "station_manager"},
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PUT", "/b/users/1/roles", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Set("user_id", int64(999)) // operator ID

	// Use actual user ID
	c.Params = gin.Params{{Key: "id", Value: string(rune(user.ID + '0'))}}

	handler.UpdateUserRoles(c)

	// Reset and try with proper ID format
	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Request, _ = http.NewRequest("PUT", "/b/users/1/roles", bytes.NewBuffer(body))
	c2.Request.Header.Set("Content-Type", "application/json")
	c2.Params = gin.Params{{Key: "id", Value: "1"}}
	c2.Set("user_id", int64(999))

	handler.UpdateUserRoles(c2)

	if w2.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w2.Code)
	}

	// Verify response
	var resp map[string]interface{}
	json.Unmarshal(w2.Body.Bytes(), &resp)

	if resp["code"].(float64) != 0 {
		t.Errorf("expected code 0, got %v", resp["code"])
	}

	data := resp["data"].(map[string]interface{})
	if data["tokens_revoked"] != true {
		t.Error("expected tokens_revoked to be true")
	}

	// Verify roles in database
	var roles []model.UserRole
	db.Where("user_id = ? AND status = ?", 1, "active").Find(&roles)

	if len(roles) != 2 {
		t.Errorf("expected 2 active roles, got %d", len(roles))
	}
}

func TestRoleHandler_UpdateUserRoles_MultipleRoles(t *testing.T) {
	db := setupRoleHandlerTestDB(t)
	userRepo := repository.NewUserRepository(db)
	userRoleService := service.NewUserRoleService(db, userRepo, nil)
	permissionService := &service.PermissionService{}

	handler := NewRoleHandler(permissionService, userRoleService)

	stationID := int64(1)
	createRoleHandlerTestUser(db, "13800000001", []string{}, stationID)

	reqBody := map[string]interface{}{
		"roles": []string{"admin", "staff", "station_manager"},
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PUT", "/b/users/1/roles", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Set("user_id", int64(999))

	handler.UpdateUserRoles(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify all 3 roles are assigned
	var roles []model.UserRole
	db.Where("user_id = ? AND status = ?", 1, "active").Find(&roles)

	if len(roles) != 3 {
		t.Errorf("expected 3 active roles, got %d", len(roles))
	}

	// Verify first role is primary
	var primaryRole model.UserRole
	db.Where("user_id = ? AND is_primary = ?", 1, true).First(&primaryRole)

	if primaryRole.Role != "admin" {
		t.Errorf("expected primary role 'admin', got '%s'", primaryRole.Role)
	}
}

func TestRoleHandler_UpdateUserRoles_InvalidUserID(t *testing.T) {
	db := setupRoleHandlerTestDB(t)
	userRepo := repository.NewUserRepository(db)
	userRoleService := service.NewUserRoleService(db, userRepo, nil)
	permissionService := &service.PermissionService{}

	handler := NewRoleHandler(permissionService, userRoleService)

	reqBody := map[string]interface{}{
		"roles": []string{"admin"},
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PUT", "/b/users/invalid/roles", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "invalid"}}

	handler.UpdateUserRoles(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestRoleHandler_UpdateUserRoles_EmptyUserID(t *testing.T) {
	db := setupRoleHandlerTestDB(t)
	userRepo := repository.NewUserRepository(db)
	userRoleService := service.NewUserRoleService(db, userRepo, nil)
	permissionService := &service.PermissionService{}

	handler := NewRoleHandler(permissionService, userRoleService)

	reqBody := map[string]interface{}{
		"roles": []string{"admin"},
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PUT", "/b/users//roles", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: ""}}

	handler.UpdateUserRoles(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestRoleHandler_UpdateUserRoles_InvalidPayload(t *testing.T) {
	db := setupRoleHandlerTestDB(t)
	userRepo := repository.NewUserRepository(db)
	userRoleService := service.NewUserRoleService(db, userRepo, nil)
	permissionService := &service.PermissionService{}

	handler := NewRoleHandler(permissionService, userRoleService)

	stationID := int64(1)
	createRoleHandlerTestUser(db, "13800000001", []string{"staff"}, stationID)

	// Missing roles field
	reqBody := map[string]interface{}{}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PUT", "/b/users/1/roles", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	handler.UpdateUserRoles(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestRoleHandler_UpdateUserRoles_UserNotFound(t *testing.T) {
	db := setupRoleHandlerTestDB(t)
	userRepo := repository.NewUserRepository(db)
	userRoleService := service.NewUserRoleService(db, userRepo, nil)
	permissionService := &service.PermissionService{}

	handler := NewRoleHandler(permissionService, userRoleService)

	reqBody := map[string]interface{}{
		"roles": []string{"admin"},
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PUT", "/b/users/99999/roles", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "99999"}}
	c.Set("user_id", int64(999))

	handler.UpdateUserRoles(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestRoleHandler_GetRolePermissions_EmptyRole(t *testing.T) {
	db := setupRoleHandlerTestDB(t)
	userRepo := repository.NewUserRepository(db)
	userRoleService := service.NewUserRoleService(db, userRepo, nil)
	permissionService := &service.PermissionService{}

	handler := NewRoleHandler(permissionService, userRoleService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/b/roles//permissions", nil)
	c.Params = gin.Params{{Key: "role", Value: ""}}

	handler.GetRolePermissions(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestRoleHandler_UpdateRolePermissions_EmptyRole(t *testing.T) {
	db := setupRoleHandlerTestDB(t)
	userRepo := repository.NewUserRepository(db)
	userRoleService := service.NewUserRoleService(db, userRepo, nil)
	permissionService := &service.PermissionService{}

	handler := NewRoleHandler(permissionService, userRoleService)

	reqBody := map[string]interface{}{
		"permissions": []string{"perm_b_tasks_pool_get"},
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PUT", "/b/roles//permissions", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "role", Value: ""}}

	handler.UpdateRolePermissions(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestRoleHandler_UpdateRolePermissions_InvalidPayload(t *testing.T) {
	db := setupRoleHandlerTestDB(t)
	userRepo := repository.NewUserRepository(db)
	userRoleService := service.NewUserRoleService(db, userRepo, nil)
	permissionService := &service.PermissionService{}

	handler := NewRoleHandler(permissionService, userRoleService)

	// Missing permissions field
	reqBody := map[string]interface{}{}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PUT", "/b/roles/admin/permissions", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "role", Value: "admin"}}

	handler.UpdateRolePermissions(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestRoleHandler_UpdateUserRoles_TokenRevoked(t *testing.T) {
	db := setupRoleHandlerTestDB(t)
	userRepo := repository.NewUserRepository(db)
	// Note: blacklistService is nil, so token revocation won't actually happen
	// but the response should still indicate tokens_revoked: true
	userRoleService := service.NewUserRoleService(db, userRepo, nil)
	permissionService := &service.PermissionService{}

	handler := NewRoleHandler(permissionService, userRoleService)

	stationID := int64(1)
	createRoleHandlerTestUser(db, "13800000001", []string{"staff"}, stationID)

	reqBody := map[string]interface{}{
		"roles": []string{"admin"},
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PUT", "/b/users/1/roles", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Set("user_id", int64(999))

	handler.UpdateUserRoles(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	data := resp["data"].(map[string]interface{})
	if data["tokens_revoked"] != true {
		t.Error("expected tokens_revoked to be true in response")
	}
}
