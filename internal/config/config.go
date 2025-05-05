package config

import (
	"fmt"
	"os"
	"strconv"
	// "strings"
	// "github.com/spf13/viper" // Removed Viper dependency
)

// Config holds the application configuration.
// 应用配置结构体 (从环境变量加载)
type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	Bilibili  BilibiliConfig
	Scheduler SchedulerConfig
	GinMode   string
}

// ServerConfig holds server related configuration.
// 服务器相关配置
type ServerConfig struct {
	Port int // Env: SERVER_PORT (Default: 8080)
}

// DatabaseConfig holds database related configuration.
// 数据库相关配置 (仅 MySQL)
type DatabaseConfig struct {
	Host     string // Env: DATABASE_HOST
	Port     int    // Env: DATABASE_PORT (Default: 3306)
	User     string // Env: DATABASE_USER
	Password string // Env: DATABASE_PASSWORD
	DBName   string // Env: DATABASE_DBNAME
}

// BilibiliConfig holds Bilibili API related configuration.
// Bilibili API 相关配置
type BilibiliConfig struct {
	UID    string // Env: BILIBILI_UID
	Cookie string // Env: BILIBILI_COOKIE
}

// SchedulerConfig holds scheduler related configuration.
// 定时任务相关配置
type SchedulerConfig struct {
	Cron string // Env: SCHEDULER_CRON (Default: "0 0 * * *")
}

// LoadConfig loads configuration strictly from environment variables using os package.
// 严格从环境变量加载配置 (使用 os 包)
func LoadConfig() (*Config, error) {
	cfg := &Config{}
	var err error

	// --- Server Config ---
	serverPortStr := getEnv("BACKEND_PORT", "8080")
	cfg.Server.Port, err = strconv.Atoi(serverPortStr)
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_PORT value %q: %w", serverPortStr, err)
	}

	// --- Database Config ---
	cfg.Database.Host = os.Getenv("DATABASE_HOST")
	if cfg.Database.Host == "" {
		return nil, fmt.Errorf("required environment variable DATABASE_HOST is not set")
	}
	dbPortStr := getEnv("DATABASE_PORT", "3306")
	cfg.Database.Port, err = strconv.Atoi(dbPortStr)
	if err != nil {
		return nil, fmt.Errorf("invalid DATABASE_PORT value %q: %w", dbPortStr, err)
	}
	cfg.Database.User = os.Getenv("DATABASE_USER")
	if cfg.Database.User == "" {
		return nil, fmt.Errorf("required environment variable DATABASE_USER is not set")
	}
	cfg.Database.Password = os.Getenv("DATABASE_PASSWORD")
	if cfg.Database.Password == "" {
		// Allow empty password? Maybe return error is safer.
		return nil, fmt.Errorf("required environment variable DATABASE_PASSWORD is not set")
	}
	cfg.Database.DBName = os.Getenv("DATABASE_DBNAME")
	if cfg.Database.DBName == "" {
		return nil, fmt.Errorf("required environment variable DATABASE_DBNAME is not set")
	}

	// --- Bilibili Config ---
	cfg.Bilibili.UID = os.Getenv("BILIBILI_UID")
	cfg.Bilibili.Cookie = os.Getenv("BILIBILI_COOKIE")
	// Add validation if UID and Cookie are required
	// if cfg.Bilibili.UID == "" {
	// 	return nil, fmt.Errorf("required environment variable BILIBILI_UID is not set")
	// }
	// if cfg.Bilibili.Cookie == "" {
	// 	return nil, fmt.Errorf("required environment variable BILIBILI_COOKIE is not set")
	// }

	// --- Scheduler Config ---
	cfg.Scheduler.Cron = getEnv("SCHEDULER_CRON", "0 0 * * *")

	// --- Gin Mode ---
	cfg.GinMode = getEnv("GIN_MODE", "debug")

	return cfg, nil
}

// getEnv retrieves an environment variable or returns a default value if not set.
// 获取环境变量，如果未设置则返回默认值
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
