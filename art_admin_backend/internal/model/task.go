package model

import (
	"time"

	"gorm.io/gorm"
)

// Task 任务模型
type Task struct {
	ID              int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	ProjectID       int64          `gorm:"type:bigint;not null;index" json:"projectId"`                                             // 所属项目ID
	ParentID        int64          `gorm:"type:bigint;default:0;index" json:"parentId"`                                             // 父任务ID(0表示根任务)
	TemplateTaskID  int64          `gorm:"type:bigint;default:0" json:"templateTaskId"`                                             // 来源模板任务ID
	Name            string         `gorm:"type:varchar(200);not null" json:"name"`                                                  // 任务名称
	Description     string         `gorm:"type:text" json:"description"`                                                            // 任务描述
	Type            string         `gorm:"type:varchar(20);default:'task'" json:"type"`                                             // 任务类型(task-普通任务, milestone-里程碑)
	Priority        int            `gorm:"type:tinyint(1);default:2;comment:优先级(1-高 2-中 3-低)" json:"priority"`                      // 优先级
	Status          int            `gorm:"type:tinyint(1);default:1;index;comment:状态(1-未开始 2-进行中 3-已完成 4-已延期 5-已取消)" json:"status"` // 状态
	Progress        int            `gorm:"type:tinyint;default:0;comment:进度(0-100)" json:"progress"`                                // 进度百分比
	StartDate       time.Time      `gorm:"type:date;not null;index" json:"startDate"`                                               // 计划开始日期
	EndDate         time.Time      `gorm:"type:date;not null;index" json:"endDate"`                                                 // 计划结束日期
	ActualStartDate *time.Time     `gorm:"type:date" json:"actualStartDate"`                                                        // 实际开始日期
	ActualEndDate   *time.Time     `gorm:"type:date" json:"actualEndDate"`                                                          // 实际结束日期
	AssigneeID      int64          `gorm:"type:bigint;index" json:"assigneeId"`                                                     // 负责人ID
	AssigneeName    string         `gorm:"type:varchar(50)" json:"assigneeName"`                                                    // 负责人名称
	EstimatedHours  float64        `gorm:"type:decimal(10,2);default:0;comment:预估工时" json:"estimatedHours"`                         // 预估工时
	ActualHours     float64        `gorm:"type:decimal(10,2);default:0;comment:实际工时" json:"actualHours"`                            // 实际工时
	IsOverdue       bool           `gorm:"type:tinyint(1);default:0;comment:是否超期" json:"isOverdue"`                                 // 是否超期
	Sort            int            `gorm:"type:int;default:0" json:"sort"`                                                          // 排序
	Remark          string         `gorm:"type:varchar(500)" json:"remark"`                                                         // 备注
	CreatorID       int64          `gorm:"type:bigint;index" json:"creatorId"`                                                      // 创建人ID
	CreatorName     string         `gorm:"type:varchar(50)" json:"creatorName"`                                                     // 创建人名称
	CreatedAt       time.Time      `gorm:"type:datetime;not null" json:"createdAt"`
	UpdatedAt       time.Time      `gorm:"type:datetime;not null" json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (Task) TableName() string {
	return "pm_task"
}

// TaskDependency 任务依赖关系
type TaskDependency struct {
	ID             int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	ProjectID      int64          `gorm:"type:bigint;not null;index" json:"projectId"`         // 所属项目ID
	PredecessorID  int64          `gorm:"type:bigint;not null;index" json:"predecessorId"`     // 前置任务ID
	SuccessorID    int64          `gorm:"type:bigint;not null;index" json:"successorId"`       // 后置任务ID
	DependencyType string         `gorm:"type:varchar(10);default:'FS'" json:"dependencyType"` // 依赖类型(FS-完成到开始,SS-开始到开始,FF-完成到完成,SF-开始到完成)
	LagDays        int            `gorm:"type:int;default:0;comment:滞后天数" json:"lagDays"`      // 滞后天数
	CreatedAt      time.Time      `gorm:"type:datetime;not null" json:"createdAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (TaskDependency) TableName() string {
	return "pm_task_dependency"
}

// TaskComment 任务评论
type TaskComment struct {
	ID        int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	TaskID    int64          `gorm:"type:bigint;not null;index" json:"taskId"` // 任务ID
	UserID    int64          `gorm:"type:bigint;not null;index" json:"userId"` // 评论人ID
	UserName  string         `gorm:"type:varchar(50)" json:"userName"`         // 评论人名称
	Content   string         `gorm:"type:text;not null" json:"content"`        // 评论内容
	CreatedAt time.Time      `gorm:"type:datetime;not null" json:"createdAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (TaskComment) TableName() string {
	return "pm_task_comment"
}

// TaskAttachment 任务附件
type TaskAttachment struct {
	ID           int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	TaskID       int64          `gorm:"type:bigint;not null;index" json:"taskId"`   // 任务ID
	FileName     string         `gorm:"type:varchar(200);not null" json:"fileName"` // 文件名
	FileSize     int64          `gorm:"type:bigint" json:"fileSize"`                // 文件大小(字节)
	FileType     string         `gorm:"type:varchar(50)" json:"fileType"`           // 文件类型
	FilePath     string         `gorm:"type:varchar(500);not null" json:"filePath"` // 文件路径
	UploaderID   int64          `gorm:"type:bigint;not null" json:"uploaderId"`     // 上传人ID
	UploaderName string         `gorm:"type:varchar(50)" json:"uploaderName"`       // 上传人名称
	CreatedAt    time.Time      `gorm:"type:datetime;not null" json:"createdAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (TaskAttachment) TableName() string {
	return "pm_task_attachment"
}
