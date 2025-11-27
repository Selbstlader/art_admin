package response

import (
	"art_admin_backend/internal/model"
	"time"
)

// MoodRecordItem 情绪记录响应项
type MoodRecordItem struct {
	ID             int64     `json:"id"`
	UserID         int64     `json:"user_id"`
	MoodType       string    `json:"mood_type"`
	MoodLabel      string    `json:"mood_label"`
	MoodTypeLabel  string    `json:"mood_type_label"`
	Intensity      int       `json:"intensity"`
	Triggers       []string  `json:"triggers"`
	Activities     []string  `json:"activities"`
	Note           string    `json:"note"`
	Location       string    `json:"location"`
	AnalysisStatus string    `json:"analysis_status"`
	DifyDocumentID string    `json:"dify_document_id"`
	CreatedAt      time.Time `json:"created_at"`
}

// UserMeditationSessionItem 用户冥想会话响应项
type UserMeditationSessionItem struct {
	ID                int64      `json:"id"`
	UserID            int64      `json:"user_id"`
	ContentID         int64      `json:"content_id"`
	ContentTitle      string     `json:"content_title"`
	CompletedDuration int        `json:"completed_duration"`
	Duration          int        `json:"duration"`
	IsCompleted       bool       `json:"is_completed"`
	CompletionRate    float64    `json:"completion_rate"` // 完成率百分比
	StartedAt         *time.Time `json:"started_at"`
	CompletedAt       *time.Time `json:"completed_at"`
	CreatedAt         time.Time  `json:"created_at"`
}

// UserAchievementItem 用户成就响应项
type UserAchievementItem struct {
	ID              int64      `json:"id"`
	UserID          int64      `json:"user_id"`
	AchievementType string     `json:"achievement_type"`
	Name            string     `json:"name"`
	Description     string     `json:"description"`
	Icon            string     `json:"icon"`
	TargetValue     int        `json:"target_value"`
	CurrentValue    int        `json:"current_value"`
	Progress        float64    `json:"progress"`
	IsUnlocked      bool       `json:"is_unlocked"`
	UnlockedAt      *time.Time `json:"unlocked_at"`
	CreatedAt       time.Time  `json:"created_at"`
}

// MoodAnalyticsItem 情绪分析响应项
type MoodAnalyticsItem struct {
	ID                  int64          `json:"id"`
	UserID              int64          `json:"user_id"`
	AnalyzeType         string         `json:"analyze_type"`
	AnalyzeDate         time.Time      `json:"analyze_date"`
	MoodDistribution    map[string]int `json:"mood_distribution"`
	AverageIntensity    float64        `json:"average_intensity"`
	TotalRecords        int            `json:"total_records"`
	MostCommonMood      string         `json:"most_common_mood"`
	MostCommonMoodLabel string         `json:"most_common_mood_label"`
	CreatedAt           time.Time      `json:"created_at"`
}

// MoodStatistics 情绪统计概览
type MoodStatistics struct {
	TotalRecords        int                    `json:"total_records"`
	AverageIntensity    float64                `json:"average_intensity"`
	MoodDistribution    map[string]int         `json:"mood_distribution"`
	MostCommonMood      string                 `json:"most_common_mood"`
	MostCommonMoodLabel string                 `json:"most_common_mood_label"`
	RecentTrend         []MoodRecordItem       `json:"recent_trend"`      // 最近7天趋势
	WeeklyComparison    map[string]interface{} `json:"weekly_comparison"` // 周对比
}

// MeditationStatistics 冥想统计概览
type MeditationStatistics struct {
	TotalSessions      int                         `json:"total_sessions"`
	TotalDuration      int                         `json:"total_duration"`       // 总时长(秒)
	AverageDuration    int                         `json:"average_duration"`     // 平均时长(秒)
	CompletionRate     float64                     `json:"completion_rate"`      // 完成率
	StreakDays         int                         `json:"streak_days"`          // 连续天数
	FavoriteCategory   string                      `json:"favorite_category"`    // 最喜欢的分类
	RecentSessions     []UserMeditationSessionItem `json:"recent_sessions"`      // 最近会话
	MostUsedType       string                      `json:"most_used_type"`       // 最常用的冥想类型
	MostUsedTypeLabel  string                      `json:"most_used_type_label"` // 最常用的冥想类型标签
	AverageImprovement float64                     `json:"average_improvement"`  // 平均改善程度
	RecordingStreak    int                         `json:"recording_streak"`     // 记录连续天数
	TotalCount         int                         `json:"total_count"`          // 总记录数
	MoodCorrelation    float64                     `json:"mood_correlation"`     // 情绪相关性
	TypeDistribution   map[string]int              `json:"type_distribution"`    // 类型分布
	AnalysisDate       time.Time                   `json:"analysis_date"`        // 分析日期
}

// SentimentData 情感数据点
type SentimentData struct {
	Date           time.Time `json:"date"`
	SentimentScore float64   `json:"sentiment_score"`
	EntryCount     int       `json:"entry_count"`
}

// DashboardOverview 用户仪表板概览
type DashboardOverview struct {
	UserID               int64                 `json:"user_id"`
	MoodStatistics       MoodStatistics        `json:"mood_statistics"`
	MeditationStatistics MeditationStatistics  `json:"meditation_statistics"`
	JournalStatistics    JournalStatistics     `json:"journal_statistics"`
	RecentAchievements   []UserAchievementItem `json:"recent_achievements"`
	DailyGoals           DailyGoals            `json:"daily_goals"`
	LastUpdated          time.Time             `json:"last_updated"`
}

// DailyGoals 每日目标
type DailyGoals struct {
	MoodRecorded        bool `json:"mood_recorded"`
	MeditationCompleted bool `json:"meditation_completed"`
	JournalWritten      bool `json:"journal_written"`
	GoalsCompleted      int  `json:"goals_completed"`
	TotalGoals          int  `json:"total_goals"`
}

// FromModel 从模型转换为响应项
// FromModel 从模型转换为响应项
func (item *MoodRecordItem) FromModel(moodRecord *model.MoodRecord) {
	item.ID = moodRecord.ID
	item.UserID = moodRecord.UserID
	item.MoodType = moodRecord.MoodType
	item.MoodLabel = model.GetMoodTypeLabel(model.MoodType(moodRecord.MoodType))
	item.Intensity = moodRecord.Intensity
	// JSON字段需要解析
	// item.Triggers = parseJSON(moodRecord.Triggers)
	// item.Activities = parseJSON(moodRecord.Activities)
	item.Note = moodRecord.Note
	item.Location = moodRecord.Location
	item.CreatedAt = moodRecord.CreatedAt
}
