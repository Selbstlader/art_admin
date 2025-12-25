package model

import (
	"time"

	"gorm.io/gorm"
)

// AppUser APP用户模型
type AppUser struct {
	ID            int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserName      string         `gorm:"type:varchar(50);not null;index:idx_app_user_user_name,unique" json:"userName"`
	NickName      string         `gorm:"type:varchar(50);not null" json:"nickName"`
	Password      string         `gorm:"type:varchar(255);not null" json:"-"`
	Phone         string         `gorm:"type:varchar(20);not null;index:idx_app_user_phone,unique" json:"phone"`
	Email         string         `gorm:"type:varchar(100)" json:"email"`
	Avatar        string         `gorm:"type:varchar(500)" json:"avatar"`
	UserType      string         `gorm:"type:char(1);default:2;index;comment:用户类型 1-管理员 2-普通用户" json:"userType"`
	Status        string         `gorm:"type:char(1);default:1;index;comment:状态 1-正常 2-禁用" json:"status"`
	LastLoginTime *time.Time     `gorm:"index" json:"lastLoginTime"`
	LastLoginIP   string         `gorm:"type:varchar(50)" json:"lastLoginIp"`
	CreateBy      string         `gorm:"type:varchar(50)" json:"createBy"`
	CreateTime    time.Time      `gorm:"autoCreateTime" json:"createTime"`
	UpdateBy      string         `gorm:"type:varchar(50)" json:"updateBy"`
	UpdateTime    time.Time      `gorm:"autoUpdateTime" json:"updateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (AppUser) TableName() string {
	return "app_user"
}
