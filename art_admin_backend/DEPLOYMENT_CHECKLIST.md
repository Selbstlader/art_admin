# AI 学习系统部署检查清单

## 📋 部署前检查

### 1. 环境准备

- [ ] Go 1.21+ 已安装
- [ ] MySQL 8.0+ 已安装并运行
- [ ] Git 已安装
- [ ] 网络可以访问外部 API (DeepSeek, Dify)

### 2. 配置文件检查

#### config/config.yaml

- [ ] 数据库配置正确
  ```yaml
  database:
    host: 127.0.0.1
    port: 3306
    database: gin_admin
    username: root
    password: 123456
  ```

- [ ] DeepSeek API Key 已配置
  ```yaml
  deepseek:
    apiKey: sk-a8e4ee88516f40e6a2dc3776d3254846  # ✅ 已配置
  ```

- [ ] Dify API Key 已配置
  ```yaml
  dify:
    datasetApiKey: dataset-QQ8Y3Cxvss4RvQXvU7Gem4sb  # ✅ 已配置
    datasetId: your-dataset-id-here  # ⚠️ 需要配置
  ```

### 3. 数据库初始化

- [ ] 执行数据库迁移脚本
  ```bash
  ./scripts/init_learning_system.sh
  ```

- [ ] 验证表已创建
  ```sql
  SHOW TABLES LIKE '%learning%';
  SHOW TABLES LIKE 'subjects';
  SHOW TABLES LIKE 'audio_cache';
  ```

- [ ] 验证默认数据已插入
  ```sql
  SELECT COUNT(*) FROM subjects;  -- 应该返回 9
  ```

### 4. 依赖检查

- [ ] 安装 Go 依赖
  ```bash
  go mod tidy
  go mod download
  ```

- [ ] 验证编译无错误
  ```bash
  go build -o bin/server cmd/server/main.go
  ```

## 🚀 启动服务

### 1. 启动后端

```bash
# 开发模式
go run cmd/server/main.go

# 生产模式
go build -o bin/server cmd/server/main.go
./bin/server
```

### 2. 验证服务启动

- [ ] 健康检查
  ```bash
  curl http://localhost:48080/health
  ```
  预期响应: `{"status":"ok","message":"Art Admin Backend is running"}`

- [ ] Swagger 文档可访问
  访问: http://localhost:48080/swagger/index.html

## ✅ 功能测试

### 1. 基础功能测试

- [ ] 获取学科列表
  ```bash
  curl http://localhost:48080/api/learning/subject/list
  ```

- [ ] 用户登录
  ```bash
  curl -X POST http://localhost:48080/api/auth/login \
    -H "Content-Type: application/json" \
    -d '{"username": "admin", "password": "123456"}'
  ```

### 2. 核心功能测试

- [ ] 生成教材(数学)
  ```bash
  curl -X POST http://localhost:48080/api/learning/material/generate \
    -H "Authorization: Bearer YOUR_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
      "subject_id": 2,
      "grade": "初中",
      "topic": "二次方程的解法",
      "difficulty": 2
    }'
  ```

- [ ] 获取教材列表
  ```bash
  curl -X GET "http://localhost:48080/api/learning/material/list?page=1&limit=10" \
    -H "Authorization: Bearer YOUR_TOKEN"
  ```

- [ ] 获取教材详情
  ```bash
  curl -X GET http://localhost:48080/api/learning/material/1 \
    -H "Authorization: Bearer YOUR_TOKEN"
  ```

- [ ] 删除教材
  ```bash
  curl -X DELETE http://localhost:48080/api/learning/material/1 \
    -H "Authorization: Bearer YOUR_TOKEN"
  ```

### 3. 异常情况测试

- [ ] 未登录访问受保护接口
  ```bash
  curl -X POST http://localhost:48080/api/learning/material/generate \
    -H "Content-Type: application/json" \
    -d '{"subject_id": 2, "grade": "初中", "topic": "测试", "difficulty": 1}'
  ```
  预期: 返回 401 未登录

- [ ] 参数错误
  ```bash
  curl -X POST http://localhost:48080/api/learning/material/generate \
    -H "Authorization: Bearer YOUR_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"subject_id": 2}'
  ```
  预期: 返回 400 参数错误

- [ ] 访问不存在的教材
  ```bash
  curl -X GET http://localhost:48080/api/learning/material/99999 \
    -H "Authorization: Bearer YOUR_TOKEN"
  ```
  预期: 返回 404 教材不存在

## 🔍 日志检查

### 1. 查看日志文件

```bash
tail -f logs/app.log
```

### 2. 关键日志检查

- [ ] 服务启动日志
  ```
  服务器启动成功，监听端口: 48080
  ```

- [ ] 数据库连接日志
  ```
  数据库连接成功
  ```

- [ ] API 请求日志
  ```
  POST /api/learning/material/generate
  ```

## 🐛 常见问题排查

### 问题 1: 数据库连接失败

**检查:**
- [ ] MySQL 服务是否运行
- [ ] 数据库配置是否正确
- [ ] 数据库用户权限是否足够

**解决:**
```bash
# 检查 MySQL 状态
mysql -u root -p -e "SELECT 1"

# 检查数据库是否存在
mysql -u root -p -e "SHOW DATABASES LIKE 'gin_admin'"
```

### 问题 2: API Key 无效

**检查:**
- [ ] DeepSeek API Key 是否正确
- [ ] Dify Dataset API Key 是否正确
- [ ] API Key 是否过期

**解决:**
1. 登录 DeepSeek 平台检查 API Key
2. 登录 Dify 平台检查 API Key
3. 重新生成并更新配置文件

### 问题 3: 教材生成失败

**检查:**
- [ ] 查看后端日志获取详细错误
- [ ] 检查 Dify Dataset ID 是否配置
- [ ] 检查知识库是否有内容
- [ ] 检查网络连接

**解决:**
```bash
# 查看详细日志
tail -100 logs/app.log | grep -i error

# 测试 DeepSeek API
curl -X POST https://api.deepseek.com/chat/completions \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"model":"deepseek-chat","messages":[{"role":"user","content":"test"}]}'

# 测试 Dify API
curl -X POST https://api.dify.ai/v1/datasets/YOUR_DATASET_ID/retrieve \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"query":"test","retrieval_model":{"search_method":"semantic_search","top_k":3}}'
```

## 📊 性能监控

### 1. 系统资源

- [ ] CPU 使用率 < 80%
- [ ] 内存使用率 < 80%
- [ ] 磁盘空间充足

### 2. API 响应时间

- [ ] 学科列表 < 100ms
- [ ] 教材列表 < 500ms
- [ ] 教材生成 < 30s (取决于 AI 响应时间)

### 3. 数据库性能

```sql
-- 检查慢查询
SHOW VARIABLES LIKE 'slow_query_log';

-- 检查连接数
SHOW STATUS LIKE 'Threads_connected';

-- 检查表大小
SELECT 
  table_name,
  ROUND(((data_length + index_length) / 1024 / 1024), 2) AS size_mb
FROM information_schema.TABLES
WHERE table_schema = 'gin_admin'
  AND table_name LIKE '%learning%';
```

## 🔒 安全检查

- [ ] API Key 不在版本控制中
- [ ] 数据库密码足够复杂
- [ ] JWT Secret 已修改为生产环境密钥
- [ ] CORS 配置正确
- [ ] 日志不包含敏感信息

## 📝 部署后操作

### 1. 备份

- [ ] 备份数据库
  ```bash
  mysqldump -u root -p gin_admin > backup_$(date +%Y%m%d).sql
  ```

- [ ] 备份配置文件
  ```bash
  cp config/config.yaml config/config.yaml.backup
  ```

### 2. 监控设置

- [ ] 设置日志轮转
- [ ] 配置告警通知
- [ ] 设置性能监控

### 3. 文档更新

- [ ] 更新部署文档
- [ ] 记录配置变更
- [ ] 更新 API 文档

## ✅ 部署完成确认

- [ ] 所有测试用例通过
- [ ] 日志无错误信息
- [ ] 性能指标正常
- [ ] 安全检查通过
- [ ] 文档已更新

---

**检查人:** ___________  
**检查日期:** ___________  
**部署环境:** [ ] 开发 [ ] 测试 [ ] 生产  
**版本号:** v1.0.0

## 🎉 部署成功!

系统已成功部署并可以使用。

**下一步:**
1. 通知团队成员
2. 开始前端开发
3. 收集用户反馈
4. 规划 Phase 2 功能

**相关文档:**
- [快速开始指南](LEARNING_SYSTEM_QUICKSTART.md)
- [API 测试文档](docs/LEARNING_SYSTEM_API_TEST.md)
- [开发总结](LEARNING_SYSTEM_SUMMARY.md)
