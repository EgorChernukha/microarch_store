CREATE TABLE `user`
(
    `id`        MEDIUMINT    NOT NULL AUTO_INCREMENT,
    `login`     VARCHAR(255) NOT NULL,
    `firstname` VARCHAR(255) NOT NULL,
    `email`     VARCHAR(255),
    `phone`     VARCHAR(255),
    PRIMARY KEY (id),
    INDEX login_idx(`login`)
)
    ENGINE = InnoDB
    CHARACTER SET = utf8mb4
    COLLATE utf8mb4_unicode_ci
;