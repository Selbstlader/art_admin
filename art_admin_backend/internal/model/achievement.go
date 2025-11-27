package model

import (
	"time"
)

// Achievement 成就模板定义
type Achievement struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Type        string    `gorm:"type:varchar(50);not null;index" json:"type"` // 成就类型
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`      // 成就名称
	Description string    `gorm:"type:text" json:"description"`                // 成就描述
	Icon        string    `gorm:"type:varchar(100)" json:"icon"`               // 图标
	Tier        string    `gorm:"type:varchar(20);not null" json:"tier"`       // 等级: bronze, silver, gold, platinum
	TargetValue int       `gorm:"not null" json:"target_value"`                // 目标值
	Points      int       `gorm:"default:0" json:"points"`                     // 奖励积分
	IsHidden    bool      `gorm:"default:false" json:"is_hidden"`              // 是否隐藏成就
	Criteria    string    `gorm:"type:json" json:"criteria"`                   // 触发条件 JSON
	IsActive    bool      `gorm:"default:true" json:"is_active"`               // 是否启用
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Achievement) TableName() string {
	return "achievements"
}

// UserAchievement 用户成就记录
type UserAchievement struct {
	ID            int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        int64      `gorm:"not null;index" json:"user_id"`
	AchievementID int64      `gorm:"not null;index" json:"achievement_id"`
	CurrentValue  int        `gorm:"default:0" json:"current_value"`   // 当前进度
	IsUnlocked    bool       `gorm:"default:false" json:"is_unlocked"` // 是否已解锁
	UnlockedAt    *time.Time `json:"unlocked_at"`                      // 解锁时间
	Progress      float64    `gorm:"-" json:"progress"`                // 进度百分比 (计算字段)
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`

	// 关联
	User        *User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Achievement *Achievement `gorm:"foreignKey:AchievementID" json:"achievement,omitempty"`
}

func (UserAchievement) TableName() string {
	return "user_achievements"
}

// UserPoint 用户积分记录
type UserPoint struct {
	ID               int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID           int64     `gorm:"not null;index" json:"user_id"`
	TotalPoints      int       `gorm:"default:0" json:"total_points"`      // 总积分
	MoodPoints       int       `gorm:"default:0" json:"mood_points"`       // 情绪记录积分
	MeditationPoints int       `gorm:"default:0" json:"meditation_points"` // 冥想积分
	JournalPoints    int       `gorm:"default:0" json:"journal_points"`    // 日记积分
	AnalysisPoints   int       `gorm:"default:0" json:"analysis_points"`   // AI分析积分
	BonusPoints      int       `gorm:"default:0" json:"bonus_points"`      // 额外积分
	Level            int       `gorm:"default:1" json:"level"`             // 用户等级
	ExperiencePoints int       `gorm:"default:0" json:"experience_points"` // 经验值
	LastUpdated      time.Time `gorm:"index" json:"last_updated"`          // 最后更新时间
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	// 关联
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (UserPoint) TableName() string {
	return "user_points"
}

// AchievementEvent 成就事件记录
type AchievementEvent struct {
	ID           int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       int64      `gorm:"not null;index" json:"user_id"`
	EventType    string     `gorm:"type:varchar(50);not null;index" json:"event_type"` // 事件类型
	EventData    string     `gorm:"type:json" json:"event_data"`                       // 事件数据 JSON
	PointsEarned int        `gorm:"default:0" json:"points_earned"`                    // 获得积分
	ProcessedAt  *time.Time `json:"processed_at"`                                      // 处理时间
	CreatedAt    time.Time  `json:"created_at"`

	// 关联
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (AchievementEvent) TableName() string {
	return "achievement_events"
}

// Leaderboard 排行榜缓存
type Leaderboard struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Type        string    `gorm:"type:varchar(50);not null;index" json:"type"`   // 排行榜类型
	Period      string    `gorm:"type:varchar(20);not null;index" json:"period"` // 周期: daily, weekly, monthly, all_time
	UserRanking string    `gorm:"type:json" json:"user_ranking"`                 // 用户排名数据 JSON
	UpdatedAt   time.Time `gorm:"index" json:"updated_at"`
}

func (Leaderboard) TableName() string {
	return "leaderboards"
}

// 成就类型常量
const (
	AchievementTypeMoodRecord = "mood_record" // 情绪记录
	AchievementTypeMeditation = "meditation"  // 冥想
	AchievementTypeJournal    = "journal"     // 日记
	AchievementTypeAnalysis   = "analysis"    // AI分析
	AchievementTypeStreak     = "streak"      // 连续记录
	AchievementTypeQuality    = "quality"     // 质量指标
	AchievementTypeEngagement = "engagement"  // 参与度
	AchievementTypeGrowth     = "growth"      // 成长
)

// 成就等级常量
const (
	AchievementTierBronze   = "bronze"
	AchievementTierSilver   = "silver"
	AchievementTierGold     = "gold"
	AchievementTierPlatinum = "platinum"
)

// 事件类型常量
const (
	EventTypeMoodRecord   = "mood_record_created"
	EventTypeMeditation   = "meditation_completed"
	EventTypeJournal      = "journal_created"
	EventTypeAnalysis     = "analysis_completed"
	EventTypeStreakUpdate = "streak_updated"
	EventTypeLevelUp      = "level_up"
	EventTypeBonus        = "bonus_earned"
)

// 排行榜类型常量
const (
	LeaderboardTypePoints     = "points"
	LeaderboardTypeStreak     = "streak"
	LeaderboardTypeMeditation = "meditation"
	LeaderboardTypeJournal    = "journal"
)

// 排行榜周期常量
const (
	PeriodDaily   = "daily"
	PeriodWeekly  = "weekly"
	PeriodMonthly = "monthly"
	PeriodAllTime = "all_time"
)
