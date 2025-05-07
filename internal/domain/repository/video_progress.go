package repository

import (
	"context"
	"errors"
	"time"

	"github.com/krisxia0506/bilibili-watcher/internal/domain/model"
)

// ErrVideoProgressNotFound 表示未找到指定的视频进度记录。
var ErrVideoProgressNotFound = errors.New("video progress not found")

// VideoProgressRepository 定义视频进度数据操作的接口。
type VideoProgressRepository interface {
	// Save 保存一条视频观看进度记录。
	Save(ctx context.Context, progress *model.VideoProgress) error

	// GetLatestByAIDAndCID 获取指定视频 (稿件+分P) 的最新一条进度记录。
	GetLatestByAIDAndCID(ctx context.Context, aid, lastPlayCID int64) (*model.VideoProgress, error)

	// ListByDateRange 获取指定日期范围内的所有进度记录。
	// Deprecated: Use ListByAIDAndTimestampRange or ListByBVIDAndTimestampRange for more specific queries.
	ListByDateRange(ctx context.Context, start, end time.Time) ([]*model.VideoProgress, error)

	// FindByAID 根据 AID 查找视频进度记录。
	// 如果未找到，应返回 ErrVideoProgressNotFound 错误。
	FindByAID(ctx context.Context, aid int64) (*model.VideoProgress, error)

	// ListByAIDAndTimestampRange 获取指定 AID 在给定时间范围内的所有进度记录，按记录时间升序排序。
	ListByAIDAndTimestampRange(ctx context.Context, aid int64, startTime, endTime time.Time) ([]*model.VideoProgress, error)

	// ListByBVIDAndTimestampRange 获取指定 BVID 在给定时间范围内的所有进度记录，按记录时间升序排序。
	ListByBVIDAndTimestampRange(ctx context.Context, bvid string, startTime, endTime time.Time) ([]*model.VideoProgress, error)
}
