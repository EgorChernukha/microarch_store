CREATE TABLE `order_delivery`
(
    `id`         BINARY(16) NOT NULL,
    `order_id`   BINARY(16) NOT NULL,
    `user_id`    BINARY(16) NOT NULL,
    `status`     INT        NOT NULL,
    `updated_at` DATETIME   NOT NULL,
    PRIMARY KEY (id),
    UNIQUE INDEX order_idx (`order_id`),
    INDEX user_idx (`user_id`)
)
    ENGINE = InnoDB
    CHARACTER SET = utf8mb4
    COLLATE utf8mb4_unicode_ci
;