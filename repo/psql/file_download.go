package psql

import (
	"database/sql"
	"time"

	"enigmacamp.com/fine_dms/model"
	"enigmacamp.com/fine_dms/repo"
)

type fileDl struct {
	db *sql.DB
}

func NewPsqlFileDownloadRepo(db *sql.DB) repo.FileDownloadRepo {
	return &fileDl{db}
}

func (self *fileDl) SelectAll() ([]model.FileDownload, error) {
	rows, err := self.db.Query(`
		SELECT	 id
			,file_id
			,user_id
			,created_at
		FROM t_file_download
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	downloads := []model.FileDownload{}
	for rows.Next() {
		download := model.FileDownload{}
		err = rows.Scan(&download.ID, &download.File.ID,
			&download.User.ID, &download.CreatedAt)
		if err != nil {
			return nil, err
		}
		downloads = append(downloads, download)
	}

	if len(downloads) == 0 {
		return nil, repo.ErrRepoNoData
	}

	return downloads, nil
}

func (self *fileDl) Create(fileDl *model.FileDownload) error {
	stmt, err := self.db.Prepare(`
		INSERT INTO t_file_download(
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

	fileDl.CreatedAt = time.Now()

	_, err = stmt.Exec(fileDl.File.ID, fileDl.User.ID, fileDl.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (self *fileDl) Update(fileDl *model.FileDownload) error {
	stmt, err := self.db.Prepare(`
		UPDATE t_file_download SET
			 file_id=?
			,user_id=?
			,created_at=?
		WHERE id=?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(fileDl.File.ID, fileDl.User.ID, fileDl.CreatedAt,
		fileDl.ID)
	if err != nil {
		return err
	}

	return nil
}

func (self *fileDl) Delete(id int) error {
	stmt, err := self.db.Prepare("DELETE FROM t_file_download WHERE id=?")
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
