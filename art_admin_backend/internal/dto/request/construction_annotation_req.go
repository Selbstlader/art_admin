package request

import "art_admin_backend/internal/model"

// AnalyzeConstructionDrawingRequest 分析施工图请求
// Analyze construction drawing request
type AnalyzeConstructionDrawingRequest struct {
	ProjectID uint   `json:"projectId" binding:"required" example:"1"`                         // 项目ID / Project ID
	CadFileID uint   `json:"cadFileId" binding:"required" example:"1"`                         // CAD文件ID / CAD file ID
	ImagePath string `json:"imagePath" binding:"required" example:"/path/to/construction.jpg"` // 施工图图片路径 / Construction drawing image path
}

// UpdateAnnotationsRequest 更新标注请求
// Update annotations request
type UpdateAnnotationsRequest struct {
	Annotations []model.AnnotationItem `json:"annotations" binding:"required"` // 标注数据 / Annotation data
}

// CreateAnnotationItemRequest 创建标注项请求
// Create annotation item request
type CreateAnnotationItemRequest struct {
	Type       string                 `json:"type" binding:"required" example:"dimension"` // 标注类型 / Annotation type
	Position   model.Position         `json:"position" binding:"required"`                 // 位置 / Position
	Content    string                 `json:"content" binding:"required" example:"3000mm"` // 标注内容 / Annotation content
	Style      model.AnnotationStyle  `json:"style"`                                       // 样式 / Style
	Properties map[string]interface{} `json:"properties"`                                  // 扩展属性 / Extended properties
}

// ExportAnnotationRequest 导出标注请求
// Export annotation request
type ExportAnnotationRequest struct {
	AnnotationID uint                         `json:"annotationId" binding:"required" example:"1"` // 标注ID / Annotation ID
	Config       model.AnnotationExportConfig `json:"config"`                                      // 导出配置 / Export configuration
}

// GetAnnotationsRequest 获取标注列表请求
// Get annotations list request
type GetAnnotationsRequest struct {
	ProjectID *uint  `form:"projectId" example:"1"`      // 项目ID（可选） / Project ID (optional)
	CadFileID *uint  `form:"cadFileId" example:"1"`      // CAD文件ID（可选） / CAD file ID (optional)
	Status    string `form:"status" example:"completed"` // 分析状态（可选） / Analysis status (optional)
	Page      int    `form:"page" example:"1"`           // 页码 / Page number
	PageSize  int    `form:"pageSize" example:"10"`      // 每页大小 / Page size
}

// Validate 验证请求参数
// Validate request parameters
func (r *GetAnnotationsRequest) Validate() {
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.PageSize <= 0 || r.PageSize > 100 {
		r.PageSize = 10
	}
}

// GetOffset 计算偏移量
// Calculate offset
func (r *GetAnnotationsRequest) GetOffset() int {
	return (r.Page - 1) * r.PageSize
}
