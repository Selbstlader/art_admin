package model

import (
	"time"

	"gorm.io/gorm"
)

// ChatRoom 聊天室模型
type ChatRoom struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"size:100;not null;comment:聊天室名称" json:"name"`
	Description string         `gorm:"size:500;comment:聊天室描述" json:"description"`
	Type        string         `gorm:"size:20;not null;default:public;comment:聊天室类型(public/private)" json:"type"`
	MaxMembers  int            `gorm:"default:0;comment:最大成员数,0表示无限制" json:"maxMembers"`
	IsActive    bool           `gorm:"default:true;comment:是否启用" json:"isActive"`
	CreatedBy   uint           `gorm:"not null;comment:创建者ID" json:"createdBy"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (ChatRoom) TableName() string {
	return "chat_rooms"
}
