package middleware

import (
	"strings"

	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/database"
	"art_admin_backend/internal/pkg/jwt"
	"art_admin_backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string

		// 优先从请求头获取 Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			// 检查格式：Bearer {token}
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
			}
		}

		// 如果 Header 中没有 token，尝试从查询参数获取（用于 WebSocket）
		if token == "" {
			token = c.Query("token")
		}

		// 如果仍然没有 token，返回未授权
		if token == "" {
			response.Unauthorized(c, "未提供认证令牌")
			c.Abort()
			return
		}

		// 解析令牌
		claims, err := jwt.ParseToken(token)
		if err != nil {
			response.Unauthorized(c, "认证令牌无效或已过期")
			c.Abort()
			return
		}

		// 将用户ID存入上下文
		c.Set("userId", claims.UserID)

		// 获取用户名并存入上下文
		if username, err := getUsernameByID(claims.UserID); err == nil {
			c.Set("username", username)
		}

		c.Next()
	}
}

// getUsernameByID 根据用户ID获取用户名
func getUsernameByID(userID int64) (string, error) {
	var user model.User
	if err := database.GetDB().Select("user_name").First(&user, userID).Error; err != nil {
		return "", err
	}
	return user.UserName, nil
}

// GetUserID 从上下文获取用户ID
func GetUserID(c *gin.Context) int64 {
	if userID, exists := c.Get("userId"); exists {
		return userID.(int64)
	}
	return 0
}

// GetUsername 从上下文获取用户名
func GetUsername(c *gin.Context) string {
	if username, exists := c.Get("username"); exists {
		return username.(string)
	}
	return ""
}
