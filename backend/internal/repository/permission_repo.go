package repository

import (
	"community-elderly-care-platform/internal/dao/model"

	"gorm.io/gorm"
)

// PermissionRepository 权限数据访问层
//
// 注意：Permission 模型未纳入 GORM Gen 管理，因此使用原生 GORM 查询
// 如需迁移到 Gen，需要先在 gorm_gen.go 中添加 Permission 表的生成配置
type PermissionRepository struct {
	db *gorm.DB
}

// NewPermissionRepository 创建权限仓库
func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

// GetByCodes 根据权限码列表获取权限
func (r *PermissionRepository) GetByCodes(codes []string) ([]model.Permission, error) {
	var perms []model.Permission
	if err := r.db.Where("code IN ? AND status = ?", codes, "active").Find(&perms).Error; err != nil {
		return nil, err
	}
	return perms, nil
}

// GetByIDs 根据ID列表获取权限
func (r *PermissionRepository) GetByIDs(ids []int64) ([]model.Permission, error) {
	var perms []model.Permission
	if err := r.db.Where("id IN ? AND status = ?", ids, "active").Find(&perms).Error; err != nil {
		return nil, err
	}
	return perms, nil
}

// List 获取所有权限
func (r *PermissionRepository) List() ([]model.Permission, error) {
	var perms []model.Permission
	if err := r.db.Where("status = ?", "active").Order("sort ASC, id ASC").Find(&perms).Error; err != nil {
		return nil, err
	}
	return perms, nil
}

// GetPublicPermissions 获取所有公共权限
func (r *PermissionRepository) GetPublicPermissions() ([]model.Permission, error) {
	var perms []model.Permission
	if err := r.db.Where("is_public = ? AND status = ?", true, "active").Find(&perms).Error; err != nil {
		return nil, err
	}
	return perms, nil
}

