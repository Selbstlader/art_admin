package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// 写入等待时间
	writeWait = 10 * time.Second
	// 读取超时时间
	pongWait = 60 * time.Second
	// ping 周期
	pingPeriod = (pongWait * 9) / 10
	// 最大消息大小
	maxMessageSize = 512 * 1024 // 512KB
)

// Client WebSocket 客户端
type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	userID   uint
	username string
	roomID   uint
	mu       sync.Mutex
}

// NewClient 创建新的客户端
func NewClient(hub *Hub, conn *websocket.Conn, userID uint, username string, roomID uint) *Client {
	return &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		userID:   userID,
		username: username,
		roomID:   roomID,
	}
}

// ReadPump 从 WebSocket 连接读取消息
func (c *Client) ReadPump() {
	defer func() {
		c.hub.Unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// 处理接收到的消息
		c.handleMessage(message)
	}
}

// WritePump 向 WebSocket 连接写入消息
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub 关闭了通道
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// 将队列中的其他消息一起发送
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage 处理接收到的消息
func (c *Client) handleMessage(message []byte) {
	var msg map[string]interface{}
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return
	}

	msgType, ok := msg["type"].(string)
	if !ok {
		return
	}

	switch msgType {
	case "typing":
		// 转发正在输入通知
		c.hub.BroadcastToRoom(c.roomID, message, c.userID)
	case "ping":
		// 响应心跳
		c.SendMessage(map[string]interface{}{
			"type":      "pong",
			"timestamp": time.Now(),
		})
	}
}

// SendMessage 发送消息给客户端
func (c *Client) SendMessage(data interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	message, err := json.Marshal(data)
	if err != nil {
		return err
	}

	select {
	case c.send <- message:
	default:
		// 发送缓冲区已满，关闭连接
		close(c.send)
		c.hub.Unregister <- c
	}

	return nil
}

// GetUserID 获取用户ID
func (c *Client) GetUserID() uint {
	return c.userID
}

// GetUsername 获取用户名
func (c *Client) GetUsername() string {
	return c.username
}

// GetRoomID 获取房间ID
func (c *Client) GetRoomID() uint {
	return c.roomID
}
