# GORM 数据库自动迁移说明

## 概述

本项目使用 GORM 的 AutoMigrate 功能来管理数据库结构，无需手动创建SQL迁移文件。

## 工作原理

1. **自动检测**: GORM 会自动检测 Go 模型结构的变化
2. **增量迁移**: 只添加新的字段，不会删除现有数据
3. **类型映射**: 自动将 Go 类型映射为数据库类型
4. **索引管理**: 自动创建和管理索引

## User 模型新增字段

在 `internal/model/user.go` 中添加了以下字段：

```go
Address string `gorm:"type:varchar(200)" json:"address"`
Des     string `gorm:"type:varchar(500)" json:"des"`
```

## 自动迁移流程

### 启动应用时

1. `database.InitializeDatabase()` 被调用
2. GORM 检查 `sys_user` 表结构
3. 自动添加 `address` 和 `des` 字段
4. 保持现有数据不变

### 迁移包含的模型

在 `internal/pkg/database/database.go` 中的 `AutoMigrate()` 函数包含：

- 用户管理：User, Role, Menu, Button
- 字典管理：DictionaryType, Dictionary  
- 项目管理：ProjectTemplate, Project, Task 等
- 聊天室：ChatRoom, ChatMessage 等

## 使用方法

### 1. 修改模型

直接在 `internal/model/` 目录下修改对应的 Go 结构体：

```go
type User struct {
    // 现有字段...
    NewField string `gorm:"type:varchar(100);default:''" json:"newField"`
}
```

### 2. 重启应用

重启后端应用，GORM 会自动检测变化并迁移数据库。

### 3. 验证迁移

查看数据库表结构，确认新字段已添加。

## GORM 标签说明

常用标签：

- `type:varchar(100)` - 指定字段类型
- `default:''` - 默认值
- `not null` - 非空约束
- `uniqueIndex` - 唯一索引
- `index` - 普通索引
- `autoCreateTime` - 自动创建时间
- `autoUpdateTime` - 自动更新时间

## 注意事项

1. **数据安全**: GORM 不会删除现有字段或数据
2. **字段重命名**: 需要先添加新字段，迁移数据后再删除旧字段
3. **类型变更**: 某些类型变更可能导致数据丢失
4. **生产环境**: 建议在低峰期进行迁移操作

## 示例

### 添加新字段

```go
// 在 User 模型中添加
Age int `gorm:"type:int;default:0" json:"age"`
```

重启应用后，数据库会自动添加 `age` 字段。

### 修改字段类型

```go
// 将 varchar(50) 改为 varchar(100)
NickName string `gorm:"type:varchar(100)" json:"nickName"`
```

重启应用后，GORM 会自动调整字段类型。

## 优势

1. **简单易用**: 无需编写 SQL 迁移脚本
2. **类型安全**: 基于 Go 类型系统
3. **自动化**: 启动时自动执行迁移
4. **跨数据库**: 支持多种数据库类型
5. **版本控制**: 模型结构即迁移文档

## 最佳实践

1. 在开发环境充分测试迁移
2. 重要操作前备份数据库
3. 使用有意义的字段名和注释
4. 合理设置字段长度和约束
5. 定期检查数据库结构一致性
