package response

import "time"

// DesignStandardResponse 设计规范响应
// Design standard response
type DesignStandardResponse struct {
	ID              uint      `json:"id"`
	Code            string    `json:"code"`            // 规范编号 / Standard code
	Name            string    `json:"name"`            // 规范名称 / Standard name
	Category        string    `json:"category"`        // 规范类别 / Category
	Content         string    `json:"content"`         // 规范内容 / Content
	ApplicableTypes []string  `json:"applicableTypes"` // 适用项目类型 / Applicable types
	Version         string    `json:"version"`         // 版本 / Version
	EffectiveDate   string    `json:"effectiveDate"`   // 生效日期 / Effective date
	Source          string    `json:"source"`          // 来源 / Source
	Interpretation  string    `json:"interpretation"`  // 规范解读 / Interpretation
	Status          string    `json:"status"`          // 状态 / Status
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// DesignStandardListResponse 设计规范列表响应
// Design standard list response
type DesignStandardListResponse struct {
	Records []DesignStandardResponse `json:"records"` // 记录列表 / Record list
	Current int                      `json:"current"` // 当前页码 / Current page
	Size    int                      `json:"size"`    // 每页条数 / Page size
	Total   int64                    `json:"total"`   // 总记录数 / Total count
}

// ComplianceCheckResponse 合规检查结果响应
// Compliance check result response
type ComplianceCheckResponse struct {
	ID               uint                           `json:"id"`
	ProjectID        uint                           `json:"projectId"`        // 项目ID / Project ID
	ProjectName      string                         `json:"projectName"`      // 项目名称 / Project name
	CheckType        string                         `json:"checkType"`        // 检查类型 / Check type
	PassedItems      []PassedItemResponse           `json:"passedItems"`      // 通过项 / Passed items
	FailedItems      []FailedItemResponse           `json:"failedItems"`      // 不合规项 / Failed items
	Suggestions      []ComplianceSuggestionResponse `json:"suggestions"`      // 改进建议 / Suggestions
	OverallScore     float64                        `json:"overallScore"`     // 整体合规分数 / Overall score
	CheckStatus      string                         `json:"checkStatus"`      // 检查状态 / Check status
	ErrorMessage     string                         `json:"errorMessage"`     // 错误信息 / Error message
	CheckedStandards []uint                         `json:"checkedStandards"` // 已检查的规范ID / Checked standard IDs
	UserID           uint                           `json:"userId"`           // 执行用户ID / User ID
	CreatedAt        time.Time                      `json:"createdAt"`
	UpdatedAt        time.Time                      `json:"updatedAt"`
}

// PassedItemResponse 通过项响应
// Passed item response
type PassedItemResponse struct {
	StandardID   uint   `json:"standardId"`   // 规范ID / Standard ID
	StandardCode string `json:"standardCode"` // 规范编号 / Standard code
	StandardName string `json:"standardName"` // 规范名称 / Standard name
	Category     string `json:"category"`     // 规范类别 / Category
	Description  string `json:"description"`  // 通过说明 / Pass description
}

// FailedItemResponse 不合规项响应
// Failed item response
type FailedItemResponse struct {
	StandardID          uint   `json:"standardId"`          // 规范ID / Standard ID
	StandardCode        string `json:"standardCode"`        // 规范编号 / Standard code
	StandardName        string `json:"standardName"`        // 规范名称 / Standard name
	Category            string `json:"category"`            // 规范类别 / Category
	ViolationContent    string `json:"violationContent"`    // 违规内容 / Violation content
	OriginalRequirement string `json:"originalRequirement"` // 原始规范条款 / Original requirement
	Severity            string `json:"severity"`            // 严重程度 / Severity
	Location            string `json:"location"`            // 问题位置 / Problem location
	Suggestion          string `json:"suggestion"`          // 修改建议 / Modification suggestion
}

// ComplianceSuggestionResponse 合规建议响应
// Compliance suggestion response
type ComplianceSuggestionResponse struct {
	Content    string `json:"content"`    // 建议内容 / Suggestion content
	Priority   string `json:"priority"`   // 优先级 / Priority
	Category   string `json:"category"`   // 相关类别 / Related category
	CostImpact string `json:"costImpact"` // 成本影响 / Cost impact
	Reference  string `json:"reference"`  // 参考规范 / Reference standard
}

// ComplianceCheckListResponse 合规检查结果列表响应
// Compliance check result list response
type ComplianceCheckListResponse struct {
	Records []ComplianceCheckResponse `json:"records"` // 记录列表 / Record list
	Current int                       `json:"current"` // 当前页码 / Current page
	Size    int                       `json:"size"`    // 每页条数 / Page size
	Total   int64                     `json:"total"`   // 总记录数 / Total count
}

// ComplianceCheckSummary 合规检查摘要
// Compliance check summary
type ComplianceCheckSummary struct {
	TotalChecks   int64   `json:"totalChecks"`   // 总检查次数 / Total checks
	PassedCount   int     `json:"passedCount"`   // 通过项数量 / Passed count
	FailedCount   int     `json:"failedCount"`   // 不合规项数量 / Failed count
	AverageScore  float64 `json:"averageScore"`  // 平均分数 / Average score
	LatestCheckAt string  `json:"latestCheckAt"` // 最近检查时间 / Latest check time
}

// StandardCategoryStats 规范类别统计
// Standard category statistics
type StandardCategoryStats struct {
	Category string `json:"category"` // 类别 / Category
	Count    int64  `json:"count"`    // 数量 / Count
}

// DesignStandardStats 设计规范统计响应
// Design standard statistics response
type DesignStandardStats struct {
	TotalCount    int64                   `json:"totalCount"`    // 总数量 / Total count
	ActiveCount   int64                   `json:"activeCount"`   // 有效数量 / Active count
	CategoryStats []StandardCategoryStats `json:"categoryStats"` // 类别统计 / Category stats
}
