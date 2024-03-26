-- +goose Up

CREATE TABLE transactions (
  id bigserial PRIMARY KEY,
  amount bigint NOT NULL,
  sender_id bigint NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
  receiver_id bigint NOT NULL REFERENCES accounts(id) ON DELETE CASCADE
);

CREATE INDEX ON transactions(sender_id);

CREATE INDEX ON transactions(receiver_id);

CREATE INDEX ON transactions(sender_id, receiver_id);



-- +goose Down
DROP TABLE transactions;
