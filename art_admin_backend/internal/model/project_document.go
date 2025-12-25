package model

import (
	"time"

	"gorm.io/gorm"
)

// ProjectDocument 项目文档模型
// Project document model for storing uploaded documents and analysis results
type ProjectDocument struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	ProjectID       uint           `gorm:"index;not null" json:"projectId"`                       // 关联项目ID / Associated project ID
	FileName        string         `gorm:"size:255;not null" json:"fileName"`                     // 文件名 / File name
	FilePath        string         `gorm:"size:500;not null" json:"filePath"`                     // 存储路径 / Storage path
	FileType        string         `gorm:"size:50" json:"fileType"`                               // 文件类型 / File type (pdf/word/image)
	FileSize        int64          `gorm:"default:0" json:"fileSize"`                             // 文件大小(字节) / File size in bytes
	AnalysisStatus  string         `gorm:"size:50;default:'pending';index" json:"analysisStatus"` // 分析状态 / Analysis status (pending/processing/completed/failed)
	ErrorMessage    string         `gorm:"size:500" json:"errorMessage"`                          // 分析错误信息 / Analysis error message
	Keywords        string         `gorm:"type:json" json:"keywords"`                             // 提取的关键字(JSON) / Extracted keywords in JSON format
	Summary         string         `gorm:"type:text" json:"summary"`                              // 文档摘要 / Document summary
	ProjectName     string         `gorm:"size:200" json:"projectName"`                           // 提取的项目名称 / Extracted project name
	ExtractedArea   float64        `gorm:"default:0" json:"extractedArea"`                        // 提取的面积 / Extracted area
	ExtractedBudget float64        `gorm:"default:0" json:"extractedBudget"`                      // 提取的预算 / Extracted budget
	ExtractedStyle  string         `gorm:"size:100" json:"extractedStyle"`                        // 提取的风格 / Extracted style
	FunctionalZones string         `gorm:"type:json" json:"functionalZones"`                      // 功能分区(JSON) / Functional zones in JSON format
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系 / Associations
	Project *DesignerProject `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
}

// TableName 表名
func (ProjectDocument) TableName() string {
	return "project_documents"
}
