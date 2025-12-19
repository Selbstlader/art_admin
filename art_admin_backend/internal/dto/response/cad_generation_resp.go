package response

import "time"

// CadGenerationResp CAD生成任务响应
// CAD generation task response
type CadGenerationResp struct {
	ID              uint                     `json:"id"`
	ProjectID       uint                     `json:"projectId"`
	ProjectName     string                   `json:"projectName,omitempty"`
	DocumentIDs     []uint                   `json:"documentIds"`
	GenerationType  string                   `json:"generationType"`
	GenerationLabel string                   `json:"generationLabel"` // 生成类型标签 / Generation type label
	Prompt          string                   `json:"prompt"`
	Parameters      *CadGenerationParamsResp `json:"parameters,omitempty"`
	Status          string                   `json:"status"`
	StatusLabel     string                   `json:"statusLabel"` // 状态标签 / Status label
	Progress        int                      `json:"progress"`
	ResultFileID    *uint                    `json:"resultFileId,omitempty"`
	ResultFile      *CadFileBriefResp        `json:"resultFile,omitempty"`
	ResultFilePath  string                   `json:"resultFilePath,omitempty"`
	PreviewImageURL string                   `json:"previewImageUrl,omitempty"`
	ErrorMessage    string                   `json:"errorMessage,omitempty"`
	AIModel         string                   `json:"aiModel,omitempty"`
	ProcessingTime  int                      `json:"processingTime"`
	CreatedBy       uint                     `json:"createdBy"`
	CreatedAt       time.Time                `json:"createdAt"`
	UpdatedAt       time.Time                `json:"updatedAt"`
}

// CadGenerationParamsResp CAD生成参数响应
// CAD generation parameters response
type CadGenerationParamsResp struct {
	Width        float64  `json:"width"`
	Height       float64  `json:"height"`
	Scale        string   `json:"scale"`
	Style        string   `json:"style"`
	RoomTypes    []string `json:"roomTypes"`
	Requirements string   `json:"requirements"`
}

// CadFileBriefResp CAD文件简要响应
// CAD file brief response
type CadFileBriefResp struct {
	ID         uint   `json:"id"`
	FileName   string `json:"fileName"`
	FileFormat string `json:"fileFormat"`
	FilePath   string `json:"filePath"`
}

// CadGenerationListResp CAD生成任务列表响应
// CAD generation task list response
type CadGenerationListResp struct {
	Records []CadGenerationResp `json:"records"`
	Current int                 `json:"current"`
	Size    int                 `json:"size"`
	Total   int64               `json:"total"`
}

// CadGenerationTypeOption 生成类型选项
// Generation type option
type CadGenerationTypeOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// GetGenerationTypeOptions 获取生成类型选项列表
// Get generation type options
func GetGenerationTypeOptions() []CadGenerationTypeOption {
	return []CadGenerationTypeOption{
		{Value: "floor_plan", Label: "平面布局图"},
		{Value: "elevation", Label: "立面图"},
		{Value: "detail", Label: "节点详图"},
		{Value: "ceiling", Label: "天花图"},
		{Value: "electric", Label: "电气图"},
	}
}

// GetGenerationTypeLabel 获取生成类型标签
// Get generation type label
func GetGenerationTypeLabel(genType string) string {
	labels := map[string]string{
		"floor_plan": "平面布局图",
		"elevation":  "立面图",
		"detail":     "节点详图",
		"ceiling":    "天花图",
		"electric":   "电气图",
	}
	if label, ok := labels[genType]; ok {
		return label
	}
	return genType
}

// GetStatusLabel 获取状态标签
// Get status label
func GetStatusLabel(status string) string {
	labels := map[string]string{
		"pending":    "待处理",
		"processing": "处理中",
		"completed":  "已完成",
		"failed":     "失败",
	}
	if label, ok := labels[status]; ok {
		return label
	}
	return status
}
