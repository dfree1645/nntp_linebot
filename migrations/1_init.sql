-- +migrate Up
CREATE TABLE `groups` (
  `group_id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(256) NOT NULL COMMENT 'グループ名',
  `high` int(11) NOT NULL COMMENT '既知の最新',
  `low` int(11) NOT NULL COMMENT '既知の最古',
  PRIMARY KEY (`group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='グループ一覧';

CREATE TABLE `articles` (
  `article_id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `article_strid` varchar(256) NOT NULL COMMENT 'サーバー生成ID',
  `send_date` datetime NOT NULL COMMENT '送信日時',
  `subject` varchar(256) NOT NULL COMMENT '件名',
  `body` text NOT NULL COMMENT '本文',
  PRIMARY KEY (`article_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='記事一覧';

CREATE TABLE `users` (
  `user_id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `line_id` varchar(64) NOT NULL COMMENT 'LineID',
  `name` varchar(256) NOT NULL COMMENT 'Line表示名',
  `mail` varchar(256) NOT NULL COMMENT 'トークン送信先メールアドレス',
  PRIMARY KEY (`user_id`),
  UNIQUE(`line_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='認証済みユーザー(LINE友達)一覧';

CREATE TABLE `tokens` (
  `line_id` varchar(64) NOT NULL COMMENT 'LineID',
  `mail` varchar(256) COMMENT 'トークン送信先メールアドレス',
  `token` char(32) COMMENT 'トークン',
  `expiration_date` datetime COMMENT '有効期限',
  PRIMARY KEY (`line_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='メール認証待ちユーザー(LINE友達)一覧';

-- +migrate Down
DROP TABLE groups;
DROP TABLE articles;
DROP TABLE users;
DROP TABLE tokens;
