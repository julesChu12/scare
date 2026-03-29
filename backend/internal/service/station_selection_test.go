package service

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/repository"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var stationSelectionTestSeq int64

func setupStationSelectionTestDB(t *testing.T) (*gorm.DB, *repository.StationRepository) {
	t.Helper()

	tmpFile := t.TempDir() + "/station_selection_test.db"
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

	createServiceStationsTable(t, db)
	return db, repository.NewStationRepository(db)
}

func createServiceStationsTable(t *testing.T, db *gorm.DB) {
	t.Helper()

	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS service_stations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			code TEXT,
			address TEXT,
			phone TEXT,
			latitude REAL,
			longitude REAL,
			service_area TEXT,
			capacity INTEGER,
			work_hours TEXT,
			status TEXT,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		);
	`).Error; err != nil {
		t.Fatalf("failed to create service_stations table: %v", err)
	}
}

func createTestStation(t *testing.T, db *gorm.DB, lat, lng float64) *model.ServiceStation {
	t.Helper()

	seq := atomic.AddInt64(&stationSelectionTestSeq, 1)
	station := &model.ServiceStation{
		Name:      fmt.Sprintf("Station-%d", seq),
		Code:      fmt.Sprintf("ST-%d", seq),
		Latitude:  lat,
		Longitude: lng,
		Status:    "active",
	}

	if err := db.Create(station).Error; err != nil {
		t.Fatalf("failed to create station: %v", err)
	}

	return station
}

func TestNearestStationByHaversine_PrefersTrueGeodesicDistance(t *testing.T) {
	nearest, err := nearestStationByHaversine([]*model.ServiceStation{
		{ID: 1, Latitude: 75, Longitude: 0},
		{ID: 2, Latitude: 80, Longitude: 20},
	}, 80, 0)
	if err != nil {
		t.Fatalf("nearestStationByHaversine returned error: %v", err)
	}

	if nearest.ID != 2 {
		t.Fatalf("expected station 2 to be nearest by Haversine, got %d", nearest.ID)
	}
}
