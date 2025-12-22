package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"art_admin_backend/internal/api"
	"art_admin_backend/internal/api/middleware"
	v1 "art_admin_backend/internal/api/v1"
	"art_admin_backend/internal/pkg/config"
	"art_admin_backend/internal/pkg/database"
	"art_admin_backend/internal/pkg/logger"
	"art_admin_backend/internal/pkg/volcengine"
	ws "art_admin_backend/internal/pkg/websocket"
	"art_admin_backend/internal/repository"
	"art_admin_backend/internal/router"
	cadSvc "art_admin_backend/internal/service/cad"
	chatSvc "art_admin_backend/internal/service/chat"
	designCompareSvc "art_admin_backend/internal/service/design_compare"
	designerProjectSvc "art_admin_backend/internal/service/designer_project"
	documentSvc "art_admin_backend/internal/service/document"
	fileSvc "art_admin_backend/internal/service/file"

	_ "art_admin_backend/docs" // swagger docs

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Art Admin API
// @version 1.0
// @description Art Admin Backend API Documentation - 基础功能 + 系统功能
// @host localhost:48080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// 自动生成 Swagger 文档
	generateSwaggerDocs()

	// 加载配置
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化日志
	if err := logger.InitLogger(&cfg.Log); err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}
	defer logger.Sync()

	logger.Info("配置加载成功")

	// 完全自动化初始化数据库（创建库→连接→迁移→初始化数据）
	if err := database.InitDB(&cfg.Database); err != nil {
		logger.Fatal(fmt.Sprintf("数据库初始化失败: %v", err))
	}
	logger.Info("数据库连接成功")

	// 执行数据库迁移
	if err := database.RunMigrations(); err != nil {
		logger.Fatal(fmt.Sprintf("数据库迁移失败: %v", err))
	}
	logger.Info("数据库迁移完成")

	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 创建路由
	r := gin.New()

	// 注册全局中间件
	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	// 注册 Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Art Admin Backend is running",
		})
	})

	// 静态文件服务 - 用于访问上传的文件
	r.Static("/uploads", "./uploads")

	// 初始化文件管理服务
	fileRepo := repository.NewFileRepository(database.GetDB())
	fileService := fileSvc.NewFileService(fileRepo, "./uploads", fmt.Sprintf("http://localhost:%d", cfg.Server.Port))
	v1.SetFileService(fileService)
	logger.Info("文件管理服务初始化完成")

	// 初始化设计师项目管理服务 / Initialize designer project service
	designerProjectRepo := repository.NewDesignerProjectRepository(database.GetDB())
	designerProjectService := designerProjectSvc.NewDesignerProjectService(designerProjectRepo)
	v1.SetDesignerProjectService(designerProjectService)
	logger.Info("设计师项目管理服务初始化完成")

	// 初始化火山AI客户端 / Initialize VolcEngine AI client
	volcClient := volcengine.NewClient(
		cfg.VolcEngine.APIKey,
		cfg.VolcEngine.BaseURL,
		cfg.VolcEngine.Model,
		cfg.VolcEngine.Timeout,
		cfg.VolcEngine.MaxTokens,
		cfg.VolcEngine.Temperature,
	)
	logger.Info("火山AI客户端初始化完成")

	// 设置设计版本服务的AI客户端 / Set AI client for design version service
	v1.SetDesignVersionAIClient(volcClient)
	logger.Info("设计版本AI分析服务初始化完成")

	// 初始化文档分析服务 / Initialize document analysis service
	documentRepo := repository.NewDocumentRepository(database.GetDB())
	documentService := documentSvc.NewDocumentService(
		documentRepo,
		designerProjectRepo,
		volcClient,
		"./uploads",
		fmt.Sprintf("http://localhost:%d", cfg.Server.Port),
	)
	v1.SetDocumentService(documentService)
	logger.Info("文档分析服务初始化完成")

	// 初始化设计比对服务 / Initialize design compare service
	designCompareRepo := repository.NewDesignCompareRepository(database.GetDB())
	designCompareService := designCompareSvc.NewDesignCompareService(
		designCompareRepo,
		designerProjectRepo,
		documentRepo,
		volcClient,
		"./uploads",
		fmt.Sprintf("http://localhost:%d", cfg.Server.Port),
	)
	v1.SetDesignCompareService(designCompareService)
	logger.Info("设计比对服务初始化完成")

	// 初始化CAD文件服务 / Initialize CAD file service
	cadFileRepo := repository.NewCadFileRepository(database.GetDB())
	// 优先使用新配置，兼容旧配置 / Prefer new config, fallback to old config
	converterType := cfg.DWGConverter.Type
	odaPath := cfg.DWGConverter.ODAPath
	libredwgPath := cfg.DWGConverter.LibreDWGPath
	if odaPath == "" && cfg.ODAConverter.Path != "" {
		odaPath = cfg.ODAConverter.Path
		converterType = "oda"
	}
	cadFileService := cadSvc.NewCadService(
		cadFileRepo,
		"./uploads",
		fmt.Sprintf("http://localhost:%d", cfg.Server.Port),
		converterType,
		odaPath,
		libredwgPath,
	)
	v1.SetCadService(cadFileService)
	logger.Info("CAD文件服务初始化完成")

	// 初始化CAD渲染服务 / Initialize CAD render service
	// 使用Doubao-1.5-vision-pro生成设计方案，Doubao-Seedream-4.5生成效果图
	// Use Doubao-1.5-vision-pro for design proposal, Doubao-Seedream-4.5 for image generation
	renderService := cadSvc.NewRenderService(
		cadFileRepo,
		designerProjectRepo,
		documentRepo, // 添加文档仓库以支持参考文档 / Add document repo for reference documents
		volcClient,   // 火山引擎AI客户端(文本) / VolcEngine AI client for text
		"./uploads",
		cfg.VolcEngineImage.APIKey,  // 图片生成API Key
		cfg.VolcEngineImage.BaseURL, // 图片生成API URL
		cfg.VolcEngineImage.Model,   // Doubao-Seedream-4.5接入点ID
		cfg.VolcEngineImage.Size,    // 图片尺寸
		cfg.VolcEngineImage.Timeout, // 超时时间
	)
	v1.SetRenderService(renderService)
	logger.Info("CAD渲染服务初始化完成")

	// 初始化通知服务 / Initialize notification service
	notificationRepo := repository.NewNotificationRepository(database.GetDB())
	v1.SetNotificationRepository(notificationRepo)
	logger.Info("通知服务初始化完成")

	// 初始化效果图记录服务 / Initialize render record service
	renderRecordRepo := repository.NewRenderRecordRepository(database.GetDB())
	userQuotaRepo := repository.NewUserRenderQuotaRepository(database.GetDB())
	renderRecordService := cadSvc.NewRenderRecordService(renderRecordRepo, userQuotaRepo, renderService, notificationRepo)
	v1.SetRenderRecordService(renderRecordService)
	logger.Info("效果图记录服务初始化完成")

	// 初始化聊天室系统
	chatHub := ws.NewHub()
	go chatHub.Run() // 启动 WebSocket Hub

	chatRoomRepo := repository.NewChatRoomRepository(database.GetDB())
	chatMessageRepo := repository.NewChatMessageRepository(database.GetDB())
	chatMemberRepo := repository.NewChatRoomMemberRepository(database.GetDB())
	chatService := chatSvc.NewChatService(chatRoomRepo, chatMessageRepo, chatMemberRepo)
	chatAPI := api.NewChatAPI(chatService, chatHub)

	// 设置聊天室 API
	router.SetChatAPI(chatAPI)
	logger.Info("聊天室服务初始化完成")

	// 注册业务路由
	router.RegisterRoutes(r)

	// 启动服务器
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Info(fmt.Sprintf("服务器启动成功，监听端口: %d", cfg.Server.Port))
	logger.Info(fmt.Sprintf("Swagger 文档地址: http://localhost:%d/swagger/index.html", cfg.Server.Port))

	if err := r.Run(addr); err != nil {
		logger.Fatal(fmt.Sprintf("服务器启动失败: %v", err))
	}
}

// generateSwaggerDocs 自动生成 Swagger 文档
func generateSwaggerDocs() {
	fmt.Println("正在生成 Swagger 文档...")

	cmd := exec.Command("swag", "init", "-g", "cmd/server/main.go", "-o", "docs")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		// 如果 swag 命令不存在，提示安装
		fmt.Println("警告: 无法生成 Swagger 文档")
		fmt.Println("请安装 swag: go install github.com/swaggo/swag/cmd/swag@latest")
		fmt.Println("然后手动运行: swag init -g cmd/server/main.go -o docs")
		return
	}

	fmt.Println("Swagger 文档生成成功!")
}
