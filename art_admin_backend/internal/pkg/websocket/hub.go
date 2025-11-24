package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"
)

// Hub 维护活动客户端集合并向客户端广播消息
type Hub struct {
	// 注册的客户端
	clients map[*Client]bool

	// 按房间分组的客户端
	rooms map[uint]map[*Client]bool

	// 按用户ID索引的客户端（支持多设备）
	userClients map[uint]map[*Client]bool

	// 广播消息到所有客户端
	broadcast chan []byte

	// 注册客户端请求
	Register chan *Client

	// 注销客户端请求
	Unregister chan *Client

	// 互斥锁
	mu sync.RWMutex
}

// NewHub 创建新的 Hub
func NewHub() *Hub {
	return &Hub{
		clients:     make(map[*Client]bool),
		rooms:       make(map[uint]map[*Client]bool),
		userClients: make(map[uint]map[*Client]bool),
		broadcast:   make(chan []byte, 256),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
	}
}

// Run 启动 Hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.registerClient(client)

		case client := <-h.Unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastMessage(message)
		}
	}
}

// registerClient 注册客户端
func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[client] = true

	// 添加到房间
	if h.rooms[client.roomID] == nil {
		h.rooms[client.roomID] = make(map[*Client]bool)
	}
	h.rooms[client.roomID][client] = true

	// 添加到用户客户端列表（支持多设备）
	if h.userClients[client.userID] == nil {
		h.userClients[client.userID] = make(map[*Client]bool)
	}
	h.userClients[client.userID][client] = true

	log.Printf("Client registered: UserID=%d, Username=%s, RoomID=%d, Total clients=%d",
		client.userID, client.username, client.roomID, len(h.clients))

	// 通知房间内其他用户有新成员加入
	h.notifyUserJoined(client)
}

// unregisterClient 注销客户端
func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.send)

		// 从房间移除
		if room, ok := h.rooms[client.roomID]; ok {
			delete(room, client)
			if len(room) == 0 {
				delete(h.rooms, client.roomID)
			}
		}

		// 从用户客户端列表移除
		if userClients, ok := h.userClients[client.userID]; ok {
			delete(userClients, client)
			if len(userClients) == 0 {
				delete(h.userClients, client.userID)
			}
		}

		log.Printf("Client unregistered: UserID=%d, Username=%s, RoomID=%d, Total clients=%d",
			client.userID, client.username, client.roomID, len(h.clients))

		// 通知房间内其他用户有成员离开
		h.notifyUserLeft(client)
	}
}

// broadcastMessage 广播消息到所有客户端
func (h *Hub) broadcastMessage(message []byte) {
	h.mu.RLock()
	// 收集所有客户端
	var clients []*Client
	for client := range h.clients {
		clients = append(clients, client)
	}
	h.mu.RUnlock()

	// 在不持有锁的情况下发送消息
	for _, client := range clients {
		select {
		case client.send <- message:
		default:
			// 发送失败，标记客户端需要注销
			go func(c *Client) {
				h.Unregister <- c
			}(client)
		}
	}
}

// BroadcastToRoom 向指定房间广播消息
func (h *Hub) BroadcastToRoom(roomID uint, message []byte, excludeUserID uint) {
	h.mu.RLock()
	room, ok := h.rooms[roomID]
	if !ok {
		h.mu.RUnlock()
		return
	}

	// 收集需要发送的客户端
	var clients []*Client
	for client := range room {
		// 排除指定用户
		if excludeUserID > 0 && client.userID == excludeUserID {
			continue
		}
		clients = append(clients, client)
	}
	h.mu.RUnlock()

	// 在不持有锁的情况下发送消息
	for _, client := range clients {
		select {
		case client.send <- message:
		default:
			// 发送失败，标记客户端需要注销
			go func(c *Client) {
				h.Unregister <- c
			}(client)
		}
	}
}

// BroadcastToRoomIncludeSelf 向指定房间广播消息（包括自己）
func (h *Hub) BroadcastToRoomIncludeSelf(roomID uint, message []byte) {
	h.BroadcastToRoom(roomID, message, 0)
}

// SendToUser 向指定用户的所有设备发送消息
func (h *Hub) SendToUser(userID uint, message []byte) {
	h.mu.RLock()
	clientsMap, ok := h.userClients[userID]
	if !ok {
		h.mu.RUnlock()
		return
	}

	// 收集需要发送的客户端
	var clients []*Client
	for client := range clientsMap {
		clients = append(clients, client)
	}
	h.mu.RUnlock()

	// 在不持有锁的情况下发送消息
	for _, client := range clients {
		select {
		case client.send <- message:
		default:
			// 发送失败，标记客户端需要注销
			go func(c *Client) {
				h.Unregister <- c
			}(client)
		}
	}
}

// GetRoomOnlineCount 获取房间在线人数
func (h *Hub) GetRoomOnlineCount(roomID uint) int {
	// 使用带超时的读锁，避免长时间阻塞
	done := make(chan int)
	go func() {
		h.mu.RLock()
		defer h.mu.RUnlock()

		room, ok := h.rooms[roomID]
		if !ok {
			done <- 0
			return
		}

		// 统计不重复的用户数
		users := make(map[uint]bool)
		for client := range room {
			users[client.userID] = true
		}

		done <- len(users)
	}()

	// 等待最多100ms，避免阻塞HTTP请求
	select {
	case count := <-done:
		return count
	case <-time.After(100 * time.Millisecond):
		// 超时返回默认值，避免阻塞
		log.Printf("GetRoomOnlineCount timeout for room %d", roomID)
		return 0
	}
}

// GetRoomOnlineUsers 获取房间在线用户列表
func (h *Hub) GetRoomOnlineUsers(roomID uint) []map[string]interface{} {
	// 使用带超时的读锁，避免长时间阻塞
	done := make(chan []map[string]interface{})
	go func() {
		h.mu.RLock()
		defer h.mu.RUnlock()

		room, ok := h.rooms[roomID]
		if !ok {
			done <- []map[string]interface{}{}
			return
		}

		// 使用 map 去重
		users := make(map[uint]map[string]interface{})
		for client := range room {
			if _, exists := users[client.userID]; !exists {
				users[client.userID] = map[string]interface{}{
					"userId":   client.userID,
					"username": client.username,
				}
			}
		}

		// 转换为数组
		result := make([]map[string]interface{}, 0, len(users))
		for _, user := range users {
			result = append(result, user)
		}

		done <- result
	}()

	// 等待最多100ms，避免阻塞HTTP请求
	select {
	case users := <-done:
		return users
	case <-time.After(100 * time.Millisecond):
		// 超时返回默认值，避免阻塞
		log.Printf("GetRoomOnlineUsers timeout for room %d", roomID)
		return []map[string]interface{}{}
	}
}

// IsUserOnline 检查用户是否在线
func (h *Hub) IsUserOnline(userID uint) bool {
	// 使用带超时的读锁，避免长时间阻塞
	done := make(chan bool)
	go func() {
		h.mu.RLock()
		defer h.mu.RUnlock()

		clients, ok := h.userClients[userID]
		done <- (ok && len(clients) > 0)
	}()

	// 等待最多50ms，避免阻塞HTTP请求
	select {
	case isOnline := <-done:
		return isOnline
	case <-time.After(50 * time.Millisecond):
		// 超时返回默认值，避免阻塞
		log.Printf("IsUserOnline timeout for user %d", userID)
		return false
	}
}

// notifyUserJoined 通知用户加入
func (h *Hub) notifyUserJoined(client *Client) {
	// 检查该用户是否已经有其他设备在线
	isFirstDevice := len(h.userClients[client.userID]) == 1

	// 只有第一个设备连接时才通知其他用户
	if !isFirstDevice {
		return
	}

	notification := map[string]interface{}{
		"type": "join",
		"data": map[string]interface{}{
			"roomId":    client.roomID,
			"userId":    client.userID,
			"username":  client.username,
			"timestamp": time.Now(),
		},
	}

	message, _ := json.Marshal(notification)
	h.BroadcastToRoom(client.roomID, message, client.userID)
}

// notifyUserLeft 通知用户离开
func (h *Hub) notifyUserLeft(client *Client) {
	// 检查该用户是否还有其他设备在线
	hasOtherDevices := len(h.userClients[client.userID]) > 0

	// 如果还有其他设备在线，不通知
	if hasOtherDevices {
		return
	}

	notification := map[string]interface{}{
		"type": "leave",
		"data": map[string]interface{}{
			"roomId":    client.roomID,
			"userId":    client.userID,
			"username":  client.username,
			"timestamp": time.Now(),
		},
	}

	message, _ := json.Marshal(notification)
	h.BroadcastToRoom(client.roomID, message, 0)
}
