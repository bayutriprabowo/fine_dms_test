package model

import "time"

type FileDownload struct {
	ID        int       `json:"id"`
	File      File      `json:"file" validate:"required"`
	User      User      `json:"user" validate:"required"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
}
