package pages

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"spaced-ace/api/models"
	"spaced-ace/constants"
	"spaced-ace/context"
	"strconv"
)

type quiz struct {
	models.QuizInfo
	FromColor string
	ToColor   string
}

type quizzesResponse struct {
	Quizzes []models.QuizInfo `json:"quizzes"`
	Length  int               `json:"length"`
}

type MyQuizzesPageData struct {
	Session *context.Session
	Quizzes []quiz
}

func MyQuizzesPage(c echo.Context) error {
	cc := c.(*context.Context)
	userId := cc.Session.User.Id

	req, _ := http.NewRequest("GET", constants.BACKEND_URL+"/quizzes/user/"+userId, nil)
	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: cc.Session.Id,
	})
	client := &http.Client{}

	resp, _ := client.Do(req)
	defer resp.Body.Close()
	var responseBody quizzesResponse
	_ = json.NewDecoder(resp.Body).Decode(&responseBody)

	var quizzes []quiz
	for _, q := range responseBody.Quizzes {
		fromColor, toColor := generateColors(q.Title, q.Id)

		quizzes = append(quizzes, quiz{
			QuizInfo:  q,
			FromColor: fromColor,
			ToColor:   toColor,
		})
	}

	pageData := MyQuizzesPageData{
		Session: cc.Session,
		Quizzes: quizzes,
	}

	return c.Render(200, "my-quizzes", pageData)
}

func DeleteQuiz(c echo.Context) error {
	cc := c.(*context.Context)
	quizId := c.Param("quizId")

	req, _ := http.NewRequest("DELETE", constants.BACKEND_URL+"/quizzes/"+quizId, nil)
	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: cc.Session.Id,
	})
	client := &http.Client{}

	resp, _ := client.Do(req)
	defer resp.Body.Close()

	return c.NoContent(http.StatusOK)
}

func hashToColor(input string) string {
	shaHash := sha256.New()
	shaHash.Write([]byte(input))
	hash := hex.EncodeToString(shaHash.Sum(nil))

	// Convert the first 6 characters of the hash into an integer
	colorInt, _ := strconv.ParseInt(hash[:6], 16, 64)

	// Map the integer to a color
	colors := []string{
		"red-300", "red-400", "red-500", "red-600",
		"orange-300", "orange-400", "orange-500", "orange-600",
		"amber-300", "amber-400", "amber-500", "amber-600",
		"yellow-300", "yellow-400", "yellow-500", "yellow-600",
		"green-300", "green-400", "green-500", "green-600",
		"blue-300", "blue-400", "blue-500", "blue-600",
		"purple-300", "purple-400", "purple-500", "purple-600",
		"pink-300", "pink-400", "pink-500", "pink-600",
		"emerald-300", "emerald-400", "emerald-500", "emerald-600",
		"teal-300", "teal-400", "teal-500", "teal-600",
		"cyan-300", "cyan-400", "cyan-500", "cyan-600",
		"indigo-300", "indigo-400", "indigo-500", "indigo-600",
		"violet-300", "violet-400", "violet-500", "violet-600",
		"fuchsia-300", "fuchsia-400", "fuchsia-500", "fuchsia-600",
		"rose-300", "rose-400", "rose-500", "rose-600",
	}
	color := colors[colorInt%int64(len(colors))]

	return color
}

func generateColors(title string, id string) (string, string) {
	fromColor := hashToColor(title + id)
	toColor := hashToColor(id)

	return fromColor, toColor
}
