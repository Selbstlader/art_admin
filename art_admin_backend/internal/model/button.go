package model

import (
	"time"
)

// Button 按钮权限模型
type Button struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	MenuID     int64     `gorm:"not null;index" json:"menuId"`               // 所属菜单ID
	AuthName   string    `gorm:"type:varchar(50);not null" json:"authName"`  // 权限名称
	AuthLabel  string    `gorm:"type:varchar(100);not null;index" json:"authLabel"` // 权限标识
	AuthIcon   string    `gorm:"type:varchar(50)" json:"authIcon"`           // 权限图标
	AuthSort   int       `gorm:"default:0" json:"authSort"`                  // 权限排序
	IsEnable   bool      `gorm:"type:tinyint(1);default:1" json:"isEnable"`  // 是否启用
	CreateTime time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime time.Time `gorm:"autoUpdateTime" json:"updateTime"`

	// 关联关系
	Menu  Menu   `gorm:"foreignKey:MenuID" json:"-"`
	Roles []Role `gorm:"many2many:sys_role_button;foreignKey:ID;joinForeignKey:ButtonID;References:RoleID;joinReferences:RoleID" json:"-"`
}

// TableName 表名
func (Button) TableName() string {
	return "sys_button"
}
