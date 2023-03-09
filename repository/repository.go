package repository

import (
	"database/sql"

	"enigmacamp.com/fine_dms/model"
	psql "enigmacamp.com/fine_dms/repository/psql"
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

type Repository struct {
	User         UserRepository
	File         FileRepository
	Tag          TagsRepository
	FileDownload FileDownloadRepository
	FileUpload   FileUploadRepository
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &psql.UserRepository{DB: db}
}

func NewFileRepository(db *sql.DB) FileRepository {
	return &psql.FileRepository{DB: db}
}

func NewTagsRepository(db *sql.DB) TagsRepository {
	return &psql.TagsRepository{DB: db}
}

func NewFileDownloadRepository(db *sql.DB) FileDownloadRepository {
	return &psql.FileDownloadRepository{DB: db}
}

func NewFileUploadRepository(db *sql.DB) FileUploadRepository {
	return &psql.FileUploadRepository{DB: db}
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		User:         NewUserRepository(db),
		File:         NewFileRepository(db),
		Tag:          NewTagsRepository(db),
		FileDownload: NewFileDownloadRepository(db),
		FileUpload:   NewFileUploadRepository(db),
	}
}
