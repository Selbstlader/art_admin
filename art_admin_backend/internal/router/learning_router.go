package router

import (
	"art_admin_backend/internal/api"
	"art_admin_backend/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

// SetupLearningRoutes 设置学习系统路由
func SetupLearningRoutes(r *gin.Engine, learningAPI *api.LearningAPI) {
	learning := r.Group("/api/learning")
	{
		// 学科相关（无需认证）
		learning.GET("/subject/list", learningAPI.ListSubjects)

		// 需要认证的接口
		auth := learning.Group("")
		auth.Use(middleware.JWTAuth())
		{
			// 年级列表
			auth.GET("/grades", learningAPI.GetGradesBySubject)

			// 教材相关
			material := auth.Group("/material")
			{
				material.POST("/generate", learningAPI.GenerateMaterial)
				material.GET("/list", learningAPI.ListMaterials)
				material.GET("/:id", learningAPI.GetMaterial)
				material.DELETE("/:id", learningAPI.DeleteMaterial)

				// 任务相关
				material.GET("/tasks", learningAPI.ListTasks)
				material.GET("/task/:id", learningAPI.GetTask)
			}
		}
	}
}
