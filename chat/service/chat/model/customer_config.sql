
create table if not exists customer_config
(
    id          bigint unsigned auto_increment ,
    kf_id    varchar(128)                  not null comment '客服ID',
    kf_name varchar(50)  default '' null,
    prompt  varchar(1000) default '' null,
    embedding_enable BOOLEAN DEFAULT false  not null comment '是否启用embedding',
    embedding_mode VARCHAR(64) default '' not null comment 'embedding的搜索模式',
    score DECIMAL(3, 1) comment '分数',
    top_k smallint DEFAULT 1  not null comment 'topK',
    clear_context_time int DEFAULT 0  not null comment '需要清理上下文的时间，按分配置，默认0不清理',
    created_at  timestamp       default CURRENT_TIMESTAMP null comment '创建时间',
    updated_at  timestamp       default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP comment '更新时间',
    PRIMARY KEY (`id`),
    KEY           `idx_open_kf_id` (`kf_id`) USING BTREE
    ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;