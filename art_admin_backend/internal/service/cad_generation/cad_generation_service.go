package cad_generation

import (
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/dto/response"
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/config"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CadGenerationService CAD生成服务
// CAD generation service
type CadGenerationService struct {
	db         *gorm.DB
	uploadPath string
}

// NewCadGenerationService 创建CAD生成服务实例
// Create CAD generation service instance
func NewCadGenerationService(db *gorm.DB) *CadGenerationService {
	uploadPath := "./uploads"
	// 创建CAD生成输出目录 / Create CAD generation output directory
	cadGenPath := filepath.Join(uploadPath, "cad", "generations")
	os.MkdirAll(cadGenPath, 0755)

	return &CadGenerationService{
		db:         db,
		uploadPath: uploadPath,
	}
}

// Create 创建CAD生成任务
// Create CAD generation task
func (s *CadGenerationService) Create(req *request.CreateCadGenerationReq, userID uint) (*response.CadGenerationResp, error) {
	// 验证项目是否存在 / Verify project exists
	var project model.DesignerProject
	if err := s.db.First(&project, req.ProjectID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("项目不存在")
		}
		return nil, err
	}

	// 序列化文档ID列表 / Serialize document IDs
	documentIDsJSON, _ := json.Marshal(req.DocumentIDs)

	// 序列化参数 / Serialize parameters
	var paramsJSON []byte
	if req.Parameters != nil {
		paramsJSON, _ = json.Marshal(req.Parameters)
	} else {
		paramsJSON = []byte("{}")
	}

	// 创建任务 / Create task
	task := &model.CadGeneration{
		ProjectID:      req.ProjectID,
		DocumentIDs:    string(documentIDsJSON),
		GenerationType: req.GenerationType,
		Prompt:         req.Prompt,
		Parameters:     string(paramsJSON),
		Status:         model.CadGenStatusPending,
		Progress:       0,
		CreatedBy:      userID,
	}

	if err := s.db.Create(task).Error; err != nil {
		return nil, err
	}

	// 异步启动生成任务 / Start generation task asynchronously
	go s.processGeneration(task.ID)

	return s.toResponse(task), nil
}

// GetByID 根据ID获取任务详情
// Get task detail by ID
func (s *CadGenerationService) GetByID(id uint) (*response.CadGenerationResp, error) {
	var task model.CadGeneration
	if err := s.db.Preload("Project").Preload("ResultFile").First(&task, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("任务不存在")
		}
		return nil, err
	}
	return s.toResponse(&task), nil
}

// List 获取任务列表
// Get task list
func (s *CadGenerationService) List(req *request.CadGenerationListReq) (*response.CadGenerationListResp, error) {
	var tasks []model.CadGeneration
	var total int64

	query := s.db.Model(&model.CadGeneration{})

	// 筛选条件 / Filter conditions
	if req.ProjectID > 0 {
		query = query.Where("project_id = ?", req.ProjectID)
	}
	if req.GenerationType != "" {
		query = query.Where("generation_type = ?", req.GenerationType)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	// 统计总数 / Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询 / Paginated query
	offset := (req.Current - 1) * req.Size
	if err := query.Preload("Project").Preload("ResultFile").
		Order("created_at DESC").
		Offset(offset).Limit(req.Size).
		Find(&tasks).Error; err != nil {
		return nil, err
	}

	// 转换响应 / Convert to response
	records := make([]response.CadGenerationResp, len(tasks))
	for i, task := range tasks {
		records[i] = *s.toResponse(&task)
	}

	return &response.CadGenerationListResp{
		Records: records,
		Current: req.Current,
		Size:    req.Size,
		Total:   total,
	}, nil
}

// Delete 删除任务
// Delete task
func (s *CadGenerationService) Delete(id uint) error {
	result := s.db.Delete(&model.CadGeneration{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("任务不存在")
	}
	return nil
}

// Retry 重试生成任务
// Retry generation task
func (s *CadGenerationService) Retry(id uint) (*response.CadGenerationResp, error) {
	var task model.CadGeneration
	if err := s.db.First(&task, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("任务不存在")
		}
		return nil, err
	}

	// 只有失败的任务可以重试 / Only failed tasks can be retried
	if task.Status != model.CadGenStatusFailed {
		return nil, errors.New("只有失败的任务可以重试")
	}

	// 重置状态 / Reset status
	task.Status = model.CadGenStatusPending
	task.Progress = 0
	task.ErrorMessage = ""
	if err := s.db.Save(&task).Error; err != nil {
		return nil, err
	}

	// 异步启动生成任务 / Start generation task asynchronously
	go s.processGeneration(task.ID)

	return s.toResponse(&task), nil
}

// ConfirmAndSave 确认并保存CAD文件到项目
// Confirm and save CAD file to project
func (s *CadGenerationService) ConfirmAndSave(req *request.ConfirmCadFileReq, userID uint) (*response.CadGenerationResp, error) {
	var task model.CadGeneration
	if err := s.db.First(&task, req.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("任务不存在")
		}
		return nil, err
	}

	// 只有已完成的任务可以确认 / Only completed tasks can be confirmed
	if task.Status != model.CadGenStatusCompleted {
		return nil, errors.New("任务尚未完成")
	}

	// 如果已经有关联的CAD文件，直接返回 / If already has associated CAD file, return directly
	if task.ResultFileID != nil {
		return s.toResponse(&task), nil
	}

	// 创建CAD文件记录 / Create CAD file record
	fileName := req.FileName
	if fileName == "" {
		fileName = fmt.Sprintf("AI生成_%s_%d.dxf", response.GetGenerationTypeLabel(task.GenerationType), task.ID)
	}

	// 检查结果文件路径 / Check result file path
	if task.ResultFilePath == "" {
		return nil, errors.New("生成的CAD文件路径为空，请重试生成任务")
	}

	// 获取相对路径（去掉/uploads/前缀）/ Get relative path
	originalPath := strings.TrimPrefix(task.ResultFilePath, "/uploads/")
	if originalPath == "" {
		return nil, errors.New("无法获取CAD文件的相对路径")
	}

	// 验证文件是否存在 / Verify file exists
	fullPath := filepath.Join(s.uploadPath, originalPath)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("CAD文件不存在: %s", fullPath)
	}

	cadFile := &model.CadFile{
		ProjectID:     task.ProjectID,
		FileName:      fileName,
		FilePath:      task.ResultFilePath,
		OriginalPath:  originalPath, // 设置原始路径用于解析 / Set original path for parsing
		FileFormat:    "dxf",
		ParseStatus:   "pending", // 设置为待解析状态 / Set to pending status
		IsAIGenerated: true,
	}

	if err := s.db.Create(cadFile).Error; err != nil {
		return nil, err
	}

	// 更新任务关联 / Update task association
	task.ResultFileID = &cadFile.ID
	if err := s.db.Save(&task).Error; err != nil {
		return nil, err
	}

	// 异步解析DXF文件 / Parse DXF file asynchronously
	go s.parseSavedDxfFile(cadFile.ID, originalPath)

	return s.toResponse(&task), nil
}

// parseSavedDxfFile 解析保存的DXF文件
// Parse saved DXF file
func (s *CadGenerationService) parseSavedDxfFile(cadFileID uint, originalPath string) {
	// 更新状态为处理中 / Update status to processing
	s.db.Model(&model.CadFile{}).Where("id = ?", cadFileID).Update("parse_status", "processing")

	// 获取完整文件路径 / Get full file path
	fullPath := filepath.Join(s.uploadPath, originalPath)

	// 读取DXF文件内容 / Read DXF file content
	content, err := os.ReadFile(fullPath)
	if err != nil {
		s.db.Model(&model.CadFile{}).Where("id = ?", cadFileID).Updates(map[string]interface{}{
			"parse_status":  "failed",
			"error_message": fmt.Sprintf("读取文件失败: %v", err),
		})
		return
	}

	// 解析DXF内容 / Parse DXF content
	parseResult := s.parseGeneratedDxf(string(content))

	// 保存解析结果到JSON文件 / Save parse result to JSON file
	parsedPath, err := s.saveParsedData(cadFileID, parseResult)
	if err != nil {
		s.db.Model(&model.CadFile{}).Where("id = ?", cadFileID).Updates(map[string]interface{}{
			"parse_status":  "failed",
			"error_message": fmt.Sprintf("保存解析结果失败: %v", err),
		})
		return
	}

	// 序列化图层列表 / Serialize layer list
	layerNames := make([]string, len(parseResult.Layers))
	for i, layer := range parseResult.Layers {
		layerNames[i] = layer.Name
	}
	layersJSON, _ := json.Marshal(layerNames)

	// 更新数据库记录 / Update database record
	s.db.Model(&model.CadFile{}).Where("id = ?", cadFileID).Updates(map[string]interface{}{
		"parse_status": "completed",
		"parsed_path":  parsedPath,
		"layer_count":  len(parseResult.Layers),
		"layers":       string(layersJSON),
		"has_3d":       false,
	})
}

// parseGeneratedDxf 解析AI生成的DXF内容
// Parse AI generated DXF content
func (s *CadGenerationService) parseGeneratedDxf(content string) *DxfParseResult {
	result := &DxfParseResult{
		Layers:   []DxfLayer{},
		Entities: []DxfEntity{},
	}

	lines := strings.Split(content, "\n")
	layerMap := make(map[string]*DxfLayer)
	var currentSection string
	var entityType string
	var currentEntity *DxfEntity
	var lastCode string

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		// 解析组码 / Parse group code
		if i+1 < len(lines) {
			code := line
			value := strings.TrimSpace(lines[i+1])

			// 段标识 / Section identifier
			if code == "0" {
				if value == "SECTION" {
					// 下一行是段名 / Next line is section name
				} else if value == "ENDSEC" {
					currentSection = ""
				} else if value == "EOF" {
					break
				} else if currentSection == "ENTITIES" {
					// 新实体 / New entity
					if currentEntity != nil {
						result.Entities = append(result.Entities, *currentEntity)
					}
					entityType = value
					currentEntity = &DxfEntity{
						Type:       entityType,
						Layer:      "0",
						Color:      7,
						Properties: make(map[string]interface{}),
					}
				}
			} else if code == "2" && lastCode == "0" {
				currentSection = value
			}

			// 解析图层 / Parse layers
			if currentSection == "TABLES" && code == "2" && entityType == "LAYER" {
				if _, exists := layerMap[value]; !exists {
					layerMap[value] = &DxfLayer{Name: value, Color: 7, EntityCount: 0}
				}
			}

			// 解析实体属性 / Parse entity properties
			if currentEntity != nil {
				switch code {
				case "8": // 图层 / Layer
					currentEntity.Layer = value
					if layer, exists := layerMap[value]; exists {
						layer.EntityCount++
					} else {
						layerMap[value] = &DxfLayer{Name: value, Color: 7, EntityCount: 1}
					}
				case "62": // 颜色 / Color
					if c, err := strconv.Atoi(value); err == nil {
						currentEntity.Color = c
					}
				case "10": // X坐标 / X coordinate
					if f, err := strconv.ParseFloat(value, 64); err == nil {
						currentEntity.Properties["x0"] = f
					}
				case "20": // Y坐标 / Y coordinate
					if f, err := strconv.ParseFloat(value, 64); err == nil {
						currentEntity.Properties["y0"] = f
					}
				case "11": // X终点 / X end
					if f, err := strconv.ParseFloat(value, 64); err == nil {
						currentEntity.Properties["x1"] = f
					}
				case "21": // Y终点 / Y end
					if f, err := strconv.ParseFloat(value, 64); err == nil {
						currentEntity.Properties["y1"] = f
					}
				case "40": // 半径或文字高度 / Radius or text height
					if f, err := strconv.ParseFloat(value, 64); err == nil {
						if entityType == "CIRCLE" || entityType == "ARC" {
							currentEntity.Properties["radius"] = f
						} else if entityType == "TEXT" {
							currentEntity.Properties["height"] = f
						}
					}
				case "1": // 文字内容 / Text content
					if entityType == "TEXT" {
						currentEntity.Properties["text"] = value
					}
				case "50": // 起始角度 / Start angle
					if f, err := strconv.ParseFloat(value, 64); err == nil {
						currentEntity.Properties["startAngle"] = f
					}
				case "51": // 结束角度 / End angle
					if f, err := strconv.ParseFloat(value, 64); err == nil {
						currentEntity.Properties["endAngle"] = f
					}
				}
			}

			lastCode = code
			i++ // 跳过值行 / Skip value line
		}
	}

	// 添加最后一个实体 / Add last entity
	if currentEntity != nil {
		result.Entities = append(result.Entities, *currentEntity)
	}

	// 转换图层map为列表 / Convert layer map to list
	for _, layer := range layerMap {
		result.Layers = append(result.Layers, *layer)
	}

	// 如果没有解析到图层，添加默认图层 / Add default layer if none parsed
	if len(result.Layers) == 0 {
		result.Layers = append(result.Layers, DxfLayer{Name: "0", Color: 7, EntityCount: len(result.Entities)})
	}

	// 计算边界框 / Calculate bounding box
	result.BoundingBox = s.calculateBoundingBox(result.Entities)

	return result
}

// calculateBoundingBox 计算实体的边界框
// Calculate bounding box from entities
func (s *CadGenerationService) calculateBoundingBox(entities []DxfEntity) *BoundingBox {
	if len(entities) == 0 {
		return &BoundingBox{MinX: 0, MinY: 0, MaxX: 10000, MaxY: 10000}
	}

	minX, minY := 1e10, 1e10
	maxX, maxY := -1e10, -1e10

	for _, entity := range entities {
		props := entity.Properties

		// 获取坐标值 / Get coordinate values
		coords := []struct {
			key string
			isX bool
		}{
			{"x0", true}, {"y0", false},
			{"x1", true}, {"y1", false},
			{"centerX", true}, {"centerY", false},
		}

		for _, coord := range coords {
			if val, ok := props[coord.key]; ok {
				if f, ok := val.(float64); ok {
					if coord.isX {
						if f < minX {
							minX = f
						}
						if f > maxX {
							maxX = f
						}
					} else {
						if f < minY {
							minY = f
						}
						if f > maxY {
							maxY = f
						}
					}
				}
			}
		}

		// 处理圆的边界 / Handle circle bounds
		if entity.Type == "CIRCLE" || entity.Type == "ARC" {
			if cx, ok := props["x0"].(float64); ok {
				if cy, ok := props["y0"].(float64); ok {
					if r, ok := props["radius"].(float64); ok {
						if cx-r < minX {
							minX = cx - r
						}
						if cx+r > maxX {
							maxX = cx + r
						}
						if cy-r < minY {
							minY = cy - r
						}
						if cy+r > maxY {
							maxY = cy + r
						}
					}
				}
			}
		}
	}

	// 如果没有有效坐标，使用默认值 / Use defaults if no valid coords
	if minX > maxX || minY > maxY {
		return &BoundingBox{MinX: 0, MinY: 0, MaxX: 10000, MaxY: 10000}
	}

	return &BoundingBox{MinX: minX, MinY: minY, MaxX: maxX, MaxY: maxY}
}

// DxfParseResult DXF解析结果
type DxfParseResult struct {
	Layers      []DxfLayer   `json:"layers"`
	Entities    []DxfEntity  `json:"entities"`
	BoundingBox *BoundingBox `json:"boundingBox"`
}

// BoundingBox 边界框
type BoundingBox struct {
	MinX float64 `json:"minX"`
	MinY float64 `json:"minY"`
	MaxX float64 `json:"maxX"`
	MaxY float64 `json:"maxY"`
}

// DxfLayer DXF图层
type DxfLayer struct {
	Name        string `json:"name"`
	Color       int    `json:"color"`
	EntityCount int    `json:"entityCount"`
}

// DxfEntity DXF实体
type DxfEntity struct {
	Type       string                 `json:"type"`
	Layer      string                 `json:"layer"`
	Color      int                    `json:"color"`
	Properties map[string]interface{} `json:"properties"`
}

// saveParsedData 保存解析数据到JSON文件
// Save parsed data to JSON file
func (s *CadGenerationService) saveParsedData(cadFileID uint, result *DxfParseResult) (string, error) {
	// 生成保存路径 / Generate save path
	dateDir := time.Now().Format("2006/01/02")
	parsedDir := filepath.Join(s.uploadPath, "cad", "parsed", dateDir)
	if err := os.MkdirAll(parsedDir, 0755); err != nil {
		return "", fmt.Errorf("创建目录失败: %w", err)
	}

	fileName := fmt.Sprintf("parsed_%d_%s.json", cadFileID, uuid.New().String()[:8])
	relativePath := filepath.Join("cad", "parsed", dateDir, fileName)
	fullPath := filepath.Join(s.uploadPath, relativePath)

	// 序列化并保存 / Serialize and save
	data, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("序列化失败: %w", err)
	}

	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return "", fmt.Errorf("保存文件失败: %w", err)
	}

	return relativePath, nil
}

// processGeneration 处理生成任务
// Process generation task using AI to generate DXF CAD file
func (s *CadGenerationService) processGeneration(taskID uint) {
	var task model.CadGeneration
	if err := s.db.First(&task, taskID).Error; err != nil {
		return
	}

	// 更新状态为处理中 / Update status to processing
	task.Status = model.CadGenStatusProcessing
	task.Progress = 0
	s.db.Save(&task)

	startTime := time.Now()

	// 1. 获取项目信息 / Get project info
	task.Progress = 10
	s.db.Save(&task)

	var projectInfo string
	var project model.DesignerProject
	if err := s.db.First(&project, task.ProjectID).Error; err == nil {
		projectInfo = s.extractProjectInfo(&project)
	}

	// 2. 获取关联文档内容 / Get associated document content
	task.Progress = 20
	s.db.Save(&task)

	var documentInfo string
	if task.DocumentIDs != "" && task.DocumentIDs != "[]" {
		var docIDs []uint
		json.Unmarshal([]byte(task.DocumentIDs), &docIDs)
		documentInfo = s.extractDocumentInfo(docIDs)
	}

	// 3. 解析生成参数 / Parse generation parameters
	task.Progress = 30
	s.db.Save(&task)

	var params model.CadGenerationParams
	if task.Parameters != "" && task.Parameters != "{}" {
		json.Unmarshal([]byte(task.Parameters), &params)
	}

	// 4. 构建DXF生成提示词 / Build DXF generation prompt
	task.Progress = 40
	s.db.Save(&task)

	dxfPrompt := s.buildDxfPrompt(&task, &params, projectInfo, documentInfo)

	// 5. 调用AI生成DXF内容 / Call AI to generate DXF content
	task.Progress = 50
	s.db.Save(&task)

	dxfContent, err := s.generateDxfWithAI(dxfPrompt)
	if err != nil {
		// 生成失败 / Generation failed
		task.Status = model.CadGenStatusFailed
		task.ErrorMessage = fmt.Sprintf("DXF生成失败: %v", err)
		task.ProcessingTime = int(time.Since(startTime).Seconds())
		s.db.Save(&task)
		return
	}

	// 6. 保存DXF文件 / Save DXF file
	task.Progress = 80
	s.db.Save(&task)

	savedPath, err := s.saveDxfFile(dxfContent, taskID)
	if err != nil {
		task.Status = model.CadGenStatusFailed
		task.ErrorMessage = fmt.Sprintf("保存DXF文件失败: %v", err)
		task.ProcessingTime = int(time.Since(startTime).Seconds())
		s.db.Save(&task)
		return
	}

	// 7. 完成任务 / Complete task
	task.Status = model.CadGenStatusCompleted
	task.Progress = 100
	task.ResultFilePath = savedPath
	task.AIModel = "Doubao-pro-256k"
	task.ProcessingTime = int(time.Since(startTime).Seconds())
	s.db.Save(&task)
}

// extractProjectInfo 提取项目信息
// Extract project information for prompt
func (s *CadGenerationService) extractProjectInfo(project *model.DesignerProject) string {
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
func (s *CadGenerationService) extractDocumentInfo(documentIDs []uint) string {
	if len(documentIDs) == 0 {
		return ""
	}

	var info strings.Builder
	info.WriteString("参考文档信息:\n")

	for _, docID := range documentIDs {
		var doc model.ProjectDocument
		if err := s.db.First(&doc, docID).Error; err != nil {
			continue
		}

		info.WriteString(fmt.Sprintf("\n【%s】\n", doc.FileName))

		if doc.Summary != "" {
			info.WriteString(fmt.Sprintf("摘要: %s\n", doc.Summary))
		}
		if doc.Keywords != "" && doc.Keywords != "[]" {
			info.WriteString(fmt.Sprintf("关键词: %s\n", doc.Keywords))
		}
		if doc.ExtractedStyle != "" {
			info.WriteString(fmt.Sprintf("设计风格: %s\n", doc.ExtractedStyle))
		}
	}

	return info.String()
}

// buildDxfPrompt 构建DXF生成提示词（返回JSON参数）
// Build DXF generation prompt for AI model (return JSON parameters)
func (s *CadGenerationService) buildDxfPrompt(task *model.CadGeneration, params *model.CadGenerationParams, projectInfo, documentInfo string) string {
	// 计算尺寸（毫米）/ Calculate dimensions (mm)
	width := 10000.0 // 默认10米
	height := 8000.0 // 默认8米
	if params != nil {
		if params.Width > 0 {
			width = params.Width * 1000 // 米转毫米
		}
		if params.Height > 0 {
			height = params.Height * 1000
		}
	}

	var prompt strings.Builder

	prompt.WriteString("你是一个专业的CAD工程师，请根据以下需求生成CAD图形参数。\n\n")

	// 生成类型 / Generation type
	prompt.WriteString("## 生成类型\n")
	if task.GenerationType != "" {
		fmt.Fprintf(&prompt, "%s\n\n", task.GenerationType)
	} else {
		prompt.WriteString("室内平面布局图\n\n")
	}

	// 尺寸参数 / Dimension parameters
	prompt.WriteString("## 空间尺寸\n")
	fmt.Fprintf(&prompt, "- 宽度: %.0f 毫米\n", width)
	fmt.Fprintf(&prompt, "- 高度: %.0f 毫米\n", height)
	if params != nil && params.Style != "" {
		fmt.Fprintf(&prompt, "- 设计风格: %s\n", params.Style)
	}
	prompt.WriteString("\n")

	// 用户设计要求 / User design requirements
	if task.Prompt != "" {
		prompt.WriteString("## 设计要求\n")
		prompt.WriteString(task.Prompt)
		prompt.WriteString("\n\n")
	}

	// JSON输出格式要求 / JSON output format requirements
	fmt.Fprintf(&prompt, `## 输出要求
请严格按照以下JSON格式输出图形参数（仅输出JSON，不要其他文字，不要markdown代码块）：

{
  "units": 4,
  "layers": [
    {"name": "WALL", "color": 7},
    {"name": "DOOR", "color": 3},
    {"name": "WINDOW", "color": 5},
    {"name": "TEXT", "color": 2},
    {"name": "FURNITURE", "color": 6}
  ],
  "shapes": [
    {"type": "rectangle", "layer": "WALL", "params": {"x": 0, "y": 0, "width": %.0f, "height": %.0f}},
    {"type": "line", "layer": "WALL", "params": {"x1": 0, "y1": 0, "x2": 1000, "y2": 0}},
    {"type": "circle", "layer": "FURNITURE", "params": {"centerX": 500, "centerY": 500, "radius": 200}},
    {"type": "text", "layer": "TEXT", "params": {"x": 100, "y": 100, "height": 200, "text": "客厅"}},
    {"type": "arc", "layer": "DOOR", "params": {"centerX": 0, "centerY": 0, "radius": 800, "startAngle": 0, "endAngle": 90}}
  ]
}

图形类型说明：
- rectangle: 矩形，参数 x,y(左下角), width, height
- line: 直线，参数 x1,y1(起点), x2,y2(终点)
- circle: 圆，参数 centerX,centerY(圆心), radius(半径)
- text: 文字，参数 x,y(位置), height(文字高度), text(内容)
- arc: 圆弧，参数 centerX,centerY(圆心), radius, startAngle, endAngle(角度)
- polyline: 多段线，参数 points([[x1,y1],[x2,y2]...]), closed(是否闭合)

图层颜色：1=红，2=黄，3=绿，4=青，5=蓝，6=品红，7=白

请根据设计要求，生成合理的室内平面布局，包括：
1. 外墙轮廓（WALL图层，矩形或多段线）
2. 内部隔墙（WALL图层，直线）
3. 门的位置（DOOR图层，圆弧表示门扇开启范围）
4. 窗户位置（WINDOW图层，双线表示）
5. 功能区域标注（TEXT图层）
6. 家具示意（FURNITURE图层，简单几何图形）

坐标单位为毫米，所有图形都在第一象限（x>=0, y>=0）。
`, width, height)

	return prompt.String()
}

// generateDxfWithAI 使用AI大模型生成DXF内容
// Generate DXF content using AI to get JSON params, then generate DXF
func (s *CadGenerationService) generateDxfWithAI(prompt string) (string, error) {
	// 1. 调用AI获取JSON参数 / Call AI to get JSON parameters
	jsonContent, err := s.callAIForJSON(prompt)
	if err != nil {
		return "", fmt.Errorf("AI生成参数失败: %w", err)
	}

	// 2. 解析JSON参数 / Parse JSON parameters
	var shapeParams AIShapeParams
	if err := json.Unmarshal([]byte(jsonContent), &shapeParams); err != nil {
		return "", fmt.Errorf("解析AI返回的JSON失败: %w, 内容: %s", err, jsonContent[:min(len(jsonContent), 500)])
	}

	// 3. 使用DXF生成器生成DXF内容 / Generate DXF content using generator
	generator := NewDXFGenerator()
	dxfContent := generator.GenerateToString(&shapeParams)

	return dxfContent, nil
}

// callAIForJSON 调用AI获取JSON参数
// Call AI to get JSON parameters
func (s *CadGenerationService) callAIForJSON(prompt string) (string, error) {
	// 获取配置 / Get config
	cfg := config.GlobalConfig
	if cfg == nil {
		return "", errors.New("配置未初始化")
	}

	// 使用豆包大模型配置 / Use Doubao LLM config
	llmConfig := cfg.VolcEngine
	if llmConfig.APIKey == "" || llmConfig.Model == "" {
		return "", errors.New("火山引擎大模型API未配置，请在配置文件中设置volcengine配置")
	}

	// 构建请求体 / Build request body
	maxTokens := llmConfig.MaxTokens
	if maxTokens <= 0 || maxTokens > 8192 {
		maxTokens = 8192 // JSON参数不需要太多token
	}

	reqBody := map[string]any{
		"model": llmConfig.Model,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "你是一个专业的CAD工程师。你只输出JSON格式的图形参数，不输出任何解释说明，不使用markdown代码块。",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"max_tokens":  maxTokens,
		"temperature": 0.3,
		"stream":      false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	// 创建HTTP请求 / Create HTTP request
	apiURL := llmConfig.BaseURL
	if apiURL == "" {
		apiURL = "https://ark.cn-beijing.volces.com/api/v3"
	}
	// 确保URL以chat/completions结尾 / Ensure URL ends with chat/completions
	if !strings.HasSuffix(apiURL, "/chat/completions") {
		apiURL = strings.TrimSuffix(apiURL, "/") + "/chat/completions"
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+llmConfig.APIKey)

	// 设置超时 / Set timeout
	timeout := llmConfig.Timeout
	if timeout <= 0 {
		timeout = 60 // 1分钟超时足够
	}

	// 发送请求 / Send request
	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应 / Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	// 检查HTTP状态码 / Check HTTP status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 检查响应是否为空 / Check if response is empty
	if len(body) == 0 {
		return "", errors.New("API返回空响应")
	}

	// 解析响应 / Parse response
	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Error struct {
			Message string `json:"message"`
			Code    string `json:"code"`
		} `json:"error"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w, 原始响应: %s", err, string(body[:min(len(body), 500)]))
	}

	// 检查错误 / Check for errors
	if result.Error.Message != "" {
		return "", fmt.Errorf("API错误: %s (code: %s)", result.Error.Message, result.Error.Code)
	}

	// 获取AI返回的内容 / Get AI returned content
	if len(result.Choices) > 0 && result.Choices[0].Message.Content != "" {
		content := result.Choices[0].Message.Content
		// 清理可能的markdown代码块标记 / Clean possible markdown code block markers
		content = strings.TrimSpace(content)
		content = strings.TrimPrefix(content, "```json\n")
		content = strings.TrimPrefix(content, "```json")
		content = strings.TrimPrefix(content, "```\n")
		content = strings.TrimPrefix(content, "```")
		content = strings.TrimSuffix(content, "\n```")
		content = strings.TrimSuffix(content, "```")
		return strings.TrimSpace(content), nil
	}

	return "", errors.New("未获取到AI返回的内容")
}

// saveDxfFile 保存DXF文件
// Save DXF file to disk
func (s *CadGenerationService) saveDxfFile(dxfContent string, taskID uint) (string, error) {
	// 生成保存路径 / Generate save path
	dateDir := time.Now().Format("2006/01/02")
	genDir := filepath.Join(s.uploadPath, "cad", "generations", dateDir)
	if err := os.MkdirAll(genDir, 0755); err != nil {
		return "", fmt.Errorf("创建目录失败: %w", err)
	}

	fileName := fmt.Sprintf("gen_%d_%s.dxf", taskID, uuid.New().String()[:8])
	relativePath := filepath.Join("cad", "generations", dateDir, fileName)
	fullPath := filepath.Join(s.uploadPath, relativePath)

	// 保存DXF文件 / Save DXF file
	if err := os.WriteFile(fullPath, []byte(dxfContent), 0644); err != nil {
		return "", fmt.Errorf("保存DXF文件失败: %w", err)
	}

	// 返回相对路径（用于URL访问）/ Return relative path for URL access
	return "/uploads/" + relativePath, nil
}

// toResponse 转换为响应结构
// Convert to response structure
func (s *CadGenerationService) toResponse(task *model.CadGeneration) *response.CadGenerationResp {
	resp := &response.CadGenerationResp{
		ID:              task.ID,
		ProjectID:       task.ProjectID,
		GenerationType:  task.GenerationType,
		GenerationLabel: response.GetGenerationTypeLabel(task.GenerationType),
		Prompt:          task.Prompt,
		Status:          task.Status,
		StatusLabel:     response.GetStatusLabel(task.Status),
		Progress:        task.Progress,
		ResultFilePath:  task.ResultFilePath,
		PreviewImageURL: task.PreviewImageURL,
		ErrorMessage:    task.ErrorMessage,
		AIModel:         task.AIModel,
		ProcessingTime:  task.ProcessingTime,
		CreatedBy:       task.CreatedBy,
		CreatedAt:       task.CreatedAt,
		UpdatedAt:       task.UpdatedAt,
	}

	// 解析文档ID列表 / Parse document IDs
	if task.DocumentIDs != "" {
		var docIDs []uint
		json.Unmarshal([]byte(task.DocumentIDs), &docIDs)
		resp.DocumentIDs = docIDs
	}

	// 解析参数 / Parse parameters
	if task.Parameters != "" && task.Parameters != "{}" {
		var params response.CadGenerationParamsResp
		json.Unmarshal([]byte(task.Parameters), &params)
		resp.Parameters = &params
	}

	// 项目名称 / Project name
	if task.Project != nil {
		resp.ProjectName = task.Project.Name
	}

	// 结果文件 / Result file
	if task.ResultFileID != nil {
		resp.ResultFileID = task.ResultFileID
		if task.ResultFile != nil {
			resp.ResultFile = &response.CadFileBriefResp{
				ID:         task.ResultFile.ID,
				FileName:   task.ResultFile.FileName,
				FileFormat: task.ResultFile.FileFormat,
				FilePath:   task.ResultFile.FilePath,
			}
		}
	}

	return resp
}
