package usecase

import "enigmacamp.com/fine_dms/model"

type TagsUsecase interface {
	Select() ([]model.Tags, error)
	Create(tag *model.Tags) error
	Update(tag *model.Tags) error
	Delete(id int) error
}

type UserUsecase interface {
	Select() ([]model.User, error)
	Create(user *model.User) error
	Update(user *model.User) error
	Delete(id int) error
}
