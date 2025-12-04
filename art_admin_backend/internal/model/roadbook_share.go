package model

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"time"
)

// RoadbookShare 分享记录模型
type RoadbookShare struct {
	ID         int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	RoadbookID int64      `gorm:"type:bigint unsigned;not null;index" json:"roadbookId"`
	UserID     int64      `gorm:"type:bigint unsigned;not null;index" json:"userId"`
	ShareCode  string     `gorm:"type:varchar(32);not null;uniqueIndex" json:"shareCode"`
	ShareType  int        `gorm:"type:tinyint;default:1" json:"shareType"`
	VisitCount int        `gorm:"type:int unsigned;default:0" json:"visitCount"`
	ExpiresAt  *time.Time `gorm:"type:timestamp" json:"expiresAt"`
	CreatedAt  time.Time  `gorm:"type:timestamp;not null" json:"createdAt"`

	// 关联
	Roadbook *Roadbook `gorm:"foreignKey:RoadbookID" json:"roadbook,omitempty"`
}

// TableName 表名
func (RoadbookShare) TableName() string {
	return "roadbook_shares"
}

// ToJSON 序列化为JSON
func (s *RoadbookShare) ToJSON() ([]byte, error) {
	return json.Marshal(s)
}

// FromJSON 从JSON反序列化
func (s *RoadbookShare) FromJSON(data []byte) error {
	return json.Unmarshal(data, s)
}

// GenerateShareCode 生成唯一分享码
func GenerateShareCode() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// RoadbookTemplate 路书模板模型
type RoadbookTemplate struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	RoadbookID  int64     `gorm:"type:bigint unsigned;not null;index" json:"roadbookId"`
	UserID      int64     `gorm:"type:bigint unsigned;not null;index" json:"userId"`
	Title       string    `gorm:"type:varchar(100);not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	CoverURL    string    `gorm:"type:varchar(500)" json:"coverUrl"`
	Destination string    `gorm:"type:varchar(100)" json:"destination"`
	Days        int       `gorm:"type:int;not null" json:"days"`
	Category    int       `gorm:"type:tinyint;default:1" json:"category"`
	UseCount    int       `gorm:"type:int unsigned;default:0" json:"useCount"`
	Status      int       `gorm:"type:tinyint;default:0" json:"status"`
	CreatedAt   time.Time `gorm:"type:timestamp;not null" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"type:timestamp;not null" json:"updatedAt"`

	// 关联
	Roadbook *Roadbook `gorm:"foreignKey:RoadbookID" json:"roadbook,omitempty"`
}

// TableName 表名
func (RoadbookTemplate) TableName() string {
	return "roadbook_templates"
}

// RoadbookStatsDaily 每日统计模型
type RoadbookStatsDaily struct {
	ID            int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	RoadbookID    int64     `gorm:"type:bigint unsigned;not null;uniqueIndex:uk_roadbook_date" json:"roadbookId"`
	StatDate      time.Time `gorm:"type:date;not null;uniqueIndex:uk_roadbook_date" json:"statDate"`
	ViewCount     int       `gorm:"type:int unsigned;default:0" json:"viewCount"`
	FavoriteCount int       `gorm:"type:int unsigned;default:0" json:"favoriteCount"`
	ShareCount    int       `gorm:"type:int unsigned;default:0" json:"shareCount"`
}

// TableName 表名
func (RoadbookStatsDaily) TableName() string {
	return "roadbook_stats_daily"
}
