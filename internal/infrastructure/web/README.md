# Web Infrastructure

此目录包含与 Web 框架（例如 Gin）适配相关的基础设施代码。

## 主要职责

*   定义 HTTP 请求处理器 (Handlers/Controllers)。
*   设置和配置 Web 服务器路由。
*   处理 HTTP 请求的解析（例如，路径参数、查询参数、请求体）。
*   调用应用层服务来处理业务逻辑。
*   将应用层的返回结果（或错误）转换为 HTTP 响应（例如，JSON 响应、状态码）。
*   处理 Web 特定的关注点，如请求验证（可以使用库如 validator）、认证中间件、CORS 配置等。
*   定义数据传输对象 (DTOs)，用于在 Web 层和应用层之间传递数据，避免直接暴露领域模型。

## 当前内容

*   (目前为空)

## 未来可能的用途

*   `handlers/video_progress_handler.go`: 定义处理视频进度相关 API 请求的 Handler。
*   `router.go`: 定义和注册所有 API 路由。
*   `middleware/`: 存放认证、日志记录等中间件。
*   `dto/`: 定义用于 API 请求和响应的 DTO 结构体。

## 注意

*   Handlers 应保持"薄"，主要负责请求/响应的转换和调用应用服务，避免包含业务逻辑。
*   使用 DTO 来隔离 Web 层和应用/领域层。 