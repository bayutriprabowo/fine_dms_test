package psql

import (
	"database/sql"
	"time"

	"enigmacamp.com/fine_dms/model"
	"enigmacamp.com/fine_dms/repo"
)

type fileUp struct {
	db *sql.DB
}

func NewPsqlFileUploadRepo(db *sql.DB) repo.FileUploadRepo {
	return &fileUp{db}
}

func (self *fileUp) SelectAll() ([]model.FileUpload, error) {
	rows, err := self.db.Query(`
		SELECT	 id
			,file_id
			,user_id
			,created_at
		FROM t_file_upload
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fileUploads := []model.FileUpload{}
	for rows.Next() {
		fileUpload := model.FileUpload{}
		err = rows.Scan(&fileUpload.ID, &fileUpload.File.ID,
			&fileUpload.User.ID, &fileUpload.CreatedAt)
		if err != nil {
			return nil, err
		}
		fileUploads = append(fileUploads, fileUpload)
	}

	if len(fileUploads) == 0 {
		return nil, repo.ErrRepoNoData
	}

	return fileUploads, nil
}

func (self *fileUp) Create(fileUp *model.FileUpload) error {
	stmt, err := self.db.Prepare(`
		INSERT INTO t_file_upload(
			 file_id
			,user_id
			,created_at
		)
		VALUES (?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	fileUp.CreatedAt = time.Now()

	_, err = stmt.Exec(fileUp.File.ID, fileUp.User.ID, fileUp.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (self *fileUp) Update(fileUp *model.FileUpload) error {
	stmt, err := self.db.Prepare(`
		UPDATE t_file_upload SET
			 file_id=?
			,user_id=?
			,created_at=?
		WHERE id=?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(fileUp.File.ID, fileUp.User.ID, fileUp.CreatedAt,
		fileUp.ID)
	if err != nil {
		return err
	}

	return nil
}

func (self *fileUp) Delete(id int) error {
	stmt, err := self.db.Prepare("DELETE FROM t_file_upload WHERE id=?")
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
