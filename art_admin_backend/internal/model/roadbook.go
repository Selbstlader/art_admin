package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// 可见性常量
const (
	VisibilityPublic   = 1 // 公开
	VisibilityPrivate  = 2 // 私有
	VisibilitySelected = 3 // 指定用户
)

// 路书状态常量
const (
	RoadbookStatusDraft     = 1 // 草稿
	RoadbookStatusPublished = 2 // 已发布
	RoadbookStatusOffline   = 3 // 已下架
)

// 途经点类型常量
const (
	WaypointTypeStart    = 1 // 起点
	WaypointTypeWaypoint = 2 // 途经点
	WaypointTypeEnd      = 3 // 终点
)

// Roadbook 路书模型
type Roadbook struct {
	ID            int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        int64          `gorm:"type:bigint unsigned;not null;index" json:"userId"`
	Title         string         `gorm:"type:varchar(100);not null" json:"title"`
	Description   string         `gorm:"type:text" json:"description"`
	CoverURL      string         `gorm:"type:varchar(500)" json:"coverUrl"`
	StartDate     *time.Time     `gorm:"type:date" json:"startDate"`
	EndDate       *time.Time     `gorm:"type:date" json:"endDate"`
	Visibility    int            `gorm:"type:tinyint;default:1" json:"visibility"`
	TravelMode    string         `gorm:"type:varchar(20);default:driving" json:"travelMode"`
	TotalDistance int            `gorm:"type:int;default:0" json:"totalDistance"`
	TotalDuration int            `gorm:"type:int;default:0" json:"totalDuration"`
	TotalBudget   float64        `gorm:"type:decimal(10,2);default:0" json:"totalBudget"`
	ViewCount     int            `gorm:"type:int unsigned;default:0" json:"viewCount"`
	FavoriteCount int            `gorm:"type:int unsigned;default:0" json:"favoriteCount"`
	LikeCount     int            `gorm:"type:int unsigned;default:0" json:"likeCount"`
	CommentCount  int            `gorm:"type:int unsigned;default:0" json:"commentCount"`
	ShareCount    int            `gorm:"type:int unsigned;default:0" json:"shareCount"`
	Status        int            `gorm:"type:tinyint;default:1" json:"status"`
	CreatedAt     time.Time      `gorm:"type:timestamp;not null" json:"createdAt"`
	UpdatedAt     time.Time      `gorm:"type:timestamp;not null" json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Waypoints []RoadbookWaypoint `gorm:"foreignKey:RoadbookID" json:"waypoints,omitempty"`
	Tags      []RoadbookTag      `gorm:"many2many:roadbook_tag_relations" json:"tags,omitempty"`
}

// TableName 表名
func (Roadbook) TableName() string {
	return "roadbooks"
}

// ToJSON 序列化为JSON
func (r *Roadbook) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}

// FromJSON 从JSON反序列化
func (r *Roadbook) FromJSON(data []byte) error {
	return json.Unmarshal(data, r)
}

// StringSlice 自定义类型用于JSON数组
type StringSlice []string

// Value 实现 driver.Valuer 接口
func (s StringSlice) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	return json.Marshal(s)
}

// Scan 实现 sql.Scanner 接口
func (s *StringSlice) Scan(value interface{}) error {
	if value == nil {
		*s = []string{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, s)
}

// RoadbookWaypoint 途经点模型
type RoadbookWaypoint struct {
	ID           int64       `gorm:"primaryKey;autoIncrement" json:"id"`
	RoadbookID   int64       `gorm:"type:bigint unsigned;not null;index" json:"roadbookId"`
	Name         string      `gorm:"type:varchar(100);not null" json:"name"`
	Address      string      `gorm:"type:varchar(300)" json:"address"`
	Longitude    float64     `gorm:"type:decimal(10,7);not null" json:"longitude"`
	Latitude     float64     `gorm:"type:decimal(10,7);not null" json:"latitude"`
	PoiID        string      `gorm:"type:varchar(50)" json:"poiId"`
	PoiType      string      `gorm:"type:varchar(50)" json:"poiType"`
	DayIndex     int         `gorm:"type:int;default:1" json:"dayIndex"`
	SortOrder    int         `gorm:"type:int;default:0" json:"sortOrder"`
	StayDuration int         `gorm:"type:int;default:60" json:"stayDuration"`
	Budget       float64     `gorm:"type:decimal(10,2);default:0" json:"budget"`
	Notes        string      `gorm:"type:text" json:"notes"`
	Images       StringSlice `gorm:"type:json" json:"images"`
	WaypointType int         `gorm:"type:tinyint;default:2" json:"waypointType"`
	CreatedAt    time.Time   `gorm:"type:timestamp;not null" json:"createdAt"`
	UpdatedAt    time.Time   `gorm:"type:timestamp;not null" json:"updatedAt"`
}

// TableName 表名
func (RoadbookWaypoint) TableName() string {
	return "roadbook_waypoints"
}

// ToJSON 序列化为JSON
func (w *RoadbookWaypoint) ToJSON() ([]byte, error) {
	return json.Marshal(w)
}

// FromJSON 从JSON反序列化
func (w *RoadbookWaypoint) FromJSON(data []byte) error {
	return json.Unmarshal(data, w)
}

// RoadbookTag 标签模型
type RoadbookTag struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(50);not null;uniqueIndex" json:"name"`
	Color     string    `gorm:"type:varchar(20);default:#409EFF" json:"color"`
	UseCount  int       `gorm:"type:int unsigned;default:0" json:"useCount"`
	IsSystem  int       `gorm:"type:tinyint;default:0" json:"isSystem"`
	CreatedAt time.Time `gorm:"type:timestamp;not null" json:"createdAt"`
}

// TableName 表名
func (RoadbookTag) TableName() string {
	return "roadbook_tags"
}

// RoadbookTagRelation 标签关联模型
type RoadbookTagRelation struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	RoadbookID int64     `gorm:"type:bigint unsigned;not null;uniqueIndex:uk_roadbook_tag" json:"roadbookId"`
	TagID      int64     `gorm:"type:bigint unsigned;not null;uniqueIndex:uk_roadbook_tag" json:"tagId"`
	CreatedAt  time.Time `gorm:"type:timestamp;not null" json:"createdAt"`
}

// TableName 表名
func (RoadbookTagRelation) TableName() string {
	return "roadbook_tag_relations"
}
