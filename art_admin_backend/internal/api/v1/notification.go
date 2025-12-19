package v1

import (
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/response"
	"art_admin_backend/internal/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var notificationRepo *repository.NotificationRepository

// SetNotificationRepository 设置通知仓库
func SetNotificationRepository(repo *repository.NotificationRepository) {
	notificationRepo = repo
}

// GetNotificationList 获取通知列表
// @Summary 获取通知列表
// @Description 获取当前用户的通知列表（只返回今天的）
// @Tags 通知
// @Produce json
// @Param current query int false "页码" default(1)
// @Param size query int false "每页数量" default(20)
// @Success 200 {object} response.Response
// @Router /api/notification/list [get]
func GetNotificationList(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	current, _ := strconv.Atoi(c.DefaultQuery("current", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	notifications, total, err := notificationRepo.ListByUserID(userID.(int64), current, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, map[string]any{
		"records": notifications,
		"total":   total,
		"current": current,
		"size":    size,
	})
}

// GetUnreadCount 获取未读通知数量
// @Summary 获取未读通知数量
// @Description 获取当前用户的未读通知数量（只统计今天的）
// @Tags 通知
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/notification/unread-count [get]
func GetUnreadCount(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	count, err := notificationRepo.GetUnreadCount(userID.(int64))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, map[string]any{
		"count": count,
	})
}

// MarkNotificationRead 标记通知为已读
// @Summary 标记通知为已读
// @Description 标记指定通知为已读
// @Tags 通知
// @Produce json
// @Param id path int true "通知ID"
// @Success 200 {object} response.Response
// @Router /api/notification/read/{id} [post]
func MarkNotificationRead(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的通知ID")
		return
	}

	if err := notificationRepo.MarkAsRead(uint(id), userID.(int64)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// MarkAllNotificationsRead 标记全部通知为已读
// @Summary 标记全部通知为已读
// @Description 标记当前用户的全部通知为已读
// @Tags 通知
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/notification/read-all [post]
func MarkAllNotificationsRead(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	if err := notificationRepo.MarkAllAsRead(userID.(int64)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// DeleteNotification 删除通知
// @Summary 删除通知
// @Description 删除指定通知
// @Tags 通知
// @Produce json
// @Param id path int true "通知ID"
// @Success 200 {object} response.Response
// @Router /api/notification/{id} [delete]
func DeleteNotification(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的通知ID")
		return
	}

	if err := notificationRepo.Delete(uint(id), userID.(int64)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// CreateNotification 创建通知（内部使用）
// Create notification (internal use)
func CreateNotification(userID int64, title, content, noticeType string, relatedID uint) error {
	if notificationRepo == nil {
		return nil
	}
	notification := &model.Notification{
		UserID:    userID,
		Title:     title,
		Content:   content,
		Type:      noticeType,
		RelatedID: relatedID,
		IsRead:    false,
	}
	return notificationRepo.Create(notification)
}
