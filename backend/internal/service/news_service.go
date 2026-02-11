package service

import (
	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/repository"
)

type NewsService struct {
	newsRepo *repository.NewsRepository
}

func NewNewsService(newsRepo *repository.NewsRepository) *NewsService {
	return &NewsService{newsRepo: newsRepo}
}

func (s *NewsService) ListPublished(page, pageSize int, newsType string, stationID *int64) ([]*model.News, int64, error) {
	return s.newsRepo.ListPublished(page, pageSize, newsType, stationID)
}

func (s *NewsService) GetDetail(id int64) (*model.News, error) {
	news, err := s.newsRepo.GetPublishedByID(id)
	if err != nil {
		return nil, err
	}
	go func() {
		_ = s.newsRepo.IncrementViewCount(id)
	}()
	return news, nil
}

// GetByID 根据ID获取新闻（B端管理用，不限状态）
func (s *NewsService) GetByID(id int64) (*model.News, error) {
	return s.newsRepo.GetByID(id)
}

// Create 创建新闻
func (s *NewsService) Create(news *model.News) error {
	if news.Status == "" {
		news.Status = "draft"
	}
	return s.newsRepo.Create(news)
}

// Update 更新新闻
func (s *NewsService) Update(news *model.News) error {
	return s.newsRepo.Update(news)
}

// Delete 删除新闻
func (s *NewsService) Delete(id int64) error {
	return s.newsRepo.Delete(id)
}
