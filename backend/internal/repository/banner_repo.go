package repository

import (
	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/dao/query"

	"gorm.io/gorm"
)

type BannerRepository struct {
	q *query.Query
}

func NewBannerRepository(db *gorm.DB) *BannerRepository {
	return &BannerRepository{q: query.Use(db)}
}

// ListActive 获取激活的 Banner（站点专属 + 全局）
// 优先级：站点专属 > 全局，按 sort 降序
func (r *BannerRepository) ListActive(stationID int64) ([]*model.Banner, error) {
	b := r.q.Banner
	db := b.UnderlyingDB().Model(&model.Banner{}).
		Where("status = ?", "active").
		Where("station_id = ? OR station_id = 0", stationID).
		Order("station_id DESC, sort DESC")

	var banners []*model.Banner
	if err := db.Find(&banners).Error; err != nil {
		return nil, err
	}
	return banners, nil
}

// ListGlobal 获取全局 Banner（station_id = 0）
func (r *BannerRepository) ListGlobal() ([]*model.Banner, error) {
	b := r.q.Banner
	return b.Where(b.StationID.Eq(0), b.Status.Eq("active")).
		Order(b.Sort.Desc()).Find()
}

// List 分页获取所有 Banner（管理端）
func (r *BannerRepository) List(page, pageSize int, stationID *int64) ([]*model.Banner, int64, error) {
	b := r.q.Banner
	db := b.UnderlyingDB().Model(&model.Banner{}).
		Select("banners.id, banners.station_id, banners.title, banners.image_url, banners.link_type, banners.link_value, banners.sort, banners.status, banners.created_at, banners.updated_at, banners.deleted_at, CASE WHEN banners.station_id = 0 THEN '全局' ELSE service_stations.name END as station_name").
		Joins("left join service_stations on banners.station_id = service_stations.id")

	if stationID != nil && *stationID > 0 {
		db = db.Where("banners.station_id = ?", *stationID)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var banners []*model.Banner
	offset := (page - 1) * pageSize
	if err := db.Order("banners.station_id DESC, banners.sort DESC, banners.id DESC").
		Offset(offset).Limit(pageSize).Scan(&banners).Error; err != nil {
		return nil, 0, err
	}

	return banners, total, nil
}

// GetByID 根据 ID 获取
func (r *BannerRepository) GetByID(id int64) (*model.Banner, error) {
	b := r.q.Banner
	return b.Where(b.ID.Eq(id)).First()
}

// Create 创建
func (r *BannerRepository) Create(banner *model.Banner) error {
	return r.q.Banner.Create(banner)
}

// Update 更新
func (r *BannerRepository) Update(banner *model.Banner) error {
	b := r.q.Banner
	_, err := b.Where(b.ID.Eq(banner.ID)).Updates(banner)
	return err
}

// Delete 删除（软删除）
func (r *BannerRepository) Delete(id int64) error {
	b := r.q.Banner
	_, err := b.Where(b.ID.Eq(id)).Delete()
	return err
}
