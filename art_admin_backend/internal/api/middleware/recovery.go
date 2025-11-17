package middleware

import (
	"art_admin_backend/internal/pkg/logger"
	"art_admin_backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Recovery 异常恢复中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录错误日志
				logger.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
				)

				// 返回错误响应
				response.ServerError(c, "服务器内部错误")
				c.Abort()
			}
		}()

		c.Next()
	}
}

