package run

import (
	"petstore/internal/logs"
	"petstore/internal/modules/user/controller"
	"petstore/internal/modules/user/repository"
	"petstore/internal/modules/user/service"

	"gorm.io/gorm"
)

func NewModulesUser(db *gorm.DB, logger logs.Logger) controller.Responder {
	rep := repository.NewDb(db, logger)

	serv := service.NewAuth(logger, rep)

	ctrl := controller.NewRespond(logger, serv)

	return ctrl
}
