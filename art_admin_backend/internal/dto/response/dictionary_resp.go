package response

import "time"

// DictionaryTypeResponse 字典类型响应
type DictionaryTypeResponse struct {
	ID          int64     `json:"id"`
	TypeName    string    `json:"typeName"`
	TypeCode    string    `json:"typeCode"`
	Description string    `json:"description"`
	Enabled     bool      `json:"enabled"`
	Remark      string    `json:"remark"`
	CreateBy    string    `json:"createBy"`
	CreateTime  time.Time `json:"createTime"`
	UpdateBy    string    `json:"updateBy"`
	UpdateTime  time.Time `json:"updateTime"`
}

// DictionaryResponse 字典数据响应
type DictionaryResponse struct {
	ID         int64     `json:"id"`
	TypeCode   string    `json:"typeCode"`
	Label      string    `json:"label"`
	Value      string    `json:"value"`
	OrderNum   int       `json:"orderNum"`
	Enabled    bool      `json:"enabled"`
	Remark     string    `json:"remark"`
	CreateBy   string    `json:"createBy"`
	CreateTime time.Time `json:"createTime"`
	UpdateBy   string    `json:"updateBy"`
	UpdateTime time.Time `json:"updateTime"`
}

// DictionaryWithTypeResponse 包含类型的字典数据响应
type DictionaryWithTypeResponse struct {
	ID             int64                  `json:"id"`
	TypeCode       string                 `json:"typeCode"`
	Label          string                 `json:"label"`
	Value          string                 `json:"value"`
	OrderNum       int                    `json:"orderNum"`
	Enabled        bool                   `json:"enabled"`
	Remark         string                 `json:"remark"`
	CreateBy       string                 `json:"createBy"`
	CreateTime     time.Time              `json:"createTime"`
	UpdateBy       string                 `json:"updateBy"`
	UpdateTime     time.Time              `json:"updateTime"`
	DictionaryType DictionaryTypeResponse `json:"dictionaryType"`
}
