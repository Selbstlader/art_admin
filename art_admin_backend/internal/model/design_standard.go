package model

import (
	"time"

	"gorm.io/gorm"
)

// DesignStandard 设计规范模型
// Design standard model for storing design compliance standards
type DesignStandard struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Code            string         `gorm:"size:100;not null;index" json:"code"`          // 规范编号 / Standard code
	Name            string         `gorm:"size:255;not null" json:"name"`                // 规范名称 / Standard name
	Category        string         `gorm:"size:100;not null;index" json:"category"`      // 规范类别 / Category (fire/accessibility/environmental/safety/other)
	Content         string         `gorm:"type:text;not null" json:"content"`            // 规范内容 / Standard content
	ApplicableTypes string         `gorm:"type:json" json:"applicableTypes"`             // 适用项目类型(JSON) / Applicable project types
	Version         string         `gorm:"size:50" json:"version"`                       // 规范版本 / Version
	EffectiveDate   *time.Time     `json:"effectiveDate"`                                // 生效日期 / Effective date
	Source          string         `gorm:"size:255" json:"source"`                       // 来源 / Source (国家标准/行业标准/地方标准)
	Interpretation  string         `gorm:"type:text" json:"interpretation"`              // 规范解读 / Interpretation
	Status          string         `gorm:"size:50;default:'active';index" json:"status"` // 状态 / Status (active/deprecated)
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (DesignStandard) TableName() string {
	return "design_standards"
}

// StandardCategory 规范类别常量
// Standard category constants
const (
	StandardCategoryFire          = "fire"          // 消防规范
	StandardCategoryAccessibility = "accessibility" // 无障碍规范
	StandardCategoryEnvironmental = "environmental" // 环保规范
	StandardCategorySafety        = "safety"        // 安全规范
	StandardCategoryOther         = "other"         // 其他规范
)

// StandardStatus 规范状态常量
// Standard status constants
const (
	StandardStatusActive     = "active"     // 有效
	StandardStatusDeprecated = "deprecated" // 已废止
)

// ComplianceCheckResult 合规检查结果模型
// Compliance check result model for storing check results
type ComplianceCheckResult struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	ProjectID        uint           `gorm:"index;not null" json:"projectId"`                    // 关联项目ID / Associated project ID
	CheckType        string         `gorm:"size:100" json:"checkType"`                          // 检查类型 / Check type (full/partial)
	PassedItems      string         `gorm:"type:json" json:"passedItems"`                       // 通过项(JSON) / Passed items
	FailedItems      string         `gorm:"type:json" json:"failedItems"`                       // 不合规项(JSON) / Failed items
	Suggestions      string         `gorm:"type:json" json:"suggestions"`                       // 改进建议(JSON) / Suggestions
	OverallScore     float64        `gorm:"default:0" json:"overallScore"`                      // 整体合规分数(0-100) / Overall compliance score
	CheckStatus      string         `gorm:"size:50;default:'pending';index" json:"checkStatus"` // 检查状态 / Check status
	ErrorMessage     string         `gorm:"size:500" json:"errorMessage"`                       // 错误信息 / Error message
	CheckedStandards string         `gorm:"type:json" json:"checkedStandards"`                  // 已检查的规范ID列表(JSON) / Checked standard IDs
	UserID           uint           `gorm:"index" json:"userId"`                                // 执行检查的用户ID / User who performed the check
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系 / Associations
	Project *DesignerProject `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
}

// TableName 表名
func (ComplianceCheckResult) TableName() string {
	return "compliance_check_results"
}

// CheckStatus 检查状态常量
// Check status constants
const (
	CheckStatusPending    = "pending"    // 待检查
	CheckStatusProcessing = "processing" // 检查中
	CheckStatusCompleted  = "completed"  // 已完成
	CheckStatusFailed     = "failed"     // 检查失败
)

// CheckType 检查类型常量
// Check type constants
const (
	CheckTypeFull    = "full"    // 全面检查
	CheckTypePartial = "partial" // 部分检查
)

// PassedItem 通过项结构
// Passed item structure for compliance check results
type PassedItem struct {
	StandardID   uint   `json:"standardId"`   // 规范ID / Standard ID
	StandardCode string `json:"standardCode"` // 规范编号 / Standard code
	StandardName string `json:"standardName"` // 规范名称 / Standard name
	Category     string `json:"category"`     // 规范类别 / Category
	Description  string `json:"description"`  // 通过说明 / Pass description
}

// FailedItem 不合规项结构
// Failed item structure for compliance check results
type FailedItem struct {
	StandardID          uint   `json:"standardId"`          // 规范ID / Standard ID
	StandardCode        string `json:"standardCode"`        // 规范编号 / Standard code
	StandardName        string `json:"standardName"`        // 规范名称 / Standard name
	Category            string `json:"category"`            // 规范类别 / Category
	ViolationContent    string `json:"violationContent"`    // 违规内容 / Violation content
	OriginalRequirement string `json:"originalRequirement"` // 原始规范条款 / Original requirement
	Severity            string `json:"severity"`            // 严重程度 / Severity (low/medium/high/critical)
	Location            string `json:"location"`            // 问题位置 / Problem location
	Suggestion          string `json:"suggestion"`          // 修改建议 / Modification suggestion
}

// ComplianceSuggestion 合规建议结构
// Compliance suggestion structure
type ComplianceSuggestion struct {
	Content    string `json:"content"`    // 建议内容 / Suggestion content
	Priority   string `json:"priority"`   // 优先级 / Priority (low/medium/high)
	Category   string `json:"category"`   // 相关类别 / Related category
	CostImpact string `json:"costImpact"` // 成本影响 / Cost impact
	Reference  string `json:"reference"`  // 参考规范 / Reference standard
}

// Severity 严重程度常量
// Severity constants
const (
	SeverityLow      = "low"      // 低
	SeverityMedium   = "medium"   // 中
	SeverityHigh     = "high"     // 高
	SeverityCritical = "critical" // 严重
)
