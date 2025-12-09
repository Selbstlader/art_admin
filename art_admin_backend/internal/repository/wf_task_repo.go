package repository

import (
	"time"

	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/database"

	"gorm.io/gorm"
)

// WfTaskRepository 工作流任务仓储
type WfTaskRepository struct{}

// NewWfTaskRepository 创建工作流任务仓储实例
func NewWfTaskRepository() *WfTaskRepository {
	return &WfTaskRepository{}
}

// Create 创建任务
func (r *WfTaskRepository) Create(task *model.WfTask) error {
	return database.DB.Create(task).Error
}

// Update 更新任务
func (r *WfTaskRepository) Update(task *model.WfTask) error {
	return database.DB.Save(task).Error
}

// FindByID 根据ID查询任务
func (r *WfTaskRepository) FindByID(id int64) (*model.WfTask, error) {
	var task model.WfTask
	err := database.DB.First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// Delete 删除任务（软删除）
func (r *WfTaskRepository) Delete(id int64) error {
	return database.DB.Delete(&model.WfTask{}, id).Error
}

// UpdateStatus 更新任务状态
func (r *WfTaskRepository) UpdateStatus(id int64, status string, comment string) error {
	now := time.Now()
	return database.DB.Model(&model.WfTask{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       status,
			"comment":      comment,
			"completed_at": &now,
		}).Error
}

// FindPendingByAssignee 查询用户的待办任务（我的待办）
func (r *WfTaskRepository) FindPendingByAssignee(assigneeID int64, query map[string]interface{}, current, size int) ([]model.WfTask, int64, error) {
	var tasks []model.WfTask
	var total int64

	db := database.DB.Model(&model.WfTask{}).
		Where("assignee_id = ? AND status = ?", assigneeID, model.WfTaskStatusPending)

	// 动态条件查询
	if nodeName, ok := query["nodeName"]; ok && nodeName != "" {
		db = db.Where("node_name LIKE ?", "%"+nodeName.(string)+"%")
	}
	if processInstID, ok := query["processInstId"]; ok && processInstID != nil {
		db = db.Where("process_inst_id = ?", processInstID)
	}

	// 统计总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (current - 1) * size
	err := db.Offset(offset).Limit(size).Order("id DESC").Find(&tasks).Error

	return tasks, total, err
}

// FindCompletedByAssignee 查询用户已处理的任务（我的已办）
func (r *WfTaskRepository) FindCompletedByAssignee(assigneeID int64, query map[string]interface{}, current, size int) ([]model.WfTask, int64, error) {
	var tasks []model.WfTask
	var total int64

	db := database.DB.Model(&model.WfTask{}).
		Where("assignee_id = ? AND status IN ?", assigneeID, []string{
			model.WfTaskStatusApproved,
			model.WfTaskStatusRejected,
		})

	// 动态条件查询
	if nodeName, ok := query["nodeName"]; ok && nodeName != "" {
		db = db.Where("node_name LIKE ?", "%"+nodeName.(string)+"%")
	}
	if processInstID, ok := query["processInstId"]; ok && processInstID != nil {
		db = db.Where("process_inst_id = ?", processInstID)
	}
	if status, ok := query["status"]; ok && status != "" {
		db = db.Where("status = ?", status)
	}

	// 统计总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (current - 1) * size
	err := db.Offset(offset).Limit(size).Order("completed_at DESC").Find(&tasks).Error

	return tasks, total, err
}

// FindByProcessInstID 根据流程实例ID查询所有任务
func (r *WfTaskRepository) FindByProcessInstID(processInstID int64) ([]model.WfTask, error) {
	var tasks []model.WfTask
	err := database.DB.Where("process_inst_id = ?", processInstID).
		Order("id ASC").Find(&tasks).Error
	return tasks, err
}

// FindPendingByProcessInstID 根据流程实例ID查询待处理任务
func (r *WfTaskRepository) FindPendingByProcessInstID(processInstID int64) ([]model.WfTask, error) {
	var tasks []model.WfTask
	err := database.DB.Where("process_inst_id = ? AND status = ?", processInstID, model.WfTaskStatusPending).
		Find(&tasks).Error
	return tasks, err
}

// FindByProcessInstIDAndNodeID 根据流程实例ID和节点ID查询任务
func (r *WfTaskRepository) FindByProcessInstIDAndNodeID(processInstID int64, nodeID string) ([]model.WfTask, error) {
	var tasks []model.WfTask
	err := database.DB.Where("process_inst_id = ? AND node_id = ?", processInstID, nodeID).
		Find(&tasks).Error
	return tasks, err
}

// CountPendingByProcessInstIDAndNodeID 统计指定节点的待处理任务数
func (r *WfTaskRepository) CountPendingByProcessInstIDAndNodeID(processInstID int64, nodeID string) (int64, error) {
	var count int64
	err := database.DB.Model(&model.WfTask{}).
		Where("process_inst_id = ? AND node_id = ? AND status = ?", processInstID, nodeID, model.WfTaskStatusPending).
		Count(&count).Error
	return count, err
}

// CountCompletedByProcessInstIDAndNodeID 统计指定节点的已完成任务数
func (r *WfTaskRepository) CountCompletedByProcessInstIDAndNodeID(processInstID int64, nodeID string) (int64, error) {
	var count int64
	err := database.DB.Model(&model.WfTask{}).
		Where("process_inst_id = ? AND node_id = ? AND status IN ?", processInstID, nodeID, []string{
			model.WfTaskStatusApproved,
			model.WfTaskStatusRejected,
		}).
		Count(&count).Error
	return count, err
}

// UpdateAssignee 更新任务审批人（用于委托/转办）
func (r *WfTaskRepository) UpdateAssignee(id int64, newAssigneeID int64, newAssigneeName string) error {
	return database.DB.Model(&model.WfTask{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"assignee_id":   newAssigneeID,
			"assignee_name": newAssigneeName,
		}).Error
}

// Delegate 委托任务 - 将任务委托给其他用户处理
// 原任务状态变为delegated，记录原审批人ID，更新新审批人
func (r *WfTaskRepository) Delegate(id int64, originalAssigneeID int64, toUserID int64, toUserName string) error {
	return database.DB.Model(&model.WfTask{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"delegated_from": originalAssigneeID,
			"assignee_id":    toUserID,
			"assignee_name":  toUserName,
			"status":         model.WfTaskStatusDelegated,
		}).Error
}

// DelegateWithTx 在事务中委托任务
func (r *WfTaskRepository) DelegateWithTx(tx *gorm.DB, id int64, originalAssigneeID int64, toUserID int64, toUserName string) error {
	return tx.Model(&model.WfTask{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"delegated_from": originalAssigneeID,
			"assignee_id":    toUserID,
			"assignee_name":  toUserName,
			"status":         model.WfTaskStatusDelegated,
		}).Error
}

// Transfer 转办任务 - 将任务完全转移给其他用户
// 原任务状态变为transferred，记录原审批人ID，更新新审批人
func (r *WfTaskRepository) Transfer(id int64, originalAssigneeID int64, toUserID int64, toUserName string) error {
	return database.DB.Model(&model.WfTask{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"transferred_from": originalAssigneeID,
			"assignee_id":      toUserID,
			"assignee_name":    toUserName,
			"status":           model.WfTaskStatusTransferred,
		}).Error
}

// TransferWithTx 在事务中转办任务
func (r *WfTaskRepository) TransferWithTx(tx *gorm.DB, id int64, originalAssigneeID int64, toUserID int64, toUserName string) error {
	return tx.Model(&model.WfTask{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"transferred_from": originalAssigneeID,
			"assignee_id":      toUserID,
			"assignee_name":    toUserName,
			"status":           model.WfTaskStatusTransferred,
		}).Error
}

// CancelPendingByProcessInstID 取消流程实例的所有待处理任务
func (r *WfTaskRepository) CancelPendingByProcessInstID(processInstID int64) error {
	return database.DB.Model(&model.WfTask{}).
		Where("process_inst_id = ? AND status = ?", processInstID, model.WfTaskStatusPending).
		Delete(&model.WfTask{}).Error
}

// Transaction 执行事务
func (r *WfTaskRepository) Transaction(fn func(tx *gorm.DB) error) error {
	return database.DB.Transaction(fn)
}

// CreateWithTx 在事务中创建任务
func (r *WfTaskRepository) CreateWithTx(tx *gorm.DB, task *model.WfTask) error {
	return tx.Create(task).Error
}

// UpdateWithTx 在事务中更新任务
func (r *WfTaskRepository) UpdateWithTx(tx *gorm.DB, task *model.WfTask) error {
	return tx.Save(task).Error
}

// CreateTaskLog 创建任务日志
func (r *WfTaskRepository) CreateTaskLog(log *model.TaskLog) error {
	return database.DB.Create(log).Error
}

// CreateTaskLogWithTx 在事务中创建任务日志
func (r *WfTaskRepository) CreateTaskLogWithTx(tx *gorm.DB, log *model.TaskLog) error {
	return tx.Create(log).Error
}

// FindTaskLogsByProcessInstID 根据流程实例ID查询任务日志（审批轨迹）
func (r *WfTaskRepository) FindTaskLogsByProcessInstID(processInstID int64) ([]model.TaskLog, error) {
	var logs []model.TaskLog
	err := database.DB.Where("process_inst_id = ?", processInstID).
		Order("created_at ASC").Find(&logs).Error
	return logs, err
}

// FindTaskLogsByTaskID 根据任务ID查询任务日志
func (r *WfTaskRepository) FindTaskLogsByTaskID(taskID int64) ([]model.TaskLog, error) {
	var logs []model.TaskLog
	err := database.DB.Where("task_id = ?", taskID).
		Order("created_at ASC").Find(&logs).Error
	return logs, err
}

// FindByAssigneeAndStatus 根据审批人ID和状态查询任务
func (r *WfTaskRepository) FindByAssigneeAndStatus(assigneeID int64, status string) ([]model.WfTask, error) {
	var tasks []model.WfTask
	err := database.DB.Where("assignee_id = ? AND status = ?", assigneeID, status).
		Order("created_at DESC").Find(&tasks).Error
	return tasks, err
}

// FindDelegatedOrTransferredByOriginalAssignee 查询被委托或转办的任务（通过原审批人ID）
func (r *WfTaskRepository) FindDelegatedOrTransferredByOriginalAssignee(originalAssigneeID int64) ([]model.WfTask, error) {
	var tasks []model.WfTask
	err := database.DB.Where("delegated_from = ? OR transferred_from = ?", originalAssigneeID, originalAssigneeID).
		Order("created_at DESC").Find(&tasks).Error
	return tasks, err
}

// UpdateStatusWithTx 在事务中更新任务状态
func (r *WfTaskRepository) UpdateStatusWithTx(tx *gorm.DB, id int64, status string, comment string) error {
	now := time.Now()
	return tx.Model(&model.WfTask{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       status,
			"comment":      comment,
			"completed_at": &now,
		}).Error
}

// FindByIDWithTx 在事务中根据ID查询任务
func (r *WfTaskRepository) FindByIDWithTx(tx *gorm.DB, id int64) (*model.WfTask, error) {
	var task model.WfTask
	err := tx.First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// BatchCreate 批量创建任务
func (r *WfTaskRepository) BatchCreate(tasks []*model.WfTask) error {
	if len(tasks) == 0 {
		return nil
	}
	return database.DB.Create(&tasks).Error
}

// BatchCreateWithTx 在事务中批量创建任务
func (r *WfTaskRepository) BatchCreateWithTx(tx *gorm.DB, tasks []*model.WfTask) error {
	if len(tasks) == 0 {
		return nil
	}
	return tx.Create(&tasks).Error
}

// CountByProcessInstIDAndStatus 统计流程实例中指定状态的任务数
func (r *WfTaskRepository) CountByProcessInstIDAndStatus(processInstID int64, status string) (int64, error) {
	var count int64
	err := database.DB.Model(&model.WfTask{}).
		Where("process_inst_id = ? AND status = ?", processInstID, status).
		Count(&count).Error
	return count, err
}

// FindOverdueTasks 查询超期任务
func (r *WfTaskRepository) FindOverdueTasks() ([]model.WfTask, error) {
	var tasks []model.WfTask
	now := time.Now()
	err := database.DB.Where("status = ? AND due_at IS NOT NULL AND due_at < ?", model.WfTaskStatusPending, now).
		Find(&tasks).Error
	return tasks, err
}

// FindTasksNearDeadline 查询即将超期的任务（指定小时数内）
func (r *WfTaskRepository) FindTasksNearDeadline(hoursBeforeDeadline int) ([]model.WfTask, error) {
	var tasks []model.WfTask
	now := time.Now()
	deadline := now.Add(time.Duration(hoursBeforeDeadline) * time.Hour)
	err := database.DB.Where("status = ? AND due_at IS NOT NULL AND due_at > ? AND due_at <= ?",
		model.WfTaskStatusPending, now, deadline).
		Find(&tasks).Error
	return tasks, err
}

// CancelPendingByProcessInstIDWithTx 在事务中取消流程实例的所有待处理任务
func (r *WfTaskRepository) CancelPendingByProcessInstIDWithTx(tx *gorm.DB, processInstID int64) error {
	return tx.Where("process_inst_id = ? AND status = ?", processInstID, model.WfTaskStatusPending).
		Delete(&model.WfTask{}).Error
}

// GetTaskStatistics 获取任务统计信息
func (r *WfTaskRepository) GetTaskStatistics(processInstID int64) (map[string]int64, error) {
	stats := make(map[string]int64)

	// 统计各状态任务数
	var results []struct {
		Status string
		Count  int64
	}

	err := database.DB.Model(&model.WfTask{}).
		Select("status, count(*) as count").
		Where("process_inst_id = ?", processInstID).
		Group("status").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	for _, r := range results {
		stats[r.Status] = r.Count
	}

	return stats, nil
}
