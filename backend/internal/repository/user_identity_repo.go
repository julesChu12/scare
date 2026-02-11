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

// GetByUserID 根据用户ID获取所有身份
func (r *UserIdentityRepository) GetByUserID(userID int64) ([]*model.UserIdentity, error) {
	ui := r.q.UserIdentity
	return ui.Where(ui.UserID.Eq(userID)).Find()
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

// HasBEndIdentity 检查用户是否有 B端身份
func (r *UserIdentityRepository) HasBEndIdentity(userID int64) (bool, error) {
	ui := r.q.UserIdentity
	count, err := ui.Where(
		ui.UserID.Eq(userID),
		ui.Status.Eq("active"),
		ui.IdentityType.In("admin", "station_manager", "staff"),
	).Count()
	return count > 0, err
}

// HasCEndIdentity 检查用户是否有 C端身份
func (r *UserIdentityRepository) HasCEndIdentity(userID int64) (bool, error) {
	ui := r.q.UserIdentity
	count, err := ui.Where(
		ui.UserID.Eq(userID),
		ui.Status.Eq("active"),
		ui.IdentityType.In("elderly", "family", "pregnant", "disabled", "child"),
	).Count()
	return count > 0, err
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

// Update 更新用户身份
func (r *UserIdentityRepository) Update(identity *model.UserIdentity) error {
	ui := r.q.UserIdentity
	_, err := ui.Where(ui.ID.Eq(identity.ID)).Updates(identity)
	return err
}

// UpdateStatus 更新身份状态
func (r *UserIdentityRepository) UpdateStatus(userID int64, identityType, status string) error {
	ui := r.q.UserIdentity
	_, err := ui.Where(
		ui.UserID.Eq(userID),
		ui.IdentityType.Eq(identityType),
	).Update(ui.Status, status)
	return err
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

// Delete 删除身份（软删除）
func (r *UserIdentityRepository) Delete(userID int64, identityType string) error {
	ui := r.q.UserIdentity
	_, err := ui.Where(
		ui.UserID.Eq(userID),
		ui.IdentityType.Eq(identityType),
	).Delete()
	return err
}

// WithTx 事务支持
func (r *UserIdentityRepository) WithTx(tx *gorm.DB) *UserIdentityRepository {
	return &UserIdentityRepository{q: query.Use(tx)}
}
