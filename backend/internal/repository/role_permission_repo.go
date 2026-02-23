package repository

import (
	"community-elderly-care-platform/internal/dao/model"

	"gorm.io/gorm"
)

// RolePermissionRepository 角色权限关联数据访问层
//
// 注意：RolePermission 模型未纳入 GORM Gen 管理，因此使用原生 GORM 查询
// 如需迁移到 Gen，需要先在 gorm_gen.go 中添加 RolePermission 表的生成配置
type RolePermissionRepository struct {
	db *gorm.DB
}

// NewRolePermissionRepository 创建角色权限关联仓库
func NewRolePermissionRepository(db *gorm.DB) *RolePermissionRepository {
	return &RolePermissionRepository{db: db}
}

// GetPermissionIDsByRoleID 获取角色的所有权限ID
func (r *RolePermissionRepository) GetPermissionIDsByRoleID(roleID int64) ([]int64, error) {
	var ids []int64
	if err := r.db.Model(&model.RolePermission{}).Where("role_id = ?", roleID).Pluck("permission_id", &ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

// GetPermissionIDsByRoleIDs 获取多个角色的所有权限ID（去重）
func (r *RolePermissionRepository) GetPermissionIDsByRoleIDs(roleIDs []int64) ([]int64, error) {
	var ids []int64
	if err := r.db.Model(&model.RolePermission{}).Where("role_id IN ?", roleIDs).Distinct().Pluck("permission_id", &ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

// ReplaceRolePermissions 替换角色的所有权限（事务）
func (r *RolePermissionRepository) ReplaceRolePermissions(roleID int64, permissionIDs []int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 删除旧的关联
		if err := tx.Where("role_id = ?", roleID).Delete(&model.RolePermission{}).Error; err != nil {
			return err
		}

		// 创建新的关联
		if len(permissionIDs) > 0 {
			rps := make([]model.RolePermission, len(permissionIDs))
			for i, permID := range permissionIDs {
				rps[i] = model.RolePermission{
					RoleID:       roleID,
					PermissionID: permID,
				}
			}
			if err := tx.Create(&rps).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
