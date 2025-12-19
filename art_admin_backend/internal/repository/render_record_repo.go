package repository

import (
	"art_admin_backend/internal/model"
	"time"

	"gorm.io/gorm"
)

// RenderRecordRepository 效果图记录仓库
// Render record repository
type RenderRecordRepository struct {
	db *gorm.DB
}

// NewRenderRecordRepository 创建效果图记录仓库
func NewRenderRecordRepository(db *gorm.DB) *RenderRecordRepository {
	return &RenderRecordRepository{db: db}
}

// Create 创建效果图记录
func (r *RenderRecordRepository) Create(record *model.RenderRecord) error {
	return r.db.Create(record).Error
}

// GetByID 根据ID获取记录
func (r *RenderRecordRepository) GetByID(id uint) (*model.RenderRecord, error) {
	var record model.RenderRecord
	err := r.db.First(&record, id).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// Update 更新记录
func (r *RenderRecordRepository) Update(record *model.RenderRecord) error {
	return r.db.Save(record).Error
}

// UpdateStatus 更新任务状态
func (r *RenderRecordRepository) UpdateStatus(id uint, status string, updates map[string]interface{}) error {
	updates["status"] = status
	return r.db.Model(&model.RenderRecord{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除记录
func (r *RenderRecordRepository) Delete(id uint) error {
	return r.db.Delete(&model.RenderRecord{}, id).Error
}

// RenderRecordFilter 效果图记录筛选条件
type RenderRecordFilter struct {
	UserID    int64
	ProjectID uint
	CadFileID uint
}

// List 获取效果图记录列表
func (r *RenderRecordRepository) List(page, pageSize int, filter RenderRecordFilter) ([]model.RenderRecord, int64, error) {
	var records []model.RenderRecord
	var total int64

	query := r.db.Model(&model.RenderRecord{})

	// 应用筛选条件 / Apply filter
	if filter.UserID > 0 {
		query = query.Where("user_id = ?", filter.UserID)
	}
	if filter.ProjectID > 0 {
		query = query.Where("project_id = ?", filter.ProjectID)
	}
	if filter.CadFileID > 0 {
		query = query.Where("cad_file_id = ?", filter.CadFileID)
	}

	// 获取总数 / Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询 / Paginated query
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// GetByUserID 获取用户的所有效果图记录
func (r *RenderRecordRepository) GetByUserID(userID int64) ([]model.RenderRecord, error) {
	var records []model.RenderRecord
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&records).Error
	return records, err
}

// CountTodayByUserID 统计用户今日保存次数
func (r *RenderRecordRepository) CountTodayByUserID(userID int64) (int64, error) {
	var count int64
	today := time.Now().Format("2006-01-02")
	err := r.db.Model(&model.RenderRecord{}).
		Where("user_id = ? AND DATE(created_at) = ?", userID, today).
		Count(&count).Error
	return count, err
}

// UserRenderQuotaRepository 用户配额仓库
type UserRenderQuotaRepository struct {
	db *gorm.DB
}

// NewUserRenderQuotaRepository 创建用户配额仓库
func NewUserRenderQuotaRepository(db *gorm.DB) *UserRenderQuotaRepository {
	return &UserRenderQuotaRepository{db: db}
}

// GetOrCreate 获取或创建用户配额记录
func (r *UserRenderQuotaRepository) GetOrCreate(userID int64) (*model.UserRenderQuota, error) {
	var quota model.UserRenderQuota
	err := r.db.Where("user_id = ?", userID).First(&quota).Error

	if err == gorm.ErrRecordNotFound {
		// 创建新记录 / Create new record
		quota = model.UserRenderQuota{
			UserID:        userID,
			DailyLimit:    3, // 默认每日3次 / Default 3 times per day
			UsedToday:     0,
			LastResetDate: time.Now(),
		}
		if err := r.db.Create(&quota).Error; err != nil {
			return nil, err
		}
		return &quota, nil
	}

	if err != nil {
		return nil, err
	}

	// 检查是否需要重置今日使用次数 / Check if need to reset today's usage
	today := time.Now().Format("2006-01-02")
	lastReset := quota.LastResetDate.Format("2006-01-02")
	if today != lastReset {
		quota.UsedToday = 0
		quota.LastResetDate = time.Now()
		r.db.Save(&quota)
	}

	return &quota, nil
}

// IncrementUsage 增加使用次数
func (r *UserRenderQuotaRepository) IncrementUsage(userID int64) error {
	return r.db.Model(&model.UserRenderQuota{}).
		Where("user_id = ?", userID).
		Update("used_today", gorm.Expr("used_today + 1")).Error
}

// UpdateDailyLimit 更新每日限制
func (r *UserRenderQuotaRepository) UpdateDailyLimit(userID int64, limit int) error {
	return r.db.Model(&model.UserRenderQuota{}).
		Where("user_id = ?", userID).
		Update("daily_limit", limit).Error
}

// ResetUsage 重置使用次数
func (r *UserRenderQuotaRepository) ResetUsage(userID int64) error {
	return r.db.Model(&model.UserRenderQuota{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"used_today":      0,
			"last_reset_date": time.Now(),
		}).Error
}

// GetByUserID 根据用户ID获取配额
func (r *UserRenderQuotaRepository) GetByUserID(userID int64) (*model.UserRenderQuota, error) {
	var quota model.UserRenderQuota
	err := r.db.Where("user_id = ?", userID).First(&quota).Error
	if err != nil {
		return nil, err
	}
	return &quota, nil
}
