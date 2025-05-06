package bilibili

import (
	"context"
	"fmt"
	"log"

	"github.com/krisxia0506/bilibili-watcher/internal/application"
)

// bilibiliFetcher 使用 Bilibili client 实现 application.VideoProgressFetcher 接口。
type bilibiliFetcher struct {
	client *Client
}

// NewBilibiliFetcher 创建一个新的 fetcher 实例。
func NewBilibiliFetcher(client *Client) application.VideoProgressFetcher {
	return &bilibiliFetcher{client: client}
}

// Fetch 使用 Bilibili API client 获取视频进度。
func (f *bilibiliFetcher) Fetch(ctx context.Context, aid, cid string) (*application.FetchedProgressData, error) {
	resp, err := f.client.GetVideoProgress(aid, cid)
	if err != nil {
		// 错误已在 client 中记录，包装后返回给 service 层
		return nil, fmt.Errorf("bilibili API client error: %w", err)
	}

	// 对响应数据结构进行基本验证
	if resp == nil { // 首先检查 resp 本身是否为 nil
		log.Printf("Received nil response from Bilibili API for AID %s, CID %s", aid, cid)
		return nil, fmt.Errorf("received nil response from Bilibili API")
	}
	// 我们依赖 client 已经检查了响应码 (resp.Code == 0)

	// 即使响应结构有效，进度信息也可能无意义
	if resp.Data.LastPlayTime <= 0 {
		log.Printf("Received non-positive progress (%d) for AID %s, CID %s", resp.Data.LastPlayTime, aid, cid)
		// 决定这是错误还是仅仅表示没有进度。暂时视为非错误，
		// 但返回 nil 数据以表示未找到有意义的进度。
		// 应用服务可以决定是保存 0 进度记录还是跳过。
		// 在这里返回错误可能会阻止后续保存合法的 0 进度（如果需要的话）。
		return nil, nil // 没有错误，但没有有效的进度数据
	}

	// 将基础设施特定的响应映射到应用程序定义的数据结构
	appData := &application.FetchedProgressData{
		AID:          resp.Data.Aid,
		BVID:         resp.Data.Bvid,
		LastPlayTime: resp.Data.LastPlayTime,
		LastPlayCid:  resp.Data.LastPlayCid,
	}

	return appData, nil
}
