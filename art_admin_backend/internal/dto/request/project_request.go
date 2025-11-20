package request

import "time"

// CreateProjectRequest 创建项目请求
type CreateProjectRequest struct {
	TemplateID  int64              `json:"templateId"`
	Name        string             `json:"name" binding:"required,max=100"`
	Description string             `json:"description"`
	Category    string             `json:"category" binding:"max=50"`
	Priority    int                `json:"priority" binding:"min=1,max=3"`
	Status      int                `json:"status" binding:"min=1,max=5"`
	StartDate   time.Time          `json:"startDate" binding:"required"`
	EndDate     time.Time          `json:"endDate" binding:"required"`
	ManagerID   int64              `json:"managerId" binding:"required"`
	Budget      float64            `json:"budget" binding:"min=0"`
	Remark      string             `json:"remark" binding:"max=500"`
	Members     []ProjectMemberReq `json:"members"`
}

// ProjectMemberReq 项目成员请求
type ProjectMemberReq struct {
	UserID   int64     `json:"userId" binding:"required"`
	Role     string    `json:"role" binding:"max=50"`
	JoinDate time.Time `json:"joinDate" binding:"required"`
}

// UpdateProjectRequest 更新项目请求
type UpdateProjectRequest struct {
	ID              int64      `json:"id" binding:"required"`
	Name            string     `json:"name" binding:"required,max=100"`
	Description     string     `json:"description"`
	Category        string     `json:"category" binding:"max=50"`
	Priority        int        `json:"priority" binding:"min=1,max=3"`
	Status          int        `json:"status" binding:"min=1,max=5"`
	Progress        int        `json:"progress" binding:"min=0,max=100"`
	StartDate       time.Time  `json:"startDate" binding:"required"`
	EndDate         time.Time  `json:"endDate" binding:"required"`
	ActualStartDate *time.Time `json:"actualStartDate"`
	ActualEndDate   *time.Time `json:"actualEndDate"`
	ManagerID       int64      `json:"managerId" binding:"required"`
	Budget          float64    `json:"budget" binding:"min=0"`
	ActualCost      float64    `json:"actualCost" binding:"min=0"`
	Remark          string     `json:"remark" binding:"max=500"`
}

// GetProjectListRequest 获取项目列表请求
type GetProjectListRequest struct {
	Page      int        `form:"page" binding:"min=1"`
	PageSize  int        `form:"pageSize" binding:"min=1,max=100"`
	Name      string     `form:"name"`
	Code      string     `form:"code"`
	Category  string     `form:"category"`
	Status    *int       `form:"status"`
	ManagerID *int64     `form:"managerId"`
	StartDate *time.Time `form:"startDate"`
	EndDate   *time.Time `form:"endDate"`
}

// CreateProjectFromTemplateRequest 从模板创建项目请求
type CreateProjectFromTemplateRequest struct {
	TemplateID  int64              `json:"templateId" binding:"required"`
	Name        string             `json:"name" binding:"required,max=100"`
	Description string             `json:"description"`
	StartDate   time.Time          `json:"startDate" binding:"required"`
	ManagerID   int64              `json:"managerId" binding:"required"`
	Budget      float64            `json:"budget" binding:"min=0"`
	Members     []ProjectMemberReq `json:"members"`
}
