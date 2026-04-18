package service

import (
	"community-elderly-care-platform/internal/repository"
)

// RoleService 角色业务逻辑
type RoleService struct {
	roleRepo *repository.RoleRepository
}

// NewRoleService 创建角色服务
// RoleService 负责角色相关的业务逻辑，目前仅提供构造函数，实际操作委托给 RoleRepository
func NewRoleService(roleRepo *repository.RoleRepository) *RoleService {
	return &RoleService{roleRepo: roleRepo}
}
