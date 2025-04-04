package context

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"spaced-ace/models/business"
	"spaced-ace/service"
)

type AppContext struct {
	echo.Context
	Session    *business.Session
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
		cc.Set("cc", cc)
		return next(cc)
	}
}

func RequireSessionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc, ok := c.(*AppContext)
		if !ok || cc.Session == nil {
			if c.Request().Header.Get("HX-Request") == "true" {
				c.Response().Header().Set("HX-Redirect", "/login")
				c.Response().Header().Set("HX-Refresh", "true")
				return c.NoContent(http.StatusUnauthorized)
			}
			return c.Redirect(http.StatusFound, "/login")
		}

		// Check if email is verified
		if !cc.Session.User.EmailVerified {
			emailVerificationUrl := "/email-verification-needed?email=" + cc.Session.User.Email
			if c.Request().Header.Get("HX-Request") == "true" {
				c.Response().Header().Set("HX-Redirect", emailVerificationUrl)
				return c.NoContent(http.StatusForbidden)
			}
			return c.Redirect(http.StatusFound, emailVerificationUrl)
		}

		return next(c)
	}
}
