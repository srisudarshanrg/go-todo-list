package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/srisudarshanrg/go-todo-list/server/database"
	"github.com/srisudarshanrg/go-todo-list/server/functions"
	"github.com/srisudarshanrg/go-todo-list/server/models"
	"github.com/srisudarshanrg/go-todo-list/server/setup"
	"github.com/srisudarshanrg/go-todo-list/server/validations"
)

const portNumber = ":9000"

var session *scs.SessionManager

func main() {
	gob.Register(models.User{})

	// session
	session = scs.New()
	session.Cookie.Persist = true
	session.Lifetime = 24 * time.Hour
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	// database setup
	db, err := database.CreateDatabaseConn()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// database and access
	setup.DBAccessHandlers(db)
	setup.SessionAccessHandlers(session)
	functions.DBAccessFunctions(db)
	validations.DBAccessFormValidations(db)
	validations.SessionAccessValidations(session)

	// routes
	server := http.Server{
		Addr:    portNumber,
		Handler: routes(),
	}

	log.Println("server running on port number", portNumber)
	server.ListenAndServe()
}

func routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(SessionLoadAndSave)

	mux.Get("/home", setup.Home)
	mux.Get("/login", setup.Login)
	mux.Get("/register", setup.Register)
	mux.Get("/tasks", setup.Tasks)
	mux.Get("/habit-tracker", setup.HabitTracker)
	mux.Get("/notes", setup.Notes)
	mux.Get("/profile", setup.Profile)
	mux.Get("/logout", setup.Logout)

	mux.Post("/home", setup.HomePost)
	mux.Post("/login", setup.LoginPost)
	mux.Post("/register", setup.RegisterPost)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

func SessionLoadAndSave(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
