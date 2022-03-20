CREATE TABLE `order_position`
(
    `id`          BINARY(16) NOT NULL,
    `order_id`    BINARY(16) NOT NULL,
    `position_id` BINARY(16) NOT NULL,
    `count`      INT        NOT NULL,
    `updated_at`  DATETIME   NOT NULL,
    PRIMARY KEY (id),
    INDEX order_idx (`order_id`),
    INDEX position_idx (`position_id`)
)
    ENGINE = InnoDB
    CHARACTER SET = utf8mb4
    COLLATE utf8mb4_unicode_ci
;