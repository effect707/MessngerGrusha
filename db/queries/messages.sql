-- name: CreateMessage :exec
INSERT INTO messages (id, chat_id, channel_id, sender_id, type, content, reply_to_id, is_edited, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);

-- name: GetMessageByID :one
SELECT id, chat_id, channel_id, sender_id, type, content, reply_to_id, is_edited, created_at, updated_at
FROM messages
WHERE id = $1;

-- name: GetChatHistory :many
SELECT id, chat_id, channel_id, sender_id, type, content, reply_to_id, is_edited, created_at, updated_at
FROM messages
WHERE chat_id = $1
  AND (
    sqlc.narg('cursor_created_at')::timestamptz IS NULL
    OR created_at < sqlc.narg('cursor_created_at')::timestamptz
    OR (created_at = sqlc.narg('cursor_created_at')::timestamptz AND id < sqlc.narg('cursor_id')::uuid)
  )
ORDER BY created_at DESC, id DESC
LIMIT $2;

-- name: SearchMessages :many
SELECT id, chat_id, channel_id, sender_id, type, content, reply_to_id, is_edited, created_at, updated_at
FROM messages
WHERE chat_id = $1
  AND content_tsv @@ plainto_tsquery('russian', $2)
ORDER BY ts_rank(content_tsv, plainto_tsquery('russian', $2)) DESC
LIMIT $3;

-- name: UpdateMessage :exec
UPDATE messages
SET content = $2, is_edited = true, updated_at = $3
WHERE id = $1;

-- name: DeleteMessage :exec
DELETE FROM messages WHERE id = $1;
