package model

import (
	"time"

	"gorm.io/gorm"
)

// Material 材料模型
// Material model for material library management
type Material struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	Name             string         `gorm:"size:200;not null;index" json:"name"`          // 材料名称 / Material name
	Category         string         `gorm:"size:100;index" json:"category"`               // 分类 / Category
	Specification    string         `gorm:"size:200" json:"specification"`                // 规格 / Specification
	Unit             string         `gorm:"size:50" json:"unit"`                          // 单位 / Unit (m²/m/个/kg)
	UnitPrice        float64        `gorm:"default:0" json:"unitPrice"`                   // 单价 / Unit price
	Brand            string         `gorm:"size:100;index" json:"brand"`                  // 品牌 / Brand
	Supplier         string         `gorm:"size:200" json:"supplier"`                     // 供应商 / Supplier
	Description      string         `gorm:"type:text" json:"description"`                 // 描述 / Description
	ImageURL         string         `gorm:"size:500" json:"imageUrl"`                     // 图片URL / Image URL
	ApplicableScenes string         `gorm:"type:json" json:"applicableScenes"`            // 适用场景(JSON) / Applicable scenes in JSON format
	Status           string         `gorm:"size:50;default:'active';index" json:"status"` // 状态 / Status (active/inactive)
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (Material) TableName() string {
	return "materials"
}
