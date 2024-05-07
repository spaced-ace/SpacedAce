package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"spaced-ace/constants"
	"spaced-ace/models"
)

type SignupForm struct {
	Email         string `form:"email"`
	Name          string `form:"name"`
	Password      string `form:"password"`
	PasswordAgain string `form:"passwordAgain"`
}

type SignupRequestBody struct {
	Email         string `json:"email"`
	Name          string `json:"name"`
	Password      string `json:"password"`
	PasswordAgain string `json:"passwordAgain"`
}

func PostRegister(c echo.Context) error {
	var signupForm = SignupForm{}
	if err := c.Bind(&signupForm); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	bodyMap := SignupRequestBody{
		Name:          signupForm.Name,
		Email:         signupForm.Email,
		Password:      signupForm.Password,
		PasswordAgain: signupForm.PasswordAgain,
	}
	bodyBytes, err := json.Marshal(bodyMap)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	bodyBuffer := bytes.NewBuffer(bodyBytes)

	resp, err := http.Post(constants.BACKEND_URL+"/create-user", "application/json", bodyBuffer)
	fmt.Println(err)
	fmt.Println(resp.Status)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadGateway, err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return echo.NewHTTPError(resp.StatusCode, resp.Status)
	}

	// parse the response
	var user models.User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// get the session cookie
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
