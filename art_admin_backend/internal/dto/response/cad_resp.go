package response

import "time"

// CadFileResponse CAD文件响应
// CAD file response
type CadFileResponse struct {
	ID           uint      `json:"id"`
	ProjectID    uint      `json:"projectId"`    // 关联项目ID / Associated project ID
	FileName     string    `json:"fileName"`     // 文件名 / File name
	OriginalPath string    `json:"originalPath"` // 原始文件路径 / Original file path
	ParsedPath   string    `json:"parsedPath"`   // 解析后数据路径 / Parsed data path
	FileFormat   string    `json:"fileFormat"`   // 文件格式 / File format (dwg/dxf)
	ParseStatus  string    `json:"parseStatus"`  // 解析状态 / Parse status
	LayerCount   int       `json:"layerCount"`   // 图层数量 / Layer count
	Layers       []string  `json:"layers"`       // 图层列表 / Layer list
	Has3D        bool      `json:"has3d"`        // 是否包含3D信息 / Whether contains 3D information
	FileSize     int64     `json:"fileSize"`     // 文件大小 / File size
	ErrorMessage string    `json:"errorMessage"` // 错误信息 / Error message
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// CadFileListResponse CAD文件列表响应
// CAD file list response
type CadFileListResponse struct {
	Records []CadFileResponse `json:"records"` // 记录列表 / Record list
	Current int               `json:"current"` // 当前页码 / Current page
	Size    int               `json:"size"`    // 每页条数 / Page size
	Total   int64             `json:"total"`   // 总记录数 / Total count
}

// CadFileUploadResponse CAD文件上传响应
// CAD file upload response
type CadFileUploadResponse struct {
	ID         uint   `json:"id"`         // CAD文件ID / CAD file ID
	FileName   string `json:"fileName"`   // 文件名 / File name
	FileFormat string `json:"fileFormat"` // 文件格式 / File format
	FileSize   int64  `json:"fileSize"`   // 文件大小 / File size
	ProjectID  uint   `json:"projectId"`  // 项目ID / Project ID
	Status     string `json:"status"`     // 状态 / Status
}

// CadLayerInfo CAD图层信息
// CAD layer information
type CadLayerInfo struct {
	Name        string `json:"name"`        // 图层名称 / Layer name
	Color       int    `json:"color"`       // 图层颜色 / Layer color
	LineType    string `json:"lineType"`    // 线型 / Line type
	Visible     bool   `json:"visible"`     // 是否可见 / Visibility
	Frozen      bool   `json:"frozen"`      // 是否冻结 / Frozen status
	Locked      bool   `json:"locked"`      // 是否锁定 / Locked status
	EntityCount int    `json:"entityCount"` // 实体数量 / Entity count
}

// CadEntityInfo CAD实体信息
// CAD entity information
type CadEntityInfo struct {
	Type       string                 `json:"type"`       // 实体类型 / Entity type (LINE, CIRCLE, ARC, etc.)
	Layer      string                 `json:"layer"`      // 所属图层 / Layer name
	Color      int                    `json:"color"`      // 颜色 / Color
	LineType   string                 `json:"lineType"`   // 线型 / Line type
	Properties map[string]interface{} `json:"properties"` // 属性 / Properties (dimensions, coordinates, etc.)
}

// CadParseResult CAD解析结果
// CAD parse result
type CadParseResult struct {
	CadFileID   uint            `json:"cadFileId"`   // CAD文件ID / CAD file ID
	FileName    string          `json:"fileName"`    // 文件名 / File name
	FileFormat  string          `json:"fileFormat"`  // 文件格式 / File format
	Layers      []CadLayerInfo  `json:"layers"`      // 图层列表 / Layer list
	LayerCount  int             `json:"layerCount"`  // 图层数量 / Layer count
	Has3D       bool            `json:"has3d"`       // 是否包含3D / Has 3D content
	BoundingBox *BoundingBox    `json:"boundingBox"` // 边界框 / Bounding box
	Entities    []CadEntityInfo `json:"entities"`    // 实体列表 / Entity list
	EntityCount int             `json:"entityCount"` // 实体总数 / Total entity count
	ParsedData  string          `json:"parsedData"`  // 解析后的JSON数据路径 / Parsed JSON data path
}

// BoundingBox 边界框
// Bounding box for CAD drawing
type BoundingBox struct {
	MinX float64 `json:"minX"`
	MinY float64 `json:"minY"`
	MinZ float64 `json:"minZ"`
	MaxX float64 `json:"maxX"`
	MaxY float64 `json:"maxY"`
	MaxZ float64 `json:"maxZ"`
}

// CadLayerListResponse CAD图层列表响应
// CAD layer list response
type CadLayerListResponse struct {
	CadFileID  uint           `json:"cadFileId"`  // CAD文件ID / CAD file ID
	FileName   string         `json:"fileName"`   // 文件名 / File name
	LayerCount int            `json:"layerCount"` // 图层数量 / Layer count
	Layers     []CadLayerInfo `json:"layers"`     // 图层列表 / Layer list
}
