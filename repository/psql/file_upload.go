package repository

import (
	"database/sql"
	"errors"
	"time"

	"enigmacamp.com/fine_dms/model"
)

type FileUploadRepository struct {
	DB *sql.DB
}

func (repo *FileUploadRepository) Select() ([]model.FileUpload, error) {
	rows, err := repo.DB.Query("SELECT id, file_id, user_id, date FROM file_uploads")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fileUploads := []model.FileUpload{}
	for rows.Next() {
		fileUpload := model.FileUpload{}
		err = rows.Scan(&fileUpload.ID, &fileUpload.File.ID, &fileUpload.User.ID, &fileUpload.Date)
		if err != nil {
			return nil, err
		}
		fileUploads = append(fileUploads, fileUpload)
	}

	if len(fileUploads) == 0 {
		return nil, errors.New("no file upload found")
	}

	return fileUploads, nil
}

func (repo *FileUploadRepository) Create(fileUpload *model.FileUpload) error {
	stmt, err := repo.DB.Prepare("INSERT INTO file_uploads(file_id, user_id, date) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	fileUpload.Date = time.Now()

	_, err = stmt.Exec(fileUpload.File.ID, fileUpload.User.ID, fileUpload.Date)
	if err != nil {
		return err
	}

	return nil
}

func (repo *FileUploadRepository) Update(fileUpload *model.FileUpload) error {
	stmt, err := repo.DB.Prepare("UPDATE file_uploads SET file_id=?, user_id=?, date=? WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(fileUpload.File.ID, fileUpload.User.ID, fileUpload.Date, fileUpload.ID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *FileUploadRepository) Delete(id int) error {
	stmt, err := repo.DB.Prepare("DELETE FROM file_uploads WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}
