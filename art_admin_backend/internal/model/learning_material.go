package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// LearningMaterial 教材内容表
type LearningMaterial struct {
	ID            int64           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        int64           `gorm:"not null;index;comment:创建用户ID" json:"user_id"`
	SubjectID     int64           `gorm:"not null;index;comment:学科ID" json:"subject_id"`
	Title         string          `gorm:"type:varchar(200);not null;comment:教材标题" json:"title"`
	Topic         string          `gorm:"type:varchar(200);not null;comment:学习题材" json:"topic"`
	Grade         string          `gorm:"type:varchar(50);not null;comment:年级" json:"grade"`
	Difficulty    int8            `gorm:"default:1;comment:难度 1-基础 2-进阶 3-高级" json:"difficulty"`
	Summary       string          `gorm:"type:text;comment:内容概要" json:"summary"`
	Content       MaterialContent `gorm:"type:longtext;not null;comment:教材内容(JSON格式)" json:"content"`
	AudioURL      string          `gorm:"type:varchar(500);comment:音频URL" json:"audio_url"`
	AudioDuration int             `gorm:"comment:音频时长(秒)" json:"audio_duration"`
	TotalTime     int             `gorm:"comment:预计学习时长(分钟)" json:"total_time"`
	ViewCount     int             `gorm:"default:0;comment:浏览次数" json:"view_count"`
	FavoriteCount int             `gorm:"default:0;comment:收藏次数" json:"favorite_count"`
	Status        int8            `gorm:"default:1;index;comment:状态 1-正常 0-已删除" json:"status"`
	CreatedAt     time.Time       `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt     time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (LearningMaterial) TableName() string {
	return "learning_materials"
}

// MaterialContent 教材内容结构
type MaterialContent struct {
	Title     string            `json:"title"`
	Summary   string            `json:"summary"`
	Sections  []MaterialSection `json:"sections"`
	TotalTime int               `json:"total_time"`
}

// MaterialSection 教材章节
type MaterialSection struct {
	Type       string     `json:"type"` // knowledge, example, exercise
	Title      string     `json:"title"`
	Content    string     `json:"content,omitempty"`
	KeyPoints  []string   `json:"key_points,omitempty"`
	Difficulty string     `json:"difficulty,omitempty"`
	Question   string     `json:"question,omitempty"`
	Solution   string     `json:"solution,omitempty"`
	Answer     string     `json:"answer,omitempty"`
	Questions  []Exercise `json:"questions,omitempty"`
}

// Exercise 练习题
type Exercise struct {
	ID          int      `json:"id"`
	Type        string   `json:"type"` // choice, fill, calculate, essay
	Question    string   `json:"question"`
	Options     []string `json:"options,omitempty"`
	Answer      string   `json:"answer"`
	Explanation string   `json:"explanation"`
}

// Scan 实现 sql.Scanner 接口
func (mc *MaterialContent) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, mc)
}

// Value 实现 driver.Valuer 接口
func (mc MaterialContent) Value() (driver.Value, error) {
	return json.Marshal(mc)
}

// MaterialSection 教材章节表
type MaterialSectionModel struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	MaterialID int64     `gorm:"not null;index;comment:教材ID" json:"material_id"`
	Type       string    `gorm:"type:varchar(20);not null;index;comment:类型 knowledge|example|exercise" json:"type"`
	Title      string    `gorm:"type:varchar(200);comment:章节标题" json:"title"`
	Content    string    `gorm:"type:longtext;comment:章节内容(JSON格式)" json:"content"`
	Sort       int       `gorm:"default:0;index;comment:排序" json:"sort"`
	AudioURL   string    `gorm:"type:varchar(500);comment:章节音频URL" json:"audio_url"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TableName 指定表名
func (MaterialSectionModel) TableName() string {
	return "material_sections"
}
