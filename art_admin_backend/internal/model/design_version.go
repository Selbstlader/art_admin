package model

import (
	"time"

	"gorm.io/gorm"
)

// DesignVersion 设计方案版本模型
// Design version model for storing different versions of design schemes
type DesignVersion struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	ProjectID     uint           `gorm:"index;not null" json:"projectId"`             // 关联项目ID / Associated project ID
	VersionNumber int            `gorm:"not null" json:"versionNumber"`               // 版本号 / Version number (1, 2, 3...)
	VersionName   string         `gorm:"size:100" json:"versionName"`                 // 版本名称 / Version name (e.g., "初稿", "修改版1")
	Description   string         `gorm:"type:text" json:"description"`                // 版本说明 / Version description
	DesignImages  *string        `gorm:"type:json" json:"designImages"`               // 设计图列表(JSON) / Design images list in JSON format
	CadFileIDs    *string        `gorm:"type:json" json:"cadFileIds"`                 // 关联CAD文件ID列表(JSON) / Associated CAD file IDs
	LayoutInfo    *string        `gorm:"type:json" json:"layoutInfo"`                 // 布局信息(JSON) / Layout information in JSON format
	AreaInfo      *string        `gorm:"type:json" json:"areaInfo"`                   // 面积信息(JSON) / Area information in JSON format
	StyleInfo     *string        `gorm:"type:json" json:"styleInfo"`                  // 风格信息(JSON) / Style information in JSON format
	MaterialInfo  *string        `gorm:"type:json" json:"materialInfo"`               // 材料信息(JSON) / Material information in JSON format
	Status        string         `gorm:"size:50;default:'draft';index" json:"status"` // 状态 / Status (draft/submitted/approved/rejected)
	CreatedBy     uint           `gorm:"index" json:"createdBy"`                      // 创建者ID / Creator user ID
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系 / Associations
	Project                 *DesignerProject         `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	CadFiles                []CadFile                `gorm:"-" json:"cadFiles,omitempty"`                                   // 关联的CAD文件(非数据库字段) / Associated CAD files
	ConstructionAnnotations []ConstructionAnnotation `gorm:"foreignKey:VersionID" json:"constructionAnnotations,omitempty"` // 施工图标注 / Construction annotations
}

// TableName 表名
func (DesignVersion) TableName() string {
	return "design_versions"
}

// DesignVersionCompare 版本对比结果模型
// Design version comparison result model
type DesignVersionCompare struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	ProjectID       uint           `gorm:"index;not null" json:"projectId"`                // 关联项目ID / Associated project ID
	VersionAID      uint           `gorm:"index;not null" json:"versionAId"`               // 版本A ID / Version A ID
	VersionBID      uint           `gorm:"index;not null" json:"versionBId"`               // 版本B ID / Version B ID
	LayoutChanges   string         `gorm:"type:json" json:"layoutChanges"`                 // 布局变化(JSON) / Layout changes in JSON format
	AreaChanges     string         `gorm:"type:json" json:"areaChanges"`                   // 面积调整(JSON) / Area changes in JSON format
	ElementChanges  string         `gorm:"type:json" json:"elementChanges"`                // 元素增减(JSON) / Element changes in JSON format
	StyleChanges    string         `gorm:"type:json" json:"styleChanges"`                  // 风格变化(JSON) / Style changes in JSON format
	MaterialChanges string         `gorm:"type:json" json:"materialChanges"`               // 材料变化(JSON) / Material changes in JSON format
	Summary         string         `gorm:"type:text" json:"summary"`                       // 对比摘要 / Comparison summary
	CompareStatus   string         `gorm:"size:50;default:'pending'" json:"compareStatus"` // 对比状态 / Compare status
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系 / Associations
	Project  *DesignerProject `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	VersionA *DesignVersion   `gorm:"foreignKey:VersionAID" json:"versionA,omitempty"`
	VersionB *DesignVersion   `gorm:"foreignKey:VersionBID" json:"versionB,omitempty"`
}

// TableName 表名
func (DesignVersionCompare) TableName() string {
	return "design_version_compares"
}

// DesignImage 设计图结构
// Design image structure
type DesignImage struct {
	URL         string `json:"url"`         // 图片URL / Image URL
	Name        string `json:"name"`        // 图片名称 / Image name
	Description string `json:"description"` // 图片描述 / Image description
	Type        string `json:"type"`        // 图片类型 / Image type (floor_plan/elevation/3d_render)
}

// LayoutInfoItem 布局信息项
// Layout information item
type LayoutInfoItem struct {
	Zone     string  `json:"zone"`     // 区域名称 / Zone name
	Area     float64 `json:"area"`     // 面积 / Area
	Position string  `json:"position"` // 位置描述 / Position description
}

// AreaInfoItem 面积信息项
// Area information item
type AreaInfoItem struct {
	Name  string  `json:"name"`  // 区域名称 / Area name
	Value float64 `json:"value"` // 面积值 / Area value
	Unit  string  `json:"unit"`  // 单位 / Unit (sqm)
}

// StyleInfoItem 风格信息项
// Style information item
type StyleInfoItem struct {
	Category string `json:"category"` // 风格类别 / Style category
	Value    string `json:"value"`    // 风格值 / Style value
	Details  string `json:"details"`  // 详细描述 / Details
}

// MaterialInfoItem 材料信息项
// Material information item
type MaterialInfoItem struct {
	Name     string  `json:"name"`     // 材料名称 / Material name
	Category string  `json:"category"` // 材料类别 / Material category
	Quantity float64 `json:"quantity"` // 数量 / Quantity
	Unit     string  `json:"unit"`     // 单位 / Unit
}

// ChangeItem 变化项结构
// Change item structure for comparison results
type ChangeItem struct {
	Field       string `json:"field"`       // 变化字段 / Changed field
	OldValue    string `json:"oldValue"`    // 旧值 / Old value
	NewValue    string `json:"newValue"`    // 新值 / New value
	ChangeType  string `json:"changeType"`  // 变化类型 / Change type (added/removed/modified)
	Description string `json:"description"` // 变化描述 / Change description
}
