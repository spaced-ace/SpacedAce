package auth

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string
}

type LoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupBody struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	PasswordAgain string `json:"passwordAgain"`
}

var db []User

func init() {
	db = []User{
		{
			Id:       uuid.NewString(),
			Name:     "John Doe",
			Email:    "user@email.com",
			Password: "password",
		},
	}
}

func AuthenticateUser(c echo.Context) error {
	var request = LoginBody{}
	if err := json.NewDecoder(c.Request().Body).Decode(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}

	var user User
	for _, u := range db {
		if u.Email == request.Email {
			user = u
			break
		}
	}
	if user.Email == "" || user.Password != request.Password {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	return c.JSON(http.StatusOK, user)
}

func CreateUser(c echo.Context) error {
	var request = SignupBody{}
	if err := json.NewDecoder(c.Request().Body).Decode(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}

	if request.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "name is required")
	}
	if request.Email == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "email is required")
	}
	if request.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "password is required")
	}
	if request.PasswordAgain == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "password again is required")
	}

	var oldUser User
	for _, u := range db {
		if u.Email == request.Email {
			oldUser = u
			break
		}
	}
	if oldUser.Email != "" {
		return echo.NewHTTPError(http.StatusConflict, "user already exists with this email")
	}

	if request.Password != request.PasswordAgain {
		return echo.NewHTTPError(http.StatusBadRequest, "passwords do not match")
	}

	newUser := User{
		Id:       uuid.NewString(),
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	}
	db = append(db, newUser)

	return c.JSON(http.StatusOK, newUser)
}
