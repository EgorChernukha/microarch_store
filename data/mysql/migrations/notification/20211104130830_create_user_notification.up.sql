CREATE TABLE `user_notification`
(
    `id`         MEDIUMINT    NOT NULL AUTO_INCREMENT,
    `user_id`    BINARY(16)   NOT NULL,
    `order_id`   BINARY(16)   NOT NULL,
    `email`      VARCHAR(255) NOT NULL,
    `message`    TEXT         NOT NULL,
    `created_at` DATETIME     NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id),
    INDEX user_order_idx (`user_id`, `order_id`)
)
    ENGINE = InnoDB
    CHARACTER SET = utf8mb4
    COLLATE utf8mb4_unicode_ci
;