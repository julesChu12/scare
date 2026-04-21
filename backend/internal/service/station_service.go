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

type StationListFilter struct {
	Keyword   string
	Status    string
	StationID int64
}

// NewStationService 创建 StationService
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

// Create 创建服务站点
func (s *StationService) Create(station *model.ServiceStation) error {
	if station.Name == "" {
		return ErrInvalidStation
	}
	return s.repo.Create(station)
}

// GetByID 根据 ID 获取站点详情
func (s *StationService) GetByID(id int64) (*model.ServiceStation, error) {
	return s.repo.GetByID(id)
}

// List 分页获取站点列表
func (s *StationService) List(page, pageSize int, filter StationListFilter) ([]*model.ServiceStation, int64, error) {
	offset := (page - 1) * pageSize
	return s.repo.List(offset, pageSize, repository.StationListFilter{
		Keyword:   filter.Keyword,
		Status:    filter.Status,
		StationID: filter.StationID,
	})
}

// Update 更新服务站点
func (s *StationService) Update(station *model.ServiceStation) error {
	if station.ID == 0 || station.Name == "" {
		return ErrInvalidStation
	}
	return s.repo.Update(station)
}

// Delete 删除服务站点
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
//
// 匹配优先级（从高到低）：
//   1. 直接提供经纬度 → 射线法围栏匹配 → Haversine 兜底
//   2. 提供用户ID（有档案地址）→ 地址 geocode → 射线法/Haversine
//   3. 上述都无 → 返回默认站点（sCare系统默认服务站点）
//
// 坐标获取策略：
//   优先使用前端直接传入的经纬度（如 GPS 定位）；
//   无经纬度但有用户ID → 读取用户档案中的地址文本 → 调用 GeocodeService 转为坐标。
//
// 在射线法链路中的位置（输入有经纬度的场景）：
//   MatchStation(lat, lng)
//     → GeofenceService.Match(lat, lng)
//         → geo.Engine.Match()
//             → BoundingBox.Contains()     // 快速排除
//             → PointInPolygon()            // 射线法精确判断
//     → nearestStationByHaversine()       // 兜底
//     → getDefaultStation()               // 最终保底
//
func (s *StationService) MatchStation(input MatchStationInput) (*model.ServiceStation, error) {
	var lat, lng float64
	hasCoords := false

	// 优先级一：直接使用传入的经纬度（前端 GPS 或手动选择坐标）
	if input.Latitude != nil && input.Longitude != nil {
		lat = *input.Latitude
		lng = *input.Longitude
		hasCoords = true
	}

	// 优先级二：无经纬度但有用户ID → 从用户档案地址反推坐标
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

	// 优先级三：有坐标 → 射线法 + Haversine
	if hasCoords {
		// 第一步：射线法围栏匹配
		if s.geofenceService != nil {
			if stationID, ok := s.geofenceService.Match(lat, lng); ok && stationID > 0 {
				station, err := s.repo.GetByID(stationID)
				if err == nil && station != nil && station.Status == "active" {
					return station, nil
				}
			}
		}

		// 第二步：Haversine 找最近站点（射线法未命中时的兜底）
		stations, err := s.repo.ListActive()
		if err == nil {
			nearest, nearestErr := nearestStationByHaversine(stations, lat, lng)
			if nearestErr == nil {
				return nearest, nil
			}
		}
	}

	// 优先级四：全无 → 返回默认站点（保底策略）
	return s.getDefaultStation()
}

// getDefaultStation 获取默认服务站点（最终保底）
//
// 兜底策略：
//   1. 优先查找名为 "sCare系统默认服务站点" 的站点
//   2. 若默认站点不存在，返回任意一个活跃站点（只要有一个站点就能服务）
//   3. 若没有任何活跃站点 → 返回错误（系统配置异常）
//
func (s *StationService) getDefaultStation() (*model.ServiceStation, error) {
	station, err := s.repo.GetByName(DefaultStationName)
	if err != nil {
		// 如果默认站点不存在，返回任意一个活跃站点
		stations, _, listErr := s.repo.List(0, 1, repository.StationListFilter{Status: "active"})
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
