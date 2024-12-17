package models

import "time"

type User struct {
	ID			int64		`json:"id"`
	Username	string		`json:"username"`
	Email		string 		`json:"email"`
	Password	string		`json:"-"`
	CreatedAt   time.Time 	`json:"created_at"`
}

type UserDTO struct {
	ID			int64		`json:"id"`
	Username	string		`json:"username"`
	Email		string 		`json:"email"`
}