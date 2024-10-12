-- name: GetQuizSession :one
    SELECT * FROM quiz_sessions
    WHERE id = $1 LIMIT 1;

-- name: GetQuizSessionsByQuizIdAndUserId :many
    SELECT * FROM quiz_sessions
    WHERE quiz_id = $1
        AND user_id = $2;

-- name: GetQuizSessionsByUserId :many
    SELECT * FROM quiz_sessions
    WHERE user_id = $1;

-- name: HasOpenQuizSession :one
SELECT EXISTS (
    SELECT 1
    FROM quiz_sessions
    WHERE quiz_id = $1
        AND user_id = $2
        AND finished_at IS NOT NULL
) AS exists;

-- name: CreateQuizSession :one
    INSERT INTO quiz_sessions (id, user_id, quiz_id, started_at, finished_at, closes_at)
    VALUES ($1, $2, $3, NOW(), NULL, $4)
    RETURNING *;

-- name: UpdateQuizSessionFinishedAt :one
    UPDATE quiz_sessions
    SET finished_at = $2
    WHERE id = $1
    RETURNING *;

-- -- name: GetAuthor :one
-- SELECT * FROM authors
-- WHERE id = $1 LIMIT 1;
--
-- -- name: ListAuthors :many
-- SELECT * FROM authors
-- ORDER BY name;
--
-- -- name: CreateAuthor :one
-- INSERT INTO authors (
--     name, bio
-- ) VALUES (
--              $1, $2
--          )
-- RETURNING *;
--
-- -- name: UpdateAuthor :exec
-- UPDATE authors
-- set name = $2,
--     bio = $3
-- WHERE id = $1;
--
-- -- name: DeleteAuthor :exec
-- DELETE FROM authors
-- WHERE id = $1;