
ALTER TABLE users ADD COLUMN oauth_user_id VARCHAR(50);
ALTER TABLE users ADD COLUMN oauth_provider_name VARCHAR(20);
ALTER TABLE users ADD constraint oauth_info_uni UNIQUE(oauth_user_id, oauth_provider_name);
