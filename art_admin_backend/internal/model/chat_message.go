package model

import (
	"time"

	"gorm.io/gorm"
)

// ChatMessage 聊天消息模型
type ChatMessage struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	RoomID      uint           `gorm:"not null;index;comment:聊天室ID" json:"roomId"`
	UserID      uint           `gorm:"not null;index;comment:发送者ID" json:"userId"`
	Username    string         `gorm:"size:50;not null;comment:发送者用户名" json:"username"`
	Content     string         `gorm:"type:text;not null;comment:消息内容" json:"content"`
	MessageType string         `gorm:"size:20;not null;default:text;comment:消息类型(text/image/file/system)" json:"messageType"`
	ReplyToID   *uint          `gorm:"index;comment:回复的消息ID" json:"replyToId,omitempty"`
	IsRecalled  bool           `gorm:"default:false;comment:是否已撤回" json:"isRecalled"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Room    *ChatRoom    `gorm:"foreignKey:RoomID" json:"room,omitempty"`
	ReplyTo *ChatMessage `gorm:"foreignKey:ReplyToID" json:"replyTo,omitempty"`
}

// TableName 指定表名
func (ChatMessage) TableName() string {
	return "chat_messages"
}
