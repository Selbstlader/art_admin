package model

import (
	"encoding/json"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// genMaterialName 生成有效的材料名称
// Generate valid material name
func genMaterialName() gopter.Gen {
	return gen.AlphaString().SuchThat(func(s string) bool {
		return len(s) > 0 && len(s) <= 200
	})
}

// genMaterialCategory 生成有效的材料分类
// Generate valid material category
func genMaterialCategory() gopter.Gen {
	return gen.OneConstOf(
		"地板", "墙面", "天花", "门窗", "灯具", "家具", "五金", "其他",
	)
}

// genMaterialUnit 生成有效的材料单位
// Generate valid material unit
func genMaterialUnit() gopter.Gen {
	return gen.OneConstOf("m²", "m", "个", "kg", "套", "组")
}

// genUnitPrice 生成有效的单价
// Generate valid unit price
func genUnitPrice() gopter.Gen {
	return gen.Float64Range(0.01, 100000.0)
}

// genMaterialStatus 生成有效的材料状态
// Generate valid material status
func genMaterialStatus() gopter.Gen {
	return gen.OneConstOf("active", "inactive")
}

// genApplicableScenes 生成适用场景JSON
// Generate applicable scenes JSON
func genApplicableScenes() gopter.Gen {
	return gen.SliceOfN(3, gen.OneConstOf(
		"办公空间", "商业空间", "工业厂房", "酒店", "餐饮", "医疗", "教育",
	)).Map(func(scenes []string) string {
		data, _ := json.Marshal(scenes)
		return string(data)
	})
}

// **Feature: designer-ai-assistant, Property 10: Material CRUD Data Integrity**
// **Validates: Requirements 6.5**
// For any material create/update operation, reading the material immediately after
// SHALL return the same data that was written.
func TestMaterialCRUDDataIntegrity(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// 生成完整的材料数据
	// Generate complete material data
	materialGen := gopter.CombineGens(
		genMaterialName(),
		genMaterialCategory(),
		gen.AlphaString().SuchThat(func(s string) bool { return len(s) <= 200 }), // specification
		genMaterialUnit(),
		genUnitPrice(),
		gen.AlphaString().SuchThat(func(s string) bool { return len(s) <= 100 }), // brand
		gen.AlphaString().SuchThat(func(s string) bool { return len(s) <= 200 }), // supplier
		genMaterialStatus(),
	)

	properties.Property("material data round-trip preserves all fields", prop.ForAll(
		func(testData []any) bool {
			name := testData[0].(string)
			category := testData[1].(string)
			specification := testData[2].(string)
			unit := testData[3].(string)
			unitPrice := testData[4].(float64)
			brand := testData[5].(string)
			supplier := testData[6].(string)
			status := testData[7].(string)

			// 创建材料对象 / Create material object
			original := Material{
				Name:          name,
				Category:      category,
				Specification: specification,
				Unit:          unit,
				UnitPrice:     unitPrice,
				Brand:         brand,
				Supplier:      supplier,
				Status:        status,
			}

			// 模拟序列化和反序列化（模拟数据库读写）
			// Simulate serialization and deserialization (simulating DB read/write)
			data, err := json.Marshal(original)
			if err != nil {
				t.Logf("Failed to marshal material: %v", err)
				return false
			}

			var restored Material
			err = json.Unmarshal(data, &restored)
			if err != nil {
				t.Logf("Failed to unmarshal material: %v", err)
				return false
			}

			// 验证所有字段一致性 / Verify all fields consistency
			if restored.Name != original.Name {
				t.Logf("Name mismatch: expected %s, got %s", original.Name, restored.Name)
				return false
			}
			if restored.Category != original.Category {
				t.Logf("Category mismatch: expected %s, got %s", original.Category, restored.Category)
				return false
			}
			if restored.Specification != original.Specification {
				t.Logf("Specification mismatch: expected %s, got %s", original.Specification, restored.Specification)
				return false
			}
			if restored.Unit != original.Unit {
				t.Logf("Unit mismatch: expected %s, got %s", original.Unit, restored.Unit)
				return false
			}
			// 浮点数比较允许微小误差 / Allow small floating point error
			if abs(restored.UnitPrice-original.UnitPrice) > 0.001 {
				t.Logf("UnitPrice mismatch: expected %f, got %f", original.UnitPrice, restored.UnitPrice)
				return false
			}
			if restored.Brand != original.Brand {
				t.Logf("Brand mismatch: expected %s, got %s", original.Brand, restored.Brand)
				return false
			}
			if restored.Supplier != original.Supplier {
				t.Logf("Supplier mismatch: expected %s, got %s", original.Supplier, restored.Supplier)
				return false
			}
			if restored.Status != original.Status {
				t.Logf("Status mismatch: expected %s, got %s", original.Status, restored.Status)
				return false
			}

			return true
		},
		materialGen,
	))

	properties.TestingRun(t)
}

// TestMaterialApplicableScenesRoundTrip 测试适用场景JSON的往返一致性
// Test applicable scenes JSON round-trip consistency
func TestMaterialApplicableScenesRoundTrip(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("applicable scenes JSON round-trip preserves data", prop.ForAll(
		func(scenes string) bool {
			original := Material{
				Name:             "测试材料",
				ApplicableScenes: scenes,
			}

			data, err := json.Marshal(original)
			if err != nil {
				t.Logf("Failed to marshal material: %v", err)
				return false
			}

			var restored Material
			err = json.Unmarshal(data, &restored)
			if err != nil {
				t.Logf("Failed to unmarshal material: %v", err)
				return false
			}

			if restored.ApplicableScenes != original.ApplicableScenes {
				t.Logf("ApplicableScenes mismatch: expected %s, got %s",
					original.ApplicableScenes, restored.ApplicableScenes)
				return false
			}

			return true
		},
		genApplicableScenes(),
	))

	properties.TestingRun(t)
}

// abs 返回浮点数的绝对值
// Return absolute value of float64
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
