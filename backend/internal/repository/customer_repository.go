package repository

import (
	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/dao/query"

	"gorm.io/gorm"
)

type CustomerRepository struct {
	q *query.Query
}

func NewCustomerRepository(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{
		q: query.Use(db),
	}
}

// GetByUserID 根据用户ID获取档案
func (r *CustomerRepository) GetByUserID(userID int64) (*model.CustomerProfile, error) {
	c := r.q.CustomerProfile
	return c.Where(c.UserID.Eq(userID)).First()
}

// GetByUserIDWithUser 获取档案并预加载用户信息
func (r *CustomerRepository) GetByUserIDWithUser(userID int64) (*model.CustomerProfile, error) {
	c := r.q.CustomerProfile
	// Gen 的 Preload 需要使用原生 DB 的方式
	db := c.UnderlyingDB().Preload("User")
	var profile model.CustomerProfile
	if err := db.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

// Create 创建档案
func (r *CustomerRepository) Create(profile *model.CustomerProfile) error {
	return r.q.CustomerProfile.Create(profile)
}

// Update 更新档案
func (r *CustomerRepository) Update(profile *model.CustomerProfile) error {
	c := r.q.CustomerProfile
	_, err := c.Where(c.UserID.Eq(profile.UserID)).Updates(profile)
	return err
}

// Delete 删除档案
func (r *CustomerRepository) Delete(userID int64) error {
	c := r.q.CustomerProfile
	_, err := c.Where(c.UserID.Eq(userID)).Delete()
	return err
}

// Exists 检查用户是否有档案
func (r *CustomerRepository) Exists(userID int64) (bool, error) {
	c := r.q.CustomerProfile
	count, err := c.Where(c.UserID.Eq(userID)).Count()
	return count > 0, err
}

// ListByType 根据客户类型查询
func (r *CustomerRepository) ListByType(customerType string) ([]*model.CustomerProfile, error) {
	c := r.q.CustomerProfile
	return c.Where(c.CustomerType.Eq(customerType)).Find()
}
