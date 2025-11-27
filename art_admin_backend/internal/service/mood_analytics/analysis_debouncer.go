package mood_analytics

import (
	"sync"
	"time"
)

// AnalysisDebouncer 分析防抖器
type AnalysisDebouncer struct {
	lastAnalysisTime sync.Map // map[int64]time.Time
	minInterval      time.Duration
	cleanupTicker    *time.Ticker
	stopCleanup      chan bool
}

// NewAnalysisDebouncer 创建分析防抖器
func NewAnalysisDebouncer(minInterval time.Duration) *AnalysisDebouncer {
	debouncer := &AnalysisDebouncer{
		minInterval:   minInterval,
		cleanupTicker: time.NewTicker(time.Hour), // 每小时清理一次
		stopCleanup:   make(chan bool),
	}

	// 启动清理协程
	go debouncer.startCleanup()

	return debouncer
}

// ShouldCheck 是否应该进行分析
func (d *AnalysisDebouncer) ShouldCheck(userID int64) bool {
	// 获取用户最后分析时间
	if lastTime, exists := d.lastAnalysisTime.Load(userID); exists {
		if lastAnalysisTime, ok := lastTime.(time.Time); ok {
			// 如果距离上次分析时间小于最小间隔，则跳过
			if time.Since(lastAnalysisTime) < d.minInterval {
				return false
			}
		}
	}

	// 更新最后分析时间
	d.lastAnalysisTime.Store(userID, time.Now())
	return true
}

// startCleanup 启动清理协程
func (d *AnalysisDebouncer) startCleanup() {
	for {
		select {
		case <-d.cleanupTicker.C:
			d.cleanup()
		case <-d.stopCleanup:
			return
		}
	}
}

// cleanup 清理过期的记录
func (d *AnalysisDebouncer) cleanup() {
	now := time.Now()
	d.lastAnalysisTime.Range(func(key, value interface{}) bool {
		if lastTime, ok := value.(time.Time); ok {
			// 清理超过24小时的记录
			if now.Sub(lastTime) > 24*time.Hour {
				d.lastAnalysisTime.Delete(key)
			}
		}
		return true
	})
}

// Stop 停止防抖器
func (d *AnalysisDebouncer) Stop() {
	d.cleanupTicker.Stop()
	close(d.stopCleanup)
}

// 全局防抖器实例
var globalDebouncer *AnalysisDebouncer

// GetAnalysisDebouncer 获取全局防抖器实例
func GetAnalysisDebouncer() *AnalysisDebouncer {
	if globalDebouncer == nil {
		globalDebouncer = NewAnalysisDebouncer(5 * time.Minute) // 默认5分钟间隔
	}
	return globalDebouncer
}
