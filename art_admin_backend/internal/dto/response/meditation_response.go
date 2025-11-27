package response

import "time"

// MeditationContentItem 冥想内容项
type MeditationContentItem struct {
	ID              int64    `json:"id"`
	Title           string   `json:"title"`
	Category        string   `json:"category"`
	CategoryLabel   string   `json:"category_label"`
	Duration        int      `json:"duration"`
	AudioURL        string   `json:"audio_url"`
	CoverImage      string   `json:"cover_image"`
	DifficultyLevel int      `json:"difficulty_level"`
	Tags            []string `json:"tags"`
	IsPublic        bool     `json:"is_public"`
	ViewCount       int64    `json:"view_count"`
	LikeCount       int64    `json:"like_count"`
}

// MeditationContentDetail 冥想内容详情
type MeditationContentDetail struct {
	MeditationContentItem
	Description string    `json:"description"`
	CreatedBy   int64     `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PlayRecordItem 播放记录项
type PlayRecordItem struct {
	ID           int64     `json:"id"`
	ContentID    int64     `json:"content_id"`
	Duration     int       `json:"duration"`       // 实际播放时长（秒）
	Progress     float64   `json:"progress"`       // 播放进度（百分比）
	Completed    bool      `json:"completed"`      // 是否完成
	LastPlayedAt time.Time `json:"last_played_at"` // 最后播放时间
	PlayCount    int       `json:"play_count"`     // 播放次数
}

// PlayHistoryItem 播放历史项
type PlayHistoryItem struct {
	ID              int64     `json:"id"`
	ContentID       int64     `json:"content_id"`
	ContentTitle    string    `json:"content_title"`    // 内容标题
	ContentImage    string    `json:"content_image"`    // 内容封面
	ContentDuration int       `json:"content_duration"` // 内容总时长
	Duration        int       `json:"duration"`         // 实际播放时长
	Progress        float64   `json:"progress"`         // 播放进度
	Completed       bool      `json:"completed"`        // 是否完成
	LastPlayedAt    time.Time `json:"last_played_at"`   // 最后播放时间
	PlayCount       int       `json:"play_count"`       // 播放次数
}
