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
		client: client, // 注入 BilibiliClient
	}
}

// FetchAndSaveVideoProgress 获取指定视频的观看进度并保存到仓库。
// 这是定时任务的主要执行逻辑。
func (s *VideoProgressService) FetchAndSaveVideoProgress(ctx context.Context, aidStr, cidStr string) error {
	log.Printf("Fetching progress for AID: %s, CID: %s", aidStr, cidStr)

	// 1. 调用 Bilibili Client 获取进度
	progressDTO, err := s.client.GetVideoProgress(ctx, aidStr, cidStr)
	if err != nil {
		// 底层 Client 实现已处理 API 错误和网络错误，这里只记录并返回包装后的错误
		log.Printf("Error fetching video progress from Bilibili client for AID %s, CID %s: %v", aidStr, cidStr, err)
		return fmt.Errorf("failed to fetch video progress from bilibili client: %w", err)
	}

	// 如果 DTO 为 nil，表示 API 调用成功但没有有效的进度信息 (例如进度为0)
	if progressDTO == nil {
		log.Printf("No valid progress data returned from Bilibili API for AID %s, CID %s. Skipping save.", aidStr, cidStr)
		// 根据业务需求，这里可以选择是否保存一条进度为 0 的记录
		// 当前选择跳过，不保存记录
		return nil
	}

	log.Printf("Successfully fetched progress for AID %d (BVID: %s): LastPlayTime=%dms, LastPlayCid=%d",
		progressDTO.AID, progressDTO.BVID, progressDTO.LastPlayTime, progressDTO.LastPlayCid)

	// 2. 将 DTO 转换为领域模型
	aid := progressDTO.AID

	// 创建新记录
	progressToSave := &model.VideoProgress{
		AID:          aid,
		BVID:         progressDTO.BVID, // 保存 BVID
		LastPlayCID:  progressDTO.LastPlayCid,              // 修复大小写
		LastPlayTime: progressDTO.LastPlayTime,
	}
	log.Printf("Creating new progress record for AID %d", aid)

	// 4. 保存到仓库
	if err := s.repo.Save(ctx, progressToSave); err != nil {
		log.Printf("Error saving video progress for AID %d: %v", aid, err)
		return fmt.Errorf("failed to save video progress: %w", err)
	}

	log.Printf("Successfully saved progress for AID %d", aid)
	return nil
}

// TODO: 添加其他应用服务方法，例如计算每日观看时长等
