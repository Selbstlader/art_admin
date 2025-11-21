# Dify AI 对话功能 - 前端实现说明

## 已完成的工作

### 1. API 接口封装 (`src/api/dify.ts`)

创建了三个 API 模块:
- `difyDatasetApi` - 知识库管理
- `difyChatApi` - AI 对话
- `difyConversationApi` - 会话管理

### 2. 类型定义 (`src/types/dify.ts`)

定义了完整的 TypeScript 类型:
- `DifyDataset` - 知识库类型
- `DifyConversation` - 会话类型
- `DifyMessage` - 消息类型
- `DifyChatRequest` - 对话请求类型
- `DifyChatResponse` - 对话响应类型
- `DifyChatStreamEvent` - 流式事件类型

### 3. Composable (`src/composables/useDifyChat.ts`)

创建了 `useDifyChat` 组合式函数,提供:

**状态管理:**
- `conversations` - 会话列表
- `currentConversation` - 当前会话
- `messages` - 消息列表
- `suggestedQuestions` - 建议问题
- `loading` - 加载状态
- `streaming` - 流式状态
- `currentAnswer` - 当前回答

**方法:**
- `getConversations()` - 获取会话列表
- `getMessages()` - 获取消息历史
- `sendMessage()` - 发送消息(非流式)
- `sendMessageStream()` - 发送消息(流式)
- `stopMessage()` - 停止消息生成
- `getSuggestedQuestions()` - 获取建议问题
- `deleteConversation()` - 删除会话
- `renameConversation()` - 重命名会话
- `selectConversation()` - 选择会话
- `newConversation()` - 新建会话

### 4. 路由配置 (`src/router/modules/dify.ts`)

添加了 Dify 模块路由:
- `/dify/chat` - AI 对话页面
- `/dify/knowledge` - 知识库管理页面

### 5. 页面组件 (`src/views/dify/chat/index.vue`)

创建了完整的 AI 对话页面,包含:

**功能特性:**
- ✅ 左侧会话列表
- ✅ 右侧对话区域
- ✅ 消息历史显示
- ✅ 流式实时输出
- ✅ 建议问题展示
- ✅ 报告信息表单
- ✅ 会话管理(重命名、删除)
- ✅ 停止生成功能
- ✅ Markdown 渲染
- ✅ 自动滚动到底部

**UI 设计:**
- 采用左右分栏布局
- 用户消息和 AI 回复区分显示
- 支持实时流式输出
- 建议问题标签点击发送
- 响应式设计

## 使用方法

### 1. 安装依赖

```bash
# 需要安装 marked 用于 Markdown 渲染
pnpm install marked
pnpm install dayjs
```

### 2. 配置环境变量

在 `.env.development` 中配置:
```
VITE_API_BASE_URL=http://localhost:48080
```

### 3. 注册路由

在 `src/router/index.ts` 中导入并注册路由:
```typescript
import difyRoutes from './modules/dify'

// 添加到路由配置中
routes.push(difyRoutes)
```

### 4. 启动项目

```bash
pnpm dev
```

### 5. 访问页面

访问 `http://localhost:5173/dify/chat` 即可使用 AI 对话功能

## 功能演示

### 发起对话

1. 填写报告信息:
   - 报告标题
   - 报告内容
   - 报告日期
   - 分析重点

2. 输入问题

3. 点击"发送"按钮

4. 实时查看 AI 回复(流式输出)

### 会话管理

- **新建对话**: 点击"新对话"按钮
- **选择会话**: 点击左侧会话列表
- **重命名**: 点击会话右侧菜单 → 重命名
- **删除**: 点击会话右侧菜单 → 删除

### 建议问题

AI 回复后会显示建议问题(如果 Dify 应用启用了该功能),点击标签即可快速发送

### 停止生成

在流式输出过程中,可以点击"停止生成"按钮中断

## 技术栈

- **Vue 3** - 框架
- **TypeScript** - 类型支持
- **Ant Design Vue** - UI 组件库
- **Pinia** - 状态管理
- **Marked** - Markdown 渲染
- **Day.js** - 时间处理
- **Fetch API** - 流式请求

## 注意事项

### 1. Lint 错误

当前代码存在一些 ESLint 格式错误,主要是:
- 缩进问题
- 引号风格
- 对象属性换行

**解决方案**: 运行 `pnpm lint:fix` 自动修复

### 2. 类型错误

`useDifyChat.ts` 中的 `userStore.userInfo` 可能需要根据实际的 store 结构调整

**解决方案**: 检查 `src/store/modules/user.ts` 中的用户信息字段名

### 3. API 响应类型

需要确保后端返回的数据格式与前端定义的类型一致

### 4. 流式请求

流式请求使用原生 Fetch API,不经过 axios 封装,需要手动处理 token

## 后续优化建议

### 功能优化

1. **消息编辑**: 支持编辑已发送的消息
2. **消息复制**: 一键复制 AI 回复
3. **导出对话**: 导出为 Markdown 或 PDF
4. **语音输入**: 支持语音转文字
5. **图片上传**: 支持上传图片到对话中

### 性能优化

1. **虚拟滚动**: 消息列表过长时使用虚拟滚动
2. **防抖优化**: 输入框添加防抖
3. **缓存策略**: 缓存会话列表和消息历史
4. **懒加载**: 消息历史分页加载

### UI 优化

1. **暗黑模式**: 支持暗黑主题
2. **自定义主题**: 支持自定义颜色
3. **快捷键**: 添加键盘快捷键
4. **拖拽调整**: 支持调整侧边栏宽度

## 文件结构

```
src/
├── api/
│   └── dify.ts                 # API 接口封装
├── types/
│   └── dify.ts                 # 类型定义
├── composables/
│   └── useDifyChat.ts          # 组合式函数
├── router/
│   └── modules/
│       └── dify.ts             # 路由配置
└── views/
    └── dify/
        ├── chat/
        │   └── index.vue       # AI 对话页面
        └── knowledge/
            └── index.vue       # 知识库管理页面(待实现)
```

## 开发团队

- 后端接口: ✅ 已完成
- 前端页面: ✅ 已完成
- 测试: ⏳ 待进行
- 文档: ✅ 已完成

## 相关文档

- [后端 API 文档](../../art_admin_backend/DIFY_CONVERSATION_API.md)
- [后端测试报告](../../art_admin_backend/NEW_API_TEST_RESULTS.md)
- [Dify 官方文档](https://docs.dify.ai/)

---

**状态**: 开发完成,待测试和格式修复  
**更新时间**: 2024-11-20
