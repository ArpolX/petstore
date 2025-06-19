package repository

import "time"

type RepositoryUser struct {
	UserName  string
	FirstName string
	LastName  string
	Email     string
	Password  string
	Phone     string
}

type RepositoryUserArray struct {
	UserArray []RepositoryUser
}

type User struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Username  string `gorm:"type:text;unique;not null"`
	FirstName string `gorm:"type:text;column:firstname;not null"`
	LastName  string `gorm:"type:text;column:lastname;not null"`
	Email     string `gorm:"type:text;not null"`
	Password  string `gorm:"type:text;not null"`
	Phone     string `gorm:"type:text;not null"`
	CreateAt  time.Time
	UpdateAt  *time.Time
	DeleteAt  *time.Time
}

type BlackList struct {
	ID  uint `gorm:"primaryKey;autoIncrement"`
	Jti string
}
