package handler

import (
	"net/http"
	"testing"

	"community-elderly-care-platform/internal/repository"
	"community-elderly-care-platform/internal/service"
	"community-elderly-care-platform/pkg/geo"

	"github.com/gin-gonic/gin"
)

func TestStationHandler_List_StationManagerScopedToOwnStation(t *testing.T) {
	db := openHandlerTestDB(t, "station_handler_scope.db")
	createHandlerTables(t, db)
	stationA := seedHandlerStation(t, db, "站点A")
	seedHandlerStation(t, db, "站点B")

	handler := NewStationHandler(service.NewStationService(repository.NewStationRepository(db)))
	c, w := newJSONTestContext(t, http.MethodGet, "/b/stations?page=1&page_size=10", nil)
	setBEndClaims(c, 2, stationA.ID, "station_manager")

	handler.List(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	resp := decodeResponseMap(t, w)
	data := resp["data"].(map[string]any)
	if int(data["total"].(float64)) != 1 {
		t.Fatalf("expected total 1, got %v", data["total"])
	}
	items := data["items"].([]any)
	item := items[0].(map[string]any)
	if int64(item["id"].(float64)) != stationA.ID {
		t.Fatalf("expected station %d, got %v", stationA.ID, item["id"])
	}
}

func TestStationHandler_Get_StationManagerForbiddenForOtherStation(t *testing.T) {
	db := openHandlerTestDB(t, "station_handler_get_scope.db")
	createHandlerTables(t, db)
	stationA := seedHandlerStation(t, db, "站点A")
	stationB := seedHandlerStation(t, db, "站点B")

	handler := NewStationHandler(service.NewStationService(repository.NewStationRepository(db)))
	c, w := newJSONTestContext(t, http.MethodGet, "/b/stations/2", nil)
	c.Params = gin.Params{{Key: "id", Value: "2"}}
	setBEndClaims(c, 2, stationA.ID, "station_manager")

	handler.Get(c)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
	resp := decodeResponseMap(t, w)
	if resp["msg"] != "forbidden" {
		t.Fatalf("expected forbidden message, got %v", resp["msg"])
	}
	_ = stationB
}

func TestZoneHandler_Create_StationManagerCannotCreateOtherStationZone(t *testing.T) {
	db := openHandlerTestDB(t, "zone_handler_create_scope.db")
	createHandlerTables(t, db)
	stationA := seedHandlerStation(t, db, "站点A")
	stationB := seedHandlerStation(t, db, "站点B")

	zoneRepo := repository.NewZoneRepository(db)
	handler := NewZoneHandler(service.NewZoneService(zoneRepo, service.NewGeofenceService(zoneRepo)))
	c, w := newJSONTestContext(t, http.MethodPost, "/b/zones", map[string]any{
		"station_id": stationB.ID,
		"name":       "跨站围栏",
		"points": []geo.Point{
			{Lat: 30.0, Lng: 120.0},
			{Lat: 30.1, Lng: 120.0},
			{Lat: 30.1, Lng: 120.1},
		},
	})
	setBEndClaims(c, 2, stationA.ID, "station_manager")

	handler.Create(c)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

func TestZoneHandler_Update_StationManagerCannotMoveZoneToOtherStation(t *testing.T) {
	db := openHandlerTestDB(t, "zone_handler_update_scope.db")
	createHandlerTables(t, db)
	stationA := seedHandlerStation(t, db, "站点A")
	stationB := seedHandlerStation(t, db, "站点B")
	zone := seedHandlerZone(t, db, stationA.ID, "原围栏", `[{"lat":30,"lng":120},{"lat":30.1,"lng":120},{"lat":30.1,"lng":120.1}]`)

	zoneRepo := repository.NewZoneRepository(db)
	handler := NewZoneHandler(service.NewZoneService(zoneRepo, service.NewGeofenceService(zoneRepo)))
	c, w := newJSONTestContext(t, http.MethodPut, "/b/zones/1", map[string]any{
		"station_id": stationB.ID,
		"name":       "跨站更新",
		"points": []geo.Point{
			{Lat: 30.0, Lng: 120.0},
			{Lat: 30.1, Lng: 120.0},
			{Lat: 30.1, Lng: 120.1},
		},
	})
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	setBEndClaims(c, 2, stationA.ID, "station_manager")

	handler.Update(c)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
	_ = zone
}
