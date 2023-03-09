package repository

import (
	"database/sql"
	"errors"
	"time"

	"enigmacamp.com/fine_dms/model"
	"enigmacamp.com/fine_dms/repository"
)

type fileRepository struct {
	DB *sql.DB
}

func NewFileRepository(db *sql.DB) repository.FileRepository {
	return &fileRepository{DB: db}
}

func (repo *fileRepository) Select() ([]model.File, error) {
	rows, err := repo.DB.Query("SELECT id, path, ext, user_id, created_at, updated_at FROM files")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	files := []model.File{}
	for rows.Next() {
		file := model.File{}
		err = rows.Scan(&file.ID, &file.Path, &file.Ext, &file.User.ID, &file.CreatedAt, &file.UpdatedAt)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	if len(files) == 0 {
		return nil, errors.New("no file found")
	}

	return files, nil
}

func (repo *fileRepository) Create(file *model.File) error {
	stmt, err := repo.DB.Prepare("INSERT INTO files(path, ext, user_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	file.CreatedAt = time.Now()
	file.UpdatedAt = time.Now()

	_, err = stmt.Exec(file.Path, file.Ext, file.User.ID, file.CreatedAt, file.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (repo *fileRepository) Update(file *model.File) error {
	stmt, err := repo.DB.Prepare("UPDATE files SET path=?, ext=?, user_id=?, updated_at=? WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	file.UpdatedAt = time.Now()

	_, err = stmt.Exec(file.Path, file.Ext, file.User.ID, file.UpdatedAt, file.ID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *fileRepository) Delete(id int) error {
	stmt, err := repo.DB.Prepare("DELETE FROM files WHERE id=?")
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
