package run

import (
	"petstore/internal/logs"
	"petstore/internal/modules/user"
	userCtrl "petstore/internal/modules/user/controller"
	userRep "petstore/internal/modules/user/repository"
	userServ "petstore/internal/modules/user/service"

	petCtrl "petstore/internal/modules/pet/controller"
	petRep "petstore/internal/modules/pet/repository"
	petServ "petstore/internal/modules/pet/service"

	orderCtrl "petstore/internal/modules/order/controller"
	orderRep "petstore/internal/modules/order/repository"
	orderServ "petstore/internal/modules/order/service"

	"gorm.io/gorm"
)

// user слои
func NewModulesUser(db *gorm.DB, logger logs.Logger) (userCtrl.Responder, user.AuthMiddlewarer) {
	rep := userRep.NewDb(db, logger)

	// настройка middleware
	middleAuth := user.NewAuthMiddleware(rep)

	serv := userServ.NewAuth(logger, rep)

	ctrl := userCtrl.NewRespond(logger, serv)

	return ctrl, middleAuth
}

// pet слои
func NewModulesPet(db *gorm.DB, logger logs.Logger) petCtrl.AnimalStorer {
	rep := petRep.NewPetRepository(logger, db)

	serv := petServ.NewPetService(logger, rep)

	ctrl := petCtrl.NewAnimalStore(logger, serv)

	return ctrl
}

// order слои
func NewModulesOrder(db *gorm.DB, logger logs.Logger) orderCtrl.OrderResponder {
	rep := orderRep.NewOrderRepository(logger, db)

	serv := orderServ.NewOrderService(logger, rep)

	ctrl := orderCtrl.NewOrderRespond(logger, serv)

	return ctrl
}
