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

// CreateTask creates a new task in the database
func CreateTask(name string, duration int, completedStatus bool, userID int) error {
	addTaskQuery := `insert into tasks(name, duration, completed_status, user_id, created_at, updated_at) values($1, $2, $3, $4, $5, $6)`
	_, err := db.Exec(addTaskQuery, name, duration, completedStatus, userID, time.Now(), time.Now())
	if err != nil {
		return err
	}
	return nil
}

// CreateHabit creates a habit task in the database
func CreateHabit(name string, description string, time_start time.Time, time_end time.Time, userID int) error {
	duration := time_end.Sub(time_start)
	addHabitQuery := `insert into habits(name, description, time_start, time_end, duration, user_id, created_at, updated_at) values($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := db.Exec(addHabitQuery, name, description, time_start, time_end, duration, userID, time.Now(), time.Now())
	if err != nil {
		return err
	}
	return nil
}

// CreateNote creates a new note in the database
func CreateNote(name string, description string, userID int) error {
	addNoteQuery := `insert into notes(name, description, user_id, created_at, updated_at) values($1, $2, $3, $4, $5)`
	_, err := db.Exec(addNoteQuery, name, description, userID, time.Now(), time.Now())
	if err != nil {
		return err
	}
	return nil
}
