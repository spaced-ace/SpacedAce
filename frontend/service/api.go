package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"spaced-ace/constants"
	"spaced-ace/models"
)

type ApiService struct {
	sessionCookie *http.Cookie
	client        *http.Client
}

func NewApiService(sessionCookie *http.Cookie) *ApiService {
	return &ApiService{
		sessionCookie: sessionCookie,
		client:        &http.Client{},
	}
}

func (a *ApiService) getResponse(method, path string, requestBody any, responseBody interface{}) error {
	data, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, constants.BACKEND_URL+path, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	if a.sessionCookie != nil {
		req.AddCookie(a.sessionCookie)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var result map[string]string
		err = json.Unmarshal(bodyBytes, &result)
		if err != nil {
			return err
		}

		return errors.New("error message: " + result["message"])
	}

	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		return err
	}

	return nil
}

func (a *ApiService) GetSession() (*models.Session, error) {
	session := new(models.Session)
	if err := a.getResponse("GET", "/authenticated", nil, session); err != nil {
		return nil, err
	}
	return session, nil
}

func (a *ApiService) DeleteSession() error {
	return a.getResponse("POST", "/logout", nil, nil)
}

func (a *ApiService) GetQuiz(quizId string) (*models.Quiz, error) {
	quiz := new(models.Quiz)
	if err := a.getResponse("GET", "/quizzes/"+quizId, nil, quiz); err != nil {
		return nil, err
	}
	return quiz, nil
}
