package main

import (
	"log"
	"net/http"

	"github.com/srisudarshanrg/go-todo-list/pkg/config"
	"github.com/srisudarshanrg/go-todo-list/pkg/db"
	"github.com/srisudarshanrg/go-todo-list/pkg/handlers"
	"github.com/srisudarshanrg/go-todo-list/pkg/render"
)

const portNumber = ":2121"

var app config.AppConfig

func main() {
	app.UseTemplateCache = false
	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Println(err)
	}

	app.TemplateCache = templateCache

	render.SetConfig(&app)
	repo := handlers.SetRepository(&app)
	handlers.SetHandlers(repo)

	db, err := db.DatabaseConn()
	if err != nil {
		log.Println(err)
	}

	defer db.Close()

	srv := http.Server{
		Handler: Routes(&app),
		Addr:    portNumber,
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
