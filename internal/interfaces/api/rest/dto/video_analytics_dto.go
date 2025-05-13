package dto

import "time"

// GetWatchedSegmentsRequest 获取观看分段请求体。
type GetWatchedSegmentsRequest struct {
	AID       string `json:"aid" binding:"omitempty"`                                          // 可选，AV 号
	BVID      string `json:"bvid" binding:"omitempty"`                                         // 可选，BV 号 (aid 和 bvid 必须提供一个)
	StartTime string `json:"start_time" binding:"required,datetime=2006-01-02T15:04:05Z07:00"` // RFC3339 格式
	EndTime   string `json:"end_time" binding:"required,datetime=2006-01-02T15:04:05Z07:00"`   // RFC3339 格式
	Interval  string `json:"interval" binding:"required,oneof=10m 30m 1h 1d"`                  // 时间间隔 (10分钟, 30分钟, 1小时, 1天)
}

// WatchedSegment 观看分段信息。
type WatchedSegment struct {
	SegmentStartTime   time.Time `json:"segment_start_time"`       // 分段开始时间
	SegmentEndTime     time.Time `json:"segment_end_time"`         // 分段结束时间
	WatchedDurationSec int64     `json:"watched_duration_seconds"` // 该分段内观看的时长（秒）
}

// GetWatchedSegmentsResponse 获取观看分段响应体 (Data 部分)。
type GetWatchedSegmentsResponse struct {
	Segments                []WatchedSegment `json:"segments"`
	TotalWatchedDurationSec int64            `json:"total_watched_duration_seconds"` // 新增：总观看时长（秒）
}
