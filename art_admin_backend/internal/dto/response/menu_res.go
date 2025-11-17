package response

// MenuResponse 菜单响应（树形结构）
type MenuResponse struct {
	ID        int64           `json:"id"`
	ParentID  int64           `json:"parentId"`          // 父菜单ID
	Name      string          `json:"name"`
	Path      string          `json:"path"`
	Component string          `json:"component,omitempty"`
	Redirect  string          `json:"redirect,omitempty"`
	Meta      MenuMeta        `json:"meta"`
	Children  []*MenuResponse `json:"children,omitempty"`
}

// MenuMeta 菜单元数据
type MenuMeta struct {
	Title      string       `json:"title"`
	Icon       string       `json:"icon,omitempty"`
	IsHide     bool         `json:"isHide"`
	IsHideTab  bool         `json:"isHideTab"`
	KeepAlive  bool         `json:"keepAlive"`
	FixedTab   bool         `json:"fixedTab"`
	IsFullPage bool         `json:"isFullPage"`
	AuthList   []AuthButton `json:"authList,omitempty"`
}

// AuthButton 按钮权限
type AuthButton struct {
	Title    string `json:"title"`
	AuthMark string `json:"authMark"`
}

