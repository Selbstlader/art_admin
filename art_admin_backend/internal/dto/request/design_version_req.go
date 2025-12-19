package request

// CreateDesignVersionRequest 创建设计版本请求
// Create design version request
type CreateDesignVersionRequest struct {
	ProjectID    uint     `json:"projectId" binding:"required"`  // 项目ID / Project ID
	VersionName  string   `json:"versionName" binding:"max=100"` // 版本名称 / Version name
	Description  string   `json:"description"`                   // 版本说明 / Version description
	DesignImages []string `json:"designImages"`                  // 设计图URL列表 / Design image URLs
	CadFileIDs   []uint   `json:"cadFileIds"`                    // 关联CAD文件ID列表 / Associated CAD file IDs
	LayoutInfo   string   `json:"layoutInfo"`                    // 布局信息(JSON) / Layout info JSON
	AreaInfo     string   `json:"areaInfo"`                      // 面积信息(JSON) / Area info JSON
	StyleInfo    string   `json:"styleInfo"`                     // 风格信息(JSON) / Style info JSON
	MaterialInfo string   `json:"materialInfo"`                  // 材料信息(JSON) / Material info JSON
}

// UpdateDesignVersionRequest 更新设计版本请求
// Update design version request
type UpdateDesignVersionRequest struct {
	ID           uint     `json:"id" binding:"required"`                                              // 版本ID / Version ID
	VersionName  string   `json:"versionName" binding:"max=100"`                                      // 版本名称 / Version name
	Description  string   `json:"description"`                                                        // 版本说明 / Version description
	DesignImages []string `json:"designImages"`                                                       // 设计图URL列表 / Design image URLs
	CadFileIDs   []uint   `json:"cadFileIds"`                                                         // 关联CAD文件ID列表 / Associated CAD file IDs
	LayoutInfo   string   `json:"layoutInfo"`                                                         // 布局信息(JSON) / Layout info JSON
	AreaInfo     string   `json:"areaInfo"`                                                           // 面积信息(JSON) / Area info JSON
	StyleInfo    string   `json:"styleInfo"`                                                          // 风格信息(JSON) / Style info JSON
	MaterialInfo string   `json:"materialInfo"`                                                       // 材料信息(JSON) / Material info JSON
	Status       string   `json:"status" binding:"omitempty,oneof=draft submitted approved rejected"` // 状态 / Status
}

// DesignVersionListRequest 设计版本列表请求
// Design version list request
type DesignVersionListRequest struct {
	ProjectID uint   `form:"projectId" binding:"required"`                                       // 项目ID / Project ID
	Current   int    `form:"current" binding:"required,min=1" example:"1"`                       // 当前页码 / Current page
	Size      int    `form:"size" binding:"required,min=1,max=100" example:"10"`                 // 每页数量 / Page size
	Status    string `form:"status" binding:"omitempty,oneof=draft submitted approved rejected"` // 状态筛选 / Status filter
}

// CompareDesignVersionsRequest 版本对比请求
// Compare design versions request
type CompareDesignVersionsRequest struct {
	ProjectID  uint `json:"projectId" binding:"required"`  // 项目ID / Project ID
	VersionAID uint `json:"versionAId" binding:"required"` // 版本A ID / Version A ID
	VersionBID uint `json:"versionBId" binding:"required"` // 版本B ID / Version B ID
}

// VersionCompareListRequest 版本对比列表请求
// Version compare list request
type VersionCompareListRequest struct {
	ProjectID     uint   `form:"projectId" binding:"required"`                                     // 项目ID / Project ID
	Current       int    `form:"current" binding:"required,min=1" example:"1"`                     // 当前页码 / Current page
	Size          int    `form:"size" binding:"required,min=1,max=100" example:"10"`               // 每页数量 / Page size
	CompareStatus string `form:"compareStatus" binding:"omitempty,oneof=pending completed failed"` // 对比状态筛选 / Compare status filter
}

// BatchDeleteDesignVersionsRequest 批量删除设计版本请求
// Batch delete design versions request
type BatchDeleteDesignVersionsRequest struct {
	IDs []uint `json:"ids" binding:"required,min=1"` // 版本ID列表 / Version ID list
}

// UploadVersionImageRequest 上传版本设计图请求
// Upload version design image request
type UploadVersionImageRequest struct {
	ProjectID uint `form:"projectId" binding:"required"` // 项目ID / Project ID
	VersionID uint `form:"versionId"`                    // 版本ID(可选，不传则创建新版本) / Version ID (optional)
}
