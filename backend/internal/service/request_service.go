package service

import (
	"context"
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
	db               *gorm.DB
	repo             *repository.RequestRepository
	taskRepo         *repository.TaskRepository
	stationRepo      *repository.StationRepository
	geofenceSvc      *GeofenceService
	geocodeSvc       *GeocodeService
	notifySvc        *NotificationService
	userIdentityRepo *repository.UserIdentityRepository
}

// RequestInput 创建服务请求的输入参数
type RequestInput struct {
	UserID          int64    `json:"user_id"`           // 用户ID
	RequestNo       string   `json:"request_no"`        // 请求编号（可选，用于幂等）
	ServiceType     string   `json:"service_type"`      // 服务类型
	SubmitLat       *float64 `json:"submit_lat"`        // 提交位置纬度（可选）
	SubmitLng       *float64 `json:"submit_lng"`        // 提交位置经度（可选）
	ServiceLat      *float64 `json:"service_lat"`       // 服务位置纬度（可选）
	ServiceLng      *float64 `json:"service_lng"`       // 服务位置经度（可选）
	SourceStationID *int64   `json:"source_station_id"` // 来源站点ID（可选）
	ContactName     string   `json:"contact_name"`      // 联系人姓名
	ContactPhone    string   `json:"contact_phone"`     // 联系人电话
	Address         string   `json:"address"`           // 地址
	Images          []string `json:"images"`            // 图片列表
}

func NewRequestService(db *gorm.DB, repo *repository.RequestRepository, taskRepo *repository.TaskRepository, stationRepo *repository.StationRepository, geofenceSvc *GeofenceService, geocodeSvc *GeocodeService, notifySvc *NotificationService, userIdentityRepo *repository.UserIdentityRepository) *RequestService {
	return &RequestService{
		db:               db,
		repo:             repo,
		taskRepo:         taskRepo,
		stationRepo:      stationRepo,
		geofenceSvc:      geofenceSvc,
		geocodeSvc:       geocodeSvc,
		notifySvc:        notifySvc,
		userIdentityRepo: userIdentityRepo,
	}
}

// Create 创建服务请求
//
// 业务目的：客户提交服务需求，系统自动匹配站点并创建任务
//
// 主要流程：
//  1. 验证输入参数（用户ID、服务类型必填）
//  2. 解析服务地址（通过地址地理编码得到服务坐标）
//  3. 幂等性检查：如果提供了 RequestNo，检查是否已存在
//  4. 生成请求编号（REQ + 时间戳 + 随机数）
//  5. 匹配服务站点：
//     a. 首先检查地理围栏（点是否在多边形内）
//     b. 如果没有匹配的围栏，查找最近的站点（Haversine 距离计算）
//     c. 如果服务地址无法解析，则创建待人工复核的需求
//  6. 在事务中创建服务请求；仅在已分配站点时同步创建任务
//
// 站点匹配只基于最终服务地址解析出的坐标进行，提交位置仅用于记录
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

	decision, err := resolveDispatch(DispatchInput{
		Address:          input.Address,
		SubmitLatitude:   input.SubmitLat,
		SubmitLongitude:  input.SubmitLng,
		ServiceLatitude:  input.ServiceLat,
		ServiceLongitude: input.ServiceLng,
		SourceStationID:  input.SourceStationID,
	}, s.stationRepo, s.geofenceSvc, s.geocodeSvc)
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

	images, err := json.Marshal(input.Images)
	if err != nil {
		return nil, false, err
	}

	request := &model.ServiceRequest{
		RequestNo:          requestNo,
		UserID:             input.UserID,
		ServiceType:        input.ServiceType,
		Status:             consts.RequestStatusPending,
		SubmitLocationLat:  decision.SubmitLatitude,
		SubmitLocationLng:  decision.SubmitLongitude,
		ServiceLocationLat: decision.ServiceLatitude,
		ServiceLocationLng: decision.ServiceLongitude,
		ContactName:        input.ContactName,
		ContactPhone:       input.ContactPhone,
		Address:            decision.ResolvedAddress,
		SourceStationID:    decision.SourceStationID,
		StationID:          decision.AssignedStationID,
		DispatchBasis:      decision.DispatchBasis,
		NeedsManualVerify:  decision.NeedsManualVerify,
		Images:             string(images),
	}
	if decision.AssignedStationID > 0 {
		request.Status = consts.RequestStatusDispatched
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		reqRepo := s.repo.WithTx(tx)
		if err := reqRepo.Create(request); err != nil {
			return err
		}
		if decision.AssignedStationID == 0 {
			return nil
		}
		taskRepo := s.taskRepo.WithTx(tx)
		task := &model.TaskAssignment{
			RequestID: request.ID,
			StationID: decision.AssignedStationID,
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

	if decision.AssignedStationID > 0 {
		// 异步通知站点管理员
		s.sendRequestNotification(decision.AssignedStationID, "新服务需求", fmt.Sprintf("收到新的%s服务需求，请及时处理。", consts.GetServiceTypeName(input.ServiceType)))
	}
	return request, true, nil
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
	StationID    *int64 `json:"station_id"`
}

// Update B端编辑服务请求（仅 pending/dispatched 状态可编辑）
//
// 说明：
// 1. 仅允许更新基础资料与人工纠偏字段，不在此处重跑派单逻辑
// 2. 空值字段保持原值，只有显式传入的非空字段才会落库
// 3. station_id 允许人工改派，但不会自动生成或重建任务
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
	if input.StationID != nil {
		updates["station_id"] = *input.StationID
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
//
// 业务规则：
// 1. 已完成需求不可取消，避免破坏已闭环数据
// 2. 已取消需求视为幂等请求，直接返回当前状态
// 3. 取消时同步取消关联任务，并通知站点管理员
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
	// 异步通知站点管理员
	s.sendRequestNotification(req.StationID, "服务需求已取消", "管理员已取消服务需求。")
	return req, nil
}

// ListAll B端查询所有需求
//
// 说明：
// - admin 可通过 stationID=0 查看全局数据
// - 非 admin 的站点范围约束由 Handler 层在进入 Service 前收口
func (s *RequestService) ListAll(stationID int64, status string, page, pageSize int) ([]*repository.RequestWithStation, int64, error) {
	offset := (page - 1) * pageSize
	return s.repo.ListAll(stationID, status, offset, pageSize)
}

// UpdateStatus B端更新状态
//
// 适用场景：后台人工审核、驳回或纠偏状态时使用。
// 这里只做状态值合法性校验，具体状态落库规则由 Repository 层统一处理。
func (s *RequestService) UpdateStatus(id int64, status string, rejectReason string) error {
	if !consts.IsValidRequestStatus(status) {
		return ErrInvalidRequest
	}
	return s.repo.UpdateStatusByAdmin(id, status, rejectReason)
}

// GetByID 获取单个服务需求详情。
func (s *RequestService) GetByID(id int64) (*model.ServiceRequest, error) {
	if id == 0 {
		return nil, ErrInvalidRequest
	}
	return s.repo.GetByID(id)
}

// Cancel C端取消自己的服务请求。
//
// 业务规则：
// 1. 仅需求所属用户可取消
// 2. 已完成需求不可取消
// 3. 已取消需求视为幂等请求
// 4. 取消时同步取消任务，并通知站点管理员
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
	// 异步通知站点管理员
	s.sendRequestNotification(req.StationID, "服务需求已取消", "用户已取消服务需求。")
	return req, true, nil
}

// Rate 对已完成的服务请求进行评价。
//
// 业务规则：
// 1. 仅需求所属用户可评价
// 2. 仅 completed 状态允许评价
// 3. 每个需求只允许评价一次
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
	if s.stationRepo == nil {
		return 0, ErrNoStation
	}

	stations, err := s.stationRepo.ListActive()
	if err != nil {
		return 0, err
	}

	nearest, err := nearestStationByHaversine(stations, lat, lng)
	if err != nil {
		return 0, err
	}

	return nearest.ID, nil
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

// sendRequestNotification 异步通知站点管理员
func (s *RequestService) sendRequestNotification(stationID int64, title, body string) {
	if s.notifySvc == nil || s.userIdentityRepo == nil || stationID == 0 {
		return
	}
	go func(sid int64, t, b string) {
		identities, err := s.userIdentityRepo.GetByStationAndType(sid, consts.IdentityStationManager)
		if err != nil {
			return
		}
		for _, identity := range identities {
			_ = s.notifySvc.SendInAppWithOptions(context.Background(), identity.UserID, t, b, NotificationSendOptions{
				Type: NotificationTypeTask,
			})
		}
	}(stationID, title, body)
}
