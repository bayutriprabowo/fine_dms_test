package usecase

import (
	"errors"

	"enigmacamp.com/fine_dms/model"
	"enigmacamp.com/fine_dms/repo"
	"enigmacamp.com/fine_dms/utils"
)

var (
	ErrUsecaseInvalidTag = errors.New("`tag_name` must be alphanumeric " +
		"and has length at least (> 0 and <= 16)")
)

type tags struct {
	tagsRepo repo.TagsRepo
}

func NewTagsUsecase(tagsRepo repo.TagsRepo) TagsUsecase {
	return &tags{tagsRepo}
}

func (self *tags) GetAll() ([]model.Tags, error) {
	res, err := self.tagsRepo.SelectAll()
	if err == repo.ErrRepoNoData {
		return nil, ErrUsecaseNoData
	}

	return res, err
}

func (self *tags) GetById(id int) (*model.Tags, error) {
	res, err := self.tagsRepo.SelectById(id)
	if err == repo.ErrRepoNoData {
		return nil, ErrUsecaseNoData
	}

	return res, err
}

func (self *tags) GetByName(name string) (*model.Tags, error) {
	res, err := self.tagsRepo.SelectByName(name)
	if err == repo.ErrRepoNoData {
		return nil, ErrUsecaseNoData
	}

	return res, err
}

func (self *tags) Add(tag *model.Tags) error {
	if utils.IsValidTag(tag.Name) {
		return ErrUsecaseInvalidTag
	}

	return self.tagsRepo.Create(tag)
}
