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

