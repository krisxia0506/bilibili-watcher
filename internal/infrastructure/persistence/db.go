package persistence

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/krisxia0506/bilibili-watcher/internal/config"
	"github.com/krisxia0506/bilibili-watcher/internal/domain/model"
)

// NewDatabaseConnection 创建新的 MySQL GORM 数据库连接。
func NewDatabaseConnection(cfg *config.DatabaseConfig) (*gorm.DB, error) {

	// 从配置字段构造 MySQL DSN
	if cfg.User == "" || cfg.Password == "" || cfg.Host == "" || cfg.DBName == "" {
		return nil, fmt.Errorf("mysql config incomplete: user, password, host, and dbname are required")
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)
	dialector := mysql.Open(dsn)

	// GORM 日志记录器配置
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别 (Silent, Error, Warn, Info)
			IgnoreRecordNotFoundError: true,        // 忽略 ErrRecordNotFound 错误
			Colorful:                  true,        // 启用彩色打印
		},
	)

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 自动迁移领域模型
	err = db.AutoMigrate(
		&model.VideoProgress{},
		// 如果需要，在此添加其他模型
	)
	if err != nil {
		return nil, fmt.Errorf("failed to auto migrate database schemas: %w", err)
	}

	log.Println("Database connection established and migrations completed.")
	return db, nil
}
