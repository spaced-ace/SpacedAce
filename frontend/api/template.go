package api

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"spaced-ace/models"
)

type Template struct {
	tmpl *template.Template
}

func NewTemplate() *Template {
	return &Template{
		tmpl: template.Must(template.ParseGlob("views/**/*.html")),
	}
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	fmt.Printf("Data: %+v\n", data)
	renderData := data.(*TemplateData)

	htmxRequest := c.Request().Header.Get("HX-Request") == "true"
	if renderData.RenderType == Component || htmxRequest {
		fmt.Printf("Rendering component: `%s`, data: %+v\n", name, renderData.Data)
		return t.tmpl.ExecuteTemplate(w, name, renderData.Data)
	}

	pageRenderData := PageRenderData{
		TemplateName: name,
		PageData: PageData{
			Session: renderData.Session,
			Data:    renderData.Data,
		},
	}
	fmt.Printf("Rendering page: `%s`, data: %+v\n", name, renderData.Data)
	return t.tmpl.ExecuteTemplate(w, "layout", pageRenderData)
}

type PageData struct {
	Session *models.Session
	Data    interface{}
}
type PageRenderData struct {
	TemplateName string
	PageData     interface{}
}

type RenderType int

const (
	Component RenderType = iota
	Page
)

type TemplateData struct {
	Session    *models.Session
	RenderType RenderType
	Data       interface{}
}

func NewPageTemplate(session *models.Session, data interface{}) *TemplateData {
	return &TemplateData{
		Session:    session,
		RenderType: Page,
		Data:       data,
	}
}
func NewComponentTemplate(session *models.Session, data interface{}) *TemplateData {
	return &TemplateData{
		Session:    session,
		RenderType: Component,
		Data:       data,
	}
}
