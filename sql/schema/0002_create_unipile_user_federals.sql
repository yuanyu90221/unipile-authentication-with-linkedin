-- +goose Up
CREATE TABLE IF NOT EXISTS unipile_user_federals (
  id BIGSERIAL PRIMARY KEY,
  account_id VARCHAR(200) NOT NULL,
  provider VARCHAR(200) NOT NULL,
  user_id BIGINT NOT NULL,
  status VARCHAR(400) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP NOT NULL DEFAULT now()
);


COMMENT ON COLUMN unipile_user_federals.account_id IS 'linkedin account_id';
COMMENT ON COLUMN unipile_user_federals.provider IS 'for different provider, LinkedIn';
COMMENT ON COLUMN unipile_user_federals.status IS 'record current linked status';

ALTER TABLE unipile_user_federals ADD FOREIGN KEY (user_id) REFERENCES users(id);
-- +goose Down
DROP TABLE IF EXISTS unipile_user_federals CASCADE;