# Domain Models (Entities & Value Objects)

此目录包含领域模型的核心定义，即实体（Entities）和值对象（Value Objects）。

## 主要职责

*   定义代表业务核心概念的数据结构。
*   **实体 (Entities):** 具有唯一标识符（ID）并且在生命周期中状态可变的对象。它们封装了与其相关的业务规则和行为。例如，一个"用户"或一个"订单"。
*   **值对象 (Value Objects):** 没有唯一标识符，通过其属性值来定义的对象。它们通常是不可变的。例如，"地址"或"金额"。
*   封装简单的、与自身属性相关的验证逻辑和业务规则。

## 当前内容

*   `video_progress.go`: 定义了视频观看进度记录的实体。
    *   `VideoProgress` 结构体: 代表一个时间点的观看进度快照。
        *   包含字段：`ID`, `AID`, `BVID`, `LastPlayCID`, `LastPlayTime`, `RecordedAt`, `GmtCreate`, `GmtModified`。
        *   通过 GORM 标签定义了数据库映射（列名、索引、非空、默认值、自动时间戳）。
        *   `TableName()` 方法: 显式指定数据库表名为 `video_progress`。

## 注意

*   模型应专注于表示业务概念和规则，避免包含与持久化或外部交互相关的逻辑。
*   复杂的、跨多个模型的业务逻辑应放在领域服务中。 