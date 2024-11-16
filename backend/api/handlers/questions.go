package handlers

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"spaced-ace-backend/api/models"
	"spaced-ace-backend/auth"
	"spaced-ace-backend/constants"
	"spaced-ace-backend/question"
	"spaced-ace-backend/quiz"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// ##### responses from the llm api
type multipleChoiceResponse struct {
	Question       string   `json:"question"`
	Options        []string `json:"options"`
	CorrectOptions []string `json:"correct_options"`
}

type singleChoiceResponse struct {
	Question      string   `json:"question"`
	Options       []string `json:"options"`
	CorrectOption string   `json:"correct_option"`
}

type trueOrFalseResponse struct {
	Question      string `json:"question"`
	CorrectAnswer bool   `json:"correct_option"`
}

// #####

type textChunk struct {
	Id   string `json:"id"`
	Text string `json:"chunk"`
}

type prompt struct {
	Prompt string `json:"prompt"`
}

type cacheEntry struct {
	chunks        *[]textChunk
	IndexLastUsed int
}

type quizAccess struct {
	userId string
	quizId string
	access int
}

var cache = make(map[string]cacheEntry)

func CreateMultipleChoiceQuestionEndpoint(c echo.Context) error {
	var request = models.QuestionCreationRequestBody{}
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		fmt.Println(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}
	quizAccess, err := accessControlQuiz(c, request.QuizId)
	if quizAccess.access != 1 {
		return c.JSON(http.StatusUnauthorized, "unauthorized")
	}

	chunkToUse, err := manageChunking(request.Prompt)
	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, "error chunking prompt")
	}
	promptJson, err := json.Marshal(prompt{Prompt: chunkToUse.Text})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "error marshalling prompt")
	}
	res, err := http.Post(
		constants.LLM_API_URL+"/multiple-choice/create",
		"application/json",
		bytes.NewBuffer(promptJson),
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "error during question generation")
	}
	defer res.Body.Close()
	generated := multipleChoiceResponse{}
	err = json.NewDecoder(res.Body).Decode(&generated)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "error decoding response")
	}

	dbQuestion := question.DBMultipleChoiceQuestion{
		UUID:           uuid.New().String(),
		QuizID:         quizAccess.quizId,
		Question:       generated.Question,
		Answers:        generated.Options,
		CorrectAnswers: generated.CorrectOptions,
	}
	err = question.CreateMultipleChoiceQuestion(&dbQuestion)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	result := models.MultipleChoiceQuestion{
		ID:             dbQuestion.UUID,
		QuizID:         quizAccess.quizId,
		QuestionType:   models.MultipleChoice,
		Question:       generated.Question,
		Answers:        generated.Options,
		CorrectAnswers: generated.CorrectOptions,
	}
	return c.JSON(http.StatusOK, &result)
}

func CreateSingleChoiceQuestionEndpoint(c echo.Context) error {
	var request = models.QuestionCreationRequestBody{}
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}
	quizAccess, err := accessControlQuiz(c, request.QuizId)
	if quizAccess.access != 1 {
		return c.JSON(http.StatusUnauthorized, "unauthorized")
	}

	chunkToUse, err := manageChunking(request.Prompt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "error chunking prompt")
	}
	promptJson, err := json.Marshal(prompt{Prompt: chunkToUse.Text})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "error marshalling prompt")
	}
	res, err := http.Post(
		constants.LLM_API_URL+"/single-choice/create",
		"application/json",
		bytes.NewBuffer(promptJson),
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "error during question generation")
	}
	defer res.Body.Close()
	generated := singleChoiceResponse{}
	err = json.NewDecoder(res.Body).Decode(&generated)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "error decoding response")
	}

	dbQuestion := question.DBSingleChoiceQuestion{
		UUID:          uuid.New().String(),
		QuizID:        quizAccess.quizId,
		Question:      generated.Question,
		Answers:       generated.Options,
		CorrectAnswer: generated.CorrectOption,
	}
	err = question.CreateSingleChoiceQuestion(&dbQuestion)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	result := models.SingleChoiceQuestion{
		ID:            dbQuestion.UUID,
		QuizID:        quizAccess.quizId,
		QuestionType:  models.SingleChoice,
		Question:      generated.Question,
		Answers:       generated.Options,
		CorrectAnswer: generated.CorrectOption,
	}
	return c.JSON(http.StatusOK, &result)
}

func CreateTrueOrFalseQuestionEndpoint(c echo.Context) error {
	var request = models.QuestionCreationRequestBody{}
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}
	quizAccess, err := accessControlQuiz(c, request.QuizId)
	if quizAccess.access != 1 {
		return c.JSON(http.StatusUnauthorized, "unauthorized")
	}

	chunkToUse, err := manageChunking(request.Prompt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "error chunking prompt")
	}
	promptJson, err := json.Marshal(prompt{Prompt: chunkToUse.Text})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "error marshalling prompt")
	}
	res, err := http.Post(
		constants.LLM_API_URL+"/true-or-false/create",
		"application/json",
		bytes.NewBuffer(promptJson),
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "error during question generation")
	}
	defer res.Body.Close()
	generated := trueOrFalseResponse{}
	err = json.NewDecoder(res.Body).Decode(&generated)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "error decoding response")
	}

	dbQuestion := question.DBTrueOrFalseQuestion{
		UUID:          uuid.New().String(),
		QuizID:        quizAccess.quizId,
		Question:      generated.Question,
		CorrectAnswer: generated.CorrectAnswer,
	}
	err = question.CreateTrueOrFalseQuestion(&dbQuestion)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	result := models.TrueOrFalseQuestion{
		ID:            dbQuestion.UUID,
		QuizID:        quizAccess.quizId,
		QuestionType:  models.TrueOrFalse,
		Question:      generated.Question,
		CorrectAnswer: generated.CorrectAnswer,
	}
	return c.JSON(http.StatusOK, &result)
}

func GetMultipleChoiceEndpoint(c echo.Context) error {
	_, err := c.Cookie("session")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "unauthorized")
	}
	questionId := c.Param("id")
	q, err := question.GetMultipleChoiceQuestion(questionId)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, "question not found")
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	result := q.MapToModel()
	return c.JSON(http.StatusOK, &result)
}

func GetSingleChoiceEndpoint(c echo.Context) error {
	_, err := c.Cookie("session")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "unauthorized")
	}
	questionId := c.Param("id")
	q, err := question.GetSingleChoiceQuestion(questionId)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, "question not found")
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	result := q.MapToModel()
	return c.JSON(http.StatusOK, result)
}

func GetTrueOrFalseEndpoint(c echo.Context) error {
	_, err := c.Cookie("session")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "unauthorized")
	}
	questionId := c.Param("id")
	q, err := question.GetTrueOrFalseQuestion(questionId)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, "question not found")
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	result := q.MapToModel()
	return c.JSON(http.StatusOK, result)
}

func UpdateMultipleChoiceQuestionEndpoint(c echo.Context) error {
	questionId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid id")
	}
	request := models.MultipleChoiceUpdateRequestBody{}
	err = json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "bad request")
	}
	access, err := accessControlQuiz(c, request.QuizId)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "unauthorized")
	}
	if access.access != 1 {
		return c.JSON(http.StatusUnauthorized, "unauthorized")
	}

	questionToUpdate, err := question.GetMultipleChoiceQuestion(questionId.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, "question not found")
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if request.Question != "" {
		questionToUpdate.Question = request.Question
	}
	if len(request.Answers) > 0 {
		questionToUpdate.Answers = request.Answers
	}
	if len(request.CorrectAnswers) > 0 {
		questionToUpdate.CorrectAnswers = request.CorrectAnswers
	}
	fmt.Println(questionToUpdate)
	err = question.UpdateMultipleChoiceQuestion(&questionToUpdate)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	result := questionToUpdate.MapToModel()
	return c.JSON(http.StatusOK, result)
}

func UpdateSingleChoiceQuestionEndpoint(c echo.Context) error {
	questionId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid id")
	}
	request := models.SingleChoiceUpdateRequestBody{}
	err = json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "bad request")
	}
	access, err := accessControlQuiz(c, request.QuizId)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "unauthorized")
	}
	if access.access != 1 {
		return c.JSON(http.StatusUnauthorized, "unauthorized")
	}

	questionToUpdate, err := question.GetSingleChoiceQuestion(questionId.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, "question not found")
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if request.Question != "" {
		questionToUpdate.Question = request.Question
	}
	if len(request.Answers) > 0 {
		questionToUpdate.Answers = request.Answers
	}
	if request.CorrectAnswer != "" {
		questionToUpdate.CorrectAnswer = request.CorrectAnswer
	}
	err = question.UpdateSingleChoiceQuestion(&questionToUpdate)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	result := questionToUpdate.MapToModel()
	return c.JSON(http.StatusOK, &result)
}

func UpdateTrueOrFalseQuestionEndpoint(c echo.Context) error {
	questionId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid id")
	}
	request := models.TrueOrFalseUpdateRequestBody{}
	err = json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "bad request")
	}
	access, err := accessControlQuiz(c, request.QuizId)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "unauthorized")
	}
	if access.access != 1 {
		return c.JSON(http.StatusUnauthorized, "unauthorized")
	}

	questionToUpdate, err := question.GetTrueOrFalseQuestion(questionId.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, "question not found")
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if request.Question != "" {
		questionToUpdate.Question = request.Question
	}
	questionToUpdate.CorrectAnswer = request.CorrectAnswer
	err = question.UpdateTrueOrFalseQuestion(&questionToUpdate)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	result := questionToUpdate.MapToModel()
	return c.JSON(http.StatusOK, result)
}

func DeleteMultipleChoiceQuestionEndpoint(c echo.Context) error {
	questionId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid question id")
	}
	quizId, err := uuid.Parse(c.Param("quizId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid quiz id")
	}
	access, err := accessControlQuiz(c, quizId.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusUnauthorized, "unauthorized")
		}
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	if access.access != 1 {
		return c.JSON(http.StatusUnauthorized, "unauthorized")
	}
	err = question.DeleteMultipleChoiceQuestion(questionId.String())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "question deleted")
}

func DeleteSingleChoiceQuestionEndpoint(c echo.Context) error {
	questionId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid question id")
	}
	quizId, err := uuid.Parse(c.Param("quizId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid quiz id")
	}
	access, err := accessControlQuiz(c, quizId.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusUnauthorized, "unauthorized")
		}
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	if access.access != 1 {
		return c.JSON(http.StatusUnauthorized, "unauthorized")
	}
	err = question.DeleteSingleChoiceQuestion(questionId.String())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "question deleted")
}

func DeleteTrueOrFalseQuestionEndpoint(c echo.Context) error {
	questionId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid question id")
	}
	quizId, err := uuid.Parse(c.Param("quizId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid quiz id")
	}
	access, err := accessControlQuiz(c, quizId.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusUnauthorized, "unauthorized")
		}
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	if access.access != 1 {
		return c.JSON(http.StatusUnauthorized, "unauthorized")
	}
	err = question.DeleteTrueOrFalseQuestion(questionId.String())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "question deleted")
}

func accessControlQuiz(c echo.Context, quizId string) (*quizAccess, error) {
	_, err := uuid.Parse(quizId)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	session, err := c.Cookie("session")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	userId, err := auth.GetUserIdBySession(session.Value)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	access, err := quiz.GetQuizAccess(userId, quizId)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &quizAccess{userId: userId, quizId: quizId, access: access}, nil
}

func manageChunking(userPrompt string) (*textChunk, error) {
	promptLength := len(userPrompt)
	if promptLength == 0 || promptLength > 100_000 {
		return nil, errors.New("prompt must be between 1 and 100,000 characters")
	}
	hash := hashPrompt(userPrompt)
	existingCacheEntry, ok := cache[hash]
	if !ok {
		chunks, err := chunkPrompt(userPrompt)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		existingCacheEntry = cacheEntry{
			chunks:        chunks,
			IndexLastUsed: -1,
		}
		cache[hash] = existingCacheEntry
	}
	var chunkToUse textChunk
	if existingCacheEntry.IndexLastUsed <= len(*existingCacheEntry.chunks)-1 {
		chunkToUse = (*existingCacheEntry.chunks)[existingCacheEntry.IndexLastUsed+1]
	} else {
		chunkToUse = (*existingCacheEntry.chunks)[0]
	}
	existingCacheEntry.IndexLastUsed++
	return &chunkToUse, nil
}

func hashPrompt(userPrompt string) string {
	hash := sha256.New()
	hash.Write([]byte(userPrompt))
	return string(hash.Sum(nil))
}

func chunkPrompt(userPrompt string) (*[]textChunk, error) {
	promptRequest := prompt{Prompt: userPrompt}
	bodyjson, err := json.Marshal(promptRequest)
	if err != nil {
		return nil, err
	}
	bodyBuffer := bytes.NewBuffer(bodyjson)

	resp, err := http.Post(constants.LLM_API_URL+"/chunk", "application/json", bodyBuffer)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	result := []textChunk{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
