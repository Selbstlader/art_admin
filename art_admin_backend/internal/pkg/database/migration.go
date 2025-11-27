package database

import (
	"art_admin_backend/internal/pkg/logger"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// RunMigrations 执行数据库迁移
func RunMigrations() error {
	logger.Info("开始执行数据库迁移...")

	// 创建迁移记录表
	if err := createMigrationTable(); err != nil {
		return fmt.Errorf("创建迁移记录表失败: %w", err)
	}

	// 获取迁移文件目录
	migrationDir := "migrations"

	// 读取所有迁移文件
	files, err := os.ReadDir(migrationDir)
	if err != nil {
		return fmt.Errorf("读取迁移目录失败: %w", err)
	}

	// 按文件名排序
	var migrationFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}
	sort.Strings(migrationFiles)

	// 执行每个迁移文件
	for _, fileName := range migrationFiles {
		if err := runMigrationFile(migrationDir, fileName); err != nil {
			return fmt.Errorf("执行迁移文件 %s 失败: %w", fileName, err)
		}
	}

	logger.Info("数据库迁移执行完成")
	return nil
}

// createMigrationTable 创建迁移记录表
func createMigrationTable() error {
	sql := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		id INT AUTO_INCREMENT PRIMARY KEY,
		version VARCHAR(255) NOT NULL UNIQUE,
		applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)
	`

	return DB.Exec(sql).Error
}

// runMigrationFile 执行单个迁移文件
func runMigrationFile(migrationDir, fileName string) error {
	// 检查是否已经执行过
	var count int64
	if err := DB.Raw("SELECT COUNT(*) FROM schema_migrations WHERE version = ?", fileName).Scan(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		logger.Info(fmt.Sprintf("迁移文件 %s 已执行，跳过", fileName))
		return nil
	}

	// 读取迁移文件内容
	filePath := filepath.Join(migrationDir, fileName)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// 按分号分割SQL语句并逐个执行
	statements := strings.Split(string(content), ";")

	executedCount := 0
	for i, stmt := range statements {
		stmt = strings.TrimSpace(stmt)

		// 跳过空语句和纯注释
		if stmt == "" || (!strings.Contains(strings.ToUpper(stmt), "CREATE") &&
			!strings.Contains(strings.ToUpper(stmt), "INSERT") &&
			!strings.Contains(strings.ToUpper(stmt), "DROP") &&
			!strings.Contains(strings.ToUpper(stmt), "ALTER") &&
			!strings.Contains(strings.ToUpper(stmt), "UPDATE") &&
			!strings.Contains(strings.ToUpper(stmt), "DELETE")) {
			continue
		}

		// 执行单个SQL语句
		if err := DB.Exec(stmt).Error; err != nil {
			return fmt.Errorf("执行第%d个SQL语句失败: %w", i+1, err)
		}

		executedCount++
	}

	logger.Info(fmt.Sprintf("迁移文件 %s 执行完成，共执行%d个SQL语句", fileName, executedCount))

	// 记录迁移执行
	if err := DB.Exec("INSERT INTO schema_migrations (version) VALUES (?)", fileName).Error; err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("成功执行迁移文件: %s", fileName))
	return nil
}
