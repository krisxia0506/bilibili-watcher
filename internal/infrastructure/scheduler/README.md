# Scheduler (`internal/infrastructure/scheduler`)

此目录负责实现和管理通用的定时任务调度。

## 主要职责

*   提供一个通用的调度器 (`Scheduler`) 来注册和运行基于 Cron 表达式的定时任务。
*   封装具体的定时任务库（本项目使用 `robfig/cron/v3`）。
*   按计划执行由调用方（如 `cmd/main.go`）提供的具体作业函数。
*   管理调度器的生命周期（启动、停止）。

## 主要组件

*   `scheduler.go`:
    *   `Scheduler` 结构体: 包含 cron 实例 (`*cron.Cron`)。
    *   `NewScheduler`: 创建调度器实例。
    *   `ScheduleJob`: 允许注册一个带有名称、Cron 表达式和无参数作业函数的定时任务。
    *   `Start` / `Stop`: 控制 cron 调度器的启动和停止。

## 注意

*   调度器本身不包含任何业务逻辑，它只是一个触发器。
*   具体的作业逻辑（例如调用哪个应用服务）由 `cmd/main.go` 在注册作业时通过函数（通常是闭包）提供。
*   错误处理：传递给 `ScheduleJob` 的作业函数需要自行处理其内部逻辑可能产生的错误（例如记录日志），调度器本身只处理注册错误。 