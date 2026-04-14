package handler

import (
	"net/http"
	"strconv"
	"testing"

	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/repository"
	"community-elderly-care-platform/internal/service"

	"github.com/gin-gonic/gin"
)

func TestStationHandler_Create_PersistsAddressPhoneAndCoords(t *testing.T) {
	db := openHandlerTestDB(t, "station_handler_create.db")
	createHandlerTables(t, db)

	handler := NewStationHandler(service.NewStationService(repository.NewStationRepository(db)))
	c, w := newJSONTestContext(t, http.MethodPost, "/b/stations", map[string]any{
		"name":      "朝阳站点",
		"code":      "CY001",
		"address":   "北京市朝阳区测试路 1 号",
		"phone":     "13800138000",
		"latitude":  39.912345,
		"longitude": 116.456789,
		"status":    "active",
	})

	handler.Create(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	resp := decodeResponseMap(t, w)
	data := resp["data"].(map[string]any)
	stationID := int64(data["id"].(float64))

	var station model.ServiceStation
	if err := db.First(&station, stationID).Error; err != nil {
		t.Fatalf("failed to query station: %v", err)
	}

	if station.Address != "北京市朝阳区测试路 1 号" {
		t.Fatalf("expected address persisted, got %q", station.Address)
	}
	if station.Phone != "13800138000" {
		t.Fatalf("expected phone persisted, got %q", station.Phone)
	}
	if station.Latitude != 39.912345 {
		t.Fatalf("expected latitude persisted, got %v", station.Latitude)
	}
	if station.Longitude != 116.456789 {
		t.Fatalf("expected longitude persisted, got %v", station.Longitude)
	}
}

func TestStationHandler_Update_PersistsAddressPhoneAndCoords(t *testing.T) {
	db := openHandlerTestDB(t, "station_handler_update.db")
	createHandlerTables(t, db)
	station := seedHandlerStation(t, db, "待更新站点")

	handler := NewStationHandler(service.NewStationService(repository.NewStationRepository(db)))
	c, w := newJSONTestContext(t, http.MethodPut, "/b/stations/1", map[string]any{
		"name":      "更新后站点",
		"code":      "UP001",
		"address":   "上海市浦东新区示例路 8 号",
		"phone":     "13900139000",
		"latitude":  31.230416,
		"longitude": 121.473701,
		"status":    "inactive",
	})
	c.Params = gin.Params{{Key: "id", Value: strconv.FormatInt(station.ID, 10)}}

	handler.Update(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var updated model.ServiceStation
	if err := db.First(&updated, station.ID).Error; err != nil {
		t.Fatalf("failed to query station: %v", err)
	}

	if updated.Name != "更新后站点" {
		t.Fatalf("expected updated name, got %q", updated.Name)
	}
	if updated.Address != "上海市浦东新区示例路 8 号" {
		t.Fatalf("expected updated address, got %q", updated.Address)
	}
	if updated.Phone != "13900139000" {
		t.Fatalf("expected updated phone, got %q", updated.Phone)
	}
	if updated.Latitude != 31.230416 {
		t.Fatalf("expected updated latitude, got %v", updated.Latitude)
	}
	if updated.Longitude != 121.473701 {
		t.Fatalf("expected updated longitude, got %v", updated.Longitude)
	}
	if updated.Status != "inactive" {
		t.Fatalf("expected updated status, got %q", updated.Status)
	}
}
