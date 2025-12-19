package request

// DesignerProjectListRequest 设计师项目列表查询请求
// Designer project list query request
type DesignerProjectListRequest struct {
	Current   int     `form:"current" binding:"required,min=1" example:"1"`
	Size      int     `form:"size" binding:"required,min=1,max=100" example:"10"`
	Name      string  `form:"name" example:"办公室设计"`           // 项目名称模糊搜索 / Project name fuzzy search
	Status    string  `form:"status" example:"draft"`         // 状态筛选 / Status filter
	Style     string  `form:"style" example:"现代简约"`           // 风格筛选 / Style filter
	MinBudget float64 `form:"minBudget" example:"100000"`     // 最小预算 / Minimum budget
	MaxBudget float64 `form:"maxBudget" example:"500000"`     // 最大预算 / Maximum budget
	StartDate string  `form:"startDate" example:"2024-01-01"` // 开始日期 / Start date
	EndDate   string  `form:"endDate" example:"2024-12-31"`   // 结束日期 / End date
	Keyword   string  `form:"keyword" example:"办公"`           // 关键字搜索(名称/描述) / Keyword search
}

// CreateDesignerProjectRequest 创建设计师项目请求
// Create designer project request
type CreateDesignerProjectRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=200"`                                 // 项目名称 / Project name
	Description string  `json:"description" binding:"max=2000"`                                        // 项目描述 / Project description
	Area        float64 `json:"area" binding:"min=0"`                                                  // 面积(平方米) / Area in square meters
	Budget      float64 `json:"budget" binding:"min=0"`                                                // 预算 / Budget
	Style       string  `json:"style" binding:"max=100"`                                               // 设计风格 / Design style
	Status      string  `json:"status" binding:"omitempty,oneof=draft in_progress completed archived"` // 状态 / Status
}

// UpdateDesignerProjectRequest 更新设计师项目请求
// Update designer project request
type UpdateDesignerProjectRequest struct {
	ID          uint    `json:"id" binding:"required"`                                                 // 项目ID / Project ID
	Name        string  `json:"name" binding:"required,min=1,max=200"`                                 // 项目名称 / Project name
	Description string  `json:"description" binding:"max=2000"`                                        // 项目描述 / Project description
	Area        float64 `json:"area" binding:"min=0"`                                                  // 面积(平方米) / Area in square meters
	Budget      float64 `json:"budget" binding:"min=0"`                                                // 预算 / Budget
	Style       string  `json:"style" binding:"max=100"`                                               // 设计风格 / Design style
	Status      string  `json:"status" binding:"omitempty,oneof=draft in_progress completed archived"` // 状态 / Status
}

// DeleteDesignerProjectRequest 删除设计师项目请求
// Delete designer project request
type DeleteDesignerProjectRequest struct {
	ID uint `json:"id" binding:"required"` // 项目ID / Project ID
}

// BatchDeleteDesignerProjectRequest 批量删除设计师项目请求
// Batch delete designer projects request
type BatchDeleteDesignerProjectRequest struct {
	IDs []uint `json:"ids" binding:"required,min=1"` // 项目ID列表 / Project ID list
}
