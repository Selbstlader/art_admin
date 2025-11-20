package service

import (
	"errors"

	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/dto/response"
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/repository"
)

type DepartmentService struct {
	deptRepo *repository.DepartmentRepository
}

func NewDepartmentService() *DepartmentService {
	return &DepartmentService{
		deptRepo: repository.NewDepartmentRepository(),
	}
}

// GetDepartmentList 获取部门列表（树形结构）
func (s *DepartmentService) GetDepartmentList(req *request.DepartmentListRequest) ([]response.DepartmentListItem, error) {
	// 构建查询条件
	query := make(map[string]interface{})
	if req.DeptName != "" {
		query["deptName"] = req.DeptName
	}
	if req.DeptCode != "" {
		query["deptCode"] = req.DeptCode
	}
	if req.Status != nil {
		query["status"] = *req.Status
	}

	// 获取所有部门
	departments, err := s.deptRepo.FindAll(query)
	if err != nil {
		return nil, err
	}

	// 转换为树形结构
	return s.buildDepartmentTree(departments), nil
}

// GetDepartmentByID 根据ID获取部门
func (s *DepartmentService) GetDepartmentByID(id int64) (*response.DepartmentItem, error) {
	department, err := s.deptRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if department == nil {
		return nil, errors.New("部门不存在")
	}

	// 获取父部门信息
	var parentName string
	if department.ParentID != nil {
		parent, err := s.deptRepo.FindByID(*department.ParentID)
		if err == nil && parent != nil {
			parentName = parent.DeptName
		}
	}

	// 转换为响应格式
	resp := &response.DepartmentItem{
		DeptID:     department.DeptID,
		ParentID:   department.ParentID,
		ParentName: parentName,
		DeptName:   department.DeptName,
		DeptCode:   department.DeptCode,
		OrderNum:   department.OrderNum,
		Leader:     department.Leader,
		Phone:      department.Phone,
		Email:      department.Email,
		Status:     department.Status,
		CreateTime: department.CreateTime,
		UpdateTime: department.UpdateTime,
	}

	return resp, nil
}

// CreateDepartment 创建部门
func (s *DepartmentService) CreateDepartment(req *request.CreateDepartmentRequest) error {
	// 检查部门编码是否已存在
	existingDept, err := s.deptRepo.FindByDeptCode(req.DeptCode)
	if err != nil {
		return err
	}
	if existingDept != nil {
		return errors.New("部门编码已存在")
	}

	// 检查父部门是否存在
	if req.ParentID != nil {
		parent, err := s.deptRepo.FindByID(*req.ParentID)
		if err != nil {
			return err
		}
		if parent == nil {
			return errors.New("父部门不存在")
		}

		// 检查是否会形成循环引用
		if s.wouldCreateCycle(*req.ParentID, 0) {
			return errors.New("不能选择当前部门或其子部门作为父部门")
		}
	}

	// 创建部门
	department := &model.Department{
		ParentID: req.ParentID,
		DeptName: req.DeptName,
		DeptCode: req.DeptCode,
		OrderNum: req.OrderNum,
		Leader:   req.Leader,
		Phone:    req.Phone,
		Email:    req.Email,
		Status:   req.Status,
	}

	return s.deptRepo.Create(department)
}

// UpdateDepartment 更新部门
func (s *DepartmentService) UpdateDepartment(req *request.UpdateDepartmentRequest) error {
	// 检查部门是否存在
	department, err := s.deptRepo.FindByID(req.DeptID)
	if err != nil {
		return err
	}
	if department == nil {
		return errors.New("部门不存在")
	}

	// 检查部门编码是否与其他部门冲突
	existingDept, err := s.deptRepo.FindByDeptCode(req.DeptCode)
	if err != nil {
		return err
	}
	if existingDept != nil && existingDept.DeptID != req.DeptID {
		return errors.New("部门编码已存在")
	}

	// 检查父部门
	if req.ParentID != nil {
		// 不能选择自己作为父部门
		if *req.ParentID == req.DeptID {
			return errors.New("不能选择自己作为父部门")
		}

		// 检查父部门是否存在
		parent, err := s.deptRepo.FindByID(*req.ParentID)
		if err != nil {
			return err
		}
		if parent == nil {
			return errors.New("父部门不存在")
		}

		// 检查是否会形成循环引用
		if s.wouldCreateCycle(*req.ParentID, req.DeptID) {
			return errors.New("不能选择当前部门或其子部门作为父部门")
		}
	}

	// 更新部门
	department.ParentID = req.ParentID
	department.DeptName = req.DeptName
	department.DeptCode = req.DeptCode
	department.OrderNum = req.OrderNum
	department.Leader = req.Leader
	department.Phone = req.Phone
	department.Email = req.Email
	department.Status = req.Status

	return s.deptRepo.Update(department)
}

// DeleteDepartment 删除部门
func (s *DepartmentService) DeleteDepartment(id int64) error {
	// 检查部门是否存在
	department, err := s.deptRepo.FindByID(id)
	if err != nil {
		return err
	}
	if department == nil {
		return errors.New("部门不存在")
	}

	// 检查是否有子部门
	children, err := s.deptRepo.FindByParentID(id)
	if err != nil {
		return err
	}
	if len(children) > 0 {
		return errors.New("存在子部门，无法删除")
	}

	// 检查是否有用户关联到该部门
	users, err := s.deptRepo.FindUsersByDeptID(id)
	if err != nil {
		return err
	}
	if len(users) > 0 {
		return errors.New("部门下存在用户，无法删除")
	}

	return s.deptRepo.Delete(id)
}

// buildDepartmentTree 构建部门树形结构
func (s *DepartmentService) buildDepartmentTree(departments []model.Department) []response.DepartmentListItem {
	// 创建部门映射
	deptMap := make(map[int64]*model.Department)
	for i := range departments {
		deptMap[departments[i].DeptID] = &departments[i]
	}

	// 构建树形结构
	var roots []response.DepartmentListItem
	for i := range departments {
		dept := &departments[i]
		if dept.ParentID == nil {
			// 根部门,递归构建子树
			roots = append(roots, s.buildDepartmentNode(dept, deptMap))
		}
	}

	return roots
}

// buildDepartmentNode 递归构建部门节点
func (s *DepartmentService) buildDepartmentNode(dept *model.Department, deptMap map[int64]*model.Department) response.DepartmentListItem {
	node := response.DepartmentListItem{
		DeptID:     dept.DeptID,
		ParentID:   dept.ParentID,
		DeptName:   dept.DeptName,
		DeptCode:   dept.DeptCode,
		OrderNum:   dept.OrderNum,
		Leader:     dept.Leader,
		Phone:      dept.Phone,
		Email:      dept.Email,
		Status:     dept.Status,
		CreateTime: dept.CreateTime,
		UpdateTime: dept.UpdateTime,
		Children:   make([]response.DepartmentListItem, 0),
	}

	// 查找并添加子部门
	for _, childDept := range deptMap {
		if childDept.ParentID != nil && *childDept.ParentID == dept.DeptID {
			node.Children = append(node.Children, s.buildDepartmentNode(childDept, deptMap))
		}
	}

	return node
}

// wouldCreateCycle 检查是否会形成循环引用
func (s *DepartmentService) wouldCreateCycle(parentID, currentDeptID int64) bool {
	// 如果是新增部门（currentDeptID为0），只需要检查父部门及其子部门
	if currentDeptID == 0 {
		return s.isDescendant(parentID, 0)
	}

	// 如果是更新部门，需要检查父部门是否是当前部门或其子部门
	return s.isDescendant(parentID, currentDeptID)
}

// isDescendant 检查targetID是否是deptId的子孙部门
func (s *DepartmentService) isDescendant(targetID, deptId int64) bool {
	if targetID == deptId {
		return true
	}

	children, err := s.deptRepo.FindByParentID(targetID)
	if err != nil {
		return false
	}

	for _, child := range children {
		if s.isDescendant(child.DeptID, deptId) {
			return true
		}
	}

	return false
}
