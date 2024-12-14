package functions

import (
	"database/sql"
	"log"
	"time"

	"github.com/srisudarshanrg/go-todo-list/server/models"
)

var db *sql.DB

// DBAccess provides the handlers with access to the database
func DBAccessFunctions(dbAccess *sql.DB) {
	db = dbAccess
}

// CreateUser creates a new user in the database
func CreateUser(username string, email string, passwordHash string) error {
	createUserQuery := `insert into users(username, email, password, join_date, created_at, updated_at) values($1, $2, $3, $4, $5, $6)`
	_, err := db.Exec(createUserQuery, username, email, passwordHash, time.Now().Format("02-01-2006"), time.Now(), time.Now())
	if err != nil {
		log.Println(err)
		return nil
	}

	return nil
}

// AuthenticateUser authenticates a user trying to login
func AuthenticateUser(credential string, password string) (bool, string, models.User, error) {
	getUserQuery := `select * from users where username=$1 or email=$1`
	result, err := db.Exec(getUserQuery, credential)
	if err != nil {
		log.Println(err)
		return false, "", models.User{}, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return false, "", models.User{}, err
	}

	if affected == 0 {
		return false, "from here: Invalid Credentials", models.User{}, nil
	}

	row := db.QueryRow(getUserQuery, credential)

	var id int
	var username, email, correctPasswordHash, joinDate string
	var createdAt, updatedAt time.Time
	err = row.Scan(&id, &username, &email, &correctPasswordHash, &joinDate, &createdAt, &updatedAt)
	if err != nil {
		log.Println(err)
		return false, "", models.User{}, err
	}

	check := CheckPasswordHash(password, correctPasswordHash)
	if !check {
		return false, "Invalid Credentials", models.User{}, nil
	}

	user := models.User{
		ID:        id,
		Username:  username,
		Email:     email,
		Password:  correctPasswordHash,
		JoinDate:  joinDate,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return true, "Successfully logged in", user, err
}
