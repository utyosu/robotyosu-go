ALTER TABLE users MODIFY COLUMN created_at timestamp NULL DEFAULT NULL;
ALTER TABLE users MODIFY COLUMN updated_at timestamp NULL DEFAULT NULL;

ALTER TABLE recruitments MODIFY COLUMN created_at timestamp NULL DEFAULT NULL;
ALTER TABLE recruitments MODIFY COLUMN updated_at timestamp NULL DEFAULT NULL;

ALTER TABLE participants MODIFY COLUMN created_at timestamp NULL DEFAULT NULL;
ALTER TABLE participants MODIFY COLUMN updated_at timestamp NULL DEFAULT NULL;

ALTER TABLE channels MODIFY COLUMN created_at timestamp NULL DEFAULT NULL;
ALTER TABLE channels MODIFY COLUMN updated_at timestamp NULL DEFAULT NULL;

ALTER TABLE twitter_configs MODIFY COLUMN created_at timestamp NULL DEFAULT NULL;
ALTER TABLE twitter_configs MODIFY COLUMN updated_at timestamp NULL DEFAULT NULL;
