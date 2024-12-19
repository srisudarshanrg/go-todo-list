package functions

import (
	"database/sql"
	"log"
	"strings"
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
		return false, "Invalid Credentials", models.User{}, nil
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

// GetTasks gets all the tasks of a given user
func GetTasks(userID int) ([]models.Task, error) {
	getTasksQuery := `select * from tasks where user_id=$1`
	rows, err := db.Query(getTasksQuery, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var id, userID int
		var name string
		var completedStatus bool
		var duration int
		var createdAt, updatedAt time.Time

		err = rows.Scan(&id, &name, &duration, &completedStatus, &userID, &createdAt, &updatedAt)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		task := models.Task{
			ID:              id,
			Name:            name,
			Duration:        duration,
			CompletedStatus: completedStatus,
			UserID:          userID,
			CreatedAt:       createdAt,
			UpdatedAt:       updatedAt,
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
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

// AddCheckForTask markes the status of the task as completed in the database
func AddCheckForTask(id int, userID int) error {
	addCheckQuery := `update tasks set completed_status=$1 where id=$2 and user_id=$3`
	_, err := db.Exec(addCheckQuery, true, id, userID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// RemoveCheckForTask markes the status of the task as incomplete in the database
func RemoveCheckForTask(id int, userID int) error {
	addCheckQuery := `update tasks set completed_status=$1 where id=$2 and user_id=$3`
	_, err := db.Exec(addCheckQuery, false, id, userID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// SearchTask searches for a task in the database based on the key and userID
func SearchTask(key string, userID int) ([]models.Task, error) {
	keyLower := strings.ToLower(key)
	keyArg := "%" + keyLower + "%"
	searchTasksQuery := `select * from tasks where lower(name) like $1 and user_id=$2`
	rows, err := db.Query(searchTasksQuery, keyArg, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var searchResults []models.Task
	for rows.Next() {
		var id, duration, userID int
		var name string
		var completedStatus bool
		var createdAt, updatedAt time.Time

		err = rows.Scan(&id, &name, &duration, &completedStatus, &userID, &createdAt, &updatedAt)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		task := models.Task{
			ID:              id,
			Name:            name,
			Duration:        duration,
			CompletedStatus: completedStatus,
			CreatedAt:       createdAt,
			UpdatedAt:       updatedAt,
		}

		searchResults = append(searchResults, task)
	}

	return searchResults, nil
}

// DeleteTask deletes a task for a given id and userID from the database
func DeleteTask(id int, userID int) error {
	deleteTaskQuery := `delete from tasks where id=$1 and user_id=$2`
	_, err := db.Exec(deleteTaskQuery, id, userID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// GetHabits gets all the habits of a given user
func GetHabits(userID int) ([]models.Habit, error) {
	getHabitsQuery := `select * from habits where user_id=$1`
	rows, err := db.Query(getHabitsQuery, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var habits []models.Habit
	for rows.Next() {
		var id, userID int
		var name, description string
		var duration int
		var timeStart, timeEnd, createdAt, updatedAt time.Time

		err = rows.Scan(&id, &name, &description, &timeStart, &timeEnd, &duration, &userID, &createdAt, &updatedAt)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		habit := models.Habit{
			ID:          id,
			Name:        name,
			Description: description,
			TimeStart:   timeStart.Format("15:04"),
			TimeEnd:     timeEnd.Format("15:04"),
			Duration:    duration,
			UserID:      userID,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		}

		habits = append(habits, habit)
	}

	return habits, nil
}

// CreateHabit creates a habit task in the database
func CreateHabit(name string, description string, timeStart time.Time, timeEnd time.Time, userID int) error {
	duration := int32(timeEnd.Sub(timeStart).Minutes())
	addHabitQuery := `insert into habits(name, description, time_start, time_end, duration, user_id, created_at, updated_at) values($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := db.Exec(addHabitQuery, name, description, timeStart, timeEnd, duration, userID, time.Now(), time.Now())
	if err != nil {
		return err
	}
	return nil
}

// SearchHabit searches for a habit in the database based on the key and userID
func SearchHabit(key string, userID int) ([]models.Habit, error) {
	keyLower := strings.ToLower(key)
	keyArg := "%" + keyLower + "%"
	searchHabitsQuery := `select * from habits where lower(name) like $1 or lower(description) like $1 and user_id=$2`
	rows, err := db.Query(searchHabitsQuery, keyArg, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var searchResults []models.Habit
	for rows.Next() {
		var id, duration, userID int
		var name, description string
		var timeStart, timeEnd, createdAt, updatedAt time.Time

		err = rows.Scan(&id, &name, &description, &timeStart, &timeEnd, &duration, &userID, &createdAt, &updatedAt)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		habit := models.Habit{
			ID:          id,
			Name:        name,
			Description: description,
			TimeStart:   timeStart.Format("15:04"),
			TimeEnd:     timeEnd.Format("15:04"),
			Duration:    duration,
			UserID:      userID,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		}

		searchResults = append(searchResults, habit)
	}

	return searchResults, nil
}

// DeleteHabit deletes a habit for a given id and userID
func DeleteHabit(id int, userID int) error {
	deleteHabitQuery := `delete from habits where id=$1 and user_id=$2`
	_, err := db.Exec(deleteHabitQuery, id, userID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// GetNotes gets all the notes of a given user
func GetNotes(userID int) ([]models.Note, error) {
	getNotesQuery := `select * from notes where user_id=$1`
	rows, err := db.Query(getNotesQuery, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var id, userID int
		var name, description string
		var createdAt, updatedAt time.Time

		err = rows.Scan(&id, &name, &description, &userID, &createdAt, &updatedAt)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		note := models.Note{
			ID:          id,
			Name:        name,
			Description: description,
			UserID:      userID,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		}

		notes = append(notes, note)
	}

	return notes, nil
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

// UpdateNote updates new changes of a note in the database
func UpdateNote(id int, name string, description string, userID int) error {
	updateNoteQuery := `update notes set name=$1, description=$2 where id=$3 and user_id=$4`
	_, err := db.Exec(updateNoteQuery, name, description, id, userID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// SearchNote searches for a note in the database based on a key
func SearchNote(key string, userID int) ([]models.Note, error) {
	keyLower := strings.ToLower(key)
	keyArg := "%" + keyLower + "%"
	searchNoteQuery := `select * from notes where name like $1 or description like $1 and user_id=$2`
	rows, err := db.Query(searchNoteQuery, keyArg, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var searchResults []models.Note
	for rows.Next() {
		var id, userID int
		var name, description string
		var createdAt, updatedAt time.Time

		err = rows.Scan(&id, &name, &description, &userID, &createdAt, &updatedAt)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		note := models.Note{
			ID:          id,
			Name:        name,
			Description: description,
			UserID:      userID,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		}

		searchResults = append(searchResults, note)
	}

	return searchResults, nil
}

// DeleteNote deletes a note for a given id and userID
func DeleteNote(id int, userID int) error {
	deleteNoteQuery := `delete from notes where id=$1 and user_id=$2`
	_, err := db.Exec(deleteNoteQuery, id, userID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
