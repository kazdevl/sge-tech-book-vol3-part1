CREATE TABLE IF NOT EXISTS `user_coin` (
    `user_id` BIGINT (20) NOT NULL COMMENT 'ユーザID',
    `num` BIGINT (20) NOT NULL COMMENT '個数',
    PRIMARY KEY (`user_id`)
    ) ENGINE = InnoDB, COMMENT = 'ユーザコイン', DEFAULT CHARACTER SET = utf8mb4;


CREATE TABLE IF NOT EXISTS `user_item` (
    `user_id` BIGINT (20) NOT NULL COMMENT 'ユーザID',
    `item_id` BIGINT (20) NOT NULL COMMENT 'アイテムID',
    `count` BIGINT (20) NOT NULL COMMENT '数',
    PRIMARY KEY (`user_id`,`item_id`)
    ) ENGINE = InnoDB, COMMENT = 'ユーザアイテム', DEFAULT CHARACTER SET = utf8mb4;

CREATE TABLE IF NOT EXISTS `user_monster` (
    `user_id` BIGINT (20) NOT NULL COMMENT 'ユーザID',
    `monster_id` BIGINT (20) NOT NULL COMMENT 'モンスターID',
    `exp` BIGINT (20) NOT NULL COMMENT '経験値',
    PRIMARY KEY (`user_id`,`monster_id`)
    ) ENGINE = InnoDB, COMMENT = 'ユーザアイテム', DEFAULT CHARACTER SET = utf8mb4;
