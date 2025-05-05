package repository

import (
	"context"
	"time"

	"github.com/krisxia0506/bilibili-watcher/internal/domain/model"
)

// VideoProgressRepository defines the interface for video progress data operations.
// 视频进度数据操作接口定义
type VideoProgressRepository interface {
	// Save 保存一条视频观看进度记录
	// Saves a video progress record.
	Save(ctx context.Context, progress *model.VideoProgress) error

	// GetLatestByAIDAndCID 获取指定视频 (稿件+分P) 的最新一条进度记录
	// Gets the latest progress record for a specific video (AID + CID).
	GetLatestByAIDAndCID(ctx context.Context, aid, cid int64) (*model.VideoProgress, error)

	// ListByDateRange 获取指定日期范围内的所有进度记录
	// Lists all progress records within a specified date range.
	ListByDateRange(ctx context.Context, start, end time.Time) ([]*model.VideoProgress, error)
}
