package repository

import (
	"database/sql"
	"errors"
	"time"

	"enigmacamp.com/fine_dms/model"
	"enigmacamp.com/fine_dms/repository"
)

type tagsRepository struct {
	DB *sql.DB
}

func NewTagsRepository(db *sql.DB) repository.TagsRepository {
	return &tagsRepository{DB: db}
}

func (repo *tagsRepository) Select() ([]model.Tags, error) {
	rows, err := repo.DB.Query("SELECT id, name, created_at, updated_at FROM tags")
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
		return nil, errors.New("no tags found")
	}

	return tags, nil
}

func (repo *tagsRepository) Create(tag *model.Tags) error {
	stmt, err := repo.DB.Prepare("INSERT INTO tags(name, created_at, updated_at) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	tag.CreatedAt = time.Now()
	tag.UpdatedAt = time.Now()

	_, err = stmt.Exec(tag.Name, tag.CreatedAt, tag.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (repo *tagsRepository) Update(tag *model.Tags) error {
	stmt, err := repo.DB.Prepare("UPDATE tags SET name=?, updated_at=? WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	tag.UpdatedAt = time.Now()

	_, err = stmt.Exec(tag.Name, tag.UpdatedAt, tag.ID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *tagsRepository) Delete(id int) error {
	stmt, err := repo.DB.Prepare("DELETE FROM tags WHERE id=?")
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
