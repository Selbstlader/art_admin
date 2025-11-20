package service

import (
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/dto/response"
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/repository"
	"errors"
)

type ProjectTemplateService struct {
	templateRepo *repository.ProjectTemplateRepository
}

func NewProjectTemplateService() *ProjectTemplateService {
	return &ProjectTemplateService{
		templateRepo: repository.NewProjectTemplateRepository(),
	}
}

// GetTemplateList 获取模板列表
func (s *ProjectTemplateService) GetTemplateList(req *request.GetProjectTemplateListRequest) ([]response.ProjectTemplateResponse, int64, error) {
	templates, total, err := s.templateRepo.FindWithPagination(req)
	if err != nil {
		return nil, 0, err
	}

	result := make([]response.ProjectTemplateResponse, 0, len(templates))
	for _, t := range templates {
		result = append(result, response.ProjectTemplateResponse{
			ID:          t.ID,
			Name:        t.Name,
			Description: t.Description,
			Category:    t.Category,
			Icon:        t.Icon,
			Duration:    t.Duration,
			IsPublic:    t.IsPublic,
			CreatorID:   t.CreatorID,
			CreatorName: t.CreatorName,
			UseCount:    t.UseCount,
			Status:      t.Status,
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
		})
	}

	return result, total, nil
}

// GetTemplateDetail 获取模板详情
func (s *ProjectTemplateService) GetTemplateDetail(id int64) (*response.ProjectTemplateDetailResponse, error) {
	template, err := s.templateRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 获取模板任务
	tasks, err := s.templateRepo.FindTasksByTemplateID(id)
	if err != nil {
		return nil, err
	}

	// 获取任务依赖
	dependencies, err := s.templateRepo.FindDependenciesByTemplateID(id)
	if err != nil {
		return nil, err
	}

	// 转换为响应DTO
	taskResponses := make([]response.TemplateTaskResponse, 0, len(tasks))
	for _, task := range tasks {
		taskResponses = append(taskResponses, response.TemplateTaskResponse{
			ID:          task.ID,
			TemplateID:  task.TemplateID,
			ParentID:    task.ParentID,
			Name:        task.Name,
			Description: task.Description,
			StartDay:    task.StartDay,
			Duration:    task.Duration,
			IsMilestone: task.IsMilestone,
			Priority:    task.Priority,
			Sort:        task.Sort,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
		})
	}

	depResponses := make([]response.TemplateDependencyResponse, 0, len(dependencies))
	for _, dep := range dependencies {
		depResponses = append(depResponses, response.TemplateDependencyResponse{
			ID:             dep.ID,
			TemplateID:     dep.TemplateID,
			PredecessorID:  dep.PredecessorID,
			SuccessorID:    dep.SuccessorID,
			DependencyType: dep.DependencyType,
			LagDays:        dep.LagDays,
			CreatedAt:      dep.CreatedAt,
		})
	}

	result := &response.ProjectTemplateDetailResponse{
		ProjectTemplateResponse: response.ProjectTemplateResponse{
			ID:          template.ID,
			Name:        template.Name,
			Description: template.Description,
			Category:    template.Category,
			Icon:        template.Icon,
			Duration:    template.Duration,
			IsPublic:    template.IsPublic,
			CreatorID:   template.CreatorID,
			CreatorName: template.CreatorName,
			UseCount:    template.UseCount,
			Status:      template.Status,
			CreatedAt:   template.CreatedAt,
			UpdatedAt:   template.UpdatedAt,
		},
		Tasks:        taskResponses,
		Dependencies: depResponses,
	}

	return result, nil
}

// CreateTemplate 创建模板
func (s *ProjectTemplateService) CreateTemplate(req *request.CreateProjectTemplateRequest, creatorID int64, creatorName string) error {
	template := &model.ProjectTemplate{
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Icon:        req.Icon,
		Duration:    req.Duration,
		IsPublic:    req.IsPublic,
		CreatorID:   creatorID,
		CreatorName: creatorName,
		Status:      1,
	}

	return s.templateRepo.CreateWithTasks(template, req.Tasks)
}

// UpdateTemplate 更新模板
func (s *ProjectTemplateService) UpdateTemplate(req *request.UpdateProjectTemplateRequest) error {
	template, err := s.templateRepo.FindByID(req.ID)
	if err != nil {
		return err
	}

	template.Name = req.Name
	template.Description = req.Description
	template.Category = req.Category
	template.Icon = req.Icon
	template.Duration = req.Duration
	template.IsPublic = req.IsPublic
	template.Status = req.Status

	return s.templateRepo.Update(template)
}

// DeleteTemplate 删除模板
func (s *ProjectTemplateService) DeleteTemplate(id int64) error {
	return s.templateRepo.Delete(id)
}

// CreateDependency 创建依赖关系
func (s *ProjectTemplateService) CreateDependency(req *request.CreateTemplateDependencyRequest) error {
	// 检查是否会形成循环依赖
	if req.PredecessorID == req.SuccessorID {
		return errors.New("任务不能依赖自己")
	}

	dependency := &model.TemplateTaskDependency{
		TemplateID:     req.TemplateID,
		PredecessorID:  req.PredecessorID,
		SuccessorID:    req.SuccessorID,
		DependencyType: req.DependencyType,
		LagDays:        req.LagDays,
	}

	return s.templateRepo.CreateDependency(dependency)
}

// DeleteDependency 删除依赖关系
func (s *ProjectTemplateService) DeleteDependency(id int64) error {
	return s.templateRepo.DeleteDependency(id)
}
