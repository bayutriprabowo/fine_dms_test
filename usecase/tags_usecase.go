package usecase

import (
	"enigmacamp.com/fine_dms/model"
	"enigmacamp.com/fine_dms/repo"
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
	// TODO: validation
	return self.tagsRepo.Create(tag)
}
