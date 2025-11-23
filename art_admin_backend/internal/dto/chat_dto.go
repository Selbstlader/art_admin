package dto

import "time"

// CreateChatRoomRequest 创建聊天室请求
type CreateChatRoomRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description" binding:"max=500"`
	Type        string `json:"type" binding:"required,oneof=public private"`
	MaxMembers  int    `json:"maxMembers" binding:"min=0"`
}

// UpdateChatRoomRequest 更新聊天室请求
type UpdateChatRoomRequest struct {
	ID          uint   `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description" binding:"max=500"`
	MaxMembers  int    `json:"maxMembers" binding:"min=0"`
	IsActive    *bool  `json:"isActive"`
}

// ChatRoomListRequest 聊天室列表请求
type ChatRoomListRequest struct {
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"pageSize" binding:"min=1,max=100"`
	Keyword  string `form:"keyword"`
	Type     string `form:"type" binding:"omitempty,oneof=public private"`
	IsActive *bool  `form:"isActive"`
}

// ChatRoomResponse 聊天室响应
type ChatRoomResponse struct {
	ID            uint       `json:"id"`
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	Type          string     `json:"type"`
	MaxMembers    int        `json:"maxMembers"`
	IsActive      bool       `json:"isActive"`
	CreatedBy     uint       `json:"createdBy"`
	MemberCount   int        `json:"memberCount"`
	OnlineCount   int        `json:"onlineCount"`
	LastMessageAt *time.Time `json:"lastMessageAt,omitempty"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

// SendMessageRequest 发送消息请求
type SendMessageRequest struct {
	RoomID      uint   `json:"roomId" binding:"required"`
	Content     string `json:"content" binding:"required,max=5000"`
	MessageType string `json:"messageType" binding:"required,oneof=text image file"`
	ReplyToID   *uint  `json:"replyToId"`
}

// ChatMessageResponse 聊天消息响应
type ChatMessageResponse struct {
	ID          uint                 `json:"id"`
	RoomID      uint                 `json:"roomId"`
	UserID      uint                 `json:"userId"`
	Username    string               `json:"username"`
	Content     string               `json:"content"`
	MessageType string               `json:"messageType"`
	ReplyToID   *uint                `json:"replyToId,omitempty"`
	ReplyTo     *ChatMessageResponse `json:"replyTo,omitempty"`
	IsRecalled  bool                 `json:"isRecalled"`
	CreatedAt   time.Time            `json:"createdAt"`
}

// MessageListRequest 消息列表请求
type MessageListRequest struct {
	RoomID   uint  `form:"roomId" binding:"required"`
	Page     int   `form:"page" binding:"min=1"`
	PageSize int   `form:"pageSize" binding:"min=1,max=100"`
	BeforeID *uint `form:"beforeId"` // 用于加载历史消息
	AfterID  *uint `form:"afterId"`  // 用于加载新消息
}

// JoinRoomRequest 加入聊天室请求
type JoinRoomRequest struct {
	RoomID uint `json:"roomId" binding:"required"`
}

// ChatRoomMemberResponse 聊天室成员响应
type ChatRoomMemberResponse struct {
	ID         uint       `json:"id"`
	RoomID     uint       `json:"roomId"`
	UserID     uint       `json:"userId"`
	Username   string     `json:"username"`
	Role       string     `json:"role"`
	IsMuted    bool       `json:"isMuted"`
	IsOnline   bool       `json:"isOnline"`
	LastReadAt *time.Time `json:"lastReadAt,omitempty"`
	JoinedAt   time.Time  `json:"joinedAt"`
}

// WebSocketMessage WebSocket消息结构
type WebSocketMessage struct {
	Type      string      `json:"type"` // message/join/leave/typing/online/error
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

// TypingNotification 正在输入通知
type TypingNotification struct {
	RoomID   uint   `json:"roomId"`
	UserID   uint   `json:"userId"`
	Username string `json:"username"`
	IsTyping bool   `json:"isTyping"`
}

// OnlineStatusNotification 在线状态通知
type OnlineStatusNotification struct {
	RoomID   uint   `json:"roomId"`
	UserID   uint   `json:"userId"`
	Username string `json:"username"`
	IsOnline bool   `json:"isOnline"`
}

// RecallMessageRequest 撤回消息请求
type RecallMessageRequest struct {
	MessageID uint `json:"messageId" binding:"required"`
}
