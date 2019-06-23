package web

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
)

// compileTemplates minifies then parses each file in filenames.
func compileTemplates(filenames []string) (*template.Template, error) {
	minifier := minify.New()
	minifier.AddFunc("text/html", html.Minify)

	// Minify and parse each file
	var tmpl *template.Template
	for _, filename := range filenames {
		name := filepath.Base(filename)
		if tmpl == nil {
			tmpl = template.New(name)
		} else {
			tmpl = tmpl.New(name)
		}

		b, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}

		mb, err := minifier.Bytes("text/html", b)
		if err != nil {
			return nil, err
		}

		tmpl.Parse(string(mb))
	}

	return tmpl, nil
}

// Loads templates (files ending in .tmpl) from web/templates/ and all its
// subdirectories.
func (h *Handler) loadTemplates() (err error) {
	webPath, ok := os.LookupEnv("WEB_PATH")
	if !ok {
		log.Fatalf("environment variable not set: %v", "WEB_PATH")
	}

	// Get each file (including in subdirectories) in $WEB_PATH/templates
	var fileNames []string
	err = filepath.Walk(webPath+"/templates", func(path string, _ os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only add .tmpl files
		if strings.HasSuffix(path, ".tmpl") {
			fileNames = append(fileNames, path)
		}

		return nil
	})
	if err != nil {
		return err
	}

	// Minify and parse template files
	h.templates, err = compileTemplates(fileNames)
	return err
}

// reloadTemplates is a hacky middleware that calls loadTemplates on every request.
// It's useful to live reload templates in development but shouldn't be used in
// production.
func (h *Handler) reloadTemplates(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Reload templates
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
	// Grab a buffer from the buffer pool
	buf := h.templateBufPool.Get().(*bytes.Buffer)
	defer h.templateBufPool.Put(buf)
	defer buf.Reset()

	// Try rendering the template
	err := h.templates.ExecuteTemplate(buf, name, data)
	if err != nil {
		// If it fails, show the status code 500 template
		h.templates.ExecuteTemplate(w, "500", nil)
		return
	}

	// Otherwise write the buffer to the response
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	buf.WriteTo(w)
	return
}
