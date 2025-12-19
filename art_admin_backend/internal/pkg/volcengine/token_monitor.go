package volcengine

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"
)

// TokenMonitor Token使用监控器
// Token usage monitor
type TokenMonitor struct {
	logs          []TokenUsageLog
	mutex         sync.RWMutex
	maxLogs       int           // 最大日志数量 / Max log count
	flushInterval time.Duration // 刷新间隔 / Flush interval
	logFilePath   string        // 日志文件路径 / Log file path
	stopChan      chan struct{}
	callbacks     []TokenUsageCallback
}

// TokenMonitorConfig Token监控配置
// Token monitor configuration
type TokenMonitorConfig struct {
	MaxLogs       int           // 最大日志数量 / Max log count
	FlushInterval time.Duration // 刷新间隔 / Flush interval
	LogFilePath   string        // 日志文件路径 / Log file path
}

// DefaultTokenMonitorConfig 默认Token监控配置
// Default token monitor configuration
var DefaultTokenMonitorConfig = TokenMonitorConfig{
	MaxLogs:       10000,
	FlushInterval: 5 * time.Minute,
	LogFilePath:   "logs/token_usage.log",
}

// NewTokenMonitor 创建Token监控器
// Create token monitor
func NewTokenMonitor(config TokenMonitorConfig) *TokenMonitor {
	monitor := &TokenMonitor{
		logs:          make([]TokenUsageLog, 0),
		maxLogs:       config.MaxLogs,
		flushInterval: config.FlushInterval,
		logFilePath:   config.LogFilePath,
		stopChan:      make(chan struct{}),
		callbacks:     make([]TokenUsageCallback, 0),
	}

	// 启动定期刷新 / Start periodic flush
	if config.FlushInterval > 0 && config.LogFilePath != "" {
		go monitor.startPeriodicFlush()
	}

	return monitor
}

// NewDefaultTokenMonitor 创建默认配置的Token监控器
// Create token monitor with default configuration
func NewDefaultTokenMonitor() *TokenMonitor {
	return NewTokenMonitor(DefaultTokenMonitorConfig)
}

// AddCallback 添加回调函数
// Add callback function
func (m *TokenMonitor) AddCallback(callback TokenUsageCallback) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.callbacks = append(m.callbacks, callback)
}

// RecordUsage 记录Token使用
// Record token usage
func (m *TokenMonitor) RecordUsage(log TokenUsageLog) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 添加日志 / Add log
	m.logs = append(m.logs, log)

	// 如果超过最大数量，移除最旧的日志 / Remove oldest logs if exceeds max
	if len(m.logs) > m.maxLogs {
		m.logs = m.logs[len(m.logs)-m.maxLogs:]
	}

	// 调用回调函数 / Call callbacks
	for _, callback := range m.callbacks {
		go callback(log)
	}
}

// GetLogs 获取所有日志
// Get all logs
func (m *TokenMonitor) GetLogs() []TokenUsageLog {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	logs := make([]TokenUsageLog, len(m.logs))
	copy(logs, m.logs)
	return logs
}

// GetLogsByTimeRange 按时间范围获取日志
// Get logs by time range
func (m *TokenMonitor) GetLogsByTimeRange(start, end time.Time) []TokenUsageLog {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var result []TokenUsageLog
	for _, log := range m.logs {
		if log.Timestamp.After(start) && log.Timestamp.Before(end) {
			result = append(result, log)
		}
	}
	return result
}

// GetLogsByModel 按模型获取日志
// Get logs by model
func (m *TokenMonitor) GetLogsByModel(model string) []TokenUsageLog {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var result []TokenUsageLog
	for _, log := range m.logs {
		if log.Model == model {
			result = append(result, log)
		}
	}
	return result
}

// GetStatistics 获取统计信息
// Get statistics
func (m *TokenMonitor) GetStatistics() TokenStatistics {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	stats := TokenStatistics{
		TotalRequests: len(m.logs),
		ModelStats:    make(map[string]ModelTokenStats),
	}

	for _, log := range m.logs {
		if log.Success {
			stats.SuccessfulRequests++
			stats.TotalPromptTokens += log.PromptTokens
			stats.TotalCompletionTokens += log.CompletionTokens
			stats.TotalTokens += log.TotalTokens
			stats.TotalResponseTime += log.ResponseTime
		} else {
			stats.FailedRequests++
		}
		stats.TotalRetries += log.RetryCount

		// 按模型统计 / Statistics by model
		modelStats, exists := stats.ModelStats[log.Model]
		if !exists {
			modelStats = ModelTokenStats{Model: log.Model}
		}
		modelStats.RequestCount++
		if log.Success {
			modelStats.PromptTokens += log.PromptTokens
			modelStats.CompletionTokens += log.CompletionTokens
			modelStats.TotalTokens += log.TotalTokens
			modelStats.TotalResponseTime += log.ResponseTime
		}
		stats.ModelStats[log.Model] = modelStats
	}

	// 计算平均响应时间 / Calculate average response time
	if stats.SuccessfulRequests > 0 {
		stats.AverageResponseTime = stats.TotalResponseTime / time.Duration(stats.SuccessfulRequests)
	}

	return stats
}

// TokenStatistics Token统计信息
// Token statistics
type TokenStatistics struct {
	TotalRequests         int                        `json:"total_requests"`
	SuccessfulRequests    int                        `json:"successful_requests"`
	FailedRequests        int                        `json:"failed_requests"`
	TotalPromptTokens     int                        `json:"total_prompt_tokens"`
	TotalCompletionTokens int                        `json:"total_completion_tokens"`
	TotalTokens           int                        `json:"total_tokens"`
	TotalRetries          int                        `json:"total_retries"`
	TotalResponseTime     time.Duration              `json:"total_response_time"`
	AverageResponseTime   time.Duration              `json:"average_response_time"`
	ModelStats            map[string]ModelTokenStats `json:"model_stats"`
}

// ModelTokenStats 模型Token统计
// Model token statistics
type ModelTokenStats struct {
	Model             string        `json:"model"`
	RequestCount      int           `json:"request_count"`
	PromptTokens      int           `json:"prompt_tokens"`
	CompletionTokens  int           `json:"completion_tokens"`
	TotalTokens       int           `json:"total_tokens"`
	TotalResponseTime time.Duration `json:"total_response_time"`
}

// ClearLogs 清除日志
// Clear logs
func (m *TokenMonitor) ClearLogs() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.logs = make([]TokenUsageLog, 0)
}

// FlushToFile 将日志刷新到文件
// Flush logs to file
func (m *TokenMonitor) FlushToFile() error {
	m.mutex.RLock()
	logs := make([]TokenUsageLog, len(m.logs))
	copy(logs, m.logs)
	m.mutex.RUnlock()

	if len(logs) == 0 {
		return nil
	}

	// 打开文件（追加模式）/ Open file (append mode)
	file, err := os.OpenFile(m.logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// 写入日志 / Write logs
	encoder := json.NewEncoder(file)
	for _, log := range logs {
		if err := encoder.Encode(log); err != nil {
			return err
		}
	}

	return nil
}

// startPeriodicFlush 启动定期刷新
// Start periodic flush
func (m *TokenMonitor) startPeriodicFlush() {
	ticker := time.NewTicker(m.flushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := m.FlushToFile(); err != nil {
				log.Printf("[TokenMonitor] 刷新日志到文件失败: %v", err)
			}
		case <-m.stopChan:
			return
		}
	}
}

// Stop 停止监控器
// Stop monitor
func (m *TokenMonitor) Stop() {
	close(m.stopChan)
	// 最后一次刷新 / Final flush
	if err := m.FlushToFile(); err != nil {
		log.Printf("[TokenMonitor] 最终刷新日志失败: %v", err)
	}
}

// CreateClientCallback 创建客户端回调函数
// Create client callback function
func (m *TokenMonitor) CreateClientCallback() TokenUsageCallback {
	return func(log TokenUsageLog) {
		m.RecordUsage(log)
	}
}
