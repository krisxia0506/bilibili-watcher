# Changelog

本项目所有显著变更都将记录在此文件中。

格式基于 [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),  
并且这个项目遵循 [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [未发布]

## [f29649b]  
### 添加  
- 支持 `BILIBILI_BVID` 环境变量，并在 `docker-compose.yml` 中提供默认值以确保兼容性  
- 优化 `_index.tsx` 中 BVID 获取逻辑：优先使用环境变量中的值，保证在缺少 URL 参数时仍能正常运行  

## [6b337b7]  
### 文档更新  
- 更新 `README.md`  

## [338ecf7]  
### 添加  
- 在客户端表单提交前，将 `datetime-local` 输入的本地时间字符串转换为 UTC ISO 格式  
- 实现初始页面自动加载逻辑：若 URL 无时间参数，自动触发表单提交，加载当天（00:00–23:59）数据  
- 深色模式支持：通过 Tailwind 类配置、`ThemeProvider`/`useTheme` 管理主题，并在 `<html>` 上应用 dark 类  
- Docker 前端 API URL 配置：在 `_index.tsx` 中支持 `BACKEND_API_URL` 环境变量  

### 修复  
- 后端时区处理：Go 服务端统一接收 UTC ISO 时间，消除 CST vs. UTC 环境下的歧义  
- Go 后端 `video_analytics_handler.go` 中 `interval` 参数解析：支持 “d” 日单位转小时后解析  
- `video_analytics_service.go` 观看时长计算：修复分段初始记录缺失导致时长为零的问题  

### 优化  
- Recharts 提示框增强：根据用户本地时区格式化显示各分段的开始/结束时间及时长  

## [099fa9d]  
### 添加  
- 主题切换功能：在 `tailwind.config.ts` 启用暗黑模式；新增 `ThemeProvider`、`useTheme`；在 `root.tsx` 添加切换按钮  
- 优化页面样式以适配深色模式  

## [400b281]  
### 添加  
- 新增 `WatchTimeChart` 组件，用于展示视频观看时长数据  
- 在首页集成图表组件，并优化表单参数输入  
- 增强错误处理：数据加载失败时给出用户友好提示  

## [44fbbd7]  
### 优化  
- 重构观看时长计算逻辑：使用 `p1`、`p2` 变量分别表示分段开始/结束状态  
- 增加无记录情况处理，确保在无进度记录时返回正确时长  
- 添加数据完整性检查，避免时间戳关系混乱  

## [fd4f878]  
### 添加  
- 优化时间间隔解析逻辑：支持 “d” 天单位输入，先转为小时再解析  
- 丰富错误反馈信息  

## [63db189]  
### 文档更新  
- 更新项目结构说明：调整目录层级，增加配置管理、接口层、SQL 相关内容  
- 优化注释，明确各目录用途  

## [88fffa5]  
### 文档更新  
- 更新 `README.md`：明确项目为 Bilibili 视频观看时长追踪与分析工具  
- 列出核心功能：定时获取进度、数据持久化、观看时长分析、API 服务、健康检查  
- 添加数据处理流程及请求时序可视化图示  
- 补充技术栈与定时任务信息  

## [6f98c94]  
### 添加  
- 在 `client.go` 的 `Get` 方法中新增 `ctx` 参数以支持上下文传递  
- 更新 `video_progress.go`、`video_view.go` 调用，传递上下文  

## [4620aa5]  
### 杂项  
- 清理注释：删除 `client.go`、`video_progress.go`、`video_view.go` 中不必要或过时的 TODO 标记  

## [e6745e5]  
### 添加  
- 移除硬编码 `SESSDATA`，在 `client.go` 中标记 TODO，准备动态加载  
- 更新注释，明确客户端结构体用途  

## [0d83a77]  
### 添加  
- 优化视频进度查询逻辑：增加时间范围 buffer，提高查询准确性  
- 在无记录时返回空结果而非错误  
- 更新日志信息，明确查询范围与状态  

## [a2d0128]  
### 添加  
- 新增 `VideoAnalyticsService` 接口及实现，用于计算视频观看分段时长  
- 在 `internal/interfaces/api/rest` 中新增视频分析路由与处理  
- 更新 `router.go`：配置健康检查和视频分析路由  
- 引入统一响应包，规范 API 格式  
- 更新配置加载：从环境变量读取 `BILIBILI_SESSDATA` 与 `BILIBILI_BVID`  

## [2b3a06f]  
### 添加  
- 在配置加载逻辑中支持从环境变量获取 `WATCH_TARGET_BVID`  
- 修改 `BilibiliConfig`，新增 `TargetBVID` 字段  
- 更新进度获取逻辑，使用配置中的目标 BVID 并完善错误处理  

## [12a87f2]  
### 添加  
- 将进度获取逻辑 AID 替换为 BVID，更新相关日志  
- 修改 `BilibiliClient` 接口，新增 BVID 参数  
- 更新 `VideoProgressService` 方法签名，接受 BVID  
- GORM 模型新增 `RecordedAt` 字段，记录进度创建时间  

## [df93de0]  
### 添加  
- 调度器更新：移除对 `VideoProgressService` 依赖，简化初始化  
- 新增 `ScheduleJob` 方法，支持按 cron 表达式注册任务  
- 在 `main.go` 初始化 `WatchTimeCalculator`、`WatchTimeService` 并调度获取进度作业  

## [9ab1786]  
### 添加  
- 新增 `WatchTimeService` 接口及实现，计算指定时段观看时长  
- 引入 `VideoPage` 模型，支持分 P 信息  
- 实现 `WatchTimeCalculator` 具体逻辑  
- 更新 `VideoViewDTO`，新增 `Pages` 字段  
- 映射 Bilibili API 分 P 信息至 DTO  

## [7659062]  
### 文档更新  
- 更新 `README.md`：添加参考项目链接，增强文档信息  

## [f50a0e9]  
### 重构  
- 按 DDD 重构 Bilibili API 交互：引入通用 `BilibiliClient` 接口与 DTO，解耦应用层与实现细节  
- 在 `infrastructure/bilibili.Client` 中实现该接口，并映射响应至 DTO  
- 更新 `VideoProgressService` 依赖新接口  
- 域模型新增 `TotalDuration` 与 `FetchTime` 字段  
- 仓库接口新增 `FindByAID` 与 `ErrVideoProgressNotFound`  
- 更新 GORM 模型与仓库实现；修改 `sql/schema.sql`  
- 移除旧的 `VideoProgressFetcher` 接口及实现  
- 修复因重构导致的依赖和类型错误  
- 更新各后端目录 README  

## [3d740e6]  
### 添加  
- 初始化视频进度管理功能：定时任务调度、数据持久化  
- 新增 cron 依赖；完善数据库模型  
- 实现进度获取与记录应用服务  
- 更新项目结构与文档，完善配置加载  

## [964ba24]  
### 修复  
- 更新 Docker 配置：确保后端与前端端口映射一致  
- 修正环境变量设置；前端启动命令由 `pnpm` 改为 `npm`  

## [dcfbee4]  
### 初始化  
- 项目结构初始化：后端与前端代码骨架、Docker 配置、示例环境变量、数据库模型、接口  
- 引入 Bilibili API 客户端，设置定时任务与数据持久化逻辑  
- 前端采用 Remix 与 Tailwind CSS  
- 添加 `README.md`，描述项目功能与开发规范  