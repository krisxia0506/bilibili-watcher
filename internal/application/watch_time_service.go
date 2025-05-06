package application

import (
	"context"
	"fmt"
	"time"

	"github.com/krisxia0506/bilibili-watcher/internal/domain/model"
	"github.com/krisxia0506/bilibili-watcher/internal/domain/service"
)

// WatchTimeService 定义了处理观看时长计算相关用例的应用服务接口。
type WatchTimeService interface {
	// CalculateWatchTimeBetweenPoints 计算指定视频在两个时间点之间的观看时长。
	// aid 和 bvid 提供一个即可。
	CalculateWatchTimeBetweenPoints(ctx context.Context,
		aid, bvid string,
		startCid int64,
		startTimeInStartCidSec int64,
		endCid int64,
		endTimeInEndCidSec int64,
	) (time.Duration, error)
}

// watchTimeService 实现了 WatchTimeService 接口。
type watchTimeService struct {
	client     BilibiliClient              // 依赖 Bilibili Client 获取视频信息
	calculator service.WatchTimeCalculator // 依赖领域服务进行计算
}

// NewWatchTimeService 创建 WatchTimeService 实例。
func NewWatchTimeService(client BilibiliClient, calculator service.WatchTimeCalculator) WatchTimeService {
	return &watchTimeService{
		client:     client,
		calculator: calculator,
	}
}

// CalculateWatchTimeBetweenPoints 实现计算逻辑。
func (s *watchTimeService) CalculateWatchTimeBetweenPoints(ctx context.Context,
	aid, bvid string,
	startCid int64,
	startTimeInStartCidSec int64,
	endCid int64,
	endTimeInEndCidSec int64,
) (time.Duration, error) {

	// 1. 输入验证 (基础)
	if aid == "" && bvid == "" {
		return 0, fmt.Errorf("either aid or bvid must be provided")
	}

	// 2. 调用 Bilibili Client 获取视频视图信息 (包含 Pages)
	videoViewDTO, err := s.client.GetVideoView(ctx, aid, bvid)
	if err != nil {
		return 0, fmt.Errorf("failed to get video view info: %w", err)
	}
	if videoViewDTO == nil {
		// 理论上 client 不会返回 nil DTO 和 nil error，但加个保护
		return 0, fmt.Errorf("received nil video view info from client")
	}

	// 3. 将应用层 DTO 的 Pages 映射到领域层模型 Pages
	if len(videoViewDTO.Pages) == 0 {
		// 如果视频没有分P信息，无法计算
		return 0, fmt.Errorf("video has no page information (pages is empty)")
	}
	domainPages := make([]model.VideoPage, 0, len(videoViewDTO.Pages))
	for _, dtoPage := range videoViewDTO.Pages {
		domainPages = append(domainPages, model.VideoPage{
			Cid:      dtoPage.Cid,
			Duration: dtoPage.Duration,
			Part:     dtoPage.Part,
			Page:     dtoPage.Page,
		})
	}

	// 4. 调用领域服务进行计算
	watchDuration, err := s.calculator.CalculateWatchTime(
		domainPages,
		startCid,
		startTimeInStartCidSec,
		endCid,
		endTimeInEndCidSec,
	)
	if err != nil {
		// 直接返回领域服务计算出的错误（例如 PageNotFound, InvalidTime 等）
		return 0, err
	}

	return watchDuration, nil
}
