package repository

import (
	"fmt"

	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/dao/query"

	"gorm.io/gorm"
)

type StationRepository struct {
	q *query.Query
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
func (r *StationRepository) List(offset, limit int) ([]*model.ServiceStation, int64, error) {
	s := r.q.ServiceStation
	
	// 获取总数
	total, err := s.Count()
	if err != nil {
		return nil, 0, err
	}
	
	// 分页查询
	stations, err := s.Order(s.ID.Desc()).
		Offset(offset).
		Limit(limit).
		Find()
	
	if err != nil {
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

// FindNearest 查找距离指定坐标最近的活跃站点
func (r *StationRepository) FindNearest(lat, lng float64) (*model.ServiceStation, error) {
	s := r.q.ServiceStation
	// 使用原生 DB 处理复杂的距离计算
	db := s.UnderlyingDB()
	var station model.ServiceStation
	orderClause := fmt.Sprintf("(POW(latitude - %f, 2) + POW(longitude - %f, 2))", lat, lng)
	err := db.Where("status = ? AND latitude != 0 AND longitude != 0", "active").
		Order(orderClause).
		First(&station).Error
	if err != nil {
		return nil, err
	}
	return &station, nil
}
