# Domain Layer (`internal/domain`)

领域层是软件的核心，包含了业务逻辑和规则。它独立于其他层，不依赖于应用层、基础设施层或接口层。

## 主要职责

*   定义核心业务概念（实体 `Entities`、值对象 `Value Objects`）。
*   实现核心业务规则和逻辑。
*   定义领域事件 (`Domain Events`)（如果需要）。
*   定义仓储接口 (`Repository Interfaces`)，用于抽象数据持久化，由基础设施层实现。
*   定义领域服务 (`Domain Services`)，用于处理跨多个领域对象的复杂业务逻辑。
*   确保领域对象始终处于有效状态（通过聚合根或工厂方法）。

## 子目录

*   `model/`: 包含领域模型（实体 `Entities` 和值对象 `Value Objects`）。
    *   `video_progress.go`: 定义了 `VideoProgress` 实体，代表视频观看进度的核心信息。
*   `repository/`: 定义仓储接口，用于抽象数据访问。
    *   `video_progress.go`: 定义了 `VideoProgressRepository` 接口，规定了视频进度数据的持久化操作（如 `Save`, `FindByAID`）。
*   `service/`: (可选) 包含领域服务。如果某些业务逻辑不适合放在任何一个实体或值对象中，可以放在领域服务里。目前此项目尚未用到。

## 关键原则

*   **无外部依赖**: 领域层代码不应导入 `internal/application`, `internal/infrastructure` 或 `internal/interface` (或 `web`) 包。
*   **封装**: 业务逻辑应封装在领域对象或领域服务内部。
*   **持久化抽象**: 通过仓储接口将数据持久化细节与领域逻辑分离。

## 注意

保持领域层的纯粹性至关重要。它应该只关注业务本身，与具体的技术实现解耦。 