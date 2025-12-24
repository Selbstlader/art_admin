package document

import (
	"archive/zip"
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/dto/response"
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/volcengine"
	"art_admin_backend/internal/repository"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// DocumentService 文档服务
type DocumentService struct {
	repo        *repository.DocumentRepository
	projectRepo *repository.DesignerProjectRepository
	aiClient    *volcengine.Client
	uploadPath  string
	baseURL     string
}

// NewDocumentService 创建文档服务
func NewDocumentService(
	repo *repository.DocumentRepository,
	projectRepo *repository.DesignerProjectRepository,
	aiClient *volcengine.Client,
	uploadPath, baseURL string,
) *DocumentService {
	os.MkdirAll(filepath.Join(uploadPath, "documents"), 0755)
	return &DocumentService{
		repo: repo, projectRepo: projectRepo, aiClient: aiClient,
		uploadPath: uploadPath, baseURL: baseURL,
	}
}

// SupportedFileTypes 支持的文件类型
var SupportedFileTypes = map[string]bool{
	".pdf": true, ".doc": true, ".docx": true,
	".jpg": true, ".jpeg": true, ".png": true, ".webp": true,
}

// Upload 上传文档
func (s *DocumentService) Upload(file *multipart.FileHeader, projectID uint, userID uint) (*response.DocumentUploadResponse, error) {
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return nil, errors.New("项目不存在")
	}
	if project.UserID != userID {
		return nil, errors.New("无权访问该项目")
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !SupportedFileTypes[ext] {
		return nil, fmt.Errorf("不支持的文件格式，支持: PDF, Word, JPG, PNG, WEBP")
	}
	if file.Size > 50*1024*1024 {
		return nil, errors.New("文件大小不能超过50MB")
	}

	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer src.Close()

	storageName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	dateDir := time.Now().Format("2006/01/02")
	fullDir := filepath.Join(s.uploadPath, "documents", dateDir)
	os.MkdirAll(fullDir, 0755)

	filePath := filepath.Join("documents", dateDir, storageName)
	fullPath := filepath.Join(s.uploadPath, filePath)

	dst, err := os.Create(fullPath)
	if err != nil {
		return nil, fmt.Errorf("创建文件失败: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		os.Remove(fullPath)
		return nil, fmt.Errorf("保存文件失败: %w", err)
	}

	fileType := s.detectFileType(ext)
	doc := &model.ProjectDocument{
		ProjectID: projectID, FileName: file.Filename, FilePath: filePath,
		FileType: fileType, FileSize: file.Size, AnalysisStatus: "pending",
		Keywords: "[]", FunctionalZones: "[]", // 初始化 JSON 字段为空数组
	}

	if err := s.repo.Create(doc); err != nil {
		os.Remove(fullPath)
		return nil, fmt.Errorf("保存文档记录失败: %w", err)
	}

	return &response.DocumentUploadResponse{
		ID: doc.ID, FileName: doc.FileName, FileType: doc.FileType,
		FileSize: doc.FileSize, ProjectID: doc.ProjectID, Status: doc.AnalysisStatus,
	}, nil
}

func (s *DocumentService) detectFileType(ext string) string {
	switch ext {
	case ".pdf":
		return "pdf"
	case ".doc", ".docx":
		return "word"
	case ".jpg", ".jpeg", ".png", ".webp":
		return "image"
	default:
		return "other"
	}
}

// GetByID 根据ID获取文档
func (s *DocumentService) GetByID(id uint, userID uint) (*response.DocumentResponse, error) {
	doc, err := s.repo.GetByIDWithProject(id)
	if err != nil {
		return nil, errors.New("文档不存在")
	}
	if doc.Project != nil && doc.Project.UserID != userID {
		return nil, errors.New("无权访问该文档")
	}
	return s.toResponse(doc), nil
}

// List 获取文档列表
func (s *DocumentService) List(req *request.DocumentListRequest, userID uint) (*response.DocumentListResponse, error) {
	project, err := s.projectRepo.GetByID(req.ProjectID)
	if err != nil {
		return nil, errors.New("项目不存在")
	}
	if project.UserID != userID {
		return nil, errors.New("无权访问该项目")
	}

	filter := repository.DocumentFilter{
		ProjectID: req.ProjectID, AnalysisStatus: req.AnalysisStatus, FileType: req.FileType,
	}
	docs, total, err := s.repo.List(req.Current, req.Size, filter)
	if err != nil {
		return nil, err
	}

	records := make([]response.DocumentResponse, len(docs))
	for i, doc := range docs {
		records[i] = *s.toResponse(&doc)
	}
	return &response.DocumentListResponse{Records: records, Current: req.Current, Size: req.Size, Total: total}, nil
}

// Delete 删除文档
func (s *DocumentService) Delete(id uint, userID uint) error {
	doc, err := s.repo.GetByIDWithProject(id)
	if err != nil {
		return errors.New("文档不存在")
	}
	if doc.Project != nil && doc.Project.UserID != userID {
		return errors.New("无权删除该文档")
	}
	fullPath := filepath.Join(s.uploadPath, doc.FilePath)
	os.Remove(fullPath)
	return s.repo.Delete(id)
}

// BatchDelete 批量删除文档
func (s *DocumentService) BatchDelete(ids []uint, userID uint) error {
	for _, id := range ids {
		s.Delete(id, userID)
	}
	return nil
}

// Analyze 分析文档（异步执行）
// Analyze document asynchronously
func (s *DocumentService) Analyze(documentID uint, userID uint) (*response.DocumentAnalysisResult, error) {
	doc, err := s.repo.GetByIDWithProject(documentID)
	if err != nil {
		return nil, errors.New("文档不存在")
	}
	if doc.Project != nil && doc.Project.UserID != userID {
		return nil, errors.New("无权分析该文档")
	}

	// 检查是否已在处理中
	if doc.AnalysisStatus == "processing" {
		return nil, errors.New("文档正在分析中，请稍候")
	}

	// 立即更新状态为处理中，清除之前的错误信息
	s.repo.UpdateAnalysisResult(documentID, map[string]any{
		"analysis_status": "processing",
		"error_message":   "",
	})

	// 异步执行分析任务
	go s.doAnalyze(documentID, doc.FilePath, doc.FileType)

	// 立即返回，告知用户分析已开始
	return &response.DocumentAnalysisResult{
		DocumentID: documentID,
		Summary:    "文档分析已开始，请稍后刷新查看结果",
	}, nil
}

// doAnalyze 执行实际的分析任务（后台运行）
// Execute actual analysis task in background
func (s *DocumentService) doAnalyze(documentID uint, filePath, fileType string) {
	fullPath := filepath.Join(s.uploadPath, filePath)

	// 提取文档内容
	content, err := s.extractContent(fullPath, fileType)
	if err != nil {
		s.repo.UpdateAnalysisResult(documentID, map[string]any{
			"analysis_status": "failed",
			"error_message":   fmt.Sprintf("提取文档内容失败: %s", err.Error()),
		})
		return
	}

	// AI 分析
	result, err := s.analyzeWithAI(content, fileType, fullPath)
	if err != nil {
		s.repo.UpdateAnalysisResult(documentID, map[string]any{
			"analysis_status": "failed",
			"error_message":   fmt.Sprintf("AI分析失败: %s", err.Error()),
		})
		return
	}

	// 保存分析结果
	keywordsJSON, _ := json.Marshal(result.Keywords)
	zonesJSON, _ := json.Marshal(result.FunctionalZones)

	// 构建建议信息（如果有）
	var suggestionMsg string
	if len(result.Suggestions) > 0 {
		suggestionMsg = strings.Join(result.Suggestions, "; ")
	}

	updates := map[string]any{
		"analysis_status":  "completed",
		"error_message":    suggestionMsg, // 将建议信息存入，便于前端展示
		"keywords":         string(keywordsJSON),
		"summary":          result.Summary,
		"project_name":     result.ProjectName,
		"extracted_area":   result.Area,
		"extracted_budget": result.Budget,
		"extracted_style":  result.Style,
		"functional_zones": string(zonesJSON),
	}

	s.repo.UpdateAnalysisResult(documentID, updates)
}

// extractContent 根据文件类型提取文档内容
// Extract document content based on file type
func (s *DocumentService) extractContent(filePath, fileType string) (string, error) {
	if fileType == "image" {
		return "", nil
	}

	switch fileType {
	case "word":
		return s.extractDocxContent(filePath)
	case "pdf":
		return s.extractPdfContent(filePath)
	default:
		// 尝试作为纯文本读取
		data, err := os.ReadFile(filePath)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}
}

// extractDocxContent 从 .docx 文件中提取文本内容
// Extract text content from .docx file (which is a ZIP archive containing XML)
func (s *DocumentService) extractDocxContent(filePath string) (string, error) {
	// .docx 文件实际上是 ZIP 压缩包
	r, err := zip.OpenReader(filePath)
	if err != nil {
		return "", fmt.Errorf("无法打开docx文件: %w", err)
	}
	defer r.Close()

	var textContent strings.Builder

	// 遍历 ZIP 中的文件，查找 word/document.xml（主文档内容）
	for _, f := range r.File {
		if f.Name == "word/document.xml" {
			rc, err := f.Open()
			if err != nil {
				return "", fmt.Errorf("无法读取document.xml: %w", err)
			}
			defer rc.Close()

			content, err := io.ReadAll(rc)
			if err != nil {
				return "", fmt.Errorf("读取document.xml失败: %w", err)
			}

			// 解析 XML 提取文本
			text := s.extractTextFromDocxXML(content)
			textContent.WriteString(text)
			break
		}
	}

	if textContent.Len() == 0 {
		return "", errors.New("docx文件中未找到有效文本内容")
	}

	return textContent.String(), nil
}

// extractTextFromDocxXML 从 docx 的 XML 内容中提取纯文本
// Extract plain text from docx XML content
func (s *DocumentService) extractTextFromDocxXML(xmlContent []byte) string {
	var result strings.Builder
	decoder := xml.NewDecoder(bytes.NewReader(xmlContent))

	var inText bool
	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}

		switch t := token.(type) {
		case xml.StartElement:
			// <w:t> 标签包含实际文本内容
			if t.Name.Local == "t" {
				inText = true
			}
			// <w:p> 标签表示段落，添加换行
			if t.Name.Local == "p" && result.Len() > 0 {
				result.WriteString("\n")
			}
		case xml.EndElement:
			if t.Name.Local == "t" {
				inText = false
			}
		case xml.CharData:
			if inText {
				result.Write(t)
			}
		}
	}

	return result.String()
}

// extractPdfContent 从 PDF 文件中提取文本内容
// Extract text content from PDF file
func (s *DocumentService) extractPdfContent(filePath string) (string, error) {
	// 读取 PDF 文件
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("无法读取PDF文件: %w", err)
	}

	// 简单的 PDF 文本提取（提取 stream 中的文本）
	text := s.extractTextFromPdfData(data)

	if strings.TrimSpace(text) == "" {
		// 如果简单提取失败，返回提示信息让 AI 知道这是 PDF
		return "[PDF文档] 无法直接提取文本内容，可能是扫描版PDF或包含复杂格式。建议：1.上传Word版本文档 2.上传文档截图进行OCR识别", nil
	}

	return text, nil
}

// extractTextFromPdfData 从 PDF 二进制数据中提取文本
// Extract text from PDF binary data using simple pattern matching
func (s *DocumentService) extractTextFromPdfData(data []byte) string {
	var result strings.Builder
	content := string(data)

	// 方法1: 提取 BT...ET 块中的文本（PDF文本对象）
	btPattern := regexp.MustCompile(`BT\s*(.*?)\s*ET`)
	matches := btPattern.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) > 1 {
			// 提取 Tj 和 TJ 操作符中的文本
			textPattern := regexp.MustCompile(`\((.*?)\)\s*Tj|\[(.*?)\]\s*TJ`)
			textMatches := textPattern.FindAllStringSubmatch(match[1], -1)
			for _, tm := range textMatches {
				if len(tm) > 1 && tm[1] != "" {
					result.WriteString(s.decodePdfString(tm[1]))
				}
				if len(tm) > 2 && tm[2] != "" {
					// TJ 数组格式
					result.WriteString(s.extractTJArrayText(tm[2]))
				}
			}
		}
	}

	// 方法2: 尝试提取 stream 中的纯文本
	if result.Len() == 0 {
		streamPattern := regexp.MustCompile(`stream\s*([\s\S]*?)\s*endstream`)
		streamMatches := streamPattern.FindAllStringSubmatch(content, -1)
		for _, sm := range streamMatches {
			if len(sm) > 1 {
				// 过滤出可打印的中英文字符
				filtered := s.filterPrintableText(sm[1])
				if len(filtered) > 10 { // 只保留有意义的文本
					result.WriteString(filtered)
					result.WriteString("\n")
				}
			}
		}
	}

	return result.String()
}

// decodePdfString 解码 PDF 字符串（处理转义字符）
func (s *DocumentService) decodePdfString(str string) string {
	// 简单处理常见转义
	str = strings.ReplaceAll(str, "\\n", "\n")
	str = strings.ReplaceAll(str, "\\r", "\r")
	str = strings.ReplaceAll(str, "\\t", "\t")
	str = strings.ReplaceAll(str, "\\(", "(")
	str = strings.ReplaceAll(str, "\\)", ")")
	str = strings.ReplaceAll(str, "\\\\", "\\")
	return str
}

// extractTJArrayText 从 TJ 数组中提取文本
func (s *DocumentService) extractTJArrayText(tjArray string) string {
	var result strings.Builder
	// 提取括号中的文本
	pattern := regexp.MustCompile(`\((.*?)\)`)
	matches := pattern.FindAllStringSubmatch(tjArray, -1)
	for _, m := range matches {
		if len(m) > 1 {
			result.WriteString(s.decodePdfString(m[1]))
		}
	}
	return result.String()
}

// filterPrintableText 过滤出可打印的文本字符
func (s *DocumentService) filterPrintableText(data string) string {
	var result strings.Builder
	for _, r := range data {
		// 保留中文、英文、数字、常用标点
		if (r >= 0x4e00 && r <= 0x9fff) || // 中文
			(r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			r == ' ' || r == '\n' || r == '\t' ||
			r == '，' || r == '。' || r == '、' || r == '：' || r == '；' ||
			r == ',' || r == '.' || r == ':' || r == ';' ||
			r == '（' || r == '）' || r == '(' || r == ')' ||
			r == '"' || r == '"' || r == '\'' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func (s *DocumentService) analyzeWithAI(content, fileType, filePath string) (*response.DocumentAnalysisResult, error) {
	systemPrompt := `你是专业的工装设计项目文档分析助手。请分析文档内容，提取：
1. 项目名称 2. 面积(平方米) 3. 预算(元) 4. 设计风格 5. 功能分区 6. 关键字(5-10个) 7. 摘要(不超过500字)
返回JSON格式：{"projectName":"","area":0,"budget":0,"style":"","functionalZones":[],"keywords":[],"summary":"","missingFields":[],"suggestions":[]}`

	var aiResponse string
	var err error

	if fileType == "image" {
		imageBase64, encErr := volcengine.EncodeImageToBase64(filePath)
		if encErr != nil {
			return nil, fmt.Errorf("编码图片失败: %w", encErr)
		}
		aiResponse, err = s.aiClient.ChatWithImage(systemPrompt, "请分析这张项目文档图片。", imageBase64)
	} else {
		userPrompt := fmt.Sprintf("请分析以下项目文档：\n\n%s", content)
		if len(userPrompt) > 10000 {
			userPrompt = userPrompt[:10000] + "\n...(已截断)"
		}
		aiResponse, err = s.aiClient.ChatWithText(systemPrompt, userPrompt)
	}
	if err != nil {
		return nil, err
	}
	return s.parseAIResponse(aiResponse)
}

func (s *DocumentService) parseAIResponse(aiResponse string) (*response.DocumentAnalysisResult, error) {
	jsonStr := aiResponse
	if idx := strings.Index(aiResponse, "{"); idx != -1 {
		if endIdx := strings.LastIndex(aiResponse, "}"); endIdx != -1 {
			jsonStr = aiResponse[idx : endIdx+1]
		}
	}

	var result struct {
		ProjectName     string   `json:"projectName"`
		Area            float64  `json:"area"`
		Budget          float64  `json:"budget"`
		Style           string   `json:"style"`
		FunctionalZones []string `json:"functionalZones"`
		Keywords        []string `json:"keywords"`
		Summary         string   `json:"summary"`
		MissingFields   []string `json:"missingFields"`
		Suggestions     []string `json:"suggestions"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return s.extractFromText(aiResponse), nil
	}

	summary := result.Summary
	if len(summary) > 500 {
		summary = summary[:500]
	}

	return &response.DocumentAnalysisResult{
		ProjectName: result.ProjectName, Area: result.Area, Budget: result.Budget,
		Style: result.Style, FunctionalZones: result.FunctionalZones, Keywords: result.Keywords,
		Summary: summary, MissingFields: result.MissingFields, Suggestions: result.Suggestions,
	}, nil
}

func (s *DocumentService) extractFromText(text string) *response.DocumentAnalysisResult {
	result := &response.DocumentAnalysisResult{
		MissingFields: []string{}, Suggestions: []string{"建议重新上传更清晰的文档"},
	}
	areaRe := regexp.MustCompile(`(\d+(?:\.\d+)?)\s*(?:平方米|㎡|平米|m2)`)
	if matches := areaRe.FindStringSubmatch(text); len(matches) > 1 {
		result.Area, _ = strconv.ParseFloat(matches[1], 64)
	}
	budgetRe := regexp.MustCompile(`(?:预算|造价)[^\d]*(\d+(?:\.\d+)?)\s*(?:万|元)`)
	if matches := budgetRe.FindStringSubmatch(text); len(matches) > 1 {
		budget, _ := strconv.ParseFloat(matches[1], 64)
		if strings.Contains(matches[0], "万") {
			budget *= 10000
		}
		result.Budget = budget
	}
	if len(text) > 500 {
		result.Summary = text[:500]
	} else {
		result.Summary = text
	}
	return result
}

// GetSummary 获取文档摘要
// Get document summary - returns structured summary with overview, requirements, special
func (s *DocumentService) GetSummary(documentID uint, userID uint) (*response.DocumentSummaryResponse, error) {
	doc, err := s.repo.GetByIDWithProject(documentID)
	if err != nil {
		return nil, errors.New("文档不存在")
	}
	if doc.Project != nil && doc.Project.UserID != userID {
		return nil, errors.New("无权访问该文档")
	}
	if doc.AnalysisStatus != "completed" {
		return nil, errors.New("文档尚未完成分析")
	}

	// 解析摘要内容，提取项目概述、核心需求、特殊要求
	// Parse summary content to extract overview, requirements, special
	summary := doc.Summary
	result := &response.DocumentSummaryResponse{}

	// 尝试按段落分割摘要内容
	// Try to split summary by sections
	if strings.Contains(summary, "项目概述") || strings.Contains(summary, "核心需求") || strings.Contains(summary, "特殊要求") {
		// 结构化摘要格式
		lines := strings.Split(summary, "\n")
		currentSection := ""
		var sectionContent []string

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			// 检测段落标题
			if strings.Contains(line, "项目概述") || strings.Contains(line, "概述") {
				if currentSection != "" && len(sectionContent) > 0 {
					s.assignSectionContent(result, currentSection, strings.Join(sectionContent, "\n"))
				}
				currentSection = "overview"
				sectionContent = []string{}
				// 如果标题后有内容，提取它
				if idx := strings.Index(line, "："); idx != -1 {
					content := strings.TrimSpace(line[idx+3:])
					if content != "" {
						sectionContent = append(sectionContent, content)
					}
				} else if idx := strings.Index(line, ":"); idx != -1 {
					content := strings.TrimSpace(line[idx+1:])
					if content != "" {
						sectionContent = append(sectionContent, content)
					}
				}
			} else if strings.Contains(line, "核心需求") || strings.Contains(line, "需求") {
				if currentSection != "" && len(sectionContent) > 0 {
					s.assignSectionContent(result, currentSection, strings.Join(sectionContent, "\n"))
				}
				currentSection = "requirements"
				sectionContent = []string{}
				if idx := strings.Index(line, "："); idx != -1 {
					content := strings.TrimSpace(line[idx+3:])
					if content != "" {
						sectionContent = append(sectionContent, content)
					}
				} else if idx := strings.Index(line, ":"); idx != -1 {
					content := strings.TrimSpace(line[idx+1:])
					if content != "" {
						sectionContent = append(sectionContent, content)
					}
				}
			} else if strings.Contains(line, "特殊要求") || strings.Contains(line, "特殊") {
				if currentSection != "" && len(sectionContent) > 0 {
					s.assignSectionContent(result, currentSection, strings.Join(sectionContent, "\n"))
				}
				currentSection = "special"
				sectionContent = []string{}
				if idx := strings.Index(line, "："); idx != -1 {
					content := strings.TrimSpace(line[idx+3:])
					if content != "" {
						sectionContent = append(sectionContent, content)
					}
				} else if idx := strings.Index(line, ":"); idx != -1 {
					content := strings.TrimSpace(line[idx+1:])
					if content != "" {
						sectionContent = append(sectionContent, content)
					}
				}
			} else if currentSection != "" {
				sectionContent = append(sectionContent, line)
			}
		}

		// 处理最后一个段落
		if currentSection != "" && len(sectionContent) > 0 {
			s.assignSectionContent(result, currentSection, strings.Join(sectionContent, "\n"))
		}
	}

	// 如果没有解析到结构化内容，将整个摘要作为概述
	// If no structured content parsed, use entire summary as overview
	if result.Overview == "" && result.Requirements == "" && result.Special == "" {
		result.Overview = summary
	}

	return result, nil
}

// assignSectionContent 分配段落内容到对应字段
// Assign section content to corresponding field
func (s *DocumentService) assignSectionContent(result *response.DocumentSummaryResponse, section string, content string) {
	switch section {
	case "overview":
		result.Overview = content
	case "requirements":
		result.Requirements = content
	case "special":
		result.Special = content
	}
}

// GetKeywords 获取文档关键字
func (s *DocumentService) GetKeywords(documentID uint, userID uint) ([]string, error) {
	doc, err := s.repo.GetByIDWithProject(documentID)
	if err != nil {
		return nil, errors.New("文档不存在")
	}
	if doc.Project != nil && doc.Project.UserID != userID {
		return nil, errors.New("无权访问该文档")
	}
	if doc.AnalysisStatus != "completed" {
		return nil, errors.New("文档尚未完成分析")
	}
	var keywords []string
	if doc.Keywords != "" {
		json.Unmarshal([]byte(doc.Keywords), &keywords)
	}
	return keywords, nil
}

func (s *DocumentService) toResponse(doc *model.ProjectDocument) *response.DocumentResponse {
	var keywords, zones []string
	if doc.Keywords != "" {
		json.Unmarshal([]byte(doc.Keywords), &keywords)
	}
	if doc.FunctionalZones != "" {
		json.Unmarshal([]byte(doc.FunctionalZones), &zones)
	}
	return &response.DocumentResponse{
		ID: doc.ID, ProjectID: doc.ProjectID, FileName: doc.FileName, FilePath: doc.FilePath,
		FileType: doc.FileType, FileSize: doc.FileSize, AnalysisStatus: doc.AnalysisStatus,
		ErrorMessage: doc.ErrorMessage, Keywords: keywords, Summary: doc.Summary,
		ProjectName: doc.ProjectName, ExtractedArea: doc.ExtractedArea,
		ExtractedBudget: doc.ExtractedBudget, ExtractedStyle: doc.ExtractedStyle,
		FunctionalZones: zones, CreatedAt: doc.CreatedAt, UpdatedAt: doc.UpdatedAt,
	}
}
