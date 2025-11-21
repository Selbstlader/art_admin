package repository

import (
	"context"

	"art_admin_backend/internal/model"

	"gorm.io/gorm"
)

// MaterialGenerateTaskRepository 教材生成任务仓储
type MaterialGenerateTaskRepository struct {
	db *gorm.DB
}

// NewMaterialGenerateTaskRepository 创建教材生成任务仓储
func NewMaterialGenerateTaskRepository(db *gorm.DB) *MaterialGenerateTaskRepository {
	return &MaterialGenerateTaskRepository{db: db}
}

// Create 创建任务
func (r *MaterialGenerateTaskRepository) Create(ctx context.Context, task *model.MaterialGenerateTask) error {
	return r.db.WithContext(ctx).Create(task).Error
}

// GetByID 根据 ID 获取任务
func (r *MaterialGenerateTaskRepository) GetByID(ctx context.Context, id int64) (*model.MaterialGenerateTask, error) {
	var task model.MaterialGenerateTask
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&task).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

// UpdateStatus 更新任务状态
func (r *MaterialGenerateTaskRepository) UpdateStatus(ctx context.Context, id int64, status int8, materialID int64, errorMsg string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if materialID > 0 {
		updates["material_id"] = materialID
	}
	if errorMsg != "" {
		updates["error_msg"] = errorMsg
	}
	return r.db.WithContext(ctx).Model(&model.MaterialGenerateTask{}).Where("id = ?", id).Updates(updates).Error
}

// ListByUserID 获取用户的任务列表
func (r *MaterialGenerateTaskRepository) ListByUserID(ctx context.Context, userID int64, statuses []int8) ([]*model.MaterialGenerateTask, error) {
	var tasks []*model.MaterialGenerateTask
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if len(statuses) > 0 {
		query = query.Where("status IN ?", statuses)
	}
	if err := query.Order("created_at DESC").Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

// GetPendingTask 获取待处理的任务
func (r *MaterialGenerateTaskRepository) GetPendingTask(ctx context.Context) (*model.MaterialGenerateTask, error) {
	var task model.MaterialGenerateTask
	if err := r.db.WithContext(ctx).
		Where("status = ?", model.TaskStatusPending).
		Order("created_at ASC").
		First(&task).Error; err != nil {
		return nil, err
	}
	return &task, nil
}
