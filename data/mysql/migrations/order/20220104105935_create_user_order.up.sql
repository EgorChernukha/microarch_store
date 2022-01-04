CREATE TABLE `user_order`
(
    `id`         MEDIUMINT      NOT NULL AUTO_INCREMENT,
    `user_id`    BINARY(16)     NOT NULL,
    `price`      DECIMAL(10, 2) NOT NULL DEFAULT 0,
    `status`     INT            NOT NULL,
    `created_at` DATETIME       NOT NULL,
    `updated_at` DATETIME       NOT NULL,
    PRIMARY KEY (id),
    INDEX user_idx (`user_id`),
    INDEX status_idx (`status`)
)
    ENGINE = InnoDB
    CHARACTER SET = utf8mb4
    COLLATE utf8mb4_unicode_ci
;