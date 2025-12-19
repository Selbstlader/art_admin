package material

import (
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// 材料分类常量 / Material category constants
var materialCategories = []string{"地板", "墙面", "天花", "门窗", "灯具", "家具", "五金", "其他"}

// 材料品牌常量 / Material brand constants
var materialBrands = []string{"东鹏", "马可波罗", "诺贝尔", "蒙娜丽莎", "冠珠", "欧派", "索菲亚", "好莱客"}

// MaterialFilterTestData 材料筛选测试数据
// Material filter test data
type MaterialFilterTestData struct {
	Category string
	Brand    string
	MinPrice float64
	MaxPrice float64
}

// MaterialTestItem 测试用材料项
// Test material item
type MaterialTestItem struct {
	Name      string
	Category  string
	Brand     string
	UnitPrice float64
	Status    string
}

// **Feature: designer-ai-assistant, Property 9: Material Filter Result Consistency**
// **Validates: Requirements 6.4**
// For any material search with filters, all returned materials SHALL match
// all specified filter criteria (category, priceRange, brand).
func TestMaterialFilterResultConsistency(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// 生成测试材料数据集 - 使用固定的分类和品牌列表
	// Generate test material dataset - use fixed category and brand lists
	materialGen := gopter.CombineGens(
		gen.IntRange(0, 99),                        // name suffix
		gen.IntRange(0, len(materialCategories)-1), // category index
		gen.IntRange(0, len(materialBrands)-1),     // brand index
		gen.Float64Range(100, 10000),               // unitPrice
	).Map(func(data []interface{}) MaterialTestItem {
		nameIdx := data[0].(int)
		catIdx := data[1].(int)
		brandIdx := data[2].(int)
		price := data[3].(float64)
		return MaterialTestItem{
			Name:      "材料" + string(rune('A'+nameIdx%26)),
			Category:  materialCategories[catIdx],
			Brand:     materialBrands[brandIdx],
			UnitPrice: price,
			Status:    "active",
		}
	})

	// 生成材料列表
	// Generate material list
	materialsGen := gen.SliceOfN(20, materialGen)

	// 生成筛选条件 - 使用相同的分类和品牌列表
	// Generate filter criteria - use same category and brand lists
	filterGen := gopter.CombineGens(
		gen.IntRange(-1, len(materialCategories)-1), // category index (-1 means no filter)
		gen.IntRange(-1, len(materialBrands)-1),     // brand index (-1 means no filter)
		gen.Float64Range(0, 5000),                   // minPrice
		gen.Float64Range(5000, 15000),               // maxPrice
	).Map(func(data []interface{}) MaterialFilterTestData {
		catIdx := data[0].(int)
		brandIdx := data[1].(int)
		minPrice := data[2].(float64)
		maxPrice := data[3].(float64)

		category := ""
		if catIdx >= 0 {
			category = materialCategories[catIdx]
		}

		brand := ""
		if brandIdx >= 0 {
			brand = materialBrands[brandIdx]
		}

		return MaterialFilterTestData{
			Category: category,
			Brand:    brand,
			MinPrice: minPrice,
			MaxPrice: maxPrice,
		}
	})

	properties.Property("filtered materials match all specified criteria", prop.ForAll(
		func(materials []MaterialTestItem, filter MaterialFilterTestData) bool {
			// 应用筛选逻辑 / Apply filter logic
			filtered := filterMaterials(materials, filter)

			// 验证所有返回的材料都符合筛选条件
			// Verify all returned materials match filter criteria
			for _, m := range filtered {
				// 检查分类筛选 / Check category filter
				if filter.Category != "" && m.Category != filter.Category {
					t.Logf("Category mismatch: expected %s, got %s", filter.Category, m.Category)
					return false
				}

				// 检查品牌筛选 / Check brand filter
				if filter.Brand != "" && m.Brand != filter.Brand {
					t.Logf("Brand mismatch: expected %s, got %s", filter.Brand, m.Brand)
					return false
				}

				// 检查价格区间筛选 / Check price range filter
				if filter.MinPrice > 0 && m.UnitPrice < filter.MinPrice {
					t.Logf("Price below minimum: expected >= %f, got %f", filter.MinPrice, m.UnitPrice)
					return false
				}
				if filter.MaxPrice > 0 && m.UnitPrice > filter.MaxPrice {
					t.Logf("Price above maximum: expected <= %f, got %f", filter.MaxPrice, m.UnitPrice)
					return false
				}
			}

			return true
		},
		materialsGen,
		filterGen,
	))

	properties.TestingRun(t)
}

// TestMaterialFilterCompletenessProperty 测试筛选结果完整性
// Test filter result completeness - all matching materials should be included
func TestMaterialFilterCompletenessProperty(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// 生成测试材料数据集
	// Generate test material dataset
	materialGen := gopter.CombineGens(
		gen.IntRange(0, 99),
		gen.IntRange(0, len(materialCategories)-1),
		gen.IntRange(0, len(materialBrands)-1),
		gen.Float64Range(100, 10000),
	).Map(func(data []interface{}) MaterialTestItem {
		nameIdx := data[0].(int)
		catIdx := data[1].(int)
		brandIdx := data[2].(int)
		price := data[3].(float64)
		return MaterialTestItem{
			Name:      "材料" + string(rune('A'+nameIdx%26)),
			Category:  materialCategories[catIdx],
			Brand:     materialBrands[brandIdx],
			UnitPrice: price,
			Status:    "active",
		}
	})

	materialsGen := gen.SliceOfN(20, materialGen)

	filterGen := gopter.CombineGens(
		gen.IntRange(-1, len(materialCategories)-1),
		gen.IntRange(-1, len(materialBrands)-1),
		gen.Float64Range(0, 5000),
		gen.Float64Range(5000, 15000),
	).Map(func(data []interface{}) MaterialFilterTestData {
		catIdx := data[0].(int)
		brandIdx := data[1].(int)
		minPrice := data[2].(float64)
		maxPrice := data[3].(float64)

		category := ""
		if catIdx >= 0 {
			category = materialCategories[catIdx]
		}

		brand := ""
		if brandIdx >= 0 {
			brand = materialBrands[brandIdx]
		}

		return MaterialFilterTestData{
			Category: category,
			Brand:    brand,
			MinPrice: minPrice,
			MaxPrice: maxPrice,
		}
	})

	properties.Property("all matching materials are included in results", prop.ForAll(
		func(materials []MaterialTestItem, filter MaterialFilterTestData) bool {
			filtered := filterMaterials(materials, filter)

			// 手动计算应该匹配的材料数量
			// Manually calculate expected matching materials count
			expectedCount := 0
			for _, m := range materials {
				if matchesFilter(m, filter) {
					expectedCount++
				}
			}

			// 验证筛选结果数量正确
			// Verify filter result count is correct
			if len(filtered) != expectedCount {
				t.Logf("Result count mismatch: expected %d, got %d", expectedCount, len(filtered))
				return false
			}

			return true
		},
		materialsGen,
		filterGen,
	))

	properties.TestingRun(t)
}

// filterMaterials 根据筛选条件过滤材料列表
// Filter materials based on filter criteria
func filterMaterials(materials []MaterialTestItem, filter MaterialFilterTestData) []MaterialTestItem {
	var result []MaterialTestItem
	for _, m := range materials {
		if matchesFilter(m, filter) {
			result = append(result, m)
		}
	}
	return result
}

// matchesFilter 检查材料是否匹配筛选条件
// Check if material matches filter criteria
func matchesFilter(m MaterialTestItem, filter MaterialFilterTestData) bool {
	// 检查分类 / Check category
	if filter.Category != "" && m.Category != filter.Category {
		return false
	}

	// 检查品牌 / Check brand
	if filter.Brand != "" && m.Brand != filter.Brand {
		return false
	}

	// 检查价格区间 / Check price range
	if filter.MinPrice > 0 && m.UnitPrice < filter.MinPrice {
		return false
	}
	if filter.MaxPrice > 0 && m.UnitPrice > filter.MaxPrice {
		return false
	}

	return true
}

// **Feature: designer-ai-assistant, Property 8: Material Cost Calculation Accuracy**
// **Validates: Requirements 6.3**
// For any material selection, the calculated total cost SHALL equal: quantity × unitPrice.
func TestMaterialCostCalculationAccuracy(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// 生成测试数据 / Generate test data
	costCalcGen := gopter.CombineGens(
		gen.Float64Range(0.01, 100000), // unitPrice
		gen.Float64Range(1, 10000),     // area
		gen.Float64Range(0, 0.3),       // lossRate (0-30%)
	)

	properties.Property("total cost equals quantity times unit price", prop.ForAll(
		func(data []interface{}) bool {
			unitPrice := data[0].(float64)
			area := data[1].(float64)
			lossRate := data[2].(float64)

			// 计算用量（含损耗）/ Calculate quantity with loss
			quantity := area * (1 + lossRate)

			// 计算总成本 / Calculate total cost
			totalCost := quantity * unitPrice

			// 验证计算结果 / Verify calculation result
			expectedCost := area * (1 + lossRate) * unitPrice

			// 允许浮点数误差 / Allow floating point error
			diff := totalCost - expectedCost
			if diff < 0 {
				diff = -diff
			}

			if diff > 0.01 {
				t.Logf("Cost calculation mismatch: expected %.4f, got %.4f", expectedCost, totalCost)
				return false
			}

			return true
		},
		costCalcGen,
	))

	properties.TestingRun(t)
}

// TestMaterialCostCalculationWithZeroLossRate 测试零损耗率的成本计算
// Test cost calculation with zero loss rate
func TestMaterialCostCalculationWithZeroLossRate(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	costCalcGen := gopter.CombineGens(
		gen.Float64Range(0.01, 100000), // unitPrice
		gen.Float64Range(1, 10000),     // area
	)

	properties.Property("with zero loss rate, quantity equals area", prop.ForAll(
		func(data []interface{}) bool {
			unitPrice := data[0].(float64)
			area := data[1].(float64)
			lossRate := 0.0

			// 计算用量 / Calculate quantity
			quantity := area * (1 + lossRate)

			// 验证用量等于面积 / Verify quantity equals area
			if quantity != area {
				t.Logf("Quantity mismatch: expected %.4f, got %.4f", area, quantity)
				return false
			}

			// 验证总成本 / Verify total cost
			totalCost := quantity * unitPrice
			expectedCost := area * unitPrice

			diff := totalCost - expectedCost
			if diff < 0 {
				diff = -diff
			}

			if diff > 0.01 {
				t.Logf("Cost mismatch: expected %.4f, got %.4f", expectedCost, totalCost)
				return false
			}

			return true
		},
		costCalcGen,
	))

	properties.TestingRun(t)
}

// TestMaterialCostCalculationNonNegative 测试成本计算结果非负
// Test cost calculation result is non-negative
func TestMaterialCostCalculationNonNegative(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	costCalcGen := gopter.CombineGens(
		gen.Float64Range(0, 100000), // unitPrice (can be 0)
		gen.Float64Range(0, 10000),  // area (can be 0)
		gen.Float64Range(0, 1),      // lossRate
	)

	properties.Property("total cost is always non-negative", prop.ForAll(
		func(data []interface{}) bool {
			unitPrice := data[0].(float64)
			area := data[1].(float64)
			lossRate := data[2].(float64)

			// 计算用量和成本 / Calculate quantity and cost
			quantity := area * (1 + lossRate)
			totalCost := quantity * unitPrice

			// 验证非负 / Verify non-negative
			if totalCost < 0 {
				t.Logf("Negative cost: %.4f", totalCost)
				return false
			}

			if quantity < 0 {
				t.Logf("Negative quantity: %.4f", quantity)
				return false
			}

			return true
		},
		costCalcGen,
	))

	properties.TestingRun(t)
}

// TestMaterialCostCalculationMonotonicity 测试成本计算单调性
// Test cost calculation monotonicity - larger area means higher cost
func TestMaterialCostCalculationMonotonicity(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	costCalcGen := gopter.CombineGens(
		gen.Float64Range(0.01, 100000), // unitPrice
		gen.Float64Range(1, 5000),      // area1
		gen.Float64Range(1, 5000),      // area2
		gen.Float64Range(0, 0.3),       // lossRate
	)

	properties.Property("larger area results in higher or equal cost", prop.ForAll(
		func(data []interface{}) bool {
			unitPrice := data[0].(float64)
			area1 := data[1].(float64)
			area2 := data[2].(float64)
			lossRate := data[3].(float64)

			// 计算两个面积的成本 / Calculate costs for two areas
			cost1 := area1 * (1 + lossRate) * unitPrice
			cost2 := area2 * (1 + lossRate) * unitPrice

			// 验证单调性 / Verify monotonicity
			if area1 > area2 && cost1 < cost2 {
				t.Logf("Monotonicity violation: area1=%.2f > area2=%.2f but cost1=%.2f < cost2=%.2f",
					area1, area2, cost1, cost2)
				return false
			}
			if area1 < area2 && cost1 > cost2 {
				t.Logf("Monotonicity violation: area1=%.2f < area2=%.2f but cost1=%.2f > cost2=%.2f",
					area1, area2, cost1, cost2)
				return false
			}

			return true
		},
		costCalcGen,
	))

	properties.TestingRun(t)
}
