package model

import (
	"time"
)

// JournalEntry 日记条目模型
type JournalEntry struct {
	ID             int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         int64     `gorm:"not null;index:idx_user_id" json:"user_id"`
	Title          string    `gorm:"type:varchar(255)" json:"title"`
	Content        string    `gorm:"type:text;not null" json:"content"`
	MoodRecordID   *int64    `gorm:"index:idx_mood_record_id" json:"mood_record_id"`                                      // 关联的情绪记录ID
	SentimentScore *float64  `gorm:"type:decimal(3,2)" json:"sentiment_score"`                                            // 情感分析得分 -1.0 到 1.0
	IsPrivate      bool      `gorm:"default:true" json:"is_private"`                                                      // 是否私密
	Tags           string    `gorm:"type:json" json:"tags"`                                                               // 标签 JSON数组
	Images         string    `gorm:"type:json" json:"images"`                                                             // 图片URLs JSON数组
	AnalysisStatus string    `gorm:"type:varchar(20);default:'pending';index:idx_analysis_status" json:"analysis_status"` // 分析状态: pending, analyzing, completed, failed
	AnalysisResult string    `gorm:"type:json" json:"analysis_result"`                                                    // AI分析结果 JSON
	DifyDocumentID string    `gorm:"type:varchar(255);index:idx_dify_document_id" json:"dify_document_id"`                // Dify知识库文档ID
	CreatedAt      time.Time `gorm:"index:idx_created_at" json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	// 关联
	MoodRecord *MoodRecord `gorm:"foreignKey:MoodRecordID" json:"mood_record,omitempty"`
}

func (JournalEntry) TableName() string {
	return "journal_entries"
}

// AchievementType 成就类型枚举
type AchievementType string

const (
	AchievementMoodStreak       AchievementType = "mood_streak"       // 连续记录情绪
	AchievementMeditationStreak AchievementType = "meditation_streak" // 连续冥想
	AchievementMeditationTotal  AchievementType = "meditation_total"  // 总冥想时长
	AchievementJournalStreak    AchievementType = "journal_streak"    // 连续写日记
	AchievementJournalCount     AchievementType = "journal_count"     // 日记总数
	AchievementFirstMood        AchievementType = "first_mood"        // 第一次记录情绪
	AchievementFirstMeditation  AchievementType = "first_meditation"  // 第一次冥想
	AchievementFirstJournal     AchievementType = "first_journal"     // 第一篇日记
)

// GetAchievementTypes 获取所有成就类型
func GetAchievementTypes() []AchievementType {
	return []AchievementType{
		AchievementMoodStreak, AchievementMeditationStreak, AchievementMeditationTotal,
		AchievementJournalStreak, AchievementJournalCount, AchievementFirstMood,
		AchievementFirstMeditation, AchievementFirstJournal,
	}
}

// GetAchievementConfig 获取成就配置
func GetAchievementConfig(achievementType AchievementType) map[string]interface{} {
	configs := map[AchievementType]map[string]interface{}{
		AchievementMoodStreak: {
			"name":        "情绪记录达人",
			"description": "连续7天记录情绪",
			"target":      7,
			"icon":        "mood",
		},
		AchievementMeditationStreak: {
			"name":        "冥想坚持者",
			"description": "连续30天冥想",
			"target":      30,
			"icon":        "meditation",
		},
		AchievementMeditationTotal: {
			"name":        "冥想大师",
			"description": "累计冥想1000分钟",
			"target":      1000,
			"icon":        "clock",
		},
		AchievementJournalStreak: {
			"name":        "日记达人",
			"description": "连续30天写日记",
			"target":      30,
			"icon":        "journal",
		},
		AchievementJournalCount: {
			"name":        "写作爱好者",
			"description": "写下100篇日记",
			"target":      100,
			"icon":        "write",
		},
		AchievementFirstMood: {
			"name":        "情绪初体验",
			"description": "第一次记录情绪",
			"target":      1,
			"icon":        "star",
		},
		AchievementFirstMeditation: {
			"name":        "冥想初体验",
			"description": "完成第一次冥想",
			"target":      1,
			"icon":        "star",
		},
		AchievementFirstJournal: {
			"name":        "日记初体验",
			"description": "写下第一篇日记",
			"target":      1,
			"icon":        "star",
		},
	}
	return configs[achievementType]
}

// MoodAnalytics 情绪分析聚合模型
type MoodAnalytics struct {
	ID               int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID           int64     `gorm:"not null;index" json:"user_id"`
	AnalyzeType      string    `gorm:"type:varchar(20);not null" json:"analyze_type"` // daily, weekly, monthly
	AnalyzeDate      time.Time `gorm:"not null;index" json:"analyze_date"`            // 分析日期
	MoodDistribution string    `gorm:"type:json" json:"mood_distribution"`            // 情绪分布 JSON
	AverageIntensity float64   `gorm:"type:decimal(3,1)" json:"average_intensity"`    // 平均强度
	TotalRecords     int       `gorm:"default:0" json:"total_records"`                // 总记录数
	MostCommonMood   string    `gorm:"type:varchar(50)" json:"most_common_mood"`      // 最常见情绪
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (MoodAnalytics) TableName() string {
	return "mood_analytics"
}
