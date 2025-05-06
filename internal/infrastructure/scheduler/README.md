# Scheduler Infrastructure

此目录包含定时任务调度器的基础设施实现。

## 主要职责

*   提供定时执行任务的能力。
*   封装具体的 cron 库（例如 `github.com/robfig/cron/v3`）的实现细节。
*   负责注册、启动和停止定时任务。
*   调用应用层的服务来执行实际的业务逻辑。

## 当前内容

*   `scheduler.go`:
    *   `Scheduler` 结构体: 管理 `cron.Cron` 实例和对应用服务的引用。
    *   `NewScheduler(progressSvc application.VideoProgressService)`: 创建调度器实例。
    *   `RegisterJobs(schedule string) error`: 根据给定的 cron 表达式注册一个任务，该任务会调用 `runRecordProgressJob`。
    *   `runRecordProgressJob()`: 定时任务执行的函数，它创建一个 `context.Background()` 并调用 `progressSvc.RecordProgressForTargetVideo`。
    *   `Start()`: 启动 cron 调度器。
    *   `Stop()`: 优雅地停止 cron 调度器，并返回一个 context 用于等待任务完成。

## 注意

*   调度器本身不包含业务逻辑，它只是一个触发器，负责按时调用应用层服务的方法。
*   错误处理：`runRecordProgressJob` 目前只记录应用服务返回的错误，然后任务结束。如果需要更复杂的重试或错误通知机制，可以在这里扩展。 