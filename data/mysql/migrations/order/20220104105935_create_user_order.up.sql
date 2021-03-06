CREATE TABLE `user_order`
(
    `id`         BINARY(16)     NOT NULL,
    `user_id`    BINARY(16)     NOT NULL,
    `order_id`   BINARY(16)     NOT NULL,
    `price`      DECIMAL(10, 2) NOT NULL DEFAULT 0,
    `status`     INT            NOT NULL,
    `created_at` DATETIME       NOT NULL DEFAULT NOW(),
    `updated_at` DATETIME       NOT NULL,
    PRIMARY KEY (id),
    INDEX user_order_idx (`user_id`, `order_id`),
    INDEX status_idx (`status`)
)
    ENGINE = InnoDB
    CHARACTER SET = utf8mb4
    COLLATE utf8mb4_unicode_ci
;