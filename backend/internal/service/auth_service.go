package service

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"community-elderly-care-platform/internal/consts"
	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/repository"
	"community-elderly-care-platform/pkg/jwt"

	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserInactive       = errors.New("user inactive")
	ErrUserNotFound       = errors.New("user not found")
	ErrNoRoleForBEnd      = errors.New("user has no B-end identity")
	ErrNoCustomerProfile  = errors.New("user has no customer profile")
	ErrNoStation          = errors.New("no station available")
	ErrNoCEndIdentity     = errors.New("user has no C-end identity")
)

type AuthService struct {
	userRepo     *repository.UserRepository
	identityRepo *repository.UserIdentityRepository
	customerRepo *repository.CustomerRepository
	stationRepo  *repository.StationRepository
	jwtManager   *jwt.Manager
	smsService   *SMSService
	geofenceSvc  *GeofenceService
	db           *gorm.DB
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewAuthService(userRepo *repository.UserRepository, identityRepo *repository.UserIdentityRepository, customerRepo *repository.CustomerRepository, jwtManager *jwt.Manager, smsService *SMSService, db *gorm.DB) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		identityRepo: identityRepo,
		customerRepo: customerRepo,
		jwtManager:   jwtManager,
		smsService:   smsService,
		db:           db,
	}
}

// SetGeofenceService sets the geofence service (optional dependency)
func (s *AuthService) SetGeofenceService(geofenceSvc *GeofenceService) {
	s.geofenceSvc = geofenceSvc
}

// SetStationRepo sets the station repository (optional dependency)
func (s *AuthService) SetStationRepo(stationRepo *repository.StationRepository) {
	s.stationRepo = stationRepo
}

// LoginBEnd B端用户登录
func (s *AuthService) LoginBEnd(phone, password string) (*Tokens, *model.User, error) {
	// 获取用户基本信息
	user, err := s.userRepo.GetByPhone(phone)
	if err != nil {
		return nil, nil, ErrInvalidCredentials
	}
	if user.Status != "active" {
		return nil, nil, ErrUserInactive
	}
	if err := VerifyPassword(user.PasswordHash, password); err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	// 获取用户的 B端身份
	bEndIdentities, err := s.identityRepo.GetBEndIdentities(user.ID)
	if err != nil {
		return nil, nil, err
	}
	if len(bEndIdentities) == 0 {
		return nil, nil, ErrNoRoleForBEnd
	}

	// 提取身份类型
	identityTypes := GetIdentityTypes(bEndIdentities)

	// 获取主身份（如果主身份是 B端身份）
	var primary string
	primaryIdentity, err := s.identityRepo.GetPrimaryIdentity(user.ID)
	if err == nil && consts.IsBEndIdentity(primaryIdentity.IdentityType) {
		primary = primaryIdentity.IdentityType
	} else {
		// 使用第一个 B端身份作为主身份
		primary = identityTypes[0]
	}

	stationID := user.StationID

	// 生成 B端 JWT Token
	access, err := s.jwtManager.GenerateToken(user.ID, "b_end", stationID, identityTypes, primary)
	if err != nil {
		return nil, nil, err
	}
	refresh, err := s.jwtManager.GenerateRefreshToken(user.ID, "b_end", stationID, identityTypes, primary)
	if err != nil {
		return nil, nil, err
	}

	return &Tokens{AccessToken: access, RefreshToken: refresh}, user, nil
}

// LoginCEnd C端用户登录（密码）
func (s *AuthService) LoginCEnd(phone, password string) (*Tokens, *model.User, error) {
	user, err := s.userRepo.GetByPhone(phone)
	if err != nil {
		return nil, nil, ErrInvalidCredentials
	}
	if user.Status != "active" {
		return nil, nil, ErrUserInactive
	}
	if err := VerifyPassword(user.PasswordHash, password); err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	// 检查是否有客户档案
	exists, err := s.customerRepo.Exists(user.ID)
	if err != nil {
		return nil, nil, err
	}
	if !exists {
		return nil, nil, ErrNoCustomerProfile
	}

	// 获取 C端身份
	cEndIdentities, err := s.identityRepo.GetCEndIdentities(user.ID)
	if err != nil {
		return nil, nil, err
	}

	var identityTypes []string
	var primary string
	if len(cEndIdentities) > 0 {
		identityTypes = GetIdentityTypes(cEndIdentities)
		// 获取主身份
		primaryIdentity, err := s.identityRepo.GetPrimaryIdentity(user.ID)
		if err == nil && consts.IsCEndIdentity(primaryIdentity.IdentityType) {
			primary = primaryIdentity.IdentityType
		} else {
			primary = identityTypes[0]
		}
	}

	// 生成 C端 JWT Token
	access, err := s.jwtManager.GenerateToken(user.ID, "c_end", 0, identityTypes, primary)
	if err != nil {
		return nil, nil, err
	}
	refresh, err := s.jwtManager.GenerateRefreshToken(user.ID, "c_end", 0, identityTypes, primary)
	if err != nil {
		return nil, nil, err
	}

	return &Tokens{AccessToken: access, RefreshToken: refresh}, user, nil
}

// LoginCEndByCode C端用户验证码登录
func (s *AuthService) LoginCEndByCode(phone, code string) (*Tokens, *model.User, error) {
	if err := s.smsService.VerifyCode(phone, code); err != nil {
		return nil, nil, err
	}

	user, err := s.userRepo.GetByPhone(phone)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrUserNotFound
		}
		return nil, nil, err
	}
	if user.Status != "active" {
		return nil, nil, ErrUserInactive
	}

	// 检查是否有客户档案
	exists, err := s.customerRepo.Exists(user.ID)
	if err != nil {
		return nil, nil, err
	}
	if !exists {
		return nil, nil, ErrNoCustomerProfile
	}

	// 获取 C端身份
	cEndIdentities, err := s.identityRepo.GetCEndIdentities(user.ID)
	if err != nil {
		return nil, nil, err
	}

	var identityTypes []string
	var primary string
	if len(cEndIdentities) > 0 {
		identityTypes = GetIdentityTypes(cEndIdentities)
		primaryIdentity, err := s.identityRepo.GetPrimaryIdentity(user.ID)
		if err == nil && consts.IsCEndIdentity(primaryIdentity.IdentityType) {
			primary = primaryIdentity.IdentityType
		} else {
			primary = identityTypes[0]
		}
	}

	access, err := s.jwtManager.GenerateToken(user.ID, "c_end", 0, identityTypes, primary)
	if err != nil {
		return nil, nil, err
	}
	refresh, err := s.jwtManager.GenerateRefreshToken(user.ID, "c_end", 0, identityTypes, primary)
	if err != nil {
		return nil, nil, err
	}

	return &Tokens{AccessToken: access, RefreshToken: refresh}, user, nil
}

func (s *AuthService) Refresh(refreshToken string) (*Tokens, error) {
	claims, err := s.jwtManager.ParseToken(refreshToken)
	if err != nil {
		return nil, err
	}

	access, err := s.jwtManager.GenerateToken(claims.UserID, claims.Type, claims.StationID, claims.Identities, claims.Primary)
	if err != nil {
		return nil, err
	}
	refresh, err := s.jwtManager.GenerateRefreshToken(claims.UserID, claims.Type, claims.StationID, claims.Identities, claims.Primary)
	if err != nil {
		return nil, err
	}

	return &Tokens{AccessToken: access, RefreshToken: refresh}, nil
}

// QuickStartInput 快速开通的输入参数
type QuickStartInput struct {
	Phone       string
	Code        string
	Name        string
	Address     string
	Latitude    *float64
	Longitude   *float64
	ServiceType string
	Description *string
}

// QuickStartResult 快速开通的返回结果
type QuickStartResult struct {
	Token        string
	RefreshToken string
	User         *model.User
	Profile      *model.CustomerProfile
	Request      *model.ServiceRequest
}

// QuickStart 快速开通服务
func (s *AuthService) QuickStart(input QuickStartInput) (*QuickStartResult, error) {
	if err := s.smsService.VerifyCode(input.Phone, input.Code); err != nil {
		return nil, err
	}

	existingUser, err := s.userRepo.GetByPhone(input.Phone)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	lat := 0.0
	lng := 0.0
	if input.Latitude != nil {
		lat = *input.Latitude
	}
	if input.Longitude != nil {
		lng = *input.Longitude
	}

	var stationID int64
	if lat != 0 && lng != 0 && s.geofenceSvc != nil {
		if matchedID, matched := s.geofenceSvc.Match(lat, lng); matched {
			stationID = matchedID
		} else if s.stationRepo != nil {
			nearestID, err := s.findNearestStation(lat, lng)
			if err != nil {
				return nil, err
			}
			stationID = nearestID
		}
	} else if s.stationRepo != nil {
		stations, err := s.stationRepo.ListActive()
		if err != nil {
			return nil, err
		}
		if len(stations) == 0 {
			return nil, ErrNoStation
		}
		stationID = stations[0].ID
	}

	if stationID == 0 {
		return nil, ErrNoStation
	}

	var user *model.User
	var profile *model.CustomerProfile
	var request *model.ServiceRequest

	err = s.db.Transaction(func(tx *gorm.DB) error {
		userRepoTx := repository.NewUserRepository(tx)
		customerRepoTx := repository.NewCustomerRepository(tx)
		identityRepoTx := repository.NewUserIdentityRepository(tx)

		if existingUser == nil {
			newUser := &model.User{
				Phone:  input.Phone,
				Name:   input.Name,
				Status: "active",
			}
			if err := userRepoTx.Create(newUser); err != nil {
				return err
			}
			user = newUser

			// 创建 elderly 身份
			identity := &model.UserIdentity{
				UserID:       user.ID,
				IdentityType: consts.IdentityElderly,
				IsPrimary:    true,
				Status:       "active",
				GrantedAt:    time.Now(),
			}
			if err := identityRepoTx.Create(identity); err != nil {
				return err
			}
		} else {
			user = existingUser
			if input.Name != "" && user.Name != input.Name {
				user.Name = input.Name
				if err := userRepoTx.Update(user); err != nil {
					return err
				}
			}
		}

		existingProfile, err := customerRepoTx.GetByUserID(user.ID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if existingProfile == nil {
			newProfile := &model.CustomerProfile{
				UserID:  user.ID,
				Address: input.Address,
			}
			if err := customerRepoTx.Create(newProfile); err != nil {
				return err
			}
			profile = newProfile
		} else {
			if input.Address != "" {
				existingProfile.Address = input.Address
			}
			if err := customerRepoTx.Update(existingProfile); err != nil {
				return err
			}
			profile = existingProfile
		}

		requestRepo := repository.NewRequestRepository(tx)
		taskRepo := repository.NewTaskRepository(tx)

		requestNo := fmt.Sprintf("REQ%d%04d", time.Now().Unix(), rand.Intn(10000))

		newRequest := &model.ServiceRequest{
			RequestNo:         requestNo,
			UserID:            user.ID,
			ServiceType:       input.ServiceType,
			Status:            consts.RequestStatusDispatched,
			SubmitLocationLat: lat,
			SubmitLocationLng: lng,
			Address:           input.Address,
			StationID:         stationID,
		}

		if input.Description != nil {
			newRequest.Description = *input.Description
		}

		if err := requestRepo.Create(newRequest); err != nil {
			return err
		}

		task := &model.TaskAssignment{
			RequestID: newRequest.ID,
			StationID: stationID,
			Status:    consts.TaskStatusDispatched,
		}
		if err := taskRepo.Create(task); err != nil {
			return err
		}

		request = newRequest
		return nil
	})

	if err != nil {
		return nil, err
	}

	// 获取身份信息
	identities, _ := s.identityRepo.GetCEndIdentities(user.ID)
	var identityTypes []string
	var primary string
	if len(identities) > 0 {
		identityTypes = GetIdentityTypes(identities)
		primary = identityTypes[0]
	}

	access, err := s.jwtManager.GenerateToken(user.ID, "c_end", 0, identityTypes, primary)
	if err != nil {
		return nil, err
	}
	refresh, err := s.jwtManager.GenerateRefreshToken(user.ID, "c_end", 0, identityTypes, primary)
	if err != nil {
		return nil, err
	}

	return &QuickStartResult{
		Token:        access,
		RefreshToken: refresh,
		User:         user,
		Profile:      profile,
		Request:      request,
	}, nil
}

func (s *AuthService) findNearestStation(lat, lng float64) (int64, error) {
	if s.stationRepo == nil {
		return 0, ErrNoStation
	}
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
