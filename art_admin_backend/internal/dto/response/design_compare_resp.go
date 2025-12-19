package response

import "time"

// DesignImageInfoResponse 设计图信息响应
// Design image info response
type DesignImageInfoResponse struct {
	FileName string `json:"fileName"` // 原始文件名 / Original file name
	FilePath string `json:"filePath"` // 存储路径 / Storage path
	FileURL  string `json:"fileUrl"`  // 完整URL / Full URL
	FileType string `json:"fileType"` // 文件类型: image/cad/pdf / File type
	FileSize int64  `json:"fileSize"` // 文件大小 / File size in bytes
}

// DesignCompareResponse 设计比对响应（支持多图片、多文档）
// Design compare response supporting multiple images and documents
type DesignCompareResponse struct {
	ID             uint                      `json:"id"`
	ProjectID      uint                      `json:"projectId"`      // 关联项目ID / Associated project ID
	Name           string                    `json:"name"`           // 比对任务名称 / Compare task name
	DocumentIDs    []uint                    `json:"documentIds"`    // 关联需求文档ID列表 / Associated document IDs
	DesignImages   []DesignImageInfoResponse `json:"designImages"`   // 设计图列表 / Design images list
	MatchItems     []MatchItemResponse       `json:"matchItems"`     // 匹配项 / Match items
	DeviationItems []DeviationItemResponse   `json:"deviationItems"` // 偏差项 / Deviation items
	Suggestions    []SuggestionItemResponse  `json:"suggestions"`    // 建议项 / Suggestion items
	OverallScore   float64                   `json:"overallScore"`   // 整体匹配度 / Overall match score
	AnalysisStatus string                    `json:"analysisStatus"` // 分析状态 / Analysis status
	ErrorMessage   string                    `json:"errorMessage"`   // 错误信息 / Error message
	CreatedAt      time.Time                 `json:"createdAt"`
	UpdatedAt      time.Time                 `json:"updatedAt"`
}

// MatchItemResponse 匹配项响应
// Match item response
type MatchItemResponse struct {
	Requirement string `json:"requirement"` // 需求描述 / Requirement description
	DesignMatch string `json:"designMatch"` // 设计匹配点 / Design match point
	Score       int    `json:"score"`       // 匹配分数(0-100) / Match score (0-100)
	SourceImage string `json:"sourceImage"` // 来源设计图 / Source design image
	SourceDoc   string `json:"sourceDoc"`   // 来源文档 / Source document
}

// DeviationItemResponse 偏差项响应
// Deviation item response
type DeviationItemResponse struct {
	Location            string `json:"location"`            // 偏差位置 / Deviation location
	Content             string `json:"content"`             // 偏差内容 / Deviation content
	OriginalRequirement string `json:"originalRequirement"` // 原始需求条款 / Original requirement clause
	Severity            string `json:"severity"`            // 严重程度 / Severity (low/medium/high)
	SourceImage         string `json:"sourceImage"`         // 来源设计图 / Source design image
	SourceDoc           string `json:"sourceDoc"`           // 来源文档 / Source document
}

// SuggestionItemResponse 建议项响应
// Suggestion item response
type SuggestionItemResponse struct {
	Content    string `json:"content"`    // 建议内容 / Suggestion content
	Priority   string `json:"priority"`   // 优先级 / Priority (low/medium/high)
	CostImpact string `json:"costImpact"` // 成本影响 / Cost impact
}

// DesignCompareListResponse 设计比对列表响应
// Design compare list response
type DesignCompareListResponse struct {
	Records []DesignCompareResponse `json:"records"` // 记录列表 / Record list
	Current int                     `json:"current"` // 当前页码 / Current page
	Size    int                     `json:"size"`    // 每页条数 / Page size
	Total   int64                   `json:"total"`   // 总记录数 / Total count
}

// DesignCompareUploadResponse 设计图上传响应（支持多文件）
// Design images upload response supporting multiple files
type DesignCompareUploadResponse struct {
	ID             uint                      `json:"id"`             // 比对记录ID / Compare record ID
	ProjectID      uint                      `json:"projectId"`      // 项目ID / Project ID
	Name           string                    `json:"name"`           // 比对任务名称 / Compare task name
	DocumentIDs    []uint                    `json:"documentIds"`    // 文档ID列表 / Document IDs
	DesignImages   []DesignImageInfoResponse `json:"designImages"`   // 设计图列表 / Design images
	AnalysisStatus string                    `json:"analysisStatus"` // 分析状态 / Analysis status
}

// DesignCompareAnalysisResult 设计比对分析结果
// Design compare analysis result
type DesignCompareAnalysisResult struct {
	CompareID      uint                     `json:"compareId"`      // 比对记录ID / Compare record ID
	MatchItems     []MatchItemResponse      `json:"matchItems"`     // 匹配项 / Match items
	DeviationItems []DeviationItemResponse  `json:"deviationItems"` // 偏差项 / Deviation items
	Suggestions    []SuggestionItemResponse `json:"suggestions"`    // 建议项 / Suggestion items
	OverallScore   float64                  `json:"overallScore"`   // 整体匹配度 / Overall match score
}
