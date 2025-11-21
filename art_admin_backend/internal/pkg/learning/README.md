# AI 学习系统模块

## 概述

AI 学习系统是一个独立的功能模块,提供基于 AI 的个性化学习教材生成功能。

## 技术架构

```
Dify 知识库(向量检索) + 自研 RAG + DeepSeek API
```

## 核心功能

### Week 1 已完成功能 ✅

1. **数据库表设计和创建**
   - `subjects` - 学科表
   - `learning_materials` - 教材内容表
   - `material_sections` - 教材章节表
   - `audio_cache` - 音频缓存表

2. **DeepSeek 客户端封装**
   - 位置: `internal/pkg/deepseek/client.go`
   - 功能: 封装 DeepSeek API 调用,支持对话生成

3. **Dify 知识库检索集成**
   - 位置: `internal/pkg/dify/knowledge.go`
   - 功能: 从 Dify 知识库检索相关教学资料

4. **教材生成核心逻辑**
   - 位置: `internal/service/learning_material_service.go`
   - 流程:
     1. 从 Dify 知识库检索相关资料
     2. 构建 RAG Prompt
     3. 调用 DeepSeek API 生成教材
     4. 解析并保存教材内容

## 目录结构

```
internal/
├── model/
│   ├── learning_subject.go       # 学科模型
│   ├── learning_material.go      # 教材模型
│   └── audio_cache.go            # 音频缓存模型
├── repository/
│   ├── subject_repository.go     # 学科仓储
│   └── learning_material_repository.go  # 教材仓储
├── service/
│   ├── subject_service.go        # 学科服务
│   └── learning_material_service.go     # 教材服务
├── api/
│   └── learning_api.go           # 学习系统 API
├── dto/
│   └── learning_dto.go           # 数据传输对象
├── router/
│   └── learning_router.go        # 路由配置
└── pkg/
    ├── deepseek/
    │   └── client.go             # DeepSeek 客户端
    ├── dify/
    │   └── knowledge.go          # Dify 知识库客户端
    └── learning/
        ├── init.go               # 依赖注入容器
        └── README.md             # 本文档
```

## API 接口

### 学科管理

- `GET /api/learning/subject/list` - 获取学科列表(无需认证)

### 教材管理(需要 JWT 认证)

- `POST /api/learning/material/generate` - 生成教材
- `GET /api/learning/material/list` - 获取教材列表
- `GET /api/learning/material/:id` - 获取教材详情
- `DELETE /api/learning/material/:id` - 删除教材

## 配置说明

在 `config/config.yaml` 中添加以下配置:

```yaml
# DeepSeek API 配置
deepseek:
  apiKey: sk-a8e4ee88516f40e6a2dc3776d3254846
  baseUrl: https://api.deepseek.com
  model: deepseek-chat
  timeout: 60
  maxTokens: 4096
  temperature: 0.7

# Dify 知识库配置
dify:
  datasetApiKey: dataset-QQ8Y3Cxvss4RvQXvU7Gem4sb
  chatApiKey: app-D7hSHUnPoMt5CQD2IEoBsPJz
  baseUrl: https://api.dify.ai/v1
  timeout: 30
```

## 数据库初始化

执行以下 SQL 脚本创建数据库表:

```bash
mysql -u root -p gin_admin < scripts/migrations/learning_system.sql
```

或者使用 GORM 自动迁移(已在主程序中集成)。

## 使用示例

### 1. 生成教材

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

### 2. 获取教材列表

```bash
curl -X GET "http://localhost:48080/api/learning/material/list?page=1&limit=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 3. 获取教材详情

```bash
curl -X GET http://localhost:48080/api/learning/material/123 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 开发计划

### Phase 1: 基础功能(Week 1-2) ✅

- [x] 数据库表设计和创建
- [x] DeepSeek 客户端封装
- [x] Dify 知识库检索集成
- [x] 教材生成核心逻辑
- [ ] 教材生成 API 接口测试
- [ ] 前端教材生成页面
- [ ] 前端教材详情页面
- [ ] 基础音频播放功能

### Phase 2: 增强功能(Week 3-4)

- [ ] TTS 服务集成
- [ ] 音频缓存机制
- [ ] 收藏功能完整实现
- [ ] 收藏夹管理
- [ ] 练习题生成和展示
- [ ] AI 批改功能
- [ ] 错题本功能
- [ ] 笔记功能

### Phase 3: 优化和完善(Week 5)

- [ ] 性能优化
- [ ] 缓存策略优化
- [ ] UI/UX 优化
- [ ] 测试和 Bug 修复

## 注意事项

1. **API Key 安全**
   - 不要将 API Key 提交到版本控制
   - 使用环境变量或配置文件
   - 定期轮换 API Key

2. **内容质量控制**
   - 定期审核 AI 生成的内容
   - 建立内容评分机制
   - 收集用户反馈优化 Prompt

3. **性能优化**
   - 实现请求限流
   - 添加缓存机制
   - 异步处理耗时任务

## 相关文档

- [AI 学习系统设计文档](../../../AI_LEARNING_SYSTEM_DESIGN.md)
- [DeepSeek API 文档](https://platform.deepseek.com/api-docs/)
- [Dify 官方文档](https://docs.dify.ai/)
