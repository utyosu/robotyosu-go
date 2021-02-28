CREATE TABLE nicknames (
  id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  user_id INT UNSIGNED NOT NULL,
  discord_guild_id BIGINT NOT NULL,
  name varchar(100) NOT NULL,
  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at timestamp NULL DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE (user_id, discord_guild_id),
  INDEX (user_id, discord_guild_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
