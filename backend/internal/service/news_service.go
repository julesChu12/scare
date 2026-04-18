package service

import (
	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/repository"
)

// NewsService 新闻服务
type NewsService struct {
	repo *repository.NewsRepository
}

// NewNewsService 创建新闻服务
func NewNewsService(repo *repository.NewsRepository) *NewsService {
	return &NewsService{repo: repo}
}

// ListPublished 获取已发布的新闻列表
func (s *NewsService) ListPublished(page, pageSize int, newsType string, stationID *int64) ([]*model.News, int64, error) {
	return s.repo.ListPublished(page, pageSize, newsType, stationID)
}

// GetByID 根据 ID 获取新闻详情
func (s *NewsService) GetByID(id int64) (*model.News, error) {
	return s.repo.GetByID(id)
}

// Create 创建新闻
func (s *NewsService) Create(news *model.News) error {
	return s.repo.Create(news)
}

// Update 更新新闻
func (s *NewsService) Update(news *model.News) error {
	return s.repo.Update(news)
}

// Delete 删除新闻
func (s *NewsService) Delete(id int64) error {
	return s.repo.Delete(id)
}
