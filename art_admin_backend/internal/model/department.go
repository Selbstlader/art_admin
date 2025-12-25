package model

import (
	"time"

	"gorm.io/gorm"
)

// Department 部门模型
type Department struct {
	DeptID     int64          `gorm:"column:id;primaryKey;autoIncrement" json:"deptId"`
	ParentID   *int64         `gorm:"index" json:"parentId"`
	DeptName   string         `gorm:"type:varchar(50);not null" json:"deptName"`
	DeptCode   string         `gorm:"type:varchar(50);not null;index:idx_sys_department_dept_code,unique" json:"deptCode"`
	OrderNum   int            `gorm:"column:sort;type:int;default:0" json:"orderNum"`
	Leader     string         `gorm:"type:varchar(20)" json:"leader"`
	Phone      string         `gorm:"type:varchar(11)" json:"phone"`
	Email      string         `gorm:"type:varchar(50)" json:"email"`
	Status     int            `gorm:"type:tinyint(1);default:1;index" json:"status"` // 1-正常 0-停用
	CreateTime time.Time      `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime time.Time      `gorm:"autoUpdateTime" json:"updateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	Parent   *Department  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []Department `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

// TableName 表名
func (Department) TableName() string {
	return "sys_department"
}
