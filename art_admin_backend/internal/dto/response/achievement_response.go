package response

import "time"

// AchievementItem 成就项
type AchievementItem struct {
	ID          int64  `json:"id"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Tier        string `json:"tier"`
	TargetValue int    `json:"target_value"`
	Points      int    `json:"points"`
	IsHidden    bool   `json:"is_hidden"`
}

// AchievementProgressItem 成就进度项
type AchievementProgressItem struct {
	AchievementItem
	CurrentProgress int        `json:"current_progress"`
	ProgressPercent float64    `json:"progress_percent"`
	IsUnlocked      bool       `json:"is_unlocked"`
	UnlockedAt      *time.Time `json:"unlocked_at"`
}

// UserStatsItem 用户统计项
type UserStatsItem struct {
	TotalPoints        int                       `json:"total_points"`
	Level              int                       `json:"level"`
	ExperiencePoints   int                       `json:"experience_points"`
	UnlockedCount      int64                     `json:"unlocked_count"`
	TierCounts         map[string]int64          `json:"tier_counts"`
	CategoryPoints     map[string]int            `json:"category_points"`
	RecentAchievements []AchievementProgressItem `json:"recent_achievements"`
	NextLevelPoints    int                       `json:"next_level_points"` // 下一等级所需积分
	LevelProgress      float64                   `json:"level_progress"`    // 等级进度百分比
}

// LeaderboardItem 排行榜项
type LeaderboardItem struct {
	Rank          int    `json:"rank"`
	UserID        int64  `json:"user_id"`
	Username      string `json:"username"`
	Avatar        string `json:"avatar"`
	Value         int    `json:"value"` // 排行值（积分、连续天数等）
	Level         int    `json:"level"`
	Tier          string `json:"tier"` // 段位
	IsCurrentUser bool   `json:"is_current_user"`
}

// LeaderboardResponse 排行榜响应
type LeaderboardResponse struct {
	Type        string            `json:"type"`
	Period      string            `json:"period"`
	TotalCount  int               `json:"total_count"`
	CurrentUser *LeaderboardItem  `json:"current_user,omitempty"`
	TopUsers    []LeaderboardItem `json:"top_users"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// AchievementNotification 成就通知
type AchievementNotification struct {
	AchievementID int64  `json:"achievement_id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Icon          string `json:"icon"`
	Tier          string `json:"tier"`
	Points        int    `json:"points"`
	Message       string `json:"message"`
}

// LevelUpNotification 升级通知
type LevelUpNotification struct {
	OldLevel int      `json:"old_level"`
	NewLevel int      `json:"new_level"`
	Message  string   `json:"message"`
	Rewards  []string `json:"rewards"`
}

// CategoryStatsItem 分类统计项
type CategoryStatsItem struct {
	Category string  `json:"category"`
	Name     string  `json:"name"`
	Points   int     `json:"points"`
	Progress float64 `json:"progress"`
	Icon     string  `json:"icon"`
}

// AchievementTierStats 成就等级统计
type AchievementTierStats struct {
	Tier        string `json:"tier"`
	Name        string `json:"name"`
	Count       int64  `json:"count"`
	TotalPoints int    `json:"total_points"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
}
