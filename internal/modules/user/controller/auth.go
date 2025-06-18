package controller

import (
	"net/http"
	"petstore/internal/logs"
	"petstore/internal/modules/user/service"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Respond struct {
	Log  logs.Logger
	Auth service.Auther
}

type Responder interface {
	Register(w http.ResponseWriter, r *http.Request)
	RegisterArray(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

func NewRespond(logger logs.Logger, auth service.Auther) Responder {
	return &Respond{
		Log:  logger,
		Auth: auth,
	}
}

func (re *Respond) Register(w http.ResponseWriter, r *http.Request) {
	user := User{}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Неверный формат", http.StatusBadRequest)
		return
	}

	serviceUser := service.ServiceUser{
		UserName:   user.UserName,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Email:      user.Email,
		Password:   user.Password,
		Phone:      user.Phone,
		UserStatus: user.UserStatus,
	}

	resp, err := re.Auth.RegisterUser(serviceUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))
}

func (re *Respond) RegisterArray(w http.ResponseWriter, r *http.Request) {
	user := UserArray{}

	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		http.Error(w, "Неверный формат", http.StatusBadRequest)
		return
	}

	serviceArrayUser := service.ServiceUserArray{}
	for _, user := range user.UserArray {
		serviceUser := service.ServiceUser{
			UserName:   user.UserName,
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			Email:      user.Email,
			Password:   user.Password,
			Phone:      user.Phone,
			UserStatus: user.UserStatus,
		}

		serviceArrayUser.UserArray = append(serviceArrayUser.UserArray, serviceUser)
	}

	resp, err := re.Auth.RegisterArrayUser(serviceArrayUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))
}

func (re *Respond) Login(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")

	tokenString, resp, err := re.Auth.LoginUser(username, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", "Bearer "+tokenString)
	w.Write([]byte(resp))
}

func (re *Respond) Logout(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")
	AuthorizatUser := r.Header.Get("Authorization")

	resp, err := re.Auth.LogoutUser(username, password, AuthorizatUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))
}

func (re *Respond) Update(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Path

	user := User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Неверный формат", http.StatusBadRequest)
		return
	}

	serviceUser := service.ServiceUser{
		UserName:   user.UserName,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Email:      user.Email,
		Password:   user.Password,
		Phone:      user.Phone,
		UserStatus: user.UserStatus,
	}

	resp, err := re.Auth.UpdateUser(username, serviceUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))
}

func (re *Respond) Get(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Path

	user, err := re.Auth.GetUser(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user.UserName == "" {
		http.Error(w, "Такого пользователя не существует", http.StatusOK)
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, "Неожиданная ошибка", http.StatusInternalServerError)
		return
	}
}

func (re *Respond) Delete(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Path

	resp, err := re.Auth.DeleteUser(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))
}
