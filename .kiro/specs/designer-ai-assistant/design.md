# Design Document - 工装设计师AI辅助系统

## Overview

本系统是面向工装设计师的AI辅助设计平台，基于现有的 art-design-pro（Vue 3 + TypeScript）前端和 art_admin_backend（Go + Gin）后端架构进行扩展。系统通过集成火山AI能力，提供项目文档智能分析、设计图比对、CAD在线预览、成本估算等核心功能。

### 技术栈

- **前端**: Vue 3.5 + TypeScript + Element Plus + Pinia + Three.js
- **后端**: Go 1.24 + Gin 1.9 + GORM + MySQL
- **AI服务**: 火山引擎AI（文本理解、视觉分析、多模态对话）
- **CAD处理**: dxf-parser（前端）+ ODA File Converter（后端）
- **测试框架**: Vitest + fast-check（前端）、gopter（后端）

## Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                         Frontend (Vue 3)                            │
├─────────────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐              │
│  │ 项目管理模块  │  │ AI对话模块   │  │ CAD预览模块  │              │
│  └──────────────┘  └──────────────┘  └──────────────┘              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐              │
│  │ 文档分析模块  │  │ 成本估算模块  │  │ 材料库模块   │              │
│  └──────────────┘  └──────────────┘  └──────────────┘              │
├─────────────────────────────────────────────────────────────────────┤
│                         API Layer (Axios)                           │
└─────────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────────┐
│                      Backend (Go + Gin)                             │
├─────────────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐              │
│  │ 项目服务     │  │ 文档服务     │  │ AI服务       │              │
│  └──────────────┘  └──────────────┘  └──────────────┘              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐              │
│  │ CAD解析服务  │  │ 材料库服务   │  │ 成本计算服务  │              │
│  └──────────────┘  └──────────────┘  └──────────────┘              │
├─────────────────────────────────────────────────────────────────────┤
│                      Data Layer (GORM + MySQL)                      │
└─────────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────────┐
│                      External Services                              │
├─────────────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐              │
│  │ 火山AI服务   │  │ 文件存储     │  │ ODA Converter │              │
│  └──────────────┘  └──────────────┘  └──────────────┘              │
└─────────────────────────────────────────────────────────────────────┘
```

## Components and Interfaces

### 前端组件结构

```
src/views/designer-assistant/
├── index.vue                    # 模块入口页面
├── project/                     # 项目管理
│   ├── ProjectList.vue          # 项目列表
│   ├── ProjectDetail.vue        # 项目详情
│   └── ProjectCreate.vue        # 创建项目
├── document/                    # 文档分析
│   ├── DocumentUpload.vue       # 文档上传
│   ├── DocumentAnalysis.vue     # 分析结果展示
│   └── KeywordExtraction.vue    # 关键字提取
├── design-compare/              # 设计比对
│   ├── CompareUpload.vue        # 上传比对
│   └── CompareResult.vue        # 比对结果
├── cad-viewer/                  # CAD预览
│   ├── CadUpload.vue            # CAD上传
│   ├── CadViewer.vue            # Three.js渲染器
│   └── LayerControl.vue         # 图层控制
├── chat/                        # AI对话
│   └── DesignerChat.vue         # 设计师对话界面
├── material/                    # 材料库
│   ├── MaterialList.vue         # 材料列表
│   ├── MaterialRecommend.vue    # 智能推荐
│   └── MaterialSelect.vue       # 材料选择
├── cost/                        # 成本估算
│   ├── CostConfig.vue           # 成本配置
│   └── CostReport.vue           # 成本报告
└── compliance/                  # 合规检查
    ├── ComplianceCheck.vue      # 合规检查
    └── ComplianceReport.vue     # 合规报告
```

### 后端API接口

#### 项目管理 API
```
POST   /api/designer/projects              # 创建项目
GET    /api/designer/projects              # 获取项目列表
GET    /api/designer/projects/:id          # 获取项目详情
PUT    /api/designer/projects/:id          # 更新项目
DELETE /api/designer/projects/:id          # 删除项目
```

#### 文档分析 API
```
POST   /api/designer/documents/upload      # 上传文档
POST   /api/designer/documents/analyze     # 分析文档
GET    /api/designer/documents/:id/keywords # 获取关键字
GET    /api/designer/documents/:id/summary  # 获取摘要
```

#### 设计比对 API
```
POST   /api/designer/compare/upload        # 上传设计图
POST   /api/designer/compare/analyze       # 执行比对分析
GET    /api/designer/compare/:id/result    # 获取比对结果
```

#### CAD预览 API
```
POST   /api/designer/cad/upload            # 上传CAD文件
GET    /api/designer/cad/:id/parse         # 获取解析数据
GET    /api/designer/cad/:id/layers        # 获取图层列表
```

#### 材料库 API
```
GET    /api/designer/materials             # 获取材料列表
GET    /api/designer/materials/:id         # 获取材料详情
POST   /api/designer/materials/recommend   # 智能推荐
POST   /api/designer/materials/calculate   # 计算用量
```

#### 成本估算 API
```
POST   /api/designer/cost/estimate         # 估算成本
GET    /api/designer/cost/:projectId       # 获取成本报告
PUT    /api/designer/cost/:projectId       # 更新成本配置
POST   /api/designer/cost/export           # 导出报告
```

#### AI对话 API
```
POST   /api/designer/chat                  # 发送消息
POST   /api/designer/chat/stream           # 流式对话
GET    /api/designer/chat/history          # 获取历史
```

## Data Models

### 项目模型 (Project)
```go
type DesignerProject struct {
    ID            uint      `gorm:"primaryKey"`
    Name          string    `gorm:"size:200;not null"`           // 项目名称
    Description   string    `gorm:"type:text"`                   // 项目描述
    Area          float64   `gorm:"default:0"`                   // 面积(平方米)
    Budget        float64   `gorm:"default:0"`                   // 预算
    Style         string    `gorm:"size:100"`                    // 设计风格
    Status        string    `gorm:"size:50;default:'draft'"`     // 状态
    UserID        uint      `gorm:"index"`                       // 所属用户
    CreatedAt     time.Time
    UpdatedAt     time.Time
}
```

### 文档模型 (Document)
```go
type ProjectDocument struct {
    ID            uint      `gorm:"primaryKey"`
    ProjectID     uint      `gorm:"index;not null"`              // 关联项目
    FileName      string    `gorm:"size:255;not null"`           // 文件名
    FilePath      string    `gorm:"size:500;not null"`           // 存储路径
    FileType      string    `gorm:"size:50"`                     // 文件类型
    FileSize      int64     `gorm:"default:0"`                   // 文件大小
    AnalysisStatus string   `gorm:"size:50;default:'pending'"`   // 分析状态
    Keywords      string    `gorm:"type:json"`                   // 提取的关键字
    Summary       string    `gorm:"type:text"`                   // 文档摘要
    CreatedAt     time.Time
    UpdatedAt     time.Time
}
```

### CAD文件模型 (CadFile)
```go
type CadFile struct {
    ID            uint      `gorm:"primaryKey"`
    ProjectID     uint      `gorm:"index;not null"`              // 关联项目
    FileName      string    `gorm:"size:255;not null"`           // 文件名
    OriginalPath  string    `gorm:"size:500;not null"`           // 原始文件路径
    ParsedPath    string    `gorm:"size:500"`                    // 解析后数据路径
    FileFormat    string    `gorm:"size:20"`                     // 文件格式(dwg/dxf)
    ParseStatus   string    `gorm:"size:50;default:'pending'"`   // 解析状态
    LayerCount    int       `gorm:"default:0"`                   // 图层数量
    Has3D         bool      `gorm:"default:false"`               // 是否包含3D信息
    CreatedAt     time.Time
    UpdatedAt     time.Time
}
```

### 材料模型 (Material)
```go
type Material struct {
    ID            uint      `gorm:"primaryKey"`
    Name          string    `gorm:"size:200;not null"`           // 材料名称
    Category      string    `gorm:"size:100;index"`              // 分类
    Specification string    `gorm:"size:200"`                    // 规格
    Unit          string    `gorm:"size:50"`                     // 单位
    UnitPrice     float64   `gorm:"default:0"`                   // 单价
    Brand         string    `gorm:"size:100"`                    // 品牌
    Supplier      string    `gorm:"size:200"`                    // 供应商
    Description   string    `gorm:"type:text"`                   // 描述
    ImageURL      string    `gorm:"size:500"`                    // 图片
    ApplicableScenes string `gorm:"type:json"`                   // 适用场景
    Status        string    `gorm:"size:50;default:'active'"`    // 状态
    CreatedAt     time.Time
    UpdatedAt     time.Time
}
```

### 成本估算模型 (CostEstimate)
```go
type CostEstimate struct {
    ID            uint      `gorm:"primaryKey"`
    ProjectID     uint      `gorm:"index;not null;uniqueIndex"`  // 关联项目
    MaterialCost  float64   `gorm:"default:0"`                   // 材料费
    LaborCost     float64   `gorm:"default:0"`                   // 人工费
    EquipmentCost float64   `gorm:"default:0"`                   // 设备费
    ManagementCost float64  `gorm:"default:0"`                   // 管理费
    TotalCost     float64   `gorm:"default:0"`                   // 总计
    BudgetLimit   float64   `gorm:"default:0"`                   // 预算上限
    Items         string    `gorm:"type:json"`                   // 明细项
    CreatedAt     time.Time
    UpdatedAt     time.Time
}
```

### 设计比对结果模型 (CompareResult)
```go
type DesignCompareResult struct {
    ID            uint      `gorm:"primaryKey"`
    ProjectID     uint      `gorm:"index;not null"`              // 关联项目
    DocumentID    uint      `gorm:"index"`                       // 关联需求文档
    DesignImagePath string  `gorm:"size:500;not null"`           // 设计图路径
    MatchItems    string    `gorm:"type:json"`                   // 匹配项
    DeviationItems string   `gorm:"type:json"`                   // 偏差项
    Suggestions   string    `gorm:"type:json"`                   // 建议项
    OverallScore  float64   `gorm:"default:0"`                   // 整体匹配度
    CreatedAt     time.Time
    UpdatedAt     time.Time
}
```



## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system-essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

Based on the prework analysis, the following correctness properties have been identified for property-based testing:

### Property 1: Document Summary Length Constraint
*For any* document analysis result, the generated summary length SHALL NOT exceed 500 characters.
**Validates: Requirements 1.3**

### Property 2: Document Analysis Result Structure
*For any* successfully parsed document, the analysis result SHALL contain all required fields: projectName, area, budget, style, and functionalZones.
**Validates: Requirements 1.2**

### Property 3: Compare Report Structure Completeness
*For any* design comparison result, the report SHALL contain three categories: matchItems, deviationItems, and suggestions.
**Validates: Requirements 2.2**

### Property 4: Deviation Item Structure
*For any* deviation item in a comparison result, it SHALL contain: location, content, and originalRequirement fields.
**Validates: Requirements 2.3**

### Property 5: Design Suggestion Count and Structure
*For any* design suggestion request, the system SHALL return at least 3 suggestions, and each suggestion SHALL contain: content, applicableScene, and costImpact fields.
**Validates: Requirements 3.2**

### Property 6: Project Cascade Delete
*For any* project deletion, all associated documents, analysis results, and CAD files SHALL be deleted from the database.
**Validates: Requirements 4.4**

### Property 7: Project Search Result Relevance
*For any* project search query, all returned projects SHALL match at least one of the search criteria (name, keyword, or time range).
**Validates: Requirements 4.3**

### Property 8: Material Cost Calculation Accuracy
*For any* material selection, the calculated total cost SHALL equal: quantity × unitPrice.
**Validates: Requirements 6.3**

### Property 9: Material Filter Result Consistency
*For any* material search with filters, all returned materials SHALL match all specified filter criteria (category, priceRange, brand).
**Validates: Requirements 6.4**

### Property 10: Material CRUD Data Integrity
*For any* material create/update operation, reading the material immediately after SHALL return the same data that was written.
**Validates: Requirements 6.5**

### Property 11: Version Association Correctness
*For any* project with multiple design versions, all versions SHALL be correctly associated with the same project ID.
**Validates: Requirements 7.1**

### Property 12: CAD Layer Data Completeness
*For any* successfully parsed CAD file, the layer list SHALL contain all layers present in the original file.
**Validates: Requirements 8.3**

### Property 13: CAD 3D Detection Accuracy
*For any* CAD file, the has3D flag SHALL be true if and only if the file contains 3D geometry data.
**Validates: Requirements 8.4**

### Property 14: Cost Report Sum Consistency
*For any* cost estimate, the totalCost SHALL equal: materialCost + laborCost + equipmentCost + managementCost.
**Validates: Requirements 10.3**

### Property 15: Budget Warning Trigger
*For any* cost estimate where totalCost > budgetLimit, the system SHALL return a budget exceeded warning.
**Validates: Requirements 10.4**

### Property 16: Compliance Report Structure
*For any* compliance check result, the report SHALL contain three categories: passedItems, failedItems, and suggestions.
**Validates: Requirements 11.2**

### Property 17: Role-Based Project Access
*For any* customer user, the returned project list SHALL only contain projects where the customer has been granted access.
**Validates: Requirements 12.2**

### Property 18: AI Retry Mechanism
*For any* AI service call that fails, the system SHALL retry up to 3 times with increasing intervals before returning an error.
**Validates: Requirements 13.2**

### Property 19: Token Usage Logging
*For any* AI service call, the system SHALL log the token consumption and response time.
**Validates: Requirements 13.4**

## Error Handling

### 前端错误处理

```typescript
// 统一错误处理类型
interface ApiError {
  code: number
  message: string
  details?: Record<string, any>
}

// 错误处理策略
const errorHandlers = {
  // 文件上传错误
  FILE_UPLOAD_ERROR: {
    unsupportedFormat: '不支持的文件格式，请上传PDF、Word或图片文件',
    fileTooLarge: '文件大小超过限制（最大50MB）',
    fileCorrupted: '文件已损坏，请重新上传'
  },
  // CAD解析错误
  CAD_PARSE_ERROR: {
    unsupportedVersion: 'CAD文件版本不支持，请转换为DXF格式后重试',
    parseTimeout: 'CAD文件解析超时，请尝试上传较小的文件',
    invalidGeometry: 'CAD文件包含无效几何数据'
  },
  // AI服务错误
  AI_SERVICE_ERROR: {
    connectionFailed: 'AI服务连接失败，请稍后重试',
    tokenExceeded: '请求内容过长，已自动分段处理',
    rateLimited: '请求过于频繁，请稍后重试'
  }
}
```

### 后端错误处理

```go
// 统一错误响应结构
type ErrorResponse struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Details interface{} `json:"details,omitempty"`
}

// 业务错误码定义
const (
    ErrCodeFileUpload     = 40001  // 文件上传错误
    ErrCodeFileParse      = 40002  // 文件解析错误
    ErrCodeCADParse       = 40003  // CAD解析错误
    ErrCodeAIService      = 40004  // AI服务错误
    ErrCodeValidation     = 40005  // 参数验证错误
    ErrCodePermission     = 40006  // 权限错误
    ErrCodeNotFound       = 40007  // 资源不存在
    ErrCodeBudgetExceeded = 40008  // 预算超出
)
```

## Testing Strategy

### 单元测试

使用 Vitest（前端）和 Go testing（后端）进行单元测试：

- 前端组件测试：验证组件渲染、用户交互、状态管理
- 后端服务测试：验证业务逻辑、数据验证、错误处理
- API接口测试：验证请求参数、响应格式、状态码

### 属性测试

使用 fast-check（前端）和 gopter（后端）进行属性测试：

**前端属性测试示例：**
```typescript
import fc from 'fast-check'
import { describe, it, expect } from 'vitest'

describe('Cost Calculation Properties', () => {
  /**
   * Feature: designer-ai-assistant, Property 8: Material Cost Calculation Accuracy
   * Validates: Requirements 6.3
   */
  it('should calculate material cost correctly for any quantity and price', () => {
    fc.assert(
      fc.property(
        fc.float({ min: 0, max: 10000 }),  // quantity
        fc.float({ min: 0, max: 100000 }), // unitPrice
        (quantity, unitPrice) => {
          const result = calculateMaterialCost(quantity, unitPrice)
          const expected = quantity * unitPrice
          return Math.abs(result - expected) < 0.01 // 允许浮点误差
        }
      ),
      { numRuns: 100 }
    )
  })

  /**
   * Feature: designer-ai-assistant, Property 14: Cost Report Sum Consistency
   * Validates: Requirements 10.3
   */
  it('should sum all cost categories correctly', () => {
    fc.assert(
      fc.property(
        fc.record({
          materialCost: fc.float({ min: 0, max: 1000000 }),
          laborCost: fc.float({ min: 0, max: 1000000 }),
          equipmentCost: fc.float({ min: 0, max: 1000000 }),
          managementCost: fc.float({ min: 0, max: 1000000 })
        }),
        (costs) => {
          const result = calculateTotalCost(costs)
          const expected = costs.materialCost + costs.laborCost + 
                          costs.equipmentCost + costs.managementCost
          return Math.abs(result - expected) < 0.01
        }
      ),
      { numRuns: 100 }
    )
  })
})
```

**后端属性测试示例：**
```go
package service_test

import (
    "testing"
    "github.com/leanovate/gopter"
    "github.com/leanovate/gopter/gen"
    "github.com/leanovate/gopter/prop"
)

/**
 * Feature: designer-ai-assistant, Property 6: Project Cascade Delete
 * Validates: Requirements 4.4
 */
func TestProjectCascadeDelete(t *testing.T) {
    parameters := gopter.DefaultTestParameters()
    parameters.MinSuccessfulTests = 100
    
    properties := gopter.NewProperties(parameters)
    
    properties.Property("deleting project removes all associated data", prop.ForAll(
        func(projectID uint) bool {
            // Setup: Create project with documents and CAD files
            project := createTestProject(projectID)
            docs := createTestDocuments(projectID, 3)
            cadFiles := createTestCadFiles(projectID, 2)
            
            // Action: Delete project
            err := projectService.DeleteProject(projectID)
            if err != nil {
                return false
            }
            
            // Verify: All associated data should be deleted
            remainingDocs := countDocumentsByProject(projectID)
            remainingCadFiles := countCadFilesByProject(projectID)
            
            return remainingDocs == 0 && remainingCadFiles == 0
        },
        gen.UInt().SuchThat(func(id uint) bool { return id > 0 }),
    ))
    
    properties.TestingRun(t)
}

/**
 * Feature: designer-ai-assistant, Property 9: Material Filter Result Consistency
 * Validates: Requirements 6.4
 */
func TestMaterialFilterConsistency(t *testing.T) {
    parameters := gopter.DefaultTestParameters()
    parameters.MinSuccessfulTests = 100
    
    properties := gopter.NewProperties(parameters)
    
    properties.Property("filtered materials match all criteria", prop.ForAll(
        func(category string, minPrice, maxPrice float64) bool {
            filter := MaterialFilter{
                Category: category,
                MinPrice: minPrice,
                MaxPrice: maxPrice,
            }
            
            results, _ := materialService.SearchMaterials(filter)
            
            for _, material := range results {
                if material.Category != category {
                    return false
                }
                if material.UnitPrice < minPrice || material.UnitPrice > maxPrice {
                    return false
                }
            }
            return true
        },
        gen.AnyString(),
        gen.Float64Range(0, 1000),
        gen.Float64Range(1000, 10000),
    ))
    
    properties.TestingRun(t)
}
```

### 测试配置

**前端测试配置 (vitest.config.ts):**
```typescript
export default defineConfig({
  test: {
    environment: 'happy-dom',
    include: ['src/**/*.{test,spec}.{js,ts}'],
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html']
    }
  }
})
```

**后端测试运行:**
```bash
# 运行所有测试
go test ./... -v

# 运行属性测试
go test ./internal/service/... -v -run Property

# 生成覆盖率报告
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

