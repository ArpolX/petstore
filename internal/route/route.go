package route

import (
	"net/http"
	petCtrl "petstore/internal/modules/pet/controller"
	userCtrl "petstore/internal/modules/user/controller"

	"github.com/go-chi/chi"
)

func HandlerPetStore(userResp userCtrl.Responder, petResp petCtrl.AnimalStorer) http.Handler {
	r := chi.NewRouter()

	//token := jwtauth.New("HS256", []byte("ho-ho"), nil)

	r.Route("/user", func(r chi.Router) {
		r.Post("/createWithArray", userResp.RegisterArray)
		r.Post("/", userResp.Register)
		r.Get("/login", userResp.Login)
		r.Get("/logout", userResp.Logout)
		r.Get("/{username}", userResp.Get)
		r.Put("/{username}", userResp.Update)
		r.Delete("/{username}", userResp.Delete)
	})

	r.Group(func(r chi.Router) {
		//r.Use(jwtauth.Verifier(token))
		//r.Use(jwtauth.Authenticator)

		r.Route("/pet", func(r chi.Router) {
			r.Post("/", petResp.RegisterPet)
			r.Put("/", petResp.UpdatePet)
			r.Get("/findByStatus", petResp.GetPetByStatus)
			r.Get("/{petId}", petResp.GetPet)
			r.Post("/{petId}", petResp.UpdateNameStatusPet)
			r.Delete("/{petId}", petResp.DeletePet)
		})
	})

	return r
}
