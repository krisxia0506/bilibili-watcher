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

		// 查找此分段的开始和结束进度记录
		var startProgress, endProgress *model.VideoProgress
		// 找到 recorded_at <= segmentStart 的最后一条记录 (在所有获取到的记录里找)
		for i := len(progressRecords) - 1; i >= 0; i-- {
			if !progressRecords[i].RecordedAt.After(segmentStart) {
				startProgress = progressRecords[i]
				break
			}
		}
		// 找到 recorded_at <= segmentEnd 的最后一条记录 (在所有获取到的记录里找)
		for i := len(progressRecords) - 1; i >= 0; i-- {
			if !progressRecords[i].RecordedAt.After(segmentEnd) {
				endProgress = progressRecords[i]
				break
			}
		}

		var watchedDuration time.Duration
		var calcErr error

		if startProgress != nil && endProgress != nil && startProgress.ID != endProgress.ID {
			if !startProgress.RecordedAt.After(endProgress.RecordedAt) {
				watchedDuration, calcErr = s.calculator.CalculateWatchTime(
					domainPages,
					startProgress.LastPlayCID,
					startProgress.LastPlayTime/1000, // ms to s
					endProgress.LastPlayCID,
					endProgress.LastPlayTime/1000, // ms to s
				)
				if calcErr != nil {
					log.Printf("Error calculating watch time for segment [%s, %s], AID %d: %v", segmentStart, segmentEnd, actualAID, calcErr)
					watchedDuration = 0
				}
			} else {
				watchedDuration = 0
			}
		} else {
			watchedDuration = 0
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
