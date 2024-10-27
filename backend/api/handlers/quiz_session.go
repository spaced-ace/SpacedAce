package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"spaced-ace-backend/api/models"
	"spaced-ace-backend/auth"
	"spaced-ace-backend/db"
	"spaced-ace-backend/question"
	"spaced-ace-backend/utils"
	"strings"
	"time"
)

type StartQuizSessionRequestBody struct {
	QuizId string `json:"quizId"`
}
type GetQuizSessionRequestBody struct {
	QuizId string `json:"quizId"`
	UserId string `json:"userId"`
}

func StartQuizSession(c echo.Context) error {
	var request StartQuizSessionRequestBody
	if err := json.NewDecoder(c.Request().Body).Decode(&request); err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	session, err := c.Cookie("session")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}
	userId, err := auth.GetUserIdBySession(session.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	sqlcQuerier := utils.GetQuerier()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// List quiz sessions
	openQuizSessions, err := sqlcQuerier.GetQuizSessionsByQuizIdAndUserId(
		ctx,
		db.GetQuizSessionsByQuizIdAndUserIdParams{
			QuizID: request.QuizId,
			UserID: userId,
		},
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error getting open quiz sessions: "+err.Error())
	}

	// Close the open quiz sessions
	for _, openQuizSession := range openQuizSessions {
		// Skip already closed sessions
		if openQuizSession.FinishedAt.Valid {
			continue
		}

		_, err := sqlcQuerier.UpdateQuizSessionFinishedAt(
			ctx,
			db.UpdateQuizSessionFinishedAtParams{
				ID: openQuizSession.ID,
				FinishedAt: pgtype.Timestamp{
					Time:             time.Now(),
					InfinityModifier: pgtype.Finite,
					Valid:            true,
				},
			},
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("error closing quiz session: %s", openQuizSession.ID))
		}
		log.Default().Printf("quiz session closed: %s", openQuizSession.ID)
	}

	// Start a new quiz session
	dbQuizSession, err := sqlcQuerier.CreateQuizSession(
		ctx,
		db.CreateQuizSessionParams{
			ID:     uuid.NewString(),
			UserID: userId,
			QuizID: request.QuizId,
		},
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error starting a new quiz session")
	}
	log.Default().Printf("quiz session started: %s", dbQuizSession.ID)

	quizSession, err := models.MapQuizSession(dbQuizSession)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, quizSession)
}

func GetQuizSessions(c echo.Context) error {
	userId := c.QueryParam("userId")
	quizId := c.QueryParam("quizId")
	open := c.QueryParam("open")

	sqlcQuerier := utils.GetQuerier()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbQuizSessions, err := sqlcQuerier.GetQuizSessionsByQuizIdAndUserId(
		ctx,
		db.GetQuizSessionsByQuizIdAndUserIdParams{
			QuizID: quizId,
			UserID: userId,
		},
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	filter := func(s db.QuizSession) bool {
		if open == "true" {
			return !s.FinishedAt.Valid
		} else if open == "false" {
			return s.FinishedAt.Valid
		}
		return true
	}

	resultCount := 0
	for _, dbQuizSession := range dbQuizSessions {
		if filter(dbQuizSession) {
			resultCount++
		}
	}

	quizSessions := make([]models.QuizSession, resultCount)
	nextIndex := 0
	for _, dbQuizSession := range dbQuizSessions {
		if !filter(dbQuizSession) {
			continue
		}

		quizSession, err := models.MapQuizSession(dbQuizSession)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		quizSessions[nextIndex] = *quizSession
		nextIndex++
	}

	response := models.GetQuizSessionsResponseBody{
		QuizSessions: quizSessions,
		Length:       resultCount,
	}
	return c.JSON(http.StatusOK, response)
}

func GetQuizSession(c echo.Context) error {
	quizSessionId := c.Param("quizSessionId")

	sqlcQuerier := utils.GetQuerier()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbQuizSession, err := sqlcQuerier.GetQuizSession(ctx, quizSessionId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	quizSession, err := models.MapQuizSession(dbQuizSession)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, quizSession)
}

func PostSubmitQuiz(c echo.Context) error {
	quizSessionId := c.Param("quizSessionId")

	session, err := c.Cookie("session")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}
	userId, err := auth.GetUserIdBySession(session.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	sqlcQuerier := utils.GetQuerier()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	quizSession, err := sqlcQuerier.GetQuizSession(
		ctx,
		quizSessionId,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("quiz session '%s' is not found", quizSessionId))
	}

	if quizSession.FinishedAt.Valid {
		return echo.NewHTTPError(http.StatusConflict, fmt.Sprintf("quiz session '%s' is already submitted", quizSessionId))
	}

	if quizSession.UserID != userId {
		return echo.NewHTTPError(http.StatusForbidden, "cannot submit other user's quiz session")
	}

	_, err = sqlcQuerier.UpdateQuizSessionFinishedAt(
		ctx,
		db.UpdateQuizSessionFinishedAtParams{
			ID: quizSessionId,
			FinishedAt: pgtype.Timestamp{
				Time:             time.Now(),
				InfinityModifier: pgtype.Finite,
				Valid:            true,
			},
		},
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("error submitting quiz session: %s", quizSessionId))
	}
	log.Default().Printf("quiz session submitted: %s", quizSessionId)

	dbQuizResult, err := sqlcQuerier.GetQuizResultByQuizSessionId(ctx, quizSessionId)
	if err != nil {
		quizResult, err := calculateAndStoreQuizResult(ctx, quizSessionId, quizSession.QuizID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "error calculating the quiz result: "+err.Error())
		}
		return c.JSON(http.StatusOK, quizResult)
	}

	quizResult, err := models.MapQuizResult(dbQuizResult)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("error parsing quiz result: %s", dbQuizResult.ID))
	}

	dbAnswerScores, err := sqlcQuerier.GetAnswerScores(ctx, dbQuizResult.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("error getting answer scores for the quiz result: %s", dbQuizResult.ID))
	}

	answerScores := make([]models.AnswerScore, len(dbAnswerScores))
	for i, dbScore := range dbAnswerScores {
		score, err := models.MapAnswerScore(dbScore)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("error parsing answer score: %s", dbScore.ID))
		}
		answerScores[i] = *score
	}

	quizResult.AnswerScores = answerScores
	return c.JSON(http.StatusOK, quizResult)
}

func GetQuizResult(c echo.Context) error {
	quizSessionId := c.Param("quizSessionId")

	sqlcQuerier := utils.GetQuerier()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	dbQuizResult, err := sqlcQuerier.GetQuizResultByQuizSessionId(ctx, quizSessionId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("quiz result is not found for ID `%s`, error: %s", quizSessionId, err.Error()))
	}
	quizResult, err := models.MapQuizResult(dbQuizResult)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("error mapping quiz result: %s", err.Error()))
	}

	dbAnswerScores, err := sqlcQuerier.GetAnswerScores(ctx, quizResult.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("error finding answer scores for quiz result with ID `%s`, error: %s", quizResult.ID, err))
	}
	answerScores := make([]models.AnswerScore, 0, len(dbAnswerScores))
	for _, dbAnswerScore := range dbAnswerScores {
		answerScore, err := models.MapAnswerScore(dbAnswerScore)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("error mapping answer score with ID `%s`", dbAnswerScore.ID))
		}
		answerScores = append(answerScores, *answerScore)
	}
	quizResult.AnswerScores = answerScores

	return c.JSON(http.StatusOK, quizResult)
}

func calculateAndStoreQuizResult(ctx context.Context, sessionID, quizID string) (*models.QuizResult, error) {
	sqlcQuerier := utils.GetQuerier()

	// Create a quiz result record with initial scores
	dbQuizResult, err := sqlcQuerier.CreateQuizResult(
		ctx,
		db.CreateQuizResultParams{
			ID:        uuid.NewString(),
			SessionID: sessionID,
			MaxScore:  0.0,
			Score:     0.0,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error creating quiz result: %w", err)
	}

	// Calculate the scores for the single choice questions
	singleChoiceAnswerScores, err := calculateSingleChoiceQuestionScores(ctx, dbQuizResult.ID, sessionID, quizID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate single choice scores: %w", err)
	}

	// Calculate the scores for the multiple choice questions
	multipleChoiceAnswerScores, err := calculateMultipleChoiceQuestionScores(ctx, dbQuizResult.ID, sessionID, quizID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate multiple choice scores: %w", err)
	}

	// Calculate the scores for the true or false questions
	trueOrFalseAnswerScores, err := calculateTrueOrFalseQuestionScores(ctx, dbQuizResult.ID, sessionID, quizID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate true or false scores: %w", err)
	}

	// Collect all the answer scores
	answerScores := make([]models.AnswerScore, 0, len(singleChoiceAnswerScores)+len(multipleChoiceAnswerScores)+len(trueOrFalseAnswerScores))
	answerScores = append(answerScores, singleChoiceAnswerScores...)
	answerScores = append(answerScores, multipleChoiceAnswerScores...)
	answerScores = append(answerScores, trueOrFalseAnswerScores...)

	// Calculate the score and the max score
	var maxScore, score float64
	for _, answerScore := range answerScores {
		maxScore += answerScore.MaxScore
		score += answerScore.Score
	}

	// Update quiz result with the calculated scores
	updatedDbQuizResult, err := sqlcQuerier.UpdateQuizResultScores(
		ctx,
		db.UpdateQuizResultScoresParams{
			ID:       dbQuizResult.ID,
			MaxScore: maxScore,
			Score:    score,
		},
	)

	// Map the results
	quizResult, err := models.MapQuizResult(updatedDbQuizResult)
	if err != nil {
		return nil, fmt.Errorf("failed to update quiz result scores: %w", err)
	}
	quizResult.AnswerScores = answerScores

	return quizResult, nil
}

func calculateSingleChoiceQuestionScores(ctx context.Context, quizResultID, sessionID, quizID string) ([]models.AnswerScore, error) {
	sqlcQuerier := utils.GetQuerier()

	questions, err := question.GetSingleChoiceQuestions(quizID)
	if err != nil {
		return []models.AnswerScore{}, err
	}

	answers, err := sqlcQuerier.GetSingleChoiceAnswers(ctx, sessionID)
	if err != nil {
		return []models.AnswerScore{}, err
	}

	dbAnswerScores := make([]db.AnswerScore, len(questions))

	for _, q := range questions {
		singleChoiceQuestion := models.SingleChoiceQuestion{
			ID:            q.UUID,
			QuizID:        q.QuizID,
			QuestionType:  models.SingleChoice,
			Question:      q.Question,
			Answers:       q.Answers,
			CorrectAnswer: q.CorrectAnswer,
		}

		userAnswer, err := findSingleChoiceAnswer(answers, q.UUID)
		if err != nil {
			return []models.AnswerScore{}, err
		}

		score := 0.0
		for _, a := range singleChoiceQuestion.Answers {
			if a == userAnswer.Answer {
				score += 1
				break
			}
		}

		dbAnswerScore, err := sqlcQuerier.CreateSingleChoiceAnswerScore(
			ctx,
			db.CreateSingleChoiceAnswerScoreParams{
				ID:                   uuid.NewString(),
				QuizResultID:         quizResultID,
				SingleChoiceAnswerID: userAnswer.ID,
				MaxScore:             1,
				Score:                score,
			},
		)
		if err != nil {
			return nil, err
		}
		dbAnswerScores = append(dbAnswerScores, dbAnswerScore)
	}

	answerScores := make([]models.AnswerScore, len(dbAnswerScores))
	for i, dbs := range dbAnswerScores {
		answerScore, err := models.MapAnswerScore(dbs)
		if err != nil {
			return []models.AnswerScore{}, err
		}
		answerScores[i] = *answerScore
	}

	return answerScores, nil
}
func calculateMultipleChoiceQuestionScores(ctx context.Context, quizResultID, sessionID, quizID string) ([]models.AnswerScore, error) {
	sqlcQuerier := utils.GetQuerier()

	questions, err := question.GetMultipleChoiceQuestions(quizID)
	if err != nil {
		return []models.AnswerScore{}, err
	}

	answers, err := sqlcQuerier.GetMultipleChoiceAnswers(ctx, sessionID)
	if err != nil {
		return []models.AnswerScore{}, err
	}

	dbAnswerScores := make([]db.AnswerScore, len(questions))

	for _, q := range questions {
		multipleChoiceQuestion := models.MultipleChoiceQuestion{
			ID:             q.UUID,
			QuizID:         q.QuizID,
			QuestionType:   models.MultipleChoice,
			Question:       q.Question,
			Answers:        q.Answers,
			CorrectAnswers: q.CorrectAnswers,
		}

		userAnswer, err := findMultipleChoiceAnswer(answers, q.UUID)
		if err != nil {
			return []models.AnswerScore{}, err
		}

		positiveScore := 1.0 / float64(len(multipleChoiceQuestion.Answers))
		negativeScore := 1.0 / float64(4-len(multipleChoiceQuestion.Answers))

		score := 0.0
		for _, correctAnswer := range multipleChoiceQuestion.Answers {
			if strings.Contains(userAnswer.Answers, correctAnswer) {
				score += positiveScore
			} else {
				score -= negativeScore
			}
		}
		if score <= 0.0 {
			score = 0
		}

		dbAnswerScore, err := sqlcQuerier.CreateSingleChoiceAnswerScore(
			ctx,
			db.CreateSingleChoiceAnswerScoreParams{
				ID:                   uuid.NewString(),
				QuizResultID:         quizResultID,
				SingleChoiceAnswerID: userAnswer.ID,
				MaxScore:             1,
				Score:                score,
			},
		)
		if err != nil {
			return nil, err
		}
		dbAnswerScores = append(dbAnswerScores, dbAnswerScore)
	}

	answerScores := make([]models.AnswerScore, len(dbAnswerScores))
	for i, dbs := range dbAnswerScores {
		answerScore, err := models.MapAnswerScore(dbs)
		if err != nil {
			return []models.AnswerScore{}, err
		}
		answerScores[i] = *answerScore
	}

	return answerScores, nil
}
func calculateTrueOrFalseQuestionScores(ctx context.Context, quizResultID, sessionID, quizID string) ([]models.AnswerScore, error) {
	sqlcQuerier := utils.GetQuerier()

	questions, err := question.GetTrueOrFalseQuestions(quizID)
	if err != nil {
		return []models.AnswerScore{}, err
	}

	answers, err := sqlcQuerier.GetTrueOrFalseAnswers(ctx, sessionID)
	if err != nil {
		return []models.AnswerScore{}, err
	}

	dbAnswerScores := make([]db.AnswerScore, len(questions))

	for _, q := range questions {
		trueOrFalseQuestion := models.TrueOrFalseQuestion{
			ID:            q.UUID,
			QuizID:        q.QuizID,
			QuestionType:  models.TrueOrFalse,
			Question:      q.Question,
			CorrectAnswer: q.CorrectAnswer,
		}

		userAnswer, err := findTrueOrFalseAnswer(answers, q.UUID)
		if err != nil {
			return []models.AnswerScore{}, err
		}

		score := 0.0
		if trueOrFalseQuestion.CorrectAnswer == userAnswer.Answer {
			score = 1.0
		}

		dbAnswerScore, err := sqlcQuerier.CreateSingleChoiceAnswerScore(
			ctx,
			db.CreateSingleChoiceAnswerScoreParams{
				ID:                   uuid.NewString(),
				QuizResultID:         quizResultID,
				SingleChoiceAnswerID: userAnswer.ID,
				MaxScore:             1,
				Score:                score,
			},
		)
		if err != nil {
			return nil, err
		}
		dbAnswerScores = append(dbAnswerScores, dbAnswerScore)
	}

	answerScores := make([]models.AnswerScore, len(dbAnswerScores))
	for i, dbs := range dbAnswerScores {
		answerScore, err := models.MapAnswerScore(dbs)
		if err != nil {
			return []models.AnswerScore{}, err
		}
		answerScores[i] = *answerScore
	}

	return answerScores, nil
}

func findSingleChoiceAnswer(answers []db.SingleChoiceAnswer, questionID string) (*models.SingleChoiceAnswer, error) {
	for _, dbAnswer := range answers {
		if dbAnswer.QuestionID == questionID {
			answer, err := models.MapSingleChoiceAnswer(dbAnswer)
			if err != nil {
				return nil, fmt.Errorf("error mapping answer with ID `%s`", dbAnswer.ID)
			}
			return answer, err
		}
	}
	return nil, fmt.Errorf("answer not found for question with ID `%s`", questionID)
}
func findMultipleChoiceAnswer(answers []db.MultipleChoiceAnswer, questionID string) (*models.MultipleChoiceAnswer, error) {
	for _, dbAnswer := range answers {
		if dbAnswer.QuestionID == questionID {
			answer, err := models.MapMultipleChoiceAnswer(dbAnswer)
			if err != nil {
				return nil, fmt.Errorf("error mapping answer with ID `%s`", dbAnswer.ID)
			}
			return answer, err
		}
	}
	return nil, fmt.Errorf("answer not found for question with ID `%s`", questionID)
}
func findTrueOrFalseAnswer(answers []db.TrueOrFalseAnswer, questionID string) (*models.TrueOrFalseAnswer, error) {
	for _, dbAnswer := range answers {
		if dbAnswer.QuestionID == questionID {
			answer, err := models.MapTrueOrFalseAnswer(dbAnswer)
			if err != nil {
				return nil, fmt.Errorf("error mapping answer with ID `%s`", dbAnswer.ID)
			}
			return answer, err
		}
	}
	return nil, fmt.Errorf("answer not found for question with ID `%s`", questionID)
}

func HasOpenQuizSession(c echo.Context) error {
	userId := c.QueryParam("userId")
	quizId := c.QueryParam("quizId")

	sqlcQuerier := utils.GetQuerier()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	hasOpenSession, err := sqlcQuerier.HasOpenQuizSession(
		ctx,
		db.HasOpenQuizSessionParams{
			QuizID: quizId,
			UserID: userId,
		},
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if hasOpenSession {
		return c.NoContent(http.StatusOK)
	} else {
		return c.NoContent(http.StatusNotFound)
	}
}
