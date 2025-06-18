package repository

import (
	"petstore/internal/logs"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Db struct {
	db  *gorm.DB
	Log logs.Logger
}

type Database interface {
	Create(user RepositoryUser) error
	Update(login string, user RepositoryUser) error
	Delete(username string) error
	GetUsernamePassword(username, password string) (RepositoryUser, error)
	GetUsername(username string) (RepositoryUser, error)
	TokenValid(jti string) error
}

func NewDb(db *gorm.DB, logger logs.Logger) Database {
	return &Db{
		db:  db,
		Log: logger,
	}
}

func (d *Db) Create(user RepositoryUser) error {
	u := User{
		Username:   user.UserName,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Email:      user.Email,
		Password:   user.Password,
		Phone:      user.Phone,
		UserStatus: user.UserStatus,
	}

	err := d.db.Create(&u).Error
	if err != nil {
		d.Log.Error("Ошибка в Create запросе", zap.String("err", err.Error()))
		return err
	}

	return nil
}

func (d *Db) Update(login string, user RepositoryUser) error {
	err := d.db.Model(&User{}).
		Where("username = ?", login).
		Updates(User{
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			Email:      user.Email,
			Password:   user.Password,
			Phone:      user.Phone,
			UserStatus: user.UserStatus,
		}).Error
	if err != nil {
		d.Log.Error("Ошибка в Update запросе", zap.String("err", err.Error()))
		return err
	}

	return nil
}

func (d *Db) Delete(username string) error {
	err := d.db.Delete(&User{}, "username = ?", username).Error
	if err != nil {
		d.Log.Error("Ошибка в Delete запросе", zap.String("err", err.Error()))
		return err
	}

	return nil
}

func (d *Db) GetUsernamePassword(username, password string) (RepositoryUser, error) {
	u := RepositoryUser{}

	err := d.db.First(&u, "username = ?", username, "password = ?", password).Error
	if err != nil {
		d.Log.Error("Ошибка в Get запросе", zap.String("err", err.Error()))
		return RepositoryUser{}, err
	}

	return u, nil
}

func (d *Db) GetUsername(username string) (RepositoryUser, error) {
	u := RepositoryUser{}

	err := d.db.First(&u, "username = ?", username).Error
	if err != nil {
		d.Log.Error("Ошибка в Get запросе", zap.String("err", err.Error()))
		return RepositoryUser{}, err
	}

	return u, nil
}

func (d *Db) TokenValid(ti string) error {
	return nil
}
