-- +goose Up
CREATE TABLE channels (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug        VARCHAR(64)  NOT NULL UNIQUE,
    name        VARCHAR(128) NOT NULL,
    description TEXT         NOT NULL DEFAULT '',
    avatar_url  TEXT,
    owner_id    UUID         NOT NULL REFERENCES users(id),
    is_private  BOOLEAN      NOT NULL DEFAULT false,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX idx_channels_slug ON channels (slug);
CREATE INDEX idx_channels_owner ON channels (owner_id);

CREATE TABLE channel_subscribers (
    channel_id    UUID           NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    user_id       UUID           NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role          VARCHAR(32)    NOT NULL DEFAULT 'subscriber',
    subscribed_at TIMESTAMPTZ    NOT NULL DEFAULT now(),
    PRIMARY KEY (channel_id, user_id)
);

-- +goose Down
DROP TABLE IF EXISTS channel_subscribers;
DROP TABLE IF EXISTS channels;
