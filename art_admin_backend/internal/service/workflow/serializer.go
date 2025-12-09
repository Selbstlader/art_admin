package workflow

import (
	"encoding/json"
	"fmt"

	"art_admin_backend/internal/model"
)

// Serializer 序列化服务接口
type Serializer interface {
	// SerializeProcessGraph 将ProcessGraph序列化为JSON字符串
	SerializeProcessGraph(graph *model.ProcessGraph) (string, error)
	// DeserializeProcessGraph 从JSON字符串反序列化为ProcessGraph
	DeserializeProcessGraph(jsonStr string) (*model.ProcessGraph, error)
	// SerializeFormSchema 将FormSchema序列化为JSON字符串
	SerializeFormSchema(schema *model.FormSchema) (string, error)
	// DeserializeFormSchema 从JSON字符串反序列化为FormSchema
	DeserializeFormSchema(jsonStr string) (*model.FormSchema, error)
}

// serializer 序列化服务实现
type serializer struct{}

// NewSerializer 创建序列化服务实例
func NewSerializer() Serializer {
	return &serializer{}
}

// SerializeProcessGraph 将ProcessGraph序列化为JSON字符串
// Requirements: 1.6 - WHEN 流程定义被序列化存储 THEN Workflow_Engine SHALL 使用JSON格式编码流程结构
func (s *serializer) SerializeProcessGraph(graph *model.ProcessGraph) (string, error) {
	if graph == nil {
		return "", fmt.Errorf("process graph cannot be nil")
	}

	data, err := json.Marshal(graph)
	if err != nil {
		return "", fmt.Errorf("failed to serialize process graph: %w", err)
	}

	return string(data), nil
}

// DeserializeProcessGraph 从JSON字符串反序列化为ProcessGraph
// Requirements: 1.7 - WHEN 从存储加载流程定义 THEN Workflow_Engine SHALL 解析JSON并还原完整的Process_Definition对象
func (s *serializer) DeserializeProcessGraph(jsonStr string) (*model.ProcessGraph, error) {
	if jsonStr == "" {
		return nil, fmt.Errorf("json string cannot be empty")
	}

	var graph model.ProcessGraph
	if err := json.Unmarshal([]byte(jsonStr), &graph); err != nil {
		return nil, fmt.Errorf("failed to deserialize process graph: %w", err)
	}

	return &graph, nil
}

// SerializeFormSchema 将FormSchema序列化为JSON字符串
// Requirements: 4.6 - WHEN 表单模板被序列化存储 THEN Form_Engine SHALL 使用JSON格式编码表单结构
func (s *serializer) SerializeFormSchema(schema *model.FormSchema) (string, error) {
	if schema == nil {
		return "", fmt.Errorf("form schema cannot be nil")
	}

	data, err := json.Marshal(schema)
	if err != nil {
		return "", fmt.Errorf("failed to serialize form schema: %w", err)
	}

	return string(data), nil
}

// DeserializeFormSchema 从JSON字符串反序列化为FormSchema
// Requirements: 4.7 - WHEN 从存储加载表单模板 THEN Form_Engine SHALL 解析JSON并还原完整的Form_Template对象
func (s *serializer) DeserializeFormSchema(jsonStr string) (*model.FormSchema, error) {
	if jsonStr == "" {
		return nil, fmt.Errorf("json string cannot be empty")
	}

	var schema model.FormSchema
	if err := json.Unmarshal([]byte(jsonStr), &schema); err != nil {
		return nil, fmt.Errorf("failed to deserialize form schema: %w", err)
	}

	return &schema, nil
}
