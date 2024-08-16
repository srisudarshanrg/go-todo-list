package handlers

import (
	"net/http"

	"github.com/srisudarshanrg/go-todo-list/pkg/config"
	"github.com/srisudarshanrg/go-todo-list/pkg/models"
	"github.com/srisudarshanrg/go-todo-list/pkg/render"
)

type Repository struct {
	app *config.AppConfig
}

var Repo *Repository

func SetRepository(a *config.AppConfig) *Repository {
	return &Repository{
		app: a,
	}
}

func SetHandlers(repo *Repository) {
	Repo = repo
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
}
