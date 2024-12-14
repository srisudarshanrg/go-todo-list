package setup

import (
	"log"
	"net/http"
)

// LoginPost handles the post requests to the login page
func LoginPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
}

// RegisterPost handles the post requests to the login page
func RegisterPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
}

// HomePost handles the post requests to the login page
func HomePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
}
