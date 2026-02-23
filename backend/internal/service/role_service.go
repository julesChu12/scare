package service

import (
	"community-elderly-care-platform/internal/repository"
)

// RoleService 角色业务逻辑
type RoleService struct {
	roleRepo *repository.RoleRepository
}

// NewRoleService 创建角色服务
func NewRoleService(roleRepo *repository.RoleRepository) *RoleService {
	return &RoleService{roleRepo: roleRepo}
}
