-- +goose Up
-- +goose StatementBegin
CREATE TABLE profile_view_logs (
  id SERIAL PRIMARY KEY,
  profile_id INT NOT NULL,
  viewer_id INT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX profile_view_logs_profile_id_idx ON profile_view_logs (profile_id);
CREATE INDEX profile_view_logs_viewer_id_idx ON profile_view_logs (viewer_id);

ALTER TABLE profile_view_logs ADD CONSTRAINT profile_view_logs_profile_id_fkey FOREIGN KEY (profile_id) REFERENCES profiles (id);
ALTER TABLE profile_view_logs ADD CONSTRAINT profile_view_logs_viewer_id_fkey FOREIGN KEY (viewer_id) REFERENCES profiles (id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE profile_view_logs;
-- +goose StatementEnd
