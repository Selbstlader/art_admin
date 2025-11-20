package model

import (
	"time"

	"gorm.io/gorm"
)

// ProjectTemplate 项目模板模型
type ProjectTemplate struct {
	ID          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string         `gorm:"type:varchar(100);not null;index" json:"name"`           // 模板名称
	Description string         `gorm:"type:varchar(500)" json:"description"`                   // 模板描述
	Category    string         `gorm:"type:varchar(50);index" json:"category"`                 // 模板分类
	Icon        string         `gorm:"type:varchar(100)" json:"icon"`                          // 模板图标
	Duration    int            `gorm:"type:int;comment:预计工期(天)" json:"duration"`               // 预计工期(天)
	IsPublic    bool           `gorm:"type:tinyint(1);default:1;comment:是否公开" json:"isPublic"` // 是否公开
	CreatorID   int64          `gorm:"type:bigint;index" json:"creatorId"`                     // 创建人ID
	CreatorName string         `gorm:"type:varchar(50)" json:"creatorName"`                    // 创建人名称
	UseCount    int            `gorm:"type:int;default:0;comment:使用次数" json:"useCount"`        // 使用次数
	Status      int            `gorm:"type:tinyint(1);default:1;index" json:"status"`          // 状态(1-启用 0-禁用)
	CreatedAt   time.Time      `gorm:"type:datetime;not null" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"type:datetime;not null" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (ProjectTemplate) TableName() string {
	return "pm_project_template"
}

// TemplateTask 模板任务模型
type TemplateTask struct {
	ID          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	TemplateID  int64          `gorm:"type:bigint;not null;index" json:"templateId"`                       // 所属模板ID
	ParentID    int64          `gorm:"type:bigint;default:0;index" json:"parentId"`                        // 父任务ID(0表示根任务)
	Name        string         `gorm:"type:varchar(200);not null" json:"name"`                             // 任务名称
	Description string         `gorm:"type:text" json:"description"`                                       // 任务描述
	StartDay    int            `gorm:"type:int;not null;comment:开始天数(相对项目开始)" json:"startDay"`             // 开始天数(相对项目开始)
	Duration    int            `gorm:"type:int;not null;comment:持续天数" json:"duration"`                     // 持续天数
	IsMilestone bool           `gorm:"type:tinyint(1);default:0;comment:是否里程碑" json:"isMilestone"`         // 是否里程碑
	Priority    int            `gorm:"type:tinyint(1);default:2;comment:优先级(1-高 2-中 3-低)" json:"priority"` // 优先级
	Sort        int            `gorm:"type:int;default:0" json:"sort"`                                     // 排序
	CreatedAt   time.Time      `gorm:"type:datetime;not null" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"type:datetime;not null" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (TemplateTask) TableName() string {
	return "pm_template_task"
}

// TemplateTaskDependency 模板任务依赖关系
type TemplateTaskDependency struct {
	ID             int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	TemplateID     int64          `gorm:"type:bigint;not null;index" json:"templateId"`        // 所属模板ID
	PredecessorID  int64          `gorm:"type:bigint;not null;index" json:"predecessorId"`     // 前置任务ID
	SuccessorID    int64          `gorm:"type:bigint;not null;index" json:"successorId"`       // 后置任务ID
	DependencyType string         `gorm:"type:varchar(10);default:'FS'" json:"dependencyType"` // 依赖类型(FS-完成到开始,SS-开始到开始,FF-完成到完成,SF-开始到完成)
	LagDays        int            `gorm:"type:int;default:0;comment:滞后天数" json:"lagDays"`      // 滞后天数
	CreatedAt      time.Time      `gorm:"type:datetime;not null" json:"createdAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (TemplateTaskDependency) TableName() string {
	return "pm_template_task_dependency"
}
