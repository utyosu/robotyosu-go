ALTER TABLE recruitments ADD discord_channel_id BIGINT NOT NULL DEFAULT 0 AFTER channel_id;
ALTER TABLE participants ADD discord_user_id BIGINT NOT NULL DEFAULT 0 AFTER user_id;
ALTER TABLE nicknames ADD discord_user_id BIGINT NOT NULL DEFAULT 0 AFTER user_id;

BEGIN;

UPDATE
  recruitments
LEFT JOIN
  channels
  ON recruitments.channel_id = channels.id
SET
  recruitments.discord_channel_id = channels.discord_channel_id;

UPDATE
  participants
LEFT JOIN
  users
  ON participants.user_id = users.id
SET
  participants.discord_user_id = users.discord_user_id;

UPDATE
  nicknames
LEFT JOIN
  users
  ON nicknames.user_id = users.id
SET
  nicknames.discord_user_id = users.discord_user_id;

COMMIT;

ALTER TABLE recruitments DROP KEY channel_id;
ALTER TABLE recruitments DROP channel_id;
ALTER TABLE recruitments ADD INDEX index_discord_channel_id_and_active(discord_channel_id, active);

ALTER TABLE participants DROP user_id;

ALTER TABLE nicknames DROP KEY user_id;
ALTER TABLE nicknames DROP KEY user_id_2;
ALTER TABLE nicknames DROP user_id;
ALTER TABLE nicknames ADD UNIQUE KEY unique_discord_user_id_and_discord_guild_id(discord_user_id, discord_guild_id);
