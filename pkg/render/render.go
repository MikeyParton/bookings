package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/MikeyParton/bookings/pkg/config"
	"github.com/MikeyParton/bookings/pkg/models"
)

var app *config.AppConfig
var functions template.FuncMap

func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultTemplateData(td models.TemplateData) models.TemplateData {
	return td
}

func RenderTemplate(w http.ResponseWriter, name string, data *models.TemplateData) {
	var tc map[string]*template.Template
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[name]

	if !ok {
		log.Fatal("Couldn't find template")
	}

	buf := new(bytes.Buffer)

	_ = t.Execute(buf, data)
	_, err := buf.WriteTo(w)

	if err != nil {
		fmt.Println("Error writing template", err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	var cache = map[string]*template.Template{}
	pages, err := filepath.Glob("./templates/*.page.html")
	if err != nil {
		fmt.Println("Error finding pages", err)
		return cache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)

		if err != nil {
			return cache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.html")

		if err != nil {
			return cache, err
		}

		if len(matches) > 0 {
			ts, err := ts.ParseGlob("./templates/*.layout.html")
			if err != nil {
				return cache, err
			}

			cache[name] = ts
		}
	}

	return cache, nil
}
