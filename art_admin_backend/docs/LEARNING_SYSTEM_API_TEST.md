# AI 学习系统 API 测试文档

## 测试环境

- 后端地址: http://localhost:48080
- API 前缀: /api/learning
- 认证方式: JWT Bearer Token

## 前置准备

### 1. 初始化数据库

```bash
cd /Users/xinyoucai/code/art_admin/art_admin_backend
./scripts/init_learning_system.sh
```

### 2. 启动后端服务

```bash
go run cmd/server/main.go
```

### 3. 获取 JWT Token

```bash
# 登录获取 token
curl -X POST http://localhost:48080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "123456"
  }'

# 响应示例:
# {
#   "code": 200,
#   "data": {
#     "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
#     "token_type": "Bearer",
#     "expires_in": 7200
#   }
# }
```

将获取到的 `access_token` 替换到后续请求中的 `YOUR_JWT_TOKEN`。

---

## API 测试用例

### 1. 获取学科列表(无需认证)

**请求:**

```bash
curl -X GET http://localhost:48080/api/learning/subject/list
```

**预期响应:**

```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "语文",
      "icon": "📚",
      "description": "中文语言文学学科",
      "grade_levels": ["小学", "初中", "高中"],
      "sort": 1
    },
    {
      "id": 2,
      "name": "数学",
      "icon": "🔢",
      "description": "数学逻辑与计算学科",
      "grade_levels": ["小学", "初中", "高中"],
      "sort": 2
    }
    // ... 更多学科
  ]
}
```

---

### 2. 生成教材(需要认证)

**请求:**

```bash
curl -X POST http://localhost:48080/api/learning/material/generate \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "subject_id": 2,
    "grade": "初中",
    "topic": "二次方程的解法",
    "difficulty": 2
  }'
```

**参数说明:**

- `subject_id`: 学科 ID(1-语文, 2-数学, 3-英语...)
- `grade`: 年级(小学/初中/高中)
- `topic`: 学习题材
- `difficulty`: 难度等级(1-基础, 2-进阶, 3-高级)

**预期响应:**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "user_id": 1,
    "subject_id": 2,
    "title": "二次方程的解法",
    "topic": "二次方程的解法",
    "grade": "初中",
    "difficulty": 2,
    "summary": "本教材详细讲解二次方程的定义、解法和应用...",
    "content": {
      "title": "二次方程的解法",
      "summary": "...",
      "sections": [
        {
          "type": "knowledge",
          "title": "二次方程的定义",
          "content": "...",
          "key_points": ["重点1", "重点2"],
          "difficulty": "进阶"
        },
        {
          "type": "example",
          "title": "例题1",
          "question": "...",
          "solution": "...",
          "answer": "..."
        },
        {
          "type": "exercise",
          "title": "练习题",
          "questions": [
            {
              "id": 1,
              "type": "choice",
              "question": "...",
              "options": ["A", "B", "C", "D"],
              "answer": "A",
              "explanation": "..."
            }
          ]
        }
      ],
      "total_time": 45
    },
    "audio_url": "",
    "audio_duration": 0,
    "total_time": 45,
    "view_count": 0,
    "favorite_count": 0,
    "created_at": "2024-11-21 13:15:00",
    "updated_at": "2024-11-21 13:15:00"
  }
}
```

**可能的错误响应:**

```json
// 未登录
{
  "code": 401,
  "message": "未登录"
}

// 参数错误
{
  "code": 400,
  "message": "参数错误: Key: 'GenerateMaterialRequest.SubjectID' Error:Field validation for 'SubjectID' failed on the 'required' tag"
}

// 生成失败
{
  "code": 500,
  "message": "生成教材失败: ..."
}
```

---

### 3. 获取教材列表(需要认证)

**请求:**

```bash
curl -X GET "http://localhost:48080/api/learning/material/list?page=1&limit=10&subject_id=2" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**参数说明:**

- `page`: 页码(必填)
- `limit`: 每页数量(必填, 1-100)
- `subject_id`: 学科 ID(可选, 用于筛选)

**预期响应:**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "user_id": 1,
        "subject_id": 2,
        "title": "二次方程的解法",
        "topic": "二次方程的解法",
        "grade": "初中",
        "difficulty": 2,
        "summary": "...",
        "content": {...},
        "total_time": 45,
        "view_count": 5,
        "favorite_count": 2,
        "created_at": "2024-11-21 13:15:00",
        "updated_at": "2024-11-21 13:15:00"
      }
    ],
    "total": 1,
    "page": 1,
    "limit": 10
  }
}
```

---

### 4. 获取教材详情(需要认证)

**请求:**

```bash
curl -X GET http://localhost:48080/api/learning/material/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**预期响应:**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "user_id": 1,
    "subject_id": 2,
    "title": "二次方程的解法",
    "content": {...},
    "view_count": 6,
    // ... 完整教材信息
  }
}
```

**注意:** 每次访问详情会自动增加浏览次数。

---

### 5. 删除教材(需要认证)

**请求:**

```bash
curl -X DELETE http://localhost:48080/api/learning/material/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**预期响应:**

```json
{
  "code": 200,
  "message": "删除成功"
}
```

**可能的错误响应:**

```json
// 无权删除
{
  "code": 500,
  "message": "删除失败: 无权删除该教材"
}

// 教材不存在
{
  "code": 500,
  "message": "删除失败: 获取教材失败: record not found"
}
```

---

## 测试场景

### 场景 1: 完整流程测试

1. 获取学科列表
2. 选择数学学科(ID=2)
3. 生成初中数学教材
4. 查看教材列表
5. 查看教材详情
6. 删除教材

### 场景 2: 不同难度测试

生成同一题材的不同难度教材:

```bash
# 基础难度
curl -X POST http://localhost:48080/api/learning/material/generate \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"subject_id": 2, "grade": "初中", "topic": "二次方程", "difficulty": 1}'

# 进阶难度
curl -X POST http://localhost:48080/api/learning/material/generate \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"subject_id": 2, "grade": "初中", "topic": "二次方程", "difficulty": 2}'

# 高级难度
curl -X POST http://localhost:48080/api/learning/material/generate \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"subject_id": 2, "grade": "初中", "topic": "二次方程", "difficulty": 3}'
```

### 场景 3: 不同学科测试

```bash
# 语文
curl -X POST http://localhost:48080/api/learning/material/generate \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"subject_id": 1, "grade": "初中", "topic": "古诗词鉴赏", "difficulty": 2}'

# 英语
curl -X POST http://localhost:48080/api/learning/material/generate \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"subject_id": 3, "grade": "初中", "topic": "现在完成时", "difficulty": 2}'

# 物理
curl -X POST http://localhost:48080/api/learning/material/generate \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"subject_id": 4, "grade": "初中", "topic": "牛顿第一定律", "difficulty": 2}'
```

---

## 性能测试

### 并发测试

使用 Apache Bench 进行并发测试:

```bash
# 安装 ab (如果未安装)
# macOS: brew install httpd
# Linux: apt-get install apache2-utils

# 并发 10 个请求,总共 100 个请求
ab -n 100 -c 10 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -p request.json \
  http://localhost:48080/api/learning/material/generate
```

request.json 内容:
```json
{"subject_id": 2, "grade": "初中", "topic": "测试题材", "difficulty": 1}
```

---

## 常见问题

### Q1: 教材生成失败

**可能原因:**
1. DeepSeek API Key 无效或过期
2. Dify 知识库 API Key 无效
3. 网络连接问题
4. API 配额不足

**解决方法:**
1. 检查 `config/config.yaml` 中的 API Key
2. 查看后端日志获取详细错误信息
3. 确认网络可以访问 DeepSeek 和 Dify API

### Q2: 生成的教材格式不正确

**可能原因:**
1. DeepSeek 返回的内容不是标准 JSON
2. Prompt 设计不够明确

**解决方法:**
1. 检查 `learning_material_service.go` 中的 Prompt 模板
2. 优化 Prompt,明确要求输出 JSON 格式
3. 增加 JSON 解析的容错处理

### Q3: 知识库检索无结果

**可能原因:**
1. Dify 知识库中没有相关内容
2. 检索关键词不准确

**解决方法:**
1. 在 Dify 平台上传相关教学资料
2. 优化检索关键词的构建逻辑

---

## 下一步计划

- [ ] 添加音频生成功能(TTS)
- [ ] 实现收藏功能
- [ ] 添加练习题批改功能
- [ ] 实现错题本功能
- [ ] 开发前端页面
