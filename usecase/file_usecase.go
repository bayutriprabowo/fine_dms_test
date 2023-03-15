package usecase

import (
	"strings"
	"time"

	"enigmacamp.com/fine_dms/model"
	"enigmacamp.com/fine_dms/repo"
)

type file struct {
	fileRepo repo.FileRepo
}

func NewFileUsecase(fileRepo repo.FileRepo) FileUsecase {
	return &file{fileRepo: fileRepo}
}

func (self *file) GetFilesByUserId(id int) ([]model.File, error) {
	files, err := self.fileRepo.SelectAllByUserId(id)
	if err != nil {
		if err == repo.ErrRepoNoData {
			return nil, ErrUsecaseNoData
		}
		return nil, err
	}

	return files, nil
}

func (self *file) UpdateFile(id int, path, ext string) error {
	if len(path) == 0 || len(ext) == 0 {
		return ErrInvalidFileData
	}

	files, err := self.fileRepo.SelectAllByUserId(id)
	if err != nil {
		if err == repo.ErrRepoNoData {
			return ErrUsecaseNoData
		}
		return err
	}

	if len(files) == 0 {
		return ErrUsecaseNoData
	}

	file := &files[0]
	file.Path = path
	file.Ext = ext
	file.UpdatedAt = time.Now()

	if err := self.fileRepo.Update(file); err != nil {
		return err
	}

	return nil
}

func (self *file) DeleteFile(userId int, fileId int) error {
	files, err := self.fileRepo.SelectAllByUserId(fileId)
	if err != nil {
		if err == repo.ErrRepoNoData {
			return ErrUsecaseNoData
		}
		return err
	}

	var file *model.File
	for _, f := range files {
		if f.ID == fileId {
			file = &f
			break
		}
	}

	if file == nil {
		return ErrUsecaseNoData
	}

	if file.User.ID != userId {
		return ErrUsecaseInvalidAuth
	}

	if err := self.fileRepo.Delete(file.ID); err != nil {
		return err
	}

	return nil
}

func (self *file) SearchByUserId(userID int, query string) ([]model.File, error) {
	files, err := self.GetFilesByUserId(userID)
	if err != nil {
		return nil, err
	}

	if query == "" {
		return files, nil
	}

	filteredFiles := make([]model.File, 0)
	for _, file := range files {
		if strings.Contains(strings.ToLower(file.User.Username), strings.ToLower(query)) {
			filteredFiles = append(filteredFiles, file)
		}
	}

	if len(filteredFiles) == 0 {
		return nil, ErrInvalidUserID
	}

	return filteredFiles, nil
}

func (self *file) SearchByName(name string) ([]model.File, error) {
	if len(name) == 0 {
		return nil, ErrInvalidQuery
	}

	files, err := self.fileRepo.SearchByName(name)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (self *file) SearchByTags(tags []string) ([]model.File, error) {
	if len(tags) == 0 {
		return nil, ErrUsecaseInvalidTag
	}

	files, err := self.fileRepo.SearchByTags(tags)
	if err != nil {
		return nil, err
	}

	return files, nil
}
