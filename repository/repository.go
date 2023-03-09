package repository

import (
	"enigmacamp.com/fine_dms/model"
)

type UserRepository interface {
	SelectUser() ([]model.User, error)
	CreateUser(user *model.User) error
	UpdateUser(user *model.User) error
	DeleteUser(id int) error
}

type FileRepository interface {
	Select() ([]model.File, error)
	Create(file *model.File) error
	Update(file *model.File) error
	Delete(id int) error
}

type TagsRepository interface {
	Select() ([]model.Tags, error)
	Create(tag *model.Tags) error
	Update(tag *model.Tags) error
	Delete(id int) error
}

type FileDownloadRepository interface {
	Select() ([]model.FileDownload, error)
	Create(fileDownload *model.FileDownload) error
	Update(fileDownload *model.FileDownload) error
	Delete(id int) error
}

type FileUploadRepository interface {
	Select() ([]model.FileUpload, error)
	Create(fileUpload *model.FileUpload) error
	Update(fileUpload *model.FileUpload) error
	Delete(id int) error
}
