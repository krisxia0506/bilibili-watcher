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

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/krisxia0506/bilibili-watcher/internal/application"
	"github.com/krisxia0506/bilibili-watcher/internal/config"
	"github.com/krisxia0506/bilibili-watcher/internal/domain/service"
	"github.com/krisxia0506/bilibili-watcher/internal/infrastructure/bilibili"
	"github.com/krisxia0506/bilibili-watcher/internal/infrastructure/persistence"
	"github.com/krisxia0506/bilibili-watcher/internal/infrastructure/scheduler"
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
	// Bilibili 客户端 (底层 API 交互, 同时实现了 application.BilibiliClient)
	biliClient := bilibili.NewClient()
	log.Println("Bilibili client initialized.")
	// Repository 实现
	videoProgressRepo := persistence.NewGormVideoProgressRepository(db)
	log.Println("Video progress repository initialized.")

	// --- 初始化领域服务 ---
	watchTimeCalculator := service.NewWatchTimeCalculator()
	log.Println("Watch time calculator initialized.")

	// --- 初始化应用服务 ---
	// VideoProgressService
	videoProgressService := application.NewVideoProgressService(videoProgressRepo, biliClient)
	log.Println("Video progress service initialized.")
	// WatchTimeService
	watchTimeSvc := application.NewWatchTimeService(biliClient, watchTimeCalculator)
	log.Println("Watch time service initialized.")
	// (确保 watchTimeSvc 被使用，否则 linter 会报错。可以在这里或 API handler 中使用它)
	_ = watchTimeSvc // Placeholder to use watchTimeSvc, remove if used elsewhere

	// --- 设置 Gin 和路由 ---
	gin.SetMode(cfg.GinMode)
	router := gin.Default()
	setupRoutes(router, db) // 传入 db 用于潜在的数据库健康检查
	// TODO: 初始化并注册 Web handlers, 包括可能使用 watchTimeSvc 的 handler

	// --- 初始化并启动调度器 ---
	appScheduler := scheduler.NewScheduler() // NewScheduler no longer takes arguments

	// 定义获取视频进度的作业逻辑
	fetchVideoProgressJob := func() {
		// TODO: 将目标 AID 从配置加载，或从数据库查询需要追踪的视频列表
		targetAID := "114102919764678" // 示例 AID
		var targetCID int64 = 0        // 初始化 targetCID

		jobName := fmt.Sprintf("FetchVideoProgress(AID: %s)", targetAID)
		log.Printf("Cron job starting: %s", jobName)
		ctx := context.Background()

		// 1. 获取视频的第一个 CID
		videoView, err := biliClient.GetVideoView(ctx, targetAID, "") // bvid 为空，使用 aid
		if err != nil {
			log.Printf("Error fetching video view for AID %s in job '%s': %v. Skipping progress fetch.", targetAID, jobName, err)
			return
		}
		if videoView == nil || len(videoView.Pages) == 0 {
			log.Printf("No pages found for video AID %s in job '%s'. Skipping progress fetch.", targetAID, jobName)
			return
		}
		targetCID = videoView.Pages[0].Cid // 使用第一个分P的 CID
		log.Printf("Determined targetCID: %d for AID: %s", targetCID, targetAID)

		// 2. 获取并保存进度
		if err := videoProgressService.FetchAndSaveVideoProgress(ctx, targetAID, strconv.FormatInt(targetCID, 10)); err != nil {
			log.Printf("Error executing FetchAndSaveVideoProgress for AID %s, CID %d in job '%s': %v", targetAID, targetCID, jobName, err)
		}
		log.Printf("Cron job finished: %s (processed AID: %s, CID: %d)", jobName, targetAID, targetCID)
	}

	// 调度作业
	if err := appScheduler.ScheduleJob("FetchVideoProgress", cfg.Scheduler.Cron, fetchVideoProgressJob); err != nil {
		log.Fatalf("Failed to schedule job 'FetchVideoProgress': %v", err)
	}

	go appScheduler.Start() // 在单独的 goroutine 中启动调度器

	// --- 启动 Gin 服务器 ---
	serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	go func() {
		// 服务连接
		log.Printf("Starting server on %s", serverAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// --- 等待中断信号以优雅关闭服务器和调度器 ---
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
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

// setupRoutes 配置 Gin 路由。
func setupRoutes(router *gin.Engine, db *gorm.DB) {
	// 健康检查端点
	router.GET("/health", func(c *gin.Context) {
		// 简单的健康检查
		healthStatus := gin.H{"status": "UP"}

		// 检查数据库连接
		sqlDB, err := db.DB()
		if err != nil {
			healthStatus["db"] = "error getting DB instance"
			c.JSON(http.StatusInternalServerError, healthStatus)
			return
		}
		if err := sqlDB.Ping(); err != nil {
			healthStatus["db"] = "down"
			c.JSON(http.StatusServiceUnavailable, healthStatus)
			return
		}
		healthStatus["db"] = "up"

		c.JSON(http.StatusOK, healthStatus)
	})

	// TODO: 在此添加真实的 API 端点
	// api := router.Group("/api/v1")
	// videoProgressHandler.RegisterRoutes(api)
}
