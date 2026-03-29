package service

import (
	"errors"
	"testing"
	"time"

	"community-elderly-care-platform/internal/consts"
	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/repository"

	"gorm.io/gorm"
)

func setupUserServiceWithDBForTest(t *testing.T) (*UserService, *gorm.DB) {
	t.Helper()

	db := openServiceTestDB(t, "user_service_more_test.db")
	createUsersTable(t, db)
	createUserIdentitiesTable(t, db)

	userRepo := repository.NewUserRepository(db)
	userIdentityRepo := repository.NewUserIdentityRepository(db)
	return NewUserService(userRepo, userIdentityRepo), db
}

func seedUserServiceUser(t *testing.T, db *gorm.DB, phone, password, name, identityType string, stationID int64) *model.User {
	t.Helper()

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	user := &model.User{
		Phone:        phone,
		PasswordHash: hash,
		Name:         name,
		Status:       "active",
		StationID:    stationID,
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	identity := &model.UserIdentity{
		UserID:       user.ID,
		IdentityType: identityType,
		IsPrimary:    true,
		StationID:    stationID,
		Status:       "active",
		GrantedAt:    time.Now(),
	}
	if err := db.Create(identity).Error; err != nil {
		t.Fatalf("failed to create identity: %v", err)
	}

	return user
}

func TestUserService_Update_Success(t *testing.T) {
	svc, db := setupUserServiceWithDBForTest(t)
	user := seedUserServiceUser(t, db, "13810000020", "Test@123", "原姓名", consts.IdentityStaff, 1)

	birthDate := time.Date(1950, 1, 2, 0, 0, 0, 0, time.Local)
	updated, err := svc.Update(UserInput{
		ID:           user.ID,
		Name:         "新姓名",
		Email:        "new@example.com",
		Gender:       "male",
		BirthDate:    birthDate,
		IDCard:       "330101195001020011",
		IDCardHMAC:   "hmac",
		IDCardMasked: "330101********0011",
		StationID:    2,
		Status:       "inactive",
		Password:     "NewPass@123",
	})
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	if updated.Name != "新姓名" || updated.Email != "new@example.com" || updated.StationID != 2 {
		t.Fatalf("unexpected updated user payload: %+v", updated.User)
	}
	if updated.Status != "inactive" || updated.IDCardHmac != "hmac" || updated.IDCardMasked != "330101********0011" {
		t.Fatalf("unexpected updated user sensitive fields: %+v", updated.User)
	}
	if !updated.BirthDate.Equal(birthDate) {
		t.Fatalf("expected birth date %v, got %v", birthDate, updated.BirthDate)
	}
	if err := VerifyPassword(updated.PasswordHash, "NewPass@123"); err != nil {
		t.Fatalf("password hash was not updated: %v", err)
	}
}

func TestUserService_ListWithFilterAndGetByID(t *testing.T) {
	svc, db := setupUserServiceWithDBForTest(t)
	staffUser := seedUserServiceUser(t, db, "13810000021", "Test@123", "张三", consts.IdentityStaff, 1)
	seedUserServiceUser(t, db, "13810000022", "Test@123", "李四", consts.IdentityFamily, 0)

	if _, err := svc.GetByID(0); !errors.Is(err, ErrInvalidUser) {
		t.Fatalf("expected ErrInvalidUser for zero ID, got %v", err)
	}

	gotUser, err := svc.GetByID(staffUser.ID)
	if err != nil {
		t.Fatalf("GetByID returned error: %v", err)
	}
	if gotUser.PrimaryIdentity == nil || gotUser.PrimaryIdentity.IdentityType != consts.IdentityStaff {
		t.Fatalf("unexpected primary identity: %+v", gotUser.PrimaryIdentity)
	}

	list, total, err := svc.ListWithFilter(1, 10, UserFilter{
		Role:      consts.IdentityStaff,
		Status:    "active",
		StationID: 1,
		Keyword:   "张",
	})
	if err != nil {
		t.Fatalf("ListWithFilter returned error: %v", err)
	}
	if total != 1 || len(list) != 1 || list[0].User.ID != staffUser.ID {
		t.Fatalf("unexpected list result: total=%d len=%d", total, len(list))
	}
}
