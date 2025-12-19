package request

// CreateCostEstimateRequest 创建成本估算请求
// Create cost estimate request
type CreateCostEstimateRequest struct {
	ProjectID   uint       `json:"projectId" binding:"required"` // 项目ID / Project ID
	BudgetLimit float64    `json:"budgetLimit"`                  // 预算上限 / Budget limit
	Items       []CostItem `json:"items"`                        // 成本明细项 / Cost items
	LaborRate   float64    `json:"laborRate"`                    // 人工费率(每平米) / Labor rate per square meter
	EquipRate   float64    `json:"equipRate"`                    // 设备费率(每平米) / Equipment rate per square meter
	MgmtRate    float64    `json:"mgmtRate"`                     // 管理费率(占总成本比例) / Management rate (percentage of total)
}

// UpdateCostEstimateRequest 更新成本估算请求
// Update cost estimate request
type UpdateCostEstimateRequest struct {
	ID          uint       `json:"id" binding:"required"` // 成本估算ID / Cost estimate ID
	BudgetLimit float64    `json:"budgetLimit"`           // 预算上限 / Budget limit
	Items       []CostItem `json:"items"`                 // 成本明细项 / Cost items
	LaborRate   float64    `json:"laborRate"`             // 人工费率(每平米) / Labor rate per square meter
	EquipRate   float64    `json:"equipRate"`             // 设备费率(每平米) / Equipment rate per square meter
	MgmtRate    float64    `json:"mgmtRate"`              // 管理费率(占总成本比例) / Management rate (percentage of total)
}

// CostItem 成本明细项
// Cost item for detailed cost breakdown
type CostItem struct {
	MaterialID   uint    `json:"materialId"`   // 材料ID / Material ID
	MaterialName string  `json:"materialName"` // 材料名称 / Material name
	Quantity     float64 `json:"quantity"`     // 数量 / Quantity
	UnitPrice    float64 `json:"unitPrice"`    // 单价 / Unit price
	TotalPrice   float64 `json:"totalPrice"`   // 总价 / Total price
	Category     string  `json:"category"`     // 分类 / Category (material/labor/equipment/management)
	Unit         string  `json:"unit"`         // 单位 / Unit
}

// CalculateCostRequest 计算成本请求
// Calculate cost request
type CalculateCostRequest struct {
	ProjectID   uint       `json:"projectId" binding:"required"` // 项目ID / Project ID
	Area        float64    `json:"area"`                         // 面积(平方米) / Area in square meters
	Items       []CostItem `json:"items"`                        // 材料明细项 / Material items
	LaborRate   float64    `json:"laborRate"`                    // 人工费率(每平米) / Labor rate per square meter
	EquipRate   float64    `json:"equipRate"`                    // 设备费率(每平米) / Equipment rate per square meter
	MgmtRate    float64    `json:"mgmtRate"`                     // 管理费率(占总成本比例) / Management rate (percentage of total)
	BudgetLimit float64    `json:"budgetLimit"`                  // 预算上限 / Budget limit
}

// ExportCostReportRequest 导出成本报告请求
// Export cost report request
type ExportCostReportRequest struct {
	ProjectID uint   `json:"projectId" binding:"required"` // 项目ID / Project ID
	Format    string `json:"format" binding:"required"`    // 导出格式(excel/pdf) / Export format (excel/pdf)
}

// GetCostEstimateRequest 获取成本估算请求
// Get cost estimate request
type GetCostEstimateRequest struct {
	ProjectID uint `json:"projectId" form:"projectId" binding:"required"` // 项目ID / Project ID
}
