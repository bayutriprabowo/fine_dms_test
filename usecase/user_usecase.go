package usecase

import (
	"errors"
	"strconv"

	"enigmacamp.com/fine_dms/model"
	"enigmacamp.com/fine_dms/repo"
	"enigmacamp.com/fine_dms/utils"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUsecaseEmptyEmail     = errors.New("`email` cannot be empty")
	ErrUsecaseEmptyUsername  = errors.New("`username` cannot be empty")
	ErrUsecaseEmptyPassword  = errors.New("`password` cannot be empty")
	ErrUsecaseEmptyFname     = errors.New("`first_name` cannot be empty")
	ErrUsecaseExistsUsername = errors.New("`username` already exists")
	ErrUsecaseExistsEmail    = errors.New("`email` already exists")
	ErrUsecaseInvalidEmail   = errors.New("`email` invalid format")
	ErrUsecaseInvalidAuth    = errors.New("`username` or `password` wrong")
)

type user struct {
	userRepo repo.UserRepo
}

func NewUserUsecase(userRepo repo.UserRepo) UserUsecase {
	return &user{userRepo}
}

func (self *user) GetAll() ([]model.User, error) {
	res, err := self.userRepo.SelectAll()
	if err == repo.ErrRepoNoData {
		return nil, ErrUsecaseNoData
	}

	return res, nil
}

func (self *user) GetById(id int) (*model.User, error) {
	res, err := self.userRepo.SelectById(id)
	if err == repo.ErrRepoNoData {
		return nil, ErrUsecaseNoData
	}

	return res, nil
}

func (self *user) GetByUsername(uname string) (*model.User, error) {
	res, err := self.userRepo.SelectByUsername(uname)
	if err == repo.ErrRepoNoData {
		return nil, ErrUsecaseNoData
	}

	return res, nil
}

func (self *user) Add(user *model.User) error {
	if err := self.validateEmpty(user); err != nil {
		return err
	}

	if err := self.validateDuplicate(user); err != nil {
		return err
	}
	// TODO: hashing password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	return self.userRepo.Create(user)
}

func (self *user) Edit(user *model.User) error {
	if err := self.validateEmpty(user); err != nil {
		return err
	}

	if err := self.validateDuplicate(user); err != nil {
		return err
	}

	// TODO: hashing password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	return self.userRepo.Update(user)
}

func (self *user) Del(id int) error {
	return self.userRepo.Delete(id)
}

// private
func (self *user) validateEmpty(user *model.User) error {
	if len(user.Email) == 0 {
		return ErrUsecaseEmptyEmail
	}

	if len(user.Password) == 0 {
		return ErrUsecaseEmptyPassword
	}

	if len(user.FirstName) == 0 {
		return ErrUsecaseEmptyFname
	}

	if !utils.ValidateEmail(user.Email) {
		return ErrUsecaseInvalidEmail
	}

	return nil
}

func (self *user) validateDuplicate(user *model.User) error {
	_, err := self.userRepo.SelectByUsername(user.Username)
	if err != nil && err != repo.ErrRepoNoData {
		return ErrUsecaseInternal
	}
	if err == nil {
		return ErrUsecaseExistsUsername
	}

	_, err = self.userRepo.SelectByEmail(user.Email)
	if err != nil && err != repo.ErrRepoNoData {
		return ErrUsecaseInternal
	}
	if err == nil {
		return ErrUsecaseExistsEmail
	}

	return nil
}

func (self *user) AuthenticateUser(username string, password string) (int64, error) {
	user, err := self.GetByUsername(username)
	if err != nil {
		return 0, errors.New("Username atau password salah")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return 0, errors.New("Username atau password salah")
	}

	return int64(user.ID), nil
}

func (self *user) Login(username, password string) (*model.User, error) {
	user, err := self.userRepo.SelectByUsername(username)
	if err != nil {
		return nil, ErrUsecaseInvalidAuth
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, ErrUsecaseInvalidAuth
	}

	return user, nil
}

func (u *user) GetUserIdFromToken(tokenStr string, secret []byte) (int64, error) {
	userIdStr, err := utils.ValidateToken(tokenStr, secret)
	if err != nil {
		return 0, err
	}

	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return userId, nil
}
