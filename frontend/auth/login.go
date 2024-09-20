package auth

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"spaced-ace/constants"
	"spaced-ace/render"
	"spaced-ace/views/forms"
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
	errors := map[string]string{}

	var loginForm = LoginForm{}
	if err := c.Bind(&loginForm); err != nil {
		errors["other"] = err.Error()
		return render.TemplRender(c, 200, forms.LoginForm(errors))
	}

	bodyMap := LoginRequestBody{
		Email:    loginForm.Email,
		Password: loginForm.Password,
	}
	bodyBytes, err := json.Marshal(bodyMap)
	if err != nil {
		errors["other"] = "Internal server error"
		return render.TemplRender(c, 200, forms.LoginForm(errors))
	}
	bodyBuffer := bytes.NewBuffer(bodyBytes)

	resp, err := http.Post(constants.BACKEND_URL+"/authenticate-user", "application/json", bodyBuffer)
	if err != nil {
		errors["other"] = "Error: Bad gateway"
		return render.TemplRender(c, 200, forms.LoginForm(errors))
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		errors["other"] = "Invalid e-mail or password"
		return render.TemplRender(c, 200, forms.LoginForm(errors))
	}

	var sessionCookie *http.Cookie
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "session" {
			sessionCookie = cookie
			break
		}
	}
	if sessionCookie == nil {
		errors["other"] = "session cookie not found"
		return render.TemplRender(c, 200, forms.LoginForm(errors))
	}

	c.SetCookie(sessionCookie)
	c.Response().Header().Set("HX-Redirect", "/my-quizzes")
	return c.String(http.StatusOK, "login successful")
}
