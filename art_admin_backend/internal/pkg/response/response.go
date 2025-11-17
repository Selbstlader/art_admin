package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code int         `json:"code"` // 状态码
	Msg  string      `json:"msg"`  // 消息
	Data interface{} `json:"data"` // 数据
}

// 状态码常量
const (
	CodeSuccess      = 200 // 成功
	CodeBadRequest   = 400 // 参数错误
	CodeUnauthorized = 401 // 未授权
	CodeForbidden    = 403 // 权限不足
	CodeNotFound     = 404 // 资源不存在
	CodeServerError  = 500 // 服务器错误
)

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: CodeSuccess,
		Msg:  "success",
		Data: data,
	})
}

// SuccessWithMsg 成功响应（自定义消息）
func SuccessWithMsg(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: CodeSuccess,
		Msg:  msg,
		Data: data,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}

// BadRequest 参数错误响应
func BadRequest(c *gin.Context, msg string) {
	Error(c, CodeBadRequest, msg)
}

// Unauthorized 未授权响应
func Unauthorized(c *gin.Context, msg string) {
	Error(c, CodeUnauthorized, msg)
}

// Forbidden 权限不足响应
func Forbidden(c *gin.Context, msg string) {
	Error(c, CodeForbidden, msg)
}

// NotFound 资源不存在响应
func NotFound(c *gin.Context, msg string) {
	Error(c, CodeNotFound, msg)
}

// ServerError 服务器错误响应
func ServerError(c *gin.Context, msg string) {
	Error(c, CodeServerError, msg)
}

// PaginatedData 分页数据结构
type PaginatedData struct {
	Records interface{} `json:"records"` // 记录列表
	Current int         `json:"current"` // 当前页码
	Size    int         `json:"size"`    // 每页条数
	Total   int64       `json:"total"`   // 总记录数
}

// SuccessWithPagination 分页成功响应
func SuccessWithPagination(c *gin.Context, records interface{}, current, size int, total int64) {
	Success(c, PaginatedData{
		Records: records,
		Current: current,
		Size:    size,
		Total:   total,
	})
}

