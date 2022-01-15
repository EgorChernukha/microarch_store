CREATE TABLE `processed_event`
(
    `event_id`     BINARY(16) NOT NULL,
    `created_at` DATETIME   NOT NULL DEFAULT NOW(),
    PRIMARY KEY (event_id)
)
    ENGINE = InnoDB
    CHARACTER SET = utf8mb4
    COLLATE utf8mb4_unicode_ci
;