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

// NewDatabaseConnection creates a new GORM database connection for MySQL.
// 创建新的 MySQL GORM 数据库连接。
func NewDatabaseConnection(cfg *config.DatabaseConfig) (*gorm.DB, error) {

	// Construct MySQL DSN from config fields
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

	// GORM logger configuration
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level (Silent, Error, Warn, Info)
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto migrate domain models
	err = db.AutoMigrate(
		&model.VideoProgress{},
		// Add other models here if needed
	)
	if err != nil {
		return nil, fmt.Errorf("failed to auto migrate database schemas: %w", err)
	}

	log.Println("Database connection established and migrations completed.")
	return db, nil
}
