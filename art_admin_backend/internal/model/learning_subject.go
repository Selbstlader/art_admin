package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// Subject 学科表
type Subject struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"type:varchar(50);not null;comment:学科名称" json:"name"`
	Icon        string    `gorm:"type:varchar(200);comment:学科图标" json:"icon"`
	Description string    `gorm:"type:text;comment:学科描述" json:"description"`
	GradeLevels JSONArray `gorm:"type:json;comment:适用年级" json:"grade_levels"`
	Sort        int       `gorm:"default:0;comment:排序" json:"sort"`
	Status      int8      `gorm:"default:1;comment:状态 1-启用 0-禁用" json:"status"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (Subject) TableName() string {
	return "subjects"
}

// JSONArray 自定义 JSON 数组类型
type JSONArray []string

// Scan 实现 sql.Scanner 接口
func (j *JSONArray) Scan(value interface{}) error {
	if value == nil {
		*j = []string{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

// Value 实现 driver.Valuer 接口
func (j JSONArray) Value() (driver.Value, error) {
	if len(j) == 0 {
		return "[]", nil
	}
	return json.Marshal(j)
}
