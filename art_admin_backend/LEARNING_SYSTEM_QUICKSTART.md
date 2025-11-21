# AI 学习系统 - 快速开始指南

## 📋 概述

AI 学习系统是一个基于 DeepSeek + Dify 知识库的智能教材生成系统,支持多学科、多年级、多难度的个性化学习内容生成。

## ✅ Week 1 已完成功能

- ✅ 数据库表设计和创建
- ✅ DeepSeek 客户端封装
- ✅ Dify 知识库检索集成
- ✅ 教材生成核心逻辑
- ✅ RESTful API 接口
- ✅ 完整的文档和测试用例

## 🚀 快速开始

### 1. 配置 API Keys

编辑 `config/config.yaml`:

```yaml
# DeepSeek API 配置
deepseek:
  apiKey: sk-a8e4ee88516f40e6a2dc3776d3254846  # 已配置
  baseUrl: https://api.deepseek.com
  model: deepseek-chat
  timeout: 60
  maxTokens: 4096
  temperature: 0.7

# Dify 知识库配置
dify:
  datasetApiKey: dataset-QQ8Y3Cxvss4RvQXvU7Gem4sb  # 已配置
  baseUrl: https://api.dify.ai/v1
  timeout: 30
  datasetId: your-dataset-id-here  # ⚠️ 需要配置你的 Dataset ID
```

**获取 Dify Dataset ID:**
1. 登录 Dify 平台: https://cloud.dify.ai
2. 进入知识库管理
3. 选择或创建一个知识库
4. 在 API 设置中找到 Dataset ID

### 2. 初始化数据库

```bash
cd /Users/xinyoucai/code/art_admin/art_admin_backend

# 方式1: 使用初始化脚本(推荐)
./scripts/init_learning_system.sh

# 方式2: 手动执行 SQL
mysql -u root -p gin_admin < scripts/migrations/learning_system.sql
```

**创建的表:**
- `subjects` - 学科表(已预置 9 个学科)
- `learning_materials` - 教材内容表
- `material_sections` - 教材章节表
- `audio_cache` - 音频缓存表

### 3. 启动服务

```bash
# 启动后端服务
go run cmd/server/main.go

# 服务启动后访问:
# - API 文档: http://localhost:48080/swagger/index.html
# - 健康检查: http://localhost:48080/health
```

### 4. 测试 API

#### 4.1 获取学科列表(无需认证)

```bash
curl http://localhost:48080/api/learning/subject/list
```

#### 4.2 登录获取 Token

```bash
curl -X POST http://localhost:48080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "123456"}'
```

保存返回的 `access_token`。

#### 4.3 生成教材

```bash
curl -X POST http://localhost:48080/api/learning/material/generate \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "subject_id": 2,
    "grade": "初中",
    "topic": "二次方程的解法",
    "difficulty": 2
  }'
```

**参数说明:**
- `subject_id`: 学科 ID
  - 1: 语文, 2: 数学, 3: 英语
  - 4: 物理, 5: 化学, 6: 生物
  - 7: 历史, 8: 地理, 9: 政治
- `grade`: 年级(小学/初中/高中)
- `topic`: 学习题材
- `difficulty`: 难度(1-基础, 2-进阶, 3-高级)

#### 4.4 查看教材列表

```bash
curl -X GET "http://localhost:48080/api/learning/material/list?page=1&limit=10" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## 📁 项目结构

```
art_admin_backend/
├── cmd/server/main.go                          # 主程序(已集成学习系统)
├── config/config.yaml                          # 配置文件
├── scripts/
│   ├── migrations/learning_system.sql          # 数据库迁移脚本
│   └── init_learning_system.sh                 # 初始化脚本
├── docs/
│   └── LEARNING_SYSTEM_API_TEST.md            # API 测试文档
└── internal/
    ├── model/
    │   ├── learning_subject.go                 # 学科模型
    │   ├── learning_material.go                # 教材模型
    │   └── audio_cache.go                      # 音频缓存模型
    ├── repository/
    │   ├── subject_repository.go               # 学科仓储
    │   └── learning_material_repository.go     # 教材仓储
    ├── service/
    │   ├── subject_service.go                  # 学科服务
    │   └── learning_material_service.go        # 教材服务(核心逻辑)
    ├── api/
    │   └── learning_api.go                     # API 控制器
    ├── dto/
    │   └── learning_dto.go                     # 数据传输对象
    ├── router/
    │   └── learning_router.go                  # 路由配置
    └── pkg/
        ├── deepseek/
        │   └── client.go                       # DeepSeek 客户端
        ├── dify/
        │   └── knowledge.go                    # Dify 知识库客户端
        └── learning/
            ├── init.go                         # 依赖注入
            └── README.md                       # 模块文档
```

## 🔧 核心功能说明

### 1. 教材生成流程

```
用户请求
  ↓
① Dify 知识库检索相关资料
  ↓
② 构建 RAG Prompt(系统提示词 + 检索结果 + 用户输入)
  ↓
③ DeepSeek API 生成结构化教材
  ↓
④ 解析 JSON 并保存到数据库
  ↓
返回教材内容
```

### 2. 教材内容结构

生成的教材包含:
- **知识点讲解** - 清晰易懂的概念说明
- **重点难点** - 标注关键知识点
- **例题演示** - 至少 2 个示例
- **练习题** - 至少 5 个练习题(含答案和解析)

### 3. API 接口列表

| 接口 | 方法 | 路径 | 认证 | 说明 |
|------|------|------|------|------|
| 获取学科列表 | GET | `/api/learning/subject/list` | ❌ | 获取所有学科 |
| 生成教材 | POST | `/api/learning/material/generate` | ✅ | 生成新教材 |
| 教材列表 | GET | `/api/learning/material/list` | ✅ | 分页查询教材 |
| 教材详情 | GET | `/api/learning/material/:id` | ✅ | 获取教材详情 |
| 删除教材 | DELETE | `/api/learning/material/:id` | ✅ | 删除教材 |

## 📚 相关文档

- [完整设计文档](AI_LEARNING_SYSTEM_DESIGN.md)
- [API 测试文档](docs/LEARNING_SYSTEM_API_TEST.md)
- [模块说明](internal/pkg/learning/README.md)

## 🐛 常见问题

### Q1: 教材生成失败

**检查清单:**
- [ ] DeepSeek API Key 是否正确
- [ ] Dify Dataset API Key 是否正确
- [ ] Dify Dataset ID 是否配置
- [ ] 网络是否可以访问外部 API
- [ ] 查看后端日志获取详细错误

### Q2: 知识库检索无结果

**解决方法:**
1. 登录 Dify 平台检查知识库是否有内容
2. 上传相关教学资料到知识库
3. 优化检索关键词

### Q3: 生成的内容格式不正确

**解决方法:**
1. 检查 DeepSeek 返回的原始内容
2. 优化 `learning_material_service.go` 中的 Prompt
3. 增强 JSON 解析的容错处理

## 🎯 下一步计划

### Phase 2: 增强功能(Week 3-4)

- [ ] 集成阿里云 TTS 服务
- [ ] 实现音频缓存机制
- [ ] 添加收藏功能
- [ ] 实现收藏夹管理
- [ ] 练习题生成和展示
- [ ] AI 批改功能
- [ ] 错题本功能
- [ ] 笔记功能

### Phase 3: 优化和完善(Week 5)

- [ ] 性能优化
- [ ] 缓存策略优化
- [ ] UI/UX 优化
- [ ] 测试和 Bug 修复

## 💡 使用建议

1. **首次使用:**
   - 先测试获取学科列表
   - 选择一个简单的题材测试生成
   - 查看生成的教材结构

2. **优化 Prompt:**
   - 根据生成效果调整 `buildSystemPrompt` 函数
   - 可以针对不同学科定制不同的 Prompt

3. **知识库管理:**
   - 定期更新 Dify 知识库内容
   - 按学科分类组织教学资料
   - 确保资料的准确性和时效性

## 📞 技术支持

如有问题,请查看:
1. 后端日志: `logs/app.log`
2. API 文档: http://localhost:48080/swagger/index.html
3. 设计文档: `AI_LEARNING_SYSTEM_DESIGN.md`

---

**开发完成时间:** 2024-11-21  
**版本:** v1.0.0 (Week 1)  
**状态:** ✅ 核心功能已完成
