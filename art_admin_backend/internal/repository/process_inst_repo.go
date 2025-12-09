package repository

import (
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/database"

	"gorm.io/gorm"
)

// ProcessInstanceRepository 流程实例仓储
type ProcessInstanceRepository struct{}

// NewProcessInstanceRepository 创建流程实例仓储实例
func NewProcessInstanceRepository() *ProcessInstanceRepository {
	return &ProcessInstanceRepository{}
}

// Create 创建流程实例
func (r *ProcessInstanceRepository) Create(inst *model.ProcessInstance) error {
	return database.DB.Create(inst).Error
}

// Update 更新流程实例
func (r *ProcessInstanceRepository) Update(inst *model.ProcessInstance) error {
	return database.DB.Save(inst).Error
}

// FindByID 根据ID查询流程实例
func (r *ProcessInstanceRepository) FindByID(id int64) (*model.ProcessInstance, error) {
	var inst model.ProcessInstance
	err := database.DB.First(&inst, id).Error
	if err != nil {
		return nil, err
	}
	return &inst, nil
}

// Delete 删除流程实例（软删除）
func (r *ProcessInstanceRepository) Delete(id int64) error {
	return database.DB.Delete(&model.ProcessInstance{}, id).Error
}

// UpdateStatus 更新流程实例状态
func (r *ProcessInstanceRepository) UpdateStatus(id int64, status string) error {
	return database.DB.Model(&model.ProcessInstance{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// UpdateCurrentNode 更新当前节点
func (r *ProcessInstanceRepository) UpdateCurrentNode(id int64, nodeID string) error {
	return database.DB.Model(&model.ProcessInstance{}).
		Where("id = ?", id).
		Update("current_node_id", nodeID).Error
}

// FindByInitiator 查询用户发起的流程实例（我发起的）
func (r *ProcessInstanceRepository) FindByInitiator(initiatorID int64, query map[string]interface{}, current, size int) ([]model.ProcessInstance, int64, error) {
	var instances []model.ProcessInstance
	var total int64

	db := database.DB.Model(&model.ProcessInstance{}).Where("initiator_id = ?", initiatorID)

	// 动态条件查询
	if title, ok := query["title"]; ok && title != "" {
		db = db.Where("title LIKE ?", "%"+title.(string)+"%")
	}
	if status, ok := query["status"]; ok && status != "" {
		db = db.Where("status = ?", status)
	}
	if processDefID, ok := query["processDefId"]; ok && processDefID != nil {
		db = db.Where("process_def_id = ?", processDefID)
	}

	// 统计总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (current - 1) * size
	err := db.Offset(offset).Limit(size).Order("id DESC").Find(&instances).Error

	return instances, total, err
}

// FindWithPagination 分页查询流程实例列表
func (r *ProcessInstanceRepository) FindWithPagination(query map[string]interface{}, current, size int) ([]model.ProcessInstance, int64, error) {
	var instances []model.ProcessInstance
	var total int64

	db := database.DB.Model(&model.ProcessInstance{})

	// 动态条件查询
	if title, ok := query["title"]; ok && title != "" {
		db = db.Where("title LIKE ?", "%"+title.(string)+"%")
	}
	if status, ok := query["status"]; ok && status != "" {
		db = db.Where("status = ?", status)
	}
	if processDefID, ok := query["processDefId"]; ok && processDefID != nil {
		db = db.Where("process_def_id = ?", processDefID)
	}
	if initiatorID, ok := query["initiatorId"]; ok && initiatorID != nil {
		db = db.Where("initiator_id = ?", initiatorID)
	}

	// 统计总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (current - 1) * size
	err := db.Offset(offset).Limit(size).Order("id DESC").Find(&instances).Error

	return instances, total, err
}

// CountByProcessDefID 统计指定流程定义的实例数量
func (r *ProcessInstanceRepository) CountByProcessDefID(processDefID int64) (int64, error) {
	var count int64
	err := database.DB.Model(&model.ProcessInstance{}).
		Where("process_def_id = ?", processDefID).
		Count(&count).Error
	return count, err
}

// CountByStatus 按状态统计流程实例数量
func (r *ProcessInstanceRepository) CountByStatus(status string) (int64, error) {
	var count int64
	err := database.DB.Model(&model.ProcessInstance{}).
		Where("status = ?", status).
		Count(&count).Error
	return count, err
}

// FindRunningByProcessDefID 查询指定流程定义的运行中实例
func (r *ProcessInstanceRepository) FindRunningByProcessDefID(processDefID int64) ([]model.ProcessInstance, error) {
	var instances []model.ProcessInstance
	err := database.DB.Where("process_def_id = ? AND status = ?", processDefID, model.ProcessInstStatusRunning).
		Find(&instances).Error
	return instances, err
}

// Transaction 执行事务
func (r *ProcessInstanceRepository) Transaction(fn func(tx *gorm.DB) error) error {
	return database.DB.Transaction(fn)
}

// CreateWithTx 在事务中创建流程实例
func (r *ProcessInstanceRepository) CreateWithTx(tx *gorm.DB, inst *model.ProcessInstance) error {
	return tx.Create(inst).Error
}

// UpdateWithTx 在事务中更新流程实例
func (r *ProcessInstanceRepository) UpdateWithTx(tx *gorm.DB, inst *model.ProcessInstance) error {
	return tx.Save(inst).Error
}
