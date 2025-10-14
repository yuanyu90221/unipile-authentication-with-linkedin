-- +goose Up
CREATE TABLE IF NOT EXISTS users (
  id BIGSERIAL PRIMARY KEY,
  account VARCHAR(100) UNIQUE NOT NULL,
  hashed_password VARCHAR(400) NOT NULL,
  refresh_token VARCHAR(400),
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP NOT NULL DEFAULT now()
);

COMMENT ON COLUMN users.account IS 'user account for login';
COMMENT ON COLUMN users.hashed_password IS 'hashed password for login';
COMMENT ON COLUMN users.refresh_token IS 'refresh token for login session';
CREATE INDEX IF NOT EXISTS account ON users (account);
-- +goose Down
DROP INDEX IF EXISTS account CASCADE;
DROP TABLE IF EXISTS users CASCADE;