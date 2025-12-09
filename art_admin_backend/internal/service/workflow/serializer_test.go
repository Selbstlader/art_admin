package workflow

import (
	"reflect"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"

	"art_admin_backend/internal/model"
)

// ============================================================================
// Generators for Property-Based Testing
// ============================================================================

// genNodeType generates random NodeType values
func genNodeType() gopter.Gen {
	return gen.OneConstOf(
		model.NodeTypeStart,
		model.NodeTypeEnd,
		model.NodeTypeApproval,
		model.NodeTypeCondition,
		model.NodeTypeParallel,
	)
}

// genApprovalMode generates random ApprovalMode values
func genApprovalMode() gopter.Gen {
	return gen.OneConstOf(
		model.ApprovalModeOrSign,
		model.ApprovalModeAndSign,
	)
}

// genRejectAction generates random RejectAction values
func genRejectAction() gopter.Gen {
	return gen.OneConstOf(
		model.RejectActionTerminate,
		model.RejectActionReturnPrev,
	)
}

// genAssigneeType generates random AssigneeType values
func genAssigneeType() gopter.Gen {
	return gen.OneConstOf(
		model.AssigneeTypeUser,
		model.AssigneeTypeRole,
		model.AssigneeTypeDeptLeader,
		model.AssigneeTypeInitiatorLeader,
	)
}

// genPosition generates random Position values
func genPosition() gopter.Gen {
	return gopter.CombineGens(
		gen.Float64Range(-1000, 1000),
		gen.Float64Range(-1000, 1000),
	).Map(func(vals []interface{}) model.Position {
		return model.Position{
			X: vals[0].(float64),
			Y: vals[1].(float64),
		}
	})
}

// genAssigneeRule generates random AssigneeRule values
func genAssigneeRule() gopter.Gen {
	return gopter.CombineGens(
		genAssigneeType(),
		gen.SliceOf(gen.Int64Range(1, 1000)),
	).Map(func(vals []interface{}) *model.AssigneeRule {
		return &model.AssigneeRule{
			Type:   vals[0].(model.AssigneeType),
			Values: vals[1].([]int64),
		}
	})
}

// genNodeProperties generates random NodeProperties values
func genNodeProperties() gopter.Gen {
	return gopter.CombineGens(
		gen.PtrOf(genAssigneeRule()),
		genApprovalMode(),
		gen.IntRange(0, 168), // 0-168 hours (1 week)
		genRejectAction(),
	).Map(func(vals []interface{}) *model.NodeProperties {
		props := &model.NodeProperties{
			ApprovalMode: vals[1].(model.ApprovalMode),
			TimeoutHours: vals[2].(int),
			RejectAction: vals[3].(model.RejectAction),
		}
		if rule, ok := vals[0].(*model.AssigneeRule); ok && rule != nil {
			props.AssigneeRule = rule
		}
		// Generate field permissions
		props.FieldPermissions = map[string]string{
			"field1": "visible",
			"field2": "editable",
			"field3": "hidden",
		}
		return props
	})
}

// genProcessNode generates random ProcessNode values
func genProcessNode() gopter.Gen {
	return gopter.CombineGens(
		gen.AlphaString().SuchThat(func(s string) bool { return len(s) > 0 && len(s) <= 20 }),
		genNodeType(),
		gen.AlphaString().SuchThat(func(s string) bool { return len(s) > 0 && len(s) <= 50 }),
		genPosition(),
		gen.PtrOf(genNodeProperties()),
	).Map(func(vals []interface{}) model.ProcessNode {
		node := model.ProcessNode{
			ID:       vals[0].(string),
			Type:     vals[1].(model.NodeType),
			Name:     vals[2].(string),
			Position: vals[3].(model.Position),
		}
		if props, ok := vals[4].(*model.NodeProperties); ok {
			node.Properties = props
		}
		return node
	})
}

// genProcessEdge generates random ProcessEdge values
func genProcessEdge() gopter.Gen {
	return gopter.CombineGens(
		gen.AlphaString().SuchThat(func(s string) bool { return len(s) > 0 && len(s) <= 20 }),
		gen.AlphaString().SuchThat(func(s string) bool { return len(s) > 0 && len(s) <= 20 }),
		gen.AlphaString().SuchThat(func(s string) bool { return len(s) > 0 && len(s) <= 20 }),
		gen.PtrOf(gen.AlphaString()),
	).Map(func(vals []interface{}) model.ProcessEdge {
		edge := model.ProcessEdge{
			ID:     vals[0].(string),
			Source: vals[1].(string),
			Target: vals[2].(string),
		}
		if cond, ok := vals[3].(*string); ok {
			edge.Condition = cond
		}
		return edge
	})
}

// genGlobalProperties generates random GlobalProperties values
func genGlobalProperties() gopter.Gen {
	return gen.Bool().Map(func(b bool) *model.GlobalProperties {
		return &model.GlobalProperties{
			AllowWithdraw: b,
		}
	})
}

// genProcessGraph generates random ProcessGraph values
func genProcessGraph() gopter.Gen {
	return gopter.CombineGens(
		gen.SliceOf(genProcessNode()).SuchThat(func(nodes []model.ProcessNode) bool {
			return len(nodes) >= 1 && len(nodes) <= 10
		}),
		gen.SliceOf(genProcessEdge()).SuchThat(func(edges []model.ProcessEdge) bool {
			return len(edges) <= 15
		}),
		gen.PtrOf(genGlobalProperties()),
	).Map(func(vals []interface{}) *model.ProcessGraph {
		graph := &model.ProcessGraph{
			Nodes: vals[0].([]model.ProcessNode),
			Edges: vals[1].([]model.ProcessEdge),
		}
		if props, ok := vals[2].(*model.GlobalProperties); ok {
			graph.GlobalProps = props
		}
		return graph
	})
}

// ============================================================================
// Property-Based Tests
// ============================================================================

// **Feature: oa-workflow-engine, Property 1: 流程定义序列化往返一致性**
// **Validates: Requirements 1.6, 1.7**
// *For any* 有效的ProcessGraph对象，序列化为JSON后再反序列化，应产生与原始对象等价的ProcessGraph。
func TestProcessGraphSerializationRoundTrip(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	parameters.Rng.Seed(42) // For reproducibility

	properties := gopter.NewProperties(parameters)

	serializer := NewSerializer()

	properties.Property("ProcessGraph serialization round-trip preserves data", prop.ForAll(
		func(graph *model.ProcessGraph) bool {
			// Serialize
			jsonStr, err := serializer.SerializeProcessGraph(graph)
			if err != nil {
				t.Logf("Serialization error: %v", err)
				return false
			}

			// Deserialize
			restored, err := serializer.DeserializeProcessGraph(jsonStr)
			if err != nil {
				t.Logf("Deserialization error: %v", err)
				return false
			}

			// Compare
			return processGraphsEqual(graph, restored)
		},
		genProcessGraph(),
	))

	properties.TestingRun(t)
}

// processGraphsEqual compares two ProcessGraph objects for equality
func processGraphsEqual(a, b *model.ProcessGraph) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Compare nodes
	if len(a.Nodes) != len(b.Nodes) {
		return false
	}
	for i := range a.Nodes {
		if !processNodesEqual(&a.Nodes[i], &b.Nodes[i]) {
			return false
		}
	}

	// Compare edges
	if len(a.Edges) != len(b.Edges) {
		return false
	}
	for i := range a.Edges {
		if !processEdgesEqual(&a.Edges[i], &b.Edges[i]) {
			return false
		}
	}

	// Compare global properties
	if !globalPropertiesEqual(a.GlobalProps, b.GlobalProps) {
		return false
	}

	return true
}

// processNodesEqual compares two ProcessNode objects for equality
func processNodesEqual(a, b *model.ProcessNode) bool {
	if a.ID != b.ID || a.Type != b.Type || a.Name != b.Name {
		return false
	}
	if a.Position.X != b.Position.X || a.Position.Y != b.Position.Y {
		return false
	}
	return nodePropertiesEqual(a.Properties, b.Properties)
}

// nodePropertiesEqual compares two NodeProperties objects for equality
func nodePropertiesEqual(a, b *model.NodeProperties) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	if a.ApprovalMode != b.ApprovalMode ||
		a.TimeoutHours != b.TimeoutHours ||
		a.RejectAction != b.RejectAction {
		return false
	}

	// Compare AssigneeRule
	if !assigneeRulesEqual(a.AssigneeRule, b.AssigneeRule) {
		return false
	}

	// Compare FieldPermissions
	if !reflect.DeepEqual(a.FieldPermissions, b.FieldPermissions) {
		return false
	}

	return true
}

// assigneeRulesEqual compares two AssigneeRule objects for equality
func assigneeRulesEqual(a, b *model.AssigneeRule) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if a.Type != b.Type {
		return false
	}
	if len(a.Values) != len(b.Values) {
		return false
	}
	for i := range a.Values {
		if a.Values[i] != b.Values[i] {
			return false
		}
	}
	return true
}

// processEdgesEqual compares two ProcessEdge objects for equality
func processEdgesEqual(a, b *model.ProcessEdge) bool {
	if a.ID != b.ID || a.Source != b.Source || a.Target != b.Target {
		return false
	}
	// Compare Condition pointers
	if a.Condition == nil && b.Condition == nil {
		return true
	}
	if a.Condition == nil || b.Condition == nil {
		return false
	}
	return *a.Condition == *b.Condition
}

// globalPropertiesEqual compares two GlobalProperties objects for equality
func globalPropertiesEqual(a, b *model.GlobalProperties) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.AllowWithdraw == b.AllowWithdraw
}

// ============================================================================
// FormSchema Generators
// ============================================================================

// genFieldType generates random FieldType values
func genFieldType() gopter.Gen {
	return gen.OneConstOf(
		model.FieldTypeText,
		model.FieldTypeNumber,
		model.FieldTypeDate,
		model.FieldTypeSelect,
		model.FieldTypeFile,
		model.FieldTypeTextarea,
	)
}

// genSelectOption generates random SelectOption values
func genSelectOption() gopter.Gen {
	return gopter.CombineGens(
		gen.AlphaString().SuchThat(func(s string) bool { return len(s) > 0 && len(s) <= 20 }),
		gen.AlphaString(), // Value can be any string
	).Map(func(vals []interface{}) model.SelectOption {
		return model.SelectOption{
			Label: vals[0].(string),
			Value: vals[1].(string),
		}
	})
}

// genFieldValidation generates random FieldValidation values
func genFieldValidation() gopter.Gen {
	return gopter.CombineGens(
		gen.PtrOf(gen.IntRange(0, 100)),
		gen.PtrOf(gen.IntRange(1, 500)),
		gen.PtrOf(gen.Float64Range(-1000, 1000)),
		gen.PtrOf(gen.Float64Range(-1000, 1000)),
		gen.PtrOf(gen.AlphaString()),
	).Map(func(vals []interface{}) *model.FieldValidation {
		validation := &model.FieldValidation{}
		if minLen, ok := vals[0].(*int); ok {
			validation.MinLength = minLen
		}
		if maxLen, ok := vals[1].(*int); ok {
			validation.MaxLength = maxLen
		}
		if min, ok := vals[2].(*float64); ok {
			validation.Min = min
		}
		if max, ok := vals[3].(*float64); ok {
			validation.Max = max
		}
		if pattern, ok := vals[4].(*string); ok {
			validation.Pattern = pattern
		}
		return validation
	})
}

// genFormField generates random FormField values
func genFormField() gopter.Gen {
	return gopter.CombineGens(
		gen.AlphaString().SuchThat(func(s string) bool { return len(s) > 0 && len(s) <= 30 }),
		gen.AlphaString().SuchThat(func(s string) bool { return len(s) > 0 && len(s) <= 50 }),
		genFieldType(),
		gen.Bool(),
		gen.AlphaString(),
		gen.SliceOf(genSelectOption()).SuchThat(func(opts []model.SelectOption) bool {
			return len(opts) <= 10
		}),
		gen.PtrOf(genFieldValidation()),
	).Map(func(vals []interface{}) model.FormField {
		field := model.FormField{
			Key:         vals[0].(string),
			Label:       vals[1].(string),
			Type:        vals[2].(model.FieldType),
			Required:    vals[3].(bool),
			Placeholder: vals[4].(string),
			Options:     vals[5].([]model.SelectOption),
		}
		if validation, ok := vals[6].(*model.FieldValidation); ok {
			field.Validation = validation
		}
		// DefaultValue is interface{}, use a simple string for testing
		field.DefaultValue = "default"
		return field
	})
}

// genFormSchema generates random FormSchema values
func genFormSchema() gopter.Gen {
	return gen.SliceOf(genFormField()).SuchThat(func(fields []model.FormField) bool {
		return len(fields) >= 1 && len(fields) <= 20
	}).Map(func(fields []model.FormField) *model.FormSchema {
		return &model.FormSchema{
			Fields: fields,
		}
	})
}

// ============================================================================
// FormSchema Property-Based Tests
// ============================================================================

// **Feature: oa-workflow-engine, Property 2: 表单模板序列化往返一致性**
// **Validates: Requirements 4.6, 4.7**
// *For any* 有效的FormSchema对象，序列化为JSON后再反序列化，应产生与原始对象等价的FormSchema。
func TestFormSchemaSerializationRoundTrip(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	parameters.Rng.Seed(42) // For reproducibility

	properties := gopter.NewProperties(parameters)

	serializer := NewSerializer()

	properties.Property("FormSchema serialization round-trip preserves data", prop.ForAll(
		func(schema *model.FormSchema) bool {
			// Serialize
			jsonStr, err := serializer.SerializeFormSchema(schema)
			if err != nil {
				t.Logf("Serialization error: %v", err)
				return false
			}

			// Deserialize
			restored, err := serializer.DeserializeFormSchema(jsonStr)
			if err != nil {
				t.Logf("Deserialization error: %v", err)
				return false
			}

			// Compare
			return formSchemasEqual(schema, restored)
		},
		genFormSchema(),
	))

	properties.TestingRun(t)
}

// formSchemasEqual compares two FormSchema objects for equality
func formSchemasEqual(a, b *model.FormSchema) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	if len(a.Fields) != len(b.Fields) {
		return false
	}

	for i := range a.Fields {
		if !formFieldsEqual(&a.Fields[i], &b.Fields[i]) {
			return false
		}
	}

	return true
}

// formFieldsEqual compares two FormField objects for equality
func formFieldsEqual(a, b *model.FormField) bool {
	if a.Key != b.Key || a.Label != b.Label || a.Type != b.Type {
		return false
	}
	if a.Required != b.Required || a.Placeholder != b.Placeholder {
		return false
	}

	// Compare DefaultValue (both should be "default" string in our generator)
	if !reflect.DeepEqual(a.DefaultValue, b.DefaultValue) {
		return false
	}

	// Compare Options
	if len(a.Options) != len(b.Options) {
		return false
	}
	for i := range a.Options {
		if a.Options[i].Label != b.Options[i].Label {
			return false
		}
		// Value comparison - both are strings in our generator
		if !reflect.DeepEqual(a.Options[i].Value, b.Options[i].Value) {
			return false
		}
	}

	// Compare Validation
	return fieldValidationsEqual(a.Validation, b.Validation)
}

// fieldValidationsEqual compares two FieldValidation objects for equality
func fieldValidationsEqual(a, b *model.FieldValidation) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Compare MinLength
	if !intPtrsEqual(a.MinLength, b.MinLength) {
		return false
	}
	// Compare MaxLength
	if !intPtrsEqual(a.MaxLength, b.MaxLength) {
		return false
	}
	// Compare Min
	if !float64PtrsEqual(a.Min, b.Min) {
		return false
	}
	// Compare Max
	if !float64PtrsEqual(a.Max, b.Max) {
		return false
	}
	// Compare Pattern
	if !stringPtrsEqual(a.Pattern, b.Pattern) {
		return false
	}

	return true
}

// intPtrsEqual compares two *int values
func intPtrsEqual(a, b *int) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

// float64PtrsEqual compares two *float64 values
func float64PtrsEqual(a, b *float64) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

// stringPtrsEqual compares two *string values
func stringPtrsEqual(a, b *string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}
