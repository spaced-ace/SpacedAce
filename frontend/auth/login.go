package auth

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"spaced-ace/constants"
)

type LoginForm struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

type LoginRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func PostLogin(c echo.Context) error {
	var loginForm = LoginForm{}
	if err := c.Bind(&loginForm); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	bodyMap := LoginRequestBody{
		Email:    loginForm.Email,
		Password: loginForm.Password,
	}
	bodyBytes, err := json.Marshal(bodyMap)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	bodyBuffer := bytes.NewBuffer(bodyBytes)

	resp, err := http.Post(constants.BACKEND_URL+"/authenticate-user", "application/json", bodyBuffer)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadGateway, err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return echo.NewHTTPError(resp.StatusCode, resp.Status)
	}

	var sessionCookie *http.Cookie
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "session" {
			sessionCookie = cookie
			break
		}
	}
	if sessionCookie == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "session cookie not found")
	}

	c.SetCookie(sessionCookie)
	c.Response().Header().Set("HX-Redirect", "/my-quizzes")
	return c.String(http.StatusOK, "login successful")
}
