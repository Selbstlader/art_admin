package model

import (
	"time"

	"gorm.io/gorm"
)

// ChatRoomMember 聊天室成员模型
type ChatRoomMember struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	RoomID     uint           `gorm:"not null;index;comment:聊天室ID" json:"roomId"`
	UserID     uint           `gorm:"not null;index;comment:用户ID" json:"userId"`
	Username   string         `gorm:"size:50;not null;comment:用户名" json:"username"`
	Role       string         `gorm:"size:20;not null;default:member;comment:角色(owner/admin/member)" json:"role"`
	IsMuted    bool           `gorm:"default:false;comment:是否被禁言" json:"isMuted"`
	LastReadAt *time.Time     `gorm:"comment:最后阅读时间" json:"lastReadAt,omitempty"`
	JoinedAt   time.Time      `gorm:"not null;comment:加入时间" json:"joinedAt"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Room *ChatRoom `gorm:"foreignKey:RoomID" json:"room,omitempty"`
}

// TableName 指定表名
func (ChatRoomMember) TableName() string {
	return "chat_room_members"
}
