SET
  @discord_channel_id = 123456789,
  @discord_guild_id = 987654321,
  @discord_user_id = 11111111;

INSERT channels (id, discord_channel_id, discord_guild_id, recruitment, timezone, twitter_config_id, created_at, updated_at) VALUES(1, @discord_channel_id, @discord_guild_id, 1, "Asia/Tokyo", 1, CURRENT_TIME(), CURRENT_TIME());
INSERT recruitments (id, channel_id, label, title, capacity, active, notified, reserve_at, expire_at) VALUES(1, 1, 1, "hogehoge@1", 2, 1, 0, ADDTIME(CURRENT_TIME(), "0 0:10:0"), ADDTIME(CURRENT_TIME(), "0 0:30:0"));
INSERT users (id, discord_user_id, name) VALUES(1, @discord_user_id, "テストさん1");
INSERT participants (id, user_id, recruitment_id) VALUES(1, 1, 1);
