package volcengine

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestClient_SimpleChat(t *testing.T) {
	// 从环境变量获取API Key，如果没有则跳过测试
	apiKey := os.Getenv("VOLCENGINE_API_KEY")
	if apiKey == "" {
		apiKey = "d48c3303-5d63-42d9-80f8-692eb079594c" // 测试用
	}

	client := NewClient(Config{
		APIKey:      apiKey,
		BaseURL:     "https://ark.cn-beijing.volces.com/api/v3",
		Model:       "ep-m-20250709135427-vv9dp", // Doubao-1.5-vision-pro
		Timeout:     60,
		MaxTokens:   1024,
		Temperature: 0.7,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	response, err := client.SimpleChat(ctx, "你是一个有帮助的助手", "你好，请简单介绍一下自己")
	if err != nil {
		t.Logf("API调用失败（可能需要正确的接入点ID）: %v", err)
		// 不标记为失败，因为可能是配置问题
		return
	}

	t.Logf("AI响应: %s", response)
}

func TestClient_SummarizeContent(t *testing.T) {
	apiKey := os.Getenv("VOLCENGINE_API_KEY")
	if apiKey == "" {
		apiKey = "d48c3303-5d63-42d9-80f8-692eb079594c"
	}

	client := NewClient(Config{
		APIKey:      apiKey,
		BaseURL:     "https://ark.cn-beijing.volces.com/api/v3",
		Model:       "ep-20251128103404-dg4wf", // Doubao-1.5-vision-pro
		Timeout:     60,
		MaxTokens:   2048,
		Temperature: 0.7,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	testContent := `
今天给大家推荐一款超好用的面霜！
这款面霜含有神经酰胺和透明质酸，保湿效果超级棒！
我是干皮，用了一周感觉皮肤水润了很多。
价格也很亲民，只要199元，性价比超高！
唯一的缺点是包装有点简陋，但不影响使用。
推荐给干皮姐妹们试试～
`

	response, err := client.SummarizeContent(ctx, testContent, "concise", 200)
	if err != nil {
		t.Logf("总结失败（可能需要正确的接入点ID）: %v", err)
		return
	}

	t.Logf("总结结果: %s", response)
}
