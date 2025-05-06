# Scheduler (`internal/infrastructure/scheduler`)

此目录负责实现和管理定时任务。

## 主要职责

*   提供一个调度器 (`Scheduler`) 来注册和运行基于 Cron 表达式的定时任务。
*   封装具体的定时任务库（本项目使用 `robfig/cron/v3`）。
*   按计划触发应用层的服务方法来执行业务逻辑。
*   管理调度器的生命周期（启动、停止）。

## 主要组件

*   `scheduler.go`:
    *   `Scheduler` 结构体: 包含 cron 实例 (`*cron.Cron`) 和应用层服务的依赖（当前为 `*application.VideoProgressService` 指针）。
    *   `NewScheduler`: 创建调度器实例，注入应用服务依赖。
    *   `RegisterJobs`: 根据传入的 Cron 表达式注册要执行的函数 (`runRecordProgressJob`)。
    *   `runRecordProgressJob`: 定时任务实际执行的函数。目前硬编码了目标视频的 AID 和 CID，并调用 `progressSvc.FetchAndSaveVideoProgress` 来执行获取和保存进度的应用层逻辑。
    *   `Start` / `Stop`: 控制 cron 调度器的启动和停止。

## 注意

*   `runRecordProgressJob` 中目前硬编码了目标视频，未来应改为从配置或数据库动态加载。
*   错误处理：`runRecordProgressJob` 会捕获并记录应用服务执行过程中的错误，但不会中断调度器的运行，以便任务可以在下一个周期重试。 