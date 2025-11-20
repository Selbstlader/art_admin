package repository

import (
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/database"
)

type ProjectTemplateRepository struct{}

func NewProjectTemplateRepository() *ProjectTemplateRepository {
	return &ProjectTemplateRepository{}
}

// FindWithPagination 分页查询模板列表
func (r *ProjectTemplateRepository) FindWithPagination(req *request.GetProjectTemplateListRequest) ([]model.ProjectTemplate, int64, error) {
	var templates []model.ProjectTemplate
	var total int64

	db := database.DB.Model(&model.ProjectTemplate{})

	// 条件查询
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Category != "" {
		db = db.Where("category = ?", req.Category)
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}

	// 统计总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	if err := db.Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&templates).Error; err != nil {
		return nil, 0, err
	}

	return templates, total, nil
}

// FindByID 根据ID查询模板
func (r *ProjectTemplateRepository) FindByID(id int64) (*model.ProjectTemplate, error) {
	var template model.ProjectTemplate
	if err := database.DB.First(&template, id).Error; err != nil {
		return nil, err
	}
	return &template, nil
}

// FindTasksByTemplateID 查询模板任务
func (r *ProjectTemplateRepository) FindTasksByTemplateID(templateID int64) ([]model.TemplateTask, error) {
	var tasks []model.TemplateTask
	if err := database.DB.Where("template_id = ?", templateID).Order("sort ASC, id ASC").Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

// FindDependenciesByTemplateID 查询模板依赖关系
func (r *ProjectTemplateRepository) FindDependenciesByTemplateID(templateID int64) ([]model.TemplateTaskDependency, error) {
	var dependencies []model.TemplateTaskDependency
	if err := database.DB.Where("template_id = ?", templateID).Find(&dependencies).Error; err != nil {
		return nil, err
	}
	return dependencies, nil
}

// CreateWithTasks 创建模板及任务
func (r *ProjectTemplateRepository) CreateWithTasks(template *model.ProjectTemplate, tasks []request.CreateTemplateTaskReq) error {
	tx := database.DB.Begin()

	// 创建模板
	if err := tx.Create(template).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 创建任务
	for _, taskReq := range tasks {
		task := &model.TemplateTask{
			TemplateID:  template.ID,
			ParentID:    taskReq.ParentID,
			Name:        taskReq.Name,
			Description: taskReq.Description,
			StartDay:    taskReq.StartDay,
			Duration:    taskReq.Duration,
			IsMilestone: taskReq.IsMilestone,
			Priority:    taskReq.Priority,
			Sort:        taskReq.Sort,
		}
		if err := tx.Create(task).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// Update 更新模板
func (r *ProjectTemplateRepository) Update(template *model.ProjectTemplate) error {
	return database.DB.Save(template).Error
}

// Delete 删除模板
func (r *ProjectTemplateRepository) Delete(id int64) error {
	tx := database.DB.Begin()

	// 删除模板任务
	if err := tx.Where("template_id = ?", id).Delete(&model.TemplateTask{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 删除模板依赖
	if err := tx.Where("template_id = ?", id).Delete(&model.TemplateTaskDependency{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 删除模板
	if err := tx.Delete(&model.ProjectTemplate{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// CreateDependency 创建依赖关系
func (r *ProjectTemplateRepository) CreateDependency(dependency *model.TemplateTaskDependency) error {
	return database.DB.Create(dependency).Error
}

// DeleteDependency 删除依赖关系
func (r *ProjectTemplateRepository) DeleteDependency(id int64) error {
	return database.DB.Delete(&model.TemplateTaskDependency{}, id).Error
}
