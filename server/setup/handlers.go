package setup

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/srisudarshanrg/go-todo-list/server/functions"
	"github.com/srisudarshanrg/go-todo-list/server/models"
)

var db *sql.DB
var session *scs.SessionManager
var data = map[string]interface{}{}

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

	tasks, _, err := functions.GetTasks(user.ID)
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

// Tasks is the handler for the tasks page
func Tasks(w http.ResponseWriter, r *http.Request) {
	userInterface := session.Get(r.Context(), "user")
	user, check := userInterface.(models.User)
	if !check {
		msg := "Login is required to access this page"
		http.Redirect(w, r, "/login?msg="+msg, http.StatusSeeOther)
	}

	session.Put(r.Context(), "linkTasks", "/tasks")
	session.Put(r.Context(), "pathTasks", "tasks.page.tmpl")

	tasks, completed, err := functions.GetTasks(user.ID)
	if err != nil {
		log.Println(err)
	}

	completedString, err := json.Marshal(completed)
	if err != nil {
		log.Println(err)
	}

	data["tasks"] = tasks
	data["completedTasks"] = string(completedString)

	err = RenderTemplate(w, r, "tasks.page.tmpl", models.TemplateData{
		Data: data,
	})
	if err != nil {
		log.Println(err)
	}
}

// TasksListView is the handler for the tasks list page
func TasksListView(w http.ResponseWriter, r *http.Request) {
	userInterface := session.Get(r.Context(), "user")
	user, check := userInterface.(models.User)
	if !check {
		msg := "Login is required to access this page"
		http.Redirect(w, r, "/login?msg="+msg, http.StatusSeeOther)
	}

	session.Put(r.Context(), "linkTasks", "/tasks-list")
	session.Put(r.Context(), "pathTasks", "tasks-list.page.tmpl")

	tasks, completed, err := functions.GetTasks(user.ID)
	if err != nil {
		log.Println(err)
	}

	completedString, err := json.Marshal(completed)
	if err != nil {
		log.Println(err)
	}

	data["tasks"] = tasks
	data["completedTasks"] = string(completedString)

	err = RenderTemplate(w, r, "tasks-list.page.tmpl", models.TemplateData{
		Data: data,
	})
	if err != nil {
		log.Println(err)
	}
}

// HabitTracker is the handler for the habit tracker page
func HabitTracker(w http.ResponseWriter, r *http.Request) {
	userInterface := session.Get(r.Context(), "user")
	user, check := userInterface.(models.User)
	if !check {
		msg := "Login is required to access this page"
		http.Redirect(w, r, "/login?msg="+msg, http.StatusSeeOther)
	}

	session.Put(r.Context(), "linkHabits", "/habit-tracker")
	session.Put(r.Context(), "pathHabits", "habit-tracker.page.tmpl")

	habits, err := functions.GetHabits(user.ID)
	if err != nil {
		log.Println(err)
	}

	data["habits"] = habits

	err = RenderTemplate(w, r, "habit-tracker.page.tmpl", models.TemplateData{
		Data: data,
	})
	if err != nil {
		log.Println(err)
	}
}

// HabitTrackerListView is the handler for the habits list page
func HabitTrackerListView(w http.ResponseWriter, r *http.Request) {
	userInterface := session.Get(r.Context(), "user")
	user, check := userInterface.(models.User)
	if !check {
		msg := "Login is required to access this page"
		http.Redirect(w, r, "/login?msg="+msg, http.StatusSeeOther)
	}

	session.Put(r.Context(), "linkHabits", "/habit-tracker-list")
	session.Put(r.Context(), "pathHabits", "habits-list.page.tmpl")

	habits, err := functions.GetHabits(user.ID)
	if err != nil {
		log.Println(err)
	}

	data["habits"] = habits

	err = RenderTemplate(w, r, "habits-list.page.tmpl", models.TemplateData{
		Data: data,
	})
	if err != nil {
		log.Println(err)
	}
}

// Notes is the handler for the notes page
func Notes(w http.ResponseWriter, r *http.Request) {
	userInterface := session.Get(r.Context(), "user")
	user, check := userInterface.(models.User)
	if !check {
		msg := "Login is required to access this page"
		http.Redirect(w, r, "/login?msg="+msg, http.StatusSeeOther)
	}

	notes, err := functions.GetNotes(user.ID)
	if err != nil {
		log.Println(err)
	}

	data["notes"] = notes

	err = RenderTemplate(w, r, "notes.page.tmpl", models.TemplateData{
		Data: data,
	})
	if err != nil {
		log.Println(err)
	}
}

// Profile is the handler for the profile page
func Profile(w http.ResponseWriter, r *http.Request) {
	userInterface := session.Get(r.Context(), "user")
	user, check := userInterface.(models.User)
	if !check {
		msg := "Login is required to access this page"
		http.Redirect(w, r, "/login?msg="+msg, http.StatusSeeOther)
	}

	err := RenderTemplate(w, r, "profile.page.tmpl", models.TemplateData{
		Data: user,
	})
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
