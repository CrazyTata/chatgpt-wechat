CREATE TABLE `customer_prompt` (
                                   `id` int(11) NOT NULL AUTO_INCREMENT,
                                   `kf_name` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT '',
                                   `kf_id` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT '',
                                   `prompt` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '',
                                   PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci