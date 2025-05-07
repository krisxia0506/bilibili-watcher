# REST API Layer (`internal/interfaces/api/rest`)

此目录包含处理 RESTful API 请求的具体实现。

## 主要职责

*   设置 Gin Web 路由引擎。
*   定义 API 路由、路由组 (如 `/api/v1`) 和中间件。
*   实现具体的 HTTP 请求处理程序 (Handlers)。
*   将 HTTP 请求（路径参数、查询参数、请求体）解析并绑定到 DTO 或 Go 类型。
*   调用相应的应用层服务来执行业务逻辑。
*   使用统一的响应包 (`pkg/response`) 将应用层的返回结果或错误格式化为 JSON 响应。
*   定义用于 API 请求/响应的 DTO。

## 子目录和文件

*   `router.go`: 包含 `SetupRouter` 函数，负责初始化 Gin 引擎、设置路由分组、注册健康检查端点 (`/healthz`) 和将路由委托给具体的 Handlers。
*   `dto/`: 存放 API 请求和响应的 DTO。
    *   `video_analytics_dto.go`: 定义了 `/video/watch-segments` 端点的请求和响应结构。
*   `video_analytics_handler.go`: 包含 `VideoAnalyticsHandler` 的实现。
    *   `NewVideoAnalyticsHandler`: 创建 Handler 实例，注入应用服务依赖。
    *   `RegisterRoutes`: 在 Gin 路由组上注册此 Handler 负责的路由。
    *   `GetWatchedSegments`: 处理 `POST /api/v1/video/watch-segments` 请求，解析请求体，调用 `VideoAnalyticsService`，并返回分段观看时长结果。
*   (未来可能添加更多 handler 文件，如 `user_handler.go` 等)

## 注意

*   Handlers 应保持轻量，主要负责请求/响应处理和调用应用服务，不应包含复杂业务逻辑。
*   使用 `pkg/response` 统一 API 响应格式。
*   DTOs 用于接口层与外部的数据交换，与领域模型解耦。 