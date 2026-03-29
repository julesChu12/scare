//go:build integration

package service

import (
	"testing"

	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/repository"
	"community-elderly-care-platform/pkg/jwt"
	"gorm.io/gorm"
)

// NOTE: 这些测试需要兼容的数据库环境
// 运行方式: go test -tags=integration ./...

func setupAuthTestDB(t *testing.T) *gorm.DB {
	db := openServiceTestDB(t, "auth_service_integration_test.db")
	createUsersTable(t, db)
	createUserIdentitiesTable(t, db)
	createCustomerProfilesTable(t, db)
	return db
}

func createTestUser(db *gorm.DB, phone string, stationID int64) *model.User {
	hash, _ := HashPassword("Test@123")
	user := &model.User{
		Phone:        phone,
		PasswordHash: hash,
		Name:         "Test User",
		Status:       "active",
		StationID:    stationID,
	}
	db.Create(user)
	return user
}

func createTestBEndUser(db *gorm.DB, phone string, identityType string, stationID int64) *model.User {
	user := createTestUser(db, phone, stationID)

	// Create B-end identity
	identity := &model.UserIdentity{
		UserID:       user.ID,
		IdentityType: identityType,
		IsPrimary:    true,
		StationID:    stationID,
		Status:       "active",
	}
	db.Create(identity)

	return user
}

func createTestCEndUser(db *gorm.DB, phone string, customerType string) *model.User {
	user := createTestUser(db, phone, 0)

	// Create customer profile
	profile := &model.CustomerProfile{
		UserID:       user.ID,
		CustomerType: customerType,
	}
	db.Create(profile)

	return user
}

func createTestCrossEndUser(db *gorm.DB, phone string, bEndIdentity string, customerType string, stationID int64) *model.User {
	user := createTestUser(db, phone, stationID)

	// Create B-end identity (not primary since user is primarily C-end)
	identity := &model.UserIdentity{
		UserID:       user.ID,
		IdentityType: bEndIdentity,
		IsPrimary:    false,
		StationID:    stationID,
		Status:       "active",
	}
	db.Create(identity)

	// Create customer profile
	profile := &model.CustomerProfile{
		UserID:       user.ID,
		CustomerType: customerType,
	}
	db.Create(profile)

	return user
}

func TestAuthService_LoginBEnd_Success(t *testing.T) {
	db := setupAuthTestDB(t)
	userRepo := repository.NewUserRepository(db)
	userIdentityRepo := repository.NewUserIdentityRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	jwtManager := jwt.NewManager("test-secret", 1, 2)
	authService := NewAuthService(userRepo, userIdentityRepo, customerRepo, jwtManager, nil, db)

	stationID := int64(1)
	createTestBEndUser(db, "13800000001", "admin", stationID)

	tokens, user, err := authService.LoginBEnd("13800000001", "Test@123")
	if err != nil {
		t.Fatalf("LoginBEnd failed: %v", err)
	}

	if tokens == nil {
		t.Fatal("expected tokens to be returned")
	}

	if tokens.AccessToken == "" {
		t.Error("expected access token to be set")
	}

	if tokens.RefreshToken == "" {
		t.Error("expected refresh token to be set")
	}

	if user == nil {
		t.Fatal("expected user to be returned")
	}

	if user.Phone != "13800000001" {
		t.Errorf("expected phone 13800000001, got %s", user.Phone)
	}
}

func TestAuthService_LoginBEnd_NoIdentity(t *testing.T) {
	db := setupAuthTestDB(t)
	userRepo := repository.NewUserRepository(db)
	userIdentityRepo := repository.NewUserIdentityRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	jwtManager := jwt.NewManager("test-secret", 1, 2)
	authService := NewAuthService(userRepo, userIdentityRepo, customerRepo, jwtManager, nil, db)

	// Create user without B-end identity
	createTestUser(db, "13800000002", 0)

	_, _, err := authService.LoginBEnd("13800000002", "Test@123")
	if err != ErrNoRoleForBEnd {
		t.Errorf("expected ErrNoRoleForBEnd, got %v", err)
	}
}

func TestAuthService_LoginBEnd_InvalidPassword(t *testing.T) {
	db := setupAuthTestDB(t)
	userRepo := repository.NewUserRepository(db)
	userIdentityRepo := repository.NewUserIdentityRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	jwtManager := jwt.NewManager("test-secret", 1, 2)
	authService := NewAuthService(userRepo, userIdentityRepo, customerRepo, jwtManager, nil, db)

	stationID := int64(1)
	createTestBEndUser(db, "13800000003", "admin", stationID)

	_, _, err := authService.LoginBEnd("13800000003", "WrongPassword")
	if err != ErrInvalidCredentials {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthService_LoginBEnd_UserNotFound(t *testing.T) {
	db := setupAuthTestDB(t)
	userRepo := repository.NewUserRepository(db)
	userIdentityRepo := repository.NewUserIdentityRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	jwtManager := jwt.NewManager("test-secret", 1, 2)
	authService := NewAuthService(userRepo, userIdentityRepo, customerRepo, jwtManager, nil, db)

	_, _, err := authService.LoginBEnd("13899999999", "Test@123")
	if err != ErrInvalidCredentials {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthService_LoginBEnd_UserInactive(t *testing.T) {
	db := setupAuthTestDB(t)
	userRepo := repository.NewUserRepository(db)
	userIdentityRepo := repository.NewUserIdentityRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	jwtManager := jwt.NewManager("test-secret", 1, 2)
	authService := NewAuthService(userRepo, userIdentityRepo, customerRepo, jwtManager, nil, db)

	// Create inactive user with B-end identity
	hash, _ := HashPassword("Test@123")
	user := &model.User{
		Phone:        "13800000004",
		PasswordHash: hash,
		Name:         "Inactive User",
		Status:       "inactive",
		StationID:    1,
	}
	db.Create(user)

	identity := &model.UserIdentity{
		UserID:       user.ID,
		IdentityType: "admin",
		IsPrimary:    true,
		StationID:    1,
		Status:       "active",
	}
	db.Create(identity)

	_, _, err := authService.LoginBEnd("13800000004", "Test@123")
	if err != ErrUserInactive {
		t.Errorf("expected ErrUserInactive, got %v", err)
	}
}

func TestAuthService_LoginCEnd_Success(t *testing.T) {
	db := setupAuthTestDB(t)
	userRepo := repository.NewUserRepository(db)
	userIdentityRepo := repository.NewUserIdentityRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	jwtManager := jwt.NewManager("test-secret", 1, 2)
	authService := NewAuthService(userRepo, userIdentityRepo, customerRepo, jwtManager, nil, db)

	createTestCEndUser(db, "13800000005", "elderly")

	tokens, user, err := authService.LoginCEnd("13800000005", "Test@123")
	if err != nil {
		t.Fatalf("LoginCEnd failed: %v", err)
	}

	if tokens == nil {
		t.Fatal("expected tokens to be returned")
	}

	if tokens.AccessToken == "" {
		t.Error("expected access token to be set")
	}

	if user == nil {
		t.Fatal("expected user to be returned")
	}
}

func TestAuthService_LoginCEnd_NoProfile(t *testing.T) {
	db := setupAuthTestDB(t)
	userRepo := repository.NewUserRepository(db)
	userIdentityRepo := repository.NewUserIdentityRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	jwtManager := jwt.NewManager("test-secret", 1, 2)
	authService := NewAuthService(userRepo, userIdentityRepo, customerRepo, jwtManager, nil, db)

	// Create user without customer profile
	createTestUser(db, "13800000006", 0)

	_, _, err := authService.LoginCEnd("13800000006", "Test@123")
	if err != ErrNoCustomerProfile {
		t.Errorf("expected ErrNoCustomerProfile, got %v", err)
	}
}

func TestAuthService_Refresh(t *testing.T) {
	db := setupAuthTestDB(t)
	userRepo := repository.NewUserRepository(db)
	userIdentityRepo := repository.NewUserIdentityRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	jwtManager := jwt.NewManager("test-secret", 1, 2)
	authService := NewAuthService(userRepo, userIdentityRepo, customerRepo, jwtManager, nil, db)

	stationID := int64(1)
	createTestBEndUser(db, "13800000007", "admin", stationID)

	tokens, _, err := authService.LoginBEnd("13800000007", "Test@123")
	if err != nil {
		t.Fatalf("LoginBEnd failed: %v", err)
	}

	// Refresh tokens
	newTokens, err := authService.Refresh(tokens.RefreshToken)
	if err != nil {
		t.Fatalf("Refresh failed: %v", err)
	}

	if newTokens.AccessToken == "" {
		t.Error("expected new access token to be set")
	}

	if newTokens.RefreshToken == "" {
		t.Error("expected new refresh token to be set")
	}
}

func TestAuthService_CrossEndUser(t *testing.T) {
	db := setupAuthTestDB(t)
	userRepo := repository.NewUserRepository(db)
	userIdentityRepo := repository.NewUserIdentityRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	jwtManager := jwt.NewManager("test-secret", 1, 2)
	authService := NewAuthService(userRepo, userIdentityRepo, customerRepo, jwtManager, nil, db)

	stationID := int64(1)
	createTestCrossEndUser(db, "13800000008", "staff", "elderly", stationID)

	// B-end login should succeed
	bTokens, _, err := authService.LoginBEnd("13800000008", "Test@123")
	if err != nil {
		t.Fatalf("B-end login failed: %v", err)
	}
	if bTokens == nil {
		t.Error("expected B-end tokens")
	}

	// C-end login should also succeed
	cTokens, _, err := authService.LoginCEnd("13800000008", "Test@123")
	if err != nil {
		t.Fatalf("C-end login failed: %v", err)
	}
	if cTokens == nil {
		t.Error("expected C-end tokens")
	}
}
