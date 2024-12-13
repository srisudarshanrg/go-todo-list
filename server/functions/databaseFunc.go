package functions

import "database/sql"

var db *sql.DB

// DBAccess provides the handlers with access to the database
func DBAccessFunctions(dbAccess *sql.DB) {
	db = dbAccess
}

// CreateUser creates a new user in the database
func CreateUser(username string, email string, password string) (bool, error) {
	return true, nil
}

// AuthenticateUser authenticates a user trying to login
func AuthenticateUser(credential string, password string) (bool, error) {
	return true, nil
}
