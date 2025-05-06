-- Schema for bilibili-watcher
-- Target: MySQL 8
-- Applying Alibaba spec (mandatory fields, NOT NULL, defaults, singular table name).

-- 视频观看进度表 (Video Progress Table)
CREATE TABLE IF NOT EXISTS `video_progress` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `aid` bigint NOT NULL DEFAULT 0 COMMENT '视频稿件 ID (AV 号)',
  `bvid` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '视频 BV 号',
  `last_play_cid` bigint NOT NULL DEFAULT 0 COMMENT '上次播放的视频分 P ID',
  `last_play_time` int NOT NULL DEFAULT 0 COMMENT '上次播放时间/进度 (毫秒)',
  `recorded_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '记录时间',
  `gmt_create` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `gmt_modified` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  INDEX `idx_video_progress_aid` (`aid`),
  INDEX `idx_video_progress_last_play_cid` (`last_play_cid`),
  INDEX `idx_video_progress_recorded_at` (`recorded_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='视频观看进度记录';

-- Note:
-- Mandatory fields: id, gmt_create, gmt_modified.
-- All fields are NOT NULL with defaults.
-- GORM needs to be configured for singular table names.
-- GORM automatically handles `gmt_create` and `gmt_modified` if column names are mapped correctly in the model. 