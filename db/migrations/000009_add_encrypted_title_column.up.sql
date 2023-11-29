ALTER TABLE recruitments ADD encrypted_title VARBINARY(256) NOT NULL DEFAULT '' AFTER title;
