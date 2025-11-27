package project

import (
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/dto/response"
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/repository"
	"fmt"
	"time"
)

type ProjectService struct {
	projectRepo  *repository.ProjectRepository
	templateRepo *repository.ProjectTemplateRepository
	taskRepo     *repository.TaskRepository
	userRepo     *repository.UserRepository
}

func NewProjectService() *ProjectService {
	return &ProjectService{
		projectRepo:  repository.NewProjectRepository(),
		templateRepo: repository.NewProjectTemplateRepository(),
		taskRepo:     repository.NewTaskRepository(),
		userRepo:     repository.NewUserRepository(),
	}
}

// generateProjectCode 生成项目编号
func (s *ProjectService) generateProjectCode() (string, error) {
	now := time.Now()
	year := now.Year()

	// 获取当年项目数量
	count, err := s.projectRepo.CountByYear(year)
	if err != nil {
		return "", err
	}

	// 生成编号格式: PROJ-YYYY-XXX
	sequence := count + 1
	code := fmt.Sprintf("PROJ-%d-%03d", year, sequence)

	// 确保编号唯一
	for {
		exists, err := s.projectRepo.ExistsByCode(code)
		if err != nil {
			return "", err
		}
		if !exists {
			break
		}
		sequence++
		code = fmt.Sprintf("PROJ-%d-%03d", year, sequence)
	}

	return code, nil
}

// GetProjectList 获取项目列表
func (s *ProjectService) GetProjectList(req *request.GetProjectListRequest) ([]response.ProjectResponse, int64, error) {
	projects, total, err := s.projectRepo.FindWithPagination(req)
	if err != nil {
		return nil, 0, err
	}

	result := make([]response.ProjectResponse, 0, len(projects))
	for _, p := range projects {
		// 查询经理姓名
		var managerName string
		if p.ManagerID > 0 {
			manager, err := s.userRepo.FindByID(p.ManagerID)
			if err == nil {
				managerName = manager.NickName
			}
		}

		result = append(result, response.ProjectResponse{
			ID:              p.ID,
			TemplateID:      p.TemplateID,
			Name:            p.Name,
			Code:            p.Code,
			Description:     p.Description,
			Category:        p.Category,
			Priority:        p.Priority,
			Status:          p.Status,
			Progress:        p.Progress,
			StartDate:       p.StartDate,
			EndDate:         p.EndDate,
			ActualStartDate: p.ActualStartDate,
			ActualEndDate:   p.ActualEndDate,
			ManagerID:       p.ManagerID,
			ManagerName:     managerName,
			CreatorID:       p.CreatorID,
			CreatorName:     p.CreatorName,
			Budget:          p.Budget,
			ActualCost:      p.ActualCost,
			Remark:          p.Remark,
			CreatedAt:       p.CreatedAt,
			UpdatedAt:       p.UpdatedAt,
		})
	}

	return result, total, nil
}

// GetProjectDetail 获取项目详情
func (s *ProjectService) GetProjectDetail(id int64) (*response.ProjectDetailResponse, error) {
	project, err := s.projectRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 查询经理姓名
	var managerName string
	if project.ManagerID > 0 {
		manager, err := s.userRepo.FindByID(project.ManagerID)
		if err == nil {
			managerName = manager.NickName
		}
	}

	// 获取项目成员
	members, err := s.projectRepo.FindMembersByProjectID(id)
	if err != nil {
		return nil, err
	}

	// 获取任务统计
	taskCount, completedCount, overdueCount, err := s.taskRepo.GetTaskStatistics(id)
	if err != nil {
		return nil, err
	}

	memberResponses := make([]response.ProjectMemberResponse, 0, len(members))
	for _, m := range members {
		memberResponses = append(memberResponses, response.ProjectMemberResponse{
			ID:        m.ID,
			ProjectID: m.ProjectID,
			UserID:    m.UserID,
			UserName:  m.UserName,
			Role:      m.Role,
			JoinDate:  m.JoinDate,
			CreatedAt: m.CreatedAt,
		})
	}

	result := &response.ProjectDetailResponse{
		ProjectResponse: response.ProjectResponse{
			ID:              project.ID,
			TemplateID:      project.TemplateID,
			Name:            project.Name,
			Code:            project.Code,
			Description:     project.Description,
			Category:        project.Category,
			Priority:        project.Priority,
			Status:          project.Status,
			Progress:        project.Progress,
			StartDate:       project.StartDate,
			EndDate:         project.EndDate,
			ActualStartDate: project.ActualStartDate,
			ActualEndDate:   project.ActualEndDate,
			ManagerID:       project.ManagerID,
			ManagerName:     managerName,
			CreatorID:       project.CreatorID,
			CreatorName:     project.CreatorName,
			Budget:          project.Budget,
			ActualCost:      project.ActualCost,
			Remark:          project.Remark,
			CreatedAt:       project.CreatedAt,
			UpdatedAt:       project.UpdatedAt,
		},
		Members:            memberResponses,
		TaskCount:          taskCount,
		CompletedTaskCount: completedCount,
		OverdueTaskCount:   overdueCount,
	}

	return result, nil
}

// CreateProject 创建项目
func (s *ProjectService) CreateProject(req *request.CreateProjectRequest, creatorID int64, creatorName string) error {
	// 自动生成项目编号
	code, err := s.generateProjectCode()
	if err != nil {
		return err
	}

	// 查询经理姓名
	var managerName string
	if req.ManagerID > 0 {
		manager, err := s.userRepo.FindByID(req.ManagerID)
		if err == nil {
			managerName = manager.NickName
		}
	}

	project := &model.Project{
		TemplateID:  req.TemplateID,
		Name:        req.Name,
		Code:        code,
		Description: req.Description,
		Category:    req.Category,
		Priority:    req.Priority,
		Status:      req.Status,
		Progress:    0,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		ManagerID:   req.ManagerID,
		ManagerName: managerName,
		CreatorID:   creatorID,
		CreatorName: creatorName,
		Budget:      req.Budget,
		Remark:      req.Remark,
	}

	return s.projectRepo.CreateWithMembers(project, req.Members)
}

// CreateProjectFromTemplate 从模板创建项目
func (s *ProjectService) CreateProjectFromTemplate(req *request.CreateProjectFromTemplateRequest, creatorID int64, creatorName string) error {
	// 自动生成项目编号
	code, err := s.generateProjectCode()
	if err != nil {
		return err
	}

	// 查询经理姓名
	var managerName string
	if req.ManagerID > 0 {
		manager, err := s.userRepo.FindByID(req.ManagerID)
		if err == nil {
			managerName = manager.NickName
		}
	}

	// 获取模板信息
	template, err := s.templateRepo.FindByID(req.TemplateID)
	if err != nil {
		return err
	}

	// 获取模板任务
	templateTasks, err := s.templateRepo.FindTasksByTemplateID(req.TemplateID)
	if err != nil {
		return err
	}

	// 获取模板依赖
	templateDeps, err := s.templateRepo.FindDependenciesByTemplateID(req.TemplateID)
	if err != nil {
		return err
	}

	// 创建项目
	project := &model.Project{
		TemplateID:  req.TemplateID,
		Name:        req.Name,
		Code:        code,
		Description: req.Description,
		Category:    template.Category,
		Priority:    2, // 默认中等优先级
		Status:      1, // 计划中
		Progress:    0,
		StartDate:   req.StartDate,
		EndDate:     req.StartDate.AddDate(0, 0, template.Duration),
		ManagerID:   req.ManagerID,
		ManagerName: managerName,
		CreatorID:   creatorID,
		CreatorName: creatorName,
		Budget:      req.Budget,
	}

	// 创建项目及其任务和依赖
	return s.projectRepo.CreateFromTemplate(project, req.Members, templateTasks, templateDeps)
}

// UpdateProject 更新项目
func (s *ProjectService) UpdateProject(req *request.UpdateProjectRequest) error {
	project, err := s.projectRepo.FindByID(req.ID)
	if err != nil {
		return err
	}

	// 查询经理姓名
	var managerName string
	if req.ManagerID > 0 {
		manager, err := s.userRepo.FindByID(req.ManagerID)
		if err == nil {
			managerName = manager.NickName
		}
	}

	project.Name = req.Name
	project.Description = req.Description
	project.Category = req.Category
	project.Priority = req.Priority
	project.Status = req.Status
	project.Progress = req.Progress
	project.StartDate = req.StartDate
	project.EndDate = req.EndDate
	project.ActualStartDate = req.ActualStartDate
	project.ActualEndDate = req.ActualEndDate
	project.ManagerID = req.ManagerID
	project.ManagerName = managerName
	project.Budget = req.Budget
	project.ActualCost = req.ActualCost
	project.Remark = req.Remark

	return s.projectRepo.Update(project)
}

// DeleteProject 删除项目
func (s *ProjectService) DeleteProject(id int64) error {
	return s.projectRepo.Delete(id)
}

// GetProjectStatistics 获取项目统计
func (s *ProjectService) GetProjectStatistics() (*response.ProjectStatisticsResponse, error) {
	total, err := s.projectRepo.Count(nil)
	if err != nil {
		return nil, err
	}

	ongoing, err := s.projectRepo.Count(map[string]interface{}{"status": 2})
	if err != nil {
		return nil, err
	}

	completed, err := s.projectRepo.Count(map[string]interface{}{"status": 3})
	if err != nil {
		return nil, err
	}

	// 统计超期项目
	now := time.Now()
	overdue, err := s.projectRepo.CountOverdue(now)
	if err != nil {
		return nil, err
	}

	return &response.ProjectStatisticsResponse{
		TotalProjects:     int(total),
		OngoingProjects:   int(ongoing),
		CompletedProjects: int(completed),
		OverdueProjects:   int(overdue),
	}, nil
}
