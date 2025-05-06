package application

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/krisxia0506/bilibili-watcher/internal/domain/model"
	"github.com/krisxia0506/bilibili-watcher/internal/domain/repository"
)

// VideoProgressService 定义了用于管理视频进度的应用服务接口。
type VideoProgressService interface {
	// RecordProgressForTargetVideo 获取并记录预定义目标视频的进度。
	RecordProgressForTargetVideo(ctx context.Context) error
	// TODO: 添加管理目标视频的方法 (增/删/查)
	// TODO: 添加获取聚合观看时长数据的方法
}

// videoProgressService 实现了 VideoProgressService 接口。
type videoProgressService struct {
	videoRepo repository.VideoProgressRepository
	fetcher   VideoProgressFetcher // 依赖于应用层定义的 fetcher 接口
	// TODO: 使目标视频可配置
	targetAID string
	targetCID string
}

// NewVideoProgressService 创建一个新的 VideoProgressService 实例。
func NewVideoProgressService(repo repository.VideoProgressRepository, fetcher VideoProgressFetcher) VideoProgressService {
	// TODO: 从配置或其他来源加载 targetAID 和 targetCID
	return &videoProgressService{
		videoRepo: repo,
		fetcher:   fetcher,
		targetAID: "114102919764678", // 从用户示例硬编码
		targetCID: "28682552616",     // 从用户示例硬编码
	}
}

// RecordProgressForTargetVideo 获取并记录硬编码目标视频的进度。
func (s *videoProgressService) RecordProgressForTargetVideo(ctx context.Context) error {
	log.Printf("Fetching progress for AID: %s, CID: %s", s.targetAID, s.targetCID)

	// 使用 fetcher 接口获取进度数据
	fetchedData, err := s.fetcher.Fetch(ctx, s.targetAID, s.targetCID)
	if err != nil {
		// 记录来自 fetcher 的错误
		log.Printf("Error fetching video progress via fetcher: %v", err)
		// 根据 fetcher 返回的错误类型，可能需要特定处理。
		return fmt.Errorf("failed to fetch video progress: %w", err)
	}

	// 检查 fetcher 是否返回了有效数据 (fetcher 实现在进度非正数时返回 nil data)
	if fetchedData == nil {
		log.Printf("Fetcher returned no valid progress data for AID %s, CID %s. Skipping save.", s.targetAID, s.targetCID)
		// 这不是一个错误，只是没有新的数据需要记录。
		return nil
	}

	log.Printf("Successfully fetched progress: %d ms (BVID: %s) for AID %s, CID %s",
		fetchedData.LastPlayTime, fetchedData.BVID, s.targetAID, s.targetCID)

	// 直接使用从 fetcher 获取的数据填充模型
	progressRecord := &model.VideoProgress{
		AID:          fetchedData.AID,
		BVID:         fetchedData.BVID,
		LastPlayCID:  fetchedData.LastPlayCid,  // 更新的字段名
		LastPlayTime: fetchedData.LastPlayTime, // 更新的字段名
		RecordedAt:   time.Now(),               // 记录获取时的时间
		// GmtCreate 和 GmtModified 将由 GORM 处理
	}

	if err := s.videoRepo.Save(ctx, progressRecord); err != nil {
		log.Printf("Error saving video progress to database for AID %s, CID %s: %v", s.targetAID, s.targetCID, err)
		return fmt.Errorf("failed to save progress record: %w", err)
	}

	log.Printf("Successfully saved progress record ID: %d for AID %s, CID %s", progressRecord.ID, s.targetAID, s.targetCID)
	return nil
}
