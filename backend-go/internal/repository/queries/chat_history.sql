-- name: CreateChatHistory :one
INSERT INTO chat_history (
  user_id,
  query,
  ai_response,
  retrieved_context
) VALUES (
  $1, $2, $3, $4
)
RETURNING id, user_id, query, ai_response, retrieved_context, created_at;

-- name: GetChatHistoryByUserID :many
SELECT id, user_id, query, ai_response, retrieved_context, created_at
FROM chat_history
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;
