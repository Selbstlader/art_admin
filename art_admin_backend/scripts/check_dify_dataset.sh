#!/bin/bash

# Dify 知识库诊断工具
# 用于检查知识库状态和可用性

set -e

echo "=========================================="
echo "Dify 知识库诊断工具"
echo "=========================================="
echo ""

# 配置
DIFY_API_KEY="${DIFY_API_KEY:-dataset-QQ8Y3Cxvss4RvQXvU7Gem4sb}"
DIFY_BASE_URL="${DIFY_BASE_URL:-https://api.dify.ai/v1}"

echo "1. 获取知识库列表..."
echo ""

RESPONSE=$(curl -s -X GET "${DIFY_BASE_URL}/datasets" \
  -H "Authorization: Bearer ${DIFY_API_KEY}" \
  -H "Accept: application/json")

echo "$RESPONSE" | jq -r '.data[] | "知识库ID: \(.id)\n名称: \(.name)\n文档总数: \(.document_count)\n可用文档数: \(.total_available_documents)\n词数: \(.word_count)\n索引状态: \(if .total_available_documents > 0 then "✅ 可用" else "⚠️  未完成索引" end)\n---"'

echo ""
echo "=========================================="
echo "2. 推荐使用的知识库"
echo "=========================================="
echo ""

# 找出可用的知识库
AVAILABLE_DATASETS=$(echo "$RESPONSE" | jq -r '.data[] | select(.total_available_documents > 0) | .id')

if [ -z "$AVAILABLE_DATASETS" ]; then
    echo "⚠️  警告: 没有找到已完成索引的知识库！"
    echo ""
    echo "建议:"
    echo "1. 登录 Dify 平台: https://cloud.dify.ai"
    echo "2. 检查知识库文档的索引状态"
    echo "3. 等待文档索引完成"
    echo "4. 或者创建新的知识库并上传文档"
else
    echo "✅ 找到可用的知识库:"
    echo ""
    for dataset_id in $AVAILABLE_DATASETS; do
        DATASET_INFO=$(echo "$RESPONSE" | jq -r ".data[] | select(.id == \"$dataset_id\") | \"ID: \(.id)\n名称: \(.name)\n可用文档: \(.total_available_documents)/\(.document_count)\"")
        echo "$DATASET_INFO"
        echo ""
        
        # 测试检索
        echo "测试检索..."
        TEST_RESULT=$(curl -s -X POST "${DIFY_BASE_URL}/datasets/${dataset_id}/retrieve" \
          -H "Authorization: Bearer ${DIFY_API_KEY}" \
          -H "Content-Type: application/json" \
          -d '{"query": "测试"}')
        
        if echo "$TEST_RESULT" | jq -e '.records' > /dev/null 2>&1; then
            RECORD_COUNT=$(echo "$TEST_RESULT" | jq '.records | length')
            echo "✅ 检索成功! 返回 $RECORD_COUNT 条记录"
            echo ""
            echo "📝 建议配置:"
            echo "在 config/config.yaml 中设置:"
            echo "dify:"
            echo "  datasetId: $dataset_id"
        else
            echo "❌ 检索失败"
            echo "$TEST_RESULT" | jq '.code, .message'
        fi
        echo ""
        echo "---"
        echo ""
    done
fi

echo ""
echo "=========================================="
echo "诊断完成"
echo "=========================================="
