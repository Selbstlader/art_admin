package v1

import (
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/pkg/response"
	"art_admin_backend/internal/service"
	"fmt"

	"github.com/gin-gonic/gin"
)

var departmentService = service.NewDepartmentService()

// GetDepartmentList 获取部门列表
// @Summary 获取部门列表
// @Description 查询部门列表（树形结构）
// @Tags 部门管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param deptName query string false "部门名称"
// @Param deptCode query string false "部门编码"
// @Param status query int false "部门状态" Enums(0,1)
// @Success 200 {object} response.Response{data=[]response.DepartmentListItem} "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/department/list [get]
func GetDepartmentList(c *gin.Context) {
	var req request.DepartmentListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	departments, err := departmentService.GetDepartmentList(&req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	// 添加调试信息
	fmt.Printf("查询到的部门数量: %d\n", len(departments))
	for i, dept := range departments {
		fmt.Printf("部门 %d: ID=%d, ParentID=%v, Children数量=%d\n", i, dept.DeptID, dept.ParentID, len(dept.Children))
	}

	response.Success(c, departments)
}

// GetDepartmentByID 根据ID获取部门
// @Summary 根据ID获取部门
// @Description 根据部门ID获取部门详情
// @Tags 部门管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "部门ID"
// @Success 200 {object} response.Response{data=response.DepartmentItem} "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "部门不存在"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/department/{id} [get]
func GetDepartmentByID(c *gin.Context) {
	var req request.DepartmentByIDRequest
	if err := c.ShouldBindUri(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	department, err := departmentService.GetDepartmentByID(req.DeptID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, department)
}

// CreateDepartment 创建部门
// @Summary 创建部门
// @Description 创建新的部门
// @Tags 部门管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body request.CreateDepartmentRequest true "部门信息"
// @Success 200 {object} response.Response "创建成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/department [post]
func CreateDepartment(c *gin.Context) {
	var req request.CreateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	err := departmentService.CreateDepartment(&req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// UpdateDepartment 更新部门
// @Summary 更新部门
// @Description 更新部门信息
// @Tags 部门管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body request.UpdateDepartmentRequest true "部门信息"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "部门不存在"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/department [put]
func UpdateDepartment(c *gin.Context) {
	var req request.UpdateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	err := departmentService.UpdateDepartment(&req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// DeleteDepartment 删除部门
// @Summary 删除部门
// @Description 删除指定部门
// @Tags 部门管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "部门ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "部门不存在"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/department/{id} [delete]
func DeleteDepartment(c *gin.Context) {
	var req request.DepartmentByIDRequest
	if err := c.ShouldBindUri(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	err := departmentService.DeleteDepartment(req.DeptID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}
