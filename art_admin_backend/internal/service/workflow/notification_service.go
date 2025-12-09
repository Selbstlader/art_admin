package workflow

import (
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/websocket"
	"context"
	"encoding/json"
	"log"
	"time"
)

// ============================================================================
// Notification Types
// ============================================================================

// NotificationType 通知类型
type NotificationType string

const (
	// NotificationTypeNewTask 新待办任务通知
	NotificationTypeNewTask NotificationType = "workflow_new_task"
	// NotificationTypeTaskReminder 任务提醒通知
	NotificationTypeTaskReminder NotificationType = "workflow_task_reminder"
	// NotificationTypeProcessStatusChange 流程状态变更通知
	NotificationTypeProcessStatusChange NotificationType = "workflow_process_status"
	// NotificationTypeTaskTimeout 任务超时通知
	NotificationTypeTaskTimeout NotificationType = "workflow_task_timeout"
	// NotificationTypeProcessCompleted 流程完成通知
	NotificationTypeProcessCompleted NotificationType = "workflow_process_completed"
	// NotificationTypeProcessRejected 流程拒绝通知
	NotificationTypeProcessRejected NotificationType = "workflow_process_rejected"
	// NotificationTypeProcessWithdrawn 流程撤回通知
	NotificationTypeProcessWithdrawn NotificationType = "workflow_process_withdrawn"
)

// WorkflowNotification 工作流通知消息结构
type WorkflowNotification struct {
	Type      NotificationType       `json:"type"`
	Title     string                 `json:"title"`
	Content   string                 `json:"content"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
}

// ============================================================================
// Interface
// ============================================================================

// NotificationService 通知服务接口
// Requirements: 6.1 - WHEN 新Task创建 THEN Notification_Service SHALL 通过WebSocket向Assignee推送待办提醒
// Requirements: 6.2 - WHEN Process_Instance状态变更 THEN Notification_Service SHALL 通知发起人审批进度更新
type NotificationService interface {
	// NotifyNewTask 发送新待办任务通知
	// Requirements: 6.1 - 新任务创建时推送通知
	NotifyNewTask(ctx context.Context, task *model.WfTask, instance *model.ProcessInstance) error

	// NotifyNewTasks 批量发送新待办任务通知
	// Requirements: 6.1 - 新任务创建时推送通知
	NotifyNewTasks(ctx context.Context, tasks []*model.WfTask, instance *model.ProcessInstance) error

	// NotifyProcessStatusChange 发送流程状态变更通知
	// Requirements: 6.2 - 状态变更时通知发起人
	NotifyProcessStatusChange(ctx context.Context, instance *model.ProcessInstance, oldStatus, newStatus string) error

	// NotifyProcessCompleted 发送流程完成通知
	// Requirements: 6.4 - 审批完成时通知发起人最终审批结果
	NotifyProcessCompleted(ctx context.Context, instance *model.ProcessInstance) error

	// NotifyProcessRejected 发送流程拒绝通知
	// Requirements: 6.4 - 审批完成时通知发起人最终审批结果
	NotifyProcessRejected(ctx context.Context, instance *model.ProcessInstance, rejectReason string) error

	// NotifyTaskTimeout 发送任务超时提醒
	// Requirements: 6.3 - 任务即将超时时提前向Assignee发送催办提醒
	NotifyTaskTimeout(ctx context.Context, task *model.WfTask, instance *model.ProcessInstance) error
}

// ============================================================================
// Implementation
// ============================================================================

// notificationService 通知服务实现
type notificationService struct {
	hub *websocket.Hub
}

// NewNotificationService 创建通知服务实例
// 复用现有 WebSocket Hub
func NewNotificationService(hub *websocket.Hub) NotificationService {
	return &notificationService{
		hub: hub,
	}
}

// NotifyNewTask 发送新待办任务通知
// Requirements: 6.1 - WHEN 新Task创建 THEN Notification_Service SHALL 通过WebSocket向Assignee推送待办提醒
func (s *notificationService) NotifyNewTask(ctx context.Context, task *model.WfTask, instance *model.ProcessInstance) error {
	if s.hub == nil {
		log.Printf("[NotificationService] WebSocket Hub not available, skipping notification for task %d", task.ID)
		return nil
	}

	notification := WorkflowNotification{
		Type:    NotificationTypeNewTask,
		Title:   "新的待办任务",
		Content: buildNewTaskContent(task, instance),
		Data: map[string]interface{}{
			"taskId":        task.ID,
			"taskNodeId":    task.NodeID,
			"taskNodeName":  task.NodeName,
			"processInstId": instance.ID,
			"processTitle":  instance.Title,
			"initiatorId":   instance.InitiatorID,
			"initiatorName": instance.InitiatorName,
			"dueAt":         task.DueAt,
		},
		Timestamp: time.Now(),
	}

	return s.sendToUser(uint(task.AssigneeID), notification)
}

// NotifyNewTasks 批量发送新待办任务通知
func (s *notificationService) NotifyNewTasks(ctx context.Context, tasks []*model.WfTask, instance *model.ProcessInstance) error {
	for _, task := range tasks {
		if err := s.NotifyNewTask(ctx, task, instance); err != nil {
			log.Printf("[NotificationService] Failed to notify task %d: %v", task.ID, err)
			// 继续处理其他任务，不中断
		}
	}
	return nil
}

// NotifyProcessStatusChange 发送流程状态变更通知
// Requirements: 6.2 - WHEN Process_Instance状态变更 THEN Notification_Service SHALL 通知发起人审批进度更新
func (s *notificationService) NotifyProcessStatusChange(ctx context.Context, instance *model.ProcessInstance, oldStatus, newStatus string) error {
	if s.hub == nil {
		log.Printf("[NotificationService] WebSocket Hub not available, skipping status change notification for instance %d", instance.ID)
		return nil
	}

	notification := WorkflowNotification{
		Type:    NotificationTypeProcessStatusChange,
		Title:   "流程状态更新",
		Content: buildStatusChangeContent(instance, oldStatus, newStatus),
		Data: map[string]interface{}{
			"processInstId": instance.ID,
			"processTitle":  instance.Title,
			"oldStatus":     oldStatus,
			"newStatus":     newStatus,
			"currentNodeId": instance.CurrentNodeID,
		},
		Timestamp: time.Now(),
	}

	return s.sendToUser(uint(instance.InitiatorID), notification)
}

// NotifyProcessCompleted 发送流程完成通知
// Requirements: 6.4 - WHEN 审批完成 THEN Notification_Service SHALL 通知发起人最终审批结果
func (s *notificationService) NotifyProcessCompleted(ctx context.Context, instance *model.ProcessInstance) error {
	if s.hub == nil {
		log.Printf("[NotificationService] WebSocket Hub not available, skipping completion notification for instance %d", instance.ID)
		return nil
	}

	notification := WorkflowNotification{
		Type:    NotificationTypeProcessCompleted,
		Title:   "审批已通过",
		Content: buildProcessCompletedContent(instance),
		Data: map[string]interface{}{
			"processInstId": instance.ID,
			"processTitle":  instance.Title,
			"status":        instance.Status,
			"completedAt":   instance.CompletedAt,
		},
		Timestamp: time.Now(),
	}

	return s.sendToUser(uint(instance.InitiatorID), notification)
}

// NotifyProcessRejected 发送流程拒绝通知
// Requirements: 6.4 - WHEN 审批完成 THEN Notification_Service SHALL 通知发起人最终审批结果
func (s *notificationService) NotifyProcessRejected(ctx context.Context, instance *model.ProcessInstance, rejectReason string) error {
	if s.hub == nil {
		log.Printf("[NotificationService] WebSocket Hub not available, skipping rejection notification for instance %d", instance.ID)
		return nil
	}

	notification := WorkflowNotification{
		Type:    NotificationTypeProcessRejected,
		Title:   "审批被拒绝",
		Content: buildProcessRejectedContent(instance, rejectReason),
		Data: map[string]interface{}{
			"processInstId": instance.ID,
			"processTitle":  instance.Title,
			"status":        instance.Status,
			"rejectReason":  rejectReason,
			"completedAt":   instance.CompletedAt,
		},
		Timestamp: time.Now(),
	}

	return s.sendToUser(uint(instance.InitiatorID), notification)
}

// NotifyTaskTimeout 发送任务超时提醒
// Requirements: 6.3 - WHEN 任务即将超时 THEN Notification_Service SHALL 提前向Assignee发送催办提醒
func (s *notificationService) NotifyTaskTimeout(ctx context.Context, task *model.WfTask, instance *model.ProcessInstance) error {
	if s.hub == nil {
		log.Printf("[NotificationService] WebSocket Hub not available, skipping timeout notification for task %d", task.ID)
		return nil
	}

	notification := WorkflowNotification{
		Type:    NotificationTypeTaskTimeout,
		Title:   "任务即将超时",
		Content: buildTaskTimeoutContent(task, instance),
		Data: map[string]interface{}{
			"taskId":        task.ID,
			"taskNodeId":    task.NodeID,
			"taskNodeName":  task.NodeName,
			"processInstId": instance.ID,
			"processTitle":  instance.Title,
			"dueAt":         task.DueAt,
		},
		Timestamp: time.Now(),
	}

	return s.sendToUser(uint(task.AssigneeID), notification)
}

// ============================================================================
// Helper Methods
// ============================================================================

// sendToUser 向指定用户发送通知
func (s *notificationService) sendToUser(userID uint, notification WorkflowNotification) error {
	message, err := json.Marshal(notification)
	if err != nil {
		log.Printf("[NotificationService] Failed to marshal notification: %v", err)
		return err
	}

	s.hub.SendToUser(userID, message)
	log.Printf("[NotificationService] Sent %s notification to user %d", notification.Type, userID)
	return nil
}

// buildNewTaskContent 构建新任务通知内容
func buildNewTaskContent(task *model.WfTask, instance *model.ProcessInstance) string {
	content := "您有一个新的待办任务需要处理"
	if instance != nil && instance.Title != "" {
		content = "您有一个新的待办任务：" + instance.Title
	}
	if task.NodeName != "" {
		content += "（" + task.NodeName + "）"
	}
	return content
}

// buildStatusChangeContent 构建状态变更通知内容
func buildStatusChangeContent(instance *model.ProcessInstance, oldStatus, newStatus string) string {
	statusMap := map[string]string{
		model.ProcessInstStatusRunning:   "进行中",
		model.ProcessInstStatusCompleted: "已完成",
		model.ProcessInstStatusRejected:  "已拒绝",
		model.ProcessInstStatusWithdrawn: "已撤回",
	}

	oldStatusText := statusMap[oldStatus]
	if oldStatusText == "" {
		oldStatusText = oldStatus
	}
	newStatusText := statusMap[newStatus]
	if newStatusText == "" {
		newStatusText = newStatus
	}

	return "您发起的流程「" + instance.Title + "」状态已从「" + oldStatusText + "」变更为「" + newStatusText + "」"
}

// buildProcessCompletedContent 构建流程完成通知内容
func buildProcessCompletedContent(instance *model.ProcessInstance) string {
	return "恭喜！您发起的流程「" + instance.Title + "」已审批通过"
}

// buildProcessRejectedContent 构建流程拒绝通知内容
func buildProcessRejectedContent(instance *model.ProcessInstance, rejectReason string) string {
	content := "您发起的流程「" + instance.Title + "」已被拒绝"
	if rejectReason != "" {
		content += "，原因：" + rejectReason
	}
	return content
}

// buildTaskTimeoutContent 构建任务超时通知内容
func buildTaskTimeoutContent(task *model.WfTask, instance *model.ProcessInstance) string {
	content := "您有一个待办任务即将超时，请尽快处理"
	if instance != nil && instance.Title != "" {
		content = "待办任务「" + instance.Title + "」即将超时，请尽快处理"
	}
	return content
}
