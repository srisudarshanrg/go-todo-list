package setup

import (
	"log"
	"net/http"
	"text/template"

	"github.com/srisudarshanrg/go-todo-list/server/models"
)

// RenderTemplate parses and executes a template
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, templateData models.TemplateData) error {
	template, err := template.ParseFiles("./templates/"+tmpl, "./templates/base.layout.tmpl", "./templates/auth.layout.tmpl")
	if err != nil {
		log.Println(err)
		return err
	}

	err = template.Execute(w, templateData)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
