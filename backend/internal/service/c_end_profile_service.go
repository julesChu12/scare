package service

import (
	"errors"

	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/repository"

	"gorm.io/gorm"
)

var ErrInvalidProfile = errors.New("invalid profile")

// CEndProfileService 负责 C 端资料更新，统一收口 user 与 customer_profile 的联动更新。
type CEndProfileService struct {
	db           *gorm.DB
	userRepo     *repository.UserRepository
	customerRepo *repository.CustomerRepository
}

// CEndProfileUpdateInput C 端资料更新输入。
type CEndProfileUpdateInput struct {
	Name     *string
	IDNumber *string
	Address  *string
	UserType *string
}

// NewCEndProfileService 创建 CEndProfileService
func NewCEndProfileService(db *gorm.DB, userRepo *repository.UserRepository, customerRepo *repository.CustomerRepository) *CEndProfileService {
	return &CEndProfileService{
		db:           db,
		userRepo:     userRepo,
		customerRepo: customerRepo,
	}
}

// Update 更新当前登录 C 端用户的资料。
//
// 业务规则：
// 1. name 写入 users 表，其余资料写入 customer_profiles 表
// 2. 两张表通过事务同时更新，避免部分成功
// 3. 用户或档案不存在时分别返回 ErrUserNotFound / ErrNoCustomerProfile
func (s *CEndProfileService) Update(userID int64, input CEndProfileUpdateInput) (*model.CustomerProfile, error) {
	if userID == 0 {
		return nil, ErrInvalidProfile
	}

	var updatedProfile *model.CustomerProfile
	err := s.db.Transaction(func(tx *gorm.DB) error {
		userRepoTx := repository.NewUserRepository(tx)
		customerRepoTx := repository.NewCustomerRepository(tx)

		user, err := userRepoTx.GetByID(userID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrUserNotFound
			}
			return err
		}

		profile, err := customerRepoTx.GetByUserID(userID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNoCustomerProfile
			}
			return err
		}

		if input.Name != nil {
			user.Name = *input.Name
			if err := userRepoTx.Update(user); err != nil {
				return err
			}
		}

		if input.IDNumber != nil {
			profile.IDCard = *input.IDNumber
		}
		if input.Address != nil {
			profile.Address = *input.Address
		}
		if input.UserType != nil {
			profile.CustomerType = *input.UserType
		}

		if err := customerRepoTx.Update(profile); err != nil {
			return err
		}
		updatedProfile = profile
		return nil
	})
	if err != nil {
		return nil, err
	}

	return updatedProfile, nil
}
