package usecase

import (
	"errors"
	"regexp"
	"strings"

	"enigmacamp.com/fine_dms/model"
	"enigmacamp.com/fine_dms/repository"
)

type TagsUsecase interface {
	Select() ([]model.Tags, error)
	Create(tag *model.Tags) error
	Update(tag *model.Tags) error
	Delete(id int) error
}

type tagsUsecase struct {
	tagsRepo repository.TagsRepository
}

func (u *tagsUsecase) Select() ([]model.Tags, error) {
	return u.tagsRepo.Select()
}

func (u *tagsUsecase) Create(tag *model.Tags) error {
	if !isValidTag(tag.Name) {
		return errors.New("Tag name must only contain alphabetical characters and maximum of 5 words")
	}
	return u.tagsRepo.Create(tag)
}

func (u *tagsUsecase) Update(tag *model.Tags) error {
	if !isValidTag(tag.Name) {
		return errors.New("Tag name must only contain alphabetical characters and maximum of 5 words")
	}
	return u.tagsRepo.Update(tag)
}

func (u *tagsUsecase) Delete(id int) error {
	return u.tagsRepo.Delete(id)
}

func NewTagsUsecase(tagsRepo repository.TagsRepository) TagsUsecase {
	return &tagsUsecase{
		tagsRepo: tagsRepo,
	}
}

func isValidTag(tagName string) bool {
	if !regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString(tagName) {
		return false
	}

	words := strings.Fields(tagName)
	if len(words) > 5 {
		return false
	}
	return true
}
