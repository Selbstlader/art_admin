package travel

import (
	"art_admin_backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

// FavoriteRoadbook 收藏路书
func FavoriteRoadbook(c *gin.Context) {
	// TODO: 实现收藏路书功能
	response.Success(c, nil)
}

// UnfavoriteRoadbook 取消收藏路书
func UnfavoriteRoadbook(c *gin.Context) {
	// TODO: 实现取消收藏路书功能
	response.Success(c, nil)
}

// GetFavoriteList 获取收藏列表
func GetFavoriteList(c *gin.Context) {
	// TODO: 实现获取收藏列表功能
	response.Success(c, nil)
}
