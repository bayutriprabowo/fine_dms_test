package usecase

import (
	"errors"
	"fmt"
	"strconv"

	"enigmacamp.com/fine_dms/model"
	"enigmacamp.com/fine_dms/repo"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUsecaseEmptyEmail     = errors.New("`email` cannot be empty")
	ErrUsecaseEmptyUsername  = errors.New("`username` cannot be empty")
	ErrUsecaseEmptyPassword  = errors.New("`password` cannot be empty")
	ErrUsecaseEmptyFname     = errors.New("`first_name` cannot be empty")
	ErrUsecaseExistsUsername = errors.New("`username` already exists")
	ErrUsecaseExistsEmail    = errors.New("`email` already exists")
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

// func (self *user) Login(username, password string) (*model.User, error) {
// 	// cari user berdasarkan username
// 	user, err := self.userRepo.SelectByUsername(username)
// 	if err != nil {
// 		return nil, ErrUsecaseInvalidAuth
// 	}

// 	// verifikasi password
// 	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
// 	if err != nil {
// 		return nil, ErrUsecaseInvalidAuth
// 	}

// 	return user, nil
// }

func (self *user) AuthenticateUser(username string, password string) (int64, error) {
	// Mendapatkan data pengguna berdasarkan username
	user, err := self.GetByUsername(username)
	if err != nil {
		return 0, errors.New("Username atau password salah")
	}

	// Membandingkan password yang dimasukkan dengan password yang tersimpan
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return 0, errors.New("Username atau password salah")
	}

	// Autentikasi berhasil, mengembalikan ID pengguna
	return int64(user.ID), nil
}

func (self *user) GetUserIdFromToken(tokenString string) (int64, error) {
	// Parsing token dari string ke object
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Memverifikasi signature token dengan secret key
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}
		return []byte("secret"), nil // Ganti "secret" dengan secret key yang digunakan
	})
	if err != nil {
		return 0, err
	}

	// Mengekstrak claims dari token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, errors.New("Invalid token")
	}

	// Mengambil ID pengguna dari claims
	userId, err := strconv.ParseInt(fmt.Sprintf("%.0f", claims["userId"]), 10, 64)
	if err != nil {
		return 0, errors.New("Invalid token")
	}

	// Token valid, mengembalikan ID pengguna
	return userId, nil
}
