package handler

import (
	"net/http"
	"testing"

	"community-elderly-care-platform/internal/repository"
	"community-elderly-care-platform/internal/service"
)

func TestCProfileHandler_UpdateProfile_UpdatesUserNameAndProfileFields(t *testing.T) {
	db := openHandlerTestDB(t, "c_profile_update.db")
	createHandlerTables(t, db)
	seedHandlerUserAndProfile(t, db, 20, "13900000020", "旧姓名", "elderly", "旧地址")

	profileService := service.NewCEndProfileService(
		db,
		repository.NewUserRepository(db),
		repository.NewCustomerRepository(db),
	)
	handler := NewCProfileHandler(profileService, service.NewGeocodeService(""))

	c, w := newJSONTestContext(t, http.MethodPut, "/c/profile", map[string]any{
		"name":      "新姓名",
		"id_number": "110101199001011234",
		"address":   "新地址",
		"user_type": "family",
	})
	setCEndClaims(c, 20)

	handler.UpdateProfile(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	user, err := repository.NewUserRepository(db).GetByID(20)
	if err != nil {
		t.Fatalf("failed to reload user: %v", err)
	}
	if user.Name != "新姓名" {
		t.Fatalf("expected user name 新姓名, got %q", user.Name)
	}

	profile, err := repository.NewCustomerRepository(db).GetByUserID(20)
	if err != nil {
		t.Fatalf("failed to reload profile: %v", err)
	}
	if profile.IDCard != "110101199001011234" {
		t.Fatalf("expected id_card updated, got %q", profile.IDCard)
	}
	if profile.Address != "新地址" {
		t.Fatalf("expected address 新地址, got %q", profile.Address)
	}
	if profile.CustomerType != "family" {
		t.Fatalf("expected customer_type family, got %q", profile.CustomerType)
	}
}
