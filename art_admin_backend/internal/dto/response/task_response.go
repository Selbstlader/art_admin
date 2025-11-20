package response

import "time"

// TaskResponse 任务响应
type TaskResponse struct {
	ID              int64      `json:"id"`
	ProjectID       int64      `json:"projectId"`
	ParentID        int64      `json:"parentId"`
	TemplateTaskID  int64      `json:"templateTaskId"`
	Name            string     `json:"name"`
	Description     string     `json:"description"`
	Type            string     `json:"type"`
	Priority        int        `json:"priority"`
	Status          int        `json:"status"`
	Progress        int        `json:"progress"`
	StartDate       time.Time  `json:"startDate"`
	EndDate         time.Time  `json:"endDate"`
	ActualStartDate *time.Time `json:"actualStartDate"`
	ActualEndDate   *time.Time `json:"actualEndDate"`
	AssigneeID      int64      `json:"assigneeId"`
	AssigneeName    string     `json:"assigneeName"`
	EstimatedHours  float64    `json:"estimatedHours"`
	ActualHours     float64    `json:"actualHours"`
	IsOverdue       bool       `json:"isOverdue"`
	Sort            int        `json:"sort"`
	Remark          string     `json:"remark"`
	CreatorID       int64      `json:"creatorId"`
	CreatorName     string     `json:"creatorName"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}

// TaskDetailResponse 任务详情响应
type TaskDetailResponse struct {
	TaskResponse
	Children     []TaskResponse           `json:"children"`
	Dependencies []TaskDependencyResponse `json:"dependencies"`
	Comments     []TaskCommentResponse    `json:"comments"`
	Attachments  []TaskAttachmentResponse `json:"attachments"`
}

// TaskDependencyResponse 任务依赖响应
type TaskDependencyResponse struct {
	ID              int64     `json:"id"`
	ProjectID       int64     `json:"projectId"`
	PredecessorID   int64     `json:"predecessorId"`
	PredecessorName string    `json:"predecessorName"`
	SuccessorID     int64     `json:"successorId"`
	SuccessorName   string    `json:"successorName"`
	DependencyType  string    `json:"dependencyType"`
	LagDays         int       `json:"lagDays"`
	CreatedAt       time.Time `json:"createdAt"`
}

// TaskCommentResponse 任务评论响应
type TaskCommentResponse struct {
	ID        int64     `json:"id"`
	TaskID    int64     `json:"taskId"`
	UserID    int64     `json:"userId"`
	UserName  string    `json:"userName"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

// TaskAttachmentResponse 任务附件响应
type TaskAttachmentResponse struct {
	ID           int64     `json:"id"`
	TaskID       int64     `json:"taskId"`
	FileName     string    `json:"fileName"`
	FileSize     int64     `json:"fileSize"`
	FileType     string    `json:"fileType"`
	FilePath     string    `json:"filePath"`
	UploaderID   int64     `json:"uploaderId"`
	UploaderName string    `json:"uploaderName"`
	CreatedAt    time.Time `json:"createdAt"`
}

// TaskGanttDataResponse 甘特图数据响应
type TaskGanttDataResponse struct {
	Tasks        []TaskGanttItem          `json:"tasks"`
	Dependencies []TaskDependencyResponse `json:"dependencies"`
}

// TaskGanttItem 甘特图任务项
type TaskGanttItem struct {
	ID           int64     `json:"id"`
	Text         string    `json:"text"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	Duration     int       `json:"duration"`
	Progress     float64   `json:"progress"`
	ParentID     int64     `json:"parent"`
	Type         string    `json:"type"`
	Priority     int       `json:"priority"`
	Status       int       `json:"status"`
	AssigneeID   int64     `json:"assigneeId"`
	AssigneeName string    `json:"assigneeName"`
	IsOverdue    bool      `json:"isOverdue"`
	Open         bool      `json:"open"`
}
