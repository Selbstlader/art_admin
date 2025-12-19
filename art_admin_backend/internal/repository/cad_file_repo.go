package repository

import (
	"art_admin_backend/internal/model"

	"gorm.io/gorm"
)

// CadFileRepository CAD文件仓库
// CAD file repository for database operations
type CadFileRepository struct {
	db *gorm.DB
}

// NewCadFileRepository 创建CAD文件仓库
// Create new CAD file repository
func NewCadFileRepository(db *gorm.DB) *CadFileRepository {
	return &CadFileRepository{db: db}
}

// Create 创建CAD文件记录
// Create CAD file record
func (r *CadFileRepository) Create(cadFile *model.CadFile) error {
	return r.db.Create(cadFile).Error
}

// GetByID 根据ID获取CAD文件
// Get CAD file by ID
func (r *CadFileRepository) GetByID(id uint) (*model.CadFile, error) {
	var cadFile model.CadFile
	err := r.db.First(&cadFile, id).Error
	if err != nil {
		return nil, err
	}
	return &cadFile, nil
}

// GetByProjectID 根据项目ID获取CAD文件列表
// Get CAD files by project ID
func (r *CadFileRepository) GetByProjectID(projectID uint) ([]model.CadFile, error) {
	var cadFiles []model.CadFile
	err := r.db.Where("project_id = ?", projectID).Find(&cadFiles).Error
	return cadFiles, err
}

// Update 更新CAD文件
// Update CAD file
func (r *CadFileRepository) Update(cadFile *model.CadFile) error {
	return r.db.Save(cadFile).Error
}

// UpdateParseStatus 更新解析状态
// Update parse status
func (r *CadFileRepository) UpdateParseStatus(id uint, status string, errorMsg string) error {
	updates := map[string]interface{}{
		"parse_status":  status,
		"error_message": errorMsg,
	}
	return r.db.Model(&model.CadFile{}).Where("id = ?", id).Updates(updates).Error
}

// UpdateParseResult 更新解析结果
// Update parse result
func (r *CadFileRepository) UpdateParseResult(id uint, parsedPath string, layerCount int, layers string, has3D bool) error {
	updates := map[string]interface{}{
		"parsed_path":  parsedPath,
		"layer_count":  layerCount,
		"layers":       layers,
		"has_3d":       has3D, // 使用数据库列名 has_3d / Use database column name
		"parse_status": "completed",
	}
	return r.db.Model(&model.CadFile{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除CAD文件
// Delete CAD file
func (r *CadFileRepository) Delete(id uint) error {
	return r.db.Delete(&model.CadFile{}, id).Error
}

// DeleteByProjectID 根据项目ID删除所有CAD文件
// Delete all CAD files by project ID
func (r *CadFileRepository) DeleteByProjectID(projectID uint) error {
	return r.db.Where("project_id = ?", projectID).Delete(&model.CadFile{}).Error
}

// List 分页获取CAD文件列表
// List CAD files with pagination
func (r *CadFileRepository) List(projectID uint, page, pageSize int, parseStatus, fileFormat string) ([]model.CadFile, int64, error) {
	var cadFiles []model.CadFile
	var total int64

	query := r.db.Model(&model.CadFile{}).Where("project_id = ?", projectID)

	// 应用筛选条件 / Apply filters
	if parseStatus != "" {
		query = query.Where("parse_status = ?", parseStatus)
	}
	if fileFormat != "" {
		query = query.Where("file_format = ?", fileFormat)
	}

	// 获取总数 / Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询 / Paginated query
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&cadFiles).Error; err != nil {
		return nil, 0, err
	}

	return cadFiles, total, nil
}

// BatchDelete 批量删除CAD文件
// Batch delete CAD files
func (r *CadFileRepository) BatchDelete(ids []uint) error {
	return r.db.Where("id IN ?", ids).Delete(&model.CadFile{}).Error
}

// GetPendingFiles 获取待解析的文件
// Get pending files for parsing
func (r *CadFileRepository) GetPendingFiles(limit int) ([]model.CadFile, error) {
	var cadFiles []model.CadFile
	err := r.db.Where("parse_status = ?", "pending").
		Order("created_at ASC").
		Limit(limit).
		Find(&cadFiles).Error
	return cadFiles, err
}

// CountByProjectID 统计项目下的CAD文件数量
// Count CAD files by project ID
func (r *CadFileRepository) CountByProjectID(projectID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.CadFile{}).Where("project_id = ?", projectID).Count(&count).Error
	return count, err
}

// GetLayersByID 获取CAD文件的图层信息
// Get layers by CAD file ID
func (r *CadFileRepository) GetLayersByID(id uint) (string, error) {
	var cadFile model.CadFile
	err := r.db.Select("layers").First(&cadFile, id).Error
	if err != nil {
		return "", err
	}
	return cadFile.Layers, nil
}
