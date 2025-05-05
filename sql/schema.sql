-- Schema for bilibili-watcher
-- Target: MySQL 8
-- Applying Alibaba spec (NOT NULL constraints, singular table name) and user preferences (no soft delete, timestamp columns at the end).

-- 视频观看进度表 (Video Progress Table)
CREATE TABLE IF NOT EXISTS `video_progress` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `aid` bigint NOT NULL COMMENT '视频稿件 ID (AV 号)',
  `cid` bigint NOT NULL COMMENT '视频分 P ID',
  `bvid` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '视频 BV 号',
  `progress` int NOT NULL DEFAULT 0 COMMENT '观看进度 (毫秒)',
  `recorded_at` datetime(3) NOT NULL COMMENT '记录时间',
  `created_at` datetime(3) NOT NULL COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  INDEX `idx_video_progress_aid` (`aid`),
  INDEX `idx_video_progress_cid` (`cid`),
  INDEX `idx_video_progress_recorded_at` (`recorded_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='视频观看进度记录';

-- Note:
-- GORM needs to be configured to use singular table names (or override for this specific model).
-- GORM needs to be configured not to use soft delete if `deleted_at` is removed from the model struct as well.
-- `created_at` and `updated_at` are now at the end of the column list.
-- All relevant columns are now NOT NULL, with defaults where appropriate. 