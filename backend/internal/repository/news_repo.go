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
	db := n.UnderlyingDB().Model(&model.News{}).Where("status = ?", "published")
	
	if newsType != "" {
		db = db.Where("type = ?", newsType)
	}
	
	if stationID != nil && *stationID > 0 {
		// 站点新闻 + 全局新闻（station_id=0）
		db = db.Where("station_id = ? OR station_id = 0", *stationID)
	} else {
		// 只显示全局新闻
		db = db.Where("station_id = 0")
	}
	
	// 获取总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 分页查询：站点专属优先，然后按发布时间降序
	var newsList []*model.News
	offset := (page - 1) * pageSize
	if err := db.Order("station_id DESC, publish_at DESC").Offset(offset).Limit(pageSize).Find(&newsList).Error; err != nil {
		return nil, 0, err
	}
	
	return newsList, total, nil
}

// GetByID 根据ID获取新闻
func (r *NewsRepository) GetByID(id int64) (*model.News, error) {
	n := r.q.News
	return n.Where(n.ID.Eq(id)).First()
}

// GetPublishedByID 获取已发布的新闻详情
func (r *NewsRepository) GetPublishedByID(id int64) (*model.News, error) {
	n := r.q.News
	return n.Where(n.ID.Eq(id), n.Status.Eq("published")).First()
}

// IncrementViewCount 增加浏览次数
func (r *NewsRepository) IncrementViewCount(id int64) error {
	n := r.q.News
	_, err := n.Where(n.ID.Eq(id)).UpdateSimple(n.ViewCount.Add(1))
	return err
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
