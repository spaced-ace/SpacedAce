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
    VALUES ($1, $2, $3, $4, NULL, $5)
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

-- name: GetQuizResultsByUserID :many
    SELECT *
    FROM quiz_results
    INNER JOIN quiz_sessions qs on qs.id = quiz_results.session_id
    WHERE true
      AND qs.user_id = $1;

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

-- name: GetAddedLearnListItems :many
    SELECT *
    FROM learn_list_added_items
    WHERE true
        AND user_id = $1;

-- name: AddQuizToLearnList :exec
    INSERT INTO learn_list_added_items(user_id, quiz_id)
    VALUES ($1, $2)
    ON CONFLICT (user_id, quiz_id) DO NOTHING;

-- name: RemoveQuizFromLearnList :exec
    DELETE FROM learn_list_added_items
    WHERE true
        AND user_id = $1
        AND quiz_id = $2;

-- name: GetReviewItems :many
    SELECT
        review_items.*,
        Q.name::text AS quiz_name,
        Q.id::text AS quiz_id,
        CASE
            WHEN SQC.question IS NOT NULL THEN SQC.question::text
            WHEN MQC.question IS NOT NULL THEN MQC.question::text
            ELSE TQC.question::text
            END AS question_name
    FROM review_items
    LEFT JOIN single_choice_questions SQC ON review_items.single_choice_question_id = SQC.uuid
    LEFT JOIN multiple_choice_questions MQC ON review_items.multiple_choice_question_id = MQC.uuid
    LEFT JOIN true_or_false_questions TQC ON review_items.true_or_false_question_id = TQC.uuid
    LEFT JOIN quizzes Q ON ( false
        OR q.id = SQC.quizid
        OR q.id = MQC.quizid
        OR q.id = TQC.quizid
    )
    WHERE true
      AND user_id = $1;

-- name: GetReviewItem :one
    SELECT
        review_items.*,
        Q.name::text AS quiz_name,
        Q.id::text AS quiz_id,
        CASE
            WHEN SQC.question IS NOT NULL THEN SQC.question::text
            WHEN MQC.question IS NOT NULL THEN MQC.question::text
            ELSE TQC.question::text
            END AS question_name
    FROM review_items
    LEFT JOIN single_choice_questions SQC ON review_items.single_choice_question_id = SQC.uuid
    LEFT JOIN multiple_choice_questions MQC ON review_items.multiple_choice_question_id = MQC.uuid
    LEFT JOIN true_or_false_questions TQC ON review_items.true_or_false_question_id = TQC.uuid
    LEFT JOIN quizzes Q ON ( false
       OR q.id = SQC.quizid
       OR q.id = MQC.quizid
       OR q.id = TQC.quizid
    )
    WHERE review_items.id = $1;

-- name: CreateSingleChoiceReviewItem :one
    INSERT INTO review_items(id, user_id, single_choice_question_id, ease_factor, difficulty, streak, next_review_date, interval_in_minutes)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    RETURNING id;

-- name: CreateMultipleChoiceReviewItem :one
    INSERT INTO review_items(id, user_id, multiple_choice_question_id, ease_factor, difficulty, streak, next_review_date, interval_in_minutes)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    RETURNING id;
-- name: CreateTrueOrFalseReviewItem :one
    INSERT INTO review_items(id, user_id, true_or_false_question_id, ease_factor, difficulty, streak, next_review_date, interval_in_minutes)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    RETURNING id;

-- name: DeleteReviewItem :exec
    DELETE FROM review_items
    WHERE true
        AND id = $1;

-- name: DeleteReviewItemsByQuizID :exec
    DELETE FROM review_items
        USING single_choice_questions, multiple_choice_questions, true_or_false_questions, quizzes
        WHERE true
            AND user_id = $1
            AND (
                review_items.single_choice_question_id = single_choice_questions.uuid
                    AND single_choice_questions.quizid = quizzes.id
                    AND quizzes.id = $2
                )
                OR (
                review_items.multiple_choice_question_id = multiple_choice_questions.uuid
                    AND multiple_choice_questions.quizid = quizzes.id
                    AND quizzes.id = $2
                ) OR (
                review_items.true_or_false_question_id = true_or_false_questions.uuid
                    AND true_or_false_questions.quizid = quizzes.id
                    AND quizzes.id = $2
                );

-- name: UpdateReviewItem :exec
    UPDATE review_items
    SET ease_factor = $2, difficulty = $3, streak = $4, next_review_date = $5, interval_in_minutes = $6
    WHERE true
        AND id = $1;

-- name: GetQuizOptions :many
    SELECT
        Q.id as quiz_id,
        Q.name as quiz_name
    FROM quizzes Q
    INNER JOIN quiz_accesses A ON A.quizid = Q.id
    WHERE true
        AND A.userid = $1;