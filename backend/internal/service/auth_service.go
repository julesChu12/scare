package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"community-elderly-care-platform/internal/consts"
	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/repository"
	"community-elderly-care-platform/pkg/jwt"

	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials      = errors.New("invalid credentials")
	ErrUserInactive            = errors.New("user inactive")
	ErrUserNotFound            = errors.New("user not found")
	ErrPasswordNotSet          = errors.New("password not set")
	ErrCurrentPasswordRequired = errors.New("current password required")
	ErrCurrentPasswordInvalid  = errors.New("current password invalid")
	ErrNoRoleForBEnd           = errors.New("user has no B-end identity")
	ErrNoCustomerProfile       = errors.New("user has no customer profile")
	ErrNoStation               = errors.New("no station available")
	ErrNoCEndIdentity          = errors.New("user has no C-end identity")
)

type AuthService struct {
	userRepo     *repository.UserRepository
	identityRepo *repository.UserIdentityRepository
	customerRepo *repository.CustomerRepository
	stationRepo  *repository.StationRepository
	jwtManager   *jwt.Manager
	smsService   *SMSService
	geofenceSvc  *GeofenceService
	geocodeSvc   *GeocodeService
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

// SetGeofenceService 设置地理围栏服务（可选依赖，用于登录时定位最近站点）
func (s *AuthService) SetGeofenceService(geofenceSvc *GeofenceService) {
	s.geofenceSvc = geofenceSvc
}

// SetGeocodeService 设置地理编码服务（可选依赖，用于解析地址坐标）
func (s *AuthService) SetGeocodeService(geocodeSvc *GeocodeService) {
	s.geocodeSvc = geocodeSvc
}

// SetStationRepo 设置站点仓库（可选依赖，用于登录时定位最近站点）
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
	if !hasPasswordHash(user.PasswordHash) {
		return nil, nil, ErrPasswordNotSet
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

// Refresh 使用刷新令牌获取新的访问令牌和刷新令牌。
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

// SetCEndPassword 为 C 端用户设置或更新密码。
func (s *AuthService) SetCEndPassword(userID int64, currentPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}
	if user.Status != "active" {
		return ErrUserInactive
	}
	if hasPasswordHash(user.PasswordHash) {
		if strings.TrimSpace(currentPassword) == "" {
			return ErrCurrentPasswordRequired
		}
		if err := VerifyPassword(user.PasswordHash, currentPassword); err != nil {
			return ErrCurrentPasswordInvalid
		}
	}

	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = hashedPassword
	return s.userRepo.Update(user)
}

// ResetCEndPassword 通过手机号验证码重置 C 端用户密码。
func (s *AuthService) ResetCEndPassword(phone, code, newPassword string) error {
	if err := s.smsService.VerifyCode(phone, code); err != nil {
		return err
	}

	user, err := s.userRepo.GetByPhone(phone)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}
	if user.Status != "active" {
		return ErrUserInactive
	}

	exists, err := s.customerRepo.Exists(user.ID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrNoCustomerProfile
	}

	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = hashedPassword
	return s.userRepo.Update(user)
}

// RegisterInput 注册输入参数
type RegisterInput struct {
	Phone    string
	Code     string
	Password string
	Name     string
}

// RegisterResult 注册结果
type RegisterResult struct {
	Token        string
	RefreshToken string
	User         *model.User
	Profile      *model.CustomerProfile
}

// Register 用户注册（创建 User + CustomerProfile + ElderlyIdentity）
func (s *AuthService) Register(input RegisterInput) (*RegisterResult, error) {
	if err := s.smsService.VerifyCode(input.Phone, input.Code); err != nil {
		return nil, err
	}

	// 检查用户是否已存在
	existingUser, err := s.userRepo.GetByPhone(input.Phone)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("该手机号已注册，请直接登录")
	}

	// 密码哈希
	hashedPassword, err := HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	var user *model.User
	var profile *model.CustomerProfile

	err = s.db.Transaction(func(tx *gorm.DB) error {
		userRepoTx := repository.NewUserRepository(tx)
		customerRepoTx := repository.NewCustomerRepository(tx)
		identityRepoTx := repository.NewUserIdentityRepository(tx)

		// 创建用户
		newUser := &model.User{
			Phone:        input.Phone,
			Name:         input.Name,
			PasswordHash: hashedPassword,
			Status:       "active",
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

		// 创建空的客户档案
		newProfile := &model.CustomerProfile{
			UserID:           user.ID,
			EmergencyContact: `{}`,
		}
		if err := customerRepoTx.Create(newProfile); err != nil {
			return err
		}
		profile = newProfile

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 生成 Token
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

	return &RegisterResult{
		Token:        access,
		RefreshToken: refresh,
		User:         user,
		Profile:      profile,
	}, nil
}

// QuickStartInput 快速开通的输入参数
type QuickStartInput struct {
	Phone            string
	Code             string
	Name             string
	Address          string
	SubmitLatitude   *float64
	SubmitLongitude  *float64
	ServiceLatitude  *float64
	ServiceLongitude *float64
	SourceStationID  *int64
	ServiceType      string
	Description      *string
	Images           []string // 图片 URL 列表
	ContactName      string   // 联系人姓名（用于服务请求）
	ContactPhone     string   // 联系人电话（用于服务请求）
}

// QuickStartResult 快速开通的返回结果
type QuickStartResult struct {
	Token        string
	RefreshToken string
	User         *model.User
	Profile      *model.CustomerProfile
	Request      *model.ServiceRequest
}

// QuickStart 快速开通服务（注册 + 建档 + 自动派单一体化）
//
// 适用场景：C 端老年用户通过手机号验证码即可快速开通服务，无需手动注册账号。
// 整个流程在一个数据库事务中完成，保证数据一致性。
//
// 执行流程：
//   1. 短信验证码校验
//   2. 调用 resolveDispatch() 做派单决策（核心：射线法围栏匹配）
//   3. 数据库事务：
//      a. 新用户 → 创建 User + ElderlyIdentity + CustomerProfile
//      b. 老用户 → 更新姓名/地址
//      c. 创建 ServiceRequest（记录派单依据 dispatch_basis）
//      d. 若自动派单成功 → 同时创建 TaskAssignment
//   4. 生成 JWT Token 返回
//
// 与射线法的关系：
//   resolveDispatch() 内部调用 resolveAssignedStation() → GeofenceService.Match()
//   → PointInPolygon()（射线法），射线法结果写入 request.DispatchBasis：
//     - "service_geofence"  → 射线法命中，自动派单
//     - "service_nearest"   → Haversine 兜底派单
//     - "service_address_manual_review" → 需人工审核派单
//
// 幂等性：
//   手机号已存在的用户 → 更新姓名和地址，不重复创建（复用已有账号）
//
func (s *AuthService) QuickStart(input QuickStartInput) (*QuickStartResult, error) {
	if err := s.smsService.VerifyCode(input.Phone, input.Code); err != nil {
		return nil, err
	}

	existingUser, err := s.userRepo.GetByPhone(input.Phone)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existingUser != nil && existingUser.Status != "active" {
		return nil, ErrUserInactive
	}

	decision, err := resolveDispatch(DispatchInput{
		Address:          input.Address,
		SubmitLatitude:   input.SubmitLatitude,
		SubmitLongitude:  input.SubmitLongitude,
		ServiceLatitude:  input.ServiceLatitude,
		ServiceLongitude: input.ServiceLongitude,
		SourceStationID:  input.SourceStationID,
	}, s.stationRepo, s.geofenceSvc, s.geocodeSvc)
	if err != nil {
		return nil, err
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
			if err := userRepoTx.CreateWithoutPassword(newUser); err != nil {
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
				UserID:           user.ID,
				Address:          decision.ResolvedAddress,
				CustomerType:     consts.IdentityElderly,
				EmergencyContact: `{}`, // JSON 列不能为空字符串，用空对象替代
			}
			if err := customerRepoTx.Create(newProfile); err != nil {
				return err
			}
			profile = newProfile
		} else {
			if decision.ResolvedAddress != "" {
				existingProfile.Address = decision.ResolvedAddress
			}
			if existingProfile.CustomerType == "" {
				existingProfile.CustomerType = consts.IdentityElderly
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
			RequestNo:          requestNo,
			UserID:             user.ID,
			ServiceType:        input.ServiceType,
			Status:             consts.RequestStatusPending,
			SubmitLocationLat:  decision.SubmitLatitude,
			SubmitLocationLng:  decision.SubmitLongitude,
			ServiceLocationLat: decision.ServiceLatitude,
			ServiceLocationLng: decision.ServiceLongitude,
			Address:            decision.ResolvedAddress,
			SourceStationID:    decision.SourceStationID,
			StationID:          decision.AssignedStationID,
			DispatchBasis:      decision.DispatchBasis,
			NeedsManualVerify:  decision.NeedsManualVerify,
			ContactName:        input.ContactName,
			ContactPhone:       input.ContactPhone,
		}
		if decision.AssignedStationID > 0 {
			newRequest.Status = consts.RequestStatusDispatched
		}

		if input.Description != nil {
			newRequest.Description = *input.Description
		}

		if len(input.Images) > 0 {
			imagesJSON, _ := json.Marshal(input.Images)
			newRequest.Images = string(imagesJSON)
		} else {
			newRequest.Images = "[]" // JSON 列不能为空字符串，用空数组替代
		}

		if err := requestRepo.Create(newRequest); err != nil {
			return err
		}

		if decision.AssignedStationID > 0 {
			task := &model.TaskAssignment{
				RequestID: newRequest.ID,
				StationID: decision.AssignedStationID,
				Status:    consts.TaskStatusDispatched,
			}
			if err := taskRepo.Create(task); err != nil {
				return err
			}
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

func hasPasswordHash(hash string) bool {
	return strings.TrimSpace(hash) != ""
}

// HasPasswordHash reports whether the user currently has password login capability.
func HasPasswordHash(hash string) bool {
	return hasPasswordHash(hash)
}

func (s *AuthService) findNearestStation(lat, lng float64) (int64, error) {
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
