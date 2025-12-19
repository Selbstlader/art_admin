package model

import (
	"time"

	"gorm.io/gorm"
)

// CadGeneration AI生成CAD任务模型
// CAD generation task model for AI-based CAD file generation
type CadGeneration struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	ProjectID       uint           `gorm:"index;not null" json:"projectId"`               // 关联项目ID / Associated project ID
	DocumentIDs     string         `gorm:"type:json" json:"documentIds"`                  // 关联文档ID列表(JSON) / Associated document IDs
	GenerationType  string         `gorm:"size:50;not null" json:"generationType"`        // 生成类型 / Generation type (floor_plan/elevation/detail)
	Prompt          string         `gorm:"type:text" json:"prompt"`                       // 生成提示词 / Generation prompt
	Parameters      string         `gorm:"type:json" json:"parameters"`                   // 生成参数(JSON) / Generation parameters
	Status          string         `gorm:"size:50;default:'pending';index" json:"status"` // 状态 / Status (pending/processing/completed/failed)
	Progress        int            `gorm:"default:0" json:"progress"`                     // 进度百分比 / Progress percentage (0-100)
	ResultFileID    *uint          `gorm:"index" json:"resultFileId"`                     // 生成的CAD文件ID / Generated CAD file ID
	ResultFilePath  string         `gorm:"size:500" json:"resultFilePath"`                // 生成的文件路径 / Generated file path
	PreviewImageURL string         `gorm:"size:500" json:"previewImageUrl"`               // 预览图URL / Preview image URL
	ErrorMessage    string         `gorm:"size:1000" json:"errorMessage"`                 // 错误信息 / Error message
	AIModel         string         `gorm:"size:100" json:"aiModel"`                       // 使用的AI模型 / AI model used
	ProcessingTime  int            `gorm:"default:0" json:"processingTime"`               // 处理耗时(秒) / Processing time in seconds
	CreatedBy       uint           `gorm:"index" json:"createdBy"`                        // 创建者ID / Creator user ID
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系 / Associations
	Project    *DesignerProject `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	ResultFile *CadFile         `gorm:"foreignKey:ResultFileID" json:"resultFile,omitempty"`
}

// TableName 表名
func (CadGeneration) TableName() string {
	return "cad_generations"
}

// CadGenerationType CAD生成类型常量
// CAD generation type constants
const (
	CadGenTypeFloorPlan = "floor_plan" // 平面布局图 / Floor plan
	CadGenTypeElevation = "elevation"  // 立面图 / Elevation
	CadGenTypeDetail    = "detail"     // 节点详图 / Detail drawing
	CadGenTypeCeiling   = "ceiling"    // 天花图 / Ceiling plan
	CadGenTypeElectric  = "electric"   // 电气图 / Electrical plan
)

// CadGenerationStatus CAD生成状态常量
// CAD generation status constants
const (
	CadGenStatusPending    = "pending"    // 待处理 / Pending
	CadGenStatusProcessing = "processing" // 处理中 / Processing
	CadGenStatusCompleted  = "completed"  // 已完成 / Completed
	CadGenStatusFailed     = "failed"     // 失败 / Failed
)

// CadGenerationParams CAD生成参数结构
// CAD generation parameters structure
type CadGenerationParams struct {
	Width        float64  `json:"width"`        // 宽度(米) / Width in meters
	Height       float64  `json:"height"`       // 高度(米) / Height in meters
	Scale        string   `json:"scale"`        // 比例 / Scale (1:50, 1:100, etc.)
	Style        string   `json:"style"`        // 设计风格 / Design style
	RoomTypes    []string `json:"roomTypes"`    // 房间类型列表 / Room types
	Requirements string   `json:"requirements"` // 特殊要求 / Special requirements
}
