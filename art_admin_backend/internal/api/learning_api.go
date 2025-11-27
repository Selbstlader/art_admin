package api

import (
	"net/http"
	"strconv"
	"strings"

	"art_admin_backend/internal/dto"
	learningSvc "art_admin_backend/internal/service/learning"

	"github.com/gin-gonic/gin"
)

// LearningAPI 学习系统 API
type LearningAPI struct {
	materialService *learningSvc.LearningMaterialService
	subjectService  *learningSvc.SubjectService
}

// NewLearningAPI 创建学习系统 API
func NewLearningAPI(
	materialService *learningSvc.LearningMaterialService,
	subjectService *learningSvc.SubjectService,
) *LearningAPI {
	return &LearningAPI{
		materialService: materialService,
		subjectService:  subjectService,
	}
}

// GenerateMaterial 生成教材（异步）
// @Summary 生成教材（异步）
// @Tags Learning
// @Accept json
// @Produce json
// @Param request body dto.GenerateMaterialRequest true "生成教材请求"
// @Success 200 {object} dto.TaskResponse
// @Router /api/learning/material/generate [post]
func (a *LearningAPI) GenerateMaterial(c *gin.Context) {
	var req dto.GenerateMaterialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	// 从上下文获取用户 ID
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未登录"})
		return
	}

	// 构建服务请求
	serviceReq := &learningSvc.GenerateMaterialRequest{
		UserID:     userID.(int64),
		SubjectID:  req.SubjectID,
		Grade:      req.Grade,
		Topic:      req.Topic,
		Difficulty: req.Difficulty,
	}

	// 创建异步生成任务
	task, err := a.materialService.CreateGenerateTask(c.Request.Context(), serviceReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建任务失败: " + err.Error()})
		return
	}

	// 构建响应
	resp := &dto.TaskResponse{
		ID:         task.ID,
		UserID:     task.UserID,
		SubjectID:  task.SubjectID,
		Grade:      task.Grade,
		Topic:      task.Topic,
		Difficulty: task.Difficulty,
		Status:     task.Status,
		MaterialID: task.MaterialID,
		ErrorMsg:   task.ErrorMsg,
		CreatedAt:  task.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  task.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "任务创建成功", "data": resp})
}

// GetMaterial 获取教材详情
// @Summary 获取教材详情
// @Tags Learning
// @Produce json
// @Param id path int true "教材ID"
// @Success 200 {object} dto.MaterialResponse
// @Router /api/learning/material/{id} [get]
func (a *LearningAPI) GetMaterial(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	material, err := a.materialService.GetMaterialByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "教材不存在"})
		return
	}

	resp := &dto.MaterialResponse{
		ID:            material.ID,
		UserID:        material.UserID,
		SubjectID:     material.SubjectID,
		Title:         material.Title,
		Topic:         material.Topic,
		Grade:         material.Grade,
		Difficulty:    material.Difficulty,
		Summary:       material.Summary,
		Content:       material.Content,
		AudioURL:      material.AudioURL,
		AudioDuration: material.AudioDuration,
		TotalTime:     material.TotalTime,
		ViewCount:     material.ViewCount,
		FavoriteCount: material.FavoriteCount,
		CreatedAt:     material.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     material.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": resp})
}

// ListMaterials 获取教材列表
// @Summary 获取教材列表
// @Tags Learning
// @Produce json
// @Param subject_id query int false "学科ID"
// @Param page query int true "页码"
// @Param limit query int true "每页数量"
// @Success 200 {object} dto.MaterialListResponse
// @Router /api/learning/material/list [get]
func (a *LearningAPI) ListMaterials(c *gin.Context) {
	var req dto.MaterialListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	// 从上下文获取用户 ID
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未登录"})
		return
	}

	materials, total, err := a.materialService.ListMaterials(
		c.Request.Context(),
		userID.(int64),
		req.SubjectID,
		req.Grade,
		req.Page,
		req.Limit,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取列表失败: " + err.Error()})
		return
	}

	// 构建响应
	list := make([]*dto.MaterialResponse, 0, len(materials))
	for _, m := range materials {
		list = append(list, &dto.MaterialResponse{
			ID:            m.ID,
			UserID:        m.UserID,
			SubjectID:     m.SubjectID,
			Title:         m.Title,
			Topic:         m.Topic,
			Grade:         m.Grade,
			Difficulty:    m.Difficulty,
			Summary:       m.Summary,
			Content:       m.Content,
			AudioURL:      m.AudioURL,
			AudioDuration: m.AudioDuration,
			TotalTime:     m.TotalTime,
			ViewCount:     m.ViewCount,
			FavoriteCount: m.FavoriteCount,
			CreatedAt:     m.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:     m.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	resp := &dto.MaterialListResponse{
		List:  list,
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": resp})
}

// DeleteMaterial 删除教材
// @Summary 删除教材
// @Tags Learning
// @Produce json
// @Param id path int true "教材ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/learning/material/{id} [delete]
func (a *LearningAPI) DeleteMaterial(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 从上下文获取用户 ID
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未登录"})
		return
	}

	if err := a.materialService.DeleteMaterial(c.Request.Context(), id, userID.(int64)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

// ListSubjects 获取学科列表
// @Summary 获取学科列表
// @Tags Learning
// @Produce json
// @Success 200 {object} []dto.SubjectResponse
// @Router /api/learning/subject/list [get]
func (a *LearningAPI) ListSubjects(c *gin.Context) {
	subjects, err := a.subjectService.ListSubjects(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取学科列表失败: " + err.Error()})
		return
	}

	// 构建响应
	list := make([]*dto.SubjectResponse, 0, len(subjects))
	for _, s := range subjects {
		list = append(list, &dto.SubjectResponse{
			ID:          s.ID,
			Name:        s.Name,
			Icon:        s.Icon,
			Description: s.Description,
			GradeLevels: s.GradeLevels,
			Sort:        s.Sort,
		})
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": list})
}

// GetGradesBySubject 获取某学科下的年级列表
// @Summary 获取年级列表
// @Tags Learning
// @Produce json
// @Param subject_id query int false "学科ID"
// @Success 200 {object} dto.GradeListResponse
// @Router /api/learning/grades [get]
func (a *LearningAPI) GetGradesBySubject(c *gin.Context) {
	subjectIDStr := c.Query("subject_id")
	var subjectID int64
	if subjectIDStr != "" {
		id, err := strconv.ParseInt(subjectIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "学科ID格式错误"})
			return
		}
		subjectID = id
	}

	// 从上下文获取用户 ID
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未登录"})
		return
	}

	gradeMap, err := a.materialService.GetGradesBySubject(c.Request.Context(), userID.(int64), subjectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取年级列表失败: " + err.Error()})
		return
	}

	// 构建响应
	grades := make([]*dto.GradeInfo, 0, len(gradeMap))
	for grade, count := range gradeMap {
		grades = append(grades, &dto.GradeInfo{
			Grade: grade,
			Count: count,
		})
	}

	resp := &dto.GradeListResponse{
		SubjectID: subjectID,
		Grades:    grades,
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": resp})
}

// GetTask 获取任务详情
// @Summary 获取任务详情
// @Tags Learning
// @Produce json
// @Param id path int true "任务ID"
// @Success 200 {object} dto.TaskResponse
// @Router /api/learning/material/task/{id} [get]
func (a *LearningAPI) GetTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	task, err := a.materialService.GetTaskByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取任务失败: " + err.Error()})
		return
	}

	// 构建响应
	resp := &dto.TaskResponse{
		ID:         task.ID,
		UserID:     task.UserID,
		SubjectID:  task.SubjectID,
		Grade:      task.Grade,
		Topic:      task.Topic,
		Difficulty: task.Difficulty,
		Status:     task.Status,
		MaterialID: task.MaterialID,
		ErrorMsg:   task.ErrorMsg,
		CreatedAt:  task.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  task.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": resp})
}

// ListTasks 获取任务列表（默认只返回进行中的任务）
// @Summary 获取任务列表
// @Tags Learning
// @Produce json
// @Param status query string false "任务状态,多个用逗号分隔"
// @Success 200 {object} []dto.TaskResponse
// @Router /api/learning/material/tasks [get]
func (a *LearningAPI) ListTasks(c *gin.Context) {
	// 从上下文获取用户 ID
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未登录"})
		return
	}

	// 解析状态参数，默认只返回待处理和生成中的任务
	var statuses []int8
	statusStr := c.Query("status")
	if statusStr != "" {
		statusParts := strings.Split(statusStr, ",")
		for _, part := range statusParts {
			status, err := strconv.ParseInt(strings.TrimSpace(part), 10, 8)
			if err == nil {
				statuses = append(statuses, int8(status))
			}
		}
	} else {
		// 默认只返回待处理(0)和生成中(1)的任务
		statuses = []int8{0, 1}
	}

	tasks, err := a.materialService.ListTasksByUserID(c.Request.Context(), userID.(int64), statuses)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取任务列表失败: " + err.Error()})
		return
	}

	// 构建响应
	list := make([]*dto.TaskResponse, 0, len(tasks))
	for _, task := range tasks {
		list = append(list, &dto.TaskResponse{
			ID:         task.ID,
			UserID:     task.UserID,
			SubjectID:  task.SubjectID,
			Grade:      task.Grade,
			Topic:      task.Topic,
			Difficulty: task.Difficulty,
			Status:     task.Status,
			MaterialID: task.MaterialID,
			ErrorMsg:   task.ErrorMsg,
			CreatedAt:  task.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:  task.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": list})
}
