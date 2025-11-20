package repository

import (
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/database"
)

type TaskRepository struct{}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{}
}

// FindWithPagination 分页查询任务列表
func (r *TaskRepository) FindWithPagination(req *request.GetTaskListRequest) ([]model.Task, int64, error) {
	var tasks []model.Task
	var total int64

	db := database.DB.Model(&model.Task{}).Where("project_id = ?", req.ProjectID)

	// 条件查询
	if req.ParentID != nil {
		db = db.Where("parent_id = ?", *req.ParentID)
	}
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Type != "" {
		db = db.Where("type = ?", req.Type)
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if req.AssigneeID != nil {
		db = db.Where("assignee_id = ?", *req.AssigneeID)
	}
	if req.IsOverdue != nil {
		db = db.Where("is_overdue = ?", *req.IsOverdue)
	}

	// 统计总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	if err := db.Offset(offset).Limit(req.PageSize).Order("sort ASC, created_at DESC").Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

// FindByID 根据ID查询任务
func (r *TaskRepository) FindByID(id int64) (*model.Task, error) {
	var task model.Task
	if err := database.DB.First(&task, id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

// FindByProjectID 根据项目ID查询所有任务
func (r *TaskRepository) FindByProjectID(projectID int64) ([]model.Task, error) {
	var tasks []model.Task
	if err := database.DB.Where("project_id = ?", projectID).Order("sort ASC, id ASC").Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

// FindByParentID 根据父任务ID查询子任务
func (r *TaskRepository) FindByParentID(parentID int64) ([]model.Task, error) {
	var tasks []model.Task
	if err := database.DB.Where("parent_id = ?", parentID).Order("sort ASC, id ASC").Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

// FindDependenciesByTaskID 查询任务的依赖关系
func (r *TaskRepository) FindDependenciesByTaskID(taskID int64) ([]model.TaskDependency, error) {
	var dependencies []model.TaskDependency
	if err := database.DB.Where("predecessor_id = ? OR successor_id = ?", taskID, taskID).Find(&dependencies).Error; err != nil {
		return nil, err
	}
	return dependencies, nil
}

// FindDependenciesByProjectID 查询项目的所有依赖关系
func (r *TaskRepository) FindDependenciesByProjectID(projectID int64) ([]model.TaskDependency, error) {
	var dependencies []model.TaskDependency
	if err := database.DB.Where("project_id = ?", projectID).Find(&dependencies).Error; err != nil {
		return nil, err
	}
	return dependencies, nil
}

// FindCommentsByTaskID 查询任务评论
func (r *TaskRepository) FindCommentsByTaskID(taskID int64) ([]model.TaskComment, error) {
	var comments []model.TaskComment
	if err := database.DB.Where("task_id = ?", taskID).Order("created_at DESC").Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

// FindAttachmentsByTaskID 查询任务附件
func (r *TaskRepository) FindAttachmentsByTaskID(taskID int64) ([]model.TaskAttachment, error) {
	var attachments []model.TaskAttachment
	if err := database.DB.Where("task_id = ?", taskID).Order("created_at DESC").Find(&attachments).Error; err != nil {
		return nil, err
	}
	return attachments, nil
}

// GetTaskStatistics 获取任务统计
func (r *TaskRepository) GetTaskStatistics(projectID int64) (taskCount, completedCount, overdueCount int, err error) {
	// 总任务数
	var total int64
	if err := database.DB.Model(&model.Task{}).Where("project_id = ?", projectID).Count(&total).Error; err != nil {
		return 0, 0, 0, err
	}
	taskCount = int(total)

	// 已完成任务数
	var completed int64
	if err := database.DB.Model(&model.Task{}).Where("project_id = ? AND status = ?", projectID, 3).Count(&completed).Error; err != nil {
		return 0, 0, 0, err
	}
	completedCount = int(completed)

	// 超期任务数
	var overdue int64
	if err := database.DB.Model(&model.Task{}).Where("project_id = ? AND is_overdue = ?", projectID, true).Count(&overdue).Error; err != nil {
		return 0, 0, 0, err
	}
	overdueCount = int(overdue)

	return taskCount, completedCount, overdueCount, nil
}

// Create 创建任务
func (r *TaskRepository) Create(task *model.Task) error {
	return database.DB.Create(task).Error
}

// Update 更新任务
func (r *TaskRepository) Update(task *model.Task) error {
	return database.DB.Save(task).Error
}

// Delete 删除任务
func (r *TaskRepository) Delete(id int64) error {
	return database.DB.Delete(&model.Task{}, id).Error
}

// CreateDependency 创建依赖关系
func (r *TaskRepository) CreateDependency(dependency *model.TaskDependency) error {
	return database.DB.Create(dependency).Error
}

// DeleteDependency 删除依赖关系
func (r *TaskRepository) DeleteDependency(id int64) error {
	return database.DB.Delete(&model.TaskDependency{}, id).Error
}

// CreateComment 创建评论
func (r *TaskRepository) CreateComment(comment *model.TaskComment) error {
	return database.DB.Create(comment).Error
}

// CreateAttachment 创建附件
func (r *TaskRepository) CreateAttachment(attachment *model.TaskAttachment) error {
	return database.DB.Create(attachment).Error
}
