package psql

import (
	"database/sql"

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
