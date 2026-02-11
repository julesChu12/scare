package service

import (
	"errors"
	"testing"
	"time"

	"community-elderly-care-platform/internal/consts"
	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/repository"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupUserServiceTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	tmpFile := t.TempDir() + "/user_service_test.db"
	dsn := tmpFile + "?_loc=auto&parseTime=true"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		NowFunc: func() time.Time {
			return time.Now().Local().Truncate(time.Second)
		},
	})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	if err := db.AutoMigrate(&model.User{}, &model.UserIdentity{}); err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	return db
}

func setupUserServiceForTest(t *testing.T) *UserService {
	t.Helper()

	db := setupUserServiceTestDB(t)
	userRepo := repository.NewUserRepository(db)
	userIdentityRepo := repository.NewUserIdentityRepository(db)
	return NewUserService(userRepo, userIdentityRepo)
}

func TestUserService_Create_InvalidInput(t *testing.T) {
	svc := setupUserServiceForTest(t)

	testCases := []UserInput{
		{Password: "Test@123", IdentityType: consts.IdentityStaff},
		{Phone: "13810000001", IdentityType: consts.IdentityStaff},
		{Phone: "13810000002", Password: "Test@123"},
		{Phone: "13810000003", Password: "Test@123", IdentityType: "unknown"},
	}

	for _, tc := range testCases {
		_, err := svc.Create(tc)
		if !errors.Is(err, ErrInvalidUser) {
			t.Fatalf("expected ErrInvalidUser, got %v", err)
		}
	}
}

func TestUserService_Create_DefaultStatusAndIdentity(t *testing.T) {
	svc := setupUserServiceForTest(t)

	user, err := svc.Create(UserInput{
		Phone:        "13810000010",
		Password:     "Test@123",
		Name:         "Service User",
		IdentityType: consts.IdentityStaff,
		StationID:    1,
	})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	if user.Status != "active" {
		t.Fatalf("expected default status active, got %s", user.Status)
	}
	if user.PasswordHash == "" {
		t.Fatal("expected password hash to be set")
	}
	if err := VerifyPassword(user.PasswordHash, "Test@123"); err != nil {
		t.Fatalf("password hash verification failed: %v", err)
	}
	if len(user.BEndIdentities) != 1 || user.BEndIdentities[0] != consts.IdentityStaff {
		t.Fatalf("expected B-end identity %s, got %+v", consts.IdentityStaff, user.BEndIdentities)
	}
	if user.PrimaryIdentity == nil || user.PrimaryIdentity.IdentityType != consts.IdentityStaff {
		t.Fatalf("expected primary identity %s, got %+v", consts.IdentityStaff, user.PrimaryIdentity)
	}
	if user.PrimaryIdentity.StationID != 1 {
		t.Fatalf("expected identity station_id=1, got %d", user.PrimaryIdentity.StationID)
	}
}

func TestUserService_Update_ValidateID(t *testing.T) {
	svc := setupUserServiceForTest(t)

	if _, err := svc.Update(UserInput{ID: 0, Name: "invalid"}); !errors.Is(err, ErrInvalidUser) {
		t.Fatalf("expected ErrInvalidUser for zero ID, got %v", err)
	}
}
