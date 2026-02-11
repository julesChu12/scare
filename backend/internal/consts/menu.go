package consts

import (
	"time"

	"gorm.io/gorm"
)

// Menu 菜单模型（前端路由配置）
type Menu struct {
	ID             int64          `gorm:"primaryKey" json:"id"`
	ParentID       int64          `gorm:"index;default:0" json:"parent_id"`
	Name           string         `gorm:"size:50;not null" json:"name"`
	Path           string         `gorm:"size:200" json:"path"`
	Component      string         `gorm:"size:200" json:"component"`
	Icon           string         `gorm:"size:100" json:"icon"`
	PermissionCode string         `gorm:"column:permission_code;size:100;index" json:"permission_code"`
	Sort           int            `gorm:"default:0" json:"sort"`
	Hidden         bool           `gorm:"default:false" json:"hidden"`
	Status         string         `gorm:"size:20;index;default:'active'" json:"status"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Children       []Menu         `gorm:"-" json:"children,omitempty"`
}

// TableName 指定表名
func (Menu) TableName() string {
	return "menus"
}
