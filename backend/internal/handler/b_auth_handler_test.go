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
	"community-elderly-care-platform/pkg/jwt"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// NOTE: 这些测试需要兼容的数据库环境
// 运行方式: go test -tags=integration ./...

func init() {
	gin.SetMode(gin.TestMode)
}

func setupBAuthTestDB(t *testing.T) *gorm.DB {
	tmpFile := t.TempDir() + "/test.db"
	dsn := tmpFile + "?_loc=Local&_parseTime=true"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	err = db.AutoMigrate(&model.User{}, &model.UserIdentity{}, &model.CustomerProfile{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	return db
}

func createBAuthTestUser(db *gorm.DB, phone string, identities []string, stationID int64) *model.User {
	hash, _ := service.HashPassword("Test@123")
	user := &model.User{
		Phone:        phone,
		PasswordHash: hash,
		Name:         "Test User",
		Status:       "active",
		StationID:    stationID,
	}
	db.Create(user)

	// Create identities
	for i, identityType := range identities {
		identity := &model.UserIdentity{
			UserID:       user.ID,
			IdentityType: identityType,
			IsPrimary:    i == 0,
			StationID:    stationID,
			Status:       "active",
		}
		db.Create(identity)
	}

	return user
}

func TestBAuthHandler_Login_Success(t *testing.T) {
	db := setupBAuthTestDB(t)
	userRepo := repository.NewUserRepository(db)
	userIdentityRepo := repository.NewUserIdentityRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	jwtManager := jwt.NewManager("test-secret", 1, 2)
	authService := service.NewAuthService(userRepo, userIdentityRepo, customerRepo, jwtManager, nil, db)
	permissionService := &service.PermissionService{} // minimal for test

	handler := NewBAuthHandler(authService, userRepo, userIdentityRepo, permissionService)

	stationID := int64(1)
	createBAuthTestUser(db, "13800000001", []string{"admin"}, stationID)

	// Create request
	reqBody := map[string]string{
		"phone":    "13800000001",
		"password": "Test@123",
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/b/auth/login", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Login(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp Response
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp.Msg != "ok" {
		t.Errorf("expected msg 'ok', got '%s'", resp.Msg)
	}
}

func TestBAuthHandler_Login_InvalidCredentials(t *testing.T) {
	db := setupBAuthTestDB(t)
	userRepo := repository.NewUserRepository(db)
	userIdentityRepo := repository.NewUserIdentityRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	jwtManager := jwt.NewManager("test-secret", 1, 2)
	authService := service.NewAuthService(userRepo, userIdentityRepo, customerRepo, jwtManager, nil, db)
	permissionService := &service.PermissionService{}

	handler := NewBAuthHandler(authService, userRepo, userIdentityRepo, permissionService)

	stationID := int64(1)
	createBAuthTestUser(db, "13800000002", []string{"admin"}, stationID)

	// Create request with wrong password
	reqBody := map[string]string{
		"phone":    "13800000002",
		"password": "WrongPassword",
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/b/auth/login", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Login(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestBAuthHandler_Login_NoIdentity(t *testing.T) {
	db := setupBAuthTestDB(t)
	userRepo := repository.NewUserRepository(db)
	userIdentityRepo := repository.NewUserIdentityRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	jwtManager := jwt.NewManager("test-secret", 1, 2)
	authService := service.NewAuthService(userRepo, userIdentityRepo, customerRepo, jwtManager, nil, db)
	permissionService := &service.PermissionService{}

	handler := NewBAuthHandler(authService, userRepo, userIdentityRepo, permissionService)

	// Create user without identities
	hash, _ := service.HashPassword("Test@123")
	user := &model.User{
		Phone:        "13800000003",
		PasswordHash: hash,
		Name:         "Test User",
		Status:       "active",
	}
	db.Create(user)

	reqBody := map[string]string{
		"phone":    "13800000003",
		"password": "Test@123",
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/b/auth/login", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Login(c)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected status %d, got %d", http.StatusForbidden, w.Code)
	}
}
