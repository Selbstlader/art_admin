package request

// CreateCadGenerationReq 创建CAD生成任务请求
// Create CAD generation task request
type CreateCadGenerationReq struct {
	ProjectID      uint                    `json:"projectId" binding:"required"`      // 项目ID / Project ID
	DocumentIDs    []uint                  `json:"documentIds"`                       // 关联文档ID列表 / Document IDs
	GenerationType string                  `json:"generationType" binding:"required"` // 生成类型 / Generation type
	Prompt         string                  `json:"prompt"`                            // 生成提示词 / Generation prompt
	Parameters     *CadGenerationParamsReq `json:"parameters"`                        // 生成参数 / Generation parameters
}

// CadGenerationParamsReq CAD生成参数请求
// CAD generation parameters request
type CadGenerationParamsReq struct {
	Width        float64  `json:"width"`        // 宽度(米) / Width in meters
	Height       float64  `json:"height"`       // 高度(米) / Height in meters
	Scale        string   `json:"scale"`        // 比例 / Scale
	Style        string   `json:"style"`        // 设计风格 / Design style
	RoomTypes    []string `json:"roomTypes"`    // 房间类型列表 / Room types
	Requirements string   `json:"requirements"` // 特殊要求 / Special requirements
}

// CadGenerationListReq CAD生成任务列表请求
// CAD generation task list request
type CadGenerationListReq struct {
	ProjectID      uint   `form:"projectId"`                        // 项目ID / Project ID
	GenerationType string `form:"generationType"`                   // 生成类型 / Generation type
	Status         string `form:"status"`                           // 状态 / Status
	Current        int    `form:"current" binding:"required,min=1"` // 当前页 / Current page
	Size           int    `form:"size" binding:"required,min=1"`    // 每页数量 / Page size
}

// RetryGenerationReq 重试生成请求
// Retry generation request
type RetryGenerationReq struct {
	ID uint `json:"id" binding:"required"` // 任务ID / Task ID
}

// ConfirmCadFileReq 确认CAD文件请求
// Confirm CAD file request
type ConfirmCadFileReq struct {
	ID       uint   `json:"id" binding:"required"` // 任务ID / Task ID
	FileName string `json:"fileName"`              // 文件名 / File name
}
