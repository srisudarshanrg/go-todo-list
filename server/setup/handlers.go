package setup

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/srisudarshanrg/go-todo-list/server/functions"
	"github.com/srisudarshanrg/go-todo-list/server/models"
)

var db *sql.DB
var session *scs.SessionManager

// DBAccess provides the handlers with access to the database
func DBAccessHandlers(dbAccess *sql.DB) {
	db = dbAccess
}

// SessionAccessHandlers provides the handlers package with access to the sessions
func SessionAccessHandlers(sessionAccess *scs.SessionManager) {
	session = sessionAccess
}

// Login is the handler for the login page
func Login(w http.ResponseWriter, r *http.Request) {
	err := session.Destroy(r.Context())
	if err != nil {
		log.Println(err)
	}

	msg := r.URL.Query().Get("msg")
	if msg != "" {
		err := RenderTemplate(w, r, "login.page.tmpl", models.TemplateData{
			Info: msg,
		})
		if err != nil {
			log.Println(err)
		}
		return
	}

	err = RenderTemplate(w, r, "login.page.tmpl", models.TemplateData{})
	if err != nil {
		log.Println(err)
	}
}

// Register is the handler for the register page
func Register(w http.ResponseWriter, r *http.Request) {
	err := RenderTemplate(w, r, "register.page.tmpl", models.TemplateData{})
	if err != nil {
		log.Println(err)
	}
}

// Home is the handler for the home page
func Home(w http.ResponseWriter, r *http.Request) {
	userInterface := session.Get(r.Context(), "user")
	user, check := userInterface.(models.User)
	if !check {
		msg := "Login is required to access this page"
		http.Redirect(w, r, "/login?msg="+msg, http.StatusSeeOther)
	}

	tasks, err := functions.GetTasks(user.ID)
	if err != nil {
		log.Println(err)
	}

	habits, err := functions.GetHabits(user.ID)
	if err != nil {
		log.Println(err)
	}

	notes, err := functions.GetNotes(user.ID)
	if err != nil {
		log.Println(err)
	}

	data := map[string]interface{}{}
	data["tasks"] = tasks
	data["habits"] = habits
	data["notes"] = notes

	// do the msg url checking after getting all the database data and doing all the logic
	msg := r.URL.Query().Get("msg")
	if msg != "" {
		err := RenderTemplate(w, r, "home.page.tmpl", models.TemplateData{
			Info: msg,
			Data: data,
		})
		if err != nil {
			log.Println(err)
		}
		return
	}

	err = RenderTemplate(w, r, "home.page.tmpl", models.TemplateData{
		Data: data,
	})
	if err != nil {
		log.Println(err)
	}
}

func Tasks(w http.ResponseWriter, r *http.Request) {
	userInterface := session.Get(r.Context(), "user")
	user, check := userInterface.(models.User)
	if !check {
		msg := "Login is required to access this page"
		http.Redirect(w, r, "/login?msg="+msg, http.StatusSeeOther)
	}
	log.Println(user)

	err := RenderTemplate(w, r, "tasks.page.tmpl", models.TemplateData{})
	if err != nil {
		log.Println(err)
	}
}

func HabitTracker(w http.ResponseWriter, r *http.Request) {
	userInterface := session.Get(r.Context(), "user")
	user, check := userInterface.(models.User)
	if !check {
		msg := "Login is required to access this page"
		http.Redirect(w, r, "/login?msg="+msg, http.StatusSeeOther)
	}
	log.Println(user)

	err := RenderTemplate(w, r, "habit-tracker.page.tmpl", models.TemplateData{})
	if err != nil {
		log.Println(err)
	}
}

func Notes(w http.ResponseWriter, r *http.Request) {
	userInterface := session.Get(r.Context(), "user")
	user, check := userInterface.(models.User)
	if !check {
		msg := "Login is required to access this page"
		http.Redirect(w, r, "/login?msg="+msg, http.StatusSeeOther)
	}
	log.Println(user)

	err := RenderTemplate(w, r, "notes.page.tmpl", models.TemplateData{})
	if err != nil {
		log.Println(err)
	}
}

func Profile(w http.ResponseWriter, r *http.Request) {
	userInterface := session.Get(r.Context(), "user")
	user, check := userInterface.(models.User)
	if !check {
		msg := "Login is required to access this page"
		http.Redirect(w, r, "/login?msg="+msg, http.StatusSeeOther)
	}
	log.Println(user)

	err := RenderTemplate(w, r, "profile.page.tmpl", models.TemplateData{})
	if err != nil {
		log.Println(err)
	}
}

// Logout is the handler for the logout functionality
func Logout(w http.ResponseWriter, r *http.Request) {
	err := session.Destroy(r.Context())
	if err != nil {
		log.Println(err)
	}
	msg := "You have been logged out"
	http.Redirect(w, r, "/login?msg="+msg, http.StatusSeeOther)
}
