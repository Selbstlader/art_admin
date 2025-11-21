#!/bin/bash

# 测试脚本 - 直接调用 Dify API 验证配置

echo "======================================"
echo "测试 Dify API 配置"
echo "======================================"
echo ""

DATASET_API_KEY="dataset-QQ8Y3Cxvss4RvQXvU7Gem4sb"
CHAT_API_KEY="app-D7hSHUnPoMt5CQD2IEoBsPJz"

echo "1. 测试知识库 API"
echo "--------------------------------------"
curl -s -X GET 'https://api.dify.ai/v1/datasets?page=1&limit=10' \
  -H "Authorization: Bearer ${DATASET_API_KEY}" \
  -H "Accept: application/json" | jq .
echo ""
echo ""

echo "2. 测试 AI 对话 API (非流式)"
echo "--------------------------------------"
curl -s -X POST 'https://api.dify.ai/v1/chat-messages' \
  -H "Authorization: Bearer ${CHAT_API_KEY}" \
  -H "Content-Type: application/json" \
  -d '{
    "inputs": {
      "report_title": "2024年市场分析报告",
      "report_content": "本报告分析了2024年市场趋势,发现市场增长率达到15%,消费者对环保产品需求增加。",
      "report_date": "2024-11-20",
      "analysis_focus": "市场趋势和消费者行为",
      "birth_date": "1990-01-01"
    },
    "query": "请分析这份报告",
    "response_mode": "blocking",
    "user": "test-user"
  }' | jq .
echo ""
echo ""

echo "======================================"
echo "测试完成"
echo "======================================"
