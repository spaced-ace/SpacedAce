package auth

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	bcrypt "golang.org/x/crypto/bcrypt"
)

type User struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
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

type ResendEmailverificationRequest struct {
	Email string `json:"email"`
}

type EmailVerificationResponse struct {
	Message string `json:"message"`
}

const EMAIL_VERIFICATION_SUCCESS = "Successfully verified email"
const EMAIL_VERIFICATION_FAIL = "Failed to verified email"
const EMAIL_VERIFICATION_ALREADY_VERIFIED = "Email already verified"
const EMAIL_VERIFICATION_RESEND = "If your email is registered, a verification link has been sent"
const EMAIL_VERIFICATION_SENT = "Email verification sent"

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

	if !user.EmailVerified {
		return echo.NewHTTPError(http.StatusForbidden, "email not verified")
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
		Id:            user.Id,
		Name:          user.Name,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
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
			Id:            dbUser.Id,
			Name:          dbUser.Name,
			Email:         dbUser.Email,
			EmailVerified: dbUser.EmailVerified,
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

	verificationToken := GenerateVerificationToken()

	newUser := DBUser{
		Id:                uuid.NewString(),
		Name:              request.Name,
		Email:             request.Email,
		Password:          string(bcryptPassword),
		EmailVerified:     false,
		VerificationToken: &verificationToken,
	}
	err = CreateUser(&newUser)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create user")
	}

	emailSvc := GetEmailVerificationService()
	err = emailSvc.SendVerificationEmail(newUser.Email, newUser.Name, verificationToken)
	if err != nil {
		// Log the error but don't fail registration
		fmt.Printf("Error sending verification email: %v\n", err)
	}

	sessionId, err := CreateSession(newUser.Id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create session")
	}
	c.SetCookie(&http.Cookie{
		Name:     "session",
		Value:    sessionId,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(59 * time.Minute),
	})

	authResponse := AuthResponse{
		Session: sessionId,
		User: User{
			Id:            newUser.Id,
			Name:          newUser.Name,
			Email:         newUser.Email,
			EmailVerified: newUser.EmailVerified,
		},
	}

	return c.JSON(http.StatusOK, authResponse)
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

func DeleteUserEndpoint(c echo.Context) error {
	session, err := c.Cookie("session")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "login required")
	}
	userid, err := GetUserIdBySession(session.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	err = DeleteUser(userid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	return c.NoContent(http.StatusOK)
}

func VerifyEmailEndpoint(c echo.Context) error {
	token := c.QueryParam("token")
	if token == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "verification token is required")
	}

	user, err := GetUserByVerificationToken(token)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "invalid or expired verification token")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	if user.EmailVerified {
		return c.JSON(http.StatusOK, EmailVerificationResponse{EMAIL_VERIFICATION_ALREADY_VERIFIED})
	}

	err = VerifyEmail(user.Id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to verify email")
	}

	return c.JSON(http.StatusOK, EmailVerificationResponse{EMAIL_VERIFICATION_SUCCESS})
}

func ResendVerificationEmailEndpoint(c echo.Context) error {
	var email ResendEmailverificationRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&email); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}

	if email.Email == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "email is required")
	}

	user, err := GetUserByEmail(email.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			// Don't reveal if email exists or not
			return c.JSON(http.StatusOK, EmailVerificationResponse{EMAIL_VERIFICATION_RESEND})
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	if user.EmailVerified {
		return c.JSON(http.StatusOK, EmailVerificationResponse{EMAIL_VERIFICATION_ALREADY_VERIFIED})
	}

	// Generate a new verification token if needed
	if *user.VerificationToken == "" {
		*user.VerificationToken = GenerateVerificationToken()
		err = UpdateUser(user)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to update user")
		}
	}

	// Send verification email
	svc := GetEmailVerificationService()
	err = svc.SendVerificationEmail(user.Email, user.Name, *user.VerificationToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to send verification email")
	}

	return c.JSON(http.StatusOK, EmailVerificationResponse{EMAIL_VERIFICATION_SENT})
}
