-- +goose Up

CREATE TABLE IF NOT EXISTS tokens (
    hash bytea PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users ON DELETE CASCADE,
    scope TEXT NOT NULL,
    expires_at TIMESTAMP(0) with time zone NOT NULL,
    created_at TIMESTAMP(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP(0) with time zone NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS tokens;
