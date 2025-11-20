package repository

import (
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/database"
	"time"
)

type ProjectRepository struct{}

func NewProjectRepository() *ProjectRepository {
	return &ProjectRepository{}
}

// FindWithPagination 分页查询项目列表
func (r *ProjectRepository) FindWithPagination(req *request.GetProjectListRequest) ([]model.Project, int64, error) {
	var projects []model.Project
	var total int64

	db := database.DB.Model(&model.Project{})

	// 条件查询
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Code != "" {
		db = db.Where("code LIKE ?", "%"+req.Code+"%")
	}
	if req.Category != "" {
		db = db.Where("category = ?", req.Category)
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if req.ManagerID != nil {
		db = db.Where("manager_id = ?", *req.ManagerID)
	}
	if req.StartDate != nil {
		db = db.Where("start_date >= ?", *req.StartDate)
	}
	if req.EndDate != nil {
		db = db.Where("end_date <= ?", *req.EndDate)
	}

	// 统计总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	if err := db.Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&projects).Error; err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}

// FindByID 根据ID查询项目
func (r *ProjectRepository) FindByID(id int64) (*model.Project, error) {
	var project model.Project
	if err := database.DB.First(&project, id).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

// FindMembersByProjectID 查询项目成员
func (r *ProjectRepository) FindMembersByProjectID(projectID int64) ([]model.ProjectMember, error) {
	var members []model.ProjectMember
	if err := database.DB.Where("project_id = ?", projectID).Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

// ExistsByCode 检查项目编号是否存在
func (r *ProjectRepository) ExistsByCode(code string) (bool, error) {
	var count int64
	if err := database.DB.Model(&model.Project{}).Where("code = ?", code).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// CreateWithMembers 创建项目及成员
func (r *ProjectRepository) CreateWithMembers(project *model.Project, members []request.ProjectMemberReq) error {
	tx := database.DB.Begin()

	// 创建项目
	if err := tx.Create(project).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 创建项目成员
	for _, memberReq := range members {
		member := &model.ProjectMember{
			ProjectID: project.ID,
			UserID:    memberReq.UserID,
			Role:      memberReq.Role,
			JoinDate:  memberReq.JoinDate,
		}
		if err := tx.Create(member).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// CreateFromTemplate 从模板创建项目
func (r *ProjectRepository) CreateFromTemplate(
	project *model.Project,
	members []request.ProjectMemberReq,
	templateTasks []model.TemplateTask,
	templateDeps []model.TemplateTaskDependency,
) error {
	tx := database.DB.Begin()

	// 创建项目
	if err := tx.Create(project).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 创建项目成员
	for _, memberReq := range members {
		member := &model.ProjectMember{
			ProjectID: project.ID,
			UserID:    memberReq.UserID,
			Role:      memberReq.Role,
			JoinDate:  memberReq.JoinDate,
		}
		if err := tx.Create(member).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// 创建任务(从模板任务转换)
	taskIDMap := make(map[int64]int64) // 模板任务ID -> 新任务ID
	for _, templateTask := range templateTasks {
		task := &model.Task{
			ProjectID:      project.ID,
			ParentID:       0, // 稍后更新
			TemplateTaskID: templateTask.ID,
			Name:           templateTask.Name,
			Description:    templateTask.Description,
			Type:           "task",
			Priority:       templateTask.Priority,
			Status:         1, // 未开始
			Progress:       0,
			StartDate:      project.StartDate.AddDate(0, 0, templateTask.StartDay),
			EndDate:        project.StartDate.AddDate(0, 0, templateTask.StartDay+templateTask.Duration),
			Sort:           templateTask.Sort,
			CreatorID:      project.CreatorID,
			CreatorName:    project.CreatorName,
		}

		if templateTask.IsMilestone {
			task.Type = "milestone"
		}

		if err := tx.Create(task).Error; err != nil {
			tx.Rollback()
			return err
		}

		taskIDMap[templateTask.ID] = task.ID
	}

	// 更新父任务ID
	for _, templateTask := range templateTasks {
		if templateTask.ParentID > 0 {
			if newParentID, ok := taskIDMap[templateTask.ParentID]; ok {
				if newTaskID, ok := taskIDMap[templateTask.ID]; ok {
					if err := tx.Model(&model.Task{}).Where("id = ?", newTaskID).Update("parent_id", newParentID).Error; err != nil {
						tx.Rollback()
						return err
					}
				}
			}
		}
	}

	// 创建任务依赖关系
	for _, templateDep := range templateDeps {
		if predecessorID, ok := taskIDMap[templateDep.PredecessorID]; ok {
			if successorID, ok := taskIDMap[templateDep.SuccessorID]; ok {
				dependency := &model.TaskDependency{
					ProjectID:      project.ID,
					PredecessorID:  predecessorID,
					SuccessorID:    successorID,
					DependencyType: templateDep.DependencyType,
					LagDays:        templateDep.LagDays,
				}
				if err := tx.Create(dependency).Error; err != nil {
					tx.Rollback()
					return err
				}
			}
		}
	}

	return tx.Commit().Error
}

// Update 更新项目
func (r *ProjectRepository) Update(project *model.Project) error {
	return database.DB.Save(project).Error
}

// Delete 删除项目
func (r *ProjectRepository) Delete(id int64) error {
	return database.DB.Delete(&model.Project{}, id).Error
}

// Count 统计项目数量
func (r *ProjectRepository) Count(conditions map[string]interface{}) (int64, error) {
	var count int64
	db := database.DB.Model(&model.Project{})

	for key, value := range conditions {
		db = db.Where(key+" = ?", value)
	}

	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// CountOverdue 统计超期项目数量
func (r *ProjectRepository) CountOverdue(now time.Time) (int64, error) {
	var count int64
	if err := database.DB.Model(&model.Project{}).
		Where("end_date < ? AND status IN (1, 2)", now).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountByYear 统计指定年份的项目数量
func (r *ProjectRepository) CountByYear(year int) (int64, error) {
	var count int64
	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)

	if err := database.DB.Model(&model.Project{}).
		Where("created_at >= ? AND created_at < ?", start, end).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
