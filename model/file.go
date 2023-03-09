package model

import "time"

type File struct {
	ID        int       `json:"id"`
	Path      string    `json:"path" validate:"required"`
	Ext       string    `json:"ext" validate:"required"`
	User      User      `json:"user" validate:"required"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
	UpdatedAt time.Time `json:"updated_at" validate:"required"`
}
