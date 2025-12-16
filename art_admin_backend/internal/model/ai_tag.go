package model

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// AITag AI标签模型 - 用于管理Dify知识库配置和AI提示词
type AITag struct {
	ID                int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Name              string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"` // 标签名称
	Description       string         `gorm:"type:varchar(500)" json:"description"`               // 标签描述
	KnowledgeBaseID   string         `gorm:"type:varchar(100)" json:"knowledge_base_id"`         // Dify知识库ID（从API获取）
	KnowledgeBaseName string         `gorm:"type:varchar(200)" json:"knowledge_base_name"`       // 知识库名称（冗余存储）
	SystemPrompt      string         `gorm:"type:text;not null" json:"system_prompt"`            // 系统提示词
	ChatAPIKey        string         `gorm:"type:varchar(200)" json:"chat_api_key"`              // Dify Chat App API Key（手动输入）
	Status            int            `gorm:"type:tinyint;default:1;index" json:"status"`         // 1=active, 0=inactive
	CreatedAt         time.Time      `gorm:"autoCreateTime" json:"created_at"`                   // 创建时间
	UpdatedAt         time.Time      `gorm:"autoUpdateTime" json:"updated_at"`                   // 更新时间
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`                                     // 软删除时间
}

// TableName 表名
func (AITag) TableName() string {
	return "ai_tags"
}

// ToJSON serializes AITag to JSON bytes
func (t *AITag) ToJSON() ([]byte, error) {
	return json.Marshal(t)
}

// FromJSON deserializes JSON bytes to AITag
func (t *AITag) FromJSON(data []byte) error {
	return json.Unmarshal(data, t)
}
