-- +goose Up
-- +goose StatementBegin
CREATE TABLE matches (
  id SERIAL PRIMARY KEY,
  profile_id INT NOT NULL,
  partner_id INT NOT NULL,
  status VARCHAR(255) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX matches_profile_id_index ON matches (profile_id);
CREATE INDEX matches_partner_id_index ON matches (partner_id);
CREATE INDEX matches_status_index ON matches (status);

ALTER TABLE matches ADD CONSTRAINT matches_profile_id_fk FOREIGN KEY (profile_id) REFERENCES profiles (id);
ALTER TABLE matches ADD CONSTRAINT matches_partner_id_fk FOREIGN KEY (partner_id) REFERENCES profiles (id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE matches;
-- +goose StatementEnd
