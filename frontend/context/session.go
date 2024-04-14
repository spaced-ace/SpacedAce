package context

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"spaced-ace/constants"
)

type User struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Session struct {
	Id   string `json:"session"`
	User User   `json:"user"`
}

type Context struct {
	echo.Context
	Session *Session
}

func SessionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sessionCookie, _ := c.Cookie("session")
		session, _ := getSession(sessionCookie)

		cc := &Context{
			Context: c,
			Session: session,
		}

		return next(cc)
	}
}

func getSession(sessionCookie *http.Cookie) (*Session, error) {
	if sessionCookie != nil {
		req, err := http.NewRequest("GET", constants.BACKEND_URL+"/authenticated", nil)
		if err != nil {
			return nil, err
		}
		req.AddCookie(sessionCookie)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Status code:" + resp.Status)
		}

		var session = Session{}
		err = json.NewDecoder(resp.Body).Decode(&session)
		if err != nil {
			return nil, err
		}
		fmt.Println("Session: ", session)
		return &session, nil
	}
	return nil, fmt.Errorf("no session cookie")
}

func RequireSessionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := c.(*Context)
		if cc.Session == nil {
			fmt.Println("No session")
			c.Response().Header().Set("HX-Redirect", "/login")
			return c.Redirect(http.StatusFound, "/login")
		}
		return next(c)
	}
}

func DeleteSession(sessionId string) error {
	req, err := http.NewRequest("POST", constants.BACKEND_URL+"/logout", nil)
	if err != nil {
		return err
	}

	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: sessionId,
	})

	client := &http.Client{}
	_, err = client.Do(req)
	return err
}
