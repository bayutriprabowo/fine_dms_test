package usecase

import (
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
	if id == 0 || path == "" || ext == "" {
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

func (self *file) SearchByUserId(userId int, query string) ([]model.File, error) {
	// if userId == 0 {
	// 	return nil, ErrInvalidUserID
	// }

	if query == "" {
		return nil, ErrInvalidQuery
	}

	files, err := self.fileRepo.SearchById(userId, query)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (self *file) SearchByName(name string) ([]model.File, error) {
	if name == "" {
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
