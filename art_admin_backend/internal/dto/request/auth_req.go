package request

// LoginRequest 登录请求
type LoginRequest struct {
	UserName string `json:"userName" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required" example:"123456"`
}
