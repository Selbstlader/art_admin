package model

import (
	"time"

	"gorm.io/gorm"
)

// CostEstimate 成本估算模型
// Cost estimate model for project budget calculation
type CostEstimate struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	ProjectID      uint           `gorm:"not null;index:idx_cost_estimates_project_id,unique" json:"projectId"` // 关联项目ID(唯一) / Associated project ID (unique)
	MaterialCost   float64        `gorm:"default:0" json:"materialCost"`                                        // 材料费 / Material cost
	LaborCost      float64        `gorm:"default:0" json:"laborCost"`                                           // 人工费 / Labor cost
	EquipmentCost  float64        `gorm:"default:0" json:"equipmentCost"`                                       // 设备费 / Equipment cost
	ManagementCost float64        `gorm:"default:0" json:"managementCost"`                                      // 管理费 / Management cost
	TotalCost      float64        `gorm:"default:0" json:"totalCost"`                                           // 总计 / Total cost
	BudgetLimit    float64        `gorm:"default:0" json:"budgetLimit"`                                         // 预算上限 / Budget limit
	Items          string         `gorm:"type:json" json:"items"`                                               // 明细项(JSON) / Cost items in JSON format
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系 / Associations
	Project *DesignerProject `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
}

// TableName 表名
func (CostEstimate) TableName() string {
	return "cost_estimates"
}

// CostItem 成本明细项结构
// Cost item structure for detailed cost breakdown
type CostItem struct {
	MaterialID   uint    `json:"materialId"`   // 材料ID / Material ID
	MaterialName string  `json:"materialName"` // 材料名称 / Material name
	Quantity     float64 `json:"quantity"`     // 数量 / Quantity
	UnitPrice    float64 `json:"unitPrice"`    // 单价 / Unit price
	TotalPrice   float64 `json:"totalPrice"`   // 总价 / Total price
	Category     string  `json:"category"`     // 分类 / Category (material/labor/equipment/management)
}
