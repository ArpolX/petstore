package user

import (
	"net/http"
	"petstore/internal/modules/user/repository"

	"github.com/go-chi/jwtauth"
)

type AuthMiddleware struct {
	Db repository.Database
}

type AuthMiddlewarer interface {
	AuthMiddlewareBlackList() func(http.Handler) http.Handler
	AuthMiddlewareRoles(allRoles []string) func(http.Handler) http.Handler
}

func NewAuthMiddleware(db repository.Database) AuthMiddlewarer {
	return &AuthMiddleware{
		Db: db,
	}
}

func (a *AuthMiddleware) AuthMiddlewareBlackList() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())

			jti, ok := claims["jti"]
			if !ok {
				http.Error(w, "Не хватает полей в полезной нагрузке токена", http.StatusUnauthorized)
				return
			}

			jtiString := jti.(string)

			validJti, err := a.Db.TokenValid(jtiString)
			if err != nil {
				http.Error(w, "Неожиданная ошибка", http.StatusInternalServerError)
				return
			}

			if validJti != "" {
				http.Error(w, "No tokens", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (a *AuthMiddleware) AuthMiddlewareRoles(allRoles []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())

			roles, ok := claims["roles"]
			if !ok {
				http.Error(w, "Не хватает полей в полезной нагрузке токена", http.StatusUnauthorized)
				return
			}

			rolesSlice := roles.([]interface{})

			for _, role := range rolesSlice {
				for _, allRole := range allRoles {
					if role.(string) == allRole {
						next.ServeHTTP(w, r)
						return
					}
				}
			}

			http.Error(w, "Вы не прошли авторизацию", http.StatusForbidden)
		})
	}
}
