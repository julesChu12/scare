package repository

import (
	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/dao/query"

	"gorm.io/gorm"
)

type ZoneRepository struct {
	q *query.Query
}

func NewZoneRepository(db *gorm.DB) *ZoneRepository {
	return &ZoneRepository{
		q: query.Use(db),
	}
}

// Create 创建围栏
func (r *ZoneRepository) Create(zone *model.ServiceZone) error {
	return r.q.ServiceZone.Create(zone)
}

// GetByID 根据ID获取围栏
func (r *ZoneRepository) GetByID(id int64) (*model.ServiceZone, error) {
	z := r.q.ServiceZone
	return z.Where(z.ID.Eq(id)).First()
}

// GetByStationID 根据站点ID获取围栏列表
func (r *ZoneRepository) GetByStationID(stationID int64) ([]*model.ServiceZone, error) {
	z := r.q.ServiceZone
	return z.Where(z.StationID.Eq(stationID), z.Status.Eq("active")).
		Order(z.Priority.Desc(), z.ID.Asc()).
		Find()
}

// ListActive 获取所有活跃围栏
func (r *ZoneRepository) ListActive() ([]*model.ServiceZone, error) {
	z := r.q.ServiceZone
	return z.Where(z.Status.Eq("active")).
		Order(z.Priority.Desc(), z.ID.Asc()).
		Find()
}

// List 获取所有围栏（分页）
func (r *ZoneRepository) List(offset, limit int) ([]*model.ServiceZone, int64, error) {
	z := r.q.ServiceZone
	
	// 获取总数
	total, err := z.Count()
	if err != nil {
		return nil, 0, err
	}
	
	// 分页查询
	zones, err := z.Order(z.ID.Desc()).
		Offset(offset).
		Limit(limit).
		Find()
	
	if err != nil {
		return nil, 0, err
	}
	
	return zones, total, nil
}

// Update 更新围栏
func (r *ZoneRepository) Update(zone *model.ServiceZone) error {
	z := r.q.ServiceZone
	_, err := z.Where(z.ID.Eq(zone.ID)).Updates(zone)
	return err
}

// Delete 删除围栏
func (r *ZoneRepository) Delete(id int64) error {
	z := r.q.ServiceZone
	_, err := z.Where(z.ID.Eq(id)).Delete()
	return err
}
