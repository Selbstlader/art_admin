package cost

import (
	"math"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// **Feature: designer-ai-assistant, Property 14: Cost Report Sum Consistency**
// **Validates: Requirements 10.3**
// For any cost estimate, the totalCost SHALL equal:
// materialCost + laborCost + equipmentCost + managementCost
func TestCostReportSumConsistency(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// 创建服务实例用于测试计算方法
	// Create service instance for testing calculation method
	service := &costService{}

	// 生成成本数据 / Generate cost data
	costGen := gopter.CombineGens(
		gen.Float64Range(0, 1000000), // materialCost
		gen.Float64Range(0, 1000000), // laborCost
		gen.Float64Range(0, 1000000), // equipmentCost
		gen.Float64Range(0, 1000000), // managementCost
	)

	properties.Property("total cost equals sum of all cost categories", prop.ForAll(
		func(data []interface{}) bool {
			materialCost := data[0].(float64)
			laborCost := data[1].(float64)
			equipmentCost := data[2].(float64)
			managementCost := data[3].(float64)

			// 使用服务方法计算总成本 / Use service method to calculate total cost
			totalCost := service.CalculateTotalCost(materialCost, laborCost, equipmentCost, managementCost)

			// 手动计算期望值 / Manually calculate expected value
			expectedTotal := materialCost + laborCost + equipmentCost + managementCost
			expectedTotal = math.Round(expectedTotal*100) / 100

			// 验证总成本等于各项之和 / Verify total equals sum of all categories
			diff := math.Abs(totalCost - expectedTotal)
			if diff > 0.01 {
				t.Logf("Sum mismatch: expected %.2f, got %.2f (diff: %.4f)", expectedTotal, totalCost, diff)
				t.Logf("Components: material=%.2f, labor=%.2f, equipment=%.2f, management=%.2f",
					materialCost, laborCost, equipmentCost, managementCost)
				return false
			}

			return true
		},
		costGen,
	))

	properties.TestingRun(t)
}

// TestCostReportSumConsistencyWithZeroValues 测试零值情况下的求和一致性
// Test sum consistency with zero values
func TestCostReportSumConsistencyWithZeroValues(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	service := &costService{}

	// 生成包含零值的成本数据 / Generate cost data with zero values
	costGen := gopter.CombineGens(
		gen.OneConstOf(0.0, 100.0, 1000.0, 10000.0), // materialCost (may be 0)
		gen.OneConstOf(0.0, 100.0, 1000.0, 10000.0), // laborCost (may be 0)
		gen.OneConstOf(0.0, 100.0, 1000.0, 10000.0), // equipmentCost (may be 0)
		gen.OneConstOf(0.0, 100.0, 1000.0, 10000.0), // managementCost (may be 0)
	)

	properties.Property("total cost equals sum even with zero values", prop.ForAll(
		func(data []interface{}) bool {
			materialCost := data[0].(float64)
			laborCost := data[1].(float64)
			equipmentCost := data[2].(float64)
			managementCost := data[3].(float64)

			totalCost := service.CalculateTotalCost(materialCost, laborCost, equipmentCost, managementCost)
			expectedTotal := math.Round((materialCost+laborCost+equipmentCost+managementCost)*100) / 100

			diff := math.Abs(totalCost - expectedTotal)
			if diff > 0.01 {
				t.Logf("Sum mismatch with zeros: expected %.2f, got %.2f", expectedTotal, totalCost)
				return false
			}

			return true
		},
		costGen,
	))

	properties.TestingRun(t)
}

// TestCostReportSumNonNegative 测试总成本非负
// Test total cost is non-negative
func TestCostReportSumNonNegative(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	service := &costService{}

	costGen := gopter.CombineGens(
		gen.Float64Range(0, 1000000),
		gen.Float64Range(0, 1000000),
		gen.Float64Range(0, 1000000),
		gen.Float64Range(0, 1000000),
	)

	properties.Property("total cost is always non-negative", prop.ForAll(
		func(data []interface{}) bool {
			materialCost := data[0].(float64)
			laborCost := data[1].(float64)
			equipmentCost := data[2].(float64)
			managementCost := data[3].(float64)

			totalCost := service.CalculateTotalCost(materialCost, laborCost, equipmentCost, managementCost)

			if totalCost < 0 {
				t.Logf("Negative total cost: %.2f", totalCost)
				return false
			}

			return true
		},
		costGen,
	))

	properties.TestingRun(t)
}

// TestCostReportSumCommutativity 测试求和交换律
// Test sum commutativity - order of addition doesn't matter
func TestCostReportSumCommutativity(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	service := &costService{}

	costGen := gopter.CombineGens(
		gen.Float64Range(0, 1000000),
		gen.Float64Range(0, 1000000),
		gen.Float64Range(0, 1000000),
		gen.Float64Range(0, 1000000),
	)

	properties.Property("sum is commutative", prop.ForAll(
		func(data []interface{}) bool {
			a := data[0].(float64)
			b := data[1].(float64)
			c := data[2].(float64)
			d := data[3].(float64)

			// 不同顺序计算 / Calculate in different orders
			total1 := service.CalculateTotalCost(a, b, c, d)
			total2 := service.CalculateTotalCost(d, c, b, a)
			total3 := service.CalculateTotalCost(b, d, a, c)

			// 验证结果相同 / Verify results are the same
			diff1 := math.Abs(total1 - total2)
			diff2 := math.Abs(total1 - total3)

			if diff1 > 0.01 || diff2 > 0.01 {
				t.Logf("Commutativity violation: total1=%.2f, total2=%.2f, total3=%.2f", total1, total2, total3)
				return false
			}

			return true
		},
		costGen,
	))

	properties.TestingRun(t)
}

// **Feature: designer-ai-assistant, Property 15: Budget Warning Trigger**
// **Validates: Requirements 10.4**
// For any cost estimate where totalCost > budgetLimit,
// the system SHALL return a budget exceeded warning.
func TestBudgetWarningTrigger(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	service := &costService{}

	// 生成总成本和预算上限 / Generate total cost and budget limit
	budgetGen := gopter.CombineGens(
		gen.Float64Range(1000, 1000000), // totalCost
		gen.Float64Range(1000, 1000000), // budgetLimit
	)

	properties.Property("budget warning triggers when totalCost exceeds budgetLimit", prop.ForAll(
		func(data []interface{}) bool {
			totalCost := data[0].(float64)
			budgetLimit := data[1].(float64)

			// 检查预算警告 / Check budget warning
			exceeded, exceededAmount, warning := service.CheckBudgetWarning(totalCost, budgetLimit)

			// 验证逻辑 / Verify logic
			if totalCost > budgetLimit {
				// 应该触发警告 / Should trigger warning
				if !exceeded {
					t.Logf("Warning not triggered: totalCost=%.2f > budgetLimit=%.2f", totalCost, budgetLimit)
					return false
				}
				// 超出金额应该正确 / Exceeded amount should be correct
				expectedExceeded := math.Round((totalCost-budgetLimit)*100) / 100
				diff := math.Abs(exceededAmount - expectedExceeded)
				if diff > 0.01 {
					t.Logf("Exceeded amount mismatch: expected %.2f, got %.2f", expectedExceeded, exceededAmount)
					return false
				}
				// 警告信息不应为空 / Warning message should not be empty
				if warning == "" {
					t.Logf("Warning message is empty when budget exceeded")
					return false
				}
			} else {
				// 不应该触发警告 / Should not trigger warning
				if exceeded {
					t.Logf("Warning triggered incorrectly: totalCost=%.2f <= budgetLimit=%.2f", totalCost, budgetLimit)
					return false
				}
				// 超出金额应该为0 / Exceeded amount should be 0
				if exceededAmount != 0 {
					t.Logf("Exceeded amount should be 0 when not exceeded, got %.2f", exceededAmount)
					return false
				}
			}

			return true
		},
		budgetGen,
	))

	properties.TestingRun(t)
}

// TestBudgetWarningWithZeroBudgetLimit 测试零预算上限情况
// Test budget warning with zero budget limit
func TestBudgetWarningWithZeroBudgetLimit(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	service := &costService{}

	// 生成总成本（预算上限为0或负数）/ Generate total cost (budget limit is 0 or negative)
	costGen := gen.Float64Range(0, 1000000)

	properties.Property("no warning when budget limit is zero or negative", prop.ForAll(
		func(totalCost float64) bool {
			// 测试预算上限为0 / Test budget limit is 0
			exceeded1, exceededAmount1, warning1 := service.CheckBudgetWarning(totalCost, 0)
			if exceeded1 || exceededAmount1 != 0 || warning1 != "" {
				t.Logf("Warning triggered with zero budget limit")
				return false
			}

			// 测试预算上限为负数 / Test budget limit is negative
			exceeded2, exceededAmount2, warning2 := service.CheckBudgetWarning(totalCost, -1000)
			if exceeded2 || exceededAmount2 != 0 || warning2 != "" {
				t.Logf("Warning triggered with negative budget limit")
				return false
			}

			return true
		},
		costGen,
	))

	properties.TestingRun(t)
}

// TestBudgetWarningExceededAmountNonNegative 测试超出金额非负
// Test exceeded amount is non-negative
func TestBudgetWarningExceededAmountNonNegative(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	service := &costService{}

	budgetGen := gopter.CombineGens(
		gen.Float64Range(0, 1000000),
		gen.Float64Range(0, 1000000),
	)

	properties.Property("exceeded amount is always non-negative", prop.ForAll(
		func(data []interface{}) bool {
			totalCost := data[0].(float64)
			budgetLimit := data[1].(float64)

			_, exceededAmount, _ := service.CheckBudgetWarning(totalCost, budgetLimit)

			if exceededAmount < 0 {
				t.Logf("Negative exceeded amount: %.2f", exceededAmount)
				return false
			}

			return true
		},
		budgetGen,
	))

	properties.TestingRun(t)
}

// TestBudgetWarningBoundaryCondition 测试边界条件
// Test boundary condition - totalCost equals budgetLimit
func TestBudgetWarningBoundaryCondition(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	service := &costService{}

	// 生成相同的总成本和预算上限 / Generate same total cost and budget limit
	budgetGen := gen.Float64Range(1000, 1000000)

	properties.Property("no warning when totalCost equals budgetLimit", prop.ForAll(
		func(amount float64) bool {
			// 当总成本等于预算上限时，不应触发警告
			// When total cost equals budget limit, should not trigger warning
			exceeded, exceededAmount, _ := service.CheckBudgetWarning(amount, amount)

			if exceeded {
				t.Logf("Warning triggered when totalCost equals budgetLimit: %.2f", amount)
				return false
			}

			if exceededAmount != 0 {
				t.Logf("Exceeded amount should be 0 when equal, got %.2f", exceededAmount)
				return false
			}

			return true
		},
		budgetGen,
	))

	properties.TestingRun(t)
}
