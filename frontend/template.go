package frontend

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/rs/zerolog/log"
)

var (
	//go:embed private
	fsTemplates embed.FS

	//go:embed assets
	fsStatic embed.FS
)

const (
	templatesBasePath = "frontend/private/"
	templatesExt      = ".tmpl"
)

func (w *Web) parseTemplate(name, path string) {
	// default path to name
	if path == "" {
		path = name
	}

	if _, ok := w.Templates[name]; ok {
		log.Panic().Msg("Trying to ready a template which has already been prepared")
		return
	}

	w.Templates[name] = template.Must(template.New(name).ParseFiles(templatesBasePath+path+templatesExt, templatesBasePath+"base"+templatesExt, templatesBasePath+"navbar"+templatesExt))
}

func (w *Web) templateGet(name string) *template.Template {
	if _, ok := w.Templates[name]; !ok {
		log.Error().Str("name", name).Msg("Trying to get a template that does not exists, returning a 404 page")
		return w.Templates["404.tmpl"]
	}

	return w.Templates[name]
}

func (w *Web) templateExec(rw http.ResponseWriter, r *http.Request, name string, data interface{}) {
	tmplData := struct {
		Errors []flashMessage
		Data   interface{}
	}{
		Errors: w.getFlash(rw, r),
		Data:   data,
	}
	if err := w.templateGet(name).ExecuteTemplate(rw, "base", tmplData); err != nil {
		log.Error().Err(err).Str("name", name).Interface("data", data).Msg("failed to view template")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}
