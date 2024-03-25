-- +goose Up
CREATE TABLE transfers(
  id bigserial PRIMARY KEY,
  amount bigint NOT NULL,
  sender_id bigint NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
  receiver_id bigint NOT NULL REFERENCES accounts(id) ON DELETE CASCADE
);



CREATE INDEX ON transfers(sender_id);

CREATE INDEX ON transfers(receiver_id);

CREATE INDEX ON transfers(sender_id, receiver_id);

--+goose Down
DROP TABLE transfers;
