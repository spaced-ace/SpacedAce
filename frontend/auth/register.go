package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"spaced-ace/constants"
	"spaced-ace/models/business"
	"spaced-ace/models/request"
	"spaced-ace/render"
	"spaced-ace/views/forms"
)

type SignupRequestBody struct {
	Email         string `json:"email"`
	Name          string `json:"name"`
	Password      string `json:"password"`
	PasswordAgain string `json:"passwordAgain"`
}

func PostRegister(c echo.Context) error {
	errors := map[string]string{}

	var signupForm = request.SignupForm{}
	if err := c.Bind(&signupForm); err != nil {
		log.Default().Println(err.Error())
		errors["other"] = "Parsing error"
		return render.TemplRender(c, 200, forms.SignUpForm(signupForm, errors))
	}

	// Remove the password and password again fields
	sanitizedSignupForm := request.SignupForm{
		Name:  signupForm.Name,
		Email: signupForm.Email,
	}

	if signupForm.Name == "" {
		errors["name"] = "Name is required"
	}
	if signupForm.Email == "" {
		errors["email"] = "Email is required"
	}
	if signupForm.Password == "" {
		errors["password"] = "Password is required"
	}
	if signupForm.Password == "" {
		errors["password_again"] = "Password again is required"
	}
	if signupForm.Password != signupForm.PasswordAgain {
		errors["password"] = "Different passwords"
		errors["password_again"] = "Different passwords"
	}
	if len(errors) > 0 {
		return render.TemplRender(c, 200, forms.SignUpForm(sanitizedSignupForm, errors))
	}

	bodyMap := SignupRequestBody{
		Name:          signupForm.Name,
		Email:         signupForm.Email,
		Password:      signupForm.Password,
		PasswordAgain: signupForm.PasswordAgain,
	}
	bodyBytes, err := json.Marshal(bodyMap)
	if err != nil {
		log.Default().Println(err.Error())
		errors["other"] = "Internal server error"
		return render.TemplRender(c, 200, forms.SignUpForm(sanitizedSignupForm, errors))
	}
	bodyBuffer := bytes.NewBuffer(bodyBytes)

	resp, err := http.Post(constants.BACKEND_URL+"/create-user", "application/json", bodyBuffer)
	fmt.Println(err)
	fmt.Println(resp.Status)
	if err != nil {
		log.Default().Println(err.Error())
		errors["other"] = "Internal server error"
		return render.TemplRender(c, 200, forms.SignUpForm(sanitizedSignupForm, errors))
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Default().Printf("Register status code: %d\n", resp.StatusCode)
		errors["other"] = "Internal server error"
		return render.TemplRender(c, 200, forms.SignUpForm(sanitizedSignupForm, errors))
	}

	// parse the response
	var user business.User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		log.Default().Println(err.Error())
		errors["other"] = "Internal server error"
		return render.TemplRender(c, 200, forms.SignUpForm(sanitizedSignupForm, errors))
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
		log.Default().Println("session cookie not found")
		errors["other"] = "Error: session cookie not found"
		return render.TemplRender(c, 200, forms.SignUpForm(sanitizedSignupForm, errors))
	}

	c.SetCookie(sessionCookie)
	c.Response().Header().Set("HX-Redirect", "/my-quizzes")
	return c.String(http.StatusOK, "login successful")
}
