package render

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/srisudarshanrg/go-todo-list/pkg/config"
	"github.com/srisudarshanrg/go-todo-list/pkg/models"
)

var app *config.AppConfig

func SetConfig(a *config.AppConfig) {
	app = a
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, templateData *models.TemplateData) {
	var templateCache map[string]*template.Template
	var err error

	if app.UseTemplateCache {
		templateCache = app.TemplateCache
	} else {
		templateCache, err = CreateTemplateCache()
		if err != nil {
			log.Fatal(err)
		}
	}

	template, exist := templateCache[tmpl]

	if !exist {
		log.Println("requested template could not be retrieved from template cache.")
	}

	err = template.Execute(w, templateData)
	if err != nil {
		log.Println(err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	templateCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./template/*.page.tmpl")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for _, page := range pages {
		templateName := filepath.Base(page)
		templateSet, err := template.New(templateName).ParseFiles(page)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		layoutExist, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			log.Println(err)
			return nil, err
		}

		if len(layoutExist) > 0 {
			templateSet, err = templateSet.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				log.Println(err)
				return nil, err
			}
		}

		templateCache[templateName] = templateSet
	}

	return templateCache, nil
}
