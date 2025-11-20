package repository

import (
	"time"

	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/database"
)

type OperationLogRepository struct{}

func NewOperationLogRepository() *OperationLogRepository {
	return &OperationLogRepository{}
}

// Create 创建操作日志
func (r *OperationLogRepository) Create(log *model.OperationLog) error {
	log.OperationTime = time.Now()
	return database.DB.Create(log).Error
}

// FindByID 根据ID查询操作日志
func (r *OperationLogRepository) FindByID(id int64) (*model.OperationLog, error) {
	var log model.OperationLog
	err := database.DB.First(&log, id).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

// Delete 删除操作日志
func (r *OperationLogRepository) Delete(id int64) error {
	return database.DB.Delete(&model.OperationLog{}, id).Error
}

// BatchDelete 批量删除操作日志
func (r *OperationLogRepository) BatchDelete(ids []int64) error {
	return database.DB.Delete(&model.OperationLog{}, ids).Error
}

// FindWithPagination 分页查询操作日志列表
func (r *OperationLogRepository) FindWithPagination(query map[string]interface{}, current, size int) ([]model.OperationLog, int64, error) {
	var logs []model.OperationLog
	var total int64

	db := database.DB.Model(&model.OperationLog{})

	// 动态条件查询
	if module, ok := query["module"]; ok && module != "" {
		db = db.Where("module LIKE ?", "%"+module.(string)+"%")
	}
	if businessType, ok := query["businessType"]; ok && businessType != "" {
		db = db.Where("business_type = ?", businessType)
	}
	if operatorName, ok := query["operatorName"]; ok && operatorName != "" {
		db = db.Where("operator_name LIKE ?", "%"+operatorName.(string)+"%")
	}
	if status, ok := query["status"]; ok && status != nil {
		db = db.Where("status = ?", status)
	}
	if startTime, ok := query["startTime"]; ok && startTime != "" {
		db = db.Where("operation_time >= ?", startTime)
	}
	if endTime, ok := query["endTime"]; ok && endTime != "" {
		db = db.Where("operation_time <= ?", endTime)
	}

	// 统计总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (current - 1) * size
	err := db.Offset(offset).Limit(size).Order("operation_time DESC").Find(&logs).Error

	return logs, total, err
}

// CleanOldLogs 清理旧日志(保留最近N天的日志)
func (r *OperationLogRepository) CleanOldLogs(days int) error {
	cutoffTime := time.Now().AddDate(0, 0, -days)
	return database.DB.Where("operation_time < ?", cutoffTime).Delete(&model.OperationLog{}).Error
}
