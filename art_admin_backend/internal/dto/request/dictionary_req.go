package request

// DictionaryTypeListRequest 字典类型列表查询请求
type DictionaryTypeListRequest struct {
	Current     int    `form:"current" binding:"required,min=1" example:"1"`
	Size        int    `form:"size" binding:"required,min=1,max=100" example:"10"`
	ID          *int64 `form:"id" example:"1"`
	TypeName    string `form:"typeName" example:"用户性别"`
	TypeCode    string `form:"typeCode" example:"user_gender"`
	Description string `form:"description" example:"用户性别字典"`
	Enabled     *bool  `form:"enabled" example:"true"`
}

// CreateDictionaryTypeRequest 创建字典类型请求
type CreateDictionaryTypeRequest struct {
	TypeName    string `json:"typeName" binding:"required,max=50"`
	TypeCode    string `json:"typeCode" binding:"required,max=50"`
	Description string `json:"description" binding:"max=200"`
	Enabled     bool   `json:"enabled"`
	Remark      string `json:"remark" binding:"max=500"`
}

// UpdateDictionaryTypeRequest 更新字典类型请求
type UpdateDictionaryTypeRequest struct {
	ID          int64  `json:"id" binding:"required"`
	TypeName    string `json:"typeName" binding:"required,max=50"`
	TypeCode    string `json:"typeCode" binding:"required,max=50"`
	Description string `json:"description" binding:"max=200"`
	Enabled     bool   `json:"enabled"`
	Remark      string `json:"remark" binding:"max=500"`
}

// DeleteDictionaryTypeRequest 删除字典类型请求
type DeleteDictionaryTypeRequest struct {
	ID int64 `json:"id" binding:"required"`
}

// DictionaryListRequest 字典数据列表查询请求
type DictionaryListRequest struct {
	Current  int    `form:"current" binding:"required,min=1" example:"1"`
	Size     int    `form:"size" binding:"required,min=1,max=100" example:"10"`
	ID       *int64 `form:"id" example:"1"`
	TypeCode string `form:"typeCode" example:"user_gender"`
	Label    string `form:"label" example:"男"`
	Value    string `form:"value" example:"1"`
	Enabled  *bool  `form:"enabled" example:"true"`
}

// CreateDictionaryRequest 创建字典数据请求
type CreateDictionaryRequest struct {
	TypeCode string `json:"typeCode" binding:"required,max=50"`
	Label    string `json:"label" binding:"required,max=100"`
	Value    string `json:"value" binding:"required,max=100"`
	OrderNum int    `json:"orderNum"`
	Enabled  bool   `json:"enabled"`
	Remark   string `json:"remark" binding:"max=500"`
}

// UpdateDictionaryRequest 更新字典数据请求
type UpdateDictionaryRequest struct {
	ID       int64  `json:"id" binding:"required"`
	TypeCode string `json:"typeCode" binding:"required,max=50"`
	Label    string `json:"label" binding:"required,max=100"`
	Value    string `json:"value" binding:"required,max=100"`
	OrderNum int    `json:"orderNum"`
	Enabled  bool   `json:"enabled"`
	Remark   string `json:"remark" binding:"max=500"`
}

// DeleteDictionaryRequest 删除字典数据请求
type DeleteDictionaryRequest struct {
	ID int64 `json:"id" binding:"required"`
}

// GetDictionaryByTypeRequest 根据类型获取字典数据请求
type GetDictionaryByTypeRequest struct {
	TypeCode string `form:"typeCode" binding:"required" example:"user_gender"`
}
