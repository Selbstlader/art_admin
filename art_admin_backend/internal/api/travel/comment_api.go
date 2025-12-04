package travel

import (
	"art_admin_backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

// CreateComment 创建评论
func CreateComment(c *gin.Context) {
	// TODO: 实现创建评论功能
	response.Success(c, nil)
}

// GetCommentList 获取评论列表
func GetCommentList(c *gin.Context) {
	// TODO: 实现获取评论列表功能
	response.Success(c, nil)
}

// DeleteComment 删除评论
func DeleteComment(c *gin.Context) {
	// TODO: 实现删除评论功能
	response.Success(c, nil)
}

// LikeComment 点赞评论
func LikeComment(c *gin.Context) {
	// TODO: 实现点赞评论功能
	response.Success(c, nil)
}
