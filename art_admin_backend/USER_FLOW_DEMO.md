# AI 学习系统 - 用户操作完整流程演示

## 📋 流程概述

用户输入 "一元二次方程" → Dify检索知识库 → DeepSeek分析整理 → 返回结构化教材

---

## 🎬 完整流程演示

### 步骤 1: 👤 用户输入

**前端界面操作：**
```
用户选择:
- 学科: 数学 📐
- 年级: 初中
- 题材: 一元二次方程
- 难度: 进阶 (2)

点击 [生成教材] 按钮
```

**发送的 API 请求：**
```bash
POST /api/learning/material/generate
Authorization: Bearer {JWT_TOKEN}
Content-Type: application/json

{
  "subject_id": 2,
  "grade": "初中",
  "topic": "一元二次方程",
  "difficulty": 2
}
```

---

### 步骤 2: 🔍 后端接收 → Dify 知识库检索

**代码位置：** `internal/service/learning_material_service.go`

```go
// 1. 获取学科信息
subject, _ := s.subjectRepo.GetByID(ctx, req.SubjectID)
// subject.Name = "数学"

// 2. 构建检索查询
query := fmt.Sprintf("%s %s %s 知识点", subject.Name, req.Grade, req.Topic)
// query = "数学 初中 一元二次方程 知识点"

// 3. 调用 Dify API 检索知识库
knowledgeContents, _ := s.difyClient.SimpleRetrieve(ctx, s.difyDatasetID, query)
```

**实际的 Dify API 调用：**
```bash
POST https://api.dify.ai/v1/datasets/7e2a67c6-00ab-46d0-a71b-886bb480d18f/retrieve
Authorization: Bearer dataset-QQ8Y3Cxvss4RvQXvU7Gem4sb
Content-Type: application/json

{
  "query": "数学 初中 一元二次方程 知识点"
}
```

**Dify 返回的知识片段：**
```json
{
  "query": {
    "content": "数学 初中 一元二次方程 知识点"
  },
  "records": [
    {
      "segment": {
        "id": "xxx",
        "content": "一元二次方程的标准形式是 ax²+bx+c=0 (a≠0)。解法包括：1. 因式分解法 2. 配方法 3. 公式法 4. 图像法...",
        "position": 1,
        "document_id": "xxx"
      },
      "score": 0.95
    },
    {
      "segment": {
        "content": "判别式 Δ=b²-4ac 决定方程根的情况：Δ>0 两个不相等实根，Δ=0 两个相等实根，Δ<0 无实根...",
        "position": 2
      },
      "score": 0.88
    }
  ]
}
```

**提取知识内容：**
```
知识片段 1: 一元二次方程的标准形式是 ax²+bx+c=0...
知识片段 2: 判别式 Δ=b²-4ac 决定方程根的情况...
```

---

### 步骤 3: 🤖 构建 RAG Prompt

**代码位置：** `internal/service/learning_material_service.go`

```go
// 构建系统提示词
systemPrompt := s.buildSystemPrompt(
    subject.Name,      // "数学"
    req.Grade,         // "初中"
    req.Topic,         // "一元二次方程"
    req.Difficulty,    // 2
    knowledgeContents  // 从 Dify 检索到的内容
)
```

**生成的完整 Prompt：**
```
你是一位专业的数学老师，擅长为初中学生设计教学内容。

请根据以下参考资料，为学生生成一份关于"一元二次方程"的教材内容。

【参考资料】
一元二次方程的标准形式是 ax²+bx+c=0 (a≠0)。解法包括：
1. 因式分解法
2. 配方法
3. 公式法
4. 图像法

判别式 Δ=b²-4ac 决定方程根的情况：
- Δ>0: 两个不相等实根
- Δ=0: 两个相等实根
- Δ<0: 无实根

【要求】
1. 内容符合初中学生的认知水平
2. 结构清晰，包含以下部分：
   - 知识点讲解（清晰易懂）
   - 重点难点标注
   - 例题演示（至少2个）
   - 练习题（至少5个，包含答案和解析）
3. 难度等级：进阶
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
      "key_points": ["重点1", "重点2"],
      "difficulty": "进阶"
    },
    ...
  ],
  "total_time": 45
}
```

---

### 步骤 4: 🧠 DeepSeek AI 分析生成

**代码位置：** `internal/service/learning_material_service.go`

```go
// 调用 DeepSeek API
response, _ := s.deepseekClient.SimpleChat(ctx, systemPrompt, userMessage)
```

**实际的 DeepSeek API 调用：**
```bash
POST https://api.deepseek.com/chat/completions
Authorization: Bearer sk-a8e4ee88516f40e6a2dc3776d3254846
Content-Type: application/json

{
  "model": "deepseek-chat",
  "messages": [
    {
      "role": "system",
      "content": "你是一位专业的数学老师..."  // 上面构建的完整 Prompt
    },
    {
      "role": "user",
      "content": "请为我生成关于'一元二次方程'的教材内容"
    }
  ],
  "temperature": 0.7,
  "max_tokens": 4096
}
```

**DeepSeek 返回的内容：**
```json
{
  "id": "xxx",
  "choices": [
    {
      "message": {
        "role": "assistant",
        "content": "```json\n{\n  \"title\": \"一元二次方程完全攻略\",\n  \"summary\": \"本教材将带你深入理解一元二次方程...\",\n  \"sections\": [\n    {\n      \"type\": \"knowledge\",\n      \"title\": \"一元二次方程的定义与标准形式\",\n      \"content\": \"一元二次方程是指只含有一个未知数，并且未知数的最高次数是2的整式方程...\",\n      \"key_points\": [\n        \"标准形式: ax²+bx+c=0 (a≠0)\",\n        \"二次项系数a不能为0\",\n        \"可以有一次项和常数项\"\n      ],\n      \"difficulty\": \"进阶\"\n    },\n    {\n      \"type\": \"knowledge\",\n      \"title\": \"一元二次方程的解法\",\n      \"content\": \"解一元二次方程有四种主要方法...\",\n      \"key_points\": [\n        \"因式分解法：适用于易分解的方程\",\n        \"配方法：万能方法\",\n        \"公式法：最常用\",\n        \"图像法：直观理解\"\n      ]\n    },\n    {\n      \"type\": \"example\",\n      \"title\": \"例题1：用因式分解法解方程\",\n      \"question\": \"解方程：x² - 5x + 6 = 0\",\n      \"solution\": \"步骤1: 分解因式 (x-2)(x-3) = 0\\n步骤2: 得 x-2=0 或 x-3=0\\n步骤3: 解得 x₁=2, x₂=3\",\n      \"answer\": \"x₁=2, x₂=3\"\n    },\n    {\n      \"type\": \"example\",\n      \"title\": \"例题2：用公式法解方程\",\n      \"question\": \"解方程：2x² + 3x - 2 = 0\",\n      \"solution\": \"步骤1: 确定 a=2, b=3, c=-2\\n步骤2: 计算判别式 Δ=b²-4ac=9+16=25\\n步骤3: 代入公式 x=(-b±√Δ)/(2a)\\n步骤4: x=(-3±5)/4\",\n      \"answer\": \"x₁=0.5, x₂=-2\"\n    },\n    {\n      \"type\": \"exercise\",\n      \"title\": \"练习题\",\n      \"questions\": [\n        {\n          \"id\": 1,\n          \"type\": \"choice\",\n          \"question\": \"方程 x²-4=0 的解是？\",\n          \"options\": [\"x=2\", \"x=-2\", \"x=±2\", \"x=4\"],\n          \"answer\": \"C\",\n          \"explanation\": \"x²=4，开平方得 x=±2\"\n        },\n        {\n          \"id\": 2,\n          \"type\": \"choice\",\n          \"question\": \"判别式 Δ<0 时，方程有几个实根？\",\n          \"options\": [\"0个\", \"1个\", \"2个\", \"无数个\"],\n          \"answer\": \"A\",\n          \"explanation\": \"Δ<0 时，方程无实数根\"\n        },\n        {\n          \"id\": 3,\n          \"type\": \"fill\",\n          \"question\": \"方程 x²-6x+9=0 可以分解为 (x-__)²=0\",\n          \"answer\": \"3\",\n          \"explanation\": \"x²-6x+9 = (x-3)²\"\n        },\n        {\n          \"id\": 4,\n          \"type\": \"calculate\",\n          \"question\": \"解方程：x²+2x-3=0\",\n          \"answer\": \"x₁=1, x₂=-3\",\n          \"explanation\": \"因式分解：(x+3)(x-1)=0\"\n        },\n        {\n          \"id\": 5,\n          \"type\": \"calculate\",\n          \"question\": \"解方程：3x²-5x+2=0\",\n          \"answer\": \"x₁=1, x₂=2/3\",\n          \"explanation\": \"用公式法：x=(-b±√Δ)/(2a)\"\n        }\n      ]\n    }\n  ],\n  \"total_time\": 50\n}\n```"
      }
    }
  ],
  "usage": {
    "prompt_tokens": 856,
    "completion_tokens": 1234,
    "total_tokens": 2090
  }
}
```

---

### 步骤 5: 📊 解析和整理数据

**代码位置：** `internal/service/learning_material_service.go`

```go
// 1. 提取 JSON 内容（去除 markdown 代码块）
response = strings.TrimPrefix(response, "```json")
response = response[:strings.LastIndex(response, "```")]

// 2. 解析 JSON
var content model.MaterialContent
json.Unmarshal([]byte(response), &content)

// 3. 构建教材对象
material := &model.LearningMaterial{
    UserID:     req.UserID,
    SubjectID:  req.SubjectID,
    Title:      content.Title,        // "一元二次方程完全攻略"
    Topic:      req.Topic,            // "一元二次方程"
    Grade:      req.Grade,            // "初中"
    Difficulty: req.Difficulty,       // 2
    Summary:    content.Summary,      // "本教材将带你深入理解..."
    Content:    content,              // 完整的结构化内容
    TotalTime:  content.TotalTime,    // 50
    Status:     1,
}
```

**整理后的数据结构：**
```json
{
  "id": 123,
  "user_id": 1,
  "subject_id": 2,
  "title": "一元二次方程完全攻略",
  "topic": "一元二次方程",
  "grade": "初中",
  "difficulty": 2,
  "summary": "本教材将带你深入理解一元二次方程...",
  "content": {
    "title": "一元二次方程完全攻略",
    "summary": "...",
    "sections": [
      // 2个知识点 + 2个例题 + 5道练习题
    ],
    "total_time": 50
  },
  "total_time": 50,
  "view_count": 0,
  "favorite_count": 0,
  "created_at": "2024-11-21 14:05:00"
}
```

---

### 步骤 6: 💾 保存到数据库

**代码位置：** `internal/service/learning_material_service.go`

```go
// 保存到数据库
s.repo.Create(ctx, material)
```

**SQL 操作：**
```sql
INSERT INTO learning_materials (
    user_id, subject_id, title, topic, grade, difficulty,
    summary, content, total_time, status, created_at
) VALUES (
    1, 2, '一元二次方程完全攻略', '一元二次方程', '初中', 2,
    '本教材将带你深入理解...', 
    '{"title":"一元二次方程完全攻略","sections":[...]}',
    50, 1, NOW()
);
```

---

### 步骤 7: 🎨 返回给前端展示

**API 响应：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 123,
    "user_id": 1,
    "subject_id": 2,
    "title": "一元二次方程完全攻略",
    "topic": "一元二次方程",
    "grade": "初中",
    "difficulty": 2,
    "summary": "本教材将带你深入理解一元二次方程，从定义到解法，从理论到实践...",
    "content": {
      "title": "一元二次方程完全攻略",
      "summary": "...",
      "sections": [
        {
          "type": "knowledge",
          "title": "一元二次方程的定义与标准形式",
          "content": "一元二次方程是指只含有一个未知数...",
          "key_points": [
            "标准形式: ax²+bx+c=0 (a≠0)",
            "二次项系数a不能为0",
            "可以有一次项和常数项"
          ],
          "difficulty": "进阶"
        },
        {
          "type": "knowledge",
          "title": "一元二次方程的解法",
          "content": "解一元二次方程有四种主要方法...",
          "key_points": [
            "因式分解法：适用于易分解的方程",
            "配方法：万能方法",
            "公式法：最常用",
            "图像法：直观理解"
          ]
        },
        {
          "type": "example",
          "title": "例题1：用因式分解法解方程",
          "question": "解方程：x² - 5x + 6 = 0",
          "solution": "步骤1: 分解因式 (x-2)(x-3) = 0\n步骤2: 得 x-2=0 或 x-3=0\n步骤3: 解得 x₁=2, x₂=3",
          "answer": "x₁=2, x₂=3"
        },
        {
          "type": "example",
          "title": "例题2：用公式法解方程",
          "question": "解方程：2x² + 3x - 2 = 0",
          "solution": "...",
          "answer": "x₁=0.5, x₂=-2"
        },
        {
          "type": "exercise",
          "title": "练习题",
          "questions": [
            {
              "id": 1,
              "type": "choice",
              "question": "方程 x²-4=0 的解是？",
              "options": ["x=2", "x=-2", "x=±2", "x=4"],
              "answer": "C",
              "explanation": "x²=4，开平方得 x=±2"
            }
            // ... 共5道题
          ]
        }
      ],
      "total_time": 50
    },
    "audio_url": "",
    "audio_duration": 0,
    "total_time": 50,
    "view_count": 0,
    "favorite_count": 0,
    "created_at": "2024-11-21 14:05:00",
    "updated_at": "2024-11-21 14:05:00"
  }
}
```

**前端展示效果：**
```
┌─────────────────────────────────────────┐
│  📚 一元二次方程完全攻略                │
│  ⏱️  预计学习时长: 50分钟               │
│  📊 难度: 进阶                          │
├─────────────────────────────────────────┤
│  📖 内容概要                            │
│  本教材将带你深入理解一元二次方程...    │
├─────────────────────────────────────────┤
│  💡 知识点讲解                          │
│  1. 一元二次方程的定义与标准形式        │
│     ⭐ 标准形式: ax²+bx+c=0 (a≠0)      │
│     ⭐ 二次项系数a不能为0               │
│                                         │
│  2. 一元二次方程的解法                  │
│     ⭐ 因式分解法：适用于易分解的方程   │
│     ⭐ 配方法：万能方法                 │
│     ⭐ 公式法：最常用                   │
├─────────────────────────────────────────┤
│  📝 例题演示                            │
│  例题1: 用因式分解法解方程              │
│  题目: 解方程：x² - 5x + 6 = 0         │
│  [查看解答]                             │
│                                         │
│  例题2: 用公式法解方程                  │
│  题目: 解方程：2x² + 3x - 2 = 0        │
│  [查看解答]                             │
├─────────────────────────────────────────┤
│  ✏️  练习题 (5道)                       │
│  1. [选择题] 方程 x²-4=0 的解是？      │
│  2. [选择题] 判别式 Δ<0 时...          │
│  3. [填空题] 方程 x²-6x+9=0...         │
│  4. [计算题] 解方程：x²+2x-3=0         │
│  5. [计算题] 解方程：3x²-5x+2=0        │
│  [开始练习]                             │
└─────────────────────────────────────────┘
```

---

## 🔄 完整数据流

```
用户输入
  ↓ (HTTP POST)
后端 API
  ↓ (查询数据库)
获取学科信息
  ↓ (构建查询)
Dify 知识库检索
  ↓ (返回知识片段)
构建 RAG Prompt
  ↓ (HTTP POST)
DeepSeek AI 生成
  ↓ (返回 JSON)
解析和整理
  ↓ (INSERT)
保存到数据库
  ↓ (HTTP Response)
返回给前端
  ↓ (渲染)
用户看到教材
```

---

## 📊 关键技术点

### 1. RAG (Retrieval-Augmented Generation)
- **检索**: Dify 向量数据库语义检索
- **增强**: 将检索结果注入到 Prompt 中
- **生成**: DeepSeek 基于增强的上下文生成内容

### 2. Prompt Engineering
- 明确的角色定位（数学老师）
- 结构化的输出要求（JSON 格式）
- 具体的内容要求（知识点+例题+练习）
- 适当的难度控制（基础/进阶/高级）

### 3. 数据处理
- JSON 解析（处理 markdown 代码块）
- 数据验证（确保必填字段）
- 错误处理（容错机制）

---

## 🎯 流程总结

| 步骤 | 操作 | 耗时 | 关键点 |
|------|------|------|--------|
| 1 | 用户输入 | < 1s | 前端表单验证 |
| 2 | Dify 检索 | 2-3s | 语义相似度匹配 |
| 3 | 构建 Prompt | < 1s | RAG 增强 |
| 4 | DeepSeek 生成 | 30-60s | AI 内容生成 |
| 5 | 数据解析 | < 1s | JSON 处理 |
| 6 | 保存数据库 | < 1s | 持久化存储 |
| 7 | 返回前端 | < 1s | HTTP 响应 |
| **总计** | **完整流程** | **~40-70s** | **端到端** |

---

## 🚀 优化方向

1. **缓存机制** - 缓存常见题材的教材
2. **流式输出** - 实时显示生成进度
3. **并行处理** - 同时检索多个知识库
4. **智能推荐** - 基于用户历史推荐相关内容

---

**文档创建时间:** 2024-11-21  
**系统版本:** v1.0.0  
**状态:** ✅ 已验证
