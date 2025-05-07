# Interfaces Layer (`internal/interfaces`)

此层负责处理与外部世界的交互，充当应用程序核心（应用层和领域层）与外部输入/输出（如用户界面、API 客户端、其他系统）之间的适配器。

## 主要职责

*   **适配外部输入**: 将来自外部的请求（例如 HTTP 请求、CLI 命令、消息队列消息）转换为对应用层服务的调用。
*   **适配外部输出**: 将应用层返回的结果或领域层的数据转换为适合外部使用的格式（例如 JSON 响应、HTML 页面、命令行输出）。
*   **隔离核心**: 保护应用层和领域层不受外部技术细节变化的影响。

## 子目录

*   `api/rest/`: 包含处理 RESTful API 请求相关的代码。
    *   `router.go`: 设置 Gin 路由引擎，定义 API 路由分组（如 `/api/v1`），注册健康检查端点 (`/healthz`)，并将路由映射到相应的处理程序 (Handlers)。
    *   `dto/`: 定义用于 API 请求体和响应体的 数据传输对象 (Data Transfer Objects)。这些 DTOs 独立于领域模型，仅用于接口层的数据交换。
        *   `video_analytics_dto.go`: 定义了获取观看分段接口的请求 (`GetWatchedSegmentsRequest`) 和响应 (`WatchedSegment`, `GetWatchedSegmentsResponse`) 结构。
    *   `video_analytics_handler.go`: 实现了处理视频分析相关 API 请求的 Handler (`VideoAnalyticsHandler`)。它负责解析请求、调用应用层服务 (`VideoAnalyticsService`) 并使用统一响应包 (`pkg/response`) 返回结果。
*   (未来可能添加)
    *   `cli/`: 命令行界面处理程序。
    *   `grpc/`: gRPC 服务定义和处理程序。
    *   `mq/`: 消息队列消费者/生产者适配器。

## 关键原则

*   **依赖方向**: 接口层依赖于应用层（调用应用服务）。它**不应该**直接依赖领域层或基础设施层（除了可能需要共享的 DTO 或基础类型）。
*   **适配和转换**: 主要工作是数据格式的适配和转换，以及将外部请求路由到正确的应用服务。
*   **简洁性**: 处理程序 (Handler) 应该保持简单，只做请求解析、调用应用服务、响应格式化的工作，不包含业务逻辑。 