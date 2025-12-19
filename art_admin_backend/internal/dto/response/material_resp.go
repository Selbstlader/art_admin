package response

import "time"

// MaterialResponse 材料响应
// Material response
type MaterialResponse struct {
	ID               uint      `json:"id"`
	Name             string    `json:"name"`             // 材料名称 / Material name
	Category         string    `json:"category"`         // 分类 / Category
	Specification    string    `json:"specification"`    // 规格 / Specification
	Unit             string    `json:"unit"`             // 单位 / Unit
	UnitPrice        float64   `json:"unitPrice"`        // 单价 / Unit price
	Brand            string    `json:"brand"`            // 品牌 / Brand
	Supplier         string    `json:"supplier"`         // 供应商 / Supplier
	Description      string    `json:"description"`      // 描述 / Description
	ImageURL         string    `json:"imageUrl"`         // 图片URL / Image URL
	ApplicableScenes []string  `json:"applicableScenes"` // 适用场景 / Applicable scenes
	Status           string    `json:"status"`           // 状态 / Status
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

// MaterialDetailResponse 材料详情响应
// Material detail response
type MaterialDetailResponse struct {
	ID               uint      `json:"id"`
	Name             string    `json:"name"`             // 材料名称 / Material name
	Category         string    `json:"category"`         // 分类 / Category
	Specification    string    `json:"specification"`    // 规格 / Specification
	Unit             string    `json:"unit"`             // 单位 / Unit
	UnitPrice        float64   `json:"unitPrice"`        // 单价 / Unit price
	Brand            string    `json:"brand"`            // 品牌 / Brand
	Supplier         string    `json:"supplier"`         // 供应商 / Supplier
	Description      string    `json:"description"`      // 描述 / Description
	ImageURL         string    `json:"imageUrl"`         // 图片URL / Image URL
	ApplicableScenes []string  `json:"applicableScenes"` // 适用场景 / Applicable scenes
	Status           string    `json:"status"`           // 状态 / Status
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

// MaterialListResponse 材料列表响应
// Material list response
type MaterialListResponse struct {
	Records []MaterialResponse `json:"records"` // 记录列表 / Record list
	Current int                `json:"current"` // 当前页码 / Current page
	Size    int                `json:"size"`    // 每页条数 / Page size
	Total   int64              `json:"total"`   // 总记录数 / Total count
}

// MaterialCategoryResponse 材料分类响应
// Material category response
type MaterialCategoryResponse struct {
	Category string `json:"category"` // 分类名称 / Category name
	Count    int64  `json:"count"`    // 该分类下材料数量 / Material count in this category
}

// MaterialBrandResponse 材料品牌响应
// Material brand response
type MaterialBrandResponse struct {
	Brand string `json:"brand"` // 品牌名称 / Brand name
	Count int64  `json:"count"` // 该品牌下材料数量 / Material count for this brand
}

// MaterialRecommendResponse 材料推荐响应
// Material recommend response
type MaterialRecommendResponse struct {
	Materials       []MaterialResponse `json:"materials"`       // 推荐材料列表 / Recommended materials
	RecommendReason string             `json:"recommendReason"` // 推荐理由 / Recommendation reason
	TotalCount      int                `json:"totalCount"`      // 推荐总数 / Total recommendation count
}

// MaterialCalculateResponse 材料用量计算响应
// Material calculate response
type MaterialCalculateResponse struct {
	MaterialID   uint    `json:"materialId"`   // 材料ID / Material ID
	MaterialName string  `json:"materialName"` // 材料名称 / Material name
	Unit         string  `json:"unit"`         // 单位 / Unit
	UnitPrice    float64 `json:"unitPrice"`    // 单价 / Unit price
	Area         float64 `json:"area"`         // 面积 / Area
	LossRate     float64 `json:"lossRate"`     // 损耗率 / Loss rate
	Quantity     float64 `json:"quantity"`     // 用量(含损耗) / Quantity with loss
	BaseQuantity float64 `json:"baseQuantity"` // 基础用量(不含损耗) / Base quantity without loss
	TotalCost    float64 `json:"totalCost"`    // 总成本 / Total cost
}

// MaterialBatchCalculateResponse 批量材料用量计算响应
// Batch material calculate response
type MaterialBatchCalculateResponse struct {
	Items         []MaterialCalculateResponse `json:"items"`         // 计算结果列表 / Calculate results
	TotalCost     float64                     `json:"totalCost"`     // 总成本 / Total cost
	TotalQuantity float64                     `json:"totalQuantity"` // 总用量 / Total quantity
}

// MaterialStatsResponse 材料统计响应
// Material stats response
type MaterialStatsResponse struct {
	TotalCount    int64                      `json:"totalCount"`    // 材料总数 / Total material count
	ActiveCount   int64                      `json:"activeCount"`   // 启用材料数 / Active material count
	InactiveCount int64                      `json:"inactiveCount"` // 停用材料数 / Inactive material count
	CategoryStats []MaterialCategoryResponse `json:"categoryStats"` // 分类统计 / Category stats
	BrandStats    []MaterialBrandResponse    `json:"brandStats"`    // 品牌统计 / Brand stats
	AvgPrice      float64                    `json:"avgPrice"`      // 平均单价 / Average unit price
	MaxPrice      float64                    `json:"maxPrice"`      // 最高单价 / Maximum unit price
	MinPrice      float64                    `json:"minPrice"`      // 最低单价 / Minimum unit price
}

// BatchImportMaterialResponse 批量导入材料响应
// Batch import material response
type BatchImportMaterialResponse struct {
	SuccessCount int      `json:"successCount"` // 成功数量 / Success count
	FailCount    int      `json:"failCount"`    // 失败数量 / Fail count
	FailReasons  []string `json:"failReasons"`  // 失败原因列表 / Fail reasons
}
