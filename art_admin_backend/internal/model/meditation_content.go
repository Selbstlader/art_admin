package model

import (
	"time"
)

// MeditationContent 冥想内容模型
type MeditationContent struct {
	ID              int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Title           string    `gorm:"type:varchar(255);not null" json:"title"`
	Category        string    `gorm:"type:varchar(100);not null" json:"category"`
	Duration        int       `gorm:"not null" json:"duration"`                                                                // 时长(秒)
	AudioURL        string    `gorm:"type:varchar(500)" json:"audio_url"`                                                      // 音频文件URL
	CoverImage      string    `gorm:"type:varchar(500)" json:"cover_image"`                                                    // 封面图片URL
	Description     string    `gorm:"type:text" json:"description"`                                                            // 描述
	DifficultyLevel int       `gorm:"default:1;check:difficulty_level >= 1 AND difficulty_level <= 5" json:"difficulty_level"` // 难度等级 1-5
	Tags            string    `gorm:"type:json" json:"tags"`                                                                   // 标签 JSON数组
	IsPublic        bool      `gorm:"default:true" json:"is_public"`                                                           // 是否公开
	ViewCount       int       `gorm:"default:0" json:"view_count"`                                                             // 查看次数
	LikeCount       int       `gorm:"default:0" json:"like_count"`                                                             // 点赞数
	CreatedBy       int64     `json:"created_by"`                                                                              // 创建者ID
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (MeditationContent) TableName() string {
	return "meditation_contents"
}

// MeditationCategory 冥想分类枚举
type MeditationCategory string

const (
	CategorySleep       MeditationCategory = "sleep"       // 睡眠
	CategoryStress      MeditationCategory = "stress"      // 减压
	CategoryFocus       MeditationCategory = "focus"       // 专注
	CategoryAnxiety     MeditationCategory = "anxiety"     // 焦虑
	CategoryBreathing   MeditationCategory = "breathing"   // 呼吸
	CategoryMindfulness MeditationCategory = "mindfulness" // 正念
	CategoryBodyScan    MeditationCategory = "body_scan"   // 身体扫描
	CategoryWalking     MeditationCategory = "walking"     // 行走冥想
)

// GetMeditationCategories 获取所有冥想分类
func GetMeditationCategories() []MeditationCategory {
	return []MeditationCategory{
		CategorySleep, CategoryStress, CategoryFocus,
		CategoryAnxiety, CategoryBreathing, CategoryMindfulness,
		CategoryBodyScan, CategoryWalking,
	}
}

// GetCategoryLabel 获取分类标签
func GetCategoryLabel(category MeditationCategory) string {
	labels := map[MeditationCategory]string{
		CategorySleep:       "睡眠冥想",
		CategoryStress:      "减压放松",
		CategoryFocus:       "专注训练",
		CategoryAnxiety:     "焦虑缓解",
		CategoryBreathing:   "呼吸练习",
		CategoryMindfulness: "正念冥想",
		CategoryBodyScan:    "身体扫描",
		CategoryWalking:     "行走冥想",
	}
	return labels[category]
}

// UserMeditationSession 用户冥想会话模型
type UserMeditationSession struct {
	ID                int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID            int64      `gorm:"not null;index" json:"user_id"`
	ContentID         int64      `gorm:"not null" json:"content_id"`
	CompletedDuration int        `gorm:"not null" json:"completed_duration"` // 实际完成时长(秒)
	IsCompleted       bool       `gorm:"default:false" json:"is_completed"`  // 是否完成
	StartedAt         *time.Time `json:"started_at"`
	CompletedAt       *time.Time `json:"completed_at"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`

	// 关联
	MeditationContent MeditationContent `gorm:"foreignKey:ContentID" json:"meditation_content,omitempty"`
}

func (UserMeditationSession) TableName() string {
	return "user_meditation_sessions"
}
