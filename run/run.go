package run

import (
	"petstore/internal/logs"
	userCtrl "petstore/internal/modules/user/controller"
	userResp "petstore/internal/modules/user/repository"
	userServ "petstore/internal/modules/user/service"

	petCtrl "petstore/internal/modules/pet/controller"
	petResp "petstore/internal/modules/pet/repository"
	petServ "petstore/internal/modules/pet/service"

	"gorm.io/gorm"
)

func NewModulesUser(db *gorm.DB, logger logs.Logger) userCtrl.Responder {
	rep := userResp.NewDb(db, logger)

	serv := userServ.NewAuth(logger, rep)

	ctrl := userCtrl.NewRespond(logger, serv)

	return ctrl
}

func NewModulesPet(db *gorm.DB, logger logs.Logger) petCtrl.AnimalStorer {
	rep := petResp.NewPetRepository(logger, db)

	serv := petServ.NewPetService(logger, rep)

	ctrl := petCtrl.NewAnimalStore(logger, serv)

	return ctrl
}
