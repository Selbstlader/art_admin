package request

// AppUserListRequest APP用户列表查询请求
type AppUserListRequest struct {
	Current  int    `form:"current" binding:"required,min=1" example:"1"`
	Size     int    `form:"size" binding:"required,min=1,max=1000" example:"10"`
	ID       *int64 `form:"id" example:"1"`
	UserName string `form:"userName" example:"user001"`
	NickName string `form:"nickName" example:"张三"`
	Phone    string `form:"phone" example:"138"`
	UserType string `form:"userType" example:"1"`
	Status   string `form:"status" example:"1"`
}

// CreateAppUserRequest 创建APP用户请求
type CreateAppUserRequest struct {
	UserName string `json:"userName" binding:"required,min=3,max=50"`
	NickName string `json:"nickName" binding:"required,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Phone    string `json:"phone" binding:"required"`
	Email    string `json:"email" binding:"omitempty,email"`
	Avatar   string `json:"avatar"`
	UserType string `json:"userType" binding:"required,oneof=1 2"`
	Status   string `json:"status" binding:"required,oneof=1 2"`
}

// UpdateAppUserRequest 更新APP用户请求
type UpdateAppUserRequest struct {
	ID       int64  `json:"id" binding:"required"`
	NickName string `json:"nickName" binding:"required,max=50"`
	Phone    string `json:"phone" binding:"required"`
	Email    string `json:"email" binding:"omitempty,email"`
	Avatar   string `json:"avatar"`
	UserType string `json:"userType" binding:"required,oneof=1 2"`
	Status   string `json:"status" binding:"required,oneof=1 2"`
}

// DeleteAppUserRequest 删除APP用户请求
type DeleteAppUserRequest struct {
	ID int64 `json:"id" binding:"required"`
}

// ResetAppUserPasswordRequest 重置APP用户密码请求
type ResetAppUserPasswordRequest struct {
	ID          int64  `json:"id" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}
