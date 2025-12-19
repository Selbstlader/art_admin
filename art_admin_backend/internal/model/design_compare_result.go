package model

import (
	"time"

	"gorm.io/gorm"
)

// DesignCompareResult 设计比对结果模型（支持多图片、多文档）
// Design comparison result model supporting multiple images and documents
type DesignCompareResult struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	ProjectID      uint           `gorm:"index;not null" json:"projectId"`                       // 关联项目ID / Associated project ID
	Name           string         `gorm:"size:200;default:''" json:"name"`                       // 比对任务名称 / Compare task name
	DocumentIDs    string         `gorm:"type:json" json:"documentIds"`                          // 关联需求文档ID列表(JSON数组) / Associated document IDs
	DesignImages   string         `gorm:"type:json" json:"designImages"`                         // 设计图列表(JSON数组) / Design images list
	MatchItems     string         `gorm:"type:json" json:"matchItems"`                           // 匹配项(JSON) / Match items in JSON format
	DeviationItems string         `gorm:"type:json" json:"deviationItems"`                       // 偏差项(JSON) / Deviation items in JSON format
	Suggestions    string         `gorm:"type:json" json:"suggestions"`                          // 建议项(JSON) / Suggestion items in JSON format
	OverallScore   float64        `gorm:"default:0" json:"overallScore"`                         // 整体匹配度(0-100) / Overall match score (0-100)
	AnalysisStatus string         `gorm:"size:50;default:'pending';index" json:"analysisStatus"` // 分析状态 / Analysis status
	ErrorMessage   string         `gorm:"size:1000" json:"errorMessage"`                         // 错误信息 / Error message if analysis failed
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// 旧字段兼容（保留以兼容未迁移的数据库）/ Legacy fields for compatibility
	// 使用指针类型允许 NULL 值，避免外键约束问题
	DocumentID      *uint  `gorm:"index" json:"-"`               // 旧：单个文档ID (nullable)
	DesignImagePath string `gorm:"size:500;default:''" json:"-"` // 旧：单个设计图路径

	// 关联关系 / Associations
	Project *DesignerProject `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
}

// TableName 表名
func (DesignCompareResult) TableName() string {
	return "design_compare_results"
}

// DesignImageInfo 设计图信息结构
// Design image info structure
type DesignImageInfo struct {
	FileName string `json:"fileName"` // 原始文件名 / Original file name
	FilePath string `json:"filePath"` // 存储路径 / Storage path
	FileType string `json:"fileType"` // 文件类型: image/cad/pdf / File type
	FileSize int64  `json:"fileSize"` // 文件大小 / File size in bytes
}

// MatchItem 匹配项结构
// Match item structure for comparison results
type MatchItem struct {
	Requirement string `json:"requirement"` // 需求描述 / Requirement description
	DesignMatch string `json:"designMatch"` // 设计匹配点 / Design match point
	Score       int    `json:"score"`       // 匹配分数(0-100) / Match score (0-100)
	SourceImage string `json:"sourceImage"` // 来源设计图 / Source design image
	SourceDoc   string `json:"sourceDoc"`   // 来源文档 / Source document
}

// DeviationItem 偏差项结构
// Deviation item structure for comparison results
type DeviationItem struct {
	Location            string `json:"location"`            // 偏差位置 / Deviation location
	Content             string `json:"content"`             // 偏差内容 / Deviation content
	OriginalRequirement string `json:"originalRequirement"` // 原始需求条款 / Original requirement clause
	Severity            string `json:"severity"`            // 严重程度 / Severity (low/medium/high)
	SourceImage         string `json:"sourceImage"`         // 来源设计图 / Source design image
	SourceDoc           string `json:"sourceDoc"`           // 来源文档 / Source document
}

// SuggestionItem 建议项结构
// Suggestion item structure for comparison results
type SuggestionItem struct {
	Content    string `json:"content"`    // 建议内容 / Suggestion content
	Priority   string `json:"priority"`   // 优先级 / Priority (low/medium/high)
	CostImpact string `json:"costImpact"` // 成本影响 / Cost impact
}
