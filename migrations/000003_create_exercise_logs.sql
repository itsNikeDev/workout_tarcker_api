-- +goose Up
CREATE TABLE exercise_logs (
    id               UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id          UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    exercise_id      UUID        NOT NULL REFERENCES exercises (id),
    performed_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    sets             INT,
    reps             INT,
    weight_kg        NUMERIC(6, 2),
    duration_seconds INT,
    notes            TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX exercise_logs_user_id_idx      ON exercise_logs (user_id);
CREATE INDEX exercise_logs_performed_at_idx ON exercise_logs (performed_at);

-- +goose Down
DROP TABLE exercise_logs;
