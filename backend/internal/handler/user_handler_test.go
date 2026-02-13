//go:build integration

package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/repository"
	"community-elderly-care-platform/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// NOTE: 这些测试需要兼容的数据库环境
// 运行方式: go test -tags=integration ./...

func setupUserHandlerTestDB(t *testing.T) *gorm.DB {
	tmpFile := t.TempDir() + "/test.db"
	dsn := tmpFile + "?_loc=Local&_parseTime=true"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	err = db.AutoMigrate(&model.User{}, &model.UserIdentity{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	return db
}

var testUserIDCounter int64 = 1000

func createUserHandlerTestUser(t *testing.T, db *gorm.DB, phone string, identities []string, stationID int64) *model.User {
	t.Helper()
	hash, _ := service.HashPassword("Test@123")
	user := &model.User{
		ID:           atomic.AddInt64(&testUserIDCounter, 1),
		Phone:        phone,
		PasswordHash: hash,
		Name:         "Test User",
		Status:       "active",
		StationID:    stationID,
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	for i, identityType := range identities {
		identity := &model.UserIdentity{
			UserID:       user.ID,
			IdentityType: identityType,
			IsPrimary:    i == 0,
			StationID:    stationID,
			Status:       "active",
		}
		if err := db.Create(identity).Error; err != nil {
			t.Fatalf("failed to create test identity: %v", err)
		}
	}

	return user
}

func TestUserHandler_List_WithIdentities(t *testing.T) {
	db := setupUserHandlerTestDB(t)
	userRepo := repository.NewUserRepository(db)
	userIdentityRepo := repository.NewUserIdentityRepository(db)
	userService := service.NewUserService(userRepo, userIdentityRepo)
	handler := NewUserHandler(userService, "test-secret", "")

	stationID := int64(1)
	createUserHandlerTestUser(t, db, "13800000001", []string{"admin"}, stationID)
	createUserHandlerTestUser(t, db, "13800000002", []string{"staff", "station_manager"}, stationID)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/b/users?page=1&page_size=10", nil)

	handler.List(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestUserHandler_Create(t *testing.T) {
	db := setupUserHandlerTestDB(t)
	userRepo := repository.NewUserRepository(db)
	userIdentityRepo := repository.NewUserIdentityRepository(db)
	userService := service.NewUserService(userRepo, userIdentityRepo)
	handler := NewUserHandler(userService, "test-secret", "")

	reqBody := map[string]interface{}{
		"phone":         "13800000003",
		"password":      "Test@123",
		"name":          "New User",
		"identity_type": "staff",
		"station_id":    1,
		"status":        "active",
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/b/users", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Create(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestUserHandler_GetByID(t *testing.T) {
	db := setupUserHandlerTestDB(t)
	userRepo := repository.NewUserRepository(db)
	userIdentityRepo := repository.NewUserIdentityRepository(db)
	userService := service.NewUserService(userRepo, userIdentityRepo)
	handler := NewUserHandler(userService, "test-secret", "")

	stationID := int64(1)
	user := createUserHandlerTestUser(t, db, "13800000004", []string{"admin"}, stationID)
	if user.ID == 0 {
		t.Fatal("expected created user to have non-zero ID")
	}
	birthDate := time.Date(1990, 2, 10, 0, 0, 0, 0, time.UTC)
	if err := db.Model(&model.User{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
		"avatar":     "https://example.com/avatar.jpg",
		"gender":     "male",
		"birth_date": birthDate,
		"id_card":    "110101199002101234",
	}).Error; err != nil {
		t.Fatalf("failed to update user profile fields: %v", err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/b/users/"+toString(user.ID), nil)
	c.Params = gin.Params{{Key: "id", Value: toString(user.ID)}}

	handler.GetByID(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}
	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected response data object, got %T", resp["data"])
	}
	if got := data["avatar"]; got != "https://example.com/avatar.jpg" {
		t.Errorf("expected avatar to be returned, got %v", got)
	}
	if got := data["gender"]; got != "male" {
		t.Errorf("expected gender to be returned, got %v", got)
	}
	if got := data["birth_date"]; got != "1990-02-10" {
		t.Errorf("expected birth_date to be returned, got %v", got)
	}
	if _, ok := data["age"]; !ok {
		t.Error("expected age field to be returned")
	}
	if got := data["id_card_masked"]; got != "1101**********1234" {
		t.Errorf("expected id_card_masked to be returned, got %v", got)
	}
	if got, ok := data["id_card_hash"].(string); !ok || got == "" {
		t.Errorf("expected non-empty id_card_hash to be returned, got %v", data["id_card_hash"])
	}
	if got, ok := data["id_card_token"].(string); !ok || got == "" {
		t.Errorf("expected non-empty id_card_token to be returned, got %v", data["id_card_token"])
	}
}

func TestUserHandler_Update_InvalidIDCardToken(t *testing.T) {
	db := setupUserHandlerTestDB(t)
	userRepo := repository.NewUserRepository(db)
	userIdentityRepo := repository.NewUserIdentityRepository(db)
	userService := service.NewUserService(userRepo, userIdentityRepo)
	handler := NewUserHandler(userService, "test-secret", "")

	user := createUserHandlerTestUser(t, db, "13800000005", []string{"admin"}, 1)
	if user.ID == 0 {
		t.Fatal("expected created user to have non-zero ID")
	}
	if err := db.Model(&model.User{}).Where("id = ?", user.ID).Update("id_card", "110101199002101234").Error; err != nil {
		t.Fatalf("failed to set user id_card: %v", err)
	}

	reqBody := map[string]interface{}{
		"name":          "Updated Name",
		"id_card_token": "invalid.token",
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PUT", "/b/users/"+toString(user.ID), bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: toString(user.ID)}}

	handler.Update(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUserHandler_Update_InvalidIDCardHash(t *testing.T) {
	db := setupUserHandlerTestDB(t)
	userRepo := repository.NewUserRepository(db)
	userIdentityRepo := repository.NewUserIdentityRepository(db)
	userService := service.NewUserService(userRepo, userIdentityRepo)
	handler := NewUserHandler(userService, "test-secret", "")

	user := createUserHandlerTestUser(t, db, "13800000006", []string{"admin"}, 1)
	if user.ID == 0 {
		t.Fatal("expected created user to have non-zero ID")
	}
	if err := db.Model(&model.User{}).Where("id = ?", user.ID).Update("id_card", "110101199002101234").Error; err != nil {
		t.Fatalf("failed to set user id_card: %v", err)
	}

	reqBody := map[string]interface{}{
		"name":         "Updated Name",
		"id_card_hash": "invalid-hash",
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PUT", "/b/users/"+toString(user.ID), bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: toString(user.ID)}}

	handler.Update(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func toString(v int64) string {
	return strconv.FormatInt(v, 10)
}
