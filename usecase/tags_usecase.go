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
	return self.GetAll()
}
