package designer_project

import (
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/dto/response"
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/repository"
	"errors"
	"time"
)

// DesignerProjectService 设计师项目服务
// Designer project service
type DesignerProjectService struct {
	repo *repository.DesignerProjectRepository
}

// NewDesignerProjectService 创建设计师项目服务
// Create designer project service
func NewDesignerProjectService(repo *repository.DesignerProjectRepository) *DesignerProjectService {
	return &DesignerProjectService{repo: repo}
}

// Create 创建项目
// Create project
func (s *DesignerProjectService) Create(req *request.CreateDesignerProjectRequest, userID uint) (*response.DesignerProjectResponse, error) {
	// 设置默认状态 / Set default status
	status := req.Status
	if status == "" {
		status = "draft"
	}

	project := &model.DesignerProject{
		Name:        req.Name,
		Description: req.Description,
		Area:        req.Area,
		Budget:      req.Budget,
		Style:       req.Style,
		Status:      status,
		UserID:      userID,
	}

	if err := s.repo.Create(project); err != nil {
		return nil, err
	}

	return s.toResponse(project), nil
}

// GetByID 根据ID获取项目
// Get project by ID
func (s *DesignerProjectService) GetByID(id uint, userID uint) (*response.DesignerProjectResponse, error) {
	project, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 验证用户权限 / Verify user permission
	if project.UserID != userID {
		return nil, errors.New("无权访问该项目")
	}

	return s.toResponse(project), nil
}

// GetDetailByID 根据ID获取项目详情（包含关联数据）
// Get project detail by ID (with associations)
func (s *DesignerProjectService) GetDetailByID(id uint, userID uint) (*response.DesignerProjectDetailResponse, error) {
	project, err := s.repo.GetByIDWithAssociations(id)
	if err != nil {
		return nil, err
	}

	// 验证用户权限 / Verify user permission
	if project.UserID != userID {
		return nil, errors.New("无权访问该项目")
	}

	return s.toDetailResponse(project), nil
}

// Update 更新项目
// Update project
func (s *DesignerProjectService) Update(req *request.UpdateDesignerProjectRequest, userID uint) (*response.DesignerProjectResponse, error) {
	// 获取现有项目 / Get existing project
	project, err := s.repo.GetByID(req.ID)
	if err != nil {
		return nil, err
	}

	// 验证用户权限 / Verify user permission
	if project.UserID != userID {
		return nil, errors.New("无权修改该项目")
	}

	// 更新字段 / Update fields
	project.Name = req.Name
	project.Description = req.Description
	project.Area = req.Area
	project.Budget = req.Budget
	project.Style = req.Style
	if req.Status != "" {
		project.Status = req.Status
	}

	if err := s.repo.Update(project); err != nil {
		return nil, err
	}

	return s.toResponse(project), nil
}

// Delete 删除项目（级联删除关联数据）
// Delete project (cascade delete associations)
func (s *DesignerProjectService) Delete(id uint, userID uint) error {
	// 获取项目 / Get project
	project, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	// 验证用户权限 / Verify user permission
	if project.UserID != userID {
		return errors.New("无权删除该项目")
	}

	// 级联删除 / Cascade delete
	return s.repo.DeleteWithAssociations(id)
}

// List 获取项目列表
// Get project list
func (s *DesignerProjectService) List(req *request.DesignerProjectListRequest, userID uint) (*response.DesignerProjectListResponse, error) {
	// 构建筛选条件 / Build filter criteria
	filter := repository.ProjectFilter{
		Name:      req.Name,
		Status:    req.Status,
		Style:     req.Style,
		MinBudget: req.MinBudget,
		MaxBudget: req.MaxBudget,
		Keyword:   req.Keyword,
		UserID:    userID,
	}

	// 解析日期 / Parse dates
	if req.StartDate != "" {
		if t, err := time.Parse("2006-01-02", req.StartDate); err == nil {
			filter.StartDate = t
		}
	}
	if req.EndDate != "" {
		if t, err := time.Parse("2006-01-02", req.EndDate); err == nil {
			// 设置为当天结束时间 / Set to end of day
			filter.EndDate = t.Add(24*time.Hour - time.Second)
		}
	}

	// 查询 / Query
	projects, total, err := s.repo.List(req.Current, req.Size, filter)
	if err != nil {
		return nil, err
	}

	// 转换响应 / Convert to response
	records := make([]response.DesignerProjectResponse, len(projects))
	for i, p := range projects {
		records[i] = *s.toResponse(&p)
	}

	return &response.DesignerProjectListResponse{
		Records: records,
		Current: req.Current,
		Size:    req.Size,
		Total:   total,
	}, nil
}

// Search 搜索项目
// Search projects
func (s *DesignerProjectService) Search(keyword string, userID uint, page, pageSize int) (*response.DesignerProjectListResponse, error) {
	projects, total, err := s.repo.Search(keyword, userID, page, pageSize)
	if err != nil {
		return nil, err
	}

	// 转换响应 / Convert to response
	records := make([]response.DesignerProjectResponse, len(projects))
	for i, p := range projects {
		records[i] = *s.toResponse(&p)
	}

	return &response.DesignerProjectListResponse{
		Records: records,
		Current: page,
		Size:    pageSize,
		Total:   total,
	}, nil
}

// BatchDelete 批量删除项目
// Batch delete projects
func (s *DesignerProjectService) BatchDelete(ids []uint, userID uint) error {
	for _, id := range ids {
		if err := s.Delete(id, userID); err != nil {
			// 继续删除其他项目，忽略单个失败 / Continue deleting others, ignore single failure
			continue
		}
	}
	return nil
}

// toResponse 转换为响应对象
// Convert to response object
func (s *DesignerProjectService) toResponse(project *model.DesignerProject) *response.DesignerProjectResponse {
	return &response.DesignerProjectResponse{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		Area:        project.Area,
		Budget:      project.Budget,
		Style:       project.Style,
		Status:      project.Status,
		UserID:      project.UserID,
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
	}
}

// toDetailResponse 转换为详情响应对象
// Convert to detail response object
func (s *DesignerProjectService) toDetailResponse(project *model.DesignerProject) *response.DesignerProjectDetailResponse {
	resp := &response.DesignerProjectDetailResponse{
		ID:            project.ID,
		Name:          project.Name,
		Description:   project.Description,
		Area:          project.Area,
		Budget:        project.Budget,
		Style:         project.Style,
		Status:        project.Status,
		UserID:        project.UserID,
		CreatedAt:     project.CreatedAt,
		UpdatedAt:     project.UpdatedAt,
		DocumentCount: len(project.Documents),
		CadFileCount:  len(project.CadFiles),
	}

	// 转换文档列表 / Convert document list
	resp.Documents = make([]response.ProjectDocumentBrief, len(project.Documents))
	for i, doc := range project.Documents {
		resp.Documents[i] = response.ProjectDocumentBrief{
			ID:             doc.ID,
			FileName:       doc.FileName,
			FileType:       doc.FileType,
			FileSize:       doc.FileSize,
			AnalysisStatus: doc.AnalysisStatus,
			CreatedAt:      doc.CreatedAt,
		}
	}

	// 转换CAD文件列表 / Convert CAD file list
	resp.CadFiles = make([]response.CadFileBrief, len(project.CadFiles))
	for i, cad := range project.CadFiles {
		resp.CadFiles[i] = response.CadFileBrief{
			ID:          cad.ID,
			FileName:    cad.FileName,
			FileFormat:  cad.FileFormat,
			ParseStatus: cad.ParseStatus,
			LayerCount:  cad.LayerCount,
			Has3D:       cad.Has3D,
			CreatedAt:   cad.CreatedAt,
		}
	}

	// 转换成本估算 / Convert cost estimate
	if project.CostEstimate != nil {
		resp.CostEstimate = &response.CostEstimateBrief{
			ID:             project.CostEstimate.ID,
			MaterialCost:   project.CostEstimate.MaterialCost,
			LaborCost:      project.CostEstimate.LaborCost,
			EquipmentCost:  project.CostEstimate.EquipmentCost,
			ManagementCost: project.CostEstimate.ManagementCost,
			TotalCost:      project.CostEstimate.TotalCost,
			BudgetLimit:    project.CostEstimate.BudgetLimit,
			UpdatedAt:      project.CostEstimate.UpdatedAt,
		}
	}

	// 转换比对结果列表 / Convert compare results list
	resp.CompareResults = make([]response.DesignCompareResultBrief, len(project.CompareResults))
	for i, cr := range project.CompareResults {
		resp.CompareResults[i] = response.DesignCompareResultBrief{
			ID:           cr.ID,
			OverallScore: cr.OverallScore,
			CreatedAt:    cr.CreatedAt,
		}
	}

	return resp
}
