package travel

import (
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/response"
	travelSvc "art_admin_backend/internal/service/travel"
	"strconv"

	"github.com/gin-gonic/gin"
)

var commentService = travelSvc.NewCommentService()

// CreateComment 创建评论
// @Summary 创建评论
// @Description 为路书创建评论或回复
// @Tags 旅游-评论
// @Accept json
// @Produce json
// @Param request body travelSvc.CreateCommentRequest true "评论信息"
// @Success 200 {object} response.Response{data=model.RoadbookComment}
// @Router /api/v1/travel/roadbooks/{id}/comments [post]
func CreateComment(c *gin.Context) {
	var req travelSvc.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 从路径获取路书ID
	roadbookIDStr := c.Param("id")
	roadbookID, err := strconv.ParseInt(roadbookIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的路书ID")
		return
	}
	req.RoadbookID = roadbookID

	// 验证评论内容
	if !model.ValidateCommentContent(req.Content) {
		response.BadRequest(c, "评论内容无效：不能为空或超过500字符")
		return
	}

	// 从上下文获取当前用户ID
	userID, _ := c.Get("userId")
	req.UserID = userID.(int64)

	comment, err := commentService.CreateComment(&req)
	if err != nil {
		switch err.Error() {
		case "invalid_content":
			response.BadRequest(c, "评论内容无效")
		case "roadbook_not_found":
			response.NotFound(c, "路书不存在")
		case "parent_comment_not_found":
			response.NotFound(c, "父评论不存在或不属于该路书")
		default:
			response.ServerError(c, "创建评论失败: "+err.Error())
		}
		return
	}

	response.Success(c, comment)
}

// GetCommentList 获取评论列表
// @Summary 获取评论列表
// @Description 获取路书的评论列表（树形结构）
// @Tags 旅游-评论
// @Accept json
// @Produce json
// @Param id path int true "路书ID"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=response.PaginatedData}
// @Router /api/v1/travel/roadbooks/{id}/comments [get]
func GetCommentList(c *gin.Context) {
	// 从路径获取路书ID
	roadbookIDStr := c.Param("id")
	roadbookID, err := strconv.ParseInt(roadbookIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的路书ID")
		return
	}

	var req travelSvc.GetCommentListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	req.RoadbookID = roadbookID

	// 设置默认值
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	comments, total, err := commentService.GetCommentList(&req)
	if err != nil {
		response.ServerError(c, "获取评论列表失败")
		return
	}

	response.SuccessWithPagination(c, comments, req.Page, req.PageSize, total)
}

// DeleteComment 删除评论
// @Summary 删除评论
// @Description 删除自己的评论
// @Tags 旅游-评论
// @Accept json
// @Produce json
// @Param id path int true "路书ID"
// @Param commentId path int true "评论ID"
// @Success 200 {object} response.Response
// @Router /api/v1/travel/roadbooks/{id}/comments/{commentId} [delete]
func DeleteComment(c *gin.Context) {
	commentIDStr := c.Param("commentId")
	commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的评论ID")
		return
	}

	// 从上下文获取当前用户ID
	userID, _ := c.Get("userId")

	err = commentService.DeleteComment(commentID, userID.(int64))
	if err != nil {
		switch err.Error() {
		case "comment_not_found":
			response.NotFound(c, "评论不存在")
		case "forbidden":
			response.Forbidden(c, "无权删除该评论")
		default:
			response.ServerError(c, "删除评论失败: "+err.Error())
		}
		return
	}

	response.Success(c, nil)
}

// LikeComment 点赞评论
// @Summary 点赞评论
// @Description 为评论点赞
// @Tags 旅游-评论
// @Accept json
// @Produce json
// @Param id path int true "路书ID"
// @Param commentId path int true "评论ID"
// @Success 200 {object} response.Response
// @Router /api/v1/travel/roadbooks/{id}/comments/{commentId}/like [post]
func LikeComment(c *gin.Context) {
	commentIDStr := c.Param("commentId")
	commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的评论ID")
		return
	}

	err = commentService.LikeComment(commentID)
	if err != nil {
		if err.Error() == "comment_not_found" {
			response.NotFound(c, "评论不存在")
			return
		}
		response.ServerError(c, "点赞失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// UnlikeComment 取消点赞评论
// @Summary 取消点赞评论
// @Description 取消评论点赞
// @Tags 旅游-评论
// @Accept json
// @Produce json
// @Param id path int true "路书ID"
// @Param commentId path int true "评论ID"
// @Success 200 {object} response.Response
// @Router /api/v1/travel/roadbooks/{id}/comments/{commentId}/like [delete]
func UnlikeComment(c *gin.Context) {
	commentIDStr := c.Param("commentId")
	commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的评论ID")
		return
	}

	err = commentService.UnlikeComment(commentID)
	if err != nil {
		if err.Error() == "comment_not_found" {
			response.NotFound(c, "评论不存在")
			return
		}
		response.ServerError(c, "取消点赞失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}
