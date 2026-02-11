package service

import (
	"community-elderly-care-platform/internal/consts"
	"community-elderly-care-platform/internal/repository"
)

type MenuService struct {
	menuRepo *repository.MenuRepository
}

func NewMenuService(menuRepo *repository.MenuRepository) *MenuService {
	return &MenuService{menuRepo: menuRepo}
}

// GetMenuTree 获取完整菜单树（管理用）
func (s *MenuService) GetMenuTree() ([]consts.Menu, error) {
	menus, err := s.menuRepo.List()
	if err != nil {
		return nil, err
	}
	return s.buildTree(menus, 0), nil
}

// GetUserMenus 获取用户可见菜单树
func (s *MenuService) GetUserMenus(permissionCodes []string) ([]consts.Menu, error) {
	menus, err := s.menuRepo.GetActiveMenusByPermissions(permissionCodes)
	if err != nil {
		return nil, err
	}
	return s.buildTree(menus, 0), nil
}

// GetByID 获取单个菜单
func (s *MenuService) GetByID(id int64) (*consts.Menu, error) {
	return s.menuRepo.GetByID(id)
}

// Create 创建菜单
func (s *MenuService) Create(menu *consts.Menu) error {
	return s.menuRepo.Create(menu)
}

// Update 更新菜单
func (s *MenuService) Update(id int64, menu *consts.Menu) error {
	existing, err := s.menuRepo.GetByID(id)
	if err != nil {
		return err
	}

	// 更新字段
	existing.ParentID = menu.ParentID
	existing.Name = menu.Name
	existing.Path = menu.Path
	existing.Component = menu.Component
	existing.Icon = menu.Icon
	existing.PermissionCode = menu.PermissionCode
	existing.Sort = menu.Sort
	existing.Hidden = menu.Hidden
	existing.Status = menu.Status

	return s.menuRepo.Update(existing)
}

// Delete 删除菜单
func (s *MenuService) Delete(id int64) error {
	return s.menuRepo.Delete(id)
}

// BatchUpdateSort 批量更新排序
func (s *MenuService) BatchUpdateSort(updates []struct {
	ID   int64
	Sort int
}) error {
	return s.menuRepo.BatchUpdateSort(updates)
}

// buildTree 构建菜单树
func (s *MenuService) buildTree(menus []consts.Menu, parentID int64) []consts.Menu {
	var tree []consts.Menu
	for _, menu := range menus {
		if menu.ParentID == parentID {
			menu.Children = s.buildTree(menus, menu.ID)
			tree = append(tree, menu)
		}
	}
	return tree
}
