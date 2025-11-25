package response

// LoginResponse 登录响应
type LoginResponse struct {
	Token        string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refreshToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// UserInfoResponse 用户信息响应
type UserInfoResponse struct {
	UserID     int64    `json:"userId" example:"1"`
	UserName   string   `json:"userName" example:"admin"`
	NickName   string   `json:"nickName" example:"管理员"`
	Email      string   `json:"email" example:"admin@example.com"`
	UserPhone  string   `json:"userPhone" example:"13800138000"`
	UserGender string   `json:"userGender" example:"male"`
	Avatar     string   `json:"avatar" example:"https://..."`
	Address    string   `json:"address" example:"广东省深圳市"`
	Des        string   `json:"des" example:"个人介绍"`
	Roles      []string `json:"roles" example:"admin,user"`
	Buttons    []string `json:"buttons" example:"user:add,user:edit"`
}
