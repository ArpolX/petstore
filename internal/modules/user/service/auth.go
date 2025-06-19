package service

import (
	"errors"
	"fmt"
	"petstore/internal/logs"
	"petstore/internal/modules/user/repository"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Auth struct {
	Log logs.Logger
	Db  repository.Database
}

type Auther interface {
	RegisterUser(user ServiceUser) (string, error)
	RegisterArrayUser(userArray ServiceUserArray) (string, error)
	LoginUser(login, password string) (string, string, error)
	LogoutUser(login, password, authUser string) (string, error)
	UpdateUser(login string, user ServiceUser) (string, error)
	GetUser(login string) (ServiceUser, error)
	DeleteUser(login string) (string, error)
}

func NewAuth(logger logs.Logger, db repository.Database) Auther {
	return &Auth{
		Log: logger,
		Db:  db,
	}
}

func (a *Auth) RegisterUser(user ServiceUser) (string, error) {
	repUser, err := a.Db.GetUsernamePassword(user.UserName, user.Password)
	if err != nil {
		return "", errors.New("Неожиданная ошибка")
	}

	if repUser.Username != "" {
		return "Такой пользователь уже существует", nil
	}

	repositoryUser := repository.RepositoryUser{
		UserName:  user.UserName,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
		Phone:     user.Phone,
	}
	err = a.Db.Create(repositoryUser)
	if err != nil {
		return "", errors.New("Неожиданная ошибка")
	}

	return fmt.Sprintf("Регистрация пользователя %v прошла успешно", user.UserName), nil
}

func (a *Auth) RegisterArrayUser(userArr ServiceUserArray) (string, error) {
	for _, u := range userArr.UserArray {
		repUser, err := a.Db.GetUsernamePassword(u.UserName, u.Password)
		if err != nil {
			return "", errors.New("Неожиданная ошибка")
		}

		if repUser.Username != "" {
			return fmt.Sprintf("Пользователь %v уже существует, процедура прервана", u.UserName), nil
		}
	}

	repositoryArrayUser := repository.RepositoryUserArray{}
	for _, user := range userArr.UserArray {
		repositoryUser := repository.RepositoryUser{
			UserName:  user.UserName,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Password:  user.Password,
			Phone:     user.Phone,
		}

		repositoryArrayUser.UserArray = append(repositoryArrayUser.UserArray, repositoryUser)
	}

	for _, user := range repositoryArrayUser.UserArray {
		err := a.Db.Create(user)
		if err != nil {
			return "", errors.New("Неожиданная ошибка")
		}
	}

	return "Регистрация нескольких пользователей прошла успешно", nil
}

func (a *Auth) LoginUser(login, password string) (string, string, error) {
	repUser, err := a.Db.GetUsernamePassword(login, password)
	if err != nil {
		return "", "", errors.New("Неожиданная ошибка")
	}

	if repUser.Username == "" {
		return "", "Неверный логин или пароль", nil
	}

	jti := uuid.New().String()
	claims := jwt.MapClaims{
		"sub":   login,
		"exp":   time.Now().Add(1 * time.Hour).Unix(),
		"roles": []string{"admin"},
		"jti":   jti,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("ho-ho"))
	if err != nil {
		a.Log.Warn("Токен не создан", zap.String("err", err.Error()))
		return "", "", errors.New("Неожиданная ошибка")
	}

	return tokenString, "Авторизация прошла успешно, токен присвоен в заголовке", nil
}

func (a *Auth) LogoutUser(login, password, authUser string) (string, error) {
	repUser, err := a.Db.GetUsernamePassword(login, password)
	if err != nil {
		return "", errors.New("Неожиданная ошибка")
	}

	if repUser.Username == "" {
		return "Неверный логин или пароль", nil
	}

	authUserSplit := strings.Split(authUser, " ")
	if len(authUserSplit) != 2 {
		a.Log.Warn("У авторизованного пользователя нет заголовка Authorization")
		return "", errors.New("Неожиданная ошибка")
	}

	token, _, err := new(jwt.Parser).ParseUnverified(authUserSplit[1], jwt.MapClaims{})
	if err != nil {
		a.Log.Warn("Ошибка парсинга jwt токена")
		return "", errors.New("Неожиданная ошибка")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		a.Log.Warn("в токене нет claims")
		return "", errors.New("Неожиданная ошибка")
	}

	jti, _ := claims["jti"]

	jtiString := jti.(string)

	err = a.Db.CreateTokenBlack(jtiString)
	if err != nil {
		return "", errors.New("Неожиданная ошибка")
	}

	return "Ваш токен аннулирован", nil
}

func (a *Auth) UpdateUser(login string, user ServiceUser) (string, error) {
	repLogin, err := a.Db.GetUsername(login)
	if err != nil {
		return "", errors.New("Неожиданная ошибка")
	}

	if repLogin.Username == "" {
		return "Такого логина не существует", nil
	}

	repUser, err := a.Db.GetUsername(user.UserName)
	if err != nil {
		return "", errors.New("Неожиданная ошибка")
	}

	if repUser.Username != "" && repUser.Username != login {
		return "Логин, на который вы хотите обновиться, уже существует", nil
	}

	repositoryUser := repository.RepositoryUser{
		UserName:  user.UserName,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
		Phone:     user.Phone,
	}

	err = a.Db.Update(login, repositoryUser)
	if err != nil {
		return "", errors.New("Неожиданная ошибка")
	}

	return fmt.Sprintf("Обновление информации пользователя %v прошло успешно", login), nil
}

func (a *Auth) GetUser(login string) (ServiceUser, error) {
	user, err := a.Db.GetUsername(login)
	if err != nil {
		return ServiceUser{}, errors.New("Неожиданная ошибка")
	}

	if user.Username == "" {
		return ServiceUser{}, nil
	}

	serviceUser := ServiceUser{
		UserName:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
		Phone:     user.Phone,
	}

	return serviceUser, nil
}

func (a *Auth) DeleteUser(login string) (string, error) {
	user, err := a.Db.GetUsername(login)
	if err != nil {
		return "", errors.New("Неожиданная ошибка")
	}

	if user.Username == "" {
		return "Такого пользователя не существует", nil
	}

	err = a.Db.Delete(login)
	if err != nil {
		return "", errors.New("Неожиданная ошибка")
	}

	return fmt.Sprintf("Удаление пользователя %v прошло успешно", login), nil
}
