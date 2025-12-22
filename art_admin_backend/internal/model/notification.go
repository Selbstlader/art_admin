package model

import (
	"time"

	"gorm.io/gorm"
)

// NotificationRelatedType 通知关联类型
// Notification related type for routing
type NotificationRelatedType string

const (
	NotificationRelatedTypeVersionCompare   NotificationRelatedType = "version_compare"   // 版本对比
	NotificationRelatedTypeRender           NotificationRelatedType = "render"            // 效果图渲染
	NotificationRelatedTypeCadGeneration    NotificationRelatedType = "cad_generation"    // CAD生成
	NotificationRelatedTypeDocumentAnalysis NotificationRelatedType = "document_analysis" // 文档分析
	NotificationRelatedTypeDesignSuggestion NotificationRelatedType = "design_suggestion" // 设计建议
	NotificationRelatedTypeGeneral          NotificationRelatedType = "general"           // 通用通知
)

// Notification 用户通知
// User notification model
type Notification struct {
	ID          uint                    `gorm:"primaryKey" json:"id"`
	UserID      int64                   `gorm:"index;not null" json:"userId"`                 // 用户ID / User ID
	Title       string                  `gorm:"size:200;not null" json:"title"`               // 通知标题 / Title
	Content     string                  `gorm:"size:500" json:"content"`                      // 通知内容 / Content
	Type        string                  `gorm:"size:20;default:'notice'" json:"type"`         // 类型: notice/message/email / Type
	IsRead      bool                    `gorm:"default:false" json:"isRead"`                  // 是否已读 / Is read
	RelatedID   uint                    `gorm:"index" json:"relatedId"`                       // 关联ID（如效果图记录ID）/ Related ID
	RelatedType NotificationRelatedType `gorm:"size:50;default:'general'" json:"relatedType"` // 关联类型 / Related type for routing
	CreatedAt   time.Time               `json:"createdAt"`
	DeletedAt   gorm.DeletedAt          `gorm:"index" json:"-"`
}

// TableName 表名
func (Notification) TableName() string {
	return "notifications"
}
