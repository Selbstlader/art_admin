package response

import "time"

// ProjectTemplateResponse 项目模板响应
type ProjectTemplateResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Icon        string    `json:"icon"`
	Duration    int       `json:"duration"`
	IsPublic    bool      `json:"isPublic"`
	CreatorID   int64     `json:"creatorId"`
	CreatorName string    `json:"creatorName"`
	UseCount    int       `json:"useCount"`
	Status      int       `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// ProjectTemplateDetailResponse 项目模板详情响应
type ProjectTemplateDetailResponse struct {
	ProjectTemplateResponse
	Tasks        []TemplateTaskResponse       `json:"tasks"`
	Dependencies []TemplateDependencyResponse `json:"dependencies"`
}

// TemplateTaskResponse 模板任务响应
type TemplateTaskResponse struct {
	ID          int64     `json:"id"`
	TemplateID  int64     `json:"templateId"`
	ParentID    int64     `json:"parentId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartDay    int       `json:"startDay"`
	Duration    int       `json:"duration"`
	IsMilestone bool      `json:"isMilestone"`
	Priority    int       `json:"priority"`
	Sort        int       `json:"sort"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// TemplateDependencyResponse 模板依赖响应
type TemplateDependencyResponse struct {
	ID             int64     `json:"id"`
	TemplateID     int64     `json:"templateId"`
	PredecessorID  int64     `json:"predecessorId"`
	SuccessorID    int64     `json:"successorId"`
	DependencyType string    `json:"dependencyType"`
	LagDays        int       `json:"lagDays"`
	CreatedAt      time.Time `json:"createdAt"`
}
