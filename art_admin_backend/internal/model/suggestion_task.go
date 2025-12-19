package model

import (
	"time"

	"gorm.io/gorm"
)

// SuggestionTask 设计建议生成任务
// Design suggestion generation task model
type SuggestionTask struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	ProjectID      uint           `gorm:"index;not null" json:"projectId"`         // 关联项目ID / Associated project ID
	UserID         uint           `gorm:"index;not null" json:"userId"`            // 用户ID / User ID
	Category       string         `gorm:"size:50" json:"category"`                 // 指定类别 / Specified category
	Count          int            `gorm:"default:5" json:"count"`                  // 请求生成数量 / Requested count
	GeneratedCount int            `gorm:"default:0" json:"generatedCount"`         // 实际生成数量 / Actually generated count
	Status         string         `gorm:"size:20;default:'pending'" json:"status"` // 状态: pending/processing/completed/failed
	ErrorMessage   string         `gorm:"size:500" json:"errorMessage"`            // 错误信息 / Error message
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系 / Associations
	Project *DesignerProject `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
}

// TableName 表名
func (SuggestionTask) TableName() string {
	return "suggestion_tasks"
}
