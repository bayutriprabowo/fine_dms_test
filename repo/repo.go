package repo

import (
	"errors"

	"enigmacamp.com/fine_dms/model"
)

var (
	ErrRepoNoData     = errors.New("no data")
	ErrRepoNoSuchData = errors.New("no such data")
	ErrRepoAlready    = errors.New("already exists")
)

type UserRepo interface {
	SelectAll() ([]model.User, error)
	SelectById(id int) (*model.User, error)
	SelectByUsername(uname string) (*model.User, error)
	SelectByEmail(email string) (*model.User, error)
	Create(user *model.User) error
	Update(user *model.User) error
	Delete(id int) error
}

type FileRepo interface {
	SelectAll() ([]model.File, error)
	Create(file *model.File) error
	Update(file *model.File) error
	Delete(id int) error
}

type TagsRepo interface {
	SelectAll() ([]model.Tags, error)
	SelectById(id int) (*model.Tags, error)
	SelectByName(name string) (*model.Tags, error)
	Create(tag *model.Tags) error
}

type TrxTagsRepo interface {
	SelectAll() ([]model.TrxTags, error)
	Create(trxTags *model.TrxTags) error
	Update(trxTags *model.TrxTags) error
	Delete(id int) error
}

type FileDownloadRepo interface {
	SelectAll() ([]model.FileDownload, error)
	Create(fileDl *model.FileDownload) error
	Update(fileDl *model.FileDownload) error
	Delete(id int) error
}

type FileUploadRepo interface {
	SelectAll() ([]model.FileUpload, error)
	Create(fileUp *model.FileUpload) error
	Update(fileUp *model.FileUpload) error
	Delete(id int) error
}
