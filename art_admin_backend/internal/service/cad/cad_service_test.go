package cad

import (
	"art_admin_backend/internal/dto/response"
	"fmt"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

/**
 * Feature: designer-ai-assistant, Property 12: CAD Layer Data Completeness
 * Validates: Requirements 8.3
 *
 * For any successfully parsed CAD file, the layer list SHALL contain all layers
 * present in the original file.
 */
func TestProperty12_CadLayerDataCompleteness(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100

	properties := gopter.NewProperties(parameters)

	// 生成随机图层名称列表 / Generate random layer names
	// 使用带索引的图层名称确保唯一性 / Use indexed layer names to ensure uniqueness
	layerNameGen := gen.IntRange(1, 10).Map(func(count int) []string {
		layers := make([]string, count)
		for i := 0; i < count; i++ {
			layers[i] = fmt.Sprintf("Layer_%d", i+1)
		}
		return layers
	})

	properties.Property("parsed layers contain all original layers", prop.ForAll(
		func(layerNames []string) bool {
			// 构建包含指定图层的DXF内容 / Build DXF content with specified layers
			dxfContent := buildDxfWithLayers(layerNames)

			// 解析DXF内容 / Parse DXF content
			parser := NewDxfParser(dxfContent)
			result, err := parser.Parse()
			if err != nil {
				return false
			}

			// 验证所有原始图层都在解析结果中 / Verify all original layers are in parse result
			parsedLayerMap := make(map[string]bool)
			for _, layer := range result.Layers {
				parsedLayerMap[layer.Name] = true
			}

			for _, originalLayer := range layerNames {
				if !parsedLayerMap[originalLayer] {
					t.Logf("Missing layer: %s", originalLayer)
					return false
				}
			}

			return true
		},
		layerNameGen,
	))

	properties.TestingRun(t)
}

/**
 * Feature: designer-ai-assistant, Property 13: CAD 3D Detection Accuracy
 * Validates: Requirements 8.4
 *
 * For any CAD file, the has3D flag SHALL be true if and only if the file
 * contains 3D geometry data.
 */
func TestProperty13_Cad3DDetectionAccuracy(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100

	properties := gopter.NewProperties(parameters)

	// 测试2D内容 / Test 2D content
	properties.Property("2D content should have has3D=false", prop.ForAll(
		func(numEntities int) bool {
			// 构建只包含2D实体的DXF内容 / Build DXF with only 2D entities
			dxfContent := buildDxf2DContent(numEntities)

			parser := NewDxfParser(dxfContent)
			result, err := parser.Parse()
			if err != nil {
				return false
			}

			// 2D内容应该has3D=false / 2D content should have has3D=false
			return !result.Has3D
		},
		gen.IntRange(1, 50),
	))

	// 测试3D内容 / Test 3D content
	properties.Property("3D content should have has3D=true", prop.ForAll(
		func(numEntities int) bool {
			// 构建包含3D实体的DXF内容 / Build DXF with 3D entities
			dxfContent := buildDxf3DContent(numEntities)

			parser := NewDxfParser(dxfContent)
			result, err := parser.Parse()
			if err != nil {
				return false
			}

			// 3D内容应该has3D=true / 3D content should have has3D=true
			return result.Has3D
		},
		gen.IntRange(1, 50),
	))

	properties.TestingRun(t)
}

// buildDxfWithLayers 构建包含指定图层的DXF内容
// Build DXF content with specified layers
func buildDxfWithLayers(layerNames []string) string {
	content := "0\nSECTION\n2\nTABLES\n0\nTABLE\n2\nLAYER\n"

	for _, name := range layerNames {
		content += "0\nLAYER\n2\n" + name + "\n70\n0\n62\n7\n6\nCONTINUOUS\n"
	}

	content += "0\nENDTAB\n0\nENDSEC\n"

	// 添加ENTITIES段，每个图层至少有一个实体 / Add ENTITIES section with at least one entity per layer
	content += "0\nSECTION\n2\nENTITIES\n"
	for _, name := range layerNames {
		content += "0\nLINE\n8\n" + name + "\n10\n0.0\n20\n0.0\n30\n0.0\n11\n100.0\n21\n100.0\n31\n0.0\n"
	}
	content += "0\nENDSEC\n0\nEOF\n"

	return content
}

// buildDxf2DContent 构建只包含2D实体的DXF内容
// Build DXF content with only 2D entities (Z=0)
func buildDxf2DContent(numEntities int) string {
	content := "0\nSECTION\n2\nTABLES\n0\nTABLE\n2\nLAYER\n"
	content += "0\nLAYER\n2\n0\n70\n0\n62\n7\n6\nCONTINUOUS\n"
	content += "0\nENDTAB\n0\nENDSEC\n"

	content += "0\nSECTION\n2\nENTITIES\n"
	for i := 0; i < numEntities; i++ {
		// 所有Z坐标都是0 / All Z coordinates are 0
		content += "0\nLINE\n8\n0\n10\n0.0\n20\n0.0\n30\n0.0\n11\n100.0\n21\n100.0\n31\n0.0\n"
	}
	content += "0\nENDSEC\n0\nEOF\n"

	return content
}

// buildDxf3DContent 构建包含3D实体的DXF内容
// Build DXF content with 3D entities (non-zero Z coordinates)
func buildDxf3DContent(numEntities int) string {
	content := "0\nSECTION\n2\nTABLES\n0\nTABLE\n2\nLAYER\n"
	content += "0\nLAYER\n2\n0\n70\n0\n62\n7\n6\nCONTINUOUS\n"
	content += "0\nENDTAB\n0\nENDSEC\n"

	content += "0\nSECTION\n2\nENTITIES\n"
	for i := 0; i < numEntities; i++ {
		// 包含非零Z坐标 / Include non-zero Z coordinates
		content += "0\nLINE\n8\n0\n10\n0.0\n20\n0.0\n30\n0.0\n11\n100.0\n21\n100.0\n31\n50.0\n"
	}
	content += "0\nENDSEC\n0\nEOF\n"

	return content
}

// TestDxfParserBasic 基础解析测试
// Basic parsing test
func TestDxfParserBasic(t *testing.T) {
	// 测试基本的DXF解析 / Test basic DXF parsing
	dxfContent := buildDxfWithLayers([]string{"Layer1", "Layer2", "Layer3"})
	parser := NewDxfParser(dxfContent)
	result, err := parser.Parse()

	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if result == nil {
		t.Fatal("Parse result is nil")
	}

	if len(result.Layers) < 3 {
		t.Errorf("Expected at least 3 layers, got %d", len(result.Layers))
	}

	// 验证图层名称 / Verify layer names
	layerMap := make(map[string]bool)
	for _, layer := range result.Layers {
		layerMap[layer.Name] = true
	}

	expectedLayers := []string{"Layer1", "Layer2", "Layer3"}
	for _, expected := range expectedLayers {
		if !layerMap[expected] {
			t.Errorf("Expected layer %s not found", expected)
		}
	}
}

// TestCadParseResultStructure 测试解析结果结构完整性
// Test parse result structure completeness
func TestCadParseResultStructure(t *testing.T) {
	dxfContent := buildDxfWithLayers([]string{"TestLayer"})
	parser := NewDxfParser(dxfContent)
	result, err := parser.Parse()

	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// 验证结果结构 / Verify result structure
	if result.Layers == nil {
		t.Error("Layers should not be nil")
	}

	if result.Entities == nil {
		t.Error("Entities should not be nil")
	}

	if result.BoundingBox == nil {
		t.Error("BoundingBox should not be nil")
	}

	if result.LayerCount != len(result.Layers) {
		t.Errorf("LayerCount mismatch: expected %d, got %d", len(result.Layers), result.LayerCount)
	}

	if result.EntityCount != len(result.Entities) {
		t.Errorf("EntityCount mismatch: expected %d, got %d", len(result.Entities), result.EntityCount)
	}
}

// verifyLayerCompleteness 验证图层完整性辅助函数
// Helper function to verify layer completeness
func verifyLayerCompleteness(original []string, parsed []response.CadLayerInfo) bool {
	parsedMap := make(map[string]bool)
	for _, layer := range parsed {
		parsedMap[layer.Name] = true
	}

	for _, name := range original {
		if !parsedMap[name] {
			return false
		}
	}
	return true
}
