-- +goose Up
CREATE TYPE message_type AS ENUM ('text', 'image', 'file', 'voice', 'system');

CREATE TABLE messages (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    chat_id     UUID REFERENCES chats(id) ON DELETE CASCADE,
    channel_id  UUID,
    sender_id   UUID         NOT NULL REFERENCES users(id),
    type        message_type NOT NULL DEFAULT 'text',
    content     TEXT         NOT NULL DEFAULT '',
    reply_to_id UUID         REFERENCES messages(id) ON DELETE SET NULL,
    is_edited   BOOLEAN      NOT NULL DEFAULT false,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),

    CONSTRAINT chk_message_target CHECK (
        (chat_id IS NOT NULL AND channel_id IS NULL) OR
        (chat_id IS NULL AND channel_id IS NOT NULL)
    ),
    CONSTRAINT chk_text_not_empty CHECK (
        type != 'text' OR (content IS NOT NULL AND content != '')
    )
);

ALTER TABLE messages ADD COLUMN content_tsv tsvector
    GENERATED ALWAYS AS (to_tsvector('russian', content)) STORED;

CREATE INDEX idx_messages_chat_created ON messages (chat_id, created_at DESC)
    WHERE chat_id IS NOT NULL;

CREATE INDEX idx_messages_channel_created ON messages (channel_id, created_at DESC)
    WHERE channel_id IS NOT NULL;

CREATE INDEX idx_messages_search ON messages USING GIN (content_tsv);

CREATE TABLE attachments (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    message_id  UUID         NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
    file_name   VARCHAR(255) NOT NULL,
    file_size   BIGINT       NOT NULL,
    mime_type   VARCHAR(127) NOT NULL,
    storage_key TEXT         NOT NULL,
    duration_ms INTEGER,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX idx_attachments_message ON attachments (message_id);

-- +goose Down
DROP TABLE IF EXISTS attachments;
DROP TABLE IF EXISTS messages;
DROP TYPE IF EXISTS message_type;
