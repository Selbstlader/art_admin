package project

import (
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/dto/response"
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/repository"
	"errors"
	"time"
)

type TaskService struct {
	taskRepo *repository.TaskRepository
}

func NewTaskService() *TaskService {
	return &TaskService{
		taskRepo: repository.NewTaskRepository(),
	}
}

// GetTaskList 获取任务列表
func (s *TaskService) GetTaskList(req *request.GetTaskListRequest) ([]response.TaskResponse, int64, error) {
	tasks, total, err := s.taskRepo.FindWithPagination(req)
	if err != nil {
		return nil, 0, err
	}

	result := make([]response.TaskResponse, 0, len(tasks))
	for _, t := range tasks {
		result = append(result, response.TaskResponse{
			ID:              t.ID,
			ProjectID:       t.ProjectID,
			ParentID:        t.ParentID,
			TemplateTaskID:  t.TemplateTaskID,
			Name:            t.Name,
			Description:     t.Description,
			Type:            t.Type,
			Priority:        t.Priority,
			Status:          t.Status,
			Progress:        t.Progress,
			StartDate:       t.StartDate,
			EndDate:         t.EndDate,
			ActualStartDate: t.ActualStartDate,
			ActualEndDate:   t.ActualEndDate,
			AssigneeID:      t.AssigneeID,
			AssigneeName:    t.AssigneeName,
			EstimatedHours:  t.EstimatedHours,
			ActualHours:     t.ActualHours,
			IsOverdue:       t.IsOverdue,
			Sort:            t.Sort,
			Remark:          t.Remark,
			CreatorID:       t.CreatorID,
			CreatorName:     t.CreatorName,
			CreatedAt:       t.CreatedAt,
			UpdatedAt:       t.UpdatedAt,
		})
	}

	return result, total, nil
}

// GetTaskDetail 获取任务详情
func (s *TaskService) GetTaskDetail(id int64) (*response.TaskDetailResponse, error) {
	task, err := s.taskRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 获取子任务
	children, err := s.taskRepo.FindByParentID(id)
	if err != nil {
		return nil, err
	}

	// 获取依赖关系
	dependencies, err := s.taskRepo.FindDependenciesByTaskID(id)
	if err != nil {
		return nil, err
	}

	// 获取评论
	comments, err := s.taskRepo.FindCommentsByTaskID(id)
	if err != nil {
		return nil, err
	}

	// 获取附件
	attachments, err := s.taskRepo.FindAttachmentsByTaskID(id)
	if err != nil {
		return nil, err
	}

	// 转换响应
	childrenResponses := make([]response.TaskResponse, 0, len(children))
	for _, c := range children {
		childrenResponses = append(childrenResponses, response.TaskResponse{
			ID:           c.ID,
			ProjectID:    c.ProjectID,
			ParentID:     c.ParentID,
			Name:         c.Name,
			Type:         c.Type,
			Priority:     c.Priority,
			Status:       c.Status,
			Progress:     c.Progress,
			StartDate:    c.StartDate,
			EndDate:      c.EndDate,
			AssigneeID:   c.AssigneeID,
			AssigneeName: c.AssigneeName,
			IsOverdue:    c.IsOverdue,
		})
	}

	depResponses := make([]response.TaskDependencyResponse, 0, len(dependencies))
	for _, d := range dependencies {
		depResponses = append(depResponses, response.TaskDependencyResponse{
			ID:             d.ID,
			ProjectID:      d.ProjectID,
			PredecessorID:  d.PredecessorID,
			SuccessorID:    d.SuccessorID,
			DependencyType: d.DependencyType,
			LagDays:        d.LagDays,
			CreatedAt:      d.CreatedAt,
		})
	}

	commentResponses := make([]response.TaskCommentResponse, 0, len(comments))
	for _, c := range comments {
		commentResponses = append(commentResponses, response.TaskCommentResponse{
			ID:        c.ID,
			TaskID:    c.TaskID,
			UserID:    c.UserID,
			UserName:  c.UserName,
			Content:   c.Content,
			CreatedAt: c.CreatedAt,
		})
	}

	attachmentResponses := make([]response.TaskAttachmentResponse, 0, len(attachments))
	for _, a := range attachments {
		attachmentResponses = append(attachmentResponses, response.TaskAttachmentResponse{
			ID:           a.ID,
			TaskID:       a.TaskID,
			FileName:     a.FileName,
			FileSize:     a.FileSize,
			FileType:     a.FileType,
			FilePath:     a.FilePath,
			UploaderID:   a.UploaderID,
			UploaderName: a.UploaderName,
			CreatedAt:    a.CreatedAt,
		})
	}

	result := &response.TaskDetailResponse{
		TaskResponse: response.TaskResponse{
			ID:              task.ID,
			ProjectID:       task.ProjectID,
			ParentID:        task.ParentID,
			TemplateTaskID:  task.TemplateTaskID,
			Name:            task.Name,
			Description:     task.Description,
			Type:            task.Type,
			Priority:        task.Priority,
			Status:          task.Status,
			Progress:        task.Progress,
			StartDate:       task.StartDate,
			EndDate:         task.EndDate,
			ActualStartDate: task.ActualStartDate,
			ActualEndDate:   task.ActualEndDate,
			AssigneeID:      task.AssigneeID,
			AssigneeName:    task.AssigneeName,
			EstimatedHours:  task.EstimatedHours,
			ActualHours:     task.ActualHours,
			IsOverdue:       task.IsOverdue,
			Sort:            task.Sort,
			Remark:          task.Remark,
			CreatorID:       task.CreatorID,
			CreatorName:     task.CreatorName,
			CreatedAt:       task.CreatedAt,
			UpdatedAt:       task.UpdatedAt,
		},
		Children:     childrenResponses,
		Dependencies: depResponses,
		Comments:     commentResponses,
		Attachments:  attachmentResponses,
	}

	return result, nil
}

// CreateTask 创建任务
func (s *TaskService) CreateTask(req *request.CreateTaskRequest, creatorID int64, creatorName string) error {
	task := &model.Task{
		ProjectID:      req.ProjectID,
		ParentID:       req.ParentID,
		Name:           req.Name,
		Description:    req.Description,
		Type:           req.Type,
		Priority:       req.Priority,
		Status:         1, // 未开始
		Progress:       0,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		AssigneeID:     req.AssigneeID,
		EstimatedHours: req.EstimatedHours,
		Sort:           req.Sort,
		Remark:         req.Remark,
		CreatorID:      creatorID,
		CreatorName:    creatorName,
	}

	// 检查是否超期
	if time.Now().After(task.EndDate) {
		task.IsOverdue = true
		task.Status = 4 // 已延期
	}

	return s.taskRepo.Create(task)
}

// UpdateTask 更新任务
func (s *TaskService) UpdateTask(req *request.UpdateTaskRequest) error {
	task, err := s.taskRepo.FindByID(req.ID)
	if err != nil {
		return err
	}

	task.Name = req.Name
	task.Description = req.Description
	task.Type = req.Type
	task.Priority = req.Priority
	task.Status = req.Status
	task.Progress = req.Progress
	task.StartDate = req.StartDate
	task.EndDate = req.EndDate
	task.ActualStartDate = req.ActualStartDate
	task.ActualEndDate = req.ActualEndDate
	task.AssigneeID = req.AssigneeID
	task.EstimatedHours = req.EstimatedHours
	task.ActualHours = req.ActualHours
	task.Sort = req.Sort
	task.Remark = req.Remark

	// 检查是否超期
	if time.Now().After(task.EndDate) && task.Status != 3 {
		task.IsOverdue = true
		if task.Status == 2 { // 如果是进行中,标记为已延期
			task.Status = 4
		}
	}

	return s.taskRepo.Update(task)
}

// BatchUpdateTask 批量更新任务(甘特图拖拽)
func (s *TaskService) BatchUpdateTask(req *request.BatchUpdateTaskRequest) error {
	for _, item := range req.Tasks {
		task, err := s.taskRepo.FindByID(item.ID)
		if err != nil {
			return err
		}

		task.StartDate = item.StartDate
		task.EndDate = item.EndDate
		task.Progress = item.Progress

		// 检查是否超期
		if time.Now().After(task.EndDate) && task.Status != 3 {
			task.IsOverdue = true
			if task.Status == 2 {
				task.Status = 4
			}
		}

		if err := s.taskRepo.Update(task); err != nil {
			return err
		}
	}

	return nil
}

// QuickUpdateTask 快速更新任务(甘特图实时更新)
func (s *TaskService) QuickUpdateTask(req *request.QuickUpdateTaskRequest) error {
	task, err := s.taskRepo.FindByID(req.ID)
	if err != nil {
		return err
	}

	// 只更新提供的字段
	if req.StartDate != nil {
		task.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		task.EndDate = *req.EndDate
	}
	if req.Progress != nil {
		task.Progress = *req.Progress
	}
	if req.Duration != nil && req.StartDate != nil {
		// 根据持续时间计算结束日期
		task.EndDate = req.StartDate.AddDate(0, 0, *req.Duration)
	}

	// 检查是否超期
	if time.Now().After(task.EndDate) && task.Status != 3 {
		task.IsOverdue = true
		if task.Status == 2 {
			task.Status = 4
		}
	} else {
		task.IsOverdue = false
	}

	return s.taskRepo.Update(task)
}

// DeleteTask 删除任务
func (s *TaskService) DeleteTask(id int64) error {
	return s.taskRepo.Delete(id)
}

// GetTaskGanttData 获取甘特图数据
func (s *TaskService) GetTaskGanttData(projectID int64) (*response.TaskGanttDataResponse, error) {
	tasks, err := s.taskRepo.FindByProjectID(projectID)
	if err != nil {
		return nil, err
	}

	dependencies, err := s.taskRepo.FindDependenciesByProjectID(projectID)
	if err != nil {
		return nil, err
	}

	// 转换为甘特图格式
	ganttTasks := make([]response.TaskGanttItem, 0, len(tasks))
	for _, t := range tasks {
		duration := int(t.EndDate.Sub(t.StartDate).Hours() / 24)
		ganttTasks = append(ganttTasks, response.TaskGanttItem{
			ID:           t.ID,
			Text:         t.Name,
			StartDate:    t.StartDate,
			EndDate:      t.EndDate,
			Duration:     duration,
			Progress:     float64(t.Progress) / 100.0,
			ParentID:     t.ParentID,
			Type:         t.Type,
			Priority:     t.Priority,
			Status:       t.Status,
			AssigneeID:   t.AssigneeID,
			AssigneeName: t.AssigneeName,
			IsOverdue:    t.IsOverdue,
			Open:         true,
		})
	}

	depResponses := make([]response.TaskDependencyResponse, 0, len(dependencies))
	for _, d := range dependencies {
		depResponses = append(depResponses, response.TaskDependencyResponse{
			ID:             d.ID,
			ProjectID:      d.ProjectID,
			PredecessorID:  d.PredecessorID,
			SuccessorID:    d.SuccessorID,
			DependencyType: d.DependencyType,
			LagDays:        d.LagDays,
			CreatedAt:      d.CreatedAt,
		})
	}

	return &response.TaskGanttDataResponse{
		Tasks:        ganttTasks,
		Dependencies: depResponses,
	}, nil
}

// CreateDependency 创建依赖关系
func (s *TaskService) CreateDependency(req *request.CreateTaskDependencyRequest) error {
	if req.PredecessorID == req.SuccessorID {
		return errors.New("任务不能依赖自己")
	}

	dependency := &model.TaskDependency{
		ProjectID:      req.ProjectID,
		PredecessorID:  req.PredecessorID,
		SuccessorID:    req.SuccessorID,
		DependencyType: req.DependencyType,
		LagDays:        req.LagDays,
	}

	return s.taskRepo.CreateDependency(dependency)
}

// DeleteDependency 删除依赖关系
func (s *TaskService) DeleteDependency(id int64) error {
	return s.taskRepo.DeleteDependency(id)
}

// CreateComment 创建评论
func (s *TaskService) CreateComment(req *request.CreateTaskCommentRequest, userID int64, userName string) error {
	comment := &model.TaskComment{
		TaskID:   req.TaskID,
		UserID:   userID,
		UserName: userName,
		Content:  req.Content,
	}

	return s.taskRepo.CreateComment(comment)
}
