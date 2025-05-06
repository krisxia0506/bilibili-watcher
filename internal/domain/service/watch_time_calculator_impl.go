package service

import (
	"fmt"
	"time"

	"github.com/krisxia0506/bilibili-watcher/internal/domain/model"
)

// watchTimeCalculatorService 实现了 WatchTimeCalculator 接口。
type watchTimeCalculatorService struct{}

// NewWatchTimeCalculator 创建 WatchTimeCalculator 服务实例。
func NewWatchTimeCalculator() WatchTimeCalculator {
	return &watchTimeCalculatorService{}
}

// findPageIndexByCid 在页面列表中查找指定 CID 的索引。
func findPageIndexByCid(pages []model.VideoPage, cid int64) (int, bool) {
	for i, p := range pages {
		if p.Cid == cid {
			return i, true
		}
	}
	return -1, false
}

// CalculateWatchTime 计算观看时长。
func (s *watchTimeCalculatorService) CalculateWatchTime(
	pages []model.VideoPage,
	startCid int64,
	startTimeInStartCidSec int64,
	endCid int64,
	endTimeInEndCidSec int64,
) (time.Duration, error) {

	startIndex, startFound := findPageIndexByCid(pages, startCid)
	endIndex, endFound := findPageIndexByCid(pages, endCid)

	if !startFound || !endFound {
		return 0, ErrPageNotFound
	}

	// 检查时间有效性
	if startTimeInStartCidSec < 0 || startTimeInStartCidSec > pages[startIndex].Duration ||
		endTimeInEndCidSec < 0 || endTimeInEndCidSec > pages[endIndex].Duration {
		return 0, ErrInvalidTime
	}

	// 检查开始点是否在结束点之前或相同
	if startIndex > endIndex || (startIndex == endIndex && startTimeInStartCidSec > endTimeInEndCidSec) {
		return 0, ErrStartAfterEnd
	}

	// 检查开始点和结束点是否完全相同
	if startIndex == endIndex && startTimeInStartCidSec == endTimeInEndCidSec {
		// return 0, ErrIdenticalPoints // 返回0时长更符合逻辑
		return 0, nil
	}

	var totalDurationSeconds int64 = 0

	if startIndex == endIndex {
		// 开始和结束在同一个分P
		totalDurationSeconds = endTimeInEndCidSec - startTimeInStartCidSec
	} else {
		// 开始和结束在不同的分P

		// 1. 计算在开始分P观看的时间
		totalDurationSeconds += pages[startIndex].Duration - startTimeInStartCidSec

		// 2. 计算中间完整观看的分P时间
		for i := startIndex + 1; i < endIndex; i++ {
			totalDurationSeconds += pages[i].Duration
		}

		// 3. 计算在结束分P观看的时间
		totalDurationSeconds += endTimeInEndCidSec
	}

	if totalDurationSeconds < 0 {
		// 理论上不应发生，但作为健壮性检查
		return 0, fmt.Errorf("internal calculation error: negative duration calculated")
	}

	return time.Duration(totalDurationSeconds) * time.Second, nil
}
