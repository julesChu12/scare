package service

import (
	"testing"

	"community-elderly-care-platform/pkg/geo"
)

func TestStationService_MatchStation_FallsBackToHaversineNearest(t *testing.T) {
	db, stationRepo := setupStationSelectionTestDB(t)
	createTestStation(t, db, 75, 0)
	expected := createTestStation(t, db, 80, 20)

	svc := NewStationService(stationRepo)
	lat := 80.0
	lng := 0.0

	station, err := svc.MatchStation(MatchStationInput{
		Latitude:  &lat,
		Longitude: &lng,
	})
	if err != nil {
		t.Fatalf("MatchStation returned error: %v", err)
	}

	if station.ID != expected.ID {
		t.Fatalf("expected nearest station %d, got %d", expected.ID, station.ID)
	}
}

func TestStationService_MatchStation_SkipsInactiveGeofenceStation(t *testing.T) {
	db, stationRepo := setupStationSelectionTestDB(t)

	inactiveStation := createTestStation(t, db, 39.953148, 116.374073)
	inactiveStation.Status = "inactive"
	if err := db.Save(inactiveStation).Error; err != nil {
		t.Fatalf("failed to update inactive station: %v", err)
	}

	activeStation := createTestStation(t, db, 39.953600, 116.372359)

	svc := NewStationService(stationRepo)
	svc.SetGeofenceService(&GeofenceService{
		engine: geo.NewEngine([]geo.Zone{
			{
				ID:        1,
				StationID: inactiveStation.ID,
				Priority:  100,
				Points: []geo.Point{
					{Lat: 39.952700, Lng: 116.373600},
					{Lat: 39.952700, Lng: 116.374500},
					{Lat: 39.953500, Lng: 116.374500},
					{Lat: 39.953500, Lng: 116.373600},
				},
			},
		}),
	})

	lat := 39.953148
	lng := 116.374073
	station, err := svc.MatchStation(MatchStationInput{
		Latitude:  &lat,
		Longitude: &lng,
	})
	if err != nil {
		t.Fatalf("MatchStation returned error: %v", err)
	}

	if station.ID != activeStation.ID {
		t.Fatalf("expected nearest active station %d, got %d", activeStation.ID, station.ID)
	}
}
