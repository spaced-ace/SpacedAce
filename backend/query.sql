-- name: GetQuizSession :one
    SELECT *
    FROM quiz_sessions
    WHERE true
        AND id = $1
    LIMIT 1;

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
        WHERE true
            AND quiz_id = $1
            AND user_id = $2
            AND finished_at IS NULL
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

-- Single Choice answers

-- name: CreateSingleChoiceAnswer :one
    INSERT INTO single_choice_answers (id, session_id, question_id, answer)
    VALUES ($1, $2, $3,  $4)
    RETURNING *;

-- name: GetSingleChoiceAnswers :many
    SELECT *
    FROM single_choice_answers
    WHERE true
        AND session_id = $1;

-- name: GetSingleChoiceAnswerById :one
    SELECT *
    FROM single_choice_answers
    WHERE true
        AND id = $1
    LIMIT 1;

-- name: GetSingleChoiceAnswerBySessionAndQuestionId :one
    SELECT *
    FROM single_choice_answers
    WHERE true
        AND session_id = $1
        AND question_id = $2
    LIMIT 1;

-- name: UpdateSingleChoiceAnswerBySessionAndQuestionId :one
    UPDATE single_choice_answers
    SET answer = $3
    WHERE true
        AND session_id = $1
        AND question_id = $2
    RETURNING *;

-- Multiple Choice answers

-- name: CreateMultipleChoiceAnswer :one
    INSERT INTO multiple_choice_answers (id, session_id, question_id, answers)
    VALUES ($1, $2, $3,  $4)
    RETURNING *;

-- name: GetMultipleChoiceAnswers :many
    SELECT *
    FROM multiple_choice_answers
    WHERE true
        AND session_id = $1;

-- name: GetMultipleChoiceAnswerById :one
    SELECT *
    FROM multiple_choice_answers
    WHERE true
        AND id = $1
    LIMIT 1;

-- name: GetMultipleChoiceAnswerBySessionAndQuestionId :one
    SELECT *
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
    RETURNING *;

-- True or false answers

-- name: CreateTrueOrFalseAnswer :one
    INSERT INTO true_or_false_answers (id, session_id, question_id, answer)
    VALUES ($1, $2, $3,  $4)
    RETURNING *;

-- name: GetTrueOrFalseAnswers :many
    SELECT *
    FROM true_or_false_answers
    WHERE true
        AND session_id = $1;

-- name: GetTrueOrFalseAnswerById :one
    SELECT *
    FROM true_or_false_answers
    WHERE true
        AND id = $1
    LIMIT 1;

-- name: GetTrueOrFalseAnswerBySessionAndQuestionId :one
    SELECT *
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

-- name: GetQuizResultByQuizSessionId :one
    SELECT *
    FROM quiz_results
    WHERE true
        AND session_id = $1
    LIMIT 1;

-- name: CreateQuizResult :one
    INSERT INTO quiz_results(id, session_id, max_score, score)
    VALUES ($1, $2, $3, $4)
    RETURNING *;

-- name: UpdateQuizResultScores :one
    UPDATE quiz_results
    SET max_score = $2, score = $3
    WHERE true
        AND id = $1
        RETURNING *;

-- name: GetAnswerScores :many
    SELECT *
    FROM answer_scores
    WHERE true
        AND quiz_result_id = $1;

-- name: CreateSingleChoiceAnswerScore :one
    INSERT INTO answer_scores(id, quiz_result_id, single_choice_answer_id, max_score, score)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING *;

-- name: CreateMultipleChoiceAnswerScore :one
    INSERT INTO answer_scores(id, quiz_result_id, multiple_choice_answer_id, max_score, score)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING *;

-- name: CreateTrueOrFalseAnswerScore :one
    INSERT INTO answer_scores(id, quiz_result_id, true_or_false_answer_id, max_score, score)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING *;
