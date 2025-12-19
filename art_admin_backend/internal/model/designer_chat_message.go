package model

import (
	"time"

	"gorm.io/gorm"
)

// DesignerChatMessage 设计师AI对话消息模型
// Designer AI chat message model
type DesignerChatMessage struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	SessionID    string         `gorm:"size:100;index;not null" json:"sessionId"` // 会话ID / Session ID
	ProjectID    uint           `gorm:"index" json:"projectId"`                   // 关联项目ID / Associated project ID
	UserID       uint           `gorm:"index;not null" json:"userId"`             // 用户ID / User ID
	Role         string         `gorm:"size:20;not null" json:"role"`             // 角色(user/assistant) / Role
	Content      string         `gorm:"type:text;not null" json:"content"`        // 消息内容 / Message content
	TokensUsed   int            `gorm:"default:0" json:"tokensUsed"`              // Token消耗 / Tokens used
	ResponseTime int64          `gorm:"default:0" json:"responseTime"`            // 响应时间(毫秒) / Response time in ms
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (DesignerChatMessage) TableName() string {
	return "designer_chat_messages"
}

// AIUsageLog AI使用日志模型
// AI usage log model for monitoring token consumption
// Requirements: 13.4 - 记录每次AI调用的token消耗和响应时间
type AIUsageLog struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"index;not null" json:"userId"`              // 用户ID / User ID
	ProjectID    uint      `gorm:"index" json:"projectId"`                    // 项目ID / Project ID
	ServiceType  string    `gorm:"size:50;index;not null" json:"serviceType"` // 服务类型 / Service type
	TokensUsed   int       `gorm:"default:0" json:"tokensUsed"`               // Token消耗 / Tokens used
	ResponseTime int64     `gorm:"default:0" json:"responseTime"`             // 响应时间(毫秒) / Response time in ms
	CreatedAt    time.Time `json:"createdAt"`
}

// TableName 表名
func (AIUsageLog) TableName() string {
	return "ai_usage_logs"
}
