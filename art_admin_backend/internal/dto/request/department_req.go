package request

// DepartmentListRequest 部门列表查询请求
type DepartmentListRequest struct {
	DeptName string `form:"deptName" example:"技术部"`
	DeptCode string `form:"deptCode" example:"TECH"`
	Status   *int   `form:"status" example:"1"`
}

// CreateDepartmentRequest 创建部门请求
type CreateDepartmentRequest struct {
	ParentID *int64 `json:"parentId" example:"1"`
	DeptName string `json:"deptName" binding:"required,max=50"`
	DeptCode string `json:"deptCode" binding:"required,max=50"`
	OrderNum int    `json:"orderNum" binding:"min=0,max=999"`
	Leader   string `json:"leader" binding:"max=20"`
	Phone    string `json:"phone" binding:"max=11"`
	Email    string `json:"email" binding:"max=50,email"`
	Status   int    `json:"status" binding:"oneof=0 1"`
}

// UpdateDepartmentRequest 更新部门请求
type UpdateDepartmentRequest struct {
	DeptID   int64  `json:"deptId" binding:"required"`
	ParentID *int64 `json:"parentId" example:"1"`
	DeptName string `json:"deptName" binding:"required,max=50"`
	DeptCode string `json:"deptCode" binding:"required,max=50"`
	OrderNum int    `json:"orderNum" binding:"min=0,max=999"`
	Leader   string `json:"leader" binding:"max=20"`
	Phone    string `json:"phone" binding:"max=11"`
	Email    string `json:"email" binding:"max=50,email"`
	Status   int    `json:"status" binding:"oneof=0 1"`
}

// DepartmentByIDRequest 根据ID获取部门请求
type DepartmentByIDRequest struct {
	DeptID int64 `uri:"id" binding:"required"`
}
