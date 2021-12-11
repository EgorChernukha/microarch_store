CREATE TABLE `user_auth`
(
    `id`       BINARY(16)          NOT NULL,
    `login`    VARCHAR(255) UNIQUE NOT NULL,
    `password` VARCHAR(255)        NOT NULL,
    PRIMARY KEY (id),
    INDEX login_idx (`login`)
)
    ENGINE = InnoDB
    CHARACTER SET = utf8mb4
    COLLATE utf8mb4_unicode_ci
;