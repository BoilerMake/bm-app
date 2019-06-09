package web

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Loads templates (files ending in .tmpl) from web/templates/ and all its
// subdirectories.
func (h *Handler) loadTemplates() (err error) {
	h.templates = template.New("")

	err = filepath.Walk("web/templates/", func(path string, _ os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, ".tmpl") {
			template.Must(h.templates.ParseFiles(path))
		}

		return nil
	})

	return err
}

// reloadTemplates is a hacky middleware that calls loadTemplates on every request.
// It's useful to live reload templates in development but shouldn't be used in
// production.
func (h *Handler) reloadTemplates(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := h.loadTemplates()
		if err != nil {
			log.Fatalf("failed to reload templates: %s", err)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// renderTemplate calls ExecuteTemplate on the passed template name.  If the
// template fails to render then render the status code 500 page
func (h *Handler) renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	buf := h.templateBufPool.Get().(*bytes.Buffer)
	defer h.templateBufPool.Put(buf)

	err := h.templates.ExecuteTemplate(buf, name, data)
	if err != nil {
		fmt.Println(err)
		h.templates.ExecuteTemplate(w, "500", nil)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	buf.WriteTo(w)
	return
}
