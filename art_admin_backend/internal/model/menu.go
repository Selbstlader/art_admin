package model

import (
	"time"

	"gorm.io/gorm"
)

// Menu 菜单模型
type Menu struct {
	ID            int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	ParentID      int64          `gorm:"default:0;index" json:"parentId"`
	Name          string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"` // 路由名称/权限标识（对应前端label）
	Path          string         `gorm:"type:varchar(200);not null" json:"path"`             // 路由地址
	Title         string         `gorm:"type:varchar(50);not null" json:"title"`             // 菜单显示名称（对应前端name）
	Component     string         `gorm:"type:varchar(200)" json:"component"`                 // 组件路径
	Redirect      string         `gorm:"type:varchar(200)" json:"redirect"`                  // 重定向路径
	Icon          string         `gorm:"type:varchar(50)" json:"icon"`                       // 图标
	IsEnable      bool           `gorm:"type:tinyint(1);default:1;index" json:"isEnable"`    // 是否启用
	Sort          int            `gorm:"default:0" json:"sort"`                              // 排序
	IsMenu        bool           `gorm:"type:tinyint(1);default:1" json:"isMenu"`            // 是否菜单（区分按钮）
	KeepAlive     bool           `gorm:"type:tinyint(1);default:1" json:"keepAlive"`         // 页面缓存
	IsHide        bool           `gorm:"type:tinyint(1);default:0" json:"isHide"`            // 隐藏菜单
	IsHideTab     bool           `gorm:"type:tinyint(1);default:0" json:"isHideTab"`         // 隐藏标签
	Link          string         `gorm:"type:varchar(500)" json:"link"`                      // 外部链接
	IsIframe      bool           `gorm:"type:tinyint(1);default:0" json:"isIframe"`          // 是否内嵌
	ShowBadge     bool           `gorm:"type:tinyint(1);default:0" json:"showBadge"`         // 显示徽章
	ShowTextBadge string         `gorm:"type:varchar(50)" json:"showTextBadge"`              // 文本徽章
	FixedTab      bool           `gorm:"type:tinyint(1);default:0" json:"fixedTab"`          // 固定标签
	ActivePath    string         `gorm:"type:varchar(200)" json:"activePath"`                // 激活路径
	Roles         []string       `gorm:"-" json:"roles"`                                     // 角色权限（前端用，不存储）
	IsFullPage    bool           `gorm:"type:tinyint(1);default:0" json:"isFullPage"`        // 全屏页面
	CreateTime    time.Time      `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime    time.Time      `gorm:"autoUpdateTime" json:"updateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	Buttons []Button `gorm:"foreignKey:MenuID" json:"-"`
	RoleRel []Role   `gorm:"many2many:sys_role_menu;foreignKey:ID;joinForeignKey:MenuID;References:RoleID;joinReferences:RoleID" json:"-"`
}

// TableName 表名
func (Menu) TableName() string {
	return "sys_menu"
}
