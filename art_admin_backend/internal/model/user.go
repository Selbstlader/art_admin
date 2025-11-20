package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID         int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserName   string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"userName"`
	NickName   string         `gorm:"type:varchar(50)" json:"nickName"`
	Password   string         `gorm:"type:varchar(255);not null" json:"-"`
	Email      string         `gorm:"type:varchar(100);index" json:"email"`
	UserPhone  string         `gorm:"type:varchar(20)" json:"userPhone"`
	UserGender string         `gorm:"type:varchar(10);default:unknown" json:"userGender"`
	Avatar     string         `gorm:"type:varchar(500)" json:"avatar"`
	Status     string         `gorm:"type:char(1);default:1;index" json:"status"`
	DeptID     *int64         `gorm:"index" json:"deptId"`
	CreateBy   string         `gorm:"type:varchar(50)" json:"createBy"`
	CreateTime time.Time      `gorm:"autoCreateTime" json:"createTime"`
	UpdateBy   string         `gorm:"type:varchar(50)" json:"updateBy"`
	UpdateTime time.Time      `gorm:"autoUpdateTime" json:"updateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	Roles      []Role      `gorm:"many2many:sys_user_role;foreignKey:ID;joinForeignKey:UserID;References:RoleID;joinReferences:RoleID" json:"-"`
	Department *Department `gorm:"foreignKey:DeptID" json:"department,omitempty"`
}

// TableName 表名
func (User) TableName() string {
	return "sys_user"
}
