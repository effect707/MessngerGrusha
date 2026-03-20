-- name: CreateChat :exec
INSERT INTO chats (id, type, name, avatar_url, created_by, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: GetChatByID :one
SELECT id, type, name, avatar_url, created_by, created_at, updated_at
FROM chats
WHERE id = $1;

-- name: GetUserChats :many
SELECT c.id, c.type, c.name, c.avatar_url, c.created_by, c.created_at, c.updated_at
FROM chats c
JOIN chat_members cm ON c.id = cm.chat_id
WHERE cm.user_id = $1
ORDER BY c.updated_at DESC;

-- name: AddChatMember :exec
INSERT INTO chat_members (chat_id, user_id, role, joined_at)
VALUES ($1, $2, $3, $4);

-- name: RemoveChatMember :exec
DELETE FROM chat_members
WHERE chat_id = $1 AND user_id = $2;

-- name: IsChatMember :one
SELECT EXISTS(
    SELECT 1 FROM chat_members WHERE chat_id = $1 AND user_id = $2
) AS is_member;

-- name: GetChatMembers :many
SELECT chat_id, user_id, role, joined_at
FROM chat_members
WHERE chat_id = $1;

-- name: CreateDirectChatLookup :exec
INSERT INTO direct_chat_lookup (chat_id, user_a, user_b)
VALUES ($1, $2, $3);

-- name: GetDirectChat :one
SELECT chat_id
FROM direct_chat_lookup
WHERE LEAST(user_a, user_b) = LEAST($1::uuid, $2::uuid)
  AND GREATEST(user_a, user_b) = GREATEST($1::uuid, $2::uuid);
