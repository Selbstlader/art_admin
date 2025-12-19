package repository

import (
	"art_admin_backend/internal/model"
	"time"

	"gorm.io/gorm"
)

// NotificationRepository 通知仓库
// Notification repository
type NotificationRepository struct {
	db *gorm.DB
}

// NewNotificationRepository 创建通知仓库
func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

// Create 创建通知
func (r *NotificationRepository) Create(notification *model.Notification) error {
	return r.db.Create(notification).Error
}

// GetByID 根据ID获取通知
func (r *NotificationRepository) GetByID(id uint) (*model.Notification, error) {
	var notification model.Notification
	err := r.db.First(&notification, id).Error
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

// ListByUserID 获取用户通知列表（只返回今天的）
func (r *NotificationRepository) ListByUserID(userID int64, page, pageSize int) ([]model.Notification, int64, error) {
	var notifications []model.Notification
	var total int64

	// 只查询今天的通知 / Only query today's notifications
	today := time.Now().Format("2006-01-02")

	query := r.db.Model(&model.Notification{}).
		Where("user_id = ?", userID).
		Where("DATE(created_at) = ?", today)

	// 获取总数 / Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询 / Paginated query
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&notifications).Error; err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

// GetUnreadCount 获取用户未读通知数量（只统计今天的）
func (r *NotificationRepository) GetUnreadCount(userID int64) (int64, error) {
	var count int64
	today := time.Now().Format("2006-01-02")

	err := r.db.Model(&model.Notification{}).
		Where("user_id = ? AND is_read = ? AND DATE(created_at) = ?", userID, false, today).
		Count(&count).Error

	return count, err
}

// MarkAsRead 标记为已读
func (r *NotificationRepository) MarkAsRead(id uint, userID int64) error {
	return r.db.Model(&model.Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_read", true).Error
}

// MarkAllAsRead 标记全部已读
func (r *NotificationRepository) MarkAllAsRead(userID int64) error {
	today := time.Now().Format("2006-01-02")
	return r.db.Model(&model.Notification{}).
		Where("user_id = ? AND DATE(created_at) = ?", userID, today).
		Update("is_read", true).Error
}

// DeleteOldNotifications 删除过期通知（保留最近N天）
func (r *NotificationRepository) DeleteOldNotifications(days int) error {
	cutoff := time.Now().AddDate(0, 0, -days)
	return r.db.Where("created_at < ?", cutoff).Delete(&model.Notification{}).Error
}

// Delete 删除通知
func (r *NotificationRepository) Delete(id uint, userID int64) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Notification{}).Error
}
