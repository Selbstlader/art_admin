package response

import (
	"art_admin_backend/internal/model"
	"encoding/json"
	"time"
)

// ConstructionAnnotationResponse 施工图标注响应
// Construction annotation response
type ConstructionAnnotationResponse struct {
	ID               uint                    `json:"id"`                // 标注ID / Annotation ID
	ProjectID        uint                    `json:"projectId"`         // 项目ID / Project ID
	CadFileID        uint                    `json:"cadFileId"`         // CAD文件ID / CAD file ID
	ImagePath        string                  `json:"imagePath"`         // 施工图图片路径 / Construction drawing image path
	AnalysisStatus   string                  `json:"analysisStatus"`    // 分析状态 / Analysis status
	ElementsDetected []model.DetectedElement `json:"elementsDetected"`  // 检测到的元素 / Detected elements
	Annotations      []model.AnnotationItem  `json:"annotations"`       // 标注数据 / Annotation data
	ErrorMessage     string                  `json:"errorMessage"`      // 错误信息 / Error message
	CreatedAt        time.Time               `json:"createdAt"`         // 创建时间 / Creation time
	UpdatedAt        time.Time               `json:"updatedAt"`         // 更新时间 / Update time
	Project          *ProjectBasicInfo       `json:"project,omitempty"` // 项目基本信息 / Project basic info
	CadFile          *CadFileBasicInfo       `json:"cadFile,omitempty"` // CAD文件基本信息 / CAD file basic info
}

// ProjectBasicInfo 项目基本信息
// Project basic information
type ProjectBasicInfo struct {
	ID          uint   `json:"id"`          // 项目ID / Project ID
	Name        string `json:"name"`        // 项目名称 / Project name
	Description string `json:"description"` // 项目描述 / Project description
}

// CadFileBasicInfo CAD文件基本信息
// CAD file basic information
type CadFileBasicInfo struct {
	ID         uint   `json:"id"`         // 文件ID / File ID
	FileName   string `json:"fileName"`   // 文件名 / File name
	FileFormat string `json:"fileFormat"` // 文件格式 / File format
}

// AnnotationListResponse 标注列表响应
// Annotation list response
type AnnotationListResponse struct {
	List       []ConstructionAnnotationResponse `json:"list"`       // 标注列表 / Annotation list
	Total      int64                            `json:"total"`      // 总数 / Total count
	Page       int                              `json:"page"`       // 当前页 / Current page
	PageSize   int                              `json:"pageSize"`   // 每页大小 / Page size
	TotalPages int                              `json:"totalPages"` // 总页数 / Total pages
}

// AnalysisResultResponse 分析结果响应
// Analysis result response
type AnalysisResultResponse struct {
	AnnotationID     uint                    `json:"annotationId"`     // 标注ID / Annotation ID
	AnalysisStatus   string                  `json:"analysisStatus"`   // 分析状态 / Analysis status
	ElementsDetected []model.DetectedElement `json:"elementsDetected"` // 检测到的元素 / Detected elements
	ElementCount     int                     `json:"elementCount"`     // 元素数量 / Element count
	ProcessingTime   string                  `json:"processingTime"`   // 处理时间 / Processing time
	ErrorMessage     string                  `json:"errorMessage"`     // 错误信息 / Error message
}

// ExportResultResponse 导出结果响应
// Export result response
type ExportResultResponse struct {
	FileName    string `json:"fileName"`    // 导出文件名 / Export file name
	FilePath    string `json:"filePath"`    // 文件路径 / File path
	FileSize    int64  `json:"fileSize"`    // 文件大小 / File size
	Format      string `json:"format"`      // 导出格式 / Export format
	DownloadURL string `json:"downloadUrl"` // 下载链接 / Download URL
}

// AnnotationStatsResponse 标注统计响应
// Annotation statistics response
type AnnotationStatsResponse struct {
	TotalAnnotations    int                        `json:"totalAnnotations"`    // 总标注数 / Total annotations
	CompletedAnalysis   int                        `json:"completedAnalysis"`   // 已完成分析数 / Completed analysis count
	PendingAnalysis     int                        `json:"pendingAnalysis"`     // 待分析数 / Pending analysis count
	FailedAnalysis      int                        `json:"failedAnalysis"`      // 失败分析数 / Failed analysis count
	ElementTypeStats    map[string]int             `json:"elementTypeStats"`    // 元素类型统计 / Element type statistics
	AnnotationTypeStats map[string]int             `json:"annotationTypeStats"` // 标注类型统计 / Annotation type statistics
	RecentActivity      []RecentAnnotationActivity `json:"recentActivity"`      // 最近活动 / Recent activity
}

// RecentAnnotationActivity 最近标注活动
// Recent annotation activity
type RecentAnnotationActivity struct {
	AnnotationID uint      `json:"annotationId"` // 标注ID / Annotation ID
	ProjectName  string    `json:"projectName"`  // 项目名称 / Project name
	Action       string    `json:"action"`       // 操作类型 / Action type
	Timestamp    time.Time `json:"timestamp"`    // 时间戳 / Timestamp
}

// ConvertToResponse 将模型转换为响应结构
// Convert model to response structure
func ConvertToResponse(annotation *model.ConstructionAnnotation) *ConstructionAnnotationResponse {
	resp := &ConstructionAnnotationResponse{
		ID:             annotation.ID,
		ProjectID:      annotation.ProjectID,
		CadFileID:      annotation.CadFileID,
		ImagePath:      annotation.ImagePath,
		AnalysisStatus: annotation.AnalysisStatus,
		ErrorMessage:   annotation.ErrorMessage,
		CreatedAt:      annotation.CreatedAt,
		UpdatedAt:      annotation.UpdatedAt,
	}

	// 解析检测到的元素 / Parse detected elements
	if annotation.ElementsDetected != "" {
		var elements []model.DetectedElement
		if err := json.Unmarshal([]byte(annotation.ElementsDetected), &elements); err == nil {
			resp.ElementsDetected = elements
		}
	}

	// 解析标注数据 / Parse annotation data
	if annotation.Annotations != "" {
		var annotations []model.AnnotationItem
		if err := json.Unmarshal([]byte(annotation.Annotations), &annotations); err == nil {
			resp.Annotations = annotations
		}
	}

	// 添加关联信息 / Add associated information
	if annotation.Project != nil {
		resp.Project = &ProjectBasicInfo{
			ID:          annotation.Project.ID,
			Name:        annotation.Project.Name,
			Description: annotation.Project.Description,
		}
	}

	if annotation.CadFile != nil {
		resp.CadFile = &CadFileBasicInfo{
			ID:         annotation.CadFile.ID,
			FileName:   annotation.CadFile.FileName,
			FileFormat: annotation.CadFile.FileFormat,
		}
	}

	return resp
}

// ConvertToListResponse 将模型列表转换为列表响应
// Convert model list to list response
func ConvertToListResponse(annotations []*model.ConstructionAnnotation, total int64, page, pageSize int) *AnnotationListResponse {
	list := make([]ConstructionAnnotationResponse, len(annotations))
	for i, annotation := range annotations {
		list[i] = *ConvertToResponse(annotation)
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return &AnnotationListResponse{
		List:       list,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}
