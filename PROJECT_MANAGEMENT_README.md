# 项目管理功能使用说明

## 功能概述

项目管理模块已完成开发,包含以下核心功能:

### 1. 后端功能 ✅
- ✅ 项目模板管理 - 创建可复用的项目流程模板
- ✅ 项目管理 - 从模板创建项目或手动创建
- ✅ 任务管理 - 支持任务分解、进度跟踪、超期标记
- ✅ 任务依赖 - 支持FS/SS/FF/SF四种依赖类型
- ✅ 甘特图数据 - 完全符合DHTMLX Gantt格式
- ✅ 批量更新 - 支持甘特图拖拽后批量保存

### 2. 前端功能 ✅
- ✅ 项目列表页面 - 查询、筛选、管理项目
- ✅ 甘特图页面 - 使用DHTMLX Gantt展示和编辑项目任务
- ✅ API服务层 - 完整的接口调用封装
- ✅ 路由配置 - 已集成到系统路由

## 快速开始

### 1. 启动后端服务

```bash
cd art_admin_backend
./run.sh
```

后端服务将在 `http://localhost:48080` 启动

### 2. 启动前端服务

```bash
cd art-design-pro
pnpm dev
```

前端服务将在 `http://localhost:5173` 启动

### 3. 访问项目管理

登录系统后,在左侧菜单中找到"项目管理"菜单,点击进入。

## 功能使用

### 创建项目模板

1. 进入项目管理 → 项目列表
2. 点击"从模板创建"按钮
3. 如果没有模板,需要先通过API创建模板

**创建模板示例:**

```bash
curl -X POST 'http://localhost:48080/api/project/template' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer YOUR_TOKEN' \
  -d '{
    "name": "软件开发项目模板",
    "description": "标准的软件开发项目流程模板",
    "category": "软件开发",
    "icon": "code",
    "duration": 90,
    "isPublic": true,
    "tasks": [
      {
        "parentId": 0,
        "name": "需求分析",
        "description": "收集和分析项目需求",
        "startDay": 0,
        "duration": 10,
        "isMilestone": true,
        "priority": 1,
        "sort": 1
      }
    ]
  }'
```

### 从模板创建项目

```bash
curl -X POST 'http://localhost:48080/api/project/from-template' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer YOUR_TOKEN' \
  -d '{
    "templateId": 1,
    "name": "电商平台开发项目",
    "code": "PROJ-2025-001",
    "description": "基于模板创建的电商平台项目",
    "startDate": "2025-11-25T00:00:00Z",
    "managerId": 1,
    "budget": 500000,
    "members": [
      {
        "userId": 1,
        "role": "项目经理",
        "joinDate": "2025-11-25T00:00:00Z"
      }
    ]
  }'
```

### 使用甘特图

1. 在项目列表中,点击某个项目的"甘特图"按钮
2. 进入甘特图页面后,可以:
   - **拖拽任务** - 调整任务的开始时间和持续时间
   - **拖拽进度条** - 调整任务完成进度
   - **创建依赖** - 从一个任务拖拽到另一个任务创建依赖关系
   - **添加任务** - 点击"添加任务"按钮
   - **保存更改** - 点击"保存"按钮批量保存所有更改

### 手动创建任务

```bash
curl -X POST 'http://localhost:48080/api/project/task' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer YOUR_TOKEN' \
  -d '{
    "projectId": 1,
    "parentId": 0,
    "name": "前端开发",
    "description": "开发前端页面",
    "type": "task",
    "priority": 2,
    "startDate": "2025-12-20T00:00:00Z",
    "endDate": "2026-01-10T00:00:00Z",
    "assigneeId": 1,
    "estimatedHours": 160,
    "sort": 1
  }'
```

## 数据库表结构

项目管理模块使用以下9个数据表:

| 表名 | 说明 |
|------|------|
| pm_project_template | 项目模板表 |
| pm_template_task | 模板任务表 |
| pm_template_task_dependency | 模板任务依赖表 |
| pm_project | 项目表 |
| pm_project_member | 项目成员表 |
| pm_task | 任务表 |
| pm_task_dependency | 任务依赖表 |
| pm_task_comment | 任务评论表 |
| pm_task_attachment | 任务附件表 |

## API 接口列表

### 项目模板接口
- `GET /api/project/template/list` - 获取模板列表
- `GET /api/project/template/:id` - 获取模板详情
- `POST /api/project/template` - 创建模板
- `PUT /api/project/template` - 更新模板
- `DELETE /api/project/template/:id` - 删除模板

### 项目接口
- `GET /api/project/list` - 获取项目列表
- `GET /api/project/:id` - 获取项目详情
- `POST /api/project` - 创建项目
- `POST /api/project/from-template` - 从模板创建项目
- `PUT /api/project` - 更新项目
- `DELETE /api/project/:id` - 删除项目
- `GET /api/project/statistics` - 获取项目统计

### 任务接口
- `GET /api/project/task/list` - 获取任务列表
- `GET /api/project/task/:id` - 获取任务详情
- `POST /api/project/task` - 创建任务
- `PUT /api/project/task` - 更新任务
- `PUT /api/project/task/batch` - 批量更新任务
- `DELETE /api/project/task/:id` - 删除任务
- `GET /api/project/task/gantt` - 获取甘特图数据
- `POST /api/project/task/dependency` - 创建任务依赖
- `DELETE /api/project/task/dependency/:id` - 删除任务依赖
- `POST /api/project/task/comment` - 创建任务评论

## 技术栈

### 后端
- Go + Gin
- GORM
- MySQL
- JWT认证

### 前端
- Vue 3 + TypeScript
- Element Plus
- DHTMLX Gantt (免费版)
- Vue Router
- Axios

## 测试报告

详细的API测试报告请查看: `art_admin_backend/API_TEST_REPORT.md`

测试结果: **10/10 接口测试通过 (100%)**

## 注意事项

1. **DHTMLX Gantt 免费版限制**
   - 免费版功能已足够使用
   - 如需高级功能(如资源视图、关键路径等),需要购买商业许可证

2. **超期检测**
   - 系统会自动检测任务是否超期
   - 超期任务会在甘特图中以橙色显示
   - 里程碑任务以红色显示

3. **任务依赖类型**
   - FS (Finish-to-Start): 完成到开始
   - SS (Start-to-Start): 开始到开始
   - FF (Finish-to-Finish): 完成到完成
   - SF (Start-to-Finish): 开始到完成

4. **权限控制**
   - 所有接口都需要JWT认证
   - 可以根据需要添加角色权限控制

## 后续优化建议

1. **功能增强**
   - 添加项目模板管理页面
   - 实现任务评论和附件上传
   - 添加项目看板视图
   - 实现项目报表和统计图表

2. **性能优化**
   - 大量任务时的甘特图性能优化
   - 实现任务数据的懒加载
   - 添加缓存机制

3. **用户体验**
   - 添加快捷键支持
   - 实现拖拽排序
   - 添加任务筛选和搜索
   - 支持导出为PDF/Excel

## 问题反馈

如有问题,请查看:
1. 后端日志: 检查 `art_admin_backend` 的控制台输出
2. 前端控制台: 检查浏览器开发者工具的Console
3. 网络请求: 检查浏览器开发者工具的Network标签

## 更新日志

### v1.0.0 (2025-11-20)
- ✅ 完成后端API开发
- ✅ 完成数据库设计和迁移
- ✅ 完成前端页面开发
- ✅ 集成DHTMLX Gantt
- ✅ 完成接口测试
- ✅ 完成路由配置
