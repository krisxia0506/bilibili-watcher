package persistence

import (
	"context"
	"errors"
	"fmt"
	"log"
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

// videoProgressGorm 是 VideoProgress 领域模型对应的 GORM 数据模型。
type videoProgressGorm struct {
	ID           uint      `gorm:"primaryKey"`
	AID          int64     `gorm:"column:aid;uniqueIndex;not null"`
	BVID         string    `gorm:"column:bvid;index;not null"`
	LastPlayCID  int64     `gorm:"column:last_play_cid;not null"`
	LastPlayTime int64     `gorm:"column:last_play_time;not null"`
	RecordedAt   time.Time `gorm:"column:recorded_at;index;not null;default:CURRENT_TIMESTAMP(3)"`
	GmtCreate    time.Time `gorm:"column:gmt_create;autoCreateTime"`
	GmtModified  time.Time `gorm:"column:gmt_modified;autoUpdateTime"`
}

// TableName 指定 GORM 应使用的表名。
func (videoProgressGorm) TableName() string {
	return "video_progress"
}

// toDomain 将 GORM 模型转换为领域模型。
func (g *videoProgressGorm) toDomain() *model.VideoProgress {
	if g == nil {
		return nil
	}
	return &model.VideoProgress{
		ID:           g.ID,
		AID:          g.AID,
		BVID:         g.BVID,
		LastPlayCID:  g.LastPlayCID,
		LastPlayTime: g.LastPlayTime,
		RecordedAt:   g.RecordedAt,
		GmtCreate:    g.GmtCreate,
		GmtModified:  g.GmtModified,
	}
}

// fromDomain 将领域模型转换为 GORM 模型。
func fromDomain(d *model.VideoProgress) *videoProgressGorm {
	if d == nil {
		return nil
	}
	return &videoProgressGorm{
		ID:           d.ID,
		AID:          d.AID,
		BVID:         d.BVID,
		LastPlayCID:  d.LastPlayCID,
		LastPlayTime: d.LastPlayTime,
		RecordedAt:   d.RecordedAt,
		GmtCreate:    d.GmtCreate,
		GmtModified:  d.GmtModified,
	}
}

// FindByAID 根据 AID 查找视频进度记录。
// 注意：此方法仍然存在，但其使用场景可能因 FetchAndSaveVideoProgress 的更改而改变。
func (r *gormVideoProgressRepository) FindByAID(ctx context.Context, aid int64) (*model.VideoProgress, error) {
	var gormModel videoProgressGorm
	result := r.db.WithContext(ctx).Where("aid = ?", aid).First(&gormModel)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, repository.ErrVideoProgressNotFound
		}
		log.Printf("Database error finding video progress by AID %d: %v", aid, result.Error)
		return nil, fmt.Errorf("database error finding progress by AID: %w", result.Error)
	}
	return gormModel.toDomain(), nil
}

// ListByAIDAndTimestampRange 获取指定 AID 在给定时间范围内的所有进度记录，按记录时间升序排序。
func (r *gormVideoProgressRepository) ListByAIDAndTimestampRange(ctx context.Context, aid int64, startTime, endTime time.Time) ([]*model.VideoProgress, error) {
	var progressesGorm []videoProgressGorm
	err := r.db.WithContext(ctx).
		Where("aid = ? AND recorded_at >= ? AND recorded_at <= ?", aid, startTime, endTime).
		Order("recorded_at ASC").
		Find(&progressesGorm).Error

	if err != nil {
		log.Printf("Database error finding video progress by AID %d and time range [%s, %s]: %v", aid, startTime, endTime, err)
		return nil, fmt.Errorf("database error finding progress by AID and time range: %w", err)
	}

	domainProgresses := make([]*model.VideoProgress, 0, len(progressesGorm))
	for _, g := range progressesGorm {
		domainProgress := g.toDomain()
		if domainProgress != nil {
			domainProgresses = append(domainProgresses, domainProgress)
		}
	}

	return domainProgresses, nil
}

// ListByBVIDAndTimestampRange 获取指定 BVID 在给定时间范围内的所有进度记录，按记录时间升序排序。
func (r *gormVideoProgressRepository) ListByBVIDAndTimestampRange(ctx context.Context, bvid string, startTime, endTime time.Time) ([]*model.VideoProgress, error) {
	var progressesGorm []videoProgressGorm
	// 查找 RecordedAt 在 [startTime, endTime] 范围内的记录
	err := r.db.WithContext(ctx).
		Where("bvid = ? AND recorded_at >= ? AND recorded_at <= ?", bvid, startTime, endTime).
		Order("recorded_at ASC"). // 按记录时间升序排序
		Find(&progressesGorm).Error

	if err != nil {
		// GORM Find 在未找到时不会返回 ErrRecordNotFound，而是返回空切片和 nil error (在某些版本/场景下), 但仍需检查错误
		log.Printf("Database error finding video progress by BVID %s and time range [%s, %s]: %v", bvid, startTime, endTime, err)
		return nil, fmt.Errorf("database error finding progress by BVID and time range: %w", err)
	}

	// 映射到领域模型
	domainProgresses := make([]*model.VideoProgress, 0, len(progressesGorm))
	for _, g := range progressesGorm {
		domainProgress := g.toDomain() // Ensure toDomain is correct
		if domainProgress != nil {
			domainProgresses = append(domainProgresses, domainProgress)
		}
	}

	return domainProgresses, nil
}
