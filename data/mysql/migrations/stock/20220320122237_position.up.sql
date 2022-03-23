CREATE TABLE `position`
(
    `id`         BINARY(16)    NOT NULL,
    `total`      INT           NOT NULL,
    `title`      VARCHAR(1024) NOT NULL,
    `updated_at` DATETIME      NOT NULL,
    PRIMARY KEY (id),
    UNIQUE INDEX title_idx (`title`)
)
    ENGINE = InnoDB
    CHARACTER SET = utf8mb4
    COLLATE utf8mb4_unicode_ci
;