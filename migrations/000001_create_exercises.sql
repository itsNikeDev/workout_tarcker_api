-- +goose Up
CREATE TABLE exercises (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX exercises_name_idx ON exercises (name);

-- +goose Down
DROP TABLE exercises;
