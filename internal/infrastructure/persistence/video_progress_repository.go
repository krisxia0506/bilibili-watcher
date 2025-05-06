package persistence

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/krisxia0506/bilibili-watcher/internal/domain/model"
	"github.com/krisxia0506/bilibili-watcher/internal/domain/repository"
)

// gormVideoProgressRepository 是 VideoProgressRepository 的 GORM 实现。
type gormVideoProgressRepository struct {
	db *gorm.DB
}

// NewGormVideoProgressRepository 创建一个新的 GORM VideoProgressRepository 实例。
func NewGormVideoProgressRepository(db *gorm.DB) repository.VideoProgressRepository {
	return &gormVideoProgressRepository{db: db}
}

// Save 保存一条视频观看进度记录。
func (r *gormVideoProgressRepository) Save(ctx context.Context, progress *model.VideoProgress) error {
	return r.db.WithContext(ctx).Create(progress).Error
}

// GetLatestByAIDAndCID 获取指定视频 (稿件+分P) 的最新一条进度记录。
func (r *gormVideoProgressRepository) GetLatestByAIDAndCID(ctx context.Context, aid, lastPlayCID int64) (*model.VideoProgress, error) {
	var progress model.VideoProgress
	// 查找给定 aid 和 last_play_cid 的最新 RecordedAt 时间戳记录
	err := r.db.WithContext(ctx).
		Where("aid = ? AND last_play_cid = ?", aid, lastPlayCID). // 更新查询中的列名
		Order("recorded_at DESC").
		First(&progress).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 如果未找到记录，则返回 nil, nil，而不是错误
		}
		return nil, err // 返回其他潜在错误
	}
	return &progress, nil
}

// ListByDateRange 获取指定日期范围内的所有进度记录。
func (r *gormVideoProgressRepository) ListByDateRange(ctx context.Context, start, end time.Time) ([]*model.VideoProgress, error) {
	var progresses []*model.VideoProgress
	err := r.db.WithContext(ctx).
		Where("recorded_at >= ? AND recorded_at < ?", start, end).
		Order("recorded_at ASC").
		Find(&progresses).Error

	if err != nil {
		return nil, err
	}
	return progresses, nil
}
