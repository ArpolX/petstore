package route

import (
	"net/http"
	"petstore/internal/modules/user/controller"

	"github.com/go-chi/chi"
)

func HandlerPetStore(uresp controller.Responder) http.Handler {
	r := chi.NewRouter()

	r.Route("/user", func(r chi.Router) {
		r.Post("/createWithArray", uresp.RegisterArray)
		r.Post("/", uresp.Register)
		r.Get("/login", uresp.Login)
		r.Get("/logout", uresp.Logout)
		r.Get("/{username}", uresp.Get)
		r.Put("/{username}", uresp.Update)
		r.Delete("/{username}", uresp.Delete)
	})

	return r
}
