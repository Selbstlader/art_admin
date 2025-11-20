package request

// CreateProjectTemplateRequest 创建项目模板请求
type CreateProjectTemplateRequest struct {
	Name        string                  `json:"name" binding:"required,max=100"`
	Description string                  `json:"description" binding:"max=500"`
	Category    string                  `json:"category" binding:"max=50"`
	Icon        string                  `json:"icon" binding:"max=100"`
	Duration    int                     `json:"duration" binding:"min=1"`
	IsPublic    bool                    `json:"isPublic"`
	Tasks       []CreateTemplateTaskReq `json:"tasks"`
}

// CreateTemplateTaskReq 创建模板任务请求
type CreateTemplateTaskReq struct {
	ParentID    int64  `json:"parentId"`
	Name        string `json:"name" binding:"required,max=200"`
	Description string `json:"description"`
	StartDay    int    `json:"startDay" binding:"min=0"`
	Duration    int    `json:"duration" binding:"min=1"`
	IsMilestone bool   `json:"isMilestone"`
	Priority    int    `json:"priority" binding:"min=1,max=3"`
	Sort        int    `json:"sort"`
}

// UpdateProjectTemplateRequest 更新项目模板请求
type UpdateProjectTemplateRequest struct {
	ID          int64  `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description" binding:"max=500"`
	Category    string `json:"category" binding:"max=50"`
	Icon        string `json:"icon" binding:"max=100"`
	Duration    int    `json:"duration" binding:"min=1"`
	IsPublic    bool   `json:"isPublic"`
	Status      int    `json:"status" binding:"min=0,max=1"`
}

// GetProjectTemplateListRequest 获取项目模板列表请求
type GetProjectTemplateListRequest struct {
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"pageSize" binding:"min=1,max=100"`
	Name     string `form:"name"`
	Category string `form:"category"`
	Status   *int   `form:"status"`
}

// CreateTemplateDependencyRequest 创建模板任务依赖请求
type CreateTemplateDependencyRequest struct {
	TemplateID     int64  `json:"templateId" binding:"required"`
	PredecessorID  int64  `json:"predecessorId" binding:"required"`
	SuccessorID    int64  `json:"successorId" binding:"required"`
	DependencyType string `json:"dependencyType" binding:"required,oneof=FS SS FF SF"`
	LagDays        int    `json:"lagDays"`
}
