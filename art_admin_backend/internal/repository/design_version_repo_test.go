package repository

import (
	"art_admin_backend/internal/model"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// setupVersionTestDB 创建版本测试用的内存数据库
// Create in-memory database for version testing
func setupVersionTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// 自动迁移所有相关模型 / Auto migrate all related models
	err = db.AutoMigrate(
		&model.DesignerProject{},
		&model.DesignVersion{},
		&model.DesignVersionCompare{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

// genVersionCount 生成版本数量（1-10）
// Generate version count (1-10)
func genVersionCount() gopter.Gen {
	return gen.IntRange(1, 10)
}

// genUserID 生成用户ID
// Generate user ID
func genUserID() gopter.Gen {
	return gen.UIntRange(1, 100)
}

// createTestProjectForVersion 创建测试项目用于版本测试
// Create test project for version testing
func createTestProjectForVersion(db *gorm.DB, userID uint) (*model.DesignerProject, error) {
	project := &model.DesignerProject{
		Name:        "测试项目",
		Description: "用于版本测试的项目",
		Area:        100.0,
		Budget:      50000,
		Style:       "现代简约",
		Status:      "in_progress",
		UserID:      userID,
	}
	err := db.Create(project).Error
	return project, err
}

// strPtr 将字符串转换为指针
// Convert string to pointer
func strPtr(s string) *string {
	return &s
}

// createTestVersions 创建测试版本
// Create test versions for a project
func createTestVersions(db *gorm.DB, projectID uint, count int, createdBy uint) ([]model.DesignVersion, error) {
	versions := make([]model.DesignVersion, 0, count)
	versionNames := []string{"初稿", "修改版1", "修改版2", "终稿", "客户确认版", "施工版", "备选方案A", "备选方案B", "优化版", "最终版"}

	for i := 0; i < count; i++ {
		nameIndex := i % len(versionNames)
		version := model.DesignVersion{
			ProjectID:     projectID,
			VersionNumber: i + 1,
			VersionName:   versionNames[nameIndex],
			Description:   "版本描述 " + versionNames[nameIndex],
			DesignImages:  strPtr(`[{"url":"/test/image.jpg","name":"设计图","type":"floor_plan"}]`),
			LayoutInfo:    strPtr(`[{"zone":"办公区","area":50,"position":"东侧"}]`),
			AreaInfo:      strPtr(`[{"name":"总面积","value":100,"unit":"sqm"}]`),
			StyleInfo:     strPtr(`[{"category":"主风格","value":"现代简约","details":"简洁明快"}]`),
			MaterialInfo:  strPtr(`[{"name":"地板","category":"地面","quantity":100,"unit":"sqm"}]`),
			Status:        "draft",
			CreatedBy:     createdBy,
		}
		if err := db.Create(&version).Error; err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}
	return versions, nil
}

// **Feature: designer-ai-assistant, Property 11: Version Association Correctness**
// **Validates: Requirements 7.1**
// For any project with multiple design versions, all versions SHALL be correctly
// associated with the same project ID.
func TestVersionAssociationCorrectness(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// 组合生成器：版本数量、用户ID
	// Combined generator: version count, user ID
	testDataGen := gopter.CombineGens(
		genVersionCount(),
		genUserID(),
	)

	properties.Property("all versions are correctly associated with the same project ID", prop.ForAll(
		func(testData []any) bool {
			// 为每个测试用例创建新的数据库实例 / Create new DB instance for each test case
			db := setupVersionTestDB(t)

			versionCount := testData[0].(int)
			userID := testData[1].(uint)

			// 创建项目 / Create project
			project, err := createTestProjectForVersion(db, userID)
			if err != nil {
				t.Logf("Failed to create project: %v", err)
				return false
			}

			// 创建多个版本 / Create multiple versions
			createdVersions, err := createTestVersions(db, project.ID, versionCount, userID)
			if err != nil {
				t.Logf("Failed to create versions: %v", err)
				return false
			}

			// 验证创建的版本数量 / Verify created version count
			if len(createdVersions) != versionCount {
				t.Logf("Version count mismatch: expected %d, got %d", versionCount, len(createdVersions))
				return false
			}

			// 使用仓库查询版本 / Query versions using repository
			repo := NewDesignVersionRepository(db)
			retrievedVersions, err := repo.GetByProjectID(project.ID)
			if err != nil {
				t.Logf("Failed to get versions by project ID: %v", err)
				return false
			}

			// 验证查询到的版本数量 / Verify retrieved version count
			if len(retrievedVersions) != versionCount {
				t.Logf("Retrieved version count mismatch: expected %d, got %d", versionCount, len(retrievedVersions))
				return false
			}

			// 验证所有版本都关联到正确的项目ID / Verify all versions are associated with correct project ID
			for _, version := range retrievedVersions {
				if version.ProjectID != project.ID {
					t.Logf("Version %d has wrong project ID: expected %d, got %d",
						version.ID, project.ID, version.ProjectID)
					return false
				}
			}

			// 验证版本号的唯一性和连续性 / Verify version number uniqueness and continuity
			versionNumbers := make(map[int]bool)
			for _, version := range retrievedVersions {
				if versionNumbers[version.VersionNumber] {
					t.Logf("Duplicate version number found: %d", version.VersionNumber)
					return false
				}
				versionNumbers[version.VersionNumber] = true
			}

			// 验证版本号范围 / Verify version number range
			for i := 1; i <= versionCount; i++ {
				if !versionNumbers[i] {
					t.Logf("Missing version number: %d", i)
					return false
				}
			}

			return true
		},
		testDataGen,
	))

	properties.TestingRun(t)
}

// TestVersionAssociationWithMultipleProjects 测试多项目版本关联隔离性
// Test version association isolation across multiple projects
func TestVersionAssociationWithMultipleProjects(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 50
	properties := gopter.NewProperties(parameters)

	// 组合生成器：项目1版本数、项目2版本数
	// Combined generator: project1 version count, project2 version count
	testDataGen := gopter.CombineGens(
		genVersionCount(),
		genVersionCount(),
	)

	properties.Property("versions from different projects are correctly isolated", prop.ForAll(
		func(testData []any) bool {
			db := setupVersionTestDB(t)

			versionCount1 := testData[0].(int)
			versionCount2 := testData[1].(int)

			// 创建两个项目 / Create two projects
			project1, err := createTestProjectForVersion(db, 1)
			if err != nil {
				t.Logf("Failed to create project1: %v", err)
				return false
			}

			project2, err := createTestProjectForVersion(db, 2)
			if err != nil {
				t.Logf("Failed to create project2: %v", err)
				return false
			}

			// 为两个项目创建版本 / Create versions for both projects
			_, err = createTestVersions(db, project1.ID, versionCount1, 1)
			if err != nil {
				t.Logf("Failed to create versions for project1: %v", err)
				return false
			}

			_, err = createTestVersions(db, project2.ID, versionCount2, 2)
			if err != nil {
				t.Logf("Failed to create versions for project2: %v", err)
				return false
			}

			repo := NewDesignVersionRepository(db)

			// 查询项目1的版本 / Query versions for project1
			versions1, err := repo.GetByProjectID(project1.ID)
			if err != nil {
				t.Logf("Failed to get versions for project1: %v", err)
				return false
			}

			// 查询项目2的版本 / Query versions for project2
			versions2, err := repo.GetByProjectID(project2.ID)
			if err != nil {
				t.Logf("Failed to get versions for project2: %v", err)
				return false
			}

			// 验证版本数量 / Verify version counts
			if len(versions1) != versionCount1 {
				t.Logf("Project1 version count mismatch: expected %d, got %d", versionCount1, len(versions1))
				return false
			}

			if len(versions2) != versionCount2 {
				t.Logf("Project2 version count mismatch: expected %d, got %d", versionCount2, len(versions2))
				return false
			}

			// 验证项目1的版本不包含项目2的版本 / Verify project1 versions don't contain project2 versions
			for _, v := range versions1 {
				if v.ProjectID != project1.ID {
					t.Logf("Project1 versions contain wrong project ID: %d", v.ProjectID)
					return false
				}
			}

			// 验证项目2的版本不包含项目1的版本 / Verify project2 versions don't contain project1 versions
			for _, v := range versions2 {
				if v.ProjectID != project2.ID {
					t.Logf("Project2 versions contain wrong project ID: %d", v.ProjectID)
					return false
				}
			}

			return true
		},
		testDataGen,
	))

	properties.TestingRun(t)
}

// TestVersionCRUDBasic 基础CRUD测试
// Basic CRUD test for design versions
func TestVersionCRUDBasic(t *testing.T) {
	db := setupVersionTestDB(t)
	repo := NewDesignVersionRepository(db)

	// 创建项目 / Create project
	project, err := createTestProjectForVersion(db, 1)
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	// 创建版本 / Create version
	version := &model.DesignVersion{
		ProjectID:     project.ID,
		VersionNumber: 1,
		VersionName:   "初稿",
		Description:   "第一版设计方案",
		Status:        "draft",
		CreatedBy:     1,
	}

	err = repo.Create(version)
	if err != nil {
		t.Fatalf("Failed to create version: %v", err)
	}

	if version.ID == 0 {
		t.Fatal("Version ID should not be 0 after creation")
	}

	// 读取版本 / Read version
	retrieved, err := repo.GetByID(version.ID)
	if err != nil {
		t.Fatalf("Failed to get version: %v", err)
	}

	if retrieved.VersionName != version.VersionName {
		t.Errorf("VersionName mismatch: expected %s, got %s", version.VersionName, retrieved.VersionName)
	}

	if retrieved.ProjectID != project.ID {
		t.Errorf("ProjectID mismatch: expected %d, got %d", project.ID, retrieved.ProjectID)
	}

	// 更新版本 / Update version
	retrieved.VersionName = "修改版1"
	retrieved.Status = "submitted"
	err = repo.Update(retrieved)
	if err != nil {
		t.Fatalf("Failed to update version: %v", err)
	}

	updated, _ := repo.GetByID(version.ID)
	if updated.VersionName != "修改版1" {
		t.Errorf("Update failed: expected '修改版1', got %s", updated.VersionName)
	}
	if updated.Status != "submitted" {
		t.Errorf("Status update failed: expected 'submitted', got %s", updated.Status)
	}

	// 删除版本 / Delete version
	err = repo.Delete(version.ID)
	if err != nil {
		t.Fatalf("Failed to delete version: %v", err)
	}

	// 验证删除（软删除） / Verify deletion (soft delete)
	var count int64
	db.Unscoped().Model(&model.DesignVersion{}).Where("id = ?", version.ID).Count(&count)
	if count != 1 {
		t.Error("Soft delete should keep the record")
	}
}

// TestGetLatestVersionNumber 测试获取最新版本号
// Test getting latest version number
func TestGetLatestVersionNumber(t *testing.T) {
	db := setupVersionTestDB(t)
	repo := NewDesignVersionRepository(db)

	// 创建项目 / Create project
	project, err := createTestProjectForVersion(db, 1)
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	// 初始状态应该返回0 / Initial state should return 0
	latestNum, err := repo.GetLatestVersionNumber(project.ID)
	if err != nil {
		t.Fatalf("Failed to get latest version number: %v", err)
	}
	if latestNum != 0 {
		t.Errorf("Expected 0 for empty project, got %d", latestNum)
	}

	// 创建版本 / Create versions
	_, err = createTestVersions(db, project.ID, 5, 1)
	if err != nil {
		t.Fatalf("Failed to create versions: %v", err)
	}

	// 应该返回5 / Should return 5
	latestNum, err = repo.GetLatestVersionNumber(project.ID)
	if err != nil {
		t.Fatalf("Failed to get latest version number: %v", err)
	}
	if latestNum != 5 {
		t.Errorf("Expected 5, got %d", latestNum)
	}
}

// TestVersionListWithFilter 测试带筛选条件的版本列表
// Test version list with filters
func TestVersionListWithFilter(t *testing.T) {
	db := setupVersionTestDB(t)
	repo := NewDesignVersionRepository(db)

	// 创建项目 / Create project
	project, err := createTestProjectForVersion(db, 1)
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	// 创建多个版本 / Create multiple versions
	versions, err := createTestVersions(db, project.ID, 5, 1)
	if err != nil {
		t.Fatalf("Failed to create versions: %v", err)
	}

	// 更新部分版本状态 / Update some version statuses
	db.Model(&model.DesignVersion{}).Where("id = ?", versions[0].ID).Update("status", "submitted")
	db.Model(&model.DesignVersion{}).Where("id = ?", versions[1].ID).Update("status", "approved")

	// 测试按项目ID筛选 / Test filter by project ID
	filter := DesignVersionFilter{ProjectID: project.ID}
	results, total, err := repo.List(1, 10, filter)
	if err != nil {
		t.Fatalf("Failed to list versions: %v", err)
	}
	if total != 5 {
		t.Errorf("Expected 5 versions, got %d", total)
	}

	// 测试按状态筛选 / Test filter by status
	filter = DesignVersionFilter{ProjectID: project.ID, Status: "draft"}
	results, total, err = repo.List(1, 10, filter)
	if err != nil {
		t.Fatalf("Failed to list versions: %v", err)
	}
	if total != 3 {
		t.Errorf("Expected 3 draft versions, got %d", total)
	}

	_ = results // 避免未使用变量警告 / Avoid unused variable warning
}

// TestCountByProjectID 测试按项目ID统计版本数量
// Test counting versions by project ID
func TestCountByProjectID(t *testing.T) {
	db := setupVersionTestDB(t)
	repo := NewDesignVersionRepository(db)

	// 创建项目 / Create project
	project, err := createTestProjectForVersion(db, 1)
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	// 初始应该为0 / Initial count should be 0
	count, err := repo.CountByProjectID(project.ID)
	if err != nil {
		t.Fatalf("Failed to count versions: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected 0, got %d", count)
	}

	// 创建版本 / Create versions
	_, err = createTestVersions(db, project.ID, 7, 1)
	if err != nil {
		t.Fatalf("Failed to create versions: %v", err)
	}

	// 应该为7 / Should be 7
	count, err = repo.CountByProjectID(project.ID)
	if err != nil {
		t.Fatalf("Failed to count versions: %v", err)
	}
	if count != 7 {
		t.Errorf("Expected 7, got %d", count)
	}
}

// TestDeleteByProjectID 测试按项目ID删除所有版本
// Test deleting all versions by project ID
func TestDeleteByProjectID(t *testing.T) {
	db := setupVersionTestDB(t)
	repo := NewDesignVersionRepository(db)

	// 创建项目 / Create project
	project, err := createTestProjectForVersion(db, 1)
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	// 创建版本 / Create versions
	_, err = createTestVersions(db, project.ID, 5, 1)
	if err != nil {
		t.Fatalf("Failed to create versions: %v", err)
	}

	// 验证版本已创建 / Verify versions were created
	count, _ := repo.CountByProjectID(project.ID)
	if count != 5 {
		t.Fatalf("Expected 5 versions, got %d", count)
	}

	// 删除所有版本 / Delete all versions
	err = repo.DeleteByProjectID(project.ID)
	if err != nil {
		t.Fatalf("Failed to delete versions by project ID: %v", err)
	}

	// 验证版本已删除 / Verify versions were deleted
	count, _ = repo.CountByProjectID(project.ID)
	if count != 0 {
		t.Errorf("Expected 0 versions after deletion, got %d", count)
	}
}
