package service

import (
	"errors"
	"time"

	"community-elderly-care-platform/internal/consts"
	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/repository"

	"gorm.io/gorm"
)

var (
	ErrIdentityNotFound     = errors.New("identity not found")
	ErrIdentityAlreadyExist = errors.New("identity already exists")
	ErrNoPrimaryIdentity    = errors.New("user has no primary identity")
	ErrCannotRemovePrimary  = errors.New("cannot remove primary identity")
)

// UserIdentityService 用户身份服务
type UserIdentityService struct {
	identityRepo     *repository.UserIdentityRepository
	blacklistService *TokenBlacklistService
	db               *gorm.DB
}

// NewUserIdentityService 创建用户身份服务实例
func NewUserIdentityService(identityRepo *repository.UserIdentityRepository, blacklistService *TokenBlacklistService, db *gorm.DB) *UserIdentityService {
	return &UserIdentityService{
		identityRepo:     identityRepo,
		blacklistService: blacklistService,
		db:               db,
	}
}

// GrantIdentity 授予用户指定身份
func (s *UserIdentityService) GrantIdentity(userID int64, identityType string, isPrimary bool, stationID int64, grantedBy int64) error {
	// 检查身份是否已存在
	exists, err := s.identityRepo.ExistsByUserIDAndType(userID, identityType)
	if err != nil {
		return err
	}
	if exists {
		return ErrIdentityAlreadyExist
	}

	identity := &model.UserIdentity{
		UserID:       userID,
		IdentityType: identityType,
		IsPrimary:    isPrimary,
		StationID:    stationID,
		Status:       "active",
		GrantedAt:    time.Now(),
		GrantedBy:    grantedBy,
	}

	// 如果是主身份，需要取消其他主身份
	if isPrimary {
		if err := s.identityRepo.SetPrimary(userID, ""); err != nil {
			return err
		}
	}

	return s.identityRepo.Create(identity)
}

// SyncIdentities 同步用户身份（替换为新列表）
func (s *UserIdentityService) SyncIdentities(userID int64, newIdentities []string, stationID int64, grantedBy int64) error {
	// 获取当前身份
	current, err := s.identityRepo.GetActiveByUserID(userID)
	if err != nil {
		return err
	}

	// 构建当前身份 map
	currentMap := make(map[string]bool)
	for _, id := range current {
		currentMap[id.IdentityType] = true
	}

	// 构建新身份 map
	newMap := make(map[string]bool)
	for _, id := range newIdentities {
		newMap[id] = true
	}

	// 计算需要移除的（当前有但新列表没有）
	for _, id := range current {
		if !newMap[id.IdentityType] {
			// 跳过主身份，或者先切换主身份
			if id.IsPrimary && len(newIdentities) > 0 {
				// 切换主身份到新列表的第一个
				_ = s.identityRepo.SetPrimary(userID, newIdentities[0])
			}
			if err := s.revokeIdentityDirect(userID, id.IdentityType); err != nil {
				return err
			}
		}
	}

	// 计算需要添加的（新列表有但当前没有）
	for i, identityType := range newIdentities {
		if !currentMap[identityType] {
			isPrimary := i == 0 && len(current) == 0 // 第一个身份设为主身份（仅当用户无身份时）
			sid := int64(0)
			if consts.IsBEndIdentity(identityType) {
				sid = stationID
			}
			if err := s.GrantIdentity(userID, identityType, isPrimary, sid, grantedBy); err != nil && err != ErrIdentityAlreadyExist {
				return err
			}
		}
	}

	return nil
}

// revokeIdentityDirect 直接撤销身份（不检查主身份）
func (s *UserIdentityService) revokeIdentityDirect(userID int64, identityType string) error {
	now := time.Now()
	return s.db.Model(&model.UserIdentity{}).
		Where("user_id = ? AND identity_type = ?", userID, identityType).
		Updates(map[string]any{
			"status":     "inactive",
			"revoked_at": now,
		}).Error
}

// GetIdentityTypes 从身份列表提取身份类型字符串数组
func GetIdentityTypes(identities []*model.UserIdentity) []string {
	types := make([]string, 0, len(identities))
	for _, identity := range identities {
		types = append(types, identity.IdentityType)
	}
	return types
}

// FilterBEndIdentities 从身份列表筛选 B端身份类型
func FilterBEndIdentities(identities []*model.UserIdentity) []string {
	types := make([]string, 0)
	for _, identity := range identities {
		if consts.IsBEndIdentity(identity.IdentityType) {
			types = append(types, identity.IdentityType)
		}
	}
	return types
}

// FilterCEndIdentities 从身份列表筛选 C端身份类型
func FilterCEndIdentities(identities []*model.UserIdentity) []string {
	types := make([]string, 0)
	for _, identity := range identities {
		if consts.IsCEndIdentity(identity.IdentityType) {
			types = append(types, identity.IdentityType)
		}
	}
	return types
}
