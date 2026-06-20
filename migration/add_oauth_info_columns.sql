
ALTER TABLE users ADD COLUMN IF NOT EXISTS oauth_user_id VARCHAR(50);
ALTER TABLE users ADD COLUMN IF NOT EXISTS oauth_provider_name VARCHAR(20);
ALTER TABLE users ADD COLUMN IF NOT EXISTS avatar_url VARCHAR(500);
ALTER TABLE users ADD CONSTRAINT oauth_info_uni UNIQUE(oauth_user_id, oauth_provider_name);
