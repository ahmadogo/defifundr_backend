-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE sessions (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  refresh_token VARCHAR(255) NOT NULL,
  user_agent VARCHAR(255) NOT NULL,
  web_oauth_client_id VARCHAR(255),
  oauth_access_token VARCHAR(1024),
  oauth_id_token VARCHAR(255),
  user_login_type VARCHAR(50) NOT NULL,
  mfa_enabled BOOLEAN NOT NULL DEFAULT false,
  client_ip VARCHAR(255) NOT NULL,
  is_blocked BOOLEAN NOT NULL DEFAULT false,
  expires_at TIMESTAMP NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT now()
);

-- Create index on user_id for faster lookups
CREATE INDEX idx_sessions_user_id ON sessions(user_id);

-- Create index on refresh_token for faster lookups
CREATE INDEX idx_sessions_refresh_token ON sessions(refresh_token);

-- Create index on expiration time for cleanup operations
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS sessions;