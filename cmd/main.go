package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/krisxia0506/bilibili-watcher/internal/config"
	"github.com/krisxia0506/bilibili-watcher/internal/infrastructure/bilibili"
	"github.com/krisxia0506/bilibili-watcher/internal/infrastructure/persistence"
	_ "github.com/krisxia0506/bilibili-watcher/internal/infrastructure/persistence" // Import for side effects (init)
)

// main 程序入口
func main() {
	// 打印所有环境变量 (用于调试)
	log.Println("--- Environment Variables ---")
	for _, env := range os.Environ() {
		log.Println(env)
	}
	log.Println("---------------------------")
	
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	db, err := persistence.NewDatabaseConnection(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Successfully connected to the database.")
	_=db
	// Initialize Bilibili client
	// In the future, you might pass config (e.g., cookie) to NewClient
	biliClient := bilibili.NewClient()
	log.Println("Bilibili client initialized.", biliClient) // TODO: Remove biliClient print later

	// Set Gin mode
	gin.SetMode(cfg.GinMode)

	// Initialize Gin engine
	router := gin.Default()

	// --- Setup Repositories, Services, Handlers (Dependency Injection) ---
	// Example (will be expanded later):
	// videoProgressRepo := persistence.NewGormVideoProgressRepository(db)
	// videoProgressService := application.NewVideoProgressService(videoProgressRepo, biliClient)
	// videoProgressHandler := web.NewVideoProgressHandler(videoProgressService)

	// --- Setup Routes ---
	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		// TODO: Add more sophisticated health checks (e.g., DB connection)
		c.JSON(http.StatusOK, gin.H{
			"status": "UP",
		})
	})

	// TODO: Add API routes here, e.g.:
	// api := router.Group("/api/v1")
	// videoProgressHandler.RegisterRoutes(api)

	// --- Start Scheduler (if needed) ---
	// TODO: Initialize and start the cron scheduler here
	// scheduler := NewScheduler(cfg.Scheduler.Cron, videoProgressService)
	// go scheduler.Start()

	// --- Start Gin Server ---
	serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Starting server on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
