package service

import (
	"encoding/json"
	"sync"

	"community-elderly-care-platform/internal/repository"
	"community-elderly-care-platform/pkg/geo"
)

// GeofenceService 地理围栏服务
type GeofenceService struct {
	zoneRepo *repository.ZoneRepository
	mu       sync.RWMutex
	engine   *geo.Engine
}

func NewGeofenceService(zoneRepo *repository.ZoneRepository) *GeofenceService {
	return &GeofenceService{
		zoneRepo: zoneRepo,
	}
}

// Reload 重新加载围栏数据
func (s *GeofenceService) Reload() error {
	zones, err := s.zoneRepo.ListActive()
	if err != nil {
		return err
	}

	engineZones := make([]geo.Zone, 0, len(zones))
	for _, zone := range zones {
		var points []geo.Point
		if err := json.Unmarshal([]byte(zone.Points), &points); err != nil {
			continue
		}
		engineZones = append(engineZones, geo.Zone{
			ID:        zone.ID,
			StationID: zone.StationID,
			Priority:  int(zone.Priority),
			Points:    points,
		})
	}

	engine := geo.NewEngine(engineZones)

	s.mu.Lock()
	s.engine = engine
	s.mu.Unlock()

	return nil
}

// Match 匹配坐标点所属的服务站点
func (s *GeofenceService) Match(lat, lng float64) (int64, bool) {
	s.mu.RLock()
	engine := s.engine
	s.mu.RUnlock()

	if engine == nil {
		return 0, false
	}
	return engine.Match(geo.Point{Lat: lat, Lng: lng})
}

// Engine 获取当前引擎实例
func (s *GeofenceService) Engine() (*geo.Engine, bool) {
	s.mu.RLock()
	engine := s.engine
	s.mu.RUnlock()

	if engine == nil {
		return nil, false
	}
	return engine, true
}
