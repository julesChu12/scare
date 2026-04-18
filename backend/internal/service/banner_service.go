package service

import (
	"errors"

	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/repository"
)

var (
	ErrBannerNotFound = errors.New("banner not found")
	ErrInvalidBanner  = errors.New("invalid banner")
)

type BannerService struct {
	repo *repository.BannerRepository
}

// NewBannerService 创建 BannerService
func NewBannerService(repo *repository.BannerRepository) *BannerService {
	return &BannerService{repo: repo}
}

// ListForStation 获取站点可见的 Banner（站点专属 + 全局）
func (s *BannerService) ListForStation(stationID int64) ([]*model.Banner, error) {
	return s.repo.ListActive(stationID)
}

// ListGlobal 获取全局 Banner（无站点时使用）
func (s *BannerService) ListGlobal() ([]*model.Banner, error) {
	return s.repo.ListGlobal()
}

// List 分页获取（管理端）
func (s *BannerService) List(page, pageSize int, stationID *int64) ([]*model.Banner, int64, error) {
	return s.repo.List(page, pageSize, stationID)
}

// BannerInput 创建/更新轮播图的输入参数
type BannerInput struct {
	ID        int64  `json:"id"`         // 轮播图ID（更新时必填）
	StationID int64  `json:"station_id"` // 站点ID（0表示全局）
	Title     string `json:"title"`      // 标题
	ImageURL  string `json:"image_url"`  // 图片URL
	LinkType  string `json:"link_type"`  // 链接类型（none/news/url）
	LinkValue string `json:"link_value"` // 链接值
	Sort      int32  `json:"sort"`       // 排序
	Status    string `json:"status"`     // 状态
}

// Create 创建轮播图
func (s *BannerService) Create(input BannerInput) (*model.Banner, error) {
	if input.ImageURL == "" {
		return nil, ErrInvalidBanner
	}

	banner := &model.Banner{
		StationID: input.StationID,
		Title:     input.Title,
		ImageURL:  input.ImageURL,
		LinkType:  input.LinkType,
		LinkValue: input.LinkValue,
		Sort:      input.Sort,
		Status:    input.Status,
	}

	if banner.LinkType == "" {
		banner.LinkType = "none"
	}
	if banner.Status == "" {
		banner.Status = "active"
	}

	if err := s.repo.Create(banner); err != nil {
		return nil, err
	}
	return banner, nil
}

// Update 更新轮播图
func (s *BannerService) Update(input BannerInput) (*model.Banner, error) {
	if input.ID == 0 {
		return nil, ErrInvalidBanner
	}

	banner, err := s.repo.GetByID(input.ID)
	if err != nil {
		return nil, ErrBannerNotFound
	}

	if input.ImageURL != "" {
		banner.ImageURL = input.ImageURL
	}
	banner.StationID = input.StationID
	banner.Title = input.Title
	banner.LinkType = input.LinkType
	banner.LinkValue = input.LinkValue
	banner.Sort = input.Sort
	if input.Status != "" {
		banner.Status = input.Status
	}

	if err := s.repo.Update(banner); err != nil {
		return nil, err
	}
	return banner, nil
}

// Delete 删除轮播图
func (s *BannerService) Delete(id int64) error {
	return s.repo.Delete(id)
}

