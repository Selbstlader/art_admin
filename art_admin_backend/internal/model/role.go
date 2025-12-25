package model

import (
	"time"

	"gorm.io/gorm"
)

// Role 角色模型
type Role struct {
	RoleID      int64          `gorm:"primaryKey;autoIncrement" json:"roleId"`
	RoleName    string         `gorm:"type:varchar(50);not null" json:"roleName"`
	RoleCode    string         `gorm:"type:varchar(50);not null;index:idx_role_code,unique" json:"roleCode"`
	Description string         `gorm:"type:varchar(200)" json:"description"`
	Enabled     bool           `gorm:"type:tinyint(1);default:1;index" json:"enabled"`
	CreateTime  time.Time      `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime  time.Time      `gorm:"autoUpdateTime" json:"updateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	Users   []User   `gorm:"many2many:sys_user_role;foreignKey:RoleID;joinForeignKey:RoleID;References:ID;joinReferences:UserID" json:"-"`
	Menus   []Menu   `gorm:"many2many:sys_role_menu;foreignKey:RoleID;joinForeignKey:RoleID;References:ID;joinReferences:MenuID" json:"-"`
	Buttons []Button `gorm:"many2many:sys_role_button;foreignKey:RoleID;joinForeignKey:RoleID;References:ID;joinReferences:ButtonID" json:"-"`
}

// TableName 表名
func (Role) TableName() string {
	return "sys_role"
}
