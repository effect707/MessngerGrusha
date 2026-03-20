-- name: CreateChannel :exec
INSERT INTO channels (id, slug, name, description, avatar_url, owner_id, is_private, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: GetChannelByID :one
SELECT id, slug, name, description, avatar_url, owner_id, is_private, created_at, updated_at
FROM channels
WHERE id = $1;

-- name: GetChannelBySlug :one
SELECT id, slug, name, description, avatar_url, owner_id, is_private, created_at, updated_at
FROM channels
WHERE slug = $1;

-- name: UpdateChannel :exec
UPDATE channels
SET name = $2, description = $3, avatar_url = $4, is_private = $5, updated_at = $6
WHERE id = $1;

-- name: DeleteChannel :exec
DELETE FROM channels WHERE id = $1;

-- name: GetPublicChannels :many
SELECT id, slug, name, description, avatar_url, owner_id, is_private, created_at, updated_at
FROM channels
WHERE is_private = false
ORDER BY created_at DESC
LIMIT $1;

-- name: GetUserChannels :many
SELECT c.id, c.slug, c.name, c.description, c.avatar_url, c.owner_id, c.is_private, c.created_at, c.updated_at
FROM channels c
JOIN channel_subscribers cs ON cs.channel_id = c.id
WHERE cs.user_id = $1
ORDER BY c.created_at DESC;

-- name: AddChannelSubscriber :exec
INSERT INTO channel_subscribers (channel_id, user_id, role, subscribed_at)
VALUES ($1, $2, $3, $4)
ON CONFLICT (channel_id, user_id) DO NOTHING;

-- name: RemoveChannelSubscriber :exec
DELETE FROM channel_subscribers
WHERE channel_id = $1 AND user_id = $2;

-- name: IsChannelSubscriber :one
SELECT EXISTS(
    SELECT 1 FROM channel_subscribers
    WHERE channel_id = $1 AND user_id = $2
) AS is_subscriber;

-- name: GetChannelSubscribers :many
SELECT channel_id, user_id, role, subscribed_at
FROM channel_subscribers
WHERE channel_id = $1
ORDER BY subscribed_at;

-- name: CountChannelSubscribers :one
SELECT COUNT(*) FROM channel_subscribers
WHERE channel_id = $1;
