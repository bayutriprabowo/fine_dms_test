package manager

import (
	"enigmacamp.com/fine_dms/repo"
	"enigmacamp.com/fine_dms/repo/psql"
)

type RepoManager interface {
	UserRepo() repo.UserRepo
	// Add other repo below
}

type repoManager struct {
	infraMgr InfraManager
}

func NewRepoManager(infr InfraManager) RepoManager {
	return &repoManager{
		infraMgr: infr,
	}
}

func (self *repoManager) UserRepo() repo.UserRepo {
	return psql.NewPsqlUserRepo(self.infraMgr.DbConn())
}
