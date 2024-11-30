-- +goose Up
-- +goose StatementBegin
CREATE TABLE profiles (
  id SERIAL PRIMARY KEY,
  user_id INT NOT NULL,
  description TEXT,
  picture TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX profiles_user_id_idx ON profiles (user_id);
ALTER TABLE profiles ADD CONSTRAINT profiles_user_id_fkey FOREIGN KEY (user_id) REFERENCES users (id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
