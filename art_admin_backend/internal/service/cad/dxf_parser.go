package cad

import (
	"art_admin_backend/internal/dto/response"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// DxfParser DXF文件解析器
// DXF file parser for extracting layer and entity information
type DxfParser struct {
	content  string
	lines    []string
	position int
}

// NewDxfParser 创建DXF解析器
// Create new DXF parser
func NewDxfParser(content string) *DxfParser {
	// 按行分割内容 / Split content by lines
	lines := strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n")
	return &DxfParser{
		content:  content,
		lines:    lines,
		position: 0,
	}
}

// Parse 解析DXF文件
// Parse DXF file and return parse result
func (p *DxfParser) Parse() (*response.CadParseResult, error) {
	result := &response.CadParseResult{
		Layers:      []response.CadLayerInfo{},
		Entities:    []response.CadEntityInfo{},
		BoundingBox: &response.BoundingBox{},
	}

	// 初始化边界框 / Initialize bounding box
	result.BoundingBox.MinX = math.MaxFloat64
	result.BoundingBox.MinY = math.MaxFloat64
	result.BoundingBox.MinZ = math.MaxFloat64
	result.BoundingBox.MaxX = -math.MaxFloat64
	result.BoundingBox.MaxY = -math.MaxFloat64
	result.BoundingBox.MaxZ = -math.MaxFloat64

	// 解析各个段 / Parse sections
	p.parseLayers(result)
	p.parseEntities(result)

	// 计算统计信息 / Calculate statistics
	result.LayerCount = len(result.Layers)
	result.EntityCount = len(result.Entities)

	// 检测是否包含3D内容 / Detect 3D content
	result.Has3D = p.detect3DContent(result)

	// 如果没有找到有效的边界框，设置默认值 / Set default bounding box if not found
	if result.BoundingBox.MinX == math.MaxFloat64 {
		result.BoundingBox = &response.BoundingBox{
			MinX: 0, MinY: 0, MinZ: 0,
			MaxX: 1000, MaxY: 1000, MaxZ: 0,
		}
	}

	return result, nil
}

// parseLayers 解析图层信息
// Parse layer information from TABLES section
func (p *DxfParser) parseLayers(result *response.CadParseResult) {
	// 查找TABLES段 / Find TABLES section
	tablesStart := strings.Index(p.content, "SECTION\n  2\nTABLES")
	if tablesStart == -1 {
		tablesStart = strings.Index(p.content, "SECTION\r\n  2\r\nTABLES")
	}
	if tablesStart == -1 {
		// 尝试其他格式 / Try alternative format
		tablesStart = strings.Index(strings.ToUpper(p.content), "TABLES")
	}

	// 使用正则表达式查找图层定义 / Use regex to find layer definitions
	layerPattern := regexp.MustCompile(`(?s)LAYER\s+2\s+([^\s]+)`)
	matches := layerPattern.FindAllStringSubmatch(p.content, -1)

	layerMap := make(map[string]bool)
	for _, match := range matches {
		if len(match) > 1 {
			layerName := strings.TrimSpace(match[1])
			if layerName != "" && !layerMap[layerName] {
				layerMap[layerName] = true
				result.Layers = append(result.Layers, response.CadLayerInfo{
					Name:    layerName,
					Visible: true,
					Color:   7, // 默认白色 / Default white
				})
			}
		}
	}

	// 如果没有找到图层，尝试从实体中提取 / Extract layers from entities if not found
	if len(result.Layers) == 0 {
		p.extractLayersFromEntities(result)
	}

	// 确保至少有一个默认图层 / Ensure at least one default layer
	if len(result.Layers) == 0 {
		result.Layers = append(result.Layers, response.CadLayerInfo{
			Name:    "0",
			Visible: true,
			Color:   7,
		})
	}
}

// extractLayersFromEntities 从实体中提取图层
// Extract layer names from entities
func (p *DxfParser) extractLayersFromEntities(result *response.CadParseResult) {
	// 查找所有图层引用 / Find all layer references
	layerRefPattern := regexp.MustCompile(`\s+8\s*\n\s*([^\s\n]+)`)
	matches := layerRefPattern.FindAllStringSubmatch(p.content, -1)

	layerMap := make(map[string]bool)
	for _, match := range matches {
		if len(match) > 1 {
			layerName := strings.TrimSpace(match[1])
			if layerName != "" && !layerMap[layerName] {
				layerMap[layerName] = true
				result.Layers = append(result.Layers, response.CadLayerInfo{
					Name:    layerName,
					Visible: true,
					Color:   7,
				})
			}
		}
	}
}

// parseEntities 解析实体信息
// Parse entity information from ENTITIES section
func (p *DxfParser) parseEntities(result *response.CadParseResult) {
	// 查找ENTITIES段 / Find ENTITIES section
	entitiesStart := strings.Index(strings.ToUpper(p.content), "ENTITIES")
	if entitiesStart == -1 {
		return
	}

	// 支持的实体类型 / Supported entity types
	entityTypes := []string{"LINE", "CIRCLE", "ARC", "POLYLINE", "LWPOLYLINE", "POINT", "TEXT", "MTEXT", "INSERT", "DIMENSION", "SPLINE", "ELLIPSE", "3DFACE", "SOLID"}

	for _, entityType := range entityTypes {
		p.parseEntityType(entityType, result)
	}

	// 更新图层实体计数 / Update layer entity counts
	layerCounts := make(map[string]int)
	for _, entity := range result.Entities {
		layerCounts[entity.Layer]++
	}

	for i := range result.Layers {
		if count, ok := layerCounts[result.Layers[i].Name]; ok {
			result.Layers[i].EntityCount = count
		}
	}
}

// parseEntityType 解析特定类型的实体
// Parse entities of a specific type
func (p *DxfParser) parseEntityType(entityType string, result *response.CadParseResult) {
	// 构建实体匹配模式 / Build entity matching pattern
	pattern := regexp.MustCompile(`(?s)` + entityType + `\s*\n(.*?)(?:ENDSEC|LINE|CIRCLE|ARC|POLYLINE|LWPOLYLINE|POINT|TEXT|MTEXT|INSERT|DIMENSION|SPLINE|ELLIPSE|3DFACE|SOLID)`)
	matches := pattern.FindAllStringSubmatch(p.content, -1)

	for _, match := range matches {
		if len(match) > 1 {
			entity := p.parseEntityData(entityType, match[1], result)
			if entity != nil {
				result.Entities = append(result.Entities, *entity)
			}
		}
	}
}

// parseEntityData 解析实体数据
// Parse entity data and extract properties
func (p *DxfParser) parseEntityData(entityType string, data string, result *response.CadParseResult) *response.CadEntityInfo {
	entity := &response.CadEntityInfo{
		Type:       entityType,
		Layer:      "0",
		Color:      7,
		Properties: make(map[string]interface{}),
	}

	lines := strings.Split(data, "\n")
	for i := 0; i < len(lines)-1; i += 2 {
		code := strings.TrimSpace(lines[i])
		if i+1 >= len(lines) {
			break
		}
		value := strings.TrimSpace(lines[i+1])

		codeNum, err := strconv.Atoi(code)
		if err != nil {
			continue
		}

		switch codeNum {
		case 8: // 图层名 / Layer name
			entity.Layer = value
		case 62: // 颜色 / Color
			if color, err := strconv.Atoi(value); err == nil {
				entity.Color = color
			}
		case 6: // 线型 / Line type
			entity.LineType = value
		case 10, 11, 12, 13: // X坐标 / X coordinates
			if x, err := strconv.ParseFloat(value, 64); err == nil {
				key := "x" + strconv.Itoa(codeNum-10)
				entity.Properties[key] = x
				p.updateBoundingBox(result.BoundingBox, x, 0, 0, true, false, false)
			}
		case 20, 21, 22, 23: // Y坐标 / Y coordinates
			if y, err := strconv.ParseFloat(value, 64); err == nil {
				key := "y" + strconv.Itoa(codeNum-20)
				entity.Properties[key] = y
				p.updateBoundingBox(result.BoundingBox, 0, y, 0, false, true, false)
			}
		case 30, 31, 32, 33: // Z坐标 / Z coordinates
			if z, err := strconv.ParseFloat(value, 64); err == nil {
				key := "z" + strconv.Itoa(codeNum-30)
				entity.Properties[key] = z
				p.updateBoundingBox(result.BoundingBox, 0, 0, z, false, false, true)
			}
		case 40: // 半径或其他尺寸 / Radius or other dimension
			if r, err := strconv.ParseFloat(value, 64); err == nil {
				entity.Properties["radius"] = r
			}
		case 50: // 起始角度 / Start angle
			if angle, err := strconv.ParseFloat(value, 64); err == nil {
				entity.Properties["startAngle"] = angle
			}
		case 51: // 结束角度 / End angle
			if angle, err := strconv.ParseFloat(value, 64); err == nil {
				entity.Properties["endAngle"] = angle
			}
		case 1: // 文本内容 / Text content
			entity.Properties["text"] = value
		}
	}

	return entity
}

// updateBoundingBox 更新边界框
// Update bounding box with new coordinates
func (p *DxfParser) updateBoundingBox(box *response.BoundingBox, x, y, z float64, updateX, updateY, updateZ bool) {
	if updateX {
		if x < box.MinX {
			box.MinX = x
		}
		if x > box.MaxX {
			box.MaxX = x
		}
	}
	if updateY {
		if y < box.MinY {
			box.MinY = y
		}
		if y > box.MaxY {
			box.MaxY = y
		}
	}
	if updateZ {
		if z < box.MinZ {
			box.MinZ = z
		}
		if z > box.MaxZ {
			box.MaxZ = z
		}
	}
}

// detect3DContent 检测是否包含3D内容
// Detect if the DXF file contains 3D content
func (p *DxfParser) detect3DContent(result *response.CadParseResult) bool {
	// 检查是否有3D实体类型 / Check for 3D entity types
	for _, entity := range result.Entities {
		if entity.Type == "3DFACE" || entity.Type == "SOLID" {
			return true
		}
		// 检查Z坐标是否有非零值 / Check for non-zero Z coordinates
		for key, value := range entity.Properties {
			if strings.HasPrefix(key, "z") {
				if z, ok := value.(float64); ok && z != 0 {
					return true
				}
			}
		}
	}

	// 检查边界框Z范围 / Check bounding box Z range
	if result.BoundingBox != nil {
		zRange := result.BoundingBox.MaxZ - result.BoundingBox.MinZ
		if zRange > 0.001 { // 允许微小误差 / Allow small tolerance
			return true
		}
	}

	return false
}

// GetLayerNames 获取所有图层名称
// Get all layer names
func (p *DxfParser) GetLayerNames() []string {
	result, _ := p.Parse()
	names := make([]string, len(result.Layers))
	for i, layer := range result.Layers {
		names[i] = layer.Name
	}
	return names
}
