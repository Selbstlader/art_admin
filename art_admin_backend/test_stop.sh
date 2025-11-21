#!/bin/bash

TOKEN="Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOjEsImlzcyI6ImFydC1hZG1pbiIsImV4cCI6MTc2MzYzNzgwOSwiaWF0IjoxNzYzNjMwNjA5fQ.bWsbDKjSQcfgCLgmNvPwlzg4t4M_YONM9JB_nr3r-3c"

echo "开始流式对话..."
curl -X POST 'http://localhost:48080/api/dify/chat/stream' \
  -H "Authorization: ${TOKEN}" \
  -H 'Content-Type: application/json' \
  -d '{
    "query": "请详细分析这份长篇报告",
    "user": "test-user-002",
    "inputs": {
      "report_title": "长篇测试报告",
      "report_content": "这是一份非常详细的测试报告,包含大量的市场分析数据和趋势预测...",
      "report_date": "2024-11-20",
      "analysis_focus": "详细分析"
    }
  }' &

# 等待1秒后尝试停止(实际使用中需要从响应中获取task_id)
sleep 1
echo -e "\n\n尝试停止消息生成..."
# 注意: 这里的task_id需要从实际响应中获取
# curl -X POST "http://localhost:48080/api/dify/chat/stop/{task_id}?user=test-user-002" \
#   -H "Authorization: ${TOKEN}"

wait
