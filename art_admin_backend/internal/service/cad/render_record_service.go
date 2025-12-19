package cad

import (
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/repository"
	"fmt"
	"time"
)

// RenderRecordService 效果图记录服务
// Render record service for managing async render tasks
type RenderRecordService struct {
	recordRepo       *repository.RenderRecordRepository
	quotaRepo        *repository.UserRenderQuotaRepository
	renderService    *RenderService
	notificationRepo *repository.NotificationRepository
}

// NewRenderRecordService 创建效果图记录服务
func NewRenderRecordService(
	recordRepo *repository.RenderRecordRepository,
	quotaRepo *repository.UserRenderQuotaRepository,
	renderService *RenderService,
	notificationRepo *repository.NotificationRepository,
) *RenderRecordService {
	return &RenderRecordService{
		recordRepo:       recordRepo,
		quotaRepo:        quotaRepo,
		renderService:    renderService,
		notificationRepo: notificationRepo,
	}
}

// SubmitRenderRequest 提交渲染任务请求
type SubmitRenderRequest struct {
	UserID      int64  `json:"userId"`
	ProjectID   uint   `json:"projectId"`
	CadFileID   uint   `json:"cadFileId"`
	DocumentIDs []uint `json:"documentIds"`
	Style       string `json:"style"`
	RoomType    string `json:"roomType"`
	Description string `json:"description"`
}

// SubmitRenderResponse 提交渲染任务响应
type SubmitRenderResponse struct {
	TaskID       uint   `json:"taskId"`       // 任务ID / Task ID
	Status       string `json:"status"`       // 状态 / Status
	Message      string `json:"message"`      // 提示信息 / Message
	RemainingUse int    `json:"remainingUse"` // 今日剩余次数 / Remaining uses today
}

// SubmitRenderTask 提交异步渲染任务
// Submit async render task with daily limit check
func (s *RenderRecordService) SubmitRenderTask(req *SubmitRenderRequest) (*SubmitRenderResponse, error) {
	// 获取或创建用户配额 / Get or create user quota
	quota, err := s.quotaRepo.GetOrCreate(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("获取用户配额失败: %w", err)
	}

	// 检查今日是否需要重置 / Check if need to reset today
	today := time.Now().Format("2006-01-02")
	lastReset := quota.LastResetDate.Format("2006-01-02")
	if today != lastReset {
		quota.UsedToday = 0
	}

	// 检查是否超过每日限制 / Check if exceeded daily limit
	if quota.UsedToday >= quota.DailyLimit {
		return nil, fmt.Errorf("今日生成次数已用完（%d/%d），请明天再试", quota.UsedToday, quota.DailyLimit)
	}

	// 创建待处理记录 / Create pending record
	record := &model.RenderRecord{
		UserID:    req.UserID,
		ProjectID: req.ProjectID,
		CadFileID: req.CadFileID,
		Status:    "pending",
		Style:     req.Style,
		RoomType:  req.RoomType,
	}

	if err := s.recordRepo.Create(record); err != nil {
		return nil, fmt.Errorf("创建任务失败: %w", err)
	}

	// 增加使用次数 / Increment usage
	if err := s.quotaRepo.IncrementUsage(req.UserID); err != nil {
		fmt.Printf("更新用户配额失败: %v\n", err)
	}

	// 计算剩余次数 / Calculate remaining uses
	remaining := quota.DailyLimit - quota.UsedToday - 1
	if remaining < 0 {
		remaining = 0
	}

	// 异步执行渲染任务 / Execute render task asynchronously
	go s.executeRenderTask(record.ID, req)

	return &SubmitRenderResponse{
		TaskID:       record.ID,
		Status:       "pending",
		Message:      "任务已提交，正在生成中，请稍后在\"我的效果图\"中查看结果",
		RemainingUse: remaining,
	}, nil
}

// executeRenderTask 执行渲染任务（异步）
// Execute render task asynchronously
func (s *RenderRecordService) executeRenderTask(taskID uint, req *SubmitRenderRequest) {
	// 更新状态为处理中 / Update status to processing
	s.recordRepo.UpdateStatus(taskID, "processing", map[string]interface{}{})

	// 调用渲染服务生成效果图 / Call render service to generate image
	result, err := s.renderService.GenerateRender(&GenerateRenderRequest{
		CadFileID:   req.CadFileID,
		ProjectID:   req.ProjectID,
		DocumentIDs: req.DocumentIDs,
		Style:       req.Style,
		RoomType:    req.RoomType,
		Description: req.Description,
	})

	if err != nil {
		// 更新状态为失败 / Update status to failed
		s.recordRepo.UpdateStatus(taskID, "failed", map[string]interface{}{
			"error_message": err.Error(),
		})
		// 发送失败通知 / Send failure notification
		s.sendNotification(req.UserID, "效果图生成失败", fmt.Sprintf("您的效果图生成任务失败：%s", err.Error()), "notice", taskID)
		return
	}

	// 更新状态为完成 / Update status to completed
	s.recordRepo.UpdateStatus(taskID, "completed", map[string]interface{}{
		"image_url":       result.ImageURL,
		"design_proposal": result.DesignProposal,
		"prompt":          result.Prompt,
	})

	// 发送成功通知 / Send success notification
	title := "效果图生成完成"
	content := fmt.Sprintf("您的%s风格效果图已生成完成，请在\"我的效果图\"中查看", req.Style)
	if req.Style == "" {
		content = "您的效果图已生成完成，请在\"我的效果图\"中查看"
	}
	s.sendNotification(req.UserID, title, content, "notice", taskID)
}

// sendNotification 发送通知
// Send notification to user
func (s *RenderRecordService) sendNotification(userID int64, title, content, noticeType string, relatedID uint) {
	if s.notificationRepo == nil {
		return
	}
	notification := &model.Notification{
		UserID:    userID,
		Title:     title,
		Content:   content,
		Type:      noticeType,
		RelatedID: relatedID,
		IsRead:    false,
	}
	if err := s.notificationRepo.Create(notification); err != nil {
		fmt.Printf("发送通知失败: %v\n", err)
	}
}

// GetUserQuota 获取用户配额信息
func (s *RenderRecordService) GetUserQuota(userID int64) (*model.UserRenderQuota, error) {
	return s.quotaRepo.GetOrCreate(userID)
}

// UpdateUserQuota 更新用户每日限制（管理员使用）
func (s *RenderRecordService) UpdateUserQuota(userID int64, dailyLimit int) error {
	_, err := s.quotaRepo.GetOrCreate(userID)
	if err != nil {
		return err
	}
	return s.quotaRepo.UpdateDailyLimit(userID, dailyLimit)
}

// ResetUserQuota 重置用户今日使用次数（管理员使用）
func (s *RenderRecordService) ResetUserQuota(userID int64) error {
	return s.quotaRepo.ResetUsage(userID)
}

// ListRenderRecords 获取效果图记录列表
// List render records with optional cadFileId filter
func (s *RenderRecordService) ListRenderRecords(page, pageSize int, userID int64, projectID uint, cadFileID uint) ([]model.RenderRecord, int64, error) {
	filter := repository.RenderRecordFilter{
		UserID:    userID,
		ProjectID: projectID,
		CadFileID: cadFileID,
	}
	return s.recordRepo.List(page, pageSize, filter)
}

// GetRenderRecord 获取单个效果图记录
func (s *RenderRecordService) GetRenderRecord(id uint) (*model.RenderRecord, error) {
	return s.recordRepo.GetByID(id)
}

// DeleteRenderRecord 删除效果图记录
func (s *RenderRecordService) DeleteRenderRecord(id uint, userID int64) error {
	record, err := s.recordRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("记录不存在")
	}
	if record.UserID != userID {
		return fmt.Errorf("无权删除此记录")
	}
	return s.recordRepo.Delete(id)
}
