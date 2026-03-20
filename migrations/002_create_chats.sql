-- +goose Up
CREATE TYPE chat_type AS ENUM ('direct', 'group');

CREATE TABLE chats (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type       chat_type    NOT NULL,
    name       VARCHAR(255),
    avatar_url TEXT,
    created_by UUID         REFERENCES users(id),
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE TABLE chat_members (
    chat_id   UUID        NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    user_id   UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role      VARCHAR(20) NOT NULL DEFAULT 'member',
    joined_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (chat_id, user_id)
);

CREATE INDEX idx_chat_members_user ON chat_members (user_id);

CREATE TABLE direct_chat_lookup (
    chat_id UUID PRIMARY KEY REFERENCES chats(id) ON DELETE CASCADE,
    user_a  UUID NOT NULL REFERENCES users(id),
    user_b  UUID NOT NULL REFERENCES users(id)
);

CREATE UNIQUE INDEX idx_direct_chat_unique_pair
    ON direct_chat_lookup (LEAST(user_a, user_b), GREATEST(user_a, user_b));

-- +goose Down
DROP TABLE IF EXISTS direct_chat_lookup;
DROP TABLE IF EXISTS chat_members;
DROP TABLE IF EXISTS chats;
DROP TYPE IF EXISTS chat_type;
