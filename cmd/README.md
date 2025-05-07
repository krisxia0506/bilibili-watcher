# Command Directory (`cmd`)

此目录包含项目的主要入口点。

## 主要组件

*   `main.go`: 后端服务的启动入口。负责：
    *   加载配置 (`internal/config`)。
    *   初始化数据库连接 (`internal/infrastructure/persistence`)。
    *   初始化基础设施组件（如 Bilibili 客户端）。
    *   初始化领域服务（如 `WatchTimeCalculator`）。
    *   初始化应用层服务（如 `VideoProgressService`, `VideoAnalyticsService`），并注入依赖。
    *   设置 Gin Web 服务器及路由（通过调用 `internal/interfaces/api/rest.SetupRouter`）。
    *   初始化并启动定时任务调度器 (`internal/infrastructure/scheduler`)，并注册具体的作业逻辑（如获取视频进度）。
    *   处理操作系统的中断信号以实现优雅停机。

## 运行

可以直接运行 `go run cmd/main.go` 来启动后端服务（需要配置好必要的环境变量）。更推荐的方式是使用 Docker Compose。 