package model

type TagsMan struct {
	ID   int  `json:"id"`
	Tags Tags `json:"tags"`
	File File `json:"file"`
}
