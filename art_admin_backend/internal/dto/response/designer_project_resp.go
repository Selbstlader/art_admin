package response

import "time"

// DesignerProjectResponse 设计师项目响应
// Designer project response
type DesignerProjectResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`        // 项目名称 / Project name
	Description string    `json:"description"` // 项目描述 / Project description
	Area        float64   `json:"area"`        // 面积(平方米) / Area in square meters
	Budget      float64   `json:"budget"`      // 预算 / Budget
	Style       string    `json:"style"`       // 设计风格 / Design style
	Status      string    `json:"status"`      // 状态 / Status
	UserID      uint      `json:"userId"`      // 所属用户ID / Owner user ID
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// DesignerProjectDetailResponse 设计师项目详情响应
// Designer project detail response
type DesignerProjectDetailResponse struct {
	ID             uint                       `json:"id"`
	Name           string                     `json:"name"`        // 项目名称 / Project name
	Description    string                     `json:"description"` // 项目描述 / Project description
	Area           float64                    `json:"area"`        // 面积(平方米) / Area in square meters
	Budget         float64                    `json:"budget"`      // 预算 / Budget
	Style          string                     `json:"style"`       // 设计风格 / Design style
	Status         string                     `json:"status"`      // 状态 / Status
	UserID         uint                       `json:"userId"`      // 所属用户ID / Owner user ID
	CreatedAt      time.Time                  `json:"createdAt"`
	UpdatedAt      time.Time                  `json:"updatedAt"`
	DocumentCount  int                        `json:"documentCount"`  // 文档数量 / Document count
	CadFileCount   int                        `json:"cadFileCount"`   // CAD文件数量 / CAD file count
	Documents      []ProjectDocumentBrief     `json:"documents"`      // 文档列表 / Document list
	CadFiles       []CadFileBrief             `json:"cadFiles"`       // CAD文件列表 / CAD file list
	CostEstimate   *CostEstimateBrief         `json:"costEstimate"`   // 成本估算 / Cost estimate
	CompareResults []DesignCompareResultBrief `json:"compareResults"` // 比对结果 / Compare results
}

// ProjectDocumentBrief 项目文档简要信息
// Project document brief info
type ProjectDocumentBrief struct {
	ID             uint      `json:"id"`
	FileName       string    `json:"fileName"`       // 文件名 / File name
	FileType       string    `json:"fileType"`       // 文件类型 / File type
	FileSize       int64     `json:"fileSize"`       // 文件大小 / File size
	AnalysisStatus string    `json:"analysisStatus"` // 分析状态 / Analysis status
	CreatedAt      time.Time `json:"createdAt"`
}

// CadFileBrief CAD文件简要信息
// CAD file brief info
type CadFileBrief struct {
	ID          uint      `json:"id"`
	FileName    string    `json:"fileName"`    // 文件名 / File name
	FileFormat  string    `json:"fileFormat"`  // 文件格式 / File format
	ParseStatus string    `json:"parseStatus"` // 解析状态 / Parse status
	LayerCount  int       `json:"layerCount"`  // 图层数量 / Layer count
	Has3D       bool      `json:"has3d"`       // 是否包含3D / Has 3D
	CreatedAt   time.Time `json:"createdAt"`
}

// CostEstimateBrief 成本估算简要信息
// Cost estimate brief info
type CostEstimateBrief struct {
	ID             uint      `json:"id"`
	MaterialCost   float64   `json:"materialCost"`   // 材料费 / Material cost
	LaborCost      float64   `json:"laborCost"`      // 人工费 / Labor cost
	EquipmentCost  float64   `json:"equipmentCost"`  // 设备费 / Equipment cost
	ManagementCost float64   `json:"managementCost"` // 管理费 / Management cost
	TotalCost      float64   `json:"totalCost"`      // 总计 / Total cost
	BudgetLimit    float64   `json:"budgetLimit"`    // 预算上限 / Budget limit
	UpdatedAt      time.Time `json:"updatedAt"`
}

// DesignCompareResultBrief 设计比对结果简要信息
// Design compare result brief info
type DesignCompareResultBrief struct {
	ID           uint      `json:"id"`
	OverallScore float64   `json:"overallScore"` // 整体匹配度 / Overall score
	CreatedAt    time.Time `json:"createdAt"`
}

// DesignerProjectListResponse 设计师项目列表响应
// Designer project list response
type DesignerProjectListResponse struct {
	Records []DesignerProjectResponse `json:"records"` // 记录列表 / Record list
	Current int                       `json:"current"` // 当前页码 / Current page
	Size    int                       `json:"size"`    // 每页条数 / Page size
	Total   int64                     `json:"total"`   // 总记录数 / Total count
}
