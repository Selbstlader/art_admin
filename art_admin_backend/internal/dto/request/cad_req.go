package request

// UploadCadFileRequest CAD文件上传请求
// Upload CAD file request
type UploadCadFileRequest struct {
	ProjectID uint `form:"projectId" binding:"required"` // 关联项目ID / Associated project ID
}

// ParseCadFileRequest 解析CAD文件请求
// Parse CAD file request
type ParseCadFileRequest struct {
	CadFileID uint `json:"cadFileId" binding:"required"` // CAD文件ID / CAD file ID
}

// CadFileListRequest CAD文件列表请求
// CAD file list request
type CadFileListRequest struct {
	ProjectID   uint   `form:"projectId" binding:"required"`                                              // 项目ID / Project ID
	Current     int    `form:"current" binding:"required,min=1" example:"1"`                              // 当前页码 / Current page
	Size        int    `form:"size" binding:"required,min=1,max=100" example:"10"`                        // 每页数量 / Page size
	ParseStatus string `form:"parseStatus" binding:"omitempty,oneof=pending processing completed failed"` // 解析状态筛选 / Parse status filter
	FileFormat  string `form:"fileFormat" binding:"omitempty,oneof=dwg dxf"`                              // 文件格式筛选 / File format filter
}

// DeleteCadFileRequest 删除CAD文件请求
// Delete CAD file request
type DeleteCadFileRequest struct {
	ID uint `json:"id" binding:"required"` // CAD文件ID / CAD file ID
}

// BatchDeleteCadFileRequest 批量删除CAD文件请求
// Batch delete CAD file request
type BatchDeleteCadFileRequest struct {
	IDs []uint `json:"ids" binding:"required,min=1"` // CAD文件ID列表 / CAD file ID list
}
