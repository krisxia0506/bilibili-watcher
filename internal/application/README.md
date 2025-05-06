# Application Layer

应用层负责编排特定的应用程序用例（Use Cases）。它作为领域层和基础设施层之间的粘合剂。

## 主要职责

*   定义应用程序的具体用例（例如，"记录视频观看进度"、"获取每日观看时长统计"）。
*   协调领域对象（实体、值对象）和领域服务来完成业务流程。
*   调用领域仓库接口来获取或持久化领域对象。
*   **定义**与外部基础设施交互的接口（例如 `VideoProgressFetcher`），但不包含具体实现。
*   处理事务管理（如果需要跨多个仓库操作）。
*   执行授权检查（例如，检查用户是否有权限执行某个操作）。
*   将领域层的错误转换为应用层特定的错误或结果。
*   **不**包含核心业务规则（那属于领域层）。
*   **不**包含基础设施的具体实现细节（例如，直接的数据库查询、HTTP 请求代码）。

## 当前内容

*   `fetcher.go`: 定义了从外部源获取视频进度的接口。
    *   `VideoProgressFetcher` 接口:
        *   `Fetch(ctx context.Context, aid, cid string) (*FetchedProgressData, error)`: 获取指定视频的最新进度。返回包含 BVID 和进度毫秒数的 `FetchedProgressData` 结构体指针，或在失败时返回错误。如果找不到有效进度（例如进度为0或负数），可能返回 `nil, nil`。
    *   `FetchedProgressData` 结构体: 包含从外部获取的基本数据 (`BVID`, `LastPlayTime`, `AID`, `LastPlayCid`)。
*   `video_progress_service.go`: 实现了管理视频进度的应用服务。
    *   `VideoProgressService` 接口:
        *   `RecordProgressForTargetVideo(ctx context.Context) error`: 获取并记录预定义目标视频的进度。
    *   `videoProgressService` 结构体: 接口的具体实现。
        *   依赖于 `repository.VideoProgressRepository` 和 `application.VideoProgressFetcher` 接口。
        *   `NewVideoProgressService(...)`: 创建服务实例，注入依赖。
        *   `RecordProgressForTargetVideo(...)`: 实现核心用例逻辑：调用 `fetcher` 获取数据，构建 `model.VideoProgress` 领域对象，然后调用 `repository` 保存。

## 注意

应用服务应该是相对较薄的一层，主要负责协调和委托，将复杂的业务规则交给领域层处理，将具体的技术实现交给基础设施层。 