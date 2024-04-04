package auth

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	bcrypt "golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type User struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
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
type AuthResponse struct {
	Session string `json:"session"`
	User    User   `json:"user"`
}

func init() {
}

func AuthenticateUser(c echo.Context) error {
	var request = LoginBody{}
	if err := json.NewDecoder(c.Request().Body).Decode(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}
	var user, err = GetUserByEmail(request.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	session, err := CreateSession(user.Id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error, failed to create session")
	}
	c.SetCookie(&http.Cookie{
		Name:     "session",
		Value:    session,
		HttpOnly: true,
		Expires:  time.Now().Add(59 * time.Minute),
	})

	var userResponse = User{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	}
	return c.JSON(http.StatusOK, userResponse)
}

func Authenticated(c echo.Context) error {
	session, err := c.Cookie("session")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	userId, err := GetUserIdBySession(session.Value)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
		}
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	dbUser, err := GetUserById(userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	authResponse := AuthResponse{
		Session: session.Value,
		User: User{
			Id:    dbUser.Id,
			Name:  dbUser.Name,
			Email: dbUser.Email,
		},
	}
	return c.JSON(http.StatusOK, authResponse)
}

func Register(c echo.Context) error {
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

	var oldUser, err = GetUserByEmail(request.Email)
	if err != nil {
		if err != sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}
	}
	if oldUser.Email != "" {
		return echo.NewHTTPError(http.StatusConflict, "user already exists with this email")
	}

	if request.Password != request.PasswordAgain {
		return echo.NewHTTPError(http.StatusBadRequest, "passwords do not match")
	}
	if len(request.Password) < 8 {
		return echo.NewHTTPError(http.StatusBadRequest, "password must be at least 8 characters long")
	}
	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	newUser := DBUser{
		Id:       uuid.NewString(),
		Name:     request.Name,
		Email:    request.Email,
		Password: string(bcryptPassword),
	}
	err = CreateUser(&newUser)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create user")
	}

	session, err := CreateSession(newUser.Id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error, failed to create session")
	}
	c.SetCookie(&http.Cookie{
		Name:     "session",
		Value:    session,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(59 * time.Minute),
	})

	return c.JSON(http.StatusOK, newUser)
}

func Logout(c echo.Context) error {
	session, err := c.Cookie("session")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	err = DeleteSession(session.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	return c.NoContent(http.StatusOK)
}
