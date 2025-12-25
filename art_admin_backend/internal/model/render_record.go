package model

import (
	"time"

	"gorm.io/gorm"
)

// RenderRecord 效果图生成记录（异步任务）
// Render record for async effect image generation
type RenderRecord struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	UserID         int64          `gorm:"index;not null" json:"userId"`                  // 用户ID / User ID
	ProjectID      uint           `gorm:"index" json:"projectId"`                        // 项目ID / Project ID
	CadFileID      uint           `gorm:"index" json:"cadFileId"`                        // CAD文件ID / CAD file ID
	Status         string         `gorm:"size:20;default:'pending';index" json:"status"` // 状态: pending/processing/completed/failed
	ImageURL       string         `gorm:"size:500" json:"imageUrl"`                      // 效果图URL / Render image URL
	DesignProposal string         `gorm:"type:text" json:"designProposal"`               // 设计方案 / Design proposal
	Style          string         `gorm:"size:100" json:"style"`                         // 设计风格 / Design style
	RoomType       string         `gorm:"size:100" json:"roomType"`                      // 空间类型 / Room type
	Prompt         string         `gorm:"type:text" json:"prompt"`                       // 生成提示词 / Generation prompt
	ErrorMessage   string         `gorm:"size:500" json:"errorMessage"`                  // 错误信息 / Error message
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// 注意：移除 GORM 关联定义以避免迁移时的索引冲突
	// 如需关联查询，请使用 Preload 或手动 Join
	// Note: Removed GORM associations to avoid index conflicts during migration
	// Use Preload or manual Join for association queries if needed
}

// TableName 表名
func (RenderRecord) TableName() string {
	return "render_records"
}

// UserRenderQuota 用户效果图生成配额
// User render quota for daily limit control
type UserRenderQuota struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        int64     `gorm:"not null;index:idx_user_render_quotas_user_id,unique" json:"userId"` // 用户ID / User ID
	DailyLimit    int       `gorm:"default:3" json:"dailyLimit"`                                        // 每日限制次数 / Daily limit
	UsedToday     int       `gorm:"default:0" json:"usedToday"`                                         // 今日已使用次数 / Used today
	LastResetDate time.Time `json:"lastResetDate"`                                                      // 上次重置日期 / Last reset date
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// TableName 表名
func (UserRenderQuota) TableName() string {
	return "user_render_quotas"
}
