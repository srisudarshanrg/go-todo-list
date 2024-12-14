package setup

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/srisudarshanrg/go-setup-template/server/models"
)

var db *sql.DB
var session *scs.SessionManager

// DBAccess provides the handlers with access to the database
func DBAccessHandlers(dbAccess *sql.DB) {
	db = dbAccess
}

// SessionAccess provides the handlers with access to the sessions
func SessionAccessHandlers(sessionAccess *scs.SessionManager) {
	session = sessionAccess
}

// Login is the handler for the login page
func Login(w http.ResponseWriter, r *http.Request) {
	err := RenderTemplate(w, "login.page.tmpl", models.TemplateData{})
	if err != nil {
		log.Println(err)
	}
}

// Register is the handler for the register page
func Register(w http.ResponseWriter, r *http.Request) {
	err := RenderTemplate(w, "register.page.tmpl", models.TemplateData{})
	if err != nil {
		log.Println(err)
	}
}

// Home is the handler for the home page
func Home(w http.ResponseWriter, r *http.Request) {
	err := RenderTemplate(w, "home.page.tmpl", models.TemplateData{})
	if err != nil {
		log.Println(err)
	}
}
