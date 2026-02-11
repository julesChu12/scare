package service

import (
	"context"
	"strings"
	"sync"

	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/repository"

	"gorm.io/gorm"
)

// PermissionNode 权限树节点（用于前端展示）
type PermissionNode struct {
	ID       string           `json:"id"`
	Code     string           `json:"code"`
	Label    string           `json:"label"`
	Type     string           `json:"type"` // menu/button/resource
	APIPath  string           `json:"api_path,omitempty"`
	Method   string           `json:"method,omitempty"`
	Children []PermissionNode `json:"children,omitempty"`
	Disabled bool             `json:"disabled,omitempty"`
	IsPublic bool             `json:"is_public,omitempty"`
}

// PermissionService 权限业务逻辑
type PermissionService struct {
	db                 *gorm.DB
	permRepo           *repository.PermissionRepository
	roleRepo           *repository.RoleRepository
	rolePermRepo       *repository.RolePermissionRepository
	blacklistService   *TokenBlacklistService

	// 缓存
	cacheMu            sync.RWMutex
	permissionsCache   []model.Permission
	publicAPIsCache    []model.Permission
}

// NewPermissionService 创建权限服务
func NewPermissionService(
	db *gorm.DB,
	permRepo *repository.PermissionRepository,
	roleRepo *repository.RoleRepository,
	rolePermRepo *repository.RolePermissionRepository,
	blacklistService *TokenBlacklistService,
) *PermissionService {
	s := &PermissionService{
		db:               db,
		permRepo:         permRepo,
		roleRepo:         roleRepo,
		rolePermRepo:     rolePermRepo,
		blacklistService: blacklistService,
	}
	// 初始化缓存
	s.RefreshCache()
	return s
}

// RefreshCache 刷新权限缓存
func (s *PermissionService) RefreshCache() error {
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()

	// 加载所有权限
	perms, err := s.permRepo.List()
	if err != nil {
		return err
	}
	s.permissionsCache = perms

	// 加载公共API权限
	publicPerms, err := s.permRepo.GetPublicPermissions()
	if err != nil {
		return err
	}
	s.publicAPIsCache = publicPerms

	return nil
}

// GetPermissionTree 获取完整权限树（从数据库读取）
func (s *PermissionService) GetPermissionTree() ([]PermissionNode, error) {
	s.cacheMu.RLock()
	perms := s.permissionsCache
	s.cacheMu.RUnlock()

	if len(perms) == 0 {
		var err error
		perms, err = s.permRepo.List()
		if err != nil {
			return nil, err
		}
	}

	return s.buildTree(perms, 0), nil
}

// buildTree 构建权限树
func (s *PermissionService) buildTree(perms []model.Permission, parentID int64) []PermissionNode {
	var nodes []PermissionNode
	for _, p := range perms {
		if p.ParentID == parentID {
			node := PermissionNode{
				ID:       p.Code,
				Code:     p.Code,
				Label:    p.Name,
				Type:     p.Type,
				APIPath:  p.APIPath,
				Method:   p.APIMethod,
				IsPublic: p.IsPublic,
				Disabled: p.IsPublic, // 公共权限不可编辑
				Children: s.buildTree(perms, p.ID),
			}
			nodes = append(nodes, node)
		}
	}
	return nodes
}

// GetRolePermissions 获取角色的权限码列表
func (s *PermissionService) GetRolePermissions(roleCode string) ([]string, error) {
	// 获取角色
	role, err := s.roleRepo.GetByCode(roleCode)
	if err != nil {
		return nil, err
	}

	// 获取角色的权限ID列表
	permIDs, err := s.rolePermRepo.GetPermissionIDsByRoleID(role.ID)
	if err != nil {
		return nil, err
	}

	if len(permIDs) == 0 {
		return []string{}, nil
	}

	// 获取权限详情
	perms, err := s.permRepo.GetByIDs(permIDs)
	if err != nil {
		return nil, err
	}

	// 提取权限码
	codes := make([]string, 0, len(perms))
	for _, p := range perms {
		codes = append(codes, p.Code)
	}

	return codes, nil
}

// UpdateRolePermissions 更新角色权限（接收权限码列表）
func (s *PermissionService) UpdateRolePermissions(ctx context.Context, roleCode string, permissionCodes []string) (int, error) {
	// 获取角色
	role, err := s.roleRepo.GetByCode(roleCode)
	if err != nil {
		return 0, err
	}

	// 获取权限ID列表
	var permIDs []int64
	if len(permissionCodes) > 0 {
		perms, err := s.permRepo.GetByCodes(permissionCodes)
		if err != nil {
			return 0, err
		}
		for _, p := range perms {
			permIDs = append(permIDs, p.ID)
		}
	}

	// 替换角色权限
	if err := s.rolePermRepo.ReplaceRolePermissions(role.ID, permIDs); err != nil {
		return 0, err
	}

	// 查找拥有该身份的所有用户
	var userIDs []int64
	err = s.db.Raw(`
		SELECT DISTINCT user_id
		FROM user_identities
		WHERE identity_type = ? AND status = 'active'
	`, roleCode).Scan(&userIDs).Error
	if err != nil {
		return 0, err
	}

	// 撤销这些用户的所有 token（B端）
	if s.blacklistService != nil {
		for _, userID := range userIDs {
			s.blacklistService.RevokeUserTokens(ctx, userID, "b_end")
		}
	}

	return len(userIDs), nil
}

// GetUserPermissionCodes 获取用户的权限码列表（根据所有角色计算并集）
//
// 业务目的：计算用户基于所有角色的权限并集，用于前端权限控制
//
// 主要流程：
// 1. 根据角色编码列表获取角色信息
// 2. 获取所有角色的权限ID（自动去重）
// 3. 根据权限ID获取权限详情
// 4. 提取并返回权限码列表
//
// 多角色处理：用户可能有多个角色，最终权限是所有角色权限的并集
func (s *PermissionService) GetUserPermissionCodes(roleCodes []string) ([]string, error) {
	if len(roleCodes) == 0 {
		return []string{}, nil
	}

	// 获取角色列表
	roles, err := s.roleRepo.GetByCodes(roleCodes)
	if err != nil {
		return nil, err
	}

	if len(roles) == 0 {
		return []string{}, nil
	}

	// 获取角色ID列表
	roleIDs := make([]int64, 0, len(roles))
	for _, r := range roles {
		roleIDs = append(roleIDs, r.ID)
	}

	// 获取所有角色的权限ID（去重）
	permIDs, err := s.rolePermRepo.GetPermissionIDsByRoleIDs(roleIDs)
	if err != nil {
		return nil, err
	}

	if len(permIDs) == 0 {
		return []string{}, nil
	}

	// 获取权限详情
	perms, err := s.permRepo.GetByIDs(permIDs)
	if err != nil {
		return nil, err
	}

	// 提取权限码（去重）
	codeSet := make(map[string]bool)
	for _, p := range perms {
		codeSet[p.Code] = true
	}

	codes := make([]string, 0, len(codeSet))
	for code := range codeSet {
		codes = append(codes, code)
	}

	return codes, nil
}

// GetUserPermissions 兼容旧接口，返回权限码列表
func (s *PermissionService) GetUserPermissions(roleCodes []string) ([]string, error) {
	return s.GetUserPermissionCodes(roleCodes)
}

// CheckAPIPermission 检查用户是否有访问指定API的权限
//
// 业务目的：在 API 请求时验证用户是否有权限访问该接口
//
// 主要流程：
// 1. 根据角色编码列表获取角色信息
// 2. 获取所有角色的权限ID
// 3. 获取权限详情（包含 API 路径和方法）
// 4. 遍历权限，检查是否有匹配的 API 路径和方法
//
// 路径匹配：支持通配符 *，例如 /api/v1/tasks/*/claim 匹配 /api/v1/tasks/123/claim
//
// 返回值：true 表示有权限，false 表示无权限
func (s *PermissionService) CheckAPIPermission(roleCodes []string, path, method string) (bool, error) {
	if len(roleCodes) == 0 {
		return false, nil
	}

	// 获取角色列表
	roles, err := s.roleRepo.GetByCodes(roleCodes)
	if err != nil {
		return false, err
	}

	if len(roles) == 0 {
		return false, nil
	}

	// 获取角色ID列表
	roleIDs := make([]int64, 0, len(roles))
	for _, r := range roles {
		roleIDs = append(roleIDs, r.ID)
	}

	// 获取所有角色的权限ID
	permIDs, err := s.rolePermRepo.GetPermissionIDsByRoleIDs(roleIDs)
	if err != nil {
		return false, err
	}

	if len(permIDs) == 0 {
		return false, nil
	}

	// 获取权限详情
	perms, err := s.permRepo.GetByIDs(permIDs)
	if err != nil {
		return false, err
	}

	// 检查是否有匹配的权限
	for _, p := range perms {
		if p.APIPath != "" && p.APIMethod != "" {
			if matchPath(p.APIPath, path) && strings.EqualFold(p.APIMethod, method) {
				return true, nil
			}
		}
	}

	return false, nil
}

// IsPublicAPI 检查是否为公共API
func (s *PermissionService) IsPublicAPI(path, method string) bool {
	s.cacheMu.RLock()
	publicPerms := s.publicAPIsCache
	s.cacheMu.RUnlock()

	for _, p := range publicPerms {
		if p.APIPath != "" && p.APIMethod != "" {
			if matchPath(p.APIPath, path) && strings.EqualFold(p.APIMethod, method) {
				return true
			}
		}
	}

	return false
}

// matchPath 路径匹配，支持通配符 *
// 例如: /api/v1/b/tasks/*/claim 匹配 /api/v1/b/tasks/123/claim
func matchPath(pattern, path string) bool {
	patternParts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")

	if len(patternParts) != len(pathParts) {
		return false
	}

	for i, part := range patternParts {
		if part == "*" {
			continue
		}
		if part != pathParts[i] {
			return false
		}
	}
	return true
}

// GetAllRoles 获取所有角色列表
func (s *PermissionService) GetAllRoles() ([]model.Role, error) {
	return s.roleRepo.List()
}
