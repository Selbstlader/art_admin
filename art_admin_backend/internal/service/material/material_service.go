package material

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/dto/response"
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/volcengine"
	"art_admin_backend/internal/repository"
)

// MaterialService 材料服务接口
// Material service interface
type MaterialService interface {
	// CRUD operations
	CreateMaterial(req *request.CreateMaterialRequest) (*response.MaterialResponse, error)
	UpdateMaterial(req *request.UpdateMaterialRequest) (*response.MaterialResponse, error)
	DeleteMaterial(id uint) error
	BatchDeleteMaterials(ids []uint) error
	GetMaterialByID(id uint) (*response.MaterialDetailResponse, error)
	GetMaterialList(req *request.MaterialListRequest) (*response.MaterialListResponse, error)

	// Query operations
	GetCategories() ([]response.MaterialCategoryResponse, error)
	GetBrands() ([]response.MaterialBrandResponse, error)
	GetStats() (*response.MaterialStatsResponse, error)

	// Batch operations
	BatchImportMaterials(req *request.BatchImportMaterialRequest) (*response.BatchImportMaterialResponse, error)

	// Calculate operations
	CalculateMaterialUsage(req *request.MaterialCalculateRequest) (*response.MaterialCalculateResponse, error)
	BatchCalculateMaterialUsage(req *request.MaterialBatchCalculateRequest) (*response.MaterialBatchCalculateResponse, error)

	// Recommend operations
	RecommendMaterials(req *request.MaterialRecommendRequest) (*response.MaterialRecommendResponse, error)
}

// materialService 材料服务实现
// Material service implementation
type materialService struct {
	repo     repository.MaterialRepository
	aiClient *volcengine.Client
}

// NewMaterialService 创建材料服务实例
// Create material service instance
func NewMaterialService() MaterialService {
	// 初始化火山AI客户端 / Initialize VolcEngine AI client
	apiKey := os.Getenv("VOLCENGINE_API_KEY")
	baseURL := os.Getenv("VOLCENGINE_BASE_URL")
	model := os.Getenv("VOLCENGINE_MODEL")

	if baseURL == "" {
		baseURL = "https://ark.cn-beijing.volces.com/api/v3"
	}
	if model == "" {
		model = "doubao-1-5-pro-32k-250115"
	}

	var aiClient *volcengine.Client
	if apiKey != "" {
		aiClient = volcengine.NewClient(apiKey, baseURL, model, 60, 4096, 0.7)
	}

	return &materialService{
		repo:     repository.NewMaterialRepository(),
		aiClient: aiClient,
	}
}

// CreateMaterial 创建材料
// Create material
func (s *materialService) CreateMaterial(req *request.CreateMaterialRequest) (*response.MaterialResponse, error) {
	// 检查材料名称是否已存在 / Check if material name already exists
	existing, _ := s.repo.FindByName(req.Name)
	if existing != nil {
		return nil, errors.New("材料名称已存在 / Material name already exists")
	}

	// 序列化适用场景 / Serialize applicable scenes
	scenesJSON, err := json.Marshal(req.ApplicableScenes)
	if err != nil {
		return nil, fmt.Errorf("序列化适用场景失败 / Failed to serialize applicable scenes: %v", err)
	}

	// 设置默认状态 / Set default status
	status := req.Status
	if status == "" {
		status = "active"
	}

	material := &model.Material{
		Name:             req.Name,
		Category:         req.Category,
		Specification:    req.Specification,
		Unit:             req.Unit,
		UnitPrice:        req.UnitPrice,
		Brand:            req.Brand,
		Supplier:         req.Supplier,
		Description:      req.Description,
		ImageURL:         req.ImageURL,
		ApplicableScenes: string(scenesJSON),
		Status:           status,
	}

	if err := s.repo.Create(material); err != nil {
		return nil, fmt.Errorf("创建材料失败 / Failed to create material: %v", err)
	}

	return s.convertToResponse(material), nil
}

// UpdateMaterial 更新材料
// Update material
func (s *materialService) UpdateMaterial(req *request.UpdateMaterialRequest) (*response.MaterialResponse, error) {
	// 查找材料 / Find material
	material, err := s.repo.FindByID(req.ID)
	if err != nil {
		return nil, errors.New("材料不存在 / Material not found")
	}

	// 检查名称是否与其他材料冲突 / Check if name conflicts with other materials
	if req.Name != material.Name {
		existing, _ := s.repo.FindByName(req.Name)
		if existing != nil && existing.ID != req.ID {
			return nil, errors.New("材料名称已存在 / Material name already exists")
		}
	}

	// 序列化适用场景 / Serialize applicable scenes
	scenesJSON, err := json.Marshal(req.ApplicableScenes)
	if err != nil {
		return nil, fmt.Errorf("序列化适用场景失败 / Failed to serialize applicable scenes: %v", err)
	}

	// 更新字段 / Update fields
	material.Name = req.Name
	material.Category = req.Category
	material.Specification = req.Specification
	material.Unit = req.Unit
	material.UnitPrice = req.UnitPrice
	material.Brand = req.Brand
	material.Supplier = req.Supplier
	material.Description = req.Description
	material.ImageURL = req.ImageURL
	material.ApplicableScenes = string(scenesJSON)
	if req.Status != "" {
		material.Status = req.Status
	}

	if err := s.repo.Update(material); err != nil {
		return nil, fmt.Errorf("更新材料失败 / Failed to update material: %v", err)
	}

	return s.convertToResponse(material), nil
}

// DeleteMaterial 删除材料
// Delete material
func (s *materialService) DeleteMaterial(id uint) error {
	// 检查材料是否存在 / Check if material exists
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("材料不存在 / Material not found")
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("删除材料失败 / Failed to delete material: %v", err)
	}

	return nil
}

// BatchDeleteMaterials 批量删除材料
// Batch delete materials
func (s *materialService) BatchDeleteMaterials(ids []uint) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的材料 / Please select materials to delete")
	}

	if err := s.repo.BatchDelete(ids); err != nil {
		return fmt.Errorf("批量删除材料失败 / Failed to batch delete materials: %v", err)
	}

	return nil
}

// GetMaterialByID 根据ID获取材料详情
// Get material detail by ID
func (s *materialService) GetMaterialByID(id uint) (*response.MaterialDetailResponse, error) {
	material, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("材料不存在 / Material not found")
	}

	return s.convertToDetailResponse(material), nil
}

// GetMaterialList 获取材料列表
// Get material list
func (s *materialService) GetMaterialList(req *request.MaterialListRequest) (*response.MaterialListResponse, error) {
	params := repository.MaterialQueryParams{
		Current:  req.Current,
		Size:     req.Size,
		Name:     req.Name,
		Category: req.Category,
		Brand:    req.Brand,
		MinPrice: req.MinPrice,
		MaxPrice: req.MaxPrice,
		Status:   req.Status,
		Keyword:  req.Keyword,
	}

	materials, total, err := s.repo.FindAll(params)
	if err != nil {
		return nil, fmt.Errorf("查询材料列表失败 / Failed to query material list: %v", err)
	}

	records := make([]response.MaterialResponse, len(materials))
	for i, m := range materials {
		records[i] = *s.convertToResponse(&m)
	}

	return &response.MaterialListResponse{
		Records: records,
		Current: req.Current,
		Size:    req.Size,
		Total:   total,
	}, nil
}

// GetCategories 获取所有分类
// Get all categories
func (s *materialService) GetCategories() ([]response.MaterialCategoryResponse, error) {
	categories, err := s.repo.GetCategories()
	if err != nil {
		return nil, fmt.Errorf("获取分类列表失败 / Failed to get categories: %v", err)
	}

	result := make([]response.MaterialCategoryResponse, len(categories))
	for i, c := range categories {
		result[i] = response.MaterialCategoryResponse{
			Category: c.Category,
			Count:    c.Count,
		}
	}

	return result, nil
}

// GetBrands 获取所有品牌
// Get all brands
func (s *materialService) GetBrands() ([]response.MaterialBrandResponse, error) {
	brands, err := s.repo.GetBrands()
	if err != nil {
		return nil, fmt.Errorf("获取品牌列表失败 / Failed to get brands: %v", err)
	}

	result := make([]response.MaterialBrandResponse, len(brands))
	for i, b := range brands {
		result[i] = response.MaterialBrandResponse{
			Brand: b.Brand,
			Count: b.Count,
		}
	}

	return result, nil
}

// GetStats 获取材料统计信息
// Get material statistics
func (s *materialService) GetStats() (*response.MaterialStatsResponse, error) {
	stats, err := s.repo.GetStats()
	if err != nil {
		return nil, fmt.Errorf("获取统计信息失败 / Failed to get statistics: %v", err)
	}

	categories, _ := s.repo.GetCategories()
	brands, _ := s.repo.GetBrands()

	categoryStats := make([]response.MaterialCategoryResponse, len(categories))
	for i, c := range categories {
		categoryStats[i] = response.MaterialCategoryResponse{
			Category: c.Category,
			Count:    c.Count,
		}
	}

	brandStats := make([]response.MaterialBrandResponse, len(brands))
	for i, b := range brands {
		brandStats[i] = response.MaterialBrandResponse{
			Brand: b.Brand,
			Count: b.Count,
		}
	}

	return &response.MaterialStatsResponse{
		TotalCount:    stats.TotalCount,
		ActiveCount:   stats.ActiveCount,
		InactiveCount: stats.InactiveCount,
		CategoryStats: categoryStats,
		BrandStats:    brandStats,
		AvgPrice:      stats.AvgPrice,
		MaxPrice:      stats.MaxPrice,
		MinPrice:      stats.MinPrice,
	}, nil
}

// BatchImportMaterials 批量导入材料
// Batch import materials
func (s *materialService) BatchImportMaterials(req *request.BatchImportMaterialRequest) (*response.BatchImportMaterialResponse, error) {
	var successCount int
	var failCount int
	var failReasons []string

	for i, item := range req.Materials {
		// 检查材料名称是否已存在 / Check if material name already exists
		existing, _ := s.repo.FindByName(item.Name)
		if existing != nil {
			failCount++
			failReasons = append(failReasons, fmt.Sprintf("第%d条: 材料名称'%s'已存在 / Row %d: Material name '%s' already exists", i+1, item.Name, i+1, item.Name))
			continue
		}

		// 序列化适用场景 / Serialize applicable scenes
		scenesJSON, err := json.Marshal(item.ApplicableScenes)
		if err != nil {
			failCount++
			failReasons = append(failReasons, fmt.Sprintf("第%d条: 序列化适用场景失败 / Row %d: Failed to serialize applicable scenes", i+1, i+1))
			continue
		}

		status := item.Status
		if status == "" {
			status = "active"
		}

		material := &model.Material{
			Name:             item.Name,
			Category:         item.Category,
			Specification:    item.Specification,
			Unit:             item.Unit,
			UnitPrice:        item.UnitPrice,
			Brand:            item.Brand,
			Supplier:         item.Supplier,
			Description:      item.Description,
			ImageURL:         item.ImageURL,
			ApplicableScenes: string(scenesJSON),
			Status:           status,
		}

		if err := s.repo.Create(material); err != nil {
			failCount++
			failReasons = append(failReasons, fmt.Sprintf("第%d条: 创建失败 - %v / Row %d: Create failed - %v", i+1, err, i+1, err))
			continue
		}

		successCount++
	}

	return &response.BatchImportMaterialResponse{
		SuccessCount: successCount,
		FailCount:    failCount,
		FailReasons:  failReasons,
	}, nil
}

// CalculateMaterialUsage 计算材料用量
// Calculate material usage
func (s *materialService) CalculateMaterialUsage(req *request.MaterialCalculateRequest) (*response.MaterialCalculateResponse, error) {
	material, err := s.repo.FindByID(req.MaterialID)
	if err != nil {
		return nil, errors.New("材料不存在 / Material not found")
	}

	// 计算基础用量 / Calculate base quantity
	baseQuantity := req.Area

	// 计算含损耗用量 / Calculate quantity with loss
	lossRate := req.LossRate
	if lossRate < 0 {
		lossRate = 0
	}
	if lossRate > 1 {
		lossRate = 1
	}
	quantity := baseQuantity * (1 + lossRate)

	// 计算总成本 / Calculate total cost
	totalCost := quantity * material.UnitPrice

	return &response.MaterialCalculateResponse{
		MaterialID:   material.ID,
		MaterialName: material.Name,
		Unit:         material.Unit,
		UnitPrice:    material.UnitPrice,
		Area:         req.Area,
		LossRate:     lossRate,
		Quantity:     quantity,
		BaseQuantity: baseQuantity,
		TotalCost:    totalCost,
	}, nil
}

// BatchCalculateMaterialUsage 批量计算材料用量
// Batch calculate material usage
func (s *materialService) BatchCalculateMaterialUsage(req *request.MaterialBatchCalculateRequest) (*response.MaterialBatchCalculateResponse, error) {
	var items []response.MaterialCalculateResponse
	var totalCost float64
	var totalQuantity float64

	for _, item := range req.Items {
		result, err := s.CalculateMaterialUsage(&item)
		if err != nil {
			return nil, err
		}
		items = append(items, *result)
		totalCost += result.TotalCost
		totalQuantity += result.Quantity
	}

	return &response.MaterialBatchCalculateResponse{
		Items:         items,
		TotalCost:     totalCost,
		TotalQuantity: totalQuantity,
	}, nil
}

// convertToResponse 转换为响应对象
// Convert to response object
func (s *materialService) convertToResponse(m *model.Material) *response.MaterialResponse {
	var scenes []string
	if m.ApplicableScenes != "" {
		_ = json.Unmarshal([]byte(m.ApplicableScenes), &scenes)
	}

	return &response.MaterialResponse{
		ID:               m.ID,
		Name:             m.Name,
		Category:         m.Category,
		Specification:    m.Specification,
		Unit:             m.Unit,
		UnitPrice:        m.UnitPrice,
		Brand:            m.Brand,
		Supplier:         m.Supplier,
		Description:      m.Description,
		ImageURL:         m.ImageURL,
		ApplicableScenes: scenes,
		Status:           m.Status,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}
}

// convertToDetailResponse 转换为详情响应对象
// Convert to detail response object
func (s *materialService) convertToDetailResponse(m *model.Material) *response.MaterialDetailResponse {
	var scenes []string
	if m.ApplicableScenes != "" {
		_ = json.Unmarshal([]byte(m.ApplicableScenes), &scenes)
	}

	return &response.MaterialDetailResponse{
		ID:               m.ID,
		Name:             m.Name,
		Category:         m.Category,
		Specification:    m.Specification,
		Unit:             m.Unit,
		UnitPrice:        m.UnitPrice,
		Brand:            m.Brand,
		Supplier:         m.Supplier,
		Description:      m.Description,
		ImageURL:         m.ImageURL,
		ApplicableScenes: scenes,
		Status:           m.Status,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}
}

// RecommendMaterials 智能推荐材料
// Intelligently recommend materials based on project requirements
func (s *materialService) RecommendMaterials(req *request.MaterialRecommendRequest) (*response.MaterialRecommendResponse, error) {
	// 设置默认推荐数量 / Set default recommendation limit
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	// 构建筛选条件 / Build filter criteria
	params := repository.MaterialQueryParams{
		Current:  1,
		Size:     100, // 获取足够多的材料用于筛选 / Get enough materials for filtering
		Category: req.Category,
		Status:   "active",
	}

	// 根据预算设置价格区间 / Set price range based on budget
	if req.Budget > 0 && req.Area > 0 {
		// 估算单位面积预算 / Estimate budget per unit area
		budgetPerUnit := req.Budget / req.Area
		params.MaxPrice = budgetPerUnit * 0.3 // 假设材料成本占总预算30% / Assume material cost is 30% of total budget
	}

	// 获取候选材料 / Get candidate materials
	materials, _, err := s.repo.FindAll(params)
	if err != nil {
		return nil, fmt.Errorf("获取材料列表失败 / Failed to get material list: %v", err)
	}

	// 如果有AI客户端，使用AI优化推荐 / If AI client available, use AI to optimize recommendations
	var recommendReason string
	var recommendedMaterials []model.Material

	if s.aiClient != nil && len(materials) > 0 {
		recommendedMaterials, recommendReason = s.aiRecommend(materials, req, limit)
	} else {
		// 无AI时使用基础推荐逻辑 / Use basic recommendation logic without AI
		recommendedMaterials, recommendReason = s.basicRecommend(materials, req, limit)
	}

	// 转换为响应格式 / Convert to response format
	records := make([]response.MaterialResponse, len(recommendedMaterials))
	for i, m := range recommendedMaterials {
		records[i] = *s.convertToResponse(&m)
	}

	return &response.MaterialRecommendResponse{
		Materials:       records,
		RecommendReason: recommendReason,
		TotalCount:      len(records),
	}, nil
}

// aiRecommend 使用AI进行智能推荐
// Use AI for intelligent recommendation
func (s *materialService) aiRecommend(materials []model.Material, req *request.MaterialRecommendRequest, limit int) ([]model.Material, string) {
	// 构建材料列表描述 / Build material list description
	var materialList strings.Builder
	for i, m := range materials {
		if i >= 50 { // 限制发送给AI的材料数量 / Limit materials sent to AI
			break
		}
		materialList.WriteString(fmt.Sprintf("%d. %s - 分类:%s, 品牌:%s, 单价:%.2f元/%s\n",
			m.ID, m.Name, m.Category, m.Brand, m.UnitPrice, m.Unit))
	}

	// 构建AI提示 / Build AI prompt
	systemPrompt := `你是一位专业的工装设计材料顾问。请根据项目需求，从给定的材料列表中推荐最合适的材料。
返回格式要求：
1. 第一行返回推荐的材料ID列表，用逗号分隔，例如：1,5,8,12,15
2. 第二行开始返回推荐理由，简洁说明为什么推荐这些材料`

	userPrompt := fmt.Sprintf(`项目信息：
- 空间类型：%s
- 设计风格：%s
- 预算：%.2f元
- 面积：%.2f平方米
- 材料分类：%s

可选材料列表：
%s

请推荐最多%d种最适合的材料。`,
		req.SpaceType, req.Style, req.Budget, req.Area, req.Category, materialList.String(), limit)

	// 调用AI / Call AI
	aiResponse, err := s.aiClient.ChatWithText(systemPrompt, userPrompt)
	if err != nil {
		// AI调用失败，回退到基础推荐 / AI call failed, fallback to basic recommendation
		return s.basicRecommend(materials, req, limit)
	}

	// 解析AI响应 / Parse AI response
	lines := strings.Split(aiResponse, "\n")
	if len(lines) < 1 {
		return s.basicRecommend(materials, req, limit)
	}

	// 解析推荐的材料ID / Parse recommended material IDs
	idLine := strings.TrimSpace(lines[0])
	idStrs := strings.Split(idLine, ",")

	var recommendedMaterials []model.Material
	materialMap := make(map[uint]model.Material)
	for _, m := range materials {
		materialMap[m.ID] = m
	}

	for _, idStr := range idStrs {
		idStr = strings.TrimSpace(idStr)
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			continue
		}
		if m, ok := materialMap[uint(id)]; ok {
			recommendedMaterials = append(recommendedMaterials, m)
			if len(recommendedMaterials) >= limit {
				break
			}
		}
	}

	// 获取推荐理由 / Get recommendation reason
	var reason string
	if len(lines) > 1 {
		reason = strings.Join(lines[1:], "\n")
	} else {
		reason = "基于项目需求和材料特性的AI智能推荐"
	}

	// 如果AI推荐结果不足，补充基础推荐 / If AI recommendation insufficient, supplement with basic
	if len(recommendedMaterials) < limit {
		basicMaterials, _ := s.basicRecommend(materials, req, limit-len(recommendedMaterials))
		for _, m := range basicMaterials {
			exists := false
			for _, rm := range recommendedMaterials {
				if rm.ID == m.ID {
					exists = true
					break
				}
			}
			if !exists {
				recommendedMaterials = append(recommendedMaterials, m)
			}
		}
	}

	return recommendedMaterials, reason
}

// basicRecommend 基础推荐逻辑（无AI时使用）
// Basic recommendation logic (used when AI is not available)
func (s *materialService) basicRecommend(materials []model.Material, req *request.MaterialRecommendRequest, limit int) ([]model.Material, string) {
	// 根据风格和空间类型进行简单匹配 / Simple matching based on style and space type
	var scored []struct {
		material model.Material
		score    int
	}

	for _, m := range materials {
		score := 0

		// 解析适用场景 / Parse applicable scenes
		var scenes []string
		if m.ApplicableScenes != "" {
			_ = json.Unmarshal([]byte(m.ApplicableScenes), &scenes)
		}

		// 空间类型匹配 / Space type matching
		for _, scene := range scenes {
			if strings.Contains(scene, req.SpaceType) {
				score += 10
			}
		}

		// 风格匹配（简单关键词匹配）/ Style matching (simple keyword matching)
		if req.Style != "" {
			if strings.Contains(m.Description, req.Style) || strings.Contains(m.Name, req.Style) {
				score += 5
			}
		}

		// 价格评分（越接近预算中位数越高）/ Price scoring (closer to budget median is better)
		if req.Budget > 0 && req.Area > 0 {
			targetPrice := (req.Budget / req.Area) * 0.3
			priceDiff := m.UnitPrice - targetPrice
			if priceDiff < 0 {
				priceDiff = -priceDiff
			}
			if priceDiff < targetPrice*0.2 {
				score += 8
			} else if priceDiff < targetPrice*0.5 {
				score += 4
			}
		}

		scored = append(scored, struct {
			material model.Material
			score    int
		}{m, score})
	}

	// 按分数排序 / Sort by score
	for i := 0; i < len(scored)-1; i++ {
		for j := i + 1; j < len(scored); j++ {
			if scored[j].score > scored[i].score {
				scored[i], scored[j] = scored[j], scored[i]
			}
		}
	}

	// 取前N个 / Take top N
	var result []model.Material
	for i := 0; i < len(scored) && i < limit; i++ {
		result = append(result, scored[i].material)
	}

	reason := "基于项目空间类型、设计风格和预算范围的智能匹配推荐"
	return result, reason
}
