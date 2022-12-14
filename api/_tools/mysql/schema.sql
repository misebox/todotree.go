CREATE TABLE IF NOT EXISTS `user` (
    `id`       BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ユーザーの識別子',
    `email`    VARCHAR(200) NOT NULL COMMENT 'e-mailアドレス',
    `name`     VARCHAR(40) NOT NULL COMMENT 'ユーザー名',
    `password` VARCHAR(160) NOT NULL COMMENT 'パスワード',
    `role`     VARCHAR(80) NOT NULL COMMENT 'ロール名',
    `created`  DATETIME NOT NULL COMMENT 'レコード作成日時',
    `modified` DATETIME NOT NULL COMMENT 'レコード修正日時',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uix_email` (`email`) USING BTREE,
    UNIQUE KEY `uix_name` (`name`) USING BTREE
) Engine = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = 'ユーザー';

CREATE TABLE IF NOT EXISTS `task` (
    `id`       BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'タスクの識別子',
    `user_id`  BIGINT UNSIGNED NOT NULL COMMENT 'タスクを作成したユーザーの識別子',
    `root_id`  BIGINT UNSIGNED NULL COMMENT 'ルートとなるタスクの識別子',
    `parent_id`  BIGINT UNSIGNED NULL COMMENT '親タスクの識別子',
    `title`    VARCHAR(200) NOT NULL COMMENT 'タスクのタイトル',
    `status`   VARCHAR(20) NOT NULL COMMENT 'タスクの状態',
    `created`  DATETIME NOT NULL COMMENT 'レコード作成日時',
    `modified` DATETIME NOT NULL COMMENT 'レコード修正日時',
    PRIMARY KEY (`id`),
    CONSTRAINT `fk_user_id`
        FOREIGN KEY (`user_id`) REFERENCES `user` (`id`)
            ON DELETE RESTRICT ON UPDATE RESTRICT
) Engine = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = 'タスク';
