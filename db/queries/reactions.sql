-- name: AddReaction :exec
INSERT INTO reactions (message_id, user_id, emoji, created_at)
VALUES ($1, $2, $3, $4)
ON CONFLICT (message_id, user_id, emoji) DO NOTHING;

-- name: RemoveReaction :exec
DELETE FROM reactions
WHERE message_id = $1 AND user_id = $2 AND emoji = $3;

-- name: GetReactionsByMessageID :many
SELECT message_id, user_id, emoji, created_at
FROM reactions
WHERE message_id = $1
ORDER BY created_at;
