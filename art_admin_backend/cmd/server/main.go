package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"art_admin_backend/internal/api/middleware"
	"art_admin_backend/internal/pkg/config"
	"art_admin_backend/internal/pkg/database"
	"art_admin_backend/internal/pkg/logger"
	"art_admin_backend/internal/router"

	_ "art_admin_backend/docs" // swagger docs

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Art Admin API
// @version 1.0
// @description Art Admin Backend API Documentation
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
	// 无需手动操作数据库，GORM会自动管理一切
	if err := database.InitializeDatabase(&cfg.Database); err != nil {
		logger.Fatal(fmt.Sprintf("初始化数据库失败: %v", err))
	}
	defer database.CloseDB()

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
