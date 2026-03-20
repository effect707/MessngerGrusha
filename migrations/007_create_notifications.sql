-- +goose Up
CREATE TABLE notifications (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type       VARCHAR(64) NOT NULL,
    payload    JSONB       NOT NULL DEFAULT '{}',
    is_read    BOOLEAN     NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_notifications_user_unread ON notifications (user_id, created_at DESC)
    WHERE is_read = false;

-- +goose Down
DROP TABLE IF EXISTS notifications;
