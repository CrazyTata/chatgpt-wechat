create table if not exists application_config
(
    id          bigint unsigned auto_increment ,
    agent_id    int unsigned default 0                 not null comment '应用ID',
    agent_secret    varchar(128) default  ''                     not null comment '应用secret',
    agent_name    varchar(128) default  ''                    not null comment '应用名',
    model  varchar(128)   default  ''                      not null comment 'model',
    post_model  varchar(128)   default  ''                      not null comment '发送请求的model',
    base_prompt VARCHAR(1000) default '' not null comment 'openai 基础设定（可选）',
    welcome VARCHAR(1000)  default '' not null comment '进入应用时的欢迎语',
    group_enable BOOLEAN DEFAULT false  not null comment '是否启用ChatGPT应用内部交流群',
    group_name VARCHAR(64) default '' not null comment 'ChatGPT群名',
    group_chat_id VARCHAR(128) default '' not null comment 'ChatGPT应用内部交流群chat_id',
    embedding_enable BOOLEAN DEFAULT false  not null comment '是否启用embedding',
    embedding_mode VARCHAR(64) default '' not null comment 'embedding的搜索模式',
    score DECIMAL(3, 1) comment '分数',
    top_k smallint DEFAULT 1  not null comment 'topK',
    clear_context_time int DEFAULT 0  not null comment '需要清理上下文的时间，按分配置，默认0不清理',
    created_at  timestamp       default CURRENT_TIMESTAMP null comment '创建时间',
    updated_at  timestamp       default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP comment '更新时间',
    is_deleted tinyint default 0 not null comment '是否删除',
    PRIMARY KEY (`id`),
    KEY           `idx_agent_id` (`agent_id`) USING BTREE
    ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
