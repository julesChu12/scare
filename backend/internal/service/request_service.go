package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"community-elderly-care-platform/internal/consts"
	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/repository"

	"gorm.io/gorm"
)

var (
	ErrInvalidRequest  = errors.New("invalid request")
	ErrRequestConflict = errors.New("request conflict")
	ErrNotCompleted    = errors.New("request not completed")
	ErrAlreadyRated    = errors.New("request already rated")
)

type RequestService struct {
	db          *gorm.DB
	repo        *repository.RequestRepository
	taskRepo    *repository.TaskRepository
	stationRepo *repository.StationRepository
	geofenceSvc *GeofenceService
	geocodeSvc  *GeocodeService
}

// RequestInput 创建服务请求的输入参数
type RequestInput struct {
	UserID       int64    `json:"user_id"`       // 用户ID
	RequestNo    string   `json:"request_no"`    // 请求编号（可选，用于幂等）
	ServiceType  string   `json:"service_type"`  // 服务类型
	Lat          *float64 `json:"lat"`           // 纬度（可选）
	Lng          *float64 `json:"lng"`           // 经度（可选）
	ContactName  string   `json:"contact_name"`  // 联系人姓名
	ContactPhone string   `json:"contact_phone"` // 联系人电话
	Address      string   `json:"address"`       // 地址
	Images       []string `json:"images"`        // 图片列表
}

func NewRequestService(db *gorm.DB, repo *repository.RequestRepository, taskRepo *repository.TaskRepository, stationRepo *repository.StationRepository, geofenceSvc *GeofenceService, geocodeSvc *GeocodeService) *RequestService {
	return &RequestService{
		db:          db,
		repo:        repo,
		taskRepo:    taskRepo,
		stationRepo: stationRepo,
		geofenceSvc: geofenceSvc,
		geocodeSvc:  geocodeSvc,
	}
}

// Create 创建服务请求
//
// 业务目的：客户提交服务需求，系统自动匹配站点并创建任务
//
// 主要流程：
// 1. 验证输入参数（用户ID、服务类型必填）
// 2. 解析坐标（优先使用传入坐标，否则通过地址反向地理编码）
// 3. 幂等性检查：如果提供了 RequestNo，检查是否已存在
// 4. 生成请求编号（REQ + 时间戳 + 随机数）
// 5. 匹配服务站点：
//    a. 首先检查地理围栏（点是否在多边形内）
//    b. 如果没有匹配的围栏，查找最近的站点（Haversine 距离计算）
// 6. 在事务中创建服务请求和对应任务（状态为 dispatched）
//
// 坐标解析优先级：显式传入坐标 > 地址反向地理编码
//
// 返回值：
// - request: 服务请求信息
// - created: 是否新创建（false 表示幂等返回已存在的请求）
// - error: 错误信息
func (s *RequestService) Create(input RequestInput) (*model.ServiceRequest, bool, error) {
	if input.UserID == 0 || input.ServiceType == "" {
		return nil, false, ErrInvalidRequest
	}
	if !consts.IsValidServiceType(input.ServiceType) {
		return nil, false, ErrInvalidRequest
	}

	lat, lng, resolvedAddress, err := s.resolveCoordinates(input)
	if err != nil {
		return nil, false, err
	}

	if input.RequestNo != "" {
		existing, err := s.repo.GetByRequestNo(input.RequestNo)
		if err == nil {
			if existing.UserID != input.UserID {
				return nil, false, ErrRequestConflict
			}
			return existing, false, nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, err
		}
	}

	requestNo := input.RequestNo
	if requestNo == "" {
		requestNo = generateRequestNo()
	}

	stationID, matched := s.geofenceSvc.Match(lat, lng)
	if !matched {
		nearest, err := s.findNearestStation(lat, lng)
		if err != nil {
			return nil, false, err
		}
		stationID = nearest
	}

	images, err := json.Marshal(input.Images)
	if err != nil {
		return nil, false, err
	}

	request := &model.ServiceRequest{
		RequestNo:         requestNo,
		UserID:            input.UserID,
		ServiceType:       input.ServiceType,
		Status:            consts.RequestStatusDispatched,
		SubmitLocationLat: lat,
		SubmitLocationLng: lng,
		ContactName:       input.ContactName,
		ContactPhone:      input.ContactPhone,
		Address:           resolvedAddress,
		StationID:         stationID,
		Images:            string(images),
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		reqRepo := s.repo.WithTx(tx)
		taskRepo := s.taskRepo.WithTx(tx)
		if err := reqRepo.Create(request); err != nil {
			return err
		}
		task := &model.TaskAssignment{
			RequestID: request.ID,
			StationID: stationID,
			Status:    consts.TaskStatusDispatched,
		}
		if err := taskRepo.Create(task); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, false, err
	}

	return request, true, nil
}

// resolveCoordinates 解析坐标
//
// 坐标解析优先级：
// 1. 优先使用显式传入的经纬度（如果提供）
// 2. 如果没有经纬度，通过地址进行反向地理编码
//
// 返回值：纬度、经度、解析后的地址、错误
func (s *RequestService) resolveCoordinates(input RequestInput) (float64, float64, string, error) {
	// Priority 1: explicit lat/lng provided
	if input.Lat != nil || input.Lng != nil {
		if input.Lat == nil || input.Lng == nil {
			return 0, 0, "", ErrInvalidRequest
		}
		lat := *input.Lat
		lng := *input.Lng
		if !validCoordinate(lat, lng) {
			return 0, 0, "", ErrInvalidRequest
		}
		return lat, lng, input.Address, nil
	}

	// Priority 2: resolve by address
	if input.Address == "" {
		return 0, 0, "", ErrInvalidRequest
	}
	if s.geocodeSvc == nil {
		return 0, 0, "", errors.New("geocode service not configured")
	}
	geo, err := s.geocodeSvc.Geocode(input.Address)
	if err != nil {
		if errors.Is(err, ErrGeocodeNotFound) {
			return 0, 0, "", ErrInvalidRequest
		}
		return 0, 0, "", err
	}
	if !validCoordinate(geo.Latitude, geo.Longitude) {
		return 0, 0, "", ErrInvalidRequest
	}
	addr := input.Address
	if geo.FormattedAddress != "" {
		addr = geo.FormattedAddress
	}
	return geo.Latitude, geo.Longitude, addr, nil
}

func (s *RequestService) ListByUser(userID int64, status string, page, pageSize int) ([]*model.ServiceRequest, int64, error) {
	if userID == 0 {
		return nil, 0, ErrInvalidRequest
	}
	offset := (page - 1) * pageSize
	return s.repo.ListByUser(userID, status, offset, pageSize)
}

// UpdateInput B端编辑服务请求的输入参数
type UpdateInput struct {
	ServiceType  string `json:"service_type"`
	ContactName  string `json:"contact_name"`
	ContactPhone string `json:"contact_phone"`
	Address      string `json:"address"`
	Description  string `json:"description"`
	Urgency      string `json:"urgency"`
}

// Update B端编辑服务请求（仅 pending/dispatched 状态可编辑）
func (s *RequestService) Update(id int64, input UpdateInput) (*model.ServiceRequest, error) {
	req, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 仅 pending 和 dispatched 状态可编辑
	if req.Status != consts.RequestStatusPending && req.Status != consts.RequestStatusDispatched {
		return nil, ErrRequestConflict
	}

	updates := map[string]interface{}{}
	if input.ServiceType != "" {
		if !consts.IsValidServiceType(input.ServiceType) {
			return nil, ErrInvalidRequest
		}
		updates["service_type"] = input.ServiceType
	}
	if input.ContactName != "" {
		updates["contact_name"] = input.ContactName
	}
	if input.ContactPhone != "" {
		updates["contact_phone"] = input.ContactPhone
	}
	if input.Address != "" {
		updates["address"] = input.Address
	}
	if input.Description != "" {
		updates["description"] = input.Description
	}
	if input.Urgency != "" {
		updates["urgency"] = input.Urgency
	}

	if len(updates) == 0 {
		return req, nil
	}

	if err := s.db.Model(&model.ServiceRequest{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return nil, err
	}

	return s.repo.GetByID(id)
}

// CancelByAdmin B端取消服务请求
func (s *RequestService) CancelByAdmin(id int64) (*model.ServiceRequest, error) {
	req, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Status == consts.RequestStatusCancelled {
		return req, nil
	}
	if req.Status == consts.RequestStatusCompleted {
		return nil, ErrRequestConflict
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.ServiceRequest{}).
			Where("id = ?", req.ID).
			Update("status", consts.RequestStatusCancelled).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.TaskAssignment{}).
			Where("request_id = ?", req.ID).
			Update("status", consts.TaskStatusCancelled).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	req.Status = consts.RequestStatusCancelled
	return req, nil
}

// ListAll B端查询所有需求
func (s *RequestService) ListAll(stationID int64, status string, page, pageSize int) ([]*repository.RequestWithStation, int64, error) {
	offset := (page - 1) * pageSize
	return s.repo.ListAll(stationID, status, offset, pageSize)
}

// UpdateStatus B端更新状态
func (s *RequestService) UpdateStatus(id int64, status string, rejectReason string) error {
	if !consts.IsValidRequestStatus(status) {
		return ErrInvalidRequest
	}
	return s.repo.UpdateStatusByAdmin(id, status, rejectReason)
}

func (s *RequestService) GetByID(id int64) (*model.ServiceRequest, error) {
	if id == 0 {
		return nil, ErrInvalidRequest
	}
	return s.repo.GetByID(id)
}

func (s *RequestService) Cancel(id int64, userID int64) (*model.ServiceRequest, bool, error) {
	req, err := s.repo.GetByID(id)
	if err != nil {
		return nil, false, err
	}
	if req.UserID != userID {
		return nil, false, ErrRequestConflict
	}
	if req.Status == consts.RequestStatusCancelled {
		return req, false, nil
	}
	if req.Status == consts.RequestStatusCompleted {
		return nil, false, ErrRequestConflict
	}

	// 使用事务同步更新 ServiceRequest 和 TaskAssignment
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 更新 ServiceRequest 状态
		if err := tx.Model(&model.ServiceRequest{}).
			Where("id = ?", req.ID).
			Update("status", consts.RequestStatusCancelled).Error; err != nil {
			return err
		}
		// 同步取消关联的 TaskAssignment
		if err := tx.Model(&model.TaskAssignment{}).
			Where("request_id = ?", req.ID).
			Update("status", consts.TaskStatusCancelled).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, false, err
	}

	req.Status = consts.RequestStatusCancelled
	return req, true, nil
}

func (s *RequestService) Rate(id int64, userID int64, rating int, feedback string) (*model.ServiceRequest, error) {
	if rating < 1 || rating > 5 {
		return nil, ErrInvalidRequest
	}

	req, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check ownership
	if req.UserID != userID {
		return nil, ErrRequestConflict
	}

	// Check if request is completed
	if req.Status != consts.RequestStatusCompleted {
		return nil, ErrNotCompleted
	}

	// Check if already rated
	if req.Rating > 0 {
		return nil, ErrAlreadyRated
	}

	// Update rating
	if err := s.repo.UpdateRating(id, rating, feedback); err != nil {
		return nil, err
	}

	req.Rating = int64(rating)
	req.Feedback = feedback
	return req, nil
}

func (s *RequestService) findNearestStation(lat, lng float64) (int64, error) {
	stations, err := s.stationRepo.ListActive()
	if err != nil {
		return 0, err
	}
	if len(stations) == 0 {
		return 0, ErrNoStation
	}
	nearestID := stations[0].ID
	minDistance := HaversineDistance(lat, lng, stations[0].Latitude, stations[0].Longitude)
	for _, station := range stations[1:] {
		distance := HaversineDistance(lat, lng, station.Latitude, station.Longitude)
		if distance < minDistance {
			minDistance = distance
			nearestID = station.ID
		}
	}
	return nearestID, nil
}

func validCoordinate(lat, lng float64) bool {
	if lat < -90 || lat > 90 {
		return false
	}
	if lng < -180 || lng > 180 {
		return false
	}
	return true
}

// generateRequestNo 生成服务请求编号
func generateRequestNo() string {
	return fmt.Sprintf("REQ%d%04d", time.Now().Unix(), rand.Intn(10000))
}
