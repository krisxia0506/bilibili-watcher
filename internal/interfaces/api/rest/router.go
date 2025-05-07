package rest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/krisxia0506/bilibili-watcher/internal/application"
	"github.com/krisxia0506/bilibili-watcher/pkg/response" // 引入统一响应包
	// Import application services and other dependencies handlers need
	// "github.com/krisxia0506/bilibili-watcher/internal/application"
)

// SetupRouter 配置并返回 Gin 引擎实例。
// 需要传入所有依赖项，以便创建和注册 handlers。
func SetupRouter(
	db *gorm.DB,
	ginMode string,
	videoAnalyticsService application.VideoAnalyticsService,
	// ... 其他需要的服务
) *gin.Engine {
	gin.SetMode(ginMode)
	router := gin.Default()

	// 健康检查路由
	router.GET("/healthz", func(c *gin.Context) {
		healthStatus := gin.H{
			"status": "UP",
			"db":     "unknown",
		}
		httpCode := http.StatusOK

		sqlDB, err := db.DB()
		if err != nil {
			healthStatus["db"] = "error getting DB instance"
			httpCode = http.StatusInternalServerError
		} else {
			if err := sqlDB.Ping(); err != nil {
				healthStatus["db"] = "down"
				httpCode = http.StatusServiceUnavailable
			} else {
				healthStatus["db"] = "up"
			}
		}
		// 使用统一响应体返回健康状态 (即使是健康检查)
		if httpCode == http.StatusOK {
			response.Success(c, healthStatus)
		} else {
			// 使用 ErrorWithData 来包含健康状态详情
			response.ErrorWithData(c, httpCode, response.CodeInternalError, fmt.Sprintf("Health check failed: %v", healthStatus["db"]), healthStatus)
		}
	})

	// API v1 分组
	apiV1 := router.Group("/api/v1")
	{
		// 初始化并注册 Video Analytics Handler
		videoAnalyticsHandler := NewVideoAnalyticsHandler(videoAnalyticsService)
		videoAnalyticsHandler.RegisterRoutes(apiV1)

		// 注册其他 handlers...
		apiV1.GET("/ping", func(c *gin.Context) {
			response.Success(c, "pong")
		})
	}

	return router
}
