# 工装设计助手 - 功能流程设计

## 完成状态

- [x] 数据模型设计
- [x] AI生成CAD后端API
- [x] AI生成CAD前端页面
- [x] 设计版本管理后端API
- [x] 设计版本管理前端页面
- [x] 施工图标注关联版本
- [x] 前端路由配置
- [ ] 实际AI生成CAD逻辑（当前为模拟实现）

## 完整业务流程

```
项目创建 → 文档分析 → AI生成CAD → CAD预览/编辑 → AI生成效果图 → 创建设计版本 → 施工图标注 → 版本对比 → 导出交付
```

## 功能模块说明

### 1. 项目管理 (ProjectList/ProjectDetail)
- 创建和管理设计项目
- 项目基本信息：名称、面积、预算、风格
- 作为所有功能的入口

### 2. 文档分析 (DocumentUpload)
- 上传需求文档（合同、设计要求等）
- AI分析提取关键信息
- 为后续CAD生成提供输入

### 3. AI生成CAD (CadGenerationList) - 新增
- 基于文档分析结果，AI生成CAD DXF文件
- 支持生成类型：
  - 平面布局图 (floor_plan)
  - 立面图 (elevation)
  - 节点详图 (detail)
  - 天花图 (ceiling)
  - 电气图 (electric)
- 生成后可预览、确认保存到项目

### 4. CAD预览/编辑 (CadViewer)
- 预览AI生成或上传的CAD文件
- 支持图层管理
- 可进行人工调整修改

### 5. AI生成效果图 (RenderHistory)
- 基于CAD文件生成效果图
- 支持多种渲染风格
- 效果图可用于设计版本

### 6. 设计版本管理 (VersionList) - 重构
- 创建设计版本，打包归档：
  - 设计图（效果图、平面图等）
  - CAD文件（关联已有CAD）
  - 版本说明
- 版本状态：草稿 → 已提交 → 已批准
- 支持版本对比

### 7. 施工图标注 (AnnotationList) - 重构
- 基于设计版本的CAD文件进行标注
- 关联关系：项目 → 版本 → CAD文件 → 标注
- AI辅助识别标注元素
- 支持导出标注结果

### 8. 版本对比 (VersionDiff)
- 选择同一项目的两个版本进行对比
- 展示差异：布局、面积、元素、风格、材料

## 数据模型关系

```
DesignerProject (项目)
    │
    ├── ProjectDocument (项目文档)
    │
    ├── CadGeneration (AI生成CAD任务)
    │       └── CadFile (生成的CAD文件)
    │
    ├── CadFile (CAD文件)
    │       └── isAiGenerated: 是否AI生成
    │
    ├── DesignVersion (设计版本)
    │       ├── designImages: 设计图列表
    │       ├── cadFileIds: 关联CAD文件ID
    │       └── ConstructionAnnotation (施工图标注)
    │
    └── DesignVersionCompare (版本对比结果)
```

## 路由配置

| 功能 | 路由 |
|------|------|
| 项目列表 | /designer/designer-assistant/project/ProjectList |
| 项目详情 | /designer/designer-assistant/project/ProjectDetail/:id |
| 文档分析 | /designer/designer-assistant/document/DocumentUpload |
| AI生成CAD | /designer/designer-assistant/cad-generation/CadGenerationList |
| CAD预览 | /designer/designer-assistant/cad-viewer/CadViewer |
| 设计版本 | /designer/designer-assistant/version-compare/VersionList |
| 版本对比 | /designer/designer-assistant/version-compare/VersionDiff |
| 施工图标注 | /designer/designer-assistant/construction-annotation/AnnotationList |

## API接口

### CAD生成接口
- POST /api/designer/cad-generations - 创建生成任务
- GET /api/designer/cad-generations - 获取任务列表
- GET /api/designer/cad-generations/:id - 获取任务详情
- POST /api/designer/cad-generations/:id/retry - 重试任务
- POST /api/designer/cad-generations/confirm - 确认保存CAD文件
- DELETE /api/designer/cad-generations/:id - 删除任务

### 设计版本接口
- POST /api/designer/versions - 创建版本
- GET /api/designer/versions - 获取版本列表
- GET /api/designer/versions/:id - 获取版本详情
- PUT /api/designer/versions/:id - 更新版本
- DELETE /api/designer/versions/:id - 删除版本
- POST /api/designer/versions/compare - 对比版本
- GET /api/designer/versions/diff - 获取版本差异

### 施工图标注接口
- POST /api/designer/construction-annotation/analyze - 分析施工图
- GET /api/designer/construction-annotation/list - 获取标注列表
- GET /api/designer/construction-annotation/:id - 获取标注详情
- PUT /api/designer/construction-annotation/:id/annotations - 更新标注
- DELETE /api/designer/construction-annotation/:id - 删除标注
