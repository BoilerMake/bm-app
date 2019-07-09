package template

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
)

// A Template wraps go's html/template to provide easier use.
type Template struct {
	Template *template.Template
	bufPool  *sync.Pool
	tmplPath string
}

// NewTemplate creates a new Template object.  It will attempt to parse the
// templates available (*.tmpl files) in the given path and its subdirectories.
func NewTemplate(tmplPath string) (t *Template, err error) {
	// "missingkey=zero" means when you try to substitute a variable that is
	// nil inside a template, it will default to its zero value instead of the
	// string "<no value>"
	tmpl := template.New("").Option("missingkey=zero")

	bufPool := &sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}

	t = &Template{
		Template: tmpl,
		bufPool:  bufPool,
		tmplPath: tmplPath,
	}

	if tmplPath != "" {
		err = t.LoadTemplates()
		if err != nil {
			return t, err
		}
	} else {
		return t, fmt.Errorf("No template path given")
	}

	return t, err
}

// CompileTemplates minifies then parses each file in filenames.
func CompileTemplates(filenames []string) (*template.Template, error) {
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

// LoadTemplates loads templates (files ending in .tmpl) from a given directory
// and all its subdirectories.
func (t *Template) LoadTemplates() (err error) {
	// Get each file (including in subdirectories) in $TEMPLATES_PATH/templates
	var fileNames []string
	err = filepath.Walk(t.tmplPath, func(path string, _ os.FileInfo, err error) error {
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
	t.Template, err = CompileTemplates(fileNames)
	return err
}

// ReloadTemplates is a middleware that calls LoadTemplates. It's useful during
// development to live reload templates in development but shouldn't be used in
// production.
func (t *Template) ReloadTemplates(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Reload templates
		err := t.LoadTemplates()
		if err != nil {
			log.Fatalf("failed to reload templates: %s", err)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RenderTemplate calls ExecuteTemplate on the passed template name.  If the
// template fails to render then render the status code 500 page
func (t *Template) RenderTemplate(w http.ResponseWriter, name string, data interface{}) {
	// Grab a buffer from the buffer pool
	buf := t.bufPool.Get().(*bytes.Buffer)
	defer t.bufPool.Put(buf)
	defer buf.Reset()

	// Try rendering the template
	err := t.Template.ExecuteTemplate(buf, name, data)
	if err != nil {
		// If it fails, show the status code 500 template
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		t.Template.ExecuteTemplate(w, "500", nil)
		return
	}

	// Otherwise write the buffer to the response
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	buf.WriteTo(w)
	return
}

// RenderToString calls ExecuteTemplate on the passed template name.
func (t *Template) RenderToString(name string, data interface{}) (result string, err error) {
	// Grab a buffer from the buffer pool
	buf := t.bufPool.Get().(*bytes.Buffer)
	defer t.bufPool.Put(buf)
	defer buf.Reset()

	// Try rendering the template
	err = t.Template.ExecuteTemplate(buf, name, data)
	if err != nil {
		return result, err
	}

	// Otherwise write the buffer to the response
	return buf.String(), err
}
