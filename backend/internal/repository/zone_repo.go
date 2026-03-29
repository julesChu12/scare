package repository

import (
	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/dao/query"

	"gorm.io/gorm"
)

type ZoneRepository struct {
	q *query.Query
}

type ZoneListFilter struct {
	StationID int64
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

// ListActive 获取所有活跃围栏
func (r *ZoneRepository) ListActive() ([]*model.ServiceZone, error) {
	z := r.q.ServiceZone
	return z.Where(z.Status.Eq("active")).
		Order(z.Priority.Desc(), z.ID.Asc()).
		Find()
}

// List 获取所有围栏（分页）
func (r *ZoneRepository) List(offset, limit int, filter ZoneListFilter) ([]*model.ServiceZone, int64, error) {
	db := r.q.ServiceZone.UnderlyingDB().Model(&model.ServiceZone{})

	if filter.StationID > 0 {
		db = db.Where("station_id = ?", filter.StationID)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var zones []*model.ServiceZone
	if err := db.Order("id DESC").Offset(offset).Limit(limit).Find(&zones).Error; err != nil {
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
