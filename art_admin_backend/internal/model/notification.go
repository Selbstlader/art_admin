package model

import (
	"time"

	"gorm.io/gorm"
)

// Notification 用户通知
// User notification model
type Notification struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    int64          `gorm:"index;not null" json:"userId"`         // 用户ID / User ID
	Title     string         `gorm:"size:200;not null" json:"title"`       // 通知标题 / Title
	Content   string         `gorm:"size:500" json:"content"`              // 通知内容 / Content
	Type      string         `gorm:"size:20;default:'notice'" json:"type"` // 类型: notice/message/email / Type
	IsRead    bool           `gorm:"default:false" json:"isRead"`          // 是否已读 / Is read
	RelatedID uint           `gorm:"index" json:"relatedId"`               // 关联ID（如效果图记录ID）/ Related ID
	CreatedAt time.Time      `json:"createdAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (Notification) TableName() string {
	return "notifications"
}
