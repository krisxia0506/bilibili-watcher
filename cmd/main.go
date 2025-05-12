package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/krisxia0506/bilibili-watcher/internal/application"
	"github.com/krisxia0506/bilibili-watcher/internal/config"
	"github.com/krisxia0506/bilibili-watcher/internal/domain/service"
	"github.com/krisxia0506/bilibili-watcher/internal/infrastructure/bilibili"
	"github.com/krisxia0506/bilibili-watcher/internal/infrastructure/persistence"
	"github.com/krisxia0506/bilibili-watcher/internal/infrastructure/scheduler"
	"github.com/krisxia0506/bilibili-watcher/internal/interfaces/api/rest"
)

// main 程序入口
func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 初始化数据库连接
	db, err := persistence.NewDatabaseConnection(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Successfully connected to the database.")

	// --- 初始化基础设施组件 ---
	biliClient := bilibili.NewClient(cfg.Bilibili.SessData)
	log.Println("Bilibili client initialized.")
	videoProgressRepo := persistence.NewGormVideoProgressRepository(db)
	log.Println("Video progress repository initialized.")

	// --- 初始化领域服务 ---
	watchTimeCalculator := service.NewWatchTimeCalculator()
	log.Println("Watch time calculator initialized.")

	// --- 初始化应用服务 ---
	videoProgressService := application.NewVideoProgressService(videoProgressRepo, biliClient)
	log.Println("Video progress service initialized.")
	log.Println("Watch time service initialized.")
	videoAnalyticsService := application.NewVideoAnalyticsService(biliClient, videoProgressRepo, watchTimeCalculator)
	log.Println("Video analytics service initialized.")

	// --- 设置 Gin 和路由 (使用新的 rest 包) ---
	router := rest.SetupRouter(db, cfg.GinMode, videoAnalyticsService /*, other services */)

	// --- 初始化并启动调度器 ---
	appScheduler := scheduler.NewScheduler()

	// 创建获取单个 BVID 视频进度的函数
	createFetchVideoProgressJobForBVID := func(bvid string) func() {
		return func() {
			if bvid == "" {
				log.Println("Error: Empty BVID provided. Skipping progress fetch job.")
				return
			}

			jobName := fmt.Sprintf("FetchVideoProgress(BVID: %s)", bvid)
			log.Printf("Cron job starting: %s", jobName)
			ctx := context.Background()

			// 1. 获取视频的 AID 和第一个 CID
			videoView, err := biliClient.GetVideoView(ctx, "", bvid)
			if err != nil {
				log.Printf("Error fetching video view for BVID %s in job '%s': %v. Skipping progress fetch.", bvid, jobName, err)
				return
			}
			if videoView == nil || len(videoView.Pages) == 0 {
				log.Printf("No pages found for video BVID %s in job '%s'. Skipping progress fetch.", bvid, jobName)
				return
			}
			targetCID := videoView.Pages[0].Cid
			log.Printf("Determined targetCID: %d for BVID: %s", targetCID, bvid)

			// 2. 获取并保存进度
			if err := videoProgressService.FetchAndSaveVideoProgress(ctx, "", bvid, strconv.FormatInt(targetCID, 10)); err != nil {
				log.Printf("Error executing FetchAndSaveVideoProgress for BVID '%s', CID %d in job '%s': %v", bvid, targetCID, jobName, err)
			}
			log.Printf("Cron job finished: %s (processed BVID: %s, CID: %d)", jobName, bvid, targetCID)
		}
	}

	// 为每个 BVID 创建单独的定时任务
	if len(cfg.Bilibili.TargetBVIDs) > 0 {
		for _, bvid := range cfg.Bilibili.TargetBVIDs {
			jobName := fmt.Sprintf("FetchVideoProgress_BVID_%s", bvid)
			if err := appScheduler.ScheduleJob(jobName, cfg.Scheduler.Cron, createFetchVideoProgressJobForBVID(bvid)); err != nil {
				log.Printf("Failed to schedule job '%s' for BVID '%s': %v", jobName, bvid, err)
			}
		}
	}

	go appScheduler.Start() // 在单独的 goroutine 中启动调度器

	// --- 启动 Gin 服务器 ---
	serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: router, // 使用 rest.SetupRouter 返回的 router
	}

	go func() {
		// 服务连接
		// 使用 http.Server 的方式是为了支持优雅停机，这是比直接使用 router.Run() 更健壮的做法
		log.Printf("Starting server on %s", serverAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// --- 等待中断信号以优雅关闭服务器和调度器 ---
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// 停止调度器
	schedulerCtx := appScheduler.Stop()
	<-schedulerCtx.Done() // 等待调度器任务完成
	log.Println("Scheduler stopped.")

	// context 用于通知服务器它有 5 秒钟时间来处理当前正在处理的请求
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
