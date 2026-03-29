package service

import "testing"

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
