CREATE TABLE `user_account`
(
    `id`         BINARY(16)     NOT NULL,
    `user_id`    BINARY(16)     NOT NULL,
    `balance`    DECIMAL(10, 2) NOT NULL DEFAULT 0,
    `updated_at` DATETIME       NOT NULL,
    PRIMARY KEY (id),
    INDEX user_idx (`user_id`)
)
    ENGINE = InnoDB
    CHARACTER SET = utf8mb4
    COLLATE utf8mb4_unicode_ci
;