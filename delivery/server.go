package delivery

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"enigmacamp.com/fine_dms/config"
	"enigmacamp.com/fine_dms/controller"
	"enigmacamp.com/fine_dms/manager"
	"github.com/gin-gonic/gin"
)

type AppServer struct {
	infra     manager.InfraManager
	ucMgr     manager.UsecaseManager
	engine    *gin.Engine
	hostPort  string
	secretKey []byte
}

func NewAppServer() AppServer {
	srv := gin.Default()
	cfg := config.NewAppConfig()

	infrMgr := manager.NewInfraManager(cfg)
	rpMgr := manager.NewRepoManager(infrMgr)
	ucMgr := manager.NewUsecaseManager(rpMgr)

	// generate a secure and random secret key
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		panic(err)
	}
	secretKey := base64.StdEncoding.EncodeToString(key)

	return AppServer{
		infra:  infrMgr,
		ucMgr:  ucMgr,
		engine: srv,
		hostPort: fmt.Sprintf("%s:%s", cfg.ApiConfig.Host,
			cfg.ApiConfig.Port),
		secretKey: []byte(secretKey),
	}
}

func (self *AppServer) Run() error {
	if err := self.infra.DbInit(); err != nil {
		return err
	}

	defer self.infra.DbConn().Close()

	self.v1()

	if err := self.engine.Run(self.hostPort); err != nil {
		return err
	}

	return nil
}

// private
func (self *AppServer) v1() {
	baseRg := self.engine.Group("/v1")
	controller.NewUserController(baseRg, self.ucMgr.UserUsecase(), self.secretKey)
	controller.NewTagsController(baseRg, self.ucMgr.TagsUsecase())
}
