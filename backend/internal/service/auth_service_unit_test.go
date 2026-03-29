package service

import (
	"errors"
	"testing"
	"time"

	"community-elderly-care-platform/internal/consts"
	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/repository"
	"community-elderly-care-platform/pkg/jwt"

	"gorm.io/gorm"
)

func setupAuthUnitService(t *testing.T) (*AuthService, *gorm.DB, *jwt.Manager) {
	t.Helper()

	db := openServiceTestDB(t, "auth_service_unit_test.db")
	createUsersTable(t, db)
	createUserIdentitiesTable(t, db)
	createCustomerProfilesTable(t, db)
	createServiceStationsTable(t, db)
	createServiceRequestsTable(t, db)
	createTaskAssignmentsTable(t, db)

	jwtManager := jwt.NewManager("test-secret", 1, 2)
	authSvc := NewAuthService(
		repository.NewUserRepository(db),
		repository.NewUserIdentityRepository(db),
		repository.NewCustomerRepository(db),
		jwtManager,
		NewSMSService(nil, "development"),
		db,
	)
	authSvc.SetStationRepo(repository.NewStationRepository(db))
	authSvc.SetGeofenceService(&GeofenceService{})

	return authSvc, db, jwtManager
}

func seedAuthUnitUser(t *testing.T, db *gorm.DB, phone, password, status string, stationID int64) *model.User {
	t.Helper()

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	user := &model.User{
		Phone:        phone,
		PasswordHash: hash,
		Name:         "测试用户",
		Status:       status,
		StationID:    stationID,
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	return user
}

func seedAuthUnitIdentity(t *testing.T, db *gorm.DB, userID int64, identityType string, isPrimary bool, stationID int64) *model.UserIdentity {
	t.Helper()

	identity := &model.UserIdentity{
		UserID:       userID,
		IdentityType: identityType,
		IsPrimary:    isPrimary,
		StationID:    stationID,
		Status:       "active",
		GrantedAt:    time.Now(),
	}
	if err := db.Create(identity).Error; err != nil {
		t.Fatalf("failed to create identity: %v", err)
	}

	return identity
}

func seedAuthUnitProfile(t *testing.T, db *gorm.DB, userID int64, address string) *model.CustomerProfile {
	t.Helper()

	profile := &model.CustomerProfile{
		UserID:           userID,
		Address:          address,
		CustomerType:     consts.IdentityElderly,
		EmergencyContact: `{}`,
	}
	if err := db.Create(profile).Error; err != nil {
		t.Fatalf("failed to create customer profile: %v", err)
	}

	return profile
}

func TestAuthService_LoginBEnd_FallsBackToFirstBEndIdentity(t *testing.T) {
	authSvc, db, jwtManager := setupAuthUnitService(t)
	user := seedAuthUnitUser(t, db, "13820000001", "Test@123", "active", 9)
	seedAuthUnitIdentity(t, db, user.ID, consts.IdentityFamily, true, 0)
	seedAuthUnitIdentity(t, db, user.ID, consts.IdentityStaff, false, 9)

	tokens, gotUser, err := authSvc.LoginBEnd("13820000001", "Test@123")
	if err != nil {
		t.Fatalf("LoginBEnd returned error: %v", err)
	}
	if gotUser.ID != user.ID {
		t.Fatalf("expected user %d, got %d", user.ID, gotUser.ID)
	}

	claims, err := jwtManager.ParseToken(tokens.AccessToken)
	if err != nil {
		t.Fatalf("failed to parse access token: %v", err)
	}
	if claims.Primary != consts.IdentityStaff {
		t.Fatalf("expected primary %s, got %s", consts.IdentityStaff, claims.Primary)
	}
	if claims.Type != "b_end" || claims.StationID != 9 {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}

func TestAuthService_LoginCEnd_FallsBackToFirstCEndIdentity(t *testing.T) {
	authSvc, db, jwtManager := setupAuthUnitService(t)
	user := seedAuthUnitUser(t, db, "13820000002", "Test@123", "active", 0)
	seedAuthUnitProfile(t, db, user.ID, "测试地址")
	seedAuthUnitIdentity(t, db, user.ID, consts.IdentityAdmin, true, 1)
	seedAuthUnitIdentity(t, db, user.ID, consts.IdentityElderly, false, 0)

	tokens, gotUser, err := authSvc.LoginCEnd("13820000002", "Test@123")
	if err != nil {
		t.Fatalf("LoginCEnd returned error: %v", err)
	}
	if gotUser.ID != user.ID {
		t.Fatalf("expected user %d, got %d", user.ID, gotUser.ID)
	}

	claims, err := jwtManager.ParseToken(tokens.AccessToken)
	if err != nil {
		t.Fatalf("failed to parse access token: %v", err)
	}
	if claims.Primary != consts.IdentityElderly {
		t.Fatalf("expected primary %s, got %s", consts.IdentityElderly, claims.Primary)
	}
	if claims.Type != "c_end" || claims.StationID != 0 {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}

func TestAuthService_LoginCEndByCode_UserNotFound(t *testing.T) {
	authSvc, _, _ := setupAuthUnitService(t)

	_, _, err := authSvc.LoginCEndByCode("13820000003", "000000")
	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestAuthService_Refresh_InvalidToken(t *testing.T) {
	authSvc, _, _ := setupAuthUnitService(t)

	if _, err := authSvc.Refresh("invalid-token"); !errors.Is(err, jwt.ErrInvalidToken) {
		t.Fatalf("expected jwt.ErrInvalidToken, got %v", err)
	}
}

func TestAuthService_QuickStart_CreatesUserProfileRequestByNearestStation(t *testing.T) {
	authSvc, db, _ := setupAuthUnitService(t)
	createTestStation(t, db, 75, 1)
	expected := createTestStation(t, db, 80, 21)

	lat := 80.0
	lng := 1.0
	description := "需要送餐上门"
	result, err := authSvc.QuickStart(QuickStartInput{
		Phone:        "13820000004",
		Code:         "000000",
		Name:         "快速用户",
		Address:      "高纬地址",
		Latitude:     &lat,
		Longitude:    &lng,
		ServiceType:  consts.ServiceTypeMeal,
		Description:  &description,
		Images:       []string{"a.png"},
		ContactName:  "联系人",
		ContactPhone: "13800138000",
	})
	if err != nil {
		t.Fatalf("QuickStart returned error: %v", err)
	}
	if result.User == nil || result.Profile == nil || result.Request == nil {
		t.Fatalf("expected user/profile/request to be created, got %+v", result)
	}
	if result.Request.StationID != expected.ID {
		t.Fatalf("expected station %d, got %d", expected.ID, result.Request.StationID)
	}
	if result.Request.Status != consts.RequestStatusDispatched {
		t.Fatalf("expected request status dispatched, got %s", result.Request.Status)
	}
	if result.Profile.Address != "高纬地址" {
		t.Fatalf("expected profile address updated, got %q", result.Profile.Address)
	}

	var identity model.UserIdentity
	if err := db.Where("user_id = ? AND identity_type = ?", result.User.ID, consts.IdentityElderly).First(&identity).Error; err != nil {
		t.Fatalf("expected elderly identity to be created: %v", err)
	}

	var task model.TaskAssignment
	if err := db.Where("request_id = ?", result.Request.ID).First(&task).Error; err != nil {
		t.Fatalf("expected task to be created: %v", err)
	}
	if task.StationID != expected.ID || task.Status != consts.TaskStatusDispatched {
		t.Fatalf("unexpected task payload: %+v", task)
	}
}

func TestAuthService_QuickStart_UpdatesExistingUserAndUsesFirstActiveStationWithoutCoords(t *testing.T) {
	authSvc, db, _ := setupAuthUnitService(t)
	expected := createTestStation(t, db, 30, 120)
	createTestStation(t, db, 31, 121)

	user := seedAuthUnitUser(t, db, "13820000005", "Test@123", "active", 0)
	seedAuthUnitProfile(t, db, user.ID, "旧地址")
	seedAuthUnitIdentity(t, db, user.ID, consts.IdentityElderly, true, 0)

	result, err := authSvc.QuickStart(QuickStartInput{
		Phone:        "13820000005",
		Code:         "000000",
		Name:         "新姓名",
		Address:      "新地址",
		ServiceType:  consts.ServiceTypeMedical,
		ContactName:  "新联系人",
		ContactPhone: "13800138001",
	})
	if err != nil {
		t.Fatalf("QuickStart returned error: %v", err)
	}

	if result.User.Name != "新姓名" {
		t.Fatalf("expected user name updated, got %q", result.User.Name)
	}
	if result.Profile.Address != "新地址" {
		t.Fatalf("expected profile address updated, got %q", result.Profile.Address)
	}
	if result.Request.StationID != expected.ID {
		t.Fatalf("expected first active station %d, got %d", expected.ID, result.Request.StationID)
	}
}

func TestAuthService_QuickStart_ReturnsNoStationWhenNoneAvailable(t *testing.T) {
	authSvc, _, _ := setupAuthUnitService(t)

	_, err := authSvc.QuickStart(QuickStartInput{
		Phone:       "13820000006",
		Code:        "000000",
		Name:        "无站点用户",
		Address:     "未知地址",
		ServiceType: consts.ServiceTypeMeal,
	})
	if !errors.Is(err, ErrNoStation) {
		t.Fatalf("expected ErrNoStation, got %v", err)
	}
}
