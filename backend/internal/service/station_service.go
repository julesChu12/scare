package service

import (
	"errors"

	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/repository"
)

var ErrInvalidStation = errors.New("invalid station")

// DefaultStationName 默认服务站点名称
const DefaultStationName = "sCare系统默认服务站点"

type StationService struct {
	repo            *repository.StationRepository
	customerRepo    *repository.CustomerRepository
	geofenceService *GeofenceService
	geocodeService  *GeocodeService
}

func NewStationService(repo *repository.StationRepository) *StationService {
	return &StationService{repo: repo}
}

// SetCustomerRepo 设置客户档案仓库（可选依赖）
func (s *StationService) SetCustomerRepo(repo *repository.CustomerRepository) {
	s.customerRepo = repo
}

// SetGeofenceService 设置围栏服务（可选依赖）
func (s *StationService) SetGeofenceService(svc *GeofenceService) {
	s.geofenceService = svc
}

// SetGeocodeService 设置地理编码服务（可选依赖）
func (s *StationService) SetGeocodeService(svc *GeocodeService) {
	s.geocodeService = svc
}

func (s *StationService) Create(station *model.ServiceStation) error {
	if station.Name == "" {
		return ErrInvalidStation
	}
	return s.repo.Create(station)
}

func (s *StationService) GetByID(id int64) (*model.ServiceStation, error) {
	return s.repo.GetByID(id)
}

func (s *StationService) List(page, pageSize int) ([]*model.ServiceStation, int64, error) {
	offset := (page - 1) * pageSize
	return s.repo.List(offset, pageSize)
}

func (s *StationService) Update(station *model.ServiceStation) error {
	if station.ID == 0 || station.Name == "" {
		return ErrInvalidStation
	}
	return s.repo.Update(station)
}

func (s *StationService) Delete(id int64) error {
	if id == 0 {
		return ErrInvalidStation
	}
	return s.repo.Delete(id)
}

// MatchStationInput 匹配站点的输入参数
type MatchStationInput struct {
	UserID    int64    `json:"user_id"`   // 用户ID（可选）
	Latitude  *float64 `json:"latitude"`  // 纬度（可选）
	Longitude *float64 `json:"longitude"` // 经度（可选）
}

// MatchStation 匹配服务站点
// 匹配优先级：
// 1. 如果提供了经纬度，使用 geofence 匹配或最近站点
// 2. 如果没有经纬度但有用户ID，获取用户 profile 中的地址，geocode 后匹配
// 3. 如果都没有，返回默认站点
func (s *StationService) MatchStation(input MatchStationInput) (*model.ServiceStation, error) {
	var lat, lng float64
	hasCoords := false

	// 优先使用传入的经纬度
	if input.Latitude != nil && input.Longitude != nil {
		lat = *input.Latitude
		lng = *input.Longitude
		hasCoords = true
	}

	// 如果没有经纬度但有用户ID，尝试从用户档案获取地址并解析
	if !hasCoords && input.UserID > 0 && s.customerRepo != nil && s.geocodeService != nil {
		profile, err := s.customerRepo.GetByUserID(input.UserID)
		if err == nil && profile != nil && profile.Address != "" {
			result, err := s.geocodeService.Geocode(profile.Address)
			if err == nil && result != nil {
				lat = result.Latitude
				lng = result.Longitude
				hasCoords = true
			}
		}
	}

	// 如果有坐标，尝试匹配
	if hasCoords {
		// 优先使用 geofence 匹配
		if s.geofenceService != nil {
			if stationID, ok := s.geofenceService.Match(lat, lng); ok && stationID > 0 {
				station, err := s.repo.GetByID(stationID)
				if err == nil && station != nil {
					return station, nil
				}
			}
		}

		// geofence 未匹配到，查找最近的站点
		station, err := s.repo.FindNearest(lat, lng)
		if err == nil && station != nil {
			return station, nil
		}
	}

	// 返回默认站点
	return s.getDefaultStation()
}

// getDefaultStation 获取默认服务站点
func (s *StationService) getDefaultStation() (*model.ServiceStation, error) {
	station, err := s.repo.GetByName(DefaultStationName)
	if err != nil {
		// 如果默认站点不存在，返回任意一个活跃站点
		stations, _, listErr := s.repo.List(0, 1)
		if listErr != nil {
			return nil, listErr
		}
		if len(stations) > 0 {
			return stations[0], nil
		}
		return nil, errors.New("no station available")
	}
	return station, nil
}
