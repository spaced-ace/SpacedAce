package context

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type User struct {
	Id    string
	Name  string
	Email string
}

type Session struct {
	Id      string
	User    User
	Expires time.Time
}

type Context struct {
	echo.Context
	Session *Session
}

var sessions = make(map[string]Session)

func getSessionFromContext(c echo.Context) *Session {
	cookie, err := c.Cookie("session")
	if err != nil {
		return nil
	}
	return GetSession(cookie.Value)
}

func SessionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		session := getSessionFromContext(c)

		// Invalid session id sent by the client
		if session == nil {
			cookie, _ := c.Cookie("session")
			if cookie != nil {
				cookie.Expires = time.Now().Add(-1 * time.Hour)
				c.SetCookie(cookie)
			}
		}

		// Session exists but has expired
		if session != nil {
			if session.Expires.Before(time.Now()) {
				DeleteSession(session.Id)
				session = nil
			}
		}

		cc := &Context{
			Context: c,
			Session: session,
		}

		return next(cc)
	}
}

func RequireSessionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := c.(*Context)
		if cc.Session == nil {
			c.Response().Header().Set("HX-Redirect", "/login")
			return c.Redirect(http.StatusFound, "/login")
		}
		return next(c)
	}
}

func CreateSession(user User) *Session {
	id := uuid.New().String()
	session := Session{
		Id:      id,
		User:    user,
		Expires: time.Now().Add(1 * time.Hour),
	}
	sessions[id] = session
	return &session
}

func GetSession(id string) *Session {
	session, ok := sessions[id]
	if !ok {
		return nil
	}
	return &session
}

func DeleteSession(id string) {
	delete(sessions, id)
}
