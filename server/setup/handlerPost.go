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

	http.Redirect(w, r, "/home?msg="+msg, http.StatusSeeOther)
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

// HomePost handles the post requests to the home page
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
	addCheckID := r.Form.Get("addCheckID")
	removeCheckID := r.Form.Get("removeCheckID")
	deleteTaskID := r.Form.Get("deleteTaskID")
	deleteHabitID := r.Form.Get("deleteHabitID")
	searchTasks := r.Form.Get("searchTasks")
	searchHabits := r.Form.Get("searchHabits")
	noteIDEdit := r.Form.Get("noteIDEdit")
	searchNotes := r.Form.Get("searchNotes")
	deleteNoteID := r.Form.Get("deleteNoteID")

	if taskName != "" {
		taskDuration, err := strconv.Atoi(r.Form.Get("taskDuration"))
		if err != nil {
			log.Println(err)
		}
		err = functions.CreateTask(taskName, taskDuration, false, user.ID)
		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/home", http.StatusSeeOther)
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
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	} else if noteName != "" {
		noteDescription := r.Form.Get("noteDescription")
		err = functions.CreateNote(noteName, noteDescription, user.ID)
		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	} else if addCheckID != "" {
		addCheckIDConverted, err := strconv.Atoi(addCheckID)
		if err != nil {
			log.Println(err)
		}
		err = functions.AddCheckForTask(addCheckIDConverted, user.ID)
		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	} else if removeCheckID != "" {
		removeCheckIDConverted, err := strconv.Atoi(removeCheckID)
		if err != nil {
			log.Println(err)
		}
		err = functions.RemoveCheckForTask(removeCheckIDConverted, user.ID)
		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	} else if deleteTaskID != "" {
		deleteTaskIDConverted, err := strconv.Atoi(deleteTaskID)
		if err != nil {
			log.Println(err)
		}
		err = functions.DeleteTask(deleteTaskIDConverted, user.ID)
		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	} else if deleteHabitID != "" {
		deleteHabitIDConverted, err := strconv.Atoi(deleteHabitID)
		if err != nil {
			log.Println(err)
		}
		err = functions.DeleteHabit(deleteHabitIDConverted, user.ID)
		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	} else if searchTasks != "" {
		log.Println(searchTasks)
		results, err := functions.SearchTask(searchTasks, user.ID)
		if err != nil {
			log.Println(err)
		}
		postData := map[string]interface{}{}
		postData["searchResultsTasks"] = results
		RenderTemplate(w, r, "home.page.tmpl", models.TemplateData{
			Data:     data,
			PostData: postData,
		})
	} else if searchHabits != "" {
		results, err := functions.SearchHabit(searchHabits, user.ID)
		if err != nil {
			log.Println(err)
		}
		postData := map[string]interface{}{}
		postData["searchResultsHabits"] = results
		RenderTemplate(w, r, "home.page.tmpl", models.TemplateData{
			Data:     data,
			PostData: postData,
		})
	} else if noteIDEdit != "" {
		noteIDEditConverted, err := strconv.Atoi(noteIDEdit)
		if err != nil {
			log.Println(err)
		}
		noteNameEdit := r.Form.Get("noteNameEdit")
		noteDescriptionEdit := r.Form.Get("noteDescriptionEdit")
		err = functions.UpdateNote(noteIDEditConverted, noteNameEdit, noteDescriptionEdit, user.ID)
		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	} else if searchNotes != "" {
		results, err := functions.SearchNote(searchNotes, user.ID)
		if err != nil {
			log.Println(err)
		}
		postData := map[string]interface{}{}
		postData["searchResultsNotes"] = results
		RenderTemplate(w, r, "home.page.tmpl", models.TemplateData{
			Data:     data,
			PostData: postData,
		})
	} else if deleteNoteID != "" {
		deleteNoteIDConverted, err := strconv.Atoi(deleteNoteID)
		if err != nil {
			log.Println(err)
		}
		err = functions.DeleteNote(deleteNoteIDConverted, user.ID)
		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
}

// TasksPost handles the post requests to the tasks page
func TasksPost(w http.ResponseWriter, r *http.Request) {
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

	link := session.Get(r.Context(), "linkTasks").(string)
	path := session.Get(r.Context(), "pathTasks").(string)

	taskName := r.Form.Get("taskName")
	addCheckID := r.Form.Get("addCheckID")
	removeCheckID := r.Form.Get("removeCheckID")
	deleteTaskID := r.Form.Get("deleteTaskID")
	searchTasks := r.Form.Get("searchTasks")

	if taskName != "" {
		taskDuration, err := strconv.Atoi(r.Form.Get("taskDuration"))
		if err != nil {
			log.Println(err)
		}
		err = functions.CreateTask(taskName, taskDuration, false, user.ID)
		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, link, http.StatusSeeOther)
	} else if addCheckID != "" {
		addCheckIDConverted, err := strconv.Atoi(addCheckID)
		if err != nil {
			log.Println(err)
		}
		err = functions.AddCheckForTask(addCheckIDConverted, user.ID)
		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, link, http.StatusSeeOther)
	} else if removeCheckID != "" {
		removeCheckIDConverted, err := strconv.Atoi(removeCheckID)
		if err != nil {
			log.Println(err)
		}
		err = functions.RemoveCheckForTask(removeCheckIDConverted, user.ID)
		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, link, http.StatusSeeOther)
	} else if deleteTaskID != "" {
		deleteTaskIDConverted, err := strconv.Atoi(deleteTaskID)
		if err != nil {
			log.Println(err)
		}
		err = functions.DeleteTask(deleteTaskIDConverted, user.ID)
		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, link, http.StatusSeeOther)
	} else if searchTasks != "" {
		results, err := functions.SearchTask(searchTasks, user.ID)
		if err != nil {
			log.Println(err)
		}

		postData := map[string]interface{}{}
		postData["searchResultsTasks"] = results
		RenderTemplate(w, r, path, models.TemplateData{
			Data:     data,
			PostData: postData,
		})
	}

	session.Remove(r.Context(), "linkTasks")
	session.Remove(r.Context(), "pathTasks")
}

// HabitsPost handles the post requests to the habits page
func HabitsPost(w http.ResponseWriter, r *http.Request) {
	userInterface := session.Get(r.Context(), "user")
	user, check := userInterface.(models.User)
	if !check {
		msg := "Login is required to access this page"
		http.Redirect(w, r, "/login?msg="+msg, http.StatusSeeOther)
	}

	link := session.Get(r.Context(), "linkHabits").(string)
	path := session.Get(r.Context(), "pathHabits").(string)

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	deleteHabitID := r.Form.Get("deleteHabitID")
	habitName := r.Form.Get("habitName")
	searchHabits := r.Form.Get("searchHabits")

	if deleteHabitID != "" {
		deleteHabitIDConverted, err := strconv.Atoi(deleteHabitID)
		if err != nil {
			log.Println(err)
		}
		err = functions.DeleteHabit(deleteHabitIDConverted, user.ID)
		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, link, http.StatusSeeOther)
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
		http.Redirect(w, r, link, http.StatusSeeOther)
	} else if searchHabits != "" {
		results, err := functions.SearchHabit(searchHabits, user.ID)
		if err != nil {
			log.Println(err)
		}

		postData := map[string]interface{}{}
		postData["searchResultsHabits"] = results
		RenderTemplate(w, r, path, models.TemplateData{
			Data:     data,
			PostData: postData,
		})
	}

	session.Remove(r.Context(), "linkHabits")
	session.Remove(r.Context(), "pathHabits")
}

// NotesPost handles the post requests to the notes page
func NotesPost(w http.ResponseWriter, r *http.Request) {
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

	noteName := r.Form.Get("noteName")
	searchNotes := r.Form.Get("searchNotes")
	noteIDEdit := r.Form.Get("noteIDEdit")
	deleteNoteID := r.Form.Get("deleteNoteID")

	if noteName != "" {
		noteDescription := r.Form.Get("noteDescription")
		err = functions.CreateNote(noteName, noteDescription, user.ID)
		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/notes", http.StatusSeeOther)
	} else if searchNotes != "" {
		results, err := functions.SearchNote(searchNotes, user.ID)
		if err != nil {
			log.Println(err)
		}

		postData := map[string]interface{}{}
		postData["searchResultsNotes"] = results
		RenderTemplate(w, r, "notes.page.tmpl", models.TemplateData{
			Data:     data,
			PostData: postData,
		})
	} else if noteIDEdit != "" {
		noteIDEditConverted, err := strconv.Atoi(noteIDEdit)
		if err != nil {
			log.Println(err)
		}
		noteNameEdit := r.Form.Get("noteNameEdit")
		noteDescriptionEdit := r.Form.Get("noteDescriptionEdit")

		err = functions.UpdateNote(noteIDEditConverted, noteNameEdit, noteDescriptionEdit, user.ID)
		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/notes", http.StatusSeeOther)
	} else if deleteNoteID != "" {
		deleteNoteIDConverted, err := strconv.Atoi(deleteNoteID)
		if err != nil {
			log.Println(err)
		}
		err = functions.DeleteNote(deleteNoteIDConverted, user.ID)
		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/notes", http.StatusSeeOther)
	}
}
