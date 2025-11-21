#!/bin/bash

# AI 学习系统初始化脚本
# 用于初始化数据库表和默认数据

set -e

echo "=========================================="
echo "AI 学习系统初始化脚本"
echo "=========================================="

# 数据库配置(从 config.yaml 读取或使用默认值)
DB_HOST=${DB_HOST:-"127.0.0.1"}
DB_PORT=${DB_PORT:-"3306"}
DB_USER=${DB_USER:-"root"}
DB_PASSWORD=${DB_PASSWORD:-"123456"}
DB_NAME=${DB_NAME:-"gin_admin"}

echo ""
echo "数据库配置:"
echo "  主机: $DB_HOST"
echo "  端口: $DB_PORT"
echo "  用户: $DB_USER"
echo "  数据库: $DB_NAME"
echo ""

# 检查 MySQL 是否可连接
echo "检查数据库连接..."
if ! mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" -e "SELECT 1" > /dev/null 2>&1; then
    echo "错误: 无法连接到数据库"
    echo "请检查数据库配置是否正确"
    exit 1
fi
echo "✓ 数据库连接成功"

# 执行 SQL 脚本
echo ""
echo "执行数据库迁移脚本..."
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" < "$(dirname "$0")/migrations/learning_system.sql"

if [ $? -eq 0 ]; then
    echo "✓ 数据库表创建成功"
else
    echo "✗ 数据库表创建失败"
    exit 1
fi

echo ""
echo "=========================================="
echo "初始化完成!"
echo "=========================================="
echo ""
echo "已创建的表:"
echo "  - subjects (学科表)"
echo "  - learning_materials (教材内容表)"
echo "  - material_sections (教材章节表)"
echo "  - audio_cache (音频缓存表)"
echo ""
echo "已插入的默认数据:"
echo "  - 9 个学科(语文、数学、英语、物理、化学、生物、历史、地理、政治)"
echo ""
echo "下一步:"
echo "  1. 启动后端服务: go run cmd/server/main.go"
echo "  2. 访问 API 文档: http://localhost:48080/swagger/index.html"
echo "  3. 测试教材生成接口"
echo ""
