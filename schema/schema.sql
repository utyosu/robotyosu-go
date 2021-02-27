CREATE TABLE users (
  id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  discord_user_id BIGINT NOT NULL,
  name varchar(100) NOT NULL,
  created_at timestamp NULL DEFAULT NULL,
  updated_at timestamp NULL DEFAULT NULL,
  deleted_at timestamp NULL DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE recruitments (
  id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  channel_id INT UNSIGNED NOT NULL,
  label INT UNSIGNED NOT NULL,
  title varchar(100) NOT NULL,
  capacity INT UNSIGNED NOT NULL,
  active BOOLEAN NOT NULL,
  notified BOOLEAN NOT NULL,
  tweet_id BIGINT NOT NULL DEFAULT 0,
  reserve_at timestamp NULL DEFAULT NULL,
  expire_at timestamp NULL DEFAULT NULL,
  created_at timestamp NULL DEFAULT NULL,
  updated_at timestamp NULL DEFAULT NULL,
  deleted_at timestamp NULL DEFAULT NULL,
  PRIMARY KEY (id),
  INDEX (channel_id, active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE participants (
  id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  user_id INT UNSIGNED NOT NULL,
  recruitment_id INT UNSIGNED NOT NULL,
  created_at timestamp NULL DEFAULT NULL,
  updated_at timestamp NULL DEFAULT NULL,
  deleted_at timestamp NULL DEFAULT NULL,
  PRIMARY KEY (id),
  INDEX (recruitment_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE channels (
  id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  discord_channel_id BIGINT NOT NULL,
  recruitment BOOLEAN NOT NULL,
  timezone varchar(30) NOT NULL,
  language varchar(10) NOT NULL DEFAULT "",
  twitter_config_id INT UNSIGNED NOT NULL DEFAULT 0,
  created_at timestamp NULL DEFAULT NULL,
  updated_at timestamp NULL DEFAULT NULL,
  deleted_at timestamp NULL DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE twitter_configs (
  id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  consumer_key varchar(100) NOT NULL,
  consumer_secret varchar(100) NOT NULL,
  access_token varchar(100) NOT NULL,
  access_token_secret varchar(100) NOT NULL,
  title varchar(100) NOT NULL,
  created_at timestamp NULL DEFAULT NULL,
  updated_at timestamp NULL DEFAULT NULL,
  deleted_at timestamp NULL DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
