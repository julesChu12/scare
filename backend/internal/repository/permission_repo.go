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

// GetByID 根据ID获取权限
func (r *PermissionRepository) GetByID(id int64) (*model.Permission, error) {
	var perm model.Permission
	if err := r.db.First(&perm, id).Error; err != nil {
		return nil, err
	}
	return &perm, nil
}

// GetByCode 根据权限码获取权限
func (r *PermissionRepository) GetByCode(code string) (*model.Permission, error) {
	var perm model.Permission
	if err := r.db.Where("code = ?", code).First(&perm).Error; err != nil {
		return nil, err
	}
	return &perm, nil
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

// GetByParentID 获取指定父级的子权限
func (r *PermissionRepository) GetByParentID(parentID int64) ([]model.Permission, error) {
	var perms []model.Permission
	if err := r.db.Where("parent_id = ? AND status = ?", parentID, "active").Order("sort ASC, id ASC").Find(&perms).Error; err != nil {
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

// GetByAPIPath 根据API路径和方法获取权限
func (r *PermissionRepository) GetByAPIPath(path, method string) (*model.Permission, error) {
	var perm model.Permission
	if err := r.db.Where("api_path = ? AND api_method = ? AND status = ?", path, method, "active").First(&perm).Error; err != nil {
		return nil, err
	}
	return &perm, nil
}

// GetAllWithAPI 获取所有有API路径的权限（用于权限检查）
func (r *PermissionRepository) GetAllWithAPI() ([]model.Permission, error) {
	var perms []model.Permission
	if err := r.db.Where("api_path != '' AND api_method != '' AND status = ?", "active").Find(&perms).Error; err != nil {
		return nil, err
	}
	return perms, nil
}

// Create 创建权限
func (r *PermissionRepository) Create(perm *model.Permission) error {
	return r.db.Create(perm).Error
}

// Update 更新权限
func (r *PermissionRepository) Update(perm *model.Permission) error {
	return r.db.Save(perm).Error
}

// Delete 删除权限（软删除）
func (r *PermissionRepository) Delete(id int64) error {
	return r.db.Delete(&model.Permission{}, id).Error
}
