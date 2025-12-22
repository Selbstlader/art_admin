package design_version

import (
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/dto/response"
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/volcengine"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"gorm.io/gorm"
)

// DesignVersionService 设计版本服务
// Design version service
type DesignVersionService struct {
	db       *gorm.DB
	aiClient *volcengine.Client // 火山引擎AI客户端 / VolcEngine AI client
}

// NewDesignVersionService 创建设计版本服务实例
// Create design version service instance
func NewDesignVersionService(db *gorm.DB) *DesignVersionService {
	return &DesignVersionService{db: db}
}

// SetAIClient 设置AI客户端
// Set AI client for version analysis
func (s *DesignVersionService) SetAIClient(client *volcengine.Client) {
	s.aiClient = client
}

// Create 创建设计版本
// Create design version
func (s *DesignVersionService) Create(req *request.CreateDesignVersionRequest, userID uint) (*response.DesignVersionResponse, error) {
	// 验证项目是否存在 / Verify project exists
	var project model.DesignerProject
	if err := s.db.First(&project, req.ProjectID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("项目不存在")
		}
		return nil, err
	}

	// 获取下一个版本号 / Get next version number
	var maxVersion int
	s.db.Model(&model.DesignVersion{}).
		Where("project_id = ?", req.ProjectID).
		Select("COALESCE(MAX(version_number), 0)").
		Scan(&maxVersion)
	nextVersion := maxVersion + 1

	// 序列化设计图列表 / Serialize design images
	var designImagesJSON *string
	if len(req.DesignImages) > 0 {
		images := make([]model.DesignImage, len(req.DesignImages))
		for i, url := range req.DesignImages {
			images[i] = model.DesignImage{
				URL:  url,
				Name: fmt.Sprintf("设计图%d", i+1),
				Type: "design",
			}
		}
		if data, err := json.Marshal(images); err == nil {
			str := string(data)
			designImagesJSON = &str
		}
	}

	// 序列化CAD文件ID列表 / Serialize CAD file IDs
	var cadFileIDsJSON *string
	if len(req.CadFileIDs) > 0 {
		if data, err := json.Marshal(req.CadFileIDs); err == nil {
			str := string(data)
			cadFileIDsJSON = &str
		}
	}

	// 生成版本名称 / Generate version name
	versionName := req.VersionName
	if versionName == "" {
		versionName = fmt.Sprintf("版本%d", nextVersion)
	}

	// 创建版本 / Create version
	version := &model.DesignVersion{
		ProjectID:     req.ProjectID,
		VersionNumber: nextVersion,
		VersionName:   versionName,
		Description:   req.Description,
		DesignImages:  designImagesJSON,
		CadFileIDs:    cadFileIDsJSON,
		LayoutInfo:    strPtr(req.LayoutInfo),
		AreaInfo:      strPtr(req.AreaInfo),
		StyleInfo:     strPtr(req.StyleInfo),
		MaterialInfo:  strPtr(req.MaterialInfo),
		Status:        "draft",
		CreatedBy:     userID,
	}

	if err := s.db.Create(version).Error; err != nil {
		return nil, err
	}

	return s.toResponse(version), nil
}

// GetByID 根据ID获取版本详情
// Get version detail by ID
func (s *DesignVersionService) GetByID(id uint) (*response.DesignVersionResponse, error) {
	var version model.DesignVersion
	if err := s.db.First(&version, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("版本不存在")
		}
		return nil, err
	}
	return s.toResponse(&version), nil
}

// List 获取版本列表
// Get version list
func (s *DesignVersionService) List(req *request.DesignVersionListRequest) (*response.DesignVersionListResponse, error) {
	var versions []model.DesignVersion
	var total int64

	query := s.db.Model(&model.DesignVersion{}).Where("project_id = ?", req.ProjectID)

	// 状态筛选 / Status filter
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	// 统计总数 / Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询 / Paginated query
	offset := (req.Current - 1) * req.Size
	if err := query.Order("version_number DESC").
		Offset(offset).Limit(req.Size).
		Find(&versions).Error; err != nil {
		return nil, err
	}

	// 转换响应 / Convert to response
	records := make([]response.DesignVersionResponse, len(versions))
	for i, v := range versions {
		records[i] = *s.toResponse(&v)
	}

	return &response.DesignVersionListResponse{
		Records: records,
		Current: req.Current,
		Size:    req.Size,
		Total:   total,
	}, nil
}

// Update 更新版本
// Update version
func (s *DesignVersionService) Update(req *request.UpdateDesignVersionRequest) (*response.DesignVersionResponse, error) {
	var version model.DesignVersion
	if err := s.db.First(&version, req.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("版本不存在")
		}
		return nil, err
	}

	// 更新字段 / Update fields
	if req.VersionName != "" {
		version.VersionName = req.VersionName
	}
	if req.Description != "" {
		version.Description = req.Description
	}
	if len(req.DesignImages) > 0 {
		images := make([]model.DesignImage, len(req.DesignImages))
		for i, url := range req.DesignImages {
			images[i] = model.DesignImage{
				URL:  url,
				Name: fmt.Sprintf("设计图%d", i+1),
				Type: "design",
			}
		}
		if data, err := json.Marshal(images); err == nil {
			str := string(data)
			version.DesignImages = &str
		}
	}
	if len(req.CadFileIDs) > 0 {
		if data, err := json.Marshal(req.CadFileIDs); err == nil {
			str := string(data)
			version.CadFileIDs = &str
		}
	}
	if req.LayoutInfo != "" {
		version.LayoutInfo = strPtr(req.LayoutInfo)
	}
	if req.AreaInfo != "" {
		version.AreaInfo = strPtr(req.AreaInfo)
	}
	if req.StyleInfo != "" {
		version.StyleInfo = strPtr(req.StyleInfo)
	}
	if req.MaterialInfo != "" {
		version.MaterialInfo = strPtr(req.MaterialInfo)
	}
	if req.Status != "" {
		version.Status = req.Status
	}

	if err := s.db.Save(&version).Error; err != nil {
		return nil, err
	}

	return s.toResponse(&version), nil
}

// Delete 删除版本
// Delete version
func (s *DesignVersionService) Delete(id uint) error {
	result := s.db.Delete(&model.DesignVersion{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("版本不存在")
	}
	return nil
}

// Compare 对比两个版本
// Compare two versions
func (s *DesignVersionService) Compare(req *request.CompareDesignVersionsRequest) (*response.VersionCompareResponse, error) {
	// 获取两个版本 / Get both versions
	var versionA, versionB model.DesignVersion
	if err := s.db.First(&versionA, req.VersionAID).Error; err != nil {
		return nil, errors.New("版本A不存在")
	}
	if err := s.db.First(&versionB, req.VersionBID).Error; err != nil {
		return nil, errors.New("版本B不存在")
	}

	// 验证属于同一项目 / Verify same project
	if versionA.ProjectID != versionB.ProjectID || versionA.ProjectID != req.ProjectID {
		return nil, errors.New("版本必须属于同一项目")
	}

	// 计算差异 / Calculate differences
	layoutChanges := s.compareJSON(ptrToStr(versionA.LayoutInfo), ptrToStr(versionB.LayoutInfo), "布局")
	areaChanges := s.compareJSON(ptrToStr(versionA.AreaInfo), ptrToStr(versionB.AreaInfo), "面积")
	styleChanges := s.compareJSON(ptrToStr(versionA.StyleInfo), ptrToStr(versionB.StyleInfo), "风格")
	materialChanges := s.compareJSON(ptrToStr(versionA.MaterialInfo), ptrToStr(versionB.MaterialInfo), "材料")

	// 生成摘要 / Generate summary
	totalChanges := len(layoutChanges) + len(areaChanges) + len(styleChanges) + len(materialChanges)
	summary := fmt.Sprintf("版本%d与版本%d对比：共发现%d处变化", versionA.VersionNumber, versionB.VersionNumber, totalChanges)

	// 保存对比记录 / Save compare record
	compare := &model.DesignVersionCompare{
		ProjectID:       req.ProjectID,
		VersionAID:      req.VersionAID,
		VersionBID:      req.VersionBID,
		LayoutChanges:   toJSON(layoutChanges),
		AreaChanges:     toJSON(areaChanges),
		StyleChanges:    toJSON(styleChanges),
		MaterialChanges: toJSON(materialChanges),
		Summary:         summary,
		CompareStatus:   "completed",
	}
	s.db.Create(compare)

	return &response.VersionCompareResponse{
		ID:         compare.ID,
		ProjectID:  req.ProjectID,
		VersionAID: req.VersionAID,
		VersionBID: req.VersionBID,
		VersionA: &response.DesignVersionBriefResponse{
			ID:            versionA.ID,
			VersionNumber: versionA.VersionNumber,
			VersionName:   versionA.VersionName,
			Status:        versionA.Status,
		},
		VersionB: &response.DesignVersionBriefResponse{
			ID:            versionB.ID,
			VersionNumber: versionB.VersionNumber,
			VersionName:   versionB.VersionName,
			Status:        versionB.Status,
		},
		LayoutChanges:   layoutChanges,
		AreaChanges:     areaChanges,
		StyleChanges:    styleChanges,
		MaterialChanges: materialChanges,
		Summary:         summary,
		CompareStatus:   "completed",
		CreatedAt:       compare.CreatedAt,
		UpdatedAt:       compare.UpdatedAt,
	}, nil
}

// GetDiff 获取版本差异（用于左右分屏展示）
// Get version diff (for side-by-side display)
func (s *DesignVersionService) GetDiff(versionAID, versionBID uint) (*response.VersionDiffResponse, error) {
	versionA, err := s.GetByID(versionAID)
	if err != nil {
		return nil, errors.New("版本A不存在")
	}
	versionB, err := s.GetByID(versionBID)
	if err != nil {
		return nil, errors.New("版本B不存在")
	}

	// 获取原始数据计算差异 / Get raw data for diff calculation
	var modelA, modelB model.DesignVersion
	s.db.First(&modelA, versionAID)
	s.db.First(&modelB, versionBID)

	layoutChanges := s.compareJSON(ptrToStr(modelA.LayoutInfo), ptrToStr(modelB.LayoutInfo), "布局")
	areaChanges := s.compareJSON(ptrToStr(modelA.AreaInfo), ptrToStr(modelB.AreaInfo), "面积")
	styleChanges := s.compareJSON(ptrToStr(modelA.StyleInfo), ptrToStr(modelB.StyleInfo), "风格")
	materialChanges := s.compareJSON(ptrToStr(modelA.MaterialInfo), ptrToStr(modelB.MaterialInfo), "材料")

	totalChanges := len(layoutChanges) + len(areaChanges) + len(styleChanges) + len(materialChanges)
	summary := fmt.Sprintf("版本%d与版本%d对比：共发现%d处变化", modelA.VersionNumber, modelB.VersionNumber, totalChanges)

	return &response.VersionDiffResponse{
		VersionA:        versionA,
		VersionB:        versionB,
		LayoutChanges:   layoutChanges,
		AreaChanges:     areaChanges,
		StyleChanges:    styleChanges,
		MaterialChanges: materialChanges,
		Summary:         summary,
	}, nil
}

// GetCompareList 获取对比记录列表
// Get compare record list
func (s *DesignVersionService) GetCompareList(req *request.VersionCompareListRequest) (*response.VersionCompareListResponse, error) {
	var compares []model.DesignVersionCompare
	var total int64

	query := s.db.Model(&model.DesignVersionCompare{}).Where("project_id = ?", req.ProjectID)

	if req.CompareStatus != "" {
		query = query.Where("compare_status = ?", req.CompareStatus)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (req.Current - 1) * req.Size
	if err := query.Preload("VersionA").Preload("VersionB").
		Order("created_at DESC").
		Offset(offset).Limit(req.Size).
		Find(&compares).Error; err != nil {
		return nil, err
	}

	records := make([]response.VersionCompareResponse, len(compares))
	for i, c := range compares {
		records[i] = *s.toCompareResponse(&c)
	}

	return &response.VersionCompareListResponse{
		Records: records,
		Current: req.Current,
		Size:    req.Size,
		Total:   total,
	}, nil
}

// GetCompareDetail 获取对比详情
// Get compare detail
func (s *DesignVersionService) GetCompareDetail(id uint) (*response.VersionCompareResponse, error) {
	var compare model.DesignVersionCompare
	if err := s.db.Preload("VersionA").Preload("VersionB").First(&compare, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("对比记录不存在")
		}
		return nil, err
	}
	return s.toCompareResponse(&compare), nil
}

// DeleteCompare 删除对比记录
// Delete compare record
func (s *DesignVersionService) DeleteCompare(id uint) error {
	result := s.db.Delete(&model.DesignVersionCompare{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("对比记录不存在")
	}
	return nil
}

// toResponse 转换为响应结构
// Convert to response structure
func (s *DesignVersionService) toResponse(v *model.DesignVersion) *response.DesignVersionResponse {
	resp := &response.DesignVersionResponse{
		ID:            v.ID,
		ProjectID:     v.ProjectID,
		VersionNumber: v.VersionNumber,
		VersionName:   v.VersionName,
		Description:   v.Description,
		Status:        v.Status,
		CreatedBy:     v.CreatedBy,
		CreatedAt:     v.CreatedAt,
		UpdatedAt:     v.UpdatedAt,
	}

	// 解析设计图 / Parse design images
	designImagesStr := ptrToStr(v.DesignImages)
	if designImagesStr != "" && designImagesStr != "[]" {
		var images []response.DesignImageResponse
		json.Unmarshal([]byte(designImagesStr), &images)
		resp.DesignImages = images
	}

	// 解析CAD文件ID / Parse CAD file IDs
	cadFileIDsStr := ptrToStr(v.CadFileIDs)
	if cadFileIDsStr != "" && cadFileIDsStr != "[]" {
		var cadFileIDs []uint
		json.Unmarshal([]byte(cadFileIDsStr), &cadFileIDs)
		resp.CadFileIDs = cadFileIDs

		// 获取CAD文件信息 / Get CAD file info
		if len(cadFileIDs) > 0 {
			var cadFiles []model.CadFile
			s.db.Where("id IN ?", cadFileIDs).Find(&cadFiles)
			resp.CadFiles = make([]response.CadFileBriefResponse, len(cadFiles))
			for i, f := range cadFiles {
				resp.CadFiles[i] = response.CadFileBriefResponse{
					ID:            f.ID,
					FileName:      f.FileName,
					FileFormat:    f.FileFormat,
					FilePath:      f.FilePath,
					IsAIGenerated: f.IsAIGenerated,
				}
			}
		}
	}

	// 解析其他JSON字段 / Parse other JSON fields
	layoutInfoStr := ptrToStr(v.LayoutInfo)
	if layoutInfoStr != "" {
		var layout []response.LayoutInfoResponse
		json.Unmarshal([]byte(layoutInfoStr), &layout)
		resp.LayoutInfo = layout
	}
	areaInfoStr := ptrToStr(v.AreaInfo)
	if areaInfoStr != "" {
		var area []response.AreaInfoResponse
		json.Unmarshal([]byte(areaInfoStr), &area)
		resp.AreaInfo = area
	}
	styleInfoStr := ptrToStr(v.StyleInfo)
	if styleInfoStr != "" {
		var style []response.StyleInfoResponse
		json.Unmarshal([]byte(styleInfoStr), &style)
		resp.StyleInfo = style
	}
	materialInfoStr := ptrToStr(v.MaterialInfo)
	if materialInfoStr != "" {
		var material []response.MaterialInfoResponse
		json.Unmarshal([]byte(materialInfoStr), &material)
		resp.MaterialInfo = material
	}

	return resp
}

// toCompareResponse 转换对比记录为响应
// Convert compare record to response
func (s *DesignVersionService) toCompareResponse(c *model.DesignVersionCompare) *response.VersionCompareResponse {
	resp := &response.VersionCompareResponse{
		ID:            c.ID,
		ProjectID:     c.ProjectID,
		VersionAID:    c.VersionAID,
		VersionBID:    c.VersionBID,
		Summary:       c.Summary,
		CompareStatus: c.CompareStatus,
		CreatedAt:     c.CreatedAt,
		UpdatedAt:     c.UpdatedAt,
	}

	if c.VersionA != nil {
		resp.VersionA = &response.DesignVersionBriefResponse{
			ID:            c.VersionA.ID,
			VersionNumber: c.VersionA.VersionNumber,
			VersionName:   c.VersionA.VersionName,
			Status:        c.VersionA.Status,
		}
	}
	if c.VersionB != nil {
		resp.VersionB = &response.DesignVersionBriefResponse{
			ID:            c.VersionB.ID,
			VersionNumber: c.VersionB.VersionNumber,
			VersionName:   c.VersionB.VersionName,
			Status:        c.VersionB.Status,
		}
	}

	// 解析变化列表 / Parse change lists
	json.Unmarshal([]byte(c.LayoutChanges), &resp.LayoutChanges)
	json.Unmarshal([]byte(c.AreaChanges), &resp.AreaChanges)
	json.Unmarshal([]byte(c.ElementChanges), &resp.ElementChanges)
	json.Unmarshal([]byte(c.StyleChanges), &resp.StyleChanges)
	json.Unmarshal([]byte(c.MaterialChanges), &resp.MaterialChanges)

	return resp
}

// compareJSON 比较两个JSON字符串的差异
// Compare two JSON strings for differences
func (s *DesignVersionService) compareJSON(jsonA, jsonB, category string) []response.ChangeItemResponse {
	var changes []response.ChangeItemResponse

	// 简单比较：如果不同则记录变化
	if jsonA != jsonB {
		changes = append(changes, response.ChangeItemResponse{
			Field:       category,
			OldValue:    jsonA,
			NewValue:    jsonB,
			ChangeType:  "modified",
			Description: fmt.Sprintf("%s信息已变更", category),
		})
	}

	return changes
}

// toJSON 将对象转换为JSON字符串
// Convert object to JSON string
func toJSON(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}

// strPtr 将字符串转换为指针（空字符串返回nil）
// Convert string to pointer (empty string returns nil)
func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// ptrToStr 将字符串指针转换为字符串（nil返回空字符串）
// Convert string pointer to string (nil returns empty string)
func ptrToStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// StartAIAnalysis 启动AI版本对比分析（异步）
// Start AI version diff analysis (async)
func (s *DesignVersionService) StartAIAnalysis(req *request.AIAnalyzeVersionDiffRequest, userID uint) (*response.AIAnalyzeVersionDiffResponse, error) {
	// 验证两个版本是否存在 / Verify both versions exist
	var versionA, versionB model.DesignVersion
	if err := s.db.First(&versionA, req.VersionAID).Error; err != nil {
		return nil, errors.New("版本A不存在")
	}
	if err := s.db.First(&versionB, req.VersionBID).Error; err != nil {
		return nil, errors.New("版本B不存在")
	}

	// 验证属于同一项目 / Verify same project
	if versionA.ProjectID != versionB.ProjectID || versionA.ProjectID != req.ProjectID {
		return nil, errors.New("版本必须属于同一项目")
	}

	// 创建对比记录（状态为pending）/ Create compare record with pending status
	// 初始化JSON字段为空数组，避免MySQL JSON字段报错
	compare := &model.DesignVersionCompare{
		ProjectID:       req.ProjectID,
		VersionAID:      req.VersionAID,
		VersionBID:      req.VersionBID,
		LayoutChanges:   "[]",
		AreaChanges:     "[]",
		ElementChanges:  "[]",
		StyleChanges:    "[]",
		MaterialChanges: "[]",
		CompareStatus:   "pending",
		Summary:         "AI分析中...",
	}
	if err := s.db.Create(compare).Error; err != nil {
		return nil, err
	}

	// 异步执行AI分析 / Execute AI analysis asynchronously
	go s.executeAIAnalysis(compare.ID, &versionA, &versionB, userID)

	return &response.AIAnalyzeVersionDiffResponse{
		TaskID:  compare.ID,
		Message: "AI分析任务已提交，完成后将通过通知告知您",
	}, nil
}

// executeAIAnalysis 执行AI分析（后台任务）
// Execute AI analysis (background task)
func (s *DesignVersionService) executeAIAnalysis(compareID uint, versionA, versionB *model.DesignVersion, userID uint) {
	// 构建分析数据 / Build analysis data
	analysisResult := s.performAIAnalysis(versionA, versionB)

	// 更新对比记录 / Update compare record
	updates := map[string]interface{}{
		"layout_changes":   toJSON(analysisResult.LayoutChanges),
		"area_changes":     toJSON(analysisResult.AreaChanges),
		"element_changes":  toJSON(analysisResult.ElementChanges),
		"style_changes":    toJSON(analysisResult.StyleChanges),
		"material_changes": toJSON(analysisResult.MaterialChanges),
		"summary":          analysisResult.Summary,
		"compare_status":   "completed",
	}
	s.db.Model(&model.DesignVersionCompare{}).Where("id = ?", compareID).Updates(updates)

	// 创建通知 / Create notification
	notification := &model.Notification{
		UserID:      int64(userID),
		Title:       fmt.Sprintf("版本对比分析完成：V%d vs V%d", versionA.VersionNumber, versionB.VersionNumber),
		Content:     analysisResult.Summary,
		Type:        "notice",
		RelatedID:   compareID,
		RelatedType: model.NotificationRelatedTypeVersionCompare,
	}
	s.db.Create(notification)
}

// AIAnalysisResult AI分析结果
// AI analysis result structure
type AIAnalysisResult struct {
	LayoutChanges   []response.ChangeItemResponse
	AreaChanges     []response.ChangeItemResponse
	ElementChanges  []response.ChangeItemResponse
	StyleChanges    []response.ChangeItemResponse
	MaterialChanges []response.ChangeItemResponse
	Summary         string
}

// performAIAnalysis 执行AI分析逻辑（使用火山引擎AI分析设计图差异）
// Perform AI analysis logic using VolcEngine AI to analyze design image differences
func (s *DesignVersionService) performAIAnalysis(versionA, versionB *model.DesignVersion) *AIAnalysisResult {
	// 如果没有AI客户端，使用基础对比 / If no AI client, use basic comparison
	if s.aiClient == nil {
		return s.performBasicAnalysis(versionA, versionB)
	}

	// 获取两个版本的设计图 / Get design images from both versions
	imagesA := s.getDesignImageURLs(versionA)
	imagesB := s.getDesignImageURLs(versionB)

	// 如果两个版本都没有设计图，使用基础对比 / If no images in both versions, use basic comparison
	if len(imagesA) == 0 && len(imagesB) == 0 {
		return s.performBasicAnalysis(versionA, versionB)
	}

	// 使用AI分析设计图差异 / Use AI to analyze design image differences
	aiResult := s.analyzeWithAI(versionA, versionB, imagesA, imagesB)
	if aiResult != nil {
		return aiResult
	}

	// AI分析失败，回退到基础对比 / AI analysis failed, fallback to basic comparison
	return s.performBasicAnalysis(versionA, versionB)
}

// getDesignImageURLs 获取设计图URL列表
// Get design image URLs from version
func (s *DesignVersionService) getDesignImageURLs(version *model.DesignVersion) []string {
	var urls []string
	imagesStr := ptrToStr(version.DesignImages)
	if imagesStr == "" || imagesStr == "[]" {
		return urls
	}

	var images []model.DesignImage
	if err := json.Unmarshal([]byte(imagesStr), &images); err != nil {
		return urls
	}

	for _, img := range images {
		if img.URL != "" {
			urls = append(urls, img.URL)
		}
	}
	return urls
}

// analyzeWithAI 使用火山AI分析设计图差异（支持多图对比）
// Analyze design differences using VolcEngine AI (supports multiple images comparison)
func (s *DesignVersionService) analyzeWithAI(versionA, versionB *model.DesignVersion, imagesA, imagesB []string) *AIAnalysisResult {
	// 构建系统提示词 / Build system prompt
	systemPrompt := `你是专业的工装设计评审专家，擅长分析设计方案的差异。请对比分析两个版本的设计方案，从以下维度进行详细分析：

1. 布局变化（layoutChanges）：空间布局、功能分区、动线设计的变化
2. 面积变化（areaChanges）：各功能区面积的调整
3. 元素变化（elementChanges）：设计元素的增减、位置变化
4. 风格变化（styleChanges）：设计风格、色彩搭配、材质选择的变化
5. 材料变化（materialChanges）：使用材料的变更

【重要】请严格按照以下JSON格式返回分析结果，不要返回其他内容：
{
  "layoutChanges": [{"field": "变化项", "oldValue": "旧值", "newValue": "新值", "changeType": "added/removed/modified", "description": "详细描述"}],
  "areaChanges": [],
  "elementChanges": [],
  "styleChanges": [],
  "materialChanges": [],
  "summary": "整体变化摘要（一句话概括主要变化）"
}

注意：
- 如果某个维度没有变化，返回空数组[]
- changeType只能是: added（新增）、removed（删除）、modified（修改）
- description要具体描述变化内容，便于设计师理解
- 必须返回有效的JSON格式，不要包含任何其他文字说明`

	// 构建用户提示词 / Build user prompt
	var userPrompt strings.Builder
	userPrompt.WriteString("请对比分析以下两个设计版本的差异。\n\n")
	userPrompt.WriteString(fmt.Sprintf("【版本A - V%d %s】\n", versionA.VersionNumber, versionA.VersionName))
	userPrompt.WriteString(fmt.Sprintf("描述：%s\n", versionA.Description))
	userPrompt.WriteString(fmt.Sprintf("设计图数量：%d张\n", len(imagesA)))
	if layoutA := ptrToStr(versionA.LayoutInfo); layoutA != "" && layoutA != "[]" {
		userPrompt.WriteString(fmt.Sprintf("布局信息：%s\n", layoutA))
	}
	if styleA := ptrToStr(versionA.StyleInfo); styleA != "" && styleA != "[]" {
		userPrompt.WriteString(fmt.Sprintf("风格信息：%s\n", styleA))
	}
	if materialA := ptrToStr(versionA.MaterialInfo); materialA != "" && materialA != "[]" {
		userPrompt.WriteString(fmt.Sprintf("材料信息：%s\n", materialA))
	}

	userPrompt.WriteString(fmt.Sprintf("\n【版本B - V%d %s】\n", versionB.VersionNumber, versionB.VersionName))
	userPrompt.WriteString(fmt.Sprintf("描述：%s\n", versionB.Description))
	userPrompt.WriteString(fmt.Sprintf("设计图数量：%d张\n", len(imagesB)))
	if layoutB := ptrToStr(versionB.LayoutInfo); layoutB != "" && layoutB != "[]" {
		userPrompt.WriteString(fmt.Sprintf("布局信息：%s\n", layoutB))
	}
	if styleB := ptrToStr(versionB.StyleInfo); styleB != "" && styleB != "[]" {
		userPrompt.WriteString(fmt.Sprintf("风格信息：%s\n", styleB))
	}
	if materialB := ptrToStr(versionB.MaterialInfo); materialB != "" && materialB != "[]" {
		userPrompt.WriteString(fmt.Sprintf("材料信息：%s\n", materialB))
	}

	// 收集所有图片并下载编码 / Collect and encode all images
	var imageContents []volcengine.ContentPart

	// 添加版本A的设计图 / Add version A design images
	for i, imgURL := range imagesA {
		if i >= 2 { // 限制每个版本最多2张图 / Limit to 2 images per version
			break
		}
		imageBase64, err := s.downloadAndEncodeImage(imgURL)
		if err != nil {
			fmt.Printf("[AI分析] 下载版本A图片失败: %v\n", err)
			continue
		}
		// 添加图片标注 / Add image label
		imageContents = append(imageContents, volcengine.ContentPart{
			Type: "text",
			Text: fmt.Sprintf("\n【版本A 设计图%d】:", i+1),
		})
		imageContents = append(imageContents, volcengine.ContentPart{
			Type: "image_url",
			ImageURL: &volcengine.ImageURL{
				URL: "data:image/jpeg;base64," + imageBase64,
			},
		})
	}

	// 添加版本B的设计图 / Add version B design images
	for i, imgURL := range imagesB {
		if i >= 2 { // 限制每个版本最多2张图 / Limit to 2 images per version
			break
		}
		imageBase64, err := s.downloadAndEncodeImage(imgURL)
		if err != nil {
			fmt.Printf("[AI分析] 下载版本B图片失败: %v\n", err)
			continue
		}
		// 添加图片标注 / Add image label
		imageContents = append(imageContents, volcengine.ContentPart{
			Type: "text",
			Text: fmt.Sprintf("\n【版本B 设计图%d】:", i+1),
		})
		imageContents = append(imageContents, volcengine.ContentPart{
			Type: "image_url",
			ImageURL: &volcengine.ImageURL{
				URL: "data:image/jpeg;base64," + imageBase64,
			},
		})
	}

	var aiResponse string
	var err error

	// 如果有图片，使用多图分析 / If has images, use multi-image analysis
	if len(imageContents) > 0 {
		// 构建包含多图的消息 / Build message with multiple images
		userContent := []volcengine.ContentPart{
			{Type: "text", Text: userPrompt.String()},
		}
		userContent = append(userContent, imageContents...)
		userContent = append(userContent, volcengine.ContentPart{
			Type: "text",
			Text: "\n\n请根据以上两个版本的设计图进行对比分析，返回JSON格式的分析结果。",
		})

		messages := []volcengine.ChatMessage{
			{
				Role: "system",
				Content: []volcengine.ContentPart{
					{Type: "text", Text: systemPrompt},
				},
			},
			{
				Role:    "user",
				Content: userContent,
			},
		}

		resp, chatErr := s.aiClient.Chat(messages)
		if chatErr != nil {
			err = chatErr
		} else if len(resp.Choices) > 0 {
			aiResponse = resp.Choices[0].Message.Content
		}
	} else {
		// 没有图片，使用纯文本分析 / No images, use text-only analysis
		aiResponse, err = s.aiClient.ChatWithText(systemPrompt, userPrompt.String())
	}

	if err != nil {
		fmt.Printf("[AI分析] 火山AI调用失败: %v\n", err)
		return nil
	}

	fmt.Printf("[AI分析] AI响应: %s\n", aiResponse)

	// 解析AI响应 / Parse AI response
	return s.parseAIAnalysisResponse(aiResponse, versionA.VersionNumber, versionB.VersionNumber)
}

// downloadAndEncodeImage 下载图片并编码为Base64
// Download image and encode to Base64
func (s *DesignVersionService) downloadAndEncodeImage(imageURL string) (string, error) {
	// 设置超时 / Set timeout
	client := &http.Client{Timeout: 30 * time.Second}

	resp, err := client.Get(imageURL)
	if err != nil {
		return "", fmt.Errorf("下载图片失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("下载图片失败，状态码: %d", resp.StatusCode)
	}

	// 限制图片大小（最大10MB）/ Limit image size (max 10MB)
	limitedReader := io.LimitReader(resp.Body, 10*1024*1024)
	imageData, err := io.ReadAll(limitedReader)
	if err != nil {
		return "", fmt.Errorf("读取图片数据失败: %w", err)
	}

	// 编码为Base64 / Encode to Base64
	return volcengine.EncodeImageToBase64FromBytes(imageData), nil
}

// parseAIAnalysisResponse 解析AI分析响应
// Parse AI analysis response
func (s *DesignVersionService) parseAIAnalysisResponse(aiResponse string, versionANum, versionBNum int) *AIAnalysisResult {
	// 提取JSON部分 / Extract JSON part
	jsonStr := aiResponse
	if idx := strings.Index(aiResponse, "{"); idx != -1 {
		if endIdx := strings.LastIndex(aiResponse, "}"); endIdx != -1 {
			jsonStr = aiResponse[idx : endIdx+1]
		}
	}

	var parsed struct {
		LayoutChanges   []response.ChangeItemResponse `json:"layoutChanges"`
		AreaChanges     []response.ChangeItemResponse `json:"areaChanges"`
		ElementChanges  []response.ChangeItemResponse `json:"elementChanges"`
		StyleChanges    []response.ChangeItemResponse `json:"styleChanges"`
		MaterialChanges []response.ChangeItemResponse `json:"materialChanges"`
		Summary         string                        `json:"summary"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		fmt.Printf("[AI分析] 解析AI响应失败: %v, 原始响应: %s\n", err, aiResponse)
		return nil
	}

	// 确保非nil / Ensure non-nil
	if parsed.LayoutChanges == nil {
		parsed.LayoutChanges = []response.ChangeItemResponse{}
	}
	if parsed.AreaChanges == nil {
		parsed.AreaChanges = []response.ChangeItemResponse{}
	}
	if parsed.ElementChanges == nil {
		parsed.ElementChanges = []response.ChangeItemResponse{}
	}
	if parsed.StyleChanges == nil {
		parsed.StyleChanges = []response.ChangeItemResponse{}
	}
	if parsed.MaterialChanges == nil {
		parsed.MaterialChanges = []response.ChangeItemResponse{}
	}

	// 生成摘要 / Generate summary
	totalChanges := len(parsed.LayoutChanges) + len(parsed.AreaChanges) +
		len(parsed.ElementChanges) + len(parsed.StyleChanges) + len(parsed.MaterialChanges)

	summary := parsed.Summary
	if summary == "" {
		if totalChanges == 0 {
			summary = fmt.Sprintf("V%d 与 V%d 对比：AI分析未发现明显差异", versionANum, versionBNum)
		} else {
			summary = fmt.Sprintf("V%d 与 V%d 对比：AI分析发现 %d 处变化（布局%d、面积%d、元素%d、风格%d、材料%d）",
				versionANum, versionBNum, totalChanges,
				len(parsed.LayoutChanges), len(parsed.AreaChanges), len(parsed.ElementChanges),
				len(parsed.StyleChanges), len(parsed.MaterialChanges))
		}
	}

	return &AIAnalysisResult{
		LayoutChanges:   parsed.LayoutChanges,
		AreaChanges:     parsed.AreaChanges,
		ElementChanges:  parsed.ElementChanges,
		StyleChanges:    parsed.StyleChanges,
		MaterialChanges: parsed.MaterialChanges,
		Summary:         summary,
	}
}

// performBasicAnalysis 执行基础对比分析（无AI时的回退方案）
// Perform basic comparison analysis (fallback when no AI)
func (s *DesignVersionService) performBasicAnalysis(versionA, versionB *model.DesignVersion) *AIAnalysisResult {
	result := &AIAnalysisResult{}

	// 分析布局变化 / Analyze layout changes
	result.LayoutChanges = s.compareJSON(ptrToStr(versionA.LayoutInfo), ptrToStr(versionB.LayoutInfo), "布局")

	// 分析面积变化 / Analyze area changes
	result.AreaChanges = s.compareJSON(ptrToStr(versionA.AreaInfo), ptrToStr(versionB.AreaInfo), "面积")

	// 分析元素变化 / Analyze element changes
	result.ElementChanges = s.analyzeElementChanges(versionA, versionB)

	// 分析风格变化 / Analyze style changes
	result.StyleChanges = s.compareJSON(ptrToStr(versionA.StyleInfo), ptrToStr(versionB.StyleInfo), "风格")

	// 分析材料变化 / Analyze material changes
	result.MaterialChanges = s.compareJSON(ptrToStr(versionA.MaterialInfo), ptrToStr(versionB.MaterialInfo), "材料")

	// 生成摘要 / Generate summary
	totalChanges := len(result.LayoutChanges) + len(result.AreaChanges) +
		len(result.ElementChanges) + len(result.StyleChanges) + len(result.MaterialChanges)

	if totalChanges == 0 {
		result.Summary = fmt.Sprintf("V%d 与 V%d 对比：未发现明显差异", versionA.VersionNumber, versionB.VersionNumber)
	} else {
		result.Summary = fmt.Sprintf("V%d 与 V%d 对比：共发现 %d 处变化（布局%d、面积%d、元素%d、风格%d、材料%d）",
			versionA.VersionNumber, versionB.VersionNumber, totalChanges,
			len(result.LayoutChanges), len(result.AreaChanges), len(result.ElementChanges),
			len(result.StyleChanges), len(result.MaterialChanges))
	}

	return result
}

// analyzeElementChanges 分析元素变化
func (s *DesignVersionService) analyzeElementChanges(versionA, versionB *model.DesignVersion) []response.ChangeItemResponse {
	// 比较设计图数量变化 / Compare design image count changes
	var changes []response.ChangeItemResponse

	imagesA := ptrToStr(versionA.DesignImages)
	imagesB := ptrToStr(versionB.DesignImages)

	if imagesA != imagesB {
		var countA, countB int
		var imgListA, imgListB []map[string]any
		json.Unmarshal([]byte(imagesA), &imgListA)
		json.Unmarshal([]byte(imagesB), &imgListB)
		countA = len(imgListA)
		countB = len(imgListB)

		if countA != countB {
			changeType := "modified"
			if countB > countA {
				changeType = "added"
			} else if countB < countA {
				changeType = "removed"
			}
			changes = append(changes, response.ChangeItemResponse{
				Field:       "设计图",
				OldValue:    fmt.Sprintf("%d张", countA),
				NewValue:    fmt.Sprintf("%d张", countB),
				ChangeType:  changeType,
				Description: fmt.Sprintf("设计图数量从%d张变为%d张", countA, countB),
			})
		}
	}

	// 比较CAD文件变化 / Compare CAD file changes
	cadA := ptrToStr(versionA.CadFileIDs)
	cadB := ptrToStr(versionB.CadFileIDs)
	if cadA != cadB {
		changes = append(changes, response.ChangeItemResponse{
			Field:       "CAD文件",
			OldValue:    cadA,
			NewValue:    cadB,
			ChangeType:  "modified",
			Description: "关联的CAD文件已变更",
		})
	}

	return changes
}
