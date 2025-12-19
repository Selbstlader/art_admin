package response

import "time"

// DocumentResponse 文档响应
// Document response
type DocumentResponse struct {
	ID              uint      `json:"id"`
	ProjectID       uint      `json:"projectId"`       // 关联项目ID / Associated project ID
	FileName        string    `json:"fileName"`        // 文件名 / File name
	FilePath        string    `json:"filePath"`        // 存储路径 / Storage path
	FileType        string    `json:"fileType"`        // 文件类型 / File type
	FileSize        int64     `json:"fileSize"`        // 文件大小 / File size
	AnalysisStatus  string    `json:"analysisStatus"`  // 分析状态 / Analysis status
	ErrorMessage    string    `json:"errorMessage"`    // 分析错误信息 / Analysis error message
	Keywords        []string  `json:"keywords"`        // 关键字列表 / Keywords list
	Summary         string    `json:"summary"`         // 文档摘要 / Document summary
	ProjectName     string    `json:"projectName"`     // 提取的项目名称 / Extracted project name
	ExtractedArea   float64   `json:"extractedArea"`   // 提取的面积 / Extracted area
	ExtractedBudget float64   `json:"extractedBudget"` // 提取的预算 / Extracted budget
	ExtractedStyle  string    `json:"extractedStyle"`  // 提取的风格 / Extracted style
	FunctionalZones []string  `json:"functionalZones"` // 功能分区 / Functional zones
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// DocumentListResponse 文档列表响应
// Document list response
type DocumentListResponse struct {
	Records []DocumentResponse `json:"records"` // 记录列表 / Record list
	Current int                `json:"current"` // 当前页码 / Current page
	Size    int                `json:"size"`    // 每页条数 / Page size
	Total   int64              `json:"total"`   // 总记录数 / Total count
}

// DocumentAnalysisResult 文档分析结果
// Document analysis result
type DocumentAnalysisResult struct {
	DocumentID      uint     `json:"documentId"`      // 文档ID / Document ID
	ProjectName     string   `json:"projectName"`     // 项目名称 / Project name
	Area            float64  `json:"area"`            // 面积 / Area
	Budget          float64  `json:"budget"`          // 预算 / Budget
	Style           string   `json:"style"`           // 设计风格 / Design style
	FunctionalZones []string `json:"functionalZones"` // 功能分区 / Functional zones
	Keywords        []string `json:"keywords"`        // 关键字 / Keywords
	Summary         string   `json:"summary"`         // 摘要 / Summary
	MissingFields   []string `json:"missingFields"`   // 缺失字段 / Missing fields
	Suggestions     []string `json:"suggestions"`     // 补充建议 / Suggestions
}

// DocumentUploadResponse 文档上传响应
// Document upload response
type DocumentUploadResponse struct {
	ID        uint   `json:"id"`        // 文档ID / Document ID
	FileName  string `json:"fileName"`  // 文件名 / File name
	FileType  string `json:"fileType"`  // 文件类型 / File type
	FileSize  int64  `json:"fileSize"`  // 文件大小 / File size
	ProjectID uint   `json:"projectId"` // 项目ID / Project ID
	Status    string `json:"status"`    // 状态 / Status
}
