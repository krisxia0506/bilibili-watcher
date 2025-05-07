# Configuration (`internal/config`)

此目录负责应用程序的配置加载和管理。

## 主要职责

*   定义配置结构体 (`Config` 及其子结构体如 `DatabaseConfig`, `BilibiliConfig` 等)。
*   从环境变量加载配置值。
*   提供默认值并在必要时进行类型转换（如端口号）。
*   验证必需的环境变量是否存在。

## 主要组件

*   `config.go`:
    *   定义了 `Config` 及其相关子结构体。
    *   `LoadConfig()`: 从环境变量读取配置，执行验证，并返回填充好的 `*Config` 实例或错误。
    *   `getEnv()`, `getEnvOrErr()`: 用于读取环境变量的辅助函数。

## 配置方式

完全通过环境变量进行配置，遵循十二因子应用 (Twelve-Factor App) 的原则。必需的环境变量包括：

*   `DATABASE_HOST`
*   `DATABASE_PORT` (默认 3306)
*   `DATABASE_USER`
*   `DATABASE_PASSWORD` (需要设置，但允许为空)
*   `DATABASE_DBNAME`
*   `BILIBILI_SESSDATA` (Bilibili Cookie)
*   `BILIBILI_BVID` (定时任务追踪的 BVID)
*   `BACKEND_PORT` (默认 8080)
*   `SCHEDULER_CRON` (默认 "0 0 * * *")
*   `GIN_MODE` (默认 "debug")

## 注意

*   此包旨在将配置加载逻辑与应用程序的其他部分分离。
*   目前配置严格从环境变量加载。
*   敏感信息（密码、Cookie）也通过环境变量加载，需要确保在部署环境中正确设置。 