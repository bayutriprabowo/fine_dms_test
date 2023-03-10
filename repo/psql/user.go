package psql

import (
	"database/sql"
	"time"

	"enigmacamp.com/fine_dms/model"
	"enigmacamp.com/fine_dms/repo"
)

type user struct {
	db *sql.DB
}

func NewPsqlUserRepo(db *sql.DB) repo.UserRepo {
	return &user{db}
}

func (self *user) SelectAll() ([]model.User, error) {
	rows, err := self.db.Query(`
		SELECT 	 id
			,username
			,password
			,email
			,first_name
			,last_name
			,created_at
			,updated_at
		FROM m_user 
	`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []model.User{}
	for rows.Next() {
		user := model.User{}
		err = rows.Scan(&user.ID, &user.Username, &user.Password,
			&user.Email, &user.FirstName, &user.LastName,
			&user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if len(users) == 0 {
		return nil, repo.ErrRepoNoData
	}

	return users, nil
}

func (self *user) SelectById(id int) (*model.User, error) {
	return self.selectOne("id", id)
}

func (self *user) SelectByUsername(uname string) (*model.User, error) {
	return self.selectOne("username", uname)
}

func (self *user) SelectByEmail(email string) (*model.User, error) {
	return self.selectOne("email", email)
}

func (self *user) Create(user *model.User) error {
	stmt, err := self.db.Prepare(`
		INSERT INTO m_user(
			 username
			,password
			,email
			,first_name
			,last_name
			,created_at
			,updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt

	_, err = stmt.Exec(user.Username, user.Password, user.Email,
		user.FirstName, user.LastName, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (self *user) Update(user *model.User) error {
	stmt, err := self.db.Prepare(`
		UPDATE m_user SET
			 username=$1
			,password=$2
			,email=$3
			,first_name=$4
			,last_name=$5
			,updated_at=$6
		WHERE id=$7
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(user.Username, user.Password, user.Email,
		user.FirstName, user.LastName, time.Now(), user.ID)
	if err != nil {
		return err
	}

	rw, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rw == 0 {
		return repo.ErrRepoNoSuchData
	}

	return nil
}

func (self *user) Delete(id int) error {
	stmt, err := self.db.Prepare("DELETE FROM m_user WHERE id=$1")
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	rw, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rw == 0 {
		return repo.ErrRepoNoSuchData
	}

	return nil
}

// private
func (self *user) selectOne(key string, val any) (*model.User, error) {
	var q = `
		SELECT	 id
			,username
			,password
			,email
			,first_name
			,last_name
			,created_at
			,updated_at
		FROM m_user
		WHERE 
	`

	q += (key + " = $1")

	row := self.db.QueryRow(q, val)

	var ret = new(model.User)

	err := row.Scan(&ret.ID, &ret.Username, &ret.Password, &ret.Email,
		&ret.FirstName, &ret.LastName, &ret.CreatedAt, &ret.UpdatedAt,
	)
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
