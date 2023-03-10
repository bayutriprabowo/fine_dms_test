package manager

import (
	"database/sql"
	"fmt"

	"enigmacamp.com/fine_dms/config"
	_ "github.com/lib/pq"
)

type InfraManager interface {
	DbInit() error
	DbConn() *sql.DB
}

type infra struct {
	db     *sql.DB
	apiCfg config.ApiConfig
	dbCfg  config.DbConfig
}

func NewInfraManager(cfg config.AppConfig) InfraManager {
	return &infra{
		apiCfg: cfg.ApiConfig,
		dbCfg:  cfg.DbConfig,
	}
}

func (self *infra) DbInit() error {
	c := &self.dbCfg

	dbParams := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Name, c.SslMode,
	)

	db, err := sql.Open("postgres", dbParams)
	if err != nil {
		return err
	}

	self.db = db
	return nil
}

func (self *infra) DbConn() *sql.DB {
	return self.db
}
