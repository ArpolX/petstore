package repository

import (
	"petstore/internal/logs"
	"time"

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
	GetUsernamePassword(username, password string) (User, error)
	GetUsername(username string) (User, error)
	CreateTokenBlack(jti string) error
	TokenValid(jti string) (string, error)
}

func NewDb(db *gorm.DB, logger logs.Logger) Database {
	return &Db{
		db:  db,
		Log: logger,
	}
}

func (d *Db) Create(user RepositoryUser) error {
	u := User{
		Username:  user.UserName,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
		Phone:     user.Phone,
		CreateAt:  time.Now(),
		UpdateAt:  nil,
		DeleteAt:  nil,
	}

	err := d.db.Create(&u).Error
	if err != nil {
		d.Log.Error("Ошибка в Create запросе", zap.String("err", err.Error()))
		return err
	}

	return nil
}

func (d *Db) Update(login string, user RepositoryUser) error {
	u := User{}
	now := time.Now()

	err := d.db.Model(&u).
		Where("username = ?", login).
		Updates(map[string]interface{}{
			"username":  user.UserName,
			"firstname": user.FirstName,
			"lastname":  user.LastName,
			"email":     user.Email,
			"password":  user.Password,
			"phone":     user.Phone,
			"update_at": &now,
		}).Error
	if err != nil {
		d.Log.Error("Ошибка в Update запросе", zap.String("err", err.Error()))
		return err
	}

	return nil
}

func (d *Db) Delete(username string) error {
	u := User{}

	err := d.db.Model(&u).Where("username = ?", username).Update("delete_at", time.Now()).Error
	if err != nil {
		d.Log.Error("Ошибка в Delete запросе", zap.String("err", err.Error()))
		return err
	}

	return nil
}

func (d *Db) GetUsernamePassword(username, password string) (User, error) {
	u := User{}

	err := d.db.First(&u, "username = ? and password = ? and delete_at is null", username, password).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		d.Log.Error("Ошибка в Get запросе", zap.String("err", err.Error()))
		return User{}, err
	}

	return u, nil
}

func (d *Db) GetUsername(username string) (User, error) {
	u := User{}

	err := d.db.First(&u, "username = ? and delete_at is null", username).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		d.Log.Error("Ошибка в Get запросе", zap.String("err", err.Error()))
		return User{}, err
	}

	return u, nil
}

func (d *Db) CreateTokenBlack(jti string) error {
	blackList := BlackList{
		Jti: jti,
	}

	err := d.db.Create(&blackList).Error
	if err != nil {
		return err
	}

	return nil
}

func (d *Db) TokenValid(jti string) (string, error) {
	blackList := BlackList{}

	err := d.db.First(&blackList, "jti = ?", jti).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		d.Log.Error("Ошибка в First запросе", zap.String("err", err.Error()))
		return "", err
	}

	return blackList.Jti, nil
}
