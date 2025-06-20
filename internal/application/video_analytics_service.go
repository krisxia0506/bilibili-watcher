package application

import (
	"context"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/krisxia0506/bilibili-watcher/internal/domain/model"
	"github.com/krisxia0506/bilibili-watcher/internal/domain/repository"
	"github.com/krisxia0506/bilibili-watcher/internal/domain/service"
)

// WatchedSegmentResult 包含单个时间分段的计算结果。
type WatchedSegmentResult struct {
	SegmentStartTime time.Time
	SegmentEndTime   time.Time
	WatchedDuration  time.Duration
}

// VideoAnalyticsResult 包含视频分析的完整结果，包括分段和总时长。
type VideoAnalyticsResult struct {
	Segments             []WatchedSegmentResult
	TotalWatchedDuration time.Duration
}

// VideoAnalyticsService 定义了视频分析相关的应用服务接口。
type VideoAnalyticsService interface {
	// GetWatchedSegments 计算并返回指定时间范围和间隔内的视频观看分段时长及总时长。
	GetWatchedSegments(ctx context.Context,
		aidStr, bvidStr string, // aid 和 bvid 提供一个
		overallStartTime, overallEndTime time.Time,
		interval time.Duration,
	) (VideoAnalyticsResult, error)
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

// getSegmentStartTime 根据记录时间、整体开始时间和间隔，计算记录所属分段的开始时间。
// 返回分段开始时间和是否有效（即记录时间是否在 [overallStartTime, overallEndTime) 范围内）。
func getSegmentStartTime(recordTime, overallStartTime, overallEndTime time.Time, interval time.Duration) (time.Time, bool) {
	// 确保记录时间在有效的查询范围内 [overallStartTime, overallEndTime)
	if recordTime.Before(overallStartTime) || !recordTime.Before(overallEndTime) { // recordTime < overallStartTime || recordTime >= overallEndTime
		return time.Time{}, false
	}

	// 计算自 overallStartTime 以来经过的时间
	elapsed := recordTime.Sub(overallStartTime)
	if elapsed < 0 {
		// 理论上不应发生，因为上面有检查 recordTime.Before(overallStartTime)
		log.Printf("[getSegmentStartTime Error] Negative elapsed time calculated for %s relative to %s", recordTime, overallStartTime)
		return time.Time{}, false
	}

	// 计算记录时间落在哪个分段索引 (向下取整)
	segmentIndex := int64(elapsed / interval)
	calculatedSegmentStart := overallStartTime.Add(time.Duration(segmentIndex) * interval)

	// 最终校验，确保计算出的开始时间不会意外地超出范围
	if calculatedSegmentStart.Before(overallStartTime) || !calculatedSegmentStart.Before(overallEndTime) {
		log.Printf("[getSegmentStartTime Warning] Calculated segment start %s is outside overall range [%s, %s). Record time: %s", calculatedSegmentStart, overallStartTime, overallEndTime, recordTime)
		return time.Time{}, false
	}

	return calculatedSegmentStart, true
}

// GetWatchedSegments 实现获取观看分段的逻辑 (基于记录点迭代和归属)。
func (s *videoAnalyticsService) GetWatchedSegments(ctx context.Context,
	aidStr, bvidStr string,
	overallStartTime, overallEndTime time.Time,
	interval time.Duration,
) (VideoAnalyticsResult, error) {

	emptyResult := VideoAnalyticsResult{Segments: []WatchedSegmentResult{}, TotalWatchedDuration: 0}

	if aidStr == "" && bvidStr == "" {
		return emptyResult, fmt.Errorf("必须提供 aid 或 bvid")
	}
	if interval <= 0 {
		return emptyResult, fmt.Errorf("interval 必须为正数")
	}
	if overallEndTime.Before(overallStartTime) {
		return emptyResult, fmt.Errorf("结束时间必须在开始时间之后")
	}

	// 1. 获取视频页面信息
	videoView, err := s.biliClient.GetVideoView(ctx, aidStr, bvidStr)
	if err != nil {
		return emptyResult, fmt.Errorf("获取视频信息失败: %w", err)
	}
	if videoView == nil || len(videoView.Pages) == 0 {
		return emptyResult, fmt.Errorf("视频没有分页信息")
	}
	domainPages := make([]model.VideoPage, 0, len(videoView.Pages))
	for _, dtoPage := range videoView.Pages {
		domainPages = append(domainPages, model.VideoPage{
			Cid: dtoPage.Cid, Duration: dtoPage.Duration, Part: dtoPage.Part, Page: dtoPage.Page,
		})
	}
	actualAID := videoView.Aid

	// 2. 获取时间范围内的所有进度记录，增加足够 buffer
	queryStartTime := overallStartTime.Add(-interval * 2)
	queryEndTime := overallEndTime.Add(interval)
	log.Printf("查询 AID %d 在扩展时间范围 [%s, %s] 内的进度记录", actualAID, queryStartTime, queryEndTime)
	progressRecords, err := s.progressRepo.ListByAIDAndTimestampRange(ctx, actualAID, queryStartTime, queryEndTime)
	if err != nil {
		return emptyResult, fmt.Errorf("列出进度记录失败: %w", err)
	}

	if len(progressRecords) < 2 {
		log.Printf("在扩展时间范围 [%s, %s] 内找到的记录少于2条 (共 %d 条)，无法计算观看时长", queryStartTime, queryEndTime, len(progressRecords))
		results := make([]WatchedSegmentResult, 0)
		segmentStart := overallStartTime
		for segmentStart.Before(overallEndTime) {
			segmentEnd := segmentStart.Add(interval)
			if segmentEnd.After(overallEndTime) {
				segmentEnd = overallEndTime
			}
			results = append(results, WatchedSegmentResult{
				SegmentStartTime: segmentStart,
				SegmentEndTime:   segmentEnd,
				WatchedDuration:  0,
			})
			segmentStart = segmentEnd
		}
		return VideoAnalyticsResult{Segments: results, TotalWatchedDuration: 0}, nil
	}

	// 3. 初始化分段时长 map
	segmentDurations := make(map[time.Time]time.Duration)
	segmentStartMapKey := overallStartTime
	for segmentStartMapKey.Before(overallEndTime) {
		segmentDurations[segmentStartMapKey] = 0
		segmentStartMapKey = segmentStartMapKey.Add(interval)
	}

	// 4. 迭代记录点，计算时长并归属到对应分段
	for i := 0; i < len(progressRecords)-1; i++ {
		pCurr := progressRecords[i]
		pNext := progressRecords[i+1]

		// 计算这两个记录点之间的观看时长
		duration, calcErr := s.calculator.CalculateWatchTime(
			domainPages,
			pCurr.LastPlayCID, pCurr.LastPlayTime/1000, // 毫秒转秒
			pNext.LastPlayCID, pNext.LastPlayTime/1000, // 毫秒转秒
		)
		if calcErr != nil {
			log.Printf("计算观看时长出错 (记录 %d -> %d): %v. P_curr: %+v, P_next: %+v",
				pCurr.ID, pNext.ID, calcErr, pCurr, pNext)
			continue // 跳过这一对记录
		}

		if duration <= 0 {
			continue // 没有时长变化，无需归属
		}

		// 确定时长发生区间的起始点 pCurr 属于哪个分段
		segmentStartTimeForAttribution, isValidSegment := getSegmentStartTime(pCurr.RecordedAt, overallStartTime, overallEndTime, interval)

		if isValidSegment {
			// 将计算出的时长累加到对应的分段
			if _, ok := segmentDurations[segmentStartTimeForAttribution]; ok {
				log.Printf("[Attribution] Attributing duration %s (from record %d at %s to record %d at %s) to segment starting at %s",
					duration, pCurr.ID, pCurr.RecordedAt, pNext.ID, pNext.RecordedAt, segmentStartTimeForAttribution)
				segmentDurations[segmentStartTimeForAttribution] += duration
			} else {
				log.Printf("[Attribution Warning] Segment start time %s (from record %d at %s) not found in map.", segmentStartTimeForAttribution, pCurr.ID, pCurr.RecordedAt)
			}
		} else {
			log.Printf("[Attribution] Skipping duration %s starting at %s (record %d) because it falls outside the requested range [%s, %s)",
				duration, pCurr.RecordedAt, pCurr.ID, overallStartTime, overallEndTime)
		}
	}

	// 5. 从 map 生成最终结果列表并计算总时长
	results := make([]WatchedSegmentResult, 0, len(segmentDurations))
	var totalWatchedDuration time.Duration
	segmentResultStart := overallStartTime
	processedSegments := make(map[time.Time]bool)

	for segmentResultStart.Before(overallEndTime) {
		segmentEnd := segmentResultStart.Add(interval)
		if segmentEnd.After(overallEndTime) {
			segmentEnd = overallEndTime
		}

		if !processedSegments[segmentResultStart] {
			duration := segmentDurations[segmentResultStart]
			results = append(results, WatchedSegmentResult{
				SegmentStartTime: segmentResultStart,
				SegmentEndTime:   segmentEnd,
				WatchedDuration:  duration,
			})
			totalWatchedDuration += duration // 累加到总时长
			log.Printf("[Final Result] Segment [%s, %s]: Duration=%s", segmentResultStart, segmentEnd, duration)
			processedSegments[segmentResultStart] = true
		}
		segmentResultStart = segmentEnd
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].SegmentStartTime.Before(results[j].SegmentStartTime)
	})

	log.Printf("[Total Duration] Calculated for range [%s, %s]: %s", overallStartTime, overallEndTime, totalWatchedDuration)
	return VideoAnalyticsResult{Segments: results, TotalWatchedDuration: totalWatchedDuration}, nil
}
