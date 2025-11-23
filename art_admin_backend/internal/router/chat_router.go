package router

import (
	"art_admin_backend/internal/api"
	"art_admin_backend/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

// SetupChatRoutes 设置聊天室路由
func SetupChatRoutes(r *gin.Engine, api *api.ChatAPI) {
	chatGroup := r.Group("/api/chat")
	chatGroup.Use(middleware.JWTAuth())
	{
		// 聊天室管理
		chatGroup.POST("/room", api.CreateRoom)
		chatGroup.PUT("/room", api.UpdateRoom)
		chatGroup.DELETE("/room/:id", api.DeleteRoom)
		chatGroup.GET("/room/:id", api.GetRoom)
		chatGroup.GET("/room/list", api.GetRoomList)

		// 成员管理
		chatGroup.POST("/room/join", api.JoinRoom)
		chatGroup.POST("/room/:id/leave", api.LeaveRoom)
		chatGroup.GET("/room/:id/members", api.GetRoomMembers)
		chatGroup.GET("/room/:id/online-users", api.GetOnlineUsers)
		chatGroup.GET("/my-rooms", api.GetUserRooms)

		// 消息管理
		chatGroup.POST("/message", api.SendMessage)
		chatGroup.GET("/message/list", api.GetMessages)
		chatGroup.POST("/message/recall", api.RecallMessage)

		// WebSocket 连接
		chatGroup.GET("/ws", api.HandleWebSocket)
	}
}
