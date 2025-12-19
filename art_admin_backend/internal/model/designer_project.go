package model

import (
	"time"

	"gorm.io/gorm"
)

// DesignerProject 工装设计项目模型
// Designer project model for tooling design assistant system
type DesignerProject struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:200;not null" json:"name"`               // 项目名称 / Project name
	Description string         `gorm:"type:text" json:"description"`                // 项目描述 / Project description
	Area        float64        `gorm:"default:0" json:"area"`                       // 面积(平方米) / Area in square meters
	Budget      float64        `gorm:"default:0" json:"budget"`                     // 预算 / Budget
	Style       string         `gorm:"size:100" json:"style"`                       // 设计风格 / Design style
	Status      string         `gorm:"size:50;default:'draft';index" json:"status"` // 状态 / Status (draft/in_progress/completed/archived)
	UserID      uint           `gorm:"index;not null" json:"userId"`                // 所属用户ID / Owner user ID
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系 / Associations
	Documents      []ProjectDocument     `gorm:"foreignKey:ProjectID" json:"documents,omitempty"`
	CadFiles       []CadFile             `gorm:"foreignKey:ProjectID" json:"cadFiles,omitempty"`
	CompareResults []DesignCompareResult `gorm:"foreignKey:ProjectID" json:"compareResults,omitempty"`
	CostEstimate   *CostEstimate         `gorm:"foreignKey:ProjectID" json:"costEstimate,omitempty"`
}

// TableName 表名
func (DesignerProject) TableName() string {
	return "designer_projects"
}
