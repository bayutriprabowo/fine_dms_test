package manager

import (
	"enigmacamp.com/fine_dms/usecase"
)

type UsecaseManager interface {
	UserUsecase() usecase.UserUsecase
	// Add other usecase below
}

type usecaseManager struct {
	repoMgr RepoManager
}

func NewUsecaseManager(repoMgr RepoManager) UsecaseManager {
	return &usecaseManager{
		repoMgr: repoMgr,
	}
}

func (self *usecaseManager) UserUsecase() usecase.UserUsecase {
	return usecase.NewUserUsecase(self.repoMgr.UserRepo())
}
