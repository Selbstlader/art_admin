package model

import (
	"time"

	"gorm.io/gorm"
)

// DesignSuggestion 设计建议模型
// Design suggestion model for AI-generated design recommendations
type DesignSuggestion struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	ProjectID       uint           `gorm:"index;not null" json:"projectId"`         // 关联项目ID / Associated project ID
	Content         string         `gorm:"type:text;not null" json:"content"`       // 建议内容 / Suggestion content
	ApplicableScene string         `gorm:"type:text" json:"applicableScene"`        // 适用场景 / Applicable scene
	CostImpact      string         `gorm:"size:200" json:"costImpact"`              // 成本影响 / Cost impact (e.g., "增加5%-10%", "节省约15%")
	Category        string         `gorm:"size:100;index" json:"category"`          // 建议类别 / Category (layout/material/style/function)
	Priority        int            `gorm:"default:0" json:"priority"`               // 优先级 / Priority (higher = more important)
	Status          string         `gorm:"size:50;default:'pending'" json:"status"` // 状态 / Status (pending/adopted/ignored)
	DetailInfo      string         `gorm:"type:json" json:"detailInfo"`             // 详细信息JSON / Detail info in JSON format
	ReferenceImages string         `gorm:"type:json" json:"referenceImages"`        // 参考图片URL列表 / Reference image URLs
	UserID          uint           `gorm:"index" json:"userId"`                     // 创建用户ID / Creator user ID
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系 / Associations
	Project *DesignerProject `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
}

// TableName 表名
// Table name
func (DesignSuggestion) TableName() string {
	return "design_suggestions"
}

// SuggestionPreference 建议偏好记录模型
// Suggestion preference record model for tracking user actions
type SuggestionPreference struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	SuggestionID uint           `gorm:"index;not null" json:"suggestionId"` // 关联建议ID / Associated suggestion ID
	UserID       uint           `gorm:"index;not null" json:"userId"`       // 用户ID / User ID
	Action       string         `gorm:"size:50;not null" json:"action"`     // 操作类型 / Action type (adopted/ignored/viewed)
	Feedback     string         `gorm:"type:text" json:"feedback"`          // 用户反馈 / User feedback
	CreatedAt    time.Time      `json:"createdAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系 / Associations
	Suggestion *DesignSuggestion `gorm:"foreignKey:SuggestionID" json:"suggestion,omitempty"`
}

// TableName 表名
// Table name
func (SuggestionPreference) TableName() string {
	return "suggestion_preferences"
}
