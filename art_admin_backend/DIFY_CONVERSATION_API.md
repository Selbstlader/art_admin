# Dify 会话管理 API 文档

## 新增接口列表

### 1. 停止消息生成

停止正在进行的AI对话消息生成(仅流式模式支持)

**接口**: `POST /api/dify/chat/stop/{task_id}`

**请求参数**:
- Path: `task_id` (string, 必填) - 任务ID,从流式响应中获取
- Query: `user` (string, 必填) - 用户标识

**请求示例**:
```bash
curl -X POST 'http://localhost:48080/api/dify/chat/stop/8afa6024-7550-4c23-96ff-6fa402317388?user=user-123' \
  -H 'Authorization: Bearer {your-jwt-token}'
```

**响应示例**:
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "result": "success"
  }
}
```

---

### 2. 获取建议问题

获取当前消息的下一轮建议问题列表

**接口**: `GET /api/dify/messages/{message_id}/suggested`

**请求参数**:
- Path: `message_id` (string, 必填) - 消息ID
- Query: `user` (string, 必填) - 用户标识

**请求示例**:
```bash
curl -X GET 'http://localhost:48080/api/dify/messages/8d678ffe-6c9d-48e3-aca0-fa66bb8c75b7/suggested?user=user-123' \
  -H 'Authorization: Bearer {your-jwt-token}'
```

**响应示例**:
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "result": "success",
    "data": [
      "如何提高市场占有率?",
      "环保产品的发展趋势是什么?",
      "新技术应用的具体案例有哪些?"
    ]
  }
}
```

---

### 3. 获取会话列表

获取当前用户的会话列表

**接口**: `GET /api/dify/conversations`

**请求参数**:
- Query: `user` (string, 必填) - 用户标识
- Query: `last_id` (string, 可选) - 最后一条记录ID,用于分页
- Query: `limit` (int, 可选) - 每页数量,默认20,最大100

**请求示例**:
```bash
curl -X GET 'http://localhost:48080/api/dify/conversations?user=user-123&limit=20' \
  -H 'Authorization: Bearer {your-jwt-token}'
```

**响应示例**:
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "limit": 20,
    "has_more": false,
    "data": [
      {
        "id": "eb5ea4b5-6db8-42ce-8b79-80e25d78a9cb",
        "name": "市场分析讨论",
        "inputs": {},
        "status": "normal",
        "introduction": "",
        "created_at": 1763631447,
        "updated_at": 1763631447
      }
    ]
  }
}
```

---

### 4. 获取会话消息历史

获取指定会话的历史消息记录

**接口**: `GET /api/dify/messages`

**请求参数**:
- Query: `conversation_id` (string, 必填) - 会话ID
- Query: `user` (string, 必填) - 用户标识
- Query: `first_id` (string, 可选) - 第一条记录ID,用于分页
- Query: `limit` (int, 可选) - 每页数量,默认20,最大100

**请求示例**:
```bash
curl -X GET 'http://localhost:48080/api/dify/messages?conversation_id=eb5ea4b5-6db8-42ce-8b79-80e25d78a9cb&user=user-123&limit=20' \
  -H 'Authorization: Bearer {your-jwt-token}'
```

**响应示例**:
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "limit": 20,
    "has_more": false,
    "data": [
      {
        "id": "8d678ffe-6c9d-48e3-aca0-fa66bb8c75b7",
        "conversation_id": "eb5ea4b5-6db8-42ce-8b79-80e25d78a9cb",
        "inputs": {
          "report_title": "2024年市场分析报告",
          "report_content": "...",
          "report_date": "2024-11-20",
          "analysis_focus": "市场趋势"
        },
        "query": "请分析这份报告",
        "answer": "关键洞察：\n- 市场增长率达到15%...",
        "message_files": [],
        "feedback": null,
        "retriever_resources": [],
        "created_at": 1763631447
      }
    ]
  }
}
```

---

### 5. 删除会话

删除指定的会话

**接口**: `DELETE /api/dify/conversations/{conversation_id}`

**请求参数**:
- Path: `conversation_id` (string, 必填) - 会话ID
- Query: `user` (string, 必填) - 用户标识

**请求示例**:
```bash
curl -X DELETE 'http://localhost:48080/api/dify/conversations/eb5ea4b5-6db8-42ce-8b79-80e25d78a9cb?user=user-123' \
  -H 'Authorization: Bearer {your-jwt-token}'
```

**响应示例**:
```json
{
  "code": 200,
  "msg": "删除成功",
  "data": null
}
```

---

### 6. 重命名会话

修改会话的名称

**接口**: `POST /api/dify/conversations/{conversation_id}/name`

**请求参数**:
- Path: `conversation_id` (string, 必填) - 会话ID
- Body:
  ```json
  {
    "user": "user-123",
    "name": "新的会话名称"
  }
  ```

**请求示例**:
```bash
curl -X POST 'http://localhost:48080/api/dify/conversations/eb5ea4b5-6db8-42ce-8b79-80e25d78a9cb/name' \
  -H 'Authorization: Bearer {your-jwt-token}' \
  -H 'Content-Type: application/json' \
  -d '{
    "user": "user-123",
    "name": "2024市场分析讨论"
  }'
```

**响应示例**:
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "id": "eb5ea4b5-6db8-42ce-8b79-80e25d78a9cb",
    "name": "2024市场分析讨论",
    "inputs": {},
    "status": "normal",
    "introduction": "",
    "created_at": 1763631447,
    "updated_at": 1763631500
  }
}
```

---

## 完整使用流程示例

### 场景: 用户进行AI对话并管理会话

#### 1. 开始新对话
```bash
# 发起流式对话
curl -X POST 'http://localhost:48080/api/dify/chat/stream' \
  -H 'Authorization: Bearer {token}' \
  -H 'Content-Type: application/json' \
  -d '{
    "query": "请分析这份报告",
    "user": "user-123",
    "inputs": {
      "report_title": "2024年市场分析报告",
      "report_content": "...",
      "report_date": "2024-11-20",
      "analysis_focus": "市场趋势"
    }
  }'

# 响应会返回 conversation_id 和 message_id
```

#### 2. 获取建议问题
```bash
# 使用返回的 message_id 获取建议问题
curl -X GET 'http://localhost:48080/api/dify/messages/{message_id}/suggested?user=user-123' \
  -H 'Authorization: Bearer {token}'
```

#### 3. 查看会话列表
```bash
# 获取所有会话
curl -X GET 'http://localhost:48080/api/dify/conversations?user=user-123&limit=20' \
  -H 'Authorization: Bearer {token}'
```

#### 4. 查看会话历史
```bash
# 获取某个会话的所有消息
curl -X GET 'http://localhost:48080/api/dify/messages?conversation_id={conversation_id}&user=user-123' \
  -H 'Authorization: Bearer {token}'
```

#### 5. 重命名会话
```bash
# 给会话起个有意义的名字
curl -X POST 'http://localhost:48080/api/dify/conversations/{conversation_id}/name' \
  -H 'Authorization: Bearer {token}' \
  -H 'Content-Type: application/json' \
  -d '{
    "user": "user-123",
    "name": "2024市场分析"
  }'
```

#### 6. 继续对话
```bash
# 使用 conversation_id 继续对话
curl -X POST 'http://localhost:48080/api/dify/chat' \
  -H 'Authorization: Bearer {token}' \
  -H 'Content-Type: application/json' \
  -d '{
    "query": "如何提高市场占有率?",
    "user": "user-123",
    "conversation_id": "{conversation_id}",
    "inputs": {
      "report_title": "2024年市场分析报告",
      "report_content": "...",
      "report_date": "2024-11-20",
      "analysis_focus": "市场策略"
    }
  }'
```

#### 7. 停止生成(流式模式)
```bash
# 如果生成时间太长,可以停止
curl -X POST 'http://localhost:48080/api/dify/chat/stop/{task_id}?user=user-123' \
  -H 'Authorization: Bearer {token}'
```

#### 8. 删除会话
```bash
# 不需要的会话可以删除
curl -X DELETE 'http://localhost:48080/api/dify/conversations/{conversation_id}?user=user-123' \
  -H 'Authorization: Bearer {token}'
```

---

## 前端集成示例 (Vue 3 + TypeScript)

### 会话管理 Composable

```typescript
// composables/useDifyConversation.ts
import { ref } from 'vue'
import { request } from '@/utils/request'

export interface Conversation {
  id: string
  name: string
  inputs: Record<string, any>
  status: string
  introduction: string
  created_at: number
  updated_at: number
}

export interface Message {
  id: string
  conversation_id: string
  inputs: Record<string, any>
  query: string
  answer: string
  created_at: number
}

export function useDifyConversation() {
  const conversations = ref<Conversation[]>([])
  const messages = ref<Message[]>([])
  const suggestedQuestions = ref<string[]>([])
  const loading = ref(false)

  // 获取会话列表
  const getConversations = async (user: string, limit = 20) => {
    loading.value = true
    try {
      const res = await request.get('/api/dify/conversations', {
        params: { user, limit }
      })
      conversations.value = res.data.data
      return res.data
    } finally {
      loading.value = false
    }
  }

  // 获取会话消息历史
  const getMessages = async (conversationId: string, user: string, limit = 20) => {
    loading.value = true
    try {
      const res = await request.get('/api/dify/messages', {
        params: { conversation_id: conversationId, user, limit }
      })
      messages.value = res.data.data
      return res.data
    } finally {
      loading.value = false
    }
  }

  // 获取建议问题
  const getSuggestedQuestions = async (messageId: string, user: string) => {
    try {
      const res = await request.get(`/api/dify/messages/${messageId}/suggested`, {
        params: { user }
      })
      suggestedQuestions.value = res.data.data
      return res.data.data
    } catch (error) {
      console.error('获取建议问题失败:', error)
      return []
    }
  }

  // 重命名会话
  const renameConversation = async (conversationId: string, user: string, name: string) => {
    loading.value = true
    try {
      const res = await request.post(`/api/dify/conversations/${conversationId}/name`, {
        user,
        name
      })
      // 更新本地会话列表
      const index = conversations.value.findIndex(c => c.id === conversationId)
      if (index !== -1) {
        conversations.value[index] = res.data
      }
      return res.data
    } finally {
      loading.value = false
    }
  }

  // 删除会话
  const deleteConversation = async (conversationId: string, user: string) => {
    loading.value = true
    try {
      await request.delete(`/api/dify/conversations/${conversationId}`, {
        params: { user }
      })
      // 从本地列表中移除
      conversations.value = conversations.value.filter(c => c.id !== conversationId)
    } finally {
      loading.value = false
    }
  }

  // 停止消息生成
  const stopMessage = async (taskId: string, user: string) => {
    try {
      await request.post(`/api/dify/chat/stop/${taskId}`, null, {
        params: { user }
      })
    } catch (error) {
      console.error('停止消息失败:', error)
    }
  }

  return {
    conversations,
    messages,
    suggestedQuestions,
    loading,
    getConversations,
    getMessages,
    getSuggestedQuestions,
    renameConversation,
    deleteConversation,
    stopMessage
  }
}
```

### 使用示例

```vue
<template>
  <div class="conversation-manager">
    <!-- 会话列表 -->
    <div class="conversation-list">
      <div
        v-for="conv in conversations"
        :key="conv.id"
        class="conversation-item"
        @click="selectConversation(conv.id)"
      >
        <div class="conv-name">{{ conv.name }}</div>
        <div class="conv-actions">
          <button @click.stop="handleRename(conv.id)">重命名</button>
          <button @click.stop="handleDelete(conv.id)">删除</button>
        </div>
      </div>
    </div>

    <!-- 消息历史 -->
    <div class="message-history">
      <div v-for="msg in messages" :key="msg.id" class="message">
        <div class="query">{{ msg.query }}</div>
        <div class="answer">{{ msg.answer }}</div>
      </div>
    </div>

    <!-- 建议问题 -->
    <div v-if="suggestedQuestions.length" class="suggested-questions">
      <h4>您可能想问:</h4>
      <button
        v-for="(q, i) in suggestedQuestions"
        :key="i"
        @click="askQuestion(q)"
      >
        {{ q }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useDifyConversation } from '@/composables/useDifyConversation'

const {
  conversations,
  messages,
  suggestedQuestions,
  getConversations,
  getMessages,
  getSuggestedQuestions,
  renameConversation,
  deleteConversation
} = useDifyConversation()

const user = 'user-123' // 从用户状态获取

onMounted(() => {
  getConversations(user)
})

const selectConversation = async (conversationId: string) => {
  await getMessages(conversationId, user)
  // 获取最后一条消息的建议问题
  if (messages.value.length > 0) {
    const lastMessage = messages.value[messages.value.length - 1]
    await getSuggestedQuestions(lastMessage.id, user)
  }
}

const handleRename = async (conversationId: string) => {
  const newName = prompt('请输入新名称:')
  if (newName) {
    await renameConversation(conversationId, user, newName)
  }
}

const handleDelete = async (conversationId: string) => {
  if (confirm('确定要删除这个会话吗?')) {
    await deleteConversation(conversationId, user)
  }
}

const askQuestion = (question: string) => {
  // 使用建议问题发起新对话
  console.log('Ask:', question)
}
</script>
```

---

## 注意事项

1. **用户标识**: 所有接口都需要提供 `user` 参数,用于隔离不同用户的会话
2. **会话隔离**: Service API 创建的会话与 WebApp 创建的会话是隔离的
3. **分页**: 会话列表和消息历史都支持分页,使用 `last_id` 或 `first_id` 进行翻页
4. **停止功能**: 停止消息生成仅在流式模式下有效
5. **建议问题**: 需要在 Dify 应用中启用"下一步问题建议"功能
6. **JWT认证**: 所有接口都需要有效的 JWT token

---

## API 接口总览

| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| 知识库列表 | GET | `/api/dify/dataset/list` | 获取知识库列表 |
| 知识库详情 | GET | `/api/dify/dataset/:id` | 获取知识库详情 |
| 上传文件 | POST | `/api/dify/dataset/upload` | 上传文件到知识库 |
| AI对话 | POST | `/api/dify/chat` | 非流式对话 |
| AI对话(流式) | POST | `/api/dify/chat/stream` | 流式对话 |
| **停止生成** | POST | `/api/dify/chat/stop/:task_id` | 停止消息生成 |
| **建议问题** | GET | `/api/dify/messages/:message_id/suggested` | 获取建议问题 |
| **会话列表** | GET | `/api/dify/conversations` | 获取会话列表 |
| **消息历史** | GET | `/api/dify/messages` | 获取消息历史 |
| **删除会话** | DELETE | `/api/dify/conversations/:conversation_id` | 删除会话 |
| **重命名会话** | POST | `/api/dify/conversations/:conversation_id/name` | 重命名会话 |
