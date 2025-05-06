package application

import (
	"context"
	"fmt"
	"log"

	"github.com/krisxia0506/bilibili-watcher/internal/domain/model"
	"github.com/krisxia0506/bilibili-watcher/internal/domain/repository"
)

// VideoProgressService 应用服务，处理视频进度相关的用例。
type VideoProgressService struct {
	repo   repository.VideoProgressRepository
	client BilibiliClient // 使用新的通用 Bilibili Client 接口
}

// NewVideoProgressService 创建 VideoProgressService 实例。
func NewVideoProgressService(repo repository.VideoProgressRepository, client BilibiliClient) *VideoProgressService {
	return &VideoProgressService{
		repo:   repo,
		client: client,
	}
}

// FetchAndSaveVideoProgress 获取指定视频的观看进度并创建一条新的进度记录。
// aidStr (视频稿件 avid) 和 bvidStr (视频稿件 bvid) 必须提供一个。
// cidStr (视频分P的 ID) 必须提供。
// 此方法总是创建新的进度记录，不包含 TotalDuration 和 FetchTime。
func (s *VideoProgressService) FetchAndSaveVideoProgress(ctx context.Context, aidStr, bvidStr, cidStr string) error {
	log.Printf("Service: Fetching progress for AID: '%s', BVID: '%s', CID: '%s'", aidStr, bvidStr, cidStr)

	// 1. 调用 Bilibili Client 获取进度
	progressDTO, err := s.client.GetVideoProgress(ctx, aidStr, bvidStr, cidStr)
	if err != nil {
		log.Printf("Error fetching video progress from Bilibili client (AID: '%s', BVID: '%s', CID: '%s'): %v", aidStr, bvidStr, cidStr, err)
		return fmt.Errorf("failed to fetch video progress from bilibili client: %w", err)
	}

	if progressDTO == nil {
		log.Printf("No valid progress data returned from Bilibili API (AID: '%s', BVID: '%s', CID: '%s'). Skipping save.", aidStr, bvidStr, cidStr)
		return nil
	}

	log.Printf("Successfully fetched progress for AID %d (BVID: %s): LastPlayTime=%dms, LastPlayCid=%d",
		progressDTO.AID, progressDTO.BVID, progressDTO.LastPlayTime, progressDTO.LastPlayCid)

	// 2. 将 DTO 转换为领域模型的核心部分
	aid := progressDTO.AID
	bvid := progressDTO.BVID
	cid := progressDTO.LastPlayCid
	progressMs := progressDTO.LastPlayTime
	// 3. 创建新的领域模型实例
	progressToSave := &model.VideoProgress{
		AID:          aid,
		BVID:         bvid,
		LastPlayCID:  cid,
		LastPlayTime: progressMs,
	}
	log.Printf("Creating new progress record for AID %d, BVID %s", aid, bvid)

	// 4. 保存到仓库
	if err := s.repo.Save(ctx, progressToSave); err != nil {
		log.Printf("Error saving new video progress for AID %d and BVID %s: %v", aid, bvid, err)
		return fmt.Errorf("failed to save new video progress: %w", err)
	}

	log.Printf("Successfully saved new progress record for AID %d and BVID %s (ID: %d)", aid, bvid, progressToSave.ID)
	return nil
}

// TODO: 添加其他应用服务方法，例如计算每日观看时长等
