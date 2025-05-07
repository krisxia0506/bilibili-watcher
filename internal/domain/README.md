# Domain Layer (`internal/domain`)

领域层是软件的核心，包含了业务逻辑和规则。它独立于其他层，不依赖于应用层、基础设施层或接口层。

## 主要职责

*   定义核心业务概念（实体 `Entities`、值对象 `Value Objects`）。
*   实现核心业务规则和逻辑。
*   定义领域事件 (`Domain Events`)（如果需要，目前未使用）。
*   定义仓储接口 (`Repository Interfaces`)，用于抽象数据持久化，由基础设施层实现。
*   定义领域服务 (`Domain Services`)，用于处理跨多个领域对象的复杂业务逻辑。
*   确保领域对象始终处于有效状态（通过聚合根或工厂方法）。

## 子目录

*   `model/`: 包含领域模型（实体 `Entities` 和值对象 `Value Objects`）。
    *   `video_progress.go`: 定义了 `VideoProgress` 实体，代表视频观看进度的核心信息。
    *   `video_page.go`: 定义了 `VideoPage` 值对象（或实体，取决于具体用法），表示视频分P信息，用于时长计算。
*   `repository/`: 定义仓储接口，用于抽象数据访问。
    *   `video_progress.go`: 定义了 `VideoProgressRepository` 接口，规定了视频进度数据的持久化和查询操作。
*   `service/`: 包含领域服务。
    *   `watch_time_calculator.go`: 定义了 `WatchTimeCalculator` 接口。
    *   `watch_time_calculator_impl.go`: 提供了 `WatchTimeCalculator` 的实现，封装了跨分P计算观看时长的复杂逻辑。

## 关键原则

*   **无外部依赖**: 领域层代码不应导入 `internal/application`, `internal/infrastructure` 或 `internal/interfaces` 包。
*   **封装**: 业务逻辑应封装在领域对象或领域服务内部。
*   **持久化抽象**: 通过仓储接口将数据持久化细节与领域逻辑分离。

## 注意

保持领域层的纯粹性至关重要。它应该只关注业务本身，与具体的技术实现解耦。 