-- +migrate Up
CREATE TABLE `groups` (
  `group_id` int(11) NOT NULL,
  `name` varchar(256) NOT NULL,
  `high` int(11) NOT NULL,
  `low` int(11) NOT NULL,
  PRIMARY KEY (`group_id`)
);

CREATE TABLE `articles` (
  `article_id` int(11) NOT NULL,
  `article_strid` varchar(256) NOT NULL,
  `send_date` datetime NOT NULL,
  `subject` varchar(256) NOT NULL,
  `body` text NOT NULL,
  PRIMARY KEY (`article_id`)
);

CREATE TABLE `users` (
  `user_id` int(11) NOT NULL,
  `line_id` varchar(64) NOT NULL,
  `name` varchar(256) NOT NULL,
  `mail` varchar(256) NOT NULL,
  PRIMARY KEY (`user_id`),
  UNIQUE(`line_id`)
);

CREATE TABLE `tokens` (
  `line_id` varchar(64) NOT NULL,
  `mail` varchar(256),
  `token` char(32),
  `expiration_date` datetime,
  PRIMARY KEY (`line_id`)
);

-- +migrate Down
DROP TABLE groups;
DROP TABLE articles;
DROP TABLE users;
DROP TABLE tokens;
