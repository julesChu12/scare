package repository

import (
	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/dao/query"

	"gorm.io/gorm"
)

// UserIdentityRepository 用户身份仓储
type UserIdentityRepository struct {
	q *query.Query
}

// NewUserIdentityRepository 创建用户身份仓储实例
func NewUserIdentityRepository(db *gorm.DB) *UserIdentityRepository {
	return &UserIdentityRepository{
		q: query.Use(db),
	}
}

// GetActiveByUserID 根据用户ID获取所有激活的身份
func (r *UserIdentityRepository) GetActiveByUserID(userID int64) ([]*model.UserIdentity, error) {
	ui := r.q.UserIdentity
	return ui.Where(ui.UserID.Eq(userID), ui.Status.Eq("active")).Find()
}

// GetPrimaryIdentity 获取用户的主身份
func (r *UserIdentityRepository) GetPrimaryIdentity(userID int64) (*model.UserIdentity, error) {
	ui := r.q.UserIdentity
	return ui.Where(
		ui.UserID.Eq(userID),
		ui.IsPrimary.Is(true),
		ui.Status.Eq("active"),
	).First()
}

// GetBEndIdentities 获取用户的 B端身份（admin/station_manager/staff）
func (r *UserIdentityRepository) GetBEndIdentities(userID int64) ([]*model.UserIdentity, error) {
	ui := r.q.UserIdentity
	return ui.Where(
		ui.UserID.Eq(userID),
		ui.Status.Eq("active"),
		ui.IdentityType.In("admin", "station_manager", "staff"),
	).Find()
}

// GetCEndIdentities 获取用户的 C端身份（elderly/family/pregnant/disabled/child）
func (r *UserIdentityRepository) GetCEndIdentities(userID int64) ([]*model.UserIdentity, error) {
	ui := r.q.UserIdentity
	return ui.Where(
		ui.UserID.Eq(userID),
		ui.Status.Eq("active"),
		ui.IdentityType.In("elderly", "family", "pregnant", "disabled", "child"),
	).Find()
}

// ExistsByUserIDAndType 检查用户是否有指定身份
func (r *UserIdentityRepository) ExistsByUserIDAndType(userID int64, identityType string) (bool, error) {
	ui := r.q.UserIdentity
	count, err := ui.Where(
		ui.UserID.Eq(userID),
		ui.IdentityType.Eq(identityType),
		ui.Status.Eq("active"),
	).Count()
	return count > 0, err
}

// Create 创建用户身份
func (r *UserIdentityRepository) Create(identity *model.UserIdentity) error {
	return r.q.UserIdentity.Omit(r.q.UserIdentity.RevokedAt).Create(identity)
}

// SetPrimary 设置主身份（会取消其他身份的主身份标记）
func (r *UserIdentityRepository) SetPrimary(userID int64, identityType string) error {
	ui := r.q.UserIdentity

	// 先取消所有主身份
	_, err := ui.Where(ui.UserID.Eq(userID)).Update(ui.IsPrimary, false)
	if err != nil {
		return err
	}

	// 设置新的主身份
	_, err = ui.Where(
		ui.UserID.Eq(userID),
		ui.IdentityType.Eq(identityType),
	).Update(ui.IsPrimary, true)
	return err
}

// GetByStationAndType 根据站点ID和身份类型获取用户身份列表
func (r *UserIdentityRepository) GetByStationAndType(stationID int64, identityType string) ([]*model.UserIdentity, error) {
	ui := r.q.UserIdentity
	return ui.Where(
		ui.StationID.Eq(stationID),
		ui.IdentityType.Eq(identityType),
		ui.Status.Eq("active"),
	).Find()
}

