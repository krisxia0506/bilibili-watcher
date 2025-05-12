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
	biliClient   BilibiliClient                     // Bilibili 客户端接口，用于获取视频页面信息
	progressRepo repository.VideoProgressRepository // 视频进度仓库接口，用于获取进度记录
	calculator   service.WatchTimeCalculator        // 观看时长计算器服务
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
		return nil, fmt.Errorf("必须提供 aid 或 bvid")
	}
	if interval <= 0 {
		return nil, fmt.Errorf("interval 必须为正数")
	}
	if overallEndTime.Before(overallStartTime) {
		return nil, fmt.Errorf("结束时间必须在开始时间之后")
	}

	// 1. 获取视频页面信息 (只需要执行一次)
	videoView, err := s.biliClient.GetVideoView(ctx, aidStr, bvidStr)
	if err != nil {
		return nil, fmt.Errorf("获取视频信息失败: %w", err)
	}
	if videoView == nil || len(videoView.Pages) == 0 {
		return nil, fmt.Errorf("视频没有分页信息")
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

	// 2. 获取时间范围内的所有进度记录，增加 buffer 时间
	// 查询 [overallStartTime - interval, overallEndTime] 的记录
	queryStartTime := overallStartTime.Add(-interval)
	log.Printf("查询 AID %d 在时间范围 [%s, %s] 内的进度记录", actualAID, queryStartTime, overallEndTime)
	progressRecords, err := s.progressRepo.ListByAIDAndTimestampRange(ctx, actualAID, queryStartTime, overallEndTime)
	if err != nil {
		return nil, fmt.Errorf("列出进度记录失败: %w", err)
	}

	if len(progressRecords) < 1 { // 如果整个扩展范围内都没有记录，则无法计算
		log.Printf("在扩展时间范围 [%s, %s] 内未找到 AID %d 的进度记录", queryStartTime, overallEndTime, actualAID)
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

		var p1 *model.VideoProgress // 代表 segmentStart 时间点（或之前）的状态
		for i := len(progressRecords) - 1; i >= 0; i-- {
			if !progressRecords[i].RecordedAt.After(segmentStart) {
				p1 = progressRecords[i]
				break
			}
		}

		var p2 *model.VideoProgress // 代表 segmentEnd 时间点（或之前）的状态
		for i := len(progressRecords) - 1; i >= 0; i-- {
			if !progressRecords[i].RecordedAt.After(segmentEnd) {
				p2 = progressRecords[i]
				break
			}
		}

		var watchedDuration time.Duration
		var calcErr error

		if p2 == nil {
			// 在 segmentEnd 或之前没有记录。没有与此分段结束相关的活动。
			watchedDuration = 0
		} else if p1 == nil {
			// 在 segmentStart 或之前未找到记录。
			// 查找第一个时间戳大于等于 segmentStart 的记录，作为此分段活动的开始。
			var firstRecordInSegment *model.VideoProgress
			for _, rec := range progressRecords { // 假设 progressRecords 按 RecordedAt 排序
				if !rec.RecordedAt.Before(segmentStart) { // 等价于 rec.RecordedAt >= segmentStart
					firstRecordInSegment = rec
					break
				}
			}

			if firstRecordInSegment != nil && !firstRecordInSegment.RecordedAt.After(p2.RecordedAt) {
				// 确保在分段内开始的活动实际上发生在分段结束状态 p2 之前或同时。
				// 并且重要的是，确保 firstRecordInSegment 发生在 segmentEnd 边界之前。
				if firstRecordInSegment.RecordedAt.Before(segmentEnd) {
					watchedDuration, calcErr = s.calculator.CalculateWatchTime(
						domainPages,
						firstRecordInSegment.LastPlayCID,
						firstRecordInSegment.LastPlayTime/1000, // 毫秒转秒
						p2.LastPlayCID,
						p2.LastPlayTime/1000, // 毫秒转秒
					)
					if calcErr != nil {
						log.Printf("计算观看时长出错 (p1 为 nil 情况，已修正) 分段 [%s, %s], AID %d: %v. 段内首记录: %+v, P2: %+v",
							segmentStart, segmentEnd, actualAID, calcErr, firstRecordInSegment, p2)
						watchedDuration = 0
					}
				} else {
					// 第一个记录发生在 segmentEnd 或之后，此分段内没有观看时长。
					watchedDuration = 0
				}

			} else {
				// 分段内没有相关的活动开始，或者 p2 无效/更早。
				watchedDuration = 0
			}
		} else if p1.ID == p2.ID {
			// 同一个记录定义了 segmentStart 和 segmentEnd 边界检查的状态。
			// 这意味着在 (p1.RecordedAt, segmentEnd] 区间内没有新的进度记录会改变 p2。
			watchedDuration = 0
		} else if p1.RecordedAt.After(p2.RecordedAt) {
			// 数据完整性问题：对于分段 [%s, %s]，p1.RecordedAt (%v) 在 p2.RecordedAt (%v) 之后
			log.Printf("数据完整性问题: p1.RecordedAt (%v) 在 p2.RecordedAt (%v) 之后，分段 [%s, %s]",
				p1.RecordedAt, p2.RecordedAt, segmentStart, segmentEnd)
			watchedDuration = 0
		} else {
			// 正常情况：p1 和 p2 都有效，不同，且 p1 不在 p2 之后。
			watchedDuration, calcErr = s.calculator.CalculateWatchTime(
				domainPages,
				p1.LastPlayCID,
				p1.LastPlayTime/1000, // 毫秒转秒
				p2.LastPlayCID,
				p2.LastPlayTime/1000, // 毫秒转秒
			)
			if calcErr != nil {
				log.Printf("计算观看时长出错 (正常情况) 分段 [%s, %s], AID %d: %v. P1: %+v, P2: %+v",
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
