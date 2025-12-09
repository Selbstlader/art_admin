package workflow

import (
	"context"
	"errors"
	"fmt"

	"art_admin_backend/internal/model"
)

// ============================================================================
// Errors
// ============================================================================

var (
	// ErrNoAssigneeFound 审批人解析结果为空
	ErrNoAssigneeFound = errors.New("审批人规则解析结果为空")
	// ErrInvalidAssigneeRule 无效的审批人规则
	ErrInvalidAssigneeRule = errors.New("无效的审批人规则")
	// ErrUserNotFound 用户不存在
	ErrUserNotFound = errors.New("用户不存在")
	// ErrDepartmentNotFound 部门不存在
	ErrDepartmentNotFound = errors.New("部门不存在")
)

// ============================================================================
// Interfaces
// ============================================================================

// ProcessContext 流程上下文，包含解析审批人所需的信息
type ProcessContext struct {
	InitiatorID   int64                  // 发起人ID
	InitiatorName string                 // 发起人名称
	DeptID        *int64                 // 发起人部门ID
	FormData      map[string]interface{} // 表单数据
}

// UserProvider 用户数据提供者接口
type UserProvider interface {
	// GetUserByID 根据ID获取用户
	GetUserByID(ctx context.Context, userID int64) (*model.User, error)
	// GetUsersByIDs 根据ID列表获取用户
	GetUsersByIDs(ctx context.Context, userIDs []int64) ([]*model.User, error)
	// GetUsersByRoleID 根据角色ID获取用户列表
	GetUsersByRoleID(ctx context.Context, roleID int64) ([]*model.User, error)
}

// DepartmentProvider 部门数据提供者接口
type DepartmentProvider interface {
	// GetDepartmentByID 根据ID获取部门
	GetDepartmentByID(ctx context.Context, deptID int64) (*model.Department, error)
	// GetDepartmentLeader 获取部门负责人用户ID
	GetDepartmentLeader(ctx context.Context, deptID int64) (int64, error)
}

// OrganizationProvider 组织架构数据提供者接口
type OrganizationProvider interface {
	// GetUserSuperior 获取用户的直接上级
	GetUserSuperior(ctx context.Context, userID int64) (int64, error)
}

// AssigneeResolver 审批人解析器接口
// Requirements: 5.1, 5.2, 5.3, 5.4 - 支持四种审批人规则类型解析
type AssigneeResolver interface {
	// Resolve 解析审批人规则，返回用户ID列表
	// Requirements: 5.1 - 指定用户规则应返回配置的用户ID列表
	// Requirements: 5.2 - 指定角色规则应返回拥有该角色的所有用户
	// Requirements: 5.3 - 部门负责人规则应返回发起人所属部门的负责人
	// Requirements: 5.4 - 发起人上级规则应返回发起人的直接上级
	Resolve(ctx context.Context, rule *model.AssigneeRule, processCtx *ProcessContext) ([]int64, error)
}

// ============================================================================
// Implementation
// ============================================================================

// assigneeResolver 审批人解析器实现
type assigneeResolver struct {
	userProvider UserProvider
	deptProvider DepartmentProvider
	orgProvider  OrganizationProvider
}

// NewAssigneeResolver 创建审批人解析器实例
func NewAssigneeResolver(
	userProvider UserProvider,
	deptProvider DepartmentProvider,
	orgProvider OrganizationProvider,
) AssigneeResolver {
	return &assigneeResolver{
		userProvider: userProvider,
		deptProvider: deptProvider,
		orgProvider:  orgProvider,
	}
}

// Resolve 解析审批人规则，返回用户ID列表
// **Feature: oa-workflow-engine, Property 11: 审批人规则解析正确性**
// **Validates: Requirements 5.1, 5.2, 5.3, 5.4**
func (r *assigneeResolver) Resolve(ctx context.Context, rule *model.AssigneeRule, processCtx *ProcessContext) ([]int64, error) {
	if rule == nil {
		return nil, ErrInvalidAssigneeRule
	}

	var userIDs []int64
	var err error

	switch rule.Type {
	case model.AssigneeTypeUser:
		// Requirements: 5.1 - 指定用户规则应返回配置的用户ID列表
		userIDs, err = r.resolveUserRule(ctx, rule)
	case model.AssigneeTypeRole:
		// Requirements: 5.2 - 指定角色规则应返回拥有该角色的所有用户
		userIDs, err = r.resolveRoleRule(ctx, rule)
	case model.AssigneeTypeDeptLeader:
		// Requirements: 5.3 - 部门负责人规则应返回发起人所属部门的负责人
		userIDs, err = r.resolveDeptLeaderRule(ctx, processCtx)
	case model.AssigneeTypeInitiatorLeader:
		// Requirements: 5.4 - 发起人上级规则应返回发起人的直接上级
		userIDs, err = r.resolveInitiatorLeaderRule(ctx, processCtx)
	default:
		return nil, fmt.Errorf("%w: 未知的规则类型 %s", ErrInvalidAssigneeRule, rule.Type)
	}

	if err != nil {
		return nil, err
	}

	// Requirements: 5.5 - 审批人规则解析结果为空时返回错误
	// **Feature: oa-workflow-engine, Property 12: 审批人规则空结果处理**
	// **Validates: Requirements 5.5**
	if len(userIDs) == 0 {
		return nil, ErrNoAssigneeFound
	}

	return userIDs, nil
}

// resolveUserRule 解析指定用户规则
// Requirements: 5.1 - WHEN 审批人规则配置为指定用户 THEN Workflow_Engine SHALL 直接分配Task给配置的用户列表
func (r *assigneeResolver) resolveUserRule(ctx context.Context, rule *model.AssigneeRule) ([]int64, error) {
	if len(rule.Values) == 0 {
		return nil, nil
	}

	// 验证用户是否存在
	users, err := r.userProvider.GetUsersByIDs(ctx, rule.Values)
	if err != nil {
		return nil, fmt.Errorf("获取用户列表失败: %w", err)
	}

	// 返回存在的用户ID
	userIDs := make([]int64, 0, len(users))
	for _, user := range users {
		userIDs = append(userIDs, user.ID)
	}

	return userIDs, nil
}

// resolveRoleRule 解析指定角色规则
// Requirements: 5.2 - WHEN 审批人规则配置为指定角色 THEN Workflow_Engine SHALL 查询拥有该角色的所有用户并分配Task
func (r *assigneeResolver) resolveRoleRule(ctx context.Context, rule *model.AssigneeRule) ([]int64, error) {
	if len(rule.Values) == 0 {
		return nil, nil
	}

	// 收集所有角色下的用户
	userIDSet := make(map[int64]struct{})
	for _, roleID := range rule.Values {
		users, err := r.userProvider.GetUsersByRoleID(ctx, roleID)
		if err != nil {
			return nil, fmt.Errorf("获取角色用户失败: %w", err)
		}
		for _, user := range users {
			userIDSet[user.ID] = struct{}{}
		}
	}

	// 转换为切片
	userIDs := make([]int64, 0, len(userIDSet))
	for userID := range userIDSet {
		userIDs = append(userIDs, userID)
	}

	return userIDs, nil
}

// resolveDeptLeaderRule 解析部门负责人规则
// Requirements: 5.3 - WHEN 审批人规则配置为部门负责人 THEN Workflow_Engine SHALL 根据发起人所属部门查找负责人并分配Task
func (r *assigneeResolver) resolveDeptLeaderRule(ctx context.Context, processCtx *ProcessContext) ([]int64, error) {
	if processCtx == nil || processCtx.DeptID == nil {
		return nil, nil
	}

	leaderID, err := r.deptProvider.GetDepartmentLeader(ctx, *processCtx.DeptID)
	if err != nil {
		return nil, fmt.Errorf("获取部门负责人失败: %w", err)
	}

	if leaderID == 0 {
		return nil, nil
	}

	return []int64{leaderID}, nil
}

// resolveInitiatorLeaderRule 解析发起人上级规则
// Requirements: 5.4 - WHEN 审批人规则配置为发起人上级 THEN Workflow_Engine SHALL 根据组织架构查找发起人的直接上级并分配Task
func (r *assigneeResolver) resolveInitiatorLeaderRule(ctx context.Context, processCtx *ProcessContext) ([]int64, error) {
	if processCtx == nil || processCtx.InitiatorID == 0 {
		return nil, nil
	}

	superiorID, err := r.orgProvider.GetUserSuperior(ctx, processCtx.InitiatorID)
	if err != nil {
		return nil, fmt.Errorf("获取发起人上级失败: %w", err)
	}

	if superiorID == 0 {
		return nil, nil
	}

	return []int64{superiorID}, nil
}
