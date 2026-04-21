package repository

import (
	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/dao/query"

	"gorm.io/gorm"
)

type NewsRepository struct {
	q *query.Query
}

func NewNewsRepository(db *gorm.DB) *NewsRepository {
	return &NewsRepository{
		q: query.Use(db),
	}
}

// ListPublished 获取已发布的新闻列表（分页）
// 优先级：站点专属 > 全局（station_id=0），按发布时间降序
func (r *NewsRepository) ListPublished(page, pageSize int, newsType string, stationID *int64) ([]*model.News, int64, error) {
	n := r.q.News

	// 对于复杂的 OR 查询，使用原生 DB
	db := n.UnderlyingDB().Model(&model.News{}).
		Select("news.id, news.station_id, news.title, news.summary, news.content, news.cover_url, news.type, news.status, news.author_id, news.publish_at, news.view_count, news.created_at, news.updated_at, news.deleted_at, CASE WHEN news.station_id = 0 THEN '全局' ELSE service_stations.name END as station_name").
		Joins("left join service_stations on news.station_id = service_stations.id").
		Where("news.status = ?", "published")

	if newsType != "" {
		db = db.Where("news.type = ?", newsType)
	}

	if stationID != nil && *stationID > 0 {
		// 站点新闻 + 全局新闻（station_id=0）
		db = db.Where("news.station_id = ? OR news.station_id = 0", *stationID)
	} else {
		// 只显示全局新闻
		db = db.Where("news.station_id = 0")
	}

	// 获取总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询：站点专属优先，然后按发布时间降序
	var newsList []*model.News
	offset := (page - 1) * pageSize
	if err := db.Order("news.station_id DESC, news.publish_at DESC").Offset(offset).Limit(pageSize).Scan(&newsList).Error; err != nil {
		return nil, 0, err
	}

	return newsList, total, nil
}

// List 获取新闻列表（管理端）
func (r *NewsRepository) List(page, pageSize int, newsType string, status string, stationID *int64) ([]*model.News, int64, error) {
	n := r.q.News
	db := n.UnderlyingDB().Model(&model.News{}).
		Select("news.id, news.station_id, news.title, news.summary, news.content, news.cover_url, news.type, news.status, news.author_id, news.publish_at, news.view_count, news.created_at, news.updated_at, news.deleted_at, CASE WHEN news.station_id = 0 THEN '全局' ELSE service_stations.name END as station_name").
		Joins("left join service_stations on news.station_id = service_stations.id")

	if newsType != "" {
		db = db.Where("news.type = ?", newsType)
	}
	if status != "" {
		db = db.Where("news.status = ?", status)
	}
	if stationID != nil && *stationID > 0 {
		db = db.Where("news.station_id = ?", *stationID)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var newsList []*model.News
	offset := (page - 1) * pageSize
	if err := db.Order("news.id DESC").Offset(offset).Limit(pageSize).Scan(&newsList).Error; err != nil {
		return nil, 0, err
	}

	return newsList, total, nil
}

// GetByID 根据ID获取新闻
func (r *NewsRepository) GetByID(id int64) (*model.News, error) {
	n := r.q.News
	return n.Where(n.ID.Eq(id)).First()
}

// Create 创建新闻
func (r *NewsRepository) Create(news *model.News) error {
	return r.q.News.Create(news)
}

// Update 更新新闻
func (r *NewsRepository) Update(news *model.News) error {
	n := r.q.News
	_, err := n.Where(n.ID.Eq(news.ID)).Updates(news)
	return err
}

// Delete 删除新闻（软删除）
func (r *NewsRepository) Delete(id int64) error {
	n := r.q.News
	_, err := n.Where(n.ID.Eq(id)).Delete()
	return err
}
