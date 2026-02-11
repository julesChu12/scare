//go:build integration

package repository

import (
	"testing"

	"community-elderly-care-platform/internal/dao/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// NOTE: 这些测试需要兼容的数据库环境
// 运行方式: go test -tags=integration ./...

func setupTestDB(t *testing.T) *gorm.DB {
	tmpFile := t.TempDir() + "/test.db"
	dsn := tmpFile + "?_loc=Local&_parseTime=true"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	err = db.AutoMigrate(&model.User{}, &model.CustomerProfile{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	return db
}

func TestCustomerRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCustomerRepository(db)

	customerType := "elderly"
	gender := "male"
	address := "北京市朝阳区幸福小区1号楼101"
	healthStatus := "良好"
	disabilityLevel := "自理"

	profile := &model.CustomerProfile{
		UserID:          1,
		CustomerType:    customerType,
		Gender:          gender,
		Address:         address,
		HealthStatus:    healthStatus,
		DisabilityLevel: disabilityLevel,
	}

	err := repo.Create(profile)
	if err != nil {
		t.Fatalf("failed to create customer profile: %v", err)
	}

	if profile.ID == 0 {
		t.Error("expected profile ID to be set after creation")
	}
}

func TestCustomerRepository_GetByUserID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCustomerRepository(db)

	// Create a test profile
	profile := &model.CustomerProfile{
		UserID:       1,
		CustomerType: "elderly",
	}
	repo.Create(profile)

	// Retrieve the profile
	result, err := repo.GetByUserID(1)
	if err != nil {
		t.Fatalf("failed to get customer profile: %v", err)
	}

	if result.UserID != 1 {
		t.Errorf("expected user_id 1, got %d", result.UserID)
	}

	if result.CustomerType != "elderly" {
		t.Errorf("expected customer_type 'elderly', got %v", result.CustomerType)
	}
}

func TestCustomerRepository_GetByUserID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCustomerRepository(db)

	_, err := repo.GetByUserID(999)
	if err == nil {
		t.Error("expected error for non-existent profile, got nil")
	}
}

func TestCustomerRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCustomerRepository(db)

	// Create initial profile
	profile := &model.CustomerProfile{
		UserID:       1,
		CustomerType: "elderly",
	}
	repo.Create(profile)

	// Update customer type
	profile.CustomerType = "disabled"

	err := repo.Update(profile)
	if err != nil {
		t.Fatalf("failed to update customer profile: %v", err)
	}

	// Verify update
	updated, _ := repo.GetByUserID(1)
	if updated.CustomerType != "disabled" {
		t.Errorf("expected customer_type 'disabled', got %v", updated.CustomerType)
	}
}

func TestCustomerRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCustomerRepository(db)

	// Create a test profile
	profile := &model.CustomerProfile{
		UserID:       1,
		CustomerType: "elderly",
	}
	repo.Create(profile)

	// Delete the profile
	err := repo.Delete(1)
	if err != nil {
		t.Fatalf("failed to delete customer profile: %v", err)
	}

	// Verify deletion
	_, err = repo.GetByUserID(1)
	if err == nil {
		t.Error("expected error for deleted profile, got nil")
	}
}

func TestCustomerRepository_Exists(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCustomerRepository(db)

	// Check non-existent profile
	exists, err := repo.Exists(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exists {
		t.Error("expected profile to not exist")
	}

	// Create a profile
	profile := &model.CustomerProfile{
		UserID:       1,
		CustomerType: "elderly",
	}
	repo.Create(profile)

	// Check existing profile
	exists, err = repo.Exists(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !exists {
		t.Error("expected profile to exist")
	}
}

func TestCustomerRepository_ListByType(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCustomerRepository(db)

	// Create multiple profiles of different types
	profiles := []*model.CustomerProfile{
		{UserID: 1, CustomerType: "elderly"},
		{UserID: 2, CustomerType: "elderly"},
		{UserID: 3, CustomerType: "disabled"},
		{UserID: 4, CustomerType: "pregnant"},
	}

	for _, p := range profiles {
		repo.Create(p)
	}

	// List elderly profiles
	elderlyProfiles, err := repo.ListByType("elderly")
	if err != nil {
		t.Fatalf("failed to list profiles: %v", err)
	}

	if len(elderlyProfiles) != 2 {
		t.Errorf("expected 2 elderly profiles, got %d", len(elderlyProfiles))
	}

	// List disabled profiles
	disabledProfiles, err := repo.ListByType("disabled")
	if err != nil {
		t.Fatalf("failed to list profiles: %v", err)
	}

	if len(disabledProfiles) != 1 {
		t.Errorf("expected 1 disabled profile, got %d", len(disabledProfiles))
	}
}
