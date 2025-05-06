# Bilibili Infrastructure

此目录包含与 Bilibili API 交互相关的基础设施代码。

## 主要职责

*   封装调用 Bilibili HTTP API 的底层细节。
*   处理 HTTP 请求的构建、发送和响应解析。
*   定义与特定 Bilibili API 端点对应的请求和响应数据结构。
*   实现应用层定义的 `VideoProgressFetcher` 接口，将具体的 API 调用结果转换为应用层所需的数据格式。

## 当前内容

*   `client.go`: 定义了底层的 Bilibili API 客户端。
    *   `Client` 结构体: 包含 `http.Client` 和基础配置。
    *   `VideoProgressResponse`, `VideoProgressData` 及嵌套结构体: 定义了 `/x/player/wbi/v2` 端点的响应 JSON 结构。
    *   `NewClient()`: 创建客户端实例。
    *   `GetVideoProgress(aid, cid string) (*VideoProgressResponse, error)`: 调用 Bilibili API 获取指定视频的观看进度，返回原始的 API 响应结构体。
*   `fetcher.go`: 实现了应用层的 `VideoProgressFetcher` 接口。
    *   `bilibiliFetcher` 结构体: 包装了 `Client`。
    *   `NewBilibiliFetcher(client *Client)`: 创建 fetcher 实例。
    *   `Fetch(ctx context.Context, aid, cid string) (*application.FetchedProgressData, error)`: 调用 `client.GetVideoProgress`，并将结果（如果有效）映射到应用层定义的 `application.FetchedProgressData` 结构。

## 注意

*   `Client` 负责处理原始的 HTTP 通信和 JSON 解析，保持与 Bilibili API 的一致性。
*   `Fetcher` 负责适配应用层的需求，将基础设施的细节（如复杂的 API 响应结构）与应用层隔离。
*   TODO: 实现动态加载 Cookie、处理 WBI 签名（如果需要）。 