package request

// UploadDocumentRequest 上传文档请求
// Upload document request
type UploadDocumentRequest struct {
	ProjectID uint `form:"projectId" binding:"required"` // 关联项目ID / Associated project ID
}

// AnalyzeDocumentRequest 分析文档请求
// Analyze document request
type AnalyzeDocumentRequest struct {
	DocumentID uint `json:"documentId" binding:"required"` // 文档ID / Document ID
}

// BatchAnalyzeDocumentRequest 批量分析文档请求
// Batch analyze document request
type BatchAnalyzeDocumentRequest struct {
	DocumentIDs []uint `json:"documentIds" binding:"required,min=1"` // 文档ID列表 / Document ID list
}

// DocumentListRequest 文档列表请求
// Document list request
type DocumentListRequest struct {
	ProjectID      uint   `form:"projectId" binding:"required"`                                                 // 项目ID / Project ID
	Current        int    `form:"current" binding:"required,min=1" example:"1"`                                 // 当前页码 / Current page
	Size           int    `form:"size" binding:"required,min=1,max=100" example:"10"`                           // 每页数量 / Page size
	AnalysisStatus string `form:"analysisStatus" binding:"omitempty,oneof=pending processing completed failed"` // 分析状态筛选 / Analysis status filter
	FileType       string `form:"fileType" binding:"omitempty"`                                                 // 文件类型筛选 / File type filter
}

// DeleteDocumentRequest 删除文档请求
// Delete document request
type DeleteDocumentRequest struct {
	ID uint `json:"id" binding:"required"` // 文档ID / Document ID
}

// BatchDeleteDocumentRequest 批量删除文档请求
// Batch delete document request
type BatchDeleteDocumentRequest struct {
	IDs []uint `json:"ids" binding:"required,min=1"` // 文档ID列表 / Document ID list
}
