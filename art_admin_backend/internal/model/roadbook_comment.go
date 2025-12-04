package model

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// RoadbookComment 评论模型
type RoadbookComment struct {
	ID         int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	RoadbookID int64          `gorm:"type:bigint unsigned;not null;index" json:"roadbookId"`
	UserID     int64          `gorm:"type:bigint unsigned;not null;index" json:"userId"`
	ParentID   int64          `gorm:"type:bigint unsigned;default:0;index" json:"parentId"`
	Content    string         `gorm:"type:varchar(500);not null" json:"content"`
	LikeCount  int            `gorm:"type:int unsigned;default:0" json:"likeCount"`
	CreatedAt  time.Time      `gorm:"type:timestamp;not null" json:"createdAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// 非数据库字段
	UserName string            `gorm:"-" json:"userName,omitempty"`
	Avatar   string            `gorm:"-" json:"avatar,omitempty"`
	Replies  []RoadbookComment `gorm:"-" json:"replies,omitempty"`
}

// TableName 表名
func (RoadbookComment) TableName() string {
	return "roadbook_comments"
}

// ToJSON 序列化为JSON
func (c *RoadbookComment) ToJSON() ([]byte, error) {
	return json.Marshal(c)
}

// FromJSON 从JSON反序列化
func (c *RoadbookComment) FromJSON(data []byte) error {
	return json.Unmarshal(data, c)
}

// ValidateCommentContent 验证评论内容
func ValidateCommentContent(content string) bool {
	// 检查是否为空或仅包含空白字符
	trimmed := ""
	for _, r := range content {
		if r != ' ' && r != '\t' && r != '\n' && r != '\r' {
			trimmed += string(r)
		}
	}
	if len(trimmed) == 0 {
		return false
	}
	// 检查长度是否超过500字符
	if len([]rune(content)) > 500 {
		return false
	}
	return true
}
