package repository

import (
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/database"

	"gorm.io/gorm"
)

// MaterialRepository 材料仓库接口
// Material repository interface
type MaterialRepository interface {
	Create(material *model.Material) error
	Update(material *model.Material) error
	Delete(id uint) error
	BatchDelete(ids []uint) error
	FindByID(id uint) (*model.Material, error)
	FindByName(name string) (*model.Material, error)
	FindAll(params MaterialQueryParams) ([]model.Material, int64, error)
	FindByCategory(category string) ([]model.Material, error)
	FindByBrand(brand string) ([]model.Material, error)
	GetCategories() ([]CategoryCount, error)
	GetBrands() ([]BrandCount, error)
	GetStats() (*MaterialStats, error)
	BatchCreate(materials []model.Material) (int, error)
}

// MaterialQueryParams 材料查询参数
// Material query parameters
type MaterialQueryParams struct {
	Current  int
	Size     int
	Name     string
	Category string
	Brand    string
	MinPrice float64
	MaxPrice float64
	Status   string
	Keyword  string
}

// CategoryCount 分类统计
// Category count
type CategoryCount struct {
	Category string
	Count    int64
}

// BrandCount 品牌统计
// Brand count
type BrandCount struct {
	Brand string
	Count int64
}

// MaterialStats 材料统计
// Material stats
type MaterialStats struct {
	TotalCount    int64
	ActiveCount   int64
	InactiveCount int64
	AvgPrice      float64
	MaxPrice      float64
	MinPrice      float64
}

// materialRepository 材料仓库实现
// Material repository implementation
type materialRepository struct {
	db *gorm.DB
}

// NewMaterialRepository 创建材料仓库实例
// Create material repository instance
func NewMaterialRepository() MaterialRepository {
	return &materialRepository{
		db: database.DB,
	}
}

// Create 创建材料
// Create material
func (r *materialRepository) Create(material *model.Material) error {
	return r.db.Create(material).Error
}

// Update 更新材料
// Update material
func (r *materialRepository) Update(material *model.Material) error {
	return r.db.Save(material).Error
}

// Delete 删除材料(软删除)
// Delete material (soft delete)
func (r *materialRepository) Delete(id uint) error {
	return r.db.Delete(&model.Material{}, id).Error
}

// BatchDelete 批量删除材料
// Batch delete materials
func (r *materialRepository) BatchDelete(ids []uint) error {
	return r.db.Delete(&model.Material{}, ids).Error
}

// FindByID 根据ID查找材料
// Find material by ID
func (r *materialRepository) FindByID(id uint) (*model.Material, error) {
	var material model.Material
	err := r.db.First(&material, id).Error
	if err != nil {
		return nil, err
	}
	return &material, nil
}

// FindByName 根据名称查找材料
// Find material by name
func (r *materialRepository) FindByName(name string) (*model.Material, error) {
	var material model.Material
	err := r.db.Where("name = ?", name).First(&material).Error
	if err != nil {
		return nil, err
	}
	return &material, nil
}

// FindAll 查询材料列表(支持分页和筛选)
// Find all materials with pagination and filters
func (r *materialRepository) FindAll(params MaterialQueryParams) ([]model.Material, int64, error) {
	var materials []model.Material
	var total int64

	query := r.db.Model(&model.Material{})

	// 名称模糊搜索 / Name fuzzy search
	if params.Name != "" {
		query = query.Where("name LIKE ?", "%"+params.Name+"%")
	}

	// 分类筛选 / Category filter
	if params.Category != "" {
		query = query.Where("category = ?", params.Category)
	}

	// 品牌筛选 / Brand filter
	if params.Brand != "" {
		query = query.Where("brand = ?", params.Brand)
	}

	// 价格区间筛选 / Price range filter
	if params.MinPrice > 0 {
		query = query.Where("unit_price >= ?", params.MinPrice)
	}
	if params.MaxPrice > 0 {
		query = query.Where("unit_price <= ?", params.MaxPrice)
	}

	// 状态筛选 / Status filter
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	// 关键字搜索(名称/描述/规格) / Keyword search
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("name LIKE ? OR description LIKE ? OR specification LIKE ?", keyword, keyword, keyword)
	}

	// 获取总数 / Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询 / Pagination query
	offset := (params.Current - 1) * params.Size
	if err := query.Order("created_at DESC").Offset(offset).Limit(params.Size).Find(&materials).Error; err != nil {
		return nil, 0, err
	}

	return materials, total, nil
}

// FindByCategory 根据分类查找材料
// Find materials by category
func (r *materialRepository) FindByCategory(category string) ([]model.Material, error) {
	var materials []model.Material
	err := r.db.Where("category = ? AND status = ?", category, "active").Find(&materials).Error
	return materials, err
}

// FindByBrand 根据品牌查找材料
// Find materials by brand
func (r *materialRepository) FindByBrand(brand string) ([]model.Material, error) {
	var materials []model.Material
	err := r.db.Where("brand = ? AND status = ?", brand, "active").Find(&materials).Error
	return materials, err
}

// GetCategories 获取所有分类及数量
// Get all categories with count
func (r *materialRepository) GetCategories() ([]CategoryCount, error) {
	var results []CategoryCount
	err := r.db.Model(&model.Material{}).
		Select("category, COUNT(*) as count").
		Where("status = ?", "active").
		Group("category").
		Order("count DESC").
		Scan(&results).Error
	return results, err
}

// GetBrands 获取所有品牌及数量
// Get all brands with count
func (r *materialRepository) GetBrands() ([]BrandCount, error) {
	var results []BrandCount
	err := r.db.Model(&model.Material{}).
		Select("brand, COUNT(*) as count").
		Where("status = ? AND brand != ''", "active").
		Group("brand").
		Order("count DESC").
		Scan(&results).Error
	return results, err
}

// GetStats 获取材料统计信息
// Get material statistics
func (r *materialRepository) GetStats() (*MaterialStats, error) {
	var stats MaterialStats

	// 总数 / Total count
	if err := r.db.Model(&model.Material{}).Count(&stats.TotalCount).Error; err != nil {
		return nil, err
	}

	// 启用数量 / Active count
	if err := r.db.Model(&model.Material{}).Where("status = ?", "active").Count(&stats.ActiveCount).Error; err != nil {
		return nil, err
	}

	// 停用数量 / Inactive count
	stats.InactiveCount = stats.TotalCount - stats.ActiveCount

	// 价格统计 / Price statistics
	type PriceStats struct {
		AvgPrice float64
		MaxPrice float64
		MinPrice float64
	}
	var priceStats PriceStats
	if err := r.db.Model(&model.Material{}).
		Select("AVG(unit_price) as avg_price, MAX(unit_price) as max_price, MIN(unit_price) as min_price").
		Where("status = ?", "active").
		Scan(&priceStats).Error; err != nil {
		return nil, err
	}
	stats.AvgPrice = priceStats.AvgPrice
	stats.MaxPrice = priceStats.MaxPrice
	stats.MinPrice = priceStats.MinPrice

	return &stats, nil
}

// BatchCreate 批量创建材料
// Batch create materials
func (r *materialRepository) BatchCreate(materials []model.Material) (int, error) {
	result := r.db.Create(&materials)
	return int(result.RowsAffected), result.Error
}
