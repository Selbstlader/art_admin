# Dify 功能测试总结

## 测试时间
2025-11-20 17:29

## 配置信息

### API Keys
- **知识库 API Key**: `dataset-QQ8Y3Cxvss4RvQXvU7Gem4sb`
- **AI对话 API Key**: `app-D7hSHUnPoMt5CQD2IEoBsPJz`
- **Base URL**: `https://api.dify.ai/v1`

## 测试结果

### ✅ 1. 知识库列表接口 - 成功

**接口**: `GET /api/dify/dataset/list`

**测试命令**:
```bash
curl -X GET 'http://localhost:48080/api/dify/dataset/list?page=1&limit=10' \
  -H 'Authorization: Bearer {jwt-token}'
```

**结果**: ✅ 成功返回 3 个知识库
- 测试知识库 1 (711b1a11-2a29-4920-a379-6ee6251a97fa)
- 酒店知识库 (9ea335ce-11cb-4cb5-8ad9-af30a735ec2a)
- 测试知识库 (7e2a67c6-00ab-46d0-a71b-886bb480d18f)

### ✅ 2. 知识库详情接口 - 成功

**接口**: `GET /api/dify/dataset/{id}`

**测试命令**:
```bash
curl -X GET 'http://localhost:48080/api/dify/dataset/711b1a11-2a29-4920-a379-6ee6251a97fa' \
  -H 'Authorization: Bearer {jwt-token}'
```

**结果**: ✅ 成功返回知识库详细信息

### ⚠️ 3. AI 对话接口 - 需要配置输入参数

**接口**: `POST /api/dify/chat`

**问题**: Dify 应用配置了必填输入变量 `birth_date`

**解决方案**: 在请求中添加 inputs 参数

**正确的请求示例**:
```bash
curl -X POST 'http://localhost:48080/api/dify/chat' \
  -H 'Authorization: Bearer {jwt-token}' \
  -H 'Content-Type: application/json' \
  -d '{
    "query": "你好",
    "user": "user-123",
    "inputs": {
      "birth_date": "1990-01-01"
    }
  }'
```

### 📝 4. 文件上传接口 - 待测试

**接口**: `POST /api/dify/dataset/upload`

**测试命令**:
```bash
curl -X POST 'http://localhost:48080/api/dify/dataset/upload' \
  -H 'Authorization: Bearer {jwt-token}' \
  -F 'dataset_id=711b1a11-2a29-4920-a379-6ee6251a97fa' \
  -F 'file=@test.txt'
```

### 📝 5. 流式对话接口 - 待测试

**接口**: `POST /api/dify/chat/stream`

**测试命令**:
```bash
curl -X POST 'http://localhost:48080/api/dify/chat/stream' \
  -H 'Authorization: Bearer {jwt-token}' \
  -H 'Content-Type: application/json' \
  -d '{
    "query": "介绍一下人工智能",
    "user": "user-123",
    "inputs": {
      "birth_date": "1990-01-01"
    }
  }'
```

## 已实现的功能

### 后端实现
1. ✅ Dify 客户端封装 (`internal/pkg/dify/client.go`)
2. ✅ 业务逻辑层 (`internal/service/dify_service.go`)
3. ✅ API 控制器 (`internal/api/v1/dify.go`)
4. ✅ 路由配置 (`internal/router/router.go`)
5. ✅ DTO 定义 (`internal/dto/request/dify_request.go`, `internal/dto/response/dify_response.go`)
6. ✅ 配置文件支持 (`config/config.yaml`)

### API 接口
1. ✅ `GET /api/dify/dataset/list` - 获取知识库列表
2. ✅ `GET /api/dify/dataset/:id` - 获取知识库详情
3. ✅ `POST /api/dify/dataset/upload` - 上传文件到知识库
4. ✅ `POST /api/dify/chat` - AI对话(非流式)
5. ✅ `POST /api/dify/chat/stream` - AI对话(流式)

## 注意事项

### 1. 应用输入变量
不同的 Dify 应用可能配置了不同的必填输入变量。使用前需要:
- 在 Dify 控制台查看应用的输入变量配置
- 在请求的 `inputs` 字段中提供所有必填变量

### 2. API Key 管理
- 知识库操作使用 `datasetApiKey`
- AI对话操作使用 `chatApiKey`
- 两个 Key 分别对应不同的功能模块

### 3. 流式响应
- 流式对话返回 SSE (Server-Sent Events) 格式
- 客户端需要支持流式数据接收
- 适合长文本生成场景

### 4. 文件上传
- 支持多种文档格式 (txt, md, pdf, docx, xlsx, csv 等)
- 文件会临时保存在 `./temp/uploads` 目录
- 上传后会自动删除临时文件

## 下一步工作

1. ⏳ 重启服务器使代码生效
2. ⏳ 测试 AI 对话接口(添加正确的 inputs 参数)
3. ⏳ 测试文件上传功能
4. ⏳ 测试流式对话功能
5. ⏳ 前端集成

## 相关文档

- [Dify 官方文档](https://docs.dify.ai/)
- [后端集成文档](./DIFY_INTEGRATION.md)
- [API 测试脚本](./test_dify.sh)
