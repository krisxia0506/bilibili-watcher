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
// 它包含了 GORM 特定的标签和可能的基础字段（如 gorm.Model）。
type videoProgressGorm struct {
	ID            uint      `gorm:"primaryKey"`
	AID           int64     `gorm:"column:aid;uniqueIndex;not null"`
	BVID          string    `gorm:"column:bvid;index;not null"`
	LastPlayCID   int64     `gorm:"column:last_play_cid;not null"` // 注意大小写
	LastPlayTime  int64     `gorm:"column:last_play_time;not null"`
	TotalDuration int64     `gorm:"column:total_duration;not null"` // 新增
	FetchTime     time.Time `gorm:"column:fetch_time;not null"`     // 新增
	GmtCreate     time.Time `gorm:"column:gmt_create;autoCreateTime"`
	GmtModified   time.Time `gorm:"column:gmt_modified;autoUpdateTime"`
	// gorm.Model // 可以包含 gorm.Model 来获取 ID, CreatedAt, UpdatedAt, DeletedAt
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
		ID:            g.ID,
		AID:           g.AID,
		BVID:          g.BVID,
		LastPlayCID:   g.LastPlayCID,
		LastPlayTime:  g.LastPlayTime,
		GmtCreate:     g.GmtCreate,
		GmtModified:   g.GmtModified,
	}
}

// fromDomain 将领域模型转换为 GORM 模型。
func fromDomain(d *model.VideoProgress) *videoProgressGorm {
	if d == nil {
		return nil
	}
	return &videoProgressGorm{
		ID:            d.ID,
		AID:           d.AID,
		BVID:          d.BVID,
		LastPlayCID:   d.LastPlayCID,
		LastPlayTime:  d.LastPlayTime,
		GmtCreate:     d.GmtCreate,     // GORM 会在创建时处理
		GmtModified:   d.GmtModified,   // GORM 会在创建/更新时处理
	}
}

// FindByAID 根据 AID 查找视频进度记录。
func (r *gormVideoProgressRepository) FindByAID(ctx context.Context, aid int64) (*model.VideoProgress, error) {
	var gormModel videoProgressGorm
	// 使用 First 而不是 Find，因为我们期望最多一条记录，并且希望在未找到时 GORM 返回 ErrRecordNotFound
	result := r.db.WithContext(ctx).Where("aid = ?", aid).First(&gormModel)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// 将 GORM 的错误转换为领域定义的错误
			return nil, repository.ErrVideoProgressNotFound
		}
		// 其他数据库错误
		log.Printf("Database error finding video progress by AID %d: %v", aid, result.Error)
		return nil, fmt.Errorf("database error finding progress by AID: %w", result.Error)
	}

	// 转换并返回领域模型
	return gormModel.toDomain(), nil
}
