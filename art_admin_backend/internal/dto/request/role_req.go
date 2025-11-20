package request

// RoleListRequest 角色列表查询请求
type RoleListRequest struct {
	Current     int    `form:"current" binding:"required,min=1" example:"1"`
	Size        int    `form:"size" binding:"required,min=1,max=1000" example:"10"`
	RoleID      *int64 `form:"roleId" example:"1"`
	RoleName    string `form:"roleName" example:"管理员"`
	RoleCode    string `form:"roleCode" example:"admin"`
	Description string `form:"description" example:"系统"`
	Enabled     *bool  `form:"enabled" example:"true"`
}

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	RoleName    string  `json:"roleName" binding:"required,max=50"`
	RoleCode    string  `json:"roleCode" binding:"required,max=50"`
	Description string  `json:"description" binding:"max=200"`
	Enabled     bool    `json:"enabled"`
	MenuIDs     []int64 `json:"menuIds"`
	ButtonIDs   []int64 `json:"buttonIds"`
}

// UpdateRoleRequest 更新角色请求
type UpdateRoleRequest struct {
	RoleID      int64   `json:"roleId" binding:"required"`
	RoleName    string  `json:"roleName" binding:"required,max=50"`
	RoleCode    string  `json:"roleCode" binding:"required,max=50"`
	Description string  `json:"description" binding:"max=200"`
	Enabled     bool    `json:"enabled"`
	MenuIDs     []int64 `json:"menuIds"`
	ButtonIDs   []int64 `json:"buttonIds"`
}

// DeleteRoleRequest 删除角色请求
type DeleteRoleRequest struct {
	RoleID int64 `json:"roleId" binding:"required"`
}

// UpdateRolePermissionsRequest 更新角色权限请求
type UpdateRolePermissionsRequest struct {
	MenuIDs   []int64 `json:"menuIds" binding:"required"`
	ButtonIDs []int64 `json:"buttonIds" binding:"required"`
}
