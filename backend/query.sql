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

-- Multiple Choice answers

-- name: GetMultipleChoiceAnswers :many
    SELECT *
    FROM multiple_choice_answers
    WHERE true
        AND session_id = $1;

-- name: GetMultipleChoiceAnswerById :one
    SELECT 1
    FROM multiple_choice_answers
    WHERE true
        AND id = $1
    LIMIT 1;

-- name: GetMultipleChoiceAnswerBySessionAndQuestionId :one
    SELECT 1
    FROM multiple_choice_answers
    WHERE true
        AND session_id = $1
        AND question_id = $2
    LIMIT 1;

-- name: UpdateMultipleChoiceAnswerBySessionAndQuestionId :one
    UPDATE multiple_choice_answers
    SET answers = $3
    WHERE true
        AND session_id = $1
        AND question_id = $2
    RETURNING *;-- True or False answers

-- True or false answers

-- name: GetTrueOrFalseAnswers :many
    SELECT *
    FROM true_or_false_answers
    WHERE true
        AND session_id = $1;

-- name: GetTrueOrFalseAnswerById :one
    SELECT 1
    FROM true_or_false_answers
    WHERE true
        AND id = $1
    LIMIT 1;

-- name: GetTrueOrFalseAnswerBySessionAndQuestionId :one
    SELECT 1
    FROM true_or_false_answers
    WHERE true
        AND session_id = $1
        AND question_id = $2
    LIMIT 1;

-- name: UpdateTrueOrFalseAnswerBySessionAndQuestionId :one
    UPDATE true_or_false_answers
    SET answer = $3
    WHERE true
        AND session_id = $1
        AND question_id = $2
    RETURNING *;
