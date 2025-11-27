package model

import (
	"time"
)

// MeditationRecord 冥想记录模型
type MeditationRecord struct {
	ID                int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID            int64     `gorm:"not null;index:idx_user_time" json:"user_id"`
	ContentID         *int64    `gorm:"index" json:"content_id"`                                                                     // 关联的冥想内容ID
	Duration          int       `gorm:"not null" json:"duration"`                                                                    // 冥想时长（分钟）
	MeditationType    string    `gorm:"type:varchar(50);not null" json:"meditation_type"`                                            // 冥想类型
	DifficultyLevel   int       `gorm:"default:1" json:"difficulty_level"`                                                           // 难度等级 1-5
	PreMoodType       string    `gorm:"type:varchar(50)" json:"pre_mood_type"`                                                       // 冥想前情绪
	PostMoodType      string    `gorm:"type:varchar(50)" json:"post_mood_type"`                                                      // 冥想后情绪
	PreMoodIntensity  int       `gorm:"check:pre_mood_intensity >= 1 AND pre_mood_intensity <= 10" json:"pre_mood_intensity"`        // 冥想前情绪强度
	PostMoodIntensity int       `gorm:"check:post_mood_intensity >= 1 AND post_mood_intensity <= 10" json:"post_mood_intensity"`     // 冥想后情绪强度
	Experience        string    `gorm:"type:text" json:"experience"`                                                                 // 冥想体验描述
	DistractionLevel  int       `gorm:"default:0;check:distraction_level >= 0 AND distraction_level <= 10" json:"distraction_level"` // 分心程度 0-10
	Techniques        string    `gorm:"type:json" json:"techniques"`                                                                 // 使用的技巧 JSON数组
	Notes             string    `gorm:"type:text" json:"notes"`                                                                      // 备注
	AnalysisStatus    string    `gorm:"type:varchar(20);default:'pending'" json:"analysis_status"`                                   // 分析状态: pending, analyzing, completed, failed
	AnalysisResult    string    `gorm:"type:json" json:"analysis_result"`                                                            // AI分析结果 JSON
	DifyDocumentID    string    `gorm:"type:varchar(255)" json:"dify_document_id"`                                                   // Dify知识库文档ID
	MoodImprovement   float64   `gorm:"type:decimal(5,2);default:0.0" json:"mood_improvement"`                                       // 情绪改善程度
	CreatedAt         time.Time `gorm:"index:idx_user_time" json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`

	// 关联
	Content *MeditationContent `gorm:"foreignKey:ContentID" json:"content,omitempty"`
}

func (MeditationRecord) TableName() string {
	return "meditation_records"
}

// MeditationType 冥想类型枚举
type MeditationType string

const (
	MeditationBreathing      MeditationType = "breathing"       // 呼吸冥想
	MeditationBodyScan       MeditationType = "body_scan"       // 身体扫描
	MeditationMindfulness    MeditationType = "mindfulness"     // 正念冥想
	MeditationLovingKindness MeditationType = "loving_kindness" // 慈心冥想
	MeditationWalking        MeditationType = "walking"         // 行走冥想
	MeditationVisualization  MeditationType = "visualization"   // 想象冥想
	MeditationMantra         MeditationType = "mantra"          // 真言冥想
	MeditationYoga           MeditationType = "yoga"            // 瑜伽冥想
)

// GetMeditationTypes 获取所有冥想类型
func GetMeditationTypes() []MeditationType {
	return []MeditationType{
		MeditationBreathing, MeditationBodyScan, MeditationMindfulness,
		MeditationLovingKindness, MeditationWalking, MeditationVisualization,
		MeditationMantra, MeditationYoga,
	}
}

// GetMeditationTypeLabel 获取冥想类型标签
func GetMeditationTypeLabel(meditationType MeditationType) string {
	labels := map[MeditationType]string{
		MeditationBreathing:      "呼吸冥想",
		MeditationBodyScan:       "身体扫描",
		MeditationMindfulness:    "正念冥想",
		MeditationLovingKindness: "慈心冥想",
		MeditationWalking:        "行走冥想",
		MeditationVisualization:  "想象冥想",
		MeditationMantra:         "真言冥想",
		MeditationYoga:           "瑜伽冥想",
	}
	return labels[meditationType]
}

// AnalysisStatus 分析状态枚举
type AnalysisStatus string

const (
	AnalysisStatusPending   AnalysisStatus = "pending"   // 等待分析
	AnalysisStatusAnalyzing AnalysisStatus = "analyzing" // 分析中
	AnalysisStatusCompleted AnalysisStatus = "completed" // 分析完成
	AnalysisStatusFailed    AnalysisStatus = "failed"    // 分析失败
)

// GetAnalysisStatuses 获取所有分析状态
func GetAnalysisStatuses() []AnalysisStatus {
	return []AnalysisStatus{
		AnalysisStatusPending, AnalysisStatusAnalyzing,
		AnalysisStatusCompleted, AnalysisStatusFailed,
	}
}

// MeditationAnalytics 冥想分析聚合模型
type MeditationAnalytics struct {
	ID                 int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID             int64     `gorm:"not null;index" json:"user_id"`
	AnalyzeType        string    `gorm:"type:varchar(20);not null" json:"analyze_type"` // daily, weekly, monthly
	AnalyzeDate        time.Time `gorm:"not null;index" json:"analyze_date"`            // 分析日期
	TotalSessions      int       `gorm:"default:0" json:"total_sessions"`               // 总冥想次数
	TotalDuration      int       `gorm:"default:0" json:"total_duration"`               // 总冥想时长（分钟）
	AverageDuration    float64   `gorm:"type:decimal(5,1)" json:"average_duration"`     // 平均时长
	MostUsedType       string    `gorm:"type:varchar(50)" json:"most_used_type"`        // 最常用类型
	AverageImprovement float64   `gorm:"type:decimal(3,1)" json:"average_improvement"`  // 平均情绪改善程度
	StreakDays         int       `gorm:"default:0" json:"streak_days"`                  // 连续天数
	Distribution       string    `gorm:"type:json" json:"distribution"`                 // 类型分布 JSON
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func (MeditationAnalytics) TableName() string {
	return "meditation_analytics"
}

// MeditationFavorite 冥想内容收藏模型
type MeditationFavorite struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"not null;index:idx_user_content,unique" json:"user_id"`
	ContentID int64     `gorm:"not null;index:idx_user_content,unique" json:"content_id"`
	CreatedAt time.Time `json:"created_at"`

	// 关联
	Content *MeditationContent `gorm:"foreignKey:ContentID" json:"content,omitempty"`
}

func (MeditationFavorite) TableName() string {
	return "meditation_favorites"
}

// MeditationPlayRecord 冥想播放记录模型
type MeditationPlayRecord struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       int64     `gorm:"not null;index:idx_user_id" json:"user_id"`
	ContentID    int64     `gorm:"not null;index:idx_content_id" json:"content_id"`
	Duration     int       `gorm:"default:0" json:"duration"`                   // 实际播放时长（秒）
	Progress     float64   `gorm:"type:decimal(5,2);default:0" json:"progress"` // 播放进度（百分比）
	Completed    bool      `gorm:"default:false" json:"completed"`              // 是否完成
	LastPlayedAt time.Time `gorm:"index:idx_last_played" json:"last_played_at"` // 最后播放时间
	PlayCount    int       `gorm:"default:1" json:"play_count"`                 // 播放次数
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// 关联
	Content *MeditationContent `gorm:"foreignKey:ContentID" json:"content,omitempty"`
}

func (MeditationPlayRecord) TableName() string {
	return "meditation_play_records"
}
