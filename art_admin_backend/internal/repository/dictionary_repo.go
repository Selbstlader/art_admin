package repository

import (
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/database"
)

type DictionaryRepository struct{}

func NewDictionaryRepository() *DictionaryRepository {
	return &DictionaryRepository{}
}

// ========== 字典类型相关 ==========

// FindDictionaryTypeByID 根据ID查询字典类型
func (r *DictionaryRepository) FindDictionaryTypeByID(id int64) (*model.DictionaryType, error) {
	var dictType model.DictionaryType
	err := database.DB.First(&dictType, id).Error
	if err != nil {
		return nil, err
	}
	return &dictType, nil
}

// FindDictionaryTypeByCode 根据编码查询字典类型
func (r *DictionaryRepository) FindDictionaryTypeByCode(code string) (*model.DictionaryType, error) {
	var dictType model.DictionaryType
	err := database.DB.Where("type_code = ?", code).First(&dictType).Error
	if err != nil {
		return nil, err
	}
	return &dictType, nil
}

// FindAllDictionaryTypes 查询所有字典类型
func (r *DictionaryRepository) FindAllDictionaryTypes() ([]model.DictionaryType, error) {
	var dictTypes []model.DictionaryType
	err := database.DB.Order("id ASC").Find(&dictTypes).Error
	return dictTypes, err
}

// FindDictionaryTypesWithPagination 分页查询字典类型
func (r *DictionaryRepository) FindDictionaryTypesWithPagination(query map[string]interface{}, current, size int) ([]model.DictionaryType, int64, error) {
	var dictTypes []model.DictionaryType
	var total int64

	db := database.DB.Model(&model.DictionaryType{})

	// 构建查询条件
	for key, value := range query {
		switch key {
		case "typeName":
			db = db.Where("type_name LIKE ?", "%"+value.(string)+"%")
		case "typeCode":
			db = db.Where("type_code LIKE ?", "%"+value.(string)+"%")
		case "description":
			db = db.Where("description LIKE ?", "%"+value.(string)+"%")
		case "enabled":
			db = db.Where("enabled = ?", value)
		default:
			db = db.Where(key+" = ?", value)
		}
	}

	// 查询总数
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (current - 1) * size
	err = db.Order("id ASC").Offset(offset).Limit(size).Find(&dictTypes).Error
	if err != nil {
		return nil, 0, err
	}

	return dictTypes, total, nil
}

// CreateDictionaryType 创建字典类型
func (r *DictionaryRepository) CreateDictionaryType(dictType *model.DictionaryType) error {
	return database.DB.Create(dictType).Error
}

// UpdateDictionaryType 更新字典类型
func (r *DictionaryRepository) UpdateDictionaryType(dictType *model.DictionaryType) error {
	return database.DB.Save(dictType).Error
}

// DeleteDictionaryType 删除字典类型
func (r *DictionaryRepository) DeleteDictionaryType(id int64) error {
	return database.DB.Delete(&model.DictionaryType{}, id).Error
}

// ========== 字典数据相关 ==========

// FindDictionaryByID 根据ID查询字典数据
func (r *DictionaryRepository) FindDictionaryByID(id int64) (*model.Dictionary, error) {
	var dict model.Dictionary
	err := database.DB.First(&dict, id).Error
	if err != nil {
		return nil, err
	}

	// 手动查询关联的字典类型
	if dict.TypeCode != "" {
		var dictType model.DictionaryType
		if err := database.DB.Where("type_code = ?", dict.TypeCode).First(&dictType).Error; err == nil {
			dict.DictionaryType = dictType
		}
	}

	return &dict, nil
}

// FindDictionariesByTypeCode 根据类型编码查询字典数据
func (r *DictionaryRepository) FindDictionariesByTypeCode(typeCode string) ([]model.Dictionary, error) {
	var dictionaries []model.Dictionary
	err := database.DB.Where("type_code = ? AND enabled = ?", typeCode, true).Order("order_num ASC, id ASC").Find(&dictionaries).Error
	return dictionaries, err
}

// FindAllDictionaries 查询所有字典数据
func (r *DictionaryRepository) FindAllDictionaries() ([]model.Dictionary, error) {
	var dictionaries []model.Dictionary
	err := database.DB.Order("type_code ASC, order_num ASC, id ASC").Find(&dictionaries).Error
	if err != nil {
		return dictionaries, err
	}

	// 手动查询关联的字典类型
	for i := range dictionaries {
		if dictionaries[i].TypeCode != "" {
			var dictType model.DictionaryType
			if err := database.DB.Where("type_code = ?", dictionaries[i].TypeCode).First(&dictType).Error; err == nil {
				dictionaries[i].DictionaryType = dictType
			}
		}
	}

	return dictionaries, err
}

// FindDictionariesWithPagination 分页查询字典数据
func (r *DictionaryRepository) FindDictionariesWithPagination(query map[string]interface{}, current, size int) ([]model.Dictionary, int64, error) {
	var dictionaries []model.Dictionary
	var total int64

	db := database.DB.Model(&model.Dictionary{})

	// 构建查询条件
	for key, value := range query {
		switch key {
		case "label":
			db = db.Where("label LIKE ?", "%"+value.(string)+"%")
		case "value":
			db = db.Where("value LIKE ?", "%"+value.(string)+"%")
		case "typeCode":
			db = db.Where("type_code LIKE ?", "%"+value.(string)+"%")
		case "enabled":
			db = db.Where("enabled = ?", value)
		default:
			db = db.Where(key+" = ?", value)
		}
	}

	// 查询总数
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (current - 1) * size
	err = db.Order("type_code ASC, order_num ASC, id ASC").Offset(offset).Limit(size).Find(&dictionaries).Error
	if err != nil {
		return nil, 0, err
	}

	// 手动查询关联的字典类型
	for i := range dictionaries {
		if dictionaries[i].TypeCode != "" {
			var dictType model.DictionaryType
			if err := database.DB.Where("type_code = ?", dictionaries[i].TypeCode).First(&dictType).Error; err == nil {
				dictionaries[i].DictionaryType = dictType
			}
		}
	}

	return dictionaries, total, nil
}

// CreateDictionary 创建字典数据
func (r *DictionaryRepository) CreateDictionary(dict *model.Dictionary) error {
	return database.DB.Create(dict).Error
}

// UpdateDictionary 更新字典数据
func (r *DictionaryRepository) UpdateDictionary(dict *model.Dictionary) error {
	return database.DB.Save(dict).Error
}

// DeleteDictionary 删除字典数据
func (r *DictionaryRepository) DeleteDictionary(id int64) error {
	return database.DB.Delete(&model.Dictionary{}, id).Error
}

// CountDictionariesByTypeCode 统计指定类型的字典数据数量
func (r *DictionaryRepository) CountDictionariesByTypeCode(typeCode string) (int64, error) {
	var count int64
	err := database.DB.Model(&model.Dictionary{}).Where("type_code = ?", typeCode).Count(&count).Error
	return count, err
}
