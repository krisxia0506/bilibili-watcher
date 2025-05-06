package service

import (
	"errors"
	"time"

	"github.com/krisxia0506/bilibili-watcher/internal/domain/model"
)

var (
	ErrPageNotFound    = errors.New("start or end page (CID) not found in video")
	ErrInvalidTime     = errors.New("start or end time is invalid (negative or exceeds duration)")
	ErrStartAfterEnd   = errors.New("start point is chronologically after end point")
	ErrIdenticalPoints = errors.New("start and end points are identical")
)

// WatchTimeCalculator 定义了计算视频观看时长的领域服务接口。
type WatchTimeCalculator interface {
	// CalculateWatchTime 计算从 startCid 的 startTimeInStartCidSec 到 endCid 的 endTimeInEndCidSec 之间的总观看时长。
	// pages: 视频的所有分P信息列表 (必须按播放顺序排列)。
	// startCid: 开始观看的分P ID。
	// startTimeInStartCidSec: 在 startCid 中开始观看的时间点（秒）。
	// endCid: 结束观看的分P ID。
	// endTimeInEndCidSec: 在 endCid 中结束观看的时间点（秒）。
	// 返回计算出的总观看时长 (time.Duration) 和可能的错误。
	// 可能的错误包括：找不到页面、时间无效、开始点晚于结束点。
	CalculateWatchTime(
		pages []model.VideoPage,
		startCid int64,
		startTimeInStartCidSec int64,
		endCid int64,
		endTimeInEndCidSec int64,
	) (time.Duration, error)
}
