package service

import (
	"database/sql"
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

func TestAuthService_LoginCEnd_PasswordNotSet(t *testing.T) {
	authSvc, db, _ := setupAuthUnitService(t)
	user := &model.User{
		Phone:  "13820000009",
		Name:   "未设密码用户",
		Status: "active",
	}
	if err := db.Omit("PasswordHash").Create(user).Error; err != nil {
		t.Fatalf("failed to create user without password: %v", err)
	}
	seedAuthUnitProfile(t, db, user.ID, "测试地址")
	seedAuthUnitIdentity(t, db, user.ID, consts.IdentityElderly, true, 0)

	_, _, err := authSvc.LoginCEnd("13820000009", "whatever")
	if !errors.Is(err, ErrPasswordNotSet) {
		t.Fatalf("expected ErrPasswordNotSet, got %v", err)
	}
}

func TestAuthService_SetCEndPassword_FirstTimeWithoutCurrentPassword(t *testing.T) {
	authSvc, db, _ := setupAuthUnitService(t)
	user := &model.User{
		Phone:  "13820000010",
		Name:   "首次设密用户",
		Status: "active",
	}
	if err := db.Omit("PasswordHash").Create(user).Error; err != nil {
		t.Fatalf("failed to create user without password: %v", err)
	}

	if err := authSvc.SetCEndPassword(user.ID, "", "NewPass@123"); err != nil {
		t.Fatalf("expected first-time password set to succeed, got %v", err)
	}

	updated, err := repository.NewUserRepository(db).GetByID(user.ID)
	if err != nil {
		t.Fatalf("failed to reload user: %v", err)
	}
	if !hasPasswordHash(updated.PasswordHash) {
		t.Fatal("expected password hash to be stored")
	}
	if err := VerifyPassword(updated.PasswordHash, "NewPass@123"); err != nil {
		t.Fatalf("stored password hash verification failed: %v", err)
	}
}

func TestAuthService_SetCEndPassword_RequiresCurrentPasswordWhenAlreadySet(t *testing.T) {
	authSvc, db, _ := setupAuthUnitService(t)
	user := seedAuthUnitUser(t, db, "13820000011", "Test@123", "active", 0)

	err := authSvc.SetCEndPassword(user.ID, "", "NewPass@123")
	if !errors.Is(err, ErrCurrentPasswordRequired) {
		t.Fatalf("expected ErrCurrentPasswordRequired, got %v", err)
	}

	err = authSvc.SetCEndPassword(user.ID, "Wrong@123", "NewPass@123")
	if !errors.Is(err, ErrCurrentPasswordInvalid) {
		t.Fatalf("expected ErrCurrentPasswordInvalid, got %v", err)
	}

	if err := authSvc.SetCEndPassword(user.ID, "Test@123", "NewPass@123"); err != nil {
		t.Fatalf("expected password update to succeed, got %v", err)
	}
}

func TestAuthService_ResetCEndPassword_Success(t *testing.T) {
	authSvc, db, _ := setupAuthUnitService(t)
	user := seedAuthUnitUser(t, db, "13820000012", "Test@123", "active", 0)
	seedAuthUnitProfile(t, db, user.ID, "测试地址")
	authSvc.smsService.SetTestCode("13820000012", "123456")

	if err := authSvc.ResetCEndPassword("13820000012", "123456", "ResetPass@123"); err != nil {
		t.Fatalf("expected reset password to succeed, got %v", err)
	}

	updated, err := repository.NewUserRepository(db).GetByID(user.ID)
	if err != nil {
		t.Fatalf("failed to reload user: %v", err)
	}
	if err := VerifyPassword(updated.PasswordHash, "ResetPass@123"); err != nil {
		t.Fatalf("updated password verification failed: %v", err)
	}
}

func TestAuthService_ResetCEndPassword_UserNotFound(t *testing.T) {
	authSvc, _, _ := setupAuthUnitService(t)
	authSvc.smsService.SetTestCode("13820000013", "123456")

	err := authSvc.ResetCEndPassword("13820000013", "123456", "ResetPass@123")
	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
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
	authSvc.SetGeocodeService(NewGeocodeService(""))
	createTestStation(t, db, 39.5, 116.0)
	expected := createTestStation(t, db, 39.908823, 116.397470)

	description := "需要送餐上门"
	result, err := authSvc.QuickStart(QuickStartInput{
		Phone:        "13820000004",
		Code:         "000000",
		Name:         "快速用户",
		Address:      "高纬地址",
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
	if result.Profile.CustomerType != consts.IdentityElderly {
		t.Fatalf("expected customer_type %q, got %q", consts.IdentityElderly, result.Profile.CustomerType)
	}

	var passwordHash sql.NullString
	if err := db.Raw("SELECT password_hash FROM users WHERE id = ?", result.User.ID).Scan(&passwordHash).Error; err != nil {
		t.Fatalf("failed to query password hash: %v", err)
	}
	if passwordHash.Valid {
		t.Fatalf("expected password_hash to be NULL for quick-start user, got %q", passwordHash.String)
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

func TestAuthService_QuickStart_UpdatesExistingUserAndCreatesManualReviewRequestWhenAddressCannotBeResolved(t *testing.T) {
	authSvc, db, _ := setupAuthUnitService(t)
	sourceStation := createTestStation(t, db, 30, 120)
	createTestStation(t, db, 31, 121)

	user := seedAuthUnitUser(t, db, "13820000005", "Test@123", "active", 0)
	seedAuthUnitProfile(t, db, user.ID, "旧地址")
	seedAuthUnitIdentity(t, db, user.ID, consts.IdentityElderly, true, 0)
	sourceStationID := sourceStation.ID

	result, err := authSvc.QuickStart(QuickStartInput{
		Phone:           "13820000005",
		Code:            "000000",
		Name:            "新姓名",
		Address:         "新地址",
		SourceStationID: &sourceStationID,
		ServiceType:     consts.ServiceTypeMedical,
		ContactName:     "新联系人",
		ContactPhone:    "13800138001",
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
	if result.Request.Status != consts.RequestStatusPending {
		t.Fatalf("expected pending request, got %s", result.Request.Status)
	}
	if result.Request.StationID != 0 || !result.Request.NeedsManualVerify || result.Request.DispatchBasis != DispatchBasisAddressManualReview {
		t.Fatalf("expected manual review request, got %+v", result.Request)
	}
	if result.Request.SourceStationID != sourceStationID {
		t.Fatalf("expected source station %d to be recorded, got %d", sourceStationID, result.Request.SourceStationID)
	}

	var taskCount int64
	if err := db.Model(&model.TaskAssignment{}).Where("request_id = ?", result.Request.ID).Count(&taskCount).Error; err != nil {
		t.Fatalf("failed to count tasks: %v", err)
	}
	if taskCount != 0 {
		t.Fatalf("expected no task to be created, got %d", taskCount)
	}
}

func TestAuthService_QuickStart_RequiresAddress(t *testing.T) {
	authSvc, _, _ := setupAuthUnitService(t)

	_, err := authSvc.QuickStart(QuickStartInput{
		Phone:       "13820000006",
		Code:        "000000",
		Name:        "无站点用户",
		ServiceType: consts.ServiceTypeMeal,
	})
	if !errors.Is(err, ErrAddressRequired) {
		t.Fatalf("expected ErrAddressRequired, got %v", err)
	}
}

func TestAuthService_QuickStart_RejectsInactiveExistingUser(t *testing.T) {
	authSvc, db, _ := setupAuthUnitService(t)
	user := seedAuthUnitUser(t, db, "13820000007", "Test@123", "inactive", 0)
	seedAuthUnitProfile(t, db, user.ID, "旧地址")
	seedAuthUnitIdentity(t, db, user.ID, consts.IdentityElderly, true, 0)

	_, err := authSvc.QuickStart(QuickStartInput{
		Phone:       "13820000007",
		Code:        "000000",
		Name:        "停用用户",
		Address:     "测试地址",
		ServiceType: consts.ServiceTypeMeal,
	})
	if !errors.Is(err, ErrUserInactive) {
		t.Fatalf("expected ErrUserInactive, got %v", err)
	}

	var requestCount int64
	if err := db.Model(&model.ServiceRequest{}).Where("user_id = ?", user.ID).Count(&requestCount).Error; err != nil {
		t.Fatalf("failed to count requests: %v", err)
	}
	if requestCount != 0 {
		t.Fatalf("expected no request to be created, got %d", requestCount)
	}
}

func TestAuthService_Register_Success(t *testing.T) {
	authSvc, db, _ := setupAuthUnitService(t)

	// 设置测试验证码
	authSvc.smsService.SetTestCode("13900000001", "123456")

	// 注册新用户
	result, err := authSvc.Register(RegisterInput{
		Phone:    "13900000001",
		Code:     "123456",
		Password: "Test@123",
		Name:     "测试用户",
	})
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}
	if result.User == nil {
		t.Fatal("expected user to be returned")
	}
	if result.User.Phone != "13900000001" {
		t.Fatalf("expected phone 13900000001, got %s", result.User.Phone)
	}
	if result.User.Name != "测试用户" {
		t.Fatalf("expected name 测试用户, got %s", result.User.Name)
	}
	if result.Token == "" {
		t.Fatal("expected token to be returned")
	}
	if result.RefreshToken == "" {
		t.Fatal("expected refresh token to be returned")
	}

	// 验证用户和身份已创建
	var user model.User
	if err := db.Where("phone = ?", "13900000001").First(&user).Error; err != nil {
		t.Fatalf("user not found in db: %v", err)
	}

	var identity model.UserIdentity
	if err := db.Where("user_id = ? AND identity_type = ?", user.ID, consts.IdentityElderly).First(&identity).Error; err != nil {
		t.Fatalf("elderly identity not found: %v", err)
	}

	var profile model.CustomerProfile
	if err := db.Where("user_id = ?", user.ID).First(&profile).Error; err != nil {
		t.Fatalf("customer profile not found: %v", err)
	}
}

func TestAuthService_Register_DuplicatePhone(t *testing.T) {
	authSvc, db, _ := setupAuthUnitService(t)

	// 创建已存在的用户
	existing := seedAuthUnitUser(t, db, "13900000002", "Test@123", "active", 0)

	// 设置测试验证码
	authSvc.smsService.SetTestCode("13900000002", "123456")

	// 尝试注册相同手机号
	_, err := authSvc.Register(RegisterInput{
		Phone:    "13900000002",
		Code:     "123456",
		Password: "Test@123",
		Name:     "重复用户",
	})
	if err == nil {
		t.Fatal("expected error for duplicate phone")
	}
	if err.Error() != "该手机号已注册，请直接登录" {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = existing // avoid unused variable
}

func TestAuthService_Register_InvalidCode(t *testing.T) {
	authSvc, _, _ := setupAuthUnitService(t)

	// 设置了正确验证码，但注册时使用错误验证码
	authSvc.smsService.SetTestCode("13900000003", "123456")

	_, err := authSvc.Register(RegisterInput{
		Phone:    "13900000003",
		Code:     "999999", // 错误验证码
		Password: "Test@123",
		Name:     "无效验证码",
	})
	if err == nil {
		t.Fatal("expected error for invalid code")
	}
	if err != ErrCodeInvalid {
		t.Fatalf("expected ErrCodeInvalid, got %v", err)
	}
}
