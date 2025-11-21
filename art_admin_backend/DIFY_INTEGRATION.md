# Dify 功能集成文档

## 概述

本项目已完成 Dify AI 平台的后端对接,包括知识库管理和 AI 对话两大功能模块。

## 配置说明

### 1. 配置文件设置

在 `config/config.yaml` 中添加 Dify 配置:

```yaml
dify:
  apiKey: your-dify-api-key-here  # 替换为你的 Dify API Key
  baseUrl: https://api.dify.ai/v1  # Dify API 基础URL
  timeout: 30                       # 请求超时时间(秒)
```

**重要**: 请将 `apiKey` 替换为你在 Dify 平台获取的真实 API Key。

### 2. 获取 API Key

1. 登录 [Dify Cloud](https://cloud.dify.ai)
2. 进入你的应用或知识库设置
3. 在 API 访问设置中获取 API Key

## API 接口说明

所有接口都需要 JWT 认证,请在请求头中添加:
```
Authorization: Bearer {your-jwt-token}
```

### 知识库管理

#### 1. 获取知识库列表

**接口**: `GET /api/dify/dataset/list`

**请求参数**:
- `page` (可选): 页码,默认 1
- `limit` (可选): 每页数量,默认 20,最大 100
- `keyword` (可选): 搜索关键词

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "data": [
      {
        "id": "3c90c3cc-0d44-4b50-8888-8dd25736052a",
        "name": "我的知识库",
        "description": "知识库描述",
        "document_count": 10,
        "word_count": 5000,
        "created_at": 1699999999
      }
    ],
    "has_more": false,
    "limit": 20,
    "total": 1,
    "page": 1
  }
}
```

#### 2. 获取知识库详情

**接口**: `GET /api/dify/dataset/{id}`

**路径参数**:
- `id`: 知识库ID

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": "3c90c3cc-0d44-4b50-8888-8dd25736052a",
    "name": "我的知识库",
    "description": "知识库描述",
    "document_count": 10,
    "word_count": 5000,
    "indexing_technique": "high_quality",
    "embedding_model": "text-embedding-ada-002"
  }
}
```

#### 3. 上传文件到知识库

**接口**: `POST /api/dify/dataset/upload`

**请求类型**: `multipart/form-data`

**请求参数**:
- `dataset_id`: 知识库ID (必填)
- `file`: 上传的文件 (必填)

**支持的文件格式**:
- 文档: txt, md, pdf, docx, xlsx, csv
- 其他格式请参考 Dify 文档

**响应示例**:
```json
{
  "code": 200,
  "message": "上传成功",
  "data": {
    "document": {
      "id": "doc-123456",
      "name": "example.pdf",
      "indexing_status": "parsing",
      "created_at": 1699999999
    },
    "batch": "batch-123456"
  }
}
```

### AI 对话

#### 1. 非流式对话

**接口**: `POST /api/dify/chat`

**请求体**:
```json
{
  "query": "你好,请介绍一下你自己",
  "user": "user-123",
  "conversation_id": "",
  "inputs": {},
  "auto_generate_name": true
}
```

**参数说明**:
- `query`: 用户问题 (必填)
- `user`: 用户标识 (必填)
- `conversation_id`: 会话ID,首次对话留空
- `inputs`: 输入变量,根据应用配置
- `auto_generate_name`: 是否自动生成会话标题

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "event": "message",
    "message_id": "msg-123456",
    "conversation_id": "conv-123456",
    "answer": "你好!我是 AI 助手...",
    "created_at": 1699999999
  }
}
```

#### 2. 流式对话 (推荐)

**接口**: `POST /api/dify/chat/stream`

**请求体**: 同非流式对话

**响应类型**: `text/event-stream` (SSE)

**响应格式**:
```
data: {"event":"message","answer":"你","message_id":"msg-123"}

data: {"event":"message","answer":"好","message_id":"msg-123"}

data: {"event":"message_end","message_id":"msg-123","conversation_id":"conv-123"}
```

**事件类型**:
- `message`: 消息内容片段
- `message_end`: 消息结束
- `error`: 错误事件

## 使用示例

### cURL 示例

#### 获取知识库列表
```bash
curl -X GET "http://localhost:48080/api/dify/dataset/list?page=1&limit=20" \
  -H "Authorization: Bearer your-jwt-token"
```

#### 上传文件
```bash
curl -X POST "http://localhost:48080/api/dify/dataset/upload" \
  -H "Authorization: Bearer your-jwt-token" \
  -F "dataset_id=your-dataset-id" \
  -F "file=@/path/to/your/file.pdf"
```

#### 流式对话
```bash
curl -X POST "http://localhost:48080/api/dify/chat/stream" \
  -H "Authorization: Bearer your-jwt-token" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "你好",
    "user": "user-123",
    "conversation_id": ""
  }'
```

### JavaScript 示例

#### 流式对话
```javascript
async function chatWithAI(query, conversationId = '') {
  const response = await fetch('http://localhost:48080/api/dify/chat/stream', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${jwtToken}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      query: query,
      user: 'user-123',
      conversation_id: conversationId,
      response_mode: 'streaming'
    })
  });

  const reader = response.body.getReader();
  const decoder = new TextDecoder();

  while (true) {
    const { done, value } = await reader.read();
    if (done) break;

    const chunk = decoder.decode(value);
    const lines = chunk.split('\n');

    for (const line of lines) {
      if (line.startsWith('data: ')) {
        const data = JSON.parse(line.slice(6));
        
        if (data.event === 'message') {
          // 处理消息片段
          console.log(data.answer);
        } else if (data.event === 'message_end') {
          // 对话结束
          console.log('Conversation ID:', data.conversation_id);
        }
      }
    }
  }
}
```

## 技术架构

### 目录结构

```
art_admin_backend/
├── internal/
│   ├── api/v1/
│   │   └── dify.go              # Dify API 控制器
│   ├── dto/
│   │   ├── request/
│   │   │   └── dify_request.go  # 请求 DTO
│   │   └── response/
│   │       └── dify_response.go # 响应 DTO
│   ├── pkg/dify/
│   │   └── client.go            # Dify 客户端封装
│   ├── service/
│   │   └── dify_service.go      # Dify 业务逻辑
│   └── router/
│       └── router.go            # 路由配置
└── config/
    └── config.yaml              # 配置文件
```

### 核心组件

1. **Dify Client** (`internal/pkg/dify/client.go`)
   - 封装 Dify API 的 HTTP 调用
   - 处理请求认证和错误响应
   - 支持流式和非流式响应

2. **Dify Service** (`internal/service/dify_service.go`)
   - 业务逻辑层
   - 参数验证和数据转换
   - 调用 Dify Client

3. **API Controller** (`internal/api/v1/dify.go`)
   - 处理 HTTP 请求
   - 参数绑定和验证
   - 返回统一格式的响应

## 注意事项

1. **API Key 安全**
   - 不要将 API Key 提交到版本控制系统
   - 生产环境使用环境变量或密钥管理服务
   - 定期轮换 API Key

2. **文件上传限制**
   - 默认文件大小限制由 Gin 框架控制
   - 建议在生产环境设置合理的文件大小限制
   - 上传的文件会临时保存在 `./temp/uploads` 目录

3. **流式响应**
   - 流式响应适合长文本生成场景
   - 客户端需要支持 SSE (Server-Sent Events)
   - 注意处理连接超时和错误重试

4. **并发控制**
   - Dify API 可能有速率限制
   - 建议实现请求队列和重试机制
   - 监控 API 使用量

## 错误处理

常见错误码:
- `400`: 参数错误
- `401`: 未授权 (JWT token 无效或过期)
- `500`: 服务器错误 (Dify API 调用失败)

错误响应格式:
```json
{
  "code": 500,
  "message": "Dify API错误 [invalid_api_key]: API Key is invalid"
}
```

## 后续优化建议

1. **缓存机制**: 对知识库列表等数据添加缓存
2. **异步处理**: 文件上传和索引可以改为异步任务
3. **监控告警**: 添加 API 调用监控和告警
4. **重试机制**: 实现请求失败自动重试
5. **日志记录**: 详细记录 API 调用日志便于排查问题

## 相关文档

- [Dify 官方文档](https://docs.dify.ai/)
- [Dify API 参考](https://docs.dify.ai/api-reference)
- [Dify GitHub](https://github.com/langgenius/dify)
