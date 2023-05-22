CREATE TABLE IF NOT EXISTS `script_log`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `script_id` bigint unsigned NOT NULL COMMENT '脚本id',
    `result` text not null default '' COMMENT '执行结果',
    `execution_count` tinyint unsigned DEFAULT 1 NOT NULL COMMENT '执行次数',
    `status` tinyint unsigned DEFAULT 1 NOT NULL COMMENT '执行状态，1开始执行，2执行中，3执行成功，4执行失败',
    `end_at`  timestamp NULL COMMENT '结束时间',
    `created_at`  timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`  timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ;