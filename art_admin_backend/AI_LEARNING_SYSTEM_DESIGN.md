# AI 在线学习系统 - 开发文档

## 技术方案

### 核心架构

```
Dify 知识库（向量检索）+ 自研 RAG + DeepSeek API + 阿里云 TTS
```

**优势：**
- ✅ 利用 Dify 的向量检索能力，无需自建向量数据库
- ✅ 自己实现 RAG 逻辑，完全可控
- ✅ 使用 DeepSeek API，成本低质量高
- ✅ 集成 TTS 服务，实现音频播放

---

## 系统架构图

```
┌─────────────────────────────────────────────────────────────┐
│                    前端 (Vue 3 + TypeScript)                 │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│  │学科选择  │  │教材生成  │  │音频播放  │  │收藏管理  │   │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘   │
└────────────────────────┬────────────────────────────────────┘
                         │ HTTP/WebSocket
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                  后端 API (Go + Gin)                         │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  业务逻辑层 (Service)                                  │  │
│  │  - 教材生成服务                                        │  │
│  │  - 知识库检索服务                                      │  │
│  │  - 音频生成服务                                        │  │
│  │  - 收藏管理服务                                        │  │
│  └──────────────────────────────────────────────────────┘  │
└────────────────────────┬────────────────────────────────────┘
                         │
         ┌───────────────┼───────────────┬──────────────┐
         ▼               ▼               ▼              ▼
┌──────────────┐ ┌──────────────┐ ┌──────────┐ ┌──────────────┐
│ Dify 知识库   │ │ DeepSeek API │ │阿里云TTS │ │ MySQL 数据库 │
│ (向量检索)    │ │ (文本生成)   │ │(语音合成)│ │ (数据存储)   │
└──────────────┘ └──────────────┘ └──────────┘ └──────────────┘
```

---

## 核心功能模块

### 1. 学科管理模块

**功能：**
- 学科分类管理（语文、数学、英语等）
- 年级设置（小学、初中、高中）
- 难度等级配置

**数据库表：**
```sql
CREATE TABLE subjects (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(50) NOT NULL COMMENT '学科名称',
  icon VARCHAR(200) COMMENT '学科图标',
  description TEXT COMMENT '学科描述',
  grade_levels JSON COMMENT '适用年级 ["小学","初中","高中"]',
  sort INT DEFAULT 0 COMMENT '排序',
  status TINYINT DEFAULT 1 COMMENT '状态 1-启用 0-禁用',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_status (status),
  INDEX idx_sort (sort)
);
```

---

### 2. 教材生成模块（核心）

#### 2.1 工作流程

```
用户输入（学科 + 年级 + 题材）
    ↓
① 调用 Dify 知识库 API 检索相关教学资料
    ↓
② 构建 RAG Prompt（系统提示词 + 检索结果 + 用户输入）
    ↓
③ 调用 DeepSeek API 生成结构化教材内容
    ↓
④ 解析并保存教材内容到数据库
    ↓
⑤ 异步生成音频（调用 TTS 服务）
    ↓
返回教材内容给前端
```

#### 2.2 Dify 知识库检索

**API 接口：** `POST https://api.dify.ai/v1/datasets/{dataset_id}/retrieve`

**请求示例：**
```json
{
  "query": "初中数学 二次方程 知识点",
  "retrieval_model": {
    "search_method": "semantic_search",
    "reranking_enable": true,
    "reranking_model": {
      "reranking_provider_name": "cohere",
      "reranking_model_name": "rerank-multilingual-v2.0"
    },
    "top_k": 5,
    "score_threshold_enabled": true,
    "score_threshold": 0.5
  }
}
```

**响应示例：**
```json
{
  "query": "初中数学 二次方程 知识点",
  "records": [
    {
      "segment": {
        "content": "二次方程的定义：形如 ax²+bx+c=0 的方程...",
        "position": 1,
        "score": 0.95
      },
      "metadata": {
        "document_name": "初中数学教学大纲.pdf",
        "page": 15
      }
    }
  ]
}
```

#### 2.3 RAG Prompt 构建

**系统提示词模板：**
```
你是一位专业的{学科}老师，擅长为{年级}学生设计教学内容。

请根据以下参考资料，为学生生成一份关于"{题材}"的教材内容。

【参考资料】
{从 Dify 检索到的知识片段}

【要求】
1. 内容符合{年级}学生的认知水平
2. 结构清晰，包含以下部分：
   - 知识点讲解（清晰易懂）
   - 重点难点标注
   - 例题演示（至少2个）
   - 练习题（至少5个，包含答案和解析）
3. 难度等级：{难度}
4. 语言生动有趣，激发学习兴趣
5. 必须输出 JSON 格式

【输出格式】
请严格按照以下 JSON 格式输出：
{
  "title": "教材标题",
  "summary": "内容概要",
  "sections": [
    {
      "type": "knowledge",
      "title": "知识点标题",
      "content": "详细讲解内容",
      "keyPoints": ["重点1", "重点2"],
      "difficulty": "基础/进阶/高级"
    },
    {
      "type": "example",
      "title": "例题标题",
      "question": "题目内容",
      "solution": "解题步骤",
      "answer": "答案"
    },
    {
      "type": "exercise",
      "title": "练习题",
      "questions": [
        {
          "id": 1,
          "type": "choice|fill|calculate|essay",
          "question": "题目内容",
          "options": ["A选项", "B选项"],  // 选择题才有
          "answer": "答案",
          "explanation": "解析"
        }
      ]
    }
  ],
  "totalTime": "预计学习时长（分钟）"
}
```

#### 2.4 DeepSeek API 调用

**请求示例：**
```json
{
  "model": "deepseek-chat",
  "messages": [
    {
      "role": "system",
      "content": "{系统提示词}"
    },
    {
      "role": "user",
      "content": "请为我生成关于'{题材}'的教材内容"
    }
  ],
  "temperature": 0.7,
  "max_tokens": 4096
}
```

#### 2.5 数据库表设计

```sql
-- 教材内容表
CREATE TABLE learning_materials (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT NOT NULL COMMENT '创建用户ID',
  subject_id BIGINT NOT NULL COMMENT '学科ID',
  title VARCHAR(200) NOT NULL COMMENT '教材标题',
  topic VARCHAR(200) NOT NULL COMMENT '学习题材',
  grade VARCHAR(50) NOT NULL COMMENT '年级',
  difficulty TINYINT DEFAULT 1 COMMENT '难度 1-基础 2-进阶 3-高级',
  summary TEXT COMMENT '内容概要',
  content LONGTEXT NOT NULL COMMENT '教材内容(JSON格式)',
  audio_url VARCHAR(500) COMMENT '音频URL',
  audio_duration INT COMMENT '音频时长(秒)',
  total_time INT COMMENT '预计学习时长(分钟)',
  view_count INT DEFAULT 0 COMMENT '浏览次数',
  favorite_count INT DEFAULT 0 COMMENT '收藏次数',
  status TINYINT DEFAULT 1 COMMENT '状态 1-正常 0-已删除',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_user_id (user_id),
  INDEX idx_subject_id (subject_id),
  INDEX idx_created_at (created_at),
  INDEX idx_status (status)
);

-- 教材章节表（用于存储解析后的结构化内容）
CREATE TABLE material_sections (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  material_id BIGINT NOT NULL COMMENT '教材ID',
  type VARCHAR(20) NOT NULL COMMENT '类型 knowledge|example|exercise',
  title VARCHAR(200) COMMENT '章节标题',
  content LONGTEXT COMMENT '章节内容(JSON格式)',
  sort INT DEFAULT 0 COMMENT '排序',
  audio_url VARCHAR(500) COMMENT '章节音频URL',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_material_id (material_id),
  INDEX idx_type (type),
  INDEX idx_sort (sort)
);
```

---

### 3. 音频播放模块

#### 3.1 TTS 服务选择

**推荐：阿里云语音合成**

**优势：**
- 中文效果优秀
- 价格合理（¥4/万字符）
- 支持多种音色
- 稳定性好

**API 文档：** https://help.aliyun.com/document_detail/84435.html

#### 3.2 音频生成流程

```
教材内容保存成功
    ↓
异步任务：提取需要转音频的文本
    ↓
检查音频缓存（基于内容哈希）
    ↓
如果缓存不存在：
  - 调用阿里云 TTS API
  - 上传音频到 OSS
  - 保存音频 URL 到数据库
  - 更新缓存
    ↓
返回音频 URL
```

#### 3.3 音频缓存表

```sql
CREATE TABLE audio_cache (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  content_hash VARCHAR(64) NOT NULL COMMENT '内容MD5哈希',
  text_content TEXT COMMENT '原始文本内容',
  audio_url VARCHAR(500) NOT NULL COMMENT '音频URL',
  voice_type VARCHAR(50) DEFAULT 'xiaoyun' COMMENT '音色类型',
  language VARCHAR(20) DEFAULT 'zh-CN' COMMENT '语言',
  file_size BIGINT COMMENT '文件大小(字节)',
  duration INT COMMENT '时长(秒)',
  access_count INT DEFAULT 0 COMMENT '访问次数',
  last_access_at TIMESTAMP COMMENT '最后访问时间',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uk_content_hash (content_hash),
  INDEX idx_last_access (last_access_at)
);
```

#### 3.4 前端音频播放器

**功能要求：**
- 播放/暂停
- 进度条拖拽
- 倍速播放（0.5x, 1x, 1.5x, 2x）
- 音量调节
- 单句播放/全文播放
- 自动播放下一段

---

### 4. 收藏功能模块

#### 4.1 收藏夹管理

```sql
-- 收藏夹表
CREATE TABLE favorite_folders (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT NOT NULL COMMENT '用户ID',
  name VARCHAR(100) NOT NULL COMMENT '收藏夹名称',
  icon VARCHAR(200) COMMENT '收藏夹图标',
  description VARCHAR(500) COMMENT '收藏夹描述',
  item_count INT DEFAULT 0 COMMENT '收藏数量',
  sort INT DEFAULT 0 COMMENT '排序',
  is_default TINYINT DEFAULT 0 COMMENT '是否默认收藏夹',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_user_id (user_id),
  INDEX idx_sort (sort)
);

-- 收藏表
CREATE TABLE favorites (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT NOT NULL COMMENT '用户ID',
  folder_id BIGINT NOT NULL COMMENT '收藏夹ID',
  resource_type TINYINT NOT NULL COMMENT '资源类型 1-教材 2-知识点 3-练习题 4-笔记',
  resource_id BIGINT NOT NULL COMMENT '资源ID',
  tags JSON COMMENT '标签',
  note TEXT COMMENT '收藏备注',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uk_user_resource (user_id, resource_type, resource_id),
  INDEX idx_user_id (user_id),
  INDEX idx_folder_id (folder_id),
  INDEX idx_resource (resource_type, resource_id)
);
```

#### 4.2 功能列表

- 创建/编辑/删除收藏夹
- 收藏/取消收藏
- 移动收藏到其他收藏夹
- 批量操作
- 按标签筛选
- 导出收藏内容

---

### 5. 练习测试模块

#### 5.1 练习题表

```sql
CREATE TABLE exercises (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  material_id BIGINT COMMENT '关联教材ID',
  subject_id BIGINT NOT NULL COMMENT '学科ID',
  type TINYINT NOT NULL COMMENT '题目类型 1-选择 2-填空 3-简答 4-计算',
  question TEXT NOT NULL COMMENT '题目内容',
  options JSON COMMENT '选项(选择题)',
  answer TEXT NOT NULL COMMENT '答案',
  explanation TEXT COMMENT '解析',
  difficulty TINYINT DEFAULT 1 COMMENT '难度 1-基础 2-进阶 3-高级',
  knowledge_points JSON COMMENT '知识点标签',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_material_id (material_id),
  INDEX idx_subject_id (subject_id),
  INDEX idx_type (type),
  INDEX idx_difficulty (difficulty)
);

-- 错题本表
CREATE TABLE wrong_questions (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT NOT NULL COMMENT '用户ID',
  exercise_id BIGINT NOT NULL COMMENT '练习题ID',
  user_answer TEXT COMMENT '用户答案',
  correct_answer TEXT COMMENT '正确答案',
  times INT DEFAULT 1 COMMENT '错误次数',
  mastered TINYINT DEFAULT 0 COMMENT '是否掌握 1-是 0-否',
  last_practice_at TIMESTAMP COMMENT '最后练习时间',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_user_exercise (user_id, exercise_id),
  INDEX idx_user_id (user_id),
  INDEX idx_mastered (mastered)
);
```

#### 5.2 AI 批改功能

**流程：**
1. 用户提交答案
2. 简单题型（选择、填空）：直接比对答案
3. 复杂题型（简答、计算）：调用 DeepSeek API 批改
4. 返回批改结果和建议

**批改 Prompt：**
```
你是一位专业的{学科}老师，请批改以下学生答案。

【题目】
{题目内容}

【标准答案】
{标准答案}

【学生答案】
{学生答案}

【要求】
1. 判断答案是否正确
2. 指出错误之处
3. 给出改进建议
4. 评分（0-100分）

请以 JSON 格式输出：
{
  "correct": true/false,
  "score": 85,
  "errors": ["错误点1", "错误点2"],
  "suggestions": ["建议1", "建议2"],
  "comment": "总体评价"
}
```

---

### 6. 笔记功能模块

```sql
CREATE TABLE notes (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT NOT NULL COMMENT '用户ID',
  material_id BIGINT COMMENT '关联教材ID',
  section_id BIGINT COMMENT '关联章节ID',
  exercise_id BIGINT COMMENT '关联练习题ID',
  title VARCHAR(200) COMMENT '笔记标题',
  content LONGTEXT NOT NULL COMMENT '笔记内容',
  content_type TINYINT DEFAULT 1 COMMENT '内容类型 1-富文本 2-Markdown',
  tags JSON COMMENT '标签',
  images JSON COMMENT '图片URL数组',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_user_id (user_id),
  INDEX idx_material_id (material_id),
  INDEX idx_created_at (created_at)
);
```

---

## API 接口设计

### 1. 教材生成相关

#### 1.1 生成教材
```
POST /api/learning/material/generate

Request:
{
  "subject_id": 1,
  "grade": "初中",
  "topic": "二次方程的解法",
  "difficulty": 2
}

Response:
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 123,
    "title": "二次方程的解法",
    "summary": "...",
    "content": {...},
    "audio_url": "https://...",
    "total_time": 45
  }
}
```

#### 1.2 获取教材列表
```
GET /api/learning/material/list?page=1&limit=20&subject_id=1

Response:
{
  "code": 200,
  "data": {
    "list": [...],
    "total": 100,
    "page": 1,
    "limit": 20
  }
}
```

#### 1.3 获取教材详情
```
GET /api/learning/material/:id

Response:
{
  "code": 200,
  "data": {
    "id": 123,
    "title": "...",
    "content": {...},
    "sections": [...]
  }
}
```

### 2. 收藏相关

#### 2.1 收藏/取消收藏
```
POST /api/learning/favorite/toggle

Request:
{
  "folder_id": 1,
  "resource_type": 1,
  "resource_id": 123
}
```

#### 2.2 获取收藏列表
```
GET /api/learning/favorite/list?folder_id=1&page=1&limit=20
```

### 3. 练习相关

#### 3.1 提交答案
```
POST /api/learning/exercise/submit

Request:
{
  "exercise_id": 456,
  "answer": "用户答案"
}

Response:
{
  "code": 200,
  "data": {
    "correct": true,
    "score": 100,
    "explanation": "..."
  }
}
```

#### 3.2 获取错题本
```
GET /api/learning/wrong-questions?page=1&limit=20&mastered=0
```

---

## 前端页面设计

### 1. 页面结构

```
/learning
  ├── /subjects              # 学科选择页
  ├── /generate              # 教材生成页
  ├── /material/:id          # 教材详情页
  ├── /favorites             # 收藏夹页
  ├── /exercises             # 练习题页
  ├── /wrong-questions       # 错题本页
  └── /notes                 # 笔记页
```

### 2. 核心组件

```
components/
  ├── LearningMaterialCard.vue      # 教材卡片
  ├── AudioPlayer.vue               # 音频播放器
  ├── ExerciseItem.vue              # 练习题组件
  ├── MarkdownEditor.vue            # Markdown 编辑器
  ├── FavoriteButton.vue            # 收藏按钮
  └── KnowledgePointTag.vue         # 知识点标签
```

---

## 开发计划

### Phase 1: 基础功能（2周）

**Week 1:**
- ✅ 数据库表设计和创建
- ✅ DeepSeek 客户端封装
- ✅ Dify 知识库检索集成
- ✅ 教材生成核心逻辑

**Week 2:**
- ✅ 教材生成 API 接口
- ✅ 前端教材生成页面
- ✅ 前端教材详情页面
- ✅ 基础音频播放功能

### Phase 2: 增强功能（2周）

**Week 3:**
- ✅ TTS 服务集成
- ✅ 音频缓存机制
- ✅ 收藏功能完整实现
- ✅ 收藏夹管理

**Week 4:**
- ✅ 练习题生成和展示
- ✅ AI 批改功能
- ✅ 错题本功能
- ✅ 笔记功能

### Phase 3: 优化和完善（1周）

**Week 5:**
- ✅ 性能优化
- ✅ 缓存策略优化
- ✅ UI/UX 优化
- ✅ 测试和 Bug 修复

---

## 成本估算

### 月度运营成本

| 项目 | 说明 | 月成本 |
|------|------|--------|
| **Dify 知识库** | 云服务基础版 | ¥0-99 |
| **DeepSeek API** | 按 token 计费 | ¥200-500 |
| **阿里云 TTS** | 按字符计费 | ¥100-300 |
| **阿里云 OSS** | 音频存储 | ¥50-100 |
| **服务器** | 2核4G | ¥200-400 |
| **数据库** | MySQL | ¥0（自建） |
| **总计** | | **¥550-1400** |

### 成本优化建议

1. **音频缓存** - 相同内容只生成一次
2. **Prompt 优化** - 减少 token 消耗
3. **批量处理** - 合并 API 请求
4. **CDN 加速** - 减少 OSS 流量费用

---

## 技术栈总结

### 后端
- **语言**: Go 1.21+
- **框架**: Gin
- **数据库**: MySQL 8.0+
- **缓存**: Redis
- **AI**: DeepSeek API
- **知识库**: Dify Cloud
- **TTS**: 阿里云语音合成
- **存储**: 阿里云 OSS

### 前端
- **框架**: Vue 3 + TypeScript
- **UI**: Ant Design Vue / Element Plus
- **状态管理**: Pinia
- **路由**: Vue Router
- **HTTP**: Axios
- **富文本**: TinyMCE / Quill
- **Markdown**: marked + highlight.js
- **音频**: Howler.js / WaveSurfer.js
- **数学公式**: KaTeX

---

## 注意事项

### 1. API Key 安全
- 不要将 API Key 提交到版本控制
- 使用环境变量或配置文件
- 定期轮换 API Key

### 2. 内容质量控制
- 定期审核 AI 生成的内容
- 建立内容评分机制
- 收集用户反馈优化 Prompt

### 3. 性能优化
- 实现请求限流
- 添加缓存机制
- 异步处理耗时任务

### 4. 用户体验
- 流式输出提升响应速度
- 加载状态提示
- 错误处理和重试机制

---

## 相关文档

- [DeepSeek API 文档](https://platform.deepseek.com/api-docs/)
- [Dify 官方文档](https://docs.dify.ai/)
- [阿里云 TTS 文档](https://help.aliyun.com/document_detail/84435.html)
- [后端集成文档](./DIFY_INTEGRATION.md)

