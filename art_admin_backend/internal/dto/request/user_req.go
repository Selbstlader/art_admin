package request

// UserListRequest 用户列表查询请求
type UserListRequest struct {
	Current    int    `form:"current" binding:"required,min=1" example:"1"`
	Size       int    `form:"size" binding:"required,min=1,max=1000" example:"10"`
	ID         *int64 `form:"id" example:"1"`
	UserName   string `form:"userName" example:"admin"`
	UserGender string `form:"userGender" example:"male"`
	UserPhone  string `form:"userPhone" example:"138"`
	UserEmail  string `form:"userEmail" example:"@qq.com"`
	Status     string `form:"status" example:"1"`
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	UserName   string  `json:"userName" binding:"required,min=3,max=50"`
	NickName   string  `json:"nickName" binding:"required,max=50"`
	Password   string  `json:"password" binding:"required,min=6"`
	Email      string  `json:"email" binding:"required,email"`
	UserPhone  string  `json:"userPhone" binding:"required"`
	UserGender string  `json:"userGender" binding:"required,oneof=male female unknown"`
	Avatar     string  `json:"avatar"`
	Status     string  `json:"status" binding:"required,oneof=0 1"`
	RoleIDs    []int64 `json:"roleIds" binding:"required,min=1"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	ID         int64   `json:"id" binding:"required"`
	NickName   string  `json:"nickName" binding:"required,max=50"`
	Email      string  `json:"email" binding:"required,email"`
	UserPhone  string  `json:"userPhone" binding:"required"`
	UserGender string  `json:"userGender" binding:"required,oneof=male female unknown"`
	Avatar     string  `json:"avatar"`
	Status     string  `json:"status" binding:"required,oneof=0 1"`
	RoleIDs    []int64 `json:"roleIds" binding:"required,min=1"`
}

// DeleteUserRequest 删除用户请求
type DeleteUserRequest struct {
	ID int64 `json:"id" binding:"required"`
}

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
	ID          int64  `json:"id" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}
