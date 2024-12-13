package models

import "time"

// User is a model for a user object in the database
type User struct {
	ID        int
	Username  string
	Email     string
	Password  string
	JoinDate  time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Task is the model for a task object in the database
type Task struct {
	ID              int
	Name            string
	Duration        time.Duration
	CompletedStatus bool
	UserID          int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Note is the model for a note object in the database
type Note struct {
	ID          int
	Name        string
	Description string
	UserID      int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Habit is the model for a habit object in the database
type Habit struct {
	ID          int
	Name        string
	Description string
	TimeStart   time.Time
	TimeEnd     time.Time
	Duration    time.Duration
	UserID      int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
