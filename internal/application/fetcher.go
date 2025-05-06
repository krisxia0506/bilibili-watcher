package application

import "context"

// FetchedProgressData 保存从外部源获取的基本进度数据。
type FetchedProgressData struct {
	// 视频稿件 ID (AV 号)
	AID int64
	// 视频 BV 号
	BVID string
	// 上次播放时间/进度 (毫秒)
	LastPlayTime int64
	// 上次播放的视频分 P ID
	LastPlayCid int64
}

// VideoProgressFetcher 定义了从外部源获取视频进度的接口。
type VideoProgressFetcher interface {
	// Fetch 获取指定视频的最新进度。
	Fetch(ctx context.Context, aid, cid string) (*FetchedProgressData, error)
}
