package request

import "time"

// CreateMoodRecordRequest 创建情绪记录请求
type CreateMoodRecordRequest struct {
	MoodType   string   `json:"mood_type" binding:"required"`
	Intensity  int      `json:"intensity" binding:"required,min=1,max=10"`
	Triggers   []string `json:"triggers"`
	Activities []string `json:"activities"`
	Note       string   `json:"note"`
	Location   string   `json:"location"`
}

// UpdateMoodRecordRequest 更新情绪记录请求
type UpdateMoodRecordRequest struct {
	ID         int64    `json:"id" binding:"required"`
	MoodType   string   `json:"mood_type"`
	Intensity  *int     `json:"intensity" binding:"omitempty,min=1,max=10"`
	Triggers   []string `json:"triggers"`
	Activities []string `json:"activities"`
	Note       string   `json:"note"`
	Location   string   `json:"location"`
}

// MoodRecordListRequest 情绪记录列表查询请求
type MoodRecordListRequest struct {
	Page         int        `form:"page" binding:"required,min=1"`
	PageSize     int        `form:"pageSize" binding:"required,min=1,max=100"`
	MoodType     string     `form:"mood_type"`
	StartDate    *time.Time `form:"start_date"`
	EndDate      *time.Time `form:"end_date"`
	MinIntensity *int       `form:"min_intensity" binding:"omitempty,min=1,max=10"`
	MaxIntensity *int       `form:"max_intensity" binding:"omitempty,min=1,max=10"`
}

// MoodAnalyticsRequest 情绪分析请求
type MoodAnalyticsRequest struct {
	AnalyzeType string    `json:"analyze_type" binding:"required,oneof=daily weekly monthly"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	EndDate     time.Time `json:"end_date" binding:"required"`
}

// CreateMeditationSessionRequest 创建冥想会话请求
type CreateMeditationSessionRequest struct {
	UserID      int64     `json:"user_id" binding:"required"`
	ContentID   int64     `json:"content_id" binding:"required"`
	AnalyzeType string    `json:"analyze_type" binding:"required,oneof=daily weekly monthly"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	EndDate     time.Time `json:"end_date" binding:"required"`
}

// UpdateMeditationSessionRequest 更新冥想会话请求
type UpdateMeditationSessionRequest struct {
	ID                int64 `json:"id" binding:"required"`
	CompletedDuration int   `json:"completed_duration" binding:"required,min=0"`
	IsCompleted       *bool `json:"is_completed"`
}

// MeditationSessionListRequest 冥想会话列表请求
type MeditationSessionListRequest struct {
	Page        int        `form:"page" binding:"min=1"`
	PageSize    int        `form:"pageSize" binding:"min=1,max=100"`
	UserID      *int64     `form:"user_id"`
	ContentID   *int64     `form:"content_id"`
	IsCompleted *bool      `form:"is_completed"`
	StartDate   *time.Time `form:"start_date"`
	EndDate     *time.Time `form:"end_date"`
}

// UserAchievementListRequest 用户成就列表请求
type UserAchievementListRequest struct {
	Page            int    `form:"page" binding:"min=1"`
	PageSize        int    `form:"pageSize" binding:"min=1,max=100"`
	UserID          *int64 `form:"user_id"`
	AchievementType string `form:"achievement_type"`
	IsUnlocked      *bool  `form:"is_unlocked"`
}
