package usecase

import (
	"errors"

	"enigmacamp.com/fine_dms/model"
	"enigmacamp.com/fine_dms/repository"
	"enigmacamp.com/fine_dms/utils"
)

type tagsUsecase struct {
	tagsRepo repository.TagsRepository
}

func (usecase *tagsUsecase) Select() ([]model.Tags, error) {
	return usecase.tagsRepo.Select()
}

func (usecase *tagsUsecase) Create(tag *model.Tags) error {
	if !utils.IsValidTag(tag.Name) {
		return errors.New("Tag name must only contain alphabetical characters and maximum of 5 words")
	}
	return usecase.tagsRepo.Create(tag)
}

func (usecase *tagsUsecase) Update(tag *model.Tags) error {
	if !utils.IsValidTag(tag.Name) {
		return errors.New("Tag name must only contain alphabetical characters and maximum of 5 words")
	}
	return usecase.tagsRepo.Update(tag)
}

func (usecase *tagsUsecase) Delete(id int) error {
	return usecase.tagsRepo.Delete(id)
}

func NewTagsUsecase(tagsRepo repository.TagsRepository) TagsUsecase {
	return &tagsUsecase{
		tagsRepo: tagsRepo,
	}
}
