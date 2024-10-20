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

	e.GET("/", func(c echo.Context) error {
		return render.TemplRender(c, http.StatusOK, pages.IndexPage())
	})

	api.RegisterRoutes(e)

	e.Logger.Fatal(e.Start(":" + constants.PORT))
}
