package response

import "time"

// DepartmentListItem 部门列表项
type DepartmentListItem struct {
	DeptID     int64                `json:"deptId"`
	ParentID   *int64               `json:"parentId"`
	ParentName string               `json:"parentName,omitempty"`
	DeptName   string               `json:"deptName"`
	DeptCode   string               `json:"deptCode"`
	OrderNum   int                  `json:"orderNum"`
	Leader     string               `json:"leader,omitempty"`
	Phone      string               `json:"phone,omitempty"`
	Email      string               `json:"email,omitempty"`
	Status     int                  `json:"status"`
	CreateTime time.Time            `json:"createTime"`
	UpdateTime time.Time            `json:"updateTime,omitempty"`
	Children   []DepartmentListItem `json:"children,omitempty"`
}

// DepartmentItem 部门详情项
type DepartmentItem = DepartmentListItem
