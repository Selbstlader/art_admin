package volcengine

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

/**
 * Feature: designer-ai-assistant, Property 18: AI Retry Mechanism
 * Validates: Requirements 13.2
 *
 * 属性测试：验证AI服务重试机制
 * Property test: Verify AI service retry mechanism
 *
 * 对于任何AI服务调用失败，系统应最多重试3次，间隔递增，然后返回错误
 * For any AI service call that fails, the system SHALL retry up to 3 times
 * with increasing intervals before returning an error.
 */

// TestRetryMechanismProperty 测试重试机制属性
// Test retry mechanism property
func TestRetryMechanismProperty(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 50

	properties := gopter.NewProperties(parameters)

	// Property 1: 重试次数不超过配置的最大值
	// Property 1: Retry count does not exceed configured maximum
	properties.Property("retry count does not exceed max retries", prop.ForAll(
		func(maxRetries int) bool {
			if maxRetries < 0 || maxRetries > 10 {
				return true // 跳过无效值 / Skip invalid values
			}

			var requestCount int32 = 0

			// 创建一个总是返回500错误的测试服务器
			// Create a test server that always returns 500 error
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				atomic.AddInt32(&requestCount, 1)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": {"message": "server error", "type": "server_error", "code": "500"}}`))
			}))
			defer server.Close()

			// 创建客户端，使用较短的超时和间隔以加快测试
			// Create client with short timeout and intervals to speed up test
			client := NewClientWithConfig(
				"test-api-key",
				server.URL,
				"test-model",
				5, // timeout
				100,
				0.7,
				RetryConfig{
					MaxRetries:      maxRetries,
					InitialInterval: 10 * time.Millisecond,
					MaxInterval:     50 * time.Millisecond,
					Multiplier:      2.0,
				},
			)
			client.SetEnableLogging(false)

			// 发送请求
			// Send request
			_, err := client.Chat([]ChatMessage{
				{Role: "user", Content: []ContentPart{{Type: "text", Text: "test"}}},
			})

			// 验证：请求次数应该等于 1（初始请求）+ maxRetries（重试次数）
			// Verify: request count should equal 1 (initial) + maxRetries (retries)
			expectedRequests := int32(1 + maxRetries)
			actualRequests := atomic.LoadInt32(&requestCount)

			if err == nil {
				return false // 应该返回错误 / Should return error
			}

			return actualRequests == expectedRequests
		},
		gen.IntRange(0, 5),
	))

	// Property 2: 重试间隔递增
	// Property 2: Retry intervals increase
	properties.Property("retry intervals increase with each attempt", prop.ForAll(
		func(initialInterval, multiplier int) bool {
			if initialInterval < 1 || initialInterval > 100 || multiplier < 1 || multiplier > 5 {
				return true // 跳过无效值 / Skip invalid values
			}

			config := RetryConfig{
				MaxRetries:      3,
				InitialInterval: time.Duration(initialInterval) * time.Millisecond,
				MaxInterval:     10 * time.Second,
				Multiplier:      float64(multiplier),
			}

			client := NewClientWithConfig("key", "url", "model", 5, 100, 0.7, config)

			// 计算各次重试的间隔
			// Calculate intervals for each retry
			intervals := make([]time.Duration, 3)
			for i := 0; i < 3; i++ {
				intervals[i] = client.calculateRetryInterval(i)
			}

			// 验证间隔递增（或达到最大值后保持不变）
			// Verify intervals increase (or stay same after reaching max)
			for i := 1; i < len(intervals); i++ {
				if intervals[i] < intervals[i-1] {
					return false // 间隔不应减少 / Interval should not decrease
				}
			}

			return true
		},
		gen.IntRange(1, 100),
		gen.IntRange(1, 5),
	))

	// Property 3: 成功请求不重试
	// Property 3: Successful requests do not retry
	properties.Property("successful requests do not trigger retries", prop.ForAll(
		func(maxRetries int) bool {
			if maxRetries < 1 || maxRetries > 5 {
				return true
			}

			var requestCount int32 = 0

			// 创建一个返回成功响应的测试服务器
			// Create a test server that returns successful response
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				atomic.AddInt32(&requestCount, 1)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{
					"id": "test-id",
					"object": "chat.completion",
					"created": 1234567890,
					"model": "test-model",
					"choices": [{"index": 0, "message": {"role": "assistant", "content": "test response"}, "finish_reason": "stop"}],
					"usage": {"prompt_tokens": 10, "completion_tokens": 5, "total_tokens": 15}
				}`))
			}))
			defer server.Close()

			client := NewClientWithConfig(
				"test-api-key",
				server.URL,
				"test-model",
				5,
				100,
				0.7,
				RetryConfig{
					MaxRetries:      maxRetries,
					InitialInterval: 10 * time.Millisecond,
					MaxInterval:     50 * time.Millisecond,
					Multiplier:      2.0,
				},
			)
			client.SetEnableLogging(false)

			_, err := client.Chat([]ChatMessage{
				{Role: "user", Content: []ContentPart{{Type: "text", Text: "test"}}},
			})

			// 成功请求应该只有1次请求，无重试
			// Successful request should have only 1 request, no retries
			return err == nil && atomic.LoadInt32(&requestCount) == 1
		},
		gen.IntRange(1, 5),
	))

	// Property 4: 4xx错误（除429外）不重试
	// Property 4: 4xx errors (except 429) do not retry
	properties.Property("4xx errors except 429 do not trigger retries", prop.ForAll(
		func(statusCode int) bool {
			// 只测试4xx错误（除429外）
			// Only test 4xx errors (except 429)
			if statusCode < 400 || statusCode >= 500 || statusCode == 429 {
				return true
			}

			var requestCount int32 = 0

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				atomic.AddInt32(&requestCount, 1)
				w.WriteHeader(statusCode)
				w.Write([]byte(`{"error": {"message": "client error", "type": "client_error", "code": "400"}}`))
			}))
			defer server.Close()

			client := NewClientWithConfig(
				"test-api-key",
				server.URL,
				"test-model",
				5,
				100,
				0.7,
				RetryConfig{
					MaxRetries:      3,
					InitialInterval: 10 * time.Millisecond,
					MaxInterval:     50 * time.Millisecond,
					Multiplier:      2.0,
				},
			)
			client.SetEnableLogging(false)

			_, err := client.Chat([]ChatMessage{
				{Role: "user", Content: []ContentPart{{Type: "text", Text: "test"}}},
			})

			// 4xx错误不应重试，只有1次请求
			// 4xx errors should not retry, only 1 request
			return err != nil && atomic.LoadInt32(&requestCount) == 1
		},
		gen.IntRange(400, 428), // 400-428, 跳过429
	))

	properties.TestingRun(t)
}

// TestRetryIntervalCalculation 测试重试间隔计算
// Test retry interval calculation
func TestRetryIntervalCalculation(t *testing.T) {
	config := RetryConfig{
		MaxRetries:      3,
		InitialInterval: 1 * time.Second,
		MaxInterval:     10 * time.Second,
		Multiplier:      2.0,
	}

	client := NewClientWithConfig("key", "url", "model", 5, 100, 0.7, config)

	// 测试间隔计算
	// Test interval calculation
	tests := []struct {
		retryCount int
		expected   time.Duration
	}{
		{0, 1 * time.Second},  // 初始间隔 / Initial interval
		{1, 2 * time.Second},  // 1 * 2 = 2
		{2, 4 * time.Second},  // 2 * 2 = 4
		{3, 8 * time.Second},  // 4 * 2 = 8
		{4, 10 * time.Second}, // 8 * 2 = 16, 但限制为10 / but capped at 10
	}

	for _, tt := range tests {
		interval := client.calculateRetryInterval(tt.retryCount)
		if interval != tt.expected {
			t.Errorf("retryCount=%d: expected %v, got %v", tt.retryCount, tt.expected, interval)
		}
	}
}

// TestIsRetryableError 测试可重试错误判断
// Test retryable error detection
func TestIsRetryableError(t *testing.T) {
	tests := []struct {
		statusCode int
		err        error
		expected   bool
	}{
		{500, nil, true}, // 5xx可重试 / 5xx retryable
		{502, nil, true},
		{503, nil, true},
		{429, nil, true},  // 429可重试 / 429 retryable
		{400, nil, false}, // 4xx不可重试 / 4xx not retryable
		{401, nil, false},
		{403, nil, false},
		{404, nil, false},
		{200, nil, false},                      // 成功不需重试 / Success no retry
		{0, errors.New("network error"), true}, // 网络错误可重试 / Network error retryable
	}

	for _, tt := range tests {
		result := isRetryableError(tt.statusCode, tt.err)
		if result != tt.expected {
			t.Errorf("statusCode=%d, err=%v: expected %v, got %v", tt.statusCode, tt.err, tt.expected, result)
		}
	}
}

// TestRetryConfigDefaults 测试默认重试配置
// Test default retry configuration
func TestRetryConfigDefaults(t *testing.T) {
	client := NewClient("key", "url", "model", 5, 100, 0.7)

	config := client.GetRetryConfig()

	if config.MaxRetries != 3 {
		t.Errorf("expected MaxRetries=3, got %d", config.MaxRetries)
	}

	if config.InitialInterval != 1*time.Second {
		t.Errorf("expected InitialInterval=1s, got %v", config.InitialInterval)
	}

	if config.MaxInterval != 10*time.Second {
		t.Errorf("expected MaxInterval=10s, got %v", config.MaxInterval)
	}

	if config.Multiplier != 2.0 {
		t.Errorf("expected Multiplier=2.0, got %f", config.Multiplier)
	}
}

/**
 * Feature: designer-ai-assistant, Property 19: Token Usage Logging
 * Validates: Requirements 13.4
 *
 * 属性测试：验证Token使用日志记录
 * Property test: Verify token usage logging
 *
 * 对于任何AI服务调用，系统应记录token消耗和响应时间
 * For any AI service call, the system SHALL log the token consumption and response time.
 */

// TestTokenUsageLoggingProperty 测试Token使用日志记录属性
// Test token usage logging property
func TestTokenUsageLoggingProperty(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 50

	properties := gopter.NewProperties(parameters)

	// Property 1: 成功请求记录完整的token信息
	// Property 1: Successful requests log complete token information
	properties.Property("successful requests log complete token info", prop.ForAll(
		func(promptTokens, completionTokens int) bool {
			if promptTokens < 0 || promptTokens > 10000 || completionTokens < 0 || completionTokens > 10000 {
				return true // 跳过无效值 / Skip invalid values
			}

			totalTokens := promptTokens + completionTokens

			// 创建返回指定token数量的测试服务器
			// Create test server that returns specified token counts
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				response := fmt.Sprintf(`{
					"id": "test-id",
					"object": "chat.completion",
					"created": 1234567890,
					"model": "test-model",
					"choices": [{"index": 0, "message": {"role": "assistant", "content": "test"}, "finish_reason": "stop"}],
					"usage": {"prompt_tokens": %d, "completion_tokens": %d, "total_tokens": %d}
				}`, promptTokens, completionTokens, totalTokens)
				w.Write([]byte(response))
			}))
			defer server.Close()

			client := NewClient("test-api-key", server.URL, "test-model", 5, 100, 0.7)
			client.ClearTokenUsageLogs()

			_, err := client.Chat([]ChatMessage{
				{Role: "user", Content: []ContentPart{{Type: "text", Text: "test"}}},
			})

			if err != nil {
				return false
			}

			logs := client.GetTokenUsageLogs()
			if len(logs) != 1 {
				return false
			}

			log := logs[0]
			// 验证日志包含正确的token信息 / Verify log contains correct token info
			return log.PromptTokens == promptTokens &&
				log.CompletionTokens == completionTokens &&
				log.TotalTokens == totalTokens &&
				log.Success == true &&
				log.ResponseTime > 0 &&
				log.RequestID != "" &&
				log.Model == "test-model"
		},
		gen.IntRange(1, 1000),
		gen.IntRange(1, 1000),
	))

	// Property 2: 失败请求也记录日志
	// Property 2: Failed requests also log
	properties.Property("failed requests also log with error message", prop.ForAll(
		func(errorCode int) bool {
			if errorCode < 400 || errorCode >= 600 {
				return true
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(errorCode)
				w.Write([]byte(`{"error": {"message": "test error", "type": "test_error", "code": "test"}}`))
			}))
			defer server.Close()

			client := NewClientWithConfig(
				"test-api-key",
				server.URL,
				"test-model",
				5,
				100,
				0.7,
				RetryConfig{
					MaxRetries:      0, // 不重试 / No retry
					InitialInterval: 10 * time.Millisecond,
					MaxInterval:     50 * time.Millisecond,
					Multiplier:      2.0,
				},
			)
			client.ClearTokenUsageLogs()

			_, err := client.Chat([]ChatMessage{
				{Role: "user", Content: []ContentPart{{Type: "text", Text: "test"}}},
			})

			if err == nil {
				return false // 应该返回错误 / Should return error
			}

			logs := client.GetTokenUsageLogs()
			if len(logs) != 1 {
				return false
			}

			log := logs[0]
			// 验证失败日志包含错误信息 / Verify failed log contains error info
			return log.Success == false &&
				log.ErrorMessage != "" &&
				log.ResponseTime > 0 &&
				log.RequestID != ""
		},
		gen.IntRange(400, 599),
	))

	// Property 3: 响应时间总是正数
	// Property 3: Response time is always positive
	properties.Property("response time is always positive", prop.ForAll(
		func(delay int) bool {
			if delay < 0 || delay > 100 {
				return true
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(time.Duration(delay) * time.Millisecond)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{
					"id": "test-id",
					"object": "chat.completion",
					"created": 1234567890,
					"model": "test-model",
					"choices": [{"index": 0, "message": {"role": "assistant", "content": "test"}, "finish_reason": "stop"}],
					"usage": {"prompt_tokens": 10, "completion_tokens": 5, "total_tokens": 15}
				}`))
			}))
			defer server.Close()

			client := NewClient("test-api-key", server.URL, "test-model", 5, 100, 0.7)
			client.ClearTokenUsageLogs()

			_, _ = client.Chat([]ChatMessage{
				{Role: "user", Content: []ContentPart{{Type: "text", Text: "test"}}},
			})

			logs := client.GetTokenUsageLogs()
			if len(logs) != 1 {
				return false
			}

			// 响应时间应该大于等于延迟时间 / Response time should be >= delay
			return logs[0].ResponseTime >= time.Duration(delay)*time.Millisecond
		},
		gen.IntRange(0, 50),
	))

	// Property 4: 多次请求累计token正确
	// Property 4: Multiple requests accumulate tokens correctly
	properties.Property("multiple requests accumulate tokens correctly", prop.ForAll(
		func(requestCount int) bool {
			if requestCount < 1 || requestCount > 10 {
				return true
			}

			tokensPerRequest := 15 // 固定每次请求的token数 / Fixed tokens per request

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{
					"id": "test-id",
					"object": "chat.completion",
					"created": 1234567890,
					"model": "test-model",
					"choices": [{"index": 0, "message": {"role": "assistant", "content": "test"}, "finish_reason": "stop"}],
					"usage": {"prompt_tokens": 10, "completion_tokens": 5, "total_tokens": 15}
				}`))
			}))
			defer server.Close()

			client := NewClient("test-api-key", server.URL, "test-model", 5, 100, 0.7)
			client.ClearTokenUsageLogs()

			for i := 0; i < requestCount; i++ {
				_, _ = client.Chat([]ChatMessage{
					{Role: "user", Content: []ContentPart{{Type: "text", Text: "test"}}},
				})
			}

			// 验证总token数 / Verify total tokens
			promptTokens, completionTokens, totalTokens := client.GetTotalTokenUsage()
			expectedTotal := requestCount * tokensPerRequest

			return totalTokens == expectedTotal &&
				promptTokens == requestCount*10 &&
				completionTokens == requestCount*5
		},
		gen.IntRange(1, 10),
	))

	properties.TestingRun(t)
}

// TestTokenMonitorProperty 测试Token监控器属性
// Test token monitor property
func TestTokenMonitorProperty(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 50

	properties := gopter.NewProperties(parameters)

	// Property 1: 监控器记录所有日志
	// Property 1: Monitor records all logs
	properties.Property("monitor records all logs", prop.ForAll(
		func(logCount int) bool {
			if logCount < 1 || logCount > 100 {
				return true
			}

			monitor := NewTokenMonitor(TokenMonitorConfig{
				MaxLogs:       1000,
				FlushInterval: 0, // 禁用自动刷新 / Disable auto flush
				LogFilePath:   "",
			})

			for i := 0; i < logCount; i++ {
				monitor.RecordUsage(TokenUsageLog{
					Timestamp:        time.Now(),
					RequestID:        fmt.Sprintf("req_%d", i),
					Model:            "test-model",
					PromptTokens:     10,
					CompletionTokens: 5,
					TotalTokens:      15,
					ResponseTime:     100 * time.Millisecond,
					Success:          true,
				})
			}

			logs := monitor.GetLogs()
			return len(logs) == logCount
		},
		gen.IntRange(1, 100),
	))

	// Property 2: 监控器限制最大日志数量
	// Property 2: Monitor limits max log count
	properties.Property("monitor limits max log count", prop.ForAll(
		func(maxLogs, logCount int) bool {
			if maxLogs < 1 || maxLogs > 100 || logCount < 1 || logCount > 200 {
				return true
			}

			monitor := NewTokenMonitor(TokenMonitorConfig{
				MaxLogs:       maxLogs,
				FlushInterval: 0,
				LogFilePath:   "",
			})

			for i := 0; i < logCount; i++ {
				monitor.RecordUsage(TokenUsageLog{
					Timestamp:        time.Now(),
					RequestID:        fmt.Sprintf("req_%d", i),
					Model:            "test-model",
					PromptTokens:     10,
					CompletionTokens: 5,
					TotalTokens:      15,
					ResponseTime:     100 * time.Millisecond,
					Success:          true,
				})
			}

			logs := monitor.GetLogs()
			expectedCount := logCount
			if expectedCount > maxLogs {
				expectedCount = maxLogs
			}

			return len(logs) == expectedCount
		},
		gen.IntRange(10, 50),
		gen.IntRange(1, 100),
	))

	// Property 3: 统计信息正确计算
	// Property 3: Statistics are calculated correctly
	properties.Property("statistics are calculated correctly", prop.ForAll(
		func(successCount, failCount int) bool {
			if successCount < 0 || successCount > 50 || failCount < 0 || failCount > 50 {
				return true
			}

			monitor := NewTokenMonitor(TokenMonitorConfig{
				MaxLogs:       1000,
				FlushInterval: 0,
				LogFilePath:   "",
			})

			tokensPerSuccess := 15

			// 添加成功日志 / Add success logs
			for i := 0; i < successCount; i++ {
				monitor.RecordUsage(TokenUsageLog{
					Timestamp:        time.Now(),
					RequestID:        fmt.Sprintf("success_%d", i),
					Model:            "test-model",
					PromptTokens:     10,
					CompletionTokens: 5,
					TotalTokens:      15,
					ResponseTime:     100 * time.Millisecond,
					Success:          true,
				})
			}

			// 添加失败日志 / Add failure logs
			for i := 0; i < failCount; i++ {
				monitor.RecordUsage(TokenUsageLog{
					Timestamp:    time.Now(),
					RequestID:    fmt.Sprintf("fail_%d", i),
					Model:        "test-model",
					ResponseTime: 50 * time.Millisecond,
					Success:      false,
					ErrorMessage: "test error",
				})
			}

			stats := monitor.GetStatistics()

			return stats.TotalRequests == successCount+failCount &&
				stats.SuccessfulRequests == successCount &&
				stats.FailedRequests == failCount &&
				stats.TotalTokens == successCount*tokensPerSuccess
		},
		gen.IntRange(0, 50),
		gen.IntRange(0, 50),
	))

	properties.TestingRun(t)
}

// TestTokenUsageLogFields 测试Token使用日志字段
// Test token usage log fields
func TestTokenUsageLogFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": "test-id",
			"object": "chat.completion",
			"created": 1234567890,
			"model": "test-model",
			"choices": [{"index": 0, "message": {"role": "assistant", "content": "test"}, "finish_reason": "stop"}],
			"usage": {"prompt_tokens": 100, "completion_tokens": 50, "total_tokens": 150}
		}`))
	}))
	defer server.Close()

	client := NewClient("test-api-key", server.URL, "test-model", 5, 100, 0.7)
	client.ClearTokenUsageLogs()

	_, err := client.Chat([]ChatMessage{
		{Role: "user", Content: []ContentPart{{Type: "text", Text: "test"}}},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	logs := client.GetTokenUsageLogs()
	if len(logs) != 1 {
		t.Fatalf("expected 1 log, got %d", len(logs))
	}

	log := logs[0]

	// 验证所有必需字段 / Verify all required fields
	if log.RequestID == "" {
		t.Error("RequestID should not be empty")
	}
	if log.Model != "test-model" {
		t.Errorf("expected model 'test-model', got '%s'", log.Model)
	}
	if log.PromptTokens != 100 {
		t.Errorf("expected PromptTokens 100, got %d", log.PromptTokens)
	}
	if log.CompletionTokens != 50 {
		t.Errorf("expected CompletionTokens 50, got %d", log.CompletionTokens)
	}
	if log.TotalTokens != 150 {
		t.Errorf("expected TotalTokens 150, got %d", log.TotalTokens)
	}
	if log.ResponseTime <= 0 {
		t.Error("ResponseTime should be positive")
	}
	if !log.Success {
		t.Error("Success should be true")
	}
	if log.Timestamp.IsZero() {
		t.Error("Timestamp should not be zero")
	}
}

// TestChunkingFunctionality 测试大文档分段处理功能
// Test large document chunking functionality
func TestChunkingFunctionality(t *testing.T) {
	// 测试token估算 / Test token estimation
	t.Run("EstimateTokenCount", func(t *testing.T) {
		// 纯英文 / Pure English
		englishText := "Hello world this is a test"
		englishTokens := EstimateTokenCount(englishText)
		if englishTokens <= 0 {
			t.Error("English token count should be positive")
		}

		// 纯中文 / Pure Chinese
		chineseText := "你好世界这是一个测试"
		chineseTokens := EstimateTokenCount(chineseText)
		if chineseTokens <= 0 {
			t.Error("Chinese token count should be positive")
		}

		// 混合文本 / Mixed text
		mixedText := "Hello 你好 world 世界"
		mixedTokens := EstimateTokenCount(mixedText)
		if mixedTokens <= 0 {
			t.Error("Mixed token count should be positive")
		}
	})

	// 测试文本分段 / Test text splitting
	t.Run("SplitTextIntoChunks", func(t *testing.T) {
		config := ChunkConfig{
			MaxTokensPerChunk: 100,
			OverlapTokens:     10,
			Separator:         "\n\n",
		}

		// 短文本不分段 / Short text should not be split
		shortText := "This is a short text."
		chunks := SplitTextIntoChunks(shortText, config)
		if len(chunks) != 1 {
			t.Errorf("Short text should result in 1 chunk, got %d", len(chunks))
		}

		// 长文本应该分段 / Long text should be split
		var longText string
		for i := 0; i < 100; i++ {
			longText += "这是一段很长的测试文本，用于测试分段功能。\n\n"
		}
		chunks = SplitTextIntoChunks(longText, config)
		if len(chunks) <= 1 {
			t.Errorf("Long text should result in multiple chunks, got %d", len(chunks))
		}
	})

	// 测试分段处理API / Test chunking API
	t.Run("ChatWithChunking", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"id": "test-id",
				"object": "chat.completion",
				"created": 1234567890,
				"model": "test-model",
				"choices": [{"index": 0, "message": {"role": "assistant", "content": "chunk response"}, "finish_reason": "stop"}],
				"usage": {"prompt_tokens": 50, "completion_tokens": 20, "total_tokens": 70}
			}`))
		}))
		defer server.Close()

		client := NewClient("test-api-key", server.URL, "test-model", 5, 100, 0.7)
		client.SetEnableLogging(false)

		// 测试短文本（不分段）/ Test short text (no chunking)
		shortResult, err := client.ChatWithChunking("system prompt", "short user prompt", DefaultChunkConfig)
		if err != nil {
			t.Fatalf("ChatWithChunking failed for short text: %v", err)
		}
		if shortResult.TotalChunks != 1 {
			t.Errorf("Short text should have 1 chunk, got %d", shortResult.TotalChunks)
		}
		if shortResult.MergedResult == "" {
			t.Error("MergedResult should not be empty")
		}
	})

	// 测试自动分段 / Test auto chunking
	t.Run("ChatWithAutoChunking", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"id": "test-id",
				"object": "chat.completion",
				"created": 1234567890,
				"model": "test-model",
				"choices": [{"index": 0, "message": {"role": "assistant", "content": "auto chunk response"}, "finish_reason": "stop"}],
				"usage": {"prompt_tokens": 30, "completion_tokens": 10, "total_tokens": 40}
			}`))
		}))
		defer server.Close()

		client := NewClient("test-api-key", server.URL, "test-model", 5, 100, 0.7)
		client.SetEnableLogging(false)

		result, err := client.ChatWithAutoChunking("system prompt", "user prompt")
		if err != nil {
			t.Fatalf("ChatWithAutoChunking failed: %v", err)
		}
		if result.TotalTokens <= 0 {
			t.Error("TotalTokens should be positive")
		}
	})
}

// TestChunkingProperty 测试分段处理属性
// Test chunking property
func TestChunkingProperty(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 30

	properties := gopter.NewProperties(parameters)

	// Property: 分段后的总内容不丢失关键信息
	// Property: Chunked content preserves key information
	properties.Property("chunking preserves content", prop.ForAll(
		func(paragraphCount int) bool {
			if paragraphCount < 1 || paragraphCount > 20 {
				return true
			}

			// 生成测试文本 / Generate test text
			var text string
			for i := 0; i < paragraphCount; i++ {
				text += fmt.Sprintf("段落%d：这是测试内容。\n\n", i+1)
			}

			config := ChunkConfig{
				MaxTokensPerChunk: 50,
				OverlapTokens:     5,
				Separator:         "\n\n",
			}

			chunks := SplitTextIntoChunks(text, config)

			// 验证所有段落都被包含 / Verify all paragraphs are included
			for i := 0; i < paragraphCount; i++ {
				marker := fmt.Sprintf("段落%d", i+1)
				found := false
				for _, chunk := range chunks {
					if strings.Contains(chunk, marker) {
						found = true
						break
					}
				}
				if !found {
					return false
				}
			}

			return true
		},
		gen.IntRange(1, 20),
	))

	// Property: 每个分段不超过最大token限制
	// Property: Each chunk does not exceed max token limit
	properties.Property("chunks respect max token limit", prop.ForAll(
		func(maxTokens int) bool {
			if maxTokens < 50 || maxTokens > 500 {
				return true
			}

			// 生成长文本 / Generate long text
			var text string
			for i := 0; i < 50; i++ {
				text += "这是一段测试文本用于验证分段功能是否正确工作。\n\n"
			}

			config := ChunkConfig{
				MaxTokensPerChunk: maxTokens,
				OverlapTokens:     10,
				Separator:         "\n\n",
			}

			chunks := SplitTextIntoChunks(text, config)

			// 验证每个分段的token数不超过限制（允许一定误差）
			// Verify each chunk's token count doesn't exceed limit (with tolerance)
			for _, chunk := range chunks {
				tokens := EstimateTokenCount(chunk)
				// 允许20%的误差，因为token估算不精确
				// Allow 20% tolerance due to imprecise token estimation
				if tokens > int(float64(maxTokens)*1.2) {
					return false
				}
			}

			return true
		},
		gen.IntRange(50, 500),
	))

	properties.TestingRun(t)
}
