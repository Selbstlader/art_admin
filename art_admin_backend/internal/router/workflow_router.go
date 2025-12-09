package router

import (
	v1 "art_admin_backend/internal/api/v1"

	"github.com/gin-gonic/gin"
)

// SetupWorkflowRoutes 设置工作流相关路由
func SetupWorkflowRoutes(authGroup *gin.RouterGroup) {
	// 工作流模块
	workflowGroup := authGroup.Group("/workflow")
	{
		// 流程定义管理
		processDefGroup := workflowGroup.Group("/process-def")
		{
			processDefGroup.GET("/list", v1.GetProcessDefList)
			processDefGroup.GET("/:id", v1.GetProcessDefDetail)
			processDefGroup.POST("", v1.CreateProcessDef)
			processDefGroup.PUT("/:id", v1.UpdateProcessDef)
			processDefGroup.DELETE("/:id", v1.DeleteProcessDef)
			processDefGroup.POST("/:id/publish", v1.PublishProcessDef)
		}

		// 流程实例管理
		processInstGroup := workflowGroup.Group("/process-inst")
		{
			processInstGroup.POST("/start", v1.StartProcess)
			processInstGroup.POST("/complete-task", v1.CompleteTask)
			processInstGroup.POST("/delegate-task", v1.DelegateTask)
			processInstGroup.POST("/transfer-task", v1.TransferTask)
			processInstGroup.POST("/:id/withdraw", v1.WithdrawProcess)
			processInstGroup.GET("/:id", v1.GetProcessInstDetail)
		}

		// 表单模板管理
		formTemplateGroup := workflowGroup.Group("/form-template")
		{
			formTemplateGroup.GET("/list", v1.GetFormTemplateList)
			formTemplateGroup.GET("/:id", v1.GetFormTemplateDetail)
			formTemplateGroup.POST("", v1.CreateFormTemplate)
			formTemplateGroup.PUT("/:id", v1.UpdateFormTemplate)
			formTemplateGroup.DELETE("/:id", v1.DeleteFormTemplate)
		}

		// 查询接口
		queryGroup := workflowGroup.Group("/query")
		{
			queryGroup.GET("/my-initiated", v1.GetMyInitiated)
			queryGroup.GET("/my-todo", v1.GetMyTodo)
			queryGroup.GET("/my-done", v1.GetMyDone)
			queryGroup.GET("/approval-trail/:id", v1.GetApprovalTrail)
			queryGroup.GET("/statistics", v1.GetWorkflowStatistics)
		}
	}
}
