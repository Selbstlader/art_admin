package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"time"

	"art_admin_backend/internal/dto/request"
	operationLogSvc "art_admin_backend/internal/service/operation_log"

	"github.com/gin-gonic/gin"
)

var operationLogService = operationLogSvc.NewOperationLogService()

// 业务类型映射
var businessTypeMap = map[string]string{
	"POST":   "新增",
	"PUT":    "修改",
	"DELETE": "删除",
	"GET":    "查询",
}

// 模块名称映射
var moduleNameMap = map[string]string{
	"/api/user":          "用户管理",
	"/api/role":          "角色管理",
	"/api/department":    "部门管理",
	"/api/menu":          "菜单管理",
	"/api/dictionary":    "字典管理",
	"/api/operation-log": "操作日志",
}

// responseWriter 自定义响应写入器，用于捕获响应数据
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// OperationLog 操作日志中间件
func OperationLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过不需要记录的路径
		if shouldSkipLog(c.Request.URL.Path) {
			c.Next()
			return
		}

		startTime := time.Now()

		// 读取请求参数
		var requestParam string
		contentType := c.Request.Header.Get("Content-Type")

		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			// 检查是否为文件上传请求（multipart/form-data）
			if strings.Contains(contentType, "multipart/form-data") {
				// 文件上传请求，不记录二进制内容，只记录基本信息
				requestParam = "[文件上传请求] " + c.Request.URL.RawQuery
			} else {
				bodyBytes, _ := io.ReadAll(c.Request.Body)
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				requestParam = string(bodyBytes)
				// 检查是否包含不可打印字符（二进制数据）
				if containsBinaryData(requestParam) {
					requestParam = "[二进制数据请求]"
				}
				// 限制参数长度
				if len(requestParam) > 2000 {
					requestParam = requestParam[:2000] + "..."
				}
			}
		} else if c.Request.Method == "GET" {
			requestParam = c.Request.URL.RawQuery
		}

		// 创建自定义响应写入器
		blw := &responseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = blw

		// 处理请求
		c.Next()

		// 计算耗时
		costTime := time.Since(startTime).Milliseconds()

		// 获取响应数据
		responseData := blw.body.String()
		if len(responseData) > 2000 {
			responseData = responseData[:2000] + "..."
		}

		// 判断操作状态
		status := 1 // 默认成功
		errorMsg := ""
		if c.Writer.Status() >= 400 {
			status = 0
			// 尝试从响应中提取错误信息
			var respData map[string]interface{}
			if err := json.Unmarshal([]byte(responseData), &respData); err == nil {
				if msg, ok := respData["msg"].(string); ok {
					errorMsg = msg
				}
			}
		}

		// 获取操作人员
		operatorName := "未知"
		if userName, exists := c.Get("userName"); exists {
			operatorName = userName.(string)
		} else if userID := GetUserID(c); userID > 0 {
			// 如果有userID但没有userName，可以从数据库查询
			// 这里简化处理，直接使用userID
			operatorName = string(rune(userID))
		}

		// 获取模块名称
		module := getModuleName(c.Request.URL.Path)

		// 获取业务类型
		businessType := getBusinessType(c.Request.Method)

		// 创建操作日志
		logReq := &request.CreateOperationLogRequest{
			Module:        module,
			BusinessType:  businessType,
			RequestMethod: c.Request.Method,
			RequestURL:    c.Request.URL.Path,
			OperatorName:  operatorName,
			OperatorIP:    c.ClientIP(),
			OperatorAddr:  "", // 可以通过IP获取地理位置
			RequestParam:  requestParam,
			ResponseData:  responseData,
			Status:        status,
			ErrorMsg:      errorMsg,
			CostTime:      costTime,
			UserAgent:     c.Request.UserAgent(),
		}

		// 异步记录日志，不影响主流程
		go func() {
			if err := operationLogService.CreateOperationLog(logReq); err != nil {
				// 记录日志失败，可以输出到系统日志
				// log.Printf("记录操作日志失败: %v", err)
			}
		}()
	}
}

// shouldSkipLog 判断是否跳过日志记录
func shouldSkipLog(path string) bool {
	// 跳过的路径列表
	skipPaths := []string{
		"/api/auth/login",         // 登录接口
		"/api/user/info",          // 获取用户信息
		"/api/operation-log/list", // 操作日志列表查询
		"/api/operation-log/",     // 操作日志详情查询
	}

	for _, skipPath := range skipPaths {
		if strings.HasPrefix(path, skipPath) {
			return true
		}
	}

	// 跳过GET请求（只记录增删改操作）
	// 如果需要记录所有操作，可以注释掉这个判断
	// return c.Request.Method == "GET"

	return false
}

// getModuleName 获取模块名称
func getModuleName(path string) string {
	for prefix, name := range moduleNameMap {
		if strings.HasPrefix(path, prefix) {
			return name
		}
	}
	return "其他"
}

// getBusinessType 获取业务类型
func getBusinessType(method string) string {
	if businessType, ok := businessTypeMap[method]; ok {
		return businessType
	}
	return "其他"
}

// containsBinaryData 检查字符串是否包含二进制数据（不可打印字符）
// Check if string contains binary data (non-printable characters)
func containsBinaryData(s string) bool {
	// 检查前100个字符是否包含不可打印字符
	checkLen := len(s)
	if checkLen > 100 {
		checkLen = 100
	}
	for i := 0; i < checkLen; i++ {
		b := s[i]
		// 允许的字符：可打印ASCII字符、换行、回车、制表符、UTF-8多字节字符
		if b < 32 && b != '\n' && b != '\r' && b != '\t' {
			return true
		}
	}
	return false
}
