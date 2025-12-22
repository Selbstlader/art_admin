package construction_annotation

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/volcengine"
	"art_admin_backend/internal/repository"
)

// AnnotationService 施工图标注服务
// Construction annotation service for element recognition and annotation generation
type AnnotationService struct {
	annotationRepo repository.ConstructionAnnotationRepository
	cadFileRepo    *repository.CadFileRepository
	volcClient     *volcengine.Client
}

// NewAnnotationService 创建标注服务实例
// Create annotation service instance
func NewAnnotationService(
	annotationRepo repository.ConstructionAnnotationRepository,
	cadFileRepo *repository.CadFileRepository,
	volcClient *volcengine.Client,
) *AnnotationService {
	return &AnnotationService{
		annotationRepo: annotationRepo,
		cadFileRepo:    cadFileRepo,
		volcClient:     volcClient,
	}
}

// AnalyzeConstructionDrawing 分析施工图并识别元素
// Analyze construction drawing and recognize elements
func (s *AnnotationService) AnalyzeConstructionDrawing(projectID, cadFileID uint, imagePath string) (*model.ConstructionAnnotation, error) {
	// 创建标注记录 / Create annotation record
	annotation := &model.ConstructionAnnotation{
		ProjectID:        projectID,
		CadFileID:        cadFileID,
		ImagePath:        imagePath,
		AnalysisStatus:   "processing",
		ElementsDetected: "[]", // 初始化为空JSON数组 / Initialize as empty JSON array
		Annotations:      "[]", // 初始化为空JSON数组 / Initialize as empty JSON array
	}

	// 保存到数据库 / Save to database
	if err := s.annotationRepo.Create(annotation); err != nil {
		return nil, fmt.Errorf("创建标注记录失败: %w", err)
	}

	// 异步处理图像分析 / Process image analysis asynchronously
	go s.processImageAnalysis(annotation.ID)

	return annotation, nil
}

// processImageAnalysis 处理图像分析
// Process image analysis
func (s *AnnotationService) processImageAnalysis(annotationID uint) {
	annotation, err := s.annotationRepo.GetByID(annotationID)
	if err != nil {
		log.Printf("获取标注记录失败: %v", err)
		return
	}

	// 将URL转换为本地文件路径 / Convert URL to local file path
	imagePath := annotation.ImagePath
	// 如果是完整URL，提取相对路径 / If full URL, extract relative path
	if strings.HasPrefix(imagePath, "http://") || strings.HasPrefix(imagePath, "https://") {
		// 提取 /uploads/ 之后的路径 / Extract path after /uploads/
		if idx := strings.Index(imagePath, "/uploads/"); idx != -1 {
			imagePath = "." + imagePath[idx:]
		}
	}

	// 编码图片为Base64 / Encode image to Base64
	imageBase64, err := volcengine.EncodeImageToBase64(imagePath)
	if err != nil {
		s.updateAnalysisStatus(annotationID, "failed", fmt.Sprintf("图片编码失败: %v", err))
		return
	}

	// 构建分析提示词 / Build analysis prompt
	systemPrompt := s.buildElementRecognitionPrompt()
	userPrompt := "请分析这张施工图，识别其中的尺寸线、材料区域、设备位置等元素。"

	// 调用火山AI进行视觉分析 / Call VolcEngine AI for visual analysis
	result, err := s.volcClient.ChatWithImage(systemPrompt, userPrompt, imageBase64)
	if err != nil {
		s.updateAnalysisStatus(annotationID, "failed", fmt.Sprintf("AI分析失败: %v", err))
		return
	}

	// 解析分析结果 / Parse analysis result
	elements, err := s.parseAnalysisResult(result)
	if err != nil {
		s.updateAnalysisStatus(annotationID, "failed", fmt.Sprintf("解析分析结果失败: %v", err))
		return
	}

	// 生成标注 / Generate annotations
	annotations := s.generateAnnotations(elements)

	// 更新数据库 / Update database
	elementsJSON, _ := json.Marshal(elements)
	annotationsJSON, _ := json.Marshal(annotations)

	updateData := map[string]interface{}{
		"analysis_status":   "completed",
		"elements_detected": string(elementsJSON),
		"annotations":       string(annotationsJSON),
		"error_message":     "",
	}

	if err := s.annotationRepo.Update(annotationID, updateData); err != nil {
		log.Printf("更新标注记录失败: %v", err)
	}

	log.Printf("施工图分析完成 - AnnotationID: %d, 检测到%d个元素", annotationID, len(elements))
}

// buildElementRecognitionPrompt 构建元素识别提示词
// Build element recognition prompt
func (s *AnnotationService) buildElementRecognitionPrompt() string {
	return `你是一个专业的施工图分析专家。请分析施工图并识别以下类型的元素：

1. 尺寸线 (dimension_line)：
   - 识别图中的尺寸标注线
   - 提取尺寸数值和单位
   - 记录尺寸线的位置和方向

2. 材料区域 (material_area)：
   - 识别不同材料的区域
   - 判断材料类型（如墙体、地面、天花等）
   - 记录材料区域的边界

3. 设备位置 (equipment)：
   - 识别设备符号和位置
   - 判断设备类型（如插座、开关、灯具等）
   - 记录设备的具体坐标

请以JSON格式返回分析结果，包含以下字段：
- type: 元素类型
- confidence: 置信度 (0-1)
- boundingBox: 边界框坐标 {x, y, width, height}
- properties: 元素属性（如尺寸值、材料类型、设备类型等）

示例格式：
{
  "elements": [
    {
      "type": "dimension_line",
      "confidence": 0.95,
      "boundingBox": {"x": 100, "y": 200, "width": 150, "height": 20},
      "properties": {"value": "3000", "unit": "mm", "direction": "horizontal"}
    }
  ]
}`
}

// parseAnalysisResult 解析AI分析结果
// Parse AI analysis result
func (s *AnnotationService) parseAnalysisResult(result string) ([]model.DetectedElement, error) {
	// 尝试提取JSON部分 / Try to extract JSON part
	jsonStart := strings.Index(result, "{")
	jsonEnd := strings.LastIndex(result, "}")
	if jsonStart == -1 || jsonEnd == -1 {
		return nil, fmt.Errorf("无法找到有效的JSON格式")
	}

	jsonStr := result[jsonStart : jsonEnd+1]

	// 解析JSON / Parse JSON
	var analysisResult struct {
		Elements []model.DetectedElement `json:"elements"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &analysisResult); err != nil {
		return nil, fmt.Errorf("JSON解析失败: %w", err)
	}

	return analysisResult.Elements, nil
}

// generateAnnotations 根据检测到的元素生成标注
// Generate annotations based on detected elements
func (s *AnnotationService) generateAnnotations(elements []model.DetectedElement) []model.AnnotationItem {
	var annotations []model.AnnotationItem
	now := time.Now()

	for i, element := range elements {
		annotation := model.AnnotationItem{
			ID:          fmt.Sprintf("auto_%d_%d", now.Unix(), i),
			Type:        s.mapElementTypeToAnnotationType(element.Type),
			Position:    s.calculateAnnotationPosition(element.BoundingBox),
			Content:     s.generateAnnotationContent(element),
			Style:       s.getDefaultAnnotationStyle(element.Type),
			Properties:  element.Properties,
			CreatedAt:   now,
			UpdatedAt:   now,
			IsGenerated: true,
		}
		annotations = append(annotations, annotation)
	}

	return annotations
}

// mapElementTypeToAnnotationType 映射元素类型到标注类型
// Map element type to annotation type
func (s *AnnotationService) mapElementTypeToAnnotationType(elementType string) string {
	switch elementType {
	case "dimension_line":
		return "dimension"
	case "material_area":
		return "material"
	case "equipment":
		return "process"
	default:
		return "general"
	}
}

// calculateAnnotationPosition 计算标注位置
// Calculate annotation position
func (s *AnnotationService) calculateAnnotationPosition(bbox model.BoundingBox) model.Position {
	// 标注位置通常在元素的右上角 / Annotation position usually at top-right of element
	return model.Position{
		X:      bbox.X + bbox.Width + 10, // 偏移10像素 / Offset 10 pixels
		Y:      bbox.Y - 5,               // 向上偏移5像素 / Offset up 5 pixels
		Anchor: "top-left",
	}
}

// generateAnnotationContent 生成标注内容
// Generate annotation content
func (s *AnnotationService) generateAnnotationContent(element model.DetectedElement) string {
	switch element.Type {
	case "dimension_line":
		if value, ok := element.Properties["value"].(string); ok {
			if unit, ok := element.Properties["unit"].(string); ok {
				return fmt.Sprintf("%s%s", value, unit)
			}
			return value
		}
		return "尺寸"
	case "material_area":
		if materialType, ok := element.Properties["materialType"].(string); ok {
			return materialType
		}
		return "材料区域"
	case "equipment":
		if equipmentType, ok := element.Properties["equipmentType"].(string); ok {
			return equipmentType
		}
		return "设备"
	default:
		return "标注"
	}
}

// getDefaultAnnotationStyle 获取默认标注样式
// Get default annotation style
func (s *AnnotationService) getDefaultAnnotationStyle(elementType string) model.AnnotationStyle {
	switch elementType {
	case "dimension_line":
		return model.AnnotationStyle{
			FontSize:   12,
			FontColor:  "#FF0000", // 红色用于尺寸标注 / Red for dimension annotations
			LineColor:  "#FF0000",
			LineWidth:  1,
			Background: "#FFFFFF",
		}
	case "material_area":
		return model.AnnotationStyle{
			FontSize:   10,
			FontColor:  "#0000FF", // 蓝色用于材料标注 / Blue for material annotations
			LineColor:  "#0000FF",
			LineWidth:  1,
			Background: "#FFFFFF",
		}
	case "equipment":
		return model.AnnotationStyle{
			FontSize:   10,
			FontColor:  "#00AA00", // 绿色用于设备标注 / Green for equipment annotations
			LineColor:  "#00AA00",
			LineWidth:  1,
			Background: "#FFFFFF",
		}
	default:
		return model.AnnotationStyle{
			FontSize:   10,
			FontColor:  "#000000",
			LineColor:  "#000000",
			LineWidth:  1,
			Background: "#FFFFFF",
		}
	}
}

// updateAnalysisStatus 更新分析状态
// Update analysis status
func (s *AnnotationService) updateAnalysisStatus(annotationID uint, status, errorMsg string) {
	updateData := map[string]interface{}{
		"analysis_status": status,
		"error_message":   errorMsg,
	}
	if err := s.annotationRepo.Update(annotationID, updateData); err != nil {
		log.Printf("更新分析状态失败: %v", err)
	}
}

// GetAnnotationByID 根据ID获取标注
// Get annotation by ID
func (s *AnnotationService) GetAnnotationByID(id uint) (*model.ConstructionAnnotation, error) {
	return s.annotationRepo.GetByID(id)
}

// GetAnnotationsByProject 获取项目的所有标注
// Get all annotations for a project
func (s *AnnotationService) GetAnnotationsByProject(projectID uint) ([]*model.ConstructionAnnotation, error) {
	return s.annotationRepo.GetByProjectID(projectID)
}

// GetAnnotationsByCadFile 获取CAD文件的所有标注
// Get all annotations for a CAD file
func (s *AnnotationService) GetAnnotationsByCadFile(cadFileID uint) ([]*model.ConstructionAnnotation, error) {
	return s.annotationRepo.GetByCadFileID(cadFileID)
}

// UpdateAnnotations 更新标注数据
// Update annotation data
func (s *AnnotationService) UpdateAnnotations(annotationID uint, annotations []model.AnnotationItem) error {
	annotationsJSON, err := json.Marshal(annotations)
	if err != nil {
		return fmt.Errorf("序列化标注数据失败: %w", err)
	}

	updateData := map[string]interface{}{
		"annotations": string(annotationsJSON),
		"updated_at":  time.Now(),
	}

	return s.annotationRepo.Update(annotationID, updateData)
}

// DeleteAnnotation 删除标注
// Delete annotation
func (s *AnnotationService) DeleteAnnotation(id uint) error {
	return s.annotationRepo.Delete(id)
}

// ValidateImageFormat 验证图片格式
// Validate image format
func (s *AnnotationService) ValidateImageFormat(imagePath string) error {
	ext := strings.ToLower(filepath.Ext(imagePath))
	validExts := []string{".jpg", ".jpeg", ".png", ".bmp", ".tiff"}

	for _, validExt := range validExts {
		if ext == validExt {
			return nil
		}
	}

	return fmt.Errorf("不支持的图片格式: %s，支持的格式: %v", ext, validExts)
}

// GetAnnotationsList 分页获取标注列表
// Get paginated annotations list
func (s *AnnotationService) GetAnnotationsList(offset, limit int) ([]*model.ConstructionAnnotation, int64, error) {
	return s.annotationRepo.List(offset, limit)
}

// GenerateAnnotationsFromElements 根据检测元素生成标注
// Generate annotations from detected elements
func (s *AnnotationService) GenerateAnnotationsFromElements(annotationID uint, elements []model.DetectedElement) error {
	// 生成标注 / Generate annotations
	annotations := s.generateAnnotations(elements)

	// 更新数据库 / Update database
	return s.UpdateAnnotations(annotationID, annotations)
}

// AddAnnotationItem 添加标注项
// Add annotation item
func (s *AnnotationService) AddAnnotationItem(annotationID uint, item model.AnnotationItem) error {
	// 获取现有标注 / Get existing annotations
	annotation, err := s.annotationRepo.GetByID(annotationID)
	if err != nil {
		return fmt.Errorf("获取标注记录失败: %w", err)
	}

	// 解析现有标注 / Parse existing annotations
	var annotations []model.AnnotationItem
	if annotation.Annotations != "" {
		if err := json.Unmarshal([]byte(annotation.Annotations), &annotations); err != nil {
			return fmt.Errorf("解析现有标注失败: %w", err)
		}
	}

	// 设置创建时间和更新时间 / Set creation and update time
	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now
	item.IsGenerated = false // 手动添加的标注 / Manually added annotation

	// 生成唯一ID / Generate unique ID
	if item.ID == "" {
		item.ID = fmt.Sprintf("manual_%d_%d", now.Unix(), len(annotations))
	}

	// 添加新标注项 / Add new annotation item
	annotations = append(annotations, item)

	// 更新数据库 / Update database
	return s.UpdateAnnotations(annotationID, annotations)
}

// UpdateAnnotationItem 更新标注项
// Update annotation item
func (s *AnnotationService) UpdateAnnotationItem(annotationID uint, itemID string, updates model.AnnotationItem) error {
	// 获取现有标注 / Get existing annotations
	annotation, err := s.annotationRepo.GetByID(annotationID)
	if err != nil {
		return fmt.Errorf("获取标注记录失败: %w", err)
	}

	// 解析现有标注 / Parse existing annotations
	var annotations []model.AnnotationItem
	if annotation.Annotations != "" {
		if err := json.Unmarshal([]byte(annotation.Annotations), &annotations); err != nil {
			return fmt.Errorf("解析现有标注失败: %w", err)
		}
	}

	// 查找并更新标注项 / Find and update annotation item
	found := false
	for i, item := range annotations {
		if item.ID == itemID {
			// 保留原有的创建时间和ID / Keep original creation time and ID
			updates.ID = item.ID
			updates.CreatedAt = item.CreatedAt
			updates.UpdatedAt = time.Now()
			annotations[i] = updates
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("标注项不存在: %s", itemID)
	}

	// 更新数据库 / Update database
	return s.UpdateAnnotations(annotationID, annotations)
}

// DeleteAnnotationItem 删除标注项
// Delete annotation item
func (s *AnnotationService) DeleteAnnotationItem(annotationID uint, itemID string) error {
	// 获取现有标注 / Get existing annotations
	annotation, err := s.annotationRepo.GetByID(annotationID)
	if err != nil {
		return fmt.Errorf("获取标注记录失败: %w", err)
	}

	// 解析现有标注 / Parse existing annotations
	var annotations []model.AnnotationItem
	if annotation.Annotations != "" {
		if err := json.Unmarshal([]byte(annotation.Annotations), &annotations); err != nil {
			return fmt.Errorf("解析现有标注失败: %w", err)
		}
	}

	// 查找并删除标注项 / Find and delete annotation item
	found := false
	for i, item := range annotations {
		if item.ID == itemID {
			annotations = append(annotations[:i], annotations[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("标注项不存在: %s", itemID)
	}

	// 更新数据库 / Update database
	return s.UpdateAnnotations(annotationID, annotations)
}

// GetAnnotationItem 获取单个标注项
// Get single annotation item
func (s *AnnotationService) GetAnnotationItem(annotationID uint, itemID string) (*model.AnnotationItem, error) {
	// 获取标注记录 / Get annotation record
	annotation, err := s.annotationRepo.GetByID(annotationID)
	if err != nil {
		return nil, fmt.Errorf("获取标注记录失败: %w", err)
	}

	// 解析标注数据 / Parse annotation data
	var annotations []model.AnnotationItem
	if annotation.Annotations != "" {
		if err := json.Unmarshal([]byte(annotation.Annotations), &annotations); err != nil {
			return nil, fmt.Errorf("解析标注数据失败: %w", err)
		}
	}

	// 查找标注项 / Find annotation item
	for _, item := range annotations {
		if item.ID == itemID {
			return &item, nil
		}
	}

	return nil, fmt.Errorf("标注项不存在: %s", itemID)
}

// GenerateDimensionAnnotation 生成尺寸标注
// Generate dimension annotation
func (s *AnnotationService) GenerateDimensionAnnotation(element model.DetectedElement) model.AnnotationItem {
	now := time.Now()

	// 提取尺寸信息 / Extract dimension information
	value := "未知尺寸"
	unit := "mm"
	if v, ok := element.Properties["value"].(string); ok {
		value = v
	}
	if u, ok := element.Properties["unit"].(string); ok {
		unit = u
	}

	return model.AnnotationItem{
		ID:       fmt.Sprintf("dim_%d", now.UnixNano()),
		Type:     "dimension",
		Position: s.calculateAnnotationPosition(element.BoundingBox),
		Content:  fmt.Sprintf("%s%s", value, unit),
		Style: model.AnnotationStyle{
			FontSize:   12,
			FontColor:  "#FF0000",
			LineColor:  "#FF0000",
			LineWidth:  1,
			Background: "#FFFFFF",
		},
		Properties: map[string]interface{}{
			"value":     value,
			"unit":      unit,
			"direction": element.Properties["direction"],
			"elementId": element.Properties["elementId"],
		},
		CreatedAt:   now,
		UpdatedAt:   now,
		IsGenerated: true,
	}
}

// GenerateMaterialAnnotation 生成材料标注
// Generate material annotation
func (s *AnnotationService) GenerateMaterialAnnotation(element model.DetectedElement) model.AnnotationItem {
	now := time.Now()

	// 提取材料信息 / Extract material information
	materialType := "未知材料"
	if mt, ok := element.Properties["materialType"].(string); ok {
		materialType = mt
	}

	return model.AnnotationItem{
		ID:       fmt.Sprintf("mat_%d", now.UnixNano()),
		Type:     "material",
		Position: s.calculateAnnotationPosition(element.BoundingBox),
		Content:  materialType,
		Style: model.AnnotationStyle{
			FontSize:   10,
			FontColor:  "#0000FF",
			LineColor:  "#0000FF",
			LineWidth:  1,
			Background: "#FFFFFF",
		},
		Properties: map[string]interface{}{
			"materialType": materialType,
			"area":         element.Properties["area"],
			"elementId":    element.Properties["elementId"],
		},
		CreatedAt:   now,
		UpdatedAt:   now,
		IsGenerated: true,
	}
}

// GenerateProcessAnnotation 生成工艺说明标注
// Generate process annotation
func (s *AnnotationService) GenerateProcessAnnotation(element model.DetectedElement) model.AnnotationItem {
	now := time.Now()

	// 提取设备信息 / Extract equipment information
	equipmentType := "设备"
	if et, ok := element.Properties["equipmentType"].(string); ok {
		equipmentType = et
	}

	// 生成工艺说明 / Generate process description
	processDescription := s.generateProcessDescription(equipmentType, element.Properties)

	return model.AnnotationItem{
		ID:       fmt.Sprintf("proc_%d", now.UnixNano()),
		Type:     "process",
		Position: s.calculateAnnotationPosition(element.BoundingBox),
		Content:  processDescription,
		Style: model.AnnotationStyle{
			FontSize:   10,
			FontColor:  "#00AA00",
			LineColor:  "#00AA00",
			LineWidth:  1,
			Background: "#FFFFFF",
		},
		Properties: map[string]interface{}{
			"equipmentType": equipmentType,
			"process":       processDescription,
			"elementId":     element.Properties["elementId"],
		},
		CreatedAt:   now,
		UpdatedAt:   now,
		IsGenerated: true,
	}
}

// generateProcessDescription 生成工艺说明
// Generate process description
func (s *AnnotationService) generateProcessDescription(equipmentType string, properties map[string]interface{}) string {
	switch equipmentType {
	case "插座", "socket":
		return "安装高度300mm，接地良好"
	case "开关", "switch":
		return "安装高度1300mm，便于操作"
	case "灯具", "light":
		return "吊装牢固，线路隐蔽"
	case "空调", "air_conditioner":
		return "预留排水管道，确保通风"
	case "消防设备", "fire_equipment":
		return "按消防规范安装，定期检查"
	default:
		return fmt.Sprintf("%s安装工艺说明", equipmentType)
	}
}

// ValidateAnnotationItem 验证标注项
// Validate annotation item
func (s *AnnotationService) ValidateAnnotationItem(item model.AnnotationItem) error {
	if item.Type == "" {
		return fmt.Errorf("标注类型不能为空")
	}

	if item.Content == "" {
		return fmt.Errorf("标注内容不能为空")
	}

	// 验证位置 / Validate position
	if item.Position.X < 0 || item.Position.Y < 0 {
		return fmt.Errorf("标注位置坐标不能为负数")
	}

	// 验证样式 / Validate style
	if item.Style.FontSize <= 0 {
		return fmt.Errorf("字体大小必须大于0")
	}

	return nil
}

// BatchUpdateAnnotations 批量更新标注
// Batch update annotations
func (s *AnnotationService) BatchUpdateAnnotations(annotationID uint, operations []AnnotationOperation) error {
	// 获取现有标注 / Get existing annotations
	annotation, err := s.annotationRepo.GetByID(annotationID)
	if err != nil {
		return fmt.Errorf("获取标注记录失败: %w", err)
	}

	// 解析现有标注 / Parse existing annotations
	var annotations []model.AnnotationItem
	if annotation.Annotations != "" {
		if err := json.Unmarshal([]byte(annotation.Annotations), &annotations); err != nil {
			return fmt.Errorf("解析现有标注失败: %w", err)
		}
	}

	// 执行批量操作 / Execute batch operations
	for _, op := range operations {
		switch op.Operation {
		case "add":
			if err := s.ValidateAnnotationItem(op.Item); err != nil {
				return fmt.Errorf("验证标注项失败: %w", err)
			}
			now := time.Now()
			op.Item.CreatedAt = now
			op.Item.UpdatedAt = now
			op.Item.IsGenerated = false
			if op.Item.ID == "" {
				op.Item.ID = fmt.Sprintf("batch_%d_%d", now.Unix(), len(annotations))
			}
			annotations = append(annotations, op.Item)

		case "update":
			found := false
			for i, item := range annotations {
				if item.ID == op.ItemID {
					if err := s.ValidateAnnotationItem(op.Item); err != nil {
						return fmt.Errorf("验证标注项失败: %w", err)
					}
					op.Item.ID = item.ID
					op.Item.CreatedAt = item.CreatedAt
					op.Item.UpdatedAt = time.Now()
					annotations[i] = op.Item
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("要更新的标注项不存在: %s", op.ItemID)
			}

		case "delete":
			found := false
			for i, item := range annotations {
				if item.ID == op.ItemID {
					annotations = append(annotations[:i], annotations[i+1:]...)
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("要删除的标注项不存在: %s", op.ItemID)
			}
		}
	}

	// 更新数据库 / Update database
	return s.UpdateAnnotations(annotationID, annotations)
}

// AnnotationOperation 标注操作
// Annotation operation
type AnnotationOperation struct {
	Operation string               `json:"operation"` // 操作类型: add, update, delete
	ItemID    string               `json:"itemId"`    // 标注项ID（update和delete时需要）
	Item      model.AnnotationItem `json:"item"`      // 标注项数据（add和update时需要）
}

// ExportAnnotation 导出标注
// Export annotation to specified format
func (s *AnnotationService) ExportAnnotation(annotationID uint, config model.AnnotationExportConfig) (*ExportResult, error) {
	// 获取标注记录 / Get annotation record
	annotation, err := s.annotationRepo.GetByID(annotationID)
	if err != nil {
		return nil, fmt.Errorf("获取标注记录失败: %w", err)
	}

	// 解析标注数据 / Parse annotation data
	var annotations []model.AnnotationItem
	if annotation.Annotations != "" {
		if err := json.Unmarshal([]byte(annotation.Annotations), &annotations); err != nil {
			return nil, fmt.Errorf("解析标注数据失败: %w", err)
		}
	}

	// 根据格式导出 / Export based on format
	switch strings.ToLower(config.Format) {
	case "pdf":
		return s.exportToPDF(annotation, annotations, config)
	case "png", "jpg", "jpeg":
		return s.exportToImage(annotation, annotations, config)
	default:
		return nil, fmt.Errorf("不支持的导出格式: %s", config.Format)
	}
}

// ExportResult 导出结果
// Export result
type ExportResult struct {
	FileName    string    `json:"fileName"`    // 文件名 / File name
	FilePath    string    `json:"filePath"`    // 文件路径 / File path
	FileSize    int64     `json:"fileSize"`    // 文件大小 / File size
	Format      string    `json:"format"`      // 导出格式 / Export format
	DownloadURL string    `json:"downloadUrl"` // 下载链接 / Download URL
	ExportedAt  time.Time `json:"exportedAt"`  // 导出时间 / Export time
}

// exportToPDF 导出为PDF格式
// Export to PDF format
func (s *AnnotationService) exportToPDF(annotation *model.ConstructionAnnotation, annotations []model.AnnotationItem, config model.AnnotationExportConfig) (*ExportResult, error) {
	// 生成文件名 / Generate file name
	timestamp := time.Now().Format("20060102_150405")
	fileName := fmt.Sprintf("annotation_%d_%s.pdf", annotation.ID, timestamp)
	filePath := filepath.Join("exports", "annotations", fileName)

	// 确保目录存在 / Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return nil, fmt.Errorf("创建导出目录失败: %w", err)
	}

	// 这里应该使用PDF生成库（如wkhtmltopdf或类似工具）
	// 为了演示，我们创建一个简单的文本文件
	// In production, use a PDF generation library like wkhtmltopdf
	content := s.generatePDFContent(annotation, annotations)

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return nil, fmt.Errorf("写入PDF文件失败: %w", err)
	}

	// 获取文件信息 / Get file info
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %w", err)
	}

	return &ExportResult{
		FileName:    fileName,
		FilePath:    filePath,
		FileSize:    fileInfo.Size(),
		Format:      "pdf",
		DownloadURL: fmt.Sprintf("/api/files/download/%s", fileName),
		ExportedAt:  time.Now(),
	}, nil
}

// exportToImage 导出为图片格式
// Export to image format
func (s *AnnotationService) exportToImage(annotation *model.ConstructionAnnotation, annotations []model.AnnotationItem, config model.AnnotationExportConfig) (*ExportResult, error) {
	// 生成文件名 / Generate file name
	timestamp := time.Now().Format("20060102_150405")
	fileName := fmt.Sprintf("annotation_%d_%s.%s", annotation.ID, timestamp, config.Format)
	filePath := filepath.Join("exports", "annotations", fileName)

	// 确保目录存在 / Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return nil, fmt.Errorf("创建导出目录失败: %w", err)
	}

	// 这里应该使用图像处理库来合成标注到原图上
	// 为了演示，我们复制原始图片
	// In production, use image processing library to overlay annotations
	if err := s.copyFile(annotation.ImagePath, filePath); err != nil {
		return nil, fmt.Errorf("复制图片文件失败: %w", err)
	}

	// 获取文件信息 / Get file info
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %w", err)
	}

	return &ExportResult{
		FileName:    fileName,
		FilePath:    filePath,
		FileSize:    fileInfo.Size(),
		Format:      config.Format,
		DownloadURL: fmt.Sprintf("/api/files/download/%s", fileName),
		ExportedAt:  time.Now(),
	}, nil
}

// generatePDFContent 生成PDF内容
// Generate PDF content
func (s *AnnotationService) generatePDFContent(annotation *model.ConstructionAnnotation, annotations []model.AnnotationItem) string {
	var content strings.Builder

	content.WriteString("施工图标注报告\n")
	content.WriteString("=================\n\n")

	content.WriteString(fmt.Sprintf("项目ID: %d\n", annotation.ProjectID))
	content.WriteString(fmt.Sprintf("CAD文件ID: %d\n", annotation.CadFileID))
	content.WriteString(fmt.Sprintf("图片路径: %s\n", annotation.ImagePath))
	content.WriteString(fmt.Sprintf("生成时间: %s\n\n", time.Now().Format("2006-01-02 15:04:05")))

	content.WriteString("标注列表:\n")
	content.WriteString("---------\n")

	for i, item := range annotations {
		content.WriteString(fmt.Sprintf("%d. %s\n", i+1, item.Content))
		content.WriteString(fmt.Sprintf("   类型: %s\n", item.Type))
		content.WriteString(fmt.Sprintf("   位置: (%.2f, %.2f)\n", item.Position.X, item.Position.Y))
		content.WriteString(fmt.Sprintf("   创建时间: %s\n", item.CreatedAt.Format("2006-01-02 15:04:05")))
		if item.IsGenerated {
			content.WriteString("   来源: 自动生成\n")
		} else {
			content.WriteString("   来源: 手动添加\n")
		}
		content.WriteString("\n")
	}

	return content.String()
}

// copyFile 复制文件
// Copy file
func (s *AnnotationService) copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// GetExportHistory 获取导出历史
// Get export history
func (s *AnnotationService) GetExportHistory(annotationID uint) ([]ExportRecord, error) {
	// 这里应该从数据库获取导出历史记录
	// 为了演示，返回空列表
	// In production, fetch from database
	return []ExportRecord{}, nil
}

// ExportRecord 导出记录
// Export record
type ExportRecord struct {
	ID           uint      `json:"id"`
	AnnotationID uint      `json:"annotationId"`
	FileName     string    `json:"fileName"`
	Format       string    `json:"format"`
	FileSize     int64     `json:"fileSize"`
	ExportedAt   time.Time `json:"exportedAt"`
	DownloadURL  string    `json:"downloadUrl"`
}

// ValidateExportConfig 验证导出配置
// Validate export configuration
func (s *AnnotationService) ValidateExportConfig(config model.AnnotationExportConfig) error {
	// 验证格式 / Validate format
	validFormats := []string{"pdf", "png", "jpg", "jpeg"}
	formatValid := false
	for _, format := range validFormats {
		if strings.ToLower(config.Format) == format {
			formatValid = true
			break
		}
	}
	if !formatValid {
		return fmt.Errorf("不支持的导出格式: %s，支持的格式: %v", config.Format, validFormats)
	}

	// 验证质量参数 / Validate quality parameter
	if config.Quality < 1 || config.Quality > 100 {
		return fmt.Errorf("图片质量必须在1-100之间")
	}

	// 验证缩放比例 / Validate scale ratio
	if config.Scale <= 0 || config.Scale > 10 {
		return fmt.Errorf("缩放比例必须在0-10之间")
	}

	return nil
}

// GetExportFormats 获取支持的导出格式
// Get supported export formats
func (s *AnnotationService) GetExportFormats() []ExportFormat {
	return []ExportFormat{
		{
			Format:      "pdf",
			Name:        "PDF文档",
			Description: "便携式文档格式，适合打印和分享",
			MimeType:    "application/pdf",
		},
		{
			Format:      "png",
			Name:        "PNG图片",
			Description: "无损压缩图片格式，支持透明背景",
			MimeType:    "image/png",
		},
		{
			Format:      "jpg",
			Name:        "JPEG图片",
			Description: "有损压缩图片格式，文件较小",
			MimeType:    "image/jpeg",
		},
	}
}

// ExportFormat 导出格式信息
// Export format information
type ExportFormat struct {
	Format      string `json:"format"`      // 格式代码 / Format code
	Name        string `json:"name"`        // 格式名称 / Format name
	Description string `json:"description"` // 格式描述 / Format description
	MimeType    string `json:"mimeType"`    // MIME类型 / MIME type
}
