# 聊天室 API 文档

## 概述

公共聊天室功能支持多设备、多账号在线实时沟通，使用 WebSocket 实现实时消息推送。

## 功能特性

- ✅ 多设备同时在线
- ✅ 实时消息推送（WebSocket）
- ✅ 消息撤回（2分钟内）
- ✅ 回复消息
- ✅ 在线状态显示
- ✅ 正在输入提示
- ✅ 聊天室成员管理
- ✅ 消息历史记录
- ✅ 分页加载

## 数据模型

### ChatRoom（聊天室）
```go
type ChatRoom struct {
    ID          uint      // 聊天室ID
    Name        string    // 聊天室名称
    Description string    // 聊天室描述
    Type        string    // 类型：public/private
    MaxMembers  int       // 最大成员数，0表示无限制
    IsActive    bool      // 是否启用
    CreatedBy   uint      // 创建者ID
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

### ChatMessage（聊天消息）
```go
type ChatMessage struct {
    ID          uint      // 消息ID
    RoomID      uint      // 聊天室ID
    UserID      uint      // 发送者ID
    Username    string    // 发送者用户名
    Content     string    // 消息内容
    MessageType string    // 消息类型：text/image/file/system
    ReplyToID   *uint     // 回复的消息ID
    IsRecalled  bool      // 是否已撤回
    CreatedAt   time.Time
}
```

### ChatRoomMember（聊天室成员）
```go
type ChatRoomMember struct {
    ID         uint       // 成员ID
    RoomID     uint       // 聊天室ID
    UserID     uint       // 用户ID
    Username   string     // 用户名
    Role       string     // 角色：owner/admin/member
    IsMuted    bool       // 是否被禁言
    LastReadAt *time.Time // 最后阅读时间
    JoinedAt   time.Time  // 加入时间
}
```

## HTTP API 接口

### 1. 创建聊天室
```
POST /api/chat/room
Authorization: Bearer {token}

Request Body:
{
    "name": "公共聊天室",
    "description": "这是一个公共聊天室",
    "type": "public",
    "maxMembers": 100
}

Response:
{
    "code": 200,
    "message": "创建成功",
    "data": {
        "id": 1,
        "name": "公共聊天室",
        "description": "这是一个公共聊天室",
        "type": "public",
        "maxMembers": 100,
        "isActive": true,
        "createdBy": 1,
        "memberCount": 1,
        "onlineCount": 0,
        "createdAt": "2024-01-01T00:00:00Z",
        "updatedAt": "2024-01-01T00:00:00Z"
    }
}
```

### 2. 获取聊天室列表
```
GET /api/chat/room/list?page=1&pageSize=10&keyword=公共&type=public&isActive=true
Authorization: Bearer {token}

Response:
{
    "code": 200,
    "data": [...],
    "total": 10,
    "page": 1,
    "pageSize": 10
}
```

### 3. 获取聊天室详情
```
GET /api/chat/room/{id}
Authorization: Bearer {token}

Response:
{
    "code": 200,
    "data": {
        "id": 1,
        "name": "公共聊天室",
        ...
    }
}
```

### 4. 更新聊天室
```
PUT /api/chat/room
Authorization: Bearer {token}

Request Body:
{
    "id": 1,
    "name": "新的聊天室名称",
    "description": "新的描述",
    "maxMembers": 200,
    "isActive": true
}

Response:
{
    "code": 200,
    "message": "更新成功"
}
```

### 5. 删除聊天室
```
DELETE /api/chat/room/{id}
Authorization: Bearer {token}

Response:
{
    "code": 200,
    "message": "删除成功"
}
```

### 6. 加入聊天室
```
POST /api/chat/room/join
Authorization: Bearer {token}

Request Body:
{
    "roomId": 1
}

Response:
{
    "code": 200,
    "message": "加入成功"
}
```

### 7. 离开聊天室
```
POST /api/chat/room/{id}/leave
Authorization: Bearer {token}

Response:
{
    "code": 200,
    "message": "离开成功"
}
```

### 8. 获取聊天室成员列表
```
GET /api/chat/room/{id}/members
Authorization: Bearer {token}

Response:
{
    "code": 200,
    "data": [
        {
            "id": 1,
            "roomId": 1,
            "userId": 1,
            "username": "张三",
            "role": "owner",
            "isMuted": false,
            "isOnline": true,
            "joinedAt": "2024-01-01T00:00:00Z"
        }
    ]
}
```

### 9. 获取在线用户列表
```
GET /api/chat/room/{id}/online-users
Authorization: Bearer {token}

Response:
{
    "code": 200,
    "data": [
        {
            "userId": 1,
            "username": "张三"
        }
    ]
}
```

### 10. 获取我的聊天室列表
```
GET /api/chat/my-rooms
Authorization: Bearer {token}

Response:
{
    "code": 200,
    "data": [...]
}
```

### 11. 发送消息
```
POST /api/chat/message
Authorization: Bearer {token}

Request Body:
{
    "roomId": 1,
    "content": "你好，大家好！",
    "messageType": "text",
    "replyToId": null
}

Response:
{
    "code": 200,
    "message": "发送成功",
    "data": {
        "id": 1,
        "roomId": 1,
        "userId": 1,
        "username": "张三",
        "content": "你好，大家好！",
        "messageType": "text",
        "isRecalled": false,
        "createdAt": "2024-01-01T00:00:00Z"
    }
}
```

### 12. 获取消息列表
```
GET /api/chat/message/list?roomId=1&page=1&pageSize=20&beforeId=100
Authorization: Bearer {token}

Response:
{
    "code": 200,
    "data": [...],
    "total": 100,
    "page": 1,
    "pageSize": 20
}
```

### 13. 撤回消息
```
POST /api/chat/message/recall
Authorization: Bearer {token}

Request Body:
{
    "messageId": 1
}

Response:
{
    "code": 200,
    "message": "撤回成功"
}
```

## WebSocket 接口

### 连接
```
ws://localhost:48080/api/chat/ws?roomId=1
Authorization: Bearer {token}
```

### 消息格式

#### 1. 接收新消息
```json
{
    "type": "message",
    "data": {
        "id": 1,
        "roomId": 1,
        "userId": 1,
        "username": "张三",
        "content": "你好！",
        "messageType": "text",
        "createdAt": "2024-01-01T00:00:00Z"
    },
    "timestamp": "2024-01-01T00:00:00Z"
}
```

#### 2. 用户加入通知
```json
{
    "type": "join",
    "data": {
        "roomId": 1,
        "userId": 2,
        "username": "李四",
        "timestamp": "2024-01-01T00:00:00Z"
    }
}
```

#### 3. 用户离开通知
```json
{
    "type": "leave",
    "data": {
        "roomId": 1,
        "userId": 2,
        "username": "李四",
        "timestamp": "2024-01-01T00:00:00Z"
    }
}
```

#### 4. 正在输入通知（客户端发送）
```json
{
    "type": "typing",
    "data": {
        "roomId": 1,
        "userId": 1,
        "username": "张三",
        "isTyping": true
    }
}
```

#### 5. 心跳（客户端发送）
```json
{
    "type": "ping"
}
```

#### 6. 心跳响应（服务端返回）
```json
{
    "type": "pong",
    "timestamp": "2024-01-01T00:00:00Z"
}
```

## 使用示例

### JavaScript 客户端示例

```javascript
// 连接 WebSocket
const token = localStorage.getItem('token');
const roomId = 1;
const ws = new WebSocket(`ws://localhost:48080/api/chat/ws?roomId=${roomId}`);

// 连接成功
ws.onopen = () => {
    console.log('WebSocket 连接成功');
    
    // 发送心跳
    setInterval(() => {
        ws.send(JSON.stringify({ type: 'ping' }));
    }, 30000);
};

// 接收消息
ws.onmessage = (event) => {
    const message = JSON.parse(event.data);
    
    switch (message.type) {
        case 'message':
            // 显示新消息
            displayMessage(message.data);
            break;
        case 'join':
            // 显示用户加入通知
            showNotification(`${message.data.username} 加入了聊天室`);
            break;
        case 'leave':
            // 显示用户离开通知
            showNotification(`${message.data.username} 离开了聊天室`);
            break;
        case 'typing':
            // 显示正在输入提示
            showTypingIndicator(message.data);
            break;
        case 'pong':
            // 心跳响应
            console.log('收到心跳响应');
            break;
    }
};

// 发送正在输入通知
function sendTypingNotification(isTyping) {
    ws.send(JSON.stringify({
        type: 'typing',
        data: {
            roomId: roomId,
            isTyping: isTyping
        }
    }));
}

// 连接关闭
ws.onclose = () => {
    console.log('WebSocket 连接关闭');
    // 可以在这里实现重连逻辑
};

// 连接错误
ws.onerror = (error) => {
    console.error('WebSocket 错误:', error);
};
```

## 注意事项

1. **认证**：所有接口都需要 JWT 认证，WebSocket 连接也需要在 URL 参数中传递 token
2. **消息撤回**：只能撤回自己发送的消息，且发送时间不超过 2 分钟
3. **多设备支持**：同一用户可以在多个设备上同时在线，消息会同步到所有设备
4. **在线状态**：只有当用户的所有设备都断开连接时，才会显示为离线
5. **消息持久化**：所有消息都会保存到数据库，支持历史消息查询
6. **分页加载**：使用 `beforeId` 参数可以实现向上滚动加载历史消息

## 错误码

- `400` - 参数错误
- `401` - 未登录或 token 无效
- `403` - 无权限
- `404` - 资源不存在
- `500` - 服务器内部错误
