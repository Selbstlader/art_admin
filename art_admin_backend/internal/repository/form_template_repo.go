package repository

import (
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/database"

	"gorm.io/gorm"
)

// FormTemplateRepository 表单模板仓储
type FormTemplateRepository struct{}

// NewFormTemplateRepository 创建表单模板仓储实例
func NewFormTemplateRepository() *FormTemplateRepository {
	return &FormTemplateRepository{}
}

// Create 创建表单模板
func (r *FormTemplateRepository) Create(template *model.FormTemplate) error {
	return database.DB.Create(template).Error
}

// Update 更新表单模板
func (r *FormTemplateRepository) Update(template *model.FormTemplate) error {
	return database.DB.Save(template).Error
}

// FindByID 根据ID查询表单模板
func (r *FormTemplateRepository) FindByID(id int64) (*model.FormTemplate, error) {
	var template model.FormTemplate
	err := database.DB.First(&template, id).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// FindByCode 根据编码查询表单模板
func (r *FormTemplateRepository) FindByCode(code string) (*model.FormTemplate, error) {
	var template model.FormTemplate
	err := database.DB.Where("code = ?", code).First(&template).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// Delete 删除表单模板（软删除）
func (r *FormTemplateRepository) Delete(id int64) error {
	return database.DB.Delete(&model.FormTemplate{}, id).Error
}

// FindWithPagination 分页查询表单模板列表
func (r *FormTemplateRepository) FindWithPagination(query map[string]interface{}, current, size int) ([]model.FormTemplate, int64, error) {
	var templates []model.FormTemplate
	var total int64

	db := database.DB.Model(&model.FormTemplate{})

	// 动态条件查询
	if name, ok := query["name"]; ok && name != "" {
		db = db.Where("name LIKE ?", "%"+name.(string)+"%")
	}
	if code, ok := query["code"]; ok && code != "" {
		db = db.Where("code LIKE ?", "%"+code.(string)+"%")
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
	err := db.Offset(offset).Limit(size).Order("id DESC").Find(&templates).Error

	return templates, total, err
}

// FindAllActive 查询所有启用的表单模板
func (r *FormTemplateRepository) FindAllActive() ([]model.FormTemplate, error) {
	var templates []model.FormTemplate
	err := database.DB.Where("status = ?", model.FormTemplateStatusActive).
		Order("id DESC").Find(&templates).Error
	return templates, err
}

// UpdateStatus 更新表单模板状态
func (r *FormTemplateRepository) UpdateStatus(id int64, status string) error {
	return database.DB.Model(&model.FormTemplate{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// ExistsByCode 检查编码是否已存在
func (r *FormTemplateRepository) ExistsByCode(code string) (bool, error) {
	var count int64
	err := database.DB.Model(&model.FormTemplate{}).
		Where("code = ?", code).
		Count(&count).Error
	return count > 0, err
}

// ExistsByCodeExcludeID 检查编码是否已存在（排除指定ID）
func (r *FormTemplateRepository) ExistsByCodeExcludeID(code string, excludeID int64) (bool, error) {
	var count int64
	err := database.DB.Model(&model.FormTemplate{}).
		Where("code = ? AND id != ?", code, excludeID).
		Count(&count).Error
	return count > 0, err
}

// IsInUse 检查表单模板是否被流程定义引用
func (r *FormTemplateRepository) IsInUse(id int64) (bool, error) {
	var count int64
	err := database.DB.Model(&model.ProcessDefinition{}).
		Where("form_template_id = ?", id).
		Count(&count).Error
	return count > 0, err
}

// Transaction 执行事务
func (r *FormTemplateRepository) Transaction(fn func(tx *gorm.DB) error) error {
	return database.DB.Transaction(fn)
}
