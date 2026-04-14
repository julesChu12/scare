package repository

import (
	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/dao/query"

	"gorm.io/gorm"
)

type UserRepository struct {
	q  *query.Query
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		q:  query.Use(db),
		db: db,
	}
}

// UserFilter 用户筛选条件
type UserFilter struct {
	Role      string // 按身份筛选（从 user_identities 表）
	Status    string // 按状态筛选
	StationID int64  // 按站点筛选
	Keyword   string // 关键词搜索（姓名/手机号）
}

// Create 创建用户
func (r *UserRepository) Create(user *model.User) error {
	return r.q.User.Omit(r.q.User.BirthDate).Create(user)
}

// CreateWithoutPassword 创建未设置密码的用户。
func (r *UserRepository) CreateWithoutPassword(user *model.User) error {
	return r.q.User.Omit(r.q.User.PasswordHash, r.q.User.BirthDate).Create(user)
}

// GetByID 根据ID获取用户
func (r *UserRepository) GetByID(id int64) (*model.User, error) {
	u := r.q.User
	return u.Where(u.ID.Eq(id)).First()
}

// GetByPhone 根据手机号获取用户
func (r *UserRepository) GetByPhone(phone string) (*model.User, error) {
	u := r.q.User
	return u.Where(u.Phone.Eq(phone)).First()
}

// ListWithFilter 获取用户列表（分页+筛选）
//
// 注意：由于筛选条件涉及子查询和复杂条件，此方法使用原生 GORM
// 以保持查询的灵活性和可读性
func (r *UserRepository) ListWithFilter(offset, limit int, filter UserFilter) ([]*model.User, int64, error) {
	db := r.db.Model(&model.User{})

	// 按身份筛选：通过子查询从 user_identities 表筛选
	if filter.Role != "" {
		db = db.Where("id IN (SELECT user_id FROM user_identities WHERE identity_type = ? AND status = 'active' AND deleted_at IS NULL)", filter.Role)
	}

	// 按状态筛选
	if filter.Status != "" {
		db = db.Where("status = ?", filter.Status)
	}

	// 按站点筛选
	if filter.StationID > 0 {
		db = db.Where("station_id = ?", filter.StationID)
	}

	// 关键词搜索（姓名或手机号）
	if filter.Keyword != "" {
		keyword := "%" + filter.Keyword + "%"
		db = db.Where("name LIKE ? OR phone LIKE ?", keyword, keyword)
	}

	// 排除软删除
	db = db.Where("deleted_at IS NULL")

	// 获取总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	var users []*model.User
	if err := db.Order("id DESC").Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// Update 更新用户
func (r *UserRepository) Update(user *model.User) error {
	u := r.q.User
	_, err := u.Where(u.ID.Eq(user.ID)).Updates(user)
	return err
}

// CountTodayNew 统计今日新增用户数量
//
// 注意：由于涉及日期函数，此方法使用原生 GORM
func (r *UserRepository) CountTodayNew() (int64, error) {
	start, end := todayRange()
	var count int64
	err := r.db.Model(&model.User{}).
		Where("created_at >= ? AND created_at < ?", start, end).
		Where("deleted_at IS NULL").
		Count(&count).Error
	return count, err
}
