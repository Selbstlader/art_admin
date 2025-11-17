package response

import "time"

// RoleListItem 角色列表项
type RoleListItem struct {
	RoleID      int64     `json:"roleId"`
	RoleName    string    `json:"roleName"`
	RoleCode    string    `json:"roleCode"`
	Description string    `json:"description"`
	Enabled     bool      `json:"enabled"`
	CreateTime  time.Time `json:"createTime"`
}

