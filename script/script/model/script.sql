CREATE TABLE IF NOT EXISTS `script`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(64) DEFAULT '' NOT NULL COMMENT '脚本名称',
    `path`  varchar(191) DEFAULT '' NOT NULL COMMENT '脚本保存路径',
    `script_type`  varchar(32) DEFAULT '' NOT NULL COMMENT '脚本类型',
    `is_delete` tinyint unsigned DEFAULT 0 NOT NULL COMMENT '是否删除，0否，1是',
    `is_enable` tinyint unsigned DEFAULT 1 NOT NULL COMMENT '是否启用，0否，1是',
    `created_at`  timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`  timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ;