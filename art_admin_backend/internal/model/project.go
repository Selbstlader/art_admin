package model

import (
	"time"

	"gorm.io/gorm"
)

// Project 项目模型
type Project struct {
	ID              int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	TemplateID      int64          `gorm:"type:bigint;default:0;index" json:"templateId"`                                           // 来源模板ID(0表示手动创建)
	Name            string         `gorm:"type:varchar(100);not null;index" json:"name"`                                            // 项目名称
	Code            string         `gorm:"type:varchar(50);uniqueIndex" json:"code"`                                                // 项目编号
	Description     string         `gorm:"type:text" json:"description"`                                                            // 项目描述
	Category        string         `gorm:"type:varchar(50);index" json:"category"`                                                  // 项目分类
	Priority        int            `gorm:"type:tinyint(1);default:2;comment:优先级(1-高 2-中 3-低)" json:"priority"`                      // 优先级
	Status          int            `gorm:"type:tinyint(1);default:1;index;comment:状态(1-计划中 2-进行中 3-已完成 4-已暂停 5-已取消)" json:"status"` // 状态
	Progress        int            `gorm:"type:tinyint;default:0;comment:进度(0-100)" json:"progress"`                                // 进度百分比
	StartDate       time.Time      `gorm:"type:date;not null;index" json:"startDate"`                                               // 开始日期
	EndDate         time.Time      `gorm:"type:date;not null;index" json:"endDate"`                                                 // 结束日期
	ActualStartDate *time.Time     `gorm:"type:date" json:"actualStartDate"`                                                        // 实际开始日期
	ActualEndDate   *time.Time     `gorm:"type:date" json:"actualEndDate"`                                                          // 实际结束日期
	ManagerID       int64          `gorm:"type:bigint;not null;index" json:"managerId"`                                             // 项目经理ID
	ManagerName     string         `gorm:"type:varchar(50)" json:"managerName"`                                                     // 项目经理名称
	CreatorID       int64          `gorm:"type:bigint;index" json:"creatorId"`                                                      // 创建人ID
	CreatorName     string         `gorm:"type:varchar(50)" json:"creatorName"`                                                     // 创建人名称
	Budget          float64        `gorm:"type:decimal(15,2);default:0" json:"budget"`                                              // 预算
	ActualCost      float64        `gorm:"type:decimal(15,2);default:0" json:"actualCost"`                                          // 实际成本
	Remark          string         `gorm:"type:varchar(500)" json:"remark"`                                                         // 备注
	CreatedAt       time.Time      `gorm:"type:datetime;not null" json:"createdAt"`
	UpdatedAt       time.Time      `gorm:"type:datetime;not null" json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (Project) TableName() string {
	return "pm_project"
}

// ProjectMember 项目成员模型
type ProjectMember struct {
	ID        int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	ProjectID int64          `gorm:"type:bigint;not null;index" json:"projectId"` // 项目ID
	UserID    int64          `gorm:"type:bigint;not null;index" json:"userId"`    // 用户ID
	UserName  string         `gorm:"type:varchar(50)" json:"userName"`            // 用户名称
	Role      string         `gorm:"type:varchar(50);comment:项目角色" json:"role"`   // 项目角色(项目经理/开发/测试等)
	JoinDate  time.Time      `gorm:"type:date;not null" json:"joinDate"`          // 加入日期
	CreatedAt time.Time      `gorm:"type:datetime;not null" json:"createdAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (ProjectMember) TableName() string {
	return "pm_project_member"
}
