#!/bin/bash

echo "======================================"
echo "Art Admin Backend 启动脚本"
echo "======================================"
echo ""

# 设置 Go 代理
echo "1. 设置 Go 代理..."
export GOPROXY=https://goproxy.io,direct

# 检查依赖
echo "2. 检查并下载依赖..."
go mod tidy

# 生成 Swagger 文档
echo "3. 生成 Swagger 文档..."
if command -v swag &> /dev/null; then
    swag init -g cmd/server/main.go -o docs
    echo "   ✅ Swagger 文档生成成功"
else
    echo "   ⚠️  swag 命令不存在，将在启动时自动生成"
    echo "   建议安装: go install github.com/swaggo/swag/cmd/swag@latest"
fi

echo ""
echo "4. 检查端口占用..."
# 检查 48080 端口是否被占用
if lsof -ti:48080 > /dev/null 2>&1; then
    echo "   ⚠️  端口 48080 已被占用，正在关闭..."
    lsof -ti:48080 | xargs kill -9 2>/dev/null
    sleep 1
    echo "   ✅ 端口已释放"
else
    echo "   ✅ 端口 48080 可用"
fi

echo ""
echo "🚀 5. 启动服务器（GORM自动管理数据库）..."
echo "   - 自动创建数据库（如果不存在）"
echo "   - 自动创建表结构"
echo "   - 自动初始化数据"
echo "   - 监听端口: 48080"
echo "   - Swagger: http://localhost:48080/swagger/index.html"
echo ""
echo "   📝 测试账号："
echo "   - 用户名: admin"
echo "   - 密码: 123456"
echo ""
echo "======================================"
echo ""

# 运行服务器
go run cmd/server/main.go

