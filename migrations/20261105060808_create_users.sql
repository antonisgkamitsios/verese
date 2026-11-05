-- +goose Up

CREATE TABLE IF NOT EXISTS users (
    id bigint GENERATED ALWAYS AS IDENTITY,
    name text NOT NULL,
    email citext UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    activated bool NOT NULL,

    version integer NOT NULL DEFAULT 1,
    created_at TIMESTAMP(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP(0) with time zone NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS users;