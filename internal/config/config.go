package config

import (
	"fmt"
	"os"
	"strconv"
	// "strings"
	// "github.com/spf13/viper" // Removed Viper dependency
)

// Config 保存应用程序配置。
// (从环境变量加载)
type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	Bilibili  BilibiliConfig
	Scheduler SchedulerConfig
	GinMode   string
}

// ServerConfig 保存服务器相关配置。
type ServerConfig struct {
	Port int // Env: SERVER_PORT (默认: 8080)
}

// DatabaseConfig 保存数据库相关配置。
// (仅 MySQL)
type DatabaseConfig struct {
	Host     string // Env: DATABASE_HOST
	Port     int    // Env: DATABASE_PORT (默认: 3306)
	User     string // Env: DATABASE_USER
	Password string // Env: DATABASE_PASSWORD
	DBName   string // Env: DATABASE_DBNAME
}

// BilibiliConfig 保存 Bilibili API 相关配置。
type BilibiliConfig struct {
	SessData   string // Env: SESSDATA (Bilibili SESSDATA Cookie, 必需)
	TargetBVID string // Env: WATCH_TARGET_BVID (用于定时任务的目标 BVID, 必需)
}

// SchedulerConfig 保存定时任务相关配置。
type SchedulerConfig struct {
	Cron string // Env: SCHEDULER_CRON (默认: "0 0 * * *")
}

// LoadConfig 使用 os 包严格从环境变量加载配置。
func LoadConfig() (*Config, error) {
	cfg := &Config{}
	var err error

	// --- 服务器配置 ---
	serverPortStr := getEnv("SERVER_PORT", "8080")
	cfg.Server.Port, err = strconv.Atoi(serverPortStr)
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_PORT value %q: %w", serverPortStr, err)
	}

	// --- 数据库配置 ---
	cfg.Database.Host = getEnvOrErr("DATABASE_HOST")
	cfg.Database.Port, err = strconv.Atoi(getEnv("DATABASE_PORT", "3306"))
	if err != nil {
		return nil, fmt.Errorf("invalid DATABASE_PORT value: %w", err)
	}
	cfg.Database.User = getEnvOrErr("DATABASE_USER")
	cfg.Database.Password = getEnvOrErr("DATABASE_PASSWORD")
	cfg.Database.DBName = getEnvOrErr("DATABASE_DBNAME")

	// --- Bilibili 配置 ---
	cfg.Bilibili.SessData = getEnvOrErr("BILIBILI_SESSDATA")
	cfg.Bilibili.TargetBVID = getEnvOrErr("BILIBILI_BVID")

	// --- 定时任务配置 ---
	cfg.Scheduler.Cron = getEnv("SCHEDULER_CRON", "0 0 * * *")

	// --- Gin 模式 ---
	cfg.GinMode = getEnv("GIN_MODE", "debug")

	// --- 检查必需的环境变量 ---
	if cfg.Database.Host == "" {
		return nil, fmt.Errorf("required environment variable DATABASE_HOST is not set")
	}
	if cfg.Database.User == "" {
		return nil, fmt.Errorf("required environment variable DATABASE_USER is not set")
	}
	// 允许空密码吗？对于本地开发可能允许，但生产通常不。
	// if cfg.Database.Password == "" {
	// 	return nil, fmt.Errorf("required environment variable DATABASE_PASSWORD is not set")
	// }
	if cfg.Database.DBName == "" {
		return nil, fmt.Errorf("required environment variable DATABASE_DBNAME is not set")
	}
	if cfg.Bilibili.SessData == "" {
		return nil, fmt.Errorf("required environment variable SESSDATA is not set")
	}
	if cfg.Bilibili.TargetBVID == "" {
		return nil, fmt.Errorf("required environment variable WATCH_TARGET_BVID is not set")
	}

	return cfg, nil
}

// getEnv 获取环境变量，如果未设置则返回默认值。
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvOrErr 获取环境变量，如果未设置则返回空字符串。
func getEnvOrErr(key string) string {
	return os.Getenv(key)
}
