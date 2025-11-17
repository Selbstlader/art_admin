package model

import (
	"time"

	"gorm.io/gorm"
)

// DictionaryType 字典类型模型
type DictionaryType struct {
	ID          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	TypeName    string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"typeName"` // 字典类型名称
	TypeCode    string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"typeCode"` // 字典类型编码
	Description string         `gorm:"type:varchar(200)" json:"description"`                  // 描述
	Enabled     bool           `gorm:"type:tinyint(1);default:1;index" json:"enabled"`        // 是否启用
	Remark      string         `gorm:"type:varchar(500)" json:"remark"`                       // 备注
	CreateBy    string         `gorm:"type:varchar(50)" json:"createBy"`                      // 创建者
	CreateTime  time.Time      `gorm:"autoCreateTime" json:"createTime"`                      // 创建时间
	UpdateBy    string         `gorm:"type:varchar(50)" json:"updateBy"`                      // 更新者
	UpdateTime  time.Time      `gorm:"autoUpdateTime" json:"updateTime"`                      // 更新时间
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`                                        // 软删除

	// 关联关系 - 移除外键约束，使用程序关联
	Dictionaries []Dictionary `gorm:"-" json:"dictionaries,omitempty"`
}

// TableName 表名
func (DictionaryType) TableName() string {
	return "sys_dictionary_type"
}

// Dictionary 字典数据模型
type Dictionary struct {
	ID         int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	TypeCode   string         `gorm:"type:varchar(50);not null;index" json:"typeCode"` // 字典类型编码
	Label      string         `gorm:"type:varchar(100);not null" json:"label"`         // 字典标签
	Value      string         `gorm:"type:varchar(100);not null" json:"value"`         // 字典值
	OrderNum   int            `gorm:"type:int;default:0" json:"orderNum"`              // 排序号
	Enabled    bool           `gorm:"type:tinyint(1);default:1;index" json:"enabled"`  // 是否启用
	Remark     string         `gorm:"type:varchar(500)" json:"remark"`                 // 备注
	CreateBy   string         `gorm:"type:varchar(50)" json:"createBy"`                // 创建者
	CreateTime time.Time      `gorm:"autoCreateTime" json:"createTime"`                // 创建时间
	UpdateBy   string         `gorm:"type:varchar(50)" json:"updateBy"`                // 更新者
	UpdateTime time.Time      `gorm:"autoUpdateTime" json:"updateTime"`                // 更新时间
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`                                  // 软删除

	// 关联关系 - 移除外键约束，使用程序关联
	DictionaryType DictionaryType `gorm:"-" json:"dictionaryType,omitempty"`
}

// TableName 表名
func (Dictionary) TableName() string {
	return "sys_dictionary"
}
