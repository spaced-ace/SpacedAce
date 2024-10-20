package render

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"spaced-ace/views/pipeline"
)

func TemplRender(ctx echo.Context, statusCode int, t templ.Component) error {
	component := t

	// Attach closed popup if the proper header is present
	component = attachClosedPopup(ctx, t)

	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := component.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}

func attachClosedPopup(c echo.Context, content templ.Component) templ.Component {
	if c.Request().Header.Get("SA-popup-action") == "close" {
		return pipeline.ClosedPopup(content)
	}
	return content
}
