package psql

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	"enigmacamp.com/fine_dms/model"
	"enigmacamp.com/fine_dms/repo"
)

type file struct {
	db *sql.DB
}

func NewPsqlFileRepo(db *sql.DB) repo.FileRepo {
	return &file{db}
}

func (self *file) SelectAllByUserId(id int) ([]model.File, error) {
	rows, err := self.db.Query(`
		SELECT	 id
			,path
			,ext
			,user_id
			,created_at
			,updated_at
		FROM m_file
		WHERE user_id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	files := []model.File{}
	for rows.Next() {
		file := model.File{}
		err = rows.Scan(&file.ID, &file.Path, &file.Ext,
			&file.User.ID, &file.CreatedAt, &file.UpdatedAt)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	if len(files) == 0 {
		return nil, repo.ErrRepoNoData
	}

	return files, nil
}

func (self *file) Create(file *model.File) error {
	stmt, err := self.db.Prepare(`
		INSERT INTO m_file(
			 path
			,ext
			,user
			,created_at
			,updated_at
		)
		VALUES ($1, $2, $3, $4, $5)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	file.CreatedAt = time.Now()
	file.UpdatedAt = file.CreatedAt

	_, err = stmt.Exec(file.Path, file.Ext, file.User.ID, file.CreatedAt,
		file.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (self *file) Update(file *model.File) error {
	stmt, err := self.db.Prepare(`
		UPDATE m_file SET
			 path=$1
			,ext=$2
			,user=$3
			,updated_at=$4
		WHERE id=$5
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	file.UpdatedAt = time.Now()

	_, err = stmt.Exec(file.Path, file.Ext, file.User.ID, file.UpdatedAt,
		file.ID)
	if err != nil {
		return err
	}

	return nil
}

// BROKEN

func (self *file) Delete(id int) error {
	stmt, err := self.db.Prepare(`
		DELETE FROM m_file WHERE id = $1
	`)
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

func (self *file) SearchById(userId int, query string) ([]model.File, error) {
	query = "%" + query + "%"
	files := []model.File{}
	rows, err := self.db.Query(`
		SELECT f.id, f.path, f.ext, f.user_id, f.created_at, f.updated_at
		FROM m_file f
		JOIN m_user u ON f.user_id = u.id
		WHERE f.id = $1	
    `, userId, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		file := model.File{}
		err = rows.Scan(&file.ID, &file.Path, &file.Ext,
			&file.User.ID, &file.CreatedAt, &file.UpdatedAt)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	if len(files) == 0 {
		return nil, repo.ErrRepoNoData
	}

	return files, nil
}

func (self *file) SearchByName(name string) ([]model.File, error) {
	files := []model.File{}
	rows, err := self.db.Query(`
        SELECT id, path, ext, user_id, created_at, updated_at
        FROM m_file
        WHERE path ILIKE $1
    `, "%"+name+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		file := model.File{}
		err = rows.Scan(&file.ID, &file.Path, &file.Ext,
			&file.User.ID, &file.CreatedAt, &file.UpdatedAt)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	if len(files) == 0 {
		return nil, repo.ErrRepoNoData
	}

	return files, nil
}

func (self *file) SearchByTags(tags []string) ([]model.File, error) {
	// create a placeholder string for each tag
	stmt := make([]string, len(tags))
	args := make([]interface{}, len(tags))

	for i, tag := range tags {
		stmt[i] = "$" + strconv.Itoa(i+1)
		args[i] = tag
	}

	query := `
		SELECT f.id, f.path, f.ext, f.user_id, f.created_at, f.updated_at
		FROM m_file f
		JOIN t_tags tt ON tt.file_id = f.id
		JOIN m_tags t ON tt.tags_id = t.id
		WHERE t.name IN (` + strings.Join(stmt, ", ") + `)
		GROUP BY f.id
	`

	files := []model.File{}
	rows, err := self.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		file := model.File{}
		err = rows.Scan(&file.ID, &file.Path, &file.Ext,
			&file.User.ID, &file.CreatedAt, &file.UpdatedAt)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	if len(files) == 0 {
		return nil, repo.ErrRepoNoData
	}

	return files, nil
}
