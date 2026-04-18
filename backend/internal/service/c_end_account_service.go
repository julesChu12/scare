package service

import (
	"errors"

	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/repository"

	"gorm.io/gorm"
)

// CEndAccountService 负责 C 端会话读取与预填充数据组装。
type CEndAccountService struct {
	userRepo     *repository.UserRepository
	customerRepo *repository.CustomerRepository
}

// CEndAccountInfo C 端“当前用户”接口使用的聚合视图。
type CEndAccountInfo struct {
	UserID       int64
	Type         string
	CustomerType string
	Name         string
	Phone        string
	Status       string
	HasPassword  bool
}

// CEndCheckPayload C 端 token 检查接口返回的聚合视图。
type CEndCheckPayload struct {
	User    CEndCheckUser
	Profile *CEndCheckProfile
}

// CEndCheckUser token 检查接口中的用户基本信息。
type CEndCheckUser struct {
	ID          int64
	Phone       string
	Role        string
	HasPassword bool
}

// CEndCheckProfile token 检查接口中的档案预填充信息。
type CEndCheckProfile struct {
	Name     string
	IDNumber string
	Address  string
	UserType string
}

// NewCEndAccountService 创建 C 端账户服务
func NewCEndAccountService(userRepo *repository.UserRepository, customerRepo *repository.CustomerRepository) *CEndAccountService {
	return &CEndAccountService{
		userRepo:     userRepo,
		customerRepo: customerRepo,
	}
}

// GetAccountInfo 获取 C 端“当前用户”接口需要的聚合信息。
//
// 业务规则：
// 1. 用户不存在时返回 ErrUserNotFound
// 2. 档案不存在时返回 ErrNoCustomerProfile
// 3. has_password 由 password_hash 是否为空统一推导
func (s *CEndAccountService) GetAccountInfo(userID int64) (*CEndAccountInfo, error) {
	user, err := s.loadUser(userID)
	if err != nil {
		return nil, err
	}

	profile, err := s.customerRepo.GetByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNoCustomerProfile
		}
		return nil, err
	}

	return &CEndAccountInfo{
		UserID:       user.ID,
		Type:         "c_end",
		CustomerType: profile.CustomerType,
		Name:         user.Name,
		Phone:        user.Phone,
		Status:       user.Status,
		HasPassword:  HasPasswordHash(user.PasswordHash),
	}, nil
}

// GetCheckPayload 获取 token 检查接口需要的聚合数据。
//
// 与 GetAccountInfo 的差异在于：当档案不存在时，仍返回 user 信息，profile 置空。
func (s *CEndAccountService) GetCheckPayload(userID int64) (*CEndCheckPayload, error) {
	user, err := s.loadUser(userID)
	if err != nil {
		return nil, err
	}

	payload := &CEndCheckPayload{
		User: CEndCheckUser{
			ID:          user.ID,
			Phone:       user.Phone,
			Role:        "c_end",
			HasPassword: HasPasswordHash(user.PasswordHash),
		},
	}

	profile, err := s.customerRepo.GetByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return payload, nil
		}
		return nil, err
	}

	payload.Profile = &CEndCheckProfile{
		Name:     user.Name,
		IDNumber: profile.IDCard,
		Address:  profile.Address,
		UserType: profile.CustomerType,
	}
	return payload, nil
}

func (s *CEndAccountService) loadUser(userID int64) (*model.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}
