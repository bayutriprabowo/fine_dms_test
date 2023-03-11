package psql

import (
	"database/sql"
	"time"

	"enigmacamp.com/fine_dms/model"
	"enigmacamp.com/fine_dms/repo"
)

type tags struct {
	db *sql.DB
}

func NewPsqlTagsRepo(db *sql.DB) repo.TagsRepo {
	return &tags{db}
}

func (self *tags) SelectAll() ([]model.Tags, error) {
	rows, err := self.db.Query(`
		SELECT	 id
			,name
			,created_at
			,updated_at
		FROM m_tags
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := []model.Tags{}
	for rows.Next() {
		tag := model.Tags{}
		err = rows.Scan(&tag.ID, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	if len(tags) == 0 {
		return nil, repo.ErrRepoNoData
	}

	return tags, nil
}

func (self *tags) SelectById(id int) (*model.Tags, error) {
	return self.selectOne("id", id)
}

func (self *tags) SelectByName(name string) (*model.Tags, error) {
	return self.selectOne("name", name)
}

func (self *tags) Create(tag *model.Tags) error {
	stmt, err := self.db.Prepare(`
		INSERT INTO m_tags(
			 name
			,created_at
			,updated_at
		)
		VALUES ($1, $2, $3)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	tag.CreatedAt = time.Now()
	tag.UpdatedAt = tag.CreatedAt

	_, err = stmt.Exec(tag.Name, tag.CreatedAt, tag.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

// private
func (self *tags) selectOne(key string, val any) (*model.Tags, error) {
	var q = `
		SELECT	 id
			,name
			,created_at
			,updated_at
		FROM m_tags
		WHERE  
	`

	q += (key + " = $1")

	row := self.db.QueryRow(q, val)

	var ret = new(model.Tags)

	err := row.Scan(&ret.ID, &ret.Name, &ret.CreatedAt, &ret.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repo.ErrRepoNoData
		}

		return nil, err
	}

	if err = row.Err(); err != nil {
		return nil, err
	}

	return ret, nil
}
