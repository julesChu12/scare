package service

import (
	"encoding/json"
	"errors"

	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/repository"
	"community-elderly-care-platform/pkg/geo"
)

var ErrInvalidZone = errors.New("invalid zone")

type ZoneService struct {
	repo        *repository.ZoneRepository
	geofenceSvc *GeofenceService
}

// ZoneInput 创建/更新服务围栏的输入参数
type ZoneInput struct {
	ID        int64       `json:"id"`         // 围栏ID（更新时必填）
	StationID int64       `json:"station_id"` // 站点ID
	Name      string      `json:"name"`       // 围栏名称
	Points    []geo.Point `json:"points"`     // 围栏顶点坐标列表
	Priority  int         `json:"priority"`   // 优先级
	Status    string      `json:"status"`     // 状态
}

type ZoneListFilter struct {
	StationID int64 `json:"station_id"`
}

// NewZoneService 创建 ZoneService
func NewZoneService(repo *repository.ZoneRepository, geofenceSvc *GeofenceService) *ZoneService {
	return &ZoneService{repo: repo, geofenceSvc: geofenceSvc}
}

// Create 创建服务围栏
func (s *ZoneService) Create(input ZoneInput) (*model.ServiceZone, error) {
	if input.StationID == 0 || input.Name == "" || len(input.Points) < 3 {
		return nil, ErrInvalidZone
	}
	points, err := json.Marshal(input.Points)
	if err != nil {
		return nil, err
	}
	zone := &model.ServiceZone{
		StationID: input.StationID,
		Name:      input.Name,
		Points:    string(points),
		Priority:  int64(input.Priority),
		Status:    input.Status,
	}
	if zone.Status == "" {
		zone.Status = "active"
	}

	if err := s.repo.Create(zone); err != nil {
		return nil, err
	}
	_ = s.geofenceSvc.Reload()
	return zone, nil
}

// Update 更新服务围栏
func (s *ZoneService) Update(input ZoneInput) (*model.ServiceZone, error) {
	if input.ID == 0 || input.StationID == 0 || input.Name == "" || len(input.Points) < 3 {
		return nil, ErrInvalidZone
	}
	zone, err := s.repo.GetByID(input.ID)
	if err != nil {
		return nil, err
	}
	points, err := json.Marshal(input.Points)
	if err != nil {
		return nil, err
	}
	zone.StationID = input.StationID
	zone.Name = input.Name
	zone.Points = string(points)
	zone.Priority = int64(input.Priority)
	if input.Status != "" {
		zone.Status = input.Status
	}

	if err := s.repo.Update(zone); err != nil {
		return nil, err
	}
	_ = s.geofenceSvc.Reload()
	return zone, nil
}

// GetByID 根据 ID 获取围栏详情
func (s *ZoneService) GetByID(id int64) (*model.ServiceZone, error) {
	if id == 0 {
		return nil, ErrInvalidZone
	}
	return s.repo.GetByID(id)
}

// Delete 删除服务围栏
func (s *ZoneService) Delete(id int64) error {
	if id == 0 {
		return ErrInvalidZone
	}
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	_ = s.geofenceSvc.Reload()
	return nil
}

// List 分页获取围栏列表
func (s *ZoneService) List(page, pageSize int, filter ZoneListFilter) ([]*model.ServiceZone, int64, error) {
	offset := (page - 1) * pageSize
	return s.repo.List(offset, pageSize, repository.ZoneListFilter{
		StationID: filter.StationID,
	})
}
