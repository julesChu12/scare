package repository

import (
	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/dao/query"

	"gorm.io/gorm"
)

type StationRepository struct {
	q *query.Query
}

type StationListFilter struct {
	Keyword   string
	Status    string
	StationID int64
}

func NewStationRepository(db *gorm.DB) *StationRepository {
	return &StationRepository{
		q: query.Use(db),
	}
}

// Create 创建站点
func (r *StationRepository) Create(station *model.ServiceStation) error {
	return r.q.ServiceStation.Create(station)
}

// GetByID 根据ID获取站点
func (r *StationRepository) GetByID(id int64) (*model.ServiceStation, error) {
	s := r.q.ServiceStation
	return s.Where(s.ID.Eq(id)).First()
}

// List 获取所有站点（分页）
func (r *StationRepository) List(offset, limit int, filter StationListFilter) ([]*model.ServiceStation, int64, error) {
	db := r.q.ServiceStation.UnderlyingDB().Model(&model.ServiceStation{})

	if filter.StationID > 0 {
		db = db.Where("id = ?", filter.StationID)
	}

	if filter.Keyword != "" {
		keyword := "%" + filter.Keyword + "%"
		db = db.Where(
			"name LIKE ? OR code LIKE ? OR address LIKE ? OR phone LIKE ?",
			keyword, keyword, keyword, keyword,
		)
	}

	if filter.Status != "" {
		db = db.Where("status = ?", filter.Status)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var stations []*model.ServiceStation
	if err := db.Order("id DESC").Offset(offset).Limit(limit).Find(&stations).Error; err != nil {
		return nil, 0, err
	}

	return stations, total, nil
}

// Update 更新站点
func (r *StationRepository) Update(station *model.ServiceStation) error {
	s := r.q.ServiceStation
	_, err := s.Where(s.ID.Eq(station.ID)).Updates(station)
	return err
}

// Delete 删除站点
func (r *StationRepository) Delete(id int64) error {
	s := r.q.ServiceStation
	_, err := s.Where(s.ID.Eq(id)).Delete()
	return err
}

// ListActive 获取所有活跃站点
func (r *StationRepository) ListActive() ([]*model.ServiceStation, error) {
	s := r.q.ServiceStation
	return s.Where(s.Status.Eq("active")).Order(s.ID.Asc()).Find()
}

// GetByName 根据名称获取站点
func (r *StationRepository) GetByName(name string) (*model.ServiceStation, error) {
	s := r.q.ServiceStation
	return s.Where(s.Name.Eq(name), s.Status.Eq("active")).First()
}
