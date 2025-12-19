package response

import "time"

// ProjectMaterialItem 项目材料项响应
// Project material item response
type ProjectMaterialItem struct {
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

// ProjectMaterialListResponse 项目材料清单响应
// Project material list response
type ProjectMaterialListResponse struct {
	ProjectID    uint                  `json:"projectId"`
	ProjectName  string                `json:"projectName"`
	Items        []ProjectMaterialItem `json:"items"`
	MaterialCost float64               `json:"materialCost"`
	ItemCount    int                   `json:"itemCount"`
}

// CustomCostItem 自定义费用项响应
// Custom cost item response
type CustomCostItem struct {
	Name   string  `json:"name"`   // 费用名称 / Cost name
	Amount float64 `json:"amount"` // 金额(支持负数) / Amount (supports negative)
}

// ProjectCostSummary 项目成本汇总响应
// Project cost summary response
type ProjectCostSummary struct {
	ProjectID       uint             `json:"projectId"`
	MaterialCost    float64          `json:"materialCost"`
	LaborCost       float64          `json:"laborCost"`
	EquipmentCost   float64          `json:"equipmentCost"`
	ManagementCost  float64          `json:"managementCost"`
	CustomCosts     []CustomCostItem `json:"customCosts"`     // 自定义费用项 / Custom cost items
	CustomCostTotal float64          `json:"customCostTotal"` // 自定义费用小计 / Custom cost subtotal
	TotalCost       float64          `json:"totalCost"`
	BudgetLimit     float64          `json:"budgetLimit"`
	BudgetExceeded  bool             `json:"budgetExceeded"`
	ExceededAmount  float64          `json:"exceededAmount"`
}

// ExportProjectMaterialsResponse 导出项目材料清单响应
// Export project materials response
type ExportProjectMaterialsResponse struct {
	FileUrl  string `json:"fileUrl"`
	FileName string `json:"fileName"`
}
