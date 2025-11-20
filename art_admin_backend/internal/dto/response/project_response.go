package response

import "time"

// ProjectResponse 项目响应
type ProjectResponse struct {
	ID              int64      `json:"id"`
	TemplateID      int64      `json:"templateId"`
	Name            string     `json:"name"`
	Code            string     `json:"code"`
	Description     string     `json:"description"`
	Category        string     `json:"category"`
	Priority        int        `json:"priority"`
	Status          int        `json:"status"`
	Progress        int        `json:"progress"`
	StartDate       time.Time  `json:"startDate"`
	EndDate         time.Time  `json:"endDate"`
	ActualStartDate *time.Time `json:"actualStartDate"`
	ActualEndDate   *time.Time `json:"actualEndDate"`
	ManagerID       int64      `json:"managerId"`
	ManagerName     string     `json:"managerName"`
	CreatorID       int64      `json:"creatorId"`
	CreatorName     string     `json:"creatorName"`
	Budget          float64    `json:"budget"`
	ActualCost      float64    `json:"actualCost"`
	Remark          string     `json:"remark"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}

// ProjectDetailResponse 项目详情响应
type ProjectDetailResponse struct {
	ProjectResponse
	Members            []ProjectMemberResponse `json:"members"`
	TaskCount          int                     `json:"taskCount"`
	CompletedTaskCount int                     `json:"completedTaskCount"`
	OverdueTaskCount   int                     `json:"overdueTaskCount"`
}

// ProjectMemberResponse 项目成员响应
type ProjectMemberResponse struct {
	ID        int64     `json:"id"`
	ProjectID int64     `json:"projectId"`
	UserID    int64     `json:"userId"`
	UserName  string    `json:"userName"`
	Role      string    `json:"role"`
	JoinDate  time.Time `json:"joinDate"`
	CreatedAt time.Time `json:"createdAt"`
}

// ProjectStatisticsResponse 项目统计响应
type ProjectStatisticsResponse struct {
	TotalProjects     int `json:"totalProjects"`
	OngoingProjects   int `json:"ongoingProjects"`
	CompletedProjects int `json:"completedProjects"`
	OverdueProjects   int `json:"overdueProjects"`
}
