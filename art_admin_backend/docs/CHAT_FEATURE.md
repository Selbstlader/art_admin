# 公共聊天室功能说明

## 功能概述

公共聊天室是一个支持多设备、多账号在线实时沟通的功能模块。使用 WebSocket 技术实现实时消息推送，提供完整的聊天室管理和消息管理功能。

## 技术架构

### 后端技术栈
- **Go 1.21+**
- **Gin** - Web 框架
- **GORM** - ORM 框架
- **gorilla/websocket** - WebSocket 支持
- **MySQL** - 数据库

### 架构设计
```
┌─────────────────────────────────────────────────────────┐
│                     客户端层                              │
│  (浏览器/移动端) - WebSocket + HTTP API                   │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│                     API 层                               │
│  - ChatAPI (HTTP 接口)                                   │
│  - WebSocket Handler (实时通信)                          │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│                   Service 层                             │
│  - ChatService (业务逻辑)                                │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│                 Repository 层                            │
│  - ChatRoomRepository                                    │
│  - ChatMessageRepository                                 │
│  - ChatRoomMemberRepository                              │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│                   数据库层                                │
│  - chat_rooms (聊天室表)                                 │
│  - chat_messages (消息表)                                │
│  - chat_room_members (成员表)                            │
└─────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│              WebSocket Hub (消息中心)                     │
│  - 连接管理 (支持多设备)                                  │
│  - 消息广播 (房间级别)                                    │
│  - 在线状态管理                                           │
└─────────────────────────────────────────────────────────┘
```

## 核心功能

### 1. 聊天室管理
- ✅ 创建聊天室（公开/私密）
- ✅ 编辑聊天室信息
- ✅ 删除聊天室
- ✅ 设置成员上限
- ✅ 启用/禁用聊天室

### 2. 成员管理
- ✅ 加入聊天室
- ✅ 离开聊天室
- ✅ 查看成员列表
- ✅ 查看在线用户
- ✅ 成员角色管理（owner/admin/member）
- ✅ 禁言功能

### 3. 消息功能
- ✅ 发送文本消息
- ✅ 发送图片/文件（预留接口）
- ✅ 回复消息
- ✅ 撤回消息（2分钟内）
- ✅ 消息历史记录
- ✅ 分页加载

### 4. 实时功能
- ✅ 实时消息推送
- ✅ 在线状态显示
- ✅ 正在输入提示
- ✅ 用户加入/离开通知
- ✅ 多设备同步

## 数据库表结构

### chat_rooms（聊天室表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键 |
| name | VARCHAR(100) | 聊天室名称 |
| description | VARCHAR(500) | 聊天室描述 |
| type | VARCHAR(20) | 类型（public/private） |
| max_members | INT | 最大成员数 |
| is_active | BOOLEAN | 是否启用 |
| created_by | BIGINT | 创建者ID |
| created_at | DATETIME | 创建时间 |
| updated_at | DATETIME | 更新时间 |
| deleted_at | DATETIME | 删除时间 |

### chat_messages（消息表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键 |
| room_id | BIGINT | 聊天室ID |
| user_id | BIGINT | 发送者ID |
| username | VARCHAR(50) | 发送者用户名 |
| content | TEXT | 消息内容 |
| message_type | VARCHAR(20) | 消息类型 |
| reply_to_id | BIGINT | 回复的消息ID |
| is_recalled | BOOLEAN | 是否已撤回 |
| created_at | DATETIME | 创建时间 |
| updated_at | DATETIME | 更新时间 |
| deleted_at | DATETIME | 删除时间 |

### chat_room_members（成员表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键 |
| room_id | BIGINT | 聊天室ID |
| user_id | BIGINT | 用户ID |
| username | VARCHAR(50) | 用户名 |
| role | VARCHAR(20) | 角色 |
| is_muted | BOOLEAN | 是否被禁言 |
| last_read_at | DATETIME | 最后阅读时间 |
| joined_at | DATETIME | 加入时间 |
| created_at | DATETIME | 创建时间 |
| updated_at | DATETIME | 更新时间 |
| deleted_at | DATETIME | 删除时间 |

## 文件结构

```
art_admin_backend/
├── internal/
│   ├── model/
│   │   ├── chat_room.go              # 聊天室模型
│   │   ├── chat_message.go           # 消息模型
│   │   └── chat_room_member.go       # 成员模型
│   ├── dto/
│   │   └── chat_dto.go               # DTO 定义
│   ├── repository/
│   │   └── chat_room_repository.go   # 数据访问层
│   ├── service/
│   │   └── chat_service.go           # 业务逻辑层
│   ├── api/
│   │   └── chat_api.go               # API 控制器
│   ├── pkg/
│   │   └── websocket/
│   │       ├── hub.go                # WebSocket 中心
│   │       └── client.go             # WebSocket 客户端
│   └── router/
│       └── chat_router.go            # 路由配置
└── docs/
    ├── CHAT_API.md                   # API 文档
    └── CHAT_FEATURE.md               # 功能说明（本文件）
```

## 启动说明

### 1. 安装依赖
```bash
cd art_admin_backend
go mod tidy
```

### 2. 配置数据库
确保 `config/config.yaml` 中的数据库配置正确。

### 3. 启动服务
```bash
go run cmd/server/main.go
```

服务启动后会自动：
- 创建数据库表
- 启动 WebSocket Hub
- 注册所有路由

### 4. 访问 Swagger 文档
```
http://localhost:48080/swagger/index.html
```

## API 使用流程

### 1. 创建聊天室
```bash
curl -X POST http://localhost:48080/api/chat/room \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "公共聊天室",
    "description": "欢迎大家",
    "type": "public",
    "maxMembers": 100
  }'
```

### 2. 加入聊天室
```bash
curl -X POST http://localhost:48080/api/chat/room/join \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "roomId": 1
  }'
```

### 3. 连接 WebSocket
```javascript
const ws = new WebSocket('ws://localhost:48080/api/chat/ws?roomId=1');
```

### 4. 发送消息
```bash
curl -X POST http://localhost:48080/api/chat/message \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "roomId": 1,
    "content": "大家好！",
    "messageType": "text"
  }'
```

## WebSocket 连接管理

### Hub 设计
- **多设备支持**：同一用户可以在多个设备上同时连接
- **房间隔离**：消息只广播到对应的聊天室
- **自动清理**：连接断开时自动清理资源
- **心跳机制**：定期发送 ping/pong 保持连接

### 连接流程
1. 客户端发起 WebSocket 连接（携带 token 和 roomId）
2. 服务端验证 token 和房间权限
3. 创建 Client 对象并注册到 Hub
4. 启动读写协程
5. 广播用户加入通知

### 断开流程
1. 客户端断开连接或超时
2. 触发 unregister 流程
3. 从 Hub 中移除 Client
4. 如果用户所有设备都断开，广播离开通知

## 性能优化

### 1. 消息广播优化
- 使用 channel 异步广播
- 房间级别隔离，避免全局广播
- 缓冲区设置合理大小

### 2. 数据库优化
- 消息表添加索引（room_id, created_at）
- 成员表添加索引（room_id, user_id）
- 使用软删除避免数据丢失

### 3. 内存优化
- 使用 sync.RWMutex 减少锁竞争
- 及时清理断开的连接
- 限制消息缓冲区大小

## 安全考虑

### 1. 认证授权
- 所有接口都需要 JWT 认证
- WebSocket 连接也需要验证 token
- 检查用户是否是聊天室成员

### 2. 权限控制
- 只有创建者可以删除聊天室
- 只能撤回自己的消息
- 禁言用户无法发送消息

### 3. 数据验证
- 消息长度限制（5000字符）
- 聊天室名称长度限制（100字符）
- 防止 SQL 注入（使用 GORM 参数化查询）

### 4. 限流保护
- WebSocket 消息大小限制（512KB）
- 可以添加发送频率限制
- 可以添加连接数限制

## 扩展建议

### 1. 功能扩展
- [ ] 图片/文件上传
- [ ] 表情包支持
- [ ] @提及功能
- [ ] 消息搜索
- [ ] 消息已读状态
- [ ] 群公告
- [ ] 聊天室分组

### 2. 性能扩展
- [ ] Redis 缓存在线用户
- [ ] 消息队列（Kafka/RabbitMQ）
- [ ] 分布式 WebSocket（多实例）
- [ ] CDN 加速文件传输

### 3. 监控告警
- [ ] 在线人数监控
- [ ] 消息发送量统计
- [ ] WebSocket 连接数监控
- [ ] 异常日志告警

## 常见问题

### Q1: WebSocket 连接失败？
A: 检查以下几点：
- token 是否有效
- roomId 是否正确
- 是否已加入聊天室
- 网络连接是否正常

### Q2: 消息没有实时推送？
A: 检查：
- WebSocket 连接是否正常
- Hub 是否正常运行
- 消息是否成功保存到数据库

### Q3: 多设备不同步？
A: 确认：
- 所有设备都连接到同一个聊天室
- Hub 的用户索引是否正确维护
- 消息广播逻辑是否正确

### Q4: 如何实现消息加密？
A: 可以在以下层面加密：
- 传输层：使用 WSS（WebSocket over TLS）
- 应用层：在发送前加密消息内容
- 数据库：加密存储敏感消息

## 联系方式

如有问题或建议，请联系开发团队。
