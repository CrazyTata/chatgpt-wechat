create table if not exists wechat_user
(
    id          bigint unsigned auto_increment ,
    user    varchar(128)                              not null comment '微信客户的external_userid',
    nickname  varchar(128)   default  ''                      not null comment '微信昵称',
    avatar VARCHAR(256) default '' not null comment '微信头像',
    gender tinyint  default 0 not null comment '性别',
    unionid VARCHAR(128) default '' not null comment 'unionid，需要绑定微信开发者帐号才能获取到',
    created_at  timestamp       default CURRENT_TIMESTAMP null comment '创建时间',
    updated_at  timestamp       default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP comment '更新时间',
    PRIMARY KEY (`id`),
    KEY           `idx_nickname` (`nickname`) USING BTREE,
    KEY           `idx_user` (`user`) USING BTREE
    ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

