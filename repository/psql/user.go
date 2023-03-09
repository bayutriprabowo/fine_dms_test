package repository

import (
	"database/sql"
	"time"

	"enigmacamp.com/fine_dms/model"
	"enigmacamp.com/fine_dms/repository"
)

type userRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepository{DB: db}
}

func (repo *userRepository) SelectUser() ([]model.User, error) {
	rows, err := repo.DB.Query("SELECT id, username, password, email, first_name, last_name, created_at, updated_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []model.User{}
	for rows.Next() {
		user := model.User{}
		err = rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if len(users) == 0 {
		return nil, repository.ErrNoData
	}

	return users, nil
}

func (repo *userRepository) CreateUser(user *model.User) error {
	stmt, err := repo.DB.Prepare("INSERT INTO users(username, password, email, first_name, last_name, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err = stmt.Exec(user.Username, user.Password, user.Email, user.FirstName, user.LastName, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (repo *userRepository) UpdateUser(user *model.User) error {
	stmt, err := repo.DB.Prepare("UPDATE users SET username=?, password=?, email=?, first_name=?, last_name=?, updated_at=? WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	user.UpdatedAt = time.Now()

	_, err = stmt.Exec(user.Username, user.Password, user.Email, user.FirstName, user.LastName, user.UpdatedAt, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *userRepository) DeleteUser(id int) error {
	stmt, err := repo.DB.Prepare("DELETE FROM users WHERE id=?")
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
