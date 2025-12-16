# Design Document: AI Tag Management

## Overview

AI标签管理系统提供一个统一的配置管理层，将 Dify 知识库、系统提示词和 API 密钥封装为可复用的"标签"。开发者在构建 AI 对话功能时，只需选择对应标签即可自动应用完整配置，无需重复编写配置代码。

系统采用前后端分离架构：
- 后端：Go + Gin 提供 RESTful API
- 前端：Vue 3 + Element Plus 提供管理界面
- 存储：MySQL 持久化标签数据

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        Frontend (Vue 3)                         │
├─────────────────────────────────────────────────────────────────┤
│  Tag Management Page    │    Tag Selector Component             │
│  - List/Search          │    - Dropdown selection               │
│  - Create/Edit/Delete   │    - Integration with chat            │
│  - Test configuration   │                                       │
└─────────────────────────┴───────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────┐
│                        Backend (Go + Gin)                       │
├─────────────────────────────────────────────────────────────────┤
│  API Layer              │    Service Layer                      │
│  - /api/v1/ai-tags/*    │    - AITagService                     │
│  - Validation           │    - Business logic                   │
│  - Auth middleware      │    - Dify integration                 │
├─────────────────────────┴───────────────────────────────────────┤
│  Repository Layer       │    Model Layer                        │
│  - AITagRepository      │    - AITag entity                     │
│  - CRUD operations      │    - GORM mapping                     │
└─────────────────────────┴───────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────┐
│                        Database (MySQL)                         │
│                        ai_tags table                            │
└─────────────────────────────────────────────────────────────────┘
```

## Components and Interfaces

### Backend Components

#### 1. AITag Model (`internal/model/ai_tag.go`)

```go
type AITag struct {
    ID              uint           `gorm:"primaryKey" json:"id"`
    Name            string         `gorm:"size:100;uniqueIndex;not null" json:"name"`
    Description     string         `gorm:"size:500" json:"description"`
    KnowledgeBaseID string         `gorm:"size:100" json:"knowledge_base_id"`      // 从 Dify 知识库列表选择
    KnowledgeBaseName string       `gorm:"size:200" json:"knowledge_base_name"`    // 知识库名称（冗余存储）
    SystemPrompt    string         `gorm:"type:text;not null" json:"system_prompt"`
    ChatAPIKey      string         `gorm:"size:200" json:"chat_api_key"`           // Dify Chat App API Key（手动输入）
    Status          int            `gorm:"default:1" json:"status"` // 1=active, 0=inactive
    CreatedAt       time.Time      `json:"created_at"`
    UpdatedAt       time.Time      `json:"updated_at"`
    DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}
```

> **说明**：
> - `KnowledgeBaseID`: 通过调用 Dify 知识库列表接口获取，前端提供下拉选择
> - `ChatAPIKey`: Dify 未提供查询接口，需要用户手动输入对应 Dify App 的 API Key

#### 2. AITag Repository (`internal/repository/ai_tag_repo.go`)

```go
type AITagRepository interface {
    Create(tag *model.AITag) error
    Update(tag *model.AITag) error
    Delete(id uint) error
    FindByID(id uint) (*model.AITag, error)
    FindByName(name string) (*model.AITag, error)
    List(page, pageSize int, keyword string) ([]*model.AITag, int64, error)
    FindActiveByID(id uint) (*model.AITag, error)
}
```

#### 2.1 Knowledge Base Integration

前端创建/编辑标签时，知识库选择流程：
1. 调用现有的 `/api/dify/dataset/list` 接口获取 Dify 知识库列表
2. 用户从下拉列表中选择知识库
3. 保存时同时存储 `knowledge_base_id` 和 `knowledge_base_name`

```typescript
// 复用现有 API
import { difyDatasetApi } from '@/api/dify'

// 获取知识库列表供选择
const knowledgeBaseList = await difyDatasetApi.getDatasetList({ page: 1, limit: 100 })
```

#### 3. AITag Service (`internal/service/ai_tag/ai_tag_service.go`)

```go
type AITagService interface {
    CreateTag(req *request.CreateAITagRequest) (*model.AITag, error)
    UpdateTag(id uint, req *request.UpdateAITagRequest) (*model.AITag, error)
    DeleteTag(id uint) error
    GetTag(id uint) (*model.AITag, error)
    ListTags(page, pageSize int, keyword string) (*response.AITagListResponse, error)
    TestTag(id uint) (*response.AITagTestResponse, error)
    GetTagConfig(id uint) (*AITagConfig, error) // For chat integration
}
```

#### 4. AITag API Handler (`internal/api/v1/ai_tag_api.go`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/v1/ai-tags | Create new tag |
| GET | /api/v1/ai-tags | List tags with pagination |
| GET | /api/v1/ai-tags/:id | Get tag details |
| PUT | /api/v1/ai-tags/:id | Update tag |
| DELETE | /api/v1/ai-tags/:id | Soft delete tag |
| POST | /api/v1/ai-tags/:id/test | Test tag configuration |

### Frontend Components

#### 1. Tag Management Page (`src/views/system/ai-tag/index.vue`)

- 标签列表表格（分页、搜索）
- 新增/编辑对话框
  - 知识库：下拉选择（调用 Dify API 获取列表）
  - Chat API Key：文本输入框（手动输入，支持密码显示/隐藏）
  - 系统提示词：多行文本框
- 删除确认
- 测试按钮（验证配置是否有效）

#### 2. Tag Selector Component (`src/components/core/forms/AITagSelector.vue`)

- 下拉选择器
- 显示标签名称和描述
- 支持 v-model 双向绑定

#### 3. API Module (`src/api/ai-tag.ts`)

```typescript
export const aiTagApi = {
  create: (data: CreateAITagRequest) => request.post({ url: '/api/v1/ai-tags', params: data }),
  list: (params: AITagListParams) => request.get({ url: '/api/v1/ai-tags', params }),
  get: (id: number) => request.get({ url: `/api/v1/ai-tags/${id}` }),
  update: (id: number, data: UpdateAITagRequest) => request.put({ url: `/api/v1/ai-tags/${id}`, params: data }),
  delete: (id: number) => request.del({ url: `/api/v1/ai-tags/${id}` }),
  test: (id: number) => request.post({ url: `/api/v1/ai-tags/${id}/test` })
}
```

## Data Models

### Database Schema

```sql
CREATE TABLE ai_tags (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE COMMENT '标签名称',
    description VARCHAR(500) DEFAULT '' COMMENT '标签描述',
    knowledge_base_id VARCHAR(100) DEFAULT '' COMMENT 'Dify知识库ID（从API获取）',
    knowledge_base_name VARCHAR(200) DEFAULT '' COMMENT '知识库名称（冗余存储）',
    system_prompt TEXT NOT NULL COMMENT '系统提示词',
    chat_api_key VARCHAR(200) DEFAULT '' COMMENT 'Dify Chat App API Key（手动输入）',
    status TINYINT DEFAULT 1 COMMENT '1=active, 0=inactive',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### Request/Response DTOs

```go
// CreateAITagRequest
type CreateAITagRequest struct {
    Name              string `json:"name" binding:"required,max=100"`
    Description       string `json:"description" binding:"max=500"`
    KnowledgeBaseID   string `json:"knowledge_base_id"`      // 从 Dify 知识库列表选择
    KnowledgeBaseName string `json:"knowledge_base_name"`    // 知识库名称
    SystemPrompt      string `json:"system_prompt" binding:"required"`
    ChatAPIKey        string `json:"chat_api_key"`           // Dify Chat App API Key（手动输入）
}

// UpdateAITagRequest
type UpdateAITagRequest struct {
    Name              string `json:"name" binding:"max=100"`
    Description       string `json:"description" binding:"max=500"`
    KnowledgeBaseID   string `json:"knowledge_base_id"`
    KnowledgeBaseName string `json:"knowledge_base_name"`
    SystemPrompt      string `json:"system_prompt"`
    ChatAPIKey        string `json:"chat_api_key"`
    Status            *int   `json:"status"`
}

// AITagListResponse
type AITagListResponse struct {
    List  []*model.AITag `json:"list"`
    Total int64          `json:"total"`
    Page  int            `json:"page"`
    Size  int            `json:"size"`
}

// AITagTestResponse
type AITagTestResponse struct {
    Success           bool   `json:"success"`
    Response          string `json:"response"`
    KnowledgeAccessed bool   `json:"knowledge_accessed"`
    ErrorMessage      string `json:"error_message,omitempty"`
}
```

## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system-essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

### Property 1: Tag creation returns complete tag data
*For any* valid tag creation request with name, knowledge_base_id, and system_prompt, the created tag returned by the service SHALL contain all submitted fields plus a valid non-zero ID and timestamps.
**Validates: Requirements 1.1**

### Property 2: Duplicate name rejection on create
*For any* existing tag name in the database, attempting to create a new tag with the same name SHALL result in a rejection error.
**Validates: Requirements 1.2**

### Property 3: Empty/whitespace validation rejection
*For any* tag creation request where name or system_prompt consists entirely of whitespace characters, the creation SHALL be rejected with a validation error.
**Validates: Requirements 1.3**

### Property 4: Tag persistence round-trip
*For any* successfully created tag, immediately querying the database by the returned ID SHALL return a tag with identical field values.
**Validates: Requirements 1.4**

### Property 5: Pagination correctness
*For any* list request with page P and pageSize S on a database with N tags, the returned list SHALL contain at most S items, and the total count SHALL equal N.
**Validates: Requirements 2.1**

### Property 6: Keyword search completeness
*For any* keyword K and set of tags, the search result SHALL include all and only tags where name OR description contains K as a substring (case-insensitive).
**Validates: Requirements 2.2**

### Property 7: Response structure completeness
*For any* tag in a list response, the tag object SHALL contain non-null values for: id, name, status, created_at, and system_prompt.
**Validates: Requirements 2.3**

### Property 8: Update reflects changes
*For any* valid update request on an existing tag, the returned tag SHALL reflect all changed fields while preserving unchanged fields.
**Validates: Requirements 3.1**

### Property 9: Duplicate name rejection on update
*For any* two distinct tags A and B, attempting to update tag A's name to match tag B's name SHALL result in a rejection error.
**Validates: Requirements 3.2**

### Property 10: Update timestamp advancement
*For any* successful tag update, the updated_at timestamp SHALL be greater than or equal to the previous updated_at value.
**Validates: Requirements 3.3**

### Property 11: Soft delete sets inactive status
*For any* deleted tag, querying the database (including soft-deleted records) SHALL show the tag with status=0 (inactive) and a non-null deleted_at timestamp.
**Validates: Requirements 4.1**

### Property 12: Chat uses tag's knowledge base ID
*For any* chat session initiated with a tag, the underlying Dify API call SHALL include the tag's knowledge_base_id in the request context.
**Validates: Requirements 5.1**

### Property 13: Chat prepends system prompt
*For any* chat session initiated with a tag, the conversation context SHALL include the tag's system_prompt as the system message.
**Validates: Requirements 5.2**

### Property 14: Chat API key selection
*For any* tag with a non-empty chat_api_key, chat requests using that tag SHALL use the tag's chat_api_key; otherwise, the system default chat API key from config SHALL be used.
**Validates: Requirements 5.3**

### Property 15: Test response structure
*For any* tag test execution, the response SHALL include: success (boolean), response (string), and knowledge_accessed (boolean) fields.
**Validates: Requirements 6.2**

### Property 16: Serialization round-trip consistency
*For any* AITag object, serializing to JSON and deserializing back SHALL produce an object with identical field values.
**Validates: Requirements 7.1, 7.2**

## Error Handling

### Backend Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| TAG_NOT_FOUND | 404 | Tag with specified ID does not exist |
| TAG_NAME_EXISTS | 400 | Tag name already exists |
| TAG_VALIDATION_ERROR | 400 | Request validation failed |
| TAG_IN_USE | 409 | Tag is in use by active sessions |
| DIFY_API_ERROR | 502 | Dify API call failed |
| INTERNAL_ERROR | 500 | Unexpected server error |

### Error Response Format

```json
{
  "code": 400,
  "message": "Tag name already exists",
  "data": null
}
```

## Testing Strategy

### Unit Testing

Unit tests will cover:
- Repository CRUD operations with mock database
- Service business logic validation
- API handler request/response handling
- Frontend component rendering and interactions

### Property-Based Testing

Property-based tests will use **fast-check** (TypeScript) for frontend and **gopter** (Go) for backend to verify the correctness properties defined above.

Each property-based test MUST:
1. Run a minimum of 100 iterations
2. Include a comment referencing the property: `// **Feature: ai-tag-management, Property N: <property_text>**`
3. Use smart generators that constrain inputs to valid ranges

### Test Categories

1. **Model Tests**: Validate GORM model serialization/deserialization
2. **Repository Tests**: Verify database operations with test database
3. **Service Tests**: Test business logic with mocked dependencies
4. **API Tests**: Integration tests for HTTP endpoints
5. **Frontend Tests**: Component tests with Vitest + Vue Test Utils
