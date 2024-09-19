-- name: GetEmailByID :one
SELECT sqlc.embed(email) FROM email
WHERE email.id = $1;

-- name: CreateEmail :one
INSERT INTO email (message_id, from_email, to_email, subject, content, sent_for) 
VALUES (sqlc.arg(message_id), sqlc.arg(from_email), sqlc.arg(to_email), sqlc.arg(subject), sqlc.arg(content), sqlc.arg(sent_for)) RETURNING *;
 