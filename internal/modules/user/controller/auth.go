package controller

import (
	"fmt"
	"net/http"
	"petstore/internal/logs"
	"petstore/internal/modules/user/service"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
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

// @Summary Регистрация нового пользователя
// @Description Создание нового пользователя с различными полями. Отсчёт будет идти от username (не должен повторяться)
// @Tags user
// @Accept json
// @Produce plain
// @Param user body User true "Заполни все поля для регистрации"
// @Success 200 {string} Info "Успешная регистрация или не ошибочное сообщение"
// @Failure 400 {string} Err "Неверной формат структуры"
// @Failure 500 {string} Err "Внутренняя ошибка сервера"
// @Router /user/ [post]
func (re *Respond) Register(w http.ResponseWriter, r *http.Request) {
	user := User{}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Неверный формат", http.StatusBadRequest)
		return
	}

	if err := IsValidStruct(user); err != nil {
		strErr := fmt.Sprintf("Невалидный запрос (неверный формат): %v", err.Error())
		http.Error(w, strErr, http.StatusBadRequest)
		return
	}

	serviceUser := service.ServiceUser{
		UserName:  user.UserName,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
		Phone:     user.Phone,
	}

	resp, err := re.Auth.RegisterUser(serviceUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))
}

// @Summary Регистрация группы пользователей
// @Description Создание новой группы пользователей с различными полями. Отсчёт будет идти от username (не должен повторяться)
// @Tags user
// @Accept json
// @Produce plain
// @Param user body UserArray true "Заполни все поля для регистрации"
// @Success 200 {string} Info "Успешная регистрация или не ошибочное сообщение"
// @Failure 400 {string} Err "Неверной формат структуры"
// @Failure 500 {string} Err "Внутренняя ошибка сервера"
// @Router /user/createWithArray [post]
func (re *Respond) RegisterArray(w http.ResponseWriter, r *http.Request) {
	user := UserArray{}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Неверный формат", http.StatusBadRequest)
		return
	}

	serviceArrayUser := service.ServiceUserArray{}
	for _, user := range user.UserArray {
		if err := IsValidStruct(user); err != nil {
			strErr := fmt.Sprintf("Невалидный запрос (неверный формат): %v", err.Error())
			http.Error(w, strErr, http.StatusBadRequest)
			return
		}
		serviceUser := service.ServiceUser{
			UserName:  user.UserName,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Password:  user.Password,
			Phone:     user.Phone,
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

// @Summary Вход в систему
// @Description Создание токена jwt и отправка в заголовке. Получите токен через postman и затем протестируйте logout, если необходимо
// @Tags user
// @Accept json
// @Produce plain
// @Param username query string true "Укажи username"
// @Param password query string true "Укажи пароль"
// @Success 200 {string} Info "Успешная регистрация или не ошибочное сообщение"
// @Failure 400 {string} Err "Неверной формат структуры"
// @Failure 500 {string} Err "Внутренняя ошибка сервера"
// @Router /user/login [get]
func (re *Respond) Login(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")
	if username == "" || password == "" {
		http.Error(w, "Неверный формат", http.StatusBadRequest)
		return
	}

	tokenString, resp, err := re.Auth.LoginUser(username, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Authorization", "Bearer "+tokenString)
	w.Write([]byte(resp))
}

// @Summary Выход из системы
// @Description Аннулирование jwt токена через black_list
// @Tags user
// @Accept json
// @Produce plain
// @Param username query string true "Укажи username"
// @Param password query string true "Укажи пароль"
// @Param Authorization header string true "Bearer токен доступа для имитации отправки браузером"
// @Success 200 {string} Info "Успешная регистрация или не ошибочное сообщение"
// @Failure 400 {string} Err "Неверной формат структуры"
// @Failure 500 {string} Err "Внутренняя ошибка сервера"
// @Router /user/logout [get]
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

// @Summary Обновить пользователя
// @Description Нельзя обновить на то имя пользователя, которое уже есть в базе
// @Tags user
// @Accept json
// @Produce plain
// @Param username path string true "Укажи username"
// @Param user body User true "Заполни все поля для обновления информации"
// @Success 200 {string} Info "Успешная регистрация или не ошибочное сообщение"
// @Failure 400 {string} Err "Неверной формат структуры"
// @Failure 500 {string} Err "Внутренняя ошибка сервера"
// @Router /user/{username} [put]
func (re *Respond) Update(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	user := User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Неверный формат", http.StatusBadRequest)
		return
	}

	if err := IsValidStruct(user); err != nil {
		strErr := fmt.Sprintf("Невалидный запрос (неверный формат): %v", err.Error())
		http.Error(w, strErr, http.StatusBadRequest)
		return
	}

	serviceUser := service.ServiceUser{
		UserName:  user.UserName,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
		Phone:     user.Phone,
	}

	resp, err := re.Auth.UpdateUser(username, serviceUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))
}

// @Summary Получить пользователя
// @Description Получить полную информацию о пользователе (для тестов)
// @Tags user
// @Accept json
// @Produce json
// @Param username path string true "Укажи username"
// @Success 200 {string} Info "Успешная регистрация или не ошибочное сообщение"
// @Failure 400 {string} Err "Неверной формат структуры"
// @Failure 500 {string} Err "Внутренняя ошибка сервера"
// @Router /user/{username} [get]
func (re *Respond) Get(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	user, err := re.Auth.GetUser(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user.UserName == "" {
		http.Error(w, "Такого пользователя не существует", http.StatusOK)
		return
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, "Неожиданная ошибка", http.StatusInternalServerError)
		return
	}
}

// @Summary Удаление пользователя
// @Description Мягкое удаление пользователя (проставление даты)
// @Tags user
// @Accept json
// @Produce plain
// @Param username path string true "Укажи username"
// @Success 200 {string} Info "Успешная регистрация или не ошибочное сообщение"
// @Failure 400 {string} Err "Неверной формат структуры"
// @Failure 500 {string} Err "Внутренняя ошибка сервера"
// @Router /user/{username} [delete]
func (re *Respond) Delete(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	resp, err := re.Auth.DeleteUser(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))
}

func IsValidStruct(u User) error {
	valid := validator.New()
	return (valid.Struct(u))
}
