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

// Create 创建档案
func (r *CustomerRepository) Create(profile *model.CustomerProfile) error {
	c := r.q.CustomerProfile
	if profile.BirthDate.IsZero() {
		return c.Omit(c.BirthDate).Create(profile)
	}
	return c.Create(profile)
}

// Update 更新档案
func (r *CustomerRepository) Update(profile *model.CustomerProfile) error {
	c := r.q.CustomerProfile
	_, err := c.Where(c.UserID.Eq(profile.UserID)).Updates(profile)
	return err
}

// Exists 检查用户是否有档案
func (r *CustomerRepository) Exists(userID int64) (bool, error) {
	c := r.q.CustomerProfile
	count, err := c.Where(c.UserID.Eq(userID)).Count()
	return count > 0, err
}
