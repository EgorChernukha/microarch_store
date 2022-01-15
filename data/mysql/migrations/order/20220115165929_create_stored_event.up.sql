CREATE TABLE `stored_event`
(
    `id`         MEDIUMINT      NOT NULL AUTO_INCREMENT,
    `uid`         BINARY(16)    NOT NULL,
    `type`       VARCHAR(255)  NOT NULL,
    `body`       VARCHAR(2047) NOT NULL,
    `confirmed`  BOOLEAN       NOT NULL DEFAULT FALSE,
    `created_at` DATETIME      NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id),
    INDEX uid_idx (`uid`)
)
    ENGINE = InnoDB
    CHARACTER SET = utf8mb4
    COLLATE utf8mb4_unicode_ci
;