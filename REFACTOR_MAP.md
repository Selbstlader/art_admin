# 模块化重构映射文档

## 项目概述
对art_admin项目进行模块化重构，将原本集中管理的文件按业务模块重新组织。

## 重构原则
1. 后端避免循环导入：model/dto/repository保持共享，service/api按模块重组
2. 前端按模块重组components/utils/types
3. 逐模块迁移，确保每步可编译
4. 保持现有功能不变

## 发现的耦合问题（重要）
### mood模块重构失败原因分析
1. **深度耦合**：mood模块与AI分析功能深度耦合，难以分离
2. **私有字段访问**：mood服务直接访问AIAnalysisService的私有字段deepSeekClient
3. **类型定义重复**：EmotionAnalysisRequest、EmotionPattern等类型在service层和dto层都有定义
4. **共享工具函数**：GetAnalysisDebouncer函数被多个模块共享使用
5. **跨模块依赖**：AIAnalysisService主要服务于mood业务但被其他模块引用

### 解决方案
- **mood+ai_analysis应合并为mood_analytics模块**：由于紧密耦合，应视为一个业务域
- **选择独立模块试点**：优先选择learning或chat等耦合度低的模块验证重构方法

## 成功的重构模式（learning模块）
### 重构步骤（可复用）
1. **依赖分析**：确认模块无内部service层耦合
2. **创建模块目录**：`mkdir -p internal/service/{module}/`
3. **移动service文件**：`mv service/*_service.go service/{module}/`
4. **修改package声明**：将`package service`改为`package {module}`
5. **更新import路径**：使用别名避免命名冲突（如`learningSvc`）
6. **逐文件更新引用**：API层、pkg层等
7. **编译测试验证**：确保每步都可编译

### 成功关键因素
- **低耦合模块**：learning模块只依赖repository层和外部客户端
- **逐步验证**：每步修改后立即测试编译
- **别名策略**：使用`{module}Svc`别名避免包名冲突
- **原子提交**：每个模块重构作为独立提交

## 后端重构映射

### Service层重组
- mood_record_service.go → service/mood/
- meditation_record_service.go → service/mood/
- meditation_content_service.go → service/mood/
- meditation_favorite_service.go → service/mood/
- journal_entry_service.go → service/mood/
- user_goal_service.go → service/mood/ (跨模块共享)

- learning_material_service.go → service/learning/
- subject_service.go → service/learning/

- project_service.go → service/project/
- project_template_service.go → service/project/
- task_service.go → service/project/
- achievement_service.go → service/project/

- chat_service.go → service/chat/

- user_service.go → service/system/
- role_service.go → service/system/
- department_service.go → service/system/
- menu_service.go → service/system/
- dictionary_service.go → service/system/
- operation_log_service.go → service/system/
- app_user_service.go → service/system/

- auth_service.go → service/auth/

### API层重组
- chat_api.go → api/chat/
- learning_api.go → api/learning/
- api/project/ → 保持现状（已模块化）
- 新建 api/mood/, api/system/, api/auth/

## 前端重构映射

### Components重组
- views/mood/下的组件 → components/mood/
- views/project/下的组件 → components/project/
- views/learning/下的组件 → components/learning/
- views/chat/下的组件 → components/chat/
- views/system/下的组件 → components/system/
- 通用组件 → components/common/

### Utils重组
- 按业务模块拆分utils/
- 创建模块专用utils文件夹
- 通用工具函数保留在utils/common/

### Types重组
- 按模块拆分types/
- 创建模块专用类型定义
- 通用类型保留在types/common/

## 执行顺序
1. 后端mood模块试点
2. 后端其他模块逐个迁移
3. 前端components重组
4. 前端utils重组
5. 前端types重组
6. 更新所有import路径
7. 编译测试验证

## 风险控制
- 每个模块迁移后立即编译测试
- 保持git提交的原子性
- 如遇问题及时回滚
