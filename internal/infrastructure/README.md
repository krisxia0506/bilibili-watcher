# Infrastructure Layer (`internal/infrastructure`)

基础设施层负责实现应用层和领域层定义的接口，提供与外部系统和技术细节交互的能力。

## 主要职责

*   实现领域层定义的仓库接口（例如，使用 GORM 与数据库交互）。
*   实现应用层定义的外部服务接口（例如，调用 Bilibili API）。
*   处理具体的技术细节，如数据库连接、HTTP 请求、文件系统访问、消息队列交互、定时任务调度等。
*   将基础设施特定的数据结构（如数据库行、API 响应）与领域模型或应用层 DTO 进行转换。

## 子目录

*   `bilibili/`: 包含与 Bilibili API 交互的具体实现。
    *   `client.go`: 提供了通用的 Bilibili HTTP 客户端，处理请求发送、认证、基础错误处理。
    *   `video_progress.go`: 实现了 `GetVideoProgress` 方法，调用 Bilibili API 获取视频进度，并将响应映射到应用层 DTO。
    *   `video_view.go`: 实现了 `GetVideoView` 方法，调用 Bilibili API 获取视频详情，并将响应映射到应用层 DTO。
*   `persistence/`: 包含数据持久化的具体实现。
    *   `db.go`: 负责初始化和管理数据库连接（GORM）。
    *   `video_progress_repository.go`: 实现了 `VideoProgressRepository` 接口，使用 GORM 将 `VideoProgress` 领域模型持久化到数据库。
*   `scheduler/`: 包含定时任务调度器的实现。
    *   `scheduler.go`: 使用 `robfig/cron` 库实现定时任务调度，负责按计划触发应用层服务。
*   `web/`: (可选，如果 Gin 相关代码放在这里) 包含与 Web 框架（如 Gin）适配的代码，例如路由设置、请求/响应处理的适配器等。（当前项目的 Gin 设置主要在 `cmd/main.go` 和潜在的 handlers 中）

## 关键原则

*   **依赖倒置**: 实现了上层（应用层、领域层）定义的接口。
*   **技术细节封装**: 将具体的技术实现细节封装在此层，使上层保持稳定。
*   **数据映射**: 负责在基础设施数据格式和上层所需的数据格式（领域模型、DTO）之间进行转换。

## 注意

*   基础设施层依赖于应用层和领域层（主要是实现它们定义的接口）。
*   领域逻辑**不应该**泄漏到基础设施层。
*   应尽量将基础设施层的具体实现细节与上层解耦。 