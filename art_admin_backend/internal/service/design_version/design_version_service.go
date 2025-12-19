package design_version

import (
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/dto/response"
	"art_admin_backend/internal/model"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// DesignVersionService 设计版本服务
// Design version service
type DesignVersionService struct {
	db *gorm.DB
}

// NewDesignVersionService 创建设计版本服务实例
// Create design version service instance
func NewDesignVersionService(db *gorm.DB) *DesignVersionService {
	return &DesignVersionService{db: db}
}

// Create 创建设计版本
// Create design version
func (s *DesignVersionService) Create(req *request.CreateDesignVersionRequest, userID uint) (*response.DesignVersionResponse, error) {
	// 验证项目是否存在 / Verify project exists
	var project model.DesignerProject
	if err := s.db.First(&project, req.ProjectID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("项目不存在")
		}
		return nil, err
	}

	// 获取下一个版本号 / Get next version number
	var maxVersion int
	s.db.Model(&model.DesignVersion{}).
		Where("project_id = ?", req.ProjectID).
		Select("COALESCE(MAX(version_number), 0)").
		Scan(&maxVersion)
	nextVersion := maxVersion + 1

	// 序列化设计图列表 / Serialize design images
	var designImagesJSON *string
	if len(req.DesignImages) > 0 {
		images := make([]model.DesignImage, len(req.DesignImages))
		for i, url := range req.DesignImages {
			images[i] = model.DesignImage{
				URL:  url,
				Name: fmt.Sprintf("设计图%d", i+1),
				Type: "design",
			}
		}
		if data, err := json.Marshal(images); err == nil {
			str := string(data)
			designImagesJSON = &str
		}
	}

	// 序列化CAD文件ID列表 / Serialize CAD file IDs
	var cadFileIDsJSON *string
	if len(req.CadFileIDs) > 0 {
		if data, err := json.Marshal(req.CadFileIDs); err == nil {
			str := string(data)
			cadFileIDsJSON = &str
		}
	}

	// 生成版本名称 / Generate version name
	versionName := req.VersionName
	if versionName == "" {
		versionName = fmt.Sprintf("版本%d", nextVersion)
	}

	// 创建版本 / Create version
	version := &model.DesignVersion{
		ProjectID:     req.ProjectID,
		VersionNumber: nextVersion,
		VersionName:   versionName,
		Description:   req.Description,
		DesignImages:  designImagesJSON,
		CadFileIDs:    cadFileIDsJSON,
		LayoutInfo:    strPtr(req.LayoutInfo),
		AreaInfo:      strPtr(req.AreaInfo),
		StyleInfo:     strPtr(req.StyleInfo),
		MaterialInfo:  strPtr(req.MaterialInfo),
		Status:        "draft",
		CreatedBy:     userID,
	}

	if err := s.db.Create(version).Error; err != nil {
		return nil, err
	}

	return s.toResponse(version), nil
}

// GetByID 根据ID获取版本详情
// Get version detail by ID
func (s *DesignVersionService) GetByID(id uint) (*response.DesignVersionResponse, error) {
	var version model.DesignVersion
	if err := s.db.First(&version, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("版本不存在")
		}
		return nil, err
	}
	return s.toResponse(&version), nil
}

// List 获取版本列表
// Get version list
func (s *DesignVersionService) List(req *request.DesignVersionListRequest) (*response.DesignVersionListResponse, error) {
	var versions []model.DesignVersion
	var total int64

	query := s.db.Model(&model.DesignVersion{}).Where("project_id = ?", req.ProjectID)

	// 状态筛选 / Status filter
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	// 统计总数 / Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询 / Paginated query
	offset := (req.Current - 1) * req.Size
	if err := query.Order("version_number DESC").
		Offset(offset).Limit(req.Size).
		Find(&versions).Error; err != nil {
		return nil, err
	}

	// 转换响应 / Convert to response
	records := make([]response.DesignVersionResponse, len(versions))
	for i, v := range versions {
		records[i] = *s.toResponse(&v)
	}

	return &response.DesignVersionListResponse{
		Records: records,
		Current: req.Current,
		Size:    req.Size,
		Total:   total,
	}, nil
}

// Update 更新版本
// Update version
func (s *DesignVersionService) Update(req *request.UpdateDesignVersionRequest) (*response.DesignVersionResponse, error) {
	var version model.DesignVersion
	if err := s.db.First(&version, req.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("版本不存在")
		}
		return nil, err
	}

	// 更新字段 / Update fields
	if req.VersionName != "" {
		version.VersionName = req.VersionName
	}
	if req.Description != "" {
		version.Description = req.Description
	}
	if len(req.DesignImages) > 0 {
		images := make([]model.DesignImage, len(req.DesignImages))
		for i, url := range req.DesignImages {
			images[i] = model.DesignImage{
				URL:  url,
				Name: fmt.Sprintf("设计图%d", i+1),
				Type: "design",
			}
		}
		if data, err := json.Marshal(images); err == nil {
			str := string(data)
			version.DesignImages = &str
		}
	}
	if len(req.CadFileIDs) > 0 {
		if data, err := json.Marshal(req.CadFileIDs); err == nil {
			str := string(data)
			version.CadFileIDs = &str
		}
	}
	if req.LayoutInfo != "" {
		version.LayoutInfo = strPtr(req.LayoutInfo)
	}
	if req.AreaInfo != "" {
		version.AreaInfo = strPtr(req.AreaInfo)
	}
	if req.StyleInfo != "" {
		version.StyleInfo = strPtr(req.StyleInfo)
	}
	if req.MaterialInfo != "" {
		version.MaterialInfo = strPtr(req.MaterialInfo)
	}
	if req.Status != "" {
		version.Status = req.Status
	}

	if err := s.db.Save(&version).Error; err != nil {
		return nil, err
	}

	return s.toResponse(&version), nil
}

// Delete 删除版本
// Delete version
func (s *DesignVersionService) Delete(id uint) error {
	result := s.db.Delete(&model.DesignVersion{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("版本不存在")
	}
	return nil
}

// Compare 对比两个版本
// Compare two versions
func (s *DesignVersionService) Compare(req *request.CompareDesignVersionsRequest) (*response.VersionCompareResponse, error) {
	// 获取两个版本 / Get both versions
	var versionA, versionB model.DesignVersion
	if err := s.db.First(&versionA, req.VersionAID).Error; err != nil {
		return nil, errors.New("版本A不存在")
	}
	if err := s.db.First(&versionB, req.VersionBID).Error; err != nil {
		return nil, errors.New("版本B不存在")
	}

	// 验证属于同一项目 / Verify same project
	if versionA.ProjectID != versionB.ProjectID || versionA.ProjectID != req.ProjectID {
		return nil, errors.New("版本必须属于同一项目")
	}

	// 计算差异 / Calculate differences
	layoutChanges := s.compareJSON(ptrToStr(versionA.LayoutInfo), ptrToStr(versionB.LayoutInfo), "布局")
	areaChanges := s.compareJSON(ptrToStr(versionA.AreaInfo), ptrToStr(versionB.AreaInfo), "面积")
	styleChanges := s.compareJSON(ptrToStr(versionA.StyleInfo), ptrToStr(versionB.StyleInfo), "风格")
	materialChanges := s.compareJSON(ptrToStr(versionA.MaterialInfo), ptrToStr(versionB.MaterialInfo), "材料")

	// 生成摘要 / Generate summary
	totalChanges := len(layoutChanges) + len(areaChanges) + len(styleChanges) + len(materialChanges)
	summary := fmt.Sprintf("版本%d与版本%d对比：共发现%d处变化", versionA.VersionNumber, versionB.VersionNumber, totalChanges)

	// 保存对比记录 / Save compare record
	compare := &model.DesignVersionCompare{
		ProjectID:       req.ProjectID,
		VersionAID:      req.VersionAID,
		VersionBID:      req.VersionBID,
		LayoutChanges:   toJSON(layoutChanges),
		AreaChanges:     toJSON(areaChanges),
		StyleChanges:    toJSON(styleChanges),
		MaterialChanges: toJSON(materialChanges),
		Summary:         summary,
		CompareStatus:   "completed",
	}
	s.db.Create(compare)

	return &response.VersionCompareResponse{
		ID:         compare.ID,
		ProjectID:  req.ProjectID,
		VersionAID: req.VersionAID,
		VersionBID: req.VersionBID,
		VersionA: &response.DesignVersionBriefResponse{
			ID:            versionA.ID,
			VersionNumber: versionA.VersionNumber,
			VersionName:   versionA.VersionName,
			Status:        versionA.Status,
		},
		VersionB: &response.DesignVersionBriefResponse{
			ID:            versionB.ID,
			VersionNumber: versionB.VersionNumber,
			VersionName:   versionB.VersionName,
			Status:        versionB.Status,
		},
		LayoutChanges:   layoutChanges,
		AreaChanges:     areaChanges,
		StyleChanges:    styleChanges,
		MaterialChanges: materialChanges,
		Summary:         summary,
		CompareStatus:   "completed",
		CreatedAt:       compare.CreatedAt,
		UpdatedAt:       compare.UpdatedAt,
	}, nil
}

// GetDiff 获取版本差异（用于左右分屏展示）
// Get version diff (for side-by-side display)
func (s *DesignVersionService) GetDiff(versionAID, versionBID uint) (*response.VersionDiffResponse, error) {
	versionA, err := s.GetByID(versionAID)
	if err != nil {
		return nil, errors.New("版本A不存在")
	}
	versionB, err := s.GetByID(versionBID)
	if err != nil {
		return nil, errors.New("版本B不存在")
	}

	// 获取原始数据计算差异 / Get raw data for diff calculation
	var modelA, modelB model.DesignVersion
	s.db.First(&modelA, versionAID)
	s.db.First(&modelB, versionBID)

	layoutChanges := s.compareJSON(ptrToStr(modelA.LayoutInfo), ptrToStr(modelB.LayoutInfo), "布局")
	areaChanges := s.compareJSON(ptrToStr(modelA.AreaInfo), ptrToStr(modelB.AreaInfo), "面积")
	styleChanges := s.compareJSON(ptrToStr(modelA.StyleInfo), ptrToStr(modelB.StyleInfo), "风格")
	materialChanges := s.compareJSON(ptrToStr(modelA.MaterialInfo), ptrToStr(modelB.MaterialInfo), "材料")

	totalChanges := len(layoutChanges) + len(areaChanges) + len(styleChanges) + len(materialChanges)
	summary := fmt.Sprintf("版本%d与版本%d对比：共发现%d处变化", modelA.VersionNumber, modelB.VersionNumber, totalChanges)

	return &response.VersionDiffResponse{
		VersionA:        versionA,
		VersionB:        versionB,
		LayoutChanges:   layoutChanges,
		AreaChanges:     areaChanges,
		StyleChanges:    styleChanges,
		MaterialChanges: materialChanges,
		Summary:         summary,
	}, nil
}

// GetCompareList 获取对比记录列表
// Get compare record list
func (s *DesignVersionService) GetCompareList(req *request.VersionCompareListRequest) (*response.VersionCompareListResponse, error) {
	var compares []model.DesignVersionCompare
	var total int64

	query := s.db.Model(&model.DesignVersionCompare{}).Where("project_id = ?", req.ProjectID)

	if req.CompareStatus != "" {
		query = query.Where("compare_status = ?", req.CompareStatus)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (req.Current - 1) * req.Size
	if err := query.Preload("VersionA").Preload("VersionB").
		Order("created_at DESC").
		Offset(offset).Limit(req.Size).
		Find(&compares).Error; err != nil {
		return nil, err
	}

	records := make([]response.VersionCompareResponse, len(compares))
	for i, c := range compares {
		records[i] = *s.toCompareResponse(&c)
	}

	return &response.VersionCompareListResponse{
		Records: records,
		Current: req.Current,
		Size:    req.Size,
		Total:   total,
	}, nil
}

// GetCompareDetail 获取对比详情
// Get compare detail
func (s *DesignVersionService) GetCompareDetail(id uint) (*response.VersionCompareResponse, error) {
	var compare model.DesignVersionCompare
	if err := s.db.Preload("VersionA").Preload("VersionB").First(&compare, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("对比记录不存在")
		}
		return nil, err
	}
	return s.toCompareResponse(&compare), nil
}

// DeleteCompare 删除对比记录
// Delete compare record
func (s *DesignVersionService) DeleteCompare(id uint) error {
	result := s.db.Delete(&model.DesignVersionCompare{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("对比记录不存在")
	}
	return nil
}

// toResponse 转换为响应结构
// Convert to response structure
func (s *DesignVersionService) toResponse(v *model.DesignVersion) *response.DesignVersionResponse {
	resp := &response.DesignVersionResponse{
		ID:            v.ID,
		ProjectID:     v.ProjectID,
		VersionNumber: v.VersionNumber,
		VersionName:   v.VersionName,
		Description:   v.Description,
		Status:        v.Status,
		CreatedBy:     v.CreatedBy,
		CreatedAt:     v.CreatedAt,
		UpdatedAt:     v.UpdatedAt,
	}

	// 解析设计图 / Parse design images
	designImagesStr := ptrToStr(v.DesignImages)
	if designImagesStr != "" && designImagesStr != "[]" {
		var images []response.DesignImageResponse
		json.Unmarshal([]byte(designImagesStr), &images)
		resp.DesignImages = images
	}

	// 解析CAD文件ID / Parse CAD file IDs
	cadFileIDsStr := ptrToStr(v.CadFileIDs)
	if cadFileIDsStr != "" && cadFileIDsStr != "[]" {
		var cadFileIDs []uint
		json.Unmarshal([]byte(cadFileIDsStr), &cadFileIDs)
		resp.CadFileIDs = cadFileIDs

		// 获取CAD文件信息 / Get CAD file info
		if len(cadFileIDs) > 0 {
			var cadFiles []model.CadFile
			s.db.Where("id IN ?", cadFileIDs).Find(&cadFiles)
			resp.CadFiles = make([]response.CadFileBriefResponse, len(cadFiles))
			for i, f := range cadFiles {
				resp.CadFiles[i] = response.CadFileBriefResponse{
					ID:            f.ID,
					FileName:      f.FileName,
					FileFormat:    f.FileFormat,
					FilePath:      f.FilePath,
					IsAIGenerated: f.IsAIGenerated,
				}
			}
		}
	}

	// 解析其他JSON字段 / Parse other JSON fields
	layoutInfoStr := ptrToStr(v.LayoutInfo)
	if layoutInfoStr != "" {
		var layout []response.LayoutInfoResponse
		json.Unmarshal([]byte(layoutInfoStr), &layout)
		resp.LayoutInfo = layout
	}
	areaInfoStr := ptrToStr(v.AreaInfo)
	if areaInfoStr != "" {
		var area []response.AreaInfoResponse
		json.Unmarshal([]byte(areaInfoStr), &area)
		resp.AreaInfo = area
	}
	styleInfoStr := ptrToStr(v.StyleInfo)
	if styleInfoStr != "" {
		var style []response.StyleInfoResponse
		json.Unmarshal([]byte(styleInfoStr), &style)
		resp.StyleInfo = style
	}
	materialInfoStr := ptrToStr(v.MaterialInfo)
	if materialInfoStr != "" {
		var material []response.MaterialInfoResponse
		json.Unmarshal([]byte(materialInfoStr), &material)
		resp.MaterialInfo = material
	}

	return resp
}

// toCompareResponse 转换对比记录为响应
// Convert compare record to response
func (s *DesignVersionService) toCompareResponse(c *model.DesignVersionCompare) *response.VersionCompareResponse {
	resp := &response.VersionCompareResponse{
		ID:            c.ID,
		ProjectID:     c.ProjectID,
		VersionAID:    c.VersionAID,
		VersionBID:    c.VersionBID,
		Summary:       c.Summary,
		CompareStatus: c.CompareStatus,
		CreatedAt:     c.CreatedAt,
		UpdatedAt:     c.UpdatedAt,
	}

	if c.VersionA != nil {
		resp.VersionA = &response.DesignVersionBriefResponse{
			ID:            c.VersionA.ID,
			VersionNumber: c.VersionA.VersionNumber,
			VersionName:   c.VersionA.VersionName,
			Status:        c.VersionA.Status,
		}
	}
	if c.VersionB != nil {
		resp.VersionB = &response.DesignVersionBriefResponse{
			ID:            c.VersionB.ID,
			VersionNumber: c.VersionB.VersionNumber,
			VersionName:   c.VersionB.VersionName,
			Status:        c.VersionB.Status,
		}
	}

	// 解析变化列表 / Parse change lists
	json.Unmarshal([]byte(c.LayoutChanges), &resp.LayoutChanges)
	json.Unmarshal([]byte(c.AreaChanges), &resp.AreaChanges)
	json.Unmarshal([]byte(c.ElementChanges), &resp.ElementChanges)
	json.Unmarshal([]byte(c.StyleChanges), &resp.StyleChanges)
	json.Unmarshal([]byte(c.MaterialChanges), &resp.MaterialChanges)

	return resp
}

// compareJSON 比较两个JSON字符串的差异
// Compare two JSON strings for differences
func (s *DesignVersionService) compareJSON(jsonA, jsonB, category string) []response.ChangeItemResponse {
	var changes []response.ChangeItemResponse

	// 简单比较：如果不同则记录变化
	if jsonA != jsonB {
		changes = append(changes, response.ChangeItemResponse{
			Field:       category,
			OldValue:    jsonA,
			NewValue:    jsonB,
			ChangeType:  "modified",
			Description: fmt.Sprintf("%s信息已变更", category),
		})
	}

	return changes
}

// toJSON 将对象转换为JSON字符串
// Convert object to JSON string
func toJSON(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}

// strPtr 将字符串转换为指针（空字符串返回nil）
// Convert string to pointer (empty string returns nil)
func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// ptrToStr 将字符串指针转换为字符串（nil返回空字符串）
// Convert string pointer to string (nil returns empty string)
func ptrToStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
