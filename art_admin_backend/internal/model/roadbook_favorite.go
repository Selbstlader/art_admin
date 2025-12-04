package model

import (
	"encoding/json"
	"time"
)

// RoadbookFavorite 收藏模型
type RoadbookFavorite struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     int64     `gorm:"type:bigint unsigned;not null;uniqueIndex:uk_user_roadbook" json:"userId"`
	RoadbookID int64     `gorm:"type:bigint unsigned;not null;uniqueIndex:uk_user_roadbook;index" json:"roadbookId"`
	CreatedAt  time.Time `gorm:"type:timestamp;not null" json:"createdAt"`

	// 关联
	Roadbook *Roadbook `gorm:"foreignKey:RoadbookID" json:"roadbook,omitempty"`
}

// TableName 表名
func (RoadbookFavorite) TableName() string {
	return "roadbook_favorites"
}

// ToJSON 序列化为JSON
func (f *RoadbookFavorite) ToJSON() ([]byte, error) {
	return json.Marshal(f)
}

// FromJSON 从JSON反序列化
func (f *RoadbookFavorite) FromJSON(data []byte) error {
	return json.Unmarshal(data, f)
}
