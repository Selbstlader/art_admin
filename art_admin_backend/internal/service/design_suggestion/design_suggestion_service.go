package design_suggestion

import (
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/dto/response"
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/database"
	"art_admin_backend/internal/pkg/volcengine"
	"art_admin_backend/internal/repository"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// DesignSuggestionService 设计建议服务
// Design suggestion service
type DesignSuggestionService struct {
	repo        *repository.DesignSuggestionRepository
	projectRepo *repository.DesignerProjectRepository
	aiClient    *volcengine.Client
	db          *gorm.DB
}

// NewDesignSuggestionService 创建设计建议服务
// Create design suggestion service
func NewDesignSuggestionService(
	repo *repository.DesignSuggestionRepository,
	projectRepo *repository.DesignerProjectRepository,
	aiClient *volcengine.Client,
) *DesignSuggestionService {
	return &DesignSuggestionService{
		repo:        repo,
		projectRepo: projectRepo,
		aiClient:    aiClient,
		db:          database.GetDB(),
	}
}

// ProjectContext 项目上下文信息，用于生成建议
// Project context for generating suggestions
type ProjectContext struct {
	Project       *model.DesignerProject
	Materials     []model.ProjectMaterial
	CostSummary   *CostSummaryInfo
	CadFilesInfo  []CadFileInfo
	DocumentsInfo []DocumentInfo
}

// CostSummaryInfo 成本汇总信息
// Cost summary information
type CostSummaryInfo struct {
	MaterialCost   float64 `json:"materialCost"`
	LaborCost      float64 `json:"laborCost"`
	EquipmentCost  float64 `json:"equipmentCost"`
	ManagementCost float64 `json:"managementCost"`
	TotalCost      float64 `json:"totalCost"`
	BudgetLimit    float64 `json:"budgetLimit"`
	BudgetExceeded bool    `json:"budgetExceeded"`
}

// CadFileInfo CAD文件信息
// CAD file information
type CadFileInfo struct {
	FileName    string   `json:"fileName"`
	FileFormat  string   `json:"fileFormat"`
	LayerCount  int      `json:"layerCount"`
	Layers      []string `json:"layers"`
	Has3D       bool     `json:"has3d"`
	ParseStatus string   `json:"parseStatus"`
}

// DocumentInfo 文档信息
// Document information
type DocumentInfo struct {
	FileName       string `json:"fileName"`
	FileType       string `json:"fileType"`
	Summary        string `json:"summary"`
	AnalysisStatus string `json:"analysisStatus"`
}

// AISuggestion AI生成的建议结构
// AI generated suggestion structure
type AISuggestion struct {
	Content         string   `json:"content"`
	ApplicableScene string   `json:"applicableScene"`
	CostImpact      string   `json:"costImpact"`
	Category        string   `json:"category"`
	Priority        int      `json:"priority"`
	DetailInfo      any      `json:"detailInfo"`
	ReferenceImages []string `json:"referenceImages"`
}

// GenerateSuggestionsAsync 异步生成设计建议
// Generate design suggestions asynchronously
// Requirements: 3.1, 3.2
func (s *DesignSuggestionService) GenerateSuggestionsAsync(req *request.GenerateDesignSuggestionRequest, userID uint) (*response.GenerateSuggestionAsyncResponse, error) {
	// 获取项目信息 / Get project info
	project, err := s.projectRepo.GetByIDWithAssociations(req.ProjectID)
	if err != nil {
		return nil, errors.New("项目不存在")
	}

	// 验证用户权限 / Verify user permission
	if project.UserID != userID {
		return nil, errors.New("无权访问该项目")
	}

	// 设置默认生成数量 / Set default count
	count := req.Count
	if count <= 0 {
		count = 3 // 默认生成3条建议 / Default 3 suggestions
	}
	if count > 10 {
		count = 10 // 最多10条 / Max 10 suggestions
	}

	// 创建生成任务记录 / Create generation task record
	task := &model.SuggestionTask{
		ProjectID: req.ProjectID,
		UserID:    userID,
		Category:  req.Category,
		Count:     count,
		Status:    "processing",
	}
	if err := s.db.Create(task).Error; err != nil {
		return nil, fmt.Errorf("创建任务失败: %w", err)
	}

	// 启动异步任务 / Start async task
	go s.doGenerateSuggestions(task.ID, project, req.Category, count, userID)

	return &response.GenerateSuggestionAsyncResponse{
		TaskID:    task.ID,
		ProjectID: req.ProjectID,
		Status:    "processing",
		Message:   "建议生成任务已提交，完成后将通过通知提醒您",
	}, nil
}

// doGenerateSuggestions 执行建议生成任务
// Execute suggestion generation task
func (s *DesignSuggestionService) doGenerateSuggestions(taskID uint, project *model.DesignerProject, category string, count int, userID uint) {
	var generatedCount int
	var errMsg string

	defer func() {
		// 更新任务状态 / Update task status
		status := "completed"
		if errMsg != "" {
			status = "failed"
		}
		s.db.Model(&model.SuggestionTask{}).Where("id = ?", taskID).Updates(map[string]interface{}{
			"status":          status,
			"generated_count": generatedCount,
			"error_message":   errMsg,
		})

		// 发送通知 / Send notification
		s.sendGenerationNotification(userID, project.ID, project.Name, generatedCount, errMsg)
	}()

	// 获取完整项目上下文 / Get full project context
	projectContext, err := s.getProjectContext(project)
	if err != nil {
		errMsg = fmt.Sprintf("获取项目上下文失败: %v", err)
		return
	}

	// 构建AI提示词 / Build AI prompt
	systemPrompt := s.buildSystemPrompt()
	userPrompt := s.buildUserPromptWithContext(projectContext, category, count)

	// 调用AI生成建议 / Call AI to generate suggestions
	aiResponse, err := s.aiClient.ChatWithText(systemPrompt, userPrompt)
	if err != nil {
		errMsg = fmt.Sprintf("AI服务调用失败: %v", err)
		return
	}

	// 解析AI响应 / Parse AI response
	suggestions, err := s.parseAIResponse(aiResponse)
	if err != nil {
		errMsg = fmt.Sprintf("解析AI响应失败: %v", err)
		return
	}

	// 保存建议到数据库 / Save suggestions to database
	for _, aiSug := range suggestions {
		detailJSON, _ := json.Marshal(aiSug.DetailInfo)
		refImagesJSON, _ := json.Marshal(aiSug.ReferenceImages)

		suggestion := model.DesignSuggestion{
			ProjectID:       project.ID,
			Content:         aiSug.Content,
			ApplicableScene: aiSug.ApplicableScene,
			CostImpact:      aiSug.CostImpact,
			Category:        aiSug.Category,
			Priority:        aiSug.Priority,
			Status:          "pending",
			DetailInfo:      string(detailJSON),
			ReferenceImages: string(refImagesJSON),
			UserID:          userID,
		}

		if err := s.repo.Create(&suggestion); err != nil {
			continue // 跳过保存失败的建议 / Skip failed saves
		}
		generatedCount++
	}

	if generatedCount == 0 {
		errMsg = "未能生成有效的设计建议"
	}
}

// sendGenerationNotification 发送生成完成通知
// Send generation completion notification
func (s *DesignSuggestionService) sendGenerationNotification(userID uint, projectID uint, projectName string, count int, errMsg string) {
	var title, content string
	noticeType := "notice"

	if errMsg != "" {
		title = "设计建议生成失败"
		content = fmt.Sprintf("项目【%s】的设计建议生成失败：%s", projectName, errMsg)
	} else {
		title = "设计建议生成完成"
		content = fmt.Sprintf("项目【%s】已成功生成 %d 条设计建议，请前往查看", projectName, count)
	}

	// 创建通知 / Create notification
	notification := &model.Notification{
		UserID:    int64(userID),
		Title:     title,
		Content:   content,
		Type:      noticeType,
		RelatedID: projectID,
		IsRead:    false,
	}
	s.db.Create(notification)
}

// GetTaskStatus 获取任务状态
// Get task status
func (s *DesignSuggestionService) GetTaskStatus(taskID uint, userID uint) (*response.SuggestionTaskStatusResponse, error) {
	var task model.SuggestionTask
	if err := s.db.First(&task, taskID).Error; err != nil {
		return nil, errors.New("任务不存在")
	}

	if task.UserID != userID {
		return nil, errors.New("无权访问该任务")
	}

	return &response.SuggestionTaskStatusResponse{
		TaskID:         task.ID,
		ProjectID:      task.ProjectID,
		Status:         task.Status,
		GeneratedCount: task.GeneratedCount,
		ErrorMessage:   task.ErrorMessage,
		CreatedAt:      task.CreatedAt,
	}, nil
}

// GenerateSuggestions 同步生成设计建议（保留兼容性）
// Generate design suggestions synchronously (kept for compatibility)
// Requirements: 3.1, 3.2
func (s *DesignSuggestionService) GenerateSuggestions(req *request.GenerateDesignSuggestionRequest, userID uint) (*response.GenerateDesignSuggestionResponse, error) {
	// 获取项目信息 / Get project info
	project, err := s.projectRepo.GetByIDWithAssociations(req.ProjectID)
	if err != nil {
		return nil, errors.New("项目不存在")
	}

	// 验证用户权限 / Verify user permission
	if project.UserID != userID {
		return nil, errors.New("无权访问该项目")
	}

	// 设置默认生成数量 / Set default count
	count := req.Count
	if count <= 0 {
		count = 3 // 默认生成3条建议 / Default 3 suggestions
	}
	if count > 10 {
		count = 10 // 最多10条 / Max 10 suggestions
	}

	// 获取完整项目上下文 / Get full project context
	projectContext, err := s.getProjectContext(project)
	if err != nil {
		return nil, fmt.Errorf("获取项目上下文失败: %w", err)
	}

	// 构建AI提示词 / Build AI prompt
	systemPrompt := s.buildSystemPrompt()
	userPrompt := s.buildUserPromptWithContext(projectContext, req.Category, count)

	// 调用AI生成建议 / Call AI to generate suggestions
	aiResponse, err := s.aiClient.ChatWithText(systemPrompt, userPrompt)
	if err != nil {
		return nil, fmt.Errorf("AI服务调用失败: %w", err)
	}

	// 解析AI响应 / Parse AI response
	suggestions, err := s.parseAIResponse(aiResponse)
	if err != nil {
		return nil, fmt.Errorf("解析AI响应失败: %w", err)
	}

	// 确保至少有3条建议 / Ensure at least 3 suggestions
	if len(suggestions) < 3 {
		return nil, errors.New("AI生成的建议数量不足，请重试")
	}

	// 保存建议到数据库 / Save suggestions to database
	var savedSuggestions []model.DesignSuggestion
	for _, aiSug := range suggestions {
		detailJSON, _ := json.Marshal(aiSug.DetailInfo)
		refImagesJSON, _ := json.Marshal(aiSug.ReferenceImages)

		suggestion := model.DesignSuggestion{
			ProjectID:       req.ProjectID,
			Content:         aiSug.Content,
			ApplicableScene: aiSug.ApplicableScene,
			CostImpact:      aiSug.CostImpact,
			Category:        aiSug.Category,
			Priority:        aiSug.Priority,
			Status:          "pending",
			DetailInfo:      string(detailJSON),
			ReferenceImages: string(refImagesJSON),
			UserID:          userID,
		}

		if err := s.repo.Create(&suggestion); err != nil {
			continue // 跳过保存失败的建议 / Skip failed saves
		}
		savedSuggestions = append(savedSuggestions, suggestion)
	}

	// 转换响应 / Convert to response
	respSuggestions := make([]response.DesignSuggestionResponse, len(savedSuggestions))
	for i, sug := range savedSuggestions {
		respSuggestions[i] = *s.toResponse(&sug)
	}

	return &response.GenerateDesignSuggestionResponse{
		Suggestions: respSuggestions,
		ProjectID:   req.ProjectID,
		Count:       len(respSuggestions),
	}, nil
}

// getProjectContext 获取项目完整上下文
// Get full project context for AI suggestion generation
func (s *DesignSuggestionService) getProjectContext(project *model.DesignerProject) (*ProjectContext, error) {
	ctx := &ProjectContext{
		Project: project,
	}

	// 获取项目材料清单 / Get project materials
	var materials []model.ProjectMaterial
	if err := s.db.Where("project_id = ?", project.ID).Find(&materials).Error; err == nil {
		ctx.Materials = materials
	}

	// 获取成本汇总 / Get cost summary
	var costEstimate model.CostEstimate
	if err := s.db.Where("project_id = ?", project.ID).First(&costEstimate).Error; err == nil {
		// 计算材料费 / Calculate material cost
		var materialCost float64
		s.db.Model(&model.ProjectMaterial{}).Where("project_id = ?", project.ID).
			Select("COALESCE(SUM(total_price), 0)").Scan(&materialCost)

		totalCost := materialCost + costEstimate.LaborCost + costEstimate.EquipmentCost + costEstimate.ManagementCost
		ctx.CostSummary = &CostSummaryInfo{
			MaterialCost:   materialCost,
			LaborCost:      costEstimate.LaborCost,
			EquipmentCost:  costEstimate.EquipmentCost,
			ManagementCost: costEstimate.ManagementCost,
			TotalCost:      totalCost,
			BudgetLimit:    costEstimate.BudgetLimit,
			BudgetExceeded: costEstimate.BudgetLimit > 0 && totalCost > costEstimate.BudgetLimit,
		}
	}

	// 获取CAD文件信息 / Get CAD files info
	if len(project.CadFiles) > 0 {
		for _, cad := range project.CadFiles {
			var layers []string
			if cad.Layers != "" {
				json.Unmarshal([]byte(cad.Layers), &layers)
			}
			ctx.CadFilesInfo = append(ctx.CadFilesInfo, CadFileInfo{
				FileName:    cad.FileName,
				FileFormat:  cad.FileFormat,
				LayerCount:  cad.LayerCount,
				Layers:      layers,
				Has3D:       cad.Has3D,
				ParseStatus: cad.ParseStatus,
			})
		}
	}

	// 获取文档信息 / Get documents info
	if len(project.Documents) > 0 {
		for _, doc := range project.Documents {
			ctx.DocumentsInfo = append(ctx.DocumentsInfo, DocumentInfo{
				FileName:       doc.FileName,
				FileType:       doc.FileType,
				Summary:        doc.Summary,
				AnalysisStatus: doc.AnalysisStatus,
			})
		}
	}

	return ctx, nil
}

// buildSystemPrompt 构建系统提示词
// Build system prompt
func (s *DesignSuggestionService) buildSystemPrompt() string {
	return `你是专业的工装设计顾问，擅长为办公空间、商业空间、工业厂房等工装项目提供设计建议。
请根据项目需求生成具体可执行的设计建议，每条建议必须包含：
1. content: 建议内容（具体、可执行的设计建议）
2. applicableScene: 适用场景（说明该建议适用于什么情况）
3. costImpact: 成本影响（如"增加约5%-10%"、"节省约15%"、"基本持平"）
4. category: 建议类别（layout-布局、material-材料、style-风格、function-功能、lighting-照明、hvac-暖通）
5. priority: 优先级（1-10，数字越大越重要）
6. detailInfo: 详细信息对象，包含实施步骤、注意事项等
7. referenceImages: 参考图片描述数组（描述建议参考的设计风格或案例）

请返回JSON数组格式，确保生成的建议数量符合要求。`
}

// buildUserPrompt 构建用户提示词（保留兼容性）
// Build user prompt (kept for compatibility)
func (s *DesignSuggestionService) buildUserPrompt(project *model.DesignerProject, category string, count int) string {
	ctx := &ProjectContext{Project: project}
	return s.buildUserPromptWithContext(ctx, category, count)
}

// buildUserPromptWithContext 基于完整上下文构建用户提示词
// Build user prompt with full project context
func (s *DesignSuggestionService) buildUserPromptWithContext(ctx *ProjectContext, category string, count int) string {
	var sb strings.Builder
	project := ctx.Project

	sb.WriteString(fmt.Sprintf("请为以下工装设计项目生成%d条设计建议：\n\n", count))
	sb.WriteString("=== 项目基本信息 ===\n")
	sb.WriteString(fmt.Sprintf("项目名称：%s\n", project.Name))
	sb.WriteString(fmt.Sprintf("项目描述：%s\n", project.Description))
	sb.WriteString(fmt.Sprintf("面积：%.2f平方米\n", project.Area))
	sb.WriteString(fmt.Sprintf("预算：%.2f元\n", project.Budget))
	sb.WriteString(fmt.Sprintf("设计风格：%s\n", project.Style))

	// 添加文档分析信息 / Add document analysis info
	if len(ctx.DocumentsInfo) > 0 {
		sb.WriteString("\n=== 项目文档分析 ===\n")
		for _, doc := range ctx.DocumentsInfo {
			if doc.AnalysisStatus == "completed" && doc.Summary != "" {
				sb.WriteString(fmt.Sprintf("- 文档【%s】(%s): %s\n", doc.FileName, doc.FileType, doc.Summary))
			}
		}
	}

	// 添加CAD文件信息 / Add CAD files info
	if len(ctx.CadFilesInfo) > 0 {
		sb.WriteString("\n=== CAD图纸信息 ===\n")
		for _, cad := range ctx.CadFilesInfo {
			sb.WriteString(fmt.Sprintf("- 图纸【%s】(格式:%s, 图层数:%d", cad.FileName, cad.FileFormat, cad.LayerCount))
			if cad.Has3D {
				sb.WriteString(", 包含3D信息")
			}
			sb.WriteString(")\n")
			if len(cad.Layers) > 0 && len(cad.Layers) <= 10 {
				sb.WriteString(fmt.Sprintf("  图层: %s\n", strings.Join(cad.Layers, ", ")))
			} else if len(cad.Layers) > 10 {
				sb.WriteString(fmt.Sprintf("  主要图层: %s 等%d个图层\n", strings.Join(cad.Layers[:10], ", "), len(cad.Layers)))
			}
		}
	}

	// 添加材料清单信息 / Add material list info
	if len(ctx.Materials) > 0 {
		sb.WriteString("\n=== 项目材料清单 ===\n")
		// 按类别分组统计 / Group by category
		categoryMap := make(map[string][]model.ProjectMaterial)
		for _, m := range ctx.Materials {
			cat := m.Category
			if cat == "" {
				cat = "其他"
			}
			categoryMap[cat] = append(categoryMap[cat], m)
		}
		for cat, materials := range categoryMap {
			var totalCost float64
			for _, m := range materials {
				totalCost += m.TotalPrice
			}
			sb.WriteString(fmt.Sprintf("- %s类材料: %d种, 合计%.2f元\n", cat, len(materials), totalCost))
			// 列出前5种材料 / List top 5 materials
			for i, m := range materials {
				if i >= 5 {
					sb.WriteString(fmt.Sprintf("  ... 等%d种材料\n", len(materials)-5))
					break
				}
				sb.WriteString(fmt.Sprintf("  · %s(%s): %.2f%s × %.2f元 = %.2f元\n",
					m.Name, m.Specification, m.Quantity, m.Unit, m.UnitPrice, m.TotalPrice))
			}
		}
	}

	// 添加成本汇总信息 / Add cost summary info
	if ctx.CostSummary != nil {
		sb.WriteString("\n=== 成本汇总 ===\n")
		sb.WriteString(fmt.Sprintf("- 材料费: %.2f元\n", ctx.CostSummary.MaterialCost))
		sb.WriteString(fmt.Sprintf("- 人工费: %.2f元\n", ctx.CostSummary.LaborCost))
		sb.WriteString(fmt.Sprintf("- 设备费: %.2f元\n", ctx.CostSummary.EquipmentCost))
		sb.WriteString(fmt.Sprintf("- 管理费: %.2f元\n", ctx.CostSummary.ManagementCost))
		sb.WriteString(fmt.Sprintf("- 总成本: %.2f元\n", ctx.CostSummary.TotalCost))
		if ctx.CostSummary.BudgetLimit > 0 {
			sb.WriteString(fmt.Sprintf("- 预算上限: %.2f元\n", ctx.CostSummary.BudgetLimit))
			if ctx.CostSummary.BudgetExceeded {
				exceeded := ctx.CostSummary.TotalCost - ctx.CostSummary.BudgetLimit
				sb.WriteString(fmt.Sprintf("⚠️ 警告：当前成本已超出预算 %.2f元，请重点关注成本优化建议！\n", exceeded))
			}
		}
	}

	// 添加类别筛选 / Add category filter
	if category != "" {
		sb.WriteString(fmt.Sprintf("\n请重点关注【%s】类别的建议。\n", s.getCategoryName(category)))
	}

	sb.WriteString("\n=== 生成要求 ===\n")
	sb.WriteString("请基于以上项目信息，结合材料清单、成本情况和CAD图纸分析，生成针对性的设计建议。\n")
	sb.WriteString("建议应考虑：\n")
	sb.WriteString("1. 材料选择的合理性和性价比\n")
	sb.WriteString("2. 空间布局与功能规划\n")
	sb.WriteString("3. 成本控制与预算优化\n")
	sb.WriteString("4. 设计风格的一致性\n")
	sb.WriteString("5. 施工可行性\n")
	sb.WriteString("\n请返回JSON数组格式的建议列表。")

	return sb.String()
}

// getCategoryName 获取类别中文名
// Get category Chinese name
func (s *DesignSuggestionService) getCategoryName(category string) string {
	names := map[string]string{
		"layout":   "空间布局",
		"material": "材料选择",
		"style":    "设计风格",
		"function": "功能规划",
		"lighting": "照明设计",
		"hvac":     "暖通空调",
	}
	if name, ok := names[category]; ok {
		return name
	}
	return category
}

// parseAIResponse 解析AI响应
// Parse AI response
func (s *DesignSuggestionService) parseAIResponse(aiResponse string) ([]AISuggestion, error) {
	// 提取JSON数组 / Extract JSON array
	jsonStr := aiResponse
	if idx := strings.Index(aiResponse, "["); idx != -1 {
		if endIdx := strings.LastIndex(aiResponse, "]"); endIdx != -1 {
			jsonStr = aiResponse[idx : endIdx+1]
		}
	}

	var suggestions []AISuggestion
	if err := json.Unmarshal([]byte(jsonStr), &suggestions); err != nil {
		// 尝试解析单个对象 / Try parsing single object
		var single AISuggestion
		if err2 := json.Unmarshal([]byte(jsonStr), &single); err2 == nil {
			suggestions = []AISuggestion{single}
		} else {
			return nil, fmt.Errorf("JSON解析失败: %w", err)
		}
	}

	// 验证并补充默认值 / Validate and set defaults
	for i := range suggestions {
		if suggestions[i].Category == "" {
			suggestions[i].Category = "general"
		}
		if suggestions[i].Priority == 0 {
			suggestions[i].Priority = 5
		}
		if suggestions[i].CostImpact == "" {
			suggestions[i].CostImpact = "待评估"
		}
	}

	return suggestions, nil
}

// GetByID 根据ID获取建议详情
// Get suggestion by ID
func (s *DesignSuggestionService) GetByID(id uint, userID uint) (*response.SuggestionDetailResponse, error) {
	suggestion, err := s.repo.GetByIDWithProject(id)
	if err != nil {
		return nil, errors.New("建议不存在")
	}

	// 验证用户权限 / Verify user permission
	if suggestion.Project != nil && suggestion.Project.UserID != userID {
		return nil, errors.New("无权访问该建议")
	}

	resp := &response.SuggestionDetailResponse{
		DesignSuggestionResponse: *s.toResponse(suggestion),
	}

	if suggestion.Project != nil {
		resp.ProjectName = suggestion.Project.Name
		resp.ProjectInfo.Area = suggestion.Project.Area
		resp.ProjectInfo.Budget = suggestion.Project.Budget
		resp.ProjectInfo.Style = suggestion.Project.Style
	}

	return resp, nil
}

// List 获取建议列表
// Get suggestion list
func (s *DesignSuggestionService) List(req *request.GetDesignSuggestionListRequest, userID uint) (*response.DesignSuggestionListResponse, error) {
	// 验证项目权限 / Verify project permission
	project, err := s.projectRepo.GetByID(req.ProjectID)
	if err != nil {
		return nil, errors.New("项目不存在")
	}
	if project.UserID != userID {
		return nil, errors.New("无权访问该项目")
	}

	filter := repository.SuggestionFilter{
		ProjectID: req.ProjectID,
		Category:  req.Category,
		Status:    req.Status,
	}

	suggestions, total, err := s.repo.List(req.Current, req.Size, filter)
	if err != nil {
		return nil, err
	}

	records := make([]response.DesignSuggestionResponse, len(suggestions))
	for i, sug := range suggestions {
		records[i] = *s.toResponse(&sug)
	}

	return &response.DesignSuggestionListResponse{
		Records: records,
		Current: req.Current,
		Size:    req.Size,
		Total:   total,
	}, nil
}

// UpdateStatus 更新建议状态
// Update suggestion status
// Requirements: 3.4
func (s *DesignSuggestionService) UpdateStatus(req *request.UpdateSuggestionStatusRequest, userID uint) error {
	suggestion, err := s.repo.GetByIDWithProject(req.ID)
	if err != nil {
		return errors.New("建议不存在")
	}

	// 验证用户权限 / Verify user permission
	if suggestion.Project != nil && suggestion.Project.UserID != userID {
		return errors.New("无权修改该建议")
	}

	// 更新状态 / Update status
	suggestion.Status = req.Status
	if err := s.repo.Update(suggestion); err != nil {
		return fmt.Errorf("更新状态失败: %w", err)
	}

	// 记录偏好 / Record preference
	pref := &model.SuggestionPreference{
		SuggestionID: req.ID,
		UserID:       userID,
		Action:       req.Status,
		Feedback:     req.Feedback,
	}
	s.repo.CreatePreference(pref)

	return nil
}

// BatchUpdateStatus 批量更新建议状态
// Batch update suggestion status
func (s *DesignSuggestionService) BatchUpdateStatus(req *request.BatchUpdateSuggestionStatusRequest, userID uint) error {
	for _, id := range req.IDs {
		singleReq := &request.UpdateSuggestionStatusRequest{
			ID:       id,
			Status:   req.Status,
			Feedback: req.Feedback,
		}
		s.UpdateStatus(singleReq, userID)
	}
	return nil
}

// Delete 删除建议
// Delete suggestion
func (s *DesignSuggestionService) Delete(id uint, userID uint) error {
	suggestion, err := s.repo.GetByIDWithProject(id)
	if err != nil {
		return errors.New("建议不存在")
	}

	// 验证用户权限 / Verify user permission
	if suggestion.Project != nil && suggestion.Project.UserID != userID {
		return errors.New("无权删除该建议")
	}

	return s.repo.Delete(id)
}

// BatchDelete 批量删除建议
// Batch delete suggestions
func (s *DesignSuggestionService) BatchDelete(ids []uint, userID uint) error {
	for _, id := range ids {
		s.Delete(id, userID)
	}
	return nil
}

// GetStats 获取建议统计
// Get suggestion statistics
func (s *DesignSuggestionService) GetStats(projectID uint, userID uint) (*response.SuggestionStatsResponse, error) {
	// 验证项目权限 / Verify project permission
	if projectID > 0 {
		project, err := s.projectRepo.GetByID(projectID)
		if err != nil {
			return nil, errors.New("项目不存在")
		}
		if project.UserID != userID {
			return nil, errors.New("无权访问该项目")
		}
	}

	stats, err := s.repo.GetStats(projectID)
	if err != nil {
		return nil, err
	}

	categoryStats, err := s.repo.GetCategoryStats(projectID)
	if err != nil {
		return nil, err
	}

	return &response.SuggestionStatsResponse{
		TotalCount:    stats["total"],
		AdoptedCount:  stats["adopted"],
		IgnoredCount:  stats["ignored"],
		PendingCount:  stats["pending"],
		CategoryStats: categoryStats,
	}, nil
}

// toResponse 转换为响应对象
// Convert to response object
func (s *DesignSuggestionService) toResponse(suggestion *model.DesignSuggestion) *response.DesignSuggestionResponse {
	var detailInfo any
	var refImages []string

	if suggestion.DetailInfo != "" {
		json.Unmarshal([]byte(suggestion.DetailInfo), &detailInfo)
	}
	if suggestion.ReferenceImages != "" {
		json.Unmarshal([]byte(suggestion.ReferenceImages), &refImages)
	}

	return &response.DesignSuggestionResponse{
		ID:              suggestion.ID,
		ProjectID:       suggestion.ProjectID,
		Content:         suggestion.Content,
		ApplicableScene: suggestion.ApplicableScene,
		CostImpact:      suggestion.CostImpact,
		Category:        suggestion.Category,
		Priority:        suggestion.Priority,
		Status:          suggestion.Status,
		DetailInfo:      detailInfo,
		ReferenceImages: refImages,
		CreatedAt:       suggestion.CreatedAt,
		UpdatedAt:       suggestion.UpdatedAt,
	}
}
