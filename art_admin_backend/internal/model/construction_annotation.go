package model

import (
	"time"

	"gorm.io/gorm"
)

// ConstructionAnnotation 施工图标注模型
// Construction drawing annotation model for storing annotation data
type ConstructionAnnotation struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	ProjectID        uint           `gorm:"index;not null" json:"projectId"`                       // 关联项目ID / Associated project ID
	VersionID        *uint          `gorm:"index" json:"versionId"`                                // 关联设计版本ID / Associated design version ID
	CadFileID        uint           `gorm:"index;not null" json:"cadFileId"`                       // 关联CAD文件ID / Associated CAD file ID
	ImagePath        string         `gorm:"size:500;not null" json:"imagePath"`                    // 施工图图片路径 / Construction drawing image path
	AnalysisStatus   string         `gorm:"size:50;default:'pending';index" json:"analysisStatus"` // 分析状态 / Analysis status (pending/processing/completed/failed)
	ElementsDetected string         `gorm:"type:json;default:'[]'" json:"elementsDetected"`        // 检测到的元素(JSON) / Detected elements in JSON format
	Annotations      string         `gorm:"type:json;default:'[]'" json:"annotations"`             // 标注数据(JSON) / Annotation data in JSON format
	ErrorMessage     string         `gorm:"size:500" json:"errorMessage"`                          // 错误信息 / Error message if analysis failed
	CreatedBy        uint           `gorm:"index" json:"createdBy"`                                // 创建者ID / Creator user ID
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系 / Associations
	Project *DesignerProject `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	Version *DesignVersion   `gorm:"foreignKey:VersionID" json:"version,omitempty"`
	CadFile *CadFile         `gorm:"foreignKey:CadFileID" json:"cadFile,omitempty"`
}

// TableName 表名
func (ConstructionAnnotation) TableName() string {
	return "construction_annotations"
}

// DetectedElement 检测到的元素
// Detected element structure
type DetectedElement struct {
	Type        string                 `json:"type"`        // 元素类型 / Element type (dimension_line, material_area, equipment)
	Confidence  float64                `json:"confidence"`  // 置信度 / Confidence score
	BoundingBox BoundingBox            `json:"boundingBox"` // 边界框 / Bounding box
	Properties  map[string]interface{} `json:"properties"`  // 元素属性 / Element properties
}

// BoundingBox 边界框
// Bounding box coordinates
type BoundingBox struct {
	X      float64 `json:"x"`      // 左上角X坐标 / Top-left X coordinate
	Y      float64 `json:"y"`      // 左上角Y坐标 / Top-left Y coordinate
	Width  float64 `json:"width"`  // 宽度 / Width
	Height float64 `json:"height"` // 高度 / Height
}

// AnnotationItem 标注项
// Annotation item structure
type AnnotationItem struct {
	ID          string                 `json:"id"`          // 标注ID / Annotation ID
	Type        string                 `json:"type"`        // 标注类型 / Annotation type (dimension, material, process)
	Position    Position               `json:"position"`    // 位置 / Position
	Content     string                 `json:"content"`     // 标注内容 / Annotation content
	Style       AnnotationStyle        `json:"style"`       // 样式 / Style
	Properties  map[string]interface{} `json:"properties"`  // 扩展属性 / Extended properties
	CreatedAt   time.Time              `json:"createdAt"`   // 创建时间 / Creation time
	UpdatedAt   time.Time              `json:"updatedAt"`   // 更新时间 / Update time
	IsGenerated bool                   `json:"isGenerated"` // 是否自动生成 / Whether auto-generated
}

// Position 位置信息
// Position information
type Position struct {
	X      float64 `json:"x"`      // X坐标 / X coordinate
	Y      float64 `json:"y"`      // Y coordinate
	Anchor string  `json:"anchor"` // 锚点 / Anchor point (top-left, center, etc.)
}

// AnnotationStyle 标注样式
// Annotation style
type AnnotationStyle struct {
	FontSize   int    `json:"fontSize"`   // 字体大小 / Font size
	FontColor  string `json:"fontColor"`  // 字体颜色 / Font color
	LineColor  string `json:"lineColor"`  // 线条颜色 / Line color
	LineWidth  int    `json:"lineWidth"`  // 线条宽度 / Line width
	Background string `json:"background"` // 背景色 / Background color
}

// AnnotationExportConfig 标注导出配置
// Annotation export configuration
type AnnotationExportConfig struct {
	Format     string  `json:"format"`     // 导出格式 / Export format (pdf, png, jpg)
	Quality    int     `json:"quality"`    // 图片质量 / Image quality (1-100)
	Scale      float64 `json:"scale"`      // 缩放比例 / Scale ratio
	ShowLayers bool    `json:"showLayers"` // 是否显示图层 / Whether to show layers
}
