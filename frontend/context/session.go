package context

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"spaced-ace/models"
	"spaced-ace/service"
)

type AppContext struct {
	echo.Context
	Session    *models.Session
	ApiService *service.ApiService
}

func SessionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := &AppContext{
			Context: c,
		}

		sessionCookie, err := c.Cookie("session")
		if err != nil {
			return next(cc)
		}

		cc.ApiService = service.NewApiService(sessionCookie)
		session, err := cc.ApiService.GetSession()
		if err != nil {
			return next(cc)
		}

		cc.Session = session
		return next(cc)
	}
}

func RequireSessionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc, ok := c.(*AppContext)
		if !ok || cc.Session == nil {
			c.Response().Header().Set("HX-Redirect", "/login")
			return c.Redirect(http.StatusFound, "/login")
		}
		return next(c)
	}
}
