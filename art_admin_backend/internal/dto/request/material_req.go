package request

// MaterialListRequest 材料列表查询请求
// Material list query request
type MaterialListRequest struct {
	Current  int     `form:"current" binding:"required,min=1" example:"1"`
	Size     int     `form:"size" binding:"required,min=1,max=100" example:"10"`
	Name     string  `form:"name" example:"大理石"`      // 材料名称模糊搜索 / Material name fuzzy search
	Category string  `form:"category" example:"地面材料"` // 分类筛选 / Category filter
	Brand    string  `form:"brand" example:"东鹏"`      // 品牌筛选 / Brand filter
	MinPrice float64 `form:"minPrice" example:"100"`  // 最低单价 / Minimum unit price
	MaxPrice float64 `form:"maxPrice" example:"1000"` // 最高单价 / Maximum unit price
	Status   string  `form:"status" example:"active"` // 状态筛选 / Status filter
	Keyword  string  `form:"keyword" example:"瓷砖"`    // 关键字搜索(名称/描述/规格) / Keyword search
}

// CreateMaterialRequest 创建材料请求
// Create material request
type CreateMaterialRequest struct {
	Name             string   `json:"name" binding:"required,min=1,max=200"`            // 材料名称 / Material name
	Category         string   `json:"category" binding:"required,max=100"`              // 分类 / Category
	Specification    string   `json:"specification" binding:"max=200"`                  // 规格 / Specification
	Unit             string   `json:"unit" binding:"required,max=50"`                   // 单位 / Unit (m²/m/个/kg)
	UnitPrice        float64  `json:"unitPrice" binding:"min=0"`                        // 单价 / Unit price
	Brand            string   `json:"brand" binding:"max=100"`                          // 品牌 / Brand
	Supplier         string   `json:"supplier" binding:"max=200"`                       // 供应商 / Supplier
	Description      string   `json:"description" binding:"max=2000"`                   // 描述 / Description
	ImageURL         string   `json:"imageUrl" binding:"max=500"`                       // 图片URL / Image URL
	ApplicableScenes []string `json:"applicableScenes"`                                 // 适用场景 / Applicable scenes
	Status           string   `json:"status" binding:"omitempty,oneof=active inactive"` // 状态 / Status
}

// UpdateMaterialRequest 更新材料请求
// Update material request
type UpdateMaterialRequest struct {
	ID               uint     `json:"id" binding:"required"`                            // 材料ID / Material ID
	Name             string   `json:"name" binding:"required,min=1,max=200"`            // 材料名称 / Material name
	Category         string   `json:"category" binding:"required,max=100"`              // 分类 / Category
	Specification    string   `json:"specification" binding:"max=200"`                  // 规格 / Specification
	Unit             string   `json:"unit" binding:"required,max=50"`                   // 单位 / Unit
	UnitPrice        float64  `json:"unitPrice" binding:"min=0"`                        // 单价 / Unit price
	Brand            string   `json:"brand" binding:"max=100"`                          // 品牌 / Brand
	Supplier         string   `json:"supplier" binding:"max=200"`                       // 供应商 / Supplier
	Description      string   `json:"description" binding:"max=2000"`                   // 描述 / Description
	ImageURL         string   `json:"imageUrl" binding:"max=500"`                       // 图片URL / Image URL
	ApplicableScenes []string `json:"applicableScenes"`                                 // 适用场景 / Applicable scenes
	Status           string   `json:"status" binding:"omitempty,oneof=active inactive"` // 状态 / Status
}

// DeleteMaterialRequest 删除材料请求
// Delete material request
type DeleteMaterialRequest struct {
	ID uint `json:"id" binding:"required"` // 材料ID / Material ID
}

// BatchDeleteMaterialRequest 批量删除材料请求
// Batch delete materials request
type BatchDeleteMaterialRequest struct {
	IDs []uint `json:"ids" binding:"required,min=1"` // 材料ID列表 / Material ID list
}

// BatchImportMaterialRequest 批量导入材料请求
// Batch import materials request
type BatchImportMaterialRequest struct {
	Materials []CreateMaterialRequest `json:"materials" binding:"required,min=1,dive"` // 材料列表 / Material list
}

// MaterialRecommendRequest 材料推荐请求
// Material recommend request
type MaterialRecommendRequest struct {
	ProjectID uint    `json:"projectId" binding:"required"` // 项目ID / Project ID
	SpaceType string  `json:"spaceType" binding:"max=100"`  // 空间类型 / Space type (办公/商业/工业)
	Style     string  `json:"style" binding:"max=100"`      // 设计风格 / Design style
	Budget    float64 `json:"budget" binding:"min=0"`       // 预算 / Budget
	Area      float64 `json:"area" binding:"min=0"`         // 面积 / Area
	Category  string  `json:"category" binding:"max=100"`   // 材料分类 / Material category
	Limit     int     `json:"limit" binding:"min=1,max=50"` // 推荐数量限制 / Recommendation limit
}

// MaterialCalculateRequest 材料用量计算请求
// Material calculate request
type MaterialCalculateRequest struct {
	MaterialID uint    `json:"materialId" binding:"required"`  // 材料ID / Material ID
	Area       float64 `json:"area" binding:"required,min=0"`  // 面积 / Area
	LossRate   float64 `json:"lossRate" binding:"min=0,max=1"` // 损耗率 / Loss rate (0-1)
}

// MaterialBatchCalculateRequest 批量材料用量计算请求
// Batch material calculate request
type MaterialBatchCalculateRequest struct {
	Items []MaterialCalculateRequest `json:"items" binding:"required,min=1,dive"` // 计算项列表 / Calculate items
}
