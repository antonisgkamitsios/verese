-- +goose Up
CREATE TABLE IF NOT EXISTS users (
  id bigint GENERATED ALWAYS AS IDENTITY,
  name text NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS users;
