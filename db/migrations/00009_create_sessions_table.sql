-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE sessions (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  refresh_token VARCHAR NOT NULL,
  user_agent VARCHAR NOT NULL,
  web_oauth_client_id VARCHAR,
  oauth_access_token VARCHAR,
  oauth_id_token VARCHAR,
  user_login_type VARCHAR NOT NULL,
  mfa_enabled BOOLEAN NOT NULL DEFAULT false,
  client_ip VARCHAR NOT NULL,
  is_blocked BOOLEAN NOT NULL DEFAULT false,
  expires_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS sessions;