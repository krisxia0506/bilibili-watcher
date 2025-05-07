# Bilibili API Client (`internal/infrastructure/bilibili`)

此目录包含与 Bilibili API 进行交互的具体实现。

## 主要职责

*   提供一个通用的 HTTP 客户端 (`Client`) 来处理对 Bilibili API 的请求。
*   封装 API 请求的细节，如构建 URL、设置必要的 Header (包括 Cookie `SESSDATA`)、发送请求、处理 HTTP 状态码和基础错误。
*   实现 `application.BilibiliClient` 接口，提供应用层所需的操作（如获取视频进度、获取视频信息）。
*   定义 Bilibili API 响应的具体数据结构。
*   将从 Bilibili API 获取的原始响应数据映射到应用层定义的 DTO。

## 主要组件

*   `client.go`: 定义了 `Client` 结构体和通用的 `Get` 方法。`Client` 负责维护 `http.Client`、基础 URL 和认证信息（如 `SESSDATA`）。`Get` 方法处理通用的 GET 请求逻辑。
*   `video_progress.go`: 包含 `GetVideoProgress` 方法的实现（作为 `*Client` 的方法）。此方法调用 `/x/player/wbi/v2` API，解析响应，并将其映射到 `application.VideoProgressDTO`。
*   `video_view.go`: 包含 `GetVideoView` 方法的实现（作为 `*Client` 的方法）。此方法调用 `/x/web-interface/view` API，解析响应，并将其映射到 `application.VideoViewDTO`。

## 注意

*   `GetVideoProgress` 和 `GetVideoView` 共同为 `*Client` 类型添加了方法，使其完整实现了 `application.BilibiliClient` 接口。
*   错误处理：底层 `Get` 方法处理 HTTP 和解码错误，各个具体的 API 方法 (`GetVideoProgress`, `GetVideoView`) 处理 Bilibili 返回的业务错误码 (`code != 0`)，并将基础设施的响应映射到应用层 DTO 或错误。 