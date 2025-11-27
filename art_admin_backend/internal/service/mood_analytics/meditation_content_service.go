package mood_analytics

import (
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/dto/response"
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/utils"
	"art_admin_backend/internal/repository"
	"errors"
	"fmt"
	"time"
)

// MeditationContentService 冥想内容服务
type MeditationContentService struct {
	meditationContentRepo *repository.MeditationContentRepository
}

// NewMeditationContentService 创建冥想内容服务
func NewMeditationContentService() *MeditationContentService {
	return &MeditationContentService{
		meditationContentRepo: repository.NewMeditationContentRepository(),
	}
}

// CreateMeditationContent 创建冥想内容
func (s *MeditationContentService) CreateMeditationContent(userID int64, req *request.CreateMeditationContentRequest) error {
	// 验证分类
	if !utils.ValidateMeditationCategory(req.Category) {
		return errors.New("无效的冥想分类")
	}

	// 验证时长
	if req.Duration < 60 {
		return errors.New("冥想时长不能少于60秒")
	}

	// 验证难度等级
	if req.DifficultyLevel < 1 || req.DifficultyLevel > 5 {
		return errors.New("难度等级必须在1-5之间")
	}

	// 创建模型
	content := &model.MeditationContent{
		Title:           req.Title,
		Category:        req.Category,
		Duration:        req.Duration,
		AudioURL:        req.AudioURL,
		CoverImage:      req.CoverImage,
		Description:     req.Description,
		DifficultyLevel: req.DifficultyLevel,
		IsPublic:        req.IsPublic,
		ViewCount:       0,
		LikeCount:       0,
		CreatedBy:       userID,
	}

	// 处理标签
	if len(req.Tags) > 0 {
		tagsJSON, err := utils.JSONMarshal(req.Tags)
		if err != nil {
			return errors.New("标签序列化失败")
		}
		content.Tags = tagsJSON
	} else {
		// 设置空JSON数组而不是nil值
		content.Tags = "[]"
	}

	// 保存到数据库
	if err := s.meditationContentRepo.Create(content); err != nil {
		fmt.Printf("创建冥想内容数据库错误: %v\n", err)
		return errors.New("创建冥想内容失败")
	}

	return nil
}

// UpdateMeditationContent 更新冥想内容
func (s *MeditationContentService) UpdateMeditationContent(userID int64, req *request.UpdateMeditationContentRequest) error {
	// 获取现有内容
	existingContent, err := s.meditationContentRepo.FindByID(req.ID)
	if err != nil {
		return errors.New("冥想内容不存在")
	}

	// 验证权限（只有创建者可以修改）
	if existingContent.CreatedBy != userID {
		return errors.New("无权限操作此内容")
	}

	// 验证分类（如果提供）
	if req.Category != nil && !utils.ValidateMeditationCategory(*req.Category) {
		return errors.New("无效的冥想分类")
	}

	// 验证时长（如果提供）
	if req.Duration != nil && *req.Duration < 60 {
		return errors.New("冥想时长不能少于60秒")
	}

	// 验证难度等级（如果提供）
	if req.DifficultyLevel != nil && (*req.DifficultyLevel < 1 || *req.DifficultyLevel > 5) {
		return errors.New("难度等级必须在1-5之间")
	}

	// 更新字段
	if req.Title != nil {
		existingContent.Title = *req.Title
	}
	if req.Category != nil {
		existingContent.Category = *req.Category
	}
	if req.Duration != nil {
		existingContent.Duration = *req.Duration
	}
	if req.AudioURL != nil {
		existingContent.AudioURL = *req.AudioURL
	}
	if req.CoverImage != nil {
		existingContent.CoverImage = *req.CoverImage
	}
	if req.Description != nil {
		existingContent.Description = *req.Description
	}
	if req.DifficultyLevel != nil {
		existingContent.DifficultyLevel = *req.DifficultyLevel
	}
	if req.IsPublic != nil {
		existingContent.IsPublic = *req.IsPublic
	}
	if req.Tags != nil {
		tagsJSON, err := utils.JSONMarshal(req.Tags)
		if err != nil {
			return errors.New("标签序列化失败")
		}
		existingContent.Tags = tagsJSON
	}

	existingContent.UpdatedAt = time.Now()

	// 更新数据库
	return s.meditationContentRepo.Update(existingContent)
}

// DeleteMeditationContent 删除冥想内容
func (s *MeditationContentService) DeleteMeditationContent(userID, id int64) error {
	// 获取内容
	content, err := s.meditationContentRepo.FindByID(id)
	if err != nil {
		return errors.New("冥想内容不存在")
	}

	// 验证权限（只有创建者可以删除）
	if content.CreatedBy != userID {
		return errors.New("无权限操作此内容")
	}

	// 删除内容
	return s.meditationContentRepo.Delete(id)
}

// GetMeditationContentList 获取冥想内容列表
func (s *MeditationContentService) GetMeditationContentList(req *request.MeditationContentListRequest) ([]response.MeditationContentItem, int64, error) {
	// 构建查询条件
	query := make(map[string]interface{})
	if req.Category != "" {
		query["category"] = req.Category
	}
	if req.DifficultyLevel != nil {
		query["difficulty_level"] = *req.DifficultyLevel
	}
	if req.IsPublic != nil {
		query["is_public"] = *req.IsPublic
	}
	if req.MinDuration != nil {
		query["min_duration"] = *req.MinDuration
	}
	if req.MaxDuration != nil {
		query["max_duration"] = *req.MaxDuration
	}

	// 查询内容
	contents, total, err := s.meditationContentRepo.FindWithPagination(query, req.Page, req.PageSize)
	if err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	items := make([]response.MeditationContentItem, len(contents))
	for i, content := range contents {
		items[i] = s.convertToMeditationContentItem(&content)
	}

	return items, total, nil
}

// GetMeditationContentDetail 获取冥想内容详情
func (s *MeditationContentService) GetMeditationContentDetail(id int64) (*response.MeditationContentDetail, error) {
	// 获取内容
	content, err := s.meditationContentRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 增加查看次数
	if err := s.meditationContentRepo.IncrementViewCount(id); err != nil {
		// 记录错误但不影响返回
	}

	// 转换为响应格式
	detail := &response.MeditationContentDetail{
		MeditationContentItem: s.convertToMeditationContentItem(content),
		Description:           content.Description,
		CreatedBy:             content.CreatedBy,
		CreatedAt:             content.CreatedAt,
		UpdatedAt:             content.UpdatedAt,
	}

	return detail, nil
}

// GetMeditationContentByCategory 根据分类获取冥想内容
func (s *MeditationContentService) GetMeditationContentByCategory(category string, limit int) ([]response.MeditationContentItem, error) {
	contents, err := s.meditationContentRepo.FindByCategory(category, limit)
	if err != nil {
		return nil, err
	}

	items := make([]response.MeditationContentItem, len(contents))
	for i, content := range contents {
		items[i] = s.convertToMeditationContentItem(&content)
	}

	return items, nil
}

// GetMeditationContentByDifficulty 根据难度等级获取冥想内容
func (s *MeditationContentService) GetMeditationContentByDifficulty(difficulty int, limit int) ([]response.MeditationContentItem, error) {
	contents, err := s.meditationContentRepo.FindByDifficulty(difficulty, limit)
	if err != nil {
		return nil, err
	}

	items := make([]response.MeditationContentItem, len(contents))
	for i, content := range contents {
		items[i] = s.convertToMeditationContentItem(&content)
	}

	return items, nil
}

// GetPopularMeditationContent 获取热门冥想内容
func (s *MeditationContentService) GetPopularMeditationContent(limit int) ([]response.MeditationContentItem, error) {
	contents, err := s.meditationContentRepo.FindPopular(limit)
	if err != nil {
		return nil, err
	}

	items := make([]response.MeditationContentItem, len(contents))
	for i, content := range contents {
		items[i] = s.convertToMeditationContentItem(&content)
	}

	return items, nil
}

// LikeMeditationContent 点赞冥想内容
func (s *MeditationContentService) LikeMeditationContent(id int64) error {
	return s.meditationContentRepo.IncrementLikeCount(id)
}

// UnlikeMeditationContent 取消点赞冥想内容
func (s *MeditationContentService) UnlikeMeditationContent(id int64) error {
	return s.meditationContentRepo.DecrementLikeCount(id)
}

// SearchMeditationContent 搜索冥想内容
func (s *MeditationContentService) SearchMeditationContent(keyword string, page, pageSize int) ([]response.MeditationContentItem, int64, error) {
	contents, total, err := s.meditationContentRepo.Search(keyword, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	items := make([]response.MeditationContentItem, len(contents))
	for i, content := range contents {
		items[i] = s.convertToMeditationContentItem(&content)
	}

	return items, total, nil
}

// GetCategories 获取所有分类
func (s *MeditationContentService) GetCategories() ([]string, error) {
	return s.meditationContentRepo.GetCategories()
}

// GetDifficultyLevels 获取所有难度等级
func (s *MeditationContentService) GetDifficultyLevels() ([]int, error) {
	return s.meditationContentRepo.GetDifficultyLevels()
}

// convertToMeditationContentItem 转换模型为响应项
func (s *MeditationContentService) convertToMeditationContentItem(content *model.MeditationContent) response.MeditationContentItem {
	item := response.MeditationContentItem{
		ID:              content.ID,
		Title:           content.Title,
		Category:        content.Category,
		CategoryLabel:   model.GetCategoryLabel(model.MeditationCategory(content.Category)),
		Duration:        content.Duration,
		AudioURL:        content.AudioURL,
		CoverImage:      content.CoverImage,
		DifficultyLevel: content.DifficultyLevel,
		IsPublic:        content.IsPublic,
		ViewCount:       int64(content.ViewCount),
		LikeCount:       int64(content.LikeCount),
	}

	// 解析标签
	if content.Tags != "" {
		tags, _ := utils.JSONUnmarshal(content.Tags)
		item.Tags = tags
	}

	return item
}
