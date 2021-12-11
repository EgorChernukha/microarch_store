CREATE TABLE `user_session`
(
    `id`         BINARY(16) NOT NULL,
    `user_id`    BINARY(16) NOT NULL,
    `valid_till` TIMESTAMP  NOT NULL,
    PRIMARY KEY (id),
    INDEX user_id_idx (`user_id`)
)
    ENGINE = InnoDB
    CHARACTER SET = utf8mb4
    COLLATE utf8mb4_unicode_ci
;