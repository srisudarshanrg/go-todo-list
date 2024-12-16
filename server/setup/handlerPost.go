package setup

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/srisudarshanrg/go-todo-list/server/functions"
	"github.com/srisudarshanrg/go-todo-list/server/models"
	"github.com/srisudarshanrg/go-todo-list/server/validations"
)

// LoginPost handles the post requests to the login page
func LoginPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	credential := r.Form.Get("credential")
	password := r.Form.Get("password")

	check, msg, user, err := functions.AuthenticateUser(credential, password)
	if !check {
		RenderTemplate(w, r, "login.page.tmpl", models.TemplateData{
			Error: msg,
		})
		if err != nil {
			log.Println(err)
		}
		return
	}

	// put user in session
	session.Put(r.Context(), "user", user)

	http.Redirect(w, r, "/?msg="+msg, http.StatusSeeOther)
}

// RegisterPost handles the post requests to the login page
func RegisterPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	username := r.Form.Get("username")
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	passwordConfirm := r.Form.Get("confirmPassword")

	// form validations
	validations.MaxLength(username, 30)
	validations.MinLength(username, 2)
	validations.ValidEmail(email)
	validations.PasswordEqualConfirmPassword(password, passwordConfirm)
	validations.UsernameExists(username)
	validations.EmailExists(email)

	// put error list in session
	validations.PutErrorListInSession(r.Context())

	errorList := session.Get(r.Context(), "errorList").([]string)
	if len(errorList) > 0 {
		RenderTemplate(w, r, "register.page.tmpl", models.TemplateData{
			Data: errorList,
		})
		log.Println("validation problem")
		session.Remove(r.Context(), "errorList")
		return
	}

	passwordHash, err := functions.HashPassword(password)
	if err != nil {
		log.Println(err)
		return
	}

	// create user
	err = functions.CreateUser(username, email, passwordHash)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("user created")

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// HomePost handles the post requests to the login page
func HomePost(w http.ResponseWriter, r *http.Request) {
	userInterface := session.Get(r.Context(), "user")
	user, check := userInterface.(models.User)
	if !check {
		msg := "Login is required to access this page"
		http.Redirect(w, r, "/login?msg="+msg, http.StatusSeeOther)
	}

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	taskName := r.Form.Get("taskName")
	habitName := r.Form.Get("habitName")
	noteName := r.Form.Get("noteName")

	if taskName != "" {
		taskDuration, err := strconv.Atoi(r.Form.Get("taskDuration"))
		if err != nil {
			log.Println(err)
		}
		err = functions.CreateTask(taskName, taskDuration, false, user.ID)
		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else if habitName != "" {
		habitDescription := r.Form.Get("habitDescription")
		habitTimeStart, err := time.Parse("15:04", r.Form.Get("habitTimeStart"))
		if err != nil {
			log.Println(err)
		}
		habitTimeEnd, err := time.Parse("15:04", r.Form.Get("habitTimeEnd"))
		if err != nil {
			log.Println(err)
		}

		err = functions.CreateHabit(habitName, habitDescription, habitTimeStart, habitTimeEnd, user.ID)
		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else if noteName != "" {
		noteDescription := r.Form.Get("noteDescription")
		err = functions.CreateNote(noteName, noteDescription, user.ID)
		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
