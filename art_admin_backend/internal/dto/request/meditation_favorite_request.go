package request

// UpdatePlayRecordRequest 更新播放记录请求
type UpdatePlayRecordRequest struct {
	ContentID int64   `json:"content_id" binding:"required"`
	Duration  int     `json:"duration" binding:"min=0"`         // 实际播放时长（秒）
	Progress  float64 `json:"progress" binding:"min=0,max=100"` // 播放进度（百分比）
	Completed bool    `json:"completed"`                        // 是否完成
}
