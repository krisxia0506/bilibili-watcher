package model

import (
	"time"

	"gorm.io/gorm"
)

// VideoProgress 记录视频观看进度
// Represents a record of video watching progress at a specific point in time.
type VideoProgress struct {
	gorm.Model           // Includes ID, CreatedAt, UpdatedAt, DeletedAt
	AID        int64     `gorm:"index;comment:视频稿件 ID (AV 号)"` // Video Archive ID (av)
	CID        int64     `gorm:"index;comment:视频分 P ID"`       // Video Content ID (page)
	BVID       string    `gorm:"comment:视频 BV 号"`              // Video BV ID
	Progress   int       `gorm:"comment:观看进度 (毫秒)"`            // Watching progress in milliseconds
	RecordedAt time.Time `gorm:"index;comment:记录时间"`           // Timestamp when the record was created
}
