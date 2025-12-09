package workflow

import (
	"fmt"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"

	"art_admin_backend/internal/model"
)

// ============================================================================
// Generators for Valid Process Graphs
// ============================================================================

// genValidProcessGraph generates a valid ProcessGraph with exactly one start node,
// at least one end node, and all nodes reachable and connected to end nodes.
func genValidProcessGraph() gopter.Gen {
	return gopter.CombineGens(
		gen.IntRange(0, 5), // number of middle nodes
		gen.Bool(),         // whether to include global props
	).Map(func(vals []interface{}) *model.ProcessGraph {
		middleNodeCount := vals[0].(int)
		includeGlobalProps := vals[1].(bool)

		nodes := make([]model.ProcessNode, 0)
		edges := make([]model.ProcessEdge, 0)

		// Add start node
		startNode := model.ProcessNode{
			ID:       "start",
			Type:     model.NodeTypeStart,
			Name:     "开始",
			Position: model.Position{X: 0, Y: 0},
		}
		nodes = append(nodes, startNode)

		// Add middle nodes (approval nodes)
		for i := 0; i < middleNodeCount; i++ {
			node := model.ProcessNode{
				ID:       fmt.Sprintf("node_%d", i),
				Type:     model.NodeTypeApproval,
				Name:     fmt.Sprintf("审批节点%d", i),
				Position: model.Position{X: float64(i+1) * 100, Y: 0},
				Properties: &model.NodeProperties{
					ApprovalMode: model.ApprovalModeOrSign,
					RejectAction: model.RejectActionTerminate,
				},
			}
			nodes = append(nodes, node)
		}

		// Add end node
		endNode := model.ProcessNode{
			ID:       "end",
			Type:     model.NodeTypeEnd,
			Name:     "结束",
			Position: model.Position{X: float64(middleNodeCount+1) * 100, Y: 0},
		}
		nodes = append(nodes, endNode)

		// Create linear edges: start -> node_0 -> node_1 -> ... -> end
		if middleNodeCount == 0 {
			// Direct connection from start to end
			edges = append(edges, model.ProcessEdge{
				ID:     "edge_start_end",
				Source: "start",
				Target: "end",
			})
		} else {
			// Start to first middle node
			edges = append(edges, model.ProcessEdge{
				ID:     "edge_start_0",
				Source: "start",
				Target: "node_0",
			})

			// Middle nodes chain
			for i := 0; i < middleNodeCount-1; i++ {
				edges = append(edges, model.ProcessEdge{
					ID:     fmt.Sprintf("edge_%d_%d", i, i+1),
					Source: fmt.Sprintf("node_%d", i),
					Target: fmt.Sprintf("node_%d", i+1),
				})
			}

			// Last middle node to end
			edges = append(edges, model.ProcessEdge{
				ID:     fmt.Sprintf("edge_%d_end", middleNodeCount-1),
				Source: fmt.Sprintf("node_%d", middleNodeCount-1),
				Target: "end",
			})
		}

		graph := &model.ProcessGraph{
			Nodes: nodes,
			Edges: edges,
		}

		if includeGlobalProps {
			graph.GlobalProps = &model.GlobalProperties{
				AllowWithdraw: true,
			}
		}

		return graph
	})
}

// genInvalidGraphNoStartNode generates a graph without a start node
func genInvalidGraphNoStartNode() gopter.Gen {
	return gopter.CombineGens(
		gen.IntRange(1, 3), // number of approval nodes
	).Map(func(vals []interface{}) *model.ProcessGraph {
		approvalCount := vals[0].(int)
		nodes := make([]model.ProcessNode, 0)
		edges := make([]model.ProcessEdge, 0)

		// Add approval nodes (no start node)
		for i := 0; i < approvalCount; i++ {
			nodes = append(nodes, model.ProcessNode{
				ID:       fmt.Sprintf("approval_%d", i),
				Type:     model.NodeTypeApproval,
				Name:     fmt.Sprintf("审批%d", i),
				Position: model.Position{X: float64(i) * 100, Y: 0},
			})
		}

		// Add end node
		nodes = append(nodes, model.ProcessNode{
			ID:       "end",
			Type:     model.NodeTypeEnd,
			Name:     "结束",
			Position: model.Position{X: float64(approvalCount) * 100, Y: 0},
		})

		// Connect approval nodes to end
		for i := 0; i < approvalCount; i++ {
			edges = append(edges, model.ProcessEdge{
				ID:     fmt.Sprintf("edge_%d", i),
				Source: fmt.Sprintf("approval_%d", i),
				Target: "end",
			})
		}

		return &model.ProcessGraph{
			Nodes: nodes,
			Edges: edges,
		}
	})
}

// genInvalidGraphMultipleStartNodes generates a graph with multiple start nodes
func genInvalidGraphMultipleStartNodes() gopter.Gen {
	return gopter.CombineGens(
		gen.IntRange(2, 4), // number of start nodes
	).Map(func(vals []interface{}) *model.ProcessGraph {
		startCount := vals[0].(int)
		nodes := make([]model.ProcessNode, 0)
		edges := make([]model.ProcessEdge, 0)

		// Add multiple start nodes
		for i := 0; i < startCount; i++ {
			nodes = append(nodes, model.ProcessNode{
				ID:       fmt.Sprintf("start_%d", i),
				Type:     model.NodeTypeStart,
				Name:     fmt.Sprintf("开始%d", i),
				Position: model.Position{X: float64(i) * 100, Y: 0},
			})
		}

		// Add end node
		nodes = append(nodes, model.ProcessNode{
			ID:       "end",
			Type:     model.NodeTypeEnd,
			Name:     "结束",
			Position: model.Position{X: float64(startCount) * 100, Y: 0},
		})

		// Connect all starts to end
		for i := 0; i < startCount; i++ {
			edges = append(edges, model.ProcessEdge{
				ID:     fmt.Sprintf("edge_%d", i),
				Source: fmt.Sprintf("start_%d", i),
				Target: "end",
			})
		}

		return &model.ProcessGraph{
			Nodes: nodes,
			Edges: edges,
		}
	})
}

// genInvalidGraphNoEndNode generates a graph without an end node
func genInvalidGraphNoEndNode() gopter.Gen {
	return gopter.CombineGens(
		gen.IntRange(1, 3), // number of approval nodes
	).Map(func(vals []interface{}) *model.ProcessGraph {
		approvalCount := vals[0].(int)
		nodes := make([]model.ProcessNode, 0)
		edges := make([]model.ProcessEdge, 0)

		// Add start node
		nodes = append(nodes, model.ProcessNode{
			ID:       "start",
			Type:     model.NodeTypeStart,
			Name:     "开始",
			Position: model.Position{X: 0, Y: 0},
		})

		// Add approval nodes (no end node)
		for i := 0; i < approvalCount; i++ {
			nodes = append(nodes, model.ProcessNode{
				ID:       fmt.Sprintf("approval_%d", i),
				Type:     model.NodeTypeApproval,
				Name:     fmt.Sprintf("审批%d", i),
				Position: model.Position{X: float64(i+1) * 100, Y: 0},
			})
		}

		// Connect start to first approval
		edges = append(edges, model.ProcessEdge{
			ID:     "edge_start",
			Source: "start",
			Target: "approval_0",
		})

		// Chain approval nodes
		for i := 0; i < approvalCount-1; i++ {
			edges = append(edges, model.ProcessEdge{
				ID:     fmt.Sprintf("edge_%d", i),
				Source: fmt.Sprintf("approval_%d", i),
				Target: fmt.Sprintf("approval_%d", i+1),
			})
		}

		return &model.ProcessGraph{
			Nodes: nodes,
			Edges: edges,
		}
	})
}

// genInvalidGraphWithUnreachableNode generates a graph with an unreachable node
func genInvalidGraphWithUnreachableNode() gopter.Gen {
	return gopter.CombineGens(
		gen.IntRange(1, 3), // number of unreachable nodes
	).Map(func(vals []interface{}) *model.ProcessGraph {
		unreachableCount := vals[0].(int)
		nodes := make([]model.ProcessNode, 0)
		edges := make([]model.ProcessEdge, 0)

		// Add start node
		nodes = append(nodes, model.ProcessNode{
			ID:       "start",
			Type:     model.NodeTypeStart,
			Name:     "开始",
			Position: model.Position{X: 0, Y: 0},
		})

		// Add reachable approval node
		nodes = append(nodes, model.ProcessNode{
			ID:       "approval_reachable",
			Type:     model.NodeTypeApproval,
			Name:     "可达审批",
			Position: model.Position{X: 100, Y: 0},
		})

		// Add unreachable nodes
		for i := 0; i < unreachableCount; i++ {
			nodes = append(nodes, model.ProcessNode{
				ID:       fmt.Sprintf("unreachable_%d", i),
				Type:     model.NodeTypeApproval,
				Name:     fmt.Sprintf("不可达%d", i),
				Position: model.Position{X: 100, Y: float64(i+1) * 100},
			})
		}

		// Add end node
		nodes = append(nodes, model.ProcessNode{
			ID:       "end",
			Type:     model.NodeTypeEnd,
			Name:     "结束",
			Position: model.Position{X: 200, Y: 0},
		})

		// Connect start -> reachable -> end
		edges = append(edges, model.ProcessEdge{
			ID:     "edge_start",
			Source: "start",
			Target: "approval_reachable",
		})
		edges = append(edges, model.ProcessEdge{
			ID:     "edge_end",
			Source: "approval_reachable",
			Target: "end",
		})

		// Unreachable nodes have no incoming edges from start

		return &model.ProcessGraph{
			Nodes: nodes,
			Edges: edges,
		}
	})
}

// genInvalidGraphWithDeadEndNode generates a graph with a dead-end node (can't reach end)
func genInvalidGraphWithDeadEndNode() gopter.Gen {
	return gopter.CombineGens(
		gen.IntRange(1, 3), // number of dead-end nodes
	).Map(func(vals []interface{}) *model.ProcessGraph {
		deadEndCount := vals[0].(int)
		nodes := make([]model.ProcessNode, 0)
		edges := make([]model.ProcessEdge, 0)

		// Add start node
		nodes = append(nodes, model.ProcessNode{
			ID:       "start",
			Type:     model.NodeTypeStart,
			Name:     "开始",
			Position: model.Position{X: 0, Y: 0},
		})

		// Add normal approval node that leads to end
		nodes = append(nodes, model.ProcessNode{
			ID:       "approval_normal",
			Type:     model.NodeTypeApproval,
			Name:     "正常审批",
			Position: model.Position{X: 100, Y: 0},
		})

		// Add dead-end nodes (reachable but can't reach end)
		for i := 0; i < deadEndCount; i++ {
			nodes = append(nodes, model.ProcessNode{
				ID:       fmt.Sprintf("deadend_%d", i),
				Type:     model.NodeTypeApproval,
				Name:     fmt.Sprintf("死胡同%d", i),
				Position: model.Position{X: 100, Y: float64(i+1) * 100},
			})
		}

		// Add end node
		nodes = append(nodes, model.ProcessNode{
			ID:       "end",
			Type:     model.NodeTypeEnd,
			Name:     "结束",
			Position: model.Position{X: 200, Y: 0},
		})

		// Connect start -> normal -> end
		edges = append(edges, model.ProcessEdge{
			ID:     "edge_start_normal",
			Source: "start",
			Target: "approval_normal",
		})
		edges = append(edges, model.ProcessEdge{
			ID:     "edge_normal_end",
			Source: "approval_normal",
			Target: "end",
		})

		// Connect start -> dead-end nodes (reachable but no path to end)
		for i := 0; i < deadEndCount; i++ {
			edges = append(edges, model.ProcessEdge{
				ID:     fmt.Sprintf("edge_start_deadend_%d", i),
				Source: "start",
				Target: fmt.Sprintf("deadend_%d", i),
			})
		}

		return &model.ProcessGraph{
			Nodes: nodes,
			Edges: edges,
		}
	})
}

// ============================================================================
// Property-Based Tests
// ============================================================================

// **Feature: oa-workflow-engine, Property 3: 流程结构完整性验证**
// **Validates: Requirements 1.1**
// *For any* ProcessGraph，验证函数应正确识别：必须有且仅有一个开始节点、至少一个结束节点、所有节点可达、无孤立节点。
func TestValidProcessGraphValidation(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	parameters.Rng.Seed(42)

	properties := gopter.NewProperties(parameters)
	validator := NewValidator()

	// Property: Valid graphs should pass validation
	properties.Property("Valid process graphs pass validation", prop.ForAll(
		func(graph *model.ProcessGraph) bool {
			result := validator.ValidateProcessGraph(graph)
			if !result.Valid {
				t.Logf("Valid graph failed validation: %v", result.Errors)
				return false
			}
			return true
		},
		genValidProcessGraph(),
	))

	properties.TestingRun(t)
}

// Test that graphs without start nodes are rejected
func TestInvalidGraphNoStartNode(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	parameters.Rng.Seed(42)

	properties := gopter.NewProperties(parameters)
	validator := NewValidator()

	properties.Property("Graphs without start node are invalid", prop.ForAll(
		func(graph *model.ProcessGraph) bool {
			result := validator.ValidateProcessGraph(graph)
			if result.Valid {
				t.Logf("Graph without start node passed validation unexpectedly")
				return false
			}
			// Check that the error is about missing start node
			hasStartError := false
			for _, err := range result.Errors {
				if err.Code == "NO_START_NODE" {
					hasStartError = true
					break
				}
			}
			return hasStartError
		},
		genInvalidGraphNoStartNode(),
	))

	properties.TestingRun(t)
}

// Test that graphs with multiple start nodes are rejected
func TestInvalidGraphMultipleStartNodes(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	parameters.Rng.Seed(42)

	properties := gopter.NewProperties(parameters)
	validator := NewValidator()

	properties.Property("Graphs with multiple start nodes are invalid", prop.ForAll(
		func(graph *model.ProcessGraph) bool {
			result := validator.ValidateProcessGraph(graph)
			if result.Valid {
				t.Logf("Graph with multiple start nodes passed validation unexpectedly")
				return false
			}
			// Check that the error is about multiple start nodes
			hasMultipleStartError := false
			for _, err := range result.Errors {
				if err.Code == "MULTIPLE_START_NODES" {
					hasMultipleStartError = true
					break
				}
			}
			return hasMultipleStartError
		},
		genInvalidGraphMultipleStartNodes(),
	))

	properties.TestingRun(t)
}

// Test that graphs without end nodes are rejected
func TestInvalidGraphNoEndNode(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	parameters.Rng.Seed(42)

	properties := gopter.NewProperties(parameters)
	validator := NewValidator()

	properties.Property("Graphs without end node are invalid", prop.ForAll(
		func(graph *model.ProcessGraph) bool {
			result := validator.ValidateProcessGraph(graph)
			if result.Valid {
				t.Logf("Graph without end node passed validation unexpectedly")
				return false
			}
			// Check that the error is about missing end node or dead end
			hasEndError := false
			for _, err := range result.Errors {
				if err.Code == "NO_END_NODE" || err.Code == "DEAD_END_NODE" {
					hasEndError = true
					break
				}
			}
			return hasEndError
		},
		genInvalidGraphNoEndNode(),
	))

	properties.TestingRun(t)
}

// Test that graphs with unreachable nodes are rejected
func TestInvalidGraphWithUnreachableNode(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	parameters.Rng.Seed(42)

	properties := gopter.NewProperties(parameters)
	validator := NewValidator()

	properties.Property("Graphs with unreachable nodes are invalid", prop.ForAll(
		func(graph *model.ProcessGraph) bool {
			result := validator.ValidateProcessGraph(graph)
			if result.Valid {
				t.Logf("Graph with unreachable node passed validation unexpectedly")
				return false
			}
			// Check that the error is about unreachable node
			hasUnreachableError := false
			for _, err := range result.Errors {
				if err.Code == "UNREACHABLE_NODE" {
					hasUnreachableError = true
					break
				}
			}
			return hasUnreachableError
		},
		genInvalidGraphWithUnreachableNode(),
	))

	properties.TestingRun(t)
}

// Test that graphs with dead-end nodes are rejected
func TestInvalidGraphWithDeadEndNode(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	parameters.Rng.Seed(42)

	properties := gopter.NewProperties(parameters)
	validator := NewValidator()

	properties.Property("Graphs with dead-end nodes are invalid", prop.ForAll(
		func(graph *model.ProcessGraph) bool {
			result := validator.ValidateProcessGraph(graph)
			if result.Valid {
				t.Logf("Graph with dead-end node passed validation unexpectedly")
				return false
			}
			// Check that the error is about dead-end node
			hasDeadEndError := false
			for _, err := range result.Errors {
				if err.Code == "DEAD_END_NODE" {
					hasDeadEndError = true
					break
				}
			}
			return hasDeadEndError
		},
		genInvalidGraphWithDeadEndNode(),
	))

	properties.TestingRun(t)
}

// Test nil graph handling
func TestNilGraphValidation(t *testing.T) {
	validator := NewValidator()
	result := validator.ValidateProcessGraph(nil)

	if result.Valid {
		t.Error("Nil graph should be invalid")
	}

	hasNilError := false
	for _, err := range result.Errors {
		if err.Code == "GRAPH_NIL" {
			hasNilError = true
			break
		}
	}

	if !hasNilError {
		t.Error("Expected GRAPH_NIL error for nil graph")
	}
}

// Test empty nodes graph handling
func TestEmptyNodesGraphValidation(t *testing.T) {
	validator := NewValidator()
	graph := &model.ProcessGraph{
		Nodes: []model.ProcessNode{},
		Edges: []model.ProcessEdge{},
	}
	result := validator.ValidateProcessGraph(graph)

	if result.Valid {
		t.Error("Graph with no nodes should be invalid")
	}

	hasNoNodesError := false
	for _, err := range result.Errors {
		if err.Code == "NO_NODES" {
			hasNoNodesError = true
			break
		}
	}

	if !hasNoNodesError {
		t.Error("Expected NO_NODES error for empty graph")
	}
}
