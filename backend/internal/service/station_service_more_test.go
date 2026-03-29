package service

import (
	"errors"
	"testing"

	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/repository"
)

func TestStationService_CreateListUpdateDelete(t *testing.T) {
	_, stationRepo := setupStationSelectionTestDB(t)
	svc := NewStationService(stationRepo)

	if err := svc.Create(&model.ServiceStation{}); !errors.Is(err, ErrInvalidStation) {
		t.Fatalf("expected ErrInvalidStation, got %v", err)
	}

	station := &model.ServiceStation{
		Name:    "测试站点",
		Code:    "ST-CRUD",
		Address: "测试地址",
		Status:  "active",
	}
	if err := svc.Create(station); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	list, total, err := svc.List(1, 10, StationListFilter{Keyword: "测试"})
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if total != 1 || len(list) != 1 || list[0].ID != station.ID {
		t.Fatalf("unexpected list result: total=%d len=%d", total, len(list))
	}

	if _, err := svc.GetByID(station.ID); err != nil {
		t.Fatalf("GetByID returned error: %v", err)
	}

	if err := svc.Update(&model.ServiceStation{ID: station.ID}); !errors.Is(err, ErrInvalidStation) {
		t.Fatalf("expected ErrInvalidStation for missing name, got %v", err)
	}

	station.Name = "更新后的站点"
	if err := svc.Update(station); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	updated, err := svc.GetByID(station.ID)
	if err != nil {
		t.Fatalf("GetByID after update returned error: %v", err)
	}
	if updated.Name != "更新后的站点" {
		t.Fatalf("expected updated name, got %q", updated.Name)
	}

	if err := svc.Delete(0); !errors.Is(err, ErrInvalidStation) {
		t.Fatalf("expected ErrInvalidStation for zero ID, got %v", err)
	}
	if err := svc.Delete(station.ID); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
	if _, err := svc.GetByID(station.ID); err == nil {
		t.Fatal("expected deleted station to be unavailable")
	}
}

func TestStationService_MatchStation_UsesUserProfileGeocode(t *testing.T) {
	db, stationRepo := setupStationSelectionTestDB(t)
	createCustomerProfilesTable(t, db)
	expected := createTestStation(t, db, 39.908823, 116.397470)

	if err := db.Create(&model.CustomerProfile{
		UserID:           2001,
		Address:          "北京市东城区",
		EmergencyContact: `{}`,
	}).Error; err != nil {
		t.Fatalf("failed to create customer profile: %v", err)
	}

	svc := NewStationService(stationRepo)
	svc.SetCustomerRepo(repository.NewCustomerRepository(db))
	svc.SetGeocodeService(NewGeocodeService(""))

	station, err := svc.MatchStation(MatchStationInput{UserID: 2001})
	if err != nil {
		t.Fatalf("MatchStation returned error: %v", err)
	}
	if station.ID != expected.ID {
		t.Fatalf("expected station %d, got %d", expected.ID, station.ID)
	}
}

func TestStationService_MatchStation_DefaultStationFallbacks(t *testing.T) {
	db, stationRepo := setupStationSelectionTestDB(t)
	defaultStation := createTestStation(t, db, 30, 120)
	defaultStation.Name = DefaultStationName
	if err := db.Save(defaultStation).Error; err != nil {
		t.Fatalf("failed to rename default station: %v", err)
	}

	svc := NewStationService(stationRepo)
	station, err := svc.MatchStation(MatchStationInput{})
	if err != nil {
		t.Fatalf("MatchStation returned error: %v", err)
	}
	if station.ID != defaultStation.ID {
		t.Fatalf("expected default station %d, got %d", defaultStation.ID, station.ID)
	}

	db2, stationRepo2 := setupStationSelectionTestDB(t)
	fallback := createTestStation(t, db2, 31, 121)
	svc2 := NewStationService(stationRepo2)

	station, err = svc2.MatchStation(MatchStationInput{})
	if err != nil {
		t.Fatalf("MatchStation fallback returned error: %v", err)
	}
	if station.ID != fallback.ID {
		t.Fatalf("expected first active station %d, got %d", fallback.ID, station.ID)
	}
}
