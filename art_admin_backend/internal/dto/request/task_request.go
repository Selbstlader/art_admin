package request

import "time"

// CreateTaskRequest 创建任务请求
type CreateTaskRequest struct {
	ProjectID      int64     `json:"projectId" binding:"required"`
	ParentID       int64     `json:"parentId"`
	Name           string    `json:"name" binding:"required,max=200"`
	Description    string    `json:"description"`
	Type           string    `json:"type" binding:"required,oneof=task milestone"`
	Priority       int       `json:"priority" binding:"min=1,max=3"`
	StartDate      time.Time `json:"startDate" binding:"required"`
	EndDate        time.Time `json:"endDate" binding:"required"`
	AssigneeID     int64     `json:"assigneeId"`
	EstimatedHours float64   `json:"estimatedHours" binding:"min=0"`
	Sort           int       `json:"sort"`
	Remark         string    `json:"remark" binding:"max=500"`
}

// UpdateTaskRequest 更新任务请求
type UpdateTaskRequest struct {
	ID              int64      `json:"id" binding:"required"`
	Name            string     `json:"name" binding:"required,max=200"`
	Description     string     `json:"description"`
	Type            string     `json:"type" binding:"required,oneof=task milestone"`
	Priority        int        `json:"priority" binding:"min=1,max=3"`
	Status          int        `json:"status" binding:"min=1,max=5"`
	Progress        int        `json:"progress" binding:"min=0,max=100"`
	StartDate       time.Time  `json:"startDate" binding:"required"`
	EndDate         time.Time  `json:"endDate" binding:"required"`
	ActualStartDate *time.Time `json:"actualStartDate"`
	ActualEndDate   *time.Time `json:"actualEndDate"`
	AssigneeID      int64      `json:"assigneeId"`
	EstimatedHours  float64    `json:"estimatedHours" binding:"min=0"`
	ActualHours     float64    `json:"actualHours" binding:"min=0"`
	Sort            int        `json:"sort"`
	Remark          string     `json:"remark" binding:"max=500"`
}

// BatchUpdateTaskRequest 批量更新任务请求(甘特图拖拽)
type BatchUpdateTaskRequest struct {
	Tasks []TaskUpdateItem `json:"tasks" binding:"required,min=1"`
}

// TaskUpdateItem 任务更新项
type TaskUpdateItem struct {
	ID        int64     `json:"id" binding:"required"`
	StartDate time.Time `json:"startDate" binding:"required"`
	EndDate   time.Time `json:"endDate" binding:"required"`
	Progress  int       `json:"progress" binding:"min=0,max=100"`
}

// GetTaskListRequest 获取任务列表请求
type GetTaskListRequest struct {
	Page       int    `form:"page" binding:"min=1"`
	PageSize   int    `form:"pageSize" binding:"min=1,max=100"`
	ProjectID  int64  `form:"projectId" binding:"required"`
	ParentID   *int64 `form:"parentId"`
	Name       string `form:"name"`
	Type       string `form:"type"`
	Status     *int   `form:"status"`
	AssigneeID *int64 `form:"assigneeId"`
	IsOverdue  *bool  `form:"isOverdue"`
}

// CreateTaskDependencyRequest 创建任务依赖请求
type CreateTaskDependencyRequest struct {
	ProjectID      int64  `json:"projectId" binding:"required"`
	PredecessorID  int64  `json:"predecessorId" binding:"required"`
	SuccessorID    int64  `json:"successorId" binding:"required"`
	DependencyType string `json:"dependencyType" binding:"required,oneof=FS SS FF SF"`
	LagDays        int    `json:"lagDays"`
}

// CreateTaskCommentRequest 创建任务评论请求
type CreateTaskCommentRequest struct {
	TaskID  int64  `json:"taskId" binding:"required"`
	Content string `json:"content" binding:"required,max=2000"`
}

// GetTaskGanttDataRequest 获取甘特图数据请求
type GetTaskGanttDataRequest struct {
	ProjectID int64 `form:"projectId" binding:"required"`
}

// QuickUpdateTaskRequest 快速更新任务请求(甘特图实时更新)
type QuickUpdateTaskRequest struct {
	ID        int64      `json:"id" binding:"required"`
	StartDate *time.Time `json:"startDate"`
	EndDate   *time.Time `json:"endDate"`
	Progress  *int       `json:"progress" binding:"omitempty,min=0,max=100"`
	Duration  *int       `json:"duration" binding:"omitempty,min=0"`
}
