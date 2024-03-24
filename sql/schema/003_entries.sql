-- +goose Up

CREATE TABLE entries (
  id bigserial PRIMARY KEY,
  account_id bigint NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
  amount bigint NOT NULL,
  created_at timestamptz NOT NULL DEFAULT (now())
);



CREATE INDEX ON entries(account_id);


-- +goose Down
DROP TABLE entries;
