-- name: CreateAttachment :exec
INSERT INTO attachments (id, message_id, file_name, file_size, mime_type, storage_key, duration_ms, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: GetAttachmentsByMessageID :many
SELECT id, message_id, file_name, file_size, mime_type, storage_key, duration_ms, created_at
FROM attachments
WHERE message_id = $1
ORDER BY created_at;

-- name: GetAttachmentByID :one
SELECT id, message_id, file_name, file_size, mime_type, storage_key, duration_ms, created_at
FROM attachments
WHERE id = $1;
