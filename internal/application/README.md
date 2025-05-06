# Application Layer (`internal/application`)

应用层负责编排领域逻辑和基础设施交互，以完成具体的业务用例 (Use Cases)。它作为领域层和外部世界（如 Web 接口、定时任务）之间的桥梁。

## 主要职责

*   定义应用服务 (`*_service.go`)，封装特定的业务流程。
*   定义与外部基础设施交互的接口（例如 `BilibiliClient`），由基础设施层实现（依赖倒置）。
*   定义数据传输对象 (DTOs) (`*_dto.go`)，用于应用层与外部或基础设施层之间的数据交换，避免暴露领域模型。
*   不包含核心领域逻辑，而是委托给领域服务或领域模型。
*   不直接依赖基础设施的具体实现，只依赖接口。

## 主要组件

*   `bilibili_client.go`: 定义了与 Bilibili API 交互的应用层接口 (`BilibiliClient`)。
*   `bilibili_dto.go`: 定义了用于 Bilibili API 交互的 DTO (`VideoProgressDTO`, `VideoViewDTO`)。
*   `video_progress_service.go`: 实现了视频进度相关的应用服务 (`VideoProgressService`)，负责编排获取 Bilibili 视频进度、转换数据和保存到仓库的流程。

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