# Repository Interfaces

此目录定义仓库（Repository）的接口。仓库是领域层的一部分，用于抽象领域对象的数据持久化机制。

## 主要职责

*   定义用于添加、获取、更新和删除领域聚合根（Aggregate Roots）或实体的操作契约（接口）。
*   隐藏数据存储的底层实现细节（如 SQL 查询、ORM 调用）。
*   提供一种面向集合的方式来访问领域对象（例如，`FindAll()`, `GetByID()`, `Save()`）。
*   接口的方法通常接受并返回领域模型对象 ([domain/model](mdc:internal/domain/model/))。

## 当前内容

*   `video_progress.go`: 定义了视频观看进度仓库的接口。
    *   `VideoProgressRepository` 接口:
        *   `Save(ctx context.Context, progress *model.VideoProgress) error`: 保存一条进度记录。
        *   `GetLatestByAIDAndCID(ctx context.Context, aid, lastPlayCID int64) (*model.VideoProgress, error)`: 获取指定视频（通过 AID 和 LastPlayCID）的最新一条进度记录。如果找不到，返回 `nil, nil`。
        *   `ListByDateRange(ctx context.Context, start, end time.Time) ([]*model.VideoProgress, error)`: 获取指定日期范围内的所有进度记录。

## 注意

*   此目录只包含接口定义，具体的实现位于基础设施层 ([infrastructure/persistence](mdc:internal/infrastructure/persistence/))。
*   仓库接口的设计应反映领域的需求，而不是数据库表结构。 