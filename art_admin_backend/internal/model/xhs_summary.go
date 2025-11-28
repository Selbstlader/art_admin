package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// XHSSummary 小红书笔记总结记录
type XHSSummary struct {
	ID             int64         `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         int64         `gorm:"index;not null" json:"user_id"`        // 用户ID
	OriginalURL    string        `gorm:"size:500" json:"original_url"`         // 原始链接（可选）
	NoteTitle      string        `gorm:"size:200" json:"note_title"`           // 笔记标题
	NoteContent    string        `gorm:"type:text" json:"note_content"`        // 笔记原文内容
	ImageURLs      StringArray   `gorm:"type:json" json:"image_urls"`          // 图片URL列表
	SummaryTitle   string        `gorm:"size:100" json:"summary_title"`        // 总结标题
	SummaryContent string        `gorm:"type:text" json:"summary_content"`     // 总结内容
	KeyPoints      StringArray   `gorm:"type:json" json:"key_points"`          // 关键信息
	Tags           StringArray   `gorm:"type:json" json:"tags"`                // 标签
	Sentiment      string        `gorm:"size:20" json:"sentiment"`             // 情感倾向
	OCRTexts       StringArray   `gorm:"type:json" json:"ocr_texts"`           // OCR识别文字
	Analysis       *AnalysisJSON `gorm:"type:json" json:"analysis"`            // 专业分析
	Style          string        `gorm:"size:20;default:concise" json:"style"` // 总结风格
	IsFavorite     bool          `gorm:"default:false" json:"is_favorite"`     // 是否收藏
	FolderID       int64         `gorm:"default:0" json:"folder_id"`           // 收藏夹ID
	CreatedAt      time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 表名
func (XHSSummary) TableName() string {
	return "xhs_summaries"
}

// StringArray 字符串数组类型（JSON存储）
type StringArray []string

func (s StringArray) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	return json.Marshal(s)
}

func (s *StringArray) Scan(value interface{}) error {
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

// AnalysisJSON 专业分析结果JSON类型
type AnalysisJSON struct {
	Category        string            `json:"category"`
	MainConclusion  string            `json:"main_conclusion"`
	DetailedData    map[string]string `json:"detailed_data"`
	RiskWarnings    []string          `json:"risk_warnings"`
	Recommendations []string          `json:"recommendations"`
}

func (a AnalysisJSON) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *AnalysisJSON) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, a)
}

// XHSFolder 小红书收藏夹
type XHSFolder struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"index;not null" json:"user_id"`
	Name      string    `gorm:"size:50;not null" json:"name"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 表名
func (XHSFolder) TableName() string {
	return "xhs_folders"
}
