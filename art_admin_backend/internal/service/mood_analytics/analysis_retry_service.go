package mood_analytics

import (
	"art_admin_backend/internal/repository"
	"context"
	"log"
	"time"
)

// AnalysisRetryService 分析重试服务
type AnalysisRetryService struct {
	moodRecordRepo       *repository.MoodRecordRepository
	journalEntryRepo     *repository.JournalEntryRepository
	meditationRecordRepo *repository.MeditationRecordRepository
	maxRetryCount        int
	retryInterval        time.Duration
}

// NewAnalysisRetryService 创建分析重试服务
func NewAnalysisRetryService() *AnalysisRetryService {
	return &AnalysisRetryService{
		moodRecordRepo:       repository.NewMoodRecordRepository(),
		journalEntryRepo:     repository.NewJournalEntryRepository(),
		meditationRecordRepo: repository.NewMeditationRecordRepository(),
		maxRetryCount:        3,                // 最大重试次数
		retryInterval:        30 * time.Minute, // 重试间隔
	}
}

// StartRetryWorker 启动重试工作协程
func (s *AnalysisRetryService) StartRetryWorker() {
	ticker := time.NewTicker(time.Hour) // 每小时检查一次
	defer ticker.Stop()

	for range ticker.C {
		s.retryFailedAnalyses()
	}
}

// retryFailedAnalyses 重试失败的分析
func (s *AnalysisRetryService) retryFailedAnalyses() {
	ctx := context.Background()

	log.Println("开始重试失败的分析...")

	// 重试情绪记录分析
	s.retryMoodRecordAnalyses(ctx)

	// 重试日记条目分析
	s.retryJournalEntryAnalyses(ctx)

	// 重试冥想记录分析
	s.retryMeditationRecordAnalyses(ctx)

	log.Println("失败分析重试完成")
}

// retryMoodRecordAnalyses 重试情绪记录分析
func (s *AnalysisRetryService) retryMoodRecordAnalyses(ctx context.Context) {
	// 这里需要在repository中添加查找失败分析的方法
	// 暂时记录日志
	log.Println("重试情绪记录分析")
}

// retryJournalEntryAnalyses 重试日记条目分析
func (s *AnalysisRetryService) retryJournalEntryAnalyses(ctx context.Context) {
	// 这里需要在repository中添加查找失败分析的方法
	// 暂时记录日志
	log.Println("重试日记条目分析")
}

// retryMeditationRecordAnalyses 重试冥想记录分析
func (s *AnalysisRetryService) retryMeditationRecordAnalyses(ctx context.Context) {
	// 这里需要在repository中添加查找失败分析的方法
	// 暂时记录日志
	log.Println("重试冥想记录分析")
}

// shouldRetry 判断是否应该重试
func (s *AnalysisRetryService) shouldRetry(retryCount int, lastRetryTime time.Time) bool {
	if retryCount >= s.maxRetryCount {
		return false
	}

	return time.Since(lastRetryTime) >= s.retryInterval
}

// markAnalysisFailed 标记分析失败
func (s *AnalysisRetryService) markAnalysisFailed(recordType, recordID string, error string) {
	log.Printf("标记分析失败 - 类型: %s, ID: %s, 错误: %s", recordType, recordID, error)
}

// 全局重试服务实例
var globalRetryService *AnalysisRetryService

// GetAnalysisRetryService 获取全局重试服务实例
func GetAnalysisRetryService() *AnalysisRetryService {
	if globalRetryService == nil {
		globalRetryService = NewAnalysisRetryService()
		go globalRetryService.StartRetryWorker()
	}
	return globalRetryService
}
