package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"spaced-ace-backend/auth"
	"spaced-ace-backend/db"
)

type StartQuizSessionRequestBody struct {
	QuizId string `json:"quizId"`
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

	uri := "user=" + "test" + " dbname=" + "postgres" + " password=" + "test" + " host=" + "172.20.10.2" + " port=" + "5432" + " sslmode=disable"
	conn, err := pgx.Connect(context.Background(), uri)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	q := db.New(conn)

	quizSession, err := q.CreateQuizSession(
		context.Background(),
		db.CreateQuizSessionParams{
			ID:     uuid.NewString(),
			UserID: userId,
			QuizID: request.QuizId,
		},
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, quizSession)
}
