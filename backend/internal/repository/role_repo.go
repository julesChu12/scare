package repository

import (
	"community-elderly-care-platform/internal/dao/model"

	"gorm.io/gorm"
)

// RoleRepository 角色数据访问层
//
// 注意：Role 模型未纳入 GORM Gen 管理，因此使用原生 GORM 查询
// 如需迁移到 Gen，需要先在 gorm_gen.go 中添加 Role 表的生成配置
type RoleRepository struct {
	db *gorm.DB
}

// NewRoleRepository 创建角色仓库
func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

// GetByID 根据ID获取角色
func (r *RoleRepository) GetByID(id int64) (*model.Role, error) {
	var role model.Role
	if err := r.db.First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// GetByCode 根据编码获取角色
func (r *RoleRepository) GetByCode(code string) (*model.Role, error) {
	var role model.Role
	if err := r.db.Where("code = ?", code).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// GetByCodes 根据编码列表获取角色
func (r *RoleRepository) GetByCodes(codes []string) ([]model.Role, error) {
	var roles []model.Role
	if err := r.db.Where("code IN ? AND status = ?", codes, "active").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// List 获取所有角色
func (r *RoleRepository) List() ([]model.Role, error) {
	var roles []model.Role
	if err := r.db.Where("status = ?", "active").Order("sort ASC, id ASC").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// Create 创建角色
func (r *RoleRepository) Create(role *model.Role) error {
	return r.db.Create(role).Error
}

// Update 更新角色
func (r *RoleRepository) Update(role *model.Role) error {
	return r.db.Save(role).Error
}

// Delete 删除角色（软删除）
func (r *RoleRepository) Delete(id int64) error {
	return r.db.Delete(&model.Role{}, id).Error
}
