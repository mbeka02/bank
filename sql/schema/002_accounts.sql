-- +goose Up 
CREATE TABLE accounts (
  id bigserial PRIMARY KEY,
  owner varchar NOT NULL,
  balance bigint NOT NULL,
  currency varchar NOT NULL,
  created_at timestamptz NOT NULL DEFAULT (now())
);
CREATE INDEX ON accounts(owner);

CREATE UNIQUE INDEX ON accounts(owner,currency);

ALTER TABLE accounts ADD FOREIGN KEY (owner) REFERENCES users(user_name);

-- +goose Down
DROP TABLE accounts;
