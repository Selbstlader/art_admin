package response

import "time"

// DesignVersionResponse 设计版本响应
// Design version response
type DesignVersionResponse struct {
	ID            uint                   `json:"id"`
	ProjectID     uint                   `json:"projectId"`     // 关联项目ID / Associated project ID
	VersionNumber int                    `json:"versionNumber"` // 版本号 / Version number
	VersionName   string                 `json:"versionName"`   // 版本名称 / Version name
	Description   string                 `json:"description"`   // 版本说明 / Version description
	DesignImages  []DesignImageResponse  `json:"designImages"`  // 设计图列表 / Design images
	CadFileIDs    []uint                 `json:"cadFileIds"`    // 关联CAD文件ID列表 / Associated CAD file IDs
	CadFiles      []CadFileBriefResponse `json:"cadFiles"`      // 关联CAD文件信息 / Associated CAD files
	LayoutInfo    []LayoutInfoResponse   `json:"layoutInfo"`    // 布局信息 / Layout info
	AreaInfo      []AreaInfoResponse     `json:"areaInfo"`      // 面积信息 / Area info
	StyleInfo     []StyleInfoResponse    `json:"styleInfo"`     // 风格信息 / Style info
	MaterialInfo  []MaterialInfoResponse `json:"materialInfo"`  // 材料信息 / Material info
	Status        string                 `json:"status"`        // 状态 / Status
	CreatedBy     uint                   `json:"createdBy"`     // 创建者ID / Creator ID
	CreatedAt     time.Time              `json:"createdAt"`
	UpdatedAt     time.Time              `json:"updatedAt"`
}

// CadFileBriefResponse CAD文件简要响应
// CAD file brief response
type CadFileBriefResponse struct {
	ID            uint   `json:"id"`
	FileName      string `json:"fileName"`      // 文件名 / File name
	FileFormat    string `json:"fileFormat"`    // 文件格式 / File format
	FilePath      string `json:"filePath"`      // 文件路径 / File path
	IsAIGenerated bool   `json:"isAiGenerated"` // 是否AI生成 / Whether AI generated
}

// DesignImageResponse 设计图响应
// Design image response
type DesignImageResponse struct {
	URL         string `json:"url"`         // 图片URL / Image URL
	Name        string `json:"name"`        // 图片名称 / Image name
	Description string `json:"description"` // 图片描述 / Image description
	Type        string `json:"type"`        // 图片类型 / Image type
}

// LayoutInfoResponse 布局信息响应
// Layout info response
type LayoutInfoResponse struct {
	Zone     string  `json:"zone"`     // 区域名称 / Zone name
	Area     float64 `json:"area"`     // 面积 / Area
	Position string  `json:"position"` // 位置描述 / Position description
}

// AreaInfoResponse 面积信息响应
// Area info response
type AreaInfoResponse struct {
	Name  string  `json:"name"`  // 区域名称 / Area name
	Value float64 `json:"value"` // 面积值 / Area value
	Unit  string  `json:"unit"`  // 单位 / Unit
}

// StyleInfoResponse 风格信息响应
// Style info response
type StyleInfoResponse struct {
	Category string `json:"category"` // 风格类别 / Style category
	Value    string `json:"value"`    // 风格值 / Style value
	Details  string `json:"details"`  // 详细描述 / Details
}

// MaterialInfoResponse 材料信息响应
// Material info response
type MaterialInfoResponse struct {
	Name     string  `json:"name"`     // 材料名称 / Material name
	Category string  `json:"category"` // 材料类别 / Material category
	Quantity float64 `json:"quantity"` // 数量 / Quantity
	Unit     string  `json:"unit"`     // 单位 / Unit
}

// DesignVersionListResponse 设计版本列表响应
// Design version list response
type DesignVersionListResponse struct {
	Records []DesignVersionResponse `json:"records"` // 记录列表 / Record list
	Current int                     `json:"current"` // 当前页码 / Current page
	Size    int                     `json:"size"`    // 每页条数 / Page size
	Total   int64                   `json:"total"`   // 总记录数 / Total count
}

// VersionCompareResponse 版本对比响应
// Version compare response
type VersionCompareResponse struct {
	ID              uint                        `json:"id"`
	ProjectID       uint                        `json:"projectId"`       // 项目ID / Project ID
	VersionAID      uint                        `json:"versionAId"`      // 版本A ID / Version A ID
	VersionBID      uint                        `json:"versionBId"`      // 版本B ID / Version B ID
	VersionA        *DesignVersionBriefResponse `json:"versionA"`        // 版本A简要信息 / Version A brief
	VersionB        *DesignVersionBriefResponse `json:"versionB"`        // 版本B简要信息 / Version B brief
	LayoutChanges   []ChangeItemResponse        `json:"layoutChanges"`   // 布局变化 / Layout changes
	AreaChanges     []ChangeItemResponse        `json:"areaChanges"`     // 面积调整 / Area changes
	ElementChanges  []ChangeItemResponse        `json:"elementChanges"`  // 元素增减 / Element changes
	StyleChanges    []ChangeItemResponse        `json:"styleChanges"`    // 风格变化 / Style changes
	MaterialChanges []ChangeItemResponse        `json:"materialChanges"` // 材料变化 / Material changes
	Summary         string                      `json:"summary"`         // 对比摘要 / Comparison summary
	CompareStatus   string                      `json:"compareStatus"`   // 对比状态 / Compare status
	CreatedAt       time.Time                   `json:"createdAt"`
	UpdatedAt       time.Time                   `json:"updatedAt"`
}

// DesignVersionBriefResponse 设计版本简要响应
// Design version brief response
type DesignVersionBriefResponse struct {
	ID            uint   `json:"id"`
	VersionNumber int    `json:"versionNumber"` // 版本号 / Version number
	VersionName   string `json:"versionName"`   // 版本名称 / Version name
	Status        string `json:"status"`        // 状态 / Status
}

// ChangeItemResponse 变化项响应
// Change item response
type ChangeItemResponse struct {
	Field       string `json:"field"`       // 变化字段 / Changed field
	OldValue    string `json:"oldValue"`    // 旧值 / Old value
	NewValue    string `json:"newValue"`    // 新值 / New value
	ChangeType  string `json:"changeType"`  // 变化类型 / Change type (added/removed/modified)
	Description string `json:"description"` // 变化描述 / Change description
}

// VersionCompareListResponse 版本对比列表响应
// Version compare list response
type VersionCompareListResponse struct {
	Records []VersionCompareResponse `json:"records"` // 记录列表 / Record list
	Current int                      `json:"current"` // 当前页码 / Current page
	Size    int                      `json:"size"`    // 每页条数 / Page size
	Total   int64                    `json:"total"`   // 总记录数 / Total count
}

// VersionDiffResponse 版本差异响应（用于左右分屏展示）
// Version diff response (for side-by-side display)
type VersionDiffResponse struct {
	VersionA        *DesignVersionResponse `json:"versionA"`        // 版本A完整信息 / Version A full info
	VersionB        *DesignVersionResponse `json:"versionB"`        // 版本B完整信息 / Version B full info
	LayoutChanges   []ChangeItemResponse   `json:"layoutChanges"`   // 布局变化 / Layout changes
	AreaChanges     []ChangeItemResponse   `json:"areaChanges"`     // 面积调整 / Area changes
	ElementChanges  []ChangeItemResponse   `json:"elementChanges"`  // 元素增减 / Element changes
	StyleChanges    []ChangeItemResponse   `json:"styleChanges"`    // 风格变化 / Style changes
	MaterialChanges []ChangeItemResponse   `json:"materialChanges"` // 材料变化 / Material changes
	Summary         string                 `json:"summary"`         // 对比摘要 / Comparison summary
}

// AIAnalyzeVersionDiffResponse AI版本对比分析响应
// AI version diff analysis response
type AIAnalyzeVersionDiffResponse struct {
	TaskID  uint   `json:"taskId"`  // 任务ID / Task ID
	Message string `json:"message"` // 提示消息 / Message
}
