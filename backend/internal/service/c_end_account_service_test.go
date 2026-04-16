package service

import (
	"errors"
	"testing"

	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/repository"
)

func TestCEndAccountService_GetAccountInfo_Success(t *testing.T) {
	db := openServiceTestDB(t, "c_end_account_info.db")
	createUsersTable(t, db)
	createCustomerProfilesTable(t, db)

	user := &model.User{
		Phone:        "13910000001",
		PasswordHash: "hashed-password",
		Name:         "张大爷",
		Status:       "active",
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	if err := db.Create(&model.CustomerProfile{
		UserID:           user.ID,
		CustomerType:     "elderly",
		Address:          "测试地址",
		EmergencyContact: `{}`,
	}).Error; err != nil {
		t.Fatalf("failed to create profile: %v", err)
	}

	svc := NewCEndAccountService(
		repository.NewUserRepository(db),
		repository.NewCustomerRepository(db),
	)

	info, err := svc.GetAccountInfo(user.ID)
	if err != nil {
		t.Fatalf("GetAccountInfo failed: %v", err)
	}
	if info.UserID != user.ID || info.Type != "c_end" {
		t.Fatalf("unexpected account info: %+v", info)
	}
	if info.CustomerType != "elderly" || info.Name != "张大爷" || info.Phone != "13910000001" {
		t.Fatalf("unexpected account info fields: %+v", info)
	}
	if !info.HasPassword {
		t.Fatalf("expected has_password=true")
	}
}

func TestCEndAccountService_GetCheckPayload_AllowsMissingProfile(t *testing.T) {
	db := openServiceTestDB(t, "c_end_check_payload.db")
	createUsersTable(t, db)
	createCustomerProfilesTable(t, db)

	user := &model.User{
		Phone:  "13910000002",
		Name:   "李阿姨",
		Status: "active",
	}
	if err := db.Omit("PasswordHash").Create(user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	svc := NewCEndAccountService(
		repository.NewUserRepository(db),
		repository.NewCustomerRepository(db),
	)

	payload, err := svc.GetCheckPayload(user.ID)
	if err != nil {
		t.Fatalf("GetCheckPayload failed: %v", err)
	}
	if payload.User.Phone != "13910000002" || payload.User.Role != "c_end" {
		t.Fatalf("unexpected payload user: %+v", payload.User)
	}
	if payload.User.HasPassword {
		t.Fatalf("expected has_password=false for empty password_hash")
	}
	if payload.Profile != nil {
		t.Fatalf("expected profile=nil when customer profile missing")
	}
}

func TestCEndProfileService_Update_UpdatesUserAndProfile(t *testing.T) {
	db := openServiceTestDB(t, "c_end_profile_update.db")
	createUsersTable(t, db)
	createCustomerProfilesTable(t, db)

	user := &model.User{
		Phone:        "13910000003",
		PasswordHash: "hashed-password",
		Name:         "旧姓名",
		Status:       "active",
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	if err := db.Create(&model.CustomerProfile{
		UserID:           user.ID,
		CustomerType:     "elderly",
		Address:          "旧地址",
		EmergencyContact: `{}`,
	}).Error; err != nil {
		t.Fatalf("failed to create profile: %v", err)
	}

	svc := NewCEndProfileService(
		db,
		repository.NewUserRepository(db),
		repository.NewCustomerRepository(db),
	)

	name := "新姓名"
	idNumber := "110101199001011234"
	address := "新地址"
	userType := "family"
	profile, err := svc.Update(user.ID, CEndProfileUpdateInput{
		Name:     &name,
		IDNumber: &idNumber,
		Address:  &address,
		UserType: &userType,
	})
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if profile.IDCard != idNumber || profile.Address != address || profile.CustomerType != userType {
		t.Fatalf("unexpected updated profile: %+v", profile)
	}

	reloadedUser, err := repository.NewUserRepository(db).GetByID(user.ID)
	if err != nil {
		t.Fatalf("failed to reload user: %v", err)
	}
	if reloadedUser.Name != name {
		t.Fatalf("expected user name %q, got %q", name, reloadedUser.Name)
	}
}

func TestCEndProfileService_Update_ProfileNotFound(t *testing.T) {
	db := openServiceTestDB(t, "c_end_profile_missing.db")
	createUsersTable(t, db)
	createCustomerProfilesTable(t, db)

	user := &model.User{
		Phone:  "13910000004",
		Name:   "无档案用户",
		Status: "active",
	}
	if err := db.Omit("PasswordHash").Create(user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	svc := NewCEndProfileService(
		db,
		repository.NewUserRepository(db),
		repository.NewCustomerRepository(db),
	)

	address := "任意地址"
	_, err := svc.Update(user.ID, CEndProfileUpdateInput{Address: &address})
	if !errors.Is(err, ErrNoCustomerProfile) {
		t.Fatalf("expected ErrNoCustomerProfile, got %v", err)
	}
}
