package usecase

import (
	"enigmacamp.com/fine_dms/model"
	"enigmacamp.com/fine_dms/repository"
	"enigmacamp.com/fine_dms/utils"
)

type userUsecase struct {
	userRepo repository.UserRepository
}

func (usecase *userUsecase) Select() ([]model.User, error) {
	return usecase.userRepo.SelectUser()
}

func (usecase *userUsecase) Create(user *model.User) error {
	if err := utils.IsValidInput(user); err != nil {
		return err
	}
	if err := utils.IsDuplicateUsername(usecase.userRepo, user); err != nil {
		return err
	}
	if err := utils.IsDuplicateEmail(usecase.userRepo, user); err != nil {
		return err
	}
	return usecase.userRepo.CreateUser(user)
}

func (usecase *userUsecase) Update(user *model.User) error {
	if err := utils.IsValidInput(user); err != nil {
		return err
	}
	if err := utils.IsDuplicateUsername(usecase.userRepo, user); err != nil {
		return err
	}
	if err := utils.IsDuplicateEmail(usecase.userRepo, user); err != nil {
		return err
	}
	return usecase.userRepo.UpdateUser(user)
}

func (usecase *userUsecase) Delete(id int) error {
	return usecase.userRepo.DeleteUser(id)
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}
