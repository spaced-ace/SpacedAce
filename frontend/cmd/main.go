package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"spaced-ace/api"
	"spaced-ace/constants"
	"spaced-ace/context"
	"spaced-ace/render"
	"spaced-ace/views/pages"
	"strings"
)

func main() {
	e := echo.New()

	// Static files
	e.Static("/static", "static")

	e.Use(middleware.Logger())
	e.Use(context.SessionMiddleware)
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5, // Compression level: 1 (fastest, least compression) to 9 (slowest, max compression)
	}))

	//e.Renderer = api.NewTemplate()

	e.GET("/", func(c echo.Context) error {
		return render.TemplRender(c, http.StatusOK, pages.IndexPage())
	})

	api.RegisterRoutes(e)

	//e.HTTPErrorHandler = func(err error, c echo.Context) {
	//	code := http.StatusInternalServerError
	//	message := "internal server error"
	//	var he *echo.HTTPError
	//	if errors.As(err, &he) {
	//		code = he.Code
	//		message = extractErrorMessage(he)
	//	}
	//
	//	data := api.NewPageTemplate(
	//		nil,
	//		map[string]string{
	//			"Message": message,
	//		},
	//	)
	//
	//	c.Response().Header().Set("HX-Push-Url", "false")
	//	if err = c.Render(code, "error-message", data); err != nil {
	//		c.Error(err)
	//	}
	//}

	e.Logger.Fatal(e.Start(":" + constants.PORT))
}

func extractErrorMessage(he *echo.HTTPError) string {
	parts := strings.Split(he.Error(), ", ")
	for _, part := range parts {
		if strings.HasPrefix(part, "message=") {
			return strings.TrimPrefix(part, "message=")
		}
	}
	return "unknown error"
}
