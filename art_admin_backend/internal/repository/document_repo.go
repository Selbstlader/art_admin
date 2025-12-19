package repository

import (
	"art_admin_backend/internal/model"

	"gorm.io/gorm"
)

// DocumentRepository 文档仓库
// Document repository
type DocumentRepository struct {
	db *gorm.DB
}

// NewDocumentRepository 创建文档仓库
// Create document repository
func NewDocumentRepository(db *gorm.DB) *DocumentRepository {
	return &DocumentRepository{db: db}
}

// Create 创建文档记录
// Create document record
func (r *DocumentRepository) Create(doc *model.ProjectDocument) error {
	return r.db.Create(doc).Error
}

// GetByID 根据ID获取文档
// Get document by ID
func (r *DocumentRepository) GetByID(id uint) (*model.ProjectDocument, error) {
	var doc model.ProjectDocument
	err := r.db.First(&doc, id).Error
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

// GetByIDWithProject 根据ID获取文档（包含项目信息）
// Get document by ID with project info
func (r *DocumentRepository) GetByIDWithProject(id uint) (*model.ProjectDocument, error) {
	var doc model.ProjectDocument
	err := r.db.Preload("Project").First(&doc, id).Error
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

// Update 更新文档
// Update document
func (r *DocumentRepository) Update(doc *model.ProjectDocument) error {
	return r.db.Save(doc).Error
}

// UpdateAnalysisResult 更新分析结果
// Update analysis result
func (r *DocumentRepository) UpdateAnalysisResult(id uint, updates map[string]interface{}) error {
	return r.db.Model(&model.ProjectDocument{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除文档
// Delete document
func (r *DocumentRepository) Delete(id uint) error {
	return r.db.Delete(&model.ProjectDocument{}, id).Error
}

// BatchDelete 批量删除文档
// Batch delete documents
func (r *DocumentRepository) BatchDelete(ids []uint) error {
	return r.db.Delete(&model.ProjectDocument{}, ids).Error
}

// DocumentFilter 文档筛选条件
// Document filter criteria
type DocumentFilter struct {
	ProjectID      uint
	AnalysisStatus string
	FileType       string
}

// List 获取文档列表
// Get document list
func (r *DocumentRepository) List(page, pageSize int, filter DocumentFilter) ([]model.ProjectDocument, int64, error) {
	var docs []model.ProjectDocument
	var total int64

	query := r.db.Model(&model.ProjectDocument{})

	// 应用筛选条件 / Apply filter criteria
	if filter.ProjectID > 0 {
		query = query.Where("project_id = ?", filter.ProjectID)
	}
	if filter.AnalysisStatus != "" {
		query = query.Where("analysis_status = ?", filter.AnalysisStatus)
	}
	if filter.FileType != "" {
		query = query.Where("file_type = ?", filter.FileType)
	}

	// 获取总数 / Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询 / Paginated query
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&docs).Error; err != nil {
		return nil, 0, err
	}

	return docs, total, nil
}

// GetByProjectID 根据项目ID获取所有文档
// Get all documents by project ID
func (r *DocumentRepository) GetByProjectID(projectID uint) ([]model.ProjectDocument, error) {
	var docs []model.ProjectDocument
	err := r.db.Where("project_id = ?", projectID).Order("created_at DESC").Find(&docs).Error
	return docs, err
}

// CountByProjectID 统计项目文档数量
// Count documents by project ID
func (r *DocumentRepository) CountByProjectID(projectID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.ProjectDocument{}).Where("project_id = ?", projectID).Count(&count).Error
	return count, err
}

// GetPendingDocuments 获取待分析的文档
// Get pending documents
func (r *DocumentRepository) GetPendingDocuments(limit int) ([]model.ProjectDocument, error) {
	var docs []model.ProjectDocument
	err := r.db.Where("analysis_status = ?", "pending").
		Order("created_at ASC").
		Limit(limit).
		Find(&docs).Error
	return docs, err
}
