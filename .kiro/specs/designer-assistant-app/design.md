# Design Document

## Overview

本设计文档描述工装设计师AI辅助系统移动端应用（Designer_App）的技术架构和实现方案。该应用基于 UniApp 框架开发，支持编译到 iOS、Android 和微信小程序三端，复用现有后端 Go/Gin API，实现与 Web 端数据完全同步。

**UI 设计风格要求：**
- C4D 风格高级 UI，3D 立体数据图表
- 蓝色渐变界面，悬浮卡片设计
- 玻璃拟态按钮，动态数据实时显示
- 纯白色背景，高细节卡片设计

**技术选型：**
- 框架：UniApp (Vue 3 + TypeScript + Vite)
- 状态管理：Pinia
- UI 组件：uni-ui + 自定义组件
- 网络请求：uni.request 封装
- 数据缓存：uni.setStorageSync / uni.getStorageSync

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        Designer App                              │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │   iOS App   │  │ Android App │  │  微信小程序  │              │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘              │
│         │                │                │                      │
│         └────────────────┼────────────────┘                      │
│                          │                                       │
│  ┌───────────────────────▼───────────────────────┐              │
│  │              UniApp Runtime                    │              │
│  │  ┌─────────────────────────────────────────┐  │              │
│  │  │           Pages (页面层)                 │  │              │
│  │  │  ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐       │  │              │
│  │  │  │Login│ │Home │ │Chat │ │Mine │       │  │              │
│  │  │  └─────┘ └─────┘ └─────┘ └─────┘       │  │              │
│  │  └─────────────────────────────────────────┘  │              │
│  │  ┌─────────────────────────────────────────┐  │              │
│  │  │         Components (组件层)              │  │              │
│  │  │  ┌──────────┐ ┌──────────┐ ┌─────────┐ │  │              │
│  │  │  │ProjectCard│ │ChatBubble│ │CostCard │ │  │              │
│  │  │  └──────────┘ └──────────┘ └─────────┘ │  │              │
│  │  └─────────────────────────────────────────┘  │              │
│  │  ┌─────────────────────────────────────────┐  │              │
│  │  │           Store (状态管理)               │  │              │
│  │  │  ┌────┐ ┌───────┐ ┌────┐ ┌──────────┐  │  │              │
│  │  │  │User│ │Project│ │Chat│ │Notification│ │  │              │
│  │  │  └────┘ └───────┘ └────┘ └──────────┘  │  │              │
│  │  └─────────────────────────────────────────┘  │              │
│  │  ┌─────────────────────────────────────────┐  │              │
│  │  │            API (接口层)                  │  │              │
│  │  │  ┌────┐ ┌───────┐ ┌────┐ ┌────────┐    │  │              │
│  │  │  │Auth│ │Project│ │Chat│ │Material│    │  │              │
│  │  │  └────┘ └───────┘ └────┘ └────────┘    │  │              │
│  │  └─────────────────────────────────────────┘  │              │
│  │  ┌─────────────────────────────────────────┐  │              │
│  │  │           Utils (工具层)                 │  │              │
│  │  │  ┌───────┐ ┌─────┐ ┌────────┐          │  │              │
│  │  │  │Request│ │Cache│ │Platform│          │  │              │
│  │  │  └───────┘ └─────┘ └────────┘          │  │              │
│  │  └─────────────────────────────────────────┘  │              │
│  └───────────────────────────────────────────────┘              │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ HTTPS
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Backend API (Go/Gin)                          │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐           │
│  │  Auth    │ │ Project  │ │   Chat   │ │ Material │           │
│  │  API     │ │   API    │ │   API    │ │   API    │           │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘           │
└─────────────────────────────────────────────────────────────────┘
```

## Components and Interfaces

### 1. 页面组件结构

```
pages/
├── login/                    # 登录模块
│   └── index.vue            # 登录页面
├── index/                    # 首页模块
│   └── index.vue            # 项目列表首页
├── project/                  # 项目模块
│   ├── detail.vue           # 项目详情
│   ├── documents.vue        # 文档分析列表
│   ├── cost.vue             # 成本报告
│   └── designs.vue          # 设计图列表
├── chat/                     # AI对话模块 (仅App)
│   └── index.vue            # 对话页面
├── material/                 # 材料库模块
│   ├── index.vue            # 材料分类列表
│   ├── list.vue             # 材料列表
│   └── detail.vue           # 材料详情
├── message/                  # 消息模块 (仅App)
│   └── index.vue            # 消息中心
└── mine/                     # 我的模块
    └── index.vue            # 个人中心
```

### 2. 核心组件接口


#### ProjectCard 组件
```typescript
/*** 项目卡片组件 - 悬浮卡片设计 ***/
interface ProjectCardProps {
  project: {
    id: number
    name: string
    status: 'draft' | 'in_progress' | 'completed' | 'archived'
    area: number
    budget: number
    style: string
    updatedAt: string
  }
}

interface ProjectCardEmits {
  (e: 'click', id: number): void
}
```

#### ChatBubble 组件
```typescript
/*** 对话气泡组件 - 玻璃拟态设计 ***/
interface ChatBubbleProps {
  message: {
    id: number
    role: 'user' | 'assistant'
    content: string
    createdAt: string
  }
  isStreaming?: boolean
}
```

#### CostSummaryCard 组件
```typescript
/*** 成本汇总卡片 - 3D立体数据图表 ***/
interface CostSummaryCardProps {
  totalCost: number
  budgetLimit: number
  categories: {
    name: string
    amount: number
    percentage: number
  }[]
}
```

#### MaterialCard 组件
```typescript
/*** 材料卡片组件 ***/
interface MaterialCardProps {
  material: {
    id: number
    name: string
    category: string
    brand: string
    price: number
    unit: string
    imageUrl: string
  }
}
```

### 3. API 接口层

```typescript
/*** API 基础配置 ***/
interface ApiConfig {
  baseUrl: string
  timeout: number
  headers: Record<string, string>
}

/*** 统一响应格式 ***/
interface ApiResponse<T = unknown> {
  code: number
  msg?: string
  data?: T
}

/*** 分页响应格式 ***/
interface PageResponse<T> {
  records: T[]
  total: number
  current: number
  size: number
}
```

#### Auth API
```typescript
/*** 认证接口 ***/
interface AuthApi {
  login(username: string, password: string): Promise<ApiResponse<LoginResponse>>
  logout(): Promise<ApiResponse<null>>
  refreshToken(): Promise<ApiResponse<TokenResponse>>
  getUserInfo(): Promise<ApiResponse<UserInfo>>
}

interface LoginResponse {
  accessToken: string
  refreshToken: string
  expiresIn: number
  user: UserInfo
}

interface UserInfo {
  id: number
  username: string
  nickname: string
  avatar: string
  role: string
}
```

#### Project API
```typescript
/*** 项目接口 ***/
interface ProjectApi {
  getList(params: ProjectListParams): Promise<ApiResponse<PageResponse<Project>>>
  getDetail(id: number): Promise<ApiResponse<ProjectDetail>>
  search(keyword: string, page: number, size: number): Promise<ApiResponse<PageResponse<Project>>>
}

interface ProjectListParams {
  current: number
  size: number
  name?: string
  status?: string
  style?: string
}

interface Project {
  id: number
  name: string
  description: string
  area: number
  budget: number
  style: string
  status: string
  createdAt: string
  updatedAt: string
}

interface ProjectDetail extends Project {
  documentCount: number
  designCount: number
  lastAnalysisTime: string
  costEstimate?: CostEstimate
}
```

#### Chat API (仅 App)
```typescript
/*** AI对话接口 ***/
interface ChatApi {
  sendMessage(data: ChatRequest): Promise<ApiResponse<ChatResponse>>
  streamMessage(data: ChatRequest, callbacks: StreamCallbacks): Promise<void>
  getHistory(params: HistoryParams): Promise<ApiResponse<PageResponse<ChatMessage>>>
  getSessions(projectId?: number): Promise<ApiResponse<ChatSession[]>>
}

interface ChatRequest {
  query: string
  projectId?: number
  sessionId?: string
}

interface ChatResponse {
  answer: string
  sessionId: string
  messageId: number
  tokensUsed: number
}

interface StreamCallbacks {
  onChunk: (content: string) => void
  onDone: (result: { sessionId: string; messageId: number }) => void
  onError: (error: string) => void
}
```

#### Document API
```typescript
/*** 文档接口 ***/
interface DocumentApi {
  getList(projectId: number, params: PageParams): Promise<ApiResponse<PageResponse<Document>>>
  getDetail(id: number): Promise<ApiResponse<DocumentDetail>>
  getKeywords(id: number): Promise<ApiResponse<string[]>>
  getSummary(id: number): Promise<ApiResponse<DocumentSummary>>
}

interface Document {
  id: number
  projectId: number
  fileName: string
  fileType: string
  fileSize: number
  analysisStatus: 'pending' | 'analyzing' | 'completed' | 'failed'
  createdAt: string
}

interface DocumentSummary {
  overview: string      // 项目概述
  requirements: string  // 核心需求
  special: string       // 特殊要求
}
```

#### Cost API
```typescript
/*** 成本接口 ***/
interface CostApi {
  getByProjectId(projectId: number): Promise<ApiResponse<CostEstimate>>
  getSummary(): Promise<ApiResponse<CostSummary>>
}

interface CostEstimate {
  id: number
  projectId: number
  totalCost: number
  budgetLimit: number
  items: CostItem[]
  createdAt: string
}

interface CostItem {
  category: string  // 材料费、人工费、设备费、管理费
  name: string
  quantity: number
  unitPrice: number
  totalPrice: number
}
```

#### Material API
```typescript
/*** 材料接口 ***/
interface MaterialApi {
  getCategories(): Promise<ApiResponse<MaterialCategory[]>>
  getList(params: MaterialListParams): Promise<ApiResponse<PageResponse<Material>>>
  getDetail(id: number): Promise<ApiResponse<MaterialDetail>>
  search(keyword: string): Promise<ApiResponse<Material[]>>
}

interface MaterialCategory {
  name: string
  count: number
  icon: string
}

interface Material {
  id: number
  name: string
  category: string
  brand: string
  price: number
  unit: string
  imageUrl: string
}

interface MaterialDetail extends Material {
  specification: string
  supplier: string
  applicableScenes: string[]
  images: string[]
}
```

## Data Models

### 1. 本地存储数据模型

```typescript
/*** 用户存储 ***/
interface UserStorage {
  token: string
  refreshToken: string
  tokenExpireTime: number
  userInfo: UserInfo
}

/*** 缓存数据 ***/
interface CacheData<T> {
  data: T
  timestamp: number
  expireTime: number  // 7天过期
}

/*** 项目列表缓存 ***/
interface ProjectListCache extends CacheData<Project[]> {
  total: number
  lastPage: number
}

/*** 离线模式标识 ***/
interface OfflineState {
  isOffline: boolean
  lastOnlineTime: number
}
```

### 2. Pinia Store 数据模型

```typescript
/*** User Store ***/
interface UserState {
  token: string
  userInfo: UserInfo | null
  isLoggedIn: boolean
}

/*** Project Store ***/
interface ProjectState {
  list: Project[]
  currentProject: ProjectDetail | null
  total: number
  loading: boolean
  searchKeyword: string
}

/*** Chat Store (仅App) ***/
interface ChatState {
  messages: ChatMessage[]
  currentSessionId: string
  isStreaming: boolean
  sessions: ChatSession[]
}

/*** Notification Store (仅App) ***/
interface NotificationState {
  list: Notification[]
  unreadCount: number
  hasPermission: boolean
}
```

### 3. 平台差异处理

```typescript
/*** 平台检测工具 ***/
interface PlatformUtils {
  isApp(): boolean           // iOS 或 Android
  isMiniProgram(): boolean   // 微信小程序
  isIOS(): boolean
  isAndroid(): boolean
  canUseAIChat(): boolean    // App 可用，小程序不可用
  canUsePush(): boolean      // App 可用，小程序使用订阅消息
}

/*** 条件编译示例 ***/
// #ifdef APP-PLUS
// App 专属代码：AI对话、推送通知
// #endif

// #ifdef MP-WEIXIN
// 小程序专属代码：订阅消息
// #endif
```


## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system-essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*


### Property 1: Token 管理一致性
*For any* 用户登录/退出操作，登录成功后本地存储应包含有效 token，退出后本地存储应不包含 token。
**Validates: Requirements 1.2, 1.4**

### Property 2: 项目列表渲染正确性
*For any* 项目数据数组，渲染后的卡片列表长度应等于数据数组长度，且每个卡片应包含项目名称、状态、面积、更新时间。
**Validates: Requirements 2.1, 3.2**

### Property 3: 分页加载数据增长
*For any* 分页加载操作（项目列表、对话历史），如果存在更多数据，加载后列表长度应大于加载前。
**Validates: Requirements 2.3, 4.7**

### Property 4: 搜索过滤结果匹配
*For any* 搜索关键字和数据列表，过滤后的结果中每一项的目标字段（项目名称/材料名称/品牌）应包含该关键字（不区分大小写）。
**Validates: Requirements 2.4, 10.2, 10.3**

### Property 5: 消息发送即时显示
*For any* 用户发送的消息内容，发送后消息列表应立即包含该消息，且消息角色为 'user'。
**Validates: Requirements 4.4**

### Property 6: 流式输出内容累积
*For any* AI 流式回复，每次 chunk 到达后显示的内容应为之前所有 chunk 的累积拼接。
**Validates: Requirements 4.5**

### Property 7: 项目上下文传递
*For any* 从项目详情页发起的对话请求，请求参数中应包含当前项目的 projectId。
**Validates: Requirements 4.6**

### Property 8: 文档摘要结构完整性
*For any* 已分析的文档摘要数据，应包含 overview（项目概述）、requirements（核心需求）、special（特殊要求）三个字段。
**Validates: Requirements 5.3**

### Property 9: 成本超支警告显示
*For any* 成本数据，当 totalCost > budgetLimit 时，UI 应显示警告样式（红色）。
**Validates: Requirements 6.2**

### Property 10: 成本明细分类正确性
*For any* 成本明细列表，按类别分组后每组内的所有项目 category 字段应相同。
**Validates: Requirements 6.3**

### Property 11: 通知列表时间排序
*For any* 通知列表，列表中每一项的 createdAt 时间应大于等于下一项的 createdAt 时间（倒序）。
**Validates: Requirements 7.5**

### Property 12: 通知已读状态更新
*For any* 点击通知操作，操作后该通知的 isRead 状态应为 true。
**Validates: Requirements 7.6**

### Property 13: 未读消息红点显示
*For any* 通知列表，当存在 isRead 为 false 的通知时，导航栏应显示红点。
**Validates: Requirements 7.7**

### Property 14: 缓存数据一致性
*For any* 成功加载的项目列表数据，应能从本地存储中读取到相同的数据。
**Validates: Requirements 8.1**

### Property 15: 缓存过期判断
*For any* 缓存数据，当 (当前时间 - 缓存时间) > 7天 时，应标记为过期需要更新。
**Validates: Requirements 8.4**

### Property 16: 设计图切换边界处理
*For any* 设计图列表和当前索引，左滑时索引应 +1（不超过最大值），右滑时索引应 -1（不小于 0）。
**Validates: Requirements 9.3**

### Property 17: 材料分类过滤正确性
*For any* 选中的材料分类，过滤后列表中每个材料的 category 字段应等于选中的分类。
**Validates: Requirements 10.2**

### Property 18: 用户信息展示完整性
*For any* 用户数据，"我的"页面应显示 avatar、nickname、role 三个字段的值。
**Validates: Requirements 11.1**

### Property 19: 平台功能可用性
*For any* 平台类型，App 平台（iOS/Android）应显示 AI 对话和消息入口，小程序平台应隐藏这些入口。
**Validates: Requirements 4.1, 4.2, 7.1**

### Property 20: 错误状态显示
*For any* API 请求失败，UI 应显示错误提示信息，且提供重试操作入口。
**Validates: Requirements 3.5, 4.8, 9.5**

## Error Handling

### 1. 网络错误处理

```typescript
/*** 网络错误类型 ***/
enum NetworkErrorType {
  TIMEOUT = 'TIMEOUT',           // 请求超时
  OFFLINE = 'OFFLINE',           // 网络离线
  SERVER_ERROR = 'SERVER_ERROR', // 服务器错误
  UNAUTHORIZED = 'UNAUTHORIZED'  // 未授权
}

/*** 错误处理策略 ***/
interface ErrorHandler {
  // 超时：提示用户检查网络，提供重试
  handleTimeout(): void
  // 离线：切换到离线模式，使用缓存数据
  handleOffline(): void
  // 服务器错误：提示服务暂时不可用
  handleServerError(): void
  // 未授权：清除 token，跳转登录页
  handleUnauthorized(): void
}
```

### 2. 业务错误处理

```typescript
/*** 业务错误码 ***/
const BusinessErrorCode = {
  PROJECT_NOT_FOUND: 40001,      // 项目不存在
  DOCUMENT_NOT_ANALYZED: 40002, // 文档未分析
  AI_SERVICE_UNAVAILABLE: 50001 // AI 服务不可用
}

/*** 错误提示映射 ***/
const ErrorMessages: Record<number, string> = {
  40001: '项目不存在或已被删除',
  40002: '该文档尚未分析，请在 Web 端进行分析',
  50001: 'AI 服务暂时不可用，请稍后重试'
}
```

### 3. 离线模式处理

```typescript
/*** 离线模式管理 ***/
interface OfflineManager {
  // 检测网络状态
  checkNetwork(): boolean
  // 进入离线模式
  enterOfflineMode(): void
  // 退出离线模式
  exitOfflineMode(): void
  // 获取缓存数据
  getCachedData<T>(key: string): T | null
  // 显示离线提示
  showOfflineToast(): void
}
```

## Testing Strategy

### 1. 单元测试

使用 Vitest 进行单元测试，覆盖以下模块：

- **工具函数测试**：日期格式化、金额格式化、缓存管理
- **Store 测试**：状态管理逻辑、数据更新
- **API 封装测试**：请求拦截、响应处理、错误处理

```typescript
/*** 测试配置 ***/
// vitest.config.ts
export default {
  test: {
    environment: 'jsdom',
    globals: true,
    coverage: {
      reporter: ['text', 'json', 'html'],
      exclude: ['node_modules/', 'dist/']
    }
  }
}
```

### 2. 属性测试

使用 fast-check 进行属性测试，验证核心业务逻辑：

```typescript
/*** 属性测试示例 ***/
import fc from 'fast-check'

// Property 4: 搜索过滤结果匹配
describe('Search Filter Property', () => {
  it('filtered results should contain search keyword', () => {
    fc.assert(
      fc.property(
        fc.string({ minLength: 1 }),
        fc.array(fc.record({ name: fc.string(), brand: fc.string() })),
        (keyword, items) => {
          const filtered = filterByKeyword(items, keyword)
          return filtered.every(item => 
            item.name.toLowerCase().includes(keyword.toLowerCase()) ||
            item.brand.toLowerCase().includes(keyword.toLowerCase())
          )
        }
      ),
      { numRuns: 100 }
    )
  })
})
```

### 3. 组件测试

使用 @vue/test-utils 进行组件测试：

- **渲染测试**：验证组件正确渲染
- **交互测试**：验证用户交互响应
- **Props 测试**：验证属性传递正确

### 4. E2E 测试

使用 uni-app 官方测试工具进行端到端测试：

- **登录流程测试**
- **项目列表加载测试**
- **AI 对话流程测试**（仅 App）
- **离线模式测试**

### 5. 平台兼容性测试

- **iOS 真机测试**：iPhone 12 及以上
- **Android 真机测试**：Android 10 及以上
- **微信小程序测试**：微信开发者工具 + 真机预览

