package model

import "time"

type FileDownload struct {
	Id   int       `json:"id"`
	File File      `json:"file"`
	User User      `json:"user"`
	Date time.Time `json:"date"`
}
