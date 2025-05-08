package application

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/krisxia0506/bilibili-watcher/internal/domain/model"
	"github.com/krisxia0506/bilibili-watcher/internal/domain/repository"
	"github.com/krisxia0506/bilibili-watcher/internal/domain/service"
)

// VideoAnalyticsService 定义了视频分析相关的应用服务接口。
type VideoAnalyticsService interface {
	// GetWatchedSegments 计算并返回指定时间范围和间隔内的视频观看分段时长。
	GetWatchedSegments(ctx context.Context,
		aidStr, bvidStr string, // aid 和 bvid 提供一个
		overallStartTime, overallEndTime time.Time,
		interval time.Duration,
	) ([]WatchedSegmentResult, error)
}

// WatchedSegmentResult 包含单个时间分段的计算结果。
type WatchedSegmentResult struct {
	SegmentStartTime time.Time
	SegmentEndTime   time.Time
	WatchedDuration  time.Duration
}

// videoAnalyticsService 实现了 VideoAnalyticsService。
type videoAnalyticsService struct {
	biliClient   BilibiliClient                     // 获取视频页面信息
	progressRepo repository.VideoProgressRepository // 获取进度记录
	calculator   service.WatchTimeCalculator        // 计算时长
}

// NewVideoAnalyticsService 创建 VideoAnalyticsService 实例。
func NewVideoAnalyticsService(
	biliClient BilibiliClient,
	progressRepo repository.VideoProgressRepository,
	calculator service.WatchTimeCalculator,
) VideoAnalyticsService {
	return &videoAnalyticsService{
		biliClient:   biliClient,
		progressRepo: progressRepo,
		calculator:   calculator,
	}
}

// GetWatchedSegments 实现获取观看分段的逻辑。
func (s *videoAnalyticsService) GetWatchedSegments(ctx context.Context,
	aidStr, bvidStr string,
	overallStartTime, overallEndTime time.Time,
	interval time.Duration,
) ([]WatchedSegmentResult, error) {

	if aidStr == "" && bvidStr == "" {
		return nil, fmt.Errorf("either aid or bvid must be provided")
	}
	if interval <= 0 {
		return nil, fmt.Errorf("interval must be positive")
	}
	if overallEndTime.Before(overallStartTime) {
		return nil, fmt.Errorf("end time must be after start time")
	}

	// 1. 获取视频页面信息 (只需要执行一次)
	videoView, err := s.biliClient.GetVideoView(ctx, aidStr, bvidStr)
	if err != nil {
		return nil, fmt.Errorf("failed to get video view info: %w", err)
	}
	if videoView == nil || len(videoView.Pages) == 0 {
		return nil, fmt.Errorf("video has no page information")
	}
	// 映射到领域模型页面
	domainPages := make([]model.VideoPage, 0, len(videoView.Pages))
	for _, dtoPage := range videoView.Pages {
		domainPages = append(domainPages, model.VideoPage{
			Cid:      dtoPage.Cid,
			Duration: dtoPage.Duration,
			Part:     dtoPage.Part,
			Page:     dtoPage.Page,
		})
	}
	actualAID := videoView.Aid // 确定视频的 AID 用于后续查询

	// 2. 获取时间范围内的所有进度记录，增加 buffer
	// 查询 [overallStartTime - interval, overallEndTime] 的记录
	queryStartTime := overallStartTime.Add(-interval)
	log.Printf("Querying progress records for AID %d in range [%s, %s]", actualAID, queryStartTime, overallEndTime)
	progressRecords, err := s.progressRepo.ListByAIDAndTimestampRange(ctx, actualAID, queryStartTime, overallEndTime)
	if err != nil {
		return nil, fmt.Errorf("failed to list progress records: %w", err)
	}

	if len(progressRecords) < 1 { // 如果整个扩展范围内都没有记录，则无法计算
		log.Printf("No progress records found for AID %d in extended range [%s, %s]", actualAID, queryStartTime, overallEndTime)
		return []WatchedSegmentResult{}, nil
	}

	// 3. 按时间间隔处理并计算
	var results []WatchedSegmentResult
	segmentStart := overallStartTime

	for segmentStart.Before(overallEndTime) {
		segmentEnd := segmentStart.Add(interval)
		if segmentEnd.After(overallEndTime) {
			segmentEnd = overallEndTime
		}

		var p1 *model.VideoProgress // Represents state at/before segmentStart
		for i := len(progressRecords) - 1; i >= 0; i-- {
			if !progressRecords[i].RecordedAt.After(segmentStart) {
				p1 = progressRecords[i]
				break
			}
		}

		var p2 *model.VideoProgress // Represents state at/before segmentEnd
		for i := len(progressRecords) - 1; i >= 0; i-- {
			if !progressRecords[i].RecordedAt.After(segmentEnd) {
				p2 = progressRecords[i]
				break
			}
		}

		var watchedDuration time.Duration
		var calcErr error

		if p2 == nil {
			// No records at or before segmentEnd. No activity relevant to the end of this segment.
			watchedDuration = 0
		} else if p1 == nil {
			// No record at or before segmentStart. All records in progressRecords are > segmentStart.
			// Playback, if any for this segment, started *after* segmentStart.
			// p2 is the last known state within this segment (or at its end).
			if len(progressRecords) == 0 { // Should have been caught by earlier check, but defensive.
				watchedDuration = 0
			} else {
				actualStartPointForCalc := progressRecords[0] // Earliest record, must be > segmentStart

				// Ensure this actualStartPointForCalc is relevant for the segment ending at p2.RecordedAt
				// and it must occur before or at p2, and start before segmentEnd.
				if !actualStartPointForCalc.RecordedAt.After(p2.RecordedAt) && actualStartPointForCalc.RecordedAt.Before(segmentEnd) {
					watchedDuration, calcErr = s.calculator.CalculateWatchTime(
						domainPages,
						actualStartPointForCalc.LastPlayCID,
						actualStartPointForCalc.LastPlayTime/1000, // ms to s
						p2.LastPlayCID,
						p2.LastPlayTime/1000, // ms to s
					)
					if calcErr != nil {
						log.Printf("Error calculating watch time (p1 nil case) for segment [%s, %s], AID %d: %v. ActualStart: %+v, P2: %+v",
							segmentStart, segmentEnd, actualAID, calcErr, actualStartPointForCalc, p2)
						watchedDuration = 0
					}
				} else {
					// actualStartPointForCalc is either after p2 or not meaningfully before segmentEnd.
					watchedDuration = 0
				}
			}
		} else if p1.ID == p2.ID {
			// The same record defines the state at segmentStart and segmentEnd boundary checks.
			// This implies no new progress records fell strictly within (p1.RecordedAt, segmentEnd) that would change p2.
			watchedDuration = 0
		} else if p1.RecordedAt.After(p2.RecordedAt) {
			log.Printf("Data integrity concern: p1.RecordedAt (%v) is after p2.RecordedAt (%v) for segment [%s, %s]",
				p1.RecordedAt, p2.RecordedAt, segmentStart, segmentEnd)
			watchedDuration = 0
		} else {
			// Normal case: p1 and p2 are valid, distinct, and p1 is not after p2.
			watchedDuration, calcErr = s.calculator.CalculateWatchTime(
				domainPages,
				p1.LastPlayCID,
				p1.LastPlayTime/1000, // ms to s
				p2.LastPlayCID,
				p2.LastPlayTime/1000, // ms to s
			)
			if calcErr != nil {
				log.Printf("Error calculating watch time (normal case) for segment [%s, %s], AID %d: %v. P1: %+v, P2: %+v",
					segmentStart, segmentEnd, actualAID, calcErr, p1, p2)
				watchedDuration = 0
			}
		}

		results = append(results, WatchedSegmentResult{
			SegmentStartTime: segmentStart,
			SegmentEndTime:   segmentEnd,
			WatchedDuration:  watchedDuration,
		})

		segmentStart = segmentEnd
	}

	return results, nil
}
