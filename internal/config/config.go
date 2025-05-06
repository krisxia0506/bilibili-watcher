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
	UID    string // Env: BILIBILI_UID
	Cookie string // Env: BILIBILI_COOKIE
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
		// 允许空密码吗？返回错误可能更安全。
		return nil, fmt.Errorf("required environment variable DATABASE_PASSWORD is not set")
	}
	cfg.Database.DBName = os.Getenv("DATABASE_DBNAME")
	if cfg.Database.DBName == "" {
		return nil, fmt.Errorf("required environment variable DATABASE_DBNAME is not set")
	}

	// --- Bilibili 配置 ---
	cfg.Bilibili.UID = os.Getenv("BILIBILI_UID")
	cfg.Bilibili.Cookie = os.Getenv("BILIBILI_COOKIE")
	// 如果 UID 和 Cookie 是必需的，添加验证
	// if cfg.Bilibili.UID == "" {
	// 	return nil, fmt.Errorf("required environment variable BILIBILI_UID is not set")
	// }
	// if cfg.Bilibili.Cookie == "" {
	// 	return nil, fmt.Errorf("required environment variable BILIBILI_COOKIE is not set")
	// }

	// --- 定时任务配置 ---
	cfg.Scheduler.Cron = getEnv("SCHEDULER_CRON", "0 0 * * *")

	// --- Gin 模式 ---
	cfg.GinMode = getEnv("GIN_MODE", "debug")

	return cfg, nil
}

// getEnv 获取环境变量，如果未设置则返回默认值。
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
