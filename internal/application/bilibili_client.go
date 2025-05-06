package application

import (
	"context"
)

// BilibiliClient 定义了应用层与 Bilibili API 交互所需的操作。
// 基础设施层需要实现此接口。
type BilibiliClient interface {
	// GetVideoProgress 获取指定视频的观看进度。
	// aid (视频稿件 avid) 和 bvid (视频稿件 bvid) 必须提供一个。
	// cid (视频分P的 ID) 必须提供。
	// 返回应用层定义的 DTO。
	GetVideoProgress(ctx context.Context, aid, bvid, cid string) (*VideoProgressDTO, error)

	// GetVideoView 获取视频的详细信息。
	// 返回应用层定义的 DTO。
	// aid: 视频的 AV 号 (可选)
	// bvid: 视频的 BV 号 (可选)
	// aid 和 bvid 必须提供一个。
	GetVideoView(ctx context.Context, aid, bvid string) (*VideoViewDTO, error)

	// TODO: 未来可以添加更多 Bilibili API 方法
}
