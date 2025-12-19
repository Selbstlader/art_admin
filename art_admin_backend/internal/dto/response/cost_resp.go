package response

import "time"

// CostEstimateResponse 成本估算响应
// Cost estimate response
type CostEstimateResponse struct {
	ID             uint       `json:"id"`
	ProjectID      uint       `json:"projectId"`      // 项目ID / Project ID
	ProjectName    string     `json:"projectName"`    // 项目名称 / Project name
	MaterialCost   float64    `json:"materialCost"`   // 材料费 / Material cost
	LaborCost      float64    `json:"laborCost"`      // 人工费 / Labor cost
	EquipmentCost  float64    `json:"equipmentCost"`  // 设备费 / Equipment cost
	ManagementCost float64    `json:"managementCost"` // 管理费 / Management cost
	TotalCost      float64    `json:"totalCost"`      // 总计 / Total cost
	BudgetLimit    float64    `json:"budgetLimit"`    // 预算上限 / Budget limit
	BudgetExceeded bool       `json:"budgetExceeded"` // 是否超预算 / Whether budget exceeded
	ExceededAmount float64    `json:"exceededAmount"` // 超出金额 / Exceeded amount
	Items          []CostItem `json:"items"`          // 明细项 / Cost items
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}

// CostItem 成本明细项响应
// Cost item response
type CostItem struct {
	MaterialID   uint    `json:"materialId"`   // 材料ID / Material ID
	MaterialName string  `json:"materialName"` // 材料名称 / Material name
	Quantity     float64 `json:"quantity"`     // 数量 / Quantity
	UnitPrice    float64 `json:"unitPrice"`    // 单价 / Unit price
	TotalPrice   float64 `json:"totalPrice"`   // 总价 / Total price
	Category     string  `json:"category"`     // 分类 / Category
	Unit         string  `json:"unit"`         // 单位 / Unit
}

// CostCalculateResponse 成本计算响应
// Cost calculate response
type CostCalculateResponse struct {
	MaterialCost   float64    `json:"materialCost"`   // 材料费 / Material cost
	LaborCost      float64    `json:"laborCost"`      // 人工费 / Labor cost
	EquipmentCost  float64    `json:"equipmentCost"`  // 设备费 / Equipment cost
	ManagementCost float64    `json:"managementCost"` // 管理费 / Management cost
	TotalCost      float64    `json:"totalCost"`      // 总计 / Total cost
	BudgetLimit    float64    `json:"budgetLimit"`    // 预算上限 / Budget limit
	BudgetExceeded bool       `json:"budgetExceeded"` // 是否超预算 / Whether budget exceeded
	ExceededAmount float64    `json:"exceededAmount"` // 超出金额 / Exceeded amount
	BudgetWarning  string     `json:"budgetWarning"`  // 预算警告信息 / Budget warning message
	Items          []CostItem `json:"items"`          // 明细项 / Cost items
}

// CostReportExportResponse 成本报告导出响应
// Cost report export response
type CostReportExportResponse struct {
	FileURL  string `json:"fileUrl"`  // 文件下载URL / File download URL
	FileName string `json:"fileName"` // 文件名 / File name
	FileSize int64  `json:"fileSize"` // 文件大小(字节) / File size in bytes
}

// CostSummaryResponse 成本汇总响应
// Cost summary response
type CostSummaryResponse struct {
	TotalProjects     int     `json:"totalProjects"`     // 项目总数 / Total projects
	TotalBudget       float64 `json:"totalBudget"`       // 总预算 / Total budget
	TotalCost         float64 `json:"totalCost"`         // 总成本 / Total cost
	AvgCostPerProject float64 `json:"avgCostPerProject"` // 平均每项目成本 / Average cost per project
	OverBudgetCount   int     `json:"overBudgetCount"`   // 超预算项目数 / Over budget project count
}

// CostBreakdownResponse 成本分解响应
// Cost breakdown response
type CostBreakdownResponse struct {
	Category   string  `json:"category"`   // 分类名称 / Category name
	Amount     float64 `json:"amount"`     // 金额 / Amount
	Percentage float64 `json:"percentage"` // 占比 / Percentage
}
