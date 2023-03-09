package model

import "time"

type Tags struct {
	ID        int       `json:"id"`
	Name      string    `json:"name" validate:"required"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
	UpdatedAt time.Time `json:"updated_at" validate:"required"`
}
