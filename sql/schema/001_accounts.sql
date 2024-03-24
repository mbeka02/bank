-- +goose Up 
CREATE TABLE accounts (
  id bigserial PRIMARY KEY,
  full_name varchar NOT NULL,
  balance bigint NOT NULL,
  currency varchar NOT NULL,
  created_at timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON accounts(full_name);

-- +goose Down
DROP TABLE accounts;
