# Command (cmd)

此目录包含项目应用程序的可执行文件入口点。

## 主要职责

*   作为应用程序的启动层。
*   负责解析命令行参数（如果需要）。
*   初始化配置加载。
*   设置日志。
*   进行依赖注入：组装应用程序的不同层（领域层、应用层、基础设施层）和组件（服务、仓库、控制器等）。
*   启动应用程序的主服务（例如，HTTP 服务器、后台任务调度器）。
*   处理操作系统的信号以实现优雅停机。

## 当前内容

*   `main.go`: Web 服务器应用程序的主要入口点。它执行以下操作：
    *   加载配置 ([config.LoadConfig](mdc:internal/config/config.go))
    *   初始化数据库连接 ([persistence.NewDatabaseConnection](mdc:internal/infrastructure/persistence/db.go))
    *   初始化 Bilibili API 客户端 ([bilibili.NewClient](mdc:internal/infrastructure/bilibili/client.go)) 和 Fetcher ([bilibili.NewBilibiliFetcher](mdc:internal/infrastructure/bilibili/fetcher.go))
    *   初始化 GORM 仓库 ([persistence.NewGormVideoProgressRepository](mdc:internal/infrastructure/persistence/video_progress_repository.go))
    *   初始化应用服务 ([application.NewVideoProgressService](mdc:internal/application/video_progress_service.go))
    *   初始化并启动定时任务调度器 ([scheduler.NewScheduler](mdc:internal/infrastructure/scheduler/scheduler.go))
    *   初始化 Gin Web 引擎并设置路由（目前仅有 /health）。
    *   启动 HTTP 服务器并监听端口。
    *   监听系统信号以实现优雅关闭。

## 注意

此目录下的代码应尽量保持简洁，主要关注应用程序的组装和启动，避免包含具体的业务逻辑或基础设施细节。 