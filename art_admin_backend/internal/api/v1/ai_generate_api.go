package v1

import (
	"context"
	"time"

	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/pkg/config"
	"art_admin_backend/internal/pkg/deepseek"
	pkgResponse "art_admin_backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

/*** AI Generate API Handler ***/
/*** Provides AI-powered content generation using DeepSeek ***/

var deepseekClient *deepseek.Client

// initDeepSeekClient initializes the DeepSeek client lazily
func initDeepSeekClient() *deepseek.Client {
	if deepseekClient == nil {
		deepseekClient = deepseek.NewClient(deepseek.Config{
			APIKey:      config.GlobalConfig.DeepSeek.APIKey,
			BaseURL:     config.GlobalConfig.DeepSeek.BaseURL,
			Model:       config.GlobalConfig.DeepSeek.Model,
			Timeout:     config.GlobalConfig.DeepSeek.Timeout,
			MaxTokens:   config.GlobalConfig.DeepSeek.MaxTokens,
			Temperature: config.GlobalConfig.DeepSeek.Temperature,
		})
	}
	return deepseekClient
}

// GeneratePrompt generates or polishes a system prompt using AI
// @Summary Generate or polish system prompt
// @Description Use AI to generate a new system prompt or polish an existing one based on keywords
// @Tags AI生成
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body request.GeneratePromptRequest true "Generation request"
// @Success 200 {object} pkgResponse.Response{data=map[string]string} "Success with generated prompt"
// @Failure 400 {object} pkgResponse.Response "Bad request"
// @Failure 401 {object} pkgResponse.Response "Unauthorized"
// @Failure 500 {object} pkgResponse.Response "Server error"
// @Router /api/ai-generate/prompt [post]
func GeneratePrompt(c *gin.Context) {
	var req request.GeneratePromptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkgResponse.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// Default mode is generate
	if req.Mode == "" {
		if req.CurrentPrompt != "" {
			req.Mode = "polish"
		} else {
			req.Mode = "generate"
		}
	}

	client := initDeepSeekClient()
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var systemPrompt, userMessage string

	if req.Mode == "polish" {
		// Polish existing prompt
		systemPrompt = `你是一位专业的AI提示词工程师，擅长优化和润色系统提示词。
请根据用户提供的关键词和现有提示词，进行优化和润色。

优化原则：
1. 保持原有提示词的核心意图和功能
2. 根据关键词补充和强化相关能力描述
3. 使语言更加专业、清晰、有条理
4. 添加必要的约束条件和输出格式要求
5. 确保提示词结构完整，包含角色定义、能力描述、行为准则

请直接输出优化后的提示词，不要添加任何解释或说明。`

		userMessage = "关键词：" + req.Keywords + "\n\n现有提示词：\n" + req.CurrentPrompt
	} else {
		// Generate new prompt
		systemPrompt = `你是一位专业的AI提示词工程师，擅长根据关键词创建高质量的系统提示词。

创建原则：
1. 根据关键词明确AI助手的角色定位
2. 详细描述AI助手应具备的能力和专业知识
3. 设定清晰的行为准则和交互风格
4. 添加必要的约束条件（如安全性、准确性要求）
5. 如适用，指定输出格式和结构

提示词结构建议：
- 角色定义：明确AI的身份和专业领域
- 能力描述：列出核心能力和知识范围
- 行为准则：规定交互方式和回复风格
- 约束条件：设定边界和限制
- 输出要求：指定回复格式（如适用）

请直接输出生成的提示词，不要添加任何解释或说明。`

		userMessage = "请根据以下关键词生成一个专业的AI助手系统提示词：\n\n" + req.Keywords
	}

	result, err := client.SimpleChat(ctx, systemPrompt, userMessage)
	if err != nil {
		pkgResponse.ServerError(c, "AI生成失败: "+err.Error())
		return
	}

	pkgResponse.Success(c, map[string]string{
		"prompt": result,
	})
}
