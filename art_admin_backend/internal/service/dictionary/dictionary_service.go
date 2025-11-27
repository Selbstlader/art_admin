package dictionary

import (
	"errors"
	"time"

	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/dto/response"
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/repository"
)

type DictionaryService struct {
	dictRepo *repository.DictionaryRepository
}

func NewDictionaryService() *DictionaryService {
	return &DictionaryService{
		dictRepo: repository.NewDictionaryRepository(),
	}
}

// ========== 字典类型相关 ==========

// GetDictionaryTypeList 获取字典类型列表
func (s *DictionaryService) GetDictionaryTypeList(req *request.DictionaryTypeListRequest) ([]response.DictionaryTypeResponse, int64, error) {
	// 构建查询条件
	query := make(map[string]interface{})
	if req.ID != nil {
		query["id"] = *req.ID
	}
	if req.TypeName != "" {
		query["typeName"] = req.TypeName
	}
	if req.TypeCode != "" {
		query["typeCode"] = req.TypeCode
	}
	if req.Description != "" {
		query["description"] = req.Description
	}
	if req.Enabled != nil {
		query["enabled"] = *req.Enabled
	}

	// 查询字典类型列表
	dictTypes, total, err := s.dictRepo.FindDictionaryTypesWithPagination(query, req.Current, req.Size)
	if err != nil {
		return nil, 0, err
	}

	// 转换为响应DTO
	result := make([]response.DictionaryTypeResponse, 0, len(dictTypes))
	for _, dictType := range dictTypes {
		result = append(result, response.DictionaryTypeResponse{
			ID:          dictType.ID,
			TypeName:    dictType.TypeName,
			TypeCode:    dictType.TypeCode,
			Description: dictType.Description,
			Enabled:     dictType.Enabled,
			Remark:      dictType.Remark,
			CreateBy:    dictType.CreateBy,
			CreateTime:  dictType.CreateTime,
			UpdateBy:    dictType.UpdateBy,
			UpdateTime:  dictType.UpdateTime,
		})
	}

	return result, total, nil
}

// GetDictionaryTypeByID 根据ID获取字典类型
func (s *DictionaryService) GetDictionaryTypeByID(id int64) (*response.DictionaryTypeResponse, error) {
	dictType, err := s.dictRepo.FindDictionaryTypeByID(id)
	if err != nil {
		return nil, err
	}

	return &response.DictionaryTypeResponse{
		ID:          dictType.ID,
		TypeName:    dictType.TypeName,
		TypeCode:    dictType.TypeCode,
		Description: dictType.Description,
		Enabled:     dictType.Enabled,
		Remark:      dictType.Remark,
		CreateBy:    dictType.CreateBy,
		CreateTime:  dictType.CreateTime,
		UpdateBy:    dictType.UpdateBy,
		UpdateTime:  dictType.UpdateTime,
	}, nil
}

// CreateDictionaryType 创建字典类型
func (s *DictionaryService) CreateDictionaryType(req *request.CreateDictionaryTypeRequest, createBy string) error {
	// 检查类型编码是否已存在
	existingDict, _ := s.dictRepo.FindDictionaryTypeByCode(req.TypeCode)
	if existingDict != nil {
		return errors.New("字典类型编码已存在")
	}

	// 创建字典类型
	dictType := &model.DictionaryType{
		TypeName:    req.TypeName,
		TypeCode:    req.TypeCode,
		Description: req.Description,
		Enabled:     req.Enabled,
		Remark:      req.Remark,
		CreateBy:    createBy,
		CreateTime:  time.Now(),
	}

	return s.dictRepo.CreateDictionaryType(dictType)
}

// UpdateDictionaryType 更新字典类型
func (s *DictionaryService) UpdateDictionaryType(req *request.UpdateDictionaryTypeRequest, updateBy string) error {
	// 检查字典类型是否存在
	dictType, err := s.dictRepo.FindDictionaryTypeByID(req.ID)
	if err != nil {
		return errors.New("字典类型不存在")
	}

	// 检查类型编码是否与其他记录冲突
	existingDict, _ := s.dictRepo.FindDictionaryTypeByCode(req.TypeCode)
	if existingDict != nil && existingDict.ID != req.ID {
		return errors.New("字典类型编码已存在")
	}

	// 更新字典类型
	dictType.TypeName = req.TypeName
	dictType.TypeCode = req.TypeCode
	dictType.Description = req.Description
	dictType.Enabled = req.Enabled
	dictType.Remark = req.Remark
	dictType.UpdateBy = updateBy
	dictType.UpdateTime = time.Now()

	return s.dictRepo.UpdateDictionaryType(dictType)
}

// DeleteDictionaryType 删除字典类型
func (s *DictionaryService) DeleteDictionaryType(id int64) error {
	// 检查字典类型是否存在
	dictType, err := s.dictRepo.FindDictionaryTypeByID(id)
	if err != nil {
		return errors.New("字典类型不存在")
	}

	// 检查是否有字典数据在使用该类型
	count, err := s.dictRepo.CountDictionariesByTypeCode(dictType.TypeCode)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("该字典类型下存在字典数据，无法删除")
	}

	return s.dictRepo.DeleteDictionaryType(id)
}

// ========== 字典数据相关 ==========

// GetDictionaryList 获取字典数据列表
func (s *DictionaryService) GetDictionaryList(req *request.DictionaryListRequest) ([]response.DictionaryWithTypeResponse, int64, error) {
	// 构建查询条件
	query := make(map[string]interface{})
	if req.ID != nil {
		query["id"] = *req.ID
	}
	if req.TypeCode != "" {
		query["typeCode"] = req.TypeCode
	}
	if req.Label != "" {
		query["label"] = req.Label
	}
	if req.Value != "" {
		query["value"] = req.Value
	}
	if req.Enabled != nil {
		query["enabled"] = *req.Enabled
	}

	// 查询字典数据列表
	dictionaries, total, err := s.dictRepo.FindDictionariesWithPagination(query, req.Current, req.Size)
	if err != nil {
		return nil, 0, err
	}

	// 转换为响应DTO
	result := make([]response.DictionaryWithTypeResponse, 0, len(dictionaries))
	for _, dict := range dictionaries {
		dictTypeResp := response.DictionaryTypeResponse{
			ID:          dict.DictionaryType.ID,
			TypeName:    dict.DictionaryType.TypeName,
			TypeCode:    dict.DictionaryType.TypeCode,
			Description: dict.DictionaryType.Description,
			Enabled:     dict.DictionaryType.Enabled,
			Remark:      dict.DictionaryType.Remark,
			CreateBy:    dict.DictionaryType.CreateBy,
			CreateTime:  dict.DictionaryType.CreateTime,
			UpdateBy:    dict.DictionaryType.UpdateBy,
			UpdateTime:  dict.DictionaryType.UpdateTime,
		}

		result = append(result, response.DictionaryWithTypeResponse{
			ID:             dict.ID,
			TypeCode:       dict.TypeCode,
			Label:          dict.Label,
			Value:          dict.Value,
			OrderNum:       dict.OrderNum,
			Enabled:        dict.Enabled,
			Remark:         dict.Remark,
			CreateBy:       dict.CreateBy,
			CreateTime:     dict.CreateTime,
			UpdateBy:       dict.UpdateBy,
			UpdateTime:     dict.UpdateTime,
			DictionaryType: dictTypeResp,
		})
	}

	return result, total, nil
}

// GetDictionaryByID 根据ID获取字典数据
func (s *DictionaryService) GetDictionaryByID(id int64) (*response.DictionaryResponse, error) {
	dict, err := s.dictRepo.FindDictionaryByID(id)
	if err != nil {
		return nil, err
	}

	return &response.DictionaryResponse{
		ID:         dict.ID,
		TypeCode:   dict.TypeCode,
		Label:      dict.Label,
		Value:      dict.Value,
		OrderNum:   dict.OrderNum,
		Enabled:    dict.Enabled,
		Remark:     dict.Remark,
		CreateBy:   dict.CreateBy,
		CreateTime: dict.CreateTime,
		UpdateBy:   dict.UpdateBy,
		UpdateTime: dict.UpdateTime,
	}, nil
}

// GetDictionaryByTypeCode 根据类型编码获取字典数据
func (s *DictionaryService) GetDictionaryByTypeCode(typeCode string) ([]response.DictionaryResponse, error) {
	dictionaries, err := s.dictRepo.FindDictionariesByTypeCode(typeCode)
	if err != nil {
		return nil, err
	}

	// 转换为响应DTO
	result := make([]response.DictionaryResponse, 0, len(dictionaries))
	for _, dict := range dictionaries {
		result = append(result, response.DictionaryResponse{
			ID:         dict.ID,
			TypeCode:   dict.TypeCode,
			Label:      dict.Label,
			Value:      dict.Value,
			OrderNum:   dict.OrderNum,
			Enabled:    dict.Enabled,
			Remark:     dict.Remark,
			CreateBy:   dict.CreateBy,
			CreateTime: dict.CreateTime,
			UpdateBy:   dict.UpdateBy,
			UpdateTime: dict.UpdateTime,
		})
	}

	return result, nil
}

// CreateDictionary 创建字典数据
func (s *DictionaryService) CreateDictionary(req *request.CreateDictionaryRequest, createBy string) error {
	// 检查字典类型是否存在
	_, err := s.dictRepo.FindDictionaryTypeByCode(req.TypeCode)
	if err != nil {
		return errors.New("字典类型不存在")
	}

	// 创建字典数据
	dict := &model.Dictionary{
		TypeCode:   req.TypeCode,
		Label:      req.Label,
		Value:      req.Value,
		OrderNum:   req.OrderNum,
		Enabled:    req.Enabled,
		Remark:     req.Remark,
		CreateBy:   createBy,
		CreateTime: time.Now(),
	}

	return s.dictRepo.CreateDictionary(dict)
}

// UpdateDictionary 更新字典数据
func (s *DictionaryService) UpdateDictionary(req *request.UpdateDictionaryRequest, updateBy string) error {
	// 检查字典数据是否存在
	dict, err := s.dictRepo.FindDictionaryByID(req.ID)
	if err != nil {
		return errors.New("字典数据不存在")
	}

	// 检查字典类型是否存在
	_, err = s.dictRepo.FindDictionaryTypeByCode(req.TypeCode)
	if err != nil {
		return errors.New("字典类型不存在")
	}

	// 更新字典数据
	dict.TypeCode = req.TypeCode
	dict.Label = req.Label
	dict.Value = req.Value
	dict.OrderNum = req.OrderNum
	dict.Enabled = req.Enabled
	dict.Remark = req.Remark
	dict.UpdateBy = updateBy
	dict.UpdateTime = time.Now()

	return s.dictRepo.UpdateDictionary(dict)
}

// DeleteDictionary 删除字典数据
func (s *DictionaryService) DeleteDictionary(id int64) error {
	// 检查字典数据是否存在
	_, err := s.dictRepo.FindDictionaryByID(id)
	if err != nil {
		return errors.New("字典数据不存在")
	}

	return s.dictRepo.DeleteDictionary(id)
}
