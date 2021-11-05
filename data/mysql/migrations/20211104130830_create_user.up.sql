CREATE TABLE `user`
(
    `id`        BINARY(16)   NOT NULL,
    `login`     VARCHAR(255) NOT NULL,
    `firstname` VARCHAR(255) NOT NULL,
    `lastname`  VARCHAR(255) NOT NULL,
    `email`     VARCHAR(255) NOT NULL,
    `phone`     VARCHAR(255) NOT NULL,
    PRIMARY KEY (id),
    INDEX login_idx (`login`)
)
    ENGINE = InnoDB
    CHARACTER SET = utf8mb4
    COLLATE utf8mb4_unicode_ci
;