package api

import (
	"art_admin_backend/internal/api/middleware"
	"art_admin_backend/internal/dto"
	ws "art_admin_backend/internal/pkg/websocket"
	"art_admin_backend/internal/service"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// 生产环境应该检查 Origin
		return true
	},
}

// ChatAPI 聊天室 API
type ChatAPI struct {
	chatService service.ChatService
	hub         *ws.Hub
}

// NewChatAPI 创建聊天室 API
func NewChatAPI(chatService service.ChatService, hub *ws.Hub) *ChatAPI {
	return &ChatAPI{
		chatService: chatService,
		hub:         hub,
	}
}

// CreateRoom 创建聊天室
// @Summary 创建聊天室
// @Tags Chat
// @Accept json
// @Produce json
// @Param request body dto.CreateChatRoomRequest true "创建聊天室请求"
// @Success 200 {object} dto.ChatRoomResponse
// @Router /api/chat/room [post]
func (a *ChatAPI) CreateRoom(c *gin.Context) {
	var req dto.CreateChatRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	userID := middleware.GetUserID(c)
	room, err := a.chatService.CreateRoom(&req, uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "创建成功", "data": room})
}

// UpdateRoom 更新聊天室
// @Summary 更新聊天室
// @Tags Chat
// @Accept json
// @Produce json
// @Param request body dto.UpdateChatRoomRequest true "更新聊天室请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/chat/room [put]
func (a *ChatAPI) UpdateRoom(c *gin.Context) {
	var req dto.UpdateChatRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	if err := a.chatService.UpdateRoom(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功"})
}

// DeleteRoom 删除聊天室
// @Summary 删除聊天室
// @Tags Chat
// @Param id path int true "聊天室ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/chat/room/{id} [delete]
func (a *ChatAPI) DeleteRoom(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if err := a.chatService.DeleteRoom(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

// GetRoom 获取聊天室详情
// @Summary 获取聊天室详情
// @Tags Chat
// @Param id path int true "聊天室ID"
// @Success 200 {object} dto.ChatRoomResponse
// @Router /api/chat/room/{id} [get]
func (a *ChatAPI) GetRoom(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	room, err := a.chatService.GetRoomByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "聊天室不存在"})
		return
	}

	// 更新在线人数
	room.OnlineCount = a.hub.GetRoomOnlineCount(room.ID)

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": room})
}

// GetRoomList 获取聊天室列表
// @Summary 获取聊天室列表
// @Tags Chat
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param keyword query string false "关键词"
// @Param type query string false "类型"
// @Param isActive query bool false "是否启用"
// @Success 200 {object} map[string]interface{}
// @Router /api/chat/room/list [get]
func (a *ChatAPI) GetRoomList(c *gin.Context) {
	var req dto.ChatRoomListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	rooms, total, err := a.chatService.GetRoomList(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败: " + err.Error()})
		return
	}

	// 更新在线人数
	for _, room := range rooms {
		room.OnlineCount = a.hub.GetRoomOnlineCount(room.ID)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":     200,
		"data":     rooms,
		"total":    total,
		"page":     req.Page,
		"pageSize": req.PageSize,
	})
}

// JoinRoom 加入聊天室
// @Summary 加入聊天室
// @Tags Chat
// @Accept json
// @Produce json
// @Param request body dto.JoinRoomRequest true "加入聊天室请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/chat/room/join [post]
func (a *ChatAPI) JoinRoom(c *gin.Context) {
	var req dto.JoinRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	userID := middleware.GetUserID(c)
	username := middleware.GetUsername(c)

	if err := a.chatService.JoinRoom(req.RoomID, uint(userID), username); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "加入成功"})
}

// LeaveRoom 离开聊天室
// @Summary 离开聊天室
// @Tags Chat
// @Param id path int true "聊天室ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/chat/room/{id}/leave [post]
func (a *ChatAPI) LeaveRoom(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	userID := middleware.GetUserID(c)
	if err := a.chatService.LeaveRoom(uint(id), uint(userID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "离开失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "离开成功"})
}

// GetRoomMembers 获取聊天室成员列表
// @Summary 获取聊天室成员列表
// @Tags Chat
// @Param id path int true "聊天室ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/chat/room/{id}/members [get]
func (a *ChatAPI) GetRoomMembers(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	members, err := a.chatService.GetRoomMembers(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败: " + err.Error()})
		return
	}

	// 更新在线状态
	for _, member := range members {
		member.IsOnline = a.hub.IsUserOnline(member.UserID)
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": members})
}

// GetUserRooms 获取用户加入的聊天室列表
// @Summary 获取用户加入的聊天室列表
// @Tags Chat
// @Success 200 {object} map[string]interface{}
// @Router /api/chat/my-rooms [get]
func (a *ChatAPI) GetUserRooms(c *gin.Context) {
	userID := middleware.GetUserID(c)
	rooms, err := a.chatService.GetUserRooms(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败: " + err.Error()})
		return
	}

	// 更新在线人数
	for _, room := range rooms {
		room.OnlineCount = a.hub.GetRoomOnlineCount(room.ID)
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": rooms})
}

// SendMessage 发送消息
// @Summary 发送消息
// @Tags Chat
// @Accept json
// @Produce json
// @Param request body dto.SendMessageRequest true "发送消息请求"
// @Success 200 {object} dto.ChatMessageResponse
// @Router /api/chat/message [post]
func (a *ChatAPI) SendMessage(c *gin.Context) {
	var req dto.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	userID := middleware.GetUserID(c)
	username := middleware.GetUsername(c)

	message, err := a.chatService.SendMessage(&req, uint(userID), username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	// 先返回 HTTP 响应
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "发送成功", "data": message})

	// 异步通过 WebSocket 广播消息，避免阻塞 HTTP 响应
	go func() {
		wsMessage := dto.WebSocketMessage{
			Type:      "message",
			Data:      message,
			Timestamp: time.Now(),
		}
		a.hub.BroadcastToRoomIncludeSelf(req.RoomID, mustMarshal(wsMessage))
	}()
}

// GetMessages 获取消息列表
// @Summary 获取消息列表
// @Tags Chat
// @Param roomId query int true "聊天室ID"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param beforeId query int false "获取该消息之前的消息"
// @Param afterId query int false "获取该消息之后的消息"
// @Success 200 {object} map[string]interface{}
// @Router /api/chat/message/list [get]
func (a *ChatAPI) GetMessages(c *gin.Context) {
	var req dto.MessageListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	// 取消分页，全量返回消息
	messages, err := a.chatService.GetAllMessages(req.RoomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": messages,
	})
}

// RecallMessage 撤回消息
// @Summary 撤回消息
// @Tags Chat
// @Accept json
// @Produce json
// @Param request body dto.RecallMessageRequest true "撤回消息请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/chat/message/recall [post]
func (a *ChatAPI) RecallMessage(c *gin.Context) {
	var req dto.RecallMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	userID := middleware.GetUserID(c)
	if err := a.chatService.RecallMessage(req.MessageID, uint(userID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "撤回成功"})
}

// HandleWebSocket 处理 WebSocket 连接
// @Summary WebSocket连接
// @Tags Chat
// @Param roomId query int true "聊天室ID"
// @Router /api/chat/ws [get]
func (a *ChatAPI) HandleWebSocket(c *gin.Context) {
	roomIDStr := c.Query("roomId")
	if roomIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少 roomId 参数"})
		return
	}

	roomID, err := strconv.ParseUint(roomIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "roomId 参数错误"})
		return
	}

	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未登录"})
		return
	}

	username := middleware.GetUsername(c)

	// 升级 HTTP 连接为 WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	// 创建客户端
	client := ws.NewClient(a.hub, conn, uint(userID), username, uint(roomID))

	// 注册客户端
	a.hub.Register <- client

	// 启动读写协程
	go client.WritePump()
	go client.ReadPump()
}

// GetOnlineUsers 获取在线用户列表
// @Summary 获取在线用户列表
// @Tags Chat
// @Param id path int true "聊天室ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/chat/room/{id}/online-users [get]
func (a *ChatAPI) GetOnlineUsers(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	users := a.hub.GetRoomOnlineUsers(uint(id))
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": users})
}

// mustMarshal JSON 序列化（忽略错误）
func mustMarshal(v interface{}) []byte {
	data, _ := json.Marshal(v)
	return data
}
