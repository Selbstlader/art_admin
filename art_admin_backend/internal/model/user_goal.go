package model

import (
	"time"
)

// UserGoal 用户目标模型
type UserGoal struct {
	ID          int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int64      `gorm:"not null;index:idx_user_id" json:"user_id"`
	GoalType    string     `gorm:"type:varchar(50);not null" json:"goal_type"` // 目标类型: mood_record, meditation, journal
	TargetValue int        `gorm:"not null" json:"target_value"`               // 目标值（次数/分钟）
	Period      string     `gorm:"type:varchar(20);not null" json:"period"`    // 周期: daily, weekly, monthly
	IsActive    bool       `gorm:"default:true" json:"is_active"`              // 是否激活
	StartDate   time.Time  `gorm:"not null" json:"start_date"`                 // 开始日期
	EndDate     *time.Time `json:"end_date"`                                   // 结束日期（可选）
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (UserGoal) TableName() string {
	return "user_goals"
}

// GoalType 目标类型枚举
const (
	GoalTypeMoodRecord = "mood_record" // 情绪记录目标
	GoalTypeMeditation = "meditation"  // 冥想目标
	GoalTypeJournal    = "journal"     // 日记目标
)

// GoalPeriod 目标周期枚举
const (
	GoalPeriodDaily   = "daily"   // 每日
	GoalPeriodWeekly  = "weekly"  // 每周
	GoalPeriodMonthly = "monthly" // 每月
)

// UserCheckin 用户打卡记录模型
type UserCheckin struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int64     `gorm:"not null;index:idx_user_date" json:"user_id"`
	GoalID      int64     `gorm:"not null;index:idx_goal_date" json:"goal_id"`
	CheckinDate string    `gorm:"type:varchar(10);not null;index:idx_user_date,idx_goal_date" json:"checkin_date"` // 打卡日期 YYYY-MM-DD
	CheckinType string    `gorm:"type:varchar(50);not null" json:"checkin_type"`                                   // 打卡类型
	Value       int       `gorm:"default:0" json:"value"`                                                          // 完成值
	Note        string    `gorm:"type:text" json:"note"`                                                           // 打卡备注
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (UserCheckin) TableName() string {
	return "user_checkins"
}

// CheckinType 打卡类型枚举
const (
	CheckinTypeMood       = "mood"       // 情绪打卡
	CheckinTypeMeditation = "meditation" // 冥想打卡
	CheckinTypeJournal    = "journal"    // 日记打卡
	CheckinTypeGeneral    = "general"    // 通用打卡
)

// GoalProgress 目标进度模型
type GoalProgress struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       int64     `gorm:"not null;index:idx_user_goal_date" json:"user_id"`
	GoalID       int64     `gorm:"not null;index:idx_user_goal_date" json:"goal_id"`
	ProgressDate string    `gorm:"type:varchar(10);not null;index:idx_user_goal_date" json:"progress_date"` // 进度日期 YYYY-MM-DD
	CurrentValue int       `gorm:"default:0" json:"current_value"`                                          // 当前完成值
	IsCompleted  bool      `gorm:"default:false" json:"is_completed"`                                       // 是否完成
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (GoalProgress) TableName() string {
	return "goal_progress"
}
