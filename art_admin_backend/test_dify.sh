#!/bin/bash

TOKEN="Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOjEsImlzcyI6ImFydC1hZG1pbiIsImV4cCI6MTc2MzYzNzgwOSwiaWF0IjoxNzYzNjMwNjA5fQ.bWsbDKjSQcfgCLgmNvPwlzg4t4M_YONM9JB_nr3r-3c"
BASE_URL="http://localhost:48080/api/dify"

echo "======================================"
echo "Dify 接口测试"
echo "======================================"
echo ""

echo "1. 测试获取知识库列表"
echo "--------------------------------------"
curl -s -X GET "${BASE_URL}/dataset/list?page=1&limit=10" \
  -H "Authorization: ${TOKEN}" \
  -H "Accept: application/json" | jq .
echo ""
echo ""

echo "2. 测试获取知识库详情"
echo "--------------------------------------"
DATASET_ID="711b1a11-2a29-4920-a379-6ee6251a97fa"
curl -s -X GET "${BASE_URL}/dataset/${DATASET_ID}" \
  -H "Authorization: ${TOKEN}" \
  -H "Accept: application/json" | jq .
echo ""
echo ""

echo "3. 测试 AI 对话 (非流式)"
echo "--------------------------------------"
curl -s -X POST "${BASE_URL}/chat" \
  -H "Authorization: ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "请分析这份报告",
    "user": "user-123",
    "inputs": {
      "report_title": "2024年市场分析报告",
      "report_content": "本报告分析了2024年市场趋势。主要发现：1) 市场年增长率达到15%，存在显著扩张机会；2) 消费者对环保产品的需求在近两年内翻倍；3) 新技术应用可降低30%的运营成本。同时也存在一些挑战：原材料价格波动、新竞争者进入市场、监管政策变化等。",
      "report_date": "2024-11-20",
      "analysis_focus": "市场趋势、消费者行为和风险评估"
    }
  }' | jq .
echo ""
echo ""

echo "4. 测试 AI 对话 (流式)"
echo "--------------------------------------"
curl -X POST "${BASE_URL}/chat/stream" \
  -H "Authorization: ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "请分析这份技术报告",
    "user": "user-123",
    "inputs": {
      "report_title": "AI技术发展报告",
      "report_content": "人工智能技术在2024年取得重大突破。大语言模型性能提升50%，应用场景扩展到医疗、教育、金融等多个领域。但也面临数据隐私、算力成本、技术伦理等挑战。",
      "report_date": "2024-11-20",
      "analysis_focus": "技术突破和应用前景"
    }
  }'
echo ""
echo ""

echo "======================================"
echo "测试完成"
echo "======================================"
