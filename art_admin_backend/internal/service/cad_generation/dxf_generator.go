package cad_generation

import (
	"fmt"
	"os"
	"strings"
)

// AIShapeParams AI返回的图形参数结构
// AI returned shape parameters structure
type AIShapeParams struct {
	Units  int           `json:"units"`  // 单位：4=毫米 / Unit: 4=mm
	Layers []AILayer     `json:"layers"` // 图层列表 / Layer list
	Shapes []AIShapeItem `json:"shapes"` // 图形列表 / Shape list
}

// AILayer 图层参数
// Layer parameters
type AILayer struct {
	Name  string `json:"name"`  // 图层名 / Layer name
	Color int    `json:"color"` // 颜色：1=红，2=黄，3=绿，4=青，5=蓝，6=品红，7=白 / Color
}

// AIShapeItem 单个图形参数
// Single shape parameters
type AIShapeItem struct {
	Type   string      `json:"type"`   // 图形类型 / Shape type
	Layer  string      `json:"layer"`  // 所属图层 / Layer name
	Params interface{} `json:"params"` // 图形参数 / Shape params
}

// RectangleParams 矩形参数
// Rectangle parameters
type RectangleParams struct {
	X      float64 `json:"x"`      // 左下角X / Bottom-left X
	Y      float64 `json:"y"`      // 左下角Y / Bottom-left Y
	Width  float64 `json:"width"`  // 宽度 / Width
	Height float64 `json:"height"` // 高度 / Height
}

// CircleParams 圆形参数
// Circle parameters
type CircleParams struct {
	CenterX float64 `json:"centerX"` // 圆心X / Center X
	CenterY float64 `json:"centerY"` // 圆心Y / Center Y
	Radius  float64 `json:"radius"`  // 半径 / Radius
}

// LineParams 直线参数
// Line parameters
type LineParams struct {
	X1 float64 `json:"x1"` // 起点X / Start X
	Y1 float64 `json:"y1"` // 起点Y / Start Y
	X2 float64 `json:"x2"` // 终点X / End X
	Y2 float64 `json:"y2"` // 终点Y / End Y
}

// TextParams 文字参数
// Text parameters
type TextParams struct {
	X      float64 `json:"x"`      // 位置X / Position X
	Y      float64 `json:"y"`      // 位置Y / Position Y
	Height float64 `json:"height"` // 文字高度 / Text height
	Text   string  `json:"text"`   // 文字内容 / Text content
}

// PolylineParams 多段线参数
// Polyline parameters
type PolylineParams struct {
	Points [][]float64 `json:"points"` // 顶点坐标 [[x1,y1],[x2,y2]...] / Vertices
	Closed bool        `json:"closed"` // 是否闭合 / Is closed
}

// ArcParams 圆弧参数
// Arc parameters
type ArcParams struct {
	CenterX    float64 `json:"centerX"`    // 圆心X / Center X
	CenterY    float64 `json:"centerY"`    // 圆心Y / Center Y
	Radius     float64 `json:"radius"`     // 半径 / Radius
	StartAngle float64 `json:"startAngle"` // 起始角度 / Start angle (degrees)
	EndAngle   float64 `json:"endAngle"`   // 结束角度 / End angle (degrees)
}

// DXFGenerator DXF文件生成器
// DXF file generator
type DXFGenerator struct {
	content strings.Builder
	handle  int // 实体句柄计数器 / Entity handle counter
}

// NewDXFGenerator 创建DXF生成器
// Create DXF generator
func NewDXFGenerator() *DXFGenerator {
	return &DXFGenerator{handle: 100}
}

// nextHandle 获取下一个句柄
// Get next handle
func (g *DXFGenerator) nextHandle() string {
	g.handle++
	return fmt.Sprintf("%X", g.handle)
}

// Generate 根据AI参数生成DXF文件
// Generate DXF file from AI parameters
func (g *DXFGenerator) Generate(params *AIShapeParams, outputPath string) error {
	// 写入HEADER段 / Write HEADER section
	g.writeHeader(params.Units)

	// 写入TABLES段 / Write TABLES section
	g.writeTables(params.Layers)

	// 写入BLOCKS段 / Write BLOCKS section
	g.writeBlocks()

	// 写入ENTITIES段 / Write ENTITIES section
	g.writeEntities(params.Shapes)

	// 写入文件结束 / Write EOF
	g.writeEOF()

	// 保存文件 / Save file
	return os.WriteFile(outputPath, []byte(g.content.String()), 0644)
}

// GenerateToString 生成DXF内容字符串
// Generate DXF content string
func (g *DXFGenerator) GenerateToString(params *AIShapeParams) string {
	g.writeHeader(params.Units)
	g.writeTables(params.Layers)
	g.writeBlocks()
	g.writeEntities(params.Shapes)
	g.writeEOF()
	return g.content.String()
}

// writeHeader 写入HEADER段
// Write HEADER section
func (g *DXFGenerator) writeHeader(units int) {
	g.content.WriteString("0\nSECTION\n")
	g.content.WriteString("2\nHEADER\n")

	// AutoCAD版本 / AutoCAD version
	g.content.WriteString("9\n$ACADVER\n1\nAC1009\n") // R12格式

	// 单位设置 / Unit setting
	g.content.WriteString(fmt.Sprintf("9\n$INSUNITS\n70\n%d\n", units))

	// 绘图范围 / Drawing limits
	g.content.WriteString("9\n$LIMMIN\n10\n0.0\n20\n0.0\n")
	g.content.WriteString("9\n$LIMMAX\n10\n100000.0\n20\n100000.0\n")

	g.content.WriteString("0\nENDSEC\n")
}

// writeTables 写入TABLES段
// Write TABLES section
func (g *DXFGenerator) writeTables(layers []AILayer) {
	g.content.WriteString("0\nSECTION\n")
	g.content.WriteString("2\nTABLES\n")

	// 图层表 / Layer table
	g.content.WriteString("0\nTABLE\n")
	g.content.WriteString("2\nLAYER\n")
	g.content.WriteString(fmt.Sprintf("70\n%d\n", len(layers)+1)) // 图层数量+默认图层

	// 默认图层0 / Default layer 0
	g.content.WriteString("0\nLAYER\n")
	g.content.WriteString("2\n0\n")
	g.content.WriteString("70\n0\n")
	g.content.WriteString("62\n7\n") // 白色
	g.content.WriteString("6\nCONTINUOUS\n")

	// 自定义图层 / Custom layers
	for _, layer := range layers {
		g.content.WriteString("0\nLAYER\n")
		g.content.WriteString(fmt.Sprintf("2\n%s\n", layer.Name))
		g.content.WriteString("70\n0\n")
		g.content.WriteString(fmt.Sprintf("62\n%d\n", layer.Color))
		g.content.WriteString("6\nCONTINUOUS\n")
	}

	g.content.WriteString("0\nENDTAB\n")
	g.content.WriteString("0\nENDSEC\n")
}

// writeBlocks 写入BLOCKS段
// Write BLOCKS section
func (g *DXFGenerator) writeBlocks() {
	g.content.WriteString("0\nSECTION\n")
	g.content.WriteString("2\nBLOCKS\n")
	g.content.WriteString("0\nENDSEC\n")
}

// writeEntities 写入ENTITIES段
// Write ENTITIES section
func (g *DXFGenerator) writeEntities(shapes []AIShapeItem) {
	g.content.WriteString("0\nSECTION\n")
	g.content.WriteString("2\nENTITIES\n")

	for _, shape := range shapes {
		switch shape.Type {
		case "line":
			g.writeLine(shape)
		case "rectangle":
			g.writeRectangle(shape)
		case "circle":
			g.writeCircle(shape)
		case "text":
			g.writeText(shape)
		case "polyline":
			g.writePolyline(shape)
		case "arc":
			g.writeArc(shape)
		}
	}

	g.content.WriteString("0\nENDSEC\n")
}

// writeLine 写入直线
// Write LINE entity
func (g *DXFGenerator) writeLine(shape AIShapeItem) {
	params, ok := shape.Params.(map[string]interface{})
	if !ok {
		return
	}

	x1 := getFloat(params, "x1")
	y1 := getFloat(params, "y1")
	x2 := getFloat(params, "x2")
	y2 := getFloat(params, "y2")

	g.content.WriteString("0\nLINE\n")
	g.content.WriteString(fmt.Sprintf("5\n%s\n", g.nextHandle()))
	g.content.WriteString(fmt.Sprintf("8\n%s\n", shape.Layer))
	g.content.WriteString(fmt.Sprintf("10\n%.6f\n", x1))
	g.content.WriteString(fmt.Sprintf("20\n%.6f\n", y1))
	g.content.WriteString("30\n0.0\n")
	g.content.WriteString(fmt.Sprintf("11\n%.6f\n", x2))
	g.content.WriteString(fmt.Sprintf("21\n%.6f\n", y2))
	g.content.WriteString("31\n0.0\n")
}

// writeRectangle 写入矩形（使用4条直线）
// Write rectangle using 4 lines
func (g *DXFGenerator) writeRectangle(shape AIShapeItem) {
	params, ok := shape.Params.(map[string]interface{})
	if !ok {
		return
	}

	x := getFloat(params, "x")
	y := getFloat(params, "y")
	w := getFloat(params, "width")
	h := getFloat(params, "height")

	// 四条边 / Four edges
	lines := []struct{ x1, y1, x2, y2 float64 }{
		{x, y, x + w, y},         // 底边
		{x + w, y, x + w, y + h}, // 右边
		{x + w, y + h, x, y + h}, // 顶边
		{x, y + h, x, y},         // 左边
	}

	for _, l := range lines {
		g.content.WriteString("0\nLINE\n")
		g.content.WriteString(fmt.Sprintf("5\n%s\n", g.nextHandle()))
		g.content.WriteString(fmt.Sprintf("8\n%s\n", shape.Layer))
		g.content.WriteString(fmt.Sprintf("10\n%.6f\n", l.x1))
		g.content.WriteString(fmt.Sprintf("20\n%.6f\n", l.y1))
		g.content.WriteString("30\n0.0\n")
		g.content.WriteString(fmt.Sprintf("11\n%.6f\n", l.x2))
		g.content.WriteString(fmt.Sprintf("21\n%.6f\n", l.y2))
		g.content.WriteString("31\n0.0\n")
	}
}

// writeCircle 写入圆
// Write CIRCLE entity
func (g *DXFGenerator) writeCircle(shape AIShapeItem) {
	params, ok := shape.Params.(map[string]interface{})
	if !ok {
		return
	}

	cx := getFloat(params, "centerX")
	cy := getFloat(params, "centerY")
	r := getFloat(params, "radius")

	g.content.WriteString("0\nCIRCLE\n")
	g.content.WriteString(fmt.Sprintf("5\n%s\n", g.nextHandle()))
	g.content.WriteString(fmt.Sprintf("8\n%s\n", shape.Layer))
	g.content.WriteString(fmt.Sprintf("10\n%.6f\n", cx))
	g.content.WriteString(fmt.Sprintf("20\n%.6f\n", cy))
	g.content.WriteString("30\n0.0\n")
	g.content.WriteString(fmt.Sprintf("40\n%.6f\n", r))
}

// writeText 写入文字
// Write TEXT entity
func (g *DXFGenerator) writeText(shape AIShapeItem) {
	params, ok := shape.Params.(map[string]interface{})
	if !ok {
		return
	}

	x := getFloat(params, "x")
	y := getFloat(params, "y")
	height := getFloat(params, "height")
	if height <= 0 {
		height = 200 // 默认文字高度200mm
	}
	text := getString(params, "text")

	g.content.WriteString("0\nTEXT\n")
	g.content.WriteString(fmt.Sprintf("5\n%s\n", g.nextHandle()))
	g.content.WriteString(fmt.Sprintf("8\n%s\n", shape.Layer))
	g.content.WriteString(fmt.Sprintf("10\n%.6f\n", x))
	g.content.WriteString(fmt.Sprintf("20\n%.6f\n", y))
	g.content.WriteString("30\n0.0\n")
	g.content.WriteString(fmt.Sprintf("40\n%.6f\n", height))
	g.content.WriteString(fmt.Sprintf("1\n%s\n", text))
}

// writePolyline 写入多段线
// Write POLYLINE entity (using LWPOLYLINE for simplicity)
func (g *DXFGenerator) writePolyline(shape AIShapeItem) {
	params, ok := shape.Params.(map[string]interface{})
	if !ok {
		return
	}

	pointsRaw, ok := params["points"].([]interface{})
	if !ok || len(pointsRaw) == 0 {
		return
	}

	closed := getBool(params, "closed")

	g.content.WriteString("0\nLWPOLYLINE\n")
	g.content.WriteString(fmt.Sprintf("5\n%s\n", g.nextHandle()))
	g.content.WriteString(fmt.Sprintf("8\n%s\n", shape.Layer))
	g.content.WriteString(fmt.Sprintf("90\n%d\n", len(pointsRaw))) // 顶点数
	if closed {
		g.content.WriteString("70\n1\n") // 闭合标志
	} else {
		g.content.WriteString("70\n0\n")
	}

	for _, p := range pointsRaw {
		point, ok := p.([]interface{})
		if !ok || len(point) < 2 {
			continue
		}
		x := toFloat(point[0])
		y := toFloat(point[1])
		g.content.WriteString(fmt.Sprintf("10\n%.6f\n", x))
		g.content.WriteString(fmt.Sprintf("20\n%.6f\n", y))
	}
}

// writeArc 写入圆弧
// Write ARC entity
func (g *DXFGenerator) writeArc(shape AIShapeItem) {
	params, ok := shape.Params.(map[string]interface{})
	if !ok {
		return
	}

	cx := getFloat(params, "centerX")
	cy := getFloat(params, "centerY")
	r := getFloat(params, "radius")
	startAngle := getFloat(params, "startAngle")
	endAngle := getFloat(params, "endAngle")

	g.content.WriteString("0\nARC\n")
	g.content.WriteString(fmt.Sprintf("5\n%s\n", g.nextHandle()))
	g.content.WriteString(fmt.Sprintf("8\n%s\n", shape.Layer))
	g.content.WriteString(fmt.Sprintf("10\n%.6f\n", cx))
	g.content.WriteString(fmt.Sprintf("20\n%.6f\n", cy))
	g.content.WriteString("30\n0.0\n")
	g.content.WriteString(fmt.Sprintf("40\n%.6f\n", r))
	g.content.WriteString(fmt.Sprintf("50\n%.6f\n", startAngle))
	g.content.WriteString(fmt.Sprintf("51\n%.6f\n", endAngle))
}

// writeEOF 写入文件结束
// Write EOF
func (g *DXFGenerator) writeEOF() {
	g.content.WriteString("0\nEOF\n")
}

// 辅助函数 / Helper functions
func getFloat(m map[string]interface{}, key string) float64 {
	if v, ok := m[key]; ok {
		return toFloat(v)
	}
	return 0
}

func toFloat(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	case int:
		return float64(val)
	case int64:
		return float64(val)
	default:
		return 0
	}
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func getBool(m map[string]interface{}, key string) bool {
	if v, ok := m[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}
