package repository

import (
	"context"
	"time"

	"github.com/krisxia0506/bilibili-watcher/internal/domain/model"
)

// VideoProgressRepository 定义视频进度数据操作的接口。
type VideoProgressRepository interface {
	// Save 保存一条视频观看进度记录。
	Save(ctx context.Context, progress *model.VideoProgress) error

	// GetLatestByAIDAndCID 获取指定视频 (稿件+分P) 的最新一条进度记录。
	GetLatestByAIDAndCID(ctx context.Context, aid, lastPlayCID int64) (*model.VideoProgress, error)

	// ListByDateRange 获取指定日期范围内的所有进度记录。
	ListByDateRange(ctx context.Context, start, end time.Time) ([]*model.VideoProgress, error)
}
