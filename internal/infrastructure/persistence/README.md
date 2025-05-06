# Persistence Layer (`internal/infrastructure/persistence`)

此目录负责实现数据的持久化和检索，通常与数据库进行交互。

## 主要职责

*   实现领域层定义的仓储接口 (`domain/repository/*`)。
*   封装与特定数据库（本项目中为 MySQL）和 ORM（本项目中为 GORM）交互的细节。
*   管理数据库连接的建立和配置。
*   定义 GORM 模型（通常与领域模型结构类似，但可能包含 GORM 特定标签或字段）。
*   在 GORM 模型和领域模型之间进行数据映射（转换）。

## 主要组件

*   `db.go`: 提供 `NewDatabaseConnection` 函数，用于根据配置建立和返回 GORM 数据库连接 (`*gorm.DB`)。
*   `video_progress_repository.go`: 实现了 `domain/repository.VideoProgressRepository` 接口。
    *   `gormVideoProgressRepository` 结构体: 包含 `*gorm.DB` 连接。
    *   `videoProgressGorm` 结构体: 定义了与 `video_progress` 表对应的 GORM 模型。
    *   `toDomain` / `fromDomain` 函数: 负责在 `videoProgressGorm` 和 `domain/model.VideoProgress` 之间进行转换。
    *   `NewGormVideoProgressRepository`: 创建仓库实例。
    *   `Save`, `GetLatestByAIDAndCID`, `ListByDateRange`, `FindByAID`: 实现了仓库接口定义的方法，执行具体的 GORM 数据库操作（Create, First, Find）。

## 关键原则

*   **接口实现**: 主要目的是实现领域层定义的持久化接口。
*   **技术细节隔离**: 将 SQL 查询、ORM 操作等数据库相关的具体实现细节限制在此层。
*   **数据映射**: 负责处理数据库记录与领域对象之间的转换。

## 注意

*   仓库的实现应忠于领域层定义的接口契约。
*   避免在仓库实现中包含业务逻辑；它只应负责数据映射和存储操作。
*   `AutoMigrate` 功能方便开发，但在生产环境中通常建议使用更专业的数据库迁移工具（如 migrate, flyway 等）配合 SQL 文件 ([sql/schema.sql](mdc:sql/schema.sql)) 进行更可控的数据库模式管理。 