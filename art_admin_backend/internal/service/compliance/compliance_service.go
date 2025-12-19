package compliance

import (
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/dto/response"
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/volcengine"
	"art_admin_backend/internal/repository"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// ComplianceService 合规检查服务
// Compliance check service
type ComplianceService struct {
	repo        *repository.ComplianceRepository
	projectRepo *repository.DesignerProjectRepository
	docRepo     *repository.DocumentRepository
	aiClient    *volcengine.Client
}

// NewComplianceService 创建合规检查服务
// Create compliance service
func NewComplianceService(
	repo *repository.ComplianceRepository,
	projectRepo *repository.DesignerProjectRepository,
	docRepo *repository.DocumentRepository,
	aiClient *volcengine.Client,
) *ComplianceService {
	return &ComplianceService{
		repo:        repo,
		projectRepo: projectRepo,
		docRepo:     docRepo,
		aiClient:    aiClient,
	}
}

// ========== 设计规范管理方法 / Design Standard Management Methods ==========

// CreateStandard 创建设计规范
// Create design standard
func (s *ComplianceService) CreateStandard(req *request.CreateDesignStandardRequest) (*response.DesignStandardResponse, error) {
	// 检查编号是否已存在 / Check if code already exists
	existing, _ := s.repo.GetStandardByCode(req.Code)
	if existing != nil {
		return nil, errors.New("规范编号已存在")
	}

	// 验证类别 / Validate category
	if !isValidCategory(req.Category) {
		return nil, errors.New("无效的规范类别")
	}

	// 解析生效日期 / Parse effective date
	var effectiveDate *time.Time
	if req.EffectiveDate != "" {
		t, err := time.Parse("2006-01-02", req.EffectiveDate)
		if err != nil {
			return nil, errors.New("生效日期格式错误，应为YYYY-MM-DD")
		}
		effectiveDate = &t
	}

	// 序列化适用类型 / Serialize applicable types
	applicableTypesJSON, _ := json.Marshal(req.ApplicableTypes)

	standard := &model.DesignStandard{
		Code:            req.Code,
		Name:            req.Name,
		Category:        req.Category,
		Content:         req.Content,
		ApplicableTypes: string(applicableTypesJSON),
		Version:         req.Version,
		EffectiveDate:   effectiveDate,
		Source:          req.Source,
		Interpretation:  req.Interpretation,
		Status:          model.StandardStatusActive,
	}

	if err := s.repo.CreateStandard(standard); err != nil {
		return nil, fmt.Errorf("创建规范失败: %w", err)
	}

	return s.toStandardResponse(standard), nil
}

// UpdateStandard 更新设计规范
// Update design standard
func (s *ComplianceService) UpdateStandard(req *request.UpdateDesignStandardRequest) (*response.DesignStandardResponse, error) {
	standard, err := s.repo.GetStandardByID(req.ID)
	if err != nil {
		return nil, errors.New("规范不存在")
	}

	// 检查编号是否与其他规范冲突 / Check if code conflicts with other standards
	if req.Code != standard.Code {
		existing, _ := s.repo.GetStandardByCode(req.Code)
		if existing != nil && existing.ID != req.ID {
			return nil, errors.New("规范编号已被其他规范使用")
		}
	}

	// 验证类别 / Validate category
	if !isValidCategory(req.Category) {
		return nil, errors.New("无效的规范类别")
	}

	// 解析生效日期 / Parse effective date
	var effectiveDate *time.Time
	if req.EffectiveDate != "" {
		t, err := time.Parse("2006-01-02", req.EffectiveDate)
		if err != nil {
			return nil, errors.New("生效日期格式错误，应为YYYY-MM-DD")
		}
		effectiveDate = &t
	}

	// 序列化适用类型 / Serialize applicable types
	applicableTypesJSON, _ := json.Marshal(req.ApplicableTypes)

	// 更新字段 / Update fields
	standard.Code = req.Code
	standard.Name = req.Name
	standard.Category = req.Category
	standard.Content = req.Content
	standard.ApplicableTypes = string(applicableTypesJSON)
	standard.Version = req.Version
	standard.EffectiveDate = effectiveDate
	standard.Source = req.Source
	standard.Interpretation = req.Interpretation
	if req.Status != "" {
		standard.Status = req.Status
	}

	if err := s.repo.UpdateStandard(standard); err != nil {
		return nil, fmt.Errorf("更新规范失败: %w", err)
	}

	return s.toStandardResponse(standard), nil
}

// GetStandardByID 根据ID获取设计规范
// Get design standard by ID
func (s *ComplianceService) GetStandardByID(id uint) (*response.DesignStandardResponse, error) {
	standard, err := s.repo.GetStandardByID(id)
	if err != nil {
		return nil, errors.New("规范不存在")
	}
	return s.toStandardResponse(standard), nil
}

// ListStandards 获取设计规范列表
// Get design standard list
func (s *ComplianceService) ListStandards(req *request.DesignStandardListRequest) (*response.DesignStandardListResponse, error) {
	filter := repository.StandardFilter{
		Category: req.Category,
		Status:   req.Status,
		Keyword:  req.Keyword,
	}

	standards, total, err := s.repo.ListStandards(req.Current, req.Size, filter)
	if err != nil {
		return nil, err
	}

	records := make([]response.DesignStandardResponse, len(standards))
	for i, standard := range standards {
		records[i] = *s.toStandardResponse(&standard)
	}

	return &response.DesignStandardListResponse{
		Records: records,
		Current: req.Current,
		Size:    req.Size,
		Total:   total,
	}, nil
}

// DeleteStandard 删除设计规范
// Delete design standard
func (s *ComplianceService) DeleteStandard(id uint) error {
	_, err := s.repo.GetStandardByID(id)
	if err != nil {
		return errors.New("规范不存在")
	}
	return s.repo.DeleteStandard(id)
}

// BatchDeleteStandards 批量删除设计规范
// Batch delete design standards
func (s *ComplianceService) BatchDeleteStandards(ids []uint) error {
	return s.repo.BatchDeleteStandards(ids)
}

// GetStandardStats 获取规范统计信息
// Get standard statistics
func (s *ComplianceService) GetStandardStats() (*response.DesignStandardStats, error) {
	total, active, err := s.repo.GetStandardStats()
	if err != nil {
		return nil, err
	}

	categoryStats, err := s.repo.GetStandardCategoryStats()
	if err != nil {
		return nil, err
	}

	stats := make([]response.StandardCategoryStats, len(categoryStats))
	for i, cs := range categoryStats {
		stats[i] = response.StandardCategoryStats{
			Category: cs.Category,
			Count:    cs.Count,
		}
	}

	return &response.DesignStandardStats{
		TotalCount:    total,
		ActiveCount:   active,
		CategoryStats: stats,
	}, nil
}

// toStandardResponse 转换为规范响应对象
// Convert to standard response object
func (s *ComplianceService) toStandardResponse(standard *model.DesignStandard) *response.DesignStandardResponse {
	var applicableTypes []string
	if standard.ApplicableTypes != "" {
		json.Unmarshal([]byte(standard.ApplicableTypes), &applicableTypes)
	}
	if applicableTypes == nil {
		applicableTypes = []string{}
	}

	effectiveDate := ""
	if standard.EffectiveDate != nil {
		effectiveDate = standard.EffectiveDate.Format("2006-01-02")
	}

	return &response.DesignStandardResponse{
		ID:              standard.ID,
		Code:            standard.Code,
		Name:            standard.Name,
		Category:        standard.Category,
		Content:         standard.Content,
		ApplicableTypes: applicableTypes,
		Version:         standard.Version,
		EffectiveDate:   effectiveDate,
		Source:          standard.Source,
		Interpretation:  standard.Interpretation,
		Status:          standard.Status,
		CreatedAt:       standard.CreatedAt,
		UpdatedAt:       standard.UpdatedAt,
	}
}

// isValidCategory 验证规范类别
// Validate standard category
func isValidCategory(category string) bool {
	validCategories := map[string]bool{
		model.StandardCategoryFire:          true,
		model.StandardCategoryAccessibility: true,
		model.StandardCategoryEnvironmental: true,
		model.StandardCategorySafety:        true,
		model.StandardCategoryOther:         true,
	}
	return validCategories[category]
}

// ========== 合规检查方法 / Compliance Check Methods ==========

// PerformComplianceCheck 执行合规检查
// Perform compliance check
// Requirements: 11.1, 11.2, 11.3
func (s *ComplianceService) PerformComplianceCheck(req *request.ComplianceCheckRequest, userID uint) (*response.ComplianceCheckResponse, error) {
	// 验证项目存在 / Verify project exists
	project, err := s.projectRepo.GetByID(req.ProjectID)
	if err != nil {
		return nil, errors.New("项目不存在")
	}

	// 验证项目权限 / Verify project permission
	if project.UserID != userID {
		return nil, errors.New("无权访问该项目")
	}

	// 确定检查类型 / Determine check type
	checkType := model.CheckTypeFull
	if req.CheckType != "" {
		checkType = req.CheckType
	}

	// 创建检查记录 / Create check record
	checkResult := &model.ComplianceCheckResult{
		ProjectID:   req.ProjectID,
		CheckType:   checkType,
		CheckStatus: model.CheckStatusProcessing,
		UserID:      userID,
	}

	if err := s.repo.CreateCheckResult(checkResult); err != nil {
		return nil, fmt.Errorf("创建检查记录失败: %w", err)
	}

	// 获取适用的规范 / Get applicable standards
	var standards []model.DesignStandard
	if len(req.Categories) > 0 {
		standards, err = s.repo.GetStandardsByCategories(req.Categories)
	} else {
		standards, err = s.repo.GetAllActiveStandards()
	}
	if err != nil {
		s.updateCheckResultFailed(checkResult.ID, "获取规范失败")
		return nil, fmt.Errorf("获取规范失败: %w", err)
	}

	if len(standards) == 0 {
		s.updateCheckResultFailed(checkResult.ID, "没有可用的规范进行检查")
		return nil, errors.New("没有可用的规范进行检查")
	}

	// 获取项目文档信息 / Get project document info
	projectInfo := s.buildProjectInfo(project)

	// 使用AI进行合规检查 / Perform compliance check with AI
	result, err := s.checkComplianceWithAI(projectInfo, standards)
	if err != nil {
		s.updateCheckResultFailed(checkResult.ID, err.Error())
		return nil, fmt.Errorf("AI合规检查失败: %w", err)
	}

	// 保存检查结果 / Save check result
	passedItemsJSON, _ := json.Marshal(result.PassedItems)
	failedItemsJSON, _ := json.Marshal(result.FailedItems)
	suggestionsJSON, _ := json.Marshal(result.Suggestions)

	standardIDs := make([]uint, len(standards))
	for i, std := range standards {
		standardIDs[i] = std.ID
	}
	checkedStandardsJSON, _ := json.Marshal(standardIDs)

	updates := map[string]any{
		"check_status":      model.CheckStatusCompleted,
		"passed_items":      string(passedItemsJSON),
		"failed_items":      string(failedItemsJSON),
		"suggestions":       string(suggestionsJSON),
		"overall_score":     result.OverallScore,
		"checked_standards": string(checkedStandardsJSON),
		"error_message":     "",
	}

	if err := s.repo.UpdateCheckResultStatus(checkResult.ID, updates); err != nil {
		return nil, fmt.Errorf("保存检查结果失败: %w", err)
	}

	// 返回完整响应 / Return complete response
	return &response.ComplianceCheckResponse{
		ID:               checkResult.ID,
		ProjectID:        req.ProjectID,
		ProjectName:      project.Name,
		CheckType:        checkType,
		PassedItems:      result.PassedItems,
		FailedItems:      result.FailedItems,
		Suggestions:      result.Suggestions,
		OverallScore:     result.OverallScore,
		CheckStatus:      model.CheckStatusCompleted,
		CheckedStandards: standardIDs,
		UserID:           userID,
		CreatedAt:        checkResult.CreatedAt,
		UpdatedAt:        time.Now(),
	}, nil
}

// buildProjectInfo 构建项目信息文本
// Build project info text for AI analysis
func (s *ComplianceService) buildProjectInfo(project *model.DesignerProject) string {
	var parts []string

	parts = append(parts, fmt.Sprintf("项目名称: %s", project.Name))
	if project.Description != "" {
		parts = append(parts, fmt.Sprintf("项目描述: %s", project.Description))
	}
	if project.Area > 0 {
		parts = append(parts, fmt.Sprintf("面积: %.2f平方米", project.Area))
	}
	if project.Budget > 0 {
		parts = append(parts, fmt.Sprintf("预算: %.2f元", project.Budget))
	}
	if project.Style != "" {
		parts = append(parts, fmt.Sprintf("设计风格: %s", project.Style))
	}

	// 获取项目文档摘要 / Get project document summaries
	docs, _ := s.docRepo.GetByProjectID(project.ID)
	if len(docs) > 0 {
		var docSummaries []string
		for _, doc := range docs {
			if doc.Summary != "" {
				docSummaries = append(docSummaries, doc.Summary)
			}
		}
		if len(docSummaries) > 0 {
			parts = append(parts, fmt.Sprintf("文档摘要: %s", strings.Join(docSummaries, "; ")))
		}
	}

	return strings.Join(parts, "\n")
}

// checkComplianceWithAI 使用AI进行合规检查
// Check compliance with AI
func (s *ComplianceService) checkComplianceWithAI(projectInfo string, standards []model.DesignStandard) (*ComplianceCheckAIResult, error) {
	// 构建规范列表文本 / Build standards list text
	var standardsText []string
	for _, std := range standards {
		standardsText = append(standardsText, fmt.Sprintf(
			"[%s] %s (%s): %s",
			std.Code, std.Name, getCategoryName(std.Category), std.Content,
		))
	}

	systemPrompt := `你是专业的工装设计合规检查专家。请根据项目信息和设计规范，进行合规性检查。

检查要求：
1. 逐条对照规范，判断项目是否符合要求
2. 对于通过的项目，说明符合原因
3. 对于不合规的项目，明确指出问题位置、违规内容、原始规范条款、严重程度和修改建议
4. 提供整体改进建议

返回JSON格式：
{
  "passedItems": [
    {"standardId": 1, "standardCode": "GB50016-2014", "standardName": "规范名称", "category": "fire", "description": "符合说明"}
  ],
  "failedItems": [
    {"standardId": 2, "standardCode": "GB50016-5.5.17", "standardName": "规范名称", "category": "fire", "violationContent": "违规内容", "originalRequirement": "原始规范条款", "severity": "high", "location": "问题位置", "suggestion": "修改建议"}
  ],
  "suggestions": [
    {"content": "建议内容", "priority": "high", "category": "fire", "costImpact": "成本影响", "reference": "参考规范"}
  ],
  "overallScore": 75.5
}

严重程度(severity): low/medium/high/critical
优先级(priority): low/medium/high
类别(category): fire/accessibility/environmental/safety/other`

	userPrompt := fmt.Sprintf(`请对以下项目进行合规检查：

【项目信息】
%s

【适用规范】
%s

请逐条检查并返回JSON格式的检查结果。`, projectInfo, strings.Join(standardsText, "\n"))

	// 使用ChatWithText方法发送纯文本请求 / Use ChatWithText method for text-only request
	aiResponse, err := s.aiClient.ChatWithText(systemPrompt, userPrompt)
	if err != nil {
		return nil, err
	}

	return s.parseComplianceAIResponse(aiResponse, standards)
}

// ComplianceCheckAIResult AI合规检查结果
// AI compliance check result
type ComplianceCheckAIResult struct {
	PassedItems  []response.PassedItemResponse
	FailedItems  []response.FailedItemResponse
	Suggestions  []response.ComplianceSuggestionResponse
	OverallScore float64
}

// parseComplianceAIResponse 解析AI合规检查响应
// Parse AI compliance check response
func (s *ComplianceService) parseComplianceAIResponse(aiResponse string, standards []model.DesignStandard) (*ComplianceCheckAIResult, error) {
	// 提取JSON部分 / Extract JSON part
	jsonStr := aiResponse
	if idx := strings.Index(aiResponse, "{"); idx != -1 {
		if endIdx := strings.LastIndex(aiResponse, "}"); endIdx != -1 {
			jsonStr = aiResponse[idx : endIdx+1]
		}
	}

	var result struct {
		PassedItems  []response.PassedItemResponse           `json:"passedItems"`
		FailedItems  []response.FailedItemResponse           `json:"failedItems"`
		Suggestions  []response.ComplianceSuggestionResponse `json:"suggestions"`
		OverallScore float64                                 `json:"overallScore"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		// 解析失败时返回默认结果 / Return default result on parse failure
		return s.generateDefaultCheckResult(standards), nil
	}

	// 确保返回非nil的切片 / Ensure non-nil slices
	if result.PassedItems == nil {
		result.PassedItems = []response.PassedItemResponse{}
	}
	if result.FailedItems == nil {
		result.FailedItems = []response.FailedItemResponse{}
	}
	if result.Suggestions == nil {
		result.Suggestions = []response.ComplianceSuggestionResponse{}
	}

	return &ComplianceCheckAIResult{
		PassedItems:  result.PassedItems,
		FailedItems:  result.FailedItems,
		Suggestions:  result.Suggestions,
		OverallScore: result.OverallScore,
	}, nil
}

// generateDefaultCheckResult 生成默认检查结果
// Generate default check result when AI parsing fails
func (s *ComplianceService) generateDefaultCheckResult(standards []model.DesignStandard) *ComplianceCheckAIResult {
	// 将所有规范标记为待检查 / Mark all standards as pending check
	suggestions := []response.ComplianceSuggestionResponse{
		{
			Content:    "AI分析结果解析失败，建议人工复核各项规范要求",
			Priority:   "high",
			Category:   "other",
			CostImpact: "无",
			Reference:  "所有适用规范",
		},
	}

	return &ComplianceCheckAIResult{
		PassedItems:  []response.PassedItemResponse{},
		FailedItems:  []response.FailedItemResponse{},
		Suggestions:  suggestions,
		OverallScore: 0,
	}
}

// updateCheckResultFailed 更新检查结果为失败
// Update check result to failed status
func (s *ComplianceService) updateCheckResultFailed(id uint, errorMsg string) {
	s.repo.UpdateCheckResultStatus(id, map[string]any{
		"check_status":  model.CheckStatusFailed,
		"error_message": errorMsg,
	})
}

// getCategoryName 获取类别中文名称
// Get category Chinese name
func getCategoryName(category string) string {
	names := map[string]string{
		model.StandardCategoryFire:          "消防规范",
		model.StandardCategoryAccessibility: "无障碍规范",
		model.StandardCategoryEnvironmental: "环保规范",
		model.StandardCategorySafety:        "安全规范",
		model.StandardCategoryOther:         "其他规范",
	}
	if name, ok := names[category]; ok {
		return name
	}
	return category
}

// ========== 合规检查结果查询方法 / Compliance Check Result Query Methods ==========

// GetCheckResultByID 根据ID获取合规检查结果
// Get compliance check result by ID
func (s *ComplianceService) GetCheckResultByID(id, userID uint) (*response.ComplianceCheckResponse, error) {
	result, err := s.repo.GetCheckResultByIDWithProject(id)
	if err != nil {
		return nil, errors.New("检查结果不存在")
	}

	// 验证权限 / Verify permission
	if result.Project != nil && result.Project.UserID != userID {
		return nil, errors.New("无权访问该检查结果")
	}

	return s.toCheckResultResponse(result), nil
}

// ListCheckResults 获取合规检查结果列表
// Get compliance check result list
func (s *ComplianceService) ListCheckResults(req *request.ComplianceCheckListRequest, userID uint) (*response.ComplianceCheckListResponse, error) {
	// 验证项目权限 / Verify project permission
	project, err := s.projectRepo.GetByID(req.ProjectID)
	if err != nil {
		return nil, errors.New("项目不存在")
	}
	if project.UserID != userID {
		return nil, errors.New("无权访问该项目")
	}

	filter := repository.CheckResultFilter{
		ProjectID:   req.ProjectID,
		CheckStatus: req.CheckStatus,
	}

	results, total, err := s.repo.ListCheckResults(req.Current, req.Size, filter)
	if err != nil {
		return nil, err
	}

	records := make([]response.ComplianceCheckResponse, len(results))
	for i, result := range results {
		records[i] = *s.toCheckResultResponse(&result)
	}

	return &response.ComplianceCheckListResponse{
		Records: records,
		Current: req.Current,
		Size:    req.Size,
		Total:   total,
	}, nil
}

// GetLatestCheckResult 获取项目最新的合规检查结果
// Get latest compliance check result for project
func (s *ComplianceService) GetLatestCheckResult(projectID, userID uint) (*response.ComplianceCheckResponse, error) {
	// 验证项目权限 / Verify project permission
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return nil, errors.New("项目不存在")
	}
	if project.UserID != userID {
		return nil, errors.New("无权访问该项目")
	}

	result, err := s.repo.GetLatestCheckResultByProjectID(projectID)
	if err != nil {
		return nil, errors.New("暂无检查结果")
	}

	resp := s.toCheckResultResponse(result)
	resp.ProjectName = project.Name
	return resp, nil
}

// GetCheckSummary 获取项目合规检查摘要
// Get compliance check summary for project
func (s *ComplianceService) GetCheckSummary(projectID, userID uint) (*response.ComplianceCheckSummary, error) {
	// 验证项目权限 / Verify project permission
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return nil, errors.New("项目不存在")
	}
	if project.UserID != userID {
		return nil, errors.New("无权访问该项目")
	}

	// 获取检查次数 / Get check count
	totalChecks, err := s.repo.GetCheckResultCountByProjectID(projectID)
	if err != nil {
		return nil, err
	}

	// 获取最新检查结果 / Get latest check result
	var passedCount, failedCount int
	var averageScore float64
	var latestCheckAt string

	latestResult, err := s.repo.GetLatestCheckResultByProjectID(projectID)
	if err == nil {
		var passedItems []response.PassedItemResponse
		var failedItems []response.FailedItemResponse
		if latestResult.PassedItems != "" {
			json.Unmarshal([]byte(latestResult.PassedItems), &passedItems)
		}
		if latestResult.FailedItems != "" {
			json.Unmarshal([]byte(latestResult.FailedItems), &failedItems)
		}
		passedCount = len(passedItems)
		failedCount = len(failedItems)
		averageScore = latestResult.OverallScore
		latestCheckAt = latestResult.CreatedAt.Format("2006-01-02 15:04:05")
	}

	return &response.ComplianceCheckSummary{
		TotalChecks:   totalChecks,
		PassedCount:   passedCount,
		FailedCount:   failedCount,
		AverageScore:  averageScore,
		LatestCheckAt: latestCheckAt,
	}, nil
}

// DeleteCheckResult 删除合规检查结果
// Delete compliance check result
func (s *ComplianceService) DeleteCheckResult(id, userID uint) error {
	result, err := s.repo.GetCheckResultByIDWithProject(id)
	if err != nil {
		return errors.New("检查结果不存在")
	}

	if result.Project != nil && result.Project.UserID != userID {
		return errors.New("无权删除该检查结果")
	}

	return s.repo.DeleteCheckResult(id)
}

// BatchDeleteCheckResults 批量删除合规检查结果
// Batch delete compliance check results
func (s *ComplianceService) BatchDeleteCheckResults(ids []uint, userID uint) error {
	for _, id := range ids {
		s.DeleteCheckResult(id, userID)
	}
	return nil
}

// toCheckResultResponse 转换为检查结果响应对象
// Convert to check result response object
func (s *ComplianceService) toCheckResultResponse(result *model.ComplianceCheckResult) *response.ComplianceCheckResponse {
	var passedItems []response.PassedItemResponse
	var failedItems []response.FailedItemResponse
	var suggestions []response.ComplianceSuggestionResponse
	var checkedStandards []uint

	if result.PassedItems != "" {
		json.Unmarshal([]byte(result.PassedItems), &passedItems)
	}
	if result.FailedItems != "" {
		json.Unmarshal([]byte(result.FailedItems), &failedItems)
	}
	if result.Suggestions != "" {
		json.Unmarshal([]byte(result.Suggestions), &suggestions)
	}
	if result.CheckedStandards != "" {
		json.Unmarshal([]byte(result.CheckedStandards), &checkedStandards)
	}

	// 确保非nil / Ensure non-nil
	if passedItems == nil {
		passedItems = []response.PassedItemResponse{}
	}
	if failedItems == nil {
		failedItems = []response.FailedItemResponse{}
	}
	if suggestions == nil {
		suggestions = []response.ComplianceSuggestionResponse{}
	}
	if checkedStandards == nil {
		checkedStandards = []uint{}
	}

	projectName := ""
	if result.Project != nil {
		projectName = result.Project.Name
	}

	return &response.ComplianceCheckResponse{
		ID:               result.ID,
		ProjectID:        result.ProjectID,
		ProjectName:      projectName,
		CheckType:        result.CheckType,
		PassedItems:      passedItems,
		FailedItems:      failedItems,
		Suggestions:      suggestions,
		OverallScore:     result.OverallScore,
		CheckStatus:      result.CheckStatus,
		ErrorMessage:     result.ErrorMessage,
		CheckedStandards: checkedStandards,
		UserID:           result.UserID,
		CreatedAt:        result.CreatedAt,
		UpdatedAt:        result.UpdatedAt,
	}
}
