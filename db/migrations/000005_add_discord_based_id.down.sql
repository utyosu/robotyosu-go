ALTER TABLE recruitments ADD channel_id INT UNSIGNED NOT NULL DEFAULT 0 AFTER id;
ALTER TABLE participants ADD user_id INT UNSIGNED NOT NULL DEFAULT 0 AFTER id;
ALTER TABLE nicknames ADD user_id INT UNSIGNED NOT NULL DEFAULT 0 AFTER id;

BEGIN;

UPDATE
  recruitments
LEFT JOIN
  channels
  ON recruitments.discord_channel_id = channels.discord_channel_id
SET
  recruitments.channel_id = channels.id;

UPDATE
  participants
LEFT JOIN
  users
  ON participants.discord_user_id = users.discord_user_id
SET
  participants.user_id = users.id;

UPDATE
  nicknames
LEFT JOIN
  users
  ON nicknames.discord_user_id = users.discord_user_id
SET
  nicknames.user_id = users.id;

COMMIT;

ALTER TABLE recruitments ADD KEY channel_id(channel_id, active);
ALTER TABLE recruitments DROP INDEX index_discord_channel_id_and_active;
ALTER TABLE recruitments DROP discord_channel_id;

ALTER TABLE participants DROP discord_user_id;

ALTER TABLE nicknames ADD UNIQUE KEY user_id(user_id, discord_guild_id);
ALTER TABLE nicknames ADD KEY user_id_2(user_id, discord_guild_id);
ALTER TABLE nicknames DROP KEY unique_discord_user_id_and_discord_guild_id;
ALTER TABLE nicknames DROP discord_user_id;

ALTER TABLE recruitments ALTER channel_id DROP DEFAULT;
ALTER TABLE participants ALTER user_id DROP DEFAULT;
ALTER TABLE nicknames ALTER user_id DROP DEFAULT;
