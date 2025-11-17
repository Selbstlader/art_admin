# Art Admin Backend

基于 Gin + GORM + MySQL 的后端管理系统

## ✨ 特性

- ✅ **完全自动化** - GORM 自动管理数据库，无需手动操作
- ✅ **自动创建数据库** - 首次运行自动创建数据库
- ✅ **自动迁移表结构** - 根据模型自动创建/更新表
- ✅ **自动初始化数据** - 自动生成测试账号和基础数据
- ✅ **JWT 认证** - 双 Token 机制（AccessToken + RefreshToken）
- ✅ **Swagger 文档** - 自动生成 API 文档
- ✅ **RBAC 权限** - 角色-菜单-按钮三级权限控制

## 🚀 快速启动

### 1. 启动后端

```bash
cd art_admin_backend
./run.sh
```

### 2. 访问 Swagger

打开浏览器访问：http://localhost:48080/swagger/index.html

### 3. 测试账号

- **用户名**: admin
- **密码**: 123456

## 📁 项目结构

```
art_admin_backend/
├── cmd/server/          # 程序入口
├── config/              # 配置文件
├── internal/
│   ├── api/            # 控制器层
│   ├── dto/            # 数据传输对象
│   ├── model/          # GORM 模型
│   ├── repository/     # 数据访问层
│   ├── service/        # 业务逻辑层
│   └── pkg/
│       ├── config/     # 配置管理
│       ├── database/   # 数据库管理（自动化）
│       ├── jwt/        # JWT 工具
│       ├── logger/     # 日志工具
│       └── utils/      # 工具函数
└── docs/               # Swagger 文档（自动生成）
```

## 🔧 GORM 自动化管理

项目使用 GORM 完全自动化管理数据库：

1. **自动创建数据库** - 检查并创建 `gin_admin` 数据库
2. **自动迁移表结构** - 根据模型自动创建 7 张表
3. **自动初始化数据** - 自动生成角色、用户、菜单、按钮权限

### 数据库表

- `sys_user` - 用户表
- `sys_role` - 角色表
- `sys_menu` - 菜单表
- `sys_button` - 按钮权限表
- `sys_user_role` - 用户角色关联表
- `sys_role_menu` - 角色菜单关联表
- `sys_role_button` - 角色按钮关联表

## 📝 API 接口

| 接口 | 方法 | 说明 | 认证 |
|-----|------|------|------|
| `/auth/login` | POST | 用户登录 | ❌ |
| `/user/info` | GET | 获取用户信息 | ✅ |
| `/api/system/menus` | GET | 获取菜单树 | ✅ |
| `/api/user/list` | GET | 用户列表 | ✅ |
| `/api/role/list` | GET | 角色列表 | ✅ |
| `/health` | GET | 健康检查 | ❌ |

## ⚙️ 配置说明

配置文件：`config/config.yaml`

```yaml
server:
  port: 48080        # 服务端口
  mode: debug        # 运行模式：debug/release

database:
  host: localhost    # 数据库主机
  port: 3306        # 数据库端口
  database: gin_admin  # 数据库名（自动创建）
  username: root    # 数据库用户名
  password: 123456  # 数据库密码
```

## 🎯 核心功能

### 1. GORM 自动化
- 自动创建数据库
- 自动迁移表结构  
- 自动初始化数据
- 代码优先，无需 SQL

### 2. JWT 认证
- AccessToken（2小时有效）
- RefreshToken（7天有效）
- 自动刷新机制

### 3. RBAC 权限
- 用户-角色多对多
- 角色-菜单多对多
- 角色-按钮多对多
- 动态菜单树

## 📦 技术栈

- **Web 框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL 8.0+
- **日志**: Zap + Lumberjack
- **配置**: Viper
- **文档**: Swagger
- **认证**: JWT
- **密码**: Bcrypt

## 🔄 开发流程

### 添加新模型

1. 在 `internal/model/` 创建模型
2. 启动项目，GORM 自动创建表
3. 无需手动写 SQL

### 添加新接口

1. 定义 DTO（`internal/dto/`）
2. 添加 Repository（`internal/repository/`）
3. 添加 Service（`internal/service/`）
4. 添加 Controller（`internal/api/v1/`）
5. 注册路由（`internal/router/`）
6. 添加 Swagger 注释
7. 启动项目，自动生成文档

## 📌 注意事项

1. 首次运行会自动创建数据库和表
2. 数据已存在时不会重复创建
3. 修改模型后重启会自动更新表结构
4. GORM 只会添加字段，不会删除字段
5. 生产环境建议使用 SQL 脚本管理表结构

## 🎉 完成

现在后端已完全自动化，无需任何手动数据库操作！
