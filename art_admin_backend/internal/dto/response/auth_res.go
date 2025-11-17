package response

// LoginResponse 登录响应
type LoginResponse struct {
	Token        string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refreshToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// UserInfoResponse 用户信息响应
type UserInfoResponse struct {
	UserID   int64    `json:"userId" example:"1"`
	UserName string   `json:"userName" example:"admin"`
	Email    string   `json:"email" example:"admin@example.com"`
	Avatar   string   `json:"avatar" example:"https://..."`
	Roles    []string `json:"roles" example:"admin,user"`
	Buttons  []string `json:"buttons" example:"user:add,user:edit"`
}

