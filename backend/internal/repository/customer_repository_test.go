//go:build integration

package repository

import (
	"testing"
	"time"

	"community-elderly-care-platform/internal/dao/model"
)

func TestCustomerRepository_CreateAndGetByUserID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCustomerRepository(db)

	now := time.Now()
	profile := &model.CustomerProfile{
		UserID:           1,
		Address:          "北京市朝阳区幸福小区1号楼101",
		CustomerType:     "elderly",
		EmergencyContact: "13800000001",
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := repo.Create(profile); err != nil {
		t.Fatalf("failed to create customer profile: %v", err)
	}
	if profile.ID == 0 {
		t.Fatal("expected profile id to be assigned")
	}

	got, err := repo.GetByUserID(1)
	if err != nil {
		t.Fatalf("failed to get customer profile: %v", err)
	}
	if got.UserID != 1 {
		t.Fatalf("expected user_id 1, got %d", got.UserID)
	}
	if got.Address != profile.Address {
		t.Fatalf("expected address %q, got %q", profile.Address, got.Address)
	}
	if got.CustomerType != profile.CustomerType {
		t.Fatalf("expected customer_type %q, got %q", profile.CustomerType, got.CustomerType)
	}
}

func TestCustomerRepository_GetByUserID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCustomerRepository(db)

	if _, err := repo.GetByUserID(999); err == nil {
		t.Fatal("expected error for missing profile")
	}
}

func TestCustomerRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCustomerRepository(db)

	now := time.Now()
	profile := &model.CustomerProfile{
		UserID:           2,
		Address:          "原地址",
		EmergencyContact: "13800000002",
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	if err := repo.Create(profile); err != nil {
		t.Fatalf("failed to create customer profile: %v", err)
	}

	profile.Address = "新地址"
	profile.UpdatedAt = now.Add(time.Minute)
	if err := repo.Update(profile); err != nil {
		t.Fatalf("failed to update customer profile: %v", err)
	}

	got, err := repo.GetByUserID(2)
	if err != nil {
		t.Fatalf("failed to get updated profile: %v", err)
	}
	if got.Address != "新地址" {
		t.Fatalf("expected updated address, got %q", got.Address)
	}
}

func TestCustomerRepository_Exists(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCustomerRepository(db)

	exists, err := repo.Exists(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exists {
		t.Fatal("expected profile to not exist")
	}

	now := time.Now()
	if err := repo.Create(&model.CustomerProfile{
		UserID:           3,
		Address:          "存在地址",
		EmergencyContact: "13800000003",
		CreatedAt:        now,
		UpdatedAt:        now,
	}); err != nil {
		t.Fatalf("failed to create customer profile: %v", err)
	}

	exists, err = repo.Exists(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !exists {
		t.Fatal("expected profile to exist")
	}
}
