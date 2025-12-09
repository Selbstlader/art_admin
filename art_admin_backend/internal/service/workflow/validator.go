package workflow

import (
	"fmt"

	"art_admin_backend/internal/model"
)

// ValidationError 验证错误
type ValidationError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// ValidationResult 验证结果
type ValidationResult struct {
	Valid  bool               `json:"valid"`
	Errors []*ValidationError `json:"errors,omitempty"`
}

// Validator 流程结构验证器接口
type Validator interface {
	// ValidateProcessGraph 验证流程图结构完整性
	// Requirements: 1.1 - WHEN 管理员在流程设计器中添加节点并连接 THEN Workflow_Engine SHALL 验证流程结构完整性
	ValidateProcessGraph(graph *model.ProcessGraph) *ValidationResult
}

// validator 验证器实现
type validator struct{}

// NewValidator 创建验证器实例
func NewValidator() Validator {
	return &validator{}
}

// ValidateProcessGraph 验证流程图结构完整性
// Property 3: 流程结构完整性验证
// *For any* ProcessGraph，验证函数应正确识别：必须有且仅有一个开始节点、至少一个结束节点、所有节点可达、无孤立节点。
func (v *validator) ValidateProcessGraph(graph *model.ProcessGraph) *ValidationResult {
	result := &ValidationResult{
		Valid:  true,
		Errors: make([]*ValidationError, 0),
	}

	if graph == nil {
		result.Valid = false
		result.Errors = append(result.Errors, &ValidationError{
			Code:    "GRAPH_NIL",
			Message: "流程图不能为空",
		})
		return result
	}

	// 1. 验证节点列表不为空
	if len(graph.Nodes) == 0 {
		result.Valid = false
		result.Errors = append(result.Errors, &ValidationError{
			Code:    "NO_NODES",
			Message: "流程图必须包含至少一个节点",
		})
		return result
	}

	// 2. 统计开始节点和结束节点数量
	startNodes := make([]string, 0)
	endNodes := make([]string, 0)
	nodeMap := make(map[string]*model.ProcessNode)

	for i := range graph.Nodes {
		node := &graph.Nodes[i]
		nodeMap[node.ID] = node

		switch node.Type {
		case model.NodeTypeStart:
			startNodes = append(startNodes, node.ID)
		case model.NodeTypeEnd:
			endNodes = append(endNodes, node.ID)
		}
	}

	// 3. 验证必须有且仅有一个开始节点
	if len(startNodes) == 0 {
		result.Valid = false
		result.Errors = append(result.Errors, &ValidationError{
			Code:    "NO_START_NODE",
			Message: "流程图必须包含一个开始节点",
		})
	} else if len(startNodes) > 1 {
		result.Valid = false
		result.Errors = append(result.Errors, &ValidationError{
			Code:    "MULTIPLE_START_NODES",
			Message: fmt.Sprintf("流程图只能有一个开始节点，当前有 %d 个", len(startNodes)),
		})
	}

	// 4. 验证至少有一个结束节点
	if len(endNodes) == 0 {
		result.Valid = false
		result.Errors = append(result.Errors, &ValidationError{
			Code:    "NO_END_NODE",
			Message: "流程图必须包含至少一个结束节点",
		})
	}

	// 如果没有开始节点，无法进行可达性分析
	if len(startNodes) != 1 {
		return result
	}

	// 5. 构建邻接表
	adjacencyList := make(map[string][]string)
	reverseAdjacencyList := make(map[string][]string)
	for _, node := range graph.Nodes {
		adjacencyList[node.ID] = make([]string, 0)
		reverseAdjacencyList[node.ID] = make([]string, 0)
	}

	for _, edge := range graph.Edges {
		// 验证边的源节点和目标节点存在
		if _, exists := nodeMap[edge.Source]; !exists {
			result.Valid = false
			result.Errors = append(result.Errors, &ValidationError{
				Code:    "INVALID_EDGE_SOURCE",
				Message: fmt.Sprintf("边 %s 的源节点 %s 不存在", edge.ID, edge.Source),
			})
			continue
		}
		if _, exists := nodeMap[edge.Target]; !exists {
			result.Valid = false
			result.Errors = append(result.Errors, &ValidationError{
				Code:    "INVALID_EDGE_TARGET",
				Message: fmt.Sprintf("边 %s 的目标节点 %s 不存在", edge.ID, edge.Target),
			})
			continue
		}

		adjacencyList[edge.Source] = append(adjacencyList[edge.Source], edge.Target)
		reverseAdjacencyList[edge.Target] = append(reverseAdjacencyList[edge.Target], edge.Source)
	}

	// 6. 从开始节点进行BFS，检查所有节点是否可达
	reachableFromStart := v.bfs(startNodes[0], adjacencyList)

	// 7. 从结束节点反向BFS，检查所有节点是否能到达结束节点
	canReachEnd := make(map[string]bool)
	for _, endNode := range endNodes {
		reachable := v.bfs(endNode, reverseAdjacencyList)
		for nodeID := range reachable {
			canReachEnd[nodeID] = true
		}
	}

	// 8. 检查孤立节点（从开始节点不可达的节点）
	for nodeID := range nodeMap {
		if !reachableFromStart[nodeID] {
			result.Valid = false
			result.Errors = append(result.Errors, &ValidationError{
				Code:    "UNREACHABLE_NODE",
				Message: fmt.Sprintf("节点 %s 从开始节点不可达", nodeID),
			})
		}
	}

	// 9. 检查死胡同节点（不能到达任何结束节点的非结束节点）
	for nodeID, node := range nodeMap {
		if node.Type != model.NodeTypeEnd && !canReachEnd[nodeID] {
			result.Valid = false
			result.Errors = append(result.Errors, &ValidationError{
				Code:    "DEAD_END_NODE",
				Message: fmt.Sprintf("节点 %s 无法到达任何结束节点", nodeID),
			})
		}
	}

	return result
}

// bfs 广度优先搜索，返回从起始节点可达的所有节点
func (v *validator) bfs(startNode string, adjacencyList map[string][]string) map[string]bool {
	visited := make(map[string]bool)
	queue := []string{startNode}
	visited[startNode] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, neighbor := range adjacencyList[current] {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, neighbor)
			}
		}
	}

	return visited
}

// IsValidGraph 快捷方法：检查流程图是否有效
func (v *validator) IsValidGraph(graph *model.ProcessGraph) bool {
	result := v.ValidateProcessGraph(graph)
	return result.Valid
}
