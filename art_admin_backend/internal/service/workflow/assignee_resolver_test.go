package workflow

import (
	"context"
	"fmt"
	"testing"

	"art_admin_backend/internal/model"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// ============================================================================
// Mock Providers for Testing
// ============================================================================

// mockUserProvider 模拟用户数据提供者
type mockUserProvider struct {
	users     map[int64]*model.User
	roleUsers map[int64][]int64 // roleID -> userIDs
}

func newMockUserProvider() *mockUserProvider {
	return &mockUserProvider{
		users:     make(map[int64]*model.User),
		roleUsers: make(map[int64][]int64),
	}
}

func (p *mockUserProvider) AddUser(user *model.User) {
	p.users[user.ID] = user
}

func (p *mockUserProvider) AddRoleUsers(roleID int64, userIDs []int64) {
	p.roleUsers[roleID] = userIDs
}

func (p *mockUserProvider) GetUserByID(ctx context.Context, userID int64) (*model.User, error) {
	if user, exists := p.users[userID]; exists {
		return user, nil
	}
	return nil, ErrUserNotFound
}

func (p *mockUserProvider) GetUsersByIDs(ctx context.Context, userIDs []int64) ([]*model.User, error) {
	result := make([]*model.User, 0)
	for _, id := range userIDs {
		if user, exists := p.users[id]; exists {
			result = append(result, user)
		}
	}
	return result, nil
}

func (p *mockUserProvider) GetUsersByRoleID(ctx context.Context, roleID int64) ([]*model.User, error) {
	userIDs, exists := p.roleUsers[roleID]
	if !exists {
		return []*model.User{}, nil
	}
	result := make([]*model.User, 0)
	for _, id := range userIDs {
		if user, exists := p.users[id]; exists {
			result = append(result, user)
		}
	}
	return result, nil
}

// mockDepartmentProvider 模拟部门数据提供者
type mockDepartmentProvider struct {
	departments map[int64]*model.Department
	leaders     map[int64]int64 // deptID -> leaderUserID
}

func newMockDepartmentProvider() *mockDepartmentProvider {
	return &mockDepartmentProvider{
		departments: make(map[int64]*model.Department),
		leaders:     make(map[int64]int64),
	}
}

func (p *mockDepartmentProvider) AddDepartment(dept *model.Department, leaderID int64) {
	p.departments[dept.DeptID] = dept
	if leaderID > 0 {
		p.leaders[dept.DeptID] = leaderID
	}
}

func (p *mockDepartmentProvider) GetDepartmentByID(ctx context.Context, deptID int64) (*model.Department, error) {
	if dept, exists := p.departments[deptID]; exists {
		return dept, nil
	}
	return nil, ErrDepartmentNotFound
}

func (p *mockDepartmentProvider) GetDepartmentLeader(ctx context.Context, deptID int64) (int64, error) {
	if leaderID, exists := p.leaders[deptID]; exists {
		return leaderID, nil
	}
	return 0, nil
}

// mockOrganizationProvider 模拟组织架构数据提供者
type mockOrganizationProvider struct {
	superiors map[int64]int64 // userID -> superiorUserID
}

func newMockOrganizationProvider() *mockOrganizationProvider {
	return &mockOrganizationProvider{
		superiors: make(map[int64]int64),
	}
}

func (p *mockOrganizationProvider) SetSuperior(userID, superiorID int64) {
	p.superiors[userID] = superiorID
}

func (p *mockOrganizationProvider) GetUserSuperior(ctx context.Context, userID int64) (int64, error) {
	if superiorID, exists := p.superiors[userID]; exists {
		return superiorID, nil
	}
	return 0, nil
}

// ============================================================================
// Test Helpers
// ============================================================================

// testAssigneeResolver 可测试的审批人解析器
type testAssigneeResolver struct {
	userProvider *mockUserProvider
	deptProvider *mockDepartmentProvider
	orgProvider  *mockOrganizationProvider
	resolver     AssigneeResolver
}

func newTestAssigneeResolver() *testAssigneeResolver {
	userProvider := newMockUserProvider()
	deptProvider := newMockDepartmentProvider()
	orgProvider := newMockOrganizationProvider()

	return &testAssigneeResolver{
		userProvider: userProvider,
		deptProvider: deptProvider,
		orgProvider:  orgProvider,
		resolver:     NewAssigneeResolver(userProvider, deptProvider, orgProvider),
	}
}

// ============================================================================
// Generators
// ============================================================================

// genUserIDs generates a slice of user IDs
func genUserIDs() gopter.Gen {
	return gen.SliceOfN(5, gen.Int64Range(1, 100)).SuchThat(func(ids []int64) bool {
		return len(ids) > 0
	})
}

// genRoleIDs generates a slice of role IDs
func genRoleIDs() gopter.Gen {
	return gen.SliceOfN(3, gen.Int64Range(1, 50)).SuchThat(func(ids []int64) bool {
		return len(ids) > 0
	})
}

// genDeptID generates a department ID
func genDeptID() gopter.Gen {
	return gen.Int64Range(1, 20)
}

// genInitiatorID generates an initiator user ID
func genInitiatorID() gopter.Gen {
	return gen.Int64Range(1, 100)
}

// ============================================================================
// Property-Based Tests
// ============================================================================

// **Feature: oa-workflow-engine, Property 11: 审批人规则解析正确性**
// **Validates: Requirements 5.1, 5.2, 5.3, 5.4**
// *For any* AssigneeRule：
// - 指定用户规则应返回配置的用户ID列表
// - 指定角色规则应返回拥有该角色的所有用户
// - 部门负责人规则应返回发起人所属部门的负责人
// - 发起人上级规则应返回发起人的直接上级

// TestAssigneeResolverUserRule tests that user rule returns configured user IDs
// Requirements: 5.1 - WHEN 审批人规则配置为指定用户 THEN Workflow_Engine SHALL 直接分配Task给配置的用户列表
func TestAssigneeResolverUserRule(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	parameters.Rng.Seed(42)

	properties := gopter.NewProperties(parameters)

	properties.Property("User rule returns exactly the configured user IDs that exist", prop.ForAll(
		func(userIDs []int64) bool {
			ctx := context.Background()
			tr := newTestAssigneeResolver()

			// Setup: Add users to the provider
			for _, id := range userIDs {
				tr.userProvider.AddUser(&model.User{
					ID:       id,
					UserName: fmt.Sprintf("user_%d", id),
				})
			}

			// Create rule with user IDs
			rule := &model.AssigneeRule{
				Type:   model.AssigneeTypeUser,
				Values: userIDs,
			}

			// Resolve
			result, err := tr.resolver.Resolve(ctx, rule, &ProcessContext{})
			if err != nil {
				t.Logf("Resolve failed: %v", err)
				return false
			}

			// Verify: result should contain all user IDs
			if len(result) != len(userIDs) {
				t.Logf("Expected %d users, got %d", len(userIDs), len(result))
				return false
			}

			// Create a set of expected IDs
			expectedSet := make(map[int64]bool)
			for _, id := range userIDs {
				expectedSet[id] = true
			}

			// Verify all returned IDs are in expected set
			for _, id := range result {
				if !expectedSet[id] {
					t.Logf("Unexpected user ID %d in result", id)
					return false
				}
			}

			return true
		},
		genUserIDs(),
	))

	properties.TestingRun(t)
}

// TestAssigneeResolverRoleRule tests that role rule returns users with the specified roles
// Requirements: 5.2 - WHEN 审批人规则配置为指定角色 THEN Workflow_Engine SHALL 查询拥有该角色的所有用户并分配Task
func TestAssigneeResolverRoleRule(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	parameters.Rng.Seed(42)

	properties := gopter.NewProperties(parameters)

	properties.Property("Role rule returns all users with the specified roles", prop.ForAll(
		func(roleID int64, userIDs []int64) bool {
			ctx := context.Background()
			tr := newTestAssigneeResolver()

			// Deduplicate user IDs first (since the map will deduplicate anyway)
			uniqueUserIDs := make(map[int64]bool)
			for _, id := range userIDs {
				uniqueUserIDs[id] = true
			}

			// Setup: Add users and associate them with the role
			for _, id := range userIDs {
				tr.userProvider.AddUser(&model.User{
					ID:       id,
					UserName: fmt.Sprintf("user_%d", id),
				})
			}
			tr.userProvider.AddRoleUsers(roleID, userIDs)

			// Create rule with role ID
			rule := &model.AssigneeRule{
				Type:   model.AssigneeTypeRole,
				Values: []int64{roleID},
			}

			// Resolve
			result, err := tr.resolver.Resolve(ctx, rule, &ProcessContext{})
			if err != nil {
				t.Logf("Resolve failed: %v", err)
				return false
			}

			// Verify: result should contain all unique users with the role
			if len(result) != len(uniqueUserIDs) {
				t.Logf("Expected %d unique users, got %d", len(uniqueUserIDs), len(result))
				return false
			}

			// Verify all returned IDs are in expected set
			for _, id := range result {
				if !uniqueUserIDs[id] {
					t.Logf("Unexpected user ID %d in result", id)
					return false
				}
			}

			return true
		},
		gen.Int64Range(1, 50),
		genUserIDs(),
	))

	properties.TestingRun(t)
}

// TestAssigneeResolverMultipleRoles tests that role rule with multiple roles returns union of users
func TestAssigneeResolverMultipleRoles(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	parameters.Rng.Seed(42)

	properties := gopter.NewProperties(parameters)

	properties.Property("Multiple role rule returns union of users from all roles", prop.ForAll(
		func(role1ID, role2ID int64, users1, users2 []int64) bool {
			ctx := context.Background()
			tr := newTestAssigneeResolver()

			// Setup: Add users for both roles
			allUsers := make(map[int64]bool)
			for _, id := range users1 {
				tr.userProvider.AddUser(&model.User{
					ID:       id,
					UserName: fmt.Sprintf("user_%d", id),
				})
				allUsers[id] = true
			}
			for _, id := range users2 {
				tr.userProvider.AddUser(&model.User{
					ID:       id,
					UserName: fmt.Sprintf("user_%d", id),
				})
				allUsers[id] = true
			}
			tr.userProvider.AddRoleUsers(role1ID, users1)
			tr.userProvider.AddRoleUsers(role2ID, users2)

			// Create rule with multiple role IDs
			rule := &model.AssigneeRule{
				Type:   model.AssigneeTypeRole,
				Values: []int64{role1ID, role2ID},
			}

			// Resolve
			result, err := tr.resolver.Resolve(ctx, rule, &ProcessContext{})
			if err != nil {
				t.Logf("Resolve failed: %v", err)
				return false
			}

			// Verify: result should contain union of all users (deduplicated)
			if len(result) != len(allUsers) {
				t.Logf("Expected %d unique users, got %d", len(allUsers), len(result))
				return false
			}

			// Verify all returned IDs are in expected set
			for _, id := range result {
				if !allUsers[id] {
					t.Logf("Unexpected user ID %d in result", id)
					return false
				}
			}

			return true
		},
		gen.Int64Range(1, 50),
		gen.Int64Range(51, 100), // Different role ID to avoid collision
		genUserIDs(),
		genUserIDs(),
	))

	properties.TestingRun(t)
}

// TestAssigneeResolverDeptLeaderRule tests that dept_leader rule returns the department leader
// Requirements: 5.3 - WHEN 审批人规则配置为部门负责人 THEN Workflow_Engine SHALL 根据发起人所属部门查找负责人并分配Task
func TestAssigneeResolverDeptLeaderRule(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	parameters.Rng.Seed(42)

	properties := gopter.NewProperties(parameters)

	properties.Property("Dept leader rule returns the department leader of initiator's department", prop.ForAll(
		func(deptID, leaderID, initiatorID int64) bool {
			ctx := context.Background()
			tr := newTestAssigneeResolver()

			// Setup: Add department with leader
			tr.deptProvider.AddDepartment(&model.Department{
				DeptID:   deptID,
				DeptName: fmt.Sprintf("dept_%d", deptID),
			}, leaderID)

			// Add leader user
			tr.userProvider.AddUser(&model.User{
				ID:       leaderID,
				UserName: fmt.Sprintf("leader_%d", leaderID),
			})

			// Create rule
			rule := &model.AssigneeRule{
				Type: model.AssigneeTypeDeptLeader,
			}

			// Create process context with initiator's department
			processCtx := &ProcessContext{
				InitiatorID: initiatorID,
				DeptID:      &deptID,
			}

			// Resolve
			result, err := tr.resolver.Resolve(ctx, rule, processCtx)
			if err != nil {
				t.Logf("Resolve failed: %v", err)
				return false
			}

			// Verify: result should contain exactly the leader ID
			if len(result) != 1 {
				t.Logf("Expected 1 user, got %d", len(result))
				return false
			}

			if result[0] != leaderID {
				t.Logf("Expected leader ID %d, got %d", leaderID, result[0])
				return false
			}

			return true
		},
		genDeptID(),
		gen.Int64Range(1, 100), // leader ID
		genInitiatorID(),
	))

	properties.TestingRun(t)
}

// TestAssigneeResolverInitiatorLeaderRule tests that initiator_leader rule returns the initiator's superior
// Requirements: 5.4 - WHEN 审批人规则配置为发起人上级 THEN Workflow_Engine SHALL 根据组织架构查找发起人的直接上级并分配Task
func TestAssigneeResolverInitiatorLeaderRule(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	parameters.Rng.Seed(42)

	properties := gopter.NewProperties(parameters)

	properties.Property("Initiator leader rule returns the initiator's direct superior", prop.ForAll(
		func(initiatorID, superiorID int64) bool {
			ctx := context.Background()
			tr := newTestAssigneeResolver()

			// Setup: Set superior for initiator
			tr.orgProvider.SetSuperior(initiatorID, superiorID)

			// Add superior user
			tr.userProvider.AddUser(&model.User{
				ID:       superiorID,
				UserName: fmt.Sprintf("superior_%d", superiorID),
			})

			// Create rule
			rule := &model.AssigneeRule{
				Type: model.AssigneeTypeInitiatorLeader,
			}

			// Create process context
			processCtx := &ProcessContext{
				InitiatorID: initiatorID,
			}

			// Resolve
			result, err := tr.resolver.Resolve(ctx, rule, processCtx)
			if err != nil {
				t.Logf("Resolve failed: %v", err)
				return false
			}

			// Verify: result should contain exactly the superior ID
			if len(result) != 1 {
				t.Logf("Expected 1 user, got %d", len(result))
				return false
			}

			if result[0] != superiorID {
				t.Logf("Expected superior ID %d, got %d", superiorID, result[0])
				return false
			}

			return true
		},
		genInitiatorID(),
		gen.Int64Range(1, 100), // superior ID
	))

	properties.TestingRun(t)
}

// ============================================================================
// Property 12: Empty Result Handling Tests
// ============================================================================

// **Feature: oa-workflow-engine, Property 12: 审批人规则空结果处理**
// **Validates: Requirements 5.5**
// *For any* AssigneeRule解析结果为空时，系统应返回错误而非创建无Assignee的Task。

// TestAssigneeResolverEmptyUserRule tests that empty user rule returns error
func TestAssigneeResolverEmptyUserRule(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	parameters.Rng.Seed(42)

	properties := gopter.NewProperties(parameters)

	properties.Property("Empty user rule returns ErrNoAssigneeFound", prop.ForAll(
		func(initiatorID int64) bool {
			ctx := context.Background()
			tr := newTestAssigneeResolver()

			// Create rule with empty user list
			rule := &model.AssigneeRule{
				Type:   model.AssigneeTypeUser,
				Values: []int64{},
			}

			// Resolve
			result, err := tr.resolver.Resolve(ctx, rule, &ProcessContext{
				InitiatorID: initiatorID,
			})

			// Should return error
			if err != ErrNoAssigneeFound {
				t.Logf("Expected ErrNoAssigneeFound, got %v", err)
				return false
			}

			// Result should be nil
			if result != nil {
				t.Logf("Expected nil result, got %v", result)
				return false
			}

			return true
		},
		genInitiatorID(),
	))

	properties.TestingRun(t)
}

// TestAssigneeResolverNonExistentUsers tests that rule with non-existent users returns error
func TestAssigneeResolverNonExistentUsers(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	parameters.Rng.Seed(42)

	properties := gopter.NewProperties(parameters)

	properties.Property("User rule with non-existent users returns ErrNoAssigneeFound", prop.ForAll(
		func(userIDs []int64) bool {
			ctx := context.Background()
			tr := newTestAssigneeResolver()

			// Don't add any users to the provider - they don't exist

			// Create rule with user IDs
			rule := &model.AssigneeRule{
				Type:   model.AssigneeTypeUser,
				Values: userIDs,
			}

			// Resolve
			result, err := tr.resolver.Resolve(ctx, rule, &ProcessContext{})

			// Should return error because no users exist
			if err != ErrNoAssigneeFound {
				t.Logf("Expected ErrNoAssigneeFound, got %v", err)
				return false
			}

			// Result should be nil
			if result != nil {
				t.Logf("Expected nil result, got %v", result)
				return false
			}

			return true
		},
		genUserIDs(),
	))

	properties.TestingRun(t)
}

// TestAssigneeResolverEmptyRoleRule tests that role with no users returns error
func TestAssigneeResolverEmptyRoleRule(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	parameters.Rng.Seed(42)

	properties := gopter.NewProperties(parameters)

	properties.Property("Role rule with no users returns ErrNoAssigneeFound", prop.ForAll(
		func(roleID int64) bool {
			ctx := context.Background()
			tr := newTestAssigneeResolver()

			// Don't add any users to the role

			// Create rule with role ID
			rule := &model.AssigneeRule{
				Type:   model.AssigneeTypeRole,
				Values: []int64{roleID},
			}

			// Resolve
			result, err := tr.resolver.Resolve(ctx, rule, &ProcessContext{})

			// Should return error because role has no users
			if err != ErrNoAssigneeFound {
				t.Logf("Expected ErrNoAssigneeFound, got %v", err)
				return false
			}

			// Result should be nil
			if result != nil {
				t.Logf("Expected nil result, got %v", result)
				return false
			}

			return true
		},
		gen.Int64Range(1, 50),
	))

	properties.TestingRun(t)
}

// TestAssigneeResolverNoDeptLeader tests that dept without leader returns error
func TestAssigneeResolverNoDeptLeader(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	parameters.Rng.Seed(42)

	properties := gopter.NewProperties(parameters)

	properties.Property("Dept leader rule with no leader returns ErrNoAssigneeFound", prop.ForAll(
		func(deptID, initiatorID int64) bool {
			ctx := context.Background()
			tr := newTestAssigneeResolver()

			// Add department without leader (leaderID = 0)
			tr.deptProvider.AddDepartment(&model.Department{
				DeptID:   deptID,
				DeptName: fmt.Sprintf("dept_%d", deptID),
			}, 0) // No leader

			// Create rule
			rule := &model.AssigneeRule{
				Type: model.AssigneeTypeDeptLeader,
			}

			// Create process context
			processCtx := &ProcessContext{
				InitiatorID: initiatorID,
				DeptID:      &deptID,
			}

			// Resolve
			result, err := tr.resolver.Resolve(ctx, rule, processCtx)

			// Should return error because no leader
			if err != ErrNoAssigneeFound {
				t.Logf("Expected ErrNoAssigneeFound, got %v", err)
				return false
			}

			// Result should be nil
			if result != nil {
				t.Logf("Expected nil result, got %v", result)
				return false
			}

			return true
		},
		genDeptID(),
		genInitiatorID(),
	))

	properties.TestingRun(t)
}

// TestAssigneeResolverNoSuperior tests that initiator without superior returns error
func TestAssigneeResolverNoSuperior(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	parameters.Rng.Seed(42)

	properties := gopter.NewProperties(parameters)

	properties.Property("Initiator leader rule with no superior returns ErrNoAssigneeFound", prop.ForAll(
		func(initiatorID int64) bool {
			ctx := context.Background()
			tr := newTestAssigneeResolver()

			// Don't set any superior for the initiator

			// Create rule
			rule := &model.AssigneeRule{
				Type: model.AssigneeTypeInitiatorLeader,
			}

			// Create process context
			processCtx := &ProcessContext{
				InitiatorID: initiatorID,
			}

			// Resolve
			result, err := tr.resolver.Resolve(ctx, rule, processCtx)

			// Should return error because no superior
			if err != ErrNoAssigneeFound {
				t.Logf("Expected ErrNoAssigneeFound, got %v", err)
				return false
			}

			// Result should be nil
			if result != nil {
				t.Logf("Expected nil result, got %v", result)
				return false
			}

			return true
		},
		genInitiatorID(),
	))

	properties.TestingRun(t)
}

// TestAssigneeResolverNilRule tests that nil rule returns error
func TestAssigneeResolverNilRule(t *testing.T) {
	ctx := context.Background()
	tr := newTestAssigneeResolver()

	// Resolve with nil rule
	result, err := tr.resolver.Resolve(ctx, nil, &ProcessContext{})

	// Should return error
	if err != ErrInvalidAssigneeRule {
		t.Errorf("Expected ErrInvalidAssigneeRule, got %v", err)
	}

	// Result should be nil
	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}
}
