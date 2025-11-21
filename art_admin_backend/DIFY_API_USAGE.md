# Dify API 使用说明

## API 接口

### 1. 获取知识库列表

```bash
curl -X GET 'http://localhost:48080/api/dify/dataset/list?page=1&limit=10' \
  -H 'Authorization: Bearer {your-jwt-token}'
```

### 2. 获取知识库详情

```bash
curl -X GET 'http://localhost:48080/api/dify/dataset/{dataset-id}' \
  -H 'Authorization: Bearer {your-jwt-token}'
```

### 3. 上传文件到知识库

```bash
curl -X POST 'http://localhost:48080/api/dify/dataset/upload' \
  -H 'Authorization: Bearer {your-jwt-token}' \
  -F 'dataset_id={dataset-id}' \
  -F 'file=@/path/to/your/file.pdf'
```

### 4. AI 对话 (非流式)

**重要**: 你的 Dify 应用配置了以下输入变量,调用时必须提供:

```bash
curl -X POST 'http://localhost:48080/api/dify/chat' \
  -H 'Authorization: Bearer {your-jwt-token}' \
  -H 'Content-Type: application/json' \
  -d '{
    "query": "请分析这份报告",
    "user": "user-123",
    "inputs": {
      "report_title": "2024年市场分析报告",
      "report_content": "本报告分析了2024年市场趋势。主要发现：1) 市场年增长率达到15%，存在显著扩张机会；2) 消费者对环保产品的需求在近两年内翻倍；3) 新技术应用可降低30%的运营成本。同时也存在一些挑战：原材料价格波动、新竞争者进入市场、监管政策变化等。",
      "report_date": "2024-11-20",
      "analysis_focus": "市场趋势、消费者行为和风险评估"
    }
  }'
```

**输入变量说明**:
- `report_title`: 报告标题 (必填)
- `report_content`: 报告内容 (必填)
- `report_date`: 报告日期 (必填)
- `analysis_focus`: 分析重点 (必填)

**注意**: 如果 Dify 应用还配置了其他必填变量(如 `birth_date`, `birth_time` 等),也需要在 `inputs` 中提供。

### 5. AI 对话 (流式)

```bash
curl -X POST 'http://localhost:48080/api/dify/chat/stream' \
  -H 'Authorization: Bearer {your-jwt-token}' \
  -H 'Content-Type: application/json' \
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
```

**流式响应格式** (SSE):
```
data: {"event":"message","answer":"关","message_id":"msg-123"}

data: {"event":"message","answer":"键","message_id":"msg-123"}

data: {"event":"message","answer":"洞","message_id":"msg-123"}

...

data: {"event":"message_end","message_id":"msg-123","conversation_id":"conv-123"}
```

## 完整的请求示例 (JavaScript)

### 非流式对话

```javascript
async function chatWithAI(reportData) {
  const response = await fetch('http://localhost:48080/api/dify/chat', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${jwtToken}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      query: "请分析这份报告",
      user: "user-123",
      inputs: {
        report_title: reportData.title,
        report_content: reportData.content,
        report_date: reportData.date,
        analysis_focus: reportData.focus
      }
    })
  });

  const result = await response.json();
  console.log(result.data.answer);
}
```

### 流式对话

```javascript
async function chatWithAIStreaming(reportData) {
  const response = await fetch('http://localhost:48080/api/dify/chat/stream', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${jwtToken}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      query: "请分析这份报告",
      user: "user-123",
      inputs: {
        report_title: reportData.title,
        report_content: reportData.content,
        report_date: reportData.date,
        analysis_focus: reportData.focus
      }
    })
  });

  const reader = response.body.getReader();
  const decoder = new TextDecoder();

  while (true) {
    const { done, value } = await reader.read();
    if (done) break;

    const chunk = decoder.decode(value);
    const lines = chunk.split('\n');

    for (const line of lines) {
      if (line.startsWith('data: ')) {
        const data = JSON.parse(line.slice(6));
        
        if (data.event === 'message') {
          // 实时显示消息片段
          console.log(data.answer);
        } else if (data.event === 'message_end') {
          // 对话结束
          console.log('Conversation ID:', data.conversation_id);
        }
      }
    }
  }
}
```

## 输出格式

根据你的 Dify 模板配置,AI 会按照以下格式输出:

```
关键洞察：
- 洞察点1
- 洞察点2
...

识别风险：
- 风险点1
- 风险点2
...

核心信息：
- 信息要点1
- 信息要点2
...
```

## 注意事项

1. **输入变量配置**: 不同的 Dify 应用可能配置了不同的输入变量,使用前请在 Dify 控制台确认所有必填字段

2. **JWT Token**: 所有接口都需要有效的 JWT token,可以通过登录接口获取

3. **会话管理**: 
   - 首次对话时 `conversation_id` 留空
   - 后续对话使用返回的 `conversation_id` 保持上下文

4. **文件大小限制**: 上传文件时注意文件大小限制(默认由服务器配置决定)

5. **流式响应**: 
   - 适合长文本生成场景
   - 客户端需要支持 SSE (Server-Sent Events)
   - 可以实时显示生成进度

## 错误处理

常见错误码:
- `400`: 参数错误或缺少必填字段
- `401`: JWT token 无效或过期
- `500`: 服务器错误或 Dify API 调用失败

错误响应示例:
```json
{
  "code": 400,
  "msg": "参数错误: report_title is required",
  "data": null
}
```
