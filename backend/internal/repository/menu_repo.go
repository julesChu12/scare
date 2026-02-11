package repository

import (
	"community-elderly-care-platform/internal/consts"
	"community-elderly-care-platform/internal/dao/query"

	"gorm.io/gorm"
)

type MenuRepository struct {
	q  *query.Query
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) *MenuRepository {
	return &MenuRepository{
		q:  query.Use(db),
		db: db,
	}
}

// List 获取所有菜单（管理用）
func (r *MenuRepository) List() ([]consts.Menu, error) {
	var menus []consts.Menu
	if err := r.db.Order("sort ASC, id ASC").Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}

// GetByID 根据ID获取菜单
func (r *MenuRepository) GetByID(id int64) (*consts.Menu, error) {
	var menu consts.Menu
	if err := r.db.First(&menu, id).Error; err != nil {
		return nil, err
	}
	return &menu, nil
}

// GetByParentID 获取指定父级的子菜单
func (r *MenuRepository) GetByParentID(parentID int64) ([]consts.Menu, error) {
	var menus []consts.Menu
	if err := r.db.Where("parent_id = ?", parentID).Order("sort ASC, id ASC").Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}

// GetActiveMenus 获取所有激活的菜单
func (r *MenuRepository) GetActiveMenus() ([]consts.Menu, error) {
	var menus []consts.Menu
	if err := r.db.Where("status = ?", "active").Order("sort ASC, id ASC").Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}

// GetActiveMenusByPermissions 根据权限码列表获取可见菜单
func (r *MenuRepository) GetActiveMenusByPermissions(permissionCodes []string) ([]consts.Menu, error) {
	var menus []consts.Menu
	query := r.db.Where("status = ?", "active").Where("hidden = ?", false)

	// 如果有权限列表，过滤需要权限的菜单
	if len(permissionCodes) > 0 {
		// 获取不需要权限的菜单 + 用户有权限的菜单
		query = query.Where("permission_code = '' OR permission_code IS NULL OR permission_code IN ?", permissionCodes)
	} else {
		// 没有权限时只获取不需要权限的菜单
		query = query.Where("permission_code = '' OR permission_code IS NULL")
	}

	if err := query.Order("sort ASC, id ASC").Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}

// Create 创建菜单
func (r *MenuRepository) Create(menu *consts.Menu) error {
	return r.db.Create(menu).Error
}

// Update 更新菜单
func (r *MenuRepository) Update(menu *consts.Menu) error {
	return r.db.Save(menu).Error
}

// Delete 删除菜单（软删除）
func (r *MenuRepository) Delete(id int64) error {
	return r.db.Delete(&consts.Menu{}, id).Error
}

// BatchUpdateSort 批量更新排序
func (r *MenuRepository) BatchUpdateSort(updates []struct {
	ID   int64
	Sort int
}) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, u := range updates {
			if err := tx.Model(&consts.Menu{}).Where("id = ?", u.ID).Update("sort", u.Sort).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
