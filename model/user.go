package model

import "time"

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username" validate:"required"`
	Password  string    `json:"password" validate:"required"`
	Email     string    `json:"email" validate:"required"`
	FirstName string    `json:"first_name" validate:"required"`
	LastName  string    `json:last_name"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
	UpdatedAt time.Time `json:"updated_at" validate:"required"`
}
