package request

// CreateDesignStandardRequest 创建设计规范请求
// Create design standard request
type CreateDesignStandardRequest struct {
	Code            string   `json:"code" binding:"required"`     // 规范编号 / Standard code
	Name            string   `json:"name" binding:"required"`     // 规范名称 / Standard name
	Category        string   `json:"category" binding:"required"` // 规范类别 / Category
	Content         string   `json:"content" binding:"required"`  // 规范内容 / Content
	ApplicableTypes []string `json:"applicableTypes"`             // 适用项目类型 / Applicable types
	Version         string   `json:"version"`                     // 版本 / Version
	EffectiveDate   string   `json:"effectiveDate"`               // 生效日期 / Effective date
	Source          string   `json:"source"`                      // 来源 / Source
	Interpretation  string   `json:"interpretation"`              // 规范解读 / Interpretation
}

// UpdateDesignStandardRequest 更新设计规范请求
// Update design standard request
type UpdateDesignStandardRequest struct {
	ID              uint     `json:"id" binding:"required"`       // 规范ID / Standard ID
	Code            string   `json:"code" binding:"required"`     // 规范编号 / Standard code
	Name            string   `json:"name" binding:"required"`     // 规范名称 / Standard name
	Category        string   `json:"category" binding:"required"` // 规范类别 / Category
	Content         string   `json:"content" binding:"required"`  // 规范内容 / Content
	ApplicableTypes []string `json:"applicableTypes"`             // 适用项目类型 / Applicable types
	Version         string   `json:"version"`                     // 版本 / Version
	EffectiveDate   string   `json:"effectiveDate"`               // 生效日期 / Effective date
	Source          string   `json:"source"`                      // 来源 / Source
	Interpretation  string   `json:"interpretation"`              // 规范解读 / Interpretation
	Status          string   `json:"status"`                      // 状态 / Status
}

// DesignStandardListRequest 设计规范列表请求
// Design standard list request
type DesignStandardListRequest struct {
	Current  int    `form:"current" binding:"required,min=1" example:"1"`       // 当前页码 / Current page
	Size     int    `form:"size" binding:"required,min=1,max=100" example:"10"` // 每页数量 / Page size
	Category string `form:"category"`                                           // 类别筛选 / Category filter
	Status   string `form:"status"`                                             // 状态筛选 / Status filter
	Keyword  string `form:"keyword"`                                            // 关键字搜索 / Keyword search
}

// ComplianceCheckRequest 合规检查请求
// Compliance check request
type ComplianceCheckRequest struct {
	ProjectID  uint     `json:"projectId" binding:"required"` // 项目ID / Project ID
	Categories []string `json:"categories"`                   // 检查类别(可选，为空则全部检查) / Categories to check
	CheckType  string   `json:"checkType"`                    // 检查类型(full/partial) / Check type
}

// ComplianceCheckListRequest 合规检查结果列表请求
// Compliance check result list request
type ComplianceCheckListRequest struct {
	ProjectID   uint   `form:"projectId" binding:"required"`                                              // 项目ID / Project ID
	Current     int    `form:"current" binding:"required,min=1" example:"1"`                              // 当前页码 / Current page
	Size        int    `form:"size" binding:"required,min=1,max=100" example:"10"`                        // 每页数量 / Page size
	CheckStatus string `form:"checkStatus" binding:"omitempty,oneof=pending processing completed failed"` // 检查状态筛选 / Check status filter
}

// BatchDeleteDesignStandardRequest 批量删除设计规范请求
// Batch delete design standard request
type BatchDeleteDesignStandardRequest struct {
	IDs []uint `json:"ids" binding:"required,min=1"` // 规范ID列表 / Standard ID list
}

// BatchDeleteComplianceCheckRequest 批量删除合规检查结果请求
// Batch delete compliance check result request
type BatchDeleteComplianceCheckRequest struct {
	IDs []uint `json:"ids" binding:"required,min=1"` // 检查结果ID列表 / Check result ID list
}
