package model

import (
	"time"

	"gorm.io/gorm"
)

// ProjectMaterial 项目材料清单模型
// Project material list model for tracking materials used in a project
type ProjectMaterial struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	ProjectID     uint           `gorm:"index;not null" json:"projectId"` // 关联项目ID / Associated project ID
	MaterialID    *uint          `gorm:"index" json:"materialId"`         // 关联材料库ID(可为空) / Associated material library ID (nullable)
	Name          string         `gorm:"size:200;not null" json:"name"`   // 材料名称 / Material name
	Category      string         `gorm:"size:100" json:"category"`        // 分类 / Category
	Specification string         `gorm:"size:200" json:"specification"`   // 规格 / Specification
	Unit          string         `gorm:"size:50;not null" json:"unit"`    // 单位 / Unit
	UnitPrice     float64        `gorm:"default:0" json:"unitPrice"`      // 单价 / Unit price
	Quantity      float64        `gorm:"default:0" json:"quantity"`       // 数量 / Quantity
	TotalPrice    float64        `gorm:"default:0" json:"totalPrice"`     // 小计 / Total price
	Brand         string         `gorm:"size:100" json:"brand"`           // 品牌 / Brand
	Supplier      string         `gorm:"size:200" json:"supplier"`        // 供应商 / Supplier
	Remark        string         `gorm:"size:500" json:"remark"`          // 备注 / Remark
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系 / Associations
	Project  *DesignerProject `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	Material *Material        `gorm:"foreignKey:MaterialID" json:"material,omitempty"`
}

// TableName 表名
func (ProjectMaterial) TableName() string {
	return "project_materials"
}

// BeforeSave 保存前计算小计
// Calculate total price before save
func (pm *ProjectMaterial) BeforeSave(tx *gorm.DB) error {
	pm.TotalPrice = pm.UnitPrice * pm.Quantity
	return nil
}
