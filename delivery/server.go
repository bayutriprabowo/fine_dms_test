package delivery

import (
	"fmt"

	"enigmacamp.com/fine_dms/config"
	"enigmacamp.com/fine_dms/controller"
	"enigmacamp.com/fine_dms/manager"
	"github.com/gin-gonic/gin"
)

type AppServer struct {
	infra    manager.InfraManager
	ucMgr    manager.UsecaseManager
	hostPort string
}

func NewAppServer() AppServer {
	cfg := config.NewAppConfig()

	infrMgr := manager.NewInfraManager(cfg)
	rpMgr := manager.NewRepoManager(infrMgr)
	ucMgr := manager.NewUsecaseManager(rpMgr)

	return AppServer{
		infra: infrMgr,
		ucMgr: ucMgr,
		hostPort: fmt.Sprintf("%s:%s", cfg.ApiConfig.Host,
			cfg.ApiConfig.Port),
	}
}

func (self *AppServer) Run() error {
	if err := self.infra.Init(); err != nil {
		return err
	}

	defer self.infra.Deinit()

	engine := gin.Default()

	self.v1(engine)

	if err := engine.Run(self.hostPort); err != nil {
		return err
	}

	return nil
}

// private
func (self *AppServer) v1(engine *gin.Engine) {
	baseRg := engine.Group("/v1")
	controller.NewUserController(baseRg, self.ucMgr.UserUsecase())
	controller.NewTagsController(baseRg, self.ucMgr.TagsUsecase())
}
