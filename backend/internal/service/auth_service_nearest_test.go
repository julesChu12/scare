package service

import "testing"

func TestAuthService_FindNearestStation_UsesHaversineDistance(t *testing.T) {
	db, stationRepo := setupStationSelectionTestDB(t)
	createTestStation(t, db, 75, 0)
	expected := createTestStation(t, db, 80, 20)

	svc := &AuthService{stationRepo: stationRepo}

	nearestID, err := svc.findNearestStation(80, 0)
	if err != nil {
		t.Fatalf("findNearestStation returned error: %v", err)
	}

	if nearestID != expected.ID {
		t.Fatalf("expected nearest station %d, got %d", expected.ID, nearestID)
	}
}
