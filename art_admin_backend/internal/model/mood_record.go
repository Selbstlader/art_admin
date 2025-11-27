package model

import (
	"time"
)

// MoodRecord 情绪记录模型
type MoodRecord struct {
	ID             int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         int64     `gorm:"not null;index:idx_user_time" json:"user_id"`
	MoodType       string    `gorm:"type:varchar(50);not null" json:"mood_type"`                         // 情绪类型：happy, sad, anxious, calm, angry, etc.
	Intensity      int       `gorm:"not null;check:intensity >= 1 AND intensity <= 10" json:"intensity"` // 强度 1-10
	Triggers       string    `gorm:"type:json" json:"triggers"`                                          // 触发因素 JSON数组
	Activities     string    `gorm:"type:json" json:"activities"`                                        // 相关活动 JSON数组
	Note           string    `gorm:"type:text" json:"note"`                                              // 备注
	Location       string    `gorm:"type:varchar(255)" json:"location"`                                  // 地理位置
	AnalysisStatus string    `gorm:"type:varchar(20);default:'pending'" json:"analysis_status"`          // 分析状态: pending, analyzing, completed, failed
	AnalysisResult string    `gorm:"type:json" json:"analysis_result"`                                   // AI分析结果 JSON
	DifyDocumentID string    `gorm:"type:varchar(255)" json:"dify_document_id"`                          // Dify知识库文档ID
	CreatedAt      time.Time `gorm:"index:idx_user_time" json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (MoodRecord) TableName() string {
	return "mood_records"
}

// MoodType 情绪类型枚举
type MoodType string

const (
	MoodHappy    MoodType = "happy"
	MoodSad      MoodType = "sad"
	MoodAnxious  MoodType = "anxious"
	MoodCalm     MoodType = "calm"
	MoodAngry    MoodType = "angry"
	MoodExcited  MoodType = "excited"
	MoodTired    MoodType = "tired"
	MoodStressed MoodType = "stressed"
)

// GetMoodTypes 获取所有情绪类型
func GetMoodTypes() []MoodType {
	return []MoodType{
		MoodHappy, MoodSad, MoodAnxious, MoodCalm,
		MoodAngry, MoodExcited, MoodTired, MoodStressed,
	}
}

// GetMoodTypeLabel 获取情绪类型标签
func GetMoodTypeLabel(moodType MoodType) string {
	labels := map[MoodType]string{
		MoodHappy:    "开心",
		MoodSad:      "难过",
		MoodAnxious:  "焦虑",
		MoodCalm:     "平静",
		MoodAngry:    "愤怒",
		MoodExcited:  "兴奋",
		MoodTired:    "疲惫",
		MoodStressed: "压力",
	}
	return labels[moodType]
}
