package service

import (
	"community-elderly-care-platform/internal/dao/model"
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

// GetByID 根据ID获取角色
func (s *RoleService) GetByID(id int64) (*model.Role, error) {
	return s.roleRepo.GetByID(id)
}

// GetByCode 根据编码获取角色
func (s *RoleService) GetByCode(code string) (*model.Role, error) {
	return s.roleRepo.GetByCode(code)
}

// GetByCodes 根据编码列表获取角色
func (s *RoleService) GetByCodes(codes []string) ([]model.Role, error) {
	return s.roleRepo.GetByCodes(codes)
}

// List 获取所有角色
func (s *RoleService) List() ([]model.Role, error) {
	return s.roleRepo.List()
}

// Create 创建角色
func (s *RoleService) Create(role *model.Role) error {
	return s.roleRepo.Create(role)
}

// Update 更新角色
func (s *RoleService) Update(role *model.Role) error {
	return s.roleRepo.Update(role)
}

// Delete 删除角色
func (s *RoleService) Delete(id int64) error {
	return s.roleRepo.Delete(id)
}
