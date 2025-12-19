package cad

import (
	"art_admin_backend/internal/dto/response"
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/volcengine"
	"art_admin_backend/internal/repository"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// RenderService CAD效果图渲染服务
// CAD rendering service for generating effect images
type RenderService struct {
	cadRepo      *repository.CadFileRepository
	projectRepo  *repository.DesignerProjectRepository
	documentRepo *repository.DocumentRepository // 文档仓库 / Document repository
	volcClient   *volcengine.Client             // 火山引擎AI客户端(文本) / VolcEngine AI client for text
	uploadPath   string
	imageAPIKey  string // 图片生成API Key / Image generation API key
	imageAPIURL  string // 图片生成API URL / Image generation API URL
	imageModel   string // 图片生成模型ID / Image generation model ID
	imageSize    string // 图片尺寸 / Image size
	imageTimeout int    // 图片生成超时时间 / Image generation timeout
}

// NewRenderService 创建渲染服务
// Create render service with all required repositories
func NewRenderService(
	cadRepo *repository.CadFileRepository,
	projectRepo *repository.DesignerProjectRepository,
	documentRepo *repository.DocumentRepository,
	volcClient *volcengine.Client,
	uploadPath string,
	imageAPIKey, imageAPIURL, imageModel, imageSize string,
	imageTimeout int,
) *RenderService {
	// 创建渲染输出目录 / Create render output directory
	renderPath := filepath.Join(uploadPath, "cad", "renders")
	os.MkdirAll(renderPath, 0755)

	// 设置默认值 / Set defaults
	if imageSize == "" {
		imageSize = "1024x1024"
	}
	if imageTimeout <= 0 {
		imageTimeout = 180
	}

	return &RenderService{
		cadRepo:      cadRepo,
		projectRepo:  projectRepo,
		documentRepo: documentRepo,
		volcClient:   volcClient,
		uploadPath:   uploadPath,
		imageAPIKey:  imageAPIKey,
		imageAPIURL:  imageAPIURL,
		imageModel:   imageModel,
		imageSize:    imageSize,
		imageTimeout: imageTimeout,
	}
}

// GenerateRenderRequest 生成效果图请求
// Generate render request with user inputs
type GenerateRenderRequest struct {
	CadFileID   uint   `json:"cadFileId"`             // CAD文件ID
	ProjectID   uint   `json:"projectId"`             // 项目ID
	DocumentIDs []uint `json:"documentIds,omitempty"` // 参考文档ID列表 / Reference document IDs
	Style       string `json:"style"`                 // 设计风格（用户自由输入）/ Design style (free text)
	RoomType    string `json:"roomType"`              // 空间类型（用户自由输入）/ Room type (free text)
	Description string `json:"description"`           // 设计要求描述 / Design requirements
}

// GenerateRender 生成效果图和设计方案
// Generate rendering and design proposal from CAD file with project documents
func (s *RenderService) GenerateRender(req *GenerateRenderRequest) (*response.RenderResult, error) {
	// 获取CAD文件信息 / Get CAD file info
	cadFile, err := s.cadRepo.GetByID(req.CadFileID)
	if err != nil {
		return nil, fmt.Errorf("CAD文件不存在: %w", err)
	}

	// 获取项目信息 / Get project info
	var projectInfo string
	if req.ProjectID > 0 {
		project, err := s.projectRepo.GetByID(req.ProjectID)
		if err == nil {
			projectInfo = s.extractProjectInfo(project)
		}
	}

	// 获取参考文档内容 / Get reference document content
	var documentInfo string
	if len(req.DocumentIDs) > 0 && s.documentRepo != nil {
		documentInfo = s.extractDocumentInfo(req.DocumentIDs)
	}

	// 生成图片提示词 / Build image prompt
	imagePrompt := s.buildImagePrompt(req, cadFile, projectInfo, documentInfo)

	// 生成设计方案提示词 / Build design proposal prompt
	designPrompt := s.buildPrompt(req, cadFile, projectInfo, documentInfo)

	// 使用Doubao-1.5-vision-pro生成设计方案 / Generate design proposal using Doubao
	designProposal, err := s.generateDesignProposal(designPrompt, cadFile.ParsedPath)
	if err != nil {
		// 如果AI调用失败，返回空设计方案 / Return empty proposal if AI fails
		designProposal = ""
	}

	// 使用Doubao-Seedream-4.5生成效果图 / Generate image using Doubao-Seedream-4.5
	imageURL, err := s.generateImage(imagePrompt)
	if err != nil {
		// 如果图片生成失败，使用占位图 / Use placeholder if image generation fails
		imageURL, _ = s.generatePlaceholder()
	}

	// 下载并保存图像 / Download and save image
	savedPath, err := s.saveImage(imageURL, req.CadFileID)
	if err != nil {
		// 如果保存失败，直接使用URL / Use URL directly if save fails
		savedPath = imageURL
	}

	return &response.RenderResult{
		CadFileID:      req.CadFileID,
		ImageURL:       savedPath,
		DesignProposal: designProposal,
		Style:          req.Style,
		RoomType:       req.RoomType,
		Prompt:         imagePrompt,
		GeneratedAt:    time.Now(),
	}, nil
}

// extractProjectInfo 提取项目信息
// Extract project information for prompt
func (s *RenderService) extractProjectInfo(project *model.DesignerProject) string {
	var info strings.Builder
	info.WriteString(fmt.Sprintf("项目名称: %s\n", project.Name))
	if project.Description != "" {
		info.WriteString(fmt.Sprintf("项目描述: %s\n", project.Description))
	}
	if project.Style != "" {
		info.WriteString(fmt.Sprintf("设计风格: %s\n", project.Style))
	}
	if project.Budget > 0 {
		info.WriteString(fmt.Sprintf("预算: %.0f元\n", project.Budget))
	}
	return info.String()
}

// extractDocumentInfo 提取文档信息
// Extract document information for prompt building
func (s *RenderService) extractDocumentInfo(documentIDs []uint) string {
	if s.documentRepo == nil || len(documentIDs) == 0 {
		return ""
	}

	var info strings.Builder
	info.WriteString("参考文档信息:\n")

	for _, docID := range documentIDs {
		doc, err := s.documentRepo.GetByID(docID)
		if err != nil {
			continue
		}

		info.WriteString(fmt.Sprintf("\n【%s】\n", doc.FileName))

		// 添加文档摘要 / Add document summary
		if doc.Summary != "" {
			info.WriteString(fmt.Sprintf("摘要: %s\n", doc.Summary))
		}

		// 添加关键词 / Add keywords
		if doc.Keywords != "" && doc.Keywords != "[]" {
			info.WriteString(fmt.Sprintf("关键词: %s\n", doc.Keywords))
		}

		// 添加提取的设计风格 / Add extracted style
		if doc.ExtractedStyle != "" {
			info.WriteString(fmt.Sprintf("设计风格: %s\n", doc.ExtractedStyle))
		}

		// 添加功能分区 / Add functional zones
		if doc.FunctionalZones != "" && doc.FunctionalZones != "[]" {
			info.WriteString(fmt.Sprintf("功能分区: %s\n", doc.FunctionalZones))
		}

		// 添加面积和预算信息 / Add area and budget info
		if doc.ExtractedArea > 0 {
			info.WriteString(fmt.Sprintf("面积: %.2f平方米\n", doc.ExtractedArea))
		}
		if doc.ExtractedBudget > 0 {
			info.WriteString(fmt.Sprintf("预算: %.0f元\n", doc.ExtractedBudget))
		}
	}

	return info.String()
}

// buildImagePrompt 构建图片生成提示词（用于Doubao-Seedream-4.5）
// Build image generation prompt for Doubao-Seedream-4.5
func (s *RenderService) buildImagePrompt(req *GenerateRenderRequest, cadFile *model.CadFile, projectInfo, documentInfo string) string {
	// 使用用户输入的风格 / Use user input style
	style := req.Style
	if style == "" {
		style = "现代简约"
	}

	// 使用用户输入的空间类型 / Use user input room type
	roomType := req.RoomType
	if roomType == "" {
		roomType = "客厅"
	}

	// 构建简洁有效的图片生成提示词 / Build concise image generation prompt
	var promptParts []string

	// 基础描述 / Base description
	promptParts = append(promptParts, fmt.Sprintf("%s室内设计效果图", roomType))
	promptParts = append(promptParts, fmt.Sprintf("%s风格", style))

	// 添加用户描述的关键词 / Add keywords from user description
	if req.Description != "" {
		promptParts = append(promptParts, req.Description)
	}

	// 添加渲染质量要求 / Add rendering quality requirements
	promptParts = append(promptParts,
		"高品质3D渲染",
		"8K超清画质",
		"专业室内摄影",
		"真实感光照",
		"高质量材质纹理",
		"透视图视角",
		"温馨舒适氛围",
		"无水印无文字",
	)

	return strings.Join(promptParts, "，")
}

// buildPrompt 构建设计方案提示词（用于Doubao-1.5-vision-pro）
// Build design proposal prompt for Doubao-1.5-vision-pro
func (s *RenderService) buildPrompt(req *GenerateRenderRequest, cadFile *model.CadFile, projectInfo, documentInfo string) string {
	// 使用用户输入的风格，如果为空则使用默认值
	// Use user input style, default if empty
	style := req.Style
	if style == "" {
		style = "现代简约风格"
	}

	// 使用用户输入的空间类型，如果为空则使用默认值
	// Use user input room type, default if empty
	roomType := req.RoomType
	if roomType == "" {
		roomType = "室内空间"
	}

	prompt := fmt.Sprintf(`请为以下室内空间生成详细的设计方案：

空间类型: %s
设计风格: %s

CAD图纸信息:
- 文件名: %s
- 图层数: %d

`, roomType, style, cadFile.FileName, cadFile.LayerCount)

	// 添加项目信息 / Add project info
	if projectInfo != "" {
		prompt += fmt.Sprintf("项目信息:\n%s\n", projectInfo)
	}

	// 添加参考文档信息 / Add reference document info
	if documentInfo != "" {
		prompt += fmt.Sprintf("\n%s\n", documentInfo)
	}

	// 添加用户描述的设计要求 / Add user's design requirements
	if req.Description != "" {
		prompt += fmt.Sprintf("\n用户设计要求:\n%s\n", req.Description)
	}

	return prompt
}

// generateDesignProposal 使用Doubao-1.5-vision-pro生成设计方案
// Generate design proposal using Doubao-1.5-vision-pro model
func (s *RenderService) generateDesignProposal(prompt, cadFilePath string) (string, error) {
	// 如果没有配置AI客户端，返回空 / Return empty if no AI client
	if s.volcClient == nil {
		return "", fmt.Errorf("AI客户端未配置")
	}

	// 构建系统提示词 / Build system prompt
	systemPrompt := `你是一位专业的室内设计师和效果图描述专家。
根据用户提供的CAD图纸信息和设计要求，生成一份详细的室内设计方案。

请按以下格式输出设计方案：

## 空间布局
- 描述整体空间布局和功能分区

## 设计风格
- 详细描述设计风格特点和元素

## 色彩方案
- 主色调、辅助色、点缀色的搭配

## 材质选择
- 地面、墙面、天花板的材质
- 家具和软装的材质

## 家具配置
- 主要家具的款式、尺寸、摆放位置

## 灯光设计
- 自然光利用
- 人工照明方案

## 装饰细节
- 绿植、艺术品、软装配饰

请提供专业、详细、可执行的设计方案。`

	// 检查是否有CAD文件可以作为参考图片 / Check if CAD file can be used as reference
	// 只支持图片格式：png, jpg, jpeg, gif, webp / Only support image formats
	var imageBase64 string
	if cadFilePath != "" {
		ext := strings.ToLower(filepath.Ext(cadFilePath))
		supportedFormats := map[string]bool{
			".png": true, ".jpg": true, ".jpeg": true, ".gif": true, ".webp": true,
		}
		if supportedFormats[ext] {
			fullPath := filepath.Join(s.uploadPath, cadFilePath)
			if data, err := os.ReadFile(fullPath); err == nil {
				imageBase64 = base64.StdEncoding.EncodeToString(data)
			}
		}
	}

	var result string
	var err error

	if imageBase64 != "" {
		// 带图片的请求 / Request with image
		result, err = s.volcClient.ChatWithImage(systemPrompt, prompt, imageBase64)
	} else {
		// 纯文本请求 / Text-only request
		result, err = s.volcClient.ChatWithText(systemPrompt, prompt)
	}

	if err != nil {
		return "", fmt.Errorf("AI生成设计方案失败: %w", err)
	}

	// 保存设计方案到文件 / Save design proposal to file
	_ = s.saveDescription(result)

	return result, nil
}

// saveDescription 保存AI生成的设计描述
// Save AI-generated design description
func (s *RenderService) saveDescription(description string) string {
	dateDir := time.Now().Format("2006/01/02")
	descDir := filepath.Join(s.uploadPath, "cad", "renders", dateDir)
	os.MkdirAll(descDir, 0755)

	fileName := fmt.Sprintf("desc_%s.txt", uuid.New().String()[:8])
	fullPath := filepath.Join(descDir, fileName)

	os.WriteFile(fullPath, []byte(description), 0644)

	return filepath.Join("cad", "renders", dateDir, fileName)
}

// generateImage 使用Doubao-Seedream-4.5生成效果图
// Generate image using Doubao-Seedream-4.5 model
func (s *RenderService) generateImage(prompt string) (string, error) {
	// 检查是否配置了图片生成API / Check if image API is configured
	if s.imageAPIKey == "" || s.imageModel == "" {
		return s.generatePlaceholder()
	}

	// 构建请求体 / Build request body
	reqBody := map[string]interface{}{
		"model":                       s.imageModel,
		"prompt":                      prompt,
		"sequential_image_generation": "disabled",
		"response_format":             "url",
		"size":                        s.imageSize,
		"stream":                      false,
		"watermark":                   false, // 不添加水印 / No watermark
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return s.generatePlaceholder()
	}

	// 创建HTTP请求 / Create HTTP request
	apiURL := s.imageAPIURL
	if apiURL == "" {
		apiURL = "https://ark.cn-beijing.volces.com/api/v3/images/generations"
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return s.generatePlaceholder()
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.imageAPIKey)

	// 发送请求 / Send request
	client := &http.Client{Timeout: time.Duration(s.imageTimeout) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return s.generatePlaceholder()
	}
	defer resp.Body.Close()

	// 读取响应 / Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return s.generatePlaceholder()
	}

	// 解析响应 / Parse response
	var result struct {
		Data []struct {
			URL string `json:"url"`
		} `json:"data"`
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return s.generatePlaceholder()
	}

	// 检查错误 / Check for errors
	if result.Error.Message != "" {
		return s.generatePlaceholder()
	}

	// 返回生成的图片URL / Return generated image URL
	if len(result.Data) > 0 && result.Data[0].URL != "" {
		return result.Data[0].URL, nil
	}

	return s.generatePlaceholder()
}

// generatePlaceholder 生成占位图（当AI服务不可用时）
// Generate placeholder image when AI service is unavailable
func (s *RenderService) generatePlaceholder() (string, error) {
	// 返回一个示例效果图URL / Return placeholder URL
	placeholderURL := "https://via.placeholder.com/1024x768/4A90D9/FFFFFF?text=AI+Rendering+Preview"
	return placeholderURL, nil
}

// saveImage 下载并保存图像
func (s *RenderService) saveImage(imageURL string, cadFileID uint) (string, error) {
	// 下载图像
	resp, err := http.Get(imageURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取图像数据
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 生成保存路径
	dateDir := time.Now().Format("2006/01/02")
	renderDir := filepath.Join(s.uploadPath, "cad", "renders", dateDir)
	os.MkdirAll(renderDir, 0755)

	fileName := fmt.Sprintf("%d_%s.png", cadFileID, uuid.New().String()[:8])
	relativePath := filepath.Join("cad", "renders", dateDir, fileName)
	fullPath := filepath.Join(s.uploadPath, relativePath)

	// 保存文件
	if err := os.WriteFile(fullPath, imageData, 0644); err != nil {
		return "", err
	}

	return relativePath, nil
}

// GetRenderHistory 获取渲染历史
func (s *RenderService) GetRenderHistory(cadFileID uint) ([]response.RenderResult, error) {
	// 扫描渲染目录获取历史记录
	var results []response.RenderResult

	renderDir := filepath.Join(s.uploadPath, "cad", "renders")
	prefix := fmt.Sprintf("%d_", cadFileID)

	filepath.Walk(renderDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if strings.HasPrefix(info.Name(), prefix) {
			relPath, _ := filepath.Rel(s.uploadPath, path)
			results = append(results, response.RenderResult{
				CadFileID:   cadFileID,
				ImageURL:    relPath,
				GeneratedAt: info.ModTime(),
			})
		}
		return nil
	})

	return results, nil
}

// encodeImageToBase64 将图像编码为base64（用于某些API）
func encodeImageToBase64(imagePath string) (string, error) {
	data, err := os.ReadFile(imagePath)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}
