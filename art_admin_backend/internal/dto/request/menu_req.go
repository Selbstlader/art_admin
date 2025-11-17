package request

// CreateMenuRequest 创建菜单请求
type CreateMenuRequest struct {
	ParentID      int64    `json:"parentId"`
	Name          string   `json:"name" binding:"required,max=50"`    // 菜单显示名称（前端form.name）
	Label         string   `json:"label" binding:"required,max=100"`  // 路由名称/权限标识（前端form.label）
	Path          string   `json:"path" binding:"required,max=200"`   // 路由地址
	Component     string   `json:"component" binding:"max=200"`       // 组件路径
	Redirect      string   `json:"redirect" binding:"max=200"`        // 重定向路径
	Icon          string   `json:"icon" binding:"max=50"`             // 图标
	IsEnable      bool     `json:"isEnable"`                          // 是否启用
	Sort          int      `json:"sort"`                              // 排序
	IsMenu        bool     `json:"isMenu"`                            // 是否菜单
	KeepAlive     bool     `json:"keepAlive"`                         // 页面缓存
	IsHide        bool     `json:"isHide"`                            // 隐藏菜单
	IsHideTab     bool     `json:"isHideTab"`                         // 隐藏标签
	Link          string   `json:"link" binding:"max=500"`            // 外部链接
	IsIframe      bool     `json:"isIframe"`                          // 是否内嵌
	ShowBadge     bool     `json:"showBadge"`                         // 显示徽章
	ShowTextBadge string   `json:"showTextBadge" binding:"max=50"`    // 文本徽章
	FixedTab      bool     `json:"fixedTab"`                          // 固定标签
	ActivePath    string   `json:"activePath" binding:"max=200"`      // 激活路径
	IsFullPage    bool     `json:"isFullPage"`                        // 全屏页面
	Roles         []string `json:"roles"`                             // 角色权限（前端用）
}

// UpdateMenuRequest 更新菜单请求
type UpdateMenuRequest struct {
	ID            int64    `json:"id" binding:"required"`
	ParentID      int64    `json:"parentId"`
	Name          string   `json:"name" binding:"required,max=50"`
	Label         string   `json:"label" binding:"required,max=100"`
	Path          string   `json:"path" binding:"required,max=200"`
	Component     string   `json:"component" binding:"max=200"`
	Redirect      string   `json:"redirect" binding:"max=200"`
	Icon          string   `json:"icon" binding:"max=50"`
	IsEnable      bool     `json:"isEnable"`
	Sort          int      `json:"sort"`
	IsMenu        bool     `json:"isMenu"`
	KeepAlive     bool     `json:"keepAlive"`
	IsHide        bool     `json:"isHide"`
	IsHideTab     bool     `json:"isHideTab"`
	Link          string   `json:"link" binding:"max=500"`
	IsIframe      bool     `json:"isIframe"`
	ShowBadge     bool     `json:"showBadge"`
	ShowTextBadge string   `json:"showTextBadge" binding:"max=50"`
	FixedTab      bool     `json:"fixedTab"`
	ActivePath    string   `json:"activePath" binding:"max=200"`
	IsFullPage    bool     `json:"isFullPage"`
	Roles         []string `json:"roles"`
}

// DeleteMenuRequest 删除菜单请求
type DeleteMenuRequest struct {
	ID int64 `json:"id" binding:"required"`
}
