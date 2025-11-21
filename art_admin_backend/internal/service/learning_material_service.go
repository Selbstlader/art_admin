package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/deepseek"
	"art_admin_backend/internal/pkg/dify"
	"art_admin_backend/internal/repository"
)

// LearningMaterialService 教材服务
type LearningMaterialService struct {
	repo           *repository.LearningMaterialRepository
	subjectRepo    *repository.SubjectRepository
	taskRepo       *repository.MaterialGenerateTaskRepository
	deepseekClient *deepseek.Client
	difyClient     *dify.KnowledgeClient
	difyDatasetID  string
}

// NewLearningMaterialService 创建教材服务
func NewLearningMaterialService(
	repo *repository.LearningMaterialRepository,
	subjectRepo *repository.SubjectRepository,
	taskRepo *repository.MaterialGenerateTaskRepository,
	deepseekClient *deepseek.Client,
	difyClient *dify.KnowledgeClient,
	difyDatasetID string,
) *LearningMaterialService {
	return &LearningMaterialService{
		repo:           repo,
		subjectRepo:    subjectRepo,
		taskRepo:       taskRepo,
		deepseekClient: deepseekClient,
		difyClient:     difyClient,
		difyDatasetID:  difyDatasetID,
	}
}

// GenerateMaterialRequest 生成教材请求
type GenerateMaterialRequest struct {
	UserID     int64  `json:"user_id" binding:"required"`
	SubjectID  int64  `json:"subject_id" binding:"required"`
	Grade      string `json:"grade" binding:"required"`
	Topic      string `json:"topic" binding:"required"`
	Difficulty int8   `json:"difficulty" binding:"required,min=1,max=3"`
}

// GenerateMaterial 生成教材
func (s *LearningMaterialService) GenerateMaterial(ctx context.Context, req *GenerateMaterialRequest) (*model.LearningMaterial, error) {
	// 1. 获取学科信息
	subject, err := s.subjectRepo.GetByID(ctx, req.SubjectID)
	if err != nil {
		return nil, fmt.Errorf("获取学科信息失败: %w", err)
	}

	// 2. 从 Dify 知识库检索相关资料
	query := fmt.Sprintf("%s %s %s 知识点", subject.Name, req.Grade, req.Topic)
	knowledgeContents, err := s.difyClient.SimpleRetrieve(ctx, s.difyDatasetID, query)
	if err != nil {
		// 如果检索失败，记录日志但继续生成（使用空的知识库内容）
		fmt.Printf("警告: 检索知识库失败: %v\n", err)
		knowledgeContents = []string{}
	}

	// 3. 构建 RAG Prompt
	systemPrompt := s.buildSystemPrompt(subject.Name, req.Grade, req.Topic, req.Difficulty, knowledgeContents)
	userMessage := fmt.Sprintf("请为我生成关于'%s'的教材内容", req.Topic)

	// 4. 调用 DeepSeek API 生成教材
	response, err := s.deepseekClient.SimpleChat(ctx, systemPrompt, userMessage)
	if err != nil {
		return nil, fmt.Errorf("生成教材失败: %w", err)
	}

	// 5. 解析 JSON 响应
	var content model.MaterialContent
	if err := s.parseJSONResponse(response, &content); err != nil {
		return nil, fmt.Errorf("解析教材内容失败: %w", err)
	}

	// 6. 保存教材到数据库
	material := &model.LearningMaterial{
		UserID:     req.UserID,
		SubjectID:  req.SubjectID,
		Title:      content.Title,
		Topic:      req.Topic,
		Grade:      req.Grade,
		Difficulty: req.Difficulty,
		Summary:    content.Summary,
		Content:    content,
		TotalTime:  content.TotalTime,
		Status:     1,
	}

	if err := s.repo.Create(ctx, material); err != nil {
		return nil, fmt.Errorf("保存教材失败: %w", err)
	}

	return material, nil
}

// buildSystemPrompt 构建系统提示词
func (s *LearningMaterialService) buildSystemPrompt(subjectName, grade, topic string, difficulty int8, knowledgeContents []string) string {
	// 难度映射
	difficultyMap := map[int8]string{
		1: "基础",
		2: "进阶",
		3: "高级",
	}
	difficultyStr := difficultyMap[difficulty]

	// 拼接知识库内容
	knowledgeText := ""
	if len(knowledgeContents) > 0 {
		knowledgeText = strings.Join(knowledgeContents, "\n\n")
	} else {
		knowledgeText = "暂无参考资料，请根据专业知识生成内容。"
	}

	prompt := fmt.Sprintf(`你是一位专业的%s老师，擅长为%s学生设计教学内容。

请根据以下参考资料，为学生生成一份关于"%s"的教材内容。

【参考资料】
%s

【要求】
1. 内容符合%s学生的认知水平
2. 结构清晰，包含以下部分：
   - 知识点讲解（清晰易懂）
   - 重点难点标注
   - 例题演示（至少2个）
   - 练习题（至少5个，包含答案和解析）
3. 难度等级：%s
4. 语言生动有趣，激发学习兴趣
5. 必须输出 JSON 格式

【输出格式】
请严格按照以下 JSON 格式输出，不要添加任何其他文字说明：
{
  "title": "教材标题",
  "summary": "内容概要",
  "sections": [
    {
      "type": "knowledge",
      "title": "知识点标题",
      "content": "详细讲解内容",
      "key_points": ["重点1", "重点2"],
      "difficulty": "基础/进阶/高级"
    },
    {
      "type": "example",
      "title": "例题标题",
      "question": "题目内容",
      "solution": "解题步骤",
      "answer": "答案"
    },
    {
      "type": "exercise",
      "title": "练习题",
      "questions": [
        {
          "id": 1,
          "type": "choice",
          "question": "题目内容",
          "options": ["A选项", "B选项", "C选项", "D选项"],
          "answer": "A",
          "explanation": "解析"
        }
      ]
    }
  ],
  "total_time": 45
}`, subjectName, grade, topic, knowledgeText, grade, difficultyStr)

	return prompt
}

// parseJSONResponse 解析 JSON 响应
func (s *LearningMaterialService) parseJSONResponse(response string, content *model.MaterialContent) error {
	// 清理响应文本，提取 JSON 部分
	response = strings.TrimSpace(response)

	// 如果响应包含 markdown 代码块，提取其中的 JSON
	if strings.HasPrefix(response, "```json") {
		// 处理 ```json 开头的情况
		response = strings.TrimPrefix(response, "```json")
		if idx := strings.LastIndex(response, "```"); idx > 0 {
			response = response[:idx]
		}
	} else if strings.HasPrefix(response, "```") {
		// 处理 ``` 开头的情况
		response = strings.TrimPrefix(response, "```")
		if idx := strings.LastIndex(response, "```"); idx > 0 {
			response = response[:idx]
		}
	}

	response = strings.TrimSpace(response)

	// 解析 JSON
	if err := json.Unmarshal([]byte(response), content); err != nil {
		return fmt.Errorf("JSON 解析失败: %w", err)
	}

	return nil
}

// GetMaterialByID 根据 ID 获取教材
func (s *LearningMaterialService) GetMaterialByID(ctx context.Context, id int64) (*model.LearningMaterial, error) {
	material, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取教材失败: %w", err)
	}

	// 增加浏览次数
	_ = s.repo.IncrementViewCount(ctx, id)

	return material, nil
}

// ListMaterials 获取教材列表
func (s *LearningMaterialService) ListMaterials(ctx context.Context, userID, subjectID int64, grade string, page, limit int) ([]*model.LearningMaterial, int64, error) {
	offset := (page - 1) * limit
	materials, total, err := s.repo.List(ctx, userID, subjectID, grade, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("获取教材列表失败: %w", err)
	}

	return materials, total, nil
}

// GetGradesBySubject 获取某学科下的年级列表
func (s *LearningMaterialService) GetGradesBySubject(ctx context.Context, userID, subjectID int64) (map[string]int64, error) {
	return s.repo.GetGradesBySubject(ctx, userID, subjectID)
}

// DeleteMaterial 删除教材（软删除）
func (s *LearningMaterialService) DeleteMaterial(ctx context.Context, id, userID int64) error {
	// 检查教材是否存在且属于当前用户
	material, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("获取教材失败: %w", err)
	}

	if material.UserID != userID {
		return fmt.Errorf("无权删除该教材")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("删除教材失败: %w", err)
	}

	return nil
}

// CreateGenerateTask 创建生成任务（异步）
func (s *LearningMaterialService) CreateGenerateTask(ctx context.Context, req *GenerateMaterialRequest) (*model.MaterialGenerateTask, error) {
	// 创建任务记录
	task := &model.MaterialGenerateTask{
		UserID:     req.UserID,
		SubjectID:  req.SubjectID,
		Grade:      req.Grade,
		Topic:      req.Topic,
		Difficulty: req.Difficulty,
		Status:     model.TaskStatusPending,
	}

	if err := s.taskRepo.Create(ctx, task); err != nil {
		return nil, fmt.Errorf("创建任务失败: %w", err)
	}

	// 启动异步生成
	go s.processGenerateTask(task.ID)

	return task, nil
}

// processGenerateTask 处理生成任务
func (s *LearningMaterialService) processGenerateTask(taskID int64) {
	ctx := context.Background()

	// 更新任务状态为生成中
	if err := s.taskRepo.UpdateStatus(ctx, taskID, model.TaskStatusProcessing, 0, ""); err != nil {
		fmt.Printf("更新任务状态失败: %v\n", err)
		return
	}

	// 获取任务详情
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		s.taskRepo.UpdateStatus(ctx, taskID, model.TaskStatusFailed, 0, fmt.Sprintf("获取任务失败: %v", err))
		return
	}

	// 构建生成请求
	req := &GenerateMaterialRequest{
		UserID:     task.UserID,
		SubjectID:  task.SubjectID,
		Grade:      task.Grade,
		Topic:      task.Topic,
		Difficulty: task.Difficulty,
	}

	// 调用原有的同步生成方法
	material, err := s.GenerateMaterial(ctx, req)
	if err != nil {
		s.taskRepo.UpdateStatus(ctx, taskID, model.TaskStatusFailed, 0, err.Error())
		return
	}

	// 更新任务状态为已完成
	if err := s.taskRepo.UpdateStatus(ctx, taskID, model.TaskStatusCompleted, material.ID, ""); err != nil {
		fmt.Printf("更新任务状态失败: %v\n", err)
	}
}

// GetTaskByID 获取任务详情
func (s *LearningMaterialService) GetTaskByID(ctx context.Context, id int64) (*model.MaterialGenerateTask, error) {
	return s.taskRepo.GetByID(ctx, id)
}

// ListTasksByUserID 获取用户的任务列表
func (s *LearningMaterialService) ListTasksByUserID(ctx context.Context, userID int64, statuses []int8) ([]*model.MaterialGenerateTask, error) {
	return s.taskRepo.ListByUserID(ctx, userID, statuses)
}
