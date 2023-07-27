CREATE TABLE IF NOT EXISTS `files`
(
    `id`          bigint unsigned NOT NULL AUTO_INCREMENT,
    `sys`         tinyint         NOT NULL DEFAULT 0 comment '',
    `type`        tinyint         NOT NULL DEFAULT 0 comment '',
    `item_id`     bigint unsigned NOT NULL DEFAULT 0,
    `user_id`     bigint unsigned NOT NULL DEFAULT 0,
    `file_type`   tinyint         NOT NULL DEFAULT 0 comment 'type: 0 default, 1 image, 2 video, 3 doc, 4 other',
    `name`        varchar(255)    NOT NULL DEFAULT '' comment 'name',
    `ext`         varchar(255)    NOT NULL DEFAULT '' comment 'ext',
    `path`        varchar(255)    NOT NULL DEFAULT '' comment 'path',
    `hash`        varchar(255)    NOT NULL DEFAULT '' comment 'path',
    `status`      tinyint         NOT NULL DEFAULT 0 COMMENT 'status:0 default, 1 active, 2 disable',
    `create_time` timestamp       NULL     DEFAULT NULL,
    `update_time` timestamp       NULL     DEFAULT NULL,
    `delete_time` timestamp       NULL     DEFAULT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_user_id` (`user_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;
