package request

import "time"

// ProjectMaterialItem 项目材料项请求
// Project material item request
type ProjectMaterialItem struct {
	MaterialID    *uint   `json:"materialId"`              // 关联材料库ID(可为空) / Associated material library ID (nullable)
	Name          string  `json:"name" binding:"required"` // 材料名称 / Material name
	Category      string  `json:"category"`                // 分类 / Category
	Specification string  `json:"specification"`           // 规格 / Specification
	Unit          string  `json:"unit" binding:"required"` // 单位 / Unit
	UnitPrice     float64 `json:"unitPrice"`               // 单价 / Unit price
	Quantity      float64 `json:"quantity"`                // 数量 / Quantity
	TotalPrice    float64 `json:"totalPrice"`              // 小计 / Total price
	Brand         string  `json:"brand"`                   // 品牌 / Brand
	Supplier      string  `json:"supplier"`                // 供应商 / Supplier
	Remark        string  `json:"remark"`                  // 备注 / Remark
}

// SaveProjectMaterialRequest 保存项目材料清单请求
// Save project material list request
type SaveProjectMaterialRequest struct {
	ProjectID uint                  `json:"projectId" binding:"required"` // 项目ID / Project ID
	Items     []ProjectMaterialItem `json:"items"`                        // 材料清单 / Material list
}

// CustomCostItem 自定义费用项
// Custom cost item
type CustomCostItem struct {
	Name   string  `json:"name"`   // 费用名称 / Cost name
	Amount float64 `json:"amount"` // 金额(支持负数) / Amount (supports negative)
}

// SaveProjectCostRequest 保存项目成本配置请求
// Save project cost config request
type SaveProjectCostRequest struct {
	ProjectID      uint             `json:"projectId" binding:"required"` // 项目ID / Project ID
	BudgetLimit    float64          `json:"budgetLimit"`                  // 预算上限 / Budget limit
	LaborCost      float64          `json:"laborCost"`                    // 人工费 / Labor cost
	EquipmentCost  float64          `json:"equipmentCost"`                // 设备费 / Equipment cost
	ManagementCost float64          `json:"managementCost"`               // 管理费 / Management cost
	CustomCosts    []CustomCostItem `json:"customCosts"`                  // 自定义费用项 / Custom cost items
}

// ImportMaterialsRequest 导入材料请求
// Import materials request
type ImportMaterialsRequest struct {
	MaterialIDs     []uint  `json:"materialIds" binding:"required"` // 材料ID列表 / Material ID list
	DefaultQuantity float64 `json:"defaultQuantity"`                // 默认数量 / Default quantity
}

// ExportProjectMaterialsRequest 导出项目材料清单请求
// Export project materials request
type ExportProjectMaterialsRequest struct {
	Format string `json:"format"` // 导出格式(excel/pdf) / Export format
}

// ProjectMaterialItemResponse 项目材料项响应(用于内部转换)
// Project material item response (for internal conversion)
type ProjectMaterialItemResponse struct {
	ID            uint      `json:"id"`
	ProjectID     uint      `json:"projectId"`
	MaterialID    *uint     `json:"materialId"`
	Name          string    `json:"name"`
	Category      string    `json:"category"`
	Specification string    `json:"specification"`
	Unit          string    `json:"unit"`
	UnitPrice     float64   `json:"unitPrice"`
	Quantity      float64   `json:"quantity"`
	TotalPrice    float64   `json:"totalPrice"`
	Brand         string    `json:"brand"`
	Supplier      string    `json:"supplier"`
	Remark        string    `json:"remark"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
