# Persistence Infrastructure

此目录包含数据持久化相关的基础设施实现代码。

## 主要职责

*   建立和管理数据库连接（例如，使用 GORM）。
*   实现领域层定义的仓库接口 ([domain/repository](mdc:internal/domain/repository/))。
*   将领域模型对象 ([domain/model](mdc:internal/domain/model/)) 与数据库记录（例如，数据库表行）进行相互映射。
*   执行实际的数据库操作（CRUD - 创建、读取、更新、删除）。
*   处理数据库事务（如果应用层需要）。
*   封装特定的 ORM 或数据库驱动程序的细节。

## 当前内容

*   `db.go`:
    *   `NewDatabaseConnection(cfg *config.DatabaseConfig) (*gorm.DB, error)`: 根据配置创建并返回一个 GORM 数据库连接实例 (MySQL)。它还配置了 GORM 日志记录器和单数表名策略，并执行数据库自动迁移 (`AutoMigrate`)。
*   `video_progress_repository.go`:
    *   `gormVideoProgressRepository` 结构体: 实现了 `repository.VideoProgressRepository` 接口，包含一个 `*gorm.DB` 实例。
    *   `NewGormVideoProgressRepository(db *gorm.DB)`: 创建仓库实例。
    *   `Save(...)`: 使用 `db.Create` 保存 `model.VideoProgress`。
    *   `GetLatestByAIDAndCID(...)`: 使用 `db.Where(...).Order(...).First(...)` 查询最新记录。
    *   `ListByDateRange(...)`: 使用 `db.Where(...).Order(...).Find(...)` 查询范围内的记录。

## 注意

*   仓库的实现应忠于领域层定义的接口契约。
*   避免在仓库实现中包含业务逻辑；它只应负责数据映射和存储操作。
*   `AutoMigrate` 功能方便开发，但在生产环境中通常建议使用更专业的数据库迁移工具（如 migrate, flyway 等）配合 SQL 文件 ([sql/schema.sql](mdc:sql/schema.sql)) 进行更可控的数据库模式管理。 