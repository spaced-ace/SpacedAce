package main

import (
	"github.com/golang-jwt/jwt"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

type LoginForm struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

const secretKey = "secret"

func createJWT(form LoginForm, expires time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": form.Email,
		"exp":   expires.Unix(),
	})
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func login(c echo.Context) error {
	var loginForm LoginForm
	err := c.Bind(&loginForm)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	if loginForm.Email != "user@email.com" && loginForm.Password != "password" {
		return c.String(http.StatusUnauthorized, "unauthorized")
	}

	// Generate the token
	expires := time.Now().Add(1 * time.Hour)
	jwtString, err := createJWT(loginForm, expires)
	if err != nil {
		return c.String(http.StatusInternalServerError, "internal server error")
	}

	// Set the token in a cookie
	cookie := new(http.Cookie)
	cookie.Name = "auth"
	cookie.Value = jwtString
	cookie.Expires = expires
	c.SetCookie(cookie)

	return c.String(http.StatusOK, "login successful")
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	authorized := e.Group("")
	authorized.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(secretKey),
		TokenLookup: "cookie:auth",
	}))
	authorized.GET("/authorized", func(c echo.Context) error {
		return c.String(http.StatusOK, "Authorized")
	})

	public := e.Group("")
	public.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	public.POST("/login", login)

	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
