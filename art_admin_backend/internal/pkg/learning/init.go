package learning

import (
	"art_admin_backend/internal/api"
	"art_admin_backend/internal/pkg/deepseek"
	"art_admin_backend/internal/pkg/dify"
	"art_admin_backend/internal/repository"
	"art_admin_backend/internal/service"

	"gorm.io/gorm"
)

// Config 学习系统配置
type Config struct {
	DeepSeek DeepSeekConfig
	Dify     DifyConfig
}

// DeepSeekConfig DeepSeek 配置
type DeepSeekConfig struct {
	APIKey      string
	BaseURL     string
	Model       string
	Timeout     int
	MaxTokens   int
	Temperature float64
}

// DifyConfig Dify 配置
type DifyConfig struct {
	DatasetAPIKey string
	BaseURL       string
	Timeout       int
	DatasetID     string
}

// Container 学习系统容器
type Container struct {
	LearningAPI *api.LearningAPI
}

// NewContainer 创建学习系统容器
func NewContainer(db *gorm.DB, config Config) *Container {
	// 创建客户端
	deepseekClient := deepseek.NewClient(deepseek.Config{
		APIKey:      config.DeepSeek.APIKey,
		BaseURL:     config.DeepSeek.BaseURL,
		Model:       config.DeepSeek.Model,
		Timeout:     config.DeepSeek.Timeout,
		MaxTokens:   config.DeepSeek.MaxTokens,
		Temperature: config.DeepSeek.Temperature,
	})

	difyClient := dify.NewKnowledgeClient(dify.KnowledgeConfig{
		APIKey:  config.Dify.DatasetAPIKey,
		BaseURL: config.Dify.BaseURL,
		Timeout: config.Dify.Timeout,
	})

	// 创建仓储
	subjectRepo := repository.NewSubjectRepository(db)
	materialRepo := repository.NewLearningMaterialRepository(db)
	taskRepo := repository.NewMaterialGenerateTaskRepository(db)

	// 创建服务
	subjectService := service.NewSubjectService(subjectRepo)
	materialService := service.NewLearningMaterialService(
		materialRepo,
		subjectRepo,
		taskRepo,
		deepseekClient,
		difyClient,
		config.Dify.DatasetID,
	)

	// 创建 API
	learningAPI := api.NewLearningAPI(materialService, subjectService)

	return &Container{
		LearningAPI: learningAPI,
	}
}
