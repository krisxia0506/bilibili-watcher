package model

import (
	"time"
	// "gorm.io/gorm" // No longer embedding gorm.Model
)

// VideoProgress 记录视频观看进度 (遵循阿里巴巴规范)
// 表示特定时间点的视频观看进度记录。
// Note: gorm.Model is not used to avoid DeletedAt field, matching the schema.
type VideoProgress struct {
	ID           uint      `gorm:"primarykey;comment:主键 ID"`
	AID          int64     `gorm:"column:aid;index;not null;default:0;comment:视频稿件 ID (AV 号)"`          // 显式列名
	BVID         string    `gorm:"column:bvid;not null;default:'';comment:视频 BV 号"`                     // 显式列名
	LastPlayCID  int64     `gorm:"column:last_play_cid;index;not null;default:0;comment:上次播放的视频分 P ID"` // 显式列名 & 重命名
	LastPlayTime int64     `gorm:"column:last_play_time;not null;default:0;comment:上次播放时间/进度 (毫秒)"`     // 重命名
	RecordedAt   time.Time `gorm:"column:recorded_at;index;not null;default:CURRENT_TIMESTAMP(3);comment:记录时间"`
	GmtCreate    time.Time `gorm:"column:gmt_create;type:datetime(3);not null;autoCreateTime;comment:创建时间"`
	GmtModified  time.Time `gorm:"column:gmt_modified;type:datetime(3);not null;autoUpdateTime;comment:更新时间"`
	// DeletedAt gorm.DeletedAt `gorm:"index"` // Removed
}

// TableName 指定 VideoProgress 的表名为 "video_progress"
// 显式设置以匹配单数表名约定和 SQL schema。
func (VideoProgress) TableName() string {
	return "video_progress"
}
