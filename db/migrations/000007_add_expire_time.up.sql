ALTER TABLE channels ADD expire_duration INT UNSIGNED DEFAULT 3600 AFTER reserve_limit_time;
ALTER TABLE channels ADD expire_duration_for_reserve INT UNSIGNED DEFAULT 1800 AFTER expire_duration;
