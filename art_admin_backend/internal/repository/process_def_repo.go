package repository

import (
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/database"

	"gorm.io/gorm"
)

// ProcessDefinitionRepository 流程定义仓储
type ProcessDefinitionRepository struct{}

// NewProcessDefinitionRepository 创建流程定义仓储实例
func NewProcessDefinitionRepository() *ProcessDefinitionRepository {
	return &ProcessDefinitionRepository{}
}

// Create 创建流程定义
func (r *ProcessDefinitionRepository) Create(def *model.ProcessDefinition) error {
	return database.DB.Create(def).Error
}

// Update 更新流程定义
func (r *ProcessDefinitionRepository) Update(def *model.ProcessDefinition) error {
	return database.DB.Save(def).Error
}

// FindByID 根据ID查询流程定义
func (r *ProcessDefinitionRepository) FindByID(id int64) (*model.ProcessDefinition, error) {
	var def model.ProcessDefinition
	err := database.DB.First(&def, id).Error
	if err != nil {
		return nil, err
	}
	return &def, nil
}

// FindByCode 根据编码查询流程定义
func (r *ProcessDefinitionRepository) FindByCode(code string) (*model.ProcessDefinition, error) {
	var def model.ProcessDefinition
	err := database.DB.Where("code = ?", code).First(&def).Error
	if err != nil {
		return nil, err
	}
	return &def, nil
}

// Delete 删除流程定义（软删除）
func (r *ProcessDefinitionRepository) Delete(id int64) error {
	return database.DB.Delete(&model.ProcessDefinition{}, id).Error
}

// FindWithPagination 分页查询流程定义列表
func (r *ProcessDefinitionRepository) FindWithPagination(query map[string]interface{}, current, size int) ([]model.ProcessDefinition, int64, error) {
	var defs []model.ProcessDefinition
	var total int64

	db := database.DB.Model(&model.ProcessDefinition{})

	// 动态条件查询
	if name, ok := query["name"]; ok && name != "" {
		db = db.Where("name LIKE ?", "%"+name.(string)+"%")
	}
	if code, ok := query["code"]; ok && code != "" {
		db = db.Where("code LIKE ?", "%"+code.(string)+"%")
	}
	if category, ok := query["category"]; ok && category != "" {
		db = db.Where("category = ?", category)
	}
	if status, ok := query["status"]; ok && status != "" {
		db = db.Where("status = ?", status)
	}
	if createdBy, ok := query["createdBy"]; ok && createdBy != nil {
		db = db.Where("created_by = ?", createdBy)
	}

	// 统计总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (current - 1) * size
	err := db.Offset(offset).Limit(size).Order("id DESC").Find(&defs).Error

	return defs, total, err
}

// FindPublishedByCode 根据编码查询已发布的流程定义
func (r *ProcessDefinitionRepository) FindPublishedByCode(code string) (*model.ProcessDefinition, error) {
	var def model.ProcessDefinition
	err := database.DB.Where("code = ? AND status = ?", code, model.ProcessDefStatusPublished).First(&def).Error
	if err != nil {
		return nil, err
	}
	return &def, nil
}

// FindLatestVersionByCode 根据编码查询最新版本的流程定义
func (r *ProcessDefinitionRepository) FindLatestVersionByCode(code string) (*model.ProcessDefinition, error) {
	var def model.ProcessDefinition
	err := database.DB.Where("code = ?", code).Order("version DESC").First(&def).Error
	if err != nil {
		return nil, err
	}
	return &def, nil
}

// GetMaxVersionByCode 获取指定编码的最大版本号
func (r *ProcessDefinitionRepository) GetMaxVersionByCode(code string) (int, error) {
	var maxVersion int
	err := database.DB.Model(&model.ProcessDefinition{}).
		Where("code = ?", code).
		Select("COALESCE(MAX(version), 0)").
		Scan(&maxVersion).Error
	return maxVersion, err
}

// UpdateStatus 更新流程定义状态
func (r *ProcessDefinitionRepository) UpdateStatus(id int64, status string) error {
	return database.DB.Model(&model.ProcessDefinition{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// ExistsByCode 检查编码是否已存在
func (r *ProcessDefinitionRepository) ExistsByCode(code string) (bool, error) {
	var count int64
	err := database.DB.Model(&model.ProcessDefinition{}).
		Where("code = ?", code).
		Count(&count).Error
	return count > 0, err
}

// ExistsByCodeExcludeID 检查编码是否已存在（排除指定ID）
func (r *ProcessDefinitionRepository) ExistsByCodeExcludeID(code string, excludeID int64) (bool, error) {
	var count int64
	err := database.DB.Model(&model.ProcessDefinition{}).
		Where("code = ? AND id != ?", code, excludeID).
		Count(&count).Error
	return count > 0, err
}

// FindByFormTemplateID 根据表单模板ID查询关联的流程定义
func (r *ProcessDefinitionRepository) FindByFormTemplateID(formTemplateID int64) ([]model.ProcessDefinition, error) {
	var defs []model.ProcessDefinition
	err := database.DB.Where("form_template_id = ?", formTemplateID).Find(&defs).Error
	return defs, err
}

// Transaction 执行事务
func (r *ProcessDefinitionRepository) Transaction(fn func(tx *gorm.DB) error) error {
	return database.DB.Transaction(fn)
}
