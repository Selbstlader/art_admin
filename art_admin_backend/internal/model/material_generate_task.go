package model

import "time"

// MaterialGenerateTask 教材生成任务表
type MaterialGenerateTask struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     int64     `gorm:"not null;index;comment:用户ID" json:"user_id"`
	SubjectID  int64     `gorm:"not null;comment:学科ID" json:"subject_id"`
	Grade      string    `gorm:"type:varchar(50);not null;comment:年级" json:"grade"`
	Topic      string    `gorm:"type:varchar(200);not null;comment:学习题材" json:"topic"`
	Difficulty int8      `gorm:"default:1;comment:难度 1-基础 2-进阶 3-高级" json:"difficulty"`
	Status     int8      `gorm:"default:0;index;comment:状态 0-待处理 1-生成中 2-已完成 3-失败" json:"status"`
	MaterialID int64     `gorm:"comment:生成的教材ID" json:"material_id"`
	ErrorMsg   string    `gorm:"type:text;comment:错误信息" json:"error_msg"`
	CreatedAt  time.Time `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (MaterialGenerateTask) TableName() string {
	return "material_generate_tasks"
}

// 任务状态常量
const (
	TaskStatusPending    = 0 // 待处理
	TaskStatusProcessing = 1 // 生成中
	TaskStatusCompleted  = 2 // 已完成
	TaskStatusFailed     = 3 // 失败
)
