package request

// UploadDesignImagesRequest 批量上传设计图请求（支持多图片、多文档、CAD文件）
// Batch upload design images request supporting multiple images, documents and CAD files
type UploadDesignImagesRequest struct {
	ProjectID         uint   `form:"projectId" binding:"required"`   // 关联项目ID / Associated project ID
	DocumentIDs       string `form:"documentIds" binding:"required"` // 关联需求文档ID列表(逗号分隔) / Document IDs (comma separated)
	Name              string `form:"name"`                           // 比对任务名称 / Compare task name
	ExistingImageUrls string `form:"existingImageUrls"`              // 已有效果图URL列表(逗号分隔) / Existing image URLs (comma separated)
}

// CreateCompareWithCadRequest 使用已有CAD文件创建比对任务请求
// Create compare task with existing CAD files request
type CreateCompareWithCadRequest struct {
	ProjectID   uint   `json:"projectId" binding:"required"`   // 关联项目ID / Associated project ID
	DocumentIDs []uint `json:"documentIds" binding:"required"` // 关联需求文档ID列表 / Document ID list
	CadFileIDs  []uint `json:"cadFileIds" binding:"required"`  // CAD文件ID列表 / CAD file ID list
	Name        string `json:"name"`                           // 比对任务名称 / Compare task name
}

// CreateCompareWithImagesRequest 使用已有效果图创建比对任务请求
// Create compare task with existing render images request
type CreateCompareWithImagesRequest struct {
	ProjectID   uint     `json:"projectId" binding:"required"`   // 关联项目ID / Associated project ID
	DocumentIDs []uint   `json:"documentIds" binding:"required"` // 关联需求文档ID列表 / Document ID list
	ImageUrls   []string `json:"imageUrls" binding:"required"`   // 效果图URL列表 / Render image URL list
	Name        string   `json:"name"`                           // 比对任务名称 / Compare task name
}

// AnalyzeDesignCompareRequest 分析设计比对请求
// Analyze design comparison request
type AnalyzeDesignCompareRequest struct {
	CompareID uint `json:"compareId" binding:"required"` // 比对记录ID / Compare record ID
}

// DesignCompareListRequest 设计比对列表请求
// Design compare list request
type DesignCompareListRequest struct {
	ProjectID      uint   `form:"projectId" binding:"required"`                                                 // 项目ID / Project ID
	Current        int    `form:"current" binding:"required,min=1" example:"1"`                                 // 当前页码 / Current page
	Size           int    `form:"size" binding:"required,min=1,max=100" example:"10"`                           // 每页数量 / Page size
	AnalysisStatus string `form:"analysisStatus" binding:"omitempty,oneof=pending processing completed failed"` // 分析状态筛选 / Analysis status filter
}

// BatchDeleteDesignCompareRequest 批量删除设计比对请求
// Batch delete design compare request
type BatchDeleteDesignCompareRequest struct {
	IDs []uint `json:"ids" binding:"required,min=1"` // 比对记录ID列表 / Compare record ID list
}
