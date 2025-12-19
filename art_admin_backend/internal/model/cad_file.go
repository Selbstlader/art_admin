package model

import (
	"time"

	"gorm.io/gorm"
)

// CadFile CAD文件模型
// CAD file model for storing uploaded CAD files and parsing results
type CadFile struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	ProjectID     uint           `gorm:"index;not null" json:"projectId"`                    // 关联项目ID / Associated project ID
	FileName      string         `gorm:"size:255;not null" json:"fileName"`                  // 文件名 / File name
	FilePath      string         `gorm:"size:500" json:"filePath"`                           // 文件路径 / File path
	OriginalPath  string         `gorm:"size:500;not null" json:"originalPath"`              // 原始文件路径 / Original file path
	ParsedPath    string         `gorm:"size:500" json:"parsedPath"`                         // 解析后数据路径 / Parsed data path
	FileFormat    string         `gorm:"size:20" json:"fileFormat"`                          // 文件格式 / File format (dwg/dxf)
	ParseStatus   string         `gorm:"size:50;default:'pending';index" json:"parseStatus"` // 解析状态 / Parse status (pending/processing/completed/failed)
	LayerCount    int            `gorm:"default:0" json:"layerCount"`                        // 图层数量 / Layer count
	Layers        string         `gorm:"type:json;default:'[]'" json:"layers"`               // 图层列表(JSON) / Layer list in JSON format
	Has3D         bool           `gorm:"column:has_3d;default:false" json:"has3d"`           // 是否包含3D信息 / Whether contains 3D information
	FileSize      int64          `gorm:"default:0" json:"fileSize"`                          // 文件大小(字节) / File size in bytes
	IsAIGenerated bool           `gorm:"default:false" json:"isAiGenerated"`                 // 是否AI生成 / Whether AI generated
	ErrorMessage  string         `gorm:"size:500" json:"errorMessage"`                       // 错误信息 / Error message if parsing failed
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系 / Associations
	Project *DesignerProject `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
}

// TableName 表名
func (CadFile) TableName() string {
	return "cad_files"
}
