package repository

import (
	"database/sql"
	"time"

	"enigmacamp.com/fine_dms/model"
	"enigmacamp.com/fine_dms/repository"
)

type fileDlRepository struct {
	DB *sql.DB
}

func NewFileDownloadRepository(db *sql.DB) repository.FileDownloadRepository {
	return &fileDlRepository{DB: db}
}

func (repo *fileDlRepository) Select() ([]model.FileDownload, error) {
	rows, err := repo.DB.Query("SELECT id, file_id, user_id, date FROM file_downloads")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	downloads := []model.FileDownload{}
	for rows.Next() {
		download := model.FileDownload{}
		err = rows.Scan(&download.ID, &download.File.ID, &download.User.ID, &download.Date)
		if err != nil {
			return nil, err
		}
		downloads = append(downloads, download)
	}

	if len(downloads) == 0 {
		return nil, repository.ErrNoData
	}

	return downloads, nil
}

func (repo *fileDlRepository) Create(download *model.FileDownload) error {
	stmt, err := repo.DB.Prepare("INSERT INTO file_downloads(file_id, user_id, date) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	download.Date = time.Now()

	_, err = stmt.Exec(download.File.ID, download.User.ID, download.Date)
	if err != nil {
		return err
	}

	return nil
}

func (repo *fileDlRepository) Update(download *model.FileDownload) error {
	stmt, err := repo.DB.Prepare("UPDATE file_downloads SET file_id=?, user_id=?, date=? WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(download.File.ID, download.User.ID, download.Date, download.ID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *fileDlRepository) Delete(id int) error {
	stmt, err := repo.DB.Prepare("DELETE FROM file_downloads WHERE id=?")
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
