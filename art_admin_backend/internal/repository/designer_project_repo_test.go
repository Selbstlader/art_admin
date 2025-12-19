package repository

import (
	"art_admin_backend/internal/model"
	"testing"
	"time"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// setupTestDB 创建测试用的内存数据库
// Create in-memory database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// 自动迁移所有相关模型 / Auto migrate all related models
	err = db.AutoMigrate(
		&model.DesignerProject{},
		&model.ProjectDocument{},
		&model.CadFile{},
		&model.DesignCompareResult{},
		&model.CostEstimate{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

// genProjectName 生成有效的项目名称
// Generate valid project name
func genProjectName() gopter.Gen {
	return gen.AlphaString().SuchThat(func(s string) bool {
		return len(s) > 0 && len(s) <= 200
	})
}

// genProjectStatus 生成有效的项目状态
// Generate valid project status
func genProjectStatus() gopter.Gen {
	return gen.OneConstOf("draft", "in_progress", "completed", "archived")
}

// genProjectStyle 生成有效的设计风格
// Generate valid design style
func genProjectStyle() gopter.Gen {
	return gen.OneConstOf("现代简约", "工业风", "北欧风", "中式", "日式", "美式")
}

// genPositiveFloat 生成正浮点数
// Generate positive float
func genPositiveFloat() gopter.Gen {
	return gen.Float64Range(0, 1000000)
}

// genDocumentCount 生成文档数量（0-5）
// Generate document count (0-5)
func genDocumentCount() gopter.Gen {
	return gen.IntRange(0, 5)
}

// genCadFileCount 生成CAD文件数量（0-3）
// Generate CAD file count (0-3)
func genCadFileCount() gopter.Gen {
	return gen.IntRange(0, 3)
}

// createTestProject 创建测试项目
// Create test project
func createTestProject(db *gorm.DB, name, status, style string, area, budget float64, userID uint) (*model.DesignerProject, error) {
	project := &model.DesignerProject{
		Name:        name,
		Description: "Test project description",
		Area:        area,
		Budget:      budget,
		Style:       style,
		Status:      status,
		UserID:      userID,
	}
	err := db.Create(project).Error
	return project, err
}

// createTestDocuments 创建测试文档
// Create test documents
func createTestDocuments(db *gorm.DB, projectID uint, count int) error {
	for i := 0; i < count; i++ {
		doc := &model.ProjectDocument{
			ProjectID:      projectID,
			FileName:       "test_doc_" + string(rune('A'+i)) + ".pdf",
			FilePath:       "/uploads/test/doc_" + string(rune('A'+i)) + ".pdf",
			FileType:       "pdf",
			FileSize:       1024 * int64(i+1),
			AnalysisStatus: "pending",
		}
		if err := db.Create(doc).Error; err != nil {
			return err
		}
	}
	return nil
}

// createTestCadFiles 创建测试CAD文件
// Create test CAD files
func createTestCadFiles(db *gorm.DB, projectID uint, count int) error {
	for i := 0; i < count; i++ {
		cad := &model.CadFile{
			ProjectID:    projectID,
			FileName:     "test_cad_" + string(rune('A'+i)) + ".dwg",
			OriginalPath: "/uploads/cad/test_" + string(rune('A'+i)) + ".dwg",
			FileFormat:   "dwg",
			ParseStatus:  "pending",
		}
		if err := db.Create(cad).Error; err != nil {
			return err
		}
	}
	return nil
}

// createTestCompareResults 创建测试比对结果
// Create test compare results
func createTestCompareResults(db *gorm.DB, projectID uint, count int) error {
	for i := 0; i < count; i++ {
		result := &model.DesignCompareResult{
			ProjectID:       projectID,
			DesignImagePath: "/uploads/design/test_" + string(rune('A'+i)) + ".jpg",
			OverallScore:    float64(80 + i),
		}
		if err := db.Create(result).Error; err != nil {
			return err
		}
	}
	return nil
}

// createTestCostEstimate 创建测试成本估算
// Create test cost estimate
func createTestCostEstimate(db *gorm.DB, projectID uint) error {
	estimate := &model.CostEstimate{
		ProjectID:      projectID,
		MaterialCost:   10000,
		LaborCost:      5000,
		EquipmentCost:  2000,
		ManagementCost: 1000,
		TotalCost:      18000,
		BudgetLimit:    20000,
	}
	return db.Create(estimate).Error
}

// **Feature: designer-ai-assistant, Property 6: Project Cascade Delete**
// **Validates: Requirements 4.4**
// For any project deletion, all associated documents, analysis results, and CAD files
// SHALL be deleted from the database.
func TestProjectCascadeDelete(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// 组合生成器：项目名称、状态、风格、面积、预算、文档数量、CAD文件数量
	// Combined generator: project name, status, style, area, budget, doc count, cad count
	testDataGen := gopter.CombineGens(
		genProjectName(),
		genProjectStatus(),
		genProjectStyle(),
		genPositiveFloat(),
		genPositiveFloat(),
		genDocumentCount(),
		genCadFileCount(),
	)

	properties.Property("deleting project removes all associated data", prop.ForAll(
		func(testData []any) bool {
			// 为每个测试用例创建新的数据库实例 / Create new DB instance for each test case
			db := setupTestDBForProperty(t)

			name := testData[0].(string)
			status := testData[1].(string)
			style := testData[2].(string)
			area := testData[3].(float64)
			budget := testData[4].(float64)
			docCount := testData[5].(int)
			cadCount := testData[6].(int)

			// 创建项目 / Create project
			project, err := createTestProject(db, name, status, style, area, budget, 1)
			if err != nil {
				t.Logf("Failed to create project: %v", err)
				return false
			}

			// 创建关联数据 / Create associated data
			if docCount > 0 {
				if err := createTestDocuments(db, project.ID, docCount); err != nil {
					t.Logf("Failed to create documents: %v", err)
					return false
				}
			}

			if cadCount > 0 {
				if err := createTestCadFiles(db, project.ID, cadCount); err != nil {
					t.Logf("Failed to create CAD files: %v", err)
					return false
				}
			}

			// 创建比对结果和成本估算 / Create compare results and cost estimate
			if err := createTestCompareResults(db, project.ID, 2); err != nil {
				t.Logf("Failed to create compare results: %v", err)
				return false
			}

			if err := createTestCostEstimate(db, project.ID); err != nil {
				t.Logf("Failed to create cost estimate: %v", err)
				return false
			}

			// 验证关联数据已创建 / Verify associated data was created
			var docCountBefore int64
			db.Model(&model.ProjectDocument{}).Where("project_id = ?", project.ID).Count(&docCountBefore)
			if int(docCountBefore) != docCount {
				t.Logf("Document count mismatch before delete: expected %d, got %d", docCount, docCountBefore)
				return false
			}

			var cadCountBefore int64
			db.Model(&model.CadFile{}).Where("project_id = ?", project.ID).Count(&cadCountBefore)
			if int(cadCountBefore) != cadCount {
				t.Logf("CAD file count mismatch before delete: expected %d, got %d", cadCount, cadCountBefore)
				return false
			}

			// 执行级联删除 / Execute cascade delete
			repo := NewDesignerProjectRepository(db)
			if err := repo.DeleteWithAssociations(project.ID); err != nil {
				t.Logf("Failed to delete project with associations: %v", err)
				return false
			}

			// 验证项目已删除 / Verify project was deleted
			var projectCount int64
			db.Model(&model.DesignerProject{}).Where("id = ?", project.ID).Count(&projectCount)
			if projectCount != 0 {
				t.Logf("Project was not deleted")
				return false
			}

			// 验证所有关联数据已删除 / Verify all associated data was deleted
			var docCountAfter int64
			db.Model(&model.ProjectDocument{}).Where("project_id = ?", project.ID).Count(&docCountAfter)
			if docCountAfter != 0 {
				t.Logf("Documents were not deleted: %d remaining", docCountAfter)
				return false
			}

			var cadCountAfter int64
			db.Model(&model.CadFile{}).Where("project_id = ?", project.ID).Count(&cadCountAfter)
			if cadCountAfter != 0 {
				t.Logf("CAD files were not deleted: %d remaining", cadCountAfter)
				return false
			}

			var compareCountAfter int64
			db.Model(&model.DesignCompareResult{}).Where("project_id = ?", project.ID).Count(&compareCountAfter)
			if compareCountAfter != 0 {
				t.Logf("Compare results were not deleted: %d remaining", compareCountAfter)
				return false
			}

			var costCountAfter int64
			db.Model(&model.CostEstimate{}).Where("project_id = ?", project.ID).Count(&costCountAfter)
			if costCountAfter != 0 {
				t.Logf("Cost estimate was not deleted: %d remaining", costCountAfter)
				return false
			}

			return true
		},
		testDataGen,
	))

	properties.TestingRun(t)
}

// setupTestDBForProperty 为属性测试创建独立的数据库实例
// Create independent database instance for property testing
func setupTestDBForProperty(t *testing.T) *gorm.DB {
	// 使用唯一的数据库名称避免并发问题 / Use unique DB name to avoid concurrency issues
	dbName := ":memory:"
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// 自动迁移 / Auto migrate
	err = db.AutoMigrate(
		&model.DesignerProject{},
		&model.ProjectDocument{},
		&model.CadFile{},
		&model.DesignCompareResult{},
		&model.CostEstimate{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

// TestProjectCRUDBasic 基础CRUD测试
// Basic CRUD test
func TestProjectCRUDBasic(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDesignerProjectRepository(db)

	// 创建项目 / Create project
	project := &model.DesignerProject{
		Name:        "测试项目",
		Description: "这是一个测试项目",
		Area:        100.5,
		Budget:      50000,
		Style:       "现代简约",
		Status:      "draft",
		UserID:      1,
	}

	err := repo.Create(project)
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	if project.ID == 0 {
		t.Fatal("Project ID should not be 0 after creation")
	}

	// 读取项目 / Read project
	retrieved, err := repo.GetByID(project.ID)
	if err != nil {
		t.Fatalf("Failed to get project: %v", err)
	}

	if retrieved.Name != project.Name {
		t.Errorf("Name mismatch: expected %s, got %s", project.Name, retrieved.Name)
	}

	// 更新项目 / Update project
	retrieved.Name = "更新后的项目名称"
	err = repo.Update(retrieved)
	if err != nil {
		t.Fatalf("Failed to update project: %v", err)
	}

	updated, _ := repo.GetByID(project.ID)
	if updated.Name != "更新后的项目名称" {
		t.Errorf("Update failed: expected '更新后的项目名称', got %s", updated.Name)
	}

	// 删除项目 / Delete project
	err = repo.Delete(project.ID)
	if err != nil {
		t.Fatalf("Failed to delete project: %v", err)
	}

	// 验证删除（软删除） / Verify deletion (soft delete)
	var count int64
	db.Unscoped().Model(&model.DesignerProject{}).Where("id = ?", project.ID).Count(&count)
	if count != 1 {
		t.Error("Soft delete should keep the record")
	}
}

// TestProjectListWithFilter 测试带筛选条件的列表查询
// Test list query with filters
func TestProjectListWithFilter(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDesignerProjectRepository(db)

	// 创建多个测试项目 / Create multiple test projects
	projects := []model.DesignerProject{
		{Name: "办公室设计A", Status: "draft", Style: "现代简约", Budget: 100000, UserID: 1},
		{Name: "办公室设计B", Status: "in_progress", Style: "工业风", Budget: 200000, UserID: 1},
		{Name: "商业空间C", Status: "completed", Style: "现代简约", Budget: 300000, UserID: 1},
		{Name: "工厂设计D", Status: "draft", Style: "工业风", Budget: 150000, UserID: 2},
	}

	for i := range projects {
		if err := db.Create(&projects[i]).Error; err != nil {
			t.Fatalf("Failed to create test project: %v", err)
		}
		// 添加小延迟确保创建时间不同 / Add small delay to ensure different creation times
		time.Sleep(10 * time.Millisecond)
	}

	// 测试按用户ID筛选 / Test filter by user ID
	filter := ProjectFilter{UserID: 1}
	results, total, err := repo.List(1, 10, filter)
	if err != nil {
		t.Fatalf("Failed to list projects: %v", err)
	}
	if total != 3 {
		t.Errorf("Expected 3 projects for user 1, got %d", total)
	}

	// 测试按状态筛选 / Test filter by status
	filter = ProjectFilter{UserID: 1, Status: "draft"}
	results, total, err = repo.List(1, 10, filter)
	if err != nil {
		t.Fatalf("Failed to list projects: %v", err)
	}
	if total != 1 {
		t.Errorf("Expected 1 draft project for user 1, got %d", total)
	}

	// 测试按风格筛选 / Test filter by style
	filter = ProjectFilter{UserID: 1, Style: "现代简约"}
	results, total, err = repo.List(1, 10, filter)
	if err != nil {
		t.Fatalf("Failed to list projects: %v", err)
	}
	if total != 2 {
		t.Errorf("Expected 2 modern style projects for user 1, got %d", total)
	}

	// 测试按预算范围筛选 / Test filter by budget range
	filter = ProjectFilter{UserID: 1, MinBudget: 150000, MaxBudget: 250000}
	results, total, err = repo.List(1, 10, filter)
	if err != nil {
		t.Fatalf("Failed to list projects: %v", err)
	}
	if total != 1 {
		t.Errorf("Expected 1 project in budget range, got %d", total)
	}

	// 测试关键字搜索 / Test keyword search
	filter = ProjectFilter{UserID: 1, Keyword: "办公室"}
	results, total, err = repo.List(1, 10, filter)
	if err != nil {
		t.Fatalf("Failed to list projects: %v", err)
	}
	if total != 2 {
		t.Errorf("Expected 2 projects with keyword '办公室', got %d", total)
	}

	_ = results // 避免未使用变量警告 / Avoid unused variable warning
}

// **Feature: designer-ai-assistant, Property 7: Project Search Result Relevance**
// **Validates: Requirements 4.3**
// For any project search query, all returned projects SHALL match at least one of
// the search criteria (name, keyword, or time range).
func TestProjectSearchResultRelevance(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// 生成搜索关键字 / Generate search keyword
	genSearchKeyword := gen.OneConstOf("办公", "商业", "工厂", "设计", "空间", "项目")

	properties.Property("search results match search criteria", prop.ForAll(
		func(keyword string) bool {
			// 为每个测试用例创建新的数据库实例 / Create new DB instance for each test case
			db := setupTestDBForProperty(t)
			repo := NewDesignerProjectRepository(db)

			// 创建测试数据 / Create test data
			testProjects := []model.DesignerProject{
				{Name: "办公室设计项目A", Description: "现代办公空间设计", Status: "draft", UserID: 1},
				{Name: "商业空间设计B", Description: "购物中心内部设计", Status: "in_progress", UserID: 1},
				{Name: "工厂车间改造C", Description: "工业厂房改造项目", Status: "completed", UserID: 1},
				{Name: "酒店大堂设计D", Description: "五星级酒店大堂", Status: "draft", UserID: 1},
				{Name: "餐厅空间E", Description: "中式餐厅设计方案", Status: "draft", UserID: 1},
				{Name: "医院门诊F", Description: "医疗空间设计", Status: "draft", UserID: 2},
			}

			for i := range testProjects {
				if err := db.Create(&testProjects[i]).Error; err != nil {
					t.Logf("Failed to create test project: %v", err)
					return false
				}
			}

			// 执行搜索 / Execute search
			results, _, err := repo.Search(keyword, 1, 1, 100)
			if err != nil {
				t.Logf("Search failed: %v", err)
				return false
			}

			// 验证所有返回结果都匹配搜索条件 / Verify all results match search criteria
			for _, project := range results {
				// 检查是否属于正确的用户 / Check if belongs to correct user
				if project.UserID != 1 {
					t.Logf("Result contains project from wrong user: %d", project.UserID)
					return false
				}

				// 检查名称或描述是否包含关键字 / Check if name or description contains keyword
				nameMatch := containsKeyword(project.Name, keyword)
				descMatch := containsKeyword(project.Description, keyword)

				if !nameMatch && !descMatch {
					t.Logf("Result '%s' does not match keyword '%s'", project.Name, keyword)
					return false
				}
			}

			return true
		},
		genSearchKeyword,
	))

	properties.TestingRun(t)
}

// containsKeyword 检查字符串是否包含关键字
// Check if string contains keyword
func containsKeyword(s, keyword string) bool {
	return len(keyword) > 0 && len(s) > 0 &&
		(contains(s, keyword) || contains(keyword, s))
}

// contains 简单的字符串包含检查
// Simple string contains check
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestProjectSearchWithTimeRange 测试按时间范围搜索
// Test search with time range
func TestProjectSearchWithTimeRange(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 50
	properties := gopter.NewProperties(parameters)

	// 生成天数偏移量 / Generate day offset
	genDayOffset := gen.IntRange(1, 30)

	properties.Property("time range filter returns projects within range", prop.ForAll(
		func(dayOffset int) bool {
			db := setupTestDBForProperty(t)
			repo := NewDesignerProjectRepository(db)

			now := time.Now()
			startDate := now.AddDate(0, 0, -dayOffset)
			endDate := now

			// 创建测试项目 / Create test projects
			project := &model.DesignerProject{
				Name:   "测试项目",
				Status: "draft",
				UserID: 1,
			}
			if err := db.Create(project).Error; err != nil {
				t.Logf("Failed to create project: %v", err)
				return false
			}

			// 使用时间范围筛选 / Filter with time range
			filter := ProjectFilter{
				UserID:    1,
				StartDate: startDate,
				EndDate:   endDate.Add(24 * time.Hour), // 包含今天 / Include today
			}

			results, total, err := repo.List(1, 10, filter)
			if err != nil {
				t.Logf("List failed: %v", err)
				return false
			}

			// 验证返回的项目在时间范围内 / Verify returned projects are within time range
			for _, p := range results {
				if p.CreatedAt.Before(startDate) || p.CreatedAt.After(endDate.Add(24*time.Hour)) {
					t.Logf("Project created at %v is outside range [%v, %v]",
						p.CreatedAt, startDate, endDate)
					return false
				}
			}

			// 刚创建的项目应该在结果中 / Newly created project should be in results
			if total < 1 {
				t.Logf("Expected at least 1 project, got %d", total)
				return false
			}

			return true
		},
		genDayOffset,
	))

	properties.TestingRun(t)
}

// TestProjectSearchCombinedFilters 测试组合筛选条件
// Test combined filter criteria
func TestProjectSearchCombinedFilters(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 50
	properties := gopter.NewProperties(parameters)

	// 组合生成器 / Combined generator
	testDataGen := gopter.CombineGens(
		genProjectStatus(),
		genProjectStyle(),
		genPositiveFloat(), // minBudget
		genPositiveFloat(), // maxBudget
	)

	properties.Property("combined filters return matching projects only", prop.ForAll(
		func(testData []any) bool {
			db := setupTestDBForProperty(t)
			repo := NewDesignerProjectRepository(db)

			status := testData[0].(string)
			style := testData[1].(string)
			minBudget := testData[2].(float64)
			maxBudget := testData[3].(float64)

			// 确保 minBudget <= maxBudget / Ensure minBudget <= maxBudget
			if minBudget > maxBudget {
				minBudget, maxBudget = maxBudget, minBudget
			}

			// 创建多个测试项目 / Create multiple test projects
			testProjects := []model.DesignerProject{
				{Name: "项目1", Status: status, Style: style, Budget: (minBudget + maxBudget) / 2, UserID: 1},
				{Name: "项目2", Status: "archived", Style: "其他风格", Budget: maxBudget + 1000, UserID: 1},
				{Name: "项目3", Status: status, Style: style, Budget: minBudget - 1, UserID: 1},
			}

			for i := range testProjects {
				if err := db.Create(&testProjects[i]).Error; err != nil {
					t.Logf("Failed to create project: %v", err)
					return false
				}
			}

			// 应用组合筛选 / Apply combined filters
			filter := ProjectFilter{
				UserID:    1,
				Status:    status,
				Style:     style,
				MinBudget: minBudget,
				MaxBudget: maxBudget,
			}

			results, _, err := repo.List(1, 10, filter)
			if err != nil {
				t.Logf("List failed: %v", err)
				return false
			}

			// 验证所有结果都匹配所有筛选条件 / Verify all results match all filter criteria
			for _, p := range results {
				if p.Status != status {
					t.Logf("Status mismatch: expected %s, got %s", status, p.Status)
					return false
				}
				if p.Style != style {
					t.Logf("Style mismatch: expected %s, got %s", style, p.Style)
					return false
				}
				if p.Budget < minBudget || p.Budget > maxBudget {
					t.Logf("Budget %f outside range [%f, %f]", p.Budget, minBudget, maxBudget)
					return false
				}
				if p.UserID != 1 {
					t.Logf("UserID mismatch: expected 1, got %d", p.UserID)
					return false
				}
			}

			return true
		},
		testDataGen,
	))

	properties.TestingRun(t)
}
