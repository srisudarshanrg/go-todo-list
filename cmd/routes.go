package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/srisudarshanrg/go-todo-list/pkg/config"
	"github.com/srisudarshanrg/go-todo-list/pkg/handlers"
)

func Routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Get("/", handlers.Repo.Home)

	return mux
}
