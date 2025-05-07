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

	// 2. 获取时间范围内的所有进度记录 (+/- 一小段时间 buffer 防止边界问题? 暂时不加)
	// 为了找到每个区间的起始点，需要查询略微超出 start time 的记录
	// 查询 [overallStartTime - buffer, overallEndTime] 的记录可能更保险，buffer 可以是 interval
	// 但更简单的方法是查询整个范围，然后在内存中查找
	progressRecords, err := s.progressRepo.ListByAIDAndTimestampRange(ctx, actualAID, overallStartTime, overallEndTime)
	if err != nil {
		return nil, fmt.Errorf("failed to list progress records: %w", err)
	}
	// 可能还需要查询 overallStartTime 之前的最后一条记录作为初始状态
	// TODO: 优化查询，只获取必要数据
	if len(progressRecords) < 1 { // 如果范围内只有0或1条记录，无法计算时长
		log.Printf("Not enough progress records found for AID %d in range [%s, %s] to calculate duration", actualAID, overallStartTime, overallEndTime)
		return []WatchedSegmentResult{}, nil // 返回空结果，不是错误
	}

	// 3. 按时间间隔处理并计算
	var results []WatchedSegmentResult
	segmentStart := overallStartTime

	for segmentStart.Before(overallEndTime) {
		segmentEnd := segmentStart.Add(interval)
		// 确保最后一个分段的结束时间不超过 overallEndTime
		if segmentEnd.After(overallEndTime) {
			segmentEnd = overallEndTime
		}

		// 查找此分段的开始和结束进度记录
		var startProgress, endProgress *model.VideoProgress
		// 找到 recorded_at <= segmentStart 的最后一条记录
		for i := len(progressRecords) - 1; i >= 0; i-- {
			if !progressRecords[i].RecordedAt.After(segmentStart) {
				startProgress = progressRecords[i]
				break
			}
		}
		// 找到 recorded_at <= segmentEnd 的最后一条记录
		for i := len(progressRecords) - 1; i >= 0; i-- {
			if !progressRecords[i].RecordedAt.After(segmentEnd) {
				endProgress = progressRecords[i]
				break
			}
		}

		var watchedDuration time.Duration
		var calcErr error

		if startProgress != nil && endProgress != nil && startProgress.ID != endProgress.ID { // 必须是不同的记录
			// 确保 startProgress 在 endProgress 之前（或同一时间）
			if !startProgress.RecordedAt.After(endProgress.RecordedAt) {
				watchedDuration, calcErr = s.calculator.CalculateWatchTime(
					domainPages,
					startProgress.LastPlayCID,
					startProgress.LastPlayTime/1000, // ms to s
					endProgress.LastPlayCID,
					endProgress.LastPlayTime/1000, // ms to s
				)
				if calcErr != nil {
					// 记录领域计算错误，但继续处理下一个分段
					log.Printf("Error calculating watch time for segment [%s, %s], AID %d: %v", segmentStart, segmentEnd, actualAID, calcErr)
					watchedDuration = 0 // 出错时时长计为0
				}
			} else {
				// 如果 start 在 end 之后，说明这段时间没有有效观看记录
				watchedDuration = 0
			}
		} else {
			// 没有找到开始或结束记录，或者它们是同一条记录，时长为 0
			watchedDuration = 0
		}

		results = append(results, WatchedSegmentResult{
			SegmentStartTime: segmentStart,
			SegmentEndTime:   segmentEnd,
			WatchedDuration:  watchedDuration,
		})

		// 移动到下一个分段的开始时间
		segmentStart = segmentEnd
	}

	return results, nil
}
