package design_compare

import (
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/dto/response"
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/volcengine"
	"art_admin_backend/internal/repository"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// DesignCompareService 设计比对服务（支持多图片、多文档、CAD文件）
// Design compare service supporting multiple images, documents and CAD files
type DesignCompareService struct {
	repo        *repository.DesignCompareRepository
	projectRepo *repository.DesignerProjectRepository
	docRepo     *repository.DocumentRepository
	cadRepo     *repository.CadFileRepository // CAD文件仓库 / CAD file repository
	aiClient    *volcengine.Client
	uploadPath  string
	baseURL     string
}

// NewDesignCompareService 创建设计比对服务
func NewDesignCompareService(
	repo *repository.DesignCompareRepository,
	projectRepo *repository.DesignerProjectRepository,
	docRepo *repository.DocumentRepository,
	cadRepo *repository.CadFileRepository,
	aiClient *volcengine.Client,
	uploadPath, baseURL string,
) *DesignCompareService {
	os.MkdirAll(filepath.Join(uploadPath, "design-images"), 0755)
	return &DesignCompareService{
		repo: repo, projectRepo: projectRepo, docRepo: docRepo,
		cadRepo:  cadRepo,
		aiClient: aiClient, uploadPath: uploadPath, baseURL: baseURL,
	}
}

// SupportedFileTypes 支持的文件类型
var SupportedFileTypes = map[string]string{
	".jpg": "image", ".jpeg": "image", ".png": "image", ".webp": "image",
	".pdf": "pdf",
}

// Upload 批量上传设计图（支持多文件、已有效果图）
// Batch upload design images supporting multiple files and existing render images
func (s *DesignCompareService) Upload(files []*multipart.FileHeader, req *request.UploadDesignImagesRequest, userID uint) (*response.DesignCompareUploadResponse, error) {
	// 验证项目权限
	project, err := s.projectRepo.GetByID(req.ProjectID)
	if err != nil {
		return nil, errors.New("项目不存在")
	}
	if project.UserID != userID {
		return nil, errors.New("无权访问该项目")
	}

	// 解析文档ID列表
	docIDs, err := s.parseDocumentIDs(req.DocumentIDs)
	if err != nil {
		return nil, err
	}
	if len(docIDs) == 0 {
		return nil, errors.New("请选择至少一份需求文档")
	}

	// 验证文档存在且属于该项目
	for _, docID := range docIDs {
		doc, err := s.docRepo.GetByIDWithProject(docID)
		if err != nil {
			return nil, fmt.Errorf("文档ID %d 不存在", docID)
		}
		if doc.ProjectID != req.ProjectID {
			return nil, fmt.Errorf("文档ID %d 不属于该项目", docID)
		}
	}

	// 解析已有效果图URL
	existingUrls := s.parseImageUrls(req.ExistingImageUrls)

	if len(files) == 0 && len(existingUrls) == 0 {
		return nil, errors.New("请上传至少一个设计文件或选择已有效果图")
	}

	// 保存所有上传的文件
	var designImages []model.DesignImageInfo
	dateDir := time.Now().Format("2006/01/02")
	fullDir := filepath.Join(s.uploadPath, "design-images", dateDir)
	os.MkdirAll(fullDir, 0755)

	for _, file := range files {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		fileType, ok := SupportedFileTypes[ext]
		if !ok {
			return nil, fmt.Errorf("不支持的文件格式: %s，支持: JPG, PNG, WEBP, PDF", ext)
		}
		if file.Size > 50*1024*1024 {
			return nil, fmt.Errorf("文件 %s 超过50MB限制", file.Filename)
		}

		// 保存文件
		storageName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
		filePath := filepath.Join("design-images", dateDir, storageName)
		fullPath := filepath.Join(s.uploadPath, filePath)

		if err := s.saveFile(file, fullPath); err != nil {
			return nil, fmt.Errorf("保存文件 %s 失败: %w", file.Filename, err)
		}

		designImages = append(designImages, model.DesignImageInfo{
			FileName: file.Filename,
			FilePath: filePath,
			FileType: fileType,
			FileSize: file.Size,
		})
	}

	// 添加已有效果图
	for i, url := range existingUrls {
		designImages = append(designImages, model.DesignImageInfo{
			FileName: fmt.Sprintf("效果图_%d", i+1),
			FilePath: url, // 直接使用URL作为路径
			FileType: "image",
			FileSize: 0,
		})
	}

	// 序列化数据
	docIDsJSON, _ := json.Marshal(docIDs)
	imagesJSON, _ := json.Marshal(designImages)

	// 生成任务名称
	name := req.Name
	if name == "" {
		name = fmt.Sprintf("设计比对_%s", time.Now().Format("20060102_150405"))
	}

	// 创建比对记录
	compare := &model.DesignCompareResult{
		ProjectID:      req.ProjectID,
		Name:           name,
		DocumentIDs:    string(docIDsJSON),
		DesignImages:   string(imagesJSON),
		AnalysisStatus: "pending",
		MatchItems:     "[]",
		DeviationItems: "[]",
		Suggestions:    "[]",
	}

	if err := s.repo.Create(compare); err != nil {
		// 清理已上传的文件
		for _, img := range designImages {
			if !strings.HasPrefix(img.FilePath, "http") {
				os.Remove(filepath.Join(s.uploadPath, img.FilePath))
			}
		}
		return nil, fmt.Errorf("保存比对记录失败: %w", err)
	}

	return &response.DesignCompareUploadResponse{
		ID:             compare.ID,
		ProjectID:      compare.ProjectID,
		Name:           compare.Name,
		DocumentIDs:    docIDs,
		DesignImages:   s.toImageInfoResponses(designImages),
		AnalysisStatus: compare.AnalysisStatus,
	}, nil
}

// parseImageUrls 解析效果图URL列表
func (s *DesignCompareService) parseImageUrls(urlsStr string) []string {
	if urlsStr == "" {
		return nil
	}
	parts := strings.Split(urlsStr, ",")
	var urls []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			urls = append(urls, p)
		}
	}
	return urls
}

func (s *DesignCompareService) parseDocumentIDs(idsStr string) ([]uint, error) {
	if idsStr == "" {
		return nil, nil
	}
	parts := strings.Split(idsStr, ",")
	var ids []uint
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		id, err := strconv.ParseUint(p, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("无效的文档ID: %s", p)
		}
		ids = append(ids, uint(id))
	}
	return ids, nil
}

// CreateWithCadFiles 使用已有CAD文件创建比对任务
// Create compare task with existing CAD files
func (s *DesignCompareService) CreateWithCadFiles(req *request.CreateCompareWithCadRequest, userID uint) (*response.DesignCompareUploadResponse, error) {
	// 验证项目权限
	project, err := s.projectRepo.GetByID(req.ProjectID)
	if err != nil {
		return nil, errors.New("项目不存在")
	}
	if project.UserID != userID {
		return nil, errors.New("无权访问该项目")
	}

	// 验证文档存在且属于该项目
	for _, docID := range req.DocumentIDs {
		doc, err := s.docRepo.GetByIDWithProject(docID)
		if err != nil {
			return nil, fmt.Errorf("文档ID %d 不存在", docID)
		}
		if doc.ProjectID != req.ProjectID {
			return nil, fmt.Errorf("文档ID %d 不属于该项目", docID)
		}
	}

	if len(req.CadFileIDs) == 0 {
		return nil, errors.New("请选择至少一个CAD文件")
	}

	// 获取CAD文件信息并构建设计图列表
	var designImages []model.DesignImageInfo
	for _, cadID := range req.CadFileIDs {
		cadFile, err := s.cadRepo.GetByID(cadID)
		if err != nil {
			return nil, fmt.Errorf("CAD文件ID %d 不存在", cadID)
		}
		if cadFile.ProjectID != req.ProjectID {
			return nil, fmt.Errorf("CAD文件ID %d 不属于该项目", cadID)
		}

		// 确定文件类型
		ext := strings.ToLower(filepath.Ext(cadFile.FileName))
		fileType := "cad"
		if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".webp" {
			fileType = "image"
		} else if ext == ".pdf" {
			fileType = "pdf"
		}

		designImages = append(designImages, model.DesignImageInfo{
			FileName: cadFile.FileName,
			FilePath: cadFile.OriginalPath,
			FileType: fileType,
			FileSize: cadFile.FileSize,
		})
	}

	// 序列化数据
	docIDsJSON, _ := json.Marshal(req.DocumentIDs)
	imagesJSON, _ := json.Marshal(designImages)

	// 生成任务名称
	name := req.Name
	if name == "" {
		name = fmt.Sprintf("设计比对_%s", time.Now().Format("20060102_150405"))
	}

	// 创建比对记录
	compare := &model.DesignCompareResult{
		ProjectID:      req.ProjectID,
		Name:           name,
		DocumentIDs:    string(docIDsJSON),
		DesignImages:   string(imagesJSON),
		AnalysisStatus: "pending",
		MatchItems:     "[]",
		DeviationItems: "[]",
		Suggestions:    "[]",
	}

	if err := s.repo.Create(compare); err != nil {
		return nil, fmt.Errorf("保存比对记录失败: %w", err)
	}

	return &response.DesignCompareUploadResponse{
		ID:             compare.ID,
		ProjectID:      compare.ProjectID,
		Name:           compare.Name,
		DocumentIDs:    req.DocumentIDs,
		DesignImages:   s.toImageInfoResponses(designImages),
		AnalysisStatus: compare.AnalysisStatus,
	}, nil
}

// CreateWithImageUrls 使用已有效果图URL创建比对任务
// Create compare task with existing render image URLs
func (s *DesignCompareService) CreateWithImageUrls(req *request.CreateCompareWithImagesRequest, userID uint) (*response.DesignCompareUploadResponse, error) {
	// 验证项目权限
	project, err := s.projectRepo.GetByID(req.ProjectID)
	if err != nil {
		return nil, errors.New("项目不存在")
	}
	if project.UserID != userID {
		return nil, errors.New("无权访问该项目")
	}

	// 验证文档存在且属于该项目
	for _, docID := range req.DocumentIDs {
		doc, err := s.docRepo.GetByIDWithProject(docID)
		if err != nil {
			return nil, fmt.Errorf("文档ID %d 不存在", docID)
		}
		if doc.ProjectID != req.ProjectID {
			return nil, fmt.Errorf("文档ID %d 不属于该项目", docID)
		}
	}

	if len(req.ImageUrls) == 0 {
		return nil, errors.New("请选择至少一张效果图")
	}

	// 构建设计图列表
	var designImages []model.DesignImageInfo
	for i, url := range req.ImageUrls {
		designImages = append(designImages, model.DesignImageInfo{
			FileName: fmt.Sprintf("效果图_%d", i+1),
			FilePath: url,
			FileType: "image",
			FileSize: 0,
		})
	}

	// 序列化数据
	docIDsJSON, _ := json.Marshal(req.DocumentIDs)
	imagesJSON, _ := json.Marshal(designImages)

	// 生成任务名称
	name := req.Name
	if name == "" {
		name = fmt.Sprintf("设计比对_%s", time.Now().Format("20060102_150405"))
	}

	// 创建比对记录
	compare := &model.DesignCompareResult{
		ProjectID:      req.ProjectID,
		Name:           name,
		DocumentIDs:    string(docIDsJSON),
		DesignImages:   string(imagesJSON),
		AnalysisStatus: "pending",
		MatchItems:     "[]",
		DeviationItems: "[]",
		Suggestions:    "[]",
	}

	if err := s.repo.Create(compare); err != nil {
		return nil, fmt.Errorf("保存比对记录失败: %w", err)
	}

	return &response.DesignCompareUploadResponse{
		ID:             compare.ID,
		ProjectID:      compare.ProjectID,
		Name:           compare.Name,
		DocumentIDs:    req.DocumentIDs,
		DesignImages:   s.toImageInfoResponses(designImages),
		AnalysisStatus: compare.AnalysisStatus,
	}, nil
}

func (s *DesignCompareService) saveFile(file *multipart.FileHeader, fullPath string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

// Analyze 分析设计比对（异步执行，支持多图多文档）
// Analyze design comparison asynchronously with multiple images and documents
func (s *DesignCompareService) Analyze(compareID, userID uint) (*response.DesignCompareAnalysisResult, error) {
	compare, err := s.repo.GetByIDWithRelations(compareID)
	if err != nil {
		return nil, errors.New("比对记录不存在")
	}

	if compare.Project != nil && compare.Project.UserID != userID {
		return nil, errors.New("无权访问该比对记录")
	}

	if compare.AnalysisStatus == "processing" {
		return nil, errors.New("比对分析正在进行中，请稍候")
	}

	// 更新状态为处理中
	s.repo.UpdateAnalysisResult(compareID, map[string]any{
		"analysis_status": "processing",
		"error_message":   "",
	})

	// 解析设计图和文档ID
	var designImages []model.DesignImageInfo
	var docIDs []uint
	json.Unmarshal([]byte(compare.DesignImages), &designImages)
	json.Unmarshal([]byte(compare.DocumentIDs), &docIDs)

	// 异步执行分析
	go s.doAnalyze(compareID, designImages, docIDs)

	return &response.DesignCompareAnalysisResult{
		CompareID:      compareID,
		MatchItems:     []response.MatchItemResponse{},
		DeviationItems: []response.DeviationItemResponse{},
		Suggestions:    []response.SuggestionItemResponse{},
		OverallScore:   0,
	}, nil
}

// doAnalyze 执行实际的分析任务
func (s *DesignCompareService) doAnalyze(compareID uint, designImages []model.DesignImageInfo, docIDs []uint) {
	// 1. 获取所有文档的需求内容
	requirements := s.buildAllRequirementsText(docIDs)

	// 2. 分析所有设计图
	var allMatchItems []response.MatchItemResponse
	var allDeviationItems []response.DeviationItemResponse
	var allSuggestions []response.SuggestionItemResponse
	var totalScore float64
	var analyzedCount int

	for _, img := range designImages {
		var imagePath string
		// 判断是 URL 还是本地文件路径
		if strings.HasPrefix(img.FilePath, "http://") || strings.HasPrefix(img.FilePath, "https://") {
			imagePath = img.FilePath // 直接使用 URL
		} else {
			imagePath = filepath.Join(s.uploadPath, img.FilePath)
		}

		result, err := s.analyzeOneImage(imagePath, img.FileName, img.FileType, requirements)
		if err != nil {
			// 记录单个文件分析失败，继续处理其他文件
			allDeviationItems = append(allDeviationItems, response.DeviationItemResponse{
				Location:            img.FileName,
				Content:             fmt.Sprintf("分析失败: %s", err.Error()),
				OriginalRequirement: "",
				Severity:            "high",
				SourceImage:         img.FileName,
			})
			continue
		}

		// 合并结果
		for i := range result.MatchItems {
			result.MatchItems[i].SourceImage = img.FileName
		}
		for i := range result.DeviationItems {
			result.DeviationItems[i].SourceImage = img.FileName
		}

		allMatchItems = append(allMatchItems, result.MatchItems...)
		allDeviationItems = append(allDeviationItems, result.DeviationItems...)
		allSuggestions = append(allSuggestions, result.Suggestions...)
		totalScore += result.OverallScore
		analyzedCount++
	}

	// 计算平均分
	var overallScore float64
	if analyzedCount > 0 {
		overallScore = totalScore / float64(analyzedCount)
	}

	// 保存结果
	matchItemsJSON, _ := json.Marshal(allMatchItems)
	deviationItemsJSON, _ := json.Marshal(allDeviationItems)
	suggestionsJSON, _ := json.Marshal(allSuggestions)

	s.repo.UpdateAnalysisResult(compareID, map[string]any{
		"analysis_status": "completed",
		"match_items":     string(matchItemsJSON),
		"deviation_items": string(deviationItemsJSON),
		"suggestions":     string(suggestionsJSON),
		"overall_score":   overallScore,
		"error_message":   "",
	})
}

// buildAllRequirementsText 构建所有文档的需求文本
func (s *DesignCompareService) buildAllRequirementsText(docIDs []uint) string {
	var allParts []string

	for _, docID := range docIDs {
		doc, err := s.docRepo.GetByIDWithProject(docID)
		if err != nil {
			continue
		}

		var parts []string
		parts = append(parts, fmt.Sprintf("【文档: %s】", doc.FileName))

		if doc.ProjectName != "" {
			parts = append(parts, fmt.Sprintf("项目名称: %s", doc.ProjectName))
		}
		if doc.ExtractedArea > 0 {
			parts = append(parts, fmt.Sprintf("面积: %.2f平方米", doc.ExtractedArea))
		}
		if doc.ExtractedBudget > 0 {
			parts = append(parts, fmt.Sprintf("预算: %.2f元", doc.ExtractedBudget))
		}
		if doc.ExtractedStyle != "" {
			parts = append(parts, fmt.Sprintf("设计风格: %s", doc.ExtractedStyle))
		}
		if doc.FunctionalZones != "" {
			var zones []string
			json.Unmarshal([]byte(doc.FunctionalZones), &zones)
			if len(zones) > 0 {
				parts = append(parts, fmt.Sprintf("功能分区: %s", strings.Join(zones, "、")))
			}
		}
		if doc.Summary != "" {
			parts = append(parts, fmt.Sprintf("需求摘要: %s", doc.Summary))
		}

		allParts = append(allParts, strings.Join(parts, "\n"))
	}

	return strings.Join(allParts, "\n\n")
}

// analyzeOneImage 分析单张设计图
func (s *DesignCompareService) analyzeOneImage(imagePath, fileName, fileType, requirements string) (*response.DesignCompareAnalysisResult, error) {
	systemPrompt := `你是专业的工装设计评审专家。请分析设计图与项目需求的匹配程度。

分析要求：
1. 识别设计图中的空间布局、功能区域、设计元素
2. 将设计图与需求进行逐项比对
3. 找出匹配项、偏差项，并给出改进建议

返回JSON格式：
{
  "matchItems": [{"requirement": "需求描述", "designMatch": "设计匹配点", "score": 85}],
  "deviationItems": [{"location": "偏差位置", "content": "偏差内容", "originalRequirement": "原始需求", "severity": "high/medium/low"}],
  "suggestions": [{"content": "建议内容", "priority": "high/medium/low", "costImpact": "成本影响说明"}],
  "overallScore": 75.5
}`

	userPrompt := fmt.Sprintf("请分析这张设计图（%s）与以下项目需求的匹配程度：\n\n%s", fileName, requirements)
	if requirements == "" {
		userPrompt = fmt.Sprintf("请分析这张设计图（%s），识别其中的空间布局、功能区域和设计元素，并给出专业评价。", fileName)
	}

	var aiResponse string
	var err error

	// 判断是 URL 还是本地文件
	isURL := strings.HasPrefix(imagePath, "http://") || strings.HasPrefix(imagePath, "https://")

	switch fileType {
	case "image":
		var imageBase64 string
		if isURL {
			// 从 URL 获取图片并编码
			imageBase64, err = volcengine.EncodeImageFromURL(imagePath)
			if err != nil {
				return nil, fmt.Errorf("从URL获取图片失败: %w", err)
			}
		} else {
			imageBase64, err = volcengine.EncodeImageToBase64(imagePath)
			if err != nil {
				return nil, fmt.Errorf("编码图片失败: %w", err)
			}
		}
		aiResponse, err = s.aiClient.ChatWithImage(systemPrompt, userPrompt, imageBase64)

	case "pdf":
		if isURL {
			return nil, fmt.Errorf("暂不支持URL类型的PDF文件")
		}
		// PDF 文件也尝试作为图片处理（如果是设计图PDF）
		imageBase64, encErr := volcengine.EncodeImageToBase64(imagePath)
		if encErr != nil {
			return nil, fmt.Errorf("编码PDF失败: %w", encErr)
		}
		aiResponse, err = s.aiClient.ChatWithImage(systemPrompt, userPrompt, imageBase64)

	case "cad":
		// CAD 文件暂时返回提示信息，后续可接入专业CAD解析服务
		return &response.DesignCompareAnalysisResult{
			MatchItems: []response.MatchItemResponse{},
			DeviationItems: []response.DeviationItemResponse{
				{
					Location:            fileName,
					Content:             "CAD文件(.dwg/.dxf)需要专业软件解析，建议导出为PDF或图片格式后重新上传",
					OriginalRequirement: "",
					Severity:            "medium",
				},
			},
			Suggestions: []response.SuggestionItemResponse{
				{
					Content:    "建议将CAD文件导出为PDF或高清图片格式，以便进行更准确的AI分析",
					Priority:   "high",
					CostImpact: "无额外成本",
				},
			},
			OverallScore: 0,
		}, nil

	default:
		return nil, fmt.Errorf("不支持的文件类型: %s", fileType)
	}

	if err != nil {
		return nil, err
	}

	return s.parseAIResponse(aiResponse)
}

// parseAIResponse 解析AI响应
func (s *DesignCompareService) parseAIResponse(aiResponse string) (*response.DesignCompareAnalysisResult, error) {
	jsonStr := aiResponse
	if idx := strings.Index(aiResponse, "{"); idx != -1 {
		if endIdx := strings.LastIndex(aiResponse, "}"); endIdx != -1 {
			jsonStr = aiResponse[idx : endIdx+1]
		}
	}

	var result struct {
		MatchItems     []response.MatchItemResponse      `json:"matchItems"`
		DeviationItems []response.DeviationItemResponse  `json:"deviationItems"`
		Suggestions    []response.SuggestionItemResponse `json:"suggestions"`
		OverallScore   float64                           `json:"overallScore"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return &response.DesignCompareAnalysisResult{
			MatchItems:     []response.MatchItemResponse{},
			DeviationItems: []response.DeviationItemResponse{},
			Suggestions: []response.SuggestionItemResponse{
				{Content: "AI分析结果解析失败，建议重新上传更清晰的设计图", Priority: "high", CostImpact: "无"},
			},
			OverallScore: 0,
		}, nil
	}

	if result.MatchItems == nil {
		result.MatchItems = []response.MatchItemResponse{}
	}
	if result.DeviationItems == nil {
		result.DeviationItems = []response.DeviationItemResponse{}
	}
	if result.Suggestions == nil {
		result.Suggestions = []response.SuggestionItemResponse{}
	}

	return &response.DesignCompareAnalysisResult{
		MatchItems:     result.MatchItems,
		DeviationItems: result.DeviationItems,
		Suggestions:    result.Suggestions,
		OverallScore:   result.OverallScore,
	}, nil
}

// GetByID 根据ID获取比对记录
func (s *DesignCompareService) GetByID(id, userID uint) (*response.DesignCompareResponse, error) {
	compare, err := s.repo.GetByIDWithRelations(id)
	if err != nil {
		return nil, errors.New("比对记录不存在")
	}

	if compare.Project != nil && compare.Project.UserID != userID {
		return nil, errors.New("无权访问该比对记录")
	}

	return s.toResponse(compare), nil
}

// List 获取比对记录列表
func (s *DesignCompareService) List(req *request.DesignCompareListRequest, userID uint) (*response.DesignCompareListResponse, error) {
	project, err := s.projectRepo.GetByID(req.ProjectID)
	if err != nil {
		return nil, errors.New("项目不存在")
	}
	if project.UserID != userID {
		return nil, errors.New("无权访问该项目")
	}

	filter := repository.DesignCompareFilter{
		ProjectID:      req.ProjectID,
		AnalysisStatus: req.AnalysisStatus,
	}

	compares, total, err := s.repo.List(req.Current, req.Size, filter)
	if err != nil {
		return nil, err
	}

	records := make([]response.DesignCompareResponse, len(compares))
	for i, compare := range compares {
		records[i] = *s.toResponse(&compare)
	}

	return &response.DesignCompareListResponse{
		Records: records,
		Current: req.Current,
		Size:    req.Size,
		Total:   total,
	}, nil
}

// Delete 删除比对记录
func (s *DesignCompareService) Delete(id, userID uint) error {
	compare, err := s.repo.GetByIDWithRelations(id)
	if err != nil {
		return errors.New("比对记录不存在")
	}

	if compare.Project != nil && compare.Project.UserID != userID {
		return errors.New("无权删除该比对记录")
	}

	// 删除设计图文件
	var designImages []model.DesignImageInfo
	json.Unmarshal([]byte(compare.DesignImages), &designImages)
	for _, img := range designImages {
		os.Remove(filepath.Join(s.uploadPath, img.FilePath))
	}

	return s.repo.Delete(id)
}

// BatchDelete 批量删除比对记录
func (s *DesignCompareService) BatchDelete(ids []uint, userID uint) error {
	for _, id := range ids {
		s.Delete(id, userID)
	}
	return nil
}

// toResponse 转换为响应对象
func (s *DesignCompareService) toResponse(compare *model.DesignCompareResult) *response.DesignCompareResponse {
	var docIDs []uint
	var designImages []model.DesignImageInfo
	var matchItems []response.MatchItemResponse
	var deviationItems []response.DeviationItemResponse
	var suggestions []response.SuggestionItemResponse

	json.Unmarshal([]byte(compare.DocumentIDs), &docIDs)
	json.Unmarshal([]byte(compare.DesignImages), &designImages)
	json.Unmarshal([]byte(compare.MatchItems), &matchItems)
	json.Unmarshal([]byte(compare.DeviationItems), &deviationItems)
	json.Unmarshal([]byte(compare.Suggestions), &suggestions)

	// 确保非nil
	if docIDs == nil {
		docIDs = []uint{}
	}
	if matchItems == nil {
		matchItems = []response.MatchItemResponse{}
	}
	if deviationItems == nil {
		deviationItems = []response.DeviationItemResponse{}
	}
	if suggestions == nil {
		suggestions = []response.SuggestionItemResponse{}
	}

	return &response.DesignCompareResponse{
		ID:             compare.ID,
		ProjectID:      compare.ProjectID,
		Name:           compare.Name,
		DocumentIDs:    docIDs,
		DesignImages:   s.toImageInfoResponses(designImages),
		MatchItems:     matchItems,
		DeviationItems: deviationItems,
		Suggestions:    suggestions,
		OverallScore:   compare.OverallScore,
		AnalysisStatus: compare.AnalysisStatus,
		ErrorMessage:   compare.ErrorMessage,
		CreatedAt:      compare.CreatedAt,
		UpdatedAt:      compare.UpdatedAt,
	}
}

func (s *DesignCompareService) toImageInfoResponses(images []model.DesignImageInfo) []response.DesignImageInfoResponse {
	result := make([]response.DesignImageInfoResponse, len(images))
	for i, img := range images {
		result[i] = response.DesignImageInfoResponse{
			FileName: img.FileName,
			FilePath: img.FilePath,
			FileURL:  fmt.Sprintf("%s/uploads/%s", s.baseURL, img.FilePath),
			FileType: img.FileType,
			FileSize: img.FileSize,
		}
	}
	return result
}
