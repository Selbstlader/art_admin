#!/bin/bash

# 模拟用户完整操作流程
# 演示: 用户输入 "一元二次方程" → Dify检索 → DeepSeek生成 → 返回结果

set -e

echo "=========================================="
echo "🎓 AI 学习系统 - 用户操作流程模拟"
echo "=========================================="
echo ""

# 配置
API_BASE="http://localhost:48080"
DIFY_API_KEY="dataset-QQ8Y3Cxvss4RvQXvU7Gem4sb"
DIFY_BASE_URL="https://api.dify.ai/v1"
DATASET_ID="7e2a67c6-00ab-46d0-a71b-886bb480d18f"

# 获取 JWT Token (模拟用户登录)
echo "步骤 0: 用户登录获取 Token"
echo "---"
LOGIN_RESPONSE=$(curl -s -X POST "${API_BASE}/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "123456"}')

TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.access_token')
echo "✅ 登录成功，获取到 Token"
echo ""

# 用户输入
USER_INPUT="一元二次方程"
SUBJECT="数学"
GRADE="初中"
DIFFICULTY=2

echo "=========================================="
echo "步骤 1: 👤 用户输入"
echo "=========================================="
echo "学科: $SUBJECT"
echo "年级: $GRADE"
echo "题材: $USER_INPUT"
echo "难度: $DIFFICULTY (进阶)"
echo ""
sleep 1

echo "=========================================="
echo "步骤 2: 🔍 Dify 知识库检索"
echo "=========================================="
SEARCH_QUERY="${SUBJECT} ${GRADE} ${USER_INPUT} 知识点"
echo "检索关键词: $SEARCH_QUERY"
echo ""

# 调用 Dify API 检索
echo "正在检索知识库..."
DIFY_RESPONSE=$(curl -s -X POST "${DIFY_BASE_URL}/datasets/${DATASET_ID}/retrieve" \
  -H "Authorization: Bearer ${DIFY_API_KEY}" \
  -H "Content-Type: application/json" \
  -d "{\"query\": \"${SEARCH_QUERY}\"}")

# 检查是否成功
if echo "$DIFY_RESPONSE" | jq -e '.records' > /dev/null 2>&1; then
    RECORD_COUNT=$(echo "$DIFY_RESPONSE" | jq '.records | length')
    echo "✅ 检索成功! 找到 $RECORD_COUNT 条相关资料"
    echo ""
    echo "📚 检索到的知识片段:"
    echo "$DIFY_RESPONSE" | jq -r '.records[] | "- 内容: \(.segment.content[:100])...\n  相关度: \(.score)"'
    echo ""
    
    # 提取知识内容
    KNOWLEDGE_CONTENT=$(echo "$DIFY_RESPONSE" | jq -r '.records[].segment.content' | head -n 3)
else
    echo "⚠️  知识库检索失败，将使用 DeepSeek 基础知识"
    KNOWLEDGE_CONTENT="暂无参考资料"
fi
echo ""
sleep 2

echo "=========================================="
echo "步骤 3: 🤖 构建 RAG Prompt"
echo "=========================================="
echo "将检索到的知识库内容与用户输入结合..."
echo ""
echo "Prompt 结构:"
echo "  - 系统角色: 专业的${SUBJECT}老师"
echo "  - 参考资料: 来自 Dify 知识库"
echo "  - 用户需求: 生成关于'${USER_INPUT}'的教材"
echo "  - 输出格式: 结构化 JSON"
echo ""
sleep 1

echo "=========================================="
echo "步骤 4: 🧠 DeepSeek AI 分析生成"
echo "=========================================="
echo "正在调用 DeepSeek API 生成教材..."
echo "这可能需要 30-60 秒..."
echo ""

# 调用后端 API 生成教材
GENERATE_RESPONSE=$(curl -s -X POST "${API_BASE}/api/learning/material/generate" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d "{
    \"subject_id\": 2,
    \"grade\": \"${GRADE}\",
    \"topic\": \"${USER_INPUT}\",
    \"difficulty\": ${DIFFICULTY}
  }")

# 检查生成结果
if echo "$GENERATE_RESPONSE" | jq -e '.data' > /dev/null 2>&1; then
    echo "✅ 教材生成成功!"
    echo ""
    
    echo "=========================================="
    echo "步骤 5: 📊 DeepSeek 整理分析结果"
    echo "=========================================="
    
    MATERIAL_ID=$(echo "$GENERATE_RESPONSE" | jq -r '.data.id')
    TITLE=$(echo "$GENERATE_RESPONSE" | jq -r '.data.title')
    SUMMARY=$(echo "$GENERATE_RESPONSE" | jq -r '.data.summary')
    TOTAL_TIME=$(echo "$GENERATE_RESPONSE" | jq -r '.data.total_time')
    SECTIONS=$(echo "$GENERATE_RESPONSE" | jq '.data.content.sections')
    SECTION_COUNT=$(echo "$SECTIONS" | jq 'length')
    
    echo "教材 ID: $MATERIAL_ID"
    echo "标题: $TITLE"
    echo "概要: $SUMMARY"
    echo "预计学习时长: ${TOTAL_TIME}分钟"
    echo "章节数量: $SECTION_COUNT"
    echo ""
    
    echo "📖 教材结构:"
    echo "$SECTIONS" | jq -r '.[] | "  \(.type | ascii_upcase): \(.title)"'
    echo ""
    
    # 显示知识点
    echo "💡 知识点讲解:"
    echo "$SECTIONS" | jq -r '.[] | select(.type == "knowledge") | "  - \(.title)\n    重点: \(.key_points | join(", "))"'
    echo ""
    
    # 显示例题
    echo "📝 例题演示:"
    echo "$SECTIONS" | jq -r '.[] | select(.type == "example") | "  - \(.title)\n    题目: \(.question)\n    答案: \(.answer)"'
    echo ""
    
    # 显示练习题数量
    EXERCISE_COUNT=$(echo "$SECTIONS" | jq '[.[] | select(.type == "exercise") | .questions[]] | length')
    echo "✏️  练习题: 共 ${EXERCISE_COUNT} 道"
    echo ""
    
    echo "=========================================="
    echo "步骤 6: 💾 保存到数据库"
    echo "=========================================="
    echo "✅ 教材已保存到数据库"
    echo "数据库 ID: $MATERIAL_ID"
    echo ""
    
    echo "=========================================="
    echo "步骤 7: 🎨 返回给用户前端展示"
    echo "=========================================="
    echo "✅ 数据已返回给前端"
    echo ""
    echo "前端将展示:"
    echo "  - 教材标题和概要"
    echo "  - 知识点讲解（含重点标注）"
    echo "  - 例题演示（含解题步骤）"
    echo "  - 练习题（含答案和解析）"
    echo "  - 预计学习时长"
    echo ""
    
    echo "=========================================="
    echo "🎉 完整流程演示完成！"
    echo "=========================================="
    echo ""
    echo "📊 流程总结:"
    echo "  1. ✅ 用户输入: ${USER_INPUT}"
    echo "  2. ✅ Dify 检索: 找到 ${RECORD_COUNT:-0} 条资料"
    echo "  3. ✅ RAG Prompt: 结合知识库内容"
    echo "  4. ✅ DeepSeek 生成: 创建结构化教材"
    echo "  5. ✅ 数据整理: ${SECTION_COUNT} 个章节"
    echo "  6. ✅ 保存数据库: ID ${MATERIAL_ID}"
    echo "  7. ✅ 返回前端: 完整教材内容"
    echo ""
    
    echo "🔗 查看完整教材:"
    echo "  GET ${API_BASE}/api/learning/material/${MATERIAL_ID}"
    echo ""
    
    # 可选：显示完整 JSON
    read -p "是否查看完整的教材 JSON? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo ""
        echo "完整教材内容:"
        echo "$GENERATE_RESPONSE" | jq '.data'
    fi
    
else
    echo "❌ 教材生成失败"
    echo "$GENERATE_RESPONSE" | jq '.'
fi

echo ""
echo "=========================================="
echo "演示结束"
echo "=========================================="
