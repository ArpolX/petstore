package route

import (
	"net/http"
	petCtrl "petstore/internal/modules/pet/controller"
	"petstore/internal/modules/user"
	userCtrl "petstore/internal/modules/user/controller"

	orderCtrl "petstore/internal/modules/order/controller"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
)

func HandlerPetStore(middleAuth user.AuthMiddlewarer, userСtrl userCtrl.Responder, petCtrl petCtrl.AnimalStorer, orderCtrl orderCtrl.OrderResponder) http.Handler {
	r := chi.NewRouter()

	token := jwtauth.New("HS256", []byte("ho-ho"), nil)

	r.Route("/user", func(r chi.Router) {
		r.Post("/createWithArray", userСtrl.RegisterArray)
		r.Post("/", userСtrl.Register)
		r.Get("/login", userСtrl.Login)
		r.Get("/logout", userСtrl.Logout)
		r.Get("/{username}", userСtrl.Get)
		r.Put("/{username}", userСtrl.Update)
		r.Delete("/{username}", userСtrl.Delete)
	})

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(token))
		r.Use(jwtauth.Authenticator)
		r.Use(middleAuth.AuthMiddlewareBlackList())
		r.Use(middleAuth.AuthMiddlewareRoles([]string{"admin"}))

		r.Route("/pet", func(r chi.Router) {
			r.Post("/", petCtrl.RegisterPet)
			r.Put("/", petCtrl.UpdatePet)
			r.Get("/findByStatus", petCtrl.GetPetByStatus)
			r.Get("/{petId}", petCtrl.GetPet)
			r.Post("/{petId}", petCtrl.UpdateNameStatusPet)
			r.Delete("/{petId}", petCtrl.DeletePet)
		})
	})

	r.Route("/store", func(r chi.Router) {
		r.Post("/order", orderCtrl.PlaceOrder)
		r.Get("/order/{orderId}", orderCtrl.GetOrder)
		r.Delete("/order/{orderId}", orderCtrl.DeleteOrder)
	})

	return r
}
