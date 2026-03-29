package service

import (
	"errors"
	"testing"

	"community-elderly-care-platform/internal/repository"
	"community-elderly-care-platform/pkg/geo"
)

func setupZoneServiceForTest(t *testing.T) (*ZoneService, *GeofenceService) {
	t.Helper()

	db := openServiceTestDB(t, "zone_service_test.db")
	createServiceZonesTable(t, db)
	zoneRepo := repository.NewZoneRepository(db)
	geofenceSvc := NewGeofenceService(zoneRepo)

	return NewZoneService(zoneRepo, geofenceSvc), geofenceSvc
}

func TestZoneService_CreateUpdateDeleteAndList(t *testing.T) {
	svc, geofenceSvc := setupZoneServiceForTest(t)
	points := []geo.Point{
		{Lat: 30.0, Lng: 120.0},
		{Lat: 30.1, Lng: 120.0},
		{Lat: 30.1, Lng: 120.1},
		{Lat: 30.0, Lng: 120.1},
	}

	if _, err := svc.Create(ZoneInput{StationID: 1, Name: "无效围栏", Points: points[:2]}); !errors.Is(err, ErrInvalidZone) {
		t.Fatalf("expected ErrInvalidZone, got %v", err)
	}
	if _, err := svc.GetByID(0); !errors.Is(err, ErrInvalidZone) {
		t.Fatalf("expected ErrInvalidZone for zero ID, got %v", err)
	}
	if err := svc.Delete(0); !errors.Is(err, ErrInvalidZone) {
		t.Fatalf("expected ErrInvalidZone for zero ID delete, got %v", err)
	}

	zone, err := svc.Create(ZoneInput{
		StationID: 1,
		Name:      "测试围栏",
		Points:    points,
		Priority:  10,
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if zone.Status != "active" {
		t.Fatalf("expected default status active, got %s", zone.Status)
	}

	if stationID, matched := geofenceSvc.Match(30.05, 120.05); !matched || stationID != 1 {
		t.Fatalf("expected created zone to be reloaded into geofence engine, got stationID=%d matched=%v", stationID, matched)
	}

	list, total, err := svc.List(1, 10, ZoneListFilter{StationID: 1})
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if total != 1 || len(list) != 1 || list[0].ID != zone.ID {
		t.Fatalf("unexpected list result: total=%d len=%d", total, len(list))
	}

	got, err := svc.GetByID(zone.ID)
	if err != nil {
		t.Fatalf("GetByID returned error: %v", err)
	}
	if got.Name != "测试围栏" {
		t.Fatalf("expected zone name 测试围栏, got %s", got.Name)
	}

	if _, err := svc.Update(ZoneInput{ID: zone.ID, StationID: 1, Name: "无效", Points: points[:2]}); !errors.Is(err, ErrInvalidZone) {
		t.Fatalf("expected ErrInvalidZone for short polygon, got %v", err)
	}

	updated, err := svc.Update(ZoneInput{
		ID:        zone.ID,
		StationID: 2,
		Name:      "已停用围栏",
		Points:    points,
		Priority:  99,
		Status:    "inactive",
	})
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if updated.StationID != 2 || updated.Status != "inactive" || updated.Priority != 99 {
		t.Fatalf("unexpected updated zone: %+v", updated)
	}

	if _, matched := geofenceSvc.Match(30.05, 120.05); matched {
		t.Fatal("expected inactive zone to be removed from geofence engine")
	}

	if err := svc.Delete(zone.ID); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
	if _, err := svc.GetByID(zone.ID); err == nil {
		t.Fatal("expected deleted zone to be unavailable")
	}
}
